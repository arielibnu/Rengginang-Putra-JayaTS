-- ============================================================================
-- RENGGINANG PUTRA JAYA
-- WBS 1.2.3.3 - Penyamaan Format Payload Data Checkout Bagi Pembeli Guest & Pelanggan
-- Skema Database MySQL untuk Tabel Pesanan Checkout
-- ============================================================================

USE rengginang_putra_jaya;

-- 1. Tabel Utama Pesanan Checkout (Menampung transaksi dari Guest & Pelanggan Login)
CREATE TABLE IF NOT EXISTS pesanan_checkout (
    id_pesanan BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
    tipe_checkout ENUM('guest', 'pelanggan') NOT NULL DEFAULT 'guest',
    is_guest BOOLEAN NOT NULL DEFAULT TRUE,
    id_pelanggan BIGINT UNSIGNED NULL,
    nama_penerima VARCHAR(150) NOT NULL,
    nomor_whatsapp VARCHAR(30) NOT NULL,
    email VARCHAR(191) NULL,
    catatan TEXT NULL,
    total_item INT UNSIGNED NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    ongkos_kirim DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_pembayaran DECIMAL(12,2) NOT NULL,
    status_pesanan VARCHAR(50) NOT NULL DEFAULT 'MENUNGGU_KONFIRMASI',
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pesanan_pelanggan
        FOREIGN KEY (id_pelanggan)
        REFERENCES pelanggan(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Tabel Rincian Item Pesanan
CREATE TABLE IF NOT EXISTS detail_pesanan_checkout (
    id_detail BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_pesanan BIGINT UNSIGNED NOT NULL,
    id_varian BIGINT UNSIGNED NOT NULL,
    nama_produk_snapshot VARCHAR(150) NOT NULL,
    nama_varian_snapshot VARCHAR(100) NOT NULL,
    harga_satuan_snapshot DECIMAL(12,2) NOT NULL,
    jumlah INT UNSIGNED NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,

    CONSTRAINT fk_detail_pesanan
        FOREIGN KEY (id_pesanan)
        REFERENCES pesanan_checkout(id_pesanan)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_detail_varian_checkout
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indeks untuk pencarian cepat nomor pesanan dan riwayat pelanggan
CREATE INDEX IF NOT EXISTS idx_pesanan_nomor ON pesanan_checkout (nomor_pesanan);
CREATE INDEX IF NOT EXISTS idx_pesanan_pelanggan ON pesanan_checkout (id_pelanggan);
