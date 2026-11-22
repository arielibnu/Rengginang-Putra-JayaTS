-- ============================================================================
-- WBS 1.2.8.3 - SKEMA DATABASE SERVER-SIDE ORDER VALIDATOR
-- Tabel validasi integritas transaksi, pesanan terverifikasi, dan rincian item
-- Database: rengginang_putra_jaya
-- ============================================================================

-- 1. TABEL AUDIT VALIDASI ORDER (Menyimpan setiap upaya verifikasi integritas nominal & kelengkapan)
CREATE TABLE IF NOT EXISTS validasi_order_checkout (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_referensi VARCHAR(64) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'valid' (Gambar 1) atau 'invalid' (Gambar 2)
    nama_pelanggan VARCHAR(100) NOT NULL,
    nomor_whatsapp VARCHAR(30) NOT NULL,
    subtotal_klien INT NOT NULL,
    subtotal_server INT NOT NULL,
    ongkir_klien INT NOT NULL,
    ongkir_server INT NOT NULL,
    total_klien INT NOT NULL,
    total_server INT NOT NULL,
    selisih_nominal INT DEFAULT 0,
    kode_penolakan VARCHAR(50) DEFAULT 'DATA_VALID',
    pesan_keterangan TEXT,
    payload_order JSON,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. TABEL PESANAN CHECKOUT (Menyimpan pesanan yang berhasil lolos validasi server)
CREATE TABLE IF NOT EXISTS pesanan_checkout (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
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
    subtotal INT NOT NULL,
    ongkos_kirim INT NOT NULL,
    total_pembayaran INT NOT NULL,
    status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu_pembayaran',
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_nomor_pesanan (nomor_pesanan),
    INDEX idx_nomor_wa (nomor_whatsapp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. TABEL ITEM PESANAN CHECKOUT (Menyimpan rincian varian produk pesanan)
CREATE TABLE IF NOT EXISTS pesanan_item_checkout (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    id_pesanan BIGINT NOT NULL,
    nomor_pesanan VARCHAR(64) NOT NULL,
    id_produk VARCHAR(50) NOT NULL,
    nama_produk VARCHAR(150) NOT NULL,
    jumlah INT NOT NULL,
    harga_satuan INT NOT NULL,
    subtotal INT NOT NULL,
    INDEX idx_id_pesanan (id_pesanan),
    INDEX idx_nomor_pesanan_item (nomor_pesanan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
