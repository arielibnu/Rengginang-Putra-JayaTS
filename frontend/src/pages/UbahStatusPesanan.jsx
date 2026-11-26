import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useParams, useLocation } from 'react-router-dom';
import {
  Circle,
  Square,
  LayoutList,
  Diamond,
  ShoppingBag,
  Store,
  FileText,
  ArrowLeft,
  CheckCircle2,
  AlertCircle,
  Clock,
  Truck
} from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

// Data fallback demo sesuai mockup Figma PPL23_02_Modal_Ubah_Status__OVERLAY
const DEMO_STATUS_LOOKUP = {
  'PJ-DEMO-0002': {
    nomor_pesanan: 'PJ-DEMO-0002',
    status_pembayaran: 'Dibayar',
    status_pemrosesan: 'Diproses',
    default_transisi: 'Diproses -> Dikirim'
  },
  'PJ-DEMO-0001': {
    nomor_pesanan: 'PJ-DEMO-0001',
    status_pembayaran: 'Menunggu Pembayaran',
    status_pemrosesan: 'Baru',
    default_transisi: 'Menunggu Diproses -> Diproses'
  },
  'PJ-DEMO-0003': {
    nomor_pesanan: 'PJ-DEMO-0003',
    status_pembayaran: 'Dibayar',
    status_pemrosesan: 'Dikirim',
    default_transisi: 'Dikirim -> Selesai'
  }
};

export default function UbahStatusPesanan() {
  const { id } = useParams();
  const location = useLocation();
  const navigate = useNavigate();

  // Parsing parameter nomor pesanan
  const queryParams = new URLSearchParams(location.search);
  const targetId = id || queryParams.get('nomor_pesanan') || queryParams.get('id') || 'PJ-DEMO-0002';

  // State data pesanan
  const [nomorPesanan, setNomorPesanan] = useState(targetId);
  const [statusPembayaran, setStatusPembayaran] = useState('Dibayar');
  const [statusPemrosesanLama, setStatusPemrosesanLama] = useState('Diproses');
  const [statusPemrosesanBaru, setStatusPemrosesanBaru] = useState('Diproses -> Dikirim');
  const [nomorResi, setNomorResi] = useState('');

  // State feedback & loading
  const [isAuthorized, setIsAuthorized] = useState(true);
  const [isCheckingAuth, setIsCheckingAuth] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [pesanSukses, setPesanSukses] = useState('');
  const [pesanError, setPesanError] = useState('');

  // Validasi sesi Owner
  useEffect(() => {
    const controller = new AbortController();

    const cekSesi = async () => {
      try {
        const res = await fetch(`${API_URL}/api/pemilik/profil`, {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        });
        if (res.status === 401 && !localStorage.getItem('ownerAuth')) {
          navigate('/login-owner', { replace: true });
          return;
        }
        setIsAuthorized(true);
      } catch (err) {
        if (err.name !== 'AbortError') {
          setIsAuthorized(true);
        }
      }
    };

    cekSesi();
    return () => controller.abort();
  }, [navigate]);

  // Ambil data status saat ini dari backend
  useEffect(() => {
    let isMounted = true;

    const muatDetail = async () => {
      try {
        let res = await fetch(`${API_URL}/api/admin/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetId)}`, {
          credentials: 'include'
        });

        if (!res.ok) {
          res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetId)}`);
        }

        if (res.ok) {
          const json = await res.json();
          const d = json.data || {};
          const bayarStr = d.status_pembayaran || (d.status_pesanan === 'dibayar' ? 'Dibayar' : 'Menunggu Pembayaran');
          const prosesStr = d.status_pemrosesan || d.status || 'Diproses';

          if (isMounted) {
            setNomorPesanan(d.nomor_pesanan || targetId);
            setStatusPembayaran(bayarStr);
            setStatusPemrosesanLama(prosesStr);

            // Tentukan pilihan transisi default
            if (prosesStr.toLowerCase().includes('baru') || prosesStr.toLowerCase().includes('menunggu')) {
              setStatusPemrosesanBaru('Menunggu Diproses -> Diproses');
            } else if (prosesStr.toLowerCase().includes('diproses')) {
              setStatusPemrosesanBaru('Diproses -> Dikirim');
            } else if (prosesStr.toLowerCase().includes('dikirim')) {
              setStatusPemrosesanBaru('Dikirim -> Selesai');
            } else {
              setStatusPemrosesanBaru('Diproses -> Dikirim');
            }
            return;
          }
        }
      } catch (e) {
        console.warn('Fallback ke data demo:', e);
      }

      if (isMounted) {
        const demo = DEMO_STATUS_LOOKUP[targetId] || DEMO_STATUS_LOOKUP['PJ-DEMO-0002'];
        setNomorPesanan(demo.nomor_pesanan);
        setStatusPembayaran(demo.status_pembayaran);
        setStatusPemrosesanLama(demo.status_pemrosesan);
        setStatusPemrosesanBaru(demo.default_transisi);
      }
    };

    muatDetail();

    return () => {
      isMounted = false;
    };
  }, [targetId]);

  // Handler simpan status pemrosesan baru ke backend
  const handleSimpanPerubahan = async (e) => {
    e.preventDefault();

    if (isLoading) return; // Mencegah double click
    setPesanError('');
    setPesanSukses('');
    setIsLoading(true);

    // Konversi label transisi menjadi status murni untuk backend
    let statusBaruClean = 'Diproses';
    if (statusPemrosesanBaru.includes('Dikirim')) {
      statusBaruClean = 'Dikirim';
    } else if (statusPemrosesanBaru.includes('Selesai')) {
      statusBaruClean = 'Selesai';
    } else if (statusPemrosesanBaru.includes('Batal')) {
      statusBaruClean = 'Batal';
    } else if (statusPemrosesanBaru.includes('Diproses')) {
      statusBaruClean = 'Diproses';
    }

    try {
      const response = await fetch(`${API_URL}/api/admin/pesanan/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          nomor_pesanan: nomorPesanan,
          status_pemrosesan: statusBaruClean,
          nomor_resi: nomorResi || undefined
        })
      });

      const data = await response.json().catch(() => ({}));

      if (response.ok) {
        setStatusPemrosesanLama(statusBaruClean);
        setPesanSukses(`Perubahan status pesanan ${nomorPesanan} berhasil disimpan menjadi "${statusBaruClean}".`);
      } else {
        // Jika status code bukan 200, tampilkan pesan error dari backend
        setPesanError(data.pesan || `Gagal mengubah status: ${response.statusText}`);
      }
    } catch (err) {
      // Fallback offline / mockup mode
      console.warn('Backend offline fallback:', err);
      setStatusPemrosesanLama(statusBaruClean);
      setPesanSukses(`Status pesanan ${nomorPesanan} berhasil disimpan menjadi "${statusBaruClean}" (Mode Demo).`);
    } finally {
      setIsLoading(false);
    }
  };

  // Handler logout Owner
  const handleLogout = async () => {
    if (isLoggingOut) return;
    setIsLoggingOut(true);
    try {
      await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });
      localStorage.removeItem('ownerAuth');
      navigate('/login-owner', { replace: true });
    } catch {
      localStorage.removeItem('ownerAuth');
      navigate('/login-owner', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  return (
    <div className="admin-layout">
      {/* SIDEBAR OWNER SESUAI FIGMA */}
      <aside className="admin-sidebar">
        <div>
          {/* Logo Brand Putra Jaya */}
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

          {/* Menu Navigasi Admin */}
          <ul className="sidebar-menu">
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
              >
                <Circle size={16} />
                <span>Ringkasan</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/produk')}
              >
                <Square size={16} />
                <span>Kelola Produk</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/stok')}
              >
                <LayoutList size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/biaya')}
              >
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
            <li>
              {/* Menu Pesanan Aktif Sesuai Frame Figma */}
              <div
                className="sidebar-item active"
                onClick={() => navigate('/admin/pesanan')}
              >
                <ShoppingBag size={16} />
                <span>Pesanan</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/penjualan-offline')}
              >
                <Store size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
                style={{ opacity: 0.7 }}
              >
                <FileText size={16} />
                <span>Laporan</span>
              </div>
            </li>
          </ul>
        </div>

        {/* Footer Sidebar */}
        <div className="sidebar-footer">
          <Link to="/" className="sidebar-footer-link">
            Buka website publik &rsaquo;
          </Link>
          <button
            type="button"
            onClick={handleLogout}
            disabled={isLoggingOut}
            className="sidebar-footer-link"
            style={{ color: '#A92631' }}
          >
            {isLoggingOut ? 'Memproses...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* KONTEN UTAMA: UBAH STATUS PEMROSESAN */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma Frame PPL23_02_Modal_Ubah_Status__OVERLAY */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Ubah Status Pemrosesan
        </div>

        {/* Tombol Kembali ke Detail Pesanan */}
        <button
          type="button"
          className="btn-back-link"
          onClick={() => navigate(`/admin/pesanan/${nomorPesanan}`)}
          style={{ marginBottom: '1.25rem' }}
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Detail Pesanan</span>
        </button>

        {/* JUDUL DAN DESKRIPSI SESUAI FIGMA */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="admin-page-title">Ubah Status Pemrosesan</h1>
          <p className="admin-page-subtitle">
            PPL-23: perubahan status sesuai transisi yang diizinkan backend.
          </p>
        </div>

        {/* NOTIFIKASI SUKSES */}
        {pesanSukses && (
          <div
            style={{
              backgroundColor: '#E8F5E9',
              border: '1px solid #C8E6C9',
              color: '#2E7D32',
              padding: '1rem 1.25rem',
              borderRadius: '12px',
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              fontSize: '0.925rem',
              fontWeight: 600
            }}
          >
            <CheckCircle2 size={20} />
            <div style={{ flex: 1 }}>{pesanSukses}</div>
            <Link
              to={`/admin/pesanan/${nomorPesanan}`}
              style={{
                color: '#1B5E20',
                textDecoration: 'underline',
                fontSize: '0.875rem'
              }}
            >
              Lihat Detail Pesanan &rarr;
            </Link>
          </div>
        )}

        {/* NOTIFIKASI ERROR */}
        {pesanError && (
          <div
            style={{
              backgroundColor: '#FFEBEE',
              border: '1px solid #FFCDD2',
              color: '#C62828',
              padding: '1rem 1.25rem',
              borderRadius: '12px',
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              fontSize: '0.925rem'
            }}
          >
            <AlertCircle size={20} />
            <span>{pesanError}</span>
          </div>
        )}

        {/* KARTU FORMULIR UBAH STATUS SESUAI FIGMA */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            borderRadius: '16px',
            border: '1px solid #EFE8DF',
            padding: '2.5rem',
            maxWidth: '680px',
            boxShadow: '0 4px 20px rgba(67, 40, 28, 0.03)'
          }}
        >
          <form onSubmit={handleSimpanPerubahan}>
            {/* 1. Header Nomor Pesanan */}
            <div style={{ marginBottom: '1.25rem' }}>
              <h2
                style={{
                  fontSize: '1.35rem',
                  fontWeight: 800,
                  color: '#2D231E',
                  margin: 0
                }}
              >
                Pesanan: {nomorPesanan}
              </h2>
            </div>

            {/* 2. Badge Status Pembayaran (Terpisah dari Status Pemrosesan Sesuai Figma) */}
            <div style={{ marginBottom: '2rem' }}>
              <span
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  backgroundColor: statusPembayaran.toLowerCase().includes('dibayar') ? '#E8F5E9' : '#FFF3E0',
                  color: statusPembayaran.toLowerCase().includes('dibayar') ? '#2E7D32' : '#E65100',
                  padding: '0.45rem 1.15rem',
                  borderRadius: '9999px',
                  fontSize: '0.85rem',
                  fontWeight: 700
                }}
              >
                Pembayaran: {statusPembayaran}
              </span>
            </div>

            {/* 3. Input / Dropdown Status Pemrosesan Baru */}
            <div style={{ marginBottom: '1.75rem' }}>
              <label
                htmlFor="select-transisi"
                style={{
                  display: 'block',
                  fontSize: '0.95rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  marginBottom: '0.65rem'
                }}
              >
                Status pemrosesan baru
              </label>

              <select
                id="select-transisi"
                value={statusPemrosesanBaru}
                onChange={(e) => setStatusPemrosesanBaru(e.target.value)}
                style={{
                  width: '100%',
                  padding: '0.85rem 1.15rem',
                  borderRadius: '10px',
                  border: '1px solid #D8C7B8',
                  backgroundColor: '#FFFFFF',
                  color: '#2D231E',
                  fontSize: '0.975rem',
                  fontWeight: 600,
                  outline: 'none',
                  cursor: 'pointer'
                }}
              >
                <option value="Diproses -> Dikirim">Diproses &rarr; Dikirim</option>
                <option value="Menunggu Diproses -> Diproses">Menunggu Diproses &rarr; Diproses</option>
                <option value="Dikirim -> Selesai">Dikirim &rarr; Selesai</option>
                <option value="Batal (Kembalikan Reservasi Stok)">Batal (Kembalikan Reservasi Stok)</option>
              </select>
            </div>

            {/* Input Tambahan Nomor Resi jika status Dikirim (WBS 1.3.4.3) */}
            {statusPemrosesanBaru.includes('Dikirim') && (
              <div style={{ marginBottom: '1.75rem' }}>
                <label
                  htmlFor="input-resi"
                  style={{
                    display: 'block',
                    fontSize: '0.875rem',
                    fontWeight: 700,
                    color: '#2D231E',
                    marginBottom: '0.5rem'
                  }}
                >
                  Nomor Resi Pengiriman (Opsional):
                </label>
                <div style={{ position: 'relative' }}>
                  <Truck
                    size={16}
                    style={{
                      position: 'absolute',
                      left: '1rem',
                      top: '50%',
                      transform: 'translateY(-50%)',
                      color: '#8E7D73'
                    }}
                  />
                  <input
                    id="input-resi"
                    type="text"
                    placeholder="Contoh: JNT-8892019920"
                    value={nomorResi}
                    onChange={(e) => setNomorResi(e.target.value)}
                    style={{
                      width: '100%',
                      padding: '0.8rem 1rem 0.8rem 2.75rem',
                      borderRadius: '10px',
                      border: '1px solid #D8C7B8',
                      fontSize: '0.9rem',
                      color: '#2D231E',
                      outline: 'none'
                    }}
                  />
                </div>
              </div>
            )}

            {/* 4. Kotak Peringatan / Notice Box Sesuai Figma */}
            <div
              style={{
                backgroundColor: '#FFF6EB',
                border: '1px solid #FCDDC4',
                borderRadius: '12px',
                padding: '1.25rem 1.5rem',
                marginBottom: '2.25rem'
              }}
            >
              <div
                style={{
                  fontSize: '0.95rem',
                  fontWeight: 800,
                  color: '#8D4004',
                  marginBottom: '0.35rem'
                }}
              >
                Pembatalan tidak otomatis berarti refund berhasil.
              </div>
              <div
                style={{
                  fontSize: '0.875rem',
                  color: '#7A6E65',
                  lineHeight: 1.4
                }}
              >
                Stok reservasi tetap dikendalikan backend.
              </div>
            </div>

            {/* 5. Tombol Simpan Perubahan Sesuai Figma */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
              <button
                type="button"
                onClick={() => navigate(`/admin/pesanan/${nomorPesanan}`)}
                style={{
                  padding: '0.85rem 1.5rem',
                  borderRadius: '10px',
                  border: '1px solid #D8C7B8',
                  backgroundColor: '#FFFFFF',
                  color: '#43281C',
                  fontSize: '0.95rem',
                  fontWeight: 700,
                  cursor: 'pointer'
                }}
              >
                Batal
              </button>

              <button
                type="submit"
                disabled={isLoading}
                style={{
                  padding: '0.85rem 2rem',
                  borderRadius: '10px',
                  border: 'none',
                  backgroundColor: isLoading ? '#8E7D73' : '#4A2C17',
                  color: '#FFFFFF',
                  fontSize: '0.95rem',
                  fontWeight: 700,
                  cursor: isLoading ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: '0 4px 12px rgba(74, 44, 23, 0.15)'
                }}
              >
                {isLoading ? 'Menyimpan...' : 'Simpan Perubahan'}
              </button>
            </div>
          </form>
        </div>

        {/* Tautan ke Invoice Digital untuk kemudahan Owner */}
        <div style={{ marginTop: '2rem' }}>
          <Link
            to={`/invoice/${nomorPesanan}`}
            target="_blank"
            rel="noopener noreferrer"
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: '0.4rem',
              color: '#43281C',
              fontSize: '0.875rem',
              fontWeight: 600,
              textDecoration: 'none'
            }}
          >
            <FileText size={16} />
            <span>Buka Invoice Digital Pesanan Ini &rsaquo;</span>
          </Link>
        </div>

      </main>
    </div>
  );
}
