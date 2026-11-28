package auditkeuntungan

import "time"

// FilterPeriode mewakili rentang bulan dan tahun laporan
type FilterPeriode struct {
	BulanNama  string `json:"bulan_nama"`
	BulanAngka int    `json:"bulan_angka"`
	Tahun      int    `json:"tahun"`
}

// AuditAntiDuplikasi memuat statistik audit integritas transaksi
type AuditAntiDuplikasi struct {
	JumlahTransaksiValid         int     `json:"jumlah_transaksi_valid"`
	JumlahTransaksiTerdeduplikasi int     `json:"jumlah_transaksi_terdeduplikasi"`
	TotalNominalPenjualanValid   float64 `json:"total_nominal_penjualan_valid"`
	TotalNominalFormat           string  `json:"total_nominal_format"`
	OngkirTitipanDikecualikan    float64 `json:"ongkir_titipan_dikecualikan"`
	SandboxDikecualikan          float64 `json:"sandbox_dikecualikan"`
	StatusIntegritas             string  `json:"status_integritas"` // VALID_ANTI_DUPLIKASI
}

// HasilAuditKeuntungan memuat metrik WBS 1.4.5.3 (BE)
type HasilAuditKeuntungan struct {
	Periode              FilterPeriode      `json:"periode"`
	KelengkapanBiaya     bool               `json:"kelengkapan_biaya"`
	StatusKelengkapan    string             `json:"status_kelengkapan"`
	DisclaimerIndikator  string             `json:"disclaimer_indikator"`
	AuditAntiDuplikasi   AuditAntiDuplikasi `json:"audit_anti_duplikasi"`
	CatatanTransparansi  string             `json:"catatan_transparansi"`
	DiverifikasiPada     time.Time          `json:"diverifikasi_pada"`
}

// ResponsAuditKeuntungan format JSON respons API
type ResponsAuditKeuntungan struct {
	Status string                `json:"status"`
	Pesan  string                `json:"pesan,omitempty"`
	Data   *HasilAuditKeuntungan `json:"data,omitempty"`
}
