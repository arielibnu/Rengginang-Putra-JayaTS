================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.2.4
PPL-20: Verifikasi Status Transaksi & Idempotency Webhook Payment Gateway
Tugas: Penanganan Notifikasi Kedaluwarsa dan Pemutusan Reservasi Stok Agar Tidak Terblokir
Deliverable: Auto-Release Engine & Stock Unlocking Service
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan pelepasan otomatis stok yang sedang direservasi (holding stock)
   ketika transaksi pembayaran payment gateway dinyatakan 'expire' (kedaluwarsa) atau 'cancel':
   - Mendeteksi kuota varian produk yang sedang ditahan pada tabel `reservasi_stok` untuk pesanan terkait.
   - Mengubah status penahanan dari 'menahan_stok' menjadi 'dilepas_otomatis'.
   - Mengembalikan kuantitas stok kembali ke inventaris aktif sehingga pembeli lain dapat memesan produk tersebut (menghindari deadlock / stok terblokir).
   - Memperbarui status pesanan menjadi 'kedaluwarsa' / 'dibatalkan'.
   - Melindungi pesanan yang sudah berstatus 'dibayar' agar tidak dapat dibatalkan atau dilepas stoknya secara keliru.
   - Mencatat seluruh riwayat pemutusan ke tabel `log_pemutusan_reservasi_stok` untuk transparansi audit pergudangan.

2. STRUKTUR FILE
   - model.go                : Struct item dilepas, payload permintaan putus reservasi, dan respon hasil.
   - kedaluwarsa_stok.go     : Layanan pemutusan reservasi stok, pengembalian stok, dan HTTP handler.
   - kedaluwarsa_stok_test.go: Unit testing pelepasan stok holding, pemulihan kuantitas, dan proteksi status.
   - schema_kedaluwarsa.sql  : Skema database MySQL log_pemutusan_reservasi_stok.
   - README.txt              : Dokumentasi teknis WBS 1.3.2.4.

3. ENDPOINT API
   - POST /api/pembayaran/kedaluwarsa/putus-reservasi
     Payload: {"nomor_pesanan": "PJ-...", "alasan": "..."}
     Respon: HTTP 200 OK detail item dan total kuantitas stok yang berhasil diputus dan dikembalikan.
   - GET /api/pembayaran/kedaluwarsa/riwayat-pemutusan
     Respon: HTTP 200 OK riwayat jejak audit pemutusan reservasi stok.
================================================================================
