package validatororder

import "time"

// ============================================================================
// WBS 1.2.8.3 - VALIDASI KETAT INTEGRITAS SUBTOTAL, ONGKIR, TOTAL NOMINAL,
// DAN KELENGKAPAN DATA DI SISI SERVER (Server-Side Order Validator)
// ============================================================================

// ItemOrder merepresentasikan rincian produk pesanan
type ItemOrder struct {
	IDProduk    string `json:"id_produk"`
	NamaProduk  string `json:"nama_produk"`
	Jumlah      int    `json:"jumlah"`
	HargaSatuan int    `json:"harga_satuan"`
	Subtotal    int    `json:"subtotal"`
}

// KontakPemesan data pemesan penerima barang
type KontakPemesan struct {
	Nama          string `json:"nama"`
	NomorWhatsApp string `json:"nomor_whatsapp"`
	Email         string `json:"email,omitempty"`
}

// AlamatPengiriman wilayah dan alamat penerima barang
type AlamatPengiriman struct {
	Provinsi      string `json:"provinsi"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Kecamatan     string `json:"kecamatan"`
	AlamatLengkap string `json:"alamat_lengkap"`
	KodePos       string `json:"kode_pos,omitempty"`
}

// KurirPilihan kurir dan tarif yang telah dipilih pada tahap sebelumnya (PPL-16)
type KurirPilihan struct {
	IDKurir          string `json:"id_kurir"`           // "jne", "jnt", "sicepat"
	NamaKurir        string `json:"nama_kurir"`         // "JNE REG", "J&T Express", "SiCepat"
	KodeLayanan      string `json:"kode_layanan"`       // "REG", "EZ"
	TarifOngkir      int    `json:"tarif_ongkir"`       // Tarif dalam rupiah
	EstimasiHari     string `json:"estimasi_hari"`      // misal "1-2 hari", "2-3 hari"
	TerverifikasiAPI bool   `json:"terverifikasi_api"`  // true jika berasal dari API Biteship resmi
}

// RingkasanKlien kalkulasi biaya yang dikirim dari browser/frontend
type RingkasanKlien struct {
	Subtotal    int `json:"subtotal"`
	OngkosKirim int `json:"ongkos_kirim"`
	TotalAkhir  int `json:"total_akhir"`
}

// PermintaanValidasiOrder payload yang dikirim saat checkout ringkasan divalidasi
type PermintaanValidasiOrder struct {
	IDPelanggan        *uint64          `json:"id_pelanggan,omitempty"`
	TipeCheckout       string           `json:"tipe_checkout"`      // "guest" atau "pelanggan"
	MetodePenerimaan   string           `json:"metode_penerimaan"`  // "Kirim via Kurir"
	Kontak             KontakPemesan    `json:"kontak"`
	Alamat             AlamatPengiriman `json:"alamat"`
	Kurir              KurirPilihan     `json:"kurir"`
	Items              []ItemOrder      `json:"items"`
	RingkasanKlien     RingkasanKlien   `json:"ringkasan_klien"`
	SimulasiBelumValid bool             `json:"simulasi_belum_valid,omitempty"` // Flag simulasi Gambar 2
}

// RincianIntegritas hasil pengecekan integritas finansial dan kelengkapan
type RincianIntegritas struct {
	KelengkapanDataValid    bool `json:"kelengkapan_data_valid"`
	IntegritasSubtotalValid bool `json:"integritas_subtotal_valid"`
	IntegritasOngkirValid   bool `json:"integritas_ongkir_valid"`
	IntegritasTotalValid    bool `json:"integritas_total_valid"`
	StokTersedia            bool `json:"stok_tersedia"`
	SubtotalServer          int  `json:"subtotal_server"`
	OngkirServer            int  `json:"ongkir_server"`
	TotalAkhirServer        int  `json:"total_akhir_server"`
	SelisihSubtotal         int  `json:"selisih_subtotal"`
	SelisihTotal            int  `json:"selisih_total"`
}

// HasilValidasiOrder respon resmi dari server-side validator
// Mendukung 2 state visual Figma:
// Gambar 1: Status = "valid", PesanJudul = "Ringkasan siap dikonfirmasi"
// Gambar 2: Status = "invalid", PesanJudul = "Belum dapat membuat pesanan"
type HasilValidasiOrder struct {
	Status           string            `json:"status"`        // "valid" (Gambar 1) atau "invalid" (Gambar 2)
	DataValid        bool              `json:"data_valid"`    // true / false
	BadgeTeks        string            `json:"badge_teks"`    // "Data valid" (Gambar 1) atau ""
	PesanJudul       string            `json:"pesan_judul"`   // "Ringkasan siap dikonfirmasi" / "Belum dapat membuat pesanan"
	PesanDetail      string            `json:"pesan_detail"`  // Pesan keterangan sesuai figma
	PesanPeringatan  string            `json:"pesan_peringatan,omitempty"` // "Tarif ongkir belum terverifikasi oleh sistem." (Gambar 2)
	KodePenolakan    string            `json:"kode_penolakan,omitempty"`
	Integritas       RincianIntegritas `json:"integritas"`
	TokenValidasi    string            `json:"token_validasi,omitempty"` // Hash keamanan anti-tampering
	DaftarPeringatan []string          `json:"daftar_peringatan,omitempty"`
	WaktuValidasi    time.Time         `json:"waktu_validasi"`
}

// PermintaanBuatPesanan payload saat user mengklik "Buat Pesanan"
type PermintaanBuatPesanan struct {
	TokenValidasi string                  `json:"token_validasi"`
	OrderData     PermintaanValidasiOrder `json:"order_data"`
}

// ResponsBuatPesanan respon saat pesanan berhasil dibuat dan tersimpan di database
type ResponsBuatPesanan struct {
	Status                string              `json:"status"` // "sukses" / "ditolak"
	NomorPesanan          string              `json:"nomor_pesanan"`
	IDPesanan             int64               `json:"id_pesanan"`
	TotalPembayaran       int                 `json:"total_pembayaran"`
	TotalPembayaranFormat string              `json:"total_pembayaran_format"`
	BatasWaktuPembayaran  string              `json:"batas_waktu_pembayaran"`
	Pesan                 string              `json:"pesan"`
	HasilValidasi         *HasilValidasiOrder `json:"hasil_validasi,omitempty"`
	WaktuDibuat           time.Time           `json:"waktu_dibuat"`
}

// AuditValidasiOrder entitas pencatatan log audit validasi ke database MySQL
type AuditValidasiOrder struct {
	ID             int64     `json:"id"`
	NomorPesanan   string    `json:"nomor_pesanan"`
	Status         string    `json:"status"`
	NamaPelanggan  string    `json:"nama_pelanggan"`
	NomorWhatsApp  string    `json:"nomor_whatsapp"`
	SubtotalKlien  int       `json:"subtotal_klien"`
	SubtotalServer int       `json:"subtotal_server"`
	OngkirKlien    int       `json:"ongkir_klien"`
	OngkirServer   int       `json:"ongkir_server"`
	TotalKlien     int       `json:"total_klien"`
	TotalServer    int       `json:"total_server"`
	KodePenolakan  string    `json:"kode_penolakan"`
	Keterangan     string    `json:"keterangan"`
	DibuatPada     time.Time `json:"dibuat_pada"`
}
