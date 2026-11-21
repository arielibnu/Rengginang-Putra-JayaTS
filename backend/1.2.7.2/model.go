package gatewayongkir

import "time"

// LokasiAsal merepresentasikan data toko asal pengiriman
type LokasiAsal struct {
	Provinsi      string `json:"provinsi"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Kecamatan     string `json:"kecamatan"`
	KodePos       string `json:"kode_pos"`
	NamaToko      string `json:"nama_toko"`
}

// LokasiTujuan merepresentasikan data wilayah penerima paket
type LokasiTujuan struct {
	Provinsi      string `json:"provinsi"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Kecamatan     string `json:"kecamatan"`
	KodePos       string `json:"kode_pos"`
	AlamatLengkap string `json:"alamat_lengkap"`
	NamaPenerima  string `json:"nama_penerima"`
	NomorWhatsApp string `json:"nomor_whatsapp"`
}

// ItemPaket merepresentasikan produk pesanan yang dihitung beratnya
type ItemPaket struct {
	IDProduk         string `json:"id_produk"`
	NamaProduk       string `json:"nama_produk"`
	Jumlah           int    `json:"jumlah"`
	BeratSatuanGram  int    `json:"berat_satuan_gram"`
	TotalBeratGram   int    `json:"total_berat_gram"`
	HargaSatuan      int    `json:"harga_satuan"`
	Subtotal         int    `json:"subtotal"`
}

// PermintaanHitungOngkir merepresentasikan request payload dari frontend
type PermintaanHitungOngkir struct {
	Asal          *LokasiAsal   `json:"asal,omitempty"`
	Tujuan        LokasiTujuan  `json:"tujuan"`
	BeratGram     int           `json:"berat_gram"`
	KurirPilihan  []string      `json:"kurir_pilihan,omitempty"` // jne, jnt, sicepat
	Items         []ItemPaket   `json:"items,omitempty"`
	SimulasiError bool          `json:"simulasi_error,omitempty"` // Flag simulasi untuk pengetesan WBS 1.2.7.4
}

// LayananKurir merepresentasikan detail tarif dan estimasi layanan kurir sandbox
type LayananKurir struct {
	IDKurir          string `json:"id_kurir"`           // "jne", "jnt", "sicepat"
	KodeLayanan      string `json:"kode_layanan"`       // "REG", "EZ", "REG"
	NamaKurir        string `json:"nama_kurir"`         // "JNE REG", "J&T Express", "SiCepat"
	Tarif            int    `json:"tarif"`              // Nominal rupiah
	TarifFormat      string `json:"tarif_format"`       // "Rp18.000"
	EstimasiHari     string `json:"estimasi_hari"`      // "2-3 hari"
	Keterangan       string `json:"keterangan"`         // "Estimasi dan harga: mengikuti API pengiriman"
	BeratDihitung    int    `json:"berat_dihitung_gram"` // Gram
	Terpilih         bool   `json:"terpilih"`
}

// ResponsKalkulasiOngkir merepresentasikan hasil perhitungan ongkos kirim sandbox
type ResponsKalkulasiOngkir struct {
	Status           string         `json:"status"` // "sukses" atau "error"
	Pesan            string         `json:"pesan"`
	Asal             LokasiAsal     `json:"asal"`
	Tujuan           LokasiTujuan   `json:"tujuan"`
	BeratTotalGram   int            `json:"berat_total_gram"`
	OpsiLayanan      []LayananKurir `json:"opsi_layanan"`
	SumberAPI        string         `json:"sumber_api"` // "Biteship API Sandbox Gateway"
	WaktuKalkulasi   time.Time      `json:"waktu_kalkulasi"`
}

// Struct payload Biteship API (https://api.biteship.com/v1/rates/couriers)
type BiteshipRatesRequest struct {
	OriginPostalCode      int            `json:"origin_postal_code"`
	DestinationPostalCode int            `json:"destination_postal_code"`
	Couriers              string         `json:"couriers"` // "jne,jnt,sicepat"
	Items                 []BiteshipItem `json:"items"`
}

type BiteshipItem struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	Quantity int    `json:"quantity"`
	Weight   int    `json:"weight"`
}

type BiteshipRatesResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Object  string            `json:"object"`
	Pricing []BiteshipPricing `json:"pricing"`
	Code    int               `json:"code,omitempty"`
}

type BiteshipPricing struct {
	AvailableForCashOnDelivery bool   `json:"available_for_cash_on_delivery"`
	Company                    string `json:"company"`              // "jne", "jnt", "sicepat"
	CourierName                string `json:"courier_name"`         // "JNE", "J&T", "SiCepat"
	CourierServiceName         string `json:"courier_service_name"`  // "Reguler", "EZ", "SIUNTUNG"
	CourierServiceCode         string `json:"courier_service_code"`  // "reg", "ez", "siunt"
	Description                string `json:"description"`
	Duration                   string `json:"duration"` // "2 - 3 hari"
	Price                      int    `json:"price"`
	Type                       string `json:"type"`
}

// PermintaanPilihKurir merepresentasikan request pemilihan kurir final
type PermintaanPilihKurir struct {
	IDKurir      string `json:"id_kurir"`
	NamaKurir    string `json:"nama_kurir"`
	Tarif        int    `json:"tarif"`
	EstimasiHari string `json:"estimasi_hari"`
	AlamatTujuan string `json:"alamat_tujuan"`
	NamaPenerima string `json:"nama_penerima"`
}

// ResponsUmumOngkir struktur respons standar JSON
type ResponsUmumOngkir struct {
	Status string      `json:"status"`
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}
