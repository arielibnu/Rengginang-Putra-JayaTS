package reservasisstok

import (
	"time"

	ordergenerator "rengginang-putra-jaya/auth/1.2.9.1"
)

// ============================================================================
// WBS 1.2.9.2 - MEKANISME RESERVASI STOK OTOMATIS (HOLDING TIME 24 JAM)
// Deliverable: Stock Reservation & Auto-Release
// ============================================================================

// ItemReservasi item yang ingin direservasi stoknya
type ItemReservasi struct {
	IDVarian   uint64 `json:"id_varian"`
	IDProduk   string `json:"id_produk"`
	NamaVarian string `json:"nama_varian"`
	Jumlah     int    `json:"jumlah"`
}

// CatatanReservasi data record penahanan stok di database
type CatatanReservasi struct {
	ID                int64     `json:"id"`
	NomorPesanan      string    `json:"nomor_pesanan"`
	IDVarian          uint64    `json:"id_varian"`
	NamaVarian        string    `json:"nama_varian"`
	JumlahDireservasi int       `json:"jumlah_direservasi"`
	StatusReservasi   string    `json:"status_reservasi"` // "menahan_stok", "dikonfirmasi", "dilepas_otomatis"
	KadaluarsaPada    time.Time `json:"kadaluarsa_pada"`  // 24 jam dari pembuatan
	KadaluarsaFormat  string    `json:"kadaluarsa_format"`
	DibuatPada        time.Time `json:"dibuat_pada"`
}

// HasilReservasiStok laporan hasil proses penahanan stok
type HasilReservasiStok struct {
	Sukses           bool               `json:"sukses"`
	KodeHasil        string             `json:"kode_hasil"` // "RESERVASI_SUKSES", "STOK_TIDAK_CUKUP", "VARIAN_KOSONG"
	Pesan            string             `json:"pesan"`
	HoldingTimeJam   int                `json:"holding_time_jam"` // 24 jam
	BatasWaktu       time.Time          `json:"batas_waktu"`
	BatasWaktuFormat string             `json:"batas_waktu_format"`
	RincianReservasi []CatatanReservasi `json:"rincian_reservasi"`
}

// PermintaanBuatPesananPPL18 payload checkout final dari frontend
type PermintaanBuatPesananPPL18 struct {
	TipeCheckout       string                              `json:"tipe_checkout"` // "guest" atau "pelanggan"
	IDPelanggan        *uint64                             `json:"id_pelanggan,omitempty"`
	Kontak             ordergenerator.DataKontakPengiriman `json:"kontak"`
	Kurir              ordergenerator.DataKurirSnapshot    `json:"kurir"`
	Items              []ordergenerator.ItemSnapshotPesanan `json:"items"`
	TokenValidasi      string                              `json:"token_validasi,omitempty"`
	SimulasiStokHabis  bool                                `json:"simulasi_stok_habis,omitempty"` // Pengujian error stok
}

// ResponsKonfirmasiPesananPPL18 respon resmi untuk Halaman Konfirmasi PPL-18
type ResponsKonfirmasiPesananPPL18 struct {
	Status                string                         `json:"status"` // "sukses" / "gagal"
	Pesan                 string                         `json:"pesan"`
	NomorPesanan          string                         `json:"nomor_pesanan"`
	NomorInvoice          string                         `json:"nomor_invoice"`
	StatusPesanan         string                         `json:"status_pesanan"`        // "menunggu_pembayaran"
	StatusPesananLabel    string                         `json:"status_pesanan_label"`  // "Menunggu pembayaran" (badge hijau Figma)
	StatusReservasi       string                         `json:"status_reservasi"`      // "direservasi"
	BatasWaktuPembayaran  string                         `json:"batas_waktu_pembayaran"`
	HoldingTimeJam        int                            `json:"holding_time_jam"`
	TotalPembayaran       int                            `json:"total_pembayaran"`
	TotalPembayaranFormat string                         `json:"total_pembayaran_format"`
	Penerima              string                         `json:"penerima"`
	Alamat                string                         `json:"alamat"`
	Kurir                 string                         `json:"kurir"`
	HasilReservasi        *HasilReservasiStok            `json:"hasil_reservasi,omitempty"`
	PesananDetail         *ordergenerator.EntitasPesanan `json:"pesanan_detail,omitempty"`
	WaktuDibuat           time.Time                      `json:"waktu_dibuat"`
}
