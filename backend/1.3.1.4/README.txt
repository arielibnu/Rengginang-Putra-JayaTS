================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.1.4
PPL-19: Integrasi Payment Gateway Sandbox
Tugas: Sinkronisasi Identitas Transaksi Gateway (transaction_id) dengan Nomor Pesanan di Database
Deliverable: Sinkronisasi Transaksi Midtrans, Webhook Notifikasi, dan Penautan Database
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan sinkronisasi transaksi payment gateway Midtrans:
   - Penautan identitas transaction_id Midtrans dengan nomor_pesanan di database MySQL.
   - Pengecekan status transaksi real-time ke API Midtrans Sandbox (/v2/{order_id}/status).
   - Pemrosesan Webhook HTTP Notification Callback dari Midtrans dengan verifikasi SHA512 signature.
   - Pembaruan otomatis status pesanan di database MySQL menjadi 'dibayar' saat settlement/capture.
   - Endpoint simulasi pembayaran instan untuk pengujian demo.

2. STRUKTUR FILE
   - model.go                  : Definisi struct notifikasi webhook, respons sinkronisasi, dan payload simulasi.
   - sinkronisasi_gateway.go   : Logika sinkronisasi status real-time, webhook handler, dan persistensi database.
   - sinkronisasi_gateway_test.go : Unit test sinkronisasi status, pemrosesan webhook settlement, dan simulasi.
   - schema_sinkronisasi.sql   : Indeks database penautan transaksi pesanan.
   - README.txt                : Dokumentasi teknis modul WBS 1.3.1.4.

3. ENDPOINT API
   - GET /api/pembayaran/status?nomor_pesanan=PJ-...
     Mengambil status transaksi terkini dari Midtrans dan database.
   - POST /api/pembayaran/webhook
     Menerima callback notifikasi pembayaran resmi dari server Midtrans Sandbox.
   - POST /api/pembayaran/simulasi-bayar
     Mensimulasikan pembayaran sukses secara instan untuk verifikasi integrasi.
