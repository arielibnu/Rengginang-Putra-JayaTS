
-- PPL-03 | WBS 1.1.3.4
-- Snapshot harga transaksi

USE rengginang_putra_jaya;

CREATE TABLE IF NOT EXISTS transaksi (
    id_transaksi BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS detail_transaksi (
    id_detail BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_transaksi BIGINT UNSIGNED NOT NULL,
    id_varian BIGINT UNSIGNED NOT NULL,

    nama_varian_snapshot VARCHAR(100) NOT NULL,
    harga_jual_snapshot DECIMAL(12,2) NOT NULL,
    jumlah INT UNSIGNED NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,

    CONSTRAINT fk_detail_transaksi
        FOREIGN KEY (id_transaksi)
        REFERENCES transaksi(id_transaksi)
        ON DELETE RESTRICT,

    CONSTRAINT fk_detail_varian
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE RESTRICT,

    CONSTRAINT chk_jumlah_transaksi
        CHECK (jumlah > 0),

    CONSTRAINT chk_harga_snapshot
        CHECK (harga_jual_snapshot >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
