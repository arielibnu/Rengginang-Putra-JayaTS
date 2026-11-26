-- Schema Query Pendukung WBS 1.4.2.1: Perhitungan Sisa Stok Akhir Berdasarkan Riwayat Mutasi

-- Indeks optimasi penutupan periode stok
ALTER TABLE riwayat_mutasi_stok 
ADD INDEX IF NOT EXISTS idx_mutasi_periode (id_varian, mode_stok, jenis_mutasi, dibuat_pada);

-- Indeks reservasi aktif
ALTER TABLE reservasi_stok 
ADD INDEX IF NOT EXISTS idx_reservasi_status_tanggal (id_varian, status_reservasi, dibuat_pada);
