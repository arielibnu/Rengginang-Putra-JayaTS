import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Circle,
  Square,
  LayoutList,
  Diamond,
  ShoppingBag,
  Store,
  FileText,
  AlertCircle,
  CheckCircle2,
  Calendar,
  RefreshCw,
  Info
} from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

const DAFTAR_BULAN = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
];

export default function LaporanStokAkhir() {
  const navigate = useNavigate();

  // Autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(true);
  const [isCheckingAuth, setIsCheckingAuth] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Filter Periode Sesuai Mockup Figma PPL26_01_Laporan_Stok_Akhir__TAB
  const [bulanTerpilih, setBulanTerpilih] = useState('November');
  const [tahunTerpilih, setTahunTerpilih] = useState('2026');

  // State Data Ringkasan Stok
  const [stokFisik, setStokFisik] = useState(null);
  const [reservasiAktif, setReservasiAktif] = useState(null);
  const [stokTersedia, setStokTersedia] = useState(null);

  // Status Laporan & Keterangan Sumber Data
  const [statusTransaksiValid, setStatusTransaksiValid] = useState('Mengikuti backend');
  const [statusSandbox, setStatusSandbox] = useState('Tidak dihitung sebagai kas nyata');
  const [statusDataBelumLengkap, setStatusDataBelumLengkap] = useState('Ditandai dan dijelaskan');

  // State Loading & Kebutuhan API Backend
  const [isLoading, setIsLoading] = useState(false);
  const [apiTersedia, setApiTersedia] = useState(false);
  const [catatanIntegrasi, setCatatanIntegrasi] = useState('');

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

  // Fungsi memuat data laporan stok dari backend
  const muatLaporanStok = async (bulan, tahun) => {
    setIsLoading(true);
    setCatatanIntegrasi('');

    try {
      // 1. Coba panggil endpoint khusus laporan stok akhir periodik (WBS 1.4.2.1)
      const res = await fetch(
        `${API_URL}/api/admin/laporan/stok-akhir?bulan=${encodeURIComponent(bulan)}&tahun=${encodeURIComponent(tahun)}`,
        { credentials: 'include' }
      );

      if (res.ok) {
        const json = await res.json();
        const data = json.data || {};
        setStokFisik(data.stok_fisik ?? null);
        setReservasiAktif(data.reservasi_aktif ?? null);
        setStokTersedia(data.stok_tersedia ?? (data.stok_fisik - data.reservasi_aktif));
        setApiTersedia(true);
        setStatusTransaksiValid('Data resmi terverifikasi backend');
        setIsLoading(false);
        return;
      }
    } catch (e) {
      console.warn('Endpoint khusus laporan stok bulanan belum tersedia di backend:', e);
    }

    // 2. Jika endpoint laporan bulanan belum tersedia, jangan buat data palsu!
    // Tampilkan kondisi "Menunggu data riil" persis seperti desain Figma PPL26_01.
    setStokFisik(null);
    setReservasiAktif(null);
    setStokTersedia(null);
    setApiTersedia(false);
    setStatusTransaksiValid('Mengikuti backend');
    setStatusSandbox('Tidak dihitung sebagai kas nyata');
    setStatusDataBelumLengkap('Ditandai dan dijelaskan');
    setCatatanIntegrasi(
      'Endpoint laporan riwayat penutupan stok bulanan (GET /api/admin/laporan/stok-akhir) belum diimplementasikan di backend (WBS 1.4.2.1). Sesuai ketentuan PPL-26, tidak ada angka dummy yang dibuat-buat, tampilan menampilkan tanda "—" (Menunggu data riil).'
    );
    setIsLoading(false);
  };

  // Muat data saat halaman pertama kali dibuka
  useEffect(() => {
    muatLaporanStok(bulanTerpilih, tahunTerpilih);
  }, []);

  // Handler tombol "Terapkan"
  const handleTerapkan = (e) => {
    e.preventDefault();
    muatLaporanStok(bulanTerpilih, tahunTerpilih);
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
                <Store size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              {/* Menu Laporan Aktif Sesuai Frame Figma PPL26_01 */}
              <div
                className="sidebar-item active"
                onClick={() => navigate('/admin/laporan')}
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

      {/* KONTEN UTAMA: LAPORAN STOK AKHIR */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Laporan Stok Akhir
        </div>

        {/* Tab Navigasi Antar Laporan (Laporan Stok & Laporan Pengeluaran) */}
        <div
          style={{
            display: 'flex',
            gap: '1rem',
            marginBottom: '1.5rem',
            borderBottom: '1px solid #EFE8DF',
            paddingBottom: '0.5rem'
          }}
        >
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
            Laporan Stok Akhir
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

        {/* Judul & Subtitle Sesuai Figma */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="admin-page-title">Laporan Stok Akhir</h1>
          <p className="admin-page-subtitle">
            Tampilkan data berdasarkan bulan dan tahun. Simulasi dipisahkan dari transaksi nyata.
          </p>
        </div>

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

        {/* 3 KARTU RINGKASAN STOK SESUAI FIGMA */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
            gap: '1.25rem',
            marginBottom: '2.5rem'
          }}
        >
          {/* Kartu 1: Stok fisik */}
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
              Stok fisik
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
              {stokFisik !== null ? stokFisik : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {stokFisik !== null ? 'Jumlah tercatat di database' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 2: Reservasi aktif */}
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
              Reservasi aktif
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
              {reservasiAktif !== null ? reservasiAktif : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {reservasiAktif !== null ? 'Stok ditahan pesanan aktif' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 3: Stok tersedia */}
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
              Stok tersedia
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
              {stokTersedia !== null ? stokTersedia : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {stokTersedia !== null ? 'Stok fisik - reservasi aktif' : 'Menunggu data riil'}
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
          {/* Judul Bagian Rekap */}
          <h2
            style={{
              fontSize: '1.25rem',
              fontWeight: 800,
              color: '#2D231E',
              margin: '0 0 1.5rem 0'
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
