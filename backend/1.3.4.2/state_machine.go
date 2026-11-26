package statemachine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	pengembalianstok "rengginang-putra-jaya/auth/1.3.4.4"
)

// DataPesananLokal untuk state management in-memory jika DB offline / demo
type DataPesananLokal struct {
	NomorPesanan     string
	StatusPemrosesan string
	StatusPembayaran string
	NomorResi        string
}

// LayananStateMachine mengelola transisi status pemrosesan pesanan yang terpisah dari status pembayaran
type LayananStateMachine struct {
	DB               *sql.DB
	LayananStok      *pengembalianstok.LayananPengembalianStok
	MemoriPesanan    map[string]*DataPesananLokal
	Mu               sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananStateMachine
func InisialisasiLayanan(db *sql.DB, layananStok *pengembalianstok.LayananPengembalianStok) *LayananStateMachine {
	s := &LayananStateMachine{
		DB:            db,
		LayananStok:   layananStok,
		MemoriPesanan: make(map[string]*DataPesananLokal),
	}

	// Data default demo
	s.MemoriPesanan["PJ-DEMO-0001"] = &DataPesananLokal{
		NomorPesanan:     "PJ-DEMO-0001",
		StatusPemrosesan: StatusProsesBaru,
		StatusPembayaran: StatusBayarMenunggu,
	}
	s.MemoriPesanan["PJ-DEMO-0002"] = &DataPesananLokal{
		NomorPesanan:     "PJ-DEMO-0002",
		StatusPemrosesan: StatusProsesDiproses,
		StatusPembayaran: StatusBayarDibayar,
	}
	s.MemoriPesanan["PJ-DEMO-0003"] = &DataPesananLokal{
		NomorPesanan:     "PJ-DEMO-0003",
		StatusPemrosesan: StatusProsesSelesai,
		StatusPembayaran: StatusBayarDibayar,
	}

	s.InisialisasiTabel()
	return s
}

// InisialisasiTabel memastikan kolom dan tabel audit tersedia
func (l *LayananStateMachine) InisialisasiTabel() {
	if l.DB == nil {
		return
	}

	// 1. Pastikan kolom status_pemrosesan & nomor_resi ada di tabel pesanan
	_, _ = l.DB.Exec(`
		ALTER TABLE pesanan 
		ADD COLUMN IF NOT EXISTS status_pemrosesan VARCHAR(30) NOT NULL DEFAULT 'menunggu_diproses',
		ADD COLUMN IF NOT EXISTS nomor_resi VARCHAR(100) NULL DEFAULT '';
	`)

	// 2. Tabel audit riwayat perubahan status pemrosesan
	query := `
	CREATE TABLE IF NOT EXISTS riwayat_status_pemrosesan (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		status_lama VARCHAR(50) NOT NULL,
		status_baru VARCHAR(50) NOT NULL,
		nomor_resi VARCHAR(100) NULL DEFAULT '',
		diubah_oleh VARCHAR(50) NOT NULL DEFAULT 'Owner',
		catatan VARCHAR(255) NULL DEFAULT '',
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_riwayat_pemrosesan_pesanan (nomor_pesanan)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
	if _, err := l.DB.Exec(query); err != nil {
		log.Println("Peringatan inisialisasi tabel riwayat_status_pemrosesan:", err)
	}
}

// NormalisasiStatus menyamakan string status dari frontend (misal: "Diproses -> Dikirim" -> "dikirim")
func (l *LayananStateMachine) NormalisasiStatus(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))

	if strings.Contains(s, "batal") {
		return StatusProsesBatal
	}
	if strings.Contains(s, "selesai") {
		return StatusProsesSelesai
	}
	if strings.Contains(s, "dikirim") {
		return StatusProsesDikirim
	}
	if strings.Contains(s, "diproses") {
		return StatusProsesDiproses
	}
	if strings.Contains(s, "baru") || strings.Contains(s, "menunggu") {
		return StatusProsesMenungguDiproses
	}

	return s
}

// ValidasiTransisiStatus memvalidasi apakah transisi dari status lama ke status baru diizinkan oleh state machine
func (l *LayananStateMachine) ValidasiTransisiStatus(statusLamaRaw, statusBaruRaw string) (bool, string) {
	lama := l.NormalisasiStatus(statusLamaRaw)
	baru := l.NormalisasiStatus(statusBaruRaw)

	// Jika status sama, tidak ada perubahan
	if lama == baru {
		return true, ""
	}

	// Aturan State Machine Pesanan:
	switch lama {
	case StatusProsesBaru, StatusProsesMenungguDiproses:
		// Dari baru / menunggu_diproses hanya boleh ke: diproses atau batal
		if baru == StatusProsesDiproses || baru == StatusProsesBatal {
			return true, ""
		}
		return false, fmt.Sprintf("Pesanan baru tidak bisa langsung diubah menjadi '%s'. Harus 'diproses' atau 'batal'.", baru)

	case StatusProsesDiproses:
		// Dari diproses boleh ke: dikirim atau batal
		if baru == StatusProsesDikirim || baru == StatusProsesBatal {
			return true, ""
		}
		return false, fmt.Sprintf("Pesanan sedang diproses tidak bisa langsung diubah menjadi '%s'. Harus 'dikirim' atau 'batal'.", baru)

	case StatusProsesDikirim:
		// Dari dikirim hanya boleh ke: selesai (tidak boleh langsung batal karena barang sudah dalam perjalanan kurir)
		if baru == StatusProsesSelesai {
			return true, ""
		}
		return false, fmt.Sprintf("Pesanan yang sudah dikirim tidak bisa diubah menjadi '%s'. Harus diselesaikan ('selesai').", baru)

	case StatusProsesSelesai:
		// Terminal state: pesanan selesai tidak boleh diubah lagi
		return false, "Pesanan yang sudah 'selesai' merupakan status final dan tidak dapat diubah lagi."

	case StatusProsesBatal:
		// Terminal state: pesanan batal tidak boleh diaktifkan lagi
		return false, "Pesanan yang sudah 'batal' merupakan status final dan tidak dapat diubah lagi."

	default:
		// Status tidak dikenal
		return false, fmt.Sprintf("Status awal '%s' tidak valid dalam state machine.", statusLamaRaw)
	}
}

// UbahStatusPemrosesan mengeksekusi perubahan status pemrosesan dengan jaminan pemisahan status pembayaran
func (l *LayananStateMachine) UbahStatusPemrosesan(req PermintaanUbahStatus) (*ResponsUbahStatus, error) {
	nomorPesanan := strings.TrimSpace(req.NomorPesanan)
	if nomorPesanan == "" {
		return nil, fmt.Errorf("nomor pesanan wajib diisi")
	}

	statusBaruClean := l.NormalisasiStatus(req.StatusPemrosesan)
	if statusBaruClean == "" {
		return nil, fmt.Errorf("status pemrosesan baru tidak boleh kosong")
	}

	l.Mu.Lock()
	defer l.Mu.Unlock()

	// 1. Ambil status saat ini dari database atau memori
	var statusProsesLama string
	var statusPembayaranTetap string
	var pesananAda bool

	if l.DB != nil {
		var statusDB, bayarDB sql.NullString
		errQ := l.DB.QueryRow(`
			SELECT status_pemrosesan, status_pesanan 
			FROM pesanan 
			WHERE nomor_pesanan = ?`,
			nomorPesanan,
		).Scan(&statusDB, &bayarDB)

		if errQ == nil {
			pesananAda = true
			if statusDB.Valid && statusDB.String != "" {
				statusProsesLama = statusDB.String
			} else {
				statusProsesLama = StatusProsesMenungguDiproses
			}
			if bayarDB.Valid && bayarDB.String != "" {
				statusPembayaranTetap = bayarDB.String
			} else {
				statusPembayaranTetap = StatusBayarMenunggu
			}
		}
	}

	if !pesananAda {
		if memo, ada := l.MemoriPesanan[nomorPesanan]; ada {
			statusProsesLama = memo.StatusPemrosesan
			statusPembayaranTetap = memo.StatusPembayaran
			pesananAda = true
		} else {
			// Demo pesanan baru
			statusProsesLama = StatusProsesDiproses
			statusPembayaranTetap = StatusBayarDibayar
			pesananAda = true
		}
	}

	// 2. Validasi aturan transisi state machine
	valid, alasanDitolak := l.ValidasiTransisiStatus(statusProsesLama, statusBaruClean)
	if !valid {
		return nil, fmt.Errorf("transisi status ditolak: %s", alasanDitolak)
	}

	stokDikembalikan := false
	totalBungkusKembali := 0

	// 3. Jika status baru adalah 'batal', otomatis picu pengembalian stok (WBS 1.3.4.4)
	if statusBaruClean == StatusProsesBatal && l.LayananStok != nil {
		hasilKembali, errKembali := l.LayananStok.KembalikanStokPesananBatal(
			nomorPesanan,
			fmt.Sprintf("Pembatalan oleh Owner: %s", req.Catatan),
		)
		if errKembali == nil && hasilKembali != nil {
			stokDikembalikan = true
			totalBungkusKembali = hasilKembali.TotalBungkusKembali
		}
	}

	// 4. Update status_pemrosesan di database pesanan (STATUS PEMBAYARAN TETAP UTUH TIDAK DIUBAH)
	if l.DB != nil {
		query := `
			UPDATE pesanan 
			SET status_pemrosesan = ?, 
			    nomor_resi = CASE WHEN ? != '' THEN ? ELSE nomor_resi END,
			    diperbarui_pada = NOW()
			WHERE nomor_pesanan = ?`
		_, errUpd := l.DB.Exec(query, statusBaruClean, req.NomorResi, req.NomorResi, nomorPesanan)
		if errUpd != nil {
			return nil, fmt.Errorf("gagal memperbarui status pemrosesan di database: %v", errUpd)
		}

		// Catat ke tabel audit riwayat_status_pemrosesan
		_, _ = l.DB.Exec(`
			INSERT INTO riwayat_status_pemrosesan (
				nomor_pesanan, status_lama, status_baru, nomor_resi, diubah_oleh, catatan
			) VALUES (?, ?, ?, ?, 'Owner', ?)`,
			nomorPesanan, statusProsesLama, statusBaruClean, req.NomorResi, req.Catatan,
		)
	}

	// Update data memori lokal
	if memo, ada := l.MemoriPesanan[nomorPesanan]; ada {
		memo.StatusPemrosesan = statusBaruClean
		if req.NomorResi != "" {
			memo.NomorResi = req.NomorResi
		}
	} else {
		l.MemoriPesanan[nomorPesanan] = &DataPesananLokal{
			NomorPesanan:     nomorPesanan,
			StatusPemrosesan: statusBaruClean,
			StatusPembayaran: statusPembayaranTetap,
			NomorResi:        req.NomorResi,
		}
	}

	pesanSukses := fmt.Sprintf("Status pemrosesan pesanan %s berhasil diperbarui menjadi '%s'. Status pembayaran tetap '%s'.", nomorPesanan, statusBaruClean, statusPembayaranTetap)
	if stokDikembalikan {
		pesanSukses += fmt.Sprintf(" Sebanyak %d bungkus stok telah dikembalikan ke persediaan aktif.", totalBungkusKembali)
	}

	return &ResponsUbahStatus{
		Status:                "sukses",
		Pesan:                 pesanSukses,
		NomorPesanan:          nomorPesanan,
		StatusPemrosesanLama:  statusProsesLama,
		StatusPemrosesanBaru:  statusBaruClean,
		StatusPembayaranTetap: statusPembayaranTetap,
		NomorResi:             req.NomorResi,
		StokDikembalikan:      stokDikembalikan,
		TotalBungkusKembali:   totalBungkusKembali,
		WaktuPembaruan:        time.Now(),
	}, nil
}

// HandlerUbahStatus (PUT / POST /api/admin/pesanan/status)
func (l *LayananStateMachine) HandlerUbahStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode harus PUT atau POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanUbahStatus
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	respons, err := l.UbahStatusPemrosesan(req)
	if err != nil {
		// Validasi transisi ditolak menghasilkan 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"pesan":  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respons)
}

// HandlerOpsiTransisi (GET /api/admin/pesanan/transisi-status?status=...)
func (l *LayananStateMachine) HandlerOpsiTransisi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	statusSaatIni := l.NormalisasiStatus(r.URL.Query().Get("status"))

	var opsi []map[string]string

	switch statusSaatIni {
	case StatusProsesBaru, StatusProsesMenungguDiproses:
		opsi = []map[string]string{
			{"nilai": StatusProsesDiproses, "label": "Menunggu Diproses -> Diproses"},
			{"nilai": StatusProsesBatal, "label": "Batal (Kembalikan Reservasi Stok)"},
		}
	case StatusProsesDiproses:
		opsi = []map[string]string{
			{"nilai": StatusProsesDikirim, "label": "Diproses -> Dikirim"},
			{"nilai": StatusProsesBatal, "label": "Batal (Kembalikan Reservasi Stok)"},
		}
	case StatusProsesDikirim:
		opsi = []map[string]string{
			{"nilai": StatusProsesSelesai, "label": "Dikirim -> Selesai"},
		}
	default:
		opsi = []map[string]string{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "sukses",
		"status_sekarang": statusSaatIni,
		"opsi_transisi":   opsi,
	})
}
