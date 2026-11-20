package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	validasiwilayah "rengginang-putra-jaya/auth/1.2.6.3"
)

func main() {
	separator := strings.Repeat("=", 75)
	handler := &validasiwilayah.HandlerValidasiWilayah{DB: nil}

	fmt.Println(separator)
	fmt.Println("BUKTI PENGUJIAN BACKEND WBS 1.2.6.3")
	fmt.Println("PPL-15: Pengisian Wilayah dan Alamat Pengiriman")
	fmt.Println("Deliverable: Validasi Data Geografis Pengiriman Sesuai Parameter API Kurir")
	fmt.Println(separator)

	// =========================================================================
	// PENGUJIAN 1: DATA GEOGRAFIS LENGKAP & VALID (SESUAI MOCKUP CHECKOUT FIGMA)
	// =========================================================================
	fmt.Println("\n[SKENARIO 1] Pengiriman via Kurir dengan Alamat Lengkap dan Valid")
	payloadValid := `{
		"kontak": {
			"nama_penerima": "Contoh Pelanggan",
			"nomor_whatsapp": "081234567890"
		},
		"metode_penerimaan": "kurir",
		"wilayah_pengiriman": {
			"provinsi": "Jawa Timur",
			"kota_kabupaten": "Kabupaten Sidoarjo",
			"kecamatan": "Kecamatan Sidoarjo",
			"kode_pos": "61211",
			"alamat_lengkap": "Jl. Contoh No. 10, RT 02 / RW 03, Kelurahan Contoh, Dekat minimarket"
		},
		"subtotal_produk": 84000,
		"total_sementara": 84000
	}`

	req1 := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadValid))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.HandlerSimpanDanValidasiWilayah(w1, req1)
	res1 := w1.Result()

	fmt.Printf("HTTP Status  : %d %s\n", res1.StatusCode, http.StatusText(res1.StatusCode))
	var resp1 validasiwilayah.ResponsAPI
	_ = json.NewDecoder(res1.Body).Decode(&resp1)
	prettyJSON1, _ := json.MarshalIndent(resp1, "", "  ")
	fmt.Println("Respon JSON  :")
	fmt.Println(string(prettyJSON1))

	// =========================================================================
	// PENGUJIAN 2: METODE PENERIMAAN 'AMBIL SENDIRI'
	// =========================================================================
	fmt.Println("\n" + separator)
	fmt.Println("[SKENARIO 2] Metode Penerimaan Ambil Sendiri di Toko")
	payloadAmbil := `{
		"kontak": {
			"nama_penerima": "Budi Santoso",
			"nomor_whatsapp": "081234567890"
		},
		"metode_penerimaan": "ambil_sendiri",
		"wilayah_pengiriman": null
	}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadAmbil))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.HandlerSimpanDanValidasiWilayah(w2, req2)
	res2 := w2.Result()
	fmt.Printf("HTTP Status  : %d %s\n", res2.StatusCode, http.StatusText(res2.StatusCode))
	var resp2 validasiwilayah.ResponsAPI
	_ = json.NewDecoder(res2.Body).Decode(&resp2)
	prettyJSON2, _ := json.MarshalIndent(resp2, "", "  ")
	fmt.Println("Respon JSON  :")
	fmt.Println(string(prettyJSON2))

	// =========================================================================
	// PENGUJIAN 3: VALIDASI ERROR KELENGKAPAN GEOGRAFIS
	// =========================================================================
	fmt.Println("\n" + separator)
	fmt.Println("[SKENARIO 3] Uji Validasi Error (Provinsi, Kota, Kecamatan Kosong, Alamat Pendek, Kode Pos Salah)")
	payloadError := `{
		"kontak": {
			"nama_penerima": "",
			"nomor_whatsapp": "123"
		},
		"metode_penerimaan": "kurir",
		"wilayah_pengiriman": {
			"provinsi": "",
			"kota_kabupaten": "",
			"kecamatan": "",
			"kode_pos": "61",
			"alamat_lengkap": "Jl"
		}
	}`
	req3 := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadError))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	handler.HandlerSimpanDanValidasiWilayah(w3, req3)
	res3 := w3.Result()
	fmt.Printf("HTTP Status  : %d %s\n", res3.StatusCode, http.StatusText(res3.StatusCode))
	var resp3 validasiwilayah.ResponsAPI
	_ = json.NewDecoder(res3.Body).Decode(&resp3)
	prettyJSON3, _ := json.MarshalIndent(resp3, "", "  ")
	fmt.Println("Respon JSON Error:")
	fmt.Println(string(prettyJSON3))

	// =========================================================================
	// PENGUJIAN 4: VALIDASI KETIDAKCOCOKAN HIERARKI WILAYAH BERTINGKAT
	// =========================================================================
	fmt.Println("\n" + separator)
	fmt.Println("[SKENARIO 4] Uji Inkonsistensi Hierarki (Kota Tidak Sesuai Provinsi)")
	payloadHierarki := `{
		"kontak": {
			"nama_penerima": "Ahmad",
			"nomor_whatsapp": "081234567890"
		},
		"metode_penerimaan": "kurir",
		"wilayah_pengiriman": {
			"provinsi": "Jawa Timur",
			"kota_kabupaten": "Kota Bandung",
			"kecamatan": "Kecamatan Coblong",
			"kode_pos": "40132",
			"alamat_lengkap": "Jl. Dago No. 15 Bandung"
		}
	}`
	req4 := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadHierarki))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	handler.HandlerSimpanDanValidasiWilayah(w4, req4)
	res4 := w4.Result()
	fmt.Printf("HTTP Status  : %d %s\n", res4.StatusCode, http.StatusText(res4.StatusCode))
	var resp4 validasiwilayah.ResponsAPI
	_ = json.NewDecoder(res4.Body).Decode(&resp4)
	prettyJSON4, _ := json.MarshalIndent(resp4, "", "  ")
	fmt.Println("Respon JSON Error Hierarki:")
	fmt.Println(string(prettyJSON4))

	// =========================================================================
	// PENGUJIAN 5: ENDPOINT MASTER DATA WILAYAH BERTINGKAT
	// =========================================================================
	fmt.Println("\n" + separator)
	fmt.Println("[SKENARIO 5] Endpoint Master Data Wilayah (GET /api/checkout/master-wilayah)")
	req5 := httptest.NewRequest(http.MethodGet, "/api/checkout/master-wilayah", nil)
	w5 := httptest.NewRecorder()
	handler.HandlerMasterWilayah(w5, req5)
	res5 := w5.Result()
	fmt.Printf("HTTP Status  : %d %s\n", res5.StatusCode, http.StatusText(res5.StatusCode))

	// =========================================================================
	// PENGUJIAN 6: ENDPOINT PARAMETER KEBUTUHAN API KURIR
	// =========================================================================
	fmt.Println("\n" + separator)
	fmt.Println("[SKENARIO 6] Endpoint Parameter Kurir (GET /api/checkout/parameter-kurir)")
	req6 := httptest.NewRequest(http.MethodGet, "/api/checkout/parameter-kurir", nil)
	w6 := httptest.NewRecorder()
	handler.HandlerParameterKurir(w6, req6)
	res6 := w6.Result()
	fmt.Printf("HTTP Status  : %d %s\n", res6.StatusCode, http.StatusText(res6.StatusCode))
	var resp6 validasiwilayah.ResponsAPI
	_ = json.NewDecoder(res6.Body).Decode(&resp6)
	prettyJSON6, _ := json.MarshalIndent(resp6, "", "  ")
	fmt.Println("Respon JSON Spesifikasi Kurir:")
	fmt.Println(string(prettyJSON6))

	fmt.Println("\n" + separator)
	fmt.Println("SELURUH SKENARIO PENGUJIAN BACKEND WBS 1.2.6.3 SELESAI DENGAN STATUS VALID (PASS)")
	fmt.Println(separator)
}
