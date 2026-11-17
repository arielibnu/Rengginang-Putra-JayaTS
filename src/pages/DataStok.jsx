import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, Plus } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_KEY = 'rengginang_data_stok';

// Data awal persediaan stok 4 varian sesuai mockup resmi Figma PPL-04
const INITIAL_STOCK_DATA = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    fisik: 34,
    reservasi: 2,
    tersedia: 32
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    fisik: 40,
    reservasi: 2,
    tersedia: 38
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    fisik: 14,
    reservasi: 2,
    tersedia: 12
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    fisik: 40,
    reservasi: 2,
    tersedia: 38
  }
];

export default function DataStok() {
  const navigate = useNavigate();

  // State daftar stok
  const [daftarStok, setDaftarStok] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      return saved ? JSON.parse(saved) : INITIAL_STOCK_DATA;
    } catch (e) {
      return INITIAL_STOCK_DATA;
    }
  });

  // State status akses Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Simpan data stok ke localStorage jika ada pembaruan
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(daftarStok));
    } catch (e) {
      console.warn('Gagal menyimpan ke localStorage:', e);
    }
  }, [daftarStok]);

  // Validasi sesi login Owner
  useEffect(() => {
    const controller = new AbortController();

    const cekSesiOwner = async () => {
      try {
        const response = await fetch(`${API_URL}/api/pemilik/profil`, {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        });

        if (response.status === 401) {
          navigate('/login-owner', { replace: true });
          return;
        }

        setIsAuthorized(true);
      } catch (err) {
        // Jika backend belum aktif, tetap izinkan pratinjau frontend mandiri
        setIsAuthorized(true);
      } finally {
        setIsChecking(false);
      }
    };

    cekSesiOwner();
    return () => controller.abort();
  }, [navigate]);

  // Hitung total ringkasan metrik
  const totalFisik = daftarStok.reduce((total, item) => total + Number(item.fisik || 0), 0);
  const totalReservasi = daftarStok.reduce((total, item) => total + Number(item.reservasi || 0), 0);
  const totalTersedia = daftarStok.reduce((total, item) => total + Number(item.tersedia || 0), 0);

  // Fungsi logout
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

  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memuat data persediaan</h1>
          <p className="login-subtitle">Menyiapkan tabel stok...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Kiri */}
      <aside className="admin-sidebar">
        <div>
          {/* Logo Brand Putra Jaya */}
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

          {/* Daftar Menu Admin */}
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
              <div className="sidebar-item active">
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

        {/* Menu Footer Sidebar */}
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

      {/* Konten Utama Persediaan Stok */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Persediaan
        </div>

        {/* Header Halaman & Tombol Catat Mutasi */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">Persediaan stok</h1>
            <p className="admin-page-subtitle">
              Pantau stok fisik, reservasi, dan stok yang siap dijual.
            </p>
          </div>

          <button
            type="button"
            className="btn-catat-mutasi"
            onClick={() => navigate('/admin/stok/mutasi')}
          >
            <Plus size={16} />
            <span>Catat mutasi</span>
          </button>
        </div>

        {/* 3 Kartu Ringkasan Metrik Sesuai Figma */}
        <div className="metrics-row">
          {/* Total Fisik */}
          <div className="metric-card">
            <div className="metric-label">Total fisik</div>
            <div className="metric-number">{totalFisik}</div>
            <div className="metric-sub">Total {daftarStok.length} varian</div>
          </div>

          {/* Direservasi */}
          <div className="metric-card">
            <div className="metric-label">Direservasi</div>
            <div className="metric-number">{totalReservasi}</div>
            <div className="metric-sub">Belum tersedia untuk dijual</div>
          </div>

          {/* Stok Tersedia */}
          <div className="metric-card">
            <div className="metric-label">Stok tersedia</div>
            <div className="metric-number">{totalTersedia}</div>
            <div className="metric-sub">Fisik dikurangi reservasi</div>
          </div>
        </div>

        {/* Banner Penjelasan Rumus Stok Sesuai Figma */}
        <div className="stock-formula-banner">
          <div className="stock-formula-title">
            Stok tersedia = stok fisik - stok reservasi
          </div>
          <div className="stock-formula-desc">
            Angka reservasi ditampilkan sebagai informasi. Pembuatan reservasi pesanan merupakan fitur PPL-18.
          </div>
        </div>

        {/* Tabel Data Persediaan Stok Sesuai Figma */}
        <div className="stock-table-card">
          <table className="stock-table">
            <thead>
              <tr>
                <th style={{ width: '45%' }}>Produk</th>
                <th style={{ textAlign: 'center', width: '18%' }}>Fisik</th>
                <th style={{ textAlign: 'center', width: '18%' }}>Reservasi</th>
                <th style={{ textAlign: 'center', width: '19%' }}>Tersedia</th>
              </tr>
            </thead>
            <tbody>
              {daftarStok.map((item) => (
                <tr key={item.id}>
                  <td className="stock-item-name">{item.nama}</td>
                  <td style={{ textAlign: 'center' }}>{item.fisik}</td>
                  <td style={{ textAlign: 'center' }}>{item.reservasi}</td>
                  <td style={{ textAlign: 'center' }} className="stock-num-green">
                    {item.tersedia}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Bagian Riwayat Perubahan Stok Sesuai Figma */}
        <div className="stock-history-section">
          <div>
            <div className="stock-history-title">Riwayat perubahan stok</div>
            <div className="stock-history-sub">
              Setiap mutasi masuk, keluar, dan penyesuaian tersimpan di riwayat.
            </div>
          </div>

          <button
            type="button"
            className="btn-lihat-riwayat"
            onClick={() => navigate('/admin/stok/riwayat')}
          >
            Lihat riwayat
          </button>
        </div>

        {/* Catatan Kaki Prototype */}
        <p className="stock-footnote">
          Data pada layar ini adalah contoh untuk prototype Figma, bukan stok aktual.
        </p>
      </main>
    </div>
  );
}
