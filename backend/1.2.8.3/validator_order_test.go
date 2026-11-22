package validatororder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestValidasiOrder_Gambar1_DataValid menguji kasus sukses yang sesuai dengan Figma Gambar 1
func TestValidasiOrder_Gambar1_DataValid(t *testing.T) {
	req := &PermintaanValidasiOrder{
		TipeCheckout:     "guest",
		MetodePenerimaan: "Kirim via Kurir",
		Kontak: KontakPemesan{
			Nama:          "Pelanggan Contoh",
			NomorWhatsApp: "081234567890",
		},
		Alamat: AlamatPengiriman{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "Kabupaten Jember",
			Kecamatan:     "Kecamatan Sumbersari",
			AlamatLengkap: "Jl. Contoh No. 12",
		},
		Kurir: KurirPilihan{
			IDKurir:          "jne",
			NamaKurir:        "JNE REG",
			KodeLayanan:      "REG",
			TarifOngkir:      18000,
			EstimasiHari:     "2-3 hari",
			TerverifikasiAPI: true,
		},
		Items: []ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
			{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
		},
		RingkasanKlien: RingkasanKlien{
			Subtotal:    84000,
			OngkosKirim: 18000,
			TotalAkhir:  102000,
		},
	}

	hasil := ValidasiOrder(req)

	if hasil.Status != "valid" {
		t.Fatalf("Ekspektasi status 'valid', didapat: %s", hasil.Status)
	}
	if !hasil.DataValid {
		t.Fatal("Ekspektasi data_valid bernilai true")
	}
	if hasil.BadgeTeks != "Data valid" {
		t.Errorf("Ekspektasi badge 'Data valid', didapat: %s", hasil.BadgeTeks)
	}
	if hasil.PesanJudul != "Ringkasan siap dikonfirmasi" {
		t.Errorf("Ekspektasi judul 'Ringkasan siap dikonfirmasi', didapat: %s", hasil.PesanJudul)
	}
	if hasil.Integritas.SubtotalServer != 84000 {
		t.Errorf("Ekspektasi subtotal server 84000, didapat: %d", hasil.Integritas.SubtotalServer)
	}
	if hasil.Integritas.TotalAkhirServer != 102000 {
		t.Errorf("Ekspektasi total server 102000, didapat: %d", hasil.Integritas.TotalAkhirServer)
	}
	if hasil.TokenValidasi == "" {
		t.Error("Ekspektasi token validasi HMAC terisi, tetapi kosong")
	}
}

// TestValidasiOrder_Gambar2_OngkirBelumTerverifikasi menguji penolakan sesuai Figma Gambar 2
func TestValidasiOrder_Gambar2_OngkirBelumTerverifikasi(t *testing.T) {
	req := &PermintaanValidasiOrder{
		TipeCheckout:     "guest",
		MetodePenerimaan: "Kirim via Kurir",
		Kontak: KontakPemesan{
			Nama:          "Pelanggan Contoh",
			NomorWhatsApp: "081234567890",
		},
		Alamat: AlamatPengiriman{
			AlamatLengkap: "Jl. Contoh No. 12",
			KotaKabupaten: "Kabupaten Jember",
		},
		Kurir: KurirPilihan{
			IDKurir:          "", // Belum pilih kurir
			TarifOngkir:      0,
			TerverifikasiAPI: false, // Ongkir belum terverifikasi API
		},
		Items: []ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
			{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
		},
		RingkasanKlien: RingkasanKlien{
			Subtotal:    84000,
			OngkosKirim: 0,
			TotalAkhir:  84000,
		},
	}

	hasil := ValidasiOrder(req)

	if hasil.Status != "invalid" {
		t.Fatalf("Ekspektasi status 'invalid', didapat: %s", hasil.Status)
	}
	if hasil.DataValid {
		t.Fatal("Ekspektasi data_valid bernilai false")
	}
	if hasil.PesanPeringatan != "Tarif ongkir belum terverifikasi oleh sistem." {
		t.Errorf("Ekspektasi pesan peringatan persis Figma Gambar 2, didapat: %s", hasil.PesanPeringatan)
	}
	if hasil.KodePenolakan != "ONGKIR_BELUM_TERVERIFIKASI" {
		t.Errorf("Ekspektasi kode penolakan ONGKIR_BELUM_TERVERIFIKASI, didapat: %s", hasil.KodePenolakan)
	}
}

// TestValidasiOrder_PenolakanTamperingSubtotal menguji anti manipulasi subtotal
func TestValidasiOrder_PenolakanTamperingSubtotal(t *testing.T) {
	req := &PermintaanValidasiOrder{
		Kontak: KontakPemesan{Nama: "Pembeli Curang", NomorWhatsApp: "08123456789"},
		Kurir: KurirPilihan{
			IDKurir: "jne", TarifOngkir: 18000, TerverifikasiAPI: true,
		},
		Items: []ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
		},
		RingkasanKlien: RingkasanKlien{
			Subtotal:    10000, // Dimanipulasi dari 56.000 menjadi 10.000!
			OngkosKirim: 18000,
			TotalAkhir:  28000,
		},
	}

	hasil := ValidasiOrder(req)

	if hasil.DataValid {
		t.Fatal("Ekspektasi manipulasi subtotal ditolak oleh server!")
	}
	if hasil.Integritas.IntegritasSubtotalValid {
		t.Error("Ekspektasi IntegritasSubtotalValid false karena ada tampering harga")
	}
}

// TestValidasiOrder_PenolakanTamperingTotalAkhir menguji anti manipulasi total akhir
func TestValidasiOrder_PenolakanTamperingTotalAkhir(t *testing.T) {
	req := &PermintaanValidasiOrder{
		Kontak: KontakPemesan{Nama: "Pembeli Curang", NomorWhatsApp: "08123456789"},
		Kurir: KurirPilihan{
			IDKurir: "jne", TarifOngkir: 18000, TerverifikasiAPI: true,
		},
		Items: []ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
			{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
		},
		RingkasanKlien: RingkasanKlien{
			Subtotal:    84000,
			OngkosKirim: 18000,
			TotalAkhir:  50000, // Seharusnya 102.000 tapi dimanipulasi jadi 50.000!
		},
	}

	hasil := ValidasiOrder(req)

	if hasil.DataValid {
		t.Fatal("Ekspektasi manipulasi total akhir ditolak oleh server!")
	}
	if hasil.Integritas.IntegritasTotalValid {
		t.Error("Ekspektasi IntegritasTotalValid false karena total pembayaran tidak cocok")
	}
}

// TestHandlerValidasiRingkasan_EndpointHTTP menguji endpoint POST /api/checkout/validasi-ringkasan
func TestHandlerValidasiRingkasan_EndpointHTTP(t *testing.T) {
	handler := &HandlerValidatorOrder{DB: nil}

	payload := PermintaanValidasiOrder{
		Kontak: KontakPemesan{Nama: "Pelanggan Tes", NomorWhatsApp: "08123456789"},
		Kurir:  KurirPilihan{IDKurir: "jne", TarifOngkir: 18000, TerverifikasiAPI: true},
		RingkasanKlien: RingkasanKlien{Subtotal: 84000, OngkosKirim: 18000, TotalAkhir: 102000},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout/validasi-ringkasan", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.HandlerValidasiRingkasan(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Ekspektasi HTTP 200 OK, didapat: %d", rec.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}
	if res["status"] != "sukses" {
		t.Errorf("Ekspektasi status respons sukses, didapat: %v", res["status"])
	}
}

// TestHandlerBuatPesanan_EndpointHTTP menguji endpoint POST /api/checkout/buat-pesanan
func TestHandlerBuatPesanan_EndpointHTTP(t *testing.T) {
	handler := &HandlerValidatorOrder{DB: nil}

	// 1. Uji Kasus Sukses Buat Pesanan
	payloadValid := PermintaanBuatPesanan{
		TokenValidasi: "token-dummy",
		OrderData: PermintaanValidasiOrder{
			TipeCheckout: "guest",
			Kontak:       KontakPemesan{Nama: "Pelanggan Sukses", NomorWhatsApp: "08123456789"},
			Kurir:        KurirPilihan{IDKurir: "jne", NamaKurir: "JNE REG", KodeLayanan: "REG", TarifOngkir: 18000, TerverifikasiAPI: true},
			Items: []ItemOrder{
				{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
				{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
			},
			RingkasanKlien: RingkasanKlien{Subtotal: 84000, OngkosKirim: 18000, TotalAkhir: 102000},
		},
	}
	bodyValid, _ := json.Marshal(payloadValid)

	reqSukses := httptest.NewRequest(http.MethodPost, "/api/checkout/buat-pesanan", bytes.NewReader(bodyValid))
	recSukses := httptest.NewRecorder()

	handler.HandlerBuatPesanan(recSukses, reqSukses)

	if recSukses.Code != http.StatusCreated {
		t.Fatalf("Ekspektasi HTTP 201 Created saat validasi lolos, didapat: %d. Body: %s", recSukses.Code, recSukses.Body.String())
	}

	var respSukses ResponsBuatPesanan
	_ = json.NewDecoder(recSukses.Body).Decode(&respSukses)
	if respSukses.NomorPesanan == "" {
		t.Error("Ekspektasi nomor pesanan terbuat, namun kosong")
	}

	// 2. Uji Kasus Gagal Buat Pesanan (Tarif Ongkir Belum Valid)
	payloadInvalid := PermintaanBuatPesanan{
		OrderData: PermintaanValidasiOrder{
			SimulasiBelumValid: true, // Simulasi Gambar 2
		},
	}
	bodyInvalid, _ := json.Marshal(payloadInvalid)

	reqGagal := httptest.NewRequest(http.MethodPost, "/api/checkout/buat-pesanan", bytes.NewReader(bodyInvalid))
	recGagal := httptest.NewRecorder()

	handler.HandlerBuatPesanan(recGagal, reqGagal)

	if recGagal.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Ekspektasi HTTP 422 Unprocessable Entity saat data invalid, didapat: %d", recGagal.Code)
	}
}
