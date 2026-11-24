package sinkronisasi

import "time"

// ============================================================================
// WBS 1.3.1.4 - SINKRONISASI IDENTITAS TRANSAKSI GATEWAY DENGAN DATABASE
// Deliverable: Sinkronisasi Identitas Transaksi Gateway (transaction_id) & Webhook
// ============================================================================

// MidtransNotification payload HTTP notification callback webhook dari Midtrans
type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"` // "settlement", "pending", "expire", "cancel", "deny"
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	Currency          string `json:"currency,omitempty"`
	MerchantID        string `json:"merchant_id,omitempty"`
}

// HasilSinkronisasiTransaksi hasil sinkronisasi status transaksi untuk dikonsumsi frontend / API
type HasilSinkronisasiTransaksi struct {
	StatusPesanan         string    `json:"status_pesanan"`       // "menunggu_pembayaran", "dibayar", "dibatalkan"
	StatusPesananLabel    string    `json:"status_pesanan_label"` // "Menunggu bayar" / "Sudah dibayar"
	NomorPesanan          string    `json:"nomor_pesanan"`
	NomorInvoice          string    `json:"nomor_invoice"`
	TransactionID         string    `json:"transaction_id"`
	StatusTransaksi       string    `json:"status_transaksi"`     // "settlement", "pending", "expire", "cancel"
	PaymentType           string    `json:"payment_type"`
	Bank                  string    `json:"bank,omitempty"`
	VANumber              string    `json:"va_number,omitempty"`
	QRURL                 string    `json:"qr_url,omitempty"`
	GrossAmount           int64     `json:"gross_amount"`
	GrossAmountFormat     string    `json:"gross_amount_format"`
	TerakhirDisinkronkan  time.Time `json:"terakhir_disinkronkan"`
}

// PermintaanSimulasiBayar payload untuk mensimulasikan status pembayaran saat demo
type PermintaanSimulasiBayar struct {
	NomorPesanan string `json:"nomor_pesanan"`
	Status       string `json:"status"` // "settlement" atau "dibayar"
}
