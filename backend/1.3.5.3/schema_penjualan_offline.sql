-- ============================================================================
-- SKEMA DATABASE WBS 1.3.5.3 & 1.3.5.4: PENJUALAN OFFLINE & FLAG PEMISAH
-- Database: rengginang_putra_jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS penjualan_offline (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_transaksi VARCHAR(64) NOT NULL UNIQUE,
    tanggal_penjualan VARCHAR(50) NOT NULL,
    id_varian BIGINT NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah INT NOT NULL,
    harga_satuan DECIMAL(12,2) NOT NULL,
    total_harga DECIMAL(15,2) NOT NULL,
    metode_pembayaran VARCHAR(50) NOT NULL DEFAULT 'Tunai',
    sumber_transaksi VARCHAR(50) NOT NULL DEFAULT 'offline_kasir',
    flag_lingkungan VARCHAR(50) NOT NULL DEFAULT 'PRODUKSI_NYATA',
    is_sandbox TINYINT(1) NOT NULL DEFAULT 0,
    stok_sebelum INT NOT NULL DEFAULT 0,
    stok_sesudah INT NOT NULL DEFAULT 0,
    kasir_nama VARCHAR(100) NOT NULL DEFAULT 'Owner Toko',
    catatan TEXT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_nomor (nomor_transaksi),
    INDEX idx_varian (id_varian),
    INDEX idx_flag (flag_lingkungan),
    INDEX idx_sumber (sumber_transaksi),
    INDEX idx_dibuat (dibuat_pada)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
