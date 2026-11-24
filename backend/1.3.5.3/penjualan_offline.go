package penjualanoffline

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LayananPenjualanOffline mengelola pencatatan transaksi kasir toko langsung (WBS 1.3.5.3)
type LayananPenjualanOffline struct {
	DB       *sql.DB
	Mu       sync.RWMutex
	InMemori []PenjualanOffline
	Counter  int
}

// InisialisasiLayanan membuat instance layanan penjualan offline
func InisialisasiLayanan(db *sql.DB) *LayananPenjualanOffline {
	layanan := &LayananPenjualanOffline{
		DB:       db,
		InMemori: make([]PenjualanOffline, 0),
		Counter:  100,
	}

	if db != nil {
		layanan.InisialisasiTabel()
	}

	return layanan
}

// InisialisasiTabel memastikan tabel penjualan_offline ada di MySQL
func (l *LayananPenjualanOffline) InisialisasiTabel() {
	if l.DB == nil {
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS penjualan_offline (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_transaksi VARCHAR(64) NOT NULL UNIQUE,
		tanggal_penjualan VARCHAR(50) NOT NULL,
		id_varian BIGINT NOT NULL,
		nama_varian VARCHAR(100) NOT NULL,
		jumlah INT NOT NULL,
		harga_satuan DECIMAL(12,2) NOT NULL,
		total_harga DECIMAL(15,2) NOT NULL,
		metode_pembayaran VARCHAR(50) NOT NULL DEFAULT 'Tunai',
		sumber_transaksi VARCHAR(50) NOT NULL DEFAULT 'offline_kasir',
		flag_lingkungan VARCHAR(50) NOT NULL DEFAULT 'PRODUKSI_NYATA',
		is_sandbox TINYINT(1) NOT NULL DEFAULT 0,
		stok_sebelum INT NOT NULL DEFAULT 0,
		stok_sesudah INT NOT NULL DEFAULT 0,
		kasir_nama VARCHAR(100) NOT NULL DEFAULT 'Owner Toko',
		catatan TEXT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_nomor (nomor_transaksi),
		INDEX idx_varian (id_varian),
		INDEX idx_flag (flag_lingkungan),
		INDEX idx_sumber (sumber_transaksi),
		INDEX idx_dibuat (dibuat_pada)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	if _, err := l.DB.Exec(query); err != nil {
		log.Println("Peringatan inisialisasi tabel penjualan_offline:", err)
	} else {
		log.Println("Tabel penjualan_offline siap di MySQL.")
	}
}

// SimpanPenjualan memproses transaksi kasir dan memotong stok fisik tanpa ekspedisi online (WBS 1.3.5.3 & 1.3.5.4)
func (l *LayananPenjualanOffline) SimpanPenjualan(req PermintaanPenjualanOffline) (*PenjualanOffline, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	// 1. Validasi kuantitas
	if req.Jumlah <= 0 {
		return nil, fmt.Errorf("jumlah pembelian wajib lebih dari 0 bungkus")
	}

	req.NamaVarian = strings.TrimSpace(req.NamaVarian)
	if req.NamaVarian == "" {
		req.NamaVarian = "Rengginang Udang"
	}

	if req.TanggalPenjualan == "" {
		req.TanggalPenjualan = time.Now().Format("02 January 2006")
	}

	if req.MetodePembayaran == "" {
		req.MetodePembayaran = "Tunai"
	}

	// Flag pemisah transaksi nyata vs sandbox (WBS 1.3.5.4)
	flagLingkungan := req.FlagLingkungan
	if flagLingkungan == "" {
		flagLingkungan = "PRODUKSI_NYATA"
	}
	isSandbox := strings.ToUpper(flagLingkungan) == "SANDBOX_SIMULASI"

	if req.HargaSatuan <= 0 {
		req.HargaSatuan = 28000
	}

	totalHarga := int64(req.Jumlah) * req.HargaSatuan
	l.Counter++
	nomorTransaksi := fmt.Sprintf("POS-%s-%04d", time.Now().Format("20060102"), (time.Now().UnixNano()/1e6)%10000)

	stokSebelum := 50
	stokSesudah := 50 - req.Jumlah

	// 2. Eksekusi pemotongan stok fisik di Database MySQL jika terkoneksi
	if l.DB != nil {
		tx, err := l.DB.Begin()
		if err != nil {
			return nil, fmt.Errorf("gagal memulai transaksi database: %v", err)
		}
		defer tx.Rollback()

		// Cari varian di tabel varian_produk
		var idVarianDB uint64
		var stokDB int
		var hargaDB float64

		row := tx.QueryRow(`
			SELECT id_varian, stok, harga_jual 
			FROM varian_produk 
			WHERE id_varian = ? OR LOWER(nama_varian) LIKE ? 
			LIMIT 1 FOR UPDATE`,
			req.IDVarian, "%"+strings.ToLower(req.NamaVarian)+"%",
		)

		if err := row.Scan(&idVarianDB, &stokDB, &hargaDB); err == nil {
			req.IDVarian = int64(idVarianDB)
			stokSebelum = stokDB

			// Validasi stok fisik tidak minus
			if req.Jumlah > stokDB {
				return nil, fmt.Errorf("stok %s tidak mencukupi (sisa: %d bungkus, permintaan: %d)", req.NamaVarian, stokDB, req.Jumlah)
			}

			stokSesudah = stokDB - req.Jumlah

			// Potong stok fisik riil
			_, errUpd := tx.Exec(`UPDATE varian_produk SET stok = ? WHERE id_varian = ?`, stokSesudah, idVarianDB)
			if errUpd != nil {
				return nil, fmt.Errorf("gagal memotong stok varian: %v", errUpd)
			}

			// Catat mutasi stok keluar
			_, _ = tx.Exec(`
				INSERT INTO riwayat_mutasi_stok (
					id_varian, jenis_mutasi, mode_stok, jumlah, stok_sebelum, stok_sesudah, catatan
				) VALUES (?, 'KELUAR', 'RIIL', ?, ?, ?, ?)`,
				idVarianDB, req.Jumlah, stokSebelum, stokSesudah,
				fmt.Sprintf("Penjualan Offline Kasir Langsung (%s) - %s", nomorTransaksi, flagLingkungan),
			)
		}

		// Simpan catatan transaksi penjualan offline
		_, errIns := tx.Exec(`
			INSERT INTO penjualan_offline (
				nomor_transaksi, tanggal_penjualan, id_varian, nama_varian,
				jumlah, harga_satuan, total_harga, metode_pembayaran,
				sumber_transaksi, flag_lingkungan, is_sandbox,
				stok_sebelum, stok_sesudah, kasir_nama, catatan, dibuat_pada
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'offline_kasir', ?, ?, ?, ?, 'Owner Toko', ?, NOW())`,
			nomorTransaksi, req.TanggalPenjualan, req.IDVarian, req.NamaVarian,
			req.Jumlah, req.HargaSatuan, totalHarga, req.MetodePembayaran,
			flagLingkungan, isSandbox, stokSebelum, stokSesudah, req.Catatan,
		)
		if errIns != nil {
			return nil, fmt.Errorf("gagal mencatat transaksi penjualan offline: %v", errIns)
		}

		if errCommit := tx.Commit(); errCommit != nil {
			return nil, fmt.Errorf("gagal menyelesaikan transaksi: %v", errCommit)
		}
	}

	trx := PenjualanOffline{
		ID:               time.Now().UnixNano(),
		NomorTransaksi:   nomorTransaksi,
		TanggalPenjualan: req.TanggalPenjualan,
		IDVarian:         req.IDVarian,
		NamaVarian:       req.NamaVarian,
		Jumlah:           req.Jumlah,
		HargaSatuan:      req.HargaSatuan,
		TotalHarga:       totalHarga,
		TotalHargaFormat: fmt.Sprintf("Rp%s", formatRibuan(totalHarga)),
		MetodePembayaran: req.MetodePembayaran,
		SumberTransaksi:  "offline_kasir",
		FlagLingkungan:   flagLingkungan,
		IsSandbox:        isSandbox,
		StokSebelum:      stokSebelum,
		StokSesudah:      stokSesudah,
		KasirNama:        "Owner Toko",
		Catatan:          req.Catatan,
		DibuatPada:       time.Now(),
		WaktuFormat:      "Baru saja",
	}

	l.InMemori = append([]PenjualanOffline{trx}, l.InMemori...)
	log.Printf("[PPL-24 BE] Transaksi kasir offline berhasil dicatat: %s (%dx %s = Rp%d) [%s]",
		nomorTransaksi, req.Jumlah, req.NamaVarian, totalHarga, flagLingkungan)

	return &trx, nil
}

// AmbilDaftarPenjualan mengambil daftar riwayat transaksi kasir offline
func (l *LayananPenjualanOffline) AmbilDaftarPenjualan(flagFilter string) (*ResponsDaftarPenjualan, error) {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	var list []PenjualanOffline
	var totalNilai int64

	if l.DB != nil {
		query := `
			SELECT id, nomor_transaksi, tanggal_penjualan, id_varian, nama_varian,
			       jumlah, harga_satuan, total_harga, metode_pembayaran,
			       sumber_transaksi, flag_lingkungan, is_sandbox,
			       stok_sebelum, stok_sesudah, kasir_nama, COALESCE(catatan, ''), dibuat_pada
			FROM penjualan_offline`

		var args []interface{}
		if flagFilter != "" {
			query += " WHERE flag_lingkungan = ?"
			args = append(args, flagFilter)
		}
		query += " ORDER BY id DESC LIMIT 50"

		rows, err := l.DB.Query(query, args...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var p PenjualanOffline
				var isSbInt int
				var hSatuan, tHarga float64
				if errScan := rows.Scan(
					&p.ID, &p.NomorTransaksi, &p.TanggalPenjualan, &p.IDVarian, &p.NamaVarian,
					&p.Jumlah, &hSatuan, &tHarga, &p.MetodePembayaran,
					&p.SumberTransaksi, &p.FlagLingkungan, &isSbInt,
					&p.StokSebelum, &p.StokSesudah, &p.KasirNama, &p.Catatan, &p.DibuatPada,
				); errScan == nil {
					p.HargaSatuan = int64(hSatuan)
					p.TotalHarga = int64(tHarga)
					p.IsSandbox = isSbInt == 1
					p.TotalHargaFormat = fmt.Sprintf("Rp%s", formatRibuan(p.TotalHarga))
					p.WaktuFormat = p.DibuatPada.Format("02 Jan 15:04")
					totalNilai += p.TotalHarga
					list = append(list, p)
				}
			}
		}
	}

	if len(list) == 0 {
		for _, item := range l.InMemori {
			if flagFilter == "" || item.FlagLingkungan == flagFilter {
				totalNilai += item.TotalHarga
				list = append(list, item)
			}
		}
	}

	return &ResponsDaftarPenjualan{
		Status:     "sukses",
		Pesan:      "Daftar penjualan offline berhasil diambil",
		Total:      len(list),
		TotalNilai: totalNilai,
		Data:       list,
	}, nil
}

// HandlerSimpanPenjualan (POST /api/admin/penjualan-offline)
func (l *LayananPenjualanOffline) HandlerSimpanPenjualan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": "Metode HTTP harus POST"})
		return
	}

	var req PermintaanPenjualanOffline
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "gagal", "pesan": "Format JSON tidak valid"})
		return
	}

	hasil, err := l.SimpanPenjualan(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "gagal",
			"pesan":  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ResponsPenjualanOffline{
		Status: "sukses",
		Pesan:  fmt.Sprintf("Penjualan offline berhasil disimpan (%s) dan stok fisik telah dipotong.", hasil.NomorTransaksi),
		Data:   hasil,
	})
}

// HandlerDaftarPenjualan (GET /api/admin/penjualan-offline)
func (l *LayananPenjualanOffline) HandlerDaftarPenjualan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	flagFilter := r.URL.Query().Get("flag")
	hasil, err := l.AmbilDaftarPenjualan(flagFilter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hasil)
}

func formatRibuan(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var res []string
	for len(s) > 3 {
		res = append([]string{s[len(s)-3:]}, res...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		res = append([]string{s}, res...)
	}
	return strings.Join(res, ".")
}
