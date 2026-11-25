-- ============================================================================
-- WBS 1.3.2.4 - SKEMA DATABASE LOG PEMUTUSAN RESERVASI STOK KEDALUWARSA
-- Aplikasi: E-Commerce Rengginang Putra Jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS log_pemutusan_reservasi_stok (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    id_varian BIGINT NOT NULL DEFAULT 1,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah_dilepas INT NOT NULL,
    status_sebelumnya VARCHAR(30) NOT NULL DEFAULT 'menahan_stok',
    status_akhir VARCHAR(30) NOT NULL DEFAULT 'dilepas_otomatis',
    alasan VARCHAR(255) NOT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_log_reservasi_pesanan (nomor_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
