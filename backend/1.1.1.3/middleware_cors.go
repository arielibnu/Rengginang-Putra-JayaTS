package admin

import (
	"net/http"
	"os"
)

// MiddlewareCORS mengatur akses frontend ke API backend.
func MiddlewareCORS(berikutnya http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asalDiizinkan := os.Getenv("ASAL_FRONTEND")
		if asalDiizinkan == "" {
			asalDiizinkan = "http://localhost:5173"
		}

		asalPermintaan := r.Header.Get("Origin")

		// Memeriksa asal permintaan dari browser.
		if asalPermintaan != "" {
			if asalPermintaan != asalDiizinkan {
				http.Error(
					w,
					"Asal permintaan tidak diizinkan",
					http.StatusForbidden,
				)
				return
			}

			// Mengizinkan frontend yang sudah ditentukan.
			w.Header().Set(
				"Access-Control-Allow-Origin",
				asalDiizinkan,
			)

			// Mengizinkan cookie sesi Owner.
			w.Header().Set(
				"Access-Control-Allow-Credentials",
				"true",
			)

			// Mengizinkan seluruh metode API CRUD.
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)

			// Mengizinkan header JSON dari frontend.
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type",
			)

			w.Header().Add("Vary", "Origin")
		}

		// Menangani permintaan preflight dari browser.
		if r.Method == http.MethodOptions {
			if asalPermintaan == "" {
				http.Error(
					w,
					"Origin diperlukan",
					http.StatusForbidden,
				)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Meneruskan permintaan ke endpoint berikutnya.
		berikutnya.ServeHTTP(w, r)
	})
}
