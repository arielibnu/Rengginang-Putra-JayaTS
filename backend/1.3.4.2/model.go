package statemachine

import "time"

// Konstanta Status Pemrosesan Pesanan (State Machine)
const (
	StatusProsesBaru             = "baru"
	StatusProsesMenungguDiproses = "menunggu_diproses"
	StatusProsesDiproses         = "diproses"
	StatusProsesDikirim          = "dikirim"
	StatusProsesSelesai          = "selesai"
	StatusProsesBatal            = "batal"
)

// Konstanta Status Pembayaran (Terpisah dari Status Pemrosesan)
const (
	StatusBayarMenunggu     = "menunggu_pembayaran"
	StatusBayarDibayar      = "dibayar"
	StatusBayarGagal        = "gagal"
	StatusBayarKedaluwarsa  = "kadaluwarsa"
)

// PermintaanUbahStatus payload permintaan ubah status dari frontend
type PermintaanUbahStatus struct {
	NomorPesanan     string `json:"nomor_pesanan"`
	StatusPemrosesan string `json:"status_pemrosesan"`
	NomorResi        string `json:"nomor_resi,omitempty"`
	Catatan          string `json:"catatan,omitempty"`
}

// ResponsUbahStatus respons balik ke frontend
type ResponsUbahStatus struct {
	Status                  string    `json:"status"` // sukses, error
	Pesan                   string    `json:"pesan"`
	NomorPesanan            string    `json:"nomor_pesanan"`
	StatusPemrosesanLama    string    `json:"status_pemrosesan_lama"`
	StatusPemrosesanBaru    string    `json:"status_pemrosesan_baru"`
	StatusPembayaranTetap   string    `json:"status_pembayaran_tetap"`
	NomorResi               string    `json:"nomor_resi,omitempty"`
	StokDikembalikan        bool      `json:"stok_dikembalikan"`
	TotalBungkusKembali     int       `json:"total_bungkus_kembali,omitempty"`
	WaktuPembaruan          time.Time `json:"waktu_pembaruan"`
}

// RiwayatStatusPemrosesan data audit log transisi status
type RiwayatStatusPemrosesan struct {
	ID           int64     `json:"id"`
	NomorPesanan string    `json:"nomor_pesanan"`
	StatusLama   string    `json:"status_lama"`
	StatusBaru   string    `json:"status_baru"`
	NomorResi    string    `json:"nomor_resi"`
	DiubahOleh   string    `json:"diubah_oleh"`
	Catatan      string    `json:"catatan"`
	DibuatPada   time.Time `json:"dibuat_pada"`
}
