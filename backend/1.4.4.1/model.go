package laporanpengeluaran

import "time"

// PeriodeLaporan menyimpan informasi rentang tanggal laporan bulanan
type PeriodeLaporan struct {
	BulanNama    string    `json:"bulan_nama"`
	BulanAngka   int       `json:"bulan_angka"`
	Tahun        int       `json:"tahun"`
	AwalPeriode  time.Time `json:"awal_periode"`
	AkhirPeriode time.Time `json:"akhir_periode"`
}

// RekapSumberData menyimpan audit transparansi sumber data
type RekapSumberData struct {
	TransaksiValid   string `json:"transaksi_valid"`
	Sandbox          string `json:"sandbox"`
	DataBelumLengkap string `json:"data_belum_lengkap"`
}

// ItemPengeluaran menyimpan rincian pengeluaran per baris transaksi
type ItemPengeluaran struct {
	IDPengeluaran      uint64  `json:"id_pengeluaran"`
	TanggalPengeluaran string  `json:"tanggal_pengeluaran"`
	KategoriAsli       string  `json:"kategori_asli"`
	KategoriBeban      string  `json:"kategori_beban"` // OPERASIONAL, BELANJA_PERSEDIAAN, BIAYA_LAIN
	Nominal            float64 `json:"nominal"`
	NominalFormat      string  `json:"nominal_format"`
	Keterangan         string  `json:"keterangan"`
}

// DataLaporanPengeluaran adalah payload utama laporan pengeluaran bulanan
type DataLaporanPengeluaran struct {
	Periode            PeriodeLaporan    `json:"periode"`
	Operasional        float64           `json:"operasional"`
	OperasionalFormat  string            `json:"operasional_format"`
	BelanjaPersediaan  float64           `json:"belanja_persediaan"`
	BelanjaFormat      string            `json:"belanja_persediaan_format"`
	BiayaLain          float64           `json:"biaya_lain"`
	BiayaLainFormat    string            `json:"biaya_lain_format"`
	TotalPengeluaran   float64           `json:"total_pengeluaran"`
	TotalFormat        string            `json:"total_pengeluaran_format"`
	SumberData         RekapSumberData   `json:"sumber_data"`
	DaftarPengeluaran  []ItemPengeluaran `json:"daftar_pengeluaran"`
	CatatanAudit       string            `json:"catatan_audit"`
	DibuatPada         time.Time         `json:"dibuat_pada"`
}

// ResponseLaporanPengeluaran adalah pembungkus response JSON standar
type ResponseLaporanPengeluaran struct {
	Status string                  `json:"status"`
	Pesan  string                  `json:"pesan,omitempty"`
	Data   *DataLaporanPengeluaran `json:"data,omitempty"`
}
