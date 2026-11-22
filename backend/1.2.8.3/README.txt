================================================================================
DOKUMENTASI TEKNIS BACKEND WBS 1.2.8.3
Server-Side Order Validator (Validasi Ketat Integritas Subtotal, Ongkir, Total Nominal)
Aplikasi: E-Commerce Rengginang Putra Jaya
================================================================================

1. DESKRIPSI WORK PACKAGE
- WBS ID      : 1.2.8.3
- Level       : Level 4 (Work Package)
- Durasi      : 1 hari
- Peran       : BE (Backend Engineer)
- Deliverable : Server-Side Order Validator
- Deskripsi   : Validasi ketat integritas subtotal, ongkir, total nominal,
                dan kelengkapan data di sisi server sebelum pesanan dibuat.

2. FITUR UTAMA VALIDATOR ORDER
a. Pengecekan Kelengkapan Data:
   - Data kontak (Nama, Nomor WhatsApp valid Indonesia minimal 9-12 digit).
   - Data alamat tujuan penerima (Alamat lengkap, Kota/Kabupaten, Provinsi).
   - Metode penerimaan ("Kirim via Kurir").
   - Pilihan kurir dan layanan logistik.
   - Minimal 1 item pesanan.

b. Validasi Integritas Finansial Server-Side (Anti-Tampering):
   - Subtotal dihitung ulang dari harga resmi katalog database.
     Mencegah manipulasi harga/subtotal dari DevTools/Client-side inspect.
   - Validasi Ongkos Kirim:
     Ongkir diverifikasi dari API kurir resmi. Menolak estimasi palsu atau ongkir nol.
     Jika tarif ongkir belum terverifikasi oleh sistem, mengembalikan status INVALID
     sesuai desain Figma Gambar 2:
     "Belum dapat membuat pesanan: Tarif ongkir belum terverifikasi oleh sistem."
   - Integritas Total Akhir:
     Total Akhir = Subtotal Server + Ongkir Server.
     Jika total klien tidak cocok dengan total server, pesanan ditolak.

c. Validasi Ketersediaan Stok Produk:
   - Memeriksa kuota stok riil produk di database sebelum mengonfirmasi order.

d. Tanda Tangan Keamanan (HMAC-SHA256 Token):
   - Jika lolos validasi (Figma Gambar 1: Data valid), server menerbitkan
     TokenValidasi bertanda tangan digital.
   - Tombol "Buat Pesanan" mengirimkan token ini ke endpoint checkout server.

3. STRUKTUR ENDPOINT HTTP
- POST /api/checkout/validasi-ringkasan
  Header  : Content-Type: application/json
  Body    : PermintaanValidasiOrder (JSON)
  Respons : Status 200 OK dengan state Gambar 1 (Valid) atau Gambar 2 (Invalid).

- POST /api/checkout/buat-pesanan
  Header  : Content-Type: application/json
  Body    : PermintaanBuatPesanan (JSON)
  Respons : Status 201 Created jika berhasil, atau 422 Unprocessable Entity jika ditolak.

- GET  /api/checkout/audit-validasi
  Respons : Daftar audit log validasi pesanan dari database / memori server.

4. TABEL DATABASE MYSQL
- validasi_order_checkout : Log audit upaya validasi, integritas selisih, dan status penolakan.
- pesanan_checkout        : Transaksi pesanan yang berhasil divalidasi dan siap dibayar.
- pesanan_item_checkout   : Varian produk yang terkait dengan setiap pesanan.
================================================================================
