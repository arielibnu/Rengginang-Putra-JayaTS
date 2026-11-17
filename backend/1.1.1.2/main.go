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
	stokapi "rengginang-putra-jaya/auth/1.1.4.2"
	peringatanstok "rengginang-putra-jaya/auth/1.1.5.2"
)

func main() {
	// ==========================================
	// KONEKSI DATABASE
	// ==========================================

	if err := hubungkanDatabase(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ==========================================
	// ENDPOINT KESEHATAN SERVER
	// ==========================================

	http.HandleFunc("/api/kesehatan", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			kirimJSON(w, http.StatusMethodNotAllowed,
				map[string]string{
					"pesan": "Metode tidak diizinkan",
				})
			return
		}

		kirimJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"pesan":  "Backend Rengginang Putra Jaya berjalan",
		})
	})

	// ==========================================
	// WBS 1.1.1.2 - LOGIN OWNER
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
					kirimJSON(w, http.StatusMethodNotAllowed,
						map[string]string{
							"pesan": "Metode tidak diizinkan",
						})
					return
				}

				idPemilik, ada := ambilIDPemilik(r)
				if !ada {
					kirimJSON(w, http.StatusUnauthorized,
						map[string]string{
							"pesan": "Sesi pemilik tidak valid",
						})
					return
				}

				kirimJSON(w, http.StatusOK,
					map[string]interface{}{
						"pesan":      "Sesi pemilik valid",
						"id_pemilik": idPemilik,
					})
			}),
		),
	)

	// ==========================================
	// WBS 1.1.1.3 - HAK AKSES ADMIN
	// ==========================================

	http.Handle(
		"/api/admin/beranda",
		middlewareSesiPemilik(
			http.HandlerFunc(admin.TanganiBerandaAdmin),
		),
	)

	// ==========================================
	// WBS 1.1.2.2 - REGISTRASI DAN LOGIN PELANGGAN
	// ==========================================

	handlerPelanggan := &pelanggan.Handler{
		DB: db,
	}

	http.Handle(
		"/api/pelanggan/daftar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(handlerPelanggan.Registrasi),
		),
	)

	http.Handle(
		"/api/pelanggan/masuk",
		middlewarePembatasLogin(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(handlerPelanggan.Login),
			),
		),
	)

	http.Handle(
		"/api/pelanggan/profil",
		http.HandlerFunc(handlerPelanggan.Profil),
	)

	http.Handle(
		"/api/pelanggan/keluar",
		middlewarePerlindunganCSRF(
			http.HandlerFunc(handlerPelanggan.Logout),
		),
	)

	// ==========================================
	// WBS 1.1.3.2 - CRUD PRODUK DAN VARIAN
	// ==========================================

	handlerProduk := &produkapi.HandlerProduk{
		DB: db,
	}

	handlerVarian := &produkapi.HandlerVarian{
		DB: db,
	}

	// GET dan POST produk.
	http.Handle(
		"/api/admin/produk",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerProduk,
			),
		),
	)

	// GET, PUT, DELETE produk dan pengelolaan varian.
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

					// Endpoint varian berdasarkan ID produk.
					if len(bagian) == 2 &&
						bagian[1] == "varian" {

						idProduk, err := strconv.ParseUint(
							bagian[0],
							10,
							64,
						)

						if err != nil || idProduk == 0 {
							kirimJSON(w,
								http.StatusBadRequest,
								map[string]string{
									"pesan": "ID produk tidak valid",
								})
							return
						}

						handlerVarian.HandlerVarianProduk(
							w,
							r,
							idProduk,
						)
						return
					}

					// Endpoint CRUD produk.
					handlerProduk.ServeHTTP(w, r)
				}),
			),
		),
	)

	// PUT dan DELETE varian.
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

	handlerSnapshot := &snapshot.Handler{
		DB: db,
	}

	http.Handle(
		"/api/admin/transaksi",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	http.Handle(
		"/api/admin/transaksi/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerSnapshot,
			),
		),
	)

	// ==========================================
	// WBS 1.1.4.2 - MUTASI DAN RIWAYAT STOK
	// ==========================================

	handlerStok := &stokapi.Handler{
		DB: db,
	}

	// POST /api/admin/stok/mutasi
	// GET  /api/admin/stok/riwayat
	http.Handle(
		"/api/admin/stok/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerStok,
			),
		),
	)

	// ==========================================
	// WBS 1.1.5.2 - PERINGATAN STOK MENIPIS & BATAS MINIMUM
	// ==========================================

	handlerPeringatan := &peringatanstok.Handler{
		DB: db,
	}

	// GET, POST, PUT /api/admin/stok/batas-minimum
	http.Handle(
		"/api/admin/stok/batas-minimum",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(handlerPeringatan.HandlerBatasMinimum),
			),
		),
	)

	// GET /api/admin/stok/peringatan
	http.Handle(
		"/api/admin/stok/peringatan",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				http.HandlerFunc(handlerPeringatan.HandlerPeringatanStok),
			),
		),
	)

	// ==========================================
	// KONFIGURASI SERVER
	// ==========================================

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Middleware CORS untuk komunikasi frontend.
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
