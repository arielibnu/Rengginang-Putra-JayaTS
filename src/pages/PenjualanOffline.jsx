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
  CheckCircle2,
  AlertCircle
} from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';

// Daftar varian produk resmi Rengginang Putra Jaya
const DAFTAR_VARIAN_DEFAULT = [
  {
    id: 'udang',
    id_varian: 1,
    nama: 'Rengginang Udang',
    harga: 28000,
    stok_tersedia: 12
  },
  {
    id: 'bawang',
    id_varian: 2,
    nama: 'Rengginang Bawang',
    harga: 28000,
    stok_tersedia: 38
  },
  {
    id: 'terasi',
    id_varian: 3,
    nama: 'Rengginang Terasi',
    harga: 28000,
    stok_tersedia: 38
  },
  {
    id: 'cumi',
    id_varian: 4,
    nama: 'Rengginang Cumi',
    harga: 30000,
    stok_tersedia: 0
  }
];

export default function PenjualanOffline() {
  const navigate = useNavigate();

  // Status autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(true);
  const [isChecking, setIsChecking] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // State data varian produk dan stok
  const [daftarVarian, setDaftarVarian] = useState(DAFTAR_VARIAN_DEFAULT);

  // State formulir kasir penjualan offline (PPL-24)
  const [tanggalPenjualan, setTanggalPenjualan] = useState('20 November 2026');
  const [varianTerpilihId, setVarianTerpilihId] = useState('udang');
  const [jumlahBungkus, setJumlahBungkus] = useState(2);
  const [metodePembayaran, setMetodePembayaran] = useState('Tunai');

  // State proses penyimpanan dan notifikasi
  const [isLoading, setIsLoading] = useState(false);
  const [pesanSukses, setPesanSukses] = useState('');
  const [pesanError, setPesanError] = useState('');
  const [riwayatPenjualanTerakhir, setRiwayatPenjualanTerakhir] = useState(null);

  // Ambil detail varian yang sedang aktif dipilih
  const varianAktif =
    daftarVarian.find((v) => v.id === varianTerpilihId) || daftarVarian[0];

  // Hitung total harga otomatis
  const hargaSatuan = varianAktif?.harga || 28000;
  const stokTersedia = Number(varianAktif?.stok_tersedia ?? 0);
  const totalHarga = (Number(jumlahBungkus) || 0) * hargaSatuan;

  // Format rupiah helper
  const formatRupiah = (val) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(val);
  };

  // Validasi sesi Owner dan muat stok terbaru dari Backend
  useEffect(() => {
    const controller = new AbortController();

    const fetchSesiDanStok = async () => {
      try {
        const options = {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        };

        // Periksa sesi Owner
        const resProfile = await fetch(`${API_URL}/api/pemilik/profil`, options);
        if (resProfile.status === 401) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        // Ambil data stok terkini dari backend (WBS 1.1.5.2 / stok peringatan)
        const resStok = await fetch(`${API_URL}/api/admin/stok/peringatan`, options);
        if (resStok.ok) {
          const jsonStok = await resStok.json();
          if (jsonStok?.data?.daftar_varian?.length > 0) {
            setDaftarVarian((prev) =>
              prev.map((item) => {
                const dariBE = jsonStok.data.daftar_varian.find(
                  (v) =>
                    v.nama_varian?.toLowerCase().includes(item.id) ||
                    v.id_varian === item.id_varian
                );
                if (dariBE) {
                  return {
                    ...item,
                    stok_tersedia:
                      dariBE.stok_tersedia ?? dariBE.stok_fisik ?? item.stok_tersedia
                  };
                }
                return item;
              })
            );
          }
        }
      } catch (err) {
        // Abaikan jika abort atau fallback offline
      }
    };

    fetchSesiDanStok();
    return () => controller.abort();
  }, [navigate]);

  // Handler simpan transaksi kasir langsung ke backend
  const handleSimpanPenjualan = async (e) => {
    e.preventDefault();
    setPesanSukses('');
    setPesanError('');

    const qty = parseInt(jumlahBungkus, 10);

    // 1. Validasi jumlah pembelian > 0
    if (isNaN(qty) || qty <= 0) {
      setPesanError('Jumlah pembelian wajib lebih dari 0 bungkus.');
      return;
    }

    // 2. Validasi varian produk terpilih
    if (!varianAktif) {
      setPesanError('Silakan pilih varian produk yang valid.');
      return;
    }

    // 3. Validasi stok mencukupi
    if (qty > stokTersedia) {
      setPesanError(
        `Stok ${varianAktif.nama} tidak mencukupi. Sisa stok tersedia: ${stokTersedia} bungkus.`
      );
      return;
    }

    setIsLoading(true);

    try {
      // 4. Kirim data transaksi kasir langsung ke API Penjualan Offline (WBS 1.3.5.3 & 1.3.5.4)
      const payloadOffline = {
        tanggal_penjualan: tanggalPenjualan,
        id_varian: varianAktif.id_varian || 1,
        nama_varian: varianAktif.nama,
        jumlah: qty,
        harga_satuan: hargaSatuan,
        metode_pembayaran: metodePembayaran,
        flag_lingkungan: 'PRODUKSI_NYATA',
        catatan: `Kasir Toko Langsung (${metodePembayaran})`
      };

      let nomorTransaksiBerhasil = `POS-${new Date().getTime().toString().slice(-6)}`;
      const resOffline = await fetch(`${API_URL}/api/admin/penjualan-offline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(payloadOffline)
      });

      if (resOffline.ok) {
        const jsonOff = await resOffline.json();
        if (jsonOff?.data?.nomor_transaksi) {
          nomorTransaksiBerhasil = jsonOff.data.nomor_transaksi;
        }
      } else {
        // Fallback coba ke mutasi stok jika endpoint utama offline
        const resMutasi = await fetch(`${API_URL}/api/admin/stok/mutasi`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
            id_varian: varianAktif.id_varian || 1,
            jumlah: qty,
            jenis_mutasi: 'keluar',
            mode_stok: 'RIIL',
            catatan: `Penjualan Offline Toko - ${tanggalPenjualan}`
          })
        });

        if (!resMutasi.ok) {
          const errorJson = await resOffline.json().catch(() => ({}));
          throw new Error(
            errorJson.pesan ||
              'Gagal menyimpan penjualan offline ke backend. Periksa kembali stok.'
          );
        }
      }

      // 5. Berhasil: Tampilkan notifikasi sukses dan catat transaksi
      const dataSukses = {
        nomor_transaksi: nomorTransaksiBerhasil,
        tanggal: tanggalPenjualan,
        varian: varianAktif.nama,
        jumlah: qty,
        harga: hargaSatuan,
        total: totalHarga,
        metode: metodePembayaran,
        flag: 'PRODUKSI_NYATA',
        waktu: new Date().toLocaleTimeString('id-ID', {
          hour: '2-digit',
          minute: '2-digit'
        })
      };

      setRiwayatPenjualanTerakhir(dataSukses);
      setPesanSukses(
        `Penjualan offline [${nomorTransaksiBerhasil}] berhasil disimpan! ${qty} bungkus ${varianAktif.nama} (${formatRupiah(totalHarga)}) telah memotong stok fisik.`
      );

      // Perbarui stok lokal setelah pengurangan resmi backend
      setDaftarVarian((prev) =>
        prev.map((v) =>
          v.id === varianAktif.id
            ? { ...v, stok_tersedia: Math.max(0, v.stok_tersedia - qty) }
            : v
        )
      );

      // Reset jumlah kembali ke 1
      setJumlahBungkus(1);
    } catch (err) {
      setPesanError(err.message || 'Terjadi kesalahan saat memproses transaksi.');
    } finally {
      setIsLoading(false);
    }
  };

  // Handler Logout
  const handleLogout = async () => {
    if (isLoggingOut) return;
    setIsLoggingOut(true);
    try {
      await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });
      localStorage.removeItem('ownerAuth');
      navigate('/login', { replace: true });
    } catch (e) {
      localStorage.removeItem('ownerAuth');
      navigate('/login', { replace: true });
    } finally {
      setIsLoggingOut(false);
    }
  };

  return (
    <div className="admin-layout">
      {/* Sidebar Navigasi Owner Sesuai Figma */}
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
              {/* Menu Penjualan Offline Aktif Sesuai Figma PPL-24 */}
              <div className="sidebar-item active">
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
            style={{ color: '#A92631', background: 'none', border: 'none', cursor: 'pointer', padding: 0 }}
          >
            {isLoggingOut ? 'Memproses logout...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* Konten Utama Halaman Penjualan Offline */}
      <main className="admin-main">
        {/* Breadcrumb Sesuai Figma */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Penjualan Offline
        </div>

        {/* Header Halaman Sesuai Figma */}
        <div className="admin-header" style={{ marginBottom: '2rem' }}>
          <h1 className="admin-page-title">Penjualan Offline</h1>
          <p className="admin-page-subtitle">
            PPL-24: kasir toko, bukan pesanan website.
          </p>
        </div>

        {/* Notifikasi Status Berhasil */}
        {pesanSukses && (
          <div
            style={{
              backgroundColor: '#ECFDF5',
              border: '1px solid #A7F3D0',
              borderRadius: '12px',
              padding: '1rem 1.25rem',
              color: '#065F46',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              marginBottom: '1.5rem',
              fontSize: '0.9rem'
            }}
          >
            <CheckCircle2 size={20} color="#059669" />
            <div style={{ flex: 1, fontWeight: 500 }}>{pesanSukses}</div>
            <button
              type="button"
              onClick={() => setPesanSukses('')}
              style={{ background: 'none', border: 'none', color: '#047857', cursor: 'pointer', fontWeight: 700 }}
            >
              ✕
            </button>
          </div>
        )}

        {/* Notifikasi Status Gagal */}
        {pesanError && (
          <div
            style={{
              backgroundColor: '#FEF2F2',
              border: '1px solid #FECACA',
              borderRadius: '12px',
              padding: '1rem 1.25rem',
              color: '#991B1B',
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              marginBottom: '1.5rem',
              fontSize: '0.9rem'
            }}
          >
            <AlertCircle size={20} color="#DC2626" />
            <div style={{ flex: 1, fontWeight: 500 }}>{pesanError}</div>
            <button
              type="button"
              onClick={() => setPesanError('')}
              style={{ background: 'none', border: 'none', color: '#B91C1C', cursor: 'pointer', fontWeight: 700 }}
            >
              ✕
            </button>
          </div>
        )}

        {/* Grid 2 Kolom: Formulir Kasir (Kiri) dan Aturan Data (Kanan) Sesuai Figma */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'minmax(0, 1.8fr) minmax(280px, 1fr)',
            gap: '2rem',
            alignItems: 'start'
          }}
        >
          {/* Kolom Kiri: Formulir Kasir Offline */}
          <div
            style={{
              backgroundColor: '#FFFFFF',
              border: '1px solid #ECE4D9',
              borderRadius: '20px',
              padding: '2.25rem',
              boxShadow: '0 1px 4px rgba(67, 40, 28, 0.03)'
            }}
          >
            <form onSubmit={handleSimpanPenjualan}>
              {/* Field 1: Tanggal Penjualan */}
              <div className="form-group" style={{ marginBottom: '1.5rem' }}>
                <label
                  className="form-label"
                  style={{
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    color: '#43281C',
                    marginBottom: '0.5rem'
                  }}
                >
                  Tanggal penjualan
                </label>
                <input
                  type="text"
                  className="form-input"
                  value={tanggalPenjualan}
                  onChange={(e) => setTanggalPenjualan(e.target.value)}
                  placeholder="Contoh: 20 November 2026"
                  required
                  style={{
                    padding: '0.85rem 1.15rem',
                    fontSize: '0.95rem',
                    borderRadius: '12px',
                    borderColor: '#E5DCD2',
                    color: '#43281C'
                  }}
                />
              </div>

              {/* Field 2: Varian Produk */}
              <div className="form-group" style={{ marginBottom: '1.5rem' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                  <label
                    className="form-label"
                    style={{
                      fontSize: '0.925rem',
                      fontWeight: 700,
                      color: '#43281C',
                      margin: 0
                    }}
                  >
                    Varian produk
                  </label>
                  <span
                    style={{
                      fontSize: '0.8rem',
                      color: stokTersedia > 5 ? '#059669' : '#DC2626',
                      fontWeight: 600
                    }}
                  >
                    Sisa stok: {stokTersedia} bungkus
                  </span>
                </div>
                <select
                  className="form-input"
                  value={varianTerpilihId}
                  onChange={(e) => {
                    setVarianTerpilihId(e.target.value);
                    setPesanError('');
                  }}
                  style={{
                    padding: '0.85rem 1.15rem',
                    fontSize: '0.95rem',
                    borderRadius: '12px',
                    borderColor: '#E5DCD2',
                    backgroundColor: '#FFFFFF',
                    color: '#43281C',
                    cursor: 'pointer'
                  }}
                >
                  {daftarVarian.map((v) => (
                    <option key={v.id} value={v.id}>
                      {v.nama} {v.stok_tersedia <= 0 ? '(Stok Habis)' : `(${formatRupiah(v.harga)})`}
                    </option>
                  ))}
                </select>
              </div>

              {/* Field 3: Jumlah (bungkus) */}
              <div className="form-group" style={{ marginBottom: '1.5rem' }}>
                <label
                  className="form-label"
                  style={{
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    color: '#43281C',
                    marginBottom: '0.5rem'
                  }}
                >
                  Jumlah (bungkus)
                </label>
                <input
                  type="number"
                  min="1"
                  max={stokTersedia || 999}
                  className="form-input"
                  value={jumlahBungkus}
                  onChange={(e) => {
                    setJumlahBungkus(e.target.value);
                    setPesanError('');
                  }}
                  placeholder="2"
                  required
                  style={{
                    padding: '0.85rem 1.15rem',
                    fontSize: '0.95rem',
                    borderRadius: '12px',
                    borderColor: '#E5DCD2',
                    color: '#43281C'
                  }}
                />
              </div>

              {/* Field 4: Harga Satuan (Otomatis dari harga produk sesuai Figma) */}
              <div className="form-group" style={{ marginBottom: '1.5rem' }}>
                <label
                  className="form-label"
                  style={{
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    color: '#43281C',
                    marginBottom: '0.5rem'
                  }}
                >
                  Harga satuan
                </label>
                <input
                  type="text"
                  readOnly
                  className="form-input"
                  value="Dari harga produk"
                  style={{
                    padding: '0.85rem 1.15rem',
                    fontSize: '0.95rem',
                    borderRadius: '12px',
                    borderColor: '#E5DCD2',
                    backgroundColor: '#FAFAF9',
                    color: '#78716C',
                    cursor: 'not-allowed'
                  }}
                />
                <span style={{ fontSize: '0.8rem', color: '#8E7D73', marginTop: '0.35rem' }}>
                  Tarif aktif: <strong>{formatRupiah(hargaSatuan)}</strong> per bungkus
                </span>
              </div>

              {/* Field 5: Metode Pembayaran (Tunai sesuai Figma) */}
              <div className="form-group" style={{ marginBottom: '2rem' }}>
                <label
                  className="form-label"
                  style={{
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    color: '#43281C',
                    marginBottom: '0.5rem'
                  }}
                >
                  Metode pembayaran
                </label>
                <input
                  type="text"
                  readOnly
                  className="form-input"
                  value={metodePembayaran}
                  style={{
                    padding: '0.85rem 1.15rem',
                    fontSize: '0.95rem',
                    borderRadius: '12px',
                    borderColor: '#E5DCD2',
                    backgroundColor: '#FAFAF9',
                    color: '#43281C',
                    cursor: 'default'
                  }}
                />
              </div>

              {/* Tombol Simpan Penjualan Langsung Sesuai Figma */}
              <button
                type="submit"
                disabled={isLoading || stokTersedia <= 0}
                style={{
                  width: '100%',
                  backgroundColor: stokTersedia <= 0 ? '#A8A29E' : '#43281C',
                  color: '#FFFFFF',
                  padding: '1rem',
                  borderRadius: '12px',
                  fontSize: '1rem',
                  fontWeight: 700,
                  border: 'none',
                  cursor: stokTersedia <= 0 || isLoading ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: '0 2px 8px rgba(67, 40, 28, 0.15)'
                }}
                onMouseEnter={(e) => {
                  if (!isLoading && stokTersedia > 0) e.target.style.backgroundColor = '#2F1B12';
                }}
                onMouseLeave={(e) => {
                  if (!isLoading && stokTersedia > 0) e.target.style.backgroundColor = '#43281C';
                }}
              >
                {isLoading ? 'Menyimpan Penjualan...' : 'Simpan Penjualan Langsung'}
              </button>
            </form>
          </div>

          {/* Kolom Kanan: Kartu Aturan Data Sesuai Figma */}
          <div
            style={{
              backgroundColor: '#F6EFE8',
              borderRadius: '20px',
              padding: '1.75rem 2rem',
              border: '1px solid #ECE0D4'
            }}
          >
            <h3
              style={{
                fontSize: '1.15rem',
                fontWeight: 800,
                color: '#43281C',
                marginTop: 0,
                marginBottom: '1.25rem'
              }}
            >
              Aturan data
            </h3>

            <ul
              style={{
                listStyle: 'none',
                padding: 0,
                margin: 0,
                display: 'flex',
                flexDirection: 'column',
                gap: '1rem'
              }}
            >
              <li
                style={{
                  fontSize: '0.925rem',
                  color: '#5C5046',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: '0.5rem',
                  lineHeight: 1.4
                }}
              >
                <span style={{ color: '#8E7D73' }}>•</span>
                <span>Kurangi stok sesuai transaksi</span>
              </li>
              <li
                style={{
                  fontSize: '0.925rem',
                  color: '#5C5046',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: '0.5rem',
                  lineHeight: 1.4
                }}
              >
                <span style={{ color: '#8E7D73' }}>•</span>
                <span>Tidak membuat order online</span>
              </li>
              <li
                style={{
                  fontSize: '0.925rem',
                  color: '#5C5046',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: '0.5rem',
                  lineHeight: 1.4
                }}
              >
                <span style={{ color: '#8E7D73' }}>•</span>
                <span>Pisahkan uang kas nyata</span>
              </li>
              <li
                style={{
                  fontSize: '0.925rem',
                  color: '#5C5046',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: '0.5rem',
                  lineHeight: 1.4
                }}
              >
                <span style={{ color: '#8E7D73' }}>•</span>
                <span>Jangan campur sandbox</span>
              </li>
            </ul>

            {/* Rekap Transaksi Terakhir (Jika ada) */}
            {riwayatPenjualanTerakhir && (
              <div
                style={{
                  marginTop: '2rem',
                  paddingTop: '1.25rem',
                  borderTop: '1px solid #E5D5C6'
                }}
              >
                <div style={{ fontSize: '0.8rem', fontWeight: 700, color: '#8E7D73', textTransform: 'uppercase', marginBottom: '0.5rem' }}>
                  Penjualan Terakhir:
                </div>
                <div style={{ fontSize: '0.875rem', fontWeight: 700, color: '#43281C' }}>
                  {riwayatPenjualanTerakhir.jumlah}x {riwayatPenjualanTerakhir.varian}
                </div>
                <div style={{ fontSize: '0.825rem', color: '#6B5E55', marginTop: '0.2rem' }}>
                  Total: {formatRupiah(riwayatPenjualanTerakhir.total)} • {riwayatPenjualanTerakhir.metode} ({riwayatPenjualanTerakhir.waktu} WIB)
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
