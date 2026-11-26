package kedaluwarsastok

import (
	"testing"
	"time"
)

func TestPutusReservasiStokKedaluwarsa(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	orderID := "PJ-20261125-9999"
	items := []CatatanReservasiLokal{
		{
			ID:              1,
			NomorPesanan:    orderID,
			IDVarian:        1,
			NamaVarian:      "Rengginang Udang Gurih",
			Jumlah:          3,
			StatusReservasi: "menahan_stok",
			KadaluarsaPada:  time.Now().Add(-1 * time.Hour), // sudah kadaluarsa
		},
		{
			ID:              2,
			NomorPesanan:    orderID,
			IDVarian:        2,
			NamaVarian:      "Rengginang Bawang Renyah",
			Jumlah:          2,
			StatusReservasi: "menahan_stok",
			KadaluarsaPada:  time.Now().Add(-1 * time.Hour),
		},
	}

	layanan.TambahReservasiLokal(orderID, items)

	hasil, err := layanan.PutusReservasiStokKedaluwarsa(orderID, "Uji coba kedaluwarsa Midtrans")
	if err != nil {
		t.Fatalf("Gagal putus reservasi stok: %v", err)
	}

	if !hasil.Sukses {
		t.Errorf("Ekspektasi hasil sukses bernilai true")
	}

	if hasil.StatusPesananBaru != "kedaluwarsa" {
		t.Errorf("Ekspektasi status pesanan 'kedaluwarsa', didapat: %s", hasil.StatusPesananBaru)
	}

	if hasil.StatusReservasiBaru != "dilepas_otomatis" {
		t.Errorf("Ekspektasi status reservasi 'dilepas_otomatis', didapat: %s", hasil.StatusReservasiBaru)
	}

	if hasil.TotalKuantitasStok != 5 {
		t.Errorf("Ekspektasi 5 bungkus stok dilepas kembali, didapat: %d", hasil.TotalKuantitasStok)
	}
}
