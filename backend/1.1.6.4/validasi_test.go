package validasibiaya

import (
	"testing"
)

// ==========================================
// TEST KATEGORI PENGELUARAN (ANTI DOUBLE-COUNTING)
// ==========================================

func TestCekKategoriPengeluaran_KategoriHPPDitolak(t *testing.T) {
	// Kategori yang termasuk komponen modal/HPP
	// harus ditolak supaya tidak dihitung dua kali.
	kategoriHPP := []string{
		"Bahan Baku",
		"bahan mentah",
		"Ketan",
		"Beras Ketan",
		"Udang segar",
		"Cumi kering",
		"Terasi super",
		"Bawang Putih",
		"Bumbu Produksi",
	}

	for _, kat := range kategoriHPP {
		err := CekKategoriPengeluaran(kat)

		if err == nil {
			t.Errorf(
				"Kategori %q seharusnya ditolak karena termasuk HPP",
				kat,
			)
		}
	}
}

func TestCekKategoriPengeluaran_KategoriOperasionalDiterima(t *testing.T) {
	// Kategori operasional yang sah harus diterima.
	kategoriSah := []string{
		"Listrik & Gas",
		"Plastik Kemasan",
		"Minyak Goreng",
		"Gaji Karyawan",
		"Transportasi",
		"Perawatan Alat",
		"Sewa Tempat",
	}

	for _, kat := range kategoriSah {
		err := CekKategoriPengeluaran(kat)

		if err != nil {
			t.Errorf(
				"Kategori %q seharusnya diterima: %v",
				kat, err,
			)
		}
	}
}

func TestCekKategoriPengeluaran_Kosong(t *testing.T) {
	err := CekKategoriPengeluaran("")

	if err == nil {
		t.Error("Kategori kosong seharusnya ditolak")
	}
}

// ==========================================
// TEST DUPLIKASI PENGELUARAN
// ==========================================

func TestCekDuplikasiPengeluaran_AdaDuplikat(t *testing.T) {
	daftarAda := []DataPengeluaranCek{
		{Tanggal: "2026-11-15", Kategori: "Listrik & Gas", Nominal: 750000},
		{Tanggal: "2026-11-10", Kategori: "Plastik Kemasan", Nominal: 1200000},
	}

	// Pengeluaran yang persis sama harus terdeteksi duplikat.
	err := CekDuplikasiPengeluaran(
		"2026-11-15", "Listrik & Gas", 750000, daftarAda,
	)

	if err == nil {
		t.Error("Pengeluaran duplikat seharusnya ditolak")
	}
}

func TestCekDuplikasiPengeluaran_TidakDuplikat(t *testing.T) {
	daftarAda := []DataPengeluaranCek{
		{Tanggal: "2026-11-15", Kategori: "Listrik & Gas", Nominal: 750000},
	}

	// Nominal berbeda bukan duplikat.
	err := CekDuplikasiPengeluaran(
		"2026-11-15", "Listrik & Gas", 800000, daftarAda,
	)

	if err != nil {
		t.Errorf("Bukan duplikat tapi ditolak: %v", err)
	}

	// Tanggal berbeda bukan duplikat.
	err = CekDuplikasiPengeluaran(
		"2026-11-16", "Listrik & Gas", 750000, daftarAda,
	)

	if err != nil {
		t.Errorf("Tanggal beda bukan duplikat tapi ditolak: %v", err)
	}

	// Kategori berbeda bukan duplikat.
	err = CekDuplikasiPengeluaran(
		"2026-11-15", "Air", 750000, daftarAda,
	)

	if err != nil {
		t.Errorf("Kategori beda bukan duplikat tapi ditolak: %v", err)
	}
}

func TestCekDuplikasiPengeluaran_DaftarKosong(t *testing.T) {
	err := CekDuplikasiPengeluaran(
		"2026-11-15", "Listrik", 500000, nil,
	)

	if err != nil {
		t.Errorf("Daftar kosong seharusnya tidak ada duplikat: %v", err)
	}
}

// ==========================================
// TEST KATEGORI OPERASIONAL SAH
// ==========================================

func TestApakahKategoriOperasionalSah(t *testing.T) {
	tests := []struct {
		kategori string
		sah      bool
	}{
		{"Listrik & Gas", true},
		{"plastik kemasan", true},
		{"minyak goreng", true},
		{"gaji karyawan", true},
		{"sewa tempat", true},
		{"transportasi", true},
		{"pemasaran", true},
		{"lain-lain", true},
		{"Beli laptop", false},
		{"Makan siang tim", false},
		{"Donasi", false},
	}

	for _, tt := range tests {
		hasil := ApakahKategoriOperasionalSah(tt.kategori)

		if hasil != tt.sah {
			t.Errorf(
				"Kategori %q: dapat %v, ingin %v",
				tt.kategori, hasil, tt.sah,
			)
		}
	}
}

// ==========================================
// TEST NOMINAL PENGELUARAN
// ==========================================

func TestCekNominalPengeluaran(t *testing.T) {
	tests := []struct {
		nama       string
		nominal    float64
		harusError bool
	}{
		{"Nominal valid", 750000, false},
		{"Nominal kecil valid", 1, false},
		{"Nominal nol ditolak", 0, true},
		{"Nominal negatif ditolak", -50000, true},
		{"Nominal melebihi batas", 1500000000, true},
		{"Nominal tepat batas", 1000000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			err := CekNominalPengeluaran(tt.nominal)

			if (err != nil) != tt.harusError {
				t.Errorf("Nominal %.0f: error=%v, harusError=%v",
					tt.nominal, err, tt.harusError)
			}
		})
	}
}

// ==========================================
// TEST VALIDASI LENGKAP
// ==========================================

func TestValidasiPengeluaranLengkap_Valid(t *testing.T) {
	hasil := ValidasiPengeluaranLengkap(
		"2026-11-15",
		"Listrik & Gas",
		750000,
		nil,
	)

	if !hasil.Valid {
		t.Errorf("Seharusnya valid: %s", hasil.Pesan)
	}

	if hasil.Peringatan != "" {
		t.Errorf("Seharusnya tanpa peringatan: %s", hasil.Peringatan)
	}
}

func TestValidasiPengeluaranLengkap_KategoriHPPDitolak(t *testing.T) {
	hasil := ValidasiPengeluaranLengkap(
		"2026-11-15",
		"Bahan Baku Ketan",
		500000,
		nil,
	)

	if hasil.Valid {
		t.Error("Kategori HPP seharusnya ditolak")
	}
}

func TestValidasiPengeluaranLengkap_NominalNegatifDitolak(t *testing.T) {
	hasil := ValidasiPengeluaranLengkap(
		"2026-11-15",
		"Listrik",
		-100,
		nil,
	)

	if hasil.Valid {
		t.Error("Nominal negatif seharusnya ditolak")
	}
}

func TestValidasiPengeluaranLengkap_DuplikatDitolak(t *testing.T) {
	daftarAda := []DataPengeluaranCek{
		{Tanggal: "2026-11-15", Kategori: "Listrik", Nominal: 750000},
	}

	hasil := ValidasiPengeluaranLengkap(
		"2026-11-15",
		"Listrik",
		750000,
		daftarAda,
	)

	if hasil.Valid {
		t.Error("Pengeluaran duplikat seharusnya ditolak")
	}
}

func TestValidasiPengeluaranLengkap_KategoriTidakDikenal(t *testing.T) {
	hasil := ValidasiPengeluaranLengkap(
		"2026-11-15",
		"Beli Laptop",
		5000000,
		nil,
	)

	if !hasil.Valid {
		t.Errorf("Seharusnya valid tapi dengan peringatan: %s", hasil.Pesan)
	}

	if hasil.Peringatan == "" {
		t.Error("Kategori tidak dikenal seharusnya ada peringatan")
	}
}
