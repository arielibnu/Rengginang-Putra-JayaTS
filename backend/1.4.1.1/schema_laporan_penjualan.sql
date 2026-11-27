-- Skema Referensi WBS 1.4.1.1 & 1.4.1.3
-- Perhitungan jumlah barang terjual dan omzet bulanan berdasarkan transaksi valid
-- serta pemisahan transaksi riil dari simulasi sandbox dan pengecualian pesanan batal.

USE rengginang_putra_jaya;

-- Tabel Penjualan Offline (Kas Masuk Nyata)
-- is_sandbox = 0 -> Kas nyata
-- flag_lingkungan = 'PRODUKSI_NYATA' -> Valid

-- Tabel Pesanan Checkout (Simulasi Online Sandbox)
-- status_pesanan = 'batal' | 'dibatalkan' | 'kedaluwarsa' -> Dikecualikan dari total barang terjual
