package signatureidempotency

import (
	"testing"
)

func TestVerifikasiSignatureKey(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	orderID := "PJ-20261125-0001"
	statusCode := "200"
	grossAmount := "56000.00"

	// Hitung signature valid
	validSig := layanan.HitungSignatureKey(orderID, statusCode, grossAmount)

	// 1. Uji signature valid
	hasilValid := layanan.VerifikasiSignatureKey(orderID, statusCode, grossAmount, validSig)
	if !hasilValid.Valid {
		t.Errorf("Ekspektasi signature valid bernilai true, didapat: false")
	}

	// 2. Uji signature palsu / tidak valid
	hasilPalsu := layanan.VerifikasiSignatureKey(orderID, statusCode, grossAmount, "invalid_signature_hash_123456")
	if hasilPalsu.Valid {
		t.Errorf("Ekspektasi signature palsu bernilai false, didapat: true")
	}
}

func TestProteksiIdempotency(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	orderID := "PJ-20261125-0002"
	trxID := "trx-002"
	statusTrx := "settlement"

	// 1. Pemanggilan Pertama (harus boleh diproses)
	hasil1, err := layanan.PeriksaDanKunciIdempotency(orderID, trxID, statusTrx)
	if err != nil {
		t.Fatalf("Gagal cek idempotency pertama: %v", err)
	}

	if hasil1.AdalahDuplikat {
		t.Errorf("Pemanggilan pertama tidak boleh dianggap duplikat")
	}
	if !hasil1.BolehDiproses {
		t.Errorf("Pemanggilan pertama harus boleh diproses")
	}

	// Tandai selesai diproses
	layanan.TandaiSelesaiIdempotency(orderID, statusTrx)

	// 2. Pemanggilan Kedua dengan order_id dan status yang sama (notifikasi berulang)
	hasil2, err := layanan.PeriksaDanKunciIdempotency(orderID, trxID, statusTrx)
	if err != nil {
		t.Fatalf("Gagal cek idempotency kedua: %v", err)
	}

	if !hasil2.AdalahDuplikat {
		t.Errorf("Ekspektasi pemanggilan kedua terdeteksi sebagai duplikat")
	}
	if hasil2.BolehDiproses {
		t.Errorf("Ekspektasi pemanggilan kedua TIDAK boleh diproses lagi (mencegah double payment)")
	}
	if hasil2.Catatan.JumlahPercobaan != 2 {
		t.Errorf("Ekspektasi percobaan berjumlah 2, didapat: %d", hasil2.Catatan.JumlahPercobaan)
	}
}
