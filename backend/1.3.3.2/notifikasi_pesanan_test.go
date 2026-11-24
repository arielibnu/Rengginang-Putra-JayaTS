package notifikasipesanan

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInisialisasiLayananNotifikasi(t *testing.T) {
	layanan := InisialisasiLayanan(nil)
	if layanan == nil {
		t.Fatal("Layanan tidak boleh nil")
	}

	hasil, err := layanan.AmbilPesananBaruOwner()
	if err != nil {
		t.Fatalf("Gagal mengambil pesanan baru: %v", err)
	}

	if hasil.JumlahBaru < 2 {
		t.Errorf("Ekspektasi minimal 2 pesanan baru fallback awal, didapat: %d", hasil.JumlahBaru)
	}

	if hasil.BadgeLabel == "" {
		t.Error("Badge label tidak boleh kosong")
	}
}

func TestCatatPemberitahuanPesananMasuk(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanTriggerNotifikasi{
		NomorPesanan:      "PJ-TEST-0001",
		NomorInvoice:      "INV-PJ-TEST-0001",
		TipePemberitahuan: "pesanan_masuk",
		MetodePenerimaan:  "Ambil Sendiri",
		TotalTransaksi:    56000,
	}

	notif, err := layanan.CatatPemberitahuan(req)
	if err != nil {
		t.Fatalf("Gagal mencatat pemberitahuan: %v", err)
	}

	if notif.NomorPesanan != "PJ-TEST-0001" {
		t.Errorf("Ekspektasi nomor pesanan PJ-TEST-0001, didapat: %s", notif.NomorPesanan)
	}

	if notif.MetodePenerimaan != "Ambil Sendiri" {
		t.Errorf("Ekspektasi metode Ambil Sendiri, didapat: %s", notif.MetodePenerimaan)
	}
}

func TestCatatPemberitahuanPesananTerbayar(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanTriggerNotifikasi{
		NomorPesanan:      "PJ-TEST-0002",
		NomorInvoice:      "INV-PJ-TEST-0002",
		TipePemberitahuan: "pesanan_terbayar",
		MetodePenerimaan:  "Kirim Kurir",
		TotalTransaksi:    112000,
	}

	notif, err := layanan.CatatPemberitahuan(req)
	if err != nil {
		t.Fatalf("Gagal mencatat pemberitahuan: %v", err)
	}

	if notif.TipePemberitahuan != "pesanan_terbayar" {
		t.Errorf("Ekspektasi tipe pesanan_terbayar, didapat: %s", notif.TipePemberitahuan)
	}
}

func TestTandaiSudahDibaca(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanTandaiDibaca{
		Semua: true,
	}

	err := layanan.TandaiSudahDibaca(req)
	if err != nil {
		t.Fatalf("Gagal menandai sudah dibaca: %v", err)
	}

	for _, item := range layanan.InMemori {
		if !item.SudahDibaca {
			t.Errorf("Item %s seharusnya sudah dibaca", item.NomorPesanan)
		}
	}
}

func TestHandlerHTTPPesananBaru(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/owner/notifikasi/pesanan-baru", nil)
	rr := httptest.NewRecorder()

	layanan.HandlerPesananBaru(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ekspektasi status 200, didapat: %d", rr.Code)
	}

	var resp ResponsPesananBaruOwner
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Respons JSON tidak valid: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Ekspektasi status sukses, didapat: %s", resp.Status)
	}
}

func TestHandlerHTTPTrigger(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := PermintaanTriggerNotifikasi{
		NomorPesanan:      "PJ-HTTP-0003",
		MetodePenerimaan:  "Kirim Kurir",
		TipePemberitahuan: "pesanan_masuk",
		TotalTransaksi:    84000,
	}
	bodyBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/owner/notifikasi/trigger", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	layanan.HandlerTrigger(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ekspektasi status 201 Created, didapat: %d", rr.Code)
	}
}
