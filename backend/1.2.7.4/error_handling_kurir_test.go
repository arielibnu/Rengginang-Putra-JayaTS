package errorhandlingkurir

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBangunErrorKurir_StrukturGambar2DanTanpaEstimasiPalsu(t *testing.T) {
	errResp := BangunErrorKurir(KegagalanAPIDown, "Jember, Jawa Timur", "")

	if errResp.Status != "error" {
		t.Errorf("Status harus 'error', didapat: %s", errResp.Status)
	}

	if errResp.ErrorCode != "TARIF_TIDAK_TERSEDIA" {
		t.Errorf("ErrorCode harus 'TARIF_TIDAK_TERSEDIA', didapat: %s", errResp.ErrorCode)
	}

	if errResp.Pesan != "Tarif pengiriman tidak tersedia" {
		t.Errorf("Pesan utama harus 'Tarif pengiriman tidak tersedia', didapat: %s", errResp.Pesan)
	}

	if errResp.Peringatan != "Jangan membuat tarif ongkir perkiraan sendiri." {
		t.Errorf("Peringatan harus 'Jangan membuat tarif ongkir perkiraan sendiri.', didapat: %s", errResp.Peringatan)
	}

	// Aturan Utama WBS 1.2.7.4: DILARANG MEMBUAT TARIF PALSU
	if errResp.EstimasiPalsu != false {
		t.Errorf("EstimasiPalsu WAJIB bernilai false, didapat: %v", errResp.EstimasiPalsu)
	}

	if len(errResp.AksiTindakan) != 2 || errResp.AksiTindakan[0] != "Coba lagi" || errResp.AksiTindakan[1] != "Ubah alamat" {
		t.Errorf("Aksi tindakan harus ['Coba lagi', 'Ubah alamat'], didapat: %v", errResp.AksiTindakan)
	}
}

func TestCekJangkauanWilayah_ValidDanInvalid(t *testing.T) {
	if !CekJangkauanWilayah("Jawa Timur", "Kabupaten Jember") {
		t.Errorf("Jember, Jawa Timur seharusnya terjangkau")
	}

	if CekJangkauanWilayah("", "") {
		t.Errorf("Provinsi dan kota kosong tidak boleh terjangkau")
	}

	if CekJangkauanWilayah("Luar Negeri", "Kuala Lumpur") {
		t.Errorf("Luar Negeri tidak boleh terjangkau")
	}
}

func TestHandlerSimulasiError_EndpointHTTP(t *testing.T) {
	h := &HandlerErrorKurir{DB: nil}

	payload := PermintaanSimulasiKegagalan{
		TipeError:     KegagalanAPIDown,
		WilayahTujuan: "Jember, Jawa Timur",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/kurir/simulasi-error", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandlerSimulasiError(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Status HTTP harus 503 Service Unavailable, didapat: %d", w.Code)
	}

	var resp ResponsErrorKurir
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Gagal decode JSON error: %v", err)
	}

	if resp.Pesan != "Tarif pengiriman tidak tersedia" {
		t.Errorf("Pesan harus sesuai Gambar 2, didapat: %s", resp.Pesan)
	}
}

func TestHandlerCekJangkauan_EndpointHTTP(t *testing.T) {
	h := &HandlerErrorKurir{DB: nil}

	payload := map[string]string{
		"provinsi": "Jawa Timur",
		"kota":     "Kabupaten Jember",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/kurir/cek-jangkauan", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandlerCekJangkauan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status HTTP harus 200 OK, didapat: %d", w.Code)
	}
}
