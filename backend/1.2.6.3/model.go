package validasiwilayah

// ============================================================================
// WBS 1.2.6.3 - VALIDASI KELENGKAPAN WILAYAH & ALAMAT TUJUAN SESUAI PARAMETER API KURIR
// Deliverable: Validasi Data Geografis Pengiriman
// Product Backlog: PPL-15 — Pengisian Wilayah dan Alamat Pengiriman
// ============================================================================

// InformasiKontak menyimpan data kontak pemesan/penerima paket.
type InformasiKontak struct {
	NamaPenerima  string `json:"nama_penerima"`
	NomorWhatsApp string `json:"nomor_whatsapp"`
	Email         string `json:"email,omitempty"`
	Catatan       string `json:"catatan,omitempty"`
}

// WilayahPengiriman menyimpan data hierarki wilayah dan alamat pengiriman.
type WilayahPengiriman struct {
	Provinsi      string `json:"provinsi"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Kecamatan     string `json:"kecamatan"`
	KodePos       string `json:"kode_pos"`
	AlamatLengkap string `json:"alamat_lengkap"`
	CatatanAlamat string `json:"catatan_alamat,omitempty"`
}

// PermintaanWilayahPengiriman adalah payload yang dikirim dari Frontend Checkout.
type PermintaanWilayahPengiriman struct {
	Kontak            InformasiKontak    `json:"kontak"`
	MetodePenerimaan  string             `json:"metode_penerimaan"` // "kurir" atau "ambil_sendiri"
	WilayahPengiriman *WilayahPengiriman `json:"wilayah_pengiriman"`
	SubtotalProduk    uint64             `json:"subtotal_produk"`
	TotalSementara    uint64             `json:"total_sementara"`
	TotalBeratGram    uint32             `json:"total_berat_gram,omitempty"`
}

// LokasiGeografis adalah entitas data geografis terstandarisasi untuk kurir ekspedisi.
type LokasiGeografis struct {
	Provinsi      string `json:"provinsi"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Kecamatan     string `json:"kecamatan"`
	KodePos       string `json:"kode_pos"`
	AlamatLengkap string `json:"alamat_lengkap"`
	KodeWilayahID string `json:"kode_wilayah_id,omitempty"`
}

// PenerimaKurir data kontak penerima paket terstandarisasi.
type PenerimaKurir struct {
	Nama           string `json:"nama"`
	Telepon        string `json:"telepon"`
	TeleponStandar string `json:"telepon_standar"` // format internasional E.164 (+62...)
}

// PaketKurir spesifikasi paket barang yang akan dikirim via kurir ekspedisi.
type PaketKurir struct {
	BeratGram uint32 `json:"berat_gram"`
	Kategori  string `json:"kategori"`
	IsFragile bool   `json:"is_fragile"` // Rengginang adalah makanan renyah mudah hancur
}

// ParameterKurirStandar adalah payload parameter kurir yang siap dikirim ke API kurir (PPL-16).
type ParameterKurirStandar struct {
	AsalPengiriman   LokasiGeografis `json:"asal_pengiriman"`
	TujuanPengiriman LokasiGeografis `json:"tujuan_pengiriman"`
	Penerima         PenerimaKurir   `json:"penerima"`
	Paket            PaketKurir      `json:"paket"`
	IsKurirReady     bool            `json:"is_kurir_ready"`
	StatusGeografis  string          `json:"status_geografis"` // "VALID", "INCOMPLETE", "NOT_REQUIRED"
}

// HasilValidasiWilayah adalah data yang dikembalikan pada respon sukses.
type HasilValidasiWilayah struct {
	Status           string                 `json:"status"`
	Pesan            string                 `json:"pesan"`
	Valid            bool                   `json:"valid"`
	MetodePenerimaan string                 `json:"metode_penerimaan"`
	Wilayah          *WilayahPengiriman     `json:"wilayah,omitempty"`
	ParameterKurir   *ParameterKurirStandar `json:"parameter_kurir,omitempty"`
}

// ResponsAPI standar struktur kembalian HTTP endpoint.
type ResponsAPI struct {
	Status string            `json:"status"` // "sukses" atau "error"
	Pesan  string            `json:"pesan"`
	Data   interface{}       `json:"data,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}
