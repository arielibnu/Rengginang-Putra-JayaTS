-- ================================================================================
-- SKEMA TABEL DATABASE: WBS 1.2.7.2
-- Kalkulasi Ongkos Kirim dan Pemilihan Layanan Kurir Sandbox
-- ================================================================================

CREATE TABLE IF NOT EXISTS riwayat_kalkulasi_ongkir (
    id INT AUTO_INCREMENT PRIMARY KEY,
    asal_kota VARCHAR(100) NOT NULL DEFAULT 'Kabupaten Situbondo',
    tujuan_provinsi VARCHAR(100) NOT NULL,
    tujuan_kota VARCHAR(100) NOT NULL,
    tujuan_kecamatan VARCHAR(100) NOT NULL,
    tujuan_kodepos VARCHAR(10) NOT NULL,
    berat_gram INT NOT NULL,
    opsi_layanan_json TEXT NOT NULL,
    dibuat_pada DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tujuan_kota (tujuan_kota),
    INDEX idx_dibuat_pada (dibuat_pada)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS pilihan_kurir_checkout (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_kurir VARCHAR(50) NOT NULL,
    nama_kurir VARCHAR(100) NOT NULL,
    tarif INT NOT NULL,
    estimasi_hari VARCHAR(50) NOT NULL,
    alamat_tujuan TEXT NOT NULL,
    nama_penerima VARCHAR(150) NOT NULL,
    dibuat_pada DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_kurir (id_kurir),
    INDEX idx_dibuat_pada (dibuat_pada)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
