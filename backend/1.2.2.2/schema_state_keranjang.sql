-- ============================================================================
-- SKEMA BASIS DATA WBS 1.2.2.2 (PPL-11)
-- Tabel Sinkronisasi Keranjang Belanja Pelanggan
-- ============================================================================

CREATE TABLE IF NOT EXISTS keranjang_pelanggan (
    id_keranjang BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pelanggan_id BIGINT UNSIGNED NOT NULL UNIQUE,
    total_item INT UNSIGNED NOT NULL DEFAULT 0,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_keranjang_pelanggan FOREIGN KEY (pelanggan_id) REFERENCES pelanggan(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS item_keranjang_pelanggan (
    id_item BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    keranjang_id BIGINT UNSIGNED NOT NULL,
    id_varian BIGINT UNSIGNED NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    jumlah INT UNSIGNED NOT NULL DEFAULT 1,
    harga_satuan DECIMAL(12,2) NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_item_keranjang FOREIGN KEY (keranjang_id) REFERENCES keranjang_pelanggan(id_keranjang) ON DELETE CASCADE,
    CONSTRAINT fk_item_varian FOREIGN KEY (id_varian) REFERENCES varian_produk(id_varian) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed data simulasi keranjang pelanggan awal sesuai Figma PPL11_01
INSERT INTO keranjang_pelanggan (id_keranjang, pelanggan_id, total_item, subtotal)
VALUES (1, 1, 3, 84000)
ON DUPLICATE KEY UPDATE total_item = VALUES(total_item), subtotal = VALUES(subtotal);

INSERT INTO item_keranjang_pelanggan (keranjang_id, id_varian, nama_varian, jumlah, harga_satuan, subtotal)
VALUES 
    (1, 1, 'Udang', 2, 28000, 56000),
    (1, 2, 'Bawang', 1, 28000, 28000)
ON DUPLICATE KEY UPDATE jumlah = VALUES(jumlah), subtotal = VALUES(subtotal);
