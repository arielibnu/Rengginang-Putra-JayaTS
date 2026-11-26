================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.2.1
PPL-20: Verifikasi Status Transaksi & Idempotency Webhook Payment Gateway
Tugas: Implementasi Endpoint Webhook untuk Menerima Notifikasi Status Pembayaran dari Payment Gateway
Deliverable: Endpoint Webhook HTTP Handler, Audit Logging, & Respon Standar Midtrans
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan endpoint penerima webhook resmi dari payment gateway Midtrans Sandbox:
   - Menerima payload JSON notifikasi HTTP POST dari server Midtrans.
   - Melakukan parsing dan validasi kelengkapan data notifikasi pembayaran (order_id, transaction_status, status_code, gross_amount, signature_key).
   - Menyimpan seluruh jejak riwayat payload asli notifikasi ke tabel MySQL `log_webhook_pembayaran` dan memory cache untuk audit trail.
   - Mengembalikan HTTP 200 OK dengan format respons standar agar payment gateway tidak melakukan retry bertubi-tubi.
   - Menolak metode selain POST dengan status HTTP 405 Method Not Allowed.

2. STRUKTUR FILE
   - model.go             : Definisi struct notifikasi webhook Midtrans, log entri, dan format respon standar.
   - webhook_gateway.go   : Implementasi handler HTTP webhook POST & endpoint riwayat audit log.
   - webhook_gateway_test.go : Unit test validasi notifikasi sukses, bad JSON payload, dan proteksi HTTP method.
   - schema_webhook.sql   : Skema DDL tabel MySQL log_webhook_pembayaran.
   - README.txt           : Dokumentasi teknis WBS 1.3.2.1.

3. ENDPOINT API
   - POST /api/pembayaran/webhook
     Payload: JSON Midtrans notification (order_id, transaction_status, gross_amount, signature_key, dll)
     Respon: HTTP 200 OK {"status": "sukses", "pesan": "...", "nomor_pesanan": "...", "waktu_proses": "..."}
   - GET /api/pembayaran/webhook/riwayat
     Respon: HTTP 200 OK daftar riwayat log notifikasi webhook yang pernah diterima.
================================================================================
