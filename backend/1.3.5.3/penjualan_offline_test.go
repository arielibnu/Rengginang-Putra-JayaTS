package penjualanoffline

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInisialisasiLayanan(t *testing.T) {
	layanan := InisialisasiLayanan(nil)
	if layanan == nil {
		t.Fatal("Layanan tidak boleh nil")
	}
}

func TestSimpanPenjualanOfflineValid(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanPenjualanOffline{
		TanggalPenjualan: "20 November 2026",
		IDVarian:         1,
		NamaVarian:       "Rengginang Udang",
		Jumlah:           2,
		HargaSatuan:      28000,
		MetodePembayaran: "Tunai",
		FlagLingkungan:   "PRODUKSI_NYATA",
	}

	hasil, err := layanan.SimpanPenjualan(req)
	if err != nil {
		t.Fatalf("Gagal simpan penjualan: %v", err)
	}

	if hasil.NomorTransaksi == "" {
		t.Error("Nomor transaksi tidak boleh kosong")
	}

	if hasil.TotalHarga != 56000 {
		t.Errorf("Ekspektasi total harga 56000, didapat: %d", hasil.TotalHarga)
	}

	if hasil.FlagLingkungan != "PRODUKSI_NYATA" {
		t.Errorf("Ekspektasi flag PRODUKSI_NYATA, didapat: %s", hasil.FlagLingkungan)
	}

	if hasil.IsSandbox != false {
		t.Error("Ekspektasi is_sandbox false")
	}
}

func TestSimpanPenjualanOfflineJumlahNol(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanPenjualanOffline{
		NamaVarian: "Rengginang Bawang",
		Jumlah:     0,
	}

	_, err := layanan.SimpanPenjualan(req)
	if err == nil {
		t.Error("Ekspektasi error saat jumlah pembelian 0")
	}
}

func TestHandlerSimpanPenjualanHTTP(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := PermintaanPenjualanOffline{
		TanggalPenjualan: "20 November 2026",
		IDVarian:         2,
		NamaVarian:       "Rengginang Bawang",
		Jumlah:           3,
		HargaSatuan:      28000,
		MetodePembayaran: "Tunai",
		FlagLingkungan:   "PRODUKSI_NYATA",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/penjualan-offline", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	layanan.HandlerSimpanPenjualan(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ekspektasi status 201 Created, didapat: %d", rr.Code)
	}

	var resp ResponsPenjualanOffline
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Gagal parse respons JSON: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Ekspektasi status sukses, didapat: %s", resp.Status)
	}
}
