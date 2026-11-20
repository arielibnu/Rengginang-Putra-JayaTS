package validasikeranjang

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// UNIT TEST WBS 1.2.2.4: API VALIDASI KERANJANG TERHADAP STOK
// ============================================================================

func TestValidasiKeranjang_Normal_Pass(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	// Payload sesuai kondisi Figma PPL11_01: 2 Udang + 1 Bawang (Subtotal Rp84.000)
	payload := `{
		"items": [
			{"id_varian": 1, "nama_varian": "Udang", "jumlah": 2},
			{"id_varian": 2, "nama_varian": "Bawang", "jumlah": 1}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/validasi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Validasi normal gagal: kode %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status respons: %s, harapan 'sukses'", resp.Status)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilValidasiKeranjang
	_ = json.Unmarshal(dataBytes, &hasil)

	if !hasil.BisaCheckout {
		t.Errorf("BisaCheckout harus bernilai true")
	}
	if hasil.Subtotal != 84000 {
		t.Errorf("Subtotal: %d, harapan 84000", hasil.Subtotal)
	}
}

func TestValidasiKeranjang_KondisiPPL11_02(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	// Payload sesuai kondisi Figma PPL11_02: 3 Udang + 1 Bawang (Subtotal Rp112.000)
	payload := `{
		"items": [
			{"id_varian": 1, "nama_varian": "Udang", "jumlah": 3},
			{"id_varian": 2, "nama_varian": "Bawang", "jumlah": 1}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/validasi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Validasi PPL11_02 gagal: kode %d", w.Code)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(w.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilValidasiKeranjang
	_ = json.Unmarshal(dataBytes, &hasil)

	if hasil.Subtotal != 112000 {
		t.Errorf("Subtotal: %d, harapan 112000", hasil.Subtotal)
	}
}

func TestValidasiKeranjang_MelebihiStok_Fail(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	// Udang stok cuma 12 bungkus, diminta 20 bungkus
	payload := `{
		"items": [
			{"id_varian": 1, "nama_varian": "Udang", "jumlah": 20}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/validasi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Pesanan melebihi stok: kode %d, harapan 400", w.Code)
	}
}

func TestValidasiKeranjang_StokHabis_Fail(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	// Cumi stok 0
	payload := `{
		"items": [
			{"id_varian": 3, "nama_varian": "Cumi", "jumlah": 1}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/validasi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Pesanan stok habis: kode %d, harapan 400", w.Code)
	}
}

func TestValidasiKeranjang_Kosong_Fail(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	payload := `{"items": []}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/validasi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Keranjang kosong: kode %d, harapan 400", w.Code)
	}
}

func TestMethodNotAllowed_Validasi(t *testing.T) {
	handler := &HandlerValidasiKeranjang{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/keranjang/validasi", nil)
	w := httptest.NewRecorder()

	handler.HandlerValidasiKeranjang(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Metode GET: kode %d, harapan 405", w.Code)
	}
}
