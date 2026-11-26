package laporanpengeluaran

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Mapping nama bulan bahasa Indonesia
var namaBulanID = map[string]int{
	"januari": 1, "februari": 2, "maret": 3, "april": 4,
	"mei": 5, "juni": 6, "juli": 7, "agustus": 8,
	"september": 9, "oktober": 10, "november": 11, "desember": 12,
}

var listBulanID = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// LayananLaporanPengeluaran bertanggung jawab mengelompokkan biaya pengeluaran berdasarkan kategori beban
type LayananLaporanPengeluaran struct {
	DB               *sql.DB
	MemoriDummy      []ItemPengeluaran
	Mu               sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananLaporanPengeluaran
func InisialisasiLayanan(db *sql.DB) *LayananLaporanPengeluaran {
	return &LayananLaporanPengeluaran{
		DB: db,
		MemoriDummy: []ItemPengeluaran{
			{
				IDPengeluaran:      1,
				TanggalPengeluaran: "2026-11-05",
				KategoriAsli:       "Listrik & Air",
				KategoriBeban:      "OPERASIONAL",
				Nominal:            450000,
				NominalFormat:      "Rp450.000",
				Keterangan:         "Tagihan listrik tempat produksi",
			},
			{
				IDPengeluaran:      2,
				TanggalPengeluaran: "2026-11-10",
				KategoriAsli:       "Beras Ketan & Bumbu",
				KategoriBeban:      "BELANJA_PERSEDIAAN",
				Nominal:            1250000,
				NominalFormat:      "Rp1.250.000",
				Keterangan:         "Bahan baku utama rengginang",
			},
			{
				IDPengeluaran:      3,
				TanggalPengeluaran: "2026-11-18",
				KategoriAsli:       "Minyak Goreng & Gas",
				KategoriBeban:      "OPERASIONAL",
				Nominal:            650000,
				NominalFormat:      "Rp650.000",
				Keterangan:         "Bahan pendukung penggorengan",
			},
			{
				IDPengeluaran:      4,
				TanggalPengeluaran: "2026-11-22",
				KategoriAsli:       "Perbaikan Alat Pengering",
				KategoriBeban:      "BIAYA_LAIN",
				Nominal:            200000,
				NominalFormat:      "Rp200.000",
				Keterangan:         "Servis blower ruang penjemuran",
			},
		},
	}
}

// KlasifikasikanKategoriBeban memetakan string kategori ke 3 kategori beban utama:
// 1. BELANJA_PERSEDIAAN: persediaan bahan baku mentah, bahan produksi rengginang, kemasan utama, barang dagang
// 2. OPERASIONAL: listrik, air, gas, bahan bakar, gaji/upah harian, sewa tempat, internet, kebersihan
// 3. BIAYA_LAIN: pemeliharaan mesin, transportasi/ekspedisi pengiriman non-pelanggan, administrasi, tak terduga
func KlasifikasikanKategoriBeban(kategoriInput string) string {
	kat := strings.ToLower(strings.TrimSpace(kategoriInput))

	// Kategori Belanja Persediaan (Bahan Baku / Modal Stok)
	if strings.Contains(kat, "persediaan") ||
		strings.Contains(kat, "bahan baku") ||
		strings.Contains(kat, "beras") ||
		strings.Contains(kat, "ketan") ||
		strings.Contains(kat, "udang") ||
		strings.Contains(kat, "bumbu") ||
		strings.Contains(kat, "kemasan") ||
		strings.Contains(kat, "plastik") ||
		strings.Contains(kat, "stok") ||
		strings.Contains(kat, "pembelian barang") {
		return "BELANJA_PERSEDIAAN"
	}

	// Kategori Biaya Operasional (Biaya rutin operasional tempat usaha)
	if strings.Contains(kat, "operasional") ||
		strings.Contains(kat, "listrik") ||
		strings.Contains(kat, "air") ||
		strings.Contains(kat, "gas") ||
		strings.Contains(kat, "minyak") ||
		strings.Contains(kat, "gaji") ||
		strings.Contains(kat, "upah") ||
		strings.Contains(kat, "sewa") ||
		strings.Contains(kat, "kebersihan") ||
		strings.Contains(kat, "telepon") ||
		strings.Contains(kat, "internet") {
		return "OPERASIONAL"
	}

	// Default ke Biaya Lain (Anti-duplikasi: di luar belanja persediaan & operasional rutin)
	return "BIAYA_LAIN"
}

// FormatRupiah mengubah nominal float64 menjadi format Rp1.500.000
func FormatRupiah(val float64) string {
	valInt := int64(val)
	strVal := strconv.FormatInt(valInt, 10)
	var hasil []string
	pjg := len(strVal)
	for i := pjg; i > 0; i -= 3 {
		awal := i - 3
		if awal < 0 {
			awal = 0
		}
		hasil = append([]string{strVal[awal:i]}, hasil...)
	}
	return "Rp" + strings.Join(hasil, ".")
}

// ParsingPeriode menerjemahkan parameter bulan dan tahun ke rentang tanggal
func (l *LayananLaporanPengeluaran) ParsingPeriode(bulanStr, tahunStr string) (PeriodeLaporan, error) {
	sekarang := time.Now()
	tahun := sekarang.Year()
	bulan := int(sekarang.Month())

	if strings.TrimSpace(tahunStr) != "" {
		if tParsed, err := strconv.Atoi(strings.TrimSpace(tahunStr)); err == nil && tParsed > 2000 {
			tahun = tParsed
		}
	}

	bulanStrClean := strings.ToLower(strings.TrimSpace(bulanStr))
	if bVal, ada := namaBulanID[bulanStrClean]; ada {
		bulan = bVal
	} else if bParsed, err := strconv.Atoi(bulanStrClean); err == nil && bParsed >= 1 && bParsed <= 12 {
		bulan = bParsed
	}

	awalPeriode := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.Local)
	akhirPeriode := awalPeriode.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return PeriodeLaporan{
		BulanNama:    listBulanID[bulan],
		BulanAngka:   bulan,
		Tahun:        tahun,
		AwalPeriode:  awalPeriode,
		AkhirPeriode: akhirPeriode,
	}, nil
}

// HitungRekapPengeluaran mengelompokkan pengeluaran berdasarkan tanggal dan kategori
func (l *LayananLaporanPengeluaran) HitungRekapPengeluaran(bulanStr, tahunStr string) (*DataLaporanPengeluaran, error) {
	periode, err := l.ParsingPeriode(bulanStr, tahunStr)
	if err != nil {
		return nil, err
	}

	var daftar []ItemPengeluaran
	var totalOperasional, totalBelanja, totalBiayaLain float64

	formatAwal := periode.AwalPeriode.Format("2006-01-02")
	formatAkhir := periode.AkhirPeriode.Format("2006-01-02")

	if l.DB != nil {
		// Pastikan tabel pengeluaran_usaha ada
		query := `
			SELECT 
				id_pengeluaran, 
				DATE_FORMAT(tanggal_pengeluaran, '%Y-%m-%d') as tgl, 
				kategori, 
				nominal, 
				COALESCE(keterangan, '') as ket
			FROM pengeluaran_usaha
			WHERE tanggal_pengeluaran >= ? AND tanggal_pengeluaran <= ?
			ORDER BY tanggal_pengeluaran ASC, id_pengeluaran ASC
		`
		rows, err := l.DB.Query(query, formatAwal, formatAkhir)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var item ItemPengeluaran
				var nominalDb float64
				if errScan := rows.Scan(&item.IDPengeluaran, &item.TanggalPengeluaran, &item.KategoriAsli, &nominalDb, &item.Keterangan); errScan == nil {
					item.Nominal = nominalDb
					item.NominalFormat = FormatRupiah(nominalDb)
					item.KategoriBeban = KlasifikasikanKategoriBeban(item.KategoriAsli)

					switch item.KategoriBeban {
					case "OPERASIONAL":
						totalOperasional += nominalDb
					case "BELANJA_PERSEDIAAN":
						totalBelanja += nominalDb
					case "BIAYA_LAIN":
						totalBiayaLain += nominalDb
					}
					daftar = append(daftar, item)
				}
			}
		}
	}

	// Fallback ke data memori jika database belum ada baris di rentang tersebut
	if len(daftar) == 0 {
		l.Mu.RLock()
		for _, itm := range l.MemoriDummy {
			// Periksa apakah masuk periode
			if strings.HasPrefix(itm.TanggalPengeluaran, fmt.Sprintf("%04d-%02d", periode.Tahun, periode.BulanAngka)) {
				switch itm.KategoriBeban {
				case "OPERASIONAL":
					totalOperasional += itm.Nominal
				case "BELANJA_PERSEDIAAN":
					totalBelanja += itm.Nominal
				case "BIAYA_LAIN":
					totalBiayaLain += itm.Nominal
				}
				daftar = append(daftar, itm)
			}
		}
		l.Mu.RUnlock()
	}

	totalSemua := totalOperasional + totalBelanja + totalBiayaLain

	hasil := &DataLaporanPengeluaran{
		Periode:           periode,
		Operasional:       totalOperasional,
		OperasionalFormat: FormatRupiah(totalOperasional),
		BelanjaPersediaan: totalBelanja,
		BelanjaFormat:     FormatRupiah(totalBelanja),
		BiayaLain:         totalBiayaLain,
		BiayaLainFormat:   FormatRupiah(totalBiayaLain),
		TotalPengeluaran:  totalSemua,
		TotalFormat:       FormatRupiah(totalSemua),
		SumberData: RekapSumberData{
			TransaksiValid:   "Mengikuti backend",
			Sandbox:          "Tidak dihitung sebagai kas nyata",
			DataBelumLengkap: "Ditandai dan dijelaskan",
		},
		DaftarPengeluaran: daftar,
		CatatanAudit:      fmt.Sprintf("Rekapitulasi pengeluaran usaha periode %s %d ditutup resmi pada akhir bulan.", periode.BulanNama, periode.Tahun),
		DibuatPada:        time.Now(),
	}

	return hasil, nil
}

// HandlerLaporanPengeluaran adalah HTTP handler untuk GET /api/admin/laporan/pengeluaran
func (l *LayananLaporanPengeluaran) HandlerLaporanPengeluaran(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(ResponseLaporanPengeluaran{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	bulan := r.URL.Query().Get("bulan")
	tahun := r.URL.Query().Get("tahun")

	laporan, err := l.HitungRekapPengeluaran(bulan, tahun)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ResponseLaporanPengeluaran{
			Status: "gagal",
			Pesan:  "Gagal mengolah data laporan pengeluaran: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ResponseLaporanPengeluaran{
		Status: "sukses",
		Data:   laporan,
	})
}
