package gatewayongkir

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHitungOngkirSandbox_SuksesBerdasarkanBerat(t *testing.T) {
	req := &PermintaanHitungOngkir{
		Tujuan: LokasiTujuan{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "Kabupaten Jember",
			Kecamatan:     "Kecamatan Kaliwates",
			KodePos:       "68131",
			AlamatLengkap: "Jl. Contoh No. 12",
			NamaPenerima:  "Pelanggan Contoh",
			NomorWhatsApp: "081234567890",
		},
		BeratGram: 1350, // 2 udang (1000g) + 1 bawang (350g) = 1350g => Pembulatan 2kg
	}

	resp, err := HitungOngkirSandbox(req)
	if err != nil {
		t.Fatalf("HitungOngkirSandbox gagal: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status harus 'sukses', didapat: %s", resp.Status)
	}

	if len(resp.OpsiLayanan) != 3 {
		t.Fatalf("Jumlah opsi layanan harus 3 (JNE, J&T, SiCepat), didapat: %d", len(resp.OpsiLayanan))
	}

	// JNE REG Jember dasar 18.000 * 2kg = 36.000
	jne := resp.OpsiLayanan[0]
	if jne.IDKurir != "jne" || jne.NamaKurir != "JNE REG" {
		t.Errorf("Kurir pertama harus JNE REG, didapat: %s", jne.NamaKurir)
	}
	if jne.Tarif != 36000 {
		t.Errorf("Tarif JNE harus 36000 (18.000 * 2kg), didapat: %d", jne.Tarif)
	}
	if !jne.Terpilih {
		t.Errorf("JNE REG harus default terpilih sesuai mockup Gambar 1")
	}

	// J&T Express Jember dasar 17.000 * 2kg = 34.000
	jnt := resp.OpsiLayanan[1]
	if jnt.IDKurir != "jnt" || jnt.Tarif != 34000 {
		t.Errorf("Tarif J&T harus 34000, didapat: %d", jnt.Tarif)
	}

	// SiCepat Jember dasar 19.000 * 2kg = 38.000
	scp := resp.OpsiLayanan[2]
	if scp.IDKurir != "sicepat" || scp.Tarif != 38000 {
		t.Errorf("Tarif SiCepat harus 38000, didapat: %d", scp.Tarif)
	}
}

func TestHitungOngkirSandbox_DefaultBerat1350(t *testing.T) {
	req := &PermintaanHitungOngkir{
		Tujuan: LokasiTujuan{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "Kota Surabaya",
		},
		BeratGram: 0, // Diberi 0, sistem harus fallback ke berat default 1350g
	}

	resp, err := HitungOngkirSandbox(req)
	if err != nil {
		t.Fatalf("HitungOngkirSandbox gagal: %v", err)
	}

	if resp.BeratTotalGram != 1350 {
		t.Errorf("Berat default harus 1350, didapat: %d", resp.BeratTotalGram)
	}
}

func TestHitungOngkirSandbox_KotaKosong(t *testing.T) {
	req := &PermintaanHitungOngkir{
		Tujuan: LokasiTujuan{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "",
		},
	}

	_, err := HitungOngkirSandbox(req)
	if err == nil {
		t.Errorf("Harus mengembalikan error jika kota kosong")
	}
}

func TestHandlerHitungOngkir_EndpointHTTP(t *testing.T) {
	h := &HandlerGatewayOngkir{DB: nil}

	payload := PermintaanHitungOngkir{
		Tujuan: LokasiTujuan{
			Provinsi:      "Jawa Timur",
			KotaKabupaten: "Kabupaten Jember",
			Kecamatan:     "Kecamatan Kaliwates",
		},
		BeratGram: 1350,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/kurir/hitung-ongkir", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandlerHitungOngkir(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status HTTP harus 200 OK, didapat: %d", w.Code)
	}

	var resp ResponsUmumOngkir
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON response: %v", err)
	}

	if resp.Status != "sukses" {
		t.Errorf("Status response harus sukses, didapat: %s", resp.Status)
	}
}

func TestHandlerPilihKurir_EndpointHTTP(t *testing.T) {
	h := &HandlerGatewayOngkir{DB: nil}

	payload := PermintaanPilihKurir{
		IDKurir:      "jne",
		NamaKurir:    "JNE REG",
		Tarif:        18000,
		EstimasiHari: "2-3 hari",
		AlamatTujuan: "Jember, Jawa Timur",
		NamaPenerima: "Pelanggan Contoh",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/kurir/pilih-layanan", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandlerPilihKurir(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status HTTP harus 200 OK, didapat: %d", w.Code)
	}
}

func TestHandlerDaftarEkspedisi_GET(t *testing.T) {
	h := &HandlerGatewayOngkir{DB: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/kurir/daftar-ekspedisi", nil)
	w := httptest.NewRecorder()

	h.HandlerDaftarEkspedisi(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status HTTP harus 200 OK, didapat: %d", w.Code)
	}
}
