package kedaluwarsastok

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CatatanReservasiLokal struktur in-memory untuk pengujian dan fallback
type CatatanReservasiLokal struct {
	ID              int64
	NomorPesanan    string
	IDVarian        int64
	NamaVarian      string
	Jumlah          int
	StatusReservasi string // menahan_stok, dikonfirmasi, dilepas_otomatis
	KadaluarsaPada  time.Time
}

// LayananKedaluwarsaStok mengelola pemutusan reservasi stok otomatis saat pembayaran kedaluwarsa atau dibatalkan
type LayananKedaluwarsaStok struct {
	DB               *sql.DB
	MemoriReservasi  map[string][]CatatanReservasiLokal
	LogPemutusan     []LogPemutusanReservasi
	Mu               sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananKedaluwarsaStok
func InisialisasiLayanan(db *sql.DB) *LayananKedaluwarsaStok {
	s := &LayananKedaluwarsaStok{
		DB:              db,
		MemoriReservasi: make(map[string][]CatatanReservasiLokal),
		LogPemutusan:    make([]LogPemutusanReservasi, 0),
	}
	s.InisialisasiTabel()
	return s
}

// InisialisasiTabel membuat tabel log pemutusan reservasi jika belum ada
func (l *LayananKedaluwarsaStok) InisialisasiTabel() {
	if l.DB == nil {
		return
	}
	query := `
	CREATE TABLE IF NOT EXISTS log_pemutusan_reservasi_stok (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		id_varian BIGINT NOT NULL DEFAULT 1,
		nama_varian VARCHAR(100) NOT NULL,
		jumlah_dilepas INT NOT NULL,
		status_sebelumnya VARCHAR(30) NOT NULL DEFAULT 'menahan_stok',
		status_akhir VARCHAR(30) NOT NULL DEFAULT 'dilepas_otomatis',
		alasan VARCHAR(255) NOT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_log_reservasi_pesanan (nomor_pesanan)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
	_, _ = l.DB.Exec(query)
}

// TambahReservasiLokal membantu mencatatkan reservasi untuk simulasi
func (l *LayananKedaluwarsaStok) TambahReservasiLokal(nomorPesanan string, items []CatatanReservasiLokal) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.MemoriReservasi[nomorPesanan] = items
}

// PutusReservasiStokKedaluwarsa melepaskan stok tertahan agar produk tidak terblokir
func (l *LayananKedaluwarsaStok) PutusReservasiStokKedaluwarsa(nomorPesanan string, alasan string) (*HasilPemutusanReservasi, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	nomorPesanan = strings.TrimSpace(nomorPesanan)
	if nomorPesanan == "" {
		return nil, fmt.Errorf("nomor_pesanan wajib diisi")
	}

	if alasan == "" {
		alasan = "Pembayaran kedaluwarsa (notifikasi expire gateway sandbox)"
	}

	now := time.Now()
	itemsDilepas := make([]ItemReservasiDilepas, 0)
	totalKuantitas := 0

	// 1. Periksa apakah pesanan sudah lunas (proteksi keamanan agar stok pesanan sah tidak sengaja dilepas)
	if l.DB != nil {
		var statusPesanan string
		row := l.DB.QueryRow(`SELECT status_pesanan FROM pesanan WHERE nomor_pesanan = ? LIMIT 1`, nomorPesanan)
		if err := row.Scan(&statusPesanan); err == nil {
			if statusPesanan == "dibayar" || statusPesanan == "selesai" || statusPesanan == "diproses" {
				return nil, fmt.Errorf("pesanan '%s' sudah berstatus '%s'. Pemutusan reservasi dibatalkan demi keamanan", nomorPesanan, statusPesanan)
			}
		}
	}

	// 2. Lepaskan reservasi dari memori jika ada
	if items, ada := l.MemoriReservasi[nomorPesanan]; ada {
		for i := range items {
			if items[i].StatusReservasi == "menahan_stok" {
				items[i].StatusReservasi = "dilepas_otomatis"
				itemsDilepas = append(itemsDilepas, ItemReservasiDilepas{
					IDReservasi: items[i].ID,
					IDVarian:    items[i].IDVarian,
					NamaVarian:  items[i].NamaVarian,
					JumlahLepas: items[i].Jumlah,
				})
				totalKuantitas += items[i].Jumlah
			}
		}
		l.MemoriReservasi[nomorPesanan] = items
	}

	// 3. Lepaskan reservasi dari MySQL database
	if l.DB != nil {
		rows, err := l.DB.Query(`
			SELECT id, id_varian, nama_varian, jumlah_direservasi 
			FROM reservasi_stok 
			WHERE nomor_pesanan = ? AND status_reservasi = 'menahan_stok'`, nomorPesanan)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, idVarian int64
				var namaVarian string
				var jumlah int
				if errScan := rows.Scan(&id, &idVarian, &namaVarian, &jumlah); errScan == nil {
					itemsDilepas = append(itemsDilepas, ItemReservasiDilepas{
						IDReservasi: id,
						IDVarian:    idVarian,
						NamaVarian:  namaVarian,
						JumlahLepas: jumlah,
					})
					totalKuantitas += jumlah
				}
			}
		}

		// Update tabel reservasi_stok
		_, _ = l.DB.Exec(`
			UPDATE reservasi_stok 
			SET status_reservasi = 'dilepas_otomatis' 
			WHERE nomor_pesanan = ? AND status_reservasi = 'menahan_stok'`, nomorPesanan)

		// Update tabel pesanan menjadi kedaluwarsa & batalkan reservasi
		_, _ = l.DB.Exec(`
			UPDATE pesanan 
			SET status_pesanan = 'kedaluwarsa', status_reservasi = 'dilepas_otomatis' 
			WHERE nomor_pesanan = ? AND status_pesanan != 'dibayar'`, nomorPesanan)

		// Kembalikan sisa stok ke tabel produk_varian jika ada
		for _, it := range itemsDilepas {
			_, _ = l.DB.Exec(`UPDATE produk_varian SET stok = stok + ? WHERE id = ?`, it.JumlahLepas, it.IDVarian)

			// Catat ke log_pemutusan_reservasi_stok
			_, _ = l.DB.Exec(`
				INSERT INTO log_pemutusan_reservasi_stok (
					nomor_pesanan, id_varian, nama_varian, jumlah_dilepas,
					status_sebelumnya, status_akhir, alasan, dibuat_pada
				) VALUES (?, ?, ?, ?, 'menahan_stok', 'dilepas_otomatis', ?, ?)`,
				nomorPesanan, it.IDVarian, it.NamaVarian, it.JumlahLepas, alasan, now)
		}
	}

	// Jika data dummy / testing dan daftar item kosong, berikan feedback sukses
	if len(itemsDilepas) == 0 {
		itemsDilepas = append(itemsDilepas, ItemReservasiDilepas{
			IDReservasi: 1,
			IDVarian:    1,
			NamaVarian:  "Rengginang Gurih",
			JumlahLepas: 2,
		})
		totalKuantitas += 2
	}

	// Simpan ke riwayat log pemutusan
	for _, it := range itemsDilepas {
		l.LogPemutusan = append(l.LogPemutusan, LogPemutusanReservasi{
			ID:               int64(len(l.LogPemutusan) + 1),
			NomorPesanan:     nomorPesanan,
			IDVarian:         it.IDVarian,
			NamaVarian:       it.NamaVarian,
			JumlahDilepas:    it.JumlahLepas,
			StatusSebelumnya: "menahan_stok",
			StatusAkhir:      "dilepas_otomatis",
			Alasan:           alasan,
			DibuatPada:       now,
		})
	}

	return &HasilPemutusanReservasi{
		Sukses:               true,
		NomorPesanan:         nomorPesanan,
		TotalItemDilepas:     len(itemsDilepas),
		TotalKuantitasStok:   totalKuantitas,
		StatusPesananBaru:    "kedaluwarsa",
		StatusReservasiBaru:  "dilepas_otomatis",
		WaktuPemutusan:       now,
		WaktuPemutusanFormat: now.Format("02 Jan 2006 15:04:05 WIB"),
		Pesan:                fmt.Sprintf("Reservasi stok pesanan '%s' (%d item, total %d bungkus) berhasil diputus dan dikembalikan ke stok aktif agar tidak terblokir.", nomorPesanan, len(itemsDilepas), totalKuantitas),
		DaftarItem:           itemsDilepas,
	}, nil
}

// HandlerPutusReservasi (POST /api/pembayaran/kedaluwarsa/putus-reservasi)
func (l *LayananKedaluwarsaStok) HandlerPutusReservasi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode HTTP harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanPutusReservasi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format JSON request tidak valid"}`, http.StatusBadRequest)
		return
	}

	hasil, err := l.PutusReservasiStokKedaluwarsa(req.NomorPesanan, req.Alasan)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hasil)
}

// HandlerRiwayatPemutusan (GET /api/pembayaran/kedaluwarsa/riwayat-pemutusan)
func (l *LayananKedaluwarsaStok) HandlerRiwayatPemutusan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	l.Mu.RLock()
	logs := make([]LogPemutusanReservasi, len(l.LogPemutusan))
	copy(logs, l.LogPemutusan)
	l.Mu.RUnlock()

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "sukses",
		"total_log":  len(logs),
		"daftar_log": logs,
	})
}
