import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  Package,
  Layers,
  ShoppingBag,
  DollarSign,
  TrendingDown,
  RefreshCw,
  LogOut,
  FileText,
  AlertCircle
} from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const DAFTAR_BULAN = [
  'Januari',
  'Februari',
  'Maret',
  'April',
  'Mei',
  'Juni',
  'Juli',
  'Agustus',
  'September',
  'Oktober',
  'November',
  'Desember'
];

export default function LaporanPenerimaanKas() {
  const navigate = useNavigate();

  // Autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(true);
  const [isCheckingAuth, setIsCheckingAuth] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Filter Periode Sesuai Mockup Figma PPL27_01_Laporan_Penerimaan_Kas__TAB
  const [bulanTerpilih, setBulanTerpilih] = useState('November');
  const [tahunTerpilih, setTahunTerpilih] = useState('2026');

  // State Data 3 Kartu Ringkasan Penerimaan Kas Sesuai Figma PPL27_01:
  // 1. Kas masuk nyata
  // 2. Simulasi sandbox
  // 3. Status pending
  const [kasMasukNyata, setKasMasukNyata] = useState(null);
  const [simulasiSandbox, setSimulasiSandbox] = useState(null);
  const [statusPending, setStatusPending] = useState(null);

  // Status Laporan & Keterangan Sumber Data Sesuai Figma
  const [statusTransaksiValid, setStatusTransaksiValid] = useState('Mengikuti backend');
  const [statusSandbox, setStatusSandbox] = useState('Tidak dihitung sebagai kas nyata');
  const [statusDataBelumLengkap, setStatusDataBelumLengkap] = useState('Ditandai dan dijelaskan');

  // State Loading & Error
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

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

  // Fungsi memuat data laporan penerimaan kas dari backend
  const muatLaporanPenerimaanKas = async (bulan, tahun) => {
    setIsLoading(true);
    setErrorMessage('');

    try {
      // Panggil endpoint laporan penerimaan kas berkala jika backend sudah menyediakannya
      const res = await fetch(
        `${API_URL}/api/admin/laporan/penerimaan-kas?bulan=${encodeURIComponent(bulan)}&tahun=${encodeURIComponent(tahun)}`,
        { credentials: 'include' }
      );

      if (res.ok) {
        const json = await res.json();
        const data = json.data || {};

        // Format nilai uang
        setKasMasukNyata(
          typeof data.kas_masuk_nyata === 'number'
            ? `Rp${data.kas_masuk_nyata.toLocaleString('id-ID')}`
            : (data.kas_masuk_nyata ?? null)
        );
        setSimulasiSandbox(
          typeof data.simulasi_sandbox === 'number'
            ? `Rp${data.simulasi_sandbox.toLocaleString('id-ID')}`
            : (data.simulasi_sandbox ?? null)
        );
        setStatusPending(
          typeof data.status_pending === 'number'
            ? `Rp${data.status_pending.toLocaleString('id-ID')}`
            : (data.status_pending ?? null)
        );

        setStatusTransaksiValid('Mengikuti backend');
        setStatusSandbox('Tidak dihitung sebagai kas nyata');
        setStatusDataBelumLengkap('Ditandai dan dijelaskan');
        setIsLoading(false);
        return;
      }
    } catch (e) {
      console.warn('Endpoint khusus laporan penerimaan kas belum tersedia di backend:', e);
    }

    // Jika API belum tersedia atau terjadi error, tampilkan kondisi menunggu data riil persis desain Figma
    // JANGAN membuat data angka dummy atau data sandbox sebagai kas nyata!
    setKasMasukNyata(null);
    setSimulasiSandbox(null);
    setStatusPending(null);
    setStatusTransaksiValid('Mengikuti backend');
    setStatusSandbox('Tidak dihitung sebagai kas nyata');
    setStatusDataBelumLengkap('Ditandai dan dijelaskan');
    setIsLoading(false);
  };

  // Muat data saat halaman pertama kali dibuka
  useEffect(() => {
    muatLaporanPenerimaanKas(bulanTerpilih, tahunTerpilih);
  }, []);

  // Handler tombol "Terapkan"
  const handleTerapkan = (e) => {
    e.preventDefault();
    muatLaporanPenerimaanKas(bulanTerpilih, tahunTerpilih);
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
              <div className="sidebar-title">PUTRA JAYA</div>
              <div className="sidebar-subtitle">Dashboard Owner</div>
            </div>
          </div>

          {/* Navigasi Menu Dashboard Owner */}
          <ul className="sidebar-menu">
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin')}
              >
                <Package size={16} />
                <span>Ringkasan</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/produk')}
              >
                <Layers size={16} />
                <span>Kelola Produk</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/stok')}
              >
                <Package size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/pengeluaran')}
              >
                <DollarSign size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
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
                <TrendingDown size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              {/* Menu Laporan Aktif Sesuai Frame Figma PPL27_01 */}
              <div
                className="sidebar-item active"
                onClick={() => navigate('/admin/laporan/penerimaan-kas')}
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

      {/* KONTEN UTAMA: LAPORAN PENERIMAAN KAS */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Laporan Penerimaan Kas
        </div>

        {/* Tab Navigasi Antar Laporan (Laporan Penjualan, Laporan Stok Akhir, Laporan Penerimaan Kas, Laporan Pengeluaran) */}
        <div
          style={{
            display: 'flex',
            gap: '1rem',
            marginBottom: '1.5rem',
            borderBottom: '1px solid #EFE8DF',
            paddingBottom: '0.5rem',
            flexWrap: 'wrap'
          }}
        >
          <button
            type="button"
            onClick={() => navigate('/admin/laporan/penjualan')}
            style={{
              background: 'none',
              border: 'none',
              padding: '0.5rem 0.75rem',
              fontSize: '0.95rem',
              fontWeight: 600,
              color: '#7A6E65',
              cursor: 'pointer',
              borderBottom: '2px solid transparent'
            }}
          >
            Laporan Penjualan
          </button>
          <button
            type="button"
            onClick={() => navigate('/admin/laporan/stok-akhir')}
            style={{
              background: 'none',
              border: 'none',
              padding: '0.5rem 0.75rem',
              fontSize: '0.95rem',
              fontWeight: 600,
              color: '#7A6E65',
              cursor: 'pointer',
              borderBottom: '2px solid transparent'
            }}
          >
            Laporan Stok Akhir
          </button>
          <button
            type="button"
            style={{
              background: 'none',
              border: 'none',
              padding: '0.5rem 0.75rem',
              fontSize: '0.95rem',
              fontWeight: 700,
              color: '#43281C',
              cursor: 'default',
              borderBottom: '2px solid #43281C'
            }}
          >
            Laporan Penerimaan Kas
          </button>
          <button
            type="button"
            onClick={() => navigate('/admin/laporan/pengeluaran')}
            style={{
              background: 'none',
              border: 'none',
              padding: '0.5rem 0.75rem',
              fontSize: '0.95rem',
              fontWeight: 600,
              color: '#7A6E65',
              cursor: 'pointer',
              borderBottom: '2px solid transparent'
            }}
          >
            Laporan Pengeluaran
          </button>
        </div>

        {/* Judul & Subtitle Sesuai Figma PPL27_01_Laporan_Penerimaan_Kas__TAB */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="admin-page-title">Laporan Penerimaan Kas</h1>
          <p className="admin-page-subtitle">
            Tampilkan data berdasarkan bulan dan tahun. Simulasi dipisahkan dari transaksi nyata.
          </p>
        </div>

        {/* Pesan Error Jika Ada */}
        {errorMessage && (
          <div
            style={{
              backgroundColor: '#FEE2E2',
              border: '1px solid #FCA5A5',
              borderRadius: '12px',
              padding: '1rem 1.25rem',
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              color: '#991B1B',
              fontSize: '0.9rem'
            }}
          >
            <AlertCircle size={18} />
            <span>{errorMessage}</span>
          </div>
        )}

        {/* FORM FILTER BULAN & TAHUN SESUAI FIGMA */}
        <form
          onSubmit={handleTerapkan}
          style={{
            display: 'flex',
            alignItems: 'flex-end',
            justifyContent: 'space-between',
            gap: '1.25rem',
            marginBottom: '2.5rem',
            flexWrap: 'wrap'
          }}
        >
          {/* Kelompok Input Filter (Kiri) */}
          <div
            style={{
              display: 'flex',
              alignItems: 'flex-end',
              gap: '1.25rem',
              flexWrap: 'wrap',
              flex: '1 1 auto'
            }}
          >
            {/* Filter Bulan */}
            <div style={{ flex: '1 1 200px', maxWidth: '280px' }}>
              <label
                htmlFor="filter-bulan"
                style={{
                  display: 'block',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  marginBottom: '0.5rem'
                }}
              >
                Bulan
              </label>
              <select
                id="filter-bulan"
                value={bulanTerpilih}
                onChange={(e) => setBulanTerpilih(e.target.value)}
                style={{
                  width: '100%',
                  padding: '0.8rem 1rem',
                  borderRadius: '10px',
                  border: '1px solid #7A6E65',
                  backgroundColor: '#FFFFFF',
                  color: '#2D231E',
                  fontSize: '0.95rem',
                  fontWeight: 600,
                  outline: 'none',
                  cursor: 'pointer'
                }}
              >
                {DAFTAR_BULAN.map((bln) => (
                  <option key={bln} value={bln}>
                    {bln}
                  </option>
                ))}
              </select>
            </div>

            {/* Filter Tahun */}
            <div style={{ flex: '1 1 200px', maxWidth: '280px' }}>
              <label
                htmlFor="filter-tahun"
                style={{
                  display: 'block',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  marginBottom: '0.5rem'
                }}
              >
                Tahun
              </label>
              <input
                id="filter-tahun"
                type="text"
                value={tahunTerpilih}
                onChange={(e) => setTahunTerpilih(e.target.value)}
                placeholder="Contoh: 2026"
                style={{
                  width: '100%',
                  padding: '0.8rem 1rem',
                  borderRadius: '10px',
                  border: '1px solid #7A6E65',
                  backgroundColor: '#FFFFFF',
                  color: '#2D231E',
                  fontSize: '0.95rem',
                  fontWeight: 600,
                  outline: 'none'
                }}
              />
            </div>
          </div>

          {/* Tombol Terapkan (Posisi Kanan Sesuai Desain Figma) */}
          <div style={{ marginLeft: 'auto' }}>
            <button
              type="submit"
              disabled={isLoading}
              style={{
                padding: '0.8rem 2.5rem',
                borderRadius: '10px',
                border: '1px solid #43281C',
                backgroundColor: '#FFFFFF',
                color: '#43281C',
                fontSize: '0.95rem',
                fontWeight: 700,
                cursor: isLoading ? 'not-allowed' : 'pointer',
                transition: 'all 0.2s ease',
                height: '46px',
                display: 'inline-flex',
                alignItems: 'center',
                gap: '0.5rem'
              }}
              onMouseEnter={(e) => {
                if (!isLoading) {
                  e.currentTarget.style.backgroundColor = '#FAF5EE';
                }
              }}
              onMouseLeave={(e) => {
                if (!isLoading) {
                  e.currentTarget.style.backgroundColor = '#FFFFFF';
                }
              }}
            >
              {isLoading ? (
                <>
                  <RefreshCw size={16} className="animate-spin" />
                  <span>Memuat...</span>
                </>
              ) : (
                <span>Terapkan</span>
              )}
            </button>
          </div>
        </form>

        {/* 3 KARTU RINGKASAN PENERIMAAN KAS SESUAI FIGMA PPL27_01 */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
            gap: '1.25rem',
            marginBottom: '2.5rem'
          }}
        >
          {/* Kartu 1: Kas masuk nyata */}
          <div
            style={{
              backgroundColor: '#FFFFFF',
              border: '1px solid #7A6E65',
              borderRadius: '16px',
              padding: '1.75rem',
              boxShadow: '0 2px 10px rgba(0,0,0,0.02)'
            }}
          >
            <div style={{ fontSize: '0.925rem', color: '#554A43', marginBottom: '1.25rem' }}>
              Kas masuk nyata
            </div>
            <div
              style={{
                fontSize: '2.25rem',
                fontWeight: 800,
                color: '#43281C',
                marginBottom: '1rem',
                lineHeight: 1
              }}
            >
              {kasMasukNyata !== null ? kasMasukNyata : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {kasMasukNyata !== null ? 'Total kas lunas yang benar-benar diterima' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 2: Simulasi sandbox */}
          <div
            style={{
              backgroundColor: '#FFFFFF',
              border: '1px solid #7A6E65',
              borderRadius: '16px',
              padding: '1.75rem',
              boxShadow: '0 2px 10px rgba(0,0,0,0.02)'
            }}
          >
            <div style={{ fontSize: '0.925rem', color: '#554A43', marginBottom: '1.25rem' }}>
              Simulasi sandbox
            </div>
            <div
              style={{
                fontSize: '2.25rem',
                fontWeight: 800,
                color: '#43281C',
                marginBottom: '1rem',
                lineHeight: 1
              }}
            >
              {simulasiSandbox !== null ? simulasiSandbox : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {simulasiSandbox !== null ? 'Pembayaran sandbox (tidak masuk kas nyata)' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 3: Status pending */}
          <div
            style={{
              backgroundColor: '#FFFFFF',
              border: '1px solid #7A6E65',
              borderRadius: '16px',
              padding: '1.75rem',
              boxShadow: '0 2px 10px rgba(0,0,0,0.02)'
            }}
          >
            <div style={{ fontSize: '0.925rem', color: '#554A43', marginBottom: '1.25rem' }}>
              Status pending
            </div>
            <div
              style={{
                fontSize: '2.25rem',
                fontWeight: 800,
                color: '#43281C',
                marginBottom: '1rem',
                lineHeight: 1
              }}
            >
              {statusPending !== null ? statusPending : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {statusPending !== null ? 'Menunggu penyelesaian pembayaran' : 'Menunggu data riil'}
            </div>
          </div>
        </div>

        {/* TABEL REKAP BULANAN SESUAI FIGMA */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            border: '1px solid #7A6E65',
            borderRadius: '16px',
            padding: '2rem',
            marginBottom: '1.5rem',
            boxShadow: '0 2px 10px rgba(0,0,0,0.02)'
          }}
        >
          <h2
            style={{
              fontSize: '1.25rem',
              fontWeight: 800,
              color: '#2D231E',
              marginBottom: '1.5rem'
            }}
          >
            Rekap bulanan
          </h2>

          {/* Tabel Sumber Data & Status */}
          <div style={{ width: '100%' }}>
            {/* Header Kolom */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                paddingBottom: '0.85rem',
                borderBottom: '1px solid #EFE8DF',
                fontSize: '0.95rem',
                fontWeight: 700,
                color: '#2D231E'
              }}
            >
              <span>Sumber data</span>
              <span>Status</span>
            </div>

            {/* Baris 1: Transaksi valid */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '1.15rem 0',
                borderBottom: '1px solid #F6F1EA',
                fontSize: '0.925rem',
                color: '#3B302A'
              }}
            >
              <span>Transaksi valid</span>
              <span style={{ color: '#554A43' }}>{statusTransaksiValid}</span>
            </div>

            {/* Baris 2: Sandbox */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '1.15rem 0',
                borderBottom: '1px solid #F6F1EA',
                fontSize: '0.925rem',
                color: '#3B302A'
              }}
            >
              <span>Sandbox</span>
              <span style={{ color: '#554A43' }}>{statusSandbox}</span>
            </div>

            {/* Baris 3: Data belum lengkap */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '1.15rem 0',
                fontSize: '0.925rem',
                color: '#3B302A'
              }}
            >
              <span>Data belum lengkap</span>
              <span style={{ color: '#554A43' }}>{statusDataBelumLengkap}</span>
            </div>
          </div>
        </div>

        {/* CATATAN FOOTNOTE MERAH SESUAI FIGMA */}
        <div
          style={{
            fontSize: '0.875rem',
            color: '#B91C1C',
            marginBottom: '1.5rem'
          }}
        >
          Tidak ada angka laporan yang dibuat-buat dalam desain ini.
        </div>

      </main>
    </div>
  );
}
