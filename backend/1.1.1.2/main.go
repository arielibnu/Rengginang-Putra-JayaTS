package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	admin "rengginang-putra-jaya/auth/1.1.1.3"
	produkapi "rengginang-putra-jaya/auth/1.1.3.2"
	snapshot "rengginang-putra-jaya/auth/1.1.3.4"
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

	// ==========================================
	// WBS 1.1.1.2 - LOGIN OWNER
	// ==========================================

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

	// API logout dengan pemeriksaan asal permintaan.
	http.Handle(
		"/api/pemilik/keluar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(tanganiKeluar),
		),
	)

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

	// ==========================================
	// WBS 1.1.1.3 - HAK AKSES ADMIN
	// ==========================================

	// Endpoint beranda admin dilindungi middleware sesi Owner.
	http.Handle(
		"/api/admin/beranda",
		middlewareSesiPemilik(
			http.HandlerFunc(admin.TanganiBerandaAdmin),
		),
	)

	// ==========================================
	// WBS 1.1.3.2 - CRUD PRODUK DAN VARIAN
	// ==========================================

	// Menggunakan koneksi MySQL dari backend Login Owner.
	handlerProduk := &produkapi.HandlerProduk{
		DB: db,
	}

	handlerVarian := &produkapi.HandlerVarian{
		DB: db,
	}

	// GET  /api/admin/produk
	// POST /api/admin/produk
	http.Handle(
		"/api/admin/produk",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerProduk,
			),
		),
	)

	// GET    /api/admin/produk/{id}
	// PUT    /api/admin/produk/{id}
	// DELETE /api/admin/produk/{id}
	//
	// GET  /api/admin/produk/{id}/varian
	// POST /api/admin/produk/{id}/varian
	http.Handle(
		"/api/admin/produk/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(func(
					w http.ResponseWriter,
					r *http.Request,
				) {
					path := strings.TrimPrefix(
						r.URL.Path,
						"/api/admin/produk/",
					)

					bagian := strings.Split(
						strings.Trim(path, "/"),
						"/",
					)

					// Jika URL mengarah ke varian produk.
					if len(bagian) == 2 &&
						bagian[1] == "varian" {

						idProduk, err := strconv.ParseUint(
							bagian[0],
							10,
							64,
						)

						if err != nil || idProduk == 0 {
							kirimJSON(
								w,
								http.StatusBadRequest,
								map[string]string{
									"pesan": "ID produk tidak valid",
								},
							)
							return
						}

						handlerVarian.HandlerVarianProduk(
							w,
							r,
							idProduk,
						)
						return
					}

					// Selain URL varian, proses sebagai CRUD produk.
					handlerProduk.ServeHTTP(w, r)
				}),
			),
		),
	)

	// PUT    /api/admin/varian/{id}
	// DELETE /api/admin/varian/{id}
	http.Handle(
		"/api/admin/varian/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerVarian,
			),
		),
	)

	// ==========================================
	// WBS 1.1.3.4 - SNAPSHOT HARGA TRANSAKSI
	// ==========================================

	// Menggunakan koneksi MySQL yang sama.
	handlerSnapshot := &snapshot.Handler{
		DB: db,
	}

	// POST /api/admin/transaksi
	// Membuat transaksi dan menyimpan snapshot harga.
	http.Handle(
		"/api/admin/transaksi",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	// GET /api/admin/transaksi/{id}
	// Membaca riwayat transaksi beserta harga snapshot.
	http.Handle(
		"/api/admin/transaksi/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	// ==========================================
	// KONFIGURASI DAN SERVER UTAMA
	// ==========================================

	// Mengambil port dari konfigurasi .env.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
