-- ============================================================================
-- SKEMA SINKRONISASI IDENTITAS TRANSAKSI GATEWAY DENGAN NOMOR PESANAN
-- WBS 1.3.1.4 (BE)
-- ============================================================================

-- Relasi foreign key atau indeks pencarian cepat nomor_pesanan dan transaction_id
ALTER TABLE transaksi_pembayaran_gateway 
ADD INDEX IF NOT EXISTS idx_transaksi_relasi (nomor_pesanan, transaction_id);
