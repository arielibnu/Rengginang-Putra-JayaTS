================================================================================
DOKUMENTASI TEKNIS BACKEND WBS 1.2.9.1
Pembuatan ID Pesanan Unik serta Penyimpanan Snapshot Harga Produk, Ongkir, dan Total
Aplikasi: E-Commerce Rengginang Putra Jaya
================================================================================

1. DESKRIPSI WORK PACKAGE
- WBS ID      : 1.2.9.1
- Level       : Level 4 (Work Package)
- Durasi      : 2 hari
- Peran       : BE (Backend Engineer)
- Deliverable : Logika Database Order & Invoice ID

2. SPESIFIKASI FITUR
a. Pembuatan Nomor Pesanan Unik (Format: PJ-YYYYMMDD-XXXX):
   - Menggunakan tanggal transaksi riil (20060102) dan 4 digit kode numerik acak aman.
   - Contoh: PJ-20261122-7891.
b. Pembuatan Nomor Invoice Unik (Format: INV-YYYYMMDD-XXXX):
   - Contoh: INV-20261122-7891.
c. Penyimpanan Snapshot Harga (Price Immutability):
   - Menyimpan harga satuan produk, subtotal belanja, tarif ongkos kirim resmi, dan grand total.
   - Perubahan harga produk di masa mendatang tidak akan mengubah nilai pesanan yang sudah tercatat.
d. Tabel MySQL:
   - pesanan: Data induk transaksi dan snapshot finansial.
   - snapshot_item_pesanan: Rincian varian dan harga satuan pada saat order dibuat.
================================================================================
