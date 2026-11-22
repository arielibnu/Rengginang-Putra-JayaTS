package pesananowner

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAmbilDaftarPesanan_Semua(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	filter := FilterPesananOwner{
		Status:  "Semua status",
		Rentang: "Semua tanggal",
	}

	data, total, err := layanan.AmbilDaftarPesanan(filter)
	if err != nil {
		t.Fatalf("Ekspektasi tidak error, didapat: %v", err)
	}

	if total == 0 || len(data) == 0 {
		t.Errorf("Ekspektasi ada data pesanan default Figma, didapat 0")
	}

	// Cek salah satu nomor demo Figma
	adaDemo1 := false
	for _, p := range data {
		if p.NomorPesanan == "PJ-DEMO-0001" {
			adaDemo1 = true
			if p.Kontak != "Pelanggan A" {
				t.Errorf("Ekspektasi kontak 'Pelanggan A', didapat: %s", p.Kontak)
			}
			if p.Pembayaran != "Menunggu" {
				t.Errorf("Ekspektasi pembayaran 'Menunggu', didapat: %s", p.Pembayaran)
			}
		}
	}

	if !adaDemo1 {
		t.Errorf("Ekspektasi PJ-DEMO-0001 ada di daftar pesanan")
	}
}

func TestAmbilDaftarPesanan_FilterStatus(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// Filter status "Diproses"
	filter := FilterPesananOwner{
		Status: "Diproses",
	}

	data, _, err := layanan.AmbilDaftarPesanan(filter)
	if err != nil {
		t.Fatalf("Gagal filter status: %v", err)
	}

	for _, p := range data {
		if p.Status != "Diproses" && p.Pembayaran != "Diproses" {
			t.Errorf("Ekspektasi status 'Diproses', didapat: %s", p.Status)
		}
	}
}

func TestAmbilDaftarPesanan_PencarianKataKunci(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// Cari "Pelanggan B"
	filter := FilterPesananOwner{
		KataKunci: "Pelanggan B",
	}

	data, _, err := layanan.AmbilDaftarPesanan(filter)
	if err != nil {
		t.Fatalf("Gagal cari pesanan: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Ekspektasi menemukan hasil untuk 'Pelanggan B'")
	}

	if data[0].NomorPesanan != "PJ-DEMO-0002" {
		t.Errorf("Ekspektasi nomor 'PJ-DEMO-0002', didapat: %s", data[0].NomorPesanan)
	}
}

func TestAmbilDetailPesanan_Sukses(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	detail, err := layanan.AmbilDetailPesanan("PJ-DEMO-0002")
	if err != nil {
		t.Fatalf("Gagal ambil detail pesanan: %v", err)
	}

	if detail.NomorPesanan != "PJ-DEMO-0002" {
		t.Errorf("Ekspektasi nomor 'PJ-DEMO-0002', didapat: %s", detail.NomorPesanan)
	}

	if detail.Pembayaran != "Dibayar" {
		t.Errorf("Ekspektasi status pembayaran 'Dibayar', didapat: %s", detail.Pembayaran)
	}

	if detail.Pemrosesan != "Diproses" {
		t.Errorf("Ekspektasi status pemrosesan 'Diproses', didapat: %s", detail.Pemrosesan)
	}
}

func TestPerbaruiStatusPemrosesan(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	err := layanan.PerbaruiStatusPemrosesan("PJ-DEMO-0002", "Selesai")
	if err != nil {
		t.Fatalf("Gagal update status: %v", err)
	}

	detail, _ := layanan.AmbilDetailPesanan("PJ-DEMO-0002")
	if detail.Pemrosesan != "Selesai" {
		t.Errorf("Ekspektasi status berubah ke 'Selesai', didapat: %s", detail.Pemrosesan)
	}
}

func TestHandlerHTTP_DaftarDanDetail(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// 1. Test GET /api/admin/pesanan
	reqDaftar := httptest.NewRequest(http.MethodGet, "/api/admin/pesanan?q=DEMO", nil)
	recDaftar := httptest.NewRecorder()
	layanan.HandlerDaftarPesanan(recDaftar, reqDaftar)

	if recDaftar.Code != http.StatusOK {
		t.Errorf("Ekspektasi HTTP 200 untuk daftar pesanan, didapat: %d", recDaftar.Code)
	}

	var resDaftar ResponsDaftarPesanan
	if err := json.NewDecoder(recDaftar.Body).Decode(&resDaftar); err != nil {
		t.Fatalf("Gagal decode JSON respons daftar: %v", err)
	}

	if len(resDaftar.Data) == 0 {
		t.Errorf("Ekspektasi data pesanan tidak kosong")
	}

	// 2. Test GET /api/admin/pesanan/detail
	reqDetail := httptest.NewRequest(http.MethodGet, "/api/admin/pesanan/detail?nomor_pesanan=PJ-DEMO-0002", nil)
	recDetail := httptest.NewRecorder()
	layanan.HandlerDetailPesanan(recDetail, reqDetail)

	if recDetail.Code != http.StatusOK {
		t.Errorf("Ekspektasi HTTP 200 untuk detail pesanan, didapat: %d", recDetail.Code)
	}

	// 3. Test PUT /api/admin/pesanan/status
	payloadUpdate := PermintaanUpdateStatus{
		NomorPesanan:     "PJ-DEMO-0002",
		StatusPemrosesan: "Selesai",
	}
	body, _ := json.Marshal(payloadUpdate)
	reqUpdate := httptest.NewRequest(http.MethodPut, "/api/admin/pesanan/status", bytes.NewReader(body))
	recUpdate := httptest.NewRecorder()
	layanan.HandlerUpdateStatusPesanan(recUpdate, reqUpdate)

	if recUpdate.Code != http.StatusOK {
		t.Errorf("Ekspektasi HTTP 200 untuk update status, didapat: %d", recUpdate.Code)
	}
}
