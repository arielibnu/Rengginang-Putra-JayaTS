-- ============================================================================
-- WBS 1.3.2.2 - SKEMA DATABASE IDEMPOTENSI WEBHOOK PEMBAYARAN GATEWAY
-- Aplikasi: E-Commerce Rengginang Putra Jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS idempotensi_webhook_pembayaran (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    idempotency_key VARCHAR(150) NOT NULL UNIQUE,
    nomor_pesanan VARCHAR(64) NOT NULL,
    transaction_id VARCHAR(100) NOT NULL,
    status_transaksi VARCHAR(50) NOT NULL,
    sudah_diproses TINYINT(1) NOT NULL DEFAULT 0,
    jumlah_percobaan INT NOT NULL DEFAULT 1,
    pertama_kali_diterima DATETIME DEFAULT CURRENT_TIMESTAMP,
    terakhir_diterima DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_idempotency_pesanan (nomor_pesanan),
    INDEX idx_idempotency_trx (transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
