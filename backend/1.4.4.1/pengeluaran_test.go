package laporanpengeluaran

import (
	"testing"
)

func TestKlasifikasiKategoriBeban(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Listrik dan Air Pabrik", "OPERASIONAL"},
		{"Gas Elpiji 12kg", "OPERASIONAL"},
		{"Gaji Karyawan Harian", "OPERASIONAL"},
		{"Beras Ketan Putih Super", "BELANJA_PERSEDIAAN"},
		{"Plastik Kemasan Standing Pouch", "BELANJA_PERSEDIAAN"},
		{"Bumbu Bawang Putih & Garam", "BELANJA_PERSEDIAAN"},
		{"Bahan Baku Udang Rebon", "BELANJA_PERSEDIAAN"},
		{"Servis Mesin Penjemur", "BIAYA_LAIN"},
		{"Biaya Administrasi Bank", "BIAYA_LAIN"},
		{"Lain-lain Tak Terduga", "BIAYA_LAIN"},
	}

	for _, tc := range tests {
		hasil := KlasifikasikanKategoriBeban(tc.input)
		if hasil != tc.expected {
			t.Errorf("KlasifikasikanKategoriBeban('%s') = %s; ingin %s", tc.input, hasil, tc.expected)
		}
	}
}

func TestHitungRekapPengeluaran(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	rekap, err := layanan.HitungRekapPengeluaran("November", "2026")
	if err != nil {
		t.Fatalf("Gagal menghitung rekap pengeluaran: %v", err)
	}

	if rekap.Periode.BulanNama != "November" || rekap.Periode.Tahun != 2026 {
		t.Errorf("Periode tidak sesuai: %+v", rekap.Periode)
	}

	if rekap.Operasional <= 0 {
		t.Errorf("Total operasional harus lebih besar dari 0, didapat: %f", rekap.Operasional)
	}

	if rekap.BelanjaPersediaan <= 0 {
		t.Errorf("Total belanja persediaan harus lebih besar dari 0, didapat: %f", rekap.BelanjaPersediaan)
	}

	totalManual := rekap.Operasional + rekap.BelanjaPersediaan + rekap.BiayaLain
	if rekap.TotalPengeluaran != totalManual {
		t.Errorf("Total pengeluaran tidak cocok: %f != %f", rekap.TotalPengeluaran, totalManual)
	}
}
