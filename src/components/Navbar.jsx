
import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Menu, X } from 'lucide-react';
import { BUSINESS_CONFIG } from '../config/business';

const API_URL = 'http://localhost:8080';

export default function Navbar({ onOpenContact }) {
  const [isMobileOpen, setIsMobileOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const location = useLocation();

  const closeMenu = () => setIsMobileOpen(false);

  // Periksa sesi pelanggan melalui backend Go.
  useEffect(() => {
    const controller = new AbortController();

    const checkLoginStatus = async () => {
      try {
        const response = await fetch(
          `${API_URL}/api/pelanggan/profil`,
          {
            method: 'GET',
            credentials: 'include',
            signal: controller.signal
          }
        );

        if (!controller.signal.aborted) {
          setIsLoggedIn(response.ok);
        }
      } catch (error) {
        if (
          error.name !== 'AbortError' &&
          !controller.signal.aborted
        ) {
          setIsLoggedIn(false);
        }
      }
    };

    checkLoginStatus();

    window.addEventListener(
      'pelanggan-auth-changed',
      checkLoginStatus
    );

    return () => {
      controller.abort();
      window.removeEventListener(
        'pelanggan-auth-changed',
        checkLoginStatus
      );
    };
  }, [location.pathname]);

  const scrollToAbout = () => {
    closeMenu();

    const el = document.getElementById('tentang');

    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <nav
      className="navbar"
      role="navigation"
      aria-label="Navigasi Utama"
    >
      <div className="container navbar-container">

        {/* Brand Kiri */}
        <Link
          to="/"
          className="navbar-brand"
          onClick={closeMenu}
        >
          <span className="brand-icon">PJ</span>
          <span>{BUSINESS_CONFIG.name}</span>
        </Link>

        {/* Menu Desktop Kanan */}
        <div className="navbar-menu-group">
          <ul className="navbar-links">
            <li>
              <Link
                to="/"
                className={`navbar-link ${
                  location.pathname === '/' ? 'active' : ''
                }`}
              >
                Beranda
              </Link>
            </li>

            <li>
              <Link
                to={BUSINESS_CONFIG.routes.catalog}
                className={`navbar-link ${
                  location.pathname === BUSINESS_CONFIG.routes.catalog
                    ? 'active'
                    : ''
                }`}
              >
                Produk
              </Link>
            </li>

            <li>
              <button
                type="button"
                className="navbar-link"
                onClick={scrollToAbout}
              >
                Tentang
              </button>
            </li>

            <li>
              <button
                type="button"
                className="navbar-link"
                onClick={onOpenContact}
              >
                Kontak
              </button>
            </li>
          </ul>

          {/* Login Pelanggan / Akun Saya */}
          {isLoggedIn ? (
            <Link
              to="/akun"
              className="btn-login-outline"
            >
              Akun Saya
            </Link>
          ) : (
            <Link
              to="/masuk"
              className="btn-login-outline"
            >
              Masuk
            </Link>
          )}
        </div>

        {/* Mobile Hamburger */}
        <button
          type="button"
          className="mobile-toggle"
          onClick={() => setIsMobileOpen(!isMobileOpen)}
          aria-label="Menu navigasi"
          aria-expanded={isMobileOpen}
        >
          {isMobileOpen ? (
            <X size={24} />
          ) : (
            <Menu size={24} />
          )}
        </button>
      </div>

      {/* Mobile Drawer */}
      {isMobileOpen && (
        <div className="mobile-menu-drawer">
          <Link
            to="/"
            className="mobile-nav-item active"
            onClick={closeMenu}
          >
            Beranda
          </Link>

          <Link
            to={BUSINESS_CONFIG.routes.catalog}
            className="mobile-nav-item"
            onClick={closeMenu}
          >
            Produk
          </Link>

          <button
            type="button"
            className="mobile-nav-item"
            style={{
              textAlign: 'left',
              width: '100%'
            }}
            onClick={scrollToAbout}
          >
            Tentang
          </button>

          <button
            type="button"
            className="mobile-nav-item"
            style={{
              textAlign: 'left',
              width: '100%'
            }}
            onClick={() => {
              closeMenu();
              onOpenContact();
            }}
          >
            Kontak
          </button>

          {isLoggedIn ? (
            <Link
              to="/akun"
              className="btn-login-outline"
              style={{
                marginTop: '0.5rem',
                width: '100%',
                textAlign: 'center'
              }}
              onClick={closeMenu}
            >
              Akun Saya
            </Link>
          ) : (
            <Link
              to="/masuk"
              className="btn-login-outline"
              style={{
                marginTop: '0.5rem',
                width: '100%',
                textAlign: 'center'
              }}
              onClick={closeMenu}
            >
              Masuk
            </Link>
          )}
        </div>
      )}
    </nav>
  );
}
