package flagpemisah

import "time"

// Konstanta flag lingkungan dan kategori transaksi pemisah
const (
	FlagProduksiNyata   = "PRODUKSI_NYATA"
	FlagSandboxSimulasi = "SANDBOX_SIMULASI"

	SumberOfflineKasir = "offline_kasir"
	SumberOnlineWeb    = "online_website"
	SumberSandbox      = "sandbox_midtrans"
)

// InfoPemisahTransaksi informasi metadata klasifikasi flag transaksi
type InfoPemisahTransaksi struct {
	IDTransaksi      string    `json:"id_transaksi"`
	NomorReferensi   string    `json:"nomor_referensi"`
	Sumber           string    `json:"sumber"`           // offline_kasir | online_website | sandbox_midtrans
	FlagLingkungan   string    `json:"flag_lingkungan"`  // PRODUKSI_NYATA | SANDBOX_SIMULASI
	IsSandbox        bool      `json:"is_sandbox"`
	TotalNominal     int64     `json:"total_nominal"`
	TotalFormat      string    `json:"total_format"`
	MetodePembayaran string    `json:"metode_pembayaran"`
	Keterangan       string    `json:"keterangan"`
	WaktuTransaksi   time.Time `json:"waktu_transaksi"`
}

// RekapAuditFlag hasil audit pemisahan uang kas nyata vs simulasi sandbox
type RekapAuditFlag struct {
	TotalTransaksiNyata    int64                  `json:"total_transaksi_nyata"`
	TotalNominalNyata      int64                  `json:"total_nominal_nyata"`
	TotalNominalNyataFmt   string                 `json:"total_nominal_nyata_fmt"`
	TotalTransaksiSandbox  int64                  `json:"total_transaksi_sandbox"`
	TotalNominalSandbox    int64                  `json:"total_nominal_sandbox"`
	TotalNominalSandboxFmt string                 `json:"total_nominal_sandbox_fmt"`
	DaftarPemisah          []InfoPemisahTransaksi `json:"daftar_pemisah"`
}

// PermintaanValidasiFlag request validasi flag lingkungan transaksi
type PermintaanValidasiFlag struct {
	SumberTransaksi string `json:"sumber_transaksi"`
	FlagLingkungan  string `json:"flag_lingkungan"`
	NomorReferensi  string `json:"nomor_referensi"`
}

// ResponsValidasiFlag balasan hasil validasi dan klasifikasi
type ResponsValidasiFlag struct {
	Valid          bool   `json:"valid"`
	Pesan          string `json:"pesan"`
	FlagLingkungan string `json:"flag_lingkungan"`
	IsSandbox      bool   `json:"is_sandbox"`
}
