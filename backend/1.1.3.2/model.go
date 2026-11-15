package produkapi

import "time"

// Struktur data produk
type Produk struct {
	IDProduk       uint64    `json:"id_produk"`
	NamaProduk     string    `json:"nama_produk"`
	Deskripsi      string    `json:"deskripsi"`
	StatusAktif    bool      `json:"status_aktif"`
	DibuatPada     time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

// Struktur data varian produk
type VarianProduk struct {
	IDVarian       uint64    `json:"id_varian"`
	IDProduk       uint64    `json:"id_produk"`
	NamaVarian     string    `json:"nama_varian"`
	HargaJual      float64   `json:"harga_jual"`
	BeratGram      uint64    `json:"berat_gram"`
	Stok           uint64    `json:"stok"`
	StatusAktif    bool      `json:"status_aktif"`
	DibuatPada     time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

// Format respons API
type ResponsAPI struct {
	Pesan string      `json:"pesan"`
	Data  interface{} `json:"data,omitempty"`
}
