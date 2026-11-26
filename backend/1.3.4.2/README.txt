WBS 1.3.4.2 - Penerapan State Machine Perubahan Status Terpisah dari Status Pembayaran (PPL-23)
=================================================================================================

Modul ini bertanggung jawab mengelola alur siklus hidup status pemrosesan pesanan (State Machine) secara independen dari status pembayaran gateway.

Fitur Utama:
1. Validasi Transisi Status Formal:
   - baru / menunggu_diproses -> diproses, batal
   - diproses -> dikirim, batal
   - dikirim -> selesai
   - selesai / batal -> terminal (dilarang diubah lagi)
2. Pemisahan Tegas Status Pembayaran:
   - Status pemrosesan pesanan tidak boleh memanipulasi status pembayaran.
   - Status bayar (menunggu_pembayaran, dibayar, gagal, kadaluwarsa) dikontrol oleh webhook payment gateway.
3. Otomatisasi Pembatalan:
   - Saat status berpindah ke 'batal', secara otomatis memicu modul WBS 1.3.4.4 untuk mengembalikan stok.
4. Tabel audit riwayat_status_pemrosesan untuk pelacakan perubahan oleh Owner.
