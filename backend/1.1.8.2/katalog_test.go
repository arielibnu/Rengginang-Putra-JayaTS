package katalogapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ==========================================
// TEST LOGIKA STATUS KETERSEDIAAN STOK
// ==========================================

func TestTentukanStatusKetersediaan(t *testing.T) {
	tests := []struct {
		nama       string
		stok       uint64
		ekspektasi string
	}{
		{"Stok habis (nol)", 0, "Habis"},
		{"Stok satu (menipis)", 1, "Stok Menipis"},
		{"Stok lima (batas menipis)", 5, "Stok Menipis"},
		{"Stok enam (tersedia normal)", 6, "Tersedia • contoh"},
		{"Stok banyak (tersedia normal)", 40, "Tersedia • contoh"},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			hasil := TentukanStatusKetersediaan(tt.stok)
			if hasil != tt.ekspektasi {
				t.Errorf("Stok %d: hasil %q, ingin %q", tt.stok, hasil, tt.ekspektasi)
			}
		})
	}
}

// ==========================================
// TEST MAPPING GAMBAR VARIAN
// ==========================================

func TestCariGambarVarian(t *testing.T) {
	tests := []struct {
		namaVarian     string
		harusMemuat    string
		harusVertikal  string
	}{
		{"Rengginang Bawang", "rengginang-bawang.jpg", "rengginang-bawang-vertical.jpg"},
		{"Rengginang Udang", "rengginang-udang.jpg", "rengginang-udang-vertical.jpg"},
		{"Rengginang Cumi", "rengginang-cumi.jpg", "rengginang-cumi-vertical.jpg"},
		{"Rengginang Terasi", "rengginang-terasi.jpg", "rengginang-terasi-vertical.jpg"},
		{"Rasa Lain Tidak Dikenal", "rengginang-all-hero.jpg", "rengginang-all-hero.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.namaVarian, func(t *testing.T) {
			gambar, gambarDetail := CariGambarVarian(tt.namaVarian)
			if !strings.Contains(gambar, tt.harusMemuat) {
				t.Errorf("Gambar %q tidak memuat %q", gambar, tt.harusMemuat)
			}
			if !strings.Contains(gambarDetail, tt.harusVertikal) {
				t.Errorf("Gambar detail %q tidak memuat %q", gambarDetail, tt.harusVertikal)
			}
		})
	}
}

// ==========================================
// TEST FORMAT RUPIAH
// ==========================================

func TestFormatRupiah(t *testing.T) {
	tests := []struct {
		angka      float64
		ekspektasi string
	}{
		{28000, "Rp28.000"},
		{18000, "Rp18.000"},
		{500, "Rp500"},
		{1500000, "Rp1.500.000"},
	}

	for _, tt := range tests {
		hasil := FormatRupiah(tt.angka)
		if hasil != tt.ekspektasi {
			t.Errorf("FormatRupiah(%.0f) = %q, ingin %q", tt.angka, hasil, tt.ekspektasi)
		}
	}
}

// ==========================================
// TEST ENDPOINT KATALOG PUBLIK
// ==========================================

func TestEndpointKatalogPublik(t *testing.T) {
	handler := &HandlerKatalog{DB: nil} // nil DB menguji fallback data resmi

	req := httptest.NewRequest(http.MethodGet, "/api/katalog/produk", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, ingin %d", res.StatusCode, http.StatusOK)
	}

	var respons ResponsKatalog
	if err := json.NewDecoder(res.Body).Decode(&respons); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if respons.Total != 4 {
		t.Errorf("Total produk = %d, ingin 4", respons.Total)
	}
}

// ==========================================
// TEST METODE SELAIN GET DITOLAK
// ==========================================

func TestMetodeSelainGetDitolak(t *testing.T) {
	handler := &HandlerKatalog{DB: nil}

	metodeDilarang := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, m := range metodeDilarang {
		req := httptest.NewRequest(m, "/api/katalog", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Metode %s: code = %d, ingin 405", m, w.Code)
		}
	}
}
