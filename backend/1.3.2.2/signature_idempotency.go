package signatureidempotency

import (
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultServerKey = "SB-Mid-server-a8va2L_17AjLE.77S~VbJtAN"
)

// LayananKeamananIdempotency mengelola validasi signature SHA512 dan proteksi pemrosesan ganda
type LayananKeamananIdempotency struct {
	DB                  *sql.DB
	ServerKey           string
	MemoriIdempotency   map[string]*CatatanIdempotency
	Mu                  sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananKeamananIdempotency
func InisialisasiLayanan(db *sql.DB) *LayananKeamananIdempotency {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = DefaultServerKey
	}

	s := &LayananKeamananIdempotency{
		DB:                db,
		ServerKey:         serverKey,
		MemoriIdempotency: make(map[string]*CatatanIdempotency),
	}
	s.InisialisasiTabel()
	return s
}

// InisialisasiTabel membuat tabel idempotensi_webhook_pembayaran jika belum ada
func (l *LayananKeamananIdempotency) InisialisasiTabel() {
	if l.DB == nil {
		return
	}
	query := `
	CREATE TABLE IF NOT EXISTS idempotensi_webhook_pembayaran (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		idempotency_key VARCHAR(150) NOT NULL UNIQUE,
		nomor_pesanan VARCHAR(64) NOT NULL,
		transaction_id VARCHAR(100) NOT NULL,
		status_transaksi VARCHAR(50) NOT NULL,
		sudah_diproses TINYINT(1) NOT NULL DEFAULT 0,
		jumlah_percobaan INT NOT NULL DEFAULT 1,
		pertama_kali_diterima DATETIME DEFAULT CURRENT_TIMESTAMP,
		terakhir_diterima DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_idempotency_pesanan (nomor_pesanan),
		INDEX idx_idempotency_trx (transaction_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
	_, _ = l.DB.Exec(query)
}

// HitungSignatureKey menghitung SHA512(order_id + status_code + gross_amount + server_key)
func (l *LayananKeamananIdempotency) HitungSignatureKey(orderID, statusCode, grossAmount string) string {
	raw := strings.TrimSpace(orderID) + strings.TrimSpace(statusCode) + strings.TrimSpace(grossAmount) + l.ServerKey
	hash := sha512.Sum512([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// VerifikasiSignatureKey membandingkan signature key yang diterima dengan hash SHA512 standar Midtrans
func (l *LayananKeamananIdempotency) VerifikasiSignatureKey(orderID, statusCode, grossAmount, signatureKey string) HasilVerifikasiSignature {
	orderID = strings.TrimSpace(orderID)
	statusCode = strings.TrimSpace(statusCode)
	grossAmount = strings.TrimSpace(grossAmount)
	signatureKey = strings.TrimSpace(signatureKey)

	hitung := l.HitungSignatureKey(orderID, statusCode, grossAmount)

	// Pembandingan konstan waktu (constant-time compare) untuk mencegah timing attack
	cocok := subtle.ConstantTimeCompare([]byte(strings.ToLower(hitung)), []byte(strings.ToLower(signatureKey))) == 1

	pesan := "Signature key Midtrans valid dan terverifikasi aman"
	if !cocok {
		pesan = "Signature key tidak cocok! Kemungkinan manipulasi data atau unauthorized webhook"
	}

	return HasilVerifikasiSignature{
		Valid:             cocok,
		HitunganSHA512:    hitung,
		SignatureDiterima: signatureKey,
		Pesan:             pesan,
	}
}

// PeriksaDanKunciIdempotency memvalidasi apakah notifikasi transaksi ini sudah pernah diproses sebelumnya
func (l *LayananKeamananIdempotency) PeriksaDanKunciIdempotency(orderID, trxID, statusTransaksi string) (*HasilCekIdempotency, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	orderID = strings.TrimSpace(orderID)
	trxID = strings.TrimSpace(trxID)
	statusTransaksi = strings.ToLower(strings.TrimSpace(statusTransaksi))

	key := fmt.Sprintf("%s:%s", orderID, statusTransaksi)
	now := time.Now()

	// 1. Cek dari memori cache
	if cat, ada := l.MemoriIdempotency[key]; ada {
		cat.JumlahPercobaan++
		cat.TerakhirDiterima = now

		// Update database jika terhubung
		if l.DB != nil {
			_, _ = l.DB.Exec(`
				UPDATE idempotensi_webhook_pembayaran 
				SET jumlah_percobaan = ?, terakhir_diterima = ? 
				WHERE idempotency_key = ?`, cat.JumlahPercobaan, now, key)
		}

		if cat.SudahDiproses {
			return &HasilCekIdempotency{
				AdalahDuplikat: true,
				BolehDiproses:  false,
				Pesan:          fmt.Sprintf("Notifikasi pembayaran status '%s' untuk order '%s' sudah pernah diproses tuntas. Double-processing dicegah (Idempotent).", statusTransaksi, orderID),
				Catatan:        *cat,
			}, nil
		}
	}

	// 2. Cek dari database MySQL jika belum ada di memori
	if l.DB != nil {
		var dbID int64
		var dbSudah int
		var dbPercobaan int
		row := l.DB.QueryRow(`
			SELECT id, sudah_diproses, jumlah_percobaan 
			FROM idempotensi_webhook_pembayaran 
			WHERE idempotency_key = ? LIMIT 1`, key)
		if err := row.Scan(&dbID, &dbSudah, &dbPercobaan); err == nil {
			dbPercobaan++
			_, _ = l.DB.Exec(`
				UPDATE idempotensi_webhook_pembayaran 
				SET jumlah_percobaan = ?, terakhir_diterima = ? 
				WHERE id = ?`, dbPercobaan, now, dbID)

			catatan := &CatatanIdempotency{
				ID:                  dbID,
				IdempotencyKey:      key,
				NomorPesanan:        orderID,
				TransactionID:       trxID,
				StatusTransaksi:     statusTransaksi,
				SudahDiproses:       dbSudah == 1,
				JumlahPercobaan:     dbPercobaan,
				PertamaKaliDiterima: now,
				TerakhirDiterima:    now,
			}
			l.MemoriIdempotency[key] = catatan

			if dbSudah == 1 {
				return &HasilCekIdempotency{
					AdalahDuplikat: true,
					BolehDiproses:  false,
					Pesan:          fmt.Sprintf("Notifikasi pembayaran status '%s' untuk order '%s' sudah tersimpan di database. Double-processing dicegah (Idempotent).", statusTransaksi, orderID),
					Catatan:        *catatan,
				}, nil
			}
		}
	}

	// 3. Catatan baru (belum pernah diproses)
	catatanBaru := &CatatanIdempotency{
		ID:                  int64(len(l.MemoriIdempotency) + 1),
		IdempotencyKey:      key,
		NomorPesanan:        orderID,
		TransactionID:       trxID,
		StatusTransaksi:     statusTransaksi,
		SudahDiproses:       false,
		JumlahPercobaan:     1,
		PertamaKaliDiterima: now,
		TerakhirDiterima:    now,
	}

	if l.DB != nil {
		res, err := l.DB.Exec(`
			INSERT INTO idempotensi_webhook_pembayaran (
				idempotency_key, nomor_pesanan, transaction_id, status_transaksi,
				sudah_diproses, jumlah_percobaan, pertama_kali_diterima, terakhir_diterima
			) VALUES (?, ?, ?, ?, 0, 1, ?, ?)
			ON DUPLICATE KEY UPDATE jumlah_percobaan = jumlah_percobaan + 1, terakhir_diterima = VALUES(terakhir_diterima)`,
			key, orderID, trxID, statusTransaksi, now, now,
		)
		if err == nil {
			if id, errID := res.LastInsertId(); errID == nil && id > 0 {
				catatanBaru.ID = id
			}
		}
	}

	l.MemoriIdempotency[key] = catatanBaru

	return &HasilCekIdempotency{
		AdalahDuplikat: false,
		BolehDiproses:  true,
		Pesan:          "Notifikasi pembayaran valid untuk diproses pertama kali.",
		Catatan:        *catatanBaru,
	}, nil
}

// TandaiSelesaiIdempotency menandai bahwa pembaruan status transaksi telah berhasil dieksekusi
func (l *LayananKeamananIdempotency) TandaiSelesaiIdempotency(orderID, statusTransaksi string) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	key := fmt.Sprintf("%s:%s", strings.TrimSpace(orderID), strings.ToLower(strings.TrimSpace(statusTransaksi)))
	if cat, ada := l.MemoriIdempotency[key]; ada {
		cat.SudahDiproses = true
	}

	if l.DB != nil {
		_, _ = l.DB.Exec(`
			UPDATE idempotensi_webhook_pembayaran 
			SET sudah_diproses = 1, terakhir_diterima = NOW() 
			WHERE idempotency_key = ?`, key)
	}
}

// HandlerVerifikasiKeamanan (POST /api/pembayaran/verifikasi-keamanan)
func (l *LayananKeamananIdempotency) HandlerVerifikasiKeamanan(w http.ResponseWriter, r *http.Request) {
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

	var req PermintaanVerifikasiKeamanan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format JSON request tidak valid"}`, http.StatusBadRequest)
		return
	}

	// 1. Verifikasi Signature Key
	hasilSig := l.VerifikasiSignatureKey(req.OrderID, req.StatusCode, req.GrossAmount, req.SignatureKey)

	// 2. Cek Idempotency
	hasilIdem, _ := l.PeriksaDanKunciIdempotency(req.OrderID, req.TransactionID, req.TransactionStatus)

	statusHTTP := http.StatusOK
	if !hasilSig.Valid {
		statusHTTP = http.StatusUnauthorized
	}

	w.WriteHeader(statusHTTP)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "sukses",
		"verifikasi_sig":  hasilSig,
		"cek_idempotency": hasilIdem,
	})
}

// HandlerRiwayatIdempotency (GET /api/pembayaran/idempotency/riwayat)
func (l *LayananKeamananIdempotency) HandlerRiwayatIdempotency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	l.Mu.RLock()
	daftar := make([]CatatanIdempotency, 0, len(l.MemoriIdempotency))
	for _, v := range l.MemoriIdempotency {
		daftar = append(daftar, *v)
	}
	l.Mu.RUnlock()

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"total":  len(daftar),
		"data":   daftar,
	})
}
