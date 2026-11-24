import React from 'react';
import { Link } from 'react-router-dom';
import { BUSINESS_CONFIG } from '../config/business';

export default function FlavorPreview() {
  return (
    <section className="flavor-section">
      <div className="container">
        {/* Header Baris */}
        <div className="flavor-header-row">
          <h2 className="flavor-main-title">Varian rengginang</h2>
          <Link to={BUSINESS_CONFIG.routes.catalog} className="flavor-view-all">
            Lihat semua rasa
          </Link>
        </div>

        {/* Grid 4 Produk */}
        <div className="flavor-cards-grid">
          {BUSINESS_CONFIG.flavors.map((item) => (
            <div key={item.id} className="flavor-product-card">
              <div className="card-img-container">
                <img 
                  src={item.image} 
                  alt={item.name} 
                  className="card-product-img"
                  loading="lazy"
                />
              </div>

              <div className="card-body">
                <h3 className="card-item-title">{item.name}</h3>
                <div className="card-item-weight">{item.weight}</div>

                <div className="card-bottom-row">
                  <span className="card-item-price">{item.price}</span>
                  <Link to={`${BUSINESS_CONFIG.routes.catalog}?id=${item.id}`} className="card-detail-link">
                    Lihat detail
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
