================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.2.4
PPL-11: Validasi Ulang Isi Item Keranjang Terhadap Ketersediaan Stok Sebelum Checkout
Deliverable: API Validasi Keranjang
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan logika backend untuk WBS 1.2.2.4 (BE):
   - Validasi ulang seluruh item di dalam keranjang belanja terhadap sisa stok aktual di MySQL
     sebelum pembeli diperbolehkan melangkah ke proses checkout.
   - Pengecekan status ketersediaan stok:
     * VALID              : Sisa stok mencukupi untuk kuantitas yang diminta.
     * MELEBIHI_STOK      : Kuantitas melebihi sisa persediaan (ditolak HTTP 400).
     * STOK_HABIS         : Produk kosong / stok = 0 (ditolak HTTP 400).
     * NONAKTIF           : Varian produk dinonaktifkan (ditolak HTTP 400).
     * TIDAK_DITEMUKAN    : ID / slug varian tidak terdaftar (ditolak HTTP 400).
   - Menghitung ulang secara terverifikasi total kuantitas dan grand total subtotal belanja.
   - Memberikan feedback `bisa_checkout: true/false` untuk keamanan transaksi.

2. STRUKTUR FILE
   - model.go                     : Definisi struct PermintaanValidasiKeranjang, ItemValidasiRequest,
                                    ItemValidasiResult, HasilValidasiKeranjang, ResponsAPI.
   - validasi_keranjang.go        : Logika verifikasi database MySQL tabel `varian_produk`,
                                    kalkulasi subtotal terverifikasi, dan HTTP Handlers:
                                    * HandlerValidasiKeranjang (POST /api/keranjang/validasi)
                                    * HandlerValidasiKeranjang (POST /api/keranjang/validasi-stok)
   - validasi_keranjang_test.go   : Unit test komprehensif menguji kondisi normal, PPL11_02,
                                    melebihi stok, stok habis, keranjang kosong, dan HTTP methods.
   - schema_validasi_keranjang.sql: Skema query MySQL untuk validasi sisa stok.
   - README.txt                   : Dokumentasi teknis lengkap WBS 1.2.2.4.

3. ENDPOINT API
   POST /api/keranjang/validasi
   Content-Type: application/json

   Contoh Request Body:
   {
     "items": [
       {"id_varian": 1, "nama_varian": "Udang", "jumlah": 2},
       {"id_varian": 2, "nama_varian": "Bawang", "jumlah": 1}
     ]
   }

   Contoh Respon Sukses (HTTP 200 OK):
   {
     "status": "sukses",
     "pesan": "Semua isi item di keranjang tersedia dan kuota mencukupi untuk checkout.",
     "data": {
       "bisa_checkout": true,
       "total_jenis_produk": 2,
       "total_kuantitas": 3,
       "subtotal": 84000,
       "format_subtotal": "Rp84.000",
       "status_validasi": "VALIDASI_SUKSES",
       ...
     }
   }

4. CARA MENJALANKAN UNIT TEST
   cd backend
   go test -v ./1.2.2.4
