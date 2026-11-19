package statemanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// UNIT TEST WBS 1.2.2.2: STATE MANAGER KERANJANG BELANJA
// ============================================================================

func TestSinkronisasiKeranjang_Guest(t *testing.T) {
	handler := &HandlerStateManager{DB: nil}

	payloadGuest := `{
		"is_guest": true,
		"items": [
			{
				"id_varian": 1,
				"slug": "udang",
				"nama_varian": "Udang",
				"nama_produk": "Rengginang Udang",
				"harga_satuan": 28000,
				"jumlah": 2,
				"stok_tersedia": 12
			},
			{
				"id_varian": 2,
				"slug": "bawang",
				"nama_varian": "Bawang",
				"nama_produk": "Rengginang Bawang",
				"harga_satuan": 28000,
				"jumlah": 1,
				"stok_tersedia": 38
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/sinkronisasi", bytes.NewBufferString(payloadGuest))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerSinkronisasiKeranjang(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Sinkronisasi guest gagal: didapat status %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status respons: %s, harapan 'sukses'", resp.Status)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var keranjang ResponKeranjang
	_ = json.Unmarshal(dataBytes, &keranjang)

	if !keranjang.IsGuest {
		t.Errorf("IsGuest harus bernilai true untuk sinkronisasi tamu")
	}
	if keranjang.TotalKuantitas != 3 {
		t.Errorf("Total kuantitas: %d, harapan 3", keranjang.TotalKuantitas)
	}
	if keranjang.Subtotal != 84000 {
		t.Errorf("Subtotal: %d, harapan 84000 (2 Udang + 1 Bawang)", keranjang.Subtotal)
	}
}

func TestSinkronisasiKeranjang_PelangganLogin(t *testing.T) {
	handler := &HandlerStateManager{DB: nil}

	payloadPelanggan := `{
		"is_guest": false,
		"id_pelanggan": 42,
		"items": [
			{
				"id_varian": 1,
				"nama_varian": "Udang",
				"nama_produk": "Rengginang Udang",
				"harga_satuan": 28000,
				"jumlah": 3,
				"stok_tersedia": 12
			},
			{
				"id_varian": 2,
				"nama_varian": "Bawang",
				"nama_produk": "Rengginang Bawang",
				"harga_satuan": 28000,
				"jumlah": 1,
				"stok_tersedia": 38
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/keranjang/sinkronisasi", bytes.NewBufferString(payloadPelanggan))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerSinkronisasiKeranjang(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Sinkronisasi pelanggan login gagal: kode %d", w.Code)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(w.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var keranjang ResponKeranjang
	_ = json.Unmarshal(dataBytes, &keranjang)

	if keranjang.IsGuest {
		t.Errorf("IsGuest harus bernilai false untuk pelanggan terdaftar")
	}
	if keranjang.IDPelanggan == nil || *keranjang.IDPelanggan != 42 {
		t.Errorf("IDPelanggan tidak tersimpan dengan benar")
	}
	// Kondisi PPL11_02: 3 Udang (84.000) + 1 Bawang (28.000) = 112.000
	if keranjang.Subtotal != 112000 {
		t.Errorf("Subtotal kondisi PPL11_02: %d, harapan 112000", keranjang.Subtotal)
	}
}

func TestAmbilKeranjang_Default(t *testing.T) {
	handler := &HandlerStateManager{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/keranjang", nil)
	w := httptest.NewRecorder()

	handler.HandlerAmbilKeranjang(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Ambil keranjang gagal: kode %d", w.Code)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(w.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var keranjang ResponKeranjang
	_ = json.Unmarshal(dataBytes, &keranjang)

	if len(keranjang.Items) != 2 {
		t.Errorf("Jumlah item default: %d, harapan 2 jenis produk", len(keranjang.Items))
	}
	if keranjang.Subtotal != 84000 {
		t.Errorf("Subtotal default: %d, harapan 84000", keranjang.Subtotal)
	}
}

func TestHapusItemKeranjang(t *testing.T) {
	handler := &HandlerStateManager{DB: nil}

	// Hapus item Udang (id_varian = 1)
	req := httptest.NewRequest(http.MethodDelete, "/api/keranjang/item?id_varian=1&id_pelanggan=1", nil)
	w := httptest.NewRecorder()

	handler.HandlerHapusItemKeranjang(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Hapus item gagal: kode %d", w.Code)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(w.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var keranjang ResponKeranjang
	_ = json.Unmarshal(dataBytes, &keranjang)

	for _, it := range keranjang.Items {
		if it.IDVarian == 1 {
			t.Errorf("Item id 1 masih ditemukan setelah dihapus")
		}
	}
}

func TestMethodNotAllowed_Keranjang(t *testing.T) {
	handler := &HandlerStateManager{DB: nil}

	req := httptest.NewRequest(http.MethodPut, "/api/keranjang/sinkronisasi", nil)
	w := httptest.NewRecorder()

	handler.HandlerSinkronisasiKeranjang(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Metode PUT: kode %d, harapan 405", w.Code)
	}
}
