import React from 'react';
import { X, Lock, ShieldCheck } from 'lucide-react';

export default function LoginInfoModal({ isOpen, onClose }) {
  if (!isOpen) return null;

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
            <Lock size={26} />
          </div>
          <h3 style={{ fontSize: '1.35rem', fontWeight: 800, color: '#43281C', marginBottom: '0.25rem' }}>
            Akses Masuk Akun
          </h3>
          <p style={{ fontSize: '0.9rem', color: '#66564C' }}>
            Fitur Login Owner & Pelanggan (PPL-01 & PPL-02)
          </p>
        </div>

        <div style={{
          backgroundColor: '#F8F3EA',
          borderRadius: '12px',
          padding: '1rem',
          marginBottom: '1.5rem',
          fontSize: '0.875rem',
          color: '#66564C',
          lineHeight: 1.55
        }}>
          Fitur autentikasi dan akun pelanggan sedang dalam tahap pengerjaan Sprint 1 oleh anggota tim. Pengunjung dapat langsung menikmati katalog produk tanpa perlu menginstal aplikasi atau membuat akun.
        </div>

        <button 
          type="button" 
          className="btn-primary-solid" 
          style={{ width: '100%', justifyContent: 'center' }}
          onClick={onClose}
        >
          Lanjut Jelajahi Beranda
        </button>
      </div>
    </div>
  );
}

