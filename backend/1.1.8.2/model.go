package katalogapi

import "time"

// ==========================================
// WBS 1.1.8.2
// Model Data Katalog Produk Publik
// Menyajikan data produk dan varian yang aktif
// untuk ditampilkan kepada calon pembeli.
// ==========================================

// ItemKatalog mewakili satu varian produk yang aktif
// dan siap ditampilkan pada katalog website publik.
type ItemKatalog struct {
	IDVarian     uint64    `json:"id_varian"`
	IDProduk     uint64    `json:"id_produk"`
	KodeVarian   string    `json:"kode_varian"`
	NamaProduk   string    `json:"nama_produk"`
	NamaVarian   string    `json:"nama_varian"`
	HargaJual    float64   `json:"harga_jual"`
	HargaFormat  string    `json:"harga_format"`
	BeratGram    uint64    `json:"berat_gram"`
	BeratFormat  string    `json:"berat_format"`
	Stok         uint64    `json:"stok"`
	StatusStok   string    `json:"status_stok"`
	Deskripsi    string    `json:"deskripsi"`
	Kategori     string    `json:"kategori"`
	Gambar       string    `json:"gambar"`
	GambarDetail string    `json:"gambar_detail"`
	StatusAktif  bool      `json:"status_aktif"`
	Diperbarui   time.Time `json:"diperbarui_pada"`
}

// ResponsKatalog adalah format JSON standar
// untuk respons API katalog produk.
type ResponsKatalog struct {
	Pesan      string        `json:"pesan"`
	Total      int           `json:"total"`
	Data       interface{}   `json:"data"`
}
