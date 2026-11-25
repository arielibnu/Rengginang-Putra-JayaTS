package webhookgateway

import "time"

// NotifikasiWebhookMidtrans payload JSON standar webhook notifikasi dari Midtrans Sandbox
type NotifikasiWebhookMidtrans struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"` // settlement, capture, pending, expire, cancel, deny
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"` // 200, 201, 202, 407, etc.
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"` // qris, bank_transfer, gopay, etc.
	OrderID           string `json:"order_id"`     // PJ-20261122-XXXX
	MerchantID        string `json:"merchant_id,omitempty"`
	GrossAmount       string `json:"gross_amount"` // contoh: "56000.00"
	FraudStatus       string `json:"fraud_status,omitempty"` // accept, challenge
	Currency          string `json:"currency,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

// LogWebhookEntri struktur pencatatan riwayat webhook yang masuk ke sistem
type LogWebhookEntri struct {
	ID                int64     `json:"id"`
	NomorPesanan      string    `json:"nomor_pesanan"`
	TransactionID     string    `json:"transaction_id"`
	StatusTransaksi   string    `json:"status_transaksi"`
	StatusCode        string    `json:"status_code"`
	GrossAmount       string    `json:"gross_amount"`
	PaymentType       string    `json:"payment_type"`
	SignatureKey      string    `json:"signature_key"`
	RawPayload        string    `json:"raw_payload"`
	StatusResponHTTP  int       `json:"status_respon_http"`
	PesanRespon       string    `json:"pesan_respon"`
	WaktuDiterima     time.Time `json:"waktu_diterima"`
	WaktuDiterimaFormat string  `json:"waktu_diterima_format"`
}

// ResponWebhookStandar format respon balik ke server Payment Gateway Midtrans
type ResponWebhookStandar struct {
	Status        string `json:"status"` // sukses, error
	Pesan         string `json:"pesan"`
	NomorPesanan  string `json:"nomor_pesanan,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	WaktuProses   string `json:"waktu_proses"`
}
