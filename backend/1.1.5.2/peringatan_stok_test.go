package peringatanstok

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestEvaluasiStatusStok menguji kesesuaian aturan bisnis:
// "Peringatan muncul saat stok tersedia kurang dari atau sama dengan batas minimum per varian."
func TestEvaluasiStatusStok(t *testing.T) {
	tests := []struct {
		nama          string
		stokTersedia  int64
		batasMinimum  uint64
		harapanStatus string
		harapanKritis bool
	}{
		{
			nama:          "Stok di bawah batas minimum (Kasus Cumi)",
			stokTersedia:  12,
			batasMinimum:  15,
			harapanStatus: "Stok menipis",
			harapanKritis: true,
		},
		{
			nama:          "Stok tepat sama dengan batas minimum",
			stokTersedia:  15,
			batasMinimum:  15,
			harapanStatus: "Stok menipis",
			harapanKritis: true,
		},
		{
			nama:          "Stok di atas batas minimum (Kasus Udang)",
			stokTersedia:  32,
			batasMinimum:  15,
			harapanStatus: "Normal",
			harapanKritis: false,
		},
		{
			nama:          "Stok di atas batas minimum (Kasus Bawang)",
			stokTersedia:  38,
			batasMinimum:  15,
			harapanStatus: "Normal",
			harapanKritis: false,
		},
		{
			nama:          "Stok habis (0 bungkus)",
			stokTersedia:  0,
			batasMinimum:  10,
			harapanStatus: "Stok menipis",
			harapanKritis: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			status, kritis := EvaluasiStatusStok(tt.stokTersedia, tt.batasMinimum)
			if status != tt.harapanStatus {
				t.Errorf("Status salah: didapat %q, harapan %q", status, tt.harapanStatus)
			}
			if kritis != tt.harapanKritis {
				t.Errorf("Flag kritis salah: didapat %v, harapan %v", kritis, tt.harapanKritis)
			}
		})
	}
}

// TestValidasiInputBatasMinimum menguji integritas validasi parameter pengaturan batas minimum.
func TestValidasiInputBatasMinimum(t *testing.T) {
	tests := []struct {
		nama        string
		req         PermintaanAturBatasMinimum
		harusSukses bool
	}{
		{
			nama: "Data valid standar",
			req: PermintaanAturBatasMinimum{
				IDVarian:     1,
				BatasMinimum: 15,
			},
			harusSukses: true,
		},
		{
			nama: "Data valid batas 0",
			req: PermintaanAturBatasMinimum{
				IDVarian:     2,
				BatasMinimum: 0,
			},
			harusSukses: true,
		},
		{
			nama: "ID varian 0 (tidak valid)",
			req: PermintaanAturBatasMinimum{
				IDVarian:     0,
				BatasMinimum: 15,
			},
			harusSukses: false,
		},
		{
			nama: "Batas minimum melebihi maksimum wajar",
			req: PermintaanAturBatasMinimum{
				IDVarian:     3,
				BatasMinimum: 2000000,
			},
			harusSukses: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			err := ValidasiInputBatasMinimum(tt.req)
			if tt.harusSukses && err != nil {
				t.Errorf("Diharapkan sukses tapi error: %v", err)
			}
			if !tt.harusSukses && err == nil {
				t.Errorf("Diharapkan error tetapi nil")
			}
		})
	}
}

// TestSusunRingkasanPeringatan menguji pembentukan objek ringkasan peringatan otomatis.
func TestSusunRingkasanPeringatan(t *testing.T) {
	t.Run("Ada 1 varian stok menipis (Sesuai Mockup Figma)", func(t *testing.T) {
		daftar := []StatusStokVarian{
			{
				IDVarian:     1,
				NamaVarian:   "Rengginang Udang",
				StokFisik:    34,
				Reservasi:    2,
				StokTersedia: 32,
				BatasMinimum: 15,
				Status:       "Normal",
				ApakahKritis: false,
			},
			{
				IDVarian:     2,
				NamaVarian:   "Rengginang Bawang",
				StokFisik:    40,
				Reservasi:    2,
				StokTersedia: 38,
				BatasMinimum: 15,
				Status:       "Normal",
				ApakahKritis: false,
			},
			{
				IDVarian:     3,
				NamaVarian:   "Rengginang Cumi",
				StokFisik:    14,
				Reservasi:    2,
				StokTersedia: 12,
				BatasMinimum: 15,
				Status:       "Stok menipis",
				ApakahKritis: true,
			},
			{
				IDVarian:     4,
				NamaVarian:   "Rengginang Terasi",
				StokFisik:    40,
				Reservasi:    2,
				StokTersedia: 38,
				BatasMinimum: 15,
				Status:       "Normal",
				ApakahKritis: false,
			},
		}

		ringkasan := SusunRingkasanPeringatan(daftar)

		if ringkasan.TotalVarianAktif != 4 {
			t.Errorf("Total varian aktif didapat %d, harap 4", ringkasan.TotalVarianAktif)
		}
		if ringkasan.JumlahVarianMenipis != 1 {
			t.Errorf("Jumlah varian menipis didapat %d, harap 1", ringkasan.JumlahVarianMenipis)
		}
		if !ringkasan.AdaPeringatan {
			t.Errorf("AdaPeringatan harus true")
		}
		if ringkasan.PesanUtama != "Perhatian: 1 varian stok menipis" {
			t.Errorf("PesanUtama didapat %q, harap 'Perhatian: 1 varian stok menipis'", ringkasan.PesanUtama)
		}
		if !strings.Contains(ringkasan.KeteranganRinci, "Rengginang Cumi tersedia 12 bungkus") {
			t.Errorf("KeteranganRinci tidak memuat Cumi: %q", ringkasan.KeteranganRinci)
		}
	})

	t.Run("Semua varian di atas batas minimum", func(t *testing.T) {
		daftar := []StatusStokVarian{
			{
				IDVarian:     1,
				NamaVarian:   "Rengginang Udang",
				StokFisik:    34,
				Reservasi:    2,
				StokTersedia: 32,
				BatasMinimum: 15,
				Status:       "Normal",
				ApakahKritis: false,
			},
		}

		ringkasan := SusunRingkasanPeringatan(daftar)
		if ringkasan.AdaPeringatan {
			t.Errorf("AdaPeringatan harus false")
		}
		if ringkasan.JumlahVarianMenipis != 0 {
			t.Errorf("JumlahVarianMenipis harus 0")
		}
	})
}

// TestHandlerHTTPMethodNotAllowed menguji proteksi method HTTP handler.
func TestHandlerHTTPMethodNotAllowed(t *testing.T) {
	h := &Handler{DB: nil}

	// Uji method DELETE pada handler batas minimum
	reqDelete := httptest.NewRequest(http.MethodDelete, "/api/admin/stok/batas-minimum", nil)
	recDelete := httptest.NewRecorder()
	h.HandlerBatasMinimum(recDelete, reqDelete)

	if recDelete.Code != http.StatusMethodNotAllowed {
		t.Errorf("Harus 405 Method Not Allowed, didapat: %d", recDelete.Code)
	}

	// Uji method POST pada handler peringatan stok
	reqPost := httptest.NewRequest(http.MethodPost, "/api/admin/stok/peringatan", nil)
	recPost := httptest.NewRecorder()
	h.HandlerPeringatanStok(recPost, reqPost)

	if recPost.Code != http.StatusMethodNotAllowed {
		t.Errorf("Harus 405 Method Not Allowed, didapat: %d", recPost.Code)
	}
}
