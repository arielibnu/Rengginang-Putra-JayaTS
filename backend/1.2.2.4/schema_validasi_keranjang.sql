-- ============================================================================
-- SKEMA & QUERY VALIDASI KERANJANG TERHADAP STOK WBS 1.2.2.4 (PPL-11)
-- ============================================================================

-- Query validasi stok aktual yang dieksekusi backend WBS 1.2.2.4:
-- SELECT id_varian, nama_varian, harga_jual, stok, status_aktif
-- FROM varian_produk
-- WHERE (id_varian = ? OR LOWER(nama_varian) = ?)
-- LIMIT 1;

-- Pastikan tabel varian_produk memiliki index untuk performa validasi realtime
CREATE INDEX IF NOT EXISTS idx_varian_nama_stok ON varian_produk(nama_varian, stok, status_aktif);
