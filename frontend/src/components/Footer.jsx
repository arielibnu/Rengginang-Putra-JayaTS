import React from 'react';
import { Link } from 'react-router-dom';
import { BUSINESS_CONFIG } from '../config/business';

export default function Footer({ onOpenContact }) {
  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <footer className="footer-dark" role="contentinfo">
      <div className="container">
        <div className="footer-inner">
          {/* Sisi Kiri: Identitas Usaha */}
          <div>
            <div className="footer-brand-title">{BUSINESS_CONFIG.name}</div>
            <div className="footer-brand-desc">Camilan tradisional khas Situbondo</div>
          </div>

          {/* Sisi Kanan: Navigasi Ringkas */}
          <ul className="footer-nav-row">
            <li>
              <Link to="/" className="footer-nav-item" onClick={scrollToTop}>
                Beranda
              </Link>
            </li>
            <li>
              <Link to={BUSINESS_CONFIG.routes.catalog} className="footer-nav-item">
                Produk
              </Link>
            </li>
            <li>
              <button 
                type="button" 
                className="footer-nav-item"
                onClick={onOpenContact}
              >
                Kontak
              </button>
            </li>
          </ul>
        </div>
      </div>
    </footer>
  );
}
