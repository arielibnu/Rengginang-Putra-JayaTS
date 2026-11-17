package stokapi

import "time"

type PermintaanMutasi struct {
	IDVarian    uint64  `json:"id_varian"`
	JenisMutasi string  `json:"jenis_mutasi"`
	ModeStok    string  `json:"mode_stok"`
	Jumlah      uint64  `json:"jumlah"`
	StokBaru    *uint64 `json:"stok_baru,omitempty"`
	Catatan     string  `json:"catatan"`
}

type RiwayatMutasi struct {
	IDMutasi    uint64    `json:"id_mutasi"`
	IDVarian    uint64    `json:"id_varian"`
	JenisMutasi string    `json:"jenis_mutasi"`
	ModeStok    string    `json:"mode_stok"`
	Jumlah      uint64    `json:"jumlah"`
	StokSebelum uint64    `json:"stok_sebelum"`
	StokSesudah uint64    `json:"stok_sesudah"`
	Catatan     string    `json:"catatan"`
	DibuatPada  time.Time `json:"dibuat_pada"`
}

type ResponsAPI struct {
	Pesan string      `json:"pesan"`
	Data  interface{} `json:"data,omitempty"`
}
