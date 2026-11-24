package sinkronisasi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSinkronisasiStatusPesanan(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	hasil, err := layanan.SinkronisasiStatusPesanan("PJ-TEST-0001")
	if err != nil {
		t.Fatalf("Gagal sinkronisasi status: %v", err)
	}

	if hasil.NomorPesanan != "PJ-TEST-0001" {
		t.Errorf("Ekspektasi nomor pesanan PJ-TEST-0001, didapat: %s", hasil.NomorPesanan)
	}

	if hasil.StatusPesananLabel != "Menunggu bayar" {
		t.Errorf("Ekspektasi status label 'Menunggu bayar', didapat: %s", hasil.StatusPesananLabel)
	}
}

func TestProsesWebhookNotifikasi(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	notif := MidtransNotification{
		OrderID:           "PJ-TEST-0001",
		TransactionID:     "trx-midtrans-123",
		TransactionStatus: "settlement",
		StatusCode:        "200",
		GrossAmount:       "120000.00",
	}

	if err := layanan.ProsesWebhookNotifikasi(notif); err != nil {
		t.Fatalf("Gagal proses webhook: %v", err)
	}

	// Cek status setelah webhook
	hasil, err := layanan.SinkronisasiStatusPesanan("PJ-TEST-0001")
	if err != nil {
		t.Fatalf("Gagal ambil status: %v", err)
	}

	if hasil.StatusPesanan != "dibayar" {
		t.Errorf("Ekspektasi status pesanan 'dibayar', didapat: %s", hasil.StatusPesanan)
	}

	if hasil.StatusPesananLabel != "Sudah dibayar" {
		t.Errorf("Ekspektasi label 'Sudah dibayar', didapat: %s", hasil.StatusPesananLabel)
	}
}

func TestHandlerSimulasiBayar(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := `{"nomor_pesanan":"PJ-TEST-0003","status":"settlement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/pembayaran/simulasi-bayar", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	layanan.HandlerSimulasiBayar(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Ekspektasi status 200, didapat: %d", res.StatusCode)
	}
}
