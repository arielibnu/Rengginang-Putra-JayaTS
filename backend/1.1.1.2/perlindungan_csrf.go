package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Memeriksa asal permintaan dari browser.
func middlewarePerlindunganCSRF(
	berikutnya http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method == http.MethodGet ||
			r.Method == http.MethodHead ||
			r.Method == http.MethodOptions {
			berikutnya.ServeHTTP(w, r)
			return
		}

		asalDiizinkan := os.Getenv("ASAL_FRONTEND")
		if asalDiizinkan == "" {
			asalDiizinkan = "http://localhost:5173"
		}

		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		// Tolak jika kedua header tidak tersedia.
		if origin == "" && referer == "" {
			kirimJSON(w, http.StatusForbidden,
				map[string]string{
					"pesan": "Asal permintaan tidak diketahui",
				})
			return
		}

		asal := origin
		if asal == "" {
			asal = referer
		}

		alamat, err := url.Parse(asal)
		if err != nil {
			kirimJSON(w, http.StatusForbidden,
				map[string]string{
					"pesan": "Asal permintaan tidak valid",
				})
			return
		}

		alamatAsal := strings.TrimSuffix(
			alamat.Scheme+"://"+alamat.Host, "/",
		)

		if alamatAsal != asalDiizinkan {
			kirimJSON(w, http.StatusForbidden,
				map[string]string{
					"pesan": "Permintaan ditolak karena asal tidak diizinkan",
				})
			return
		}

		berikutnya.ServeHTTP(w, r)
	})
}
