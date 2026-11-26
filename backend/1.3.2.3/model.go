package pembaruanstatus

import "time"

// Konstanta Status Transaksi Midtrans
const (
	MidtransSettlement = "settlement"
	MidtransCapture    = "capture"
	MidtransPending    = "pending"
	MidtransExpire     = "expire"
	MidtransCancel     = "cancel"
	MidtransDeny       = "deny"
	MidtransFailure    = "failure"
)

// Konstanta Status Pesanan Sistem Internal
const (
	StatusMenungguBayar = "menunggu_pembayaran" // Label: Menunggu
	StatusLunasDibayar  = "dibayar"             // Label: Lunas/Dibayar
	StatusKedaluwarsa   = "kedaluwarsa"         // Label: Kedaluwarsa
	StatusGagal         = "gagal"               // Label: Gagal
)

// PembaruanStatusRequest parameter input untuk memperbarui status pesanan
type PembaruanStatusRequest struct {
	NomorPesanan      string `json:"nomor_pesanan"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	PaymentType       string `json:"payment_type,omitempty"`
	GrossAmount       int64  `json:"gross_amount,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

// HasilPembaruanStatus hasil eksekusi pembaruan status
type HasilPembaruanStatus struct {
	Sukses           bool      `json:"sukses"`
	NomorPesanan     string    `json:"nomor_pesanan"`
	TransactionID    string    `json:"transaction_id"`
	StatusGateway    string    `json:"status_gateway"`
	StatusSebelumnya string    `json:"status_sebelumnya"`
	StatusPesanan    string    `json:"status_pesanan"`       // kode internal (dibayar, menunggu_pembayaran, kedaluwarsa, gagal)
	LabelStatus      string    `json:"label_status"`        // Menunggu, Lunas/Dibayar, Kedaluwarsa, Gagal
	WaktuPembaruan   time.Time `json:"waktu_pembaruan"`
	WaktuFormat      string    `json:"waktu_format"`
	Pesan            string    `json:"pesan"`
}

// RiwayatStatusTransaksi rekaman audit perubahan status pembayaran
type RiwayatStatusTransaksi struct {
	ID            int64     `json:"id"`
	NomorPesanan  string    `json:"nomor_pesanan"`
	TransactionID string    `json:"transaction_id"`
	StatusLama    string    `json:"status_lama"`
	StatusBaru    string    `json:"status_baru"`
	LabelStatus   string    `json:"label_status"`
	Pemicu        string    `json:"pemicu"` // webhook, manual_sync, simulator
	DibuatPada    time.Time `json:"dibuat_pada"`
}
