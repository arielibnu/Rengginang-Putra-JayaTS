================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.3.2
PPL-22: Notifikasi Pesanan Baru untuk Owner
Tugas: Pembuatan trigger catatan pemberitahuan saat pesanan baru sukses terbayar atau masuk
Deliverable: Trigger notifikasi pesanan dan API notifikasi Owner Dashboard
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan pencatatan pemberitahuan otomatis (trigger notification)
   saat pesanan baru masuk dari pelanggan (PPL-18 / Checkout) maupun saat pesanan
   berhasil dibayar (PPL-19 Midtrans payment gateway).
   Data notifikasi disimpan ke tabel MySQL `notifikasi_pesanan` dan disediakan API
   realtime polling untuk komponen Notifikasi Pesanan Baru pada Dashboard Owner Figma.

2. STRUKTUR FILE
   - model.go                  : Definisi struct data notifikasi pesanan, ringkasan, dan payload respons.
   - notifikasi_pesanan.go     : Layanan manajemen notifikasi, trigger pesanan masuk/terbayar, dan HTTP Handler.
   - notifikasi_pesanan_test.go: Unit testing otomatis untuk trigger dan HTTP Handlers.
   - schema_notifikasi.sql     : DDL skema database MySQL dan trigger otomatis.
   - README.txt                : Dokumentasi modul WBS 1.3.3.2.

3. SKEMA DATABASE
   - Tabel: `notifikasi_pesanan`
     Kolom:
     * id BIGINT AUTO_INCREMENT PRIMARY KEY
     * nomor_pesanan VARCHAR(64)
     * nomor_invoice VARCHAR(64)
     * tipe_pemberitahuan VARCHAR(50) ('pesanan_masuk', 'pesanan_terbayar')
     * judul VARCHAR(100)
     * pesan TEXT
     * metode_penerimaan VARCHAR(50) ('Ambil Sendiri', 'Kirim Kurir')
     * total_transaksi BIGINT
     * status_pesanan VARCHAR(30) ('baru', 'diproses', 'selesai')
     * sudah_dibaca TINYINT(1)
     * dibuat_pada DATETIME

4. ENDPOINT API
   - GET /api/owner/notifikasi/pesanan-baru
     Respons:
     {
       "status": "sukses",
       "jumlah_baru": 2,
       "badge_label": "2 Baru",
       "daftar_pesanan": [
         {"nomor_pesanan": "PJ-DEMO-0001", "metode": "Ambil Sendiri", "teks_ringkas": "PJ-DEMO-0001 — Ambil Sendiri"},
         {"nomor_pesanan": "PJ-DEMO-0002", "metode": "Kirim Kurir", "teks_ringkas": "PJ-DEMO-0002 — Kirim Kurir"}
       ]
     }

   - GET /api/owner/notifikasi
     Mengambil seluruh riwayat catatan pemberitahuan.

   - POST /api/owner/notifikasi/baca
     Menandai notifikasi telah dibaca oleh Owner.

   - POST /api/owner/notifikasi/trigger
     Memicu pembuatan catatan notifikasi pesanan secara programatis/event-driven.

5. INTEGRASI
   - Terhubung dengan PPL-18 (Pembuatan Pesanan /api/pesanan/buat)
   - Terhubung dengan PPL-19 (Midtrans Payment Gateway /api/pembayaran/simulasi-bayar & webhook)
   - Terhubung dengan PPL-22 Frontend Owner Dashboard untuk realtime polling (interval 8 detik)
