import React from 'react';
import { BUSINESS_CONFIG } from '../config/business';

export default function WhyChooseUs() {
  return (
    <section id="tentang" className="why-section">
      <div className="container">
        <div className="why-banner-card">
          <h2 className="why-title">Mengapa pilih Putra Jaya?</h2>

          <div className="why-points-grid">
            {BUSINESS_CONFIG.features.map((feat) => (
              <div key={feat.number} className="why-point-item">
                <span className="why-number-circle">{feat.number}</span>
                <div>
                  <h3 className="why-item-title">{feat.title}</h3>
                  <p className="why-item-desc">{feat.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

