
-- WBS 1.1.2.2
-- Database akun pelanggan

CREATE TABLE IF NOT EXISTS pelanggan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(150) NOT NULL,
    email VARCHAR(191) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    aktif BOOLEAN NOT NULL DEFAULT TRUE,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sesi_pelanggan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pelanggan_id BIGINT UNSIGNED NOT NULL,
    token_hash CHAR(64) NOT NULL UNIQUE,
    kedaluwarsa_pada DATETIME NOT NULL,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_sesi_pelanggan
        FOREIGN KEY (pelanggan_id)
        REFERENCES pelanggan(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_sesi_pelanggan_kedaluwarsa
ON sesi_pelanggan (kedaluwarsa_pada);
