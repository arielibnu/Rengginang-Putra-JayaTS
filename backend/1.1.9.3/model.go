package ketersediaanproduk

// DetailProdukResponse data detail spesifikasi produk untuk publik (PPL-09).
type DetailProdukResponse struct {
	IDVarian           uint64       `json:"id_varian"`
	Slug               string       `json:"slug"`
	NamaProduk         string       `json:"nama_produk"`
	NamaVarian         string       `json:"nama_varian"`
	HargaJual          uint64       `json:"harga_jual"`
	FormatHarga        string       `json:"format_harga"`
	BeratGram          uint32       `json:"berat_gram"`
	FormatBerat        string       `json:"format_berat"`
	Deskripsi          string       `json:"deskripsi"`
	Kategori           string       `json:"kategori"`
	StokTersedia       int64        `json:"stok_tersedia"`
	FormatStok         string       `json:"format_stok"`
	StatusKetersediaan string       `json:"status_ketersediaan"` // "TERSEDIA" atau "STOK HABIS"
	BisaDibeli         bool         `json:"bisa_dibeli"`
	PesanStatus        string       `json:"pesan_status"`
	PesanKeterangan    string       `json:"pesan_keterangan"`
	GambarUtama        string       `json:"gambar_utama"`
	GaleriFoto         []FotoGaleri `json:"galeri_foto"`
}

// FotoGaleri menyajikan daftar foto kemasan, label, dan detail isi.
type FotoGaleri struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// PermintaanValidasiBeli payload untuk memverifikasi apakah varian bisa dibeli.
type PermintaanValidasiBeli struct {
	IDVarian   uint64 `json:"id_varian"`
	Slug       string `json:"slug,omitempty"`
	JumlahBeli uint32 `json:"jumlah_beli"`
}

// HasilValidasiBeli hasil keputusan apakah tombol beli diaktifkan atau dinonaktifkan (disabled).
type HasilValidasiBeli struct {
	IDVarian     uint64 `json:"id_varian"`
	Slug         string `json:"slug"`
	BisaDibeli   bool   `json:"bisa_dibeli"`
	StatusStok   string `json:"status_stok"`
	SisaStok     int64  `json:"sisa_stok"`
	PesanValidasi string `json:"pesan_validasi"`
}

// ResponsAPI struktur standar respon JSON backend.
type ResponsAPI struct {
	Status string      `json:"status"`
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}
