package pesananowner

import "time"

// FilterPesananOwner menampung parameter pencarian dan penyaringan data pesanan.
type FilterPesananOwner struct {
	KataKunci string `json:"q"`       // Nomor pesanan atau nama kontak
	Status    string `json:"status"`  // Semua status, Baru, Diproses, Selesai, Batal
	Rentang   string `json:"rentang"` // Semua tanggal, Hari ini, 7 hari terakhir, 30 hari terakhir
	Halaman   int    `json:"halaman"` // Pagination
	Batas     int    `json:"batas"`
}

// PesananRingkasOwner data baris tabel antrean pesanan pada frame PPL21_01.
type PesananRingkasOwner struct {
	ID               int64     `json:"id"`
	NomorPesanan     string    `json:"nomor_pesanan"`
	Kontak           string    `json:"kontak"`
	WhatsApp         string    `json:"whatsapp"`
	Metode           string    `json:"metode"`     // Ambil Sendiri / Kirim Kurir
	Pembayaran       string    `json:"pembayaran"` // Menunggu / Dibayar / Gagal
	Status           string    `json:"status"`     // Baru / Diproses / Selesai / Batal
	Total            int       `json:"total"`
	TotalFormat      string    `json:"total_format"`
	ItemsRingkas     string    `json:"items"`
	Tanggal          string    `json:"tanggal"`
	WaktuDibuat      time.Time `json:"waktu_dibuat"`
}

// ItemBarangSnapshot representasi snapshot produk dari pesanan pada frame PPL21_02.
type ItemBarangSnapshot struct {
	IDVarian          int64  `json:"id_varian"`
	NamaProduk        string `json:"nama_produk"`
	NamaVarian        string `json:"nama_varian"`
	JumlahBungkus     int    `json:"jumlah"`
	HargaSatuan       int    `json:"harga_satuan"`
	Subtotal          int    `json:"subtotal"`
	HargaSatuanFormat string `json:"harga_satuan_format"`
	SubtotalFormat    string `json:"subtotal_format"`
}

// DetailPesananOwner informasi pembeli, ringkasan, dan status pesanan pada frame PPL21_02.
type DetailPesananOwner struct {
	NomorPesanan     string               `json:"nomor_pesanan"`
	NomorInvoice     string               `json:"nomor_invoice"`
	Nama             string               `json:"nama"`
	WhatsApp         string               `json:"whatsapp"`
	Alamat           string               `json:"alamat"`
	Pengiriman       string               `json:"pengiriman"`
	Metode           string               `json:"metode"`
	Subtotal         int                  `json:"subtotal"`
	SubtotalFormat   string               `json:"subtotal_format"`
	Ongkir           int                  `json:"ongkir"`
	OngkirFormat     string               `json:"ongkir_format"`
	Total            int                  `json:"total"`
	TotalFormat      string               `json:"total_format"`
	Pembayaran       string               `json:"pembayaran"` // Status pembayaran (Dibayar / Menunggu)
	Pemrosesan       string               `json:"pemrosesan"` // Status pemrosesan (Menunggu Diproses / Diproses / Selesai / Batal)
	Items            []ItemBarangSnapshot `json:"items_detail"`
	ItemsRingkasan   string               `json:"items"`
	Tanggal          string               `json:"tanggal"`
}

// PermintaanUpdateStatus payload untuk mengubah status pemrosesan pesanan.
type PermintaanUpdateStatus struct {
	NomorPesanan     string `json:"nomor_pesanan"`
	StatusPemrosesan string `json:"status_pemrosesan"`
}

// ResponsDaftarPesanan response API daftar pesanan.
type ResponsDaftarPesanan struct {
	Status string                `json:"status"`
	Pesan  string                `json:"pesan,omitempty"`
	Data   []PesananRingkasOwner `json:"data"`
	Total  int                   `json:"total"`
}

// ResponsDetailPesanan response API detail pesanan.
type ResponsDetailPesanan struct {
	Status string              `json:"status"`
	Pesan  string              `json:"pesan,omitempty"`
	Data   *DetailPesananOwner `json:"data"`
}
