================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.1.2
PPL-19: Integrasi Payment Gateway Sandbox
Tugas: Integrasi Pembuatan Tagihan Payment Gateway Sandbox (Midtrans) Berbasis ID Pesanan
Deliverable: Pembuatan Tagihan QRIS & Virtual Account Sandbox
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan pembuatan tagihan pembayaran online menggunakan
   Midtrans Sandbox Core API:
   - Pembuatan tagihan QRIS dinamis (GoPay / QRIS)
   - Pembuatan nomor Virtual Account (BCA, BNI, BRI, Permata, Mandiri)
   - Penautan ID pesanan (order_id) dengan nilai transaksi tagihan
   - Penyimpanan data tagihan ke database MySQL (tabel transaksi_pembayaran_gateway)
   - Fallback in-memory map tangguh jika database sedang offline

2. STRUKTUR FILE
   - model.go                  : Definisi struct data tagihan, request charge Midtrans, dan respons.
   - tagihan_gateway.go        : Logika pembuatan tagihan HTTP ke Midtrans API Sandbox dan Handler HTTP.
   - tagihan_gateway_test.go   : Unit testing pembuatan tagihan QRIS dan Virtual Account.
   - schema_tagihan_gateway.sql: Skema database MySQL transaksi pembayaran.
   - README.txt                : Dokumentasi modul WBS 1.3.1.2.

3. ENDPOINT API
   - POST /api/pembayaran/buat-tagihan
     Payload: {"nomor_pesanan": "PJ-...", "metode": "qris|va", "bank": "bca|bni|bri"}
     Respons: ResponsTagihanGateway dengan QR URL, VA Number, transaction_id, dan simulator URL.
