package main

import (
	"log"
	"net/http"
	"os"

	admin "rengginang-putra-jaya/auth/1.1.1.3"
)

func main() {
	// Menghubungkan Golang ke database MySQL.
	if err := hubungkanDatabase(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Endpoint untuk memeriksa kondisi server.
	http.HandleFunc("/api/kesehatan", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			kirimJSON(
				w,
				http.StatusMethodNotAllowed,
				map[string]string{
					"pesan": "Metode tidak diizinkan",
				},
			)
			return
		}

		kirimJSON(
			w,
			http.StatusOK,
			map[string]string{
				"status": "ok",
				"pesan":  "Backend Rengginang Putra Jaya berjalan",
			},
		)
	})

	// WBS 1.1.1.2
	// API login dengan pembatasan percobaan
	// dan pemeriksaan asal permintaan.
	http.Handle(
		"/api/pemilik/masuk",
		middlewarePembatasLogin(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(tanganiMasuk),
			),
		),
	)

	// WBS 1.1.1.2
	// API logout dengan pemeriksaan asal permintaan.
	http.Handle(
		"/api/pemilik/keluar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(tanganiKeluar),
		),
	)

	// WBS 1.1.1.2
	// Endpoint profil hanya dapat diakses setelah login.
	http.Handle(
		"/api/pemilik/profil",
		middlewareSesiPemilik(
			http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.Method != http.MethodGet {
					w.Header().Set("Allow", "GET")
					kirimJSON(
						w,
						http.StatusMethodNotAllowed,
						map[string]string{
							"pesan": "Metode tidak diizinkan",
						},
					)
					return
				}

				idPemilik, ada := ambilIDPemilik(r)
				if !ada {
					kirimJSON(
						w,
						http.StatusUnauthorized,
						map[string]string{
							"pesan": "Sesi pemilik tidak valid",
						},
					)
					return
				}

				kirimJSON(
					w,
					http.StatusOK,
					map[string]interface{}{
						"pesan":      "Sesi pemilik valid",
						"id_pemilik": idPemilik,
					},
				)
			}),
		),
	)

	// WBS 1.1.1.3
	// Endpoint beranda admin dilindungi middleware sesi Owner.
	http.Handle(
		"/api/admin/beranda",
		middlewareSesiPemilik(
			http.HandlerFunc(admin.TanganiBerandaAdmin),
		),
	)

	// Mengambil port dari konfigurasi .env.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// WBS 1.1.1.3
	// Mengaktifkan middleware CORS untuk komunikasi frontend.
	penanganAPI := admin.MiddlewareCORS(http.DefaultServeMux)

	// Menjalankan server backend.
	log.Println(
		"Server berjalan di http://localhost:" + port,
	)

	log.Fatal(
		http.ListenAndServe(
			"127.0.0.1:"+port,
			penanganAPI,
		),
	)
}
