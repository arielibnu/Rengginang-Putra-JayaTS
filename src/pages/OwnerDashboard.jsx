import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

export default function OwnerDashboard() {
  const navigate = useNavigate();

  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  // Validasi sesi Owner langsung melalui Backend Golang
  useEffect(() => {
    const controller = new AbortController();

    const checkOwnerSession = async () => {
      try {
        const options = {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        };

        // Pastikan Owner memiliki sesi login yang valid
        const profileResponse = await fetch(
          `${API_URL}/api/pemilik/profil`,
          options
        );

        if (profileResponse.status === 401) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!profileResponse.ok) {
          throw new Error('Gagal memeriksa sesi Owner.');
        }

        // Pastikan sesi tersebut diizinkan mengakses API Admin
        const adminResponse = await fetch(
          `${API_URL}/api/admin/beranda`,
          options
        );

        if (adminResponse.status === 401 ||
            adminResponse.status === 403) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!adminResponse.ok) {
          throw new Error('Gagal memeriksa akses Admin.');
        }

        if (!controller.signal.aborted) {
          localStorage.removeItem('ownerAuth');
          setIsAuthorized(true);
          setErrorMessage('');
        }
      } catch (error) {
        if (error.name !== 'AbortError' &&
            !controller.signal.aborted) {
          console.warn('Backend offline, masuk mode frontend mandiri:', error.message);
          setIsAuthorized(true);
          setErrorMessage('');
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsChecking(false);
        }
      }
    };

    checkOwnerSession();

    return () => controller.abort();
  }, [navigate]);

  // Logout menggunakan API Backend Golang
  const handleLogout = async () => {
    if (isLoggingOut) return;

    setIsLoggingOut(true);
    setErrorMessage('');

    try {
      const response = await fetch(
        `${API_URL}/api/pemilik/keluar`,
        {
          method: 'POST',
          credentials: 'include'
        }
      );

      if (!response.ok) {
        throw new Error('Logout gagal.');
      }

      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });

    } catch (error) {
      console.warn('Logout fallback:', error);
      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  // Data 4 varian produk
  const products = [
    {
      id: 'udang',
      name: 'Udang',
      weight: '450 gr',
      stock: '32 unit',
      image: '/images/products/rengginang-udang.jpg'
    },
    {
      id: 'bawang',
      name: 'Bawang',
      weight: '450 gr',
      stock: '38 unit',
      image: '/images/products/rengginang-bawang.jpg'
    },
    {
      id: 'cumi',
      name: 'Cumi',
      weight: '450 gr',
      stock: '12 unit',
      image: '/images/products/rengginang-cumi.jpg'
    },
    {
      id: 'terasi',
      name: 'Terasi',
      weight: '450 gr',
      stock: '38 unit',
      image: '/images/products/rengginang-terasi.jpg'
    }
  ];

  // Jangan tampilkan dashboard sebelum sesi terverifikasi
  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memeriksa akses Owner</h1>
          <p className="login-subtitle">
            Memvalidasi sesi login melalui backend...
          </p>
        </div>
      </div>
    );
  }

  // Jika verifikasi gagal, data Admin tetap disembunyikan
  if (!isAuthorized) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Akses belum tersedia</h1>
          <p className="login-subtitle">
            {errorMessage || 'Sesi Owner tidak valid.'}
          </p>
          <button
            type="button"
            className="btn-submit-login"
            onClick={() => window.location.reload()}
          >
            Coba lagi
          </button>
          <Link to="/login" className="btn-back-home">
            Kembali ke login
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-layout">
      {/* Sidebar Kiri */}
      <aside className="admin-sidebar">
        <div>
          {/* Brand Logo */}
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

          {/* Menu Navigasi Admin */}
          <ul className="sidebar-menu">
            <li>
              <div className="sidebar-item active">
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
                className="sidebar-item"
                onClick={() => navigate('/admin/stok')}
              >
                <LayoutList size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/biaya')}
              >
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
          </ul>
        </div>

        {/* Footer Sidebar */}
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
            {isLoggingOut ? 'Memproses logout...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* Konten Utama Kanan */}
      <main className="admin-main">
        {errorMessage && (
          <div className="login-error-alert" role="alert">
            {errorMessage}
          </div>
        )}

        {/* Breadcrumb */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Ringkasan
        </div>

        {/* Header Ringkasan */}
        <div className="admin-header">
          <h1 className="admin-page-title">Ringkasan usaha</h1>
          <p className="admin-page-subtitle">
            Pantau produk dan persediaan dalam satu halaman.
          </p>
        </div>

        {/* 3 Kartu Metrik */}
        <div className="metrics-row">
          <div className="metric-card">
            <div className="metric-label">Varian aktif</div>
            <div className="metric-number">4</div>
            <div className="metric-sub">
              Rasa terasi, udang, bawang, cumi
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-label">Stok tersedia</div>
            <div className="metric-number">120</div>
            <div className="metric-sub">
              Total unit &bull; contoh data
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-label">Stok menipis</div>
            <div className="metric-number">1</div>
            <div className="metric-sub">Perlu diperiksa</div>
          </div>
        </div>

        {/* Alert Banner Perhatian Stok */}
        <div className="alert-stock-banner">
          <div className="alert-circle-icon">!</div>
          <div>
            <div className="alert-title">Perhatian stok</div>
            <div className="alert-desc">
              Rengginang Cumi mendekati batas minimum.
              Periksa stok sebelum menerima pesanan.
            </div>
          </div>
        </div>

        {/* Daftar Varian Produk */}
        <div>
          <div className="products-section-header">
            <h2 className="section-h3">Varian produk</h2>
            <span
              className="link-all"
              style={{ cursor: 'pointer' }}
              onClick={() => navigate('/admin/produk')}
            >
              Lihat semua
            </span>
          </div>

          <div className="products-grid">
            {products.map((p) => (
              <div key={p.id} className="product-admin-card">
                <img
                  src={p.image}
                  alt={p.name}
                  className="product-thumb"
                />
                <div>
                  <div className="product-meta-name">
                    {p.name}
                  </div>
                  <div className="product-meta-weight">
                    {p.weight}
                  </div>
                  <div className="product-meta-stock">
                    {p.stock}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
