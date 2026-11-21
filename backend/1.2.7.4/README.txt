================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.7.4
PPL-16: Integrasi API Pengiriman Sandbox dan Ongkos Kirim
Tugas: Penanganan Kegagalan API Kurir atau Rute Tidak Terjangkau Tanpa Membuat Estimasi Tarif Palsu
Deliverable: Error Handling Koneksi API Kurir
Role: BE, FE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan sistem penanganan kegagalan koneksi API kurir
   dan penanganan rute tidak terjangkau (WBS 1.2.7.4):
   - Mencegah keras pembuatan estimasi tarif palsu / tebakan harga (EstimasiPalsu = false)
     ketika API kurir offline, down, timeout, atau wilayah tujuan belum terjangkau.
   - Mengembalikan pesan terstruktur resmi yang 100% cocok dengan mockup Figma Gambar 2:
     * Judul: "Tarif pengiriman tidak tersedia"
     * Detail: "API gagal atau layanan tidak melayani wilayah tujuan."
     * Peringatan: "Jangan membuat tarif ongkir perkiraan sendiri."
     * Aksi Tindakan: ["Coba lagi", "Ubah alamat"]
   - Menghasilkan HTTP Status Code yang tepat:
     * 503 Service Unavailable (saat API kurir down / timeout)
     * 422 Unprocessable Entity (saat rute wilayah tidak terjangkau)
   - Menyediakan fitur simulasi kegagalan untuk pengetesan otomatis dan audit log kegagalan.

2. STRUKTUR FILE
   - model.go                 : Definisi struct data ResponsErrorKurir, TipeKegagalanKurir,
                                PermintaanSimulasiKegagalan, dan LogErrorAPI.
   - error_handling_kurir.go  : Logika penanganan error, validasi keterjangkauan rute,
                                persistensi audit error ke database MySQL/in-memory, dan HTTP Handlers:
                                * HandlerSimulasiError (POST /api/kurir/simulasi-error)
                                * HandlerCekJangkauan  (POST /api/kurir/cek-jangkauan)
                                * HandlerLogKegagalan  (GET  /api/kurir/log-error)
   - error_handling_kurir_test.go : Unit test pengujian struktur error Gambar 2, aturan tanpa tarif palsu,
                                keterjangkauan wilayah, dan simulasi endpoint HTTP.
   - README.txt               : Dokumentasi teknis modul WBS 1.2.7.4.

3. ENDPOINT API
   - POST /api/kurir/simulasi-error
     Endpoint simulasi kegagalan API kurir mengembalikan HTTP 503 dan pesan Gambar 2.
   - POST /api/kurir/cek-jangkauan
     Memeriksa apakah wilayah tujuan terjangkau oleh kurir mitra pengiriman.
   - GET /api/kurir/log-error
     Melihat catatan riwayat error API kurir.
