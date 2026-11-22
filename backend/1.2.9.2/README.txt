================================================================================
DOKUMENTASI TEKNIS BACKEND WBS 1.2.9.2
Mekanisme Reservasi Stok Otomatis (Menahan Stok) dengan Batas Waktu 24 Jam
Aplikasi: E-Commerce Rengginang Putra Jaya
================================================================================

1. DESKRIPSI WORK PACKAGE
- WBS ID      : 1.2.9.2
- Level       : Level 4 (Work Package)
- Durasi      : 2 hari
- Peran       : BE (Backend Engineer)
- Deliverable : Stock Reservation & Auto-Release

2. SPESIFIKASI FITUR
a. Validasi Ketersediaan Sisa Stok Sebelum Reservasi (Proteksi Overbooking):
   - Sisa Stok = Total Stok - Total Stok Sedang Direservasi Aktif.
   - Jika kuantitas dipesan > sisa stok, transaksi otomatis ditolak (HTTP 409 Conflict).
   - Menjamin tidak ada pelanggan yang memesan melebihi stok gudang riil.
b. Penahanan Stok Otomatis (Holding Time 24 Jam):
   - Saat pesanan terbuat, status pesanan menjadi "menunggu_pembayaran".
   - Stok varian ditahan secara transaksional di tabel `reservasi_stok`.
   - Batas waktu pembayaran dihitung otomatis: waktu_pembuatan + 24 jam.
c. Auto-Release Stok Kadaluarsa:
   - Jika pembayaran tidak diselesaikan dalam 24 jam, status pesanan menjadi "dibatalkan"
     dan kuota reservasi stok otomatis dilepas kembali ke toko.
d. Endpoint HTTP:
   - POST /api/pesanan/buat: Membuat pesanan, snapshot harga, dan reservasi stok.
   - GET  /api/pesanan/detail: Mengambil detail pesanan asli untuk Halaman Konfirmasi PPL-18.
   - GET  /api/pesanan/reservasi: Audit list kuota stok yang sedang ditahan.
================================================================================
