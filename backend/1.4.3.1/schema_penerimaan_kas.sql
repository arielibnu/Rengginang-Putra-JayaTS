-- Skema Database Referensi WBS 1.4.3.1 (BE)
-- Perhitungan total dana riil yang masuk berdasarkan tanggal penyelesaian pembayaran (bukan saat order)

USE rengginang_putra_jaya;

-- Tabel sumber data:
-- 1. penjualan_offline (Kasir Toko Fisik) -> is_sandbox = 0, status lunas seketika
-- 2. pesanan_checkout (Online Digital) -> status_pesanan in ('dibayar', 'lunas', 'selesai')
-- 3. status_pending -> status_pesanan in ('menunggu_pembayaran', 'pending')
