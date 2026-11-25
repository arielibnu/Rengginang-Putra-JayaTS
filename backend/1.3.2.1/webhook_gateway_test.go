package webhookgateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWebhookSukses(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := `{
		"transaction_time": "2026-11-25 15:00:00",
		"transaction_status": "settlement",
		"transaction_id": "trx-midtrans-test-123",
		"status_message": "midtrans payment notification",
		"status_code": "200",
		"signature_key": "abc123def456",
		"payment_type": "qris",
		"order_id": "PJ-20261125-0001",
		"gross_amount": "56000.00"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/pembayaran/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	layanan.HandlerWebhook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Ekspektasi kode HTTP 200, didapat: %d", rec.Code)
	}

	var res ResponWebhookStandar
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("Gagal parse respon JSON: %v", err)
	}

	if res.Status != "sukses" {
		t.Errorf("Ekspektasi status 'sukses', didapat: %s", res.Status)
	}

	if res.NomorPesanan != "PJ-20261125-0001" {
		t.Errorf("Ekspektasi nomor pesanan PJ-20261125-0001, didapat: %s", res.NomorPesanan)
	}

	logs := layanan.AmbilSemuaLog()
	if len(logs) != 1 {
		t.Errorf("Ekspektasi 1 log webhook tercatat, didapat: %d", len(logs))
	}
}

func TestHandlerWebhookBadPayload(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/pembayaran/webhook", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	layanan.HandlerWebhook(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Ekspektasi kode HTTP 400 untuk bad JSON, didapat: %d", rec.Code)
	}
}

func TestHandlerWebhookMethodNotAllowed(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/pembayaran/webhook", nil)
	rec := httptest.NewRecorder()

	layanan.HandlerWebhook(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ekspektasi kode HTTP 405 Method Not Allowed, didapat: %d", rec.Code)
	}
}
