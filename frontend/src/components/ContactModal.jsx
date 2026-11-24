import React from 'react';
import { X, MessageSquare, AlertCircle, CheckCircle2, MapPin } from 'lucide-react';
import { BUSINESS_CONFIG } from '../config/business';

export default function ContactModal({ isOpen, onClose }) {
  if (!isOpen) return null;

  const isConfigured = Boolean(
    BUSINESS_CONFIG.whatsappNumber && BUSINESS_CONFIG.whatsappNumber.trim() !== ''
  );

  const whatsappUrl = isConfigured
    ? `https://wa.me/${BUSINESS_CONFIG.whatsappNumber.replace(/[^0-9]/g, '')}?text=${encodeURIComponent(
        BUSINESS_CONFIG.whatsappDefaultMessage
      )}`
    : null;

  return (
    <div className="modal-overlay" onClick={onClose} role="dialog" aria-modal="true">
      <div className="modal-dialog" onClick={(e) => e.stopPropagation()}>
        <button 
          type="button" 
          className="modal-close-btn" 
          onClick={onClose}
          aria-label="Tutup modal"
        >
          <X size={20} />
        </button>

        <div style={{ textAlign: 'center', marginBottom: '1.25rem' }}>
          <div style={{
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: '52px',
            height: '52px',
            borderRadius: '50%',
            backgroundColor: '#F8F3EA',
            color: '#43281C',
            marginBottom: '0.75rem'
          }}>
            <MessageSquare size={26} />
          </div>
          <h3 style={{ fontSize: '1.35rem', fontWeight: 800, color: '#43281C', marginBottom: '0.25rem' }}>
            Hubungi Rengginang Putra Jaya
          </h3>
          <p style={{ fontSize: '0.9rem', color: '#66564C' }}>
            Produksi Hj. Arifah — {BUSINESS_CONFIG.location}
          </p>
        </div>

        {!isConfigured ? (
          <div style={{
            backgroundColor: '#FFF8F0',
            border: '1px solid #F5E0CE',
            borderRadius: '12px',
            padding: '1rem',
            marginBottom: '1.5rem',
            display: 'flex',
            gap: '0.75rem',
            alignItems: 'flex-start'
          }}>
            <AlertCircle size={20} style={{ color: '#D97724', flexShrink: 0, marginTop: '2px' }} />
            <div>
              <div style={{ fontSize: '0.875rem', fontWeight: 700, color: '#8F4A12', marginBottom: '0.2rem' }}>
                Kontak WhatsApp Belum Dikonfigurasi
              </div>
              <div style={{ fontSize: '0.8125rem', color: '#6B401B', lineHeight: 1.45 }}>
                Nomor WhatsApp resmi pemilik usaha belum dimasukkan pada file konfigurasi (<code>src/config/business.js</code>).
                Tombol chat dinonaktifkan secara aman.
              </div>
            </div>
          </div>
        ) : (
          <div style={{
            backgroundColor: '#EAF7EE',
            border: '1px solid #BDE4C7',
            borderRadius: '12px',
            padding: '1rem',
            marginBottom: '1.5rem',
            display: 'flex',
            gap: '0.75rem',
            alignItems: 'center'
          }}>
            <CheckCircle2 size={20} style={{ color: '#278A48' }} />
            <div style={{ fontSize: '0.875rem', color: '#1B5B30', fontWeight: 600 }}>
              WhatsApp Resmi: +{BUSINESS_CONFIG.whatsappNumber}
            </div>
          </div>
        )}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          {isConfigured ? (
            <a 
              href={whatsappUrl} 
              target="_blank" 
              rel="noopener noreferrer" 
              className="btn-primary-solid"
              style={{ justifyContent: 'center' }}
            >
              <MessageSquare size={18} style={{ marginRight: '6px' }} />
              Chat WhatsApp Sekarang
            </a>
          ) : (
            <button 
              type="button" 
              className="btn-primary-solid" 
              disabled 
              style={{ backgroundColor: '#B0A197', cursor: 'not-allowed', justifyContent: 'center' }}
            >
              Chat WhatsApp (Belum Dikonfigurasi)
            </button>
          )}

          <button 
            type="button" 
            className="btn-white-solid" 
            style={{ justifyContent: 'center' }}
            onClick={onClose}
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  );
}

