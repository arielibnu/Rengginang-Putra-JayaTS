package notifikasipesanan

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

// LayananNotifikasi mengelola catatan notifikasi pesanan untuk Owner Dashboard
type LayananNotifikasi struct {
	DB       *sql.DB
	Mu       sync.RWMutex
	InMemori []NotifikasiPesanan
}

// InisialisasiLayanan membuat instance layanan notifikasi dengan koneksi database dan in-memory fallback
func InisialisasiLayanan(db *sql.DB) *LayananNotifikasi {
	layanan := &LayananNotifikasi{
		DB:       db,
		InMemori: make([]NotifikasiPesanan, 0),
	}

	// Buat seed fallback in-memory awal sesuai standar Figma PPL-22
	layanan.InMemori = append(layanan.InMemori,
		NotifikasiPesanan{
			ID:                   1,
			NomorPesanan:         "PJ-DEMO-0001",
			NomorInvoice:         "INV-PJ-DEMO-0001",
			TipePemberitahuan:    "pesanan_masuk",
			Judul:                "Pesanan Baru: PJ-DEMO-0001",
			Pesan:                "Pesanan baru masuk dengan metode Ambil Sendiri di Toko",
			MetodePenerimaan:     "Ambil Sendiri",
			TotalTransaksi:       84000,
			TotalTransaksiFormat: "Rp84.000",
			StatusPesanan:        "baru",
			SudahDibaca:          false,
			DibuatPada:           time.Now().Add(-15 * time.Minute),
			WaktuFormat:          "15 menit lalu",
		},
		NotifikasiPesanan{
			ID:                   2,
			NomorPesanan:         "PJ-DEMO-0002",
			NomorInvoice:         "INV-PJ-DEMO-0002",
			TipePemberitahuan:    "pesanan_terbayar",
			Judul:                "Pesanan Baru: PJ-DEMO-0002",
			Pesan:                "Pesanan baru terbayar via QRIS dengan metode Kirim Kurir",
			MetodePenerimaan:     "Kirim Kurir",
			TotalTransaksi:       120000,
			TotalTransaksiFormat: "Rp120.000",
			StatusPesanan:        "baru",
			SudahDibaca:          false,
			DibuatPada:           time.Now().Add(-5 * time.Minute),
			WaktuFormat:          "5 menit lalu",
		},
	)

	if db != nil {
		layanan.InisialisasiTabel()
	}

	return layanan
}

// InisialisasiTabel memastikan tabel notifikasi_pesanan ada di MySQL
func (l *LayananNotifikasi) InisialisasiTabel() {
	if l.DB == nil {
		return
	}

	queryTabel := `
	CREATE TABLE IF NOT EXISTS notifikasi_pesanan (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		nomor_invoice VARCHAR(64) DEFAULT NULL,
		tipe_pemberitahuan VARCHAR(50) NOT NULL DEFAULT 'pesanan_masuk',
		judul VARCHAR(100) NOT NULL,
		pesan TEXT NOT NULL,
		metode_penerimaan VARCHAR(50) NOT NULL DEFAULT 'Kirim Kurir',
		total_transaksi BIGINT NOT NULL DEFAULT 0,
		status_pesanan VARCHAR(30) NOT NULL DEFAULT 'baru',
		sudah_dibaca TINYINT(1) NOT NULL DEFAULT 0,
		waktu_dibaca DATETIME DEFAULT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_nomor_pesanan (nomor_pesanan),
		INDEX idx_tipe_pemberitahuan (tipe_pemberitahuan),
		INDEX idx_status_pesanan (status_pesanan),
		INDEX idx_sudah_dibaca (sudah_dibaca),
		INDEX idx_dibuat_pada (dibuat_pada)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	if _, err := l.DB.Exec(queryTabel); err != nil {
		log.Println("Peringatan inisialisasi tabel notifikasi_pesanan:", err)
	} else {
		log.Println("Tabel notifikasi_pesanan siap di MySQL.")
	}

	// Sinkronisasi data awal jika tabel masih kosong
	var hitung int
	_ = l.DB.QueryRow("SELECT COUNT(*) FROM notifikasi_pesanan").Scan(&hitung)
	if hitung == 0 {
		for _, seed := range l.InMemori {
			_, _ = l.DB.Exec(`
				INSERT INTO notifikasi_pesanan (
					nomor_pesanan, nomor_invoice, tipe_pemberitahuan, judul, pesan,
					metode_penerimaan, total_transaksi, status_pesanan, sudah_dibaca, dibuat_pada
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				seed.NomorPesanan, seed.NomorInvoice, seed.TipePemberitahuan, seed.Judul, seed.Pesan,
				seed.MetodePenerimaan, seed.TotalTransaksi, seed.StatusPesanan, 0, seed.DibuatPada,
			)
		}
	}
}

// CatatPemberitahuan memicu pembuatan catatan notifikasi pesanan baru atau pesanan terbayar (WBS 1.3.3.2)
func (l *LayananNotifikasi) CatatPemberitahuan(req PermintaanTriggerNotifikasi) (*NotifikasiPesanan, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	req.NomorPesanan = strings.TrimSpace(req.NomorPesanan)
	if req.NomorPesanan == "" {
		return nil, fmt.Errorf("nomor_pesanan wajib diisi")
	}

	if req.TipePemberitahuan == "" {
		req.TipePemberitahuan = "pesanan_masuk"
	}

	if req.MetodePenerimaan == "" {
		req.MetodePenerimaan = "Kirim Kurir"
	}

	if req.Judul == "" {
		if req.TipePemberitahuan == "pesanan_terbayar" {
			req.Judul = fmt.Sprintf("Pembayaran Lunas: %s", req.NomorPesanan)
		} else {
			req.Judul = fmt.Sprintf("Pesanan Baru: %s", req.NomorPesanan)
		}
	}

	if req.Pesan == "" {
		if req.TipePemberitahuan == "pesanan_terbayar" {
			req.Pesan = fmt.Sprintf("Pesanan %s telah berhasil dibayar. Segera siapkan pesanan untuk dikirim atau diambil.", req.NomorPesanan)
		} else {
			req.Pesan = fmt.Sprintf("Pesanan baru masuk (%s). Harap periksa rincian pada Dashboard.", req.MetodePenerimaan)
		}
	}

	now := time.Now()
	baru := NotifikasiPesanan{
		ID:                   time.Now().UnixNano(),
		NomorPesanan:         req.NomorPesanan,
		NomorInvoice:         req.NomorInvoice,
		TipePemberitahuan:    req.TipePemberitahuan,
		Judul:                req.Judul,
		Pesan:                req.Pesan,
		MetodePenerimaan:     req.MetodePenerimaan,
		TotalTransaksi:       req.TotalTransaksi,
		TotalTransaksiFormat: fmt.Sprintf("Rp%s", formatRibuan(req.TotalTransaksi)),
		StatusPesanan:        "baru",
		SudahDibaca:          false,
		DibuatPada:           now,
		WaktuFormat:          "Baru saja",
	}

	// 1. Simpan ke database jika koneksi tersedia
	if l.DB != nil {
		res, err := l.DB.Exec(`
			INSERT INTO notifikasi_pesanan (
				nomor_pesanan, nomor_invoice, tipe_pemberitahuan, judul, pesan,
				metode_penerimaan, total_transaksi, status_pesanan, sudah_dibaca, dibuat_pada
			) VALUES (?, ?, ?, ?, ?, ?, ?, 'baru', 0, ?)`,
			baru.NomorPesanan, baru.NomorInvoice, baru.TipePemberitahuan, baru.Judul, baru.Pesan,
			baru.MetodePenerimaan, baru.TotalTransaksi, now,
		)
		if err == nil {
			id, _ := res.LastInsertId()
			baru.ID = id
		} else {
			log.Println("Peringatan simpan notifikasi ke MySQL:", err)
		}
	}

	// 2. Simpan ke in-memory cache (taruh di urutan teratas)
	l.InMemori = append([]NotifikasiPesanan{baru}, l.InMemori...)

	log.Printf("[PPL-22 BE] Trigger notifikasi berhasil dicatat: %s (%s - %s)", baru.NomorPesanan, baru.TipePemberitahuan, baru.MetodePenerimaan)
	return &baru, nil
}

// AmbilPesananBaruOwner mengambil daftar pesanan baru yang belum dibaca/diproses khusus Dashboard Owner
func (l *LayananNotifikasi) AmbilPesananBaruOwner() (*ResponsPesananBaruOwner, error) {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	var daftarNotif []NotifikasiPesanan

	// 1. Coba ambil langsung dari MySQL jika tersedia
	if l.DB != nil {
		rows, err := l.DB.Query(`
			SELECT id, nomor_pesanan, COALESCE(nomor_invoice, ''), tipe_pemberitahuan,
			       judul, pesan, metode_penerimaan, total_transaksi, status_pesanan,
			       sudah_dibaca, dibuat_pada
			FROM notifikasi_pesanan
			WHERE sudah_dibaca = 0 OR status_pesanan = 'baru'
			ORDER BY id DESC LIMIT 10`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var n NotifikasiPesanan
				var dibacaInt int
				if errScan := rows.Scan(
					&n.ID, &n.NomorPesanan, &n.NomorInvoice, &n.TipePemberitahuan,
					&n.Judul, &n.Pesan, &n.MetodePenerimaan, &n.TotalTransaksi,
					&n.StatusPesanan, &dibacaInt, &n.DibuatPada,
				); errScan == nil {
					n.SudahDibaca = dibacaInt == 1
					n.TotalTransaksiFormat = fmt.Sprintf("Rp%s", formatRibuan(n.TotalTransaksi))
					n.WaktuFormat = formatWaktuRelatif(n.DibuatPada)
					daftarNotif = append(daftarNotif, n)
				}
			}
		}
	}

	// 2. Jika dari database belum ada atau kosong, fallback ke in-memory
	if len(daftarNotif) == 0 {
		for _, item := range l.InMemori {
			if !item.SudahDibaca || item.StatusPesanan == "baru" {
				itemCopy := item
				itemCopy.WaktuFormat = formatWaktuRelatif(item.DibuatPada)
				daftarNotif = append(daftarNotif, itemCopy)
			}
		}
	}

	// 3. Susun daftar ringkas unik sesuai Figma (Nomor — Metode)
	ringkasList := make([]ItemNotifikasiRingkas, 0)
	sudahAda := make(map[string]bool)
	daftarUnik := make([]NotifikasiPesanan, 0)

	for _, n := range daftarNotif {
		if !sudahAda[n.NomorPesanan] {
			sudahAda[n.NomorPesanan] = true
			daftarUnik = append(daftarUnik, n)
			ringkasList = append(ringkasList, ItemNotifikasiRingkas{
				NomorPesanan: n.NomorPesanan,
				Metode:       n.MetodePenerimaan,
				TeksRingkas:  fmt.Sprintf("%s — %s", n.NomorPesanan, n.MetodePenerimaan),
				Tipe:         n.TipePemberitahuan,
				Waktu:        n.WaktuFormat,
			})
		}
	}

	jumlahBaru := len(ringkasList)
	badgeLabel := fmt.Sprintf("%d Baru", jumlahBaru)

	return &ResponsPesananBaruOwner{
		Status:        "sukses",
		Pesan:         "Daftar notifikasi pesanan baru berhasil diambil",
		JumlahBaru:    jumlahBaru,
		BadgeLabel:    badgeLabel,
		DaftarPesanan: ringkasList,
		Notifikasi:    daftarUnik,
	}, nil
}

// TandaiSudahDibaca menandai notifikasi telah dibaca oleh Owner
func (l *LayananNotifikasi) TandaiSudahDibaca(req PermintaanTandaiDibaca) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	now := time.Now()

	// 1. Update ke MySQL
	if l.DB != nil {
		if req.Semua {
			_, _ = l.DB.Exec("UPDATE notifikasi_pesanan SET sudah_dibaca = 1, waktu_dibaca = ? WHERE sudah_dibaca = 0", now)
		} else if req.NomorPesanan != "" {
			_, _ = l.DB.Exec("UPDATE notifikasi_pesanan SET sudah_dibaca = 1, waktu_dibaca = ? WHERE nomor_pesanan = ?", now, req.NomorPesanan)
		} else if len(req.IDNotifikasi) > 0 {
			for _, id := range req.IDNotifikasi {
				_, _ = l.DB.Exec("UPDATE notifikasi_pesanan SET sudah_dibaca = 1, waktu_dibaca = ? WHERE id = ?", now, id)
			}
		}
	}

	// 2. Update cache in-memory
	for i := range l.InMemori {
		if req.Semua || (req.NomorPesanan != "" && l.InMemori[i].NomorPesanan == req.NomorPesanan) {
			l.InMemori[i].SudahDibaca = true
			l.InMemori[i].WaktuDibaca = &now
		} else {
			for _, id := range req.IDNotifikasi {
				if l.InMemori[i].ID == id {
					l.InMemori[i].SudahDibaca = true
					l.InMemori[i].WaktuDibaca = &now
				}
			}
		}
	}

	return nil
}

// ============================================================================
// HTTP HANDLERS (WBS 1.3.3.2)
// ============================================================================

// HandlerPesananBaru (GET /api/owner/notifikasi/pesanan-baru)
// Endpoint khusus untuk komponen Notifikasi Pesanan Baru pada Dashboard Owner Figma
func (l *LayananNotifikasi) HandlerPesananBaru(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": "Metode tidak diizinkan"})
		return
	}

	hasil, err := l.AmbilPesananBaruOwner()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hasil)
}

// HandlerDaftarSemua (GET /api/owner/notifikasi)
// Mengambil semua riwayat notifikasi dengan filter status
func (l *LayananNotifikasi) HandlerDaftarSemua(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hasil, err := l.AmbilPesananBaruOwner()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "sukses",
		"data":       hasil.Notifikasi,
		"total_baru": hasil.JumlahBaru,
	})
}

// HandlerTandaiDibaca (POST /api/owner/notifikasi/baca)
// Menandai satu atau semua notifikasi telah dibaca
func (l *LayananNotifikasi) HandlerTandaiDibaca(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanTandaiDibaca
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": "Format JSON tidak valid"})
		return
	}

	if err := l.TandaiSudahDibaca(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "sukses",
		"pesan":  "Notifikasi berhasil ditandai telah dibaca",
	})
}

// HandlerTrigger (POST /api/owner/notifikasi/trigger)
// Endpoint trigger manual/otomatis dari modul lain (PPL-19 Pembayaran / PPL-18 Pembuatan Pesanan)
func (l *LayananNotifikasi) HandlerTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanTriggerNotifikasi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": "Format JSON tidak valid"})
		return
	}

	notif, err := l.CatatPemberitahuan(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"pesan":  "Trigger catatan pemberitahuan berhasil diproses",
		"data":   notif,
	})
}

// formatRibuan helper pemformat ribuan angka rupiah
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

// formatWaktuRelatif helper waktu relatif ramah pengguna
func formatWaktuRelatif(t time.Time) string {
	dur := time.Since(t)
	if dur < time.Minute {
		return "Baru saja"
	} else if dur < time.Hour {
		return fmt.Sprintf("%d menit lalu", int(dur.Minutes()))
	} else if dur < 24*time.Hour {
		return fmt.Sprintf("%d jam lalu", int(dur.Hours()))
	}
	return t.Format("02 Jan 15:04")
}
