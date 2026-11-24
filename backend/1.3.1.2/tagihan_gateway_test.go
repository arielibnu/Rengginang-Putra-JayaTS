package tagihan

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuatTagihanQRIS(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanBuatTagihan{
		NomorPesanan: "PJ-TEST-0001",
		Metode:       "qris",
	}

	resp, err := layanan.BuatTagihanMidtrans(req)
	if err != nil {
		t.Fatalf("Gagal membuat tagihan QRIS: %v", err)
	}

	if resp.PaymentType != "qris" {
		t.Errorf("Ekspektasi payment_type 'qris', didapat: %s", resp.PaymentType)
	}

	if resp.GrossAmount <= 0 {
		t.Errorf("Gross amount harus > 0, didapat: %d", resp.GrossAmount)
	}

	if resp.TransactionStatus != "pending" {
		t.Errorf("Ekspektasi transaction_status 'pending', didapat: %s", resp.TransactionStatus)
	}
}

func TestBuatTagihanVirtualAccount(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := PermintaanBuatTagihan{
		NomorPesanan: "PJ-TEST-0002",
		Metode:       "va",
		Bank:         "bca",
	}

	resp, err := layanan.BuatTagihanMidtrans(req)
	if err != nil {
		t.Fatalf("Gagal membuat tagihan VA: %v", err)
	}

	if resp.PaymentType != "bank_transfer" {
		t.Errorf("Ekspektasi payment_type 'bank_transfer', didapat: %s", resp.PaymentType)
	}

	if resp.VANumber == "" {
		t.Error("Nomor VA tidak boleh kosong")
	}

	if resp.Bank != "BCA" {
		t.Errorf("Ekspektasi bank BCA, didapat: %s", resp.Bank)
	}
}

func TestHandlerBuatTagihanEndpoint(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := `{"nomor_pesanan":"PJ-DEMO-0001","metode":"qris"}`
	req := httptest.NewRequest(http.MethodPost, "/api/pembayaran/buat-tagihan", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	layanan.HandlerBuatTagihan(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Ekspektasi HTTP 200, didapat: %d", res.StatusCode)
	}

	var jsonResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResp); err != nil {
		t.Fatalf("Gagal decode JSON: %v", err)
	}

	if jsonResp["status"] != "sukses" {
		t.Errorf("Ekspektasi status sukses, didapat: %v", jsonResp["status"])
	}
}
