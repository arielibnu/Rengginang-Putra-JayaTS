package errorhandlingkurir

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

// In-Memory Storage untuk audit log error jika MySQL offline
var (
	memoriErrorLog = make([]ResponsErrorKurir, 0)
	muError        sync.RWMutex
)

// HandlerErrorKurir menampung koneksi DB
type HandlerErrorKurir struct {
	DB *sql.DB
}

// BangunErrorKurir membuat payload error resmi yang 100% presisi sesuai desain Figma Gambar 2
// dan memenuhi aturan ketat WBS 1.2.7.4: "tanpa membuat estimasi tarif palsu" (EstimasiPalsu = false)
func BangunErrorKurir(tipe TipeKegagalanKurir, wilayah string, detailKustom string) *ResponsErrorKurir {
	detail := "API gagal atau layanan tidak melayani wilayah tujuan."
	if detailKustom != "" {
		detail = detailKustom
	}

	return &ResponsErrorKurir{
		Status:        "error",
		ErrorCode:     "TARIF_TIDAK_TERSEDIA",
		Pesan:         "Tarif pengiriman tidak tersedia",
		Detail:        detail,
		Peringatan:    "Jangan membuat tarif ongkir perkiraan sendiri.",
		TipeError:     tipe,
		EstimasiPalsu: false, // DILARANG KERAS MEMBUAT TARIF PALSU
		WilayahTujuan: wilayah,
		AksiTindakan:  []string{"Coba lagi", "Ubah alamat"},
		WaktuError:    time.Now(),
	}
}

// CekJangkauanWilayah memeriksa apakah suatu wilayah dapat dijangkau oleh rute ekspedisi
func CekJangkauanWilayah(provinsi, kota string) bool {
	provClean := strings.TrimSpace(provinsi)
	kotaClean := strings.TrimSpace(kota)

	if provClean == "" || kotaClean == "" {
		return false
	}

	// Wilayah terpencil atau luar negeri yang tidak terjangkau kurir standar Situbondo
	ruteBlacklist := []string{
		"Luar Negeri",
		"Wilayah Terpencil",
		"Zona Non-Coverage",
	}

	for _, b := range ruteBlacklist {
		if strings.EqualFold(provClean, b) || strings.EqualFold(kotaClean, b) {
			return false
		}
	}

	return true
}

// SimpanAuditError menyimpan log insiden kegagalan API kurir ke database
func (h *HandlerErrorKurir) SimpanAuditError(errResp *ResponsErrorKurir) error {
	muError.Lock()
	memoriErrorLog = append(memoriErrorLog, *errResp)
	muError.Unlock()

	if h.DB == nil {
		return nil
	}

	query := `INSERT INTO log_kegagalan_api_kurir 
		(tipe_error, wilayah_tujuan, pesan_error, detail, dibuat_pada)
		VALUES (?, ?, ?, ?, NOW())`

	_, err := h.DB.Exec(query,
		string(errResp.TipeError),
		errResp.WilayahTujuan,
		errResp.Pesan,
		errResp.Detail,
	)
	return err
}

// HandlerSimulasiError menangani POST /api/kurir/simulasi-error (WBS 1.2.7.4)
// Endpoint ini digunakan untuk menguji respon error dan memastikan Frontend menampilkan Gambar 2
func (h *HandlerErrorKurir) HandlerSimulasiError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSONError(w, http.StatusMethodNotAllowed, map[string]string{
			"status": "error",
			"pesan":  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	var req PermintaanSimulasiKegagalan
	_ = json.NewDecoder(r.Body).Decode(&req)

	tipe := req.TipeError
	if tipe == "" {
		tipe = KegagalanAPIDown
	}

	wilayah := req.WilayahTujuan
	if wilayah == "" {
		wilayah = "Jember, Jawa Timur"
	}

	errResp := BangunErrorKurir(tipe, wilayah, "API gagal atau layanan tidak melayani wilayah tujuan.")
	_ = h.SimpanAuditError(errResp)

	// Kembalikan status HTTP 503 Service Unavailable sesuai spesifikasi Gambar 2
	kirimJSONError(w, http.StatusServiceUnavailable, errResp)
}

// HandlerCekJangkauan menangani POST /api/kurir/cek-jangkauan
func (h *HandlerErrorKurir) HandlerCekJangkauan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	type reqJangkauan struct {
		Provinsi string `json:"provinsi"`
		Kota     string `json:"kota"`
	}

	var req reqJangkauan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimJSONError(w, http.StatusBadRequest, map[string]string{
			"status": "error",
			"pesan":  "Format payload tidak valid",
		})
		return
	}

	terjangkau := CekJangkauanWilayah(req.Provinsi, req.Kota)
	if !terjangkau {
		errResp := BangunErrorKurir(KegagalanRuteTidakTerjangkau, req.Kota+", "+req.Provinsi, "Layanan ekspedisi tidak melayani pengiriman ke wilayah tujuan ini.")
		_ = h.SimpanAuditError(errResp)
		kirimJSONError(w, http.StatusUnprocessableEntity, errResp)
		return
	}

	kirimJSONError(w, http.StatusOK, map[string]interface{}{
		"status":     "sukses",
		"pesan":      "Wilayah tujuan terjangkau oleh jaringan kurir mitra.",
		"terjangkau": true,
	})
}

// HandlerLogKegagalan menangani GET /api/kurir/log-error
func (h *HandlerErrorKurir) HandlerLogKegagalan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSONError(w, http.StatusMethodNotAllowed, map[string]string{
			"status": "error",
			"pesan":  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	muError.RLock()
	defer muError.RUnlock()

	kirimJSONError(w, http.StatusOK, map[string]interface{}{
		"status": "sukses",
		"pesan":  "Audit log kegagalan koneksi API kurir",
		"data":   memoriErrorLog,
	})
}

func kirimJSONError(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
