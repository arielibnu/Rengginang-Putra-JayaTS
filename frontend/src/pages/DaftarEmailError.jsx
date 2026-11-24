import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/home.css';

export default function DaftarEmailError() {
  const navigate = useNavigate();

  // state input form
  const [nama, setNama] = useState('');
  const [email, setEmail] = useState('nama@email.com');
  const [password, setPassword] = useState('');
  const [konfirmasiPassword, setKonfirmasiPassword] = useState('');

  // tombol Ubah Email akan fokus kembali ke input email atau kembali ke form daftar
  const handleUbahEmail = (e) => {
    e.preventDefault();
    navigate('/daftar');
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

          {/* Form daftar akun di sebelah kanan dengan status error email */}
          <div className="auth-card-white">
            <span className="auth-badge-tag">AKUN PELANGGAN</span>

            <h1 className="auth-title">Buat akun baru</h1>
            <p className="auth-subtitle">
              Daftar dengan email untuk belanja lebih mudah.
            </p>

            <form onSubmit={handleUbahEmail} className="auth-form">
              <div className="auth-field">
                <label className="auth-label">Nama lengkap</label>
                <input
                  type="text"
                  className="auth-input"
                  placeholder="Masukkan nama lengkap"
                  value={nama}
                  onChange={(e) => setNama(e.target.value)}
                />
              </div>

              <div className="auth-field">
                <label className="auth-label">Email</label>
                <input
                  type="email"
                  className="auth-input has-error"
                  placeholder="nama@email.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
                <div className="auth-input-error-msg">
                  Email sudah terdaftar. Coba masuk.
                </div>
              </div>

              <div className="auth-field">
                <label className="auth-label">Password</label>
                <input
                  type="password"
                  className="auth-input"
                  placeholder="Minimal 8 karakter"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>

              <div className="auth-field">
                <label className="auth-label">Konfirmasi password</label>
                <input
                  type="password"
                  className="auth-input"
                  placeholder="Ulangi password"
                  value={konfirmasiPassword}
                  onChange={(e) => setKonfirmasiPassword(e.target.value)}
                />
              </div>

              <button type="submit" className="btn-auth-submit">
                Ubah Email
              </button>
            </form>

            <div className="auth-switch-text">
              Sudah punya akun?{' '}
              <Link to="/login" className="auth-switch-btn">
                Masuk
              </Link>
            </div>

            <hr className="auth-divider" />

            <button
              type="button"
              onClick={() => navigate('/katalog')}
              className="auth-guest-link"
            >
              Belanja tanpa akun
            </button>

            {/* Teks mengarah ke login Owner */}
            <div className="owner-login-hint">
              Pemilik usaha?{' '}
              <Link to="/login-owner">Masuk ke Dashboard Owner</Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
