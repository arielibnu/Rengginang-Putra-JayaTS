-- ============================================================================
-- WBS 1.3.2.3 - SKEMA DATABASE RIWAYAT PEMBARUAN STATUS TRANSAKSI
-- Aplikasi: E-Commerce Rengginang Putra Jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS riwayat_status_transaksi (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    transaction_id VARCHAR(100) NOT NULL,
    status_lama VARCHAR(50) NOT NULL,
    status_baru VARCHAR(50) NOT NULL,
    label_status VARCHAR(50) NOT NULL,
    pemicu VARCHAR(50) NOT NULL DEFAULT 'webhook',
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_riwayat_pesanan (nomor_pesanan),
    INDEX idx_riwayat_trx (transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
