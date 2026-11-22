-- ============================================================================
-- WBS 1.2.9.2 - SKEMA DATABASE RESERVASI STOK OTOMATIS & AUTO-RELEASE
-- Database: rengginang_putra_jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS reservasi_stok (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    id_varian BIGINT NOT NULL DEFAULT 1,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah_direservasi INT NOT NULL,
    status_reservasi VARCHAR(30) NOT NULL DEFAULT 'menahan_stok', -- 'menahan_stok', 'dikonfirmasi', 'dilepas_otomatis'
    kadaluarsa_pada DATETIME NOT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_nomor_reservasi (nomor_pesanan),
    INDEX idx_status_reservasi (status_reservasi),
    INDEX idx_kadaluarsa (kadaluarsa_pada)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
