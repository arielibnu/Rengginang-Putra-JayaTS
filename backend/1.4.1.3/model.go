package pemisahpenjualan

import "time"

// FilterPemisahTransaksi kriteria filter transaksi riil vs sandbox
type FilterPemisahTransaksi struct {
	BulanAngka int
	Tahun      int
}

// HasilPemisahanPenjualan hasil audit dan pemisahan transaksi riil vs simulasi online & batal
type HasilPemisahanPenjualan struct {
	PenjualanOfflineRiil float64   `json:"penjualan_offline_riil"`
	SimulasiSandbox      float64   `json:"simulasi_sandbox"`
	JumlahPesananBatal   int       `json:"jumlah_pesanan_batal"`
	TotalOmzetRiil       float64   `json:"total_omzet_riil"`
	Keterangan           string    `json:"keterangan"`
	WaktuAudit           time.Time `json:"waktu_audit"`
}
