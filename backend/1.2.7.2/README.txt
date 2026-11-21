================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.7.2
PPL-16: Integrasi API Pengiriman Sandbox dan Ongkos Kirim
Tugas: Integrasi Permintaan Layanan Kurir dan Tarif Pengiriman Melalui API Kurir Sandbox Berdasarkan Berat
Deliverable: Service Gateway API Ongkir Sandbox
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan Service Gateway API kurir sandbox (WBS 1.2.7.2):
   - Menerima parameter pengiriman dari modul alamat tujuan (PPL-15) dan berat paket
     berdasarkan item belanja (default 1.350 gram untuk 2 Udang + 1 Bawang).
   - Menghitung ongkos kirim secara akurat dengan formula pembulatan kilogram
     (Math.Ceil(berat_gram / 1000.0)).
   - Menghubungkan gateway ke 3 mitra kurir utama sesuai desain Figma Gambar 1:
     * JNE REG (default terpilih pada tampilan frontend)
     * J&T Express
     * SiCepat (Regular/BEST)
   - Menyimpan seluruh kalkulasi ke tabel database MySQL `riwayat_kalkulasi_ongkir`
     dan `pilihan_kurir_checkout`, serta mendukung penyimpanan in-memory jika MySQL offline.
   - Menyediakan API yang terhubung langsung ke antarmuka frontend React Checkout.

2. STRUKTUR FILE
   - model.go                : Definisi struct data PermintaanHitungOngkir, LayananKurir,
                               ResponsKalkulasiOngkir, LokasiAsal, LokasiTujuan, dan PermintaanPilihKurir.
   - gateway_ongkir.go       : Logika gateway kalkulasi ongkir sandbox berbasis berat,
                               master tarif kurir antarkota/provinsi, persistensi database,
                               dan HTTP Handlers:
                               * HandlerHitungOngkir    (POST /api/kurir/hitung-ongkir)
                               * HandlerDaftarEkspedisi (GET  /api/kurir/daftar-ekspedisi)
                               * HandlerPilihKurir      (POST /api/kurir/pilih-layanan)
                               * HandlerRiwayatKalkulasi(GET  /api/kurir/riwayat)
   - gateway_ongkir_test.go  : Unit test pengujian kalkulasi berat, multiplier tarif kurir,
                               default berat paket, endpoint HTTP, dan metode request.
   - schema_ongkir.sql       : DDL SQL MySQL tabel `riwayat_kalkulasi_ongkir` dan `pilihan_kurir_checkout`.
   - README.txt              : Dokumentasi teknis modul WBS 1.2.7.2.

3. ENDPOINT API
   - POST /api/kurir/hitung-ongkir
     Menerima payload tujuan dan berat paket, mengembalikan 3 opsi kurir beserta tarif dan estimasinya.
   - GET /api/kurir/daftar-ekspedisi
     Mendapatkan daftar ekspedisi yang didukung.
   - POST /api/kurir/pilih-layanan
     Menyimpan kurir yang dipilih user ke database saat menekan tombol "Gunakan layanan".
   - GET /api/kurir/riwayat
     Melihat audit log riwayat kalkulasi tarif.
