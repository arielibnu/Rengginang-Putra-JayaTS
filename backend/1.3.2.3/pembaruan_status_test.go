package pembaruanstatus

import (
	"testing"
)

func TestPetakanStatusGateway(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	tests := []struct {
		trxStatus    string
		fraudStatus  string
		expectedKode string
		expectedText string
	}{
		{"settlement", "accept", StatusLunasDibayar, "Lunas/Dibayar"},
		{"capture", "accept", StatusLunasDibayar, "Lunas/Dibayar"},
		{"capture", "challenge", StatusMenungguBayar, "Menunggu"},
		{"pending", "", StatusMenungguBayar, "Menunggu"},
		{"expire", "", StatusKedaluwarsa, "Kedaluwarsa"},
		{"cancel", "", StatusGagal, "Gagal"},
		{"deny", "", StatusGagal, "Gagal"},
		{"failure", "", StatusGagal, "Gagal"},
	}

	for _, tt := range tests {
		kode, label := layanan.PetakanStatusGateway(tt.trxStatus, tt.fraudStatus)
		if kode != tt.expectedKode {
			t.Errorf("TrxStatus '%s': ekspektasi kode '%s', didapat '%s'", tt.trxStatus, tt.expectedKode, kode)
		}
		if label != tt.expectedText {
			t.Errorf("TrxStatus '%s': ekspektasi label '%s', didapat '%s'", tt.trxStatus, tt.expectedText, label)
		}
	}
}

func TestPerbaruiStatusOtomatis(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// 1. Update ke Settlement -> Dibayar / Lunas
	req := PembaruanStatusRequest{
		NomorPesanan:      "PJ-20261125-0010",
		TransactionID:     "trx-0010",
		TransactionStatus: "settlement",
	}

	hasil, err := layanan.PerbaruiStatusOtomatis(req, "test_webhook")
	if err != nil {
		t.Fatalf("Gagal perbarui status otomatis: %v", err)
	}

	if hasil.StatusPesanan != StatusLunasDibayar {
		t.Errorf("Ekspektasi status pesanan '%s', didapat '%s'", StatusLunasDibayar, hasil.StatusPesanan)
	}
	if hasil.LabelStatus != "Lunas/Dibayar" {
		t.Errorf("Ekspektasi label 'Lunas/Dibayar', didapat '%s'", hasil.LabelStatus)
	}

	// 2. Cek riwayat
	statusTerkini, ada := layanan.AmbilStatusTerakhir("PJ-20261125-0010")
	if !ada {
		t.Fatalf("Status nomor pesanan tidak ditemukan di memori")
	}
	if statusTerkini.StatusPesanan != StatusLunasDibayar {
		t.Errorf("Ekspektasi status di memori '%s', didapat '%s'", StatusLunasDibayar, statusTerkini.StatusPesanan)
	}
}
