package pembaruanstatus

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LayananPembaruanStatus bertanggung jawab memperbarui status transaksi secara otomatis dan akurat
type LayananPembaruanStatus struct {
	DB            *sql.DB
	MemoriStatus  map[string]*HasilPembaruanStatus
	Riwayat       []RiwayatStatusTransaksi
	Mu            sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananPembaruanStatus
func InisialisasiLayanan(db *sql.DB) *LayananPembaruanStatus {
	s := &LayananPembaruanStatus{
		DB:           db,
		MemoriStatus: make(map[string]*HasilPembaruanStatus),
		Riwayat:      make([]RiwayatStatusTransaksi, 0),
	}
	s.InisialisasiTabel()
	return s
}

// InisialisasiTabel memastikan tabel riwayat_status_transaksi tersedia
func (l *LayananPembaruanStatus) InisialisasiTabel() {
	if l.DB == nil {
		return
	}
	query := `
	CREATE TABLE IF NOT EXISTS riwayat_status_transaksi (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		transaction_id VARCHAR(100) NOT NULL,
		status_lama VARCHAR(50) NOT NULL,
		status_baru VARCHAR(50) NOT NULL,
		label_status VARCHAR(50) NOT NULL,
		pemicu VARCHAR(50) NOT NULL DEFAULT 'webhook',
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_riwayat_pesanan (nomor_pesanan),
		INDEX idx_riwayat_trx (transaction_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
	_, _ = l.DB.Exec(query)
}

// PetakanStatusGateway memetakan status payment gateway ke status sistem dan label resmi
func (l *LayananPembaruanStatus) PetakanStatusGateway(trxStatus, fraudStatus string) (kodeStatus string, labelStatus string) {
	trxStatus = strings.ToLower(strings.TrimSpace(trxStatus))
	fraudStatus = strings.ToLower(strings.TrimSpace(fraudStatus))

	switch trxStatus {
	case MidtransSettlement:
		return StatusLunasDibayar, "Lunas/Dibayar"
	case MidtransCapture:
		if fraudStatus == "challenge" {
			return StatusMenungguBayar, "Menunggu"
		}
		return StatusLunasDibayar, "Lunas/Dibayar"
	case MidtransPending:
		return StatusMenungguBayar, "Menunggu"
	case MidtransExpire:
		return StatusKedaluwarsa, "Kedaluwarsa"
	case MidtransCancel, MidtransDeny, MidtransFailure:
		return StatusGagal, "Gagal"
	default:
		return StatusMenungguBayar, "Menunggu"
	}
}

// PerbaruiStatusOtomatis memperbarui status pesanan dan transaksi di database serta memori
func (l *LayananPembaruanStatus) PerbaruiStatusOtomatis(req PembaruanStatusRequest, pemicu string) (*HasilPembaruanStatus, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	req.NomorPesanan = strings.TrimSpace(req.NomorPesanan)
	if req.NomorPesanan == "" {
		return nil, fmt.Errorf("nomor_pesanan wajib diisi")
	}

	statusBaru, labelBaru := l.PetakanStatusGateway(req.TransactionStatus, req.FraudStatus)
	now := time.Now()
	statusLama := StatusMenungguBayar

	if exist, ada := l.MemoriStatus[req.NomorPesanan]; ada {
		statusLama = exist.StatusPesanan
	}

	if l.DB != nil {
		var curStatus sql.NullString
		_ = l.DB.QueryRow(`SELECT status_pesanan FROM pesanan WHERE nomor_pesanan = ? LIMIT 1`, req.NomorPesanan).Scan(&curStatus)
		if curStatus.Valid && curStatus.String != "" {
			statusLama = curStatus.String
		}

		// Update tabel pesanan
		_, _ = l.DB.Exec(`UPDATE pesanan SET status_pesanan = ? WHERE nomor_pesanan = ?`, statusBaru, req.NomorPesanan)

		// Update tabel transaksi_pembayaran_gateway jika ada
		_, _ = l.DB.Exec(`
			UPDATE transaksi_pembayaran_gateway 
			SET transaction_status = ?, diperbarui_pada = ? 
			WHERE nomor_pesanan = ?`, req.TransactionStatus, now, req.NomorPesanan)

		// Catat ke tabel riwayat_status_transaksi
		_, _ = l.DB.Exec(`
			INSERT INTO riwayat_status_transaksi (nomor_pesanan, transaction_id, status_lama, status_baru, label_status, pemicu, dibuat_pada)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			req.NomorPesanan, req.TransactionID, statusLama, statusBaru, labelBaru, pemicu, now)
	}

	hasil := &HasilPembaruanStatus{
		Sukses:           true,
		NomorPesanan:     req.NomorPesanan,
		TransactionID:    req.TransactionID,
		StatusGateway:    req.TransactionStatus,
		StatusSebelumnya: statusLama,
		StatusPesanan:    statusBaru,
		LabelStatus:      labelBaru,
		WaktuPembaruan:   now,
		WaktuFormat:      now.Format("02 Jan 2006 15:04:05 WIB"),
		Pesan:            fmt.Sprintf("Status pesanan '%s' berhasil diperbarui otomatis menjadi '%s' (%s).", req.NomorPesanan, statusBaru, labelBaru),
	}

	l.MemoriStatus[req.NomorPesanan] = hasil

	l.Riwayat = append(l.Riwayat, RiwayatStatusTransaksi{
		ID:            int64(len(l.Riwayat) + 1),
		NomorPesanan:  req.NomorPesanan,
		TransactionID: req.TransactionID,
		StatusLama:    statusLama,
		StatusBaru:    statusBaru,
		LabelStatus:   labelBaru,
		Pemicu:        pemicu,
		DibuatPada:    now,
	})

	return hasil, nil
}

// AmbilStatusTerakhir mengembalikan status pembayaran terkini untuk nomor pesanan
func (l *LayananPembaruanStatus) AmbilStatusTerakhir(nomorPesanan string) (*HasilPembaruanStatus, bool) {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	hasil, ada := l.MemoriStatus[nomorPesanan]
	return hasil, ada
}

// HandlerPerbaruiStatus (POST /api/pembayaran/perbarui-status)
func (l *LayananPembaruanStatus) HandlerPerbaruiStatus(w http.ResponseWriter, r *http.Request) {
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

	var req PembaruanStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format JSON request tidak valid"}`, http.StatusBadRequest)
		return
	}

	hasil, err := l.PerbaruiStatusOtomatis(req, "manual_api")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hasil)
}

// HandlerStatusPesanan (GET /api/pembayaran/status-pesanan?nomor_pesanan=PJ-...)
func (l *LayananPembaruanStatus) HandlerStatusPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	orderID := r.URL.Query().Get("nomor_pesanan")
	if orderID == "" {
		http.Error(w, `{"status":"error","pesan":"Parameter nomor_pesanan wajib diisi"}`, http.StatusBadRequest)
		return
	}

	hasil, ada := l.AmbilStatusTerakhir(orderID)
	if !ada {
		// Fallback default menunggu pembayaran
		hasil = &HasilPembaruanStatus{
			Sukses:        true,
			NomorPesanan:  orderID,
			StatusPesanan: StatusMenungguBayar,
			LabelStatus:   "Menunggu",
			Pesan:         "Pesanan sedang menunggu pembayaran",
		}
	}

	_ = json.NewEncoder(w).Encode(hasil)
}
