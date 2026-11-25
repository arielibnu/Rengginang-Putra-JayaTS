package signatureidempotency

import "time"

// HasilVerifikasiSignature hasil komparasi signature key SHA512 Midtrans
type HasilVerifikasiSignature struct {
	Valid             bool   `json:"valid"`
	HitunganSHA512    string `json:"hitungan_sha512"`
	SignatureDiterima string `json:"signature_diterima"`
	Pesan             string `json:"pesan"`
}

// CatatanIdempotency rekaman unik untuk mencegah double processing notifikasi webhook
type CatatanIdempotency struct {
	ID                  int64     `json:"id"`
	IdempotencyKey      string    `json:"idempotency_key"`
	NomorPesanan        string    `json:"nomor_pesanan"`
	TransactionID       string    `json:"transaction_id"`
	StatusTransaksi     string    `json:"status_transaksi"`
	SudahDiproses       bool      `json:"sudah_diproses"`
	JumlahPercobaan     int       `json:"jumlah_percobaan"`
	PertamaKaliDiterima time.Time `json:"pertama_kali_diterima"`
	TerakhirDiterima    time.Time `json:"terakhir_diterima"`
}

// HasilCekIdempotency hasil penentuan apakah webhook boleh dieksekusi atau dilewati
type HasilCekIdempotency struct {
	AdalahDuplikat bool               `json:"adalah_duplikat"`
	BolehDiproses  bool               `json:"boleh_diproses"`
	Pesan          string             `json:"pesan"`
	Catatan        CatatanIdempotency `json:"catatan"`
}

// PermintaanVerifikasiKeamanan payload pengujian endpoint verifikasi keamanan
type PermintaanVerifikasiKeamanan struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
}
