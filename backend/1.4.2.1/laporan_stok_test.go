package laporanbulanan

import (
	"testing"
)

func TestParsingPeriode(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// Test 1: Nama bulan bahasa Indonesia
	p1, err1 := layanan.ParsingPeriode("November", "2026")
	if err1 != nil {
		t.Fatalf("Gagal parsing periode: %v", err1)
	}
	if p1.BulanAngka != 11 || p1.Tahun != 2026 {
		t.Errorf("Bulan harus 11 dan Tahun 2026, didapat: %d / %d", p1.BulanAngka, p1.Tahun)
	}

	// Test 2: Angka bulan
	p2, err2 := layanan.ParsingPeriode("12", "2025")
	if err2 != nil {
		t.Fatalf("Gagal parsing periode angka: %v", err2)
	}
	if p2.BulanNama != "Desember" || p2.Tahun != 2025 {
		t.Errorf("Bulan harus Desember dan Tahun 2025, didapat: %s / %d", p2.BulanNama, p2.Tahun)
	}
}

func TestHitungSisaStokAkhir(t *testing.T) {
	layanan := InisialisasiLayanan(nil)
	periode, _ := layanan.ParsingPeriode("November", "2026")

	hasil, err := layanan.HitungSisaStokAkhir(periode)
	if err != nil {
		t.Fatalf("Gagal hitung sisa stok akhir: %v", err)
	}

	if hasil.StokFisik <= 0 {
		t.Errorf("Stok fisik harus > 0, didapat: %d", hasil.StokFisik)
	}

	// Pastikan stok tersedia = stok fisik - reservasi
	if hasil.StokTersedia != hasil.StokFisik-hasil.ReservasiAktif {
		t.Errorf("Stok tersedia (%d) tidak sesuai dengan stok fisik (%d) - reservasi (%d)",
			hasil.StokTersedia, hasil.StokFisik, hasil.ReservasiAktif)
	}

	// Pastikan sumber data memisahkan sandbox
	if hasil.SumberData.Sandbox != "Tidak dihitung sebagai kas nyata" {
		t.Errorf("Label sandbox tidak sesuai: %s", hasil.SumberData.Sandbox)
	}
}
