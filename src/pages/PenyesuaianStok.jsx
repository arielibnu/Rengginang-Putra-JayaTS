import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ArrowLeft, Calendar } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STOCK_KEY = 'rengginang_data_stok';
const HISTORY_KEY = 'rengginang_riwayat_mutasi';

const DEFAULT_PRODUCTS = [
  { id: 'udang', nama: 'Rengginang Udang', fisik: 34, reservasi: 2, tersedia: 32 },
  { id: 'bawang', nama: 'Rengginang Bawang', fisik: 40, reservasi: 2, tersedia: 38 },
  { id: 'cumi', nama: 'Rengginang Cumi', fisik: 14, reservasi: 2, tersedia: 12 },
  { id: 'terasi', nama: 'Rengginang Terasi', fisik: 40, reservasi: 2, tersedia: 38 }
];

export default function PenyesuaianStok() {
  const navigate = useNavigate();

  // State stok
  const [daftarStok, setDaftarStok] = useState(() => {
    try {
      const saved = localStorage.getItem(STOCK_KEY);
      return saved ? JSON.parse(saved) : DEFAULT_PRODUCTS;
    } catch (e) {
      return DEFAULT_PRODUCTS;
    }
  });

  // State form penyesuaian sesuai mockup Figma Gambar 3 (Rengginang Cumi)
  const [selectedProductId, setSelectedProductId] = useState('cumi');
  const [stokFisikAkhir, setStokFisikAkhir] = useState(16);
  const [tanggal, setTanggal] = useState('16 November 2026');
  const [keterangan, setKeterangan] = useState('Hasil stok opname');
  const [errorPesan, setErrorPesan] = useState('');
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Ambil varian produk yang dipilih (default: Rengginang Cumi)
  const produkDipilih =
    daftarStok.find((p) => p.id === selectedProductId) || daftarStok[2] || DEFAULT_PRODUCTS[2];

  // Hitung live stok sesudah penyesuaian
  const fisikSekarang = Number(produkDipilih.fisik || 0);
  const reservasi = Number(produkDipilih.reservasi || 0);
  const fisikSesudah = Number(stokFisikAkhir || 0);
  const selisih = fisikSesudah - fisikSekarang;
  const tersediaSesudah = fisikSesudah - reservasi;

  // Validasi & Simpan Penyesuaian Stok
  const handleSimpan = (e) => {
    e.preventDefault();
    setErrorPesan('');

    if (tersediaSesudah < 0) {
      setErrorPesan('Stok tersedia tidak boleh negatif.');
      return;
    }

    const updatedStok = daftarStok.map((item) => {
      if (item.id === produkDipilih.id) {
        return {
          ...item,
          fisik: fisikSesudah,
          tersedia: tersediaSesudah
        };
      }
      return item;
    });

    setDaftarStok(updatedStok);
    localStorage.setItem(STOCK_KEY, JSON.stringify(updatedStok));

    // Catat ke riwayat mutasi
    const riwayatBaru = {
      id: Date.now().toString(),
      tanggal: tanggal,
      produk: produkDipilih.nama,
      jenis: 'penyesuaian',
      jumlah: selisih,
      fisikSebelum: fisikSekarang,
      fisikSesudah: fisikSesudah,
      keterangan: keterangan || 'Hasil stok opname'
    };

    try {
      const existingHistory = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
      existingHistory.unshift(riwayatBaru);
      localStorage.setItem(HISTORY_KEY, JSON.stringify(existingHistory));
    } catch (e) {
      console.warn('Gagal mencatat riwayat:', e);
    }

    navigate('/admin/stok');
  };

  // Logout
  const handleLogout = async () => {
    if (isLoggingOut) return;
    setIsLoggingOut(true);
    try {
      await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });
    } catch (e) {
      console.log('Logout lokal...');
    } finally {
      localStorage.removeItem('ownerAuth');
      navigate('/login-owner');
    }
  };

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Kiri */}
      <aside className="admin-sidebar">
        <div>
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

          <ul className="sidebar-menu">
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
              >
                <Circle size={16} />
                <span>Ringkasan</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/produk')}
              >
                <Square size={16} />
                <span>Kelola Produk</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item active"
                onClick={() => navigate('/admin/stok')}
              >
                <LayoutList size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
              >
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
          </ul>
        </div>

        <div className="sidebar-footer">
          <Link to="/" className="sidebar-footer-link">
            Buka website publik
          </Link>
          <button
            type="button"
            onClick={handleLogout}
            disabled={isLoggingOut}
            className="sidebar-footer-link"
            style={{ color: '#A92631' }}
          >
            {isLoggingOut ? 'Memproses...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* Konten Utama Penyesuaian Stok */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Catat Mutasi
        </div>

        {/* Header Halaman & Tombol Kembali */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">Penyesuaian stok</h1>
            <p className="admin-page-subtitle">
              Samakan jumlah fisik dengan hasil pengecekan persediaan.
            </p>
          </div>

          <Link to="/admin/stok" className="btn-back-stock">
            <ArrowLeft size={16} />
            <span>Data Stok</span>
          </Link>
        </div>

        {/* Grid 2 Kolom Sesuai Figma */}
        <div className="stock-mutation-layout">
          {/* Kolom Kiri: Form Penyesuaian Stok */}
          <div className="stock-table-card" style={{ padding: '2rem 2.25rem' }}>
            {errorPesan && (
              <div className="modal-alert-error" style={{ marginBottom: '1.5rem' }}>
                {errorPesan}
              </div>
            )}

            <form onSubmit={handleSimpan} className="modal-form">
              {/* Jenis Perubahan: 3 Pill Tombol */}
              <div className="form-group">
                <label className="form-label">Jenis perubahan</label>
                <div className="mutation-type-pills">
                  <button
                    type="button"
                    className="mutation-pill-btn"
                    onClick={() => navigate('/admin/stok/mutasi')}
                  >
                    Stok masuk
                  </button>
                  <button
                    type="button"
                    className="mutation-pill-btn"
                    onClick={() => navigate('/admin/stok/keluar')}
                  >
                    Stok keluar
                  </button>
                  <button
                    type="button"
                    className="mutation-pill-btn active"
                  >
                    Penyesuaian
                  </button>
                </div>
              </div>

              {/* Pilih Varian Produk (Default: Rengginang Cumi) */}
              <div className="form-group" style={{ marginTop: '0.75rem' }}>
                <label className="form-label">Pilih varian produk</label>
                <select
                  className="form-input"
                  value={selectedProductId}
                  onChange={(e) => setSelectedProductId(e.target.value)}
                  style={{ backgroundColor: '#FFFFFF', cursor: 'pointer' }}
                >
                  {daftarStok.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.nama}
                    </option>
                  ))}
                </select>
              </div>

              {/* Baris 2 Kolom: Stok Fisik Akhir & Tanggal */}
              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: '1fr 1fr',
                  gap: '1.25rem',
                  marginTop: '0.75rem'
                }}
              >
                {/* Stok Fisik Akhir Hasil Pengecekan */}
                <div className="form-group">
                  <label className="form-label">Stok fisik akhir hasil pengecekan</label>
                  <div className="input-with-suffix">
                    <input
                      type="number"
                      className="form-input"
                      value={stokFisikAkhir}
                      onChange={(e) => setStokFisikAkhir(Number(e.target.value))}
                      min="0"
                      required
                    />
                    <span className="input-suffix-text">unit</span>
                  </div>
                </div>

                {/* Tanggal Pencatatan */}
                <div className="form-group">
                  <label className="form-label">Tanggal pencatatan</label>
                  <div className="input-with-suffix">
                    <input
                      type="text"
                      className="form-input"
                      value={tanggal}
                      onChange={(e) => setTanggal(e.target.value)}
                    />
                    <span className="input-suffix-text">
                      <Calendar size={15} />
                    </span>
                  </div>
                </div>
              </div>

              {/* Keterangan */}
              <div className="form-group" style={{ marginTop: '0.75rem' }}>
                <label className="form-label">Keterangan</label>
                <input
                  type="text"
                  className="form-input"
                  value={keterangan}
                  onChange={(e) => setKeterangan(e.target.value)}
                  placeholder="Contoh: Hasil stok opname"
                />
              </div>

              {/* Kotak Catatan Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FAF5EE',
                  border: '1px solid #EFE6D9',
                  borderRadius: '12px',
                  padding: '1rem 1.25rem',
                  marginTop: '1.25rem',
                  marginBottom: '1rem'
                }}
              >
                <div style={{ fontWeight: 800, fontSize: '0.875rem', color: '#43281C', marginBottom: '0.2rem' }}>
                  Catatan:
                </div>
                <div style={{ fontSize: '0.825rem', color: '#8E7D73' }}>
                  Setiap penyimpanan harus dicatat pada riwayat mutasi.
                </div>
              </div>

              {/* Tombol Aksi: Batal & Simpan */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.85rem', marginTop: '1.25rem' }}>
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => navigate('/admin/stok')}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-confirm">
                  Simpan
                </button>
              </div>
            </form>
          </div>

          {/* Kolom Kanan: Ringkasan Stok & Aturan Persediaan */}
          <div>
            {/* Kartu Ringkasan Stok (Rengginang Cumi) */}
            <div className="stock-summary-card">
              <h2 className="rules-card-title" style={{ marginBottom: '0.2rem' }}>
                Ringkasan stok
              </h2>
              <div style={{ fontWeight: 800, color: '#43281C', fontSize: '1rem' }}>
                {produkDipilih.nama}
              </div>

              <div className="stock-summary-rows">
                <div className="stock-summary-row">
                  <span>Stok fisik sekarang</span>
                  <span className="stock-summary-val-bold">{fisikSekarang} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Stok reservasi</span>
                  <span className="stock-summary-val-bold">{reservasi} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Stok tersedia</span>
                  <span className="stock-summary-val-bold">{fisikSekarang - reservasi} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Fisik sesudah perubahan</span>
                  <span className="stock-summary-val-green">{fisikSesudah} unit</span>
                </div>
              </div>
            </div>

            {/* Kartu Aturan Persediaan */}
            <div className="rules-card">
              <h3 className="rules-card-title">Aturan persediaan</h3>
              <div className="rules-list">
                <div>Tersedia = Fisik - Reservasi</div>
                <div>Stok tersedia tidak boleh negatif.</div>
                <div>Perubahan dicatat sebagai riwayat.</div>
              </div>
            </div>
          </div>
        </div>

        {/* Catatan Kaki Prototype */}
        <p className="stock-footnote">
          Nilai pada formulir adalah contoh; pengisian & validasi nyata dikerjakan saat implementasi.
        </p>
      </main>
    </div>
  );
}
