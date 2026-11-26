package pengembalianstok

import "time"

// ItemPengembalianStok rincian varian produk yang dikembalikan stoknya ke persediaan aktif
type ItemPengembalianStok struct {
	IDVarian     int64  `json:"id_varian"`
	NamaVarian   string `json:"nama_varian"`
	JumlahBungkus int   `json:"jumlah_bungkus"`
	StokSebelum  int    `json:"stok_sebelum"`
	StokSesudah  int    `json:"stok_sesudah"`
}

// HasilPengembalianStok hasil eksekusi pengembalian stok saat pembatalan
type HasilPengembalianStok struct {
	NomorPesanan      string                 `json:"nomor_pesanan"`
	Status            string                 `json:"status"` // sukses, sudah_dikembalikan, dilewati
	Pesan             string                 `json:"pesan"`
	TotalBungkusKembali int                  `json:"total_bungkus_kembali"`
	DaftarItem        []ItemPengembalianStok `json:"daftar_item"`
	WaktuEksekusi     time.Time              `json:"waktu_eksekusi"`
}

// LogPengembalianStok entri tabel audit pengembalian stok pembatalan
type LogPengembalianStok struct {
	ID           int64     `json:"id"`
	NomorPesanan string    `json:"nomor_pesanan"`
	IDVarian     int64     `json:"id_varian"`
	NamaVarian   string    `json:"nama_varian"`
	Jumlah       int       `json:"jumlah"`
	Alasan       string    `json:"alasan"`
	DibuatPada   time.Time `json:"dibuat_pada"`
}
