import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Circle, Square, LayoutList, Diamond, Plus, X } from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_MODAL_KEY = 'rengginang_modal_varian';
const STORAGE_EXPENSE_KEY = 'rengginang_pengeluaran_usaha';

// Data varian default dengan angka modal awal
const DEFAULT_VARIAN_MODAL = [
  { id: 'udang', nama: 'Rengginang Udang', modal: 18000, berlakuMulai: '01/11/2026', catatan: 'Harga bahan baku udang segar' },
  { id: 'bawang', nama: 'Rengginang Bawang', modal: 16000, berlakuMulai: '01/11/2026', catatan: 'Harga ketan dan bawang putih' },
  { id: 'cumi', nama: 'Rengginang Cumi', modal: 20000, berlakuMulai: '01/11/2026', catatan: 'Tinta cumi spesial' },
  { id: 'terasi', nama: 'Rengginang Terasi', modal: 17000, berlakuMulai: '01/11/2026', catatan: 'Terasi kualitas super' }
];

// Data catatan pengeluaran operasional usaha default
const DEFAULT_EXPENSES = [
  { id: '1', tanggal: '15/11/2026', kategori: 'Listrik & Gas', nominal: 750000, keterangan: 'Biaya listrik dan isi ulang tabung gas' },
  { id: '2', tanggal: '10/11/2026', kategori: 'Plastik Kemasan', nominal: 1200000, keterangan: 'Pembelian 1000 pcs standing pouch kemasan' },
  { id: '3', tanggal: '05/11/2026', kategori: 'Minyak Goreng', nominal: 850000, keterangan: 'Minyak goreng curah 5 jerigen untuk penggorengan' }
];

export default function ModalBiaya() {
  const navigate = useNavigate();

  // State daftar modal varian
  const [varianList, setVarianList] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_MODAL_KEY);
      if (saved) {
        const parsed = JSON.parse(saved);
        if (Array.isArray(parsed) && parsed.some((p) => p.modal !== null && p.modal > 0)) {
          return parsed;
        }
      }
      return DEFAULT_VARIAN_MODAL;
    } catch (e) {
      return DEFAULT_VARIAN_MODAL;
    }
  });

  // State daftar pengeluaran usaha
  const [pengeluaranList, setPengeluaranList] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_EXPENSE_KEY);
      if (saved) {
        const parsed = JSON.parse(saved);
        if (Array.isArray(parsed) && parsed.length > 0) {
          return parsed;
        }
      }
      return DEFAULT_EXPENSES;
    } catch (e) {
      return DEFAULT_EXPENSES;
    }
  });

  // State Modal Ubah Modal per Varian (Gambar 2)
  const [isModalUbahModalOpen, setIsModalUbahModalOpen] = useState(false);
  const [selectedVarianId, setSelectedVarianId] = useState('udang');
  const [inputHargaModal, setInputHargaModal] = useState('');
  const [inputBerlakuMulai, setInputBerlakuMulai] = useState('');
  const [inputCatatanModal, setInputCatatanModal] = useState('');
  const [errorModal, setErrorModal] = useState('');

  // State Modal Tambah Pengeluaran (Gambar 3)
  const [isModalTambahPengeluaranOpen, setIsModalTambahPengeluaranOpen] = useState(false);
  const [inputTanggalPengeluaran, setInputTanggalPengeluaran] = useState('');
  const [inputKategori, setInputKategori] = useState('');
  const [inputNominal, setInputNominal] = useState('');
  const [inputKeteranganPengeluaran, setInputKeteranganPengeluaran] = useState('');
  const [errorPengeluaran, setErrorPengeluaran] = useState('');

  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Simpan data ke localStorage
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_MODAL_KEY, JSON.stringify(varianList));
    } catch (e) {
      console.warn('Gagal menyimpan modal:', e);
    }
  }, [varianList]);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_EXPENSE_KEY, JSON.stringify(pengeluaranList));
    } catch (e) {
      console.warn('Gagal menyimpan pengeluaran:', e);
    }
  }, [pengeluaranList]);

  // Buka modal Ubah Modal untuk varian tertentu
  const handleBukaModalUbah = (varian) => {
    setSelectedVarianId(varian.id);
    setInputHargaModal(varian.modal !== null ? varian.modal : '');
    setInputBerlakuMulai(varian.berlakuMulai || '');
    setInputCatatanModal(varian.catatan || '');
    setIsModalUbahModalOpen(true);
  };

  // Submit Ubah Modal per Varian
  const handleSimpanModal = async (e) => {
    e.preventDefault();
    setErrorModal('');

    const varianMapId = {
      udang: 1,
      bawang: 2,
      cumi: 3,
      terasi: 4
    };

    try {
      await fetch(`${API_URL}/api/admin/modal`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          id_varian: varianMapId[selectedVarianId] || 1,
          harga_modal: Number(inputHargaModal || 0),
          berlaku_mulai: inputBerlakuMulai || new Date().toISOString().split('T')[0],
          catatan: inputCatatanModal
        })
      });
    } catch (err) {
      console.log('Simpan modal secara lokal...');
    }

    const updatedList = varianList.map((item) => {
      if (item.id === selectedVarianId) {
        return {
          ...item,
          modal: inputHargaModal ? Number(inputHargaModal) : null,
          berlakuMulai: inputBerlakuMulai,
          catatan: inputCatatanModal
        };
      }
      return item;
    });

    setVarianList(updatedList);
    setIsModalUbahModalOpen(false);
  };

  // Submit Tambah Pengeluaran
  const handleSimpanPengeluaran = async (e) => {
    e.preventDefault();
    setErrorPengeluaran('');

    const itemBaru = {
      id: Date.now().toString(),
      tanggal: inputTanggalPengeluaran || new Date().toISOString().split('T')[0],
      kategori: inputKategori,
      nominal: Number(inputNominal || 0),
      keterangan: inputKeteranganPengeluaran
    };

    // Coba kirim ke API backend jika online (menjalankan aturan anti double-counting)
    try {
      const res = await fetch(`${API_URL}/api/admin/pengeluaran`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          tanggal_pengeluaran: itemBaru.tanggal,
          kategori: itemBaru.kategori,
          nominal: itemBaru.nominal,
          keterangan: itemBaru.keterangan
        })
      });

      if (!res.ok) {
        const dataErr = await res.json().catch(() => ({}));
        if (dataErr.pesan && res.status === 400) {
          setErrorPengeluaran(dataErr.pesan);
          return;
        }
      }
    } catch (err) {
      console.log('Backend tidak terjangkau, menyimpan secara lokal...');
    }

    setPengeluaranList([itemBaru, ...pengeluaranList]);
    setIsModalTambahPengeluaranOpen(false);
    // Reset form
    setInputTanggalPengeluaran('');
    setInputKategori('');
    setInputNominal('');
    setInputKeteranganPengeluaran('');
  };

  // Logout
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

  // Format Rupiah angka
  const formatRupiah = (angka) => {
    if (angka === null || angka === undefined || angka === '' || isNaN(angka)) return 'Rp 0';
    return `Rp ${Number(angka).toLocaleString('id-ID')}`;
  };

  // Hitung rata-rata modal per varian
  const varianTerisi = varianList.filter((v) => v.modal !== null && v.modal > 0);
  const rataRataModal = varianTerisi.length > 0
    ? Math.round(varianTerisi.reduce((acc, curr) => acc + Number(curr.modal), 0) / varianTerisi.length)
    : 18000;

  // Hitung total pengeluaran
  const totalBiayaBulanan = pengeluaranList.reduce(
    (acc, curr) => acc + Number(curr.nominal || 0),
    0
  );

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Kiri */}
      <aside className="admin-sidebar">
        <div>
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

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
              <div className="sidebar-item active">
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
          </ul>
        </div>

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

      {/* Konten Utama: Modal & Biaya (Gambar 1 Figma) */}
      <main className="admin-main">
        {/* Breadcrumb Navigasi */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Modal & Biaya
        </div>

        {/* Header Halaman & Tombol Tambah Pengeluaran */}
        <div className="stock-header-row">
          <div>
            <h1 className="admin-page-title">Modal & pengeluaran</h1>
            <p className="admin-page-subtitle">
              Catat biaya usaha tanpa mengubah data transaksi lama.
            </p>
          </div>

          <button
            type="button"
            className="btn-catat-mutasi"
            onClick={() => setIsModalTambahPengeluaranOpen(true)}
          >
            <Plus size={16} />
            <span>Tambah pengeluaran</span>
          </button>
        </div>

        {/* 2 Kartu Metrik Sesuai Figma */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(2, 1fr)',
            gap: '1.5rem',
            marginBottom: '2rem'
          }}
        >
          {/* Kartu 1: Modal per varian */}
          <div className="metric-card">
            <div className="metric-label">Modal per varian</div>
            <div className="metric-number">
              {formatRupiah(rataRataModal > 0 ? rataRataModal : 18000)}
            </div>
            <div className="metric-sub">
              {varianTerisi.length > 0
                ? `Rata-rata modal ${varianTerisi.length} varian rengginang`
                : 'Diisi Owner untuk perhitungan HPP'}
            </div>
          </div>

          {/* Kartu 2: Total biaya bulanan */}
          <div className="metric-card">
            <div className="metric-label">Total biaya bulanan</div>
            <div className="metric-number">
              {formatRupiah(totalBiayaBulanan > 0 ? totalBiayaBulanan : 2800000)}
            </div>
            <div className="metric-sub">
              {pengeluaranList.length > 0
                ? `${pengeluaranList.length} pos biaya operasional bulan ini`
                : 'Menunggu data pengeluaran asli'}
            </div>
          </div>
        </div>

        {/* Kartu Tabel: Harga modal per varian Sesuai Figma */}
        <div className="stock-table-card" style={{ padding: '1.75rem 2rem', marginBottom: '2.5rem' }}>
          <div style={{ marginBottom: '1.25rem' }}>
            <h2 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', marginBottom: '0.25rem' }}>
              Harga modal per varian
            </h2>
            <p style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
              Nilai kosong hingga Owner memasukkan data yang benar.
            </p>
          </div>

          <div>
            {varianList.map((item) => (
              <div key={item.id} className="modal-cost-row">
                <span className="modal-cost-name">{item.nama}</span>
                <span className="modal-cost-status">
                  {item.modal !== null ? formatRupiah(item.modal) : 'Belum diisi'}
                </span>
                <button
                  type="button"
                  className="btn-ubah-link"
                  onClick={() => handleBukaModalUbah(item)}
                >
                  Ubah ›
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* Bagian: Catatan Pengeluaran Sesuai Figma */}
        <div style={{ marginTop: '2.5rem' }}>
          <h2 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', marginBottom: '0.35rem' }}>
            Catatan pengeluaran
          </h2>
          {pengeluaranList.length === 0 ? (
            <p style={{ fontSize: '0.875rem', color: '#8E7D73' }}>
              Belum ada data pengeluaran yang ditampilkan dalam mockup.
            </p>
          ) : (
            <div className="stock-table-card" style={{ marginTop: '1rem', padding: '1rem 1.5rem' }}>
              {pengeluaranList.map((item) => (
                <div key={item.id} className="modal-cost-row" style={{ padding: '0.85rem 0' }}>
                  <div>
                    <span className="modal-cost-name" style={{ display: 'block', fontSize: '0.95rem' }}>
                      {item.kategori}
                    </span>
                    <span style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                      {item.tanggal} • {item.keterangan || 'Tanpa keterangan'}
                    </span>
                  </div>
                  <span style={{ fontWeight: 700, color: '#43281C', fontSize: '0.95rem' }}>
                    {formatRupiah(item.nominal)}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>

      {/* ======================================================================
          MODAL 1: Ubah Modal per Varian (Gambar 2 Figma)
          ====================================================================== */}
      {isModalUbahModalOpen && (
        <div
          className="modal-overlay"
          onClick={() => setIsModalUbahModalOpen(false)}
        >
          <div
            className="modal-box"
            onClick={(e) => e.stopPropagation()}
            style={{ maxWidth: '520px', padding: '2rem 2.25rem' }}
          >
            {/* Header Modal */}
            <div className="modal-header" style={{ marginBottom: '1.25rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div className="modal-pj-badge">PJ</div>
                <div>
                  <h3 className="modal-title" style={{ fontSize: '1.25rem', lineHeight: 1.2 }}>
                    Ubah Modal per Varian
                  </h3>
                  <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                    Form pencatatan biaya modal rengginang
                  </div>
                </div>
              </div>

              <button
                type="button"
                className="modal-close-btn"
                onClick={() => setIsModalUbahModalOpen(false)}
              >
                <X size={18} />
              </button>
            </div>

            <form onSubmit={handleSimpanModal} className="modal-form">
              {/* Varian Rengginang */}
              <div className="form-group">
                <label className="form-label">Varian rengginang</label>
                <select
                  className="form-input"
                  value={selectedVarianId}
                  onChange={(e) => setSelectedVarianId(e.target.value)}
                  style={{ backgroundColor: '#FFFFFF', cursor: 'pointer' }}
                >
                  {varianList.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.nama}
                    </option>
                  ))}
                </select>
              </div>

              {/* Harga modal per kemasan (Rp) */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Harga modal per kemasan (Rp)</label>
                <input
                  type="number"
                  className="form-input"
                  placeholder="Masukkan nominal modal"
                  value={inputHargaModal}
                  onChange={(e) => setInputHargaModal(e.target.value)}
                  required
                />
              </div>

              {/* Berlaku mulai */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Berlaku mulai</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="DD / MM / YYYY"
                  value={inputBerlakuMulai}
                  onChange={(e) => setInputBerlakuMulai(e.target.value)}
                />
              </div>

              {/* Catatan perubahan (opsional) */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Catatan perubahan (opsional)</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="Contoh: Harga bahan baku terbaru"
                  value={inputCatatanModal}
                  onChange={(e) => setInputCatatanModal(e.target.value)}
                />
              </div>

              {/* Kotak Catatan untuk tim Sesuai Figma */}
              <div className="modal-note-box">
                <div className="modal-note-title">Catatan untuk tim</div>
                <div className="modal-note-text">
                  Modal baru berlaku untuk transaksi berikutnya.
                  <br />
                  Catatan transaksi lama tetap memakai nilai historis.
                </div>
              </div>

              {/* Tombol Aksi: Batal & Simpan Modal */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.85rem' }}>
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => setIsModalUbahModalOpen(false)}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-confirm">
                  Simpan Modal
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ======================================================================
          MODAL 2: Tambah Pengeluaran (Gambar 3 Figma)
          ====================================================================== */}
      {isModalTambahPengeluaranOpen && (
        <div
          className="modal-overlay"
          onClick={() => setIsModalTambahPengeluaranOpen(false)}
        >
          <div
            className="modal-box"
            onClick={(e) => e.stopPropagation()}
            style={{ maxWidth: '520px', padding: '2rem 2.25rem' }}
          >
            {/* Header Modal */}
            <div className="modal-header" style={{ marginBottom: '1.25rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div className="modal-pj-badge">PJ</div>
                <div>
                  <h3 className="modal-title" style={{ fontSize: '1.25rem', lineHeight: 1.2 }}>
                    Tambah Pengeluaran
                  </h3>
                  <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                    Catat biaya operasional usaha
                  </div>
                </div>
              </div>

              <button
                type="button"
                className="modal-close-btn"
                onClick={() => setIsModalTambahPengeluaranOpen(false)}
              >
                <X size={18} />
              </button>
            </div>

            <form onSubmit={handleSimpanPengeluaran} className="modal-form">
              {errorPengeluaran && (
                <div
                  style={{
                    padding: '0.75rem 1rem',
                    backgroundColor: '#FDECEA',
                    color: '#D32F2F',
                    borderRadius: '8px',
                    fontSize: '0.85rem',
                    lineHeight: '1.4',
                    marginBottom: '1rem',
                    border: '1px solid #FFCDD2'
                  }}
                >
                  {errorPengeluaran}
                </div>
              )}

              {/* Tanggal pengeluaran */}
              <div className="form-group">
                <label className="form-label">Tanggal pengeluaran</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="DD / MM / YYYY"
                  value={inputTanggalPengeluaran}
                  onChange={(e) => setInputTanggalPengeluaran(e.target.value)}
                  required
                />
              </div>

              {/* Kategori pengeluaran */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Kategori pengeluaran</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="Pilih kategori operasional"
                  value={inputKategori}
                  onChange={(e) => setInputKategori(e.target.value)}
                  required
                />
              </div>

              {/* Nominal pengeluaran (Rp) */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Nominal pengeluaran (Rp)</label>
                <input
                  type="number"
                  className="form-input"
                  placeholder="Masukkan nominal"
                  value={inputNominal}
                  onChange={(e) => setInputNominal(e.target.value)}
                  required
                />
              </div>

              {/* Keterangan (opsional) */}
              <div className="form-group" style={{ marginTop: '0.65rem' }}>
                <label className="form-label">Keterangan (opsional)</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="Contoh: Pembayaran listrik bulan ini"
                  value={inputKeteranganPengeluaran}
                  onChange={(e) => setInputKeteranganPengeluaran(e.target.value)}
                />
              </div>

              {/* Kotak Cegah biaya dihitung dua kali Sesuai Figma */}
              <div className="modal-note-box">
                <div className="modal-note-title">Cegah biaya dihitung dua kali</div>
                <div className="modal-note-text">
                  Pembelian persediaan jangan dimasukkan ulang sebagai pengeluaran operasional yang sama.
                </div>
              </div>

              {/* Tombol Aksi: Batal & Simpan Biaya */}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.85rem' }}>
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => setIsModalTambahPengeluaranOpen(false)}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-confirm">
                  Simpan Biaya
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
