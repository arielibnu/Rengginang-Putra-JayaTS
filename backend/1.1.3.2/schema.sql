
-- PPL-03 | WBS 1.1.3.2
-- Skema Database Produk dan Varian

USE rengginang_putra_jaya;

-- Tabel utama produk
CREATE TABLE IF NOT EXISTS produk (
    id_produk BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama_produk VARCHAR(150) NOT NULL,
    deskripsi TEXT,
    status_aktif BOOLEAN NOT NULL DEFAULT TRUE,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel varian produk
CREATE TABLE IF NOT EXISTS varian_produk (
    id_varian BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_produk BIGINT UNSIGNED NOT NULL,
    nama_varian VARCHAR(100) NOT NULL,
    harga_jual DECIMAL(12,2) NOT NULL,
    berat_gram INT UNSIGNED NOT NULL,
    stok INT UNSIGNED NOT NULL DEFAULT 0,
    status_aktif BOOLEAN NOT NULL DEFAULT TRUE,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_varian_produk
        FOREIGN KEY (id_produk)
        REFERENCES produk(id_produk)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT uq_varian_per_produk
        UNIQUE (id_produk, nama_varian),

    CONSTRAINT chk_harga_nonnegatif
        CHECK (harga_jual >= 0),

    CONSTRAINT chk_berat_positif
        CHECK (berat_gram > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
