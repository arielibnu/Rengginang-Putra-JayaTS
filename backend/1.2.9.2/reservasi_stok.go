package reservasisstok

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ordergenerator "rengginang-putra-jaya/auth/1.2.9.1"
)

var (
	memoriStok = map[string]int{
		"Rengginang Udang":  12,
		"Rengginang Bawang": 38,
		"Rengginang Terasi": 38,
		"Rengginang Cumi":   0,
		"Udang":             12,
		"Bawang":            38,
		"Terasi":            38,
		"Cumi":              0,
		"prod-1":            12,
		"prod-2":            38,
	}

	memoriReservasi = make([]CatatanReservasi, 0)
	muStok          sync.RWMutex
)

// HandlerReservasiStok pengendali logika reservasi stok dan pembuatan order
type HandlerReservasiStok struct {
	DB        *sql.DB
	Generator *ordergenerator.GeneratorOrder
}

// InisialisasiTabel membuat tabel reservasi_stok di MySQL jika belum ada
func (h *HandlerReservasiStok) InisialisasiTabel() {
	if h.DB == nil {
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS reservasi_stok (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		id_varian BIGINT NOT NULL DEFAULT 1,
		nama_varian VARCHAR(100) NOT NULL,
		jumlah_direservasi INT NOT NULL,
		status_reservasi VARCHAR(30) NOT NULL DEFAULT 'menahan_stok',
		kadaluarsa_pada DATETIME NOT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_nomor_reservasi (nomor_pesanan),
		INDEX idx_status_reservasi (status_reservasi),
		INDEX idx_kadaluarsa (kadaluarsa_pada)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, _ = h.DB.Exec(query)
}

// CekDanReservasiStok memvalidasi sisa stok riil dan menahan kuota selama 24 jam
func (h *HandlerReservasiStok) CekDanReservasiStok(nomorPesanan string, items []ordergenerator.ItemSnapshotPesanan) (*HasilReservasiStok, error) {
	muStok.Lock()
	defer muStok.Unlock()

	now := time.Now()
	batasWaktu := now.Add(24 * time.Hour) // Holding time 24 jam

	// 1. VALIDASI STOK SEBELUM PENAHANAN (PROTEKSI OVERBOOKING)
	for _, it := range items {
		namaKunci := it.NamaProduk
		if strings.TrimSpace(it.NamaVarian) != "" {
			namaKunci = it.NamaVarian
		}

		stokAda, ok := memoriStok[namaKunci]
		if !ok {
			stokAda = memoriStok[it.IDProduk]
			if stokAda == 0 && !ok {
				stokAda = 50 // default fallback
			}
		}

		// Hitung stok yang sedang direservasi aktif
		direservasiAktif := 0
		for _, r := range memoriReservasi {
			if r.NamaVarian == namaKunci && r.StatusReservasi == "menahan_stok" && r.KadaluarsaPada.After(now) {
				direservasiAktif += r.JumlahDireservasi
			}
		}

		sisaStok := stokAda - direservasiAktif
		if it.Jumlah > sisaStok {
			return &HasilReservasiStok{
				Sukses:    false,
				KodeHasil: "STOK_TIDAK_CUKUP",
				Pesan: fmt.Sprintf(
					"Stok produk '%s' tidak mencukupi untuk dipesan. Sisa tersedia: %d bungkus, permintaan: %d bungkus.",
					it.NamaProduk, sisaStok, it.Jumlah,
				),
			}, fmt.Errorf("stok tidak cukup")
		}
	}

	// 2. EKSEKUSI PENAHANAN STOK (RESERVASI RESMI)
	rincian := make([]CatatanReservasi, 0)
	for _, it := range items {
		namaV := it.NamaVarian
		if strings.TrimSpace(namaV) == "" {
			namaV = it.NamaProduk
		}

		catatan := CatatanReservasi{
			ID:                int64(len(memoriReservasi) + 1),
			NomorPesanan:      nomorPesanan,
			IDVarian:          it.IDVarian,
			NamaVarian:        namaV,
			JumlahDireservasi: it.Jumlah,
			StatusReservasi:   "menahan_stok",
			KadaluarsaPada:    batasWaktu,
			KadaluarsaFormat:  batasWaktu.Format("02 Jan 2006, 15:04 WIB"),
			DibuatPada:        now,
		}

		if h.DB != nil {
			query := `
			INSERT INTO reservasi_stok (
				nomor_pesanan, id_varian, nama_varian, jumlah_direservasi, status_reservasi, kadaluarsa_pada, dibuat_pada
			) VALUES (?, ?, ?, ?, ?, ?, ?)`
			res, err := h.DB.Exec(query, nomorPesanan, it.IDVarian, namaV, it.Jumlah, "menahan_stok", batasWaktu, now)
			if err == nil {
				idBaru, _ := res.LastInsertId()
				catatan.ID = idBaru
			}
		}

		memoriReservasi = append(memoriReservasi, catatan)
		rincian = append(rincian, catatan)
	}

	return &HasilReservasiStok{
		Sukses:           true,
		KodeHasil:        "RESERVASI_SUKSES",
		Pesan:            "Stok pesanan berhasil ditahan selama 24 jam.",
		HoldingTimeJam:   24,
		BatasWaktu:       batasWaktu,
		BatasWaktuFormat: batasWaktu.Format("02 Jan 2006, 15:04 WIB"),
		RincianReservasi: rincian,
	}, nil
}

// LepasStokKadaluarsa melepaskan stok otomatis jika batas holding time 24 jam terlewati
func (h *HandlerReservasiStok) LepasStokKadaluarsa() int {
	muStok.Lock()
	defer muStok.Unlock()

	now := time.Now()
	jumlahDilepas := 0

	for i := range memoriReservasi {
		if memoriReservasi[i].StatusReservasi == "menahan_stok" && now.After(memoriReservasi[i].KadaluarsaPada) {
			memoriReservasi[i].StatusReservasi = "dilepas_otomatis"
			jumlahDilepas++
		}
	}

	if h.DB != nil {
		query := `
		UPDATE reservasi_stok 
		SET status_reservasi = 'dilepas_otomatis' 
		WHERE status_reservasi = 'menahan_stok' AND kadaluarsa_pada < NOW()`
		_, _ = h.DB.Exec(query)

		queryPesanan := `
		UPDATE pesanan 
		SET status_reservasi = 'kadaluarsa', status_pesanan = 'dibatalkan' 
		WHERE status_reservasi = 'direservasi' AND batas_waktu_reservasi < NOW()`
		_, _ = h.DB.Exec(queryPesanan)
	}

	return jumlahDilepas
}

// ============================================================================
// HTTP ENDPOINTS HANDLERS
// ============================================================================

// HandlerBuatPesanan (POST /api/pesanan/buat)
// Memproses pembuatan pesanan resmi, snapshot harga, dan penahanan stok 24 jam (PPL-18)
func (h *HandlerReservasiStok) HandlerBuatPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode HTTP harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanBuatPesananPPL18
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format JSON request tidak valid"}`, http.StatusBadRequest)
		return
	}

	// 1. PENGUJIAN ERROR SIMULASI STOK HABIS (WBS 1.2.9.2)
	if req.SimulasiStokHabis {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "gagal",
			"kode":      "STOK_TIDAK_CUKUP",
			"pesan":     "Stok varian produk tidak mencukupi untuk membuat pesanan. Reservasi stok gagal.",
			"data_stok": map[string]int{"tersedia": 0, "diminta": 5},
		})
		return
	}

	// 2. DEFAULT ITEM JIKA KOSONG (2 Rengginang Udang, 1 Rengginang Bawang)
	if len(req.Items) == 0 {
		req.Items = []ordergenerator.ItemSnapshotPesanan{
			{
				IDProduk:            "prod-1",
				NamaProduk:          "Rengginang Udang",
				NamaVarian:          "Udang Gurih",
				Jumlah:              2,
				HargaSatuanSnapshot: 28000,
			},
			{
				IDProduk:            "prod-2",
				NamaProduk:          "Rengginang Bawang",
				NamaVarian:          "Bawang Renyah",
				Jumlah:              1,
				HargaSatuanSnapshot: 28000,
			},
		}
	}

	// 3. DEFAULT KURIR JIKA KOSONG (Kecuali Ambil Sendiri / Pickup PPL-14)
	isPickup := req.Kurir.IDKurir == "pickup" || req.TipeCheckout == "pickup" || strings.Contains(strings.ToLower(req.Kurir.NamaKurir), "ambil")
	if req.Kurir.TarifOngkir <= 0 && !isPickup {
		req.Kurir = ordergenerator.DataKurirSnapshot{
			IDKurir:          "jne",
			NamaKurir:        "JNE REG",
			KodeLayanan:      "REG",
			TarifOngkir:      36000,
			EstimasiHari:     "2-3 hari",
			TerverifikasiAPI: true,
		}
	} else if isPickup {
		req.Kurir.TarifOngkir = 0
		if req.Kurir.NamaKurir == "" {
			req.Kurir.NamaKurir = "Ambil Sendiri di Toko"
			req.Kurir.IDKurir = "pickup"
			req.Kurir.KodeLayanan = "PICKUP"
			req.Kurir.EstimasiHari = "Langsung di Toko"
		}
	}

	// 4. MEMBUAT SNAPSHOT PESANAN RESMI (WBS 1.2.9.1)
	pesanan, err := h.Generator.BuatSnapshotDanSimpanOrder(
		req.TipeCheckout,
		req.IDPelanggan,
		req.Kontak,
		req.Kurir,
		req.Items,
	)
	if err != nil {
		http.Error(w, `{"status":"error","pesan":"Gagal mencatat snapshot pesanan"}`, http.StatusInternalServerError)
		return
	}

	// 5. RESERVASI STOK ATOMIK SELAMA 24 JAM (WBS 1.2.9.2)
	hasilReservasi, errReservasi := h.CekDanReservasiStok(pesanan.NomorPesanan, req.Items)
	if errReservasi != nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "gagal",
			"kode":      "STOK_TIDAK_CUKUP",
			"pesan":     hasilReservasi.Pesan,
			"reservasi": hasilReservasi,
		})
		return
	}

	// 6. RESPON SUKSES UNTUK HALAMAN KONFIRMASI (PPL-18 GAMBAR FIGMA)
	totalFmt := fmt.Sprintf("Rp%s", formatRibuan(pesanan.SnapshotTotal))

	namaPenerima := req.Kontak.NamaPenerima
	if strings.TrimSpace(namaPenerima) == "" {
		namaPenerima = "Pelanggan Contoh"
	}

	resp := ResponsKonfirmasiPesananPPL18{
		Status:                "sukses",
		Pesan:                 "Pesanan berhasil dibuat dan stok telah direservasi.",
		NomorPesanan:          pesanan.NomorPesanan,
		NomorInvoice:          pesanan.NomorInvoice,
		StatusPesanan:         pesanan.StatusPesanan,
		StatusPesananLabel:    "Menunggu pembayaran", // Badge hijau sesuai Figma
		StatusReservasi:       "direservasi",
		BatasWaktuPembayaran:  pesanan.BatasWaktuFormat,
		HoldingTimeJam:        24,
		TotalPembayaran:       pesanan.SnapshotTotal,
		TotalPembayaranFormat: totalFmt,
		Penerima:              namaPenerima,
		Alamat:                req.Kontak.AlamatLengkap,
		Kurir:                 req.Kurir.NamaKurir,
		HasilReservasi:        hasilReservasi,
		PesananDetail:         pesanan,
		WaktuDibuat:           pesanan.DibuatPada,
	}

	// 7. Trigger catatan pemberitahuan pesanan baru masuk (WBS 1.3.3.2)
	if h.DB != nil {
		metodeStr := "Kirim Kurir"
		if isPickup {
			metodeStr = "Ambil Sendiri"
		}
		_, _ = h.DB.Exec(`
			INSERT INTO notifikasi_pesanan (
				nomor_pesanan, nomor_invoice, tipe_pemberitahuan, judul, pesan,
				metode_penerimaan, total_transaksi, status_pesanan, sudah_dibaca, dibuat_pada
			) VALUES (
				?, ?, 'pesanan_masuk', CONCAT('Pesanan Baru: ', ?),
				CONCAT('Pesanan baru masuk dari ', ?, ' dengan metode ', ?),
				?, ?, 'baru', 0, NOW()
			)`,
			pesanan.NomorPesanan, pesanan.NomorInvoice, pesanan.NomorPesanan,
			namaPenerima, metodeStr, metodeStr, pesanan.SnapshotTotal,
		)
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// HandlerDetailPesanan (GET /api/pesanan/detail)
// Mengambil detail pesanan asli dari backend berdasarkan query nomor_pesanan atau order_id
func (h *HandlerReservasiStok) HandlerDetailPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	nomorPesanan := r.URL.Query().Get("nomor_pesanan")
	if nomorPesanan == "" {
		nomorPesanan = r.URL.Query().Get("order_id")
	}

	if nomorPesanan == "" {
		http.Error(w, `{"status":"error","pesan":"Parameter nomor_pesanan wajib disertakan"}`, http.StatusBadRequest)
		return
	}

	pesanan, ada := h.Generator.AmbilPesananBerdasarkanNomor(nomorPesanan)
	if !ada || pesanan == nil {
		// Cek fallback jika nomor berformat PJ-DEMO
		if strings.HasPrefix(nomorPesanan, "PJ-DEMO") {
			pesanan = &ordergenerator.EntitasPesanan{
				NomorPesanan:     nomorPesanan,
				NomorInvoice:     "INV-DEMO-0001",
				StatusPesanan:    "menunggu_pembayaran",
				StatusReservasi:  "direservasi",
				SnapshotSubtotal: 84000,
				SnapshotOngkir:   36000,
				SnapshotTotal:    120000,
				BatasWaktuFormat: time.Now().Add(24 * time.Hour).Format("02 Jan 2006, 15:04 WIB"),
			}
		} else {
			http.Error(w, `{"status":"error","pesan":"Data pesanan tidak ditemukan di database"}`, http.StatusNotFound)
			return
		}
	}

	totalFmt := fmt.Sprintf("Rp%s", formatRibuan(pesanan.SnapshotTotal))

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"data": ResponsKonfirmasiPesananPPL18{
			Status:                "sukses",
			NomorPesanan:          pesanan.NomorPesanan,
			NomorInvoice:          pesanan.NomorInvoice,
			StatusPesanan:         pesanan.StatusPesanan,
			StatusPesananLabel:    "Menunggu pembayaran",
			StatusReservasi:       pesanan.StatusReservasi,
			BatasWaktuPembayaran:  pesanan.BatasWaktuFormat,
			HoldingTimeJam:        24,
			TotalPembayaran:       pesanan.SnapshotTotal,
			TotalPembayaranFormat: totalFmt,
			Penerima:              pesanan.Kontak.NamaPenerima,
			Alamat:                pesanan.Kontak.AlamatLengkap,
			Kurir:                 pesanan.Kurir.NamaKurir,
			PesananDetail:         pesanan,
		},
	})
}

// HandlerDaftarReservasi (GET /api/pesanan/reservasi)
func (h *HandlerReservasiStok) HandlerDaftarReservasi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	muStok.RLock()
	defer muStok.RUnlock()

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"total":  len(memoriReservasi),
		"data":   memoriReservasi,
	})
}

// formatRibuan memformat integer ke rupiah
func formatRibuan(n int) string {
	str := fmt.Sprintf("%d", n)
	var res []string
	for len(str) > 3 {
		res = append([]string{str[len(str)-3:]}, res...)
		str = str[:len(str)-3]
	}
	if len(str) > 0 {
		res = append([]string{str}, res...)
	}
	return strings.Join(res, ".")
}
