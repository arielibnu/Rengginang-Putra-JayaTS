package laporanpenjualan

import (
	"testing"
)

func TestParsingPeriodePenjualan(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	p, err := layanan.ParsingPeriode("November", "2026")
	if err != nil {
		t.Fatalf("Gagal parsing periode: %v", err)
	}

	if p.BulanNama != "November" || p.BulanAngka != 11 || p.Tahun != 2026 {
		t.Errorf("Hasil parsing tidak sesuai: %+v", p)
	}
}

func TestHitungRekapPenjualanPemisahanSandboxDanBatal(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	laporan, err := layanan.HitungRekapPenjualan("November", "2026")
	if err != nil {
		t.Fatalf("Gagal hitung rekap penjualan: %v", err)
	}

	// Validasi WBS 1.4.1.1: Total omzet riil mengikuti penjualan offline nyata
	if laporan.PenjualanOfflineNyata <= 0 {
		t.Errorf("Penjualan offline nyata harus > 0, didapat: %f", laporan.PenjualanOfflineNyata)
	}

	// Validasi WBS 1.4.1.3: Transaksi simulasi online terpisah dari kas nyata
	if laporan.TotalOmzetRiil != laporan.PenjualanOfflineNyata {
		t.Errorf("Total omzet riil tidak boleh mencampurkan simulasi online! Riil: %f, Total: %f", laporan.PenjualanOfflineNyata, laporan.TotalOmzetRiil)
	}

	// Validasi WBS 1.4.1.3: Pesanan batal tercatat terpisah
	if laporan.PesananBatal < 0 {
		t.Errorf("Pesanan batal tidak valid: %d", laporan.PesananBatal)
	}
}
