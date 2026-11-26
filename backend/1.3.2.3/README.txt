================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.3.2.3
PPL-20: Verifikasi Status Transaksi & Idempotency Webhook Payment Gateway
Tugas: Pembaruan Otomatis Status Transaksi (Menunggu, Lunas/Dibayar, Kedaluwarsa, Gagal)
Deliverable: State Machine & Auto Status Updater Midtrans ke Database
Role: BE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan pemetaan dan pembaruan otomatis status transaksi pesanan
   secara real-time berdasarkan notifikasi status dari payment gateway Midtrans Sandbox:
   - settlement / capture (accept) -> Status Pesanan: "dibayar" | Label: "Lunas/Dibayar"
   - pending                        -> Status Pesanan: "menunggu_pembayaran" | Label: "Menunggu"
   - expire                         -> Status Pesanan: "kedaluwarsa" | Label: "Kedaluwarsa"
   - cancel / deny / failure        -> Status Pesanan: "gagal" | Label: "Gagal"
   
   Pembaruan dilakukan secara sinkron ke tabel `pesanan`, tabel `transaksi_pembayaran_gateway`,
   serta mencatat audit pergantian status ke tabel `riwayat_status_transaksi`.

2. STRUKTUR FILE
   - model.go                   : Definisi konstanta status Midtrans, status pesanan, struct hasil & audit.
   - pembaruan_status.go        : Logika pemetaan status, update ke database dan in-memory cache, HTTP handlers.
   - pembaruan_status_test.go   : Unit testing keakuratan pemetaan status dan eksekusi pembaruan otomatis.
   - schema_pembaruan_status.sql: Skema database MySQL riwayat_status_transaksi.
   - README.txt                 : Dokumentasi teknis WBS 1.3.2.3.

3. ENDPOINT API
   - POST /api/pembayaran/perbarui-status
     Payload: {"nomor_pesanan": "PJ-...", "transaction_status": "settlement", "transaction_id": "trx-..."}
     Respon: HTTP 200 OK status berhasil diperbarui dengan detail label resmi.
   - GET /api/pembayaran/status-pesanan?nomor_pesanan=PJ-...
     Respon: HTTP 200 OK data status terkini dari pesanan.
================================================================================
