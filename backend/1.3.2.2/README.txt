================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.2.2
PPL-20: Verifikasi Status Transaksi & Idempotency Webhook Payment Gateway
Tugas: Verifikasi Signature Key Keamanan Notifikasi Gateway dan Proteksi Pemrosesan Ganda (Idempotency)
Deliverable: SHA512 Signature Verifier & Idempotency Lock Layer
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan dua pilar keamanan kritis pada payment gateway:
   a. Verifikasi Keaslian Signature Key SHA512:
      - Formula: SHA512(order_id + status_code + gross_amount + server_key).
      - Menjamin notifikasi berasal secara sah dari Midtrans Sandbox dan data tidak dimanipulasi di tengah jalan.
      - Menggunakan teknik subtle.ConstantTimeCompare untuk mencegah Timing Attacks.
   b. Proteksi Pemrosesan Ganda (Idempotency Protection):
      - Payment gateway seperti Midtrans sering mengirimkan notifikasi retry jika server lambat merespons.
      - Modul ini menggunakan kunci idempotency unik (order_id + status_transaksi).
      - Jika notifikasi pembayaran settlement/lunas yang sama dikirimkan berulang kali, modul ini mendeteksi duplikasi, menghentikan eksekusi duplikat (BolehDiproses: false), dan mengembalikan status idempotent aman tanpa mencatat transaksi atau menambah saldo dua kali.

2. STRUKTUR FILE
   - model.go                   : Struct verifikasi signature key, catatan idempotency, dan request payload.
   - signature_idempotency.go   : Algoritma validasi SHA512 dan manajemen tabel/memori kunci idempotency.
   - signature_idempotency_test.go : Unit test signature matching, tampering detection, dan duplicate prevention.
   - schema_idempotency.sql     : Skema database MySQL idempotensi_webhook_pembayaran.
   - README.txt                 : Dokumentasi teknis WBS 1.3.2.2.

3. ENDPOINT API
   - POST /api/pembayaran/verifikasi-keamanan
     Menguji verifikasi signature key dan kelayakan idempotency notifikasi.
   - GET /api/pembayaran/idempotency/riwayat
     Mengambil log entri kunci idempotency dan status pemrosesan ganda.
================================================================================
