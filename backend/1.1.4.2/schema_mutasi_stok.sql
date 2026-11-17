
USE rengginang_putra_jaya;

-- Tabel stok simulasi
CREATE TABLE IF NOT EXISTS stok_simulasi (
    id_varian BIGINT UNSIGNED PRIMARY KEY,
    stok_simulasi INT UNSIGNED NOT NULL DEFAULT 0,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_simulasi_varian
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel riwayat mutasi stok
CREATE TABLE IF NOT EXISTS riwayat_mutasi_stok (
    id_mutasi BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    id_varian BIGINT UNSIGNED NOT NULL,

    jenis_mutasi ENUM(
        'MASUK',
        'KELUAR',
        'PENYESUAIAN'
    ) NOT NULL,

    mode_stok ENUM(
        'RIIL',
        'SIMULASI'
    ) NOT NULL DEFAULT 'RIIL',

    jumlah INT UNSIGNED NOT NULL,
    stok_sebelum INT UNSIGNED NOT NULL,
    stok_sesudah INT UNSIGNED NOT NULL,

    catatan VARCHAR(255) NULL,

    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_mutasi_varian (id_varian),
    INDEX idx_mutasi_tanggal (dibuat_pada),
    INDEX idx_mutasi_mode (mode_stok),

    CONSTRAINT fk_mutasi_varian
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT chk_jumlah_mutasi
        CHECK (jumlah > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
