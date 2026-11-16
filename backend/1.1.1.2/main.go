
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	admin "rengginang-putra-jaya/auth/1.1.1.3"
	pelanggan "rengginang-putra-jaya/auth/1.1.2.2"
	produkapi "rengginang-putra-jaya/auth/1.1.3.2"
	snapshot "rengginang-putra-jaya/auth/1.1.3.4"
)

func main() {
	// Koneksi MySQL
	if err := hubungkanDatabase(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cek kondisi server
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
	// PPL-01: LOGIN OWNER
	// ==========================================

	http.Handle(
		"/api/pemilik/masuk",
		middlewarePembatasLogin(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(tanganiMasuk),
			),
		),
	)

	http.Handle(
		"/api/pemilik/keluar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(tanganiKeluar),
		),
	)

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
	// PPL-01: HAK AKSES ADMIN
	// ==========================================

	http.Handle(
		"/api/admin/beranda",
		middlewareSesiPemilik(
			http.HandlerFunc(admin.TanganiBerandaAdmin),
		),
	)

	// ==========================================
	// PPL-02: REGISTRASI DAN LOGIN PELANGGAN
	// WBS 1.1.2.2
	// ==========================================

	handlerPelanggan := &pelanggan.Handler{
		DB: db,
	}

	// POST - Daftar pelanggan
	http.Handle(
		"/api/pelanggan/daftar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(handlerPelanggan.Registrasi),
		),
	)

	// POST - Login pelanggan
	http.Handle(
		"/api/pelanggan/masuk",
		middlewarePembatasLogin(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(handlerPelanggan.Login),
			),
		),
	)

	// GET - Profil pelanggan
	http.Handle(
		"/api/pelanggan/profil",
		http.HandlerFunc(handlerPelanggan.Profil),
	)

	// POST - Logout pelanggan
	http.Handle(
		"/api/pelanggan/keluar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(handlerPelanggan.Logout),
		),
	)

	// ==========================================
	// PPL-03: CRUD PRODUK DAN VARIAN
	// WBS 1.1.3.2
	// ==========================================

	handlerProduk := &produkapi.HandlerProduk{
		DB: db,
	}

	handlerVarian := &produkapi.HandlerVarian{
		DB: db,
	}

	// GET / POST produk
	http.Handle(
		"/api/admin/produk",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerProduk,
			),
		),
	)

	// GET / PUT / DELETE produk
	// GET / POST varian berdasarkan produk
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

					handlerProduk.ServeHTTP(w, r)
				}),
			),
		),
	)

	// PUT / DELETE varian
	http.Handle(
		"/api/admin/varian/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerVarian,
			),
		),
	)

	// ==========================================
	// PPL-03: SNAPSHOT HARGA TRANSAKSI
	// WBS 1.1.3.4
	// ==========================================

	handlerSnapshot := &snapshot.Handler{
		DB: db,
	}

	// POST transaksi
	http.Handle(
		"/api/admin/transaksi",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	// GET riwayat transaksi
	http.Handle(
		"/api/admin/transaksi/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	// ==========================================
	// SERVER UTAMA
	// ==========================================

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORS untuk frontend React
	penanganAPI := admin.MiddlewareCORS(
		http.DefaultServeMux,
	)

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