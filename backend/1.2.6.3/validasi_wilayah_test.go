package validasiwilayah

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// UNIT TEST WBS 1.2.6.3: VALIDASI DATA GEOGRAFIS PENGIRIMAN API KURIR
// ============================================================================

func TestValidasiWilayah_LengkapDanValid(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

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

	req := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadValid))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerSimpanDanValidasiWilayah(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Validasi wilayah gagal: didapat status %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status respons: %s, harapan 'sukses'", resp.Status)
	}

	// Parsing Hasil
	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilValidasiWilayah
	_ = json.Unmarshal(dataBytes, &hasil)

	if !hasil.Valid {
		t.Errorf("Harapan valid = true, didapat false")
	}

	if hasil.ParameterKurir == nil {
		t.Fatalf("ParameterKurir tidak boleh nil")
	}

	if !hasil.ParameterKurir.IsKurirReady {
		t.Errorf("Harapan IsKurirReady = true")
	}

	if hasil.ParameterKurir.StatusGeografis != "VALID" {
		t.Errorf("StatusGeografis: %s, harapan 'VALID'", hasil.ParameterKurir.StatusGeografis)
	}

	// Verifikasi Asal Toko & Tujuan
	if hasil.ParameterKurir.AsalPengiriman.KotaKabupaten != "Kabupaten Situbondo" {
		t.Errorf("Asal pengiriman salah: %s", hasil.ParameterKurir.AsalPengiriman.KotaKabupaten)
	}
	if hasil.ParameterKurir.TujuanPengiriman.KotaKabupaten != "Kabupaten Sidoarjo" {
		t.Errorf("Tujuan pengiriman salah: %s", hasil.ParameterKurir.TujuanPengiriman.KotaKabupaten)
	}
	if hasil.ParameterKurir.TujuanPengiriman.KodePos != "61211" {
		t.Errorf("Kode pos tujuan salah: %s", hasil.ParameterKurir.TujuanPengiriman.KodePos)
	}
	if hasil.ParameterKurir.Paket.BeratGram == 0 {
		t.Errorf("Berat paket tidak boleh 0")
	}
	if !hasil.ParameterKurir.Paket.IsFragile {
		t.Errorf("Paket rengginang harus ditandai is_fragile = true")
	}
}

func TestValidasiWilayah_AmbilSendiri(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	payloadAmbilSendiri := `{
		"kontak": {
			"nama_penerima": "Budi Santoso",
			"nomor_whatsapp": "081234567890"
		},
		"metode_penerimaan": "ambil_sendiri",
		"wilayah_pengiriman": null,
		"subtotal_produk": 84000,
		"total_sementara": 84000
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(payloadAmbilSendiri))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerSimpanDanValidasiWilayah(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Metode ambil_sendiri gagal: status %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(res.Body).Decode(&resp)

	dataBytes, _ := json.Marshal(resp.Data)
	var hasil HasilValidasiWilayah
	_ = json.Unmarshal(dataBytes, &hasil)

	if hasil.ParameterKurir.IsKurirReady {
		t.Errorf("IsKurirReady harus false untuk metode ambil_sendiri")
	}
	if hasil.ParameterKurir.StatusGeografis != "NOT_REQUIRED" {
		t.Errorf("StatusGeografis: %s, harapan 'NOT_REQUIRED'", hasil.ParameterKurir.StatusGeografis)
	}
}

func TestValidasiWilayah_ParameterTidakLengkap(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	kasusUji := []struct {
		namaKasus  string
		payload    string
		fieldError string
	}{
		{
			namaKasus: "Wilayah pengiriman kosong saat kurir",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": null
			}`,
			fieldError: "wilayah_pengiriman",
		},
		{
			namaKasus: "Provinsi tujuan kosong",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "provinsi",
		},
		{
			namaKasus: "Kota/Kabupaten tujuan kosong",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "kota_kabupaten",
		},
		{
			namaKasus: "Kecamatan tujuan kosong",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "kecamatan",
		},
		{
			namaKasus: "Alamat lengkap kosong",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": ""
				}
			}`,
			fieldError: "alamat_lengkap",
		},
		{
			namaKasus: "Alamat lengkap terlalu singkat (< 10 karakter)",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. A 1"
				}
			}`,
			fieldError: "alamat_lengkap",
		},
		{
			namaKasus: "Kode pos bukan 5 digit",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "612",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "kode_pos",
		},
		{
			namaKasus: "Kota tidak cocok dengan provinsi dalam master data",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kota Bandung",
					"kecamatan": "Kecamatan Coblong",
					"kode_pos": "40132",
					"alamat_lengkap": "Jl. Dago No. 10 Bandung"
				}
			}`,
			fieldError: "kota_kabupaten",
		},
		{
			namaKasus: "Kecamatan tidak cocok dengan kota dalam master data",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Coblong",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "kecamatan",
		},
		{
			namaKasus: "Nama penerima kosong",
			payload: `{
				"kontak": {"nama_penerima": "", "nomor_whatsapp": "081234567890"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "nama_penerima",
		},
		{
			namaKasus: "Nomor WhatsApp bukan nomor valid Indonesia",
			payload: `{
				"kontak": {"nama_penerima": "Budi", "nomor_whatsapp": "12345"},
				"metode_penerimaan": "kurir",
				"wilayah_pengiriman": {
					"provinsi": "Jawa Timur",
					"kota_kabupaten": "Kabupaten Sidoarjo",
					"kecamatan": "Kecamatan Sidoarjo",
					"kode_pos": "61211",
					"alamat_lengkap": "Jl. Contoh No. 10 Sidoarjo"
				}
			}`,
			fieldError: "nomor_whatsapp",
		},
	}

	for _, tc := range kasusUji {
		t.Run(tc.namaKasus, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/checkout/wilayah-pengiriman", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandlerSimpanDanValidasiWilayah(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Fatalf("Status didapat %d, harapan 400 Bad Request", res.StatusCode)
			}

			var resp ResponsAPI
			_ = json.NewDecoder(res.Body).Decode(&resp)

			if resp.Status != "error" {
				t.Errorf("Harapan status 'error', didapat '%s'", resp.Status)
			}

			if _, ok := resp.Errors[tc.fieldError]; !ok {
				t.Errorf("Error untuk field '%s' tidak ditemukan dalam respons: %v", tc.fieldError, resp.Errors)
			}
		})
	}
}

func TestHandlerValidasiAlamatKurir_POST(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	payload := `{
		"kontak": {"nama_penerima": "Siti Aminah", "nomor_whatsapp": "085678901234"},
		"metode_penerimaan": "kurir",
		"wilayah_pengiriman": {
			"provinsi": "DKI Jakarta",
			"kota_kabupaten": "Kota Jakarta Selatan",
			"kecamatan": "Kecamatan Cilandak",
			"kode_pos": "12430",
			"alamat_lengkap": "Jl. TB Simatupang No. 25, Cilandak Barat"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/checkout/validasi-alamat", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandlerValidasiAlamatKurir(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status didapat %d, harapan 200", res.StatusCode)
	}
}

func TestHandlerMasterWilayah_GET(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/checkout/master-wilayah", nil)
	w := httptest.NewRecorder()

	handler.HandlerMasterWilayah(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status didapat %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(res.Body).Decode(&resp)

	if resp.Status != "sukses" {
		t.Errorf("Status didapat %s, harapan 'sukses'", resp.Status)
	}
}

func TestHandlerParameterKurir_GET(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/checkout/parameter-kurir", nil)
	w := httptest.NewRecorder()

	handler.HandlerParameterKurir(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status didapat %d, harapan 200", res.StatusCode)
	}

	var resp ResponsAPI
	_ = json.NewDecoder(res.Body).Decode(&resp)

	if resp.Status != "sukses" {
		t.Errorf("Status didapat %s, harapan 'sukses'", resp.Status)
	}
}

func TestMethodNotAllowedWilayah(t *testing.T) {
	handler := &HandlerValidasiWilayah{DB: nil}

	// GET pada endpoint POST
	req := httptest.NewRequest(http.MethodGet, "/api/checkout/wilayah-pengiriman", nil)
	w := httptest.NewRecorder()
	handler.HandlerSimpanDanValidasiWilayah(w, req)
	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Harapan 405 Method Not Allowed untuk GET di wilayah-pengiriman")
	}

	// POST pada endpoint GET
	req2 := httptest.NewRequest(http.MethodPost, "/api/checkout/master-wilayah", nil)
	w2 := httptest.NewRecorder()
	handler.HandlerMasterWilayah(w2, req2)
	if w2.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Harapan 405 Method Not Allowed untuk POST di master-wilayah")
	}
}

func TestFormatNomorE164(t *testing.T) {
	if got := FormatNomorE164("081234567890"); got != "+6281234567890" {
		t.Errorf("FormatNomorE164(0812...) = %s, harapan +6281234567890", got)
	}
	if got := FormatNomorE164("6281234567890"); got != "+6281234567890" {
		t.Errorf("FormatNomorE164(62812...) = %s, harapan +6281234567890", got)
	}
	if got := FormatNomorE164("+6281234567890"); got != "+6281234567890" {
		t.Errorf("FormatNomorE164(+62812...) = %s, harapan +6281234567890", got)
	}
	if got := FormatNomorE164("0812-3456-7890"); got != "+6281234567890" {
		t.Errorf("FormatNomorE164(0812-3456...) = %s, harapan +6281234567890", got)
	}
}
