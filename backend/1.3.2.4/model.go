package kedaluwarsastok

import "time"

// ItemReservasiDilepas rincian varian produk yang diputus penahanannya
type ItemReservasiDilepas struct {
	IDReservasi  int64  `json:"id_reservasi"`
	IDVarian     int64  `json:"id_varian"`
	NamaVarian   string `json:"nama_varian"`
	JumlahLepas  int    `json:"jumlah_lepas"`
}

// PermintaanPutusReservasi payload pemutusan reservasi stok karena kedaluwarsa
type PermintaanPutusReservasi struct {
	NomorPesanan  string `json:"nomor_pesanan"`
	Alasan        string `json:"alasan,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// HasilPemutusanReservasi respon rinci setelah pemutusan reservasi berhasil dieksekusi
type HasilPemutusanReservasi struct {
	Sukses               bool                   `json:"sukses"`
	NomorPesanan         string                 `json:"nomor_pesanan"`
	TotalItemDilepas     int                    `json:"total_item_dilepas"`
	TotalKuantitasStok   int                    `json:"total_kuantitas_stok"`
	StatusPesananBaru    string                 `json:"status_pesanan_baru"`
	StatusReservasiBaru  string                 `json:"status_reservasi_baru"`
	WaktuPemutusan       time.Time              `json:"waktu_pemutusan"`
	WaktuPemutusanFormat string                 `json:"waktu_pemutusan_format"`
	Pesan                string                 `json:"pesan"`
	DaftarItem           []ItemReservasiDilepas `json:"daftar_item"`
}

// LogPemutusanReservasi entri jejak audit pelepasan stok
type LogPemutusanReservasi struct {
	ID               int64     `json:"id"`
	NomorPesanan     string    `json:"nomor_pesanan"`
	IDVarian         int64     `json:"id_varian"`
	NamaVarian       string    `json:"nama_varian"`
	JumlahDilepas    int       `json:"jumlah_dilepas"`
	StatusSebelumnya string    `json:"status_sebelumnya"`
	StatusAkhir      string    `json:"status_akhir"`
	Alasan           string    `json:"alasan"`
	DibuatPada       time.Time `json:"dibuat_pada"`
}
