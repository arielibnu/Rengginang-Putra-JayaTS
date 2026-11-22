-- ============================================================================
-- WBS 1.2.9.1 - SKEMA DATABASE ORDER GENERATOR & SNAPSHOT HARGA
-- Database: rengginang_putra_jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS pesanan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
    nomor_invoice VARCHAR(64) NOT NULL UNIQUE,
    tipe_checkout VARCHAR(20) NOT NULL DEFAULT 'guest',
    id_pelanggan BIGINT NULL,
    nama_penerima VARCHAR(100) NOT NULL,
    nomor_whatsapp VARCHAR(30) NOT NULL,
    alamat_lengkap TEXT NOT NULL,
    kota_kabupaten VARCHAR(100) NOT NULL,
    provinsi VARCHAR(100) NOT NULL,
    kurir_id VARCHAR(50) NOT NULL,
    kurir_nama VARCHAR(100) NOT NULL,
    kurir_layanan VARCHAR(50) NOT NULL,
    snapshot_subtotal INT NOT NULL,
    snapshot_ongkir INT NOT NULL,
    snapshot_total INT NOT NULL,
    status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu_pembayaran',
    status_reservasi VARCHAR(30) NOT NULL DEFAULT 'direservasi',
    batas_waktu_reservasi DATETIME NOT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nomor_pesanan (nomor_pesanan),
    INDEX idx_nomor_invoice (nomor_invoice),
    INDEX idx_status_pesanan (status_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS snapshot_item_pesanan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    id_pesanan BIGINT NOT NULL,
    nomor_pesanan VARCHAR(64) NOT NULL,
    id_varian BIGINT NOT NULL DEFAULT 1,
    id_produk VARCHAR(50) NOT NULL,
    nama_produk VARCHAR(150) NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah INT NOT NULL,
    harga_satuan_snapshot INT NOT NULL,
    subtotal_snapshot INT NOT NULL,
    berat_satuan_gram INT NOT NULL DEFAULT 500,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_pesanan_snapshot (nomor_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
