package notifikasipesanan

import "time"

// NotifikasiPesanan merepresentasikan entitas catatan notifikasi pesanan di database
type NotifikasiPesanan struct {
	ID                   int64      `json:"id"`
	NomorPesanan         string     `json:"nomor_pesanan"`
	NomorInvoice         string     `json:"nomor_invoice"`
	TipePemberitahuan    string     `json:"tipe_pemberitahuan"` // pesanan_masuk | pesanan_terbayar | stok_menipis
	Judul                string     `json:"judul"`
	Pesan                string     `json:"pesan"`
	MetodePenerimaan     string     `json:"metode_penerimaan"` // Ambil Sendiri | Kirim Kurir
	TotalTransaksi       int64      `json:"total_transaksi"`
	TotalTransaksiFormat string     `json:"total_transaksi_format"`
	StatusPesanan        string     `json:"status_pesanan"` // baru | diproses | selesai
	SudahDibaca          bool       `json:"sudah_dibaca"`
	WaktuDibaca          *time.Time `json:"waktu_dibaca,omitempty"`
	DibuatPada           time.Time  `json:"dibuat_pada"`
	WaktuFormat          string     `json:"waktu_format"`
}

// ItemNotifikasiRingkas format rincian untuk card Notifikasi Pesanan Baru pada Figma PPL-22
type ItemNotifikasiRingkas struct {
	NomorPesanan string `json:"nomor_pesanan"`
	Metode       string `json:"metode"` // Ambil Sendiri | Kirim Kurir
	TeksRingkas  string `json:"teks_ringkas"` // PJ-DEMO-0001 — Ambil Sendiri
	Tipe         string `json:"tipe"`
	Waktu        string `json:"waktu"`
}

// ResponsPesananBaruOwner payload respons untuk komponen Notifikasi Pesanan Baru Owner Dashboard
type ResponsPesananBaruOwner struct {
	Status        string                  `json:"status"`
	Pesan         string                  `json:"pesan"`
	JumlahBaru    int                     `json:"jumlah_baru"`
	BadgeLabel    string                  `json:"badge_label"` // "2 Baru"
	DaftarPesanan []ItemNotifikasiRingkas `json:"daftar_pesanan"`
	Notifikasi    []NotifikasiPesanan     `json:"notifikasi"`
}

// PermintaanTriggerNotifikasi payload untuk trigger manual atau otomatis pembuatan notifikasi
type PermintaanTriggerNotifikasi struct {
	NomorPesanan      string `json:"nomor_pesanan"`
	NomorInvoice      string `json:"nomor_invoice"`
	TipePemberitahuan string `json:"tipe_pemberitahuan"` // pesanan_masuk | pesanan_terbayar
	Judul             string `json:"judul"`
	Pesan             string `json:"pesan"`
	MetodePenerimaan  string `json:"metode_penerimaan"`
	TotalTransaksi    int64  `json:"total_transaksi"`
}

// PermintaanTandaiDibaca payload saat Owner membuka detail atau mengklik notifikasi
type PermintaanTandaiDibaca struct {
	IDNotifikasi []int64 `json:"id_notifikasi"`
	NomorPesanan string  `json:"nomor_pesanan"`
	Semua        bool    `json:"semua"`
}
