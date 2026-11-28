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
  AlertCircle,
  Download,
  Printer,
  Info,
  CheckCircle2
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

export default function LaporanEstimasiKeuntungan() {
  const navigate = useNavigate();

  // Autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(true);
  const [isCheckingAuth, setIsCheckingAuth] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Filter Periode Sesuai Mockup Figma PPL-29
  const [bulanTerpilih, setBulanTerpilih] = useState('November');
  const [tahunTerpilih, setTahunTerpilih] = useState('2026');

  // State Data 3 Kartu Ringkasan Sesuai Desain Figma:
  // 1. Penerimaan nyata
  // 2. Modal barang terjual (HPP)
  // 3. Biaya operasional
  const [penerimaanNyata, setPenerimaanNyata] = useState(null);
  const [modalBarangTerjual, setModalBarangTerjual] = useState(null);
  const [biayaOperasional, setBiayaOperasional] = useState(null);
  const [estimasiLabaBersih, setEstimasiLabaBersih] = useState(null);

  // Status Laporan & Keterangan Sumber Data Sesuai Figma
  const [statusTransaksiValid, setStatusTransaksiValid] = useState('Mengikuti backend');
  const [statusSandbox, setStatusSandbox] = useState('Tidak dihitung sebagai kas nyata');
  const [statusDataBelumLengkap, setStatusDataBelumLengkap] = useState('Ditandai dan dijelaskan');

  // State Loading & Notifikasi
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

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

  // Fungsi memuat data laporan estimasi keuntungan dari backend
  const muatLaporanEstimasiKeuntungan = async (bulan, tahun) => {
    setIsLoading(true);
    setErrorMessage('');
    setSuccessMessage('');

    try {
      // Hubungi endpoint laporan estimasi keuntungan jika sudah diimplementasikan di backend
      const res = await fetch(
        `${API_URL}/api/admin/laporan/estimasi-keuntungan?bulan=${encodeURIComponent(bulan)}&tahun=${encodeURIComponent(tahun)}`,
        { credentials: 'include' }
      );

      if (res.ok) {
        const json = await res.json();
        const data = json.data || {};

        // Format nilai mata uang jika berupa angka
        setPenerimaanNyata(
          typeof data.penerimaan_nyata === 'number'
            ? `Rp${data.penerimaan_nyata.toLocaleString('id-ID')}`
            : (data.penerimaan_nyata ?? null)
        );
        setModalBarangTerjual(
          typeof data.modal_barang_terjual === 'number'
            ? `Rp${data.modal_barang_terjual.toLocaleString('id-ID')}`
            : (data.modal_barang_terjual ?? null)
        );
        setBiayaOperasional(
          typeof data.biaya_operasional === 'number'
            ? `Rp${data.biaya_operasional.toLocaleString('id-ID')}`
            : (data.biaya_operasional ?? null)
        );
        setEstimasiLabaBersih(
          typeof data.estimasi_laba_bersih === 'number'
            ? `Rp${data.estimasi_laba_bersih.toLocaleString('id-ID')}`
            : (data.estimasi_laba_bersih ?? null)
        );

        setStatusTransaksiValid('Mengikuti backend');
        setStatusSandbox('Tidak dihitung sebagai kas nyata');
        setStatusDataBelumLengkap('Ditandai dan dijelaskan');
        setIsLoading(false);
        return;
      }
    } catch (e) {
      console.warn('Endpoint khusus estimasi keuntungan belum tersedia di backend:', e);
    }

    // Jika API belum tersedia atau terjadi error, tampilkan kondisi menunggu data riil persis desain Figma
    // JANGAN membuat data angka dummy atau data sandbox sebagai laporan keuntungan nyata!
    setPenerimaanNyata(null);
    setModalBarangTerjual(null);
    setBiayaOperasional(null);
    setEstimasiLabaBersih(null);
    setStatusTransaksiValid('Mengikuti backend');
    setStatusSandbox('Tidak dihitung sebagai kas nyata');
    setStatusDataBelumLengkap('Ditandai dan dijelaskan');
    setIsLoading(false);
  };

  // Muat data saat halaman pertama kali dibuka
  useEffect(() => {
    muatLaporanEstimasiKeuntungan(bulanTerpilih, tahunTerpilih);
  }, []);

  // Handler tombol "Terapkan"
  const handleTerapkan = (e) => {
    e.preventDefault();
    muatLaporanEstimasiKeuntungan(bulanTerpilih, tahunTerpilih);
  };

  // Handler ekspor CSV/Excel (WBS 1.4.5.2)
  const handleExportCSV = () => {
    const csvContent =
      'data:text/csv;charset=utf-8,' +
      [
        ['Laporan Estimasi Keuntungan Bulanan - Rengginang Putra Jaya'],
        [`Periode: ${bulanTerpilih} ${tahunTerpilih}`],
        [''],
        ['Metrik', 'Nilai', 'Keterangan'],
        ['Penerimaan Nyata', penerimaanNyata || 'Menunggu data riil', 'Total penjualan riil valid non-sandbox'],
        ['Modal Barang Terjual (HPP)', modalBarangTerjual || 'Menunggu data riil', 'HPP rengginang terjual'],
        ['Biaya Operasional', biayaOperasional || 'Menunggu data riil', 'Biaya operasional bulanan'],
        ['Estimasi Laba Bersih', estimasiLabaBersih || 'Menunggu data riil', 'Rumus: Penerimaan - HPP - Biaya'],
        [''],
        ['Audit Sumber Data', 'Status'],
        ['Transaksi Valid', statusTransaksiValid],
        ['Sandbox Midtrans', statusSandbox],
        ['Data Belum Lengkap', statusDataBelumLengkap],
        [''],
        ['Catatan: Transaksi sandbox dan ongkir titipan tidak dihitung ke laba bersih usaha.']
      ]
        .map((e) => e.map((val) => `"${val}"`).join(','))
        .join('\n');

    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `Estimasi_Keuntungan_${bulanTerpilih}_${tahunTerpilih}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    setSuccessMessage('Data laporan berhasil diekspor ke format CSV / Excel!');
    setTimeout(() => setSuccessMessage(''), 4000);
  };

  // Handler cetak / unduh PDF (WBS 1.4.5.2)
  const handlePrintPDF = () => {
    window.print();
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
              <div className="sidebar-brand-badge">Dashboard Owner</div>
            </div>
          </div>

          {/* Menu Navigasi Sidebar */}
          <ul className="sidebar-menu">
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
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
                <ShoppingBag size={16} />
                <span>Kelola Produk</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/stok')}
              >
                <Layers size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/modal-biaya')}
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
                <TrendingDown size={16} />
                <span>Pesanan</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/penjualan-offline')}
              >
                <DollarSign size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              {/* Menu Laporan Aktif Sesuai Frame Figma */}
              <div
                className="sidebar-item active"
                onClick={() => navigate('/admin/laporan/estimasi-keuntungan')}
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

      {/* KONTEN UTAMA: ESTIMASI KEUNTUNGAN */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Estimasi Keuntungan
        </div>

        {/* Tab Navigasi Antar Laporan */}
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
            onClick={() => navigate('/admin/laporan/penerimaan-kas')}
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
            Estimasi Keuntungan
          </button>
        </div>

        {/* Judul & Subtitle Sesuai Figma */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="admin-page-title">Estimasi Keuntungan</h1>
          <p className="admin-page-subtitle">
            Tampilkan data berdasarkan bulan dan tahun. Simulasi dipisahkan dari transaksi nyata.
          </p>
        </div>

        {/* Pesan Notifikasi & Error */}
        {errorMessage && (
          <div
            style={{
              backgroundColor: '#FEE2E2',
              border: '1px solid #FCA5A5',
              borderRadius: '10px',
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

        {successMessage && (
          <div
            style={{
              backgroundColor: '#ECFDF5',
              border: '1px solid #A7F3D0',
              borderRadius: '10px',
              padding: '1rem 1.25rem',
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              color: '#065F46',
              fontSize: '0.9rem'
            }}
          >
            <CheckCircle2 size={18} />
            <span>{successMessage}</span>
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
            marginBottom: '2rem',
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
                type="number"
                value={tahunTerpilih}
                onChange={(e) => setTahunTerpilih(e.target.value)}
                min="2020"
                max="2035"
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

          {/* Tombol Terapkan & Tombol Ekspor (Kanan Sesuai Desain Figma & WBS 1.4.5.2) */}
          <div
            style={{
              marginLeft: 'auto',
              display: 'flex',
              gap: '0.75rem',
              alignItems: 'center',
              flexWrap: 'wrap'
            }}
          >
            {/* Tombol Ekspor CSV/Excel (WBS 1.4.5.2) */}
            <button
              type="button"
              onClick={handleExportCSV}
              style={{
                padding: '0.8rem 1.25rem',
                borderRadius: '10px',
                border: '1px solid #7A6E65',
                backgroundColor: '#FFFFFF',
                color: '#43281C',
                fontSize: '0.9rem',
                fontWeight: 600,
                cursor: 'pointer',
                display: 'inline-flex',
                alignItems: 'center',
                gap: '0.5rem',
                height: '46px',
                transition: 'all 0.2s'
              }}
              title="Unduh data dalam format CSV/Excel"
            >
              <Download size={16} />
              <span>Ekspor Excel/CSV</span>
            </button>

            {/* Tombol Cetak/PDF (WBS 1.4.5.2) */}
            <button
              type="button"
              onClick={handlePrintPDF}
              style={{
                padding: '0.8rem 1.25rem',
                borderRadius: '10px',
                border: '1px solid #7A6E65',
                backgroundColor: '#FFFFFF',
                color: '#43281C',
                fontSize: '0.9rem',
                fontWeight: 600,
                cursor: 'pointer',
                display: 'inline-flex',
                alignItems: 'center',
                gap: '0.5rem',
                height: '46px',
                transition: 'all 0.2s'
              }}
              title="Cetak atau simpan sebagai PDF"
            >
              <Printer size={16} />
              <span>Cetak PDF</span>
            </button>

            {/* Tombol Terapkan Utama (Sesuai Figma) */}
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

        {/* 3 KARTU RINGKASAN ESTIMASI KEUNTUNGAN SESUAI FIGMA PPL-29 */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
            gap: '1.25rem',
            marginBottom: '2rem'
          }}
        >
          {/* Kartu 1: Penerimaan nyata */}
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
              Penerimaan nyata
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
              {penerimaanNyata !== null ? penerimaanNyata : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {penerimaanNyata !== null ? 'Total penerimaan dari transaksi riil valid' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 2: Modal barang terjual (HPP) */}
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
              Modal barang terjual
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
              {modalBarangTerjual !== null ? modalBarangTerjual : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {modalBarangTerjual !== null ? 'Total HPP rengginang yang terjual' : 'Menunggu data riil'}
            </div>
          </div>

          {/* Kartu 3: Biaya operasional */}
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
              Biaya operasional
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
              {biayaOperasional !== null ? biayaOperasional : '—'}
            </div>
            <div style={{ fontSize: '0.85rem', color: '#7A6E65' }}>
              {biayaOperasional !== null ? 'Beban operasional usaha pada periode ini' : 'Menunggu data riil'}
            </div>
          </div>
        </div>

        {/* KARTU SOROTAN ESTIMASI LABA BERSIH (WBS 1.4.5.2) */}
        <div
          style={{
            backgroundColor: '#FAF5EE',
            border: '2px solid #43281C',
            borderRadius: '16px',
            padding: '1.5rem 2rem',
            marginBottom: '2rem',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            flexWrap: 'wrap',
            gap: '1rem'
          }}
        >
          <div>
            <div style={{ fontSize: '1rem', fontWeight: 700, color: '#43281C', marginBottom: '0.25rem' }}>
              Estimasi Laba Bersih Usaha
            </div>
            <div style={{ fontSize: '0.875rem', color: '#7A6E65' }}>
              Rumus: Penjualan Riil - Modal Barang Terjual (HPP) - Biaya Operasional (Pengecualian Ongkir & Sandbox)
            </div>
          </div>
          <div
            style={{
              fontSize: '2.25rem',
              fontWeight: 800,
              color: '#43281C',
              lineHeight: 1
            }}
          >
            {estimasiLabaBersih !== null ? estimasiLabaBersih : '—'}
          </div>
        </div>

        {/* TRANSPARANSI RUMUS & INDIKATOR KELENGKAPAN DATA (WBS 1.4.5.3) */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            border: '1px solid #EFE8DF',
            borderRadius: '12px',
            padding: '1.25rem 1.5rem',
            marginBottom: '2rem',
            display: 'flex',
            gap: '1rem',
            alignItems: 'flex-start'
          }}
        >
          <Info size={20} style={{ color: '#7A6E65', marginTop: '2px', flexShrink: 0 }} />
          <div style={{ fontSize: '0.875rem', color: '#554A43', lineHeight: 1.6 }}>
            <strong>Transparansi Audit Perhitungan Laba (WBS 1.4.5.1 & 1.4.5.3):</strong>
            <ul style={{ margin: '0.35rem 0 0 1.25rem', padding: 0 }}>
              <li><strong>Estimasi Laba Bersih</strong> = Penjualan Nyata - Modal Terjual (HPP) - Biaya Operasional.</li>
              <li><strong>Pengecualian Transaksi Sandbox</strong>: Pembayaran Midtrans Sandbox tidak pernah diakui sebagai keuntungan riil.</li>
              <li><strong>Pengecualian Ongkir Titipan</strong>: Ongkos kirim ekspedisi titipan langsung diteruskan ke kurir dan tidak dihitung ke laba usaha.</li>
              <li><strong>Indikator Kelengkapan Biaya</strong>: Jika komponen biaya/modal belum tercatat lengkap di database, sistem menampilkan status <em>Menunggu data riil</em> agar terhindar dari bias keputusan.</li>
            </ul>
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
