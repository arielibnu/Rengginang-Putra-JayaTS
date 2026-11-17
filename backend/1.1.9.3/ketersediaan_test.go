package ketersediaanproduk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestTentukanStatusKetersediaan menguji aturan penentuan status stok habis dan ketersediaan pembelian.
func TestTentukanStatusKetersediaan(t *testing.T) {
	tests := []struct {
		nama           string
		stok           int64
		harapanStatus  string
		harapanBisaBeli bool
	}{
		{
			nama:           "Stok 0 (Kasus Cumi pada Gambar 1 Figma)",
			stok:           0,
			harapanStatus:  "STOK HABIS",
			harapanBisaBeli: false,
		},
		{
			nama:           "Stok negatif",
			stok:           -5,
			harapanStatus:  "STOK HABIS",
			harapanBisaBeli: false,
		},
		{
			nama:           "Stok 12 bungkus (Kasus Udang pada Gambar 2 Figma)",
			stok:           12,
			harapanStatus:  "TERSEDIA",
			harapanBisaBeli: true,
		},
		{
			nama:           "Stok 38 bungkus (Kasus Bawang)",
			stok:           38,
			harapanStatus:  "TERSEDIA",
			harapanBisaBeli: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			status, bisaBeli, pesan, _ := TentukanStatusKetersediaan(tt.stok)
			if status != tt.harapanStatus {
				t.Errorf("Status salah: didapat %q, harapan %q", status, tt.harapanStatus)
			}
			if bisaBeli != tt.harapanBisaBeli {
				t.Errorf("BisaDibeli salah: didapat %v, harapan %v", bisaBeli, tt.harapanBisaBeli)
			}
			if status == "STOK HABIS" && pesan != "Produk sedang habis" {
				t.Errorf("Pesan status salah: %q", pesan)
			}
		})
	}
}

// TestValidasiInteraksiBeli menguji pembatasan tombol disabled out-of-stock.
func TestValidasiInteraksiBeli(t *testing.T) {
	tests := []struct {
		nama           string
		stokTersedia   int64
		jumlahBeli     uint32
		harapanBisaBeli bool
		harusError     bool
	}{
		{
			nama:           "Beli produk stok habis (dinonaktifkan)",
			stokTersedia:   0,
			jumlahBeli:     1,
			harapanBisaBeli: false,
			harusError:     true,
		},
		{
			nama:           "Beli produk tersedia dengan jumlah valid",
			stokTersedia:   12,
			jumlahBeli:     2,
			harapanBisaBeli: true,
			harusError:     false,
		},
		{
			nama:           "Beli melebihi stok yang tersedia",
			stokTersedia:   12,
			jumlahBeli:     15,
			harapanBisaBeli: false,
			harusError:     true,
		},
		{
			nama:           "Jumlah beli nol",
			stokTersedia:   12,
			jumlahBeli:     0,
			harapanBisaBeli: false,
			harusError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			bisaBeli, _, err := ValidasiInteraksiBeli(tt.stokTersedia, tt.jumlahBeli)
			if bisaBeli != tt.harapanBisaBeli {
				t.Errorf("BisaBeli didapat %v, harapan %v", bisaBeli, tt.harapanBisaBeli)
			}
			if tt.harusError && err == nil {
				t.Errorf("Diharapkan error tetapi nil")
			}
			if !tt.harusError && err != nil {
				t.Errorf("Diharapkan sukses tapi error: %v", err)
			}
		})
	}
}

// TestFormatRupiah menguji format nilai rupiah.
func TestFormatRupiah(t *testing.T) {
	if FormatRupiah(28000) != "Rp28.000" {
		t.Errorf("FormatRupiah salah: didapat %q", FormatRupiah(28000))
	}
	if FormatRupiah(1500000) != "Rp1.500.000" {
		t.Errorf("FormatRupiah salah: didapat %q", FormatRupiah(1500000))
	}
}

// TestHandlerDetailProduk menguji endpoint HTTP detail produk.
func TestHandlerDetailProduk(t *testing.T) {
	h := &HandlerDetail{DB: nil}

	// 1. Uji produk Cumi (Stok Habis)
	reqCumi := httptest.NewRequest(http.MethodGet, "/api/katalog/detail?id=cumi", nil)
	recCumi := httptest.NewRecorder()
	h.HandlerDetailProduk(recCumi, reqCumi)

	if recCumi.Code != http.StatusOK {
		t.Fatalf("Kode respon salah: %d", recCumi.Code)
	}

	var respCumi ResponsAPI
	if err := json.Unmarshal(recCumi.Body.Bytes(), &respCumi); err != nil {
		t.Fatalf("Gagal unmarshal JSON: %v", err)
	}

	dataCumiMap, ok := respCumi.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data response bukan map")
	}

	if dataCumiMap["status_ketersediaan"] != "STOK HABIS" {
		t.Errorf("Status ketersediaan Cumi harus 'STOK HABIS', didapat: %v", dataCumiMap["status_ketersediaan"])
	}
	if dataCumiMap["bisa_dibeli"] != false {
		t.Errorf("Bisa dibeli Cumi harus false")
	}

	// 2. Uji produk Udang (Tersedia)
	reqUdang := httptest.NewRequest(http.MethodGet, "/api/katalog/detail?id=udang", nil)
	recUdang := httptest.NewRecorder()
	h.HandlerDetailProduk(recUdang, reqUdang)

	var respUdang ResponsAPI
	_ = json.Unmarshal(recUdang.Body.Bytes(), &respUdang)
	dataUdangMap := respUdang.Data.(map[string]interface{})

	if dataUdangMap["status_ketersediaan"] != "TERSEDIA" {
		t.Errorf("Status ketersediaan Udang harus 'TERSEDIA'")
	}
	if dataUdangMap["bisa_dibeli"] != true {
		t.Errorf("Bisa dibeli Udang harus true")
	}
}

// TestHandlerValidasiBeli menguji endpoint validasi tombol beli out-of-stock.
func TestHandlerValidasiBeli(t *testing.T) {
	h := &HandlerDetail{DB: nil}

	// 1. Coba validasi beli produk yang habis (Cumi) -> Ditolak (HTTP 400)
	bodyCumi, _ := json.Marshal(PermintaanValidasiBeli{Slug: "cumi", JumlahBeli: 1})
	reqCumi := httptest.NewRequest(http.MethodPost, "/api/katalog/validasi-beli", bytes.NewReader(bodyCumi))
	recCumi := httptest.NewRecorder()
	h.HandlerValidasiBeli(recCumi, reqCumi)

	if recCumi.Code != http.StatusBadRequest {
		t.Errorf("Harus ditolak 400 Bad Request untuk produk habis, didapat: %d", recCumi.Code)
	}

	// 2. Validasi beli produk yang tersedia (Udang) -> Diterima (HTTP 200)
	bodyUdang, _ := json.Marshal(PermintaanValidasiBeli{Slug: "udang", JumlahBeli: 1})
	reqUdang := httptest.NewRequest(http.MethodPost, "/api/katalog/validasi-beli", bytes.NewReader(bodyUdang))
	recUdang := httptest.NewRecorder()
	h.HandlerValidasiBeli(recUdang, reqUdang)

	if recUdang.Code != http.StatusOK {
		t.Errorf("Harus 200 OK untuk produk tersedia, didapat: %d", recUdang.Code)
	}
}
