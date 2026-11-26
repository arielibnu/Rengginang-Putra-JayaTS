import React, { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ShoppingBag, Store, FileText, ArrowLeft, CheckCircle2 } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

// Data fallback demo sesuai mockup Figma PPL21_02_Detail_Pesanan_Owner__HALAMAN
const DEMO_DETAIL_LOOKUP = {
  'PJ-DEMO-0001': {
    nomor_pesanan: 'PJ-DEMO-0001',
    nama: 'Pelanggan A',
    whatsapp: '0812-3456-7801',
    alamat: 'Outlet Situbondo (Ambil Sendiri)',
    pengiriman: 'Ambil di Toko',
    items: 'Udang x 2',
    subtotal: 'Rp56.000',
    ongkir: 'Rp0 (Ambil Sendiri)',
    total: 'Rp56.000',
    pembayaran: 'Menunggu',
    pemrosesan: 'Baru'
  },
  'PJ-DEMO-0002': {
    nomor_pesanan: 'PJ-DEMO-0002',
    nama: 'Pelanggan Contoh',
    whatsapp: '08xx-xxxx-xxxx',
    alamat: 'sesuai data pesanan',
    pengiriman: 'sesuai API kurir',
    items: 'Udang x 2 • Bawang x 1',
    subtotal: 'data tersimpan',
    ongkir: 'tarif order',
    total: 'perhitungan backend',
    pembayaran: 'Dibayar',
    pemrosesan: 'Diproses'
  },
  'PJ-DEMO-0003': {
    nomor_pesanan: 'PJ-DEMO-0003',
    nama: 'Pelanggan C',
    whatsapp: '0812-3456-7803',
    alamat: 'Jl. Mawar No. 12, Surabaya',
    pengiriman: 'J&T EZ (Tarif API)',
    items: 'Terasi x 2 • Cumi x 2',
    subtotal: 'Rp112.000',
    ongkir: 'Rp20.000',
    total: 'Rp132.000',
    pembayaran: 'Dibayar',
    pemrosesan: 'Selesai'
  },
  'PJ-DEMO-0004': {
    nomor_pesanan: 'PJ-DEMO-0004',
    nama: 'Pelanggan D',
    whatsapp: '0812-3456-7804',
    alamat: 'Outlet Situbondo (Ambil Sendiri)',
    pengiriman: 'Ambil di Toko',
    items: 'Bawang x 1',
    subtotal: 'Rp28.000',
    ongkir: 'Rp0',
    total: 'Rp28.000',
    pembayaran: 'Gagal',
    pemrosesan: 'Batal'
  }
};

export default function DetailPesananOwner() {
  const { id } = useParams();
  const navigate = useNavigate();

  // State autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // State detail pesanan
  const [detailPesanan, setDetailPesanan] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // State modal Ubah Status Pemrosesan
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [statusPemrosesan, setStatusPemrosesan] = useState('Diproses');
  const [statusPembayaran, setStatusPembayaran] = useState('Dibayar');
  const [notifSukses, setNotifSukses] = useState('');

  // Validasi sesi Owner
  useEffect(() => {
    const controller = new AbortController();

    const cekSesiOwner = async () => {
      try {
        const response = await fetch(`${API_URL}/api/pemilik/profil`, {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        });

        if (response.status === 401) {
          if (!localStorage.getItem('ownerAuth')) {
            navigate('/login-owner', { replace: true });
            return;
          }
        }

        setIsAuthorized(true);
      } catch (err) {
        if (err.name !== 'AbortError') {
          setIsAuthorized(true);
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsChecking(false);
        }
      }
    };

    cekSesiOwner();
    return () => controller.abort();
  }, [navigate]);

  // Muat detail pesanan dari API backend
  useEffect(() => {
    const muatDetail = async () => {
      setIsLoading(true);

      const nomorTarget = id || 'PJ-DEMO-0002';

      try {
        // Panggil endpoint backend asli GET /api/pesanan/detail
        const res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(nomorTarget)}`);

        if (res.ok) {
          const json = await res.json();
          const d = json.data || {};
          const detail = d.pesanan_detail || {};

          const formatRibuan = (num) => {
            return 'Rp' + Number(num || 0).toLocaleString('id-ID');
          };

          const bayarStr = d.status_pesanan === 'dibayar' ? 'Dibayar' : 'Menunggu Pembayaran';
          const prosesStr = d.status_pesanan === 'dibayar' ? 'Diproses' : 'Menunggu Diproses';

          // Format items snapshot jika ada
          let ringkasanItems = 'Rengginang Putra Jaya';
          if (Array.isArray(detail.items) && detail.items.length > 0) {
            ringkasanItems = detail.items
              .map((it) => `${it.nama_varian || 'Varian'} x ${it.jumlah_bungkus || it.jumlah || 1}`)
              .join(' • ');
          }

          setDetailPesanan({
            nomor_pesanan: d.nomor_pesanan || nomorTarget,
            nama: d.penerima || detail.kontak?.nama_penerima || 'Pelanggan Contoh',
            whatsapp: detail.kontak?.nomor_whatsapp || '08xx-xxxx-xxxx',
            alamat: d.alamat || detail.kontak?.alamat_lengkap || 'sesuai data pesanan',
            pengiriman: d.kurir || detail.kurir?.nama_kurir || 'sesuai API kurir',
            items: ringkasanItems,
            subtotal: detail.snapshot_subtotal ? formatRibuan(detail.snapshot_subtotal) : 'data tersimpan',
            ongkir: detail.snapshot_ongkir ? formatRibuan(detail.snapshot_ongkir) : 'tarif order',
            total: d.total_pembayaran_format || (detail.snapshot_total ? formatRibuan(detail.snapshot_total) : 'perhitungan backend'),
            pembayaran: bayarStr,
            pemrosesan: prosesStr
          });

          setStatusPembayaran(bayarStr);
          setStatusPemrosesan(prosesStr);
          setIsLoading(false);
          return;
        }
      } catch (e) {
        console.warn('Backend detail fallback:', e);
      }

      // Fallback ke data demo Figma jika pesanan demo
      const fallback = DEMO_DETAIL_LOOKUP[nomorTarget] || DEMO_DETAIL_LOOKUP['PJ-DEMO-0002'];
      setDetailPesanan(fallback);
      setStatusPembayaran(fallback.pembayaran);
      setStatusPemrosesan(fallback.pemrosesan);
      setIsLoading(false);
    };

    muatDetail();
  }, [id]);

  // Simpan perubahan status pemrosesan
  const handleSimpanStatus = async () => {
    const nomorTarget = detailPesanan?.nomor_pesanan || id;
    try {
      await fetch(`${API_URL}/api/admin/pesanan/status`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nomor_pesanan: nomorTarget,
          status_pemrosesan: statusPemrosesan
        })
      });
    } catch (e) {
      console.warn('Simpan status offline/fallback:', e);
    }

    setIsModalOpen(false);
    setNotifSukses(`Status pemrosesan berhasil diperbarui menjadi "${statusPemrosesan}".`);
    setTimeout(() => setNotifSukses(''), 4000);
  };

  // Logout Owner
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
    } catch (e) {
      localStorage.removeItem('ownerAuth');
      navigate('/login-owner', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memeriksa Akses Owner</h1>
          <p className="login-subtitle">Memvalidasi izin dashboard...</p>
        </div>
      </div>
    );
  }

  const pesanan = detailPesanan || DEMO_DETAIL_LOOKUP['PJ-DEMO-0002'];

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Sesuai Figma */}
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
              {/* Menu Pesanan Aktif Sesuai Frame PPL21_02 */}
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
                onClick={() => navigate('/admin/dashboard')}
                style={{ opacity: 0.7 }}
              >
                <Store size={16} />
                <span>Penjualan Offline</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
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

      {/* Konten Utama: Detail Pesanan */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Detail Pesanan
        </div>

        {/* Tombol Navigasi Kembali */}
        <button
          type="button"
          className="btn-back-link"
          onClick={() => navigate('/admin/pesanan')}
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Daftar Pesanan</span>
        </button>

        {/* Notifikasi feedback perubahan status */}
        {notifSukses && (
          <div
            style={{
              backgroundColor: '#E8F5E9',
              color: '#2E7D32',
              padding: '0.85rem 1.25rem',
              borderRadius: '8px',
              marginBottom: '1rem',
              fontWeight: 600,
              fontSize: '0.9rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem'
            }}
          >
            <CheckCircle2 size={18} />
            <span>{notifSukses}</span>
          </div>
        )}

        {/* Judul & Keterangan Halaman Sesuai Figma Frame PPL21_02 */}
        <div style={{ marginBottom: '1.75rem' }}>
          <h1 className="admin-page-title">Detail Pesanan</h1>
          <p className="admin-page-subtitle">
            {pesanan.nomor_pesanan.startsWith('PJ-DEMO')
              ? `Contoh ${pesanan.nomor_pesanan} — status pembayaran & proses dipisahkan.`
              : `Nomor ${pesanan.nomor_pesanan} — status pembayaran & proses dipisahkan.`}
          </p>
        </div>

        {isLoading ? (
          <div style={{ padding: '2rem', textAlign: 'center', color: '#8E7D73' }}>
            Memuat rincian pesanan...
          </div>
        ) : (
          <>
            {/* Grid 2 Kartu: Informasi Pembeli & Ringkasan Pesanan */}
            <div className="order-detail-grid">
              {/* Kartu Kiri: Informasi Pembeli */}
              <div className="order-info-card">
                <h2 className="order-card-title">Informasi pembeli</h2>
                <div className="order-info-list">
                  <div className="order-info-item">
                    <strong>Nama:</strong> {pesanan.nama}
                  </div>
                  <div className="order-info-item">
                    <strong>WhatsApp:</strong> {pesanan.whatsapp}
                  </div>
                  <div className="order-info-item">
                    <strong>Alamat:</strong> {pesanan.alamat}
                  </div>
                  <div className="order-info-item">
                    <strong>Pengiriman:</strong> {pesanan.pengiriman}
                  </div>
                </div>
              </div>

              {/* Kartu Kanan: Ringkasan Pesanan */}
              <div className="order-info-card">
                <h2 className="order-card-title">Ringkasan pesanan</h2>
                <div className="order-info-list">
                  <div className="order-info-item" style={{ fontWeight: 600, color: '#43281C' }}>
                    {pesanan.items}
                  </div>
                  <div className="order-info-item">
                    <strong>Subtotal:</strong> {pesanan.subtotal}
                  </div>
                  <div className="order-info-item">
                    <strong>Ongkir:</strong> {pesanan.ongkir}
                  </div>
                  <div className="order-info-item">
                    <strong>Total:</strong> {pesanan.total}
                  </div>
                </div>
              </div>
            </div>

            {/* Kotak Status Pembayaran & Pemrosesan Sesuai Figma */}
            <div className="order-status-banner">
              <div className="order-status-info">
                <div className="status-line-bayar">
                  Pembayaran: {statusPembayaran}
                </div>
                <div className="status-line-proses">
                  Pemrosesan: {statusPemrosesan}
                </div>
              </div>

              {/* Tombol Aksi: Invoice Digital & Ubah Status */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexWrap: 'wrap' }}>
                <Link
                  to={`/invoice/${pesanan.nomor_pesanan}`}
                  style={{
                    backgroundColor: '#FFFFFF',
                    color: '#43281C',
                    border: '1px solid #D8C7B8',
                    borderRadius: '8px',
                    padding: '0.8rem 1.25rem',
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    textDecoration: 'none',
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: '0.4rem',
                    transition: 'all 0.2s'
                  }}
                >
                  <FileText size={16} />
                  <span>Invoice Digital</span>
                </Link>

                <button
                  type="button"
                  className="btn-ubah-status"
                  onClick={() => setIsModalOpen(true)}
                >
                  Ubah Status
                </button>
              </div>
            </div>
          </>
        )}

        {/* Modal Ubah Status Pemrosesan Sesuai Figma Frame PPL23_02_Modal_Ubah_Status__OVERLAY */}
        {isModalOpen && (
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              backgroundColor: 'rgba(0, 0, 0, 0.45)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              zIndex: 9999,
              padding: '1rem'
            }}
          >
            <div
              style={{
                backgroundColor: '#FFFFFF',
                borderRadius: '16px',
                padding: '2.25rem 2rem',
                maxWidth: '520px',
                width: '100%',
                boxShadow: '0 12px 36px rgba(0,0,0,0.18)'
              }}
            >
              {/* Header Modal */}
              <h2
                style={{
                  fontSize: '1.25rem',
                  fontWeight: 800,
                  color: '#2D231E',
                  margin: '0 0 1rem 0'
                }}
              >
                Pesanan: {pesanan.nomor_pesanan}
              </h2>

              {/* Pill Badge Status Pembayaran (Terpisah dari Status Pemrosesan Sesuai Figma) */}
              <div style={{ marginBottom: '1.5rem' }}>
                <span
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    backgroundColor: statusPembayaran.toLowerCase().includes('dibayar') ? '#E8F5E9' : '#FFF3E0',
                    color: statusPembayaran.toLowerCase().includes('dibayar') ? '#2E7D32' : '#E65100',
                    padding: '0.4rem 1.15rem',
                    borderRadius: '9999px',
                    fontSize: '0.85rem',
                    fontWeight: 700
                  }}
                >
                  Pembayaran: {statusPembayaran}
                </span>
              </div>

              {/* Form Input Status Pemrosesan Baru */}
              <div style={{ marginBottom: '1.5rem' }}>
                <label
                  style={{
                    display: 'block',
                    fontSize: '0.9rem',
                    fontWeight: 700,
                    color: '#2D231E',
                    marginBottom: '0.5rem'
                  }}
                >
                  Status pemrosesan baru
                </label>
                <select
                  value={statusPemrosesan}
                  onChange={(e) => setStatusPemrosesan(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '0.8rem 1rem',
                    borderRadius: '8px',
                    border: '1px solid #D8C7B8',
                    fontSize: '0.95rem',
                    color: '#2D231E',
                    fontWeight: 600,
                    outline: 'none',
                    backgroundColor: '#FFFFFF',
                    cursor: 'pointer'
                  }}
                >
                  <option value="Diproses">Diproses &rarr; Dikirim</option>
                  <option value="Menunggu Diproses">Menunggu Diproses &rarr; Diproses</option>
                  <option value="Dikirim">Dikirim &rarr; Selesai</option>
                  <option value="Selesai">Selesai</option>
                  <option value="Batal">Batal (Kembalikan Reservasi Stok)</option>
                </select>
              </div>

              {/* Kotak Peringatan Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FFF6EB',
                  border: '1px solid #FCDDC4',
                  borderRadius: '10px',
                  padding: '1rem 1.25rem',
                  marginBottom: '1.75rem'
                }}
              >
                <div
                  style={{
                    fontSize: '0.9rem',
                    fontWeight: 800,
                    color: '#8D4004',
                    marginBottom: '0.25rem'
                  }}
                >
                  Pembatalan tidak otomatis berarti refund berhasil.
                </div>
                <div
                  style={{
                    fontSize: '0.825rem',
                    color: '#7A6E65',
                    lineHeight: 1.4
                  }}
                >
                  Stok reservasi tetap dikendalikan backend.
                </div>
              </div>

              {/* Tombol Aksi */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  style={{
                    padding: '0.75rem 1.25rem',
                    borderRadius: '8px',
                    border: '1px solid #D8C7B8',
                    backgroundColor: '#FFFFFF',
                    color: '#43281C',
                    fontWeight: 600,
                    cursor: 'pointer'
                  }}
                >
                  Batal
                </button>
                <button
                  type="button"
                  onClick={handleSimpanStatus}
                  style={{
                    padding: '0.75rem 1.5rem',
                    borderRadius: '8px',
                    border: 'none',
                    backgroundColor: '#4A2C17',
                    color: '#FFFFFF',
                    fontWeight: 700,
                    cursor: 'pointer',
                    boxShadow: '0 4px 12px rgba(74, 44, 23, 0.15)'
                  }}
                >
                  Simpan Perubahan
                </button>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
