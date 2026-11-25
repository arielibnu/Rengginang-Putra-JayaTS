package webhookgateway

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LayananWebhook menangani penerimaan webhook notifikasi dari payment gateway
type LayananWebhook struct {
	DB          *sql.DB
	LogWebhook  []LogWebhookEntri
	Mu          sync.RWMutex
	OnNotifikasi func(notif NotifikasiWebhookMidtrans) error
}

// InisialisasiLayanan membuat instance baru LayananWebhook
func InisialisasiLayanan(db *sql.DB) *LayananWebhook {
	s := &LayananWebhook{
		DB:         db,
		LogWebhook: make([]LogWebhookEntri, 0),
	}
	s.InisialisasiTabel()
	return s
}

// InisialisasiTabel memastikan skema tabel log webhook tersedia di database
func (l *LayananWebhook) InisialisasiTabel() {
	if l.DB == nil {
		return
	}
	query := `
	CREATE TABLE IF NOT EXISTS log_webhook_pembayaran (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL,
		transaction_id VARCHAR(100) NOT NULL,
		status_transaksi VARCHAR(50) NOT NULL,
		status_code VARCHAR(10) NOT NULL,
		gross_amount VARCHAR(30) NOT NULL,
		payment_type VARCHAR(50) NULL,
		signature_key VARCHAR(150) NULL,
		raw_payload TEXT NOT NULL,
		status_respon_http INT NOT NULL DEFAULT 200,
		pesan_respon VARCHAR(255) NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_webhook_pesanan (nomor_pesanan),
		INDEX idx_webhook_trx (transaction_id),
		INDEX idx_webhook_status (status_transaksi)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
	_, _ = l.DB.Exec(query)
}

// CatatLogWebhook menyimpan jejak audit notifikasi webhook yang diterima
func (l *LayananWebhook) CatatLogWebhook(notif NotifikasiWebhookMidtrans, raw string, statusHTTP int, pesan string) LogWebhookEntri {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	now := time.Now()
	entri := LogWebhookEntri{
		ID:                  int64(len(l.LogWebhook) + 1),
		NomorPesanan:        notif.OrderID,
		TransactionID:       notif.TransactionID,
		StatusTransaksi:     notif.TransactionStatus,
		StatusCode:          notif.StatusCode,
		GrossAmount:         notif.GrossAmount,
		PaymentType:         notif.PaymentType,
		SignatureKey:        notif.SignatureKey,
		RawPayload:          raw,
		StatusResponHTTP:    statusHTTP,
		PesanRespon:         pesan,
		WaktuDiterima:       now,
		WaktuDiterimaFormat: now.Format("02 Jan 2006 15:04:05 WIB"),
	}

	if l.DB != nil {
		res, err := l.DB.Exec(`
			INSERT INTO log_webhook_pembayaran (
				nomor_pesanan, transaction_id, status_transaksi, status_code,
				gross_amount, payment_type, signature_key, raw_payload,
				status_respon_http, pesan_respon, dibuat_pada
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			notif.OrderID, notif.TransactionID, notif.TransactionStatus, notif.StatusCode,
			notif.GrossAmount, notif.PaymentType, notif.SignatureKey, raw,
			statusHTTP, pesan, now,
		)
		if err == nil {
			id, _ := res.LastInsertId()
			entri.ID = id
		}
	}

	l.LogWebhook = append(l.LogWebhook, entri)
	return entri
}

// AmbilSemuaLog mengambil seluruh jejak audit notifikasi webhook
func (l *LayananWebhook) AmbilSemuaLog() []LogWebhookEntri {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	hasil := make([]LogWebhookEntri, len(l.LogWebhook))
	copy(hasil, l.LogWebhook)
	return hasil
}

// HandlerWebhook menerima callback notifikasi status pembayaran dari payment gateway (POST /api/pembayaran/webhook)
func (l *LayananWebhook) HandlerWebhook(w http.ResponseWriter, r *http.Request) {
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
		_ = json.NewEncoder(w).Encode(ResponWebhookStandar{
			Status:      "error",
			Pesan:       "Metode HTTP yang diizinkan hanya POST untuk endpoint webhook",
			WaktuProses: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bytes.TrimSpace(bodyBytes)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ResponWebhookStandar{
			Status:      "error",
			Pesan:       "Body notifikasi webhook tidak boleh kosong",
			WaktuProses: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	var notif NotifikasiWebhookMidtrans
	if err := json.Unmarshal(bodyBytes, &notif); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ResponWebhookStandar{
			Status:      "error",
			Pesan:       fmt.Sprintf("Format JSON webhook tidak valid: %v", err),
			WaktuProses: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	notif.OrderID = strings.TrimSpace(notif.OrderID)
	if notif.OrderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ResponWebhookStandar{
			Status:      "error",
			Pesan:       "Field order_id wajib disertakan dalam notifikasi gateway",
			WaktuProses: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	pesanSukses := fmt.Sprintf("Notifikasi status '%s' untuk order '%s' berhasil diterima oleh webhook", notif.TransactionStatus, notif.OrderID)
	l.CatatLogWebhook(notif, string(bodyBytes), http.StatusOK, pesanSukses)

	if l.OnNotifikasi != nil {
		if err := l.OnNotifikasi(notif); err != nil {
			log.Printf("[Webhook] Peringatan penanganan callback notifikasi: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ResponWebhookStandar{
		Status:        "sukses",
		Pesan:         pesanSukses,
		NomorPesanan:  notif.OrderID,
		TransactionID: notif.TransactionID,
		WaktuProses:   time.Now().Format("2006-01-02 15:04:05"),
	})
}

// HandlerRiwayatLog (GET /api/pembayaran/webhook/riwayat)
func (l *LayananWebhook) HandlerRiwayatLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	logs := l.AmbilSemuaLog()
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "sukses",
		"total_log":  len(logs),
		"daftar_log": logs,
	})
}
