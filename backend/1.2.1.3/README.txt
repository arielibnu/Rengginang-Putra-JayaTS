================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.1.3
PPL-10: Validasi Batas Jumlah Pembelian Terhadap Sisa Stok Tersedia Secara Realtime
Deliverable: Validasi Kuota Stok Realtime
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan logika backend untuk WBS 1.2.1.3 (BE):
   - Validasi realtime sisa stok aktual dari database MySQL pada tabel `varian_produk`.
   - Pencegahan pemesanan melebihi sisa stok yang tersedia (kuota terlampaui).
   - Validasi batas jumlah pembelian minimal (minimal 1 bungkus, tidak boleh 0 atau negatif).
   - Perhitungan total harga secara realtime berdasarkan kuantitas yang dipesan.
   - Pengecekan status ketersediaan (TERSEDIA, STOK HABIS, MELEBIHI STOK, JUMLAH TIDAK VALID).
   - Sinkronisasi endpoint validasi dengan antarmuka Frontend React pada komponen stepper kuantitas.

2. STRUKTUR FILE
   - model.go                : Definisi struct data PermintaanValidasiKuota, DataVarianStok,
                               HasilValidasiKuota, dan ResponsAPI standar.
   - validasi_kuota.go       : Logika bisnis pengecekan kuota, query stok realtime ke MySQL,
                               fallback data resmi Figma, dan HTTP Handlers:
                               * HandlerCekKuotaRealtime (GET /api/katalog/cek-kuota)
                               * HandlerCekKuotaRealtime (POST /api/katalog/validasi-kuota)
                               * HandlerCekKuotaRealtime (POST /api/katalog/validasi-beli)
   - validasi_kuota_test.go  : Unit test komprehensif menguji skenario batas stok, pembelian normal,
                               pembelian melebihi kuota, stok habis, jumlah 0/negatif, dan HTTP method.
   - schema_kuota_stok.sql   : DDL dan seed data MySQL untuk tabel produk dan varian_produk.
   - README.txt              : Dokumentasi teknis WBS 1.2.1.3.

3. ENDPOINT API

   a. GET /api/katalog/cek-kuota?id={id_atau_slug}&jumlah={jumlah_diminta}
      - Deskripsi : Mengecek kuota stok secara realtime saat pengguna mengubah kuantitas pada antarmuka.
      - Parameter Query:
        * id / slug : id atau nama varian (contoh: udang, cumi, bawang, terasi, atau 1, 2, 3, 4).
        * jumlah    : kuantitas pembelian yang diminta.
      
      - Contoh Respon Sukses (Udang, Jumlah: 2 bungkus, Stok: 12 bungkus):
        HTTP 200 OK
        {
          "status": "sukses",
          "pesan": "Jumlah pembelian valid (2 bungkus) dan kuota stok tersedia.",
          "data": {
            "id_varian": 1,
            "slug": "udang",
            "nama_varian": "Udang",
            "harga_satuan": 28000,
            "format_harga": "Rp28.000",
            "stok_tersedia": 12,
            "format_stok": "12 bungkus",
            "batas_maksimal_beli": 12,
            "jumlah_diminta": 2,
            "total_harga": 56000,
            "format_total_harga": "Rp56.000",
            "sisa_setelah_beli": 10,
            "status_ketersediaan": "TERSEDIA",
            "kuota_tersedia": true,
            "bisa_dibeli": true,
            "valid": true,
            "pesan": "Jumlah pembelian valid (2 bungkus) dan kuota stok tersedia.",
            "keterangan": "Stok mencukupi untuk pesanan Anda. Sisa stok setelah pembelian: 10 bungkus."
          }
        }

      - Contoh Respon Ditolak - Melebihi Kuota Stok (Udang, Jumlah: 15 bungkus, Stok: 12 bungkus):
        HTTP 400 Bad Request
        {
          "status": "ditolak",
          "pesan": "Jumlah pembelian (15 bungkus) melebihi batas kuota sisa stok yang tersedia (12 bungkus).",
          "data": {
            "id_varian": 1,
            "slug": "udang",
            "nama_varian": "Udang",
            "harga_satuan": 28000,
            "format_harga": "Rp28.000",
            "stok_tersedia": 12,
            "format_stok": "12 bungkus",
            "batas_maksimal_beli": 12,
            "jumlah_diminta": 15,
            "total_harga": 420000,
            "format_total_harga": "Rp420.000",
            "sisa_setelah_beli": -3,
            "status_ketersediaan": "MELEBIHI STOK",
            "kuota_tersedia": false,
            "bisa_dibeli": false,
            "valid": false,
            "pesan": "Jumlah pembelian (15 bungkus) melebihi batas kuota sisa stok yang tersedia (12 bungkus).",
            "keterangan": "Maksimal pembelian untuk varian ini adalah 12 bungkus. Kurangi jumlah pesanan Anda."
          }
        }

      - Contoh Respon Ditolak - Stok Habis (Cumi, Jumlah: 1 bungkus, Stok: 0 bungkus):
        HTTP 400 Bad Request
        {
          "status": "ditolak",
          "pesan": "Produk sedang habis (out-of-stock). Pembelian dinonaktifkan.",
          "data": {
            "id_varian": 3,
            "slug": "cumi",
            "nama_varian": "Cumi",
            "stok_tersedia": 0,
            "status_ketersediaan": "STOK HABIS",
            "kuota_tersedia": false,
            "bisa_dibeli": false,
            "valid": false,
            ...
          }
        }

   b. POST /api/katalog/validasi-kuota dan POST /api/katalog/validasi-beli
      - Deskripsi : Endpoint POST yang menerima payload JSON dari form pembelian saat checkout/klik Beli.
      - Request Body:
        {
          "id_varian": "udang",
          "jumlah": 5
        }
      - Respon HTTP:
        * 200 OK jika jumlah <= sisa stok.
        * 400 Bad Request jika jumlah > sisa stok atau stok habis.

4. CARA MENJALANKAN & MENGUJI
   a. Menjalankan Unit Test:
      cd backend
      go test -v ./1.2.1.3

   b. Menjalankan Server Backend:
      cd backend
      go run ./1.1.1.2
