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
	modalapi "rengginang-putra-jaya/auth/1.1.6.3"
	katalogapi "rengginang-putra-jaya/auth/1.1.8.2"
	ketersediaanproduk "rengginang-putra-jaya/auth/1.1.9.3"
	kuotastok "rengginang-putra-jaya/auth/1.2.1.3"
	statemanager "rengginang-putra-jaya/auth/1.2.2.2"
	validasikeranjang "rengginang-putra-jaya/auth/1.2.2.4"
	checkoutapi "rengginang-putra-jaya/auth/1.2.3.3"
	validasiwilayah "rengginang-putra-jaya/auth/1.2.6.3"
	gatewayongkir "rengginang-putra-jaya/auth/1.2.7.2"
	errorhandlingkurir "rengginang-putra-jaya/auth/1.2.7.4"
)

func main() {
	// ==========================================
	// KONEKSI DATABASE
	// ==========================================

	if err := hubungkanDatabase(); err != nil {
		log.Println("Peringatan koneksi MySQL:", err)
		log.Println("Server tetap berjalan dengan mode data in-memory/fallback...")
	} else if db != nil {
		defer db.Close()
	}

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
	// WBS 1.1.6.3 - MODAL DAN PENGELUARAN USAHA
	// ==========================================

	handlerModal := &modalapi.Handler{
		DB: db,
	}

	// GET, POST /api/admin/modal
	http.Handle(
		"/api/admin/modal",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerModal,
			),
		),
	)

	// GET /api/admin/modal/ringkasan
	http.Handle(
		"/api/admin/modal/",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerModal,
			),
		),
	)

	// GET, POST /api/admin/pengeluaran
	http.Handle(
		"/api/admin/pengeluaran",
		middlewareSesiPemilik(
			middlewarePerlindunganCSRF(
				handlerModal,
			),
		),
	)
	// ==========================================
	// WBS 1.1.8.2 - API PUBLIC CATALOG PRODUK
	// ==========================================

	handlerKatalog := &katalogapi.HandlerKatalog{
		DB: db,
	}

	// GET /api/katalog dan GET /api/katalog/produk (Publik tanpa sesi)
	http.Handle(
		"/api/katalog",
		handlerKatalog,
	)

	http.Handle(
		"/api/katalog/",
		handlerKatalog,
	)

	// ==========================================
	// WBS 1.1.9.3 - DETAIL PRODUK & STATUS KETERSEDIAAN
	// ==========================================

	handlerKetersediaan := &ketersediaanproduk.HandlerDetail{
		DB: db,
	}

	// GET /api/katalog/detail (Publik tanpa sesi)
	http.HandleFunc(
		"/api/katalog/detail",
		handlerKetersediaan.HandlerDetailProduk,
	)

	// ==========================================
	// WBS 1.2.1.3 - VALIDASI KUOTA STOK REALTIME (PPL-10)
	// ==========================================

	handlerKuota := &kuotastok.HandlerValidasiKuota{
		DB: db,
	}

	// GET /api/katalog/cek-kuota (Pengecekan kuota realtime saat stepper/jumlah berubah)
	http.HandleFunc(
		"/api/katalog/cek-kuota",
		handlerKuota.HandlerCekKuotaRealtime,
	)

	// POST /api/katalog/validasi-kuota (Validasi kuota sisa stok realtime)
	http.HandleFunc(
		"/api/katalog/validasi-kuota",
		handlerKuota.HandlerCekKuotaRealtime,
	)

	// POST /api/katalog/validasi-beli (Kompatibilitas tombol beli frontend)
	http.HandleFunc(
		"/api/katalog/validasi-beli",
		handlerKuota.HandlerCekKuotaRealtime,
	)

	// ==========================================
	// WBS 1.2.2.2 - STATE MANAGER KERANJANG BELANJA (PPL-11)
	// ==========================================

	handlerStateManager := &statemanager.HandlerStateManager{
		DB: db,
	}

	http.HandleFunc(
		"/api/keranjang/sinkronisasi",
		handlerStateManager.HandlerSinkronisasiKeranjang,
	)
	http.HandleFunc(
		"/api/keranjang",
		handlerStateManager.HandlerAmbilKeranjang,
	)
	http.HandleFunc(
		"/api/keranjang/item",
		handlerStateManager.HandlerHapusItemKeranjang,
	)
	http.HandleFunc(
		"/api/keranjang/bersihkan",
		handlerStateManager.HandlerBersihkanKeranjang,
	)

	// ==========================================
	// WBS 1.2.2.4 - API VALIDASI KERANJANG TERHADAP STOK (PPL-11)
	// ==========================================

	handlerValidasiKeranjang := &validasikeranjang.HandlerValidasiKeranjang{
		DB: db,
	}

	http.HandleFunc(
		"/api/keranjang/validasi",
		handlerValidasiKeranjang.HandlerValidasiKeranjang,
	)
	http.HandleFunc(
		"/api/keranjang/validasi-stok",
		handlerValidasiKeranjang.HandlerValidasiKeranjang,
	)

	// ==========================================
	// WBS 1.2.3.3 - PENYAMAAN FORMAT PAYLOAD CHECKOUT (PPL-12)
	// ==========================================

	handlerCheckout := &checkoutapi.HandlerCheckout{
		DB: db,
	}

	// POST /api/checkout (Order Request Standar Guest & Pelanggan)
	http.HandleFunc(
		"/api/checkout",
		handlerCheckout.HandlerProsesCheckout,
	)

	// POST /api/checkout/proses (Alias endpoint pemrosesan pesanan)
	http.HandleFunc(
		"/api/checkout/proses",
		handlerCheckout.HandlerProsesCheckout,
	)

	// GET /api/checkout/skema (Spesifikasi skema format payload order request)
	http.HandleFunc(
		"/api/checkout/skema",
		handlerCheckout.HandlerSkemaPayload,
	)

	// ==========================================
	// WBS 1.2.6.3 - VALIDASI WILAYAH & ALAMAT PENGIRIMAN SESUAI KEBUTUHAN API KURIR (PPL-15)
	// ==========================================

	handlerWilayah := &validasiwilayah.HandlerValidasiWilayah{
		DB: db,
	}

	// POST /api/checkout/wilayah-pengiriman (Endpoint integrasi Frontend Checkout PPL-15)
	http.HandleFunc(
		"/api/checkout/wilayah-pengiriman",
		handlerWilayah.HandlerSimpanDanValidasiWilayah,
	)

	// POST /api/checkout/validasi-alamat (Endpoint validasi data geografis kurir)
	http.HandleFunc(
		"/api/checkout/validasi-alamat",
		handlerWilayah.HandlerValidasiAlamatKurir,
	)

	// GET /api/checkout/master-wilayah (Master data dropdown bertingkat provinsi -> kota -> kecamatan)
	http.HandleFunc(
		"/api/checkout/master-wilayah",
		handlerWilayah.HandlerMasterWilayah,
	)

	// GET /api/checkout/parameter-kurir (Skema parameter kebutuhan API kurir)
	http.HandleFunc(
		"/api/checkout/parameter-kurir",
		handlerWilayah.HandlerParameterKurir,
	)

	// ==========================================
	// WBS 1.2.7.2 - SERVICE GATEWAY API ONGKIR SANDBOX BERDASARKAN BERAT (PPL-16)
	// ==========================================

	handlerOngkir := &gatewayongkir.HandlerGatewayOngkir{
		DB: db,
	}

	// POST /api/kurir/hitung-ongkir (Endpoint utama hitung ongkir kurir sandbox berdasar berat)
	http.HandleFunc(
		"/api/kurir/hitung-ongkir",
		handlerOngkir.HandlerHitungOngkir,
	)
	http.HandleFunc(
		"/api/checkout/hitung-ongkir",
		handlerOngkir.HandlerHitungOngkir,
	)

	// GET /api/kurir/daftar-ekspedisi (Daftar opsi kurir mitra: JNE, J&T, SiCepat)
	http.HandleFunc(
		"/api/kurir/daftar-ekspedisi",
		handlerOngkir.HandlerDaftarEkspedisi,
	)

	// POST /api/kurir/pilih-layanan (Menyimpan kurir yang dipilih user ke database saat 'Gunakan layanan')
	http.HandleFunc(
		"/api/kurir/pilih-layanan",
		handlerOngkir.HandlerPilihKurir,
	)

	// GET /api/kurir/riwayat (Riwayat audit kalkulasi ongkos kirim)
	http.HandleFunc(
		"/api/kurir/riwayat",
		handlerOngkir.HandlerRiwayatKalkulasi,
	)

	// ==========================================
	// WBS 1.2.7.4 - ERROR HANDLING KONEKSI API KURIR & RUTE TANPA ESTIMASI PALSU (PPL-16)
	// ==========================================

	handlerErrorKurir := &errorhandlingkurir.HandlerErrorKurir{
		DB: db,
	}

	// POST /api/kurir/simulasi-error (Simulasi kegagalan API kurir sesuai Gambar 2)
	http.HandleFunc(
		"/api/kurir/simulasi-error",
		handlerErrorKurir.HandlerSimulasiError,
	)

	// POST /api/kurir/cek-jangkauan (Pemeriksaan keterjangkauan rute pengiriman)
	http.HandleFunc(
		"/api/kurir/cek-jangkauan",
		handlerErrorKurir.HandlerCekJangkauan,
	)

	// GET /api/kurir/log-error (Audit log insiden kegagalan koneksi API kurir)
	http.HandleFunc(
		"/api/kurir/log-error",
		handlerErrorKurir.HandlerLogKegagalan,
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
