-- ============================================================================
-- SKEMA TABEL DATABASE UNTUK WBS 1.2.6.3
-- PPL-15: Pengisian Wilayah dan Alamat Pengiriman
-- Deliverable: Validasi Data Geografis Pengiriman
-- ============================================================================

CREATE TABLE IF NOT EXISTS `alamat_pengiriman_pesanan` (
    `id_alamat` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `id_pesanan` VARCHAR(64) NULL,
    `nama_penerima` VARCHAR(150) NOT NULL,
    `nomor_whatsapp` VARCHAR(25) NOT NULL,
    `metode_penerimaan` ENUM('kurir', 'ambil_sendiri') NOT NULL DEFAULT 'kurir',
    `provinsi` VARCHAR(100) NOT NULL DEFAULT '',
    `kota_kabupaten` VARCHAR(100) NOT NULL DEFAULT '',
    `kecamatan` VARCHAR(100) NOT NULL DEFAULT '',
    `kode_pos` VARCHAR(10) NOT NULL DEFAULT '',
    `alamat_lengkap` TEXT NOT NULL,
    `catatan_alamat` TEXT NULL,
    `kurir_ready` BOOLEAN NOT NULL DEFAULT TRUE,
    `waktu_dibuat` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_pesanan` (`id_pesanan`),
    INDEX `idx_wilayah` (`provinsi`, `kota_kabupaten`, `kecamatan`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel master data wilayah Indonesia terintegrasi dengan standarisasi kode kurir
CREATE TABLE IF NOT EXISTS `master_wilayah_kurir` (
    `id_wilayah` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `kode_provinsi` VARCHAR(10) NOT NULL,
    `nama_provinsi` VARCHAR(100) NOT NULL,
    `kode_kota` VARCHAR(10) NOT NULL,
    `nama_kota` VARCHAR(100) NOT NULL,
    `kode_kecamatan` VARCHAR(10) NOT NULL,
    `nama_kecamatan` VARCHAR(100) NOT NULL,
    `kode_pos_default` VARCHAR(5) NOT NULL,
    `aktif` BOOLEAN NOT NULL DEFAULT TRUE,
    INDEX `idx_hierarki` (`nama_provinsi`, `nama_kota`, `nama_kecamatan`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
