package statemanager

import (
	"time"
)

// ============================================================================
// WBS 1.2.2.2 - PENYEDIAAN PENYIMPANAN KERANJANG BAGI GUEST & PELANGGAN LOGIN
// Deliverable: State Manager Keranjang Belanja
// ============================================================================

// ItemKeranjang merepresentasikan satu jenis varian produk yang ada di keranjang.
type ItemKeranjang struct {
	IDVarian     uint64 `json:"id_varian"`
	Slug         string `json:"slug,omitempty"`
	NamaVarian   string `json:"nama_varian"`
	NamaProduk   string `json:"nama_produk"`
	BeratBersih  string `json:"berat_bersih,omitempty"`
	HargaSatuan  uint64 `json:"harga_satuan"`
	Jumlah       uint32 `json:"jumlah"`
	StokTersedia int64  `json:"stok_tersedia"`
	Subtotal     uint64 `json:"subtotal"`
	Gambar       string `json:"gambar,omitempty"`
}

// PermintaanSinkronisasi payload request dari frontend untuk sinkronisasi state keranjang.
type PermintaanSinkronisasi struct {
	IsGuest     bool            `json:"is_guest"`
	IDPelanggan *uint64         `json:"id_pelanggan,omitempty"`
	Items       []ItemKeranjang `json:"items"`
}

// ResponKeranjang format output keranjang yang tersinkronisasi.
type ResponKeranjang struct {
	IDKeranjang     string          `json:"id_keranjang"`
	IsGuest         bool            `json:"is_guest"`
	IDPelanggan     *uint64         `json:"id_pelanggan,omitempty"`
	TotalJenis      int             `json:"total_jenis_produk"`
	TotalKuantitas  uint32          `json:"total_kuantitas"`
	Subtotal        uint64          `json:"subtotal"`
	FormatSubtotal  string          `json:"format_subtotal"`
	Items           []ItemKeranjang `json:"items"`
	Pesan           string          `json:"pesan"`
	TerakhirDiubah  time.Time       `json:"terakhir_diubah"`
}

// ResponsAPI format standar respon JSON server.
type ResponsAPI struct {
	Status string      `json:"status"` // "sukses", "gagal", "ditolak"
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}
