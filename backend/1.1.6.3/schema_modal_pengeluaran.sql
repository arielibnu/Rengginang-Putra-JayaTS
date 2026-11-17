
-- PPL-06 | WBS 1.1.6.3
-- Tabel Modal per Varian dan Pengeluaran Operasional Usaha
-- Menyimpan riwayat modal dan pengeluaran sebagai sumber
-- perhitungan laporan bulanan.

USE rengginang_putra_jaya;

-- ==========================================
-- Tabel modal per varian rengginang
-- ==========================================
-- Menyimpan harga modal (HPP) setiap varian produk.
-- Setiap kali owner mengubah modal, baris baru ditambahkan
-- supaya riwayat perubahan harga tetap tersimpan.

CREATE TABLE IF NOT EXISTS modal_varian (
    id_modal BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    id_varian BIGINT UNSIGNED NOT NULL,

    harga_modal DECIMAL(12,2) NOT NULL,

    berlaku_mulai DATE NOT NULL,

    catatan VARCHAR(255) NULL,

    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_modal_varian (id_varian),
    INDEX idx_modal_berlaku (berlaku_mulai),

    CONSTRAINT fk_modal_varian
        FOREIGN KEY (id_varian)
        REFERENCES varian_produk(id_varian)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT chk_harga_modal_positif
        CHECK (harga_modal > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- ==========================================
-- Tabel pengeluaran operasional usaha
-- ==========================================
-- Mencatat semua biaya operasional usaha
-- seperti listrik, gas, kemasan, minyak goreng, dll.
-- Digunakan sebagai sumber perhitungan laporan bulanan.

CREATE TABLE IF NOT EXISTS pengeluaran_usaha (
    id_pengeluaran BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    tanggal_pengeluaran DATE NOT NULL,

    kategori VARCHAR(100) NOT NULL,

    nominal DECIMAL(12,2) NOT NULL,

    keterangan VARCHAR(255) NULL,

    dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_pengeluaran_tanggal (tanggal_pengeluaran),
    INDEX idx_pengeluaran_kategori (kategori),

    CONSTRAINT chk_nominal_positif
        CHECK (nominal > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
