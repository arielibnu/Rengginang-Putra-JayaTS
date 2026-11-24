import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ChevronLeft, Search } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

// Data riwayat default sesuai mockup resmi Figma Gambar 4
const DEFAULT_HISTORY = [
  {
    id: '1',
    waktu: '16 Nov, 16:05',
    produk: 'Rengginang Udang',
    catatan: 'Produksi selesai',
    jenis: 'Masuk',
    perubahan: '+10',
    fisikAwal: 24,
    fisikAkhir: 34
  },
  {
    id: '2',
    waktu: '16 Nov, 14:20',
    produk: 'Rengginang Udang',
    catatan: 'Penjualan toko',
    jenis: 'Keluar',
    perubahan: '-5',
    fisikAwal: 29,
    fisikAkhir: 24
  },
  {
    id: '3',
    waktu: '15 Nov, 09:40',
    produk: 'Rengginang Bawang',
    catatan: 'Penjualan toko',
    jenis: 'Keluar',
    perubahan: '-4',
    fisikAwal: 44,
    fisikAkhir: 40
  },
  {
    id: '4',
    waktu: '14 Nov, 15:05',
    produk: 'Rengginang Cumi',
    catatan: 'Hasil stok opname',
    jenis: 'Penyesuaian',
    perubahan: '-2',
    fisikAwal: 16,
    fisikAkhir: 14
  }
];

export default function RiwayatStok() {
  const navigate = useNavigate();

  // Filter pencarian
  const [keyword, setKeyword] = useState('');
  const [filterJenis, setFilterJenis] = useState('semua');
  const [filterTanggal, setFilterTanggal] = useState('semua');
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Filter data sesuai kata kunci & jenis mutasi
  const riwayatTampil = DEFAULT_HISTORY.filter((item) => {
    const matchKeyword =
      item.produk.toLowerCase().includes(keyword.toLowerCase()) ||
      item.catatan.toLowerCase().includes(keyword.toLowerCase());

    const matchJenis =
      filterJenis === 'semua' || item.jenis.toLowerCase() === filterJenis.toLowerCase();

    return matchKeyword && matchJenis;
  });

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

      {/* Konten Utama Riwayat Stok */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Riwayat Stok
        </div>

        {/* Header Halaman & Tombol Kembali */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">Riwayat mutasi stok</h1>
            <p className="admin-page-subtitle">
              Catatan penambahan, pengeluaran, dan penyesuaian setiap varian.
            </p>
          </div>

          <Link to="/admin/stok" className="btn-back-stock">
            <ChevronLeft size={16} />
            <span>Data Stok</span>
          </Link>
        </div>

        {/* Filter Bar Pencarian Sesuai Figma */}
        <div className="history-filter-card">
          <div className="history-filter-grid">
            {/* Input Pencarian */}
            <input
              type="text"
              className="form-input"
              placeholder="Cari produk atau catatan..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
            />

            {/* Filter Jenis Mutasi */}
            <select
              className="form-input"
              value={filterJenis}
              onChange={(e) => setFilterJenis(e.target.value)}
              style={{ backgroundColor: '#FFFFFF', cursor: 'pointer' }}
            >
              <option value="semua">Semua jenis mutasi</option>
              <option value="masuk">Masuk</option>
              <option value="keluar">Keluar</option>
              <option value="penyesuaian">Penyesuaian</option>
            </select>

            {/* Filter Tanggal */}
            <select
              className="form-input"
              value={filterTanggal}
              onChange={(e) => setFilterTanggal(e.target.value)}
              style={{ backgroundColor: '#FFFFFF', cursor: 'pointer' }}
            >
              <option value="semua">Semua tanggal</option>
              <option value="hari-ini">Hari ini</option>
              <option value="minggu-ini">Minggu ini</option>
              <option value="bulan-ini">Bulan ini</option>
            </select>
          </div>
        </div>

        {/* Tabel Riwayat Mutasi Sesuai Figma */}
        <div className="stock-table-card">
          <table className="stock-table">
            <thead>
              <tr>
                <th style={{ width: '18%' }}>Waktu</th>
                <th style={{ width: '32%' }}>Produk</th>
                <th style={{ width: '15%' }}>Jenis</th>
                <th style={{ width: '13%' }}>Perubahan</th>
                <th style={{ width: '11%', textAlign: 'center' }}>Fisik awal</th>
                <th style={{ width: '11%', textAlign: 'center' }}>Fisik akhir</th>
              </tr>
            </thead>
            <tbody>
              {riwayatTampil.map((row) => (
                <tr key={row.id}>
                  {/* Waktu */}
                  <td style={{ color: '#66564C', fontSize: '0.9rem' }}>
                    {row.waktu}
                  </td>

                  {/* Produk & Catatan */}
                  <td>
                    <div style={{ fontWeight: 800, color: '#43281C' }}>
                      {row.produk}
                    </div>
                    <div className="history-item-subtitle">{row.catatan}</div>
                  </td>

                  {/* Jenis Mutasi Badge */}
                  <td>
                    <span
                      className={`badge-history-type ${
                        row.jenis === 'Masuk'
                          ? 'badge-type-masuk'
                          : row.jenis === 'Keluar'
                          ? 'badge-type-keluar'
                          : 'badge-type-penyesuaian'
                      }`}
                    >
                      {row.jenis}
                    </span>
                  </td>

                  {/* Perubahan (+10 / -5 / dst) */}
                  <td>
                    <span
                      className={
                        row.perubahan.startsWith('+')
                          ? 'history-diff-positive'
                          : 'history-diff-negative'
                      }
                    >
                      {row.perubahan}
                    </span>
                  </td>

                  {/* Fisik Awal */}
                  <td style={{ textAlign: 'center', color: '#66564C' }}>
                    {row.fisikAwal}
                  </td>

                  {/* Fisik Akhir */}
                  <td style={{ textAlign: 'center', color: '#43281C', fontWeight: 600 }}>
                    {row.fisikAkhir}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Banner Penjelasan Sesuai Figma */}
        <div className="banner-history-info">
          <div className="banner-history-info-title">
            Riwayat membantu owner menelusuri perubahan stok setiap produk.
          </div>
          <div className="banner-history-info-sub">
            Semua transaksi di layar ini merupakan data contoh untuk kebutuhan desain UI/UX.
          </div>
        </div>
      </main>
    </div>
  );
}
