================================================================================
WBS 1.2.10.2 - ENDPOINT API DAFTAR PESANAN & DETAIL BARANG OWNER (PPL-21)
================================================================================

Deskripsi:
Modul backend untuk mendukung Dashboard Pemantauan Order Masuk Owner:
1. Endpoint GET /api/admin/pesanan:
   - Pencarian berdasarkan nomor pesanan (PJ-YYYYMMDD-XXXX), nama penerima, atau WhatsApp.
   - Filter status pemrosesan (Baru, Diproses, Selesai, Batal) dan pembayaran (Menunggu, Dibayar).
   - Filter rentang waktu transaksi (Hari ini, 7 hari terakhir, 30 hari terakhir, Semua tanggal).
2. Endpoint GET /api/admin/pesanan/detail:
   - Mengambil rincian pesanan pelanggan (nama, nomor WA, alamat lengkap / pickup, kurir).
   - Mengambil snapshot barang pesanan dari tabel snapshot_item_pesanan.
   - Pemisahan data status pembayaran (dibayar / menunggu) dan status pemrosesan (diproses / selesai / batal).
3. Endpoint PUT / POST /api/admin/pesanan/status:
   - Pembaruan status pemrosesan pesanan oleh Owner secara transaksional ke database MySQL.

Struktur Berkas:
- model.go                : Definisi struct filter, ringkasan, detail, dan respons JSON.
- manajemen_pesanan.go    : Query database MySQL, filter logik, handler HTTP API.
- manajemen_pesanan_test.go: Pengujian unit test fitur filter, detail barang, dan HTTP handler (100% PASS).
- schema_pesanan_owner.sql: Skema penambahan kolom status_pemrosesan & indeks performa.
- README.txt              : Dokumentasi teknis modul backend.
================================================================================
