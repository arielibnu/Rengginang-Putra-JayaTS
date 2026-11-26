-- ============================================================================
-- WBS 1.3.2.1 - SKEMA DATABASE LOG WEBHOOK PEMBAYARAN GATEWAY
-- Aplikasi: E-Commerce Rengginang Putra Jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS log_webhook_pembayaran (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    transaction_id VARCHAR(100) NOT NULL,
    status_transaksi VARCHAR(50) NOT NULL,
    status_code VARCHAR(10) NOT NULL,
    gross_amount VARCHAR(30) NOT NULL,
    payment_type VARCHAR(50) NULL,
    signature_key VARCHAR(150) NULL,
    raw_payload TEXT NOT NULL,
    status_respon_http INT NOT NULL DEFAULT 200,
    pesan_respon VARCHAR(255) NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_webhook_pesanan (nomor_pesanan),
    INDEX idx_webhook_trx (transaction_id),
    INDEX idx_webhook_status (status_transaksi)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
