
import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/home.css';

const API_URL = 'http://localhost:8080';

export default function AkunPelanggan() {
  const navigate = useNavigate();

  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [loggingOut, setLoggingOut] = useState(false);
  const [error, setError] = useState('');

  // Mengambil data pelanggan dari sesi backend Go.
  useEffect(() => {
    const controller = new AbortController();

    const ambilProfil = async () => {
      try {
        const response = await fetch(
          `${API_URL}/api/pelanggan/profil`,
          {
            method: 'GET',
            credentials: 'include',
            signal: controller.signal
          }
        );

        if (response.status === 401 || response.status === 403) {
          navigate('/masuk', { replace: true });
          return;
        }

        if (!response.ok) {
          throw new Error('Gagal mengambil data pelanggan.');
        }

        const data = await response.json();

        if (!controller.signal.aborted) {
          setUser(data.pelanggan);
          setError('');
        }
      } catch (err) {
        if (err.name !== 'AbortError' && !controller.signal.aborted) {
          setError(
            'Tidak dapat memeriksa sesi pelanggan. Periksa koneksi backend.'
          );
        }
      } finally {
        if (!controller.signal.aborted) {
          setLoading(false);
        }
      }
    };

    ambilProfil();

    return () => controller.abort();
  }, [navigate]);

  // Menghapus sesi pelanggan melalui backend.
  const handleKeluar = async () => {
    if (loggingOut) return;

    setLoggingOut(true);
    setError('');

    try {
      const response = await fetch(
        `${API_URL}/api/pelanggan/keluar`,
        {
          method: 'POST',
          credentials: 'include'
        }
      );

      if (!response.ok) {
        throw new Error('Logout gagal. Silakan coba lagi.');
      }

      setUser(null);
      window.dispatchEvent(new Event('pelanggan-auth-changed'));
      navigate('/masuk', { replace: true });
    } catch (err) {
      setError(err.message || 'Gagal keluar dari akun.');
    } finally {
      setLoggingOut(false);
    }
  };

  // Jangan tampilkan profil palsu saat sesi diperiksa.
  if (loading) {
    return (
      <section className="auth-pelanggan-section">
        <div className="container">
          <div className="auth-container-box">
            <div className="auth-card-white">
              <h1 className="auth-title">Memeriksa akun...</h1>
              <p className="auth-subtitle">
                Memuat data pelanggan dari server.
              </p>
            </div>
          </div>
        </div>
      </section>
    );
  }

  // Jika backend gagal atau sesi tidak tersedia,
  // jangan tampilkan data pelanggan.
  if (!user) {
    return (
      <section className="auth-pelanggan-section">
        <div className="container">
          <div className="auth-container-box">
            <div className="auth-card-white">
              <h1 className="auth-title">Akun belum tersedia</h1>

              <p className="auth-subtitle">
                {error || 'Silakan masuk untuk melihat akun.'}
              </p>

              <button
                type="button"
                className="btn-auth-submit"
                onClick={() => window.location.reload()}
              >
                Coba lagi
              </button>

              <div className="auth-switch-text">
                <Link to="/masuk" className="auth-switch-btn">
                  Kembali ke Login
                </Link>
              </div>
            </div>
          </div>
        </div>
      </section>
    );
  }

  const inisial = user.nama
    ? user.nama.charAt(0).toUpperCase()
    : 'P';

  return (
    <section className="auth-pelanggan-section">
      <div className="container">
        <div className="auth-container-box">

          {/* Foto produk sebelah kiri */}
          <div className="auth-image-wrapper">
            <img
              src="/images/products/rengginang-all-hero.jpg"
              alt="Rengginang Putra Jaya"
              className="auth-hero-img"
            />
          </div>

          {/* Kartu akun pelanggan */}
          <div className="auth-card-white">
            <div className="auth-badge-success">
              <span className="auth-badge-dot">●</span>
              <span>LOGIN BERHASIL</span>
            </div>

            <h1 className="auth-title">Akun Saya</h1>

            <p className="auth-subtitle">
              Kamu sudah masuk sebagai pelanggan.
            </p>

            {error && (
              <div
                className="auth-alert-message auth-alert-error"
                role="alert"
              >
                {error}
              </div>
            )}

            {/* Data pelanggan dari database MySQL */}
            <div className="akun-profile-card">
              <div className="akun-profile-info">
                <div className="akun-avatar">{inisial}</div>

                <div>
                  <div className="akun-nama">
                    {user.nama}
                  </div>

                  <div className="akun-email">
                    {user.email}
                  </div>
                </div>
              </div>

              <span className="akun-role-pill">
                Pelanggan
              </span>
            </div>

            {/* Tombol belanja */}
            <Link
              to="/katalog"
              className="btn-auth-submit"
              style={{
                display: 'block',
                textAlign: 'center'
              }}
            >
              Lanjut Belanja
            </Link>

            <button
              type="button"
              onClick={handleKeluar}
              disabled={loggingOut}
              className="btn-akun-secondary"
            >
              {loggingOut ? 'Memproses...' : 'Keluar Akun'}
            </button>

            <p className="akun-footer-note">
              Belanja juga bisa dilakukan tanpa login.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}