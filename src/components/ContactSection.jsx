import React from 'react';
import { MessageSquare, AlertCircle, MapPin, CheckCircle2 } from 'lucide-react';
import { BUSINESS_CONFIG } from '../config/business';

export default function ContactSection() {
  const isWhatsAppConfigured = Boolean(
    BUSINESS_CONFIG.whatsappNumber && BUSINESS_CONFIG.whatsappNumber.trim() !== ''
  );

  const whatsappUrl = isWhatsAppConfigured
    ? `https://wa.me/${BUSINESS_CONFIG.whatsappNumber.replace(/[^0-9]/g, '')}?text=${encodeURIComponent(
        BUSINESS_CONFIG.whatsappDefaultMessage
      )}`
    : null;

  return (
    <section id="kontak" className="contact-section">
      <div className="container">
        <div className="contact-box">
          <div className="contact-icon-wrapper">
            <MessageSquare size={30} />
          </div>

          <h2 className="contact-title">Hubungi Kami</h2>
          <p className="contact-subtitle">
            Punya pertanyaan mengenai ketersediaan varian rasa atau informasi pemesanan Rengginang Putra Jaya?
            Kami siap melayani Anda melalui WhatsApp.
          </p>

          {/* Status WhatsApp Config */}
          {!isWhatsAppConfigured ? (
            <div className="contact-status-box" role="alert">
              <AlertCircle size={20} />
              <div>
                <div className="contact-status-title">Kontak WhatsApp Belum Dikonfigurasi</div>
                <div className="contact-status-desc">
                  Nomor WhatsApp resmi pemilik usaha belum dimasukkan pada file konfigurasi (<code>src/config/business.js</code>).
                  Tombol chat dinonaktifkan untuk mencegah tautan yang tidak valid.
                </div>
              </div>
            </div>
          ) : (
            <div className="contact-status-box" style={{ backgroundColor: '#EAF7EE', borderColor: '#BDE4C7' }}>
              <CheckCircle2 size={20} style={{ color: '#278A48' }} />
              <div>
                <div className="contact-status-title" style={{ color: '#1B5B30' }}>Kontak Resmi Terhubung</div>
                <div className="contact-status-desc" style={{ color: '#278A48' }}>
                  Nomor WhatsApp terverifikasi: +{BUSINESS_CONFIG.whatsappNumber}
                </div>
              </div>
            </div>
          )}

          <div className="contact-actions">
            {isWhatsAppConfigured ? (
              <a 
                href={whatsappUrl} 
                target="_blank" 
                rel="noopener noreferrer" 
                className="btn btn-whatsapp"
              >
                <MessageSquare size={20} />
                <span>Chat WhatsApp Sekarang</span>
              </a>
            ) : (
              <button 
                type="button" 
                className="btn btn-whatsapp" 
                disabled 
                title="Nomor WhatsApp belum dikonfigurasi pada src/config/business.js"
              >
                <MessageSquare size={20} />
                <span>Chat WhatsApp (Belum Dikonfigurasi)</span>
              </button>
            )}
          </div>

          <div className="contact-meta">
            <MapPin size={15} />
            <span>{BUSINESS_CONFIG.location}, {BUSINESS_CONFIG.country}</span>
          </div>
        </div>
      </div>
    </section>
  );
}

