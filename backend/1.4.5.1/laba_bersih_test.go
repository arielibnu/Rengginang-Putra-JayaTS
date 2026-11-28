package lababersih

import (
	"testing"
)

func TestParsingPeriodeEstimasiKeuntungan(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	periode, err := layanan.ParsePeriode("November", "2026")
	if err != nil {
		t.Fatalf("ParsePeriode gagal: %v", err)
	}

	if periode.BulanAngka != 11 {
		t.Errorf("Ekspektasi bulan 11, dapat %d", periode.BulanAngka)
	}
	if periode.Tahun != 2026 {
		t.Errorf("Ekspektasi tahun 2026, dapat %d", periode.Tahun)
	}
}

func TestHitungEstimasiLabaBersihFormula(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	periode, _ := layanan.ParsePeriode("November", "2026")
	hasil, err := layanan.HitungEstimasiKeuntungan(periode)
	if err != nil {
		t.Fatalf("HitungEstimasiKeuntungan gagal: %v", err)
	}

	// Formula WBS 1.4.5.1: Laba Bersih = Penjualan Bersih - HPP Terjual - Biaya Operasional
	ekspektasiLaba := hasil.PenerimaanNyata - hasil.ModalBarangTerjual - hasil.BiayaOperasional
	if hasil.EstimasiLabaBersih != ekspektasiLaba {
		t.Errorf("Formula laba bersih tidak cocok: dapat %.2f, ekspektasi %.2f", hasil.EstimasiLabaBersih, ekspektasiLaba)
	}

	// Verifikasi nilai tidak negatif dan dipisahkan dari sandbox
	if hasil.SimulasiSandboxDieksklusikan == 0 {
		t.Errorf("Simulasi sandbox harus dicatat terpisah dan dieksklusikan")
	}
}
