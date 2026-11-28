-- PPL-29 | WBS 1.4.5.1
-- Skrip query dan optimasi indeks untuk laporan estimasi keuntungan bulanan
USE rengginang_putra_jaya;

-- Indeks pencarian periode penjualan offline
CREATE INDEX IF NOT EXISTS idx_penjualan_off_periode 
ON penjualan_offline (is_sandbox, flag_lingkungan, dibuat_pada);

-- Indeks pencarian beban pengeluaran operasional
CREATE INDEX IF NOT EXISTS idx_pengeluaran_usaha_tgl 
ON pengeluaran_usaha (tanggal_pengeluaran);
