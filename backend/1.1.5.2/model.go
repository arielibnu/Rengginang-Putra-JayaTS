package peringatanstok

import "time"

// PermintaanAturBatasMinimum menampung payload request untuk mengubah batas minimum stok suatu varian.
type PermintaanAturBatasMinimum struct {
	IDVarian     uint64 `json:"id_varian"`
	BatasMinimum uint64 `json:"batas_minimum"`
	Catatan      string `json:"catatan,omitempty"`
}

// KonfigurasiBatasMinimum merepresentasikan data batas minimum yang tersimpan di basis data.
type KonfigurasiBatasMinimum struct {
	IDVarian      uint64    `json:"id_varian"`
	NamaVarian    string    `json:"nama_varian,omitempty"`
	BatasMinimum  uint64    `json:"batas_minimum"`
	Catatan       string    `json:"catatan,omitempty"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

// StatusStokVarian menyajikan evaluasi status persediaan suatu varian terhadap batas minimumnya.
type StatusStokVarian struct {
	IDVarian     uint64 `json:"id_varian"`
	NamaVarian   string `json:"nama_varian"`
	StokFisik    uint64 `json:"stok_fisik"`
	Reservasi    uint64 `json:"reservasi"`
	StokTersedia int64  `json:"stok_tersedia"`
	BatasMinimum uint64 `json:"batas_minimum"`
	Status       string `json:"status"` // "Normal" atau "Stok menipis"
	ApakahKritis bool   `json:"apakah_kritis"`
}

// RingkasanPeringatanStok merangkum kondisi stok seluruh varian untuk dashboard dan notifikasi Owner.
type RingkasanPeringatanStok struct {
	TotalVarianAktif    int                `json:"total_varian_aktif"`
	JumlahVarianMenipis int                `json:"jumlah_varian_menipis"`
	AdaPeringatan       bool               `json:"ada_peringatan"`
	PesanUtama          string             `json:"pesan_utama"`
	KeteranganRinci     string             `json:"keterangan_rinci,omitempty"`
	DaftarVarian        []StatusStokVarian `json:"daftar_varian"`
}

// ResponsAPI format standar respon JSON backend.
type ResponsAPI struct {
	Status string      `json:"status,omitempty"`
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}
