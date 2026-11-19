================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.2.2
PPL-11: Penyediaan Penyimpanan Keranjang Bagi Guest (localStorage) dan Pelanggan Login (Sinkronisasi)
Deliverable: State Manager Keranjang Belanja
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan logika backend untuk WBS 1.2.2.2:
   - State Manager untuk sinkronisasi data keranjang belanja.
   - Penanganan pembeli tamu (Guest): Data keranjang dikelola pada localStorage frontend,
     dan divalidasi keamanannya oleh backend saat sinkronisasi request.
   - Penanganan pelanggan terdaftar (Logged-in): Data keranjang disimpan persisten ke basis data
     MySQL pada tabel `keranjang_pelanggan` dan `item_keranjang_pelanggan`.
   - Perhitungan otomatis kuantitas, subtotal per item, dan total belanja di sisi server.
   - Fitur penghapusan item keranjang dan pengosongan keranjang belanja.

2. STRUKTUR FILE
   - model.go                  : Definisi struct ItemKeranjang, PermintaanSinkronisasi,
                                 ResponKeranjang, dan ResponsAPI standar.
   - state_keranjang.go        : Logika sinkronisasi state, rekap kuantitas dan harga,
                                 koneksi database MySQL, serta HTTP Handlers:
                                 * HandlerSinkronisasiKeranjang (POST /api/keranjang/sinkronisasi)
                                 * HandlerAmbilKeranjang        (GET  /api/keranjang)
                                 * HandlerHapusItemKeranjang    (DELETE /api/keranjang/item)
                                 * HandlerBersihkanKeranjang    (POST /api/keranjang/bersihkan)
   - state_keranjang_test.go   : Unit test komprehensif menguji skenario guest, pelanggan login,
                                 kondisi PPL11_01 & PPL11_02, hapus item, dan HTTP methods.
   - schema_state_keranjang.sql: DDL tabel `keranjang_pelanggan` dan `item_keranjang_pelanggan`.
   - README.txt                : Dokumentasi teknis lengkap WBS 1.2.2.2.

3. ENDPOINT API
   a. POST /api/keranjang/sinkronisasi
      - Deskripsi: Menerima payload keranjang dari frontend dan menyinkronkannya.
      - Body JSON:
        {
          "is_guest": true,
          "items": [
            { "id_varian": 1, "nama_varian": "Udang", "jumlah": 2, "harga_satuan": 28000 },
            { "id_varian": 2, "nama_varian": "Bawang", "jumlah": 1, "harga_satuan": 28000 }
          ]
        }
      - Respon HTTP 200 OK:
        {
          "status": "sukses",
          "data": {
            "id_keranjang": "GUEST-CART",
            "total_kuantitas": 3,
            "subtotal": 84000,
            "format_subtotal": "Rp84.000",
            ...
          }
        }

   b. GET /api/keranjang
      - Deskripsi: Mengambil data keranjang aktif pelanggan atau fallback resmi.

   c. DELETE /api/keranjang/item?id_varian={id}
      - Deskripsi: Menghapus item tertentu dari keranjang belanja.

4. CARA MENJALANKAN & MENGUJI
   cd backend
   go test -v ./1.2.2.2
