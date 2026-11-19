package checkoutapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ============================================================================
// UNIT TEST WBS 1.2.3.3: STANDARISASI SCHEMA ORDER REQUEST
// ============================================================================

func TestPenyamaanPayload_Guest(t *testing.T) {
	handler := &HandlerCheckout{DB: nil}

	payloadGuest := `{
		"tipe_checkout": "guest",
		"is_guest": true,
		"id_pelanggan": null,
		"kontak": {
			"nama_penerima": "Ahmad Dani",
			"nomor_whatsapp": "081234567890",
			"email": "dani@example.com",
			"catatan": "Kirim sore ya"
		},
		"item_pesanan": [
			{
				"id_varian": 1,
				"slug": "udang",
				"nama_varian": "Udang",
				"nama_produk": "Rengginang Udang",
				"harga_satuan": 28000,
				"jumlah": 2,
				"subtotal": 56000
			}
		],
		"ringkasan": {
			"total_item": 2,
			"subtotal": 56000,
			"ongkos_kirim": 0,
			"total_pembayaran": 56000
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString(payloadGuest))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerProsesCheckout(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Guest checkout gagal: didapat status %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status respons: %s, harapan 'sukses'", resp.Status)
	}

	// Parsing data respons
	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilCheckoutResponse
	_ = json.Unmarshal(dataBytes, &hasil)

	if !hasil.IsGuest {
		t.Errorf("IsGuest harus bernilai true untuk checkout tamu")
	}
	if hasil.TipeCheckout != "guest" {
		t.Errorf("TipeCheckout: %s, harapan 'guest'", hasil.TipeCheckout)
	}
	if !strings.HasPrefix(hasil.NomorPesanan, "PJ-GUEST-") {
		t.Errorf("Nomor pesanan %s tidak berawalan PJ-GUEST-", hasil.NomorPesanan)
	}
	if hasil.Ringkasan.TotalPembayaran != 56000 {
		t.Errorf("Total pembayaran = %d, harapan 56000", hasil.Ringkasan.TotalPembayaran)
	}
}

func TestPenyamaanPayload_PelangganLogin(t *testing.T) {
	handler := &HandlerCheckout{DB: nil}

	payloadPelanggan := `{
		"tipe_checkout": "pelanggan",
		"is_guest": false,
		"id_pelanggan": 42,
		"kontak": {
			"nama_penerima": "Siti Nurhaliza",
			"nomor_whatsapp": "085678912345",
			"email": "siti@example.com"
		},
		"item_pesanan": [
			{
				"id_varian": 2,
				"slug": "bawang",
				"nama_varian": "Bawang",
				"nama_produk": "Rengginang Bawang",
				"harga_satuan": 28000,
				"jumlah": 3,
				"subtotal": 84000
			}
		],
		"ringkasan": {
			"total_item": 3,
			"subtotal": 84000,
			"ongkos_kirim": 0,
			"total_pembayaran": 84000
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString(payloadPelanggan))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerProsesCheckout(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Pelanggan checkout gagal: didapat status %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(res.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilCheckoutResponse
	_ = json.Unmarshal(dataBytes, &hasil)

	if hasil.IsGuest {
		t.Errorf("IsGuest harus bernilai false untuk pelanggan login")
	}
	if hasil.TipeCheckout != "pelanggan" {
		t.Errorf("TipeCheckout: %s, harapan 'pelanggan'", hasil.TipeCheckout)
	}
	if hasil.IDPelanggan == nil || *hasil.IDPelanggan != 42 {
		t.Errorf("IDPelanggan tidak tersimpan dengan benar")
	}
	if !strings.HasPrefix(hasil.NomorPesanan, "PJ-CUST-") {
		t.Errorf("Nomor pesanan %s tidak berawalan PJ-CUST-", hasil.NomorPesanan)
	}
}

func TestPenyamaanPayload_ValidasiInput(t *testing.T) {
	handler := &HandlerCheckout{DB: nil}

	tests := []struct {
		nama       string
		payload    string
		statusKode int
	}{
		{
			nama: "Nama penerima kosong",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "", "nomor_whatsapp": "081234567890"},
				"item_pesanan": [{"id_varian": 1, "jumlah": 1}]
			}`,
			statusKode: http.StatusBadRequest,
		},
		{
			nama: "Nomor WhatsApp tidak valid (terlalu pendek)",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "123"},
				"item_pesanan": [{"id_varian": 1, "jumlah": 1}]
			}`,
			statusKode: http.StatusBadRequest,
		},
		{
			nama: "Keranjang kosong / tanpa item",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"item_pesanan": []
			}`,
			statusKode: http.StatusBadRequest,
		},
		{
			nama: "Jumlah pesanan nol (0)",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"item_pesanan": [{"id_varian": 1, "jumlah": 0}]
			}`,
			statusKode: http.StatusBadRequest,
		},
		{
			nama: "Melebihi sisa stok (Udang stok 12, minta 50)",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"item_pesanan": [{"id_varian": 1, "nama_varian": "Udang", "jumlah": 50}]
			}`,
			statusKode: http.StatusBadRequest,
		},
		{
			nama: "Varian stok habis (Cumi stok 0)",
			payload: `{
				"tipe_checkout": "guest",
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"item_pesanan": [{"id_varian": 3, "nama_varian": "Cumi", "jumlah": 1}]
			}`,
			statusKode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandlerProsesCheckout(w, req)
			if w.Code != tt.statusKode {
				t.Errorf("%s: didapat kode %d, harapan %d", tt.nama, w.Code, tt.statusKode)
			}
		})
	}
}

func TestHandlerSkemaPayload_GET(t *testing.T) {
	handler := &HandlerCheckout{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/checkout/skema", nil)
	w := httptest.NewRecorder()

	handler.HandlerSkemaPayload(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Endpoint skema didapat kode %d, harapan 200", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON respons: %v", err)
	}

	if resp["deliverable"] != "Standarisasi Schema Order Request" {
		t.Errorf("Deliverable tidak sesuai: %v", resp["deliverable"])
	}
}

func TestMethodNotAllowed(t *testing.T) {
	handler := &HandlerCheckout{DB: nil}

	metodeDilarang := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, m := range metodeDilarang {
		req := httptest.NewRequest(m, "/api/checkout", nil)
		w := httptest.NewRecorder()

		handler.HandlerProsesCheckout(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Metode %s: code %d, harapan 405", m, w.Code)
		}
	}
}
