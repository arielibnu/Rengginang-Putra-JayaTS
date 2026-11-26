WBS 1.3.4.4 - Mekanisme Pengembalian Stok Otomatis Saat Pesanan Dibatalkan (PPL-23)
===================================================================================

Modul ini bertanggung jawab menangani pelepasan dan pengembalian stok varian produk secara otomatis saat pesanan pelanggan dibatalkan atau kedaluwarsa.

Fitur Utama:
1. Lookup item transaksi pesanan dari tabel snapshot_item_pesanan atau reservasi_stok.
2. Penambahan stok varian secara atomik pada tabel varian_produk.
3. Pencatatan mutasi resmi ke tabel riwayat_mutasi_stok dengan mode 'MASUK' dan catatan pembatalan.
4. Idempotency pengembalian: mencegah stok dikembalikan dua kali untuk pesanan yang sama.
5. Tabel audit log_pengembalian_stok_batal untuk pelacakan transaksi owner.
