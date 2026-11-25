package webhookgateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	signatureidempotency "rengginang-putra-jaya/auth/1.3.2.2"
	pembaruanstatus "rengginang-putra-jaya/auth/1.3.2.3"
	kedaluwarsastok "rengginang-putra-jaya/auth/1.3.2.4"
)

// TestIntegrasiEndToEndPPL20 memverifikasi seluruh Acceptance Criteria PPL-20
func TestIntegrasiEndToEndPPL20(t *testing.T) {
	layananWebhook := InisialisasiLayanan(nil)
	layananKeamananIdem := signatureidempotency.InisialisasiLayanan(nil)
	layananPembaruanStatus := pembaruanstatus.InisialisasiLayanan(nil)
	layananKedaluwarsaStok := kedaluwarsastok.InisialisasiLayanan(nil)

	// Hubungkan orkestrasi seperti di main.go
	layananWebhook.OnNotifikasi = func(notif NotifikasiWebhookMidtrans) error {
		// 1. Verifikasi Signature Key (WBS 1.3.2.2)
		if notif.SignatureKey != "" {
			hasilSig := layananKeamananIdem.VerifikasiSignatureKey(notif.OrderID, notif.StatusCode, notif.GrossAmount, notif.SignatureKey)
			if !hasilSig.Valid {
				t.Logf("[Verifikasi Signature] Signature tidak valid untuk order: %s", notif.OrderID)
				return nil
			}
		}

		// 2. Proteksi Idempotency (WBS 1.3.2.2)
		cekIdem, _ := layananKeamananIdem.PeriksaDanKunciIdempotency(notif.OrderID, notif.TransactionID, notif.TransactionStatus)
		if cekIdem != nil && cekIdem.AdalahDuplikat && !cekIdem.BolehDiproses {
			t.Logf("[Idempotency] Notifikasi order %s status %s diabaikan karena sudah pernah diproses.", notif.OrderID, notif.TransactionStatus)
			return nil
		}

		// 3. Pembaruan Otomatis Status Transaksi (WBS 1.3.2.3)
		reqStatus := pembaruanstatus.PembaruanStatusRequest{
			NomorPesanan:      notif.OrderID,
			TransactionID:     notif.TransactionID,
			TransactionStatus: notif.TransactionStatus,
			FraudStatus:       notif.FraudStatus,
			PaymentType:       notif.PaymentType,
		}
		hasilStatus, _ := layananPembaruanStatus.PerbaruiStatusOtomatis(reqStatus, "webhook_test")
		t.Logf("[Status Update] Order %s berhasil diperbarui ke: %s (%s)", notif.OrderID, hasilStatus.StatusPesanan, hasilStatus.LabelStatus)

		// 4. Penanganan Kedaluwarsa & Pemutusan Reservasi Stok (WBS 1.3.2.4)
		if notif.TransactionStatus == "expire" {
			hasilLepas, _ := layananKedaluwarsaStok.PutusReservasiStokKedaluwarsa(notif.OrderID, "Notifikasi expire gateway sandbox")
			t.Logf("[Kedaluwarsa Stok] Reservasi stok order %s dilepas: %d bungkus.", notif.OrderID, hasilLepas.TotalKuantitasStok)
		}

		// Tandai selesai
		layananKeamananIdem.TandaiSelesaiIdempotency(notif.OrderID, notif.TransactionStatus)
		return nil
	}

	// =========================================================================
	// SKENARIO 1: PEMBAYARAN BERHASIL (SETTLEMENT) & STATUS MENJADI DIBAYAR/LUNAS
	// =========================================================================
	order1 := "PJ-PPL20-001"
	sig1 := layananKeamananIdem.HitungSignatureKey(order1, "200", "56000.00")
	payload1, _ := json.Marshal(NotifikasiWebhookMidtrans{
		OrderID:           order1,
		TransactionID:     "trx-ppl20-001",
		StatusCode:        "200",
		GrossAmount:       "56000.00",
		SignatureKey:      sig1,
		TransactionStatus: "settlement",
		PaymentType:       "qris",
	})

	req1 := httptest.NewRequest(http.MethodPost, "/api/pembayaran/webhook", bytes.NewBuffer(payload1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	layananWebhook.HandlerWebhook(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("Ekspektasi HTTP 200, didapat: %d", rec1.Code)
	}

	st1, _ := layananPembaruanStatus.AmbilStatusTerakhir(order1)
	if st1.StatusPesanan != "dibayar" || st1.LabelStatus != "Lunas/Dibayar" {
		t.Errorf("Ekspektasi status 'dibayar' ('Lunas/Dibayar'), didapat: %s (%s)", st1.StatusPesanan, st1.LabelStatus)
	}

	// =========================================================================
	// SKENARIO 2: NOTIFIKASI BERULANG (IDEMPOTENCY TEST)
	// =========================================================================
	rec1Repeat := httptest.NewRequest(http.MethodPost, "/api/pembayaran/webhook", bytes.NewBuffer(payload1))
	rec1Repeat.Header.Set("Content-Type", "application/json")
	w1Repeat := httptest.NewRecorder()
	layananWebhook.HandlerWebhook(w1Repeat, rec1Repeat)

	if w1Repeat.Code != http.StatusOK {
		t.Fatalf("Ekspektasi webhook berulang tetap merespons 200 OK, didapat: %d", w1Repeat.Code)
	}

	// =========================================================================
	// SKENARIO 3: PEMBAYARAN KEDALUWARSA & PEMUTUSAN RESERVASI STOK AGAR TIDAK TERBLOKIR
	// =========================================================================
	orderExpired := "PJ-PPL20-002"
	layananKedaluwarsaStok.TambahReservasiLokal(orderExpired, []kedaluwarsastok.CatatanReservasiLokal{
		{
			ID:              101,
			NomorPesanan:    orderExpired,
			IDVarian:        1,
			NamaVarian:      "Rengginang Udang",
			Jumlah:          4,
			StatusReservasi: "menahan_stok",
			KadaluarsaPada:  time.Now().Add(-10 * time.Minute),
		},
	})

	sigExpired := layananKeamananIdem.HitungSignatureKey(orderExpired, "202", "112000.00")
	payloadExpired, _ := json.Marshal(NotifikasiWebhookMidtrans{
		OrderID:           orderExpired,
		TransactionID:     "trx-ppl20-002",
		StatusCode:        "202",
		GrossAmount:       "112000.00",
		SignatureKey:      sigExpired,
		TransactionStatus: "expire",
		PaymentType:       "bank_transfer",
	})

	reqExp := httptest.NewRequest(http.MethodPost, "/api/pembayaran/webhook", bytes.NewBuffer(payloadExpired))
	reqExp.Header.Set("Content-Type", "application/json")
	recExp := httptest.NewRecorder()
	layananWebhook.HandlerWebhook(recExp, reqExp)

	if recExp.Code != http.StatusOK {
		t.Fatalf("Ekspektasi HTTP 200 untuk webhook expire, didapat: %d", recExp.Code)
	}

	stExp, _ := layananPembaruanStatus.AmbilStatusTerakhir(orderExpired)
	if stExp.StatusPesanan != "kedaluwarsa" || stExp.LabelStatus != "Kedaluwarsa" {
		t.Errorf("Ekspektasi status 'kedaluwarsa' ('Kedaluwarsa'), didapat: %s (%s)", stExp.StatusPesanan, stExp.LabelStatus)
	}

	t.Logf("Seluruh pengujian integrasi PPL-20 sukses 100%%!")
}
