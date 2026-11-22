import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate, Link } from 'react-router-dom';

const API_URL = 'http://localhost:8080';

export default function KonfirmasiPesanan() {
  const location = useLocation();
  const navigate = useNavigate();

  // State data pesanan hasil respons backend
  const [pesanan, setPesanan] = useState(null);
  const [nomorPesanan, setNomorPesanan] = useState('');
  const [statusPesanan, setStatusPesanan] = useState('menunggu_pembayaran');
  const [statusLabel, setStatusLabel] = useState('Menunggu pembayaran');
  const [loading, setLoading] = useState(true);
  const [isDemo, setIsDemo] = useState(false);
  const [tampilkanDetail, setTampilkanDetail] = useState(false);

  useEffect(() => {
    // 1. Cek data yang dikirim melalui navigasi state dari Checkout PPL-17
    const stateOrder = location.state?.pesananData || location.state?.order;
    const queryParams = new URLSearchParams(location.search);
    const queryOrderId = queryParams.get('order_id') || queryParams.get('nomor_pesanan') || queryParams.get('id');

    if (stateOrder && stateOrder.nomor_pesanan) {
      setPesanan(stateOrder);
      setNomorPesanan(stateOrder.nomor_pesanan);
      setStatusPesanan(stateOrder.status_pesanan || 'menunggu_pembayaran');
      setStatusLabel(stateOrder.status_pesanan_label || 'Menunggu pembayaran');
      setIsDemo(false);
      setLoading(false);
      return;
    }

    // 2. Jika ada query parameter order_id, ambil data asli dari backend API
    const targetId = queryOrderId || 'PJ-DEMO-0001';
    if (targetId === 'PJ-DEMO-0001' && !queryOrderId) {
      setIsDemo(true);
    }

    const fetchDetailPesanan = async () => {
      try {
        const res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetId)}`);
        if (res.ok) {
          const json = await res.json();
          if (json.data) {
            setPesanan(json.data);
            setNomorPesanan(json.data.nomor_pesanan);
            setStatusPesanan(json.data.status_pesanan || 'menunggu_pembayaran');
            setStatusLabel(json.data.status_pesanan_label || 'Menunggu pembayaran');
            setIsDemo(json.data.nomor_pesanan === 'PJ-DEMO-0001');
          } else {
            setNomorPesanan(targetId);
          }
        } else {
          // Fallback lokal jika data demo
          setNomorPesanan(targetId);
        }
      } catch {
        // Fallback offline
        setNomorPesanan(targetId);
      } finally {
        setLoading(false);
      }
    };

    fetchDetailPesanan();
  }, [location.search, location.state]);

  // Handler tombol "Lihat status pesanan"
  const handleLihatStatus = () => {
    setTampilkanDetail(!tampilkanDetail);
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '85vh', paddingBottom: '6rem' }}>
      <div
        style={{
          maxWidth: '1060px',
          margin: '0 auto',
          padding: '2.5rem 1.5rem 1rem 1.5rem'
        }}
      >
        {/* JUDUL HALAMAN SESUAI FIGMA */}
        <h1
          style={{
            fontSize: '1.75rem',
            fontWeight: 700,
            color: '#2D231E',
            margin: '0 0 2rem 0'
          }}
        >
          Pesanan berhasil dibuat
        </h1>

        {/* KARTU KONFIRMASI TENGAH SESUAI FIGMA PPL18_01_Konfirmasi_Pesanan__HALAMAN */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            borderRadius: '16px',
            border: '1px solid #ECE7E1',
            boxShadow: '0 4px 20px rgba(67, 40, 28, 0.04)',
            maxWidth: '780px',
            margin: '0 auto',
            padding: '3.75rem 2rem 3.5rem 2rem',
            textAlign: 'center'
          }}
        >
          {/* BADGE HIJAU: Menunggu pembayaran */}
          <div style={{ display: 'flex', justifyContent: 'center', marginBottom: '1.5rem' }}>
            <span
              style={{
                backgroundColor: '#EAF5EE',
                color: '#15803D',
                padding: '0.35rem 1.15rem',
                borderRadius: '9999px',
                fontSize: '0.85rem',
                fontWeight: 600,
                display: 'inline-flex',
                alignItems: 'center',
                gap: '0.35rem'
              }}
            >
              {statusLabel}
            </span>
          </div>

          {/* JUDUL DALAM KARTU: Pesanan berhasil dibuat */}
          <h2
            style={{
              fontSize: '1.65rem',
              fontWeight: 700,
              color: '#2D231E',
              margin: '0 0 0.85rem 0'
            }}
          >
            Pesanan berhasil dibuat
          </h2>

          {/* KETERANGAN: Nomor order dibuat unik dan reservasi stok ditangani backend. */}
          <p
            style={{
              fontSize: '0.925rem',
              color: '#6B5E55',
              lineHeight: 1.5,
              margin: '0 0 1.75rem 0'
            }}
          >
            Nomor order dibuat unik dan reservasi stok ditangani backend.
          </p>

          {/* NOMOR PESANAN ASLI / CONTOH */}
          <div
            style={{
              fontSize: '1.25rem',
              fontWeight: 700,
              color: '#2D231E',
              marginBottom: '0.5rem',
              letterSpacing: '0.5px'
            }}
          >
            {isDemo ? (
              <span>Contoh: {nomorPesanan || 'PJ-DEMO-0001'}</span>
            ) : (
              <span>Nomor Pesanan: {nomorPesanan || 'PJ-20261122-4109'}</span>
            )}
          </div>

          {/* CATATAN / SUB-TEKS MERAH SESUAI FIGMA */}
          <div
            style={{
              fontSize: '0.8rem',
              color: '#B91C1C',
              marginBottom: '2.5rem'
            }}
          >
            {isDemo ? (
              'Ini adalah contoh layar, bukan transaksi nyata.'
            ) : (
              'Stok telah berhasil direservasi oleh backend selama 24 jam.'
            )}
          </div>

          {/* TOMBOL COKELAT: Lihat status pesanan */}
          <div style={{ display: 'flex', justifyContent: 'center' }}>
            <button
              type="button"
              onClick={handleLihatStatus}
              style={{
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                border: 'none',
                borderRadius: '10px',
                padding: '0.95rem 3.5rem',
                fontSize: '0.95rem',
                fontWeight: 700,
                cursor: 'pointer',
                transition: 'background-color 0.2s',
                width: '100%',
                maxWidth: '320px'
              }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#5A3726')}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '#43281C')}
            >
              {tampilkanDetail ? 'Sembunyikan status' : 'Lihat status pesanan'}
            </button>
          </div>

          {/* ==================================================================
              DETAIL STATUS PESANAN (MUNCUL KETIKA TOMBOL STATUS DIKLIK)
              ================================================================== */}
          {tampilkanDetail && (
            <div
              style={{
                marginTop: '2.5rem',
                paddingTop: '2rem',
                borderTop: '1px solid #ECE7E1',
                textAlign: 'left'
              }}
            >
              <div
                style={{
                  fontWeight: 700,
                  fontSize: '1.05rem',
                  color: '#2D231E',
                  marginBottom: '1rem'
                }}
              >
                Informasi & Riwayat Transaksi
              </div>

              <div
                style={{
                  backgroundColor: '#FAF8F5',
                  borderRadius: '12px',
                  padding: '1.25rem 1.5rem',
                  fontSize: '0.875rem',
                  color: '#5C5046',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '0.65rem'
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#78716C' }}>Nomor Pesanan</span>
                  <span style={{ fontWeight: 700, color: '#2D231E' }}>{nomorPesanan}</span>
                </div>
                {pesanan?.nomor_invoice && (
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span style={{ color: '#78716C' }}>Nomor Invoice</span>
                    <span style={{ fontWeight: 600, color: '#2D231E' }}>{pesanan.nomor_invoice}</span>
                  </div>
                )}
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#78716C' }}>Status Pembayaran</span>
                  <span style={{ fontWeight: 600, color: '#15803D' }}>{statusLabel}</span>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#78716C' }}>Status Reservasi Stok</span>
                  <span style={{ fontWeight: 600, color: '#2D231E' }}>
                    Direservasi (Holding 24 Jam)
                  </span>
                </div>
                {pesanan?.batas_waktu_pembayaran && (
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span style={{ color: '#78716C' }}>Batas Waktu Pembayaran</span>
                    <span style={{ fontWeight: 600, color: '#B91C1C' }}>
                      {pesanan.batas_waktu_pembayaran}
                    </span>
                  </div>
                )}
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#78716C' }}>Total Pembayaran</span>
                  <span style={{ fontWeight: 700, color: '#43281C', fontSize: '1rem' }}>
                    {pesanan?.total_pembayaran_format || 'Rp120.000'}
                  </span>
                </div>
                {pesanan?.penerima && (
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span style={{ color: '#78716C' }}>Nama Penerima</span>
                    <span style={{ fontWeight: 600, color: '#2D231E' }}>{pesanan.penerima}</span>
                  </div>
                )}
                {pesanan?.kurir && (
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span style={{ color: '#78716C' }}>Layanan Pengiriman</span>
                    <span style={{ fontWeight: 600, color: '#2D231E' }}>{pesanan.kurir}</span>
                  </div>
                )}
              </div>

              <div style={{ textAlign: 'center', marginTop: '1.25rem' }}>
                <Link
                  to="/produk"
                  style={{
                    color: '#43281C',
                    fontWeight: 600,
                    fontSize: '0.875rem',
                    textDecoration: 'underline'
                  }}
                >
                  Kembali ke Katalog Produk
                </Link>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
