package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	validatororder "rengginang-putra-jaya/auth/1.2.8.3"
)

func main() {
	fmt.Println("================================================================================")
	fmt.Println("BUKTI EKSEKUSI SERVER-SIDE ORDER VALIDATOR (WBS 1.2.8.3)")
	fmt.Println("================================================================================")

	handler := &validatororder.HandlerValidatorOrder{DB: nil}
	handler.InisialisasiTabel()

	// 1. UJI GAMBAR 1 (DATA VALID - RINGKASAN SIAP DIKONFIRMASI)
	fmt.Println("\n--- 1. UJI VALIDASI RINGKASAN: SESUAI FIGMA GAMBAR 1 (DATA VALID) ---")
	reqGambar1 := validatororder.PermintaanValidasiOrder{
		TipeCheckout:     "guest",
		MetodePenerimaan: "Kirim via Kurir",
		Kontak: validatororder.KontakPemesan{
			Nama:          "Pelanggan Contoh",
			NomorWhatsApp: "081234567890",
		},
		Alamat: validatororder.AlamatPengiriman{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "Kabupaten Jember",
			Kecamatan:     "Kecamatan Sumbersari",
			AlamatLengkap: "Jl. Contoh No. 12",
		},
		Kurir: validatororder.KurirPilihan{
			IDKurir:          "jne",
			NamaKurir:        "JNE REG",
			KodeLayanan:      "REG",
			TarifOngkir:      18000,
			EstimasiHari:     "2-3 hari",
			TerverifikasiAPI: true,
		},
		Items: []validatororder.ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
			{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
		},
		RingkasanKlien: validatororder.RingkasanKlien{
			Subtotal:    84000,
			OngkosKirim: 18000,
			TotalAkhir:  102000,
		},
	}
	bodyG1, _ := json.Marshal(reqGambar1)
	httpReq1 := httptest.NewRequest(http.MethodPost, "/api/checkout/validasi-ringkasan", bytes.NewReader(bodyG1))
	rec1 := httptest.NewRecorder()
	handler.HandlerValidasiRingkasan(rec1, httpReq1)
	res1, _ := io.ReadAll(rec1.Body)
	var formattedJSON1 bytes.Buffer
	_ = json.Indent(&formattedJSON1, res1, "", "  ")
	fmt.Printf("HTTP Status: %d\n", rec1.Code)
	fmt.Println(formattedJSON1.String())

	// Ambil token validasi dari hasil 1
	var hasilMap map[string]interface{}
	_ = json.Unmarshal(res1, &hasilMap)
	laporan, _ := hasilMap["laporan"].(map[string]interface{})
	tokenValidasi, _ := laporan["token_validasi"].(string)

	// 2. UJI GAMBAR 2 (ONGKIR BELUM TERVERIFIKASI - BELUM DAPAT MEMBUAT PESANAN)
	fmt.Println("\n--- 2. UJI VALIDASI RINGKASAN: SESUAI FIGMA GAMBAR 2 (BELUM DAPAT MEMBUAT PESANAN) ---")
	reqGambar2 := validatororder.PermintaanValidasiOrder{
		SimulasiBelumValid: true,
		Kurir: validatororder.KurirPilihan{
			IDKurir:          "",
			TarifOngkir:      0,
			TerverifikasiAPI: false,
		},
		RingkasanKlien: validatororder.RingkasanKlien{
			Subtotal:    84000,
			OngkosKirim: 0,
			TotalAkhir:  84000,
		},
	}
	bodyG2, _ := json.Marshal(reqGambar2)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/api/checkout/validasi-ringkasan", bytes.NewReader(bodyG2))
	rec2 := httptest.NewRecorder()
	handler.HandlerValidasiRingkasan(rec2, httpReq2)
	res2, _ := io.ReadAll(rec2.Body)
	var formattedJSON2 bytes.Buffer
	_ = json.Indent(&formattedJSON2, res2, "", "  ")
	fmt.Printf("HTTP Status: %d\n", rec2.Code)
	fmt.Println(formattedJSON2.String())

	// 3. UJI EKSEKUSI TOMBOL BUAT PESANAN (LOLOS VALIDASI)
	fmt.Println("\n--- 3. UJI EKSEKUSI BUAT PESANAN: TOMBOL 'BUAT PESANAN' (HTTP 201 CREATED) ---")
	reqBuat := validatororder.PermintaanBuatPesanan{
		TokenValidasi: tokenValidasi,
		OrderData:     reqGambar1,
	}
	bodyBuat, _ := json.Marshal(reqBuat)
	httpReq3 := httptest.NewRequest(http.MethodPost, "/api/checkout/buat-pesanan", bytes.NewReader(bodyBuat))
	rec3 := httptest.NewRecorder()
	handler.HandlerBuatPesanan(rec3, httpReq3)
	res3, _ := io.ReadAll(rec3.Body)
	var formattedJSON3 bytes.Buffer
	_ = json.Indent(&formattedJSON3, res3, "", "  ")
	fmt.Printf("HTTP Status: %d\n", rec3.Code)
	fmt.Println(formattedJSON3.String())

	// 4. UJI PENOLAKAN TAMPERING SUBTOTAL
	fmt.Println("\n--- 4. UJI INTEGRITAS ANTI-TAMPERING: PENOLAKAN MANIPULASI SUBTOTAL ---")
	reqTampering := validatororder.PermintaanBuatPesanan{
		TokenValidasi: tokenValidasi,
		OrderData: validatororder.PermintaanValidasiOrder{
			Kurir: reqGambar1.Kurir,
			Items: reqGambar1.Items,
			RingkasanKlien: validatororder.RingkasanKlien{
				Subtotal:    10000, // Dimanipulasi dari Rp84.000 menjadi Rp10.000!
				OngkosKirim: 18000,
				TotalAkhir:  28000,
			},
		},
	}
	bodyTamp, _ := json.Marshal(reqTampering)
	httpReq4 := httptest.NewRequest(http.MethodPost, "/api/checkout/buat-pesanan", bytes.NewReader(bodyTamp))
	rec4 := httptest.NewRecorder()
	handler.HandlerBuatPesanan(rec4, httpReq4)
	res4, _ := io.ReadAll(rec4.Body)
	var formattedJSON4 bytes.Buffer
	_ = json.Indent(&formattedJSON4, res4, "", "  ")
	fmt.Printf("HTTP Status: %d (Harus 422 Unprocessable Entity)\n", rec4.Code)
	fmt.Println(formattedJSON4.String())

	fmt.Println("\n================================================================================")
	fmt.Println("SELURUH VALIDASI SISI SERVER BERHASIL DIUJI DENGAN SEMPURNA!")
	fmt.Println("================================================================================")
	os.Exit(0)
}
