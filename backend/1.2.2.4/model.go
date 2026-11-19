package validasikeranjang

import (
	"time"
)

// ============================================================================
// WBS 1.2.2.4 - VALIDASI ULANG ISI ITEM KERANJANG TERHADAP SISA STOK
// Deliverable: API Validasi Keranjang
// ============================================================================

// ItemValidasiRequest item yang dikirim frontend untuk divalidasi stoknya.
type ItemValidasiRequest struct {
	IDVarian   uint64 `json:"id_varian"`
	Slug       string `json:"slug,omitempty"`
	NamaVarian string `json:"nama_varian"`
	Jumlah     uint32 `json:"jumlah"`
}

// PermintaanValidasiKeranjang payload request validasi keranjang belanja.
type PermintaanValidasiKeranjang struct {
	Items []ItemValidasiRequest `json:"items"`
}

// ItemValidasiResult hasil pengecekan stok aktual per item.
type ItemValidasiResult struct {
	IDVarian       uint64 `json:"id_varian"`
	Slug           string `json:"slug"`
	NamaVarian     string `json:"nama_varian"`
	JumlahDiminta  uint32 `json:"jumlah_diminta"`
	StokTersedia   int64  `json:"stok_tersedia"`
	HargaSatuan    uint64 `json:"harga_satuan"`
	Subtotal       uint64 `json:"subtotal"`
	StatusItem     string `json:"status_item"` // "VALID", "MELEBIHI_STOK", "STOK_HABIS", "TIDAK_DITEMUKAN"
	Pesan          string `json:"pesan"`
	Valid          bool   `json:"valid"`
}

// HasilValidasiKeranjang rangkuman hasil validasi seluruh isi keranjang belanja.
type HasilValidasiKeranjang struct {
	BisaCheckout     bool                 `json:"bisa_checkout"`
	TotalJenisProduk int                  `json:"total_jenis_produk"`
	TotalKuantitas   uint32               `json:"total_kuantitas"`
	Subtotal         uint64               `json:"subtotal"`
	FormatSubtotal   string               `json:"format_subtotal"`
	StatusValidasi   string               `json:"status_validasi"` // "VALIDASI_SUKSES", "VALIDASI_DITOLAK"
	PesanUtama       string               `json:"pesan_utama"`
	RincianItem      []ItemValidasiResult `json:"rincian_item"`
	DivalidasiPada   time.Time            `json:"divalidasi_pada"`
}

// ResponsAPI format standar respon JSON server.
type ResponsAPI struct {
	Status string      `json:"status"` // "sukses", "ditolak", "gagal"
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}
