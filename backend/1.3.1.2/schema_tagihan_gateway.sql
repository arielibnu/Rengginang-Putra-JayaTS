-- ============================================================================
-- SKEMA TABEL TRANSAKSI PEMBAYARAN GATEWAY (MIDTRANS SANDBOX)
-- WBS 1.3.1.2 (BE)
-- ============================================================================

CREATE TABLE IF NOT EXISTS transaksi_pembayaran_gateway (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    id_pesanan BIGINT NULL,
    nomor_pesanan VARCHAR(64) NOT NULL,
    nomor_invoice VARCHAR(64) NULL,
    transaction_id VARCHAR(100) NOT NULL,
    payment_type VARCHAR(50) NOT NULL,
    bank VARCHAR(30) NULL,
    va_number VARCHAR(100) NULL,
    qr_url TEXT NULL,
    qr_string TEXT NULL,
    gross_amount BIGINT NOT NULL,
    transaction_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    fraud_status VARCHAR(50) NULL DEFAULT 'accept',
    raw_gateway_response TEXT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nomor_pesanan (nomor_pesanan),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_transaction_status (transaction_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
