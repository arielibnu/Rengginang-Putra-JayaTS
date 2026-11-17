
package validasistok

import (
	"encoding/json"
	"testing"
)

func TestNormalisasiMode(t *testing.T) {
	tests := []struct {
		input string
		want  string
		valid bool
	}{
		{"", "RIIL", true},
		{"riil", "RIIL", true},
		{" simulasi ", "SIMULASI", true},
		{"lainnya", "", false},
	}

	for _, tt := range tests {
		got, err := NormalisasiMode(tt.input)

		if (err == nil) != tt.valid {
			t.Errorf("NormalisasiMode(%q): status valid tidak sesuai", tt.input)
			continue
		}

		if err == nil && got != tt.want {
			t.Errorf("hasil %q, ingin %q", got, tt.want)
		}
	}
}

func TestNormalisasiJenis(t *testing.T) {
	tests := []string{
		"MASUK",
		"keluar",
		"penyesuaian",
	}

	for _, jenis := range tests {
		if _, err := NormalisasiJenis(jenis); err != nil {
			t.Errorf("Jenis %s seharusnya valid", jenis)
		}
	}

	if _, err := NormalisasiJenis("HAPUS"); err == nil {
		t.Error("Jenis mutasi tidak valid seharusnya ditolak")
	}
}

func TestHitungPerubahan(t *testing.T) {
	stok50 := uint64(50)
	stok42 := uint64(42)
	stokMelebihi := MaksimumStok + 1

	tests := []struct {
		nama       string
		jenis      string
		stokAwal   uint64
		jumlah     uint64
		stokBaru   *uint64
		wantJumlah uint64
		wantAkhir  uint64
		harusError bool
	}{
		{"Stok masuk", "MASUK", 40, 5, nil, 5, 45, false},
		{"Stok keluar", "KELUAR", 40, 3, nil, 3, 37, false},
		{"Stok keluar berlebihan", "KELUAR", 40, 41, nil, 0, 0, true},
		{"Stok keluar nol", "KELUAR", 40, 0, nil, 0, 0, true},
		{"Stok masuk nol", "MASUK", 40, 0, nil, 0, 0, true},
		{"Stok masuk overflow", "MASUK", MaksimumStok, 1, nil, 0, 0, true},
		{"Penyesuaian naik", "PENYESUAIAN", 42, 0, &stok50, 8, 50, false},
		{"Penyesuaian turun", "PENYESUAIAN", 42, 0, &stokMelebihi, 0, 0, true},
		{"Penyesuaian ke nol", "PENYESUAIAN", 42, 0, func() *uint64 { v := uint64(0); return &v }(), 42, 0, false},
		{"Penyesuaian tanpa nilai", "PENYESUAIAN", 42, 0, nil, 0, 0, true},
		{"Penyesuaian sama", "PENYESUAIAN", 42, 0, &stok42, 0, 0, true},
		{"Stok awal tidak valid", "MASUK", MaksimumStok + 1, 1, nil, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			jumlah, akhir, err := HitungPerubahan(
				tt.jenis,
				tt.stokAwal,
				tt.jumlah,
				tt.stokBaru,
			)

			if (err != nil) != tt.harusError {
				t.Fatalf("status error tidak sesuai: %v", err)
			}

			if err != nil {
				return
			}

			if jumlah != tt.wantJumlah {
				t.Errorf("jumlah = %d, ingin %d", jumlah, tt.wantJumlah)
			}

			if akhir != tt.wantAkhir {
				t.Errorf("stok akhir = %d, ingin %d", akhir, tt.wantAkhir)
			}
		})
	}
}

func TestJSONMenolakStokNegatif(t *testing.T) {
	var input struct {
		Jumlah uint64 `json:"jumlah"`
	}

	err := json.Unmarshal(
		[]byte(`{"jumlah":-5}`),
		&input,
	)

	if err == nil {
		t.Fatal("Stok negatif seharusnya ditolak")
	}
}
