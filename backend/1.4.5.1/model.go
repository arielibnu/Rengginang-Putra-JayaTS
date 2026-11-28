package lababersih

import "time"

// FilterPeriode mewakili rentang bulan dan tahun laporan
type FilterPeriode struct {
	BulanNama    string    `json:"bulan_nama"`
	BulanAngka   int       `json:"bulan_angka"`
	Tahun        int       `json:"tahun"`
	AwalPeriode  time.Time `json:"awal_periode"`
	AkhirPeriode time.Time `json:"akhir_periode"`
}

// SumberDataStatus memetakan ringkasan audit data sesuai desain Figma PPL-29
type SumberDataStatus struct {
	TransaksiValid   string `json:"transaksi_valid"`
	Sandbox          string `json:"sandbox"`
	DataBelumLengkap string `json:"data_belum_lengkap"`
}

// ItemRincianLaba mencatat baris kontribusi laba/beban
type ItemRincianLaba struct {
	Kategori      string  `json:"kategori"`
	Keterangan    string  `json:"keterangan"`
	Nominal       float64 `json:"nominal"`
	NominalFormat string  `json:"nominal_format"`
	Tipe          string  `json:"tipe"` // PENDAPATAN, HPP, BEBAN
}

// DataEstimasiKeuntungan memuat metrik WBS 1.4.5.1 (BE)
type DataEstimasiKeuntungan struct {
	Periode                     FilterPeriode    `json:"periode"`
	PenerimaanNyata             float64          `json:"penerimaan_nyata"`
	PenerimaanNyataFormat       string           `json:"penerimaan_nyata_format"`
	ModalBarangTerjual          float64          `json:"modal_barang_terjual"`
	ModalBarangTerjualFormat    string           `json:"modal_barang_terjual_format"`
	BiayaOperasional            float64          `json:"biaya_operasional"`
	BiayaOperasionalFormat      string           `json:"biaya_operasional_format"`
	EstimasiLabaBersih          float64          `json:"estimasi_laba_bersih"`
	EstimasiLabaBersihFormat    string           `json:"estimasi_laba_bersih_format"`
	OngkirTitipanDieksklusikan  float64          `json:"ongkir_titipan_dieksklusikan"`
	SimulasiSandboxDieksklusikan float64         `json:"simulasi_sandbox_dieksklusikan"`
	SumberData                  SumberDataStatus `json:"sumber_data"`
	RincianPerhitungan          []ItemRincianLaba `json:"rincian_perhitungan"`
	CatatanAudit                string           `json:"catatan_audit"`
	DibuatPada                  time.Time        `json:"dibuat_pada"`
}

// ResponsLaporanEstimasiKeuntungan format JSON standar API
type ResponsLaporanEstimasiKeuntungan struct {
	Status string                  `json:"status"`
	Pesan  string                  `json:"pesan,omitempty"`
	Data   *DataEstimasiKeuntungan `json:"data,omitempty"`
}
