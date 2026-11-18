================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.1.9.3
PPL-09: Penandaan Status Stok Habis & Pembatasan Interaksi Pembelian Out-of-Stock
================================================================================

1. DESKRIPSI
   Modul ini mengimplementasikan logika backend untuk WBS 1.1.9.3:
   - Validasi ketersediaan stok produk (TERSEDIA vs STOK HABIS).
   - Penandaan status dan pesan ketersediaan produk sesuai desain Figma PPL-09.
   - Pembatasan interaksi pembelian (validasi backend untuk tombol beli yang disabled
     pada frontend agar tidak bisa dibypass).

2. STRUKTUR FILE
   - model.go                : Definisi struct respons detail produk, galeri foto,
                               request dan response validasi beli.
   - ketersediaan.go         : Logika penentuan status ketersediaan, validasi batas beli,
                               dan HTTP Handlers:
                               * HandlerDetailProduk (GET /api/katalog/detail)
                               * HandlerValidasiBeli (POST /api/katalog/validasi-beli)
   - ketersediaan_test.go    : Pengujian unit test untuk validasi logika ketersediaan
                               dan endpoint HTTP.
   - README.txt              : Dokumentasi modul WBS 1.1.9.3.

3. ENDPOINT API

   a. GET /api/katalog/detail?id={id_varian}
      - Deskripsi : Mengambil rincian spesifikasi produk, harga, berat, stok,
                    dan status ketersediaan.
      - Parameter :
        * id : id varian produk (contoh: cumi, udang, bawang, terasi).
      - Contoh Respon (Stok Habis / Cumi):
        {
          "status": "sukses",
          "pesan": "Berhasil mengambil detail spesifikasi produk",
          "data": {
            "id_varian": "cumi",
            "nama": "Rengginang Cumi",
            "nama_varian": "Cumi",
            "kategori": "Rengginang tradisional",
            "deskripsi": "Rengginang tradisional • kemasan 450 gram",
            "berat_gram": 450,
            "berat_bersih_teks": "450 gram",
            "harga": 28000,
            "harga_teks": "Rp28.000",
            "stok": 0,
            "stok_tersedia_teks": "0 bungkus",
            "status_ketersediaan": "STOK HABIS",
            "bisa_dibeli": false,
            "pesan_status": "Produk sedang habis",
            "sub_pesan": "Pembelian dinonaktifkan. Silakan pilih varian lain.",
            "galeri_foto": [...]
          }
        }

      - Contoh Respon (Tersedia / Udang):
        {
          "status": "sukses",
          "pesan": "Berhasil mengambil detail spesifikasi produk",
          "data": {
            "id_varian": "udang",
            "nama": "Rengginang Udang",
            "nama_varian": "Udang",
            "stok": 12,
            "stok_tersedia_teks": "12 bungkus",
            "status_ketersediaan": "TERSEDIA",
            "bisa_dibeli": true,
            "pesan_status": "Stok tersedia",
            "sub_pesan": "Jumlah stok mengikuti data persediaan dari Owner.",
            ...
          }
        }

   b. POST /api/katalog/validasi-beli
      - Deskripsi : Memvalidasi permintaan pembelian apakah diperbolehkan
                    atau ditolak karena stok habis.
      - Request Body:
        {
          "id_varian": "cumi",
          "jumlah": 1
        }
      - Respon (Jika Stok Habis / HTTP 400):
        {
          "status": "gagal",
          "pesan": "Pembelian dinonaktifkan karena produk sedang habis (out-of-stock)",
          "data": {
            "bisa_dibeli": false,
            "alasan": "stok habis"
          }
        }

4. PENGUJIAN UNIT
   Jalankan:
   cd backend
   go test -v ./1.1.9.3
