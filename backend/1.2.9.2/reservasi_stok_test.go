package reservasisstok

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ordergenerator "rengginang-putra-jaya/auth/1.2.9.1"
)

func TestCekDanReservasiStok_Sukses(t *testing.T) {
	gen := &ordergenerator.GeneratorOrder{DB: nil}
	handler := &HandlerReservasiStok{DB: nil, Generator: gen}

	items := []ordergenerator.ItemSnapshotPesanan{
		{NamaProduk: "Rengginang Udang", NamaVarian: "Udang", Jumlah: 2, HargaSatuanSnapshot: 28000},
		{NamaProduk: "Rengginang Bawang", NamaVarian: "Bawang", Jumlah: 1, HargaSatuanSnapshot: 28000},
	}

	hasil, err := handler.CekDanReservasiStok("PJ-TEST-0001", items)
	if err != nil {
		t.Fatalf("Ekspektasi reservasi sukses, gagal dengan error: %v", err)
	}

	if !hasil.Sukses {
		t.Fatal("Ekspektasi hasil.Sukses true")
	}
	if hasil.HoldingTimeJam != 24 {
		t.Errorf("Ekspektasi holding time 24 jam, didapat: %d", hasil.HoldingTimeJam)
	}
	if len(hasil.RincianReservasi) != 2 {
		t.Errorf("Ekspektasi 2 rincian reservasi, didapat: %d", len(hasil.RincianReservasi))
	}
}

func TestProteksiOverbooking_StokTidakCukup(t *testing.T) {
	gen := &ordergenerator.GeneratorOrder{DB: nil}
	handler := &HandlerReservasiStok{DB: nil, Generator: gen}

	// Permintaan 999 bungkus (melebihi kuota stok tersedia)
	itemsOver := []ordergenerator.ItemSnapshotPesanan{
		{NamaProduk: "Rengginang Udang", NamaVarian: "Udang", Jumlah: 999, HargaSatuanSnapshot: 28000},
	}

	hasil, err := handler.CekDanReservasiStok("PJ-TEST-OVER", itemsOver)
	if err == nil {
		t.Fatal("Ekspektasi error overbooking saat stok tidak mencukupi, tapi tidak ada error")
	}

	if hasil.Sukses {
		t.Fatal("Ekspektasi hasil.Sukses false untuk overbooking")
	}
	if hasil.KodeHasil != "STOK_TIDAK_CUKUP" {
		t.Errorf("Ekspektasi kode 'STOK_TIDAK_CUKUP', didapat: %s", hasil.KodeHasil)
	}
}

func TestAutoRelease_StokKadaluarsa(t *testing.T) {
	gen := &ordergenerator.GeneratorOrder{DB: nil}
	handler := &HandlerReservasiStok{DB: nil, Generator: gen}

	muStok.Lock()
	memoriReservasi = append(memoriReservasi, CatatanReservasi{
		ID:                9999,
		NomorPesanan:      "PJ-EXPIRED-TEST",
		NamaVarian:        "Udang",
		JumlahDireservasi: 5,
		StatusReservasi:   "menahan_stok",
		KadaluarsaPada:    time.Now().Add(-1 * time.Hour), // Sudah kadaluarsa 1 jam lalu
	})
	muStok.Unlock()

	jumlahDilepas := handler.LepasStokKadaluarsa()
	if jumlahDilepas < 1 {
		t.Error("Ekspektasi minimal 1 stok kadaluarsa berhasil dilepas secara otomatis")
	}
}

func TestHandlerBuatPesanan_EndpointHTTP(t *testing.T) {
	gen := &ordergenerator.GeneratorOrder{DB: nil}
	handler := &HandlerReservasiStok{DB: nil, Generator: gen}

	payload := PermintaanBuatPesananPPL18{
		TipeCheckout: "guest",
		Kontak: ordergenerator.DataKontakPengiriman{
			NamaPenerima:  "Pelanggan Contoh",
			NomorWhatsApp: "081234567890",
			AlamatLengkap: "Jl. Contoh No. 12",
			KotaKabupaten: "Kabupaten Jember",
			Provinsi:      "Jawa Timur",
		},
		Kurir: ordergenerator.DataKurirSnapshot{
			IDKurir:          "jne",
			NamaKurir:        "JNE REG",
			KodeLayanan:      "REG",
			TarifOngkir:      36000,
			TerverifikasiAPI: true,
		},
		Items: []ordergenerator.ItemSnapshotPesanan{
			{NamaProduk: "Rengginang Udang", NamaVarian: "Udang", Jumlah: 2, HargaSatuanSnapshot: 28000},
			{NamaProduk: "Rengginang Bawang", NamaVarian: "Bawang", Jumlah: 1, HargaSatuanSnapshot: 28000},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/pesanan/buat", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.HandlerBuatPesanan(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Ekspektasi HTTP 201 Created, didapat: %d. Body: %s", rec.Code, rec.Body.String())
	}

	var res ResponsKonfirmasiPesananPPL18
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if !strings.HasPrefix(res.NomorPesanan, "PJ-") {
		t.Errorf("Ekspektasi nomor pesanan diawali PJ-, didapat: %s", res.NomorPesanan)
	}
	if res.StatusPesananLabel != "Menunggu pembayaran" {
		t.Errorf("Ekspektasi badge 'Menunggu pembayaran', didapat: %s", res.StatusPesananLabel)
	}
	if res.HoldingTimeJam != 24 {
		t.Errorf("Ekspektasi holding time 24 jam, didapat: %d", res.HoldingTimeJam)
	}
}

func TestHandlerBuatPesanan_SimulasiStokHabis(t *testing.T) {
	gen := &ordergenerator.GeneratorOrder{DB: nil}
	handler := &HandlerReservasiStok{DB: nil, Generator: gen}

	payload := PermintaanBuatPesananPPL18{
		SimulasiStokHabis: true,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/pesanan/buat", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.HandlerBuatPesanan(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("Ekspektasi HTTP 409 Conflict saat stok habis, didapat: %d", rec.Code)
	}
}
