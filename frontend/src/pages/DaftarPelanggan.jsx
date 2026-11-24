
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/home.css';

const API_URL = 'http://localhost:8080';

export default function DaftarPelanggan() {
  const navigate = useNavigate();

  const [nama, setNama] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [konfirmasiPassword, setKonfirmasiPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleDaftar = async (e) => {
    e.preventDefault();
    setError('');

    if (
      !nama.trim() ||
      !email.trim() ||
      !password ||
      !konfirmasiPassword
    ) {
      setError('Semua kolom wajib diisi.');
      return;
    }

    if (password !== konfirmasiPassword) {
      setError('Konfirmasi password tidak cocok.');
      return;
    }

    if (
      new TextEncoder().encode(password).length < 8 ||
      new TextEncoder().encode(password).length > 72
    ) {
      setError('Password harus berisi 8 sampai 72 byte.');
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(
        `${API_URL}/api/pelanggan/daftar`,
        {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            nama: nama.trim(),
            email: email.trim().toLowerCase(),
            password
          })
        }
      );

      const data = await response.json();

      if (response.status === 409) {
        navigate('/daftar-error');
        return;
      }

      if (!response.ok) {
        throw new Error(
          data.pesan || 'Pendaftaran gagal. Silakan coba lagi.'
        );
      }

      // Registrasi berhasil, pelanggan perlu login
      // untuk mendapatkan cookie sesi dari backend.
      navigate('/masuk', {
        replace: true,
        state: {
          pesan: 'Pendaftaran berhasil. Silakan masuk.',
          email: email.trim().toLowerCase()
        }
      });
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
          {/* Foto produk sebelah kiri */}
          <div className="auth-image-wrapper">
            <img
              src="/images/products/rengginang-all-hero.jpg"
              alt="Rengginang Putra Jaya"
              className="auth-hero-img"
            />
          </div>

          {/* Form daftar akun di sebelah kanan */}
          <div className="auth-card-white">
            <span className="auth-badge-tag">
              AKUN PELANGGAN
            </span>

            <h1 className="auth-title">Buat akun baru</h1>

            <p className="auth-subtitle">
              Daftar dengan email untuk belanja lebih mudah.
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
              onSubmit={handleDaftar}
              className="auth-form"
            >
              <div className="auth-field">
                <label className="auth-label">
                  Nama lengkap
                </label>

                <input
                  type="text"
                  className="auth-input"
                  placeholder="Masukkan nama lengkap"
                  value={nama}
                  onChange={(e) => setNama(e.target.value)}
                  autoComplete="name"
                  maxLength={150}
                  required
                />
              </div>

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
                  maxLength={191}
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
                  placeholder="Minimal 8 karakter"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="new-password"
                  minLength={8}
                  required
                />
              </div>

              <div className="auth-field">
                <label className="auth-label">
                  Konfirmasi password
                </label>

                <input
                  type="password"
                  className="auth-input"
                  placeholder="Ulangi password"
                  value={konfirmasiPassword}
                  onChange={(e) => setKonfirmasiPassword(e.target.value)}
                  autoComplete="new-password"
                  required
                />
              </div>

              <button
                type="submit"
                className="btn-auth-submit"
                disabled={loading}
              >
                {loading ? 'Memproses...' : 'Daftar Sekarang'}
              </button>
            </form>

            <div className="auth-switch-text">
              Sudah punya akun?{' '}
              <Link to="/masuk" className="auth-switch-btn">
                Masuk
              </Link>
            </div>

            <hr className="auth-divider" />

            <button
              type="button"
              className="btn-guest-checkout"
              onClick={() => navigate('/katalog')}
            >
              Belanja tanpa akun
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