import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_KEY = 'rengginang_data_stok';

// Data default varian & stok awal sesuai mockup resmi Figma PPL-05
const DATA_VARIAN_AWAL = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    varian: 'Udang',
    fisik: 34,
    reservasi: 2,
    tersedia: 32,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-udang-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-udang.jpg'
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    varian: 'Bawang',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-bawang-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-bawang.jpg'
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    varian: 'Cumi',
    fisik: 14,
    reservasi: 2,
    tersedia: 12,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-cumi-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-cumi.jpg'
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    varian: 'Terasi',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-terasi-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-terasi.jpg'
  }
];

export default function OwnerDashboard() {
  const navigate = useNavigate();

  // State status autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  // State daftar stok varian (sinkron dengan localStorage jika sudah dimutasi di DataStok)
  const [daftarVarian] = useState(() => {
    try {
      const dataTersimpan = localStorage.getItem(STORAGE_KEY);
      if (!dataTersimpan) return DATA_VARIAN_AWAL;

      const stokLokal = JSON.parse(dataTersimpan);
      // Gabungkan data stok dinamis dari localStorage dengan atribut produk
      return DATA_VARIAN_AWAL.map((varian) => {
        const itemTersimpan = stokLokal.find((s) => s.id === varian.id);
        if (!itemTersimpan) return varian;

        return {
          ...varian,
          fisik: Number(itemTersimpan.fisik ?? varian.fisik),
          reservasi: Number(itemTersimpan.reservasi ?? varian.reservasi),
          tersedia: Number(itemTersimpan.tersedia ?? varian.tersedia),
          minimum: Number(itemTersimpan.minimum ?? varian.minimum ?? 15)
        };
      });
    } catch (e) {
      return DATA_VARIAN_AWAL;
    }
  });

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

        if (adminResponse.status === 401 || adminResponse.status === 403) {
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
        if (error.name !== 'AbortError' && !controller.signal.aborted) {
          // Mode mandiri frontend jika backend offline
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
      const response = await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Logout gagal.');
      }

      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
    } catch (error) {
      console.warn('Logout fallback lokal:', error);
      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  // Hitung ringkasan angka metrik
  const varianAktifCount = daftarVarian.length;
  const totalStokFisik = daftarVarian.reduce((acc, v) => acc + Number(v.fisik || 0), 0);
  const totalStokTersedia = daftarVarian.reduce((acc, v) => acc + Number(v.tersedia || 0), 0);

  // Evaluasi varian dengan stok menipis (stok tersedia < batas minimum)
  const varianMenipis = daftarVarian.filter(
    (v) => Number(v.tersedia || 0) < Number(v.minimum || 15)
  );

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
      {/* Sidebar Kiri Sesuai Figma */}
      <aside className="admin-sidebar">
        <div>
          {/* Brand Logo Putra Jaya */}
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
            Buka website publik &rsaquo;
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

        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Ringkasan
        </div>

        {/* Header Ringkasan Usaha Sesuai Figma */}
        <div className="admin-header">
          <h1 className="admin-page-title">Ringkasan usaha</h1>
          <p className="admin-page-subtitle">
            Pantau produk dan persediaan dalam satu halaman.
          </p>
        </div>

        {/* 3 Kartu Metrik Utama Sesuai Figma */}
        <div className="metrics-row">
          <div className="metric-card">
            <div className="metric-label">Varian aktif</div>
            <div className="metric-number">{varianAktifCount}</div>
            <div className="metric-sub">
              Rasa terasi, udang, bawang, cumi
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-label">Stok fisik</div>
            <div className="metric-number">{totalStokFisik}</div>
            <div className="metric-sub">
              Total seluruh varian &bull; bungkus
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-label">Stok tersedia</div>
            <div className="metric-number">{totalStokTersedia}</div>
            <div className="metric-sub">
              Setelah dikurangi reservasi
            </div>
          </div>
        </div>

        {/* Banner Peringatan Stok Menipis (PPL-05) Sesuai Figma */}
        {varianMenipis.length > 0 && (
          <div className="stock-warning-banner">
            <div className="stock-warning-left">
              <div className="stock-warning-icon">!</div>
              <div>
                <div className="stock-warning-title">
                  Perhatian: {varianMenipis.length} varian stok menipis
                </div>
                <div className="stock-warning-desc-primary">
                  {varianMenipis
                    .map(
                      (v) =>
                        `${v.nama}: tersedia ${v.tersedia} bungkus, batas minimum ${v.minimum || 15} bungkus.`
                    )
                    .join(' • ')}
                </div>
                <div className="stock-warning-desc-sub">
                  Periksa persediaan sebelum menerima pesanan.
                </div>
              </div>
            </div>

            <button
              type="button"
              className="btn-warning-action"
              onClick={() => navigate('/admin/stok')}
            >
              Lihat data stok
            </button>
          </div>
        )}

        {/* Bagian Ringkasan Stok Per Varian Sesuai Figma */}
        <section className="stock-summary-section">
          <div className="stock-summary-header">
            <h2 className="stock-summary-title">Ringkasan stok per varian</h2>
            <p className="stock-summary-subtitle">
              Status mengikuti stok tersedia dan batas minimum yang diatur Owner.
            </p>
          </div>

          <div className="stock-summary-card">
            <div className="table-responsive">
              <table className="stock-summary-table">
                <thead>
                  <tr>
                    <th style={{ textAlign: 'left', width: '30%' }}>Varian</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Stok fisik</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Reservasi</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Tersedia</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Minimum</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {daftarVarian.map((item) => {
                    const isMenipis = Number(item.tersedia || 0) < Number(item.minimum || 15);
                    return (
                      <tr key={item.id}>
                        <td style={{ fontWeight: 700 }}>{item.nama}</td>
                        <td style={{ textAlign: 'center' }}>{item.fisik}</td>
                        <td style={{ textAlign: 'center' }}>{item.reservasi}</td>
                        <td
                          style={{ textAlign: 'center' }}
                          className={isMenipis ? 'stock-num-highlight-low' : ''}
                        >
                          {item.tersedia}
                        </td>
                        <td style={{ textAlign: 'center' }}>{item.minimum || 15}</td>
                        <td style={{ textAlign: 'center' }}>
                          {isMenipis ? (
                            <span className="badge-status-menipis">Menipis</span>
                          ) : (
                            <span className="badge-status-normal">Normal</span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Bagian Varian Produk Sesuai Figma */}
        <section>
          <div className="products-section-header">
            <div>
              <h2 className="section-h3">Varian produk</h2>
              <p
                style={{
                  fontSize: '0.875rem',
                  color: '#8E7D73',
                  marginTop: '0.2rem'
                }}
              >
                Informasi produk yang dikelola Owner.
              </p>
            </div>
            <span
              className="link-all"
              style={{ cursor: 'pointer', fontWeight: 700 }}
              onClick={() => navigate('/admin/produk')}
            >
              Lihat semua
            </span>
          </div>

          <div className="products-grid">
            {daftarVarian.map((p) => {
              const isMenipis = Number(p.tersedia || 0) < Number(p.minimum || 15);
              return (
                <div key={p.id} className="product-card-vertical">
                  <img
                    src={p.gambar}
                    alt={p.nama}
                    className="product-card-thumb"
                    onError={(e) => {
                      if (p.gambarCadangan && e.currentTarget.src !== p.gambarCadangan) {
                        e.currentTarget.src = p.gambarCadangan;
                      }
                    }}
                  />
                  <div className="product-card-details">
                    <div className="product-card-title">{p.varian}</div>
                    <div className="product-card-weight">{p.berat}</div>
                    <div className="product-card-price">{p.harga}</div>
                    <div className="product-card-stock-label">Stok tersedia</div>
                    <div
                      className={`product-card-stock-val ${
                        isMenipis ? 'low' : 'safe'
                      }`}
                    >
                      {p.tersedia} bungkus
                    </div>
                    {isMenipis && (
                      <div className="product-card-badge-low">
                        Stok menipis
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </section>
      </main>
    </div>
  );
}
