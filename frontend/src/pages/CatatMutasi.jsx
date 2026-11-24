import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, ArrowLeft, Calendar } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STOCK_KEY = 'rengginang_data_stok';
const HISTORY_KEY = 'rengginang_riwayat_mutasi';

// Data produk default untuk pilihan varian
const DEFAULT_PRODUCTS = [
  { id: 'udang', nama: 'Rengginang Udang', fisik: 34, reservasi: 2, tersedia: 32 },
  { id: 'bawang', nama: 'Rengginang Bawang', fisik: 40, reservasi: 2, tersedia: 38 },
  { id: 'cumi', nama: 'Rengginang Cumi', fisik: 14, reservasi: 2, tersedia: 12 },
  { id: 'terasi', nama: 'Rengginang Terasi', fisik: 40, reservasi: 2, tersedia: 38 }
];

export default function CatatMutasi() {
  const navigate = useNavigate();

  // State daftar stok dari localStorage
  const [daftarStok, setDaftarStok] = useState(() => {
    try {
      const saved = localStorage.getItem(STOCK_KEY);
      return saved ? JSON.parse(saved) : DEFAULT_PRODUCTS;
    } catch (e) {
      return DEFAULT_PRODUCTS;
    }
  });

  // State form catat mutasi sesuai mockup
  const [jenisMutasi, setJenisMutasi] = useState('masuk'); // 'masuk' | 'keluar' | 'penyesuaian'
  const [selectedProductId, setSelectedProductId] = useState('udang');
  const [jumlah, setJumlah] = useState(10);
  const [tanggal, setTanggal] = useState('16 November 2026');
  const [keterangan, setKeterangan] = useState('Hasil produksi selesai');
  const [errorPesan, setErrorPesan] = useState('');

  // Sesi Owner
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Ambil varian produk yang dipilih saat ini
  const produkDipilih =
    daftarStok.find((p) => p.id === selectedProductId) || daftarStok[0] || DEFAULT_PRODUCTS[0];

  // Hitung live stok sesudah perubahan
  const fisikSekarang = Number(produkDipilih.fisik || 0);
  const reservasi = Number(produkDipilih.reservasi || 0);
  const qty = Number(jumlah || 0);

  let fisikSesudah = fisikSekarang;
  if (jenisMutasi === 'masuk') {
    fisikSesudah = fisikSekarang + qty;
  } else if (jenisMutasi === 'keluar') {
    fisikSesudah = fisikSekarang - qty;
  } else if (jenisMutasi === 'penyesuaian') {
    fisikSesudah = qty;
  }

  const tersediaSesudah = fisikSesudah - reservasi;

  // Validasi & Simpan Perubahan Mutasi
  const handleSimpan = (e) => {
    e.preventDefault();
    setErrorPesan('');

    if (qty <= 0 && jenisMutasi !== 'penyesuaian') {
      setErrorPesan('Jumlah stok harus lebih dari 0.');
      return;
    }

    // Acceptance Criteria: Stok tersedia tidak boleh negatif
    if (tersediaSesudah < 0) {
      setErrorPesan('Stok tersedia tidak boleh negatif setelah perubahan.');
      return;
    }

    // Perbarui data stok
    const updatedStok = daftarStok.map((item) => {
      if (item.id === produkDipilih.id) {
        return {
          ...item,
          fisik: fisikSesudah,
          tersedia: tersediaSesudah
        };
      }
      return item;
    });

    setDaftarStok(updatedStok);
    localStorage.setItem(STOCK_KEY, JSON.stringify(updatedStok));

    // Catat ke riwayat mutasi
    const riwayatBaru = {
      id: Date.now().toString(),
      tanggal: tanggal,
      produk: produkDipilih.nama,
      jenis: jenisMutasi,
      jumlah: qty,
      fisikSebelum: fisikSekarang,
      fisikSesudah: fisikSesudah,
      keterangan: keterangan || '-'
    };

    try {
      const existingHistory = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
      existingHistory.unshift(riwayatBaru);
      localStorage.setItem(HISTORY_KEY, JSON.stringify(existingHistory));
    } catch (e) {
      console.warn('Gagal mencatat riwayat:', e);
    }

    // Arahkan ke halaman mutasi stok berhasil dicatat
    navigate('/admin/stok/tersimpan');
  };

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

  // Judul dinamis sesuai tab aktif
  let judulHalaman = 'Catat stok masuk';
  let subjudulHalaman = 'Tambahkan persediaan setelah produksi atau penerimaan barang.';
  let labelJumlah = 'Jumlah stok yang diterima';

  if (jenisMutasi === 'keluar') {
    judulHalaman = 'Catat stok keluar';
    subjudulHalaman = 'Catat pengeluaran barang tanpa memakai stok yang direservasi.';
    labelJumlah = 'Jumlah stok yang dikeluarkan';
  } else if (jenisMutasi === 'penyesuaian') {
    judulHalaman = 'Penyesuaian stok';
    subjudulHalaman = 'Sesuaikan jumlah persediaan dengan hasil opname fisik di gudang.';
    labelJumlah = 'Jumlah fisik aktual';
  }

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Kiri */}
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
                className="sidebar-item active"
                onClick={() => navigate('/admin/stok')}
              >
                <LayoutList size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
              >
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
          </ul>
        </div>

        {/* Footer Sidebar */}
        <div className="sidebar-footer">
          <Link to="/" className="sidebar-footer-link">
            Buka website publik
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

      {/* Konten Utama Form Catat Mutasi */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Data Stok / Catat Mutasi
        </div>

        {/* Header Halaman & Tombol Kembali */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">{judulHalaman}</h1>
            <p className="admin-page-subtitle">{subjudulHalaman}</p>
          </div>

          <Link to="/admin/stok" className="btn-back-stock">
            <ArrowLeft size={16} />
            <span>Data Stok</span>
          </Link>
        </div>

        {/* Grid 2 Kolom Sesuai Figma */}
        <div className="stock-mutation-layout">
          {/* Kolom Kiri: Form Catat Mutasi */}
          <div className="stock-table-card" style={{ padding: '2rem 2.25rem' }}>
            {errorPesan && (
              <div className="modal-alert-error" style={{ marginBottom: '1.5rem' }}>
                {errorPesan}
              </div>
            )}

            <form onSubmit={handleSimpan} className="modal-form">
              {/* Jenis Perubahan: 3 Pill Tombol */}
              <div className="form-group">
                <label className="form-label">Jenis perubahan</label>
                <div className="mutation-type-pills">
                  <button
                    type="button"
                    className={`mutation-pill-btn ${jenisMutasi === 'masuk' ? 'active' : ''}`}
                    onClick={() => {
                      setJenisMutasi('masuk');
                      setKeterangan('Hasil produksi selesai');
                    }}
                  >
                    Stok masuk
                  </button>
                  <button
                    type="button"
                    className={`mutation-pill-btn ${jenisMutasi === 'keluar' ? 'active' : ''}`}
                    onClick={() => {
                      setJenisMutasi('keluar');
                      setJumlah(5);
                      setKeterangan('Penjualan toko');
                    }}
                  >
                    Stok keluar
                  </button>
                  <button
                    type="button"
                    className={`mutation-pill-btn ${jenisMutasi === 'penyesuaian' ? 'active' : ''}`}
                    onClick={() => {
                      setJenisMutasi('penyesuaian');
                      setKeterangan('Penyesuaian hasil opname fisik');
                    }}
                  >
                    Penyesuaian
                  </button>
                </div>
              </div>

              {/* Pilih Varian Produk */}
              <div className="form-group" style={{ marginTop: '0.75rem' }}>
                <label className="form-label">Pilih varian produk</label>
                <select
                  className="form-input"
                  value={selectedProductId}
                  onChange={(e) => setSelectedProductId(e.target.value)}
                  style={{ backgroundColor: '#FFFFFF', cursor: 'pointer' }}
                >
                  {daftarStok.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.nama}
                    </option>
                  ))}
                </select>
              </div>

              {/* Baris 2 Kolom: Jumlah & Tanggal */}
              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: '1fr 1fr',
                  gap: '1.25rem',
                  marginTop: '0.75rem'
                }}
              >
                {/* Jumlah Stok */}
                <div className="form-group">
                  <label className="form-label">{labelJumlah}</label>
                  <div className="input-with-suffix">
                    <input
                      type="number"
                      className="form-input"
                      value={jumlah}
                      onChange={(e) => setJumlah(Number(e.target.value))}
                      min="1"
                      required
                    />
                    <span className="input-suffix-text">unit</span>
                  </div>
                </div>

                {/* Tanggal Pencatatan */}
                <div className="form-group">
                  <label className="form-label">Tanggal pencatatan</label>
                  <div className="input-with-suffix">
                    <input
                      type="text"
                      className="form-input"
                      value={tanggal}
                      onChange={(e) => setTanggal(e.target.value)}
                    />
                    <span className="input-suffix-text">
                      <Calendar size={15} />
                    </span>
                  </div>
                </div>
              </div>

              {/* Keterangan */}
              <div className="form-group" style={{ marginTop: '0.75rem' }}>
                <label className="form-label">Keterangan</label>
                <input
                  type="text"
                  className="form-input"
                  value={keterangan}
                  onChange={(e) => setKeterangan(e.target.value)}
                  placeholder="Contoh: Hasil produksi selesai"
                />
              </div>

              {/* Kotak Catatan Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FAF5EE',
                  border: '1px solid #EFE6D9',
                  borderRadius: '12px',
                  padding: '1rem 1.25rem',
                  marginTop: '1.25rem',
                  marginBottom: '1rem'
                }}
              >
                <div style={{ fontWeight: 800, fontSize: '0.875rem', color: '#43281C', marginBottom: '0.2rem' }}>
                  Catatan:
                </div>
                <div style={{ fontSize: '0.825rem', color: '#8E7D73' }}>
                  Setiap penyimpanan harus dicatat pada riwayat mutasi.
                </div>
              </div>

              {/* Tombol Aksi: Batal & Simpan */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.85rem', marginTop: '1.25rem' }}>
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => navigate('/admin/stok')}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-confirm">
                  Simpan
                </button>
              </div>
            </form>
          </div>

          {/* Kolom Kanan: Ringkasan Stok & Aturan Persediaan */}
          <div>
            {/* Kartu Ringkasan Stok */}
            <div className="stock-summary-card">
              <h2 className="rules-card-title" style={{ marginBottom: '0.2rem' }}>
                Ringkasan stok
              </h2>
              <div style={{ fontWeight: 800, color: '#43281C', fontSize: '1rem' }}>
                {produkDipilih.nama}
              </div>

              <div className="stock-summary-rows">
                <div className="stock-summary-row">
                  <span>Stok fisik sekarang</span>
                  <span className="stock-summary-val-bold">{fisikSekarang} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Stok reservasi</span>
                  <span className="stock-summary-val-bold">{reservasi} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Stok tersedia</span>
                  <span className="stock-summary-val-bold">{fisikSekarang - reservasi} unit</span>
                </div>
                <div className="stock-summary-row">
                  <span>Fisik sesudah perubahan</span>
                  <span className="stock-summary-val-green">{fisikSesudah} unit</span>
                </div>
              </div>
            </div>

            {/* Kartu Aturan Persediaan */}
            <div className="rules-card">
              <h3 className="rules-card-title">Aturan persediaan</h3>
              <div className="rules-list">
                <div>Tersedia = Fisik - Reservasi</div>
                <div>Stok tersedia tidak boleh negatif.</div>
                <div>Perubahan dicatat sebagai riwayat.</div>
              </div>
            </div>
          </div>
        </div>

        {/* Catatan Kaki Prototype */}
        <p className="stock-footnote">
          Nilai pada formulir adalah contoh; pengisian & validasi nyata dikerjakan saat implementasi.
        </p>
      </main>
    </div>
  );
}
