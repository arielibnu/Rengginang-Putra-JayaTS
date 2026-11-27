package laporanpenjualan

import "time"

// PeriodeLaporan menyimpan data periode bulan dan tahun laporan
type PeriodeLaporan struct {
	BulanNama    string    `json:"bulan_nama"`
	BulanAngka   int       `json:"bulan_angka"`
	Tahun        int       `json:"tahun"`
	AwalPeriode  time.Time `json:"awal_periode"`
	AkhirPeriode time.Time `json:"akhir_periode"`
}

// RekapSumberData menyimpan audit status data
type RekapSumberData struct {
	TransaksiValid   string `json:"transaksi_valid"`
	Sandbox          string `json:"sandbox"`
	DataBelumLengkap string `json:"data_belum_lengkap"`
}

// RincianPenjualanVarian menyimpan total barang terjual per varian
type RincianPenjualanVarian struct {
	IDVarian           int64   `json:"id_varian"`
	NamaVarian         string  `json:"nama_varian"`
	JumlahTerjualNyata int     `json:"jumlah_terjual_nyata"`
	OmzetNyata         float64 `json:"omzet_nyata"`
	OmzetNyataFormat   string  `json:"omzet_nyata_format"`
	JumlahSimulasi     int     `json:"jumlah_simulasi"`
	NilaiSimulasi      float64 `json:"nilai_simulasi"`
}

// DataLaporanPenjualan adalah payload utama response laporan penjualan
type DataLaporanPenjualan struct {
	Periode                PeriodeLaporan           `json:"periode"`
	PenjualanOfflineNyata  float64                  `json:"penjualan_offline_nyata"`
	PenjualanOfflineFormat string                   `json:"penjualan_offline_nyata_format"`
	SimulasiOnline         float64                  `json:"simulasi_online"`
	SimulasiOnlineFormat   string                   `json:"simulasi_online_format"`
	PesananBatal           int                      `json:"pesanan_batal"`
	PesananBatalFormat     string                   `json:"pesanan_batal_format"`
	TotalBarangTerjualRiil int                      `json:"total_barang_terjual_riil"`
	TotalOmzetRiil         float64                  `json:"total_omzet_riil"`
	TotalOmzetRiilFormat   string                   `json:"total_omzet_riil_format"`
	SumberData             RekapSumberData          `json:"sumber_data"`
	RincianPerVarian       []RincianPenjualanVarian `json:"rincian_per_varian"`
	CatatanAudit           string                   `json:"catatan_audit"`
	DibuatPada             time.Time                `json:"dibuat_pada"`
}

// ResponseLaporanPenjualan adalah response pembungkus JSON standar
type ResponseLaporanPenjualan struct {
	Status string                `json:"status"`
	Pesan  string                `json:"pesan,omitempty"`
	Data   *DataLaporanPenjualan `json:"data,omitempty"`
}
