package ordergenerator

import "time"

// ============================================================================
// WBS 1.2.9.1 - PEMBUATAN ID PESANAN UNIK SERTA PENYIMPANAN SNAPSHOT HARGA
// Deliverable: Logika Database Order & Invoice ID
// ============================================================================

// ItemSnapshotPesanan rincian snapshot harga produk saat transaksi dilakukan
type ItemSnapshotPesanan struct {
	IDVarian            uint64 `json:"id_varian"`
	IDProduk            string `json:"id_produk"`
	NamaProduk          string `json:"nama_produk"`
	NamaVarian          string `json:"nama_varian"`
	Jumlah              int    `json:"jumlah"`
	HargaSatuanSnapshot int    `json:"harga_satuan_snapshot"` // Snapshot harga saat checkout
	SubtotalSnapshot    int    `json:"subtotal_snapshot"`
	BeratSatuanGram     int    `json:"berat_satuan_gram"`
}

// DataKontakPengiriman snapshot data penerima paket
type DataKontakPengiriman struct {
	NamaPenerima  string `json:"nama_penerima"`
	NomorWhatsApp string `json:"nomor_whatsapp"`
	AlamatLengkap string `json:"alamat_lengkap"`
	KotaKabupaten string `json:"kota_kabupaten"`
	Provinsi      string `json:"provinsi"`
	KodePos       string `json:"kode_pos,omitempty"`
}

// DataKurirSnapshot snapshot kurir dan tarif ongkir resmi saat checkout
type DataKurirSnapshot struct {
	IDKurir          string `json:"id_kurir"`
	NamaKurir        string `json:"nama_kurir"`
	KodeLayanan      string `json:"kode_layanan"`
	TarifOngkir      int    `json:"tarif_ongkir"`
	EstimasiHari     string `json:"estimasi_hari"`
	TerverifikasiAPI bool   `json:"terverifikasi_api"`
}

// EntitasPesanan entitas transaksi pesanan resmi
type EntitasPesanan struct {
	ID                   int64                  `json:"id"`
	NomorPesanan         string                 `json:"nomor_pesanan"`         // Format PJ-YYYYMMDD-XXXX
	NomorInvoice         string                 `json:"nomor_invoice"`         // Format INV-YYYYMMDD-XXXX
	TipeCheckout         string                 `json:"tipe_checkout"`         // "guest" atau "pelanggan"
	IDPelanggan          *uint64                `json:"id_pelanggan,omitempty"`
	Kontak               DataKontakPengiriman   `json:"kontak"`
	Kurir                DataKurirSnapshot      `json:"kurir"`
	Items                []ItemSnapshotPesanan  `json:"items"`
	SnapshotSubtotal     int                    `json:"snapshot_subtotal"`
	SnapshotOngkir       int                    `json:"snapshot_ongkir"`
	SnapshotTotal        int                    `json:"snapshot_total"`
	StatusPesanan        string                 `json:"status_pesanan"`        // "menunggu_pembayaran", "dibayar", "dibatalkan"
	StatusReservasi      string                 `json:"status_reservasi"`      // "direservasi", "dikonfirmasi", "dilepas"
	BatasWaktuReservasi  time.Time              `json:"batas_waktu_reservasi"` // Holding time 24 jam
	BatasWaktuFormat     string                 `json:"batas_waktu_format"`
	DibuatPada           time.Time              `json:"dibuat_pada"`
	DiperbaruiPada       time.Time              `json:"diperbarui_pada"`
}
