-- ============================================================================
-- RENGGINANG PUTRA JAYA
-- WBS 1.2.1.3 - Validasi Batas Jumlah Pembelian Terhadap Sisa Stok Realtime
-- Skema Database MySQL dan Data Awal untuk Pengujian Realtime Stock & Quota
-- ============================================================================

CREATE DATABASE IF NOT EXISTS rengginang_putra_jaya;
USE rengginang_putra_jaya;

-- 1. Tabel Utama Produk (jika belum ada)
CREATE TABLE IF NOT EXISTS produk (
    id_produk BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama_produk VARCHAR(150) NOT NULL,
    deskripsi TEXT,
    status_aktif BOOLEAN NOT NULL DEFAULT TRUE,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Tabel Varian Produk dengan Kolom Stok Realtime
CREATE TABLE IF NOT EXISTS varian_produk (
    id_varian BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_produk BIGINT UNSIGNED NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    harga_jual DECIMAL(12,2) NOT NULL,
    berat_gram INT UNSIGNED NOT NULL,
    stok INT UNSIGNED NOT NULL DEFAULT 0,
    status_aktif BOOLEAN NOT NULL DEFAULT TRUE,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_kuota_varian_produk
        FOREIGN KEY (id_produk)
        REFERENCES produk(id_produk)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT uq_kuota_varian_per_produk
        UNIQUE (id_produk, nama_varian),

    CONSTRAINT chk_kuota_harga_nonnegatif
        CHECK (harga_jual >= 0),

    CONSTRAINT chk_kuota_berat_positif
        CHECK (berat_gram > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Data Awal Produk (Induk Rengginang)
INSERT INTO produk (id_produk, nama_produk, deskripsi, status_aktif)
VALUES (1, 'Rengginang Tradisional Putra Jaya', 'Rengginang tradisional khas gurih dan renyah kemasan 450 gram', TRUE)
ON DUPLICATE KEY UPDATE 
    nama_produk = VALUES(nama_produk),
    deskripsi = VALUES(deskripsi);

-- 4. Data Awal 4 Varian Produk Sesuai Mockup Resmi Figma (PPL-08, PPL-09, PPL-10)
-- - Udang  : stok 12 bungkus (TERSEDIA)
-- - Bawang : stok 38 bungkus (TERSEDIA)
-- - Cumi   : stok 0 bungkus  (STOK HABIS)
-- - Terasi : stok 38 bungkus (TERSEDIA)
INSERT INTO varian_produk (id_varian, id_produk, nama_varian, harga_jual, berat_gram, stok, status_aktif)
VALUES
    (1, 1, 'Udang', 28000.00, 450, 12, TRUE),
    (2, 1, 'Bawang', 28000.00, 450, 38, TRUE),
    (3, 1, 'Cumi', 28000.00, 450, 0, TRUE),
    (4, 1, 'Terasi', 28000.00, 450, 38, TRUE)
ON DUPLICATE KEY UPDATE
    harga_jual = VALUES(harga_jual),
    berat_gram = VALUES(berat_gram),
    stok = VALUES(stok),
    status_aktif = VALUES(status_aktif);
