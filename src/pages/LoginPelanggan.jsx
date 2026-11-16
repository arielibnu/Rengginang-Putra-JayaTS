
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/home.css';

const API_URL = 'http://localhost:8080';

export default function LoginPelanggan() {
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');

    if (!email.trim() || !password) {
      setError('Email dan password wajib diisi.');
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(
        `${API_URL}/api/pelanggan/masuk`,
        {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            email: email.trim(),
            password
          })
        }
      );

      const data = await response.json();

      if (!response.ok) {
        throw new Error(
          data.pesan || 'Email atau password salah.'
        );
      }

      // Backend sudah memvalidasi password dan membuat cookie sesi.
      // Jangan menyimpan password atau token di localStorage.
      window.dispatchEvent(new Event('pelanggan-auth-changed'));

      navigate('/akun', { replace: true });

    } catch (err) {
      setError(
        err.message || 'Tidak dapat terhubung ke server.'
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="auth-pelanggan-section">
      <div className="container">
        <div className="auth-container-box">

          <div className="auth-image-wrapper">
            <img
              src="/images/products/rengginang-all-hero.jpg"
              alt="Rengginang Putra Jaya"
              className="auth-hero-img"
            />
          </div>

          <div className="auth-card-white">
            <h1 className="auth-title">
              Selamat datang kembali
            </h1>

            <p className="auth-subtitle">
              Masuk untuk pengalaman belanja yang lebih mudah.
            </p>

            {error && (
              <div
                className="auth-alert-message auth-alert-error"
                role="alert"
              >
                {error}
              </div>
            )}

            <form
              onSubmit={handleLogin}
              className="auth-form"
            >
              <div className="auth-field">
                <label className="auth-label">
                  Email
                </label>

                <input
                  type="email"
                  className="auth-input"
                  placeholder="nama@email.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoComplete="email"
                  required
                />
              </div>

              <div className="auth-field">
                <label className="auth-label">
                  Password
                </label>

                <input
                  type="password"
                  className="auth-input"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  required
                />
              </div>

              <button
                type="submit"
                className="btn-auth-submit"
                disabled={loading}
              >
                {loading ? 'Memproses...' : 'Masuk'}
              </button>
            </form>

            <div className="auth-switch-text">
              Belum punya akun?{' '}
              <Link
                to="/daftar"
                className="auth-switch-btn"
              >
                Daftar
              </Link>
            </div>

            <hr className="auth-divider" />

            <button
              type="button"
              className="btn-guest-checkout"
              onClick={() => navigate('/katalog')}
            >
              Lanjut belanja tanpa akun
            </button>

            <div className="owner-login-hint">
              Pemilik usaha?{' '}
              <Link to="/login">
                Masuk ke Dashboard Owner
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}