-- ====================================================================
-- PPL-05: PERINGATAN STOK MENIPIS OTOMATIS
-- WBS 1.1.5.2: Penyediaan Pengaturan Batas Minimum Stok Per Varian
--              dan Penerapan Aturan Peringatan Stok Kritis
-- ====================================================================

USE rengginang_putra_jaya;

-- 1. Tambahkan kolom batas_minimum pada tabel varian_produk (jika belum ada)
-- Default batas minimum adalah 15 bungkus sesuai standar operasional Owner
ALTER TABLE varian_produk 
ADD COLUMN IF NOT EXISTS batas_minimum INT UNSIGNED NOT NULL DEFAULT 15 AFTER stok;

-- 2. Tabel konfigurasi batas minimum stok per varian (opsional/relasional eksplisit)
CREATE TABLE IF NOT EXISTS batas_minimum_stok (
    id_varian BIGINT UNSIGNED PRIMARY KEY,
    batas_minimum INT UNSIGNED NOT NULL DEFAULT 15,
    catatan VARCHAR(255) NULL,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL 
        DEFAULT CURRENT_TIMESTAMP 
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_batas_min_varian
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT chk_batas_minimum_nonnegatif
        CHECK (batas_minimum >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Inisialisasi batas minimum default (15 bungkus) untuk varian yang sudah ada
INSERT INTO batas_minimum_stok (id_varian, batas_minimum)
SELECT id_varian, 15
FROM varian_produk
ON DUPLICATE KEY UPDATE batas_minimum = VALUES(batas_minimum);
