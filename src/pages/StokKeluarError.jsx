import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

export default function StokKeluarError() {
  const navigate = useNavigate();
  const [jumlah, setJumlah] = useState(40);
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
          Dashboard Owner / Data Stok / Catat Mutasi
        </div>

        {/* Header Halaman */}
        <div className="admin-header">
          <h1 className="admin-page-title">Stok keluar tidak dapat disimpan</h1>
          <p className="admin-page-subtitle">
            Contoh kondisi validasi saat jumlah keluar melebihi stok tersedia.
          </p>
        </div>

        {/* Grid 2 Kolom Sesuai Figma Gambar 2 */}
        <div className="stock-mutation-layout">
          {/* Kolom Kiri: Form dengan Peringatan Error */}
          <div className="stock-table-card" style={{ padding: '2rem 2.25rem' }}>
            {/* Alert Banner Merah Validasi Stok Tidak Cukup */}
            <div className="stock-alert-card">
              <div className="stock-alert-card-title">
                ! Stok tersedia tidak mencukupi
              </div>
              <div className="stock-alert-card-desc">
                Rengginang Udang hanya memiliki 32 unit yang tersedia.
                <br />
                Permintaan keluar 40 unit tidak bisa diproses.
              </div>
            </div>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                navigate('/admin/stok/keluar');
              }}
              className="modal-form"
            >
              {/* Varian Produk */}
              <div className="form-group">
                <label className="form-label">Varian produk</label>
                <select
                  className="form-input"
                  defaultValue="udang"
                  style={{ backgroundColor: '#FFFFFF' }}
                >
                  <option value="udang">Rengginang Udang</option>
                  <option value="bawang">Rengginang Bawang</option>
                  <option value="cumi">Rengginang Cumi</option>
                  <option value="terasi">Rengginang Terasi</option>
                </select>
              </div>

              {/* Jumlah Stok Keluar (Error Border) */}
              <div className="form-group" style={{ marginTop: '0.75rem' }}>
                <label className="form-label">Jumlah stok keluar</label>
                <div className="input-with-suffix">
                  <input
                    type="number"
                    className="form-input has-error"
                    value={jumlah}
                    onChange={(e) => setJumlah(Number(e.target.value))}
                    style={{ borderColor: '#E05252' }}
                  />
                  <span className="input-suffix-text">unit</span>
                </div>
                <div className="input-helper-error">
                  Maksimal 32 unit karena 2 unit sedang direservasi.
                </div>
              </div>

              {/* Tombol Aksi: Batal & Ubah Jumlah */}
              <div style={{ display: 'flex', gap: '1rem', marginTop: '2.5rem' }}>
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => navigate('/admin/stok')}
                  style={{ flex: 1 }}
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="btn-modal-confirm"
                  style={{ flex: 1.5 }}
                >
                  Ubah jumlah
                </button>
              </div>
            </form>
          </div>

          {/* Kolom Kanan: Card Aturan Stok Sesuai Figma */}
          <div>
            <div className="rules-card">
              <h3 className="rules-card-title">Aturan stok</h3>
              <div className="rules-list" style={{ gap: '0.65rem' }}>
                <div>Fisik: 34 unit</div>
                <div>Reservasi: 2 unit</div>
                <div style={{ color: '#1E833B', fontWeight: 800 }}>
                  Tersedia: 32 unit
                </div>
                <div style={{ color: '#B91C1C', fontWeight: 600, marginTop: '0.5rem' }}>
                  Stok tidak boleh negatif.
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Catatan Kaki Prototype */}
        <p className="stock-footnote">
          Peringatan ini adalah simulasi tampilan, bukan validasi sistem yang sudah aktif.
        </p>
      </main>
    </div>
  );
}
