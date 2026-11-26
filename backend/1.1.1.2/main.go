package main

import (
	"fmt"
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
	validatororder "rengginang-putra-jaya/auth/1.2.8.3"
	ordergenerator "rengginang-putra-jaya/auth/1.2.9.1"
	reservasisstok "rengginang-putra-jaya/auth/1.2.9.2"
	pesananowner "rengginang-putra-jaya/auth/1.2.10.2"
	tagihanmidtrans "rengginang-putra-jaya/auth/1.3.1.2"
	sinkronisasimidtrans "rengginang-putra-jaya/auth/1.3.1.4"
	webhookgateway "rengginang-putra-jaya/auth/1.3.2.1"
	signatureidempotency "rengginang-putra-jaya/auth/1.3.2.2"
	pembaruanstatus "rengginang-putra-jaya/auth/1.3.2.3"
	kedaluwarsastok "rengginang-putra-jaya/auth/1.3.2.4"
	notifikasipesanan "rengginang-putra-jaya/auth/1.3.3.2"
	statemachine "rengginang-putra-jaya/auth/1.3.4.2"
	pengembalianstok "rengginang-putra-jaya/auth/1.3.4.4"
	penjualanoffline "rengginang-putra-jaya/auth/1.3.5.3"
	flagpemisah "rengginang-putra-jaya/auth/1.3.5.4"
	laporanbulanan "rengginang-putra-jaya/auth/1.4.2.1"
	laporanpengeluaran "rengginang-putra-jaya/auth/1.4.4.1"
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
	// WBS 1.2.8.3 - VALIDASI KETAT INTEGRITAS SUBTOTAL, ONGKIR, TOTAL NOMINAL (PPL-17)
	// ==========================================

	handlerValidatorOrder := &validatororder.HandlerValidatorOrder{
		DB: db,
	}
	handlerValidatorOrder.InisialisasiTabel()

	// POST /api/checkout/validasi-ringkasan (Server-side Order Validator)
	http.HandleFunc(
		"/api/checkout/validasi-ringkasan",
		handlerValidatorOrder.HandlerValidasiRingkasan,
	)
	http.HandleFunc(
		"/api/checkout/validasi-order",
		handlerValidatorOrder.HandlerValidasiRingkasan,
	)

	// POST /api/checkout/buat-pesanan (Eksekusi pesanan final setelah validasi integritas)
	http.HandleFunc(
		"/api/checkout/buat-pesanan",
		handlerValidatorOrder.HandlerBuatPesanan,
	)

	// GET /api/checkout/audit-validasi (Audit log integritas validasi order)
	http.HandleFunc(
		"/api/checkout/audit-validasi",
		handlerValidatorOrder.HandlerRiwayatValidasi,
	)

	// ==========================================
	// WBS 1.2.9 - PEMBUATAN PESANAN DAN RESERVASI STOK (PPL-18)
	// ==========================================

	generatorOrder := &ordergenerator.GeneratorOrder{
		DB: db,
	}
	generatorOrder.InisialisasiTabel()

	handlerReservasiStok := &reservasisstok.HandlerReservasiStok{
		DB:        db,
		Generator: generatorOrder,
	}
	handlerReservasiStok.InisialisasiTabel()

	// POST /api/pesanan/buat (Pembuatan pesanan resmi, snapshot harga, dan penahanan stok 24 jam)
	http.HandleFunc(
		"/api/pesanan/buat",
		handlerReservasiStok.HandlerBuatPesanan,
	)

	// GET /api/pesanan/detail (Detail pesanan asli untuk Halaman Konfirmasi PPL-18 & tracking)
	http.HandleFunc(
		"/api/pesanan/detail",
		handlerReservasiStok.HandlerDetailPesanan,
	)

	// GET /api/pesanan/reservasi (Audit list kuota stok yang sedang ditahan)
	http.HandleFunc(
		"/api/pesanan/reservasi",
		handlerReservasiStok.HandlerDaftarReservasi,
	)

	// ==========================================
	// WBS 1.2.10.2 - DAFTAR & DETAIL PESANAN OWNER (PPL-21)
	// ==========================================

	layananPesananOwner := pesananowner.InisialisasiLayanan(db)

	// GET /api/admin/pesanan (Daftar antrean pesanan masuk dengan filter status, tanggal, dan pencarian)
	http.HandleFunc(
		"/api/admin/pesanan",
		layananPesananOwner.HandlerDaftarPesanan,
	)

	// GET /api/admin/pesanan/detail (Rincian pesanan pembeli dan snapshot barang)
	http.HandleFunc(
		"/api/admin/pesanan/detail",
		layananPesananOwner.HandlerDetailPesanan,
	)

	// ==========================================
	// WBS 1.3.4.2 & 1.3.4.4 - STATE MACHINE PEMROSESAN & PENGEMBALIAN STOK (PPL-23)
	// ==========================================
	layananPengembalianStok := pengembalianstok.InisialisasiLayanan(db)
	layananStateMachine := statemachine.InisialisasiLayanan(db, layananPengembalianStok)

	// PUT / POST /api/admin/pesanan/status (Pembaruan status pemrosesan dengan State Machine formal & terpisah dari pembayaran)
	http.HandleFunc(
		"/api/admin/pesanan/status",
		layananStateMachine.HandlerUbahStatus,
	)

	// GET /api/admin/pesanan/transisi-status (Pilihan transisi status yang diizinkan untuk formulir frontend)
	http.HandleFunc(
		"/api/admin/pesanan/transisi-status",
		layananStateMachine.HandlerOpsiTransisi,
	)

	// POST /api/admin/pesanan/kembalikan-stok (Endpoint pengembalian stok pembatalan WBS 1.3.4.4)
	http.HandleFunc(
		"/api/admin/pesanan/kembalikan-stok",
		layananPengembalianStok.HandlerPengembalianStok,
	)

	// ==========================================
	// WBS 1.3.1.2 & 1.3.1.4 - MIDTRANS SANDBOX PAYMENT GATEWAY (PPL-19)
	// ==========================================

	layananTagihanMidtrans := tagihanmidtrans.InisialisasiLayanan(db)
	layananSinkronisasiMidtrans := sinkronisasimidtrans.InisialisasiLayanan(db)

	// POST /api/pembayaran/buat-tagihan (WBS 1.3.1.2)
	http.HandleFunc(
		"/api/pembayaran/buat-tagihan",
		layananTagihanMidtrans.HandlerBuatTagihan,
	)

	// GET /api/pembayaran/status (WBS 1.3.1.4)
	http.HandleFunc(
		"/api/pembayaran/status",
		layananSinkronisasiMidtrans.HandlerStatusPembayaran,
	)

	// POST /api/pembayaran/webhook (WBS 1.3.1.4)
	http.HandleFunc(
		"/api/pembayaran/webhook",
		layananSinkronisasiMidtrans.HandlerWebhook,
	)

	// POST /api/pembayaran/simulasi-bayar (WBS 1.3.1.4 Simulasi)
	http.HandleFunc(
		"/api/pembayaran/simulasi-bayar",
		layananSinkronisasiMidtrans.HandlerSimulasiBayar,
	)

	// ==========================================
	// WBS 1.3.3.2 - NOTIFIKASI PESANAN BARU OWNER (PPL-22)
	// ==========================================

	layananNotifikasi := notifikasipesanan.InisialisasiLayanan(db)

	// GET /api/owner/notifikasi/pesanan-baru (Komponen Notifikasi Pesanan Baru Figma PPL-22)
	http.HandleFunc(
		"/api/owner/notifikasi/pesanan-baru",
		layananNotifikasi.HandlerPesananBaru,
	)

	// GET /api/owner/notifikasi (Daftar semua notifikasi pesanan)
	http.HandleFunc(
		"/api/owner/notifikasi",
		layananNotifikasi.HandlerDaftarSemua,
	)

	// POST /api/owner/notifikasi/baca (Tandai dibaca saat Owner membuka pesanan)
	http.HandleFunc(
		"/api/owner/notifikasi/baca",
		layananNotifikasi.HandlerTandaiDibaca,
	)

	// POST /api/owner/notifikasi/trigger (Trigger pemberitahuan pesanan baru / pembayaran)
	http.HandleFunc(
		"/api/owner/notifikasi/trigger",
		layananNotifikasi.HandlerTrigger,
	)

	// ==========================================
	// WBS 1.3.2.1 s/d 1.3.2.4 - PAYMENT GATEWAY WEBHOOK, IDEMPOTENCY, STATUS & KEDALUWARSA (PPL-20)
	// ==========================================

	layananWebhook := webhookgateway.InisialisasiLayanan(db)
	layananKeamananIdem := signatureidempotency.InisialisasiLayanan(db)
	layananPembaruanStatus := pembaruanstatus.InisialisasiLayanan(db)
	layananKedaluwarsaStok := kedaluwarsastok.InisialisasiLayanan(db)

	// Callback orkestrator terpadu saat notifikasi webhook Midtrans diterima
	layananWebhook.OnNotifikasi = func(notif webhookgateway.NotifikasiWebhookMidtrans) error {
		// 1. Verifikasi Keamanan Signature Key SHA512 (WBS 1.3.2.2)
		if notif.SignatureKey != "" {
			hasilSig := layananKeamananIdem.VerifikasiSignatureKey(notif.OrderID, notif.StatusCode, notif.GrossAmount, notif.SignatureKey)
			if !hasilSig.Valid {
				log.Printf("[PPL-20] DITOLAK: Signature Key tidak valid untuk pesanan: %s", notif.OrderID)
				return fmt.Errorf("signature key tidak valid")
			}
		}

		// 2. Proteksi Pemrosesan Ganda / Idempotency (WBS 1.3.2.2)
		cekIdem, err := layananKeamananIdem.PeriksaDanKunciIdempotency(notif.OrderID, notif.TransactionID, notif.TransactionStatus)
		if err != nil {
			log.Printf("[PPL-20] Peringatan cek idempotency: %v", err)
		}
		if cekIdem != nil && cekIdem.AdalahDuplikat && !cekIdem.BolehDiproses {
			log.Printf("[PPL-20] IDEMPOTENT: Notifikasi pesanan %s status %s sudah diproses sebelumnya. Pemrosesan duplikat dicegah.", notif.OrderID, notif.TransactionStatus)
			return nil
		}

		// 3. Pembaruan Otomatis Status Transaksi (WBS 1.3.2.3)
		reqStatus := pembaruanstatus.PembaruanStatusRequest{
			NomorPesanan:      notif.OrderID,
			TransactionID:     notif.TransactionID,
			TransactionStatus: notif.TransactionStatus,
			FraudStatus:       notif.FraudStatus,
			PaymentType:       notif.PaymentType,
			SettlementTime:    notif.SettlementTime,
		}
		hasilStatus, err := layananPembaruanStatus.PerbaruiStatusOtomatis(reqStatus, "webhook_gateway")
		if err != nil {
			log.Printf("[PPL-20] Peringatan perbarui status: %v", err)
		} else {
			log.Printf("[PPL-20] SUKSES: Status pesanan %s berubah menjadi '%s' (%s)", notif.OrderID, hasilStatus.StatusPesanan, hasilStatus.LabelStatus)
		}

		// 4. Penanganan Notifikasi Kedaluwarsa & Pemutusan Reservasi Stok (WBS 1.3.2.4)
		if notif.TransactionStatus == "expire" || notif.TransactionStatus == "cancel" {
			hasilLepas, errLepas := layananKedaluwarsaStok.PutusReservasiStokKedaluwarsa(notif.OrderID, "Notifikasi expire gateway sandbox PPL-20")
			if errLepas != nil {
				log.Printf("[PPL-20] Peringatan pemutusan reservasi stok: %v", errLepas)
			} else {
				log.Printf("[PPL-20] SUKSES: Reservasi stok pesanan %s diputus (%d bungkus dikembalikan ke inventaris)", notif.OrderID, hasilLepas.TotalKuantitasStok)
			}
		}

		// 5. Kunci Idempotency selesai diproses (WBS 1.3.2.2)
		layananKeamananIdem.TandaiSelesaiIdempotency(notif.OrderID, notif.TransactionStatus)

		// 6. Sinkronkan dengan notifikasi owner jika status settlement / dibayar
		if notif.TransactionStatus == "settlement" || notif.TransactionStatus == "capture" {
			_, _ = layananNotifikasi.CatatPemberitahuan(notifikasipesanan.PermintaanTriggerNotifikasi{
				NomorPesanan:      notif.OrderID,
				NomorInvoice:      "INV-" + notif.OrderID,
				TipePemberitahuan: "pesanan_terbayar",
				Judul:             "Pembayaran Berhasil Diterima",
				Pesan:             "Pesanan " + notif.OrderID + " telah dibayar lunas melalui payment gateway.",
				MetodePenerimaan:  "Kirim Kurir",
				TotalTransaksi:    56000,
			})
		}

		return nil
	}

	// Endpoint WBS 1.3.2.1 (Webhook Receiver & Audit Log)
	http.HandleFunc("/api/pembayaran/webhook/v2", layananWebhook.HandlerWebhook)
	http.HandleFunc("/api/pembayaran/webhook/riwayat", layananWebhook.HandlerRiwayatLog)

	// Endpoint WBS 1.3.2.2 (Signature Verification & Idempotency Check)
	http.HandleFunc("/api/pembayaran/verifikasi-keamanan", layananKeamananIdem.HandlerVerifikasiKeamanan)
	http.HandleFunc("/api/pembayaran/idempotency/riwayat", layananKeamananIdem.HandlerRiwayatIdempotency)

	// Endpoint WBS 1.3.2.3 (Auto Update Status & Realtime Status Query)
	http.HandleFunc("/api/pembayaran/perbarui-status", layananPembaruanStatus.HandlerPerbaruiStatus)
	http.HandleFunc("/api/pembayaran/status-pesanan", layananPembaruanStatus.HandlerStatusPesanan)

	// Endpoint WBS 1.3.2.4 (Stock Reservation Release on Expiry)
	http.HandleFunc("/api/pembayaran/kedaluwarsa/putus-reservasi", layananKedaluwarsaStok.HandlerPutusReservasi)
	http.HandleFunc("/api/pembayaran/kedaluwarsa/riwayat-pemutusan", layananKedaluwarsaStok.HandlerRiwayatPemutusan)

	// ==========================================
	// WBS 1.3.5.3 & 1.3.5.4 - PENJUALAN OFFLINE & FLAG PEMISAH (PPL-24)
	// ==========================================

	layananPenjualanOffline := penjualanoffline.InisialisasiLayanan(db)
	layananPemisahFlag := flagpemisah.InisialisasiLayanan(db)

	// POST / GET /api/admin/penjualan-offline (Kasir Toko Langsung WBS 1.3.5.3)
	http.HandleFunc(
		"/api/admin/penjualan-offline",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				layananPenjualanOffline.HandlerDaftarPenjualan(w, r)
			} else {
				layananPenjualanOffline.HandlerSimpanPenjualan(w, r)
			}
		},
	)

	// GET /api/admin/transaksi/audit-flag (Audit Flag Nyata vs Sandbox WBS 1.3.5.4)
	http.HandleFunc(
		"/api/admin/transaksi/audit-flag",
		layananPemisahFlag.HandlerAudit,
	)

	// POST /api/admin/transaksi/validasi-flag (Validasi Flag WBS 1.3.5.4)
	http.HandleFunc(
		"/api/admin/transaksi/validasi-flag",
		layananPemisahFlag.HandlerValidasiFlag,
	)

	// ==========================================
	// WBS 1.4.2.1 - LAPORAN REKAP STOK AKHIR (PPL-26)
	// ==========================================

	layananLaporanStok := laporanbulanan.InisialisasiLayanan(db)

	// GET /api/admin/laporan/stok-akhir (Perhitungan sisa stok akhir periode WBS 1.4.2.1)
	http.HandleFunc(
		"/api/admin/laporan/stok-akhir",
		layananLaporanStok.HandlerLaporanStokAkhir,
	)

	// ==========================================
	// WBS 1.4.4.1 - LAPORAN REKAP PENGELUARAN USAHA (PPL-28)
	// ==========================================

	layananLaporanPengeluaran := laporanpengeluaran.InisialisasiLayanan(db)

	// GET /api/admin/laporan/pengeluaran (Pengelompokan biaya pengeluaran WBS 1.4.4.1)
	http.HandleFunc(
		"/api/admin/laporan/pengeluaran",
		layananLaporanPengeluaran.HandlerLaporanPengeluaran,
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
