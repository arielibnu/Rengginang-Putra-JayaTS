package errorhandlingkurir

import "time"

// TipeKegagalanKurir jenis-jenis kegagalan koneksi API atau rute
type TipeKegagalanKurir string

const (
	KegagalanAPIDown            TipeKegagalanKurir = "API_DOWN"
	KegagalanTimeout            TipeKegagalanKurir = "API_TIMEOUT"
	KegagalanRuteTidakTerjangkau TipeKegagalanKurir = "RUTE_TIDAK_TERJANGKAU"
	KegagalanBeratMelebihiBatas TipeKegagalanKurir = "BERAT_MELEBIHI_BATAS"
	KegagalanDataTidakLengkap   TipeKegagalanKurir = "DATA_TIDAK_LENGKAP"
)

// ResponsErrorKurir merepresentasikan payload error terstandarisasi sesuai Figma Gambar 2
type ResponsErrorKurir struct {
	Status           string             `json:"status"` // "error"
	ErrorCode        string             `json:"error_code"` // "TARIF_TIDAK_TERSEDIA"
	Pesan            string             `json:"pesan"` // "Tarif pengiriman tidak tersedia"
	Detail           string             `json:"detail"` // "API gagal atau layanan tidak melayani wilayah tujuan."
	Peringatan       string             `json:"peringatan"` // "Jangan membuat tarif ongkir perkiraan sendiri."
	TipeError        TipeKegagalanKurir `json:"tipe_error"`
	EstimasiPalsu    bool               `json:"estimasi_palsu"` // Wajib false (sesuai aturan WBS 1.2.7.4)
	WilayahTujuan    string             `json:"wilayah_tujuan,omitempty"`
	AksiTindakan     []string           `json:"aksi_tindakan"` // ["Coba lagi", "Ubah alamat"]
	WaktuError       time.Time          `json:"waktu_error"`
}

// PermintaanSimulasiKegagalan request untuk simulasi pengujian kegagalan API kurir
type PermintaanSimulasiKegagalan struct {
	TipeError     TipeKegagalanKurir `json:"tipe_error"` // API_DOWN | RUTE_TIDAK_TERJANGKAU | API_TIMEOUT
	WilayahTujuan string             `json:"wilayah_tujuan"`
}

// LogErrorAPI model audit penyimpanan kegagalan API kurir
type LogErrorAPI struct {
	ID            int                `json:"id"`
	TipeError     TipeKegagalanKurir `json:"tipe_error"`
	WilayahTujuan string             `json:"wilayah_tujuan"`
	PesanError    string             `json:"pesan_error"`
	WaktuKejadian time.Time          `json:"waktu_kejadian"`
}
