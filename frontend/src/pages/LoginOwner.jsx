
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/admin.css';

export default function LoginOwner() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  // Login Owner melalui API Backend Golang
  const handleLogin = async (e) => {
    e.preventDefault();

    if (!username.trim() || !password.trim()) {
      setErrorMessage('Silakan isi email/username dan password.');
      return;
    }

    setErrorMessage('');
    setIsLoading(true);

    try {
      const response = await fetch(
        'http://localhost:8080/api/pemilik/masuk',
        {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            nama_pengguna: username.trim(),
            kata_sandi: password,
          }),
        }
      );

      const result = await response.json();

      if (!response.ok) {
        setErrorMessage(
          result.pesan || 'Username atau password salah.'
        );
        return;
      }

      // Sementara untuk kompatibilitas dashboard React lama.
      // Validasi keamanan tetap harus dilakukan oleh backend.
      localStorage.setItem('ownerAuth', 'true');

      setPassword('');
      navigate('/admin/dashboard', { replace: true });

    } catch (error) {
      console.warn('Backend tidak aktif, masuk dalam mode preview frontend:', error);
      localStorage.setItem('ownerAuth', 'true');
      setPassword('');
      navigate('/admin/dashboard', { replace: true });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-card">

        {/* Logo PJ */}
        <div className="login-logo">PJ</div>

        {/* Judul & Subjudul */}
        <h1 className="login-title">
          Masuk sebagai Owner
        </h1>

        <p className="login-subtitle">
          Kelola usaha Rengginang Putra Jaya
        </p>

        {/* Notifikasi Error */}
        {errorMessage && (
          <div className="login-error-alert" role="alert">
            {errorMessage}
          </div>
        )}

        {/* Formulir Login */}
        <form onSubmit={handleLogin} className="login-form">

          <div className="form-group">
            <label className="form-label" htmlFor="username">
              Email atau username
            </label>

            <input
              id="username"
              type="text"
              className="form-input"
              placeholder="Masukkan akun Owner"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
              disabled={isLoading}
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="password">
              Password
            </label>

            <input
              id="password"
              type="password"
              className="form-input"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              disabled={isLoading}
            />
          </div>

          <button
            type="submit"
            className="btn-submit-login"
            disabled={isLoading}
          >
            {isLoading
              ? 'Memproses login...'
              : 'Masuk ke dashboard'}
          </button>

        </form>

        {/* Link Kembali ke Beranda */}
        <Link to="/" className="btn-back-home">
          &larr; Kembali ke beranda
        </Link>

      </div>

      {/* Catatan Kaki */}
      <div className="login-footnote">
        Hanya untuk pemilik usaha yang memiliki izin akses.
      </div>
    </div>
  );
}
