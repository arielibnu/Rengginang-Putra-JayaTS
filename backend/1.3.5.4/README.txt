================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.5.4
PPL-24: Pencatatan Penjualan Langsung (Offline)
Tugas: Penerapan flag pemisah antara transaksi nyata (online/offline) dan transaksi simulasi sandbox
Deliverable: Validasi, klasifikasi flag transaksi, dan endpoint audit pemisahan kas
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini menerapkan flag pemisah tegas antara:
   - Transaksi Nyata (PRODUKSI_NYATA / is_sandbox: false): uang kas fisik penjualan langsung di toko
   - Transaksi Simulasi (SANDBOX_SIMULASI / is_sandbox: true): pembayaran dummy sandbox Midtrans
   Dengan pemisahan ini, pembukuan kas nyata tidak tercampur dengan transaksi percobaan/testing sandbox.

2. STRUKTUR FILE
   - model.go             : Definisi konstanta flag, struct rekap audit, dan payload validasi.
   - flag_pemisah.go      : Layanan evaluasi klasifikasi flag dan HTTP Handlers.
   - flag_pemisah_test.go : Unit testing evaluasi flag nyata vs sandbox.
   - README.txt           : Dokumentasi teknis deliverable WBS 1.3.5.4.

3. ENDPOINT API
   - GET /api/admin/transaksi/audit-flag
     Menghasilkan rekap agregasi pemisahan transaksi nyata vs sandbox.

   - POST /api/admin/transaksi/validasi-flag
     Memverifikasi flag sebelum transaksi dicatat ke database.
