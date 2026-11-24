-- ============================================================================
-- SKEMA DATABASE WBS 1.3.3.2: PEMBERITAHUAN PESANAN BARU & PEMBAYARAN SUKSES
-- Database: rengginang_putra_jaya
-- ============================================================================

CREATE TABLE IF NOT EXISTS notifikasi_pesanan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nomor_pesanan VARCHAR(64) NOT NULL,
    nomor_invoice VARCHAR(64) DEFAULT NULL,
    tipe_pemberitahuan VARCHAR(50) NOT NULL DEFAULT 'pesanan_masuk', -- 'pesanan_masuk', 'pesanan_terbayar', 'stok_menipis'
    judul VARCHAR(100) NOT NULL,
    pesan TEXT NOT NULL,
    metode_penerimaan VARCHAR(50) NOT NULL DEFAULT 'Kirim Kurir', -- 'Ambil Sendiri' | 'Kirim Kurir'
    total_transaksi BIGINT NOT NULL DEFAULT 0,
    status_pesanan VARCHAR(30) NOT NULL DEFAULT 'baru', -- 'baru', 'diproses', 'selesai'
    sudah_dibaca TINYINT(1) NOT NULL DEFAULT 0,
    waktu_dibaca DATETIME DEFAULT NULL,
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nomor_pesanan (nomor_pesanan),
    INDEX idx_tipe_pemberitahuan (tipe_pemberitahuan),
    INDEX idx_status_pesanan (status_pesanan),
    INDEX idx_sudah_dibaca (sudah_dibaca),
    INDEX idx_dibuat_pada (dibuat_pada)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- TRIGGER 1: TRIGGER OTOMATIS CATATAN PEMBERITAHUAN SAAT PESANAN BARU MASUK
DROP TRIGGER IF EXISTS trg_catat_pesanan_baru;
DELIMITER $$
CREATE TRIGGER trg_catat_pesanan_baru
AFTER INSERT ON pesanan
FOR EACH ROW
BEGIN
    DECLARE v_metode VARCHAR(50);
    IF NEW.kurir_id = 'pickup' OR LOWER(NEW.kurir_nama) LIKE '%ambil%' THEN
        SET v_metode = 'Ambil Sendiri';
    ELSE
        SET v_metode = 'Kirim Kurir';
    END IF;

    INSERT INTO notifikasi_pesanan (
        nomor_pesanan,
        nomor_invoice,
        tipe_pemberitahuan,
        judul,
        pesan,
        metode_penerimaan,
        total_transaksi,
        status_pesanan,
        sudah_dibaca,
        dibuat_pada
    ) VALUES (
        NEW.nomor_pesanan,
        NEW.nomor_invoice,
        'pesanan_masuk',
        CONCAT('Pesanan Baru: ', NEW.nomor_pesanan),
        CONCAT('Pesanan baru masuk dari ', NEW.nama_penerima, ' dengan metode ', v_metode),
        v_metode,
        NEW.snapshot_total,
        'baru',
        0,
        NOW()
    );
END$$
DELIMITER ;

-- TRIGGER 2: TRIGGER OTOMATIS CATATAN PEMBERITAHUAN SAAT PESANAN SUKSES TERBAYAR
DROP TRIGGER IF EXISTS trg_catat_pesanan_terbayar;
DELIMITER $$
CREATE TRIGGER trg_catat_pesanan_terbayar
AFTER UPDATE ON transaksi_pembayaran_gateway
FOR EACH ROW
BEGIN
    IF (NEW.transaction_status IN ('settlement', 'capture') AND (OLD.transaction_status NOT IN ('settlement', 'capture') OR OLD.transaction_status IS NULL)) THEN
        INSERT INTO notifikasi_pesanan (
            nomor_pesanan,
            nomor_invoice,
            tipe_pemberitahuan,
            judul,
            pesan,
            metode_penerimaan,
            total_transaksi,
            status_pesanan,
            sudah_dibaca,
            dibuat_pada
        ) VALUES (
            NEW.nomor_pesanan,
            NEW.nomor_invoice,
            'pesanan_terbayar',
            CONCAT('Pembayaran Lunas: ', NEW.nomor_pesanan),
            CONCAT('Pembayaran pesanan ', NEW.nomor_pesanan, ' via ', NEW.payment_type, ' telah berhasil diverifikasi gateway.'),
            'Kirim Kurir',
            NEW.gross_amount,
            'baru',
            0,
            NOW()
        );
    END IF;
END$$
DELIMITER ;
