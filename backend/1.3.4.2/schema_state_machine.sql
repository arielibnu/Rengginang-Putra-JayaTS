-- Schema untuk WBS 1.3.4.2: Penerapan State Machine Perubahan Status Terpisah dari Status Pembayaran
ALTER TABLE pesanan 
ADD COLUMN IF NOT EXISTS status_pemrosesan VARCHAR(30) NOT NULL DEFAULT 'menunggu_diproses',
ADD COLUMN IF NOT EXISTS nomor_resi VARCHAR(100) NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS riwayat_status_pemrosesan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    status_lama VARCHAR(50) NOT NULL,
    status_baru VARCHAR(50) NOT NULL,
    nomor_resi VARCHAR(100) NULL DEFAULT '',
    diubah_oleh VARCHAR(50) NOT NULL DEFAULT 'Owner',
    catatan VARCHAR(255) NULL DEFAULT '',
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_riwayat_pemrosesan_pesanan (nomor_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
