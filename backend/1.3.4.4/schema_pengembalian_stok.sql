-- Schema untuk WBS 1.3.4.4: Mekanisme Pengembalian Stok Otomatis saat Pesanan Dibatalkan
CREATE TABLE IF NOT EXISTS log_pengembalian_stok_batal (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    id_varian BIGINT NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah_dikembalikan INT NOT NULL,
    alasan VARCHAR(255) NOT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_log_kembali_pesanan (nomor_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
