import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, Check } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

export default function MutasiTersimpan() {
  const navigate = useNavigate();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

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

      {/* Konten Utama */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Mutasi Tersimpan
        </div>

        {/* Header Halaman */}
        <div className="admin-header">
          <h1 className="admin-page-title">Mutasi stok berhasil dicatat</h1>
          <p className="admin-page-subtitle">
            Contoh tampilan hasil setelah owner menyimpan stok keluar.
          </p>
        </div>

        {/* Kartu Sukses Sesuai Figma Gambar 1 */}
        <div className="stock-success-card">
          {/* Ikon Centang Hijau Lingkaran */}
          <div className="success-circle-icon">
            <Check size={36} strokeWidth={3} />
          </div>

          <h2 className="stock-success-title">Stok berhasil diperbarui</h2>
          <p className="stock-success-sub">
            Perubahan tersimpan sebagai riwayat mutasi.
          </p>

          {/* Kotak Rincian Mutasi Sesuai Figma */}
          <div className="mutation-detail-box">
            <div className="mutation-detail-row">
              <div>
                <div className="mutation-detail-col-title">Produk</div>
                <div className="mutation-detail-col-val">Rengginang Udang</div>
              </div>

              <div>
                <div className="mutation-detail-col-title">Mutasi</div>
                <div className="mutation-detail-col-val" style={{ color: '#A92631' }}>
                  Keluar 5 unit
                </div>
              </div>

              <div>
                <div className="mutation-detail-col-title">Stok fisik</div>
                <div className="mutation-detail-col-val" style={{ color: '#1E833B' }}>
                  34 → 29
                </div>
              </div>
            </div>

            <div className="mutation-detail-note">
              Stok tersedia setelah mutasi: 27 unit (2 unit direservasi).
            </div>
          </div>

          {/* Tombol Aksi */}
          <div className="success-actions-row">
            <button
              type="button"
              className="btn-modal-confirm"
              onClick={() => navigate('/admin/stok/riwayat')}
              style={{ padding: '0.8rem 2rem' }}
            >
              Lihat riwayat
            </button>

            <button
              type="button"
              className="btn-modal-cancel"
              onClick={() => navigate('/admin/stok')}
              style={{ padding: '0.8rem 1.8rem' }}
            >
              Kembali ke data stok
            </button>
          </div>

          <p className="stock-footnote" style={{ marginTop: '1rem' }}>
            Simulasi prototype: tidak ada perubahan database.
          </p>
        </div>
      </main>
    </div>
  );
}
