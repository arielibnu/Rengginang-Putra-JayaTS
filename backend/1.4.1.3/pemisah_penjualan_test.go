package pemisahpenjualan

import (
	"testing"
)

func TestPemisahanTransaksiRiilDariSandbox(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	hasil, err := layanan.PisahkanTransaksiRiilDariSandbox(11, 2026)
	if err != nil {
		t.Fatalf("Error memisahkan transaksi: %v", err)
	}

	if hasil.PenjualanOfflineRiil <= 0 {
		t.Errorf("Penjualan offline riil harus > 0, didapat: %f", hasil.PenjualanOfflineRiil)
	}

	if hasil.SimulasiSandbox <= 0 {
		t.Errorf("Simulasi sandbox harus > 0, didapat: %f", hasil.SimulasiSandbox)
	}

	if hasil.TotalOmzetRiil != hasil.PenjualanOfflineRiil {
		t.Errorf("Total omzet riil tidak boleh mencampurkan simulasi online! Riil: %f, Total: %f", hasil.PenjualanOfflineRiil, hasil.TotalOmzetRiil)
	}
}
