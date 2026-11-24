import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft, ShoppingBag } from 'lucide-react';
import { BUSINESS_CONFIG } from '../config/business';

export default function CatalogPlaceholder() {
  return (
    <div style={{
      minHeight: '70vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '3rem 1.5rem',
      backgroundColor: '#FFFFFF'
    }}>
      <div style={{
        backgroundColor: '#FAF6EF',
        border: '1px solid #ECE4D9',
        borderRadius: '24px',
        padding: '3rem 2rem',
        textAlign: 'center',
        maxWidth: '560px',
        boxShadow: '0 8px 24px rgba(67, 40, 28, 0.05)'
      }}>
        <div style={{
          display: 'inline-flex',
          padding: '1rem',
          backgroundColor: '#43281C',
          borderRadius: '16px',
          color: '#FFFFFF',
          marginBottom: '1.25rem'
        }}>
          <ShoppingBag size={32} />
        </div>

        <h1 style={{ fontSize: '1.65rem', fontWeight: 800, color: '#43281C', marginBottom: '0.75rem' }}>
          Katalog Rengginang Putra Jaya
        </h1>
        <p style={{ fontSize: '0.975rem', color: '#66564C', lineHeight: 1.6, marginBottom: '2rem' }}>
          Navigasi katalog produk telah terhubung secara fungsional. Fitur katalog penuh (PPL-08: Terasi, Udang, Bawang, Cumi) sedang dalam tahap pengerjaan Sprint 1 oleh rekan tim.
        </p>

        <div style={{ display: 'flex', justifyContent: 'center', gap: '0.875rem', flexWrap: 'wrap' }}>
          <Link to="/" className="btn-primary-solid">
            <ArrowLeft size={18} style={{ marginRight: '6px' }} />
            Kembali ke Beranda
          </Link>
        </div>
      </div>
    </div>
  );
}
