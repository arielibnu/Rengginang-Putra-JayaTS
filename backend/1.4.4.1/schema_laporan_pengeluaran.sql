-- Skema Referensi WBS 1.4.4.1 (BE)
-- Pengelompokan biaya pengeluaran operasional berdasarkan tanggal, kategori beban, dan periode

USE rengginang_putra_jaya;

CREATE TABLE IF NOT EXISTS pengeluaran_usaha (
    id_pengeluaran BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tanggal_pengeluaran DATE NOT NULL,
    kategori VARCHAR(100) NOT NULL,
    nominal DECIMAL(12,2) NOT NULL,
    keterangan VARCHAR(255) NULL,
    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_pengeluaran_tanggal (tanggal_pengeluaran),
    INDEX idx_pengeluaran_kategori (kategori)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Contoh data inisialisasi riil
INSERT INTO pengeluaran_usaha (tanggal_pengeluaran, kategori, nominal, keterangan)
SELECT '2026-11-05', 'Listrik & Gas Operasional', 450000.00, 'Tagihan listrik & tabung gas ruang penggorengan'
WHERE NOT EXISTS (SELECT 1 FROM pengeluaran_usaha WHERE tanggal_pengeluaran = '2026-11-05' AND kategori = 'Listrik & Gas Operasional');

INSERT INTO pengeluaran_usaha (tanggal_pengeluaran, kategori, nominal, keterangan)
SELECT '2026-11-12', 'Belanja Bahan Baku Ketan', 1250000.00, 'Pembelian 10 karung beras ketan super dari suplier'
WHERE NOT EXISTS (SELECT 1 FROM pengeluaran_usaha WHERE tanggal_pengeluaran = '2026-11-12' AND kategori = 'Belanja Bahan Baku Ketan');

INSERT INTO pengeluaran_usaha (tanggal_pengeluaran, kategori, nominal, keterangan)
SELECT '2026-11-20', 'Biaya Pemeliharaan Mesin & Lainnya', 350000.00, 'Perbaikan seal mesin pengemas standing pouch'
WHERE NOT EXISTS (SELECT 1 FROM pengeluaran_usaha WHERE tanggal_pengeluaran = '2026-11-20' AND kategori = 'Biaya Pemeliharaan Mesin & Lainnya');
