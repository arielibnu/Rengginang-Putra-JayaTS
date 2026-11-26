package laporanbulanan

import "time"

// PeriodeLaporan informasi rentang waktu laporan stok akhir
type PeriodeLaporan struct {
	BulanNama    string    `json:"bulan_nama"`
	BulanAngka   int       `json:"bulan_angka"`
	Tahun        int       `json:"tahun"`
	AwalPeriode  time.Time `json:"awal_periode"`
	AkhirPeriode time.Time `json:"akhir_periode"`
}

// RincianStokVarian kondisi stok tiap varian pada akhir periode
type RincianStokVarian struct {
	IDVarian             int64   `json:"id_varian"`
	NamaVarian           string  `json:"nama_varian"`
	HargaSatuan          float64 `json:"harga_satuan"`
	StokFisik            int     `json:"stok_fisik"`
	ReservasiAktif       int     `json:"reservasi_aktif"`
	StokTersedia         int     `json:"stok_tersedia"`
	NilaiPersediaan      float64 `json:"nilai_persediaan"`
	NilaiPersediaanFormat string  `json:"nilai_persediaan_format"`
	TotalMasuk           int     `json:"total_masuk"`
	TotalKeluar          int     `json:"total_keluar"`
}

// RingkasanStokAkhir akumulasi seluruh varian
type RingkasanStokAkhir struct {
	TotalStokFisik        int     `json:"total_stok_fisik"`
	TotalReservasiAktif   int     `json:"total_reservasi_aktif"`
	TotalStokTersedia     int     `json:"total_stok_tersedia"`
	TotalNilaiPersediaan  float64 `json:"total_nilai_persediaan"`
	NilaiPersediaanFormat string  `json:"nilai_persediaan_format"`
}

// SumberDataStatus status kevalidan data sesuai acuan Figma
type SumberDataStatus struct {
	TransaksiValid   string `json:"transaksi_valid"`
	Sandbox          string `json:"sandbox"`
	DataBelumLengkap string `json:"data_belum_lengkap"`
}

// LaporanStokAkhirRespons objek respons resmi endpoint laporan stok akhir
type LaporanStokAkhirRespons struct {
	Periode          PeriodeLaporan       `json:"periode"`
	StokFisik        int                  `json:"stok_fisik"`
	ReservasiAktif   int                  `json:"reservasi_aktif"`
	StokTersedia     int                  `json:"stok_tersedia"`
	Ringkasan        RingkasanStokAkhir   `json:"ringkasan"`
	SumberData       SumberDataStatus     `json:"sumber_data"`
	RincianVarian    []RincianStokVarian  `json:"rincian_varian"`
	CatatanAudit     string               `json:"catatan_audit"`
	DibuatPada       time.Time            `json:"dibuat_pada"`
}
