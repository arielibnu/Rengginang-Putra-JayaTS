package ordergenerator

import (
	"strings"
	"testing"
)

func TestGenerateNomorPesananDanInvoiceUnik(t *testing.T) {
	no1 := GenerateNomorPesananUnik()
	no2 := GenerateNomorPesananUnik()

	if !strings.HasPrefix(no1, "PJ-") {
		t.Fatalf("Nomor pesanan harus diawali 'PJ-', didapat: %s", no1)
	}
	if no1 == no2 {
		t.Fatalf("Nomor pesanan harus unik: no1=%s, no2=%s", no1, no2)
	}

	inv1 := GenerateNomorInvoiceUnik()
	if !strings.HasPrefix(inv1, "INV-") {
		t.Fatalf("Nomor invoice harus diawali 'INV-', didapat: %s", inv1)
	}
}

func TestBuatSnapshotDanSimpanOrder(t *testing.T) {
	gen := &GeneratorOrder{DB: nil}

	kontak := DataKontakPengiriman{
		NamaPenerima:  "Pelanggan Contoh",
		NomorWhatsApp: "081234567890",
		AlamatLengkap: "Jl. Contoh No. 12",
		KotaKabupaten: "Kabupaten Jember",
		Provinsi:      "Jawa Timur",
	}

	kurir := DataKurirSnapshot{
		IDKurir:          "jne",
		NamaKurir:        "JNE REG",
		KodeLayanan:      "REG",
		TarifOngkir:      18000,
		EstimasiHari:     "2-3 hari",
		TerverifikasiAPI: true,
	}

	items := []ItemSnapshotPesanan{
		{
			IDProduk:            "prod-1",
			NamaProduk:          "Rengginang Udang",
			NamaVarian:          "Udang Gurih",
			Jumlah:              2,
			HargaSatuanSnapshot: 28000,
		},
		{
			IDProduk:            "prod-2",
			NamaProduk:          "Rengginang Bawang",
			NamaVarian:          "Bawang Renyah",
			Jumlah:              1,
			HargaSatuanSnapshot: 28000,
		},
	}

	pesanan, err := gen.BuatSnapshotDanSimpanOrder("guest", nil, kontak, kurir, items)
	if err != nil {
		t.Fatalf("Gagal membuat pesanan snapshot: %v", err)
	}

	if pesanan.SnapshotSubtotal != 84000 {
		t.Errorf("Ekspektasi subtotal snapshot 84000, didapat: %d", pesanan.SnapshotSubtotal)
	}
	if pesanan.SnapshotOngkir != 18000 {
		t.Errorf("Ekspektasi ongkir snapshot 18000, didapat: %d", pesanan.SnapshotOngkir)
	}
	if pesanan.SnapshotTotal != 102000 {
		t.Errorf("Ekspektasi total snapshot 102000, didapat: %d", pesanan.SnapshotTotal)
	}
	if pesanan.StatusPesanan != "menunggu_pembayaran" {
		t.Errorf("Ekspektasi status pesanan 'menunggu_pembayaran', didapat: %s", pesanan.StatusPesanan)
	}
	if pesanan.StatusReservasi != "direservasi" {
		t.Errorf("Ekspektasi status reservasi 'direservasi', didapat: %s", pesanan.StatusReservasi)
	}

	// Cari pesanan
	ditemukan, ada := gen.AmbilPesananBerdasarkanNomor(pesanan.NomorPesanan)
	if !ada || ditemukan == nil {
		t.Fatalf("Pesanan dengan nomor %s harus ditemukan", pesanan.NomorPesanan)
	}
}
