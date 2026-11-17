import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, Plus } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_KEY = 'rengginang_data_stok';

// Data awal persediaan stok 4 varian sesuai mockup resmi Figma PPL-04 & PPL-05
const INITIAL_STOCK_DATA = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    fisik: 34,
    reservasi: 2,
    tersedia: 32,
    minimum: 15
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    fisik: 14,
    reservasi: 2,
    tersedia: 12,
    minimum: 15
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    fisik: 40,
    reservasi: 2,
    tersedia: 38,
    minimum: 15
  }
];

export default function DataStok() {
  const navigate = useNavigate();

  // State daftar stok varian
  const [daftarStok, setDaftarStok] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (!saved) return INITIAL_STOCK_DATA;
      const parsed = JSON.parse(saved);

      // Pastikan setiap varian memiliki field minimum
      return INITIAL_STOCK_DATA.map((item) => {
        const found = parsed.find((p) => p.id === item.id);
        if (!found) return item;
        return {
          ...item,
          fisik: Number(found.fisik ?? item.fisik),
          reservasi: Number(found.reservasi ?? item.reservasi),
          tersedia: Number(found.tersedia ?? item.tersedia),
          minimum: Number(found.minimum ?? item.minimum ?? 15)
        };
      });
    } catch (e) {
      return INITIAL_STOCK_DATA;
    }
  });

  // State autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // State Modal Atur Batas Minimum (PPL-05)
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedVariantId, setSelectedVariantId] = useState('cumi');
  const [inputMinimum, setInputMinimum] = useState('15');

  // Simpan data stok ke localStorage setiap kali ada perubahan
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(daftarStok));
    } catch (e) {
      console.warn('Gagal menyimpan ke localStorage:', e);
    }
  }, [daftarStok]);

  // Validasi sesi login Owner
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
          navigate('/login-owner', { replace: true });
          return;
        }

        setIsAuthorized(true);
      } catch (err) {
        // Jika backend belum aktif, tetap izinkan pratinjau frontend mandiri
        setIsAuthorized(true);
      } finally {
        setIsChecking(false);
      }
    };

    cekSesiOwner();
    return () => controller.abort();
  }, [navigate]);

  // Fungsi logout
  const handleLogout = async () => {
    if (isLoggingOut) return;
    setIsLoggingOut(true);

    try {
      await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });
    } catch (e) {
      console.log('Logout lokal...');
    } finally {
      localStorage.removeItem('ownerAuth');
      navigate('/login-owner');
    }
  };

  // Evaluasi varian yang stoknya menipis (stok tersedia <= batas minimum)
  const varianMenipis = daftarStok.filter(
    (item) => Number(item.tersedia || 0) <= Number(item.minimum || 15)
  );

  // Ambil data varian yang sedang dipilih di modal
  const varianTerpilih =
    daftarStok.find((item) => item.id === selectedVariantId) || daftarStok[0];

  // Buka modal dan inisialisasi batas minimum
  const handleOpenModal = () => {
    const cumiItem = daftarStok.find((v) => v.id === 'cumi') || daftarStok[0];
    setSelectedVariantId(cumiItem.id);
    setInputMinimum(String(cumiItem.minimum || 15));
    setIsModalOpen(true);
  };

  // Saat ganti varian di dropdown modal
  const handleSelectVariantChange = (e) => {
    const newId = e.target.value;
    setSelectedVariantId(newId);
    const targetVariant = daftarStok.find((v) => v.id === newId);
    if (targetVariant) {
      setInputMinimum(String(targetVariant.minimum || 15));
    }
  };

  // Simpan batas minimum baru
  const handleSaveMinimum = (e) => {
    e.preventDefault();
    const parsedMin = parseInt(inputMinimum, 10);
    const validMin = isNaN(parsedMin) || parsedMin < 0 ? 0 : parsedMin;

    setDaftarStok((prev) =>
      prev.map((item) =>
        item.id === selectedVariantId ? { ...item, minimum: validMin } : item
      )
    );
    setIsModalOpen(false);
  };

  // Status varian terpilih di dalam modal berdasarkan input yang sedang diketik
  const currentPreviewMin = parseInt(inputMinimum, 10) || 0;
  const isSelectedVariantLow =
    Number(varianTerpilih.tersedia || 0) <= currentPreviewMin;

  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memeriksa akses Owner</h1>
          <p className="login-subtitle">Memvalidasi sesi login melalui backend...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Kiri Sesuai Figma */}
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

          {/* Daftar Menu Admin */}
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
              <div className="sidebar-item active">
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
          </ul>
        </div>

        {/* Menu Footer Sidebar */}
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

      {/* Konten Utama Persediaan Stok */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Peringatan Minimum
        </div>

        {/* Header Halaman & Tombol Aksi */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">Persediaan stok</h1>
            <p className="admin-page-subtitle">
              Pantau stok tersedia serta batas minimum setiap varian.
            </p>
          </div>

          <div className="stock-top-actions">
            <button
              type="button"
              className="btn-catat-mutasi"
              onClick={() => navigate('/admin/stok/mutasi')}
            >
              <Plus size={16} />
              <span>Catat mutasi</span>
            </button>

            <button
              type="button"
              className="btn-atur-minimum"
              onClick={handleOpenModal}
            >
              Atur batas minimum
            </button>
          </div>
        </div>

        {/* Banner Peringatan Stok Menipis (PPL-05) Sesuai Figma */}
        {varianMenipis.length > 0 && (
          <div className="stock-warning-banner" style={{ marginBottom: '2rem' }}>
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
                        `${v.nama} tersedia ${v.tersedia} bungkus, sedangkan batas minimum ${v.minimum || 15} bungkus.`
                    )
                    .join(' • ')}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Tabel Daftar Persediaan Sesuai Figma */}
        <section>
          <h2
            style={{
              fontSize: '1.35rem',
              fontWeight: 800,
              color: '#43281C',
              marginBottom: '1.25rem'
            }}
          >
            Daftar persediaan
          </h2>

          <div className="stock-table-card" style={{ padding: '0.5rem 1.75rem', marginBottom: '1.5rem' }}>
            <table className="stock-table">
              <thead>
                <tr>
                  <th style={{ width: '35%', textAlign: 'left' }}>Varian produk</th>
                  <th style={{ textAlign: 'center', width: '13%' }}>Fisik</th>
                  <th style={{ textAlign: 'center', width: '13%' }}>Reservasi</th>
                  <th style={{ textAlign: 'center', width: '13%' }}>Tersedia</th>
                  <th style={{ textAlign: 'center', width: '13%' }}>Minimum</th>
                  <th style={{ textAlign: 'center', width: '13%' }}>Status</th>
                </tr>
              </thead>
              <tbody>
                {daftarStok.map((item) => {
                  const isMenipis =
                    Number(item.tersedia || 0) <= Number(item.minimum || 15);
                  return (
                    <tr
                      key={item.id}
                      className={isMenipis ? 'row-highlight-low' : ''}
                    >
                      <td className="stock-item-name">{item.nama}</td>
                      <td style={{ textAlign: 'center' }}>{item.fisik}</td>
                      <td style={{ textAlign: 'center' }}>{item.reservasi}</td>
                      <td
                        style={{ textAlign: 'center' }}
                        className={
                          isMenipis ? 'stock-num-highlight-low' : 'stock-num-green'
                        }
                      >
                        {item.tersedia}
                      </td>
                      <td style={{ textAlign: 'center' }}>{item.minimum || 15}</td>
                      <td style={{ textAlign: 'center' }}>
                        {isMenipis ? (
                          <span className="badge-status-menipis">Stok menipis</span>
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
        </section>

        {/* Kotak Aturan Peringatan Otomatis Sesuai Figma */}
        <div className="rule-auto-warning-box">
          <div className="rule-auto-warning-title">Aturan peringatan otomatis</div>
          <div className="rule-auto-warning-desc">
            Peringatan muncul saat stok tersedia kurang dari atau sama dengan batas minimum per varian.
          </div>
        </div>

        {/* Catatan Kaki Prototype */}
        <p className="stock-footnote" style={{ marginTop: '1.25rem' }}>
          Catatan: data contoh untuk rancangan UI. Nilai aktual mengikuti API/backend PPL-04.
        </p>
      </main>

      {/* Modal Atur Batas Minimum (PPL-05) Sesuai Mockup Figma */}
      {isModalOpen && (
        <div className="modal-overlay" onClick={() => setIsModalOpen(false)}>
          <div
            className="modal-dialog"
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
          >
            {/* Header Modal */}
            <div className="modal-header-row">
              <div className="modal-icon-badge">
                <LayoutList size={22} />
              </div>
              <div>
                <h2 className="modal-title-text">Atur batas minimum</h2>
                <p className="modal-subtitle-text">
                  Pengaturan batas minimum per varian produk
                </p>
              </div>
            </div>

            {/* Form Pengaturan */}
            <form onSubmit={handleSaveMinimum}>
              <div className="modal-form-group">
                <label className="modal-form-label">Varian rengginang</label>
                <select
                  className="modal-select"
                  value={selectedVariantId}
                  onChange={handleSelectVariantChange}
                >
                  {daftarStok.map((v) => (
                    <option key={v.id} value={v.id}>
                      {v.nama}
                    </option>
                  ))}
                </select>
              </div>

              <div className="modal-form-group" style={{ marginTop: '1rem' }}>
                <label className="modal-form-label">
                  Batas minimum stok (bungkus)
                </label>
                <input
                  type="number"
                  min="0"
                  className="modal-input"
                  value={inputMinimum}
                  onChange={(e) => setInputMinimum(e.target.value)}
                  required
                />
              </div>

              {/* Box Informasi Stok Saat Ini */}
              <div className="modal-stock-preview-box">
                <div className="modal-stock-preview-title">
                  Informasi stok saat ini
                </div>
                <div className="modal-stock-preview-row">
                  <span>Stok tersedia</span>
                  <span
                    className={`modal-stock-preview-val ${
                      isSelectedVariantLow ? 'low' : 'safe'
                    }`}
                  >
                    {varianTerpilih.tersedia} bungkus
                  </span>
                </div>
                <div className="modal-stock-preview-row">
                  <span>Status</span>
                  {isSelectedVariantLow ? (
                    <span className="badge-status-menipis">Stok menipis</span>
                  ) : (
                    <span className="badge-status-normal">Normal</span>
                  )}
                </div>
              </div>

              <div className="modal-helper-text">
                Peringatan muncul otomatis saat stok tersedia &le; minimum.
              </div>

              {/* Tombol Batal & Simpan */}
              <div className="modal-actions-row">
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => setIsModalOpen(false)}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-submit">
                  Simpan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
