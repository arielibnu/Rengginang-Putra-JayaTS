package tagihan

import "time"

// ============================================================================
// WBS 1.3.1.2 - INTEGRASI PEMBUATAN TAGIHAN PAYMENT GATEWAY SANDBOX
// Deliverable: Integrasi Pembuatan Tagihan Midtrans (QRIS & Virtual Account)
// ============================================================================

// PermintaanBuatTagihan payload dari frontend saat memilih metode pembayaran
type PermintaanBuatTagihan struct {
	NomorPesanan string `json:"nomor_pesanan"`
	Metode       string `json:"metode"` // "qris" atau "va" / "bank_transfer"
	Bank         string `json:"bank,omitempty"`   // "bca", "bni", "bri", "permata", "mandiri"
}

// ResponsTagihanGateway respons resmi tagihan pembayaran untuk frontend
type ResponsTagihanGateway struct {
	Status            string    `json:"status"` // "sukses" atau "error"
	Pesan             string    `json:"pesan"`
	NomorPesanan      string    `json:"nomor_pesanan"`
	TransactionID     string    `json:"transaction_id"`
	PaymentType       string    `json:"payment_type"` // "qris" atau "bank_transfer"
	GrossAmount       int64     `json:"gross_amount"`
	GrossAmountFormat string    `json:"gross_amount_format"`
	Bank              string    `json:"bank,omitempty"`
	VANumber          string    `json:"va_number,omitempty"`
	QRURL             string    `json:"qr_url,omitempty"`
	QRString          string    `json:"qr_string,omitempty"`
	SnapToken         string    `json:"snap_token,omitempty"`
	SnapRedirectURL   string    `json:"snap_redirect_url,omitempty"`
	TransactionStatus string    `json:"transaction_status"` // "pending"
	SimulatorURL      string    `json:"simulator_url,omitempty"`
	ExpiryTime        string    `json:"expiry_time,omitempty"`
	DibuatPada        time.Time `json:"dibuat_pada"`
}

// EntitasTransaksiPembayaran record tabel transaksi pembayaran gateway di database
type EntitasTransaksiPembayaran struct {
	ID                int64     `json:"id"`
	IDPesanan         int64     `json:"id_pesanan"`
	NomorPesanan      string    `json:"nomor_pesanan"`
	NomorInvoice      string    `json:"nomor_invoice"`
	TransactionID     string    `json:"transaction_id"`
	PaymentType       string    `json:"payment_type"`
	Bank              string    `json:"bank"`
	VANumber          string    `json:"va_number"`
	QRURL             string    `json:"qr_url"`
	QRString          string    `json:"qr_string"`
	GrossAmount       int64     `json:"gross_amount"`
	TransactionStatus string    `json:"transaction_status"`
	FraudStatus       string    `json:"fraud_status"`
	RawGatewayResponse string   `json:"raw_gateway_response"`
	DibuatPada        time.Time `json:"dibuat_pada"`
	DiperbaruiPada    time.Time `json:"diperbarui_pada"`
}

// MidtransChargeRequest payload request ke Midtrans Core API /v2/charge
type MidtransChargeRequest struct {
	PaymentType        string                   `json:"payment_type"`
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	CustomerDetails    *MidtransCustomerDetails `json:"customer_details,omitempty"`
	ItemDetails        []MidtransItemDetails    `json:"item_details,omitempty"`
	BankTransfer       *MidtransBankTransfer    `json:"bank_transfer,omitempty"`
	QRIS               *MidtransQRIS            `json:"qris,omitempty"`
}

type MidtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name"`
	Phone     string `json:"phone,omitempty"`
}

type MidtransItemDetails struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type MidtransBankTransfer struct {
	Bank string `json:"bank"`
}

type MidtransQRIS struct {
	Acquirer string `json:"acquirer"`
}

// MidtransAction tindakan untuk QRIS Midtrans
type MidtransAction struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// MidtransVANumber rincian nomor VA dari Midtrans
type MidtransVANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

// MidtransChargeResponse respons JSON dari Midtrans Core API /v2/charge
type MidtransChargeResponse struct {
	StatusCode        string             `json:"status_code"`
	StatusMessage     string             `json:"status_message"`
	TransactionID     string             `json:"transaction_id"`
	OrderID           string             `json:"order_id"`
	GrossAmount       string             `json:"gross_amount"`
	Currency          string             `json:"currency"`
	PaymentType       string             `json:"payment_type"`
	TransactionTime   string             `json:"transaction_time"`
	TransactionStatus string             `json:"transaction_status"`
	FraudStatus       string             `json:"fraud_status"`
	Actions           []MidtransAction   `json:"actions,omitempty"`
	QRString          string             `json:"qr_string,omitempty"`
	VANumbers         []MidtransVANumber `json:"va_numbers,omitempty"`
	PermataVANumber   string             `json:"permata_va_number,omitempty"`
	BillKey           string             `json:"bill_key,omitempty"`
	BillerCode        string             `json:"biller_code,omitempty"`
	ExpiryTime        string             `json:"expiry_time,omitempty"`
}

// MidtransSnapRequest payload permintaan pembuatan sesi transaksi ke Midtrans Snap API
type MidtransSnapRequest struct {
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	CustomerDetails    *MidtransCustomerDetails   `json:"customer_details,omitempty"`
	ItemDetails        []MidtransItemDetails      `json:"item_details,omitempty"`
	EnabledPayments    []string                   `json:"enabled_payments,omitempty"`
}

// MidtransSnapResponse respons resmi dari Midtrans Snap API
type MidtransSnapResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}
