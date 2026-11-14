import React from 'react';
import { UserCheck, MapPin, Package, Award } from 'lucide-react';
import { BUSINESS_CONFIG } from '../config/business';

export default function AboutSection() {
  return (
    <section id="tentang-kami" className="about-section">
      <div className="container">
        <div className="section-header">
          <span className="section-tag">Profil Usaha</span>
          <h2 className="section-title">Tentang {BUSINESS_CONFIG.name}</h2>
          <p className="section-description">
            Usaha kuliner lokal yang berfokus menyajikan rengginang berkualitas dengan cita rasa tradisional khas Situbondo.
          </p>
        </div>

        <div className="about-content-wrapper">
          {/* Deskripsi Usaha */}
          <div className="about-text-column">
            {BUSINESS_CONFIG.aboutDescription.map((paragraph, index) => (
              <p key={index} className="about-paragraph">
                {paragraph}
              </p>
            ))}
          </div>

          {/* Kartu Informasi Resmi (Sesuai Label Kemasan) */}
          <div className="about-cards-grid">
            <div className="about-info-card">
              <div className="about-card-icon">
                <UserCheck size={22} />
              </div>
              <div>
                <div className="about-card-label">Produsen Resmi</div>
                <div className="about-card-val">{BUSINESS_CONFIG.producer}</div>
                <div className="about-card-sub">Diproduksi secara terstandar di Situbondo</div>
              </div>
            </div>

            <div className="about-info-card">
              <div className="about-card-icon">
                <MapPin size={22} />
              </div>
              <div>
                <div className="about-card-label">Lokasi Produksi</div>
                <div className="about-card-val">{BUSINESS_CONFIG.location}</div>
                <div className="about-card-sub">Daerah pesisir penghasil aneka sari laut</div>
              </div>
            </div>

            <div className="about-info-card">
              <div className="about-card-icon">
                <Package size={22} />
              </div>
              <div>
                <div className="about-card-label">Kemasan & Bobot</div>
                <div className="about-card-val">{BUSINESS_CONFIG.netWeight}</div>
                <div className="about-card-sub">Kemasan higienis dan rapat</div>
              </div>
            </div>

            <div className="about-info-card">
              <div className="about-card-icon">
                <Award size={22} />
              </div>
              <div>
                <div className="about-card-label">Karakteristik Produk</div>
                <div className="about-card-val">{BUSINESS_CONFIG.featureTag}</div>
                <div className="about-card-sub">Langsung digoreng renyah mekar</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

