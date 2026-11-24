================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.5.3
PPL-24: Pencatatan Penjualan Langsung (Offline)
Tugas: Pencatatan penjualan offline yang langsung memotong stok fisik tanpa melalui alur ekspedisi online
Deliverable: API dan service transaksi kasir toko langsung
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan pencatatan penjualan offline langsung oleh Owner (kasir toko).
   Fitur utama:
   - Pencatatan transaksi beli langsung di tempat (POS).
   - Pemotongan stok fisik riil di database MySQL secara atomik.
   - Tanpa melalui alur ekspedisi online dan tanpa membuat order duplikat di tabel pesanan.
   - Menghasilkan bukti nota kasir dengan nomor transaksi POS-YYYYMMDD-XXXX.

2. STRUKTUR FILE
   - model.go                   : Struct data transaksi kasir offline, permintaan, dan balasan.
   - penjualan_offline.go       : Layanan transaksi kasir, validasi stok, dan HTTP Handlers.
   - schema_penjualan_offline.sql: DDL tabel penjualan_offline.
   - penjualan_offline_test.go  : Unit testing fungsionalitas dan endpoint HTTP.
   - README.txt                 : Dokumentasi teknis deliverable WBS 1.3.5.3.

3. ENDPOINT API
   - POST /api/admin/penjualan-offline
     Menyimpan transaksi kasir langsung dan memotong stok fisik varian.
     Payload:
     {
       "tanggal_penjualan": "20 November 2026",
       "id_varian": 1,
       "nama_varian": "Rengginang Udang",
       "jumlah": 2,
       "harga_satuan": 28000,
       "metode_pembayaran": "Tunai"
     }

   - GET /api/admin/penjualan-offline
     Mengambil daftar riwayat penjualan kasir toko.
