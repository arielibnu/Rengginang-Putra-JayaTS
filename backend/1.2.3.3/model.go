package checkoutapi

import (
	"encoding/json"
	"strconv"
	"time"
)

// ============================================================================
// WBS 1.2.3.3 - PENYAMAAN FORMAT PAYLOAD DATA CHECKOUT BAGI PEMBELI GUEST & PELANGGAN
// Deliverable: Standarisasi Schema Order Request
// ============================================================================

// PayloadCheckoutRequest adalah skema standar request pemesanan checkout (Order Request).
// Skema ini digunakan secara seragam baik untuk pembeli tanpa login (Guest)
// maupun pelanggan yang sudah login (Registered Customer).
type PayloadCheckoutRequest struct {
	TipeCheckout string                 `json:"tipe_checkout"` // "guest" atau "pelanggan"
	IsGuest      bool                   `json:"is_guest"`      // true jika guest, false jika pelanggan
	IDPelanggan  *uint64                `json:"id_pelanggan"`  // null untuk guest, terisi ID untuk pelanggan
	Kontak       InformasiKontak        `json:"kontak"`        // Data kontak penerima pesanan
	ItemPesanan  []ItemPesananCheckout  `json:"item_pesanan"`  // Daftar varian produk yang dipesan
	Ringkasan    RingkasanBiayaCheckout `json:"ringkasan"`     // Kalkulasi ringkasan biaya pesanan
}

// UnmarshalJSON kustom untuk mendukung fleksibilitas penamaan field dari frontend.
func (p *PayloadCheckoutRequest) UnmarshalJSON(data []byte) error {
	type Alias PayloadCheckoutRequest
	aux := struct {
		RawIDPelanggan interface{} `json:"id_pelanggan"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalisasi IDPelanggan
	if aux.RawIDPelanggan != nil {
		switch v := aux.RawIDPelanggan.(type) {
		case float64:
			id := uint64(v)
			p.IDPelanggan = &id
		case string:
			if num, err := strconv.ParseUint(v, 10, 64); err == nil && num > 0 {
				p.IDPelanggan = &num
			}
		}
	}

	// Sinkronisasi tipe_checkout dan is_guest
	if p.TipeCheckout == "pelanggan" || (!p.IsGuest && p.IDPelanggan != nil) {
		p.TipeCheckout = "pelanggan"
		p.IsGuest = false
	} else {
		p.TipeCheckout = "guest"
		p.IsGuest = true
		p.IDPelanggan = nil
	}

	return nil
}

// InformasiKontak rincian data penerima pesanan (minimal nama & nomor WhatsApp).
type InformasiKontak struct {
	NamaPenerima  string `json:"nama_penerima"`
	NomorWhatsApp string `json:"nomor_whatsapp"`
	Email         string `json:"email,omitempty"`
	Catatan       string `json:"catatan,omitempty"`
}

// ItemPesananCheckout spesifikasi varian produk yang dipesan.
type ItemPesananCheckout struct {
	IDVarian    uint64 `json:"id_varian"`
	Slug        string `json:"slug,omitempty"`
	NamaVarian  string `json:"nama_varian"`
	NamaProduk  string `json:"nama_produk"`
	HargaSatuan uint64 `json:"harga_satuan"`
	Jumlah      uint32 `json:"jumlah"`
	Subtotal    uint64 `json:"subtotal"`
}

// RingkasanBiayaCheckout kalkulasi total harga, ongkos kirim, dan grand total.
type RingkasanBiayaCheckout struct {
	TotalItem       uint32 `json:"total_item"`
	Subtotal        uint64 `json:"subtotal"`
	OngkosKirim     uint64 `json:"ongkos_kirim"`
	TotalPembayaran uint64 `json:"total_pembayaran"`
}

// HasilCheckoutResponse format respon standar setelah payload checkout divalidasi.
type HasilCheckoutResponse struct {
	NomorPesanan   string                 `json:"nomor_pesanan"`
	TipeCheckout   string                 `json:"tipe_checkout"`
	IsGuest        bool                   `json:"is_guest"`
	IDPelanggan    *uint64                `json:"id_pelanggan"`
	Kontak         InformasiKontak        `json:"kontak"`
	ItemPesanan    []ItemPesananCheckout  `json:"item_pesanan"`
	Ringkasan      RingkasanBiayaCheckout `json:"ringkasan"`
	StatusPesanan  string                 `json:"status_pesanan"`
	PesanNotifikasi string                `json:"pesan_notifikasi"`
	DibuatPada     time.Time              `json:"dibuat_pada"`
}

// ResponsAPI format standar respon JSON server.
type ResponsAPI struct {
	Status string      `json:"status"` // "sukses", "gagal", "ditolak"
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}
