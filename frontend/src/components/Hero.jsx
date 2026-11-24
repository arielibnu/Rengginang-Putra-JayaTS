import React from 'react';
import { Link } from 'react-router-dom';
import { BUSINESS_CONFIG } from '../config/business';

export default function Hero({ onOpenContact }) {
  return (
    <section className="hero-wrapper">
      <div className="container">
        <div className="hero-banner-card">
          {/* Kolom Teks Kiri */}
          <div className="hero-text-side">
            <div className="hero-tag">{BUSINESS_CONFIG.taglineText}</div>
            <h1 className="hero-headline">
              Rengginang<br />
              tradisional,<br />
              rasa istimewa.
            </h1>
            <p className="hero-subtext">
              {BUSINESS_CONFIG.heroDescription}
            </p>
            <div className="hero-btn-row">
              <Link to={BUSINESS_CONFIG.routes.catalog} className="btn-primary-solid">
                Lihat Produk
              </Link>
              <button 
                type="button" 
                className="btn-white-solid" 
                onClick={onOpenContact}
              >
                Hubungi Kami
              </button>
            </div>
          </div>

          {/* Kolom Foto Kanan */}
          <div className="hero-image-side">
            <img 
              src="/images/products/rengginang-all-hero.jpg" 
              alt="Rengginang Putra Jaya Empat Varian Rasa" 
              className="hero-img"
            />
          </div>
        </div>
      </div>
    </section>
  );
}
