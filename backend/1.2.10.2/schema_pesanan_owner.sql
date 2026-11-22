-- ============================================================================
-- WBS 1.2.10.2 - SKEMA DATABASE MANAJEMEN PESANAN OWNER (PPL-21)
-- Database: rengginang_putra_jaya
-- ============================================================================

-- Pastikan kolom status_pemrosesan tersedia pada tabel pesanan
ALTER TABLE pesanan 
ADD COLUMN IF NOT EXISTS status_pemrosesan VARCHAR(30) NOT NULL DEFAULT 'menunggu_diproses';

-- Index untuk mempercepat filter status pemrosesan dan rentang tanggal
CREATE INDEX IF NOT EXISTS idx_status_pemrosesan ON pesanan(status_pemrosesan);
CREATE INDEX IF NOT EXISTS idx_pesanan_dibuat ON pesanan(dibuat_pada);
