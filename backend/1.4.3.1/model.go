package penerimaankas

import "time"

// PeriodeLaporan menyimpan data rentang periode laporan kas
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

// ItemArusKas menyimpan rincian bukti kas masuk per transaksi
type ItemArusKas struct {
	NomorReferensi string    `json:"nomor_referensi"`
	Sumber         string    `json:"sumber"` // kasir_offline_tunai, transfer_digital_sandbox
	WaktuLunas     time.Time `json:"waktu_lunas"`
	WaktuLunasStr  string    `json:"waktu_lunas_str"`
	Nominal        float64   `json:"nominal"`
	NominalFormat  string    `json:"nominal_format"`
	Metode         string    `json:"metode"`
	Status         string    `json:"status"` // LUNAS, PENDING, BATAL
}

// DataLaporanPenerimaanKas adalah payload utama respon laporan kas
type DataLaporanPenerimaanKas struct {
	Periode              PeriodeLaporan    `json:"periode"`
	KasMasukNyata        float64           `json:"kas_masuk_nyata"`
	KasMasukNyataFormat  string            `json:"kas_masuk_nyata_format"`
	SimulasiSandbox      float64           `json:"simulasi_sandbox"`
	SimulasiSandboxFormat string           `json:"simulasi_sandbox_format"`
	StatusPending        float64           `json:"status_pending"`
	StatusPendingFormat  string            `json:"status_pending_format"`
	TotalKasTercatat     float64           `json:"total_kas_tercatat"`
	TotalKasFormat       string            `json:"total_kas_format"`
	SumberData           RekapSumberData   `json:"sumber_data"`
	RincianArusKas       []ItemArusKas     `json:"rincian_arus_kas"`
	CatatanAudit         string            `json:"catatan_audit"`
	DibuatPada           time.Time         `json:"dibuat_pada"`
}

// ResponseLaporanPenerimaanKas pembungkus response JSON
type ResponseLaporanPenerimaanKas struct {
	Status string                    `json:"status"`
	Pesan  string                    `json:"pesan,omitempty"`
	Data   *DataLaporanPenerimaanKas `json:"data,omitempty"`
}
