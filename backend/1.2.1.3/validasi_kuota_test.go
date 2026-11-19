package kuotastok

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// PENGUJIAN LOGIKA VALIDASI KUOTA STOK
// ============================================================================

func TestValidasiKuotaStok(t *testing.T) {
	tests := []struct {
		nama           string
		stokTersedia   int64
		jumlahDiminta  int64
		harapanValid   bool
		harapanBisaBeli bool
		harapanStatus  string
	}{
		{
			nama:           "Pembelian normal dalam batas kuota (Udang: beli 2 dari stok 12)",
			stokTersedia:   12,
			jumlahDiminta:  2,
			harapanValid:   true,
			harapanBisaBeli: true,
			harapanStatus:  "TERSEDIA",
		},
		{
			nama:           "Pembelian tepat batas maksimum sisa stok (Udang: beli 12 dari stok 12)",
			stokTersedia:   12,
			jumlahDiminta:  12,
			harapanValid:   true,
			harapanBisaBeli: true,
			harapanStatus:  "TERSEDIA",
		},
		{
			nama:           "Pembelian melebihi batas kuota stok (Udang: beli 13 dari stok 12)",
			stokTersedia:   12,
			jumlahDiminta:  13,
			harapanValid:   false,
			harapanBisaBeli: false,
			harapanStatus:  "MELEBIHI STOK",
		},
		{
			nama:           "Pembelian varian stok habis (Cumi: beli 1 dari stok 0)",
			stokTersedia:   0,
			jumlahDiminta:  1,
			harapanValid:   false,
			harapanBisaBeli: false,
			harapanStatus:  "STOK HABIS",
		},
		{
			nama:           "Pembelian varian stok negatif",
			stokTersedia:   -2,
			jumlahDiminta:  1,
			harapanValid:   false,
			harapanBisaBeli: false,
			harapanStatus:  "STOK HABIS",
		},
		{
			nama:           "Jumlah pembelian nol (0)",
			stokTersedia:   12,
			jumlahDiminta:  0,
			harapanValid:   false,
			harapanBisaBeli: false,
			harapanStatus:  "JUMLAH TIDAK VALID",
		},
		{
			nama:           "Jumlah pembelian negatif (-1)",
			stokTersedia:   12,
			jumlahDiminta:  -1,
			harapanValid:   false,
			harapanBisaBeli: false,
			harapanStatus:  "JUMLAH TIDAK VALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			valid, bisaBeli, status, pesan, ket := ValidasiKuotaStok(tt.stokTersedia, tt.jumlahDiminta)

			if valid != tt.harapanValid {
				t.Errorf("Valid: didapat %v, harapan %v", valid, tt.harapanValid)
			}
			if bisaBeli != tt.harapanBisaBeli {
				t.Errorf("BisaDibeli: didapat %v, harapan %v", bisaBeli, tt.harapanBisaBeli)
			}
			if status != tt.harapanStatus {
				t.Errorf("Status: didapat %q, harapan %q", status, tt.harapanStatus)
			}
			if pesan == "" {
				t.Errorf("Pesan tidak boleh kosong")
			}
			if ket == "" {
				t.Errorf("Keterangan tidak boleh kosong")
			}
		})
	}
}

// ============================================================================
// PENGUJIAN FORMAT RUPIAH
// ============================================================================

func TestFormatRupiah(t *testing.T) {
	tests := []struct {
		nominal  uint64
		harapan  string
	}{
		{28000, "Rp28.000"},
		{56000, "Rp56.000"},
		{140000, "Rp140.000"},
		{0, "Rp0"},
		{1250000, "Rp1.250.000"},
	}

	for _, tt := range tests {
		hasil := FormatRupiah(tt.nominal)
		if hasil != tt.harapan {
			t.Errorf("FormatRupiah(%d) = %q, harapan %q", tt.nominal, hasil, tt.harapan)
		}
	}
}

// ============================================================================
// PENGUJIAN ENDPOINT GET /api/katalog/cek-kuota
// ============================================================================

func TestEndpointCekKuota_GET(t *testing.T) {
	handler := &HandlerValidasiKuota{DB: nil}

	tests := []struct {
		nama           string
		url            string
		statusKode     int
		statusRespons  string
		bisaDibeli     bool
	}{
		{
			nama:          "GET Kuota valid udang 2 bungkus (tersedia)",
			url:           "/api/katalog/cek-kuota?id=udang&jumlah=2",
			statusKode:    http.StatusOK,
			statusRespons: "sukses",
			bisaDibeli:    true,
		},
		{
			nama:          "GET Kuota melebihi stok udang 15 bungkus (stok hanya 12)",
			url:           "/api/katalog/cek-kuota?id=udang&jumlah=15",
			statusKode:    http.StatusBadRequest,
			statusRespons: "ditolak",
			bisaDibeli:    false,
		},
		{
			nama:          "GET Kuota cumi 1 bungkus (stok habis)",
			url:           "/api/katalog/cek-kuota?id=cumi&jumlah=1",
			statusKode:    http.StatusBadRequest,
			statusRespons: "ditolak",
			bisaDibeli:    false,
		},
		{
			nama:          "GET Kuota jumlah nol",
			url:           "/api/katalog/cek-kuota?id=bawang&jumlah=0",
			statusKode:    http.StatusBadRequest,
			statusRespons: "ditolak",
			bisaDibeli:    false,
		},
		{
			nama:          "GET Varian tidak dikenal",
			url:           "/api/katalog/cek-kuota?id=keju_manis&jumlah=1",
			statusKode:    http.StatusNotFound,
			statusRespons: "gagal",
			bisaDibeli:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			handler.HandlerCekKuotaRealtime(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.statusKode {
				t.Fatalf("Status code = %d, harapan %d", res.StatusCode, tt.statusKode)
			}

			var resp ResponsAPI
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("Gagal membaca respons JSON: %v", err)
			}

			if resp.Status != tt.statusRespons {
				t.Errorf("Status respons = %q, harapan %q", resp.Status, tt.statusRespons)
			}

			if tt.statusKode == http.StatusOK && resp.Data != nil {
				if resp.Data.BisaDibeli != tt.bisaDibeli {
					t.Errorf("BisaDibeli = %v, harapan %v", resp.Data.BisaDibeli, tt.bisaDibeli)
				}
				if resp.Data.TotalHarga != uint64(resp.Data.JumlahDiminta)*resp.Data.HargaSatuan {
					t.Errorf("Perhitungan total harga salah")
				}
			}
		})
	}
}

// ============================================================================
// PENGUJIAN ENDPOINT POST /api/katalog/validasi-kuota
// ============================================================================

func TestEndpointValidasiKuota_POST(t *testing.T) {
	handler := &HandlerValidasiKuota{DB: nil}

	tests := []struct {
		nama          string
		bodyJSON      string
		statusKode    int
		statusRespons string
		valid         bool
	}{
		{
			nama:          "POST Payload frontend udang 5 bungkus",
			bodyJSON:      `{"id_varian": "udang", "jumlah": 5}`,
			statusKode:    http.StatusOK,
			statusRespons: "sukses",
			valid:         true,
		},
		{
			nama:          "POST Payload numeric ID varian 1 udang 12 bungkus",
			bodyJSON:      `{"id_varian": 1, "jumlah_beli": 12}`,
			statusKode:    http.StatusOK,
			statusRespons: "sukses",
			valid:         true,
		},
		{
			nama:          "POST Payload udang 13 bungkus (melebihi stok 12)",
			bodyJSON:      `{"id_varian": "udang", "jumlah": 13}`,
			statusKode:    http.StatusBadRequest,
			statusRespons: "ditolak",
			valid:         false,
		},
		{
			nama:          "POST Payload cumi 1 bungkus (stok habis)",
			bodyJSON:      `{"id_varian": "cumi", "jumlah": 1}`,
			statusKode:    http.StatusBadRequest,
			statusRespons: "ditolak",
			valid:         false,
		},
		{
			nama:          "POST Format JSON rusak",
			bodyJSON:      `{"id_varian": `,
			statusKode:    http.StatusBadRequest,
			statusRespons: "gagal",
			valid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/katalog/validasi-kuota", bytes.NewBufferString(tt.bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandlerCekKuotaRealtime(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.statusKode {
				t.Fatalf("Status code = %d, harapan %d", res.StatusCode, tt.statusKode)
			}

			var resp ResponsAPI
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("Gagal membaca respons JSON: %v", err)
			}

			if resp.Status != tt.statusRespons {
				t.Errorf("Status respons = %q, harapan %q", resp.Status, tt.statusRespons)
			}

			if tt.valid && (resp.Data == nil || !resp.Data.Valid) {
				t.Errorf("Data hasil validasi tidak valid padahal diharapkan valid")
			}
		})
	}
}

// ============================================================================
// PENGUJIAN METODE HTTP TIDAK DIIZINKAN
// ============================================================================

func TestEndpointCekKuota_MethodNotAllowed(t *testing.T) {
	handler := &HandlerValidasiKuota{DB: nil}
	metodeDilarang := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, m := range metodeDilarang {
		req := httptest.NewRequest(m, "/api/katalog/cek-kuota", nil)
		w := httptest.NewRecorder()

		handler.HandlerCekKuotaRealtime(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Metode %s didapat code %d, harapan 405", m, w.Code)
		}
	}
}

// ============================================================================
// PENGUJIAN KOMPATIBILITAS POST /api/katalog/validasi-beli
// ============================================================================

func TestEndpointValidasiBeli_KompatibilitasFrontend(t *testing.T) {
	handler := &HandlerValidasiKuota{DB: nil}

	// 1. Kasus Sukses (Jumlah pembelian valid dalam stok)
	reqSukses := httptest.NewRequest(
		http.MethodPost,
		"/api/katalog/validasi-beli",
		bytes.NewBufferString(`{"id_varian":"udang","jumlah":2}`),
	)
	reqSukses.Header.Set("Content-Type", "application/json")
	wSukses := httptest.NewRecorder()

	handler.HandlerCekKuotaRealtime(wSukses, reqSukses)
	if wSukses.Code != http.StatusOK {
		t.Errorf("POST validasi-beli sukses didapat kode %d, harapan 200", wSukses.Code)
	}

	// 2. Kasus Ditolak (Jumlah pembelian melebihi stok yang ada)
	reqDitolak := httptest.NewRequest(
		http.MethodPost,
		"/api/katalog/validasi-beli",
		bytes.NewBufferString(`{"id_varian":"udang","jumlah":99}`),
	)
	reqDitolak.Header.Set("Content-Type", "application/json")
	wDitolak := httptest.NewRecorder()

	handler.HandlerCekKuotaRealtime(wDitolak, reqDitolak)
	if wDitolak.Code != http.StatusBadRequest {
		t.Errorf("POST validasi-beli melebihi stok didapat kode %d, harapan 400", wDitolak.Code)
	}
}

