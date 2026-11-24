package penjualanoffline

import "time"

// PenjualanOffline merepresentasikan transaksi kasir toko langsung di luar website (WBS 1.3.5.3)
type PenjualanOffline struct {
	ID                int64     `json:"id"`
	NomorTransaksi    string    `json:"nomor_transaksi"`
	TanggalPenjualan  string    `json:"tanggal_penjualan"`
	IDVarian          int64     `json:"id_varian"`
	NamaVarian        string    `json:"nama_varian"`
	Jumlah            int       `json:"jumlah"`
	HargaSatuan       int64     `json:"harga_satuan"`
	TotalHarga        int64     `json:"total_harga"`
	TotalHargaFormat  string    `json:"total_harga_format"`
	MetodePembayaran  string    `json:"metode_pembayaran"` // Tunai
	SumberTransaksi   string    `json:"sumber_transaksi"`  // offline_kasir
	FlagLingkungan    string    `json:"flag_lingkungan"`   // PRODUKSI_NYATA vs SANDBOX_SIMULASI (WBS 1.3.5.4)
	IsSandbox         bool      `json:"is_sandbox"`        // false = Transaksi Nyata
	StokSebelum       int       `json:"stok_sebelum"`
	StokSesudah       int       `json:"stok_sesudah"`
	KasirNama         string    `json:"kasir_nama"`
	Catatan           string    `json:"catatan"`
	DibuatPada        time.Time `json:"dibuat_pada"`
	WaktuFormat       string    `json:"waktu_format"`
}

// PermintaanPenjualanOffline payload request dari frontend kasir offline
type PermintaanPenjualanOffline struct {
	TanggalPenjualan string `json:"tanggal_penjualan"`
	IDVarian         int64  `json:"id_varian"`
	NamaVarian       string `json:"nama_varian"`
	Jumlah           int    `json:"jumlah"`
	HargaSatuan      int64  `json:"harga_satuan"`
	MetodePembayaran string `json:"metode_pembayaran"`
	FlagLingkungan   string `json:"flag_lingkungan,omitempty"` // PRODUKSI_NYATA (default)
	Catatan          string `json:"catatan,omitempty"`
}

// ResponsPenjualanOffline payload balasan setelah penjualan offline berhasil dicatat
type ResponsPenjualanOffline struct {
	Status string            `json:"status"`
	Pesan  string            `json:"pesan"`
	Data   *PenjualanOffline `json:"data,omitempty"`
}

// ResponsDaftarPenjualan payload daftar riwayat kasir toko
type ResponsDaftarPenjualan struct {
	Status     string             `json:"status"`
	Pesan      string             `json:"pesan"`
	Total      int                `json:"total"`
	TotalNilai int64              `json:"total_nilai"`
	Data       []PenjualanOffline `json:"data"`
}
