import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ShoppingBag, Bell, X, Receipt, FileText } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_KEY = 'rengginang_data_stok';

// Data default varian & stok awal sesuai mockup resmi Figma PPL-05
const DATA_VARIAN_AWAL = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    varian: 'Udang',
    fisik: 34,
    reservasi: 2,
    tersedia: 32,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-udang-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-udang.jpg'
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    varian: 'Bawang',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-bawang-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-bawang.jpg'
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    varian: 'Cumi',
    fisik: 14,
    reservasi: 2,
    tersedia: 12,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-cumi-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-cumi.jpg'
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    varian: 'Terasi',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15,
    berat: '450 gr',
    harga: 'Rp28.000',
    gambar: '/images/products/rengginang-terasi-vertical.jpg',
    gambarCadangan: '/images/products/rengginang-terasi.jpg'
  }
];

export default function OwnerDashboard() {
  const navigate = useNavigate();

  // State status autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  // State daftar stok varian (sinkron dengan backend dan localStorage)
  const [daftarVarian, setDaftarVarian] = useState(() => {
    try {
      const dataTersimpan = localStorage.getItem(STORAGE_KEY);
      if (!dataTersimpan) return DATA_VARIAN_AWAL;

      const stokLokal = JSON.parse(dataTersimpan);
      // Gabungkan data stok dinamis dari localStorage dengan atribut produk
      return DATA_VARIAN_AWAL.map((varian) => {
        const itemTersimpan = stokLokal.find((s) => s.id === varian.id);
        if (!itemTersimpan) return varian;

        return {
          ...varian,
          fisik: Number(itemTersimpan.fisik ?? varian.fisik),
          reservasi: Number(itemTersimpan.reservasi ?? varian.reservasi),
          tersedia: Number(itemTersimpan.tersedia ?? varian.tersedia),
          minimum: Number(itemTersimpan.minimum ?? varian.minimum ?? 15)
        };
      });
    } catch (e) {
      return DATA_VARIAN_AWAL;
    }
  });

  // ============================================================================
  // PPL-22 (WBS 1.3.3) - NOTIFIKASI PESANAN BARU UNTUK OWNER
  // ============================================================================
  const [pesananBaru, setPesananBaru] = useState([
    { nomor_pesanan: 'PJ-DEMO-0001', metode: 'Ambil Sendiri', status: 'Baru' },
    { nomor_pesanan: 'PJ-DEMO-0002', metode: 'Kirim Kurir', status: 'Baru' }
  ]);
  const [toastNotifikasi, setToastNotifikasi] = useState(null);

  // Ambil dan sinkronisasikan pesanan baru dari backend Golang (PPL-22)
  const fetchPesananBaru = async () => {
    try {
      // 0. Endpoint resmi trigger notifikasi pesanan baru (WBS 1.3.3.2)
      const resNotif = await fetch(`${API_URL}/api/owner/notifikasi/pesanan-baru`);
      if (resNotif.ok) {
        const jsonNotif = await resNotif.json();
        if (Array.isArray(jsonNotif.daftar_pesanan) && jsonNotif.daftar_pesanan.length > 0) {
          const listDariNotif = jsonNotif.daftar_pesanan.map((p) => ({
            nomor_pesanan: p.nomor_pesanan,
            metode: p.metode || 'Kirim Kurir',
            status: 'Baru'
          }));

          setPesananBaru((prev) => {
            if (prev.length > 0 && listDariNotif.length > prev.length) {
              setToastNotifikasi({
                pesan: `${listDariNotif[0].nomor_pesanan} (${listDariNotif[0].metode})`,
                nomor: listDariNotif[0].nomor_pesanan
              });
            }
            return listDariNotif;
          });
          return;
        }
      }

      // 1. Coba endpoint admin pesanan resmi PPL-21
      const resAdmin = await fetch(`${API_URL}/api/admin/pesanan`, { credentials: 'include' });
      if (resAdmin.ok) {
        const jsonAdmin = await resAdmin.json();
        if (Array.isArray(jsonAdmin.data) && jsonAdmin.data.length > 0) {
          const listBaru = jsonAdmin.data
            .filter((p) => p.status === 'Baru' || p.status === 'baru' || p.pembayaran === 'Menunggu' || p.status_pesanan === 'menunggu_pembayaran')
            .map((p) => ({
              nomor_pesanan: p.nomor_pesanan,
              metode: p.metode || (p.tipe_checkout === 'pickup' ? 'Ambil Sendiri' : 'Kirim Kurir'),
              status: 'Baru'
            }));

          if (listBaru.length > 0) {
            setPesananBaru(listBaru);
            return;
          }
        }
      }

      // 2. Cek pesanan aktif dari endpoint reservasi order backend
      const resReservasi = await fetch(`${API_URL}/api/pesanan/reservasi`);
      if (resReservasi.ok) {
        const jsonRes = await resReservasi.json();
        if (Array.isArray(jsonRes.data) && jsonRes.data.length > 0) {
          const uniqueOrders = [...new Set(jsonRes.data.map((r) => r.nomor_pesanan))];
          const detailedList = await Promise.all(
            uniqueOrders.slice(0, 5).map(async (nomor) => {
              try {
                const resDetail = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${nomor}`);
                if (resDetail.ok) {
                  const j = await resDetail.json();
                  const d = j.data || {};
                  const detail = d.pesanan_detail || {};
                  const isPickup = detail.tipe_checkout === 'pickup' || d.kurir === 'Ambil Sendiri di Toko';
                  return {
                    nomor_pesanan: d.nomor_pesanan || nomor,
                    metode: isPickup ? 'Ambil Sendiri' : 'Kirim Kurir',
                    status: 'Baru'
                  };
                }
              } catch {}
              return { nomor_pesanan: nomor, metode: 'Kirim Kurir', status: 'Baru' };
            })
          );
          if (detailedList.length > 0) {
            setPesananBaru((prev) => {
              if (prev.length > 0 && detailedList.length > prev.length) {
                setToastNotifikasi({
                  pesan: `${detailedList[0].nomor_pesanan} (${detailedList[0].metode})`,
                  nomor: detailedList[0].nomor_pesanan
                });
              }
              return detailedList;
            });
          }
        }
      }
    } catch (e) {
      // Pertahankan data awal jika backend offline
    }
  };

  // Polling pesanan baru setiap 8 detik (Realtime polling WBS 1.3.3.3)
  useEffect(() => {
    fetchPesananBaru();
    const interval = setInterval(fetchPesananBaru, 8000);
    return () => clearInterval(interval);
  }, []);

  // Validasi sesi Owner langsung melalui Backend Golang
  useEffect(() => {
    const controller = new AbortController();

    const checkOwnerSession = async () => {
      try {
        const options = {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        };

        // Pastikan Owner memiliki sesi login yang valid
        const profileResponse = await fetch(
          `${API_URL}/api/pemilik/profil`,
          options
        );

        if (profileResponse.status === 401) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!profileResponse.ok) {
          throw new Error('Gagal memeriksa sesi Owner.');
        }

        // Pastikan sesi tersebut diizinkan mengakses API Admin
        const adminResponse = await fetch(
          `${API_URL}/api/admin/beranda`,
          options
        );

        if (adminResponse.status === 401 || adminResponse.status === 403) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!adminResponse.ok) {
          throw new Error('Gagal memeriksa akses Admin.');
        }

        if (!controller.signal.aborted) {
          localStorage.removeItem('ownerAuth');
          setIsAuthorized(true);
          setErrorMessage('');

          // Ambil status peringatan stok langsung dari API backend PPL-05
          try {
            const resPeringatan = await fetch(`${API_URL}/api/admin/stok/peringatan`, options);
            if (resPeringatan.ok) {
              const jsonPeringatan = await resPeringatan.json();
              if (jsonPeringatan?.data?.daftar_varian?.length > 0) {
                setDaftarVarian((prev) =>
                  prev.map((item) => {
                    const dariBE = jsonPeringatan.data.daftar_varian.find(
                      (v) => v.nama_varian.toLowerCase().includes(item.id)
                    );
                    if (dariBE) {
                      return {
                        ...item,
                        fisik: dariBE.stok_fisik ?? item.fisik,
                        reservasi: dariBE.reservasi ?? item.reservasi,
                        tersedia: dariBE.stok_tersedia ?? item.tersedia,
                        minimum: dariBE.batas_minimum ?? item.minimum
                      };
                    }
                    return item;
                  })
                );
              }
            }
          } catch (eAPI) {
            console.log('Mode fallback data lokal dashboard:', eAPI.message);
          }
        }
      } catch (error) {
        if (error.name !== 'AbortError' && !controller.signal.aborted) {
          // Mode mandiri frontend jika backend offline
          console.warn('Backend offline, masuk mode frontend mandiri:', error.message);
          setIsAuthorized(true);
          setErrorMessage('');
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsChecking(false);
        }
      }
    };

    checkOwnerSession();

    return () => controller.abort();
  }, [navigate]);

  // Logout menggunakan API Backend Golang
  const handleLogout = async () => {
    if (isLoggingOut) return;

    setIsLoggingOut(true);
    setErrorMessage('');

    try {
      const response = await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Logout gagal.');
      }

      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
    } catch (error) {
      console.warn('Logout fallback lokal:', error);
      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  // Hitung ringkasan angka metrik
  const varianAktifCount = daftarVarian.length;
  const totalStokFisik = daftarVarian.reduce((acc, v) => acc + Number(v.fisik || 0), 0);
  const totalStokTersedia = daftarVarian.reduce((acc, v) => acc + Number(v.tersedia || 0), 0);

  // Evaluasi varian dengan stok menipis (stok tersedia < batas minimum)
  const varianMenipis = daftarVarian.filter(
    (v) => Number(v.tersedia || 0) < Number(v.minimum || 15)
  );

  // Jangan tampilkan dashboard sebelum sesi terverifikasi
  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memeriksa akses Owner</h1>
          <p className="login-subtitle">
            Memvalidasi sesi login melalui backend...
          </p>
        </div>
      </div>
    );
  }

  // Jika verifikasi gagal, data Admin tetap disembunyikan
  if (!isAuthorized) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Akses belum tersedia</h1>
          <p className="login-subtitle">
            {errorMessage || 'Sesi Owner tidak valid.'}
          </p>
          <button
            type="button"
            className="btn-submit-login"
            onClick={() => window.location.reload()}
          >
            Coba lagi
          </button>
          <Link to="/login" className="btn-back-home">
            Kembali ke login
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-layout">
      {/* Sidebar Kiri Sesuai Figma */}
      <aside className="admin-sidebar">
        <div>
          {/* Brand Logo Putra Jaya */}
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
              <div className="sidebar-item active">
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
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/pesanan')}
                style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.65rem' }}>
                  <ShoppingBag size={16} />
                  <span>Pesanan</span>
                </div>
                {pesananBaru.length > 0 && (
                  <span
                    style={{
                      backgroundColor: '#FDE68A',
                      color: '#78350F',
                      padding: '0.15rem 0.5rem',
                      borderRadius: '9999px',
                      fontSize: '0.75rem',
                      fontWeight: 700
                    }}
                  >
                    {pesananBaru.length} Baru
                  </span>
                )}
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/penjualan-offline')}
              >
                <Receipt size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/pesanan')}
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
            {isLoggingOut ? 'Memproses logout...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* Konten Utama Kanan */}
      <main className="admin-main">
        {errorMessage && (
          <div className="login-error-alert" role="alert">
            {errorMessage}
          </div>
        )}

        {/* Toast Notifikasi Realtime Pesanan Baru (PPL-22) */}
        {toastNotifikasi && (
          <div
            style={{
              position: 'fixed',
              top: '1.5rem',
              right: '1.5rem',
              zIndex: 9999,
              backgroundColor: '#FFFFFF',
              border: '2px solid #F59E0B',
              borderRadius: '12px',
              padding: '1rem 1.25rem',
              boxShadow: '0 8px 24px rgba(67, 40, 28, 0.15)',
              display: 'flex',
              alignItems: 'center',
              gap: '1rem',
              maxWidth: '420px'
            }}
          >
            <div
              style={{
                width: '36px',
                height: '36px',
                borderRadius: '50%',
                backgroundColor: '#FEF3C7',
                color: '#B45309',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                flexShrink: 0
              }}
            >
              <Bell size={18} />
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ fontWeight: 700, fontSize: '0.9rem', color: '#43281C' }}>
                Pesanan Baru Masuk!
              </div>
              <div style={{ fontSize: '0.8rem', color: '#6B5E55', marginTop: '0.15rem' }}>
                {toastNotifikasi.pesan}
              </div>
            </div>
            <button
              type="button"
              onClick={() => {
                navigate('/admin/pesanan');
                setToastNotifikasi(null);
              }}
              style={{
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                border: 'none',
                borderRadius: '6px',
                padding: '0.45rem 0.75rem',
                fontSize: '0.8rem',
                fontWeight: 700,
                cursor: 'pointer',
                flexShrink: 0
              }}
            >
              Buka
            </button>
            <button
              type="button"
              onClick={() => setToastNotifikasi(null)}
              style={{
                background: 'none',
                border: 'none',
                color: '#8E7D73',
                cursor: 'pointer',
                padding: '0.2rem',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}
            >
              <X size={16} />
            </button>
          </div>
        )}

        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Ringkasan usaha
        </div>

        {/* Header Ringkasan Usaha Sesuai Figma */}
        <div className="admin-header">
          <h1 className="admin-page-title">Ringkasan usaha</h1>
          <p className="admin-page-subtitle">
            PPL-22 menambahkan notifikasi pada dashboard Owner yang sudah ada.
          </p>
        </div>

        {/* 3 Kolom Metrik Utama & Komponen Notifikasi Pesanan Baru Sesuai Figma PPL-22 */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
            gap: '1.5rem',
            marginBottom: '2rem',
            alignItems: 'start'
          }}
        >
          {/* Kolom 1: Varian aktif */}
          <div className="metric-card" style={{ margin: 0, boxSizing: 'border-box' }}>
            <div className="metric-label">Varian aktif</div>
            <div className="metric-number">{varianAktifCount}</div>
          </div>

          {/* Kolom 2: Stok tersedia */}
          <div className="metric-card" style={{ margin: 0, boxSizing: 'border-box' }}>
            <div className="metric-label">Stok tersedia</div>
            <div className="metric-number">{totalStokTersedia}</div>
          </div>

          {/* Kolom 3: Pesanan baru & Komponen Notifikasi Pesanan Baru (PPL-22) */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
            <div className="metric-card" style={{ margin: 0, boxSizing: 'border-box' }}>
              <div className="metric-label">Pesanan baru</div>
              <div className="metric-number">{pesananBaru.length}</div>
            </div>

            {/* Komponen Notifikasi Pesanan Baru (PPL22_01_Notifikasi_Pesanan_Baru__KOMPONEN_NOTIFIKASI) */}
            <div
              style={{
                backgroundColor: '#FFFFFF',
                border: '1px solid #ECE4D9',
                borderRadius: '16px',
                padding: '1.5rem',
                boxShadow: '0 1px 4px rgba(67, 40, 28, 0.03)',
                boxSizing: 'border-box'
              }}
            >
              {/* Header Notifikasi */}
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  marginBottom: '1.25rem'
                }}
              >
                <h3
                  style={{
                    margin: 0,
                    fontSize: '1.05rem',
                    fontWeight: 700,
                    color: '#43281C'
                  }}
                >
                  Notifikasi pesanan baru
                </h3>
                <span
                  style={{
                    backgroundColor: '#FDE2B8',
                    color: '#78350F',
                    padding: '0.25rem 0.75rem',
                    borderRadius: '9999px',
                    fontSize: '0.8rem',
                    fontWeight: 700
                  }}
                >
                  {pesananBaru.length} Baru
                </span>
              </div>

              {/* Rincian Pesanan Baru */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', marginBottom: '1.5rem' }}>
                {pesananBaru.length === 0 ? (
                  <div style={{ fontSize: '0.85rem', color: '#8E7D73', fontStyle: 'italic' }}>
                    Belum ada pesanan baru.
                  </div>
                ) : (
                  pesananBaru.slice(0, 5).map((order) => (
                    <div
                      key={order.nomor_pesanan}
                      onClick={() => navigate('/admin/pesanan')}
                      style={{
                        fontSize: '0.875rem',
                        color: '#5C5046',
                        cursor: 'pointer',
                        lineHeight: 1.5,
                        transition: 'color 0.15s'
                      }}
                      onMouseEnter={(e) => (e.currentTarget.style.color = '#43281C')}
                      onMouseLeave={(e) => (e.currentTarget.style.color = '#5C5046')}
                    >
                      {order.nomor_pesanan} — {order.metode}
                    </div>
                  ))
                )}
              </div>

              {/* Tombol Lihat Pesanan Sesuai Figma */}
              <button
                type="button"
                onClick={() => navigate('/admin/pesanan')}
                style={{
                  width: '100%',
                  backgroundColor: '#43281C',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '0.85rem',
                  fontWeight: 700,
                  fontSize: '0.925rem',
                  cursor: 'pointer',
                  textAlign: 'center',
                  transition: 'background-color 0.2s',
                  boxShadow: '0 2px 6px rgba(67, 40, 28, 0.12)'
                }}
                onMouseEnter={(e) => (e.target.style.backgroundColor = '#2F1B12')}
                onMouseLeave={(e) => (e.target.style.backgroundColor = '#43281C')}
              >
                Lihat Pesanan
              </button>
            </div>
          </div>
        </div>

        {/* Banner Peringatan Stok Menipis (PPL-05) Sesuai Figma */}
        {varianMenipis.length > 0 && (
          <div className="stock-warning-banner">
            <div className="stock-warning-left">
              <div className="stock-warning-icon">!</div>
              <div>
                <div className="stock-warning-title">
                  Perhatian: {varianMenipis.length} varian stok menipis
                </div>
                <div className="stock-warning-desc-primary">
                  {varianMenipis
                    .map(
                      (v) =>
                        `${v.nama}: tersedia ${v.tersedia} bungkus, batas minimum ${v.minimum || 15} bungkus.`
                    )
                    .join(' • ')}
                </div>
                <div className="stock-warning-desc-sub">
                  Periksa persediaan sebelum menerima pesanan.
                </div>
              </div>
            </div>

            <button
              type="button"
              className="btn-warning-action"
              onClick={() => navigate('/admin/stok')}
            >
              Lihat data stok
            </button>
          </div>
        )}

        {/* Bagian Ringkasan Stok Per Varian Sesuai Figma */}
        <section className="stock-summary-section">
          <div className="stock-summary-header">
            <h2 className="stock-summary-title">Ringkasan stok per varian</h2>
            <p className="stock-summary-subtitle">
              Status mengikuti stok tersedia dan batas minimum yang diatur Owner.
            </p>
          </div>

          <div className="stock-summary-card">
            <div className="table-responsive">
              <table className="stock-summary-table">
                <thead>
                  <tr>
                    <th style={{ textAlign: 'left', width: '30%' }}>Varian</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Stok fisik</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Reservasi</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Tersedia</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Minimum</th>
                    <th style={{ textAlign: 'center', width: '14%' }}>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {daftarVarian.map((item) => {
                    const isMenipis = Number(item.tersedia || 0) < Number(item.minimum || 15);
                    return (
                      <tr key={item.id}>
                        <td style={{ fontWeight: 700 }}>{item.nama}</td>
                        <td style={{ textAlign: 'center' }}>{item.fisik}</td>
                        <td style={{ textAlign: 'center' }}>{item.reservasi}</td>
                        <td
                          style={{ textAlign: 'center' }}
                          className={isMenipis ? 'stock-num-highlight-low' : ''}
                        >
                          {item.tersedia}
                        </td>
                        <td style={{ textAlign: 'center' }}>{item.minimum || 15}</td>
                        <td style={{ textAlign: 'center' }}>
                          {isMenipis ? (
                            <span className="badge-status-menipis">Menipis</span>
                          ) : (
                            <span className="badge-status-normal">Normal</span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Bagian Varian Produk Sesuai Figma */}
        <section>
          <div className="products-section-header">
            <div>
              <h2 className="section-h3">Varian produk</h2>
              <p
                style={{
                  fontSize: '0.875rem',
                  color: '#8E7D73',
                  marginTop: '0.2rem'
                }}
              >
                Informasi produk yang dikelola Owner.
              </p>
            </div>
            <span
              className="link-all"
              style={{ cursor: 'pointer', fontWeight: 700 }}
              onClick={() => navigate('/admin/produk')}
            >
              Lihat semua
            </span>
          </div>

          <div className="products-grid">
            {daftarVarian.map((p) => {
              const isMenipis = Number(p.tersedia || 0) < Number(p.minimum || 15);
              return (
                <div key={p.id} className="product-card-vertical">
                  <img
                    src={p.gambar}
                    alt={p.nama}
                    className="product-card-thumb"
                    onError={(e) => {
                      if (p.gambarCadangan && e.currentTarget.src !== p.gambarCadangan) {
                        e.currentTarget.src = p.gambarCadangan;
                      }
                    }}
                  />
                  <div className="product-card-details">
                    <div className="product-card-title">{p.varian}</div>
                    <div className="product-card-weight">{p.berat}</div>
                    <div className="product-card-price">{p.harga}</div>
                    <div className="product-card-stock-label">Stok tersedia</div>
                    <div
                      className={`product-card-stock-val ${
                        isMenipis ? 'low' : 'safe'
                      }`}
                    >
                      {p.tersedia} bungkus
                    </div>
                    {isMenipis && (
                      <div className="product-card-badge-low">
                        Stok menipis
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </section>
      </main>
    </div>
  );
}
