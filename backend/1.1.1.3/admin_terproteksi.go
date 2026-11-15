package admin

import (
	"encoding/json"
	"net/http"
)

// Handler beranda admin.
// Pemeriksaan sesi dilakukan oleh middleware di aplikasi utama.
func TanganiBerandaAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"pesan": "Metode tidak diizinkan",
		})
		return
	}

	kirimJSON(w, http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"pesan":  "Akses admin diizinkan",
		"data": map[string]string{
			"halaman": "Beranda Admin",
			"akses":   "Owner",
		},
	})
}

// Membuat respons JSON.
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
