import React, { useEffect, useState, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ShoppingBag, Store, FileText } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

// Data pesanan default sesuai frame Figma PPL21_01_Daftar_Pesanan_Owner__HALAMAN
const DATA_PESANAN_FIGMA = [
  {
    nomor_pesanan: 'PJ-DEMO-0001',
    kontak: 'Pelanggan A',
    whatsapp: '0812-3456-7801',
    metode: 'Ambil Sendiri',
    pembayaran: 'Menunggu',
    status: 'Baru',
    alamat: 'Outlet Situbondo (Ambil Sendiri)',
    kurir: 'Ambil di Toko',
    subtotal: 56000,
    ongkir: 0,
    total: 56000,
    items: 'Udang x 2',
    tanggal: '2026-11-22'
  },
  {
    nomor_pesanan: 'PJ-DEMO-0002',
    kontak: 'Pelanggan B',
    whatsapp: '0812-3456-7802',
    metode: 'Kirim Kurir',
    pembayaran: 'Dibayar',
    status: 'Diproses',
    alamat: 'Jl. Melati No. 45, Situbondo',
    kurir: 'JNE REG (Tarif API)',
    subtotal: 84000,
    ongkir: 15000,
    total: 99000,
    items: 'Udang x 2 • Bawang x 1',
    tanggal: '2026-11-22'
  },
  {
    nomor_pesanan: 'PJ-DEMO-0003',
    kontak: 'Pelanggan C',
    whatsapp: '0812-3456-7803',
    metode: 'Kirim Kurir',
    pembayaran: 'Dibayar',
    status: 'Selesai',
    alamat: 'Jl. Mawar No. 12, Surabaya',
    kurir: 'J&T EZ (Tarif API)',
    subtotal: 112000,
    ongkir: 20000,
    total: 132000,
    items: 'Terasi x 2 • Cumi x 2',
    tanggal: '2026-11-21'
  },
  {
    nomor_pesanan: 'PJ-DEMO-0004',
    kontak: 'Pelanggan D',
    whatsapp: '0812-3456-7804',
    metode: 'Ambil Sendiri',
    pembayaran: 'Gagal',
    status: 'Batal',
    alamat: 'Outlet Situbondo (Ambil Sendiri)',
    kurir: 'Ambil di Toko',
    subtotal: 28000,
    ongkir: 0,
    total: 28000,
    items: 'Bawang x 1',
    tanggal: '2026-11-20'
  }
];

export default function DaftarPesananOwner() {
  const navigate = useNavigate();

  // State autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // State data pesanan & loading
  const [daftarPesanan, setDaftarPesanan] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [errorPesan, setErrorPesan] = useState('');

  // State filter & pencarian sesuai Figma
  const [kataKunci, setKataKunci] = useState('');
  const [filterStatus, setFilterStatus] = useState('Semua status');
  const [rentangWaktu, setRentangWaktu] = useState('Semua tanggal');

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
          // Fallback lokal jika backend offline
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

  // Muat data pesanan dari backend
  useEffect(() => {
    const muatPesanan = async () => {
      setIsLoading(true);
      setErrorPesan('');

      try {
        // Coba endpoint khusus daftar pesanan admin
        const resAdmin = await fetch(`${API_URL}/api/admin/pesanan`, {
          credentials: 'include'
        });

        if (resAdmin.ok) {
          const jsonAdmin = await resAdmin.json();
          if (Array.isArray(jsonAdmin.data) && jsonAdmin.data.length > 0) {
            setDaftarPesanan(jsonAdmin.data);
            setIsLoading(false);
            return;
          }
        }

        // Cek data pesanan riil dari reservasi / order aktif di backend
        const resReservasi = await fetch(`${API_URL}/api/pesanan/reservasi`);
        let pesananRiil = [];

        if (resReservasi.ok) {
          const jsonRes = await resReservasi.json();
          if (Array.isArray(jsonRes.data) && jsonRes.data.length > 0) {
            // Ambil nomor pesanan unik
            const uniqueOrders = [...new Set(jsonRes.data.map((r) => r.nomor_pesanan))];

            // Fetch detail setiap pesanan
            const promises = uniqueOrders.map(async (nomor) => {
              try {
                const resDetail = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${nomor}`);
                if (resDetail.ok) {
                  const j = await resDetail.json();
                  const d = j.data || {};
                  const detail = d.pesanan_detail || {};
                  return {
                    nomor_pesanan: d.nomor_pesanan || nomor,
                    kontak: d.penerima || detail.kontak?.nama_penerima || 'Pelanggan Online',
                    whatsapp: detail.kontak?.nomor_whatsapp || '-',
                    metode: detail.tipe_checkout === 'pickup' ? 'Ambil Sendiri' : 'Kirim Kurir',
                    pembayaran: d.status_pesanan === 'dibayar' ? 'Dibayar' : 'Menunggu',
                    status: d.status_pesanan === 'dibayar' ? 'Diproses' : 'Baru',
                    alamat: d.alamat || detail.kontak?.alamat_lengkap || 'Outlet Situbondo',
                    kurir: d.kurir || detail.kurir?.nama_kurir || 'Kurir Reguler',
                    subtotal: detail.snapshot_subtotal || d.total_pembayaran || 0,
                    ongkir: detail.snapshot_ongkir || 0,
                    total: d.total_pembayaran || detail.snapshot_total || 0,
                    items: jsonRes.data
                      .filter((r) => r.nomor_pesanan === nomor)
                      .map((r) => `${r.nama_varian} x ${r.jumlah_direservasi}`)
                      .join(' • ') || 'Rengginang',
                    tanggal: d.waktu_dibuat ? d.waktu_dibuat.split('T')[0] : '2026-11-22'
                  };
                }
              } catch (e) {
                return null;
              }
              return null;
            });

            const hasilRiil = await Promise.all(promises);
            pesananRiil = hasilRiil.filter(Boolean);
          }
        }

        // Gabungkan pesanan riil dengan data Figma (hindari duplikasi nomor pesanan)
        const combined = [...pesananRiil];
        DATA_PESANAN_FIGMA.forEach((itemFigma) => {
          if (!combined.some((c) => c.nomor_pesanan === itemFigma.nomor_pesanan)) {
            combined.push(itemFigma);
          }
        });

        setDaftarPesanan(combined);
      } catch (err) {
        console.warn('Gagal memuat pesanan backend, memuat data default:', err);
        setDaftarPesanan(DATA_PESANAN_FIGMA);
      } finally {
        setIsLoading(false);
      }
    };

    muatPesanan();
  }, []);

  // Filter pencarian, status, dan rentang tanggal
  const pesananTerfilter = useMemo(() => {
    return daftarPesanan.filter((p) => {
      // 1. Pencarian nomor pesanan atau nama kontak
      if (kataKunci.trim() !== '') {
        const cari = kataKunci.toLowerCase();
        const cocokNomor = p.nomor_pesanan?.toLowerCase().includes(cari);
        const cocokKontak = p.kontak?.toLowerCase().includes(cari);
        if (!cocokNomor && !cocokKontak) return false;
      }

      // 2. Filter Status
      if (filterStatus !== 'Semua status') {
        const statusMatch = p.status?.toLowerCase() === filterStatus.toLowerCase();
        const bayarMatch = p.pembayaran?.toLowerCase() === filterStatus.toLowerCase();
        if (!statusMatch && !bayarMatch) return false;
      }

      // 3. Rentang Waktu
      if (rentangWaktu !== 'Semua tanggal') {
        const today = new Date().toISOString().split('T')[0];
        if (rentangWaktu === 'Hari ini' && p.tanggal !== today) {
          return false;
        }
      }

      return true;
    });
  }, [daftarPesanan, kataKunci, filterStatus, rentangWaktu]);

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

  // Helper tampilan badge
  const renderBadgePembayaran = (pembayaran) => {
    const val = (pembayaran || '').toLowerCase();
    if (val.includes('dibayar')) {
      return <span className="order-badge-pill order-badge-dibayar">Dibayar</span>;
    }
    if (val.includes('gagal')) {
      return <span className="order-badge-pill order-badge-gagal">Gagal</span>;
    }
    return <span className="order-badge-pill order-badge-menunggu">Menunggu</span>;
  };

  const renderBadgeStatus = (status) => {
    const val = (status || '').toLowerCase();
    if (val.includes('diproses')) {
      return <span className="order-badge-pill order-badge-diproses">Diproses</span>;
    }
    if (val.includes('selesai')) {
      return <span className="order-badge-pill order-badge-selesai">Selesai</span>;
    }
    if (val.includes('batal')) {
      return <span className="order-badge-pill order-badge-batal">Batal</span>;
    }
    return <span className="order-badge-pill order-badge-baru">Baru</span>;
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

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Owner Sesuai Figma */}
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
              {/* Menu Pesanan Aktif Sesuai Frame PPL21_01 */}
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

      {/* Konten Utama: Daftar Pesanan */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Daftar Pesanan
        </div>

        {/* Judul & Deskripsi Halaman */}
        <div style={{ marginBottom: '1.75rem' }}>
          <h1 className="admin-page-title">Daftar Pesanan</h1>
          <p className="admin-page-subtitle">
            Lihat dan kelola pesanan online sesuai statusnya.
          </p>
        </div>

        {/* Baris Filter & Pencarian Sesuai Figma */}
        <div className="orders-filter-row">
          <div className="orders-filter-group">
            <label className="orders-filter-label" htmlFor="input-cari">
              Cari pesanan
            </label>
            <input
              id="input-cari"
              type="text"
              className="orders-filter-input"
              placeholder="Cari nomor atau nama"
              value={kataKunci}
              onChange={(e) => setKataKunci(e.target.value)}
            />
          </div>

          <div className="orders-filter-group">
            <label className="orders-filter-label" htmlFor="select-status">
              Filter status
            </label>
            <select
              id="select-status"
              className="orders-filter-select"
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
            >
              <option value="Semua status">Semua status</option>
              <option value="Baru">Baru</option>
              <option value="Diproses">Diproses</option>
              <option value="Selesai">Selesai</option>
              <option value="Batal">Batal</option>
              <option value="Menunggu">Menunggu Pembayaran</option>
              <option value="Dibayar">Sudah Dibayar</option>
            </select>
          </div>

          <div className="orders-filter-group">
            <label className="orders-filter-label" htmlFor="select-rentang">
              Rentang waktu
            </label>
            <select
              id="select-rentang"
              className="orders-filter-select"
              value={rentangWaktu}
              onChange={(e) => setRentangWaktu(e.target.value)}
            >
              <option value="Semua tanggal">Semua tanggal</option>
              <option value="Hari ini">Hari ini</option>
              <option value="7 hari terakhir">7 hari terakhir</option>
              <option value="30 hari terakhir">30 hari terakhir</option>
            </select>
          </div>
        </div>

        {/* Card Tabel Pesanan */}
        <div className="orders-table-card">
          <div className="table-responsive">
            <table className="orders-table">
              <thead>
                <tr>
                  <th>No. pesanan</th>
                  <th>Kontak</th>
                  <th>Metode</th>
                  <th>Pembayaran</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td colSpan="5" style={{ textAlign: 'center', padding: '2rem', color: '#8E7D73' }}>
                      Memuat daftar pesanan...
                    </td>
                  </tr>
                ) : pesananTerfilter.length === 0 ? (
                  <tr>
                    <td colSpan="5" style={{ textAlign: 'center', padding: '2.5rem', color: '#8E7D73' }}>
                      Tidak ada pesanan yang sesuai dengan kata kunci atau filter status.
                    </td>
                  </tr>
                ) : (
                  pesananTerfilter.map((p) => (
                    <tr
                      key={p.nomor_pesanan}
                      className="order-row-clickable"
                      onClick={() => navigate(`/admin/pesanan/${p.nomor_pesanan}`)}
                      title="Klik untuk melihat detail pesanan"
                    >
                      <td className="order-cell-id">{p.nomor_pesanan}</td>
                      <td>{p.kontak}</td>
                      <td>{p.metode}</td>
                      <td>{renderBadgePembayaran(p.pembayaran)}</td>
                      <td>{renderBadgeStatus(p.status)}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Catatan Footnote Sesuai Figma */}
        <p className="table-footnote" style={{ marginTop: '1rem', color: '#8E7D73', fontSize: '0.85rem' }}>
          Klik baris pesanan untuk melihat rincian. Nomor DEMO bukan pesanan riil.
        </p>
      </main>
    </div>
  );
}
