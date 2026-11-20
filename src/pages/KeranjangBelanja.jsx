import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Trash2, Minus, Plus, ShoppingCart, ArrowLeft, AlertCircle } from 'lucide-react';

// ============================================================================
// PPL-11: KERANJANG BELANJA (WBS 1.2.2)
// - 1.2.2.1 (UX): Desain antarmuka keranjang belanja responsif & preview subtotal
// - 1.2.2.2 (FE): State manager keranjang & sinkronisasi localStorage
// - 1.2.2.3 (FE): Tambah, ubah jumlah, hapus item, dan perhitungan subtotal otomatis
// - Batasan: Jumlah tidak boleh melebihi stok tersedia
// ============================================================================

const STORAGE_KEY = 'keranjang_belanja';

// Data produk bawaan sesuai Figma PPL11_01
const DEFAULT_ITEMS = [
  {
    id: 1,
    slug: 'udang',
    nama: 'Rengginang Udang',
    varian: 'Udang',
    beratBersih: '450 gr',
    harga: 28000,
    jumlah: 2, // 2 bungkus pada PPL11_01
    stokTersedia: 12,
    gambar: '/images/products/rengginang-udang.jpg'
  },
  {
    id: 2,
    slug: 'bawang',
    nama: 'Rengginang Bawang',
    varian: 'Bawang',
    beratBersih: '450 gr',
    harga: 28000,
    jumlah: 1, // 1 bungkus
    stokTersedia: 38,
    gambar: '/images/products/rengginang-bawang.jpg'
  }
];

export default function KeranjangBelanja() {
  const navigate = useNavigate();
  const location = useLocation();

  // State Keranjang Belanja dari localStorage
  const [items, setItems] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) {
        const parsed = JSON.parse(saved);
        if (Array.isArray(parsed)) return parsed;
      }
    } catch (e) {
      console.warn('Gagal membaca keranjang dari localStorage:', e);
    }
    return DEFAULT_ITEMS;
  });

  // State Modal Konfirmasi Hapus (Figma Image 4)
  const [itemAkanDihapus, setItemAkanDihapus] = useState(null);

  // State Pesan Peringatan Stok
  const [pesanPeringatan, setPesanPeringatan] = useState('');

  // Simpan perubahan ke localStorage secara otomatis (WBS 1.2.2.2)
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
    } catch (e) {
      console.warn('Gagal menyimpan keranjang ke localStorage:', e);
    }
  }, [items]);

  // Format Rupiah standar
  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  // Handler tombol Tambah (+) dengan validasi batas stok tersedia
  const handleTambahJumlah = (id) => {
    setPesanPeringatan('');
    setItems((prevItems) =>
      prevItems.map((it) => {
        if (it.id === id) {
          if (it.jumlah >= it.stokTersedia) {
            setPesanPeringatan(
              `Jumlah pembelian ${it.nama} sudah mencapai batas maksimal sisa stok (${it.stokTersedia} bungkus).`
            );
            return it;
          }
          return { ...it, jumlah: it.jumlah + 1 };
        }
        return it;
      })
    );
  };

  // Handler tombol Kurang (-)
  const handleKurangJumlah = (id) => {
    setPesanPeringatan('');
    setItems((prevItems) =>
      prevItems.map((it) => {
        if (it.id === id) {
          if (it.jumlah > 1) {
            return { ...it, jumlah: it.jumlah - 1 };
          }
          // Jika kuantitas sudah 1, tawarkan konfirmasi hapus
          return it;
        }
        return it;
      })
    );
  };

  // Handler Buka Modal Konfirmasi Hapus
  const handleBukaModalHapus = (item) => {
    setItemAkanDihapus(item);
  };

  // Handler Konfirmasi Hapus Produk dari Keranjang (WBS 1.2.2.3)
  const handleKonfirmasiHapus = () => {
    if (!itemAkanDihapus) return;
    setItems((prevItems) => prevItems.filter((it) => it.id !== itemAkanDihapus.id));
    setItemAkanDihapus(null);
    setPesanPeringatan('');
  };

  // Kalkulasi Total dan Subtotal (WBS 1.2.2.3)
  const totalSubtotal = items.reduce((acc, curr) => acc + curr.harga * curr.jumlah, 0);
  const totalJumlahItem = items.reduce((acc, curr) => acc + curr.jumlah, 0);
  const totalJenisProduk = items.length;

  // Handler Lanjutkan ke Checkout
  const handleLanjutCheckout = () => {
    if (items.length === 0) return;
    navigate('/checkout', {
      state: {
        keranjang: items,
        produk: items[0],
        totalPembayaran: totalSubtotal
      }
    });
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', paddingBottom: '5rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            BREADCRUMB SESUAI FIGMA (Beranda / Keranjang)
            ====================================================================== */}
        <div style={{ marginBottom: '0.75rem', fontSize: '0.85rem', color: '#8E7D73' }}>
          <Link to="/" style={{ color: '#8E7D73', textDecoration: 'none' }}>
            Beranda
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <span style={{ color: '#43281C', fontWeight: 600 }}>Keranjang</span>
        </div>

        {/* ======================================================================
            JUDUL HALAMAN & SUBTITLE SESUAI FIGMA
            ====================================================================== */}
        <div style={{ marginBottom: '1.75rem' }}>
          <h1
            style={{
              fontSize: '2rem',
              fontWeight: 800,
              color: '#43281C',
              marginBottom: '0.35rem',
              letterSpacing: '-0.02em'
            }}
          >
            Keranjang belanja
          </h1>
          <p style={{ fontSize: '0.925rem', color: '#8E7D73', margin: 0 }}>
            Periksa produk dan jumlah sebelum melanjutkan checkout.
          </p>
        </div>

        {/* ======================================================================
            INDIKATOR PROGRESS 3 LANGKAH SESUAI FIGMA
            1 Keranjang (Aktif) -> 2 Checkout -> 3 Konfirmasi pesanan
            ====================================================================== */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            border: '1px solid #ECE4D9',
            borderRadius: '16px',
            padding: '1.25rem 2.5rem',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            marginBottom: '2.5rem',
            boxShadow: '0 1px 4px rgba(67, 40, 28, 0.03)'
          }}
        >
          {/* Langkah 1: Keranjang (Aktif) */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <div
              style={{
                width: '32px',
                height: '32px',
                borderRadius: '50%',
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '0.9rem',
                fontWeight: 700
              }}
            >
              1
            </div>
            <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
              Keranjang
            </span>
          </div>

          {/* Garis Pemisah 1 */}
          <div style={{ flex: 1, height: '1px', backgroundColor: '#E8DED3', margin: '0 1.5rem' }} />

          {/* Langkah 2: Checkout (Inaktif) */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <div
              style={{
                width: '32px',
                height: '32px',
                borderRadius: '50%',
                backgroundColor: '#F3EDE4',
                color: '#8E7D73',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '0.9rem',
                fontWeight: 600
              }}
            >
              2
            </div>
            <span style={{ fontSize: '0.95rem', fontWeight: 500, color: '#8E7D73' }}>
              Checkout
            </span>
          </div>

          {/* Garis Pemisah 2 */}
          <div style={{ flex: 1, height: '1px', backgroundColor: '#E8DED3', margin: '0 1.5rem' }} />

          {/* Langkah 3: Konfirmasi Pesanan (Inaktif) */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <div
              style={{
                width: '32px',
                height: '32px',
                borderRadius: '50%',
                backgroundColor: '#F3EDE4',
                color: '#8E7D73',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '0.9rem',
                fontWeight: 600
              }}
            >
              3
            </div>
            <span style={{ fontSize: '0.95rem', fontWeight: 500, color: '#8E7D73' }}>
              Konfirmasi pesanan
            </span>
          </div>
        </div>

        {/* Pesan Peringatan Batas Stok Realtime */}
        {pesanPeringatan && (
          <div
            style={{
              backgroundColor: '#FEF3C7',
              border: '1px solid #F59E0B',
              color: '#92400E',
              padding: '0.85rem 1.25rem',
              borderRadius: '10px',
              fontSize: '0.875rem',
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '0.65rem'
            }}
          >
            <AlertCircle size={18} />
            <span>{pesanPeringatan}</span>
          </div>
        )}

        {/* ======================================================================
            KONDISI 1: KERANJANG KOSONG SESUAI FIGMA (Image 5)
            ====================================================================== */}
        {items.length === 0 ? (
          <div
            style={{
              backgroundColor: '#FFFFFF',
              border: '1px solid #ECE4D9',
              borderRadius: '20px',
              padding: '4.5rem 2rem',
              textAlign: 'center',
              boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)',
              maxWidth: '860px',
              margin: '0 auto'
            }}
          >
            {/* Lingkaran Ikon Keranjang */}
            <div
              style={{
                width: '96px',
                height: '96px',
                borderRadius: '50%',
                backgroundColor: '#F5EBE1',
                color: '#7A5B47',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                margin: '0 auto 1.5rem auto'
              }}
            >
              <ShoppingCart size={44} strokeWidth={1.8} />
            </div>

            <h2
              style={{
                fontSize: '1.5rem',
                fontWeight: 800,
                color: '#43281C',
                marginBottom: '0.5rem'
              }}
            >
              Keranjangmu masih kosong
            </h2>

            <p style={{ fontSize: '0.95rem', color: '#8E7D73', marginBottom: '2rem' }}>
              Yuk, pilih rengginang favoritmu terlebih dahulu.
            </p>

            <Link
              to="/katalog"
              style={{
                display: 'inline-block',
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                textDecoration: 'none',
                fontWeight: 700,
                fontSize: '0.95rem',
                padding: '0.95rem 2.25rem',
                borderRadius: '10px',
                transition: 'background-color 0.2s',
                boxShadow: '0 2px 8px rgba(67, 40, 28, 0.15)'
              }}
              onMouseEnter={(e) => (e.target.style.backgroundColor = '#2F1B12')}
              onMouseLeave={(e) => (e.target.style.backgroundColor = '#43281C')}
            >
              Lihat Produk
            </Link>

            <div style={{ marginTop: '2rem', fontSize: '0.8rem', color: '#A3968C' }}>
              Pelanggan dapat mulai belanja tanpa membuat akun.
            </div>
          </div>
        ) : (
          /* ======================================================================
              KONDISI 2: KERANJANG TERISI SESUAI FIGMA PPL11_01 & PPL11_02
              Layout 2 Kolom: Kiri (Daftar Produk) & Kanan (Ringkasan Belanja)
              ====================================================================== */
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
              gap: '2rem',
              alignItems: 'start'
            }}
          >
            {/* ==================================================================
                KOLOM KIRI: DAFTAR PRODUK DI KERANJANG
                ================================================================== */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              {/* Header Kolom Kiri */}
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', margin: 0 }}>
                  Produk di keranjang
                </h2>
                <span style={{ fontSize: '0.875rem', color: '#8E7D73' }}>
                  {totalJenisProduk} produk
                </span>
              </div>

              {/* Daftar Kartu Produk */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                {items.map((item) => {
                  const subtotalItem = item.harga * item.jumlah;
                  const isMaxStock = item.jumlah >= item.stokTersedia;
                  const isMinQty = item.jumlah <= 1;

                  return (
                    <div
                      key={item.id}
                      style={{
                        backgroundColor: '#FFFFFF',
                        border: '1px solid #ECE4D9',
                        borderRadius: '16px',
                        padding: '1.5rem',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)',
                        flexWrap: 'wrap',
                        gap: '1.25rem'
                      }}
                    >
                      {/* Bagian Kiri: Thumbnail & Info Produk */}
                      <div style={{ display: 'flex', alignItems: 'center', gap: '1.25rem' }}>
                        {/* Thumbnail Kemasan Produk */}
                        <div
                          style={{
                            width: '84px',
                            height: '84px',
                            borderRadius: '10px',
                            overflow: 'hidden',
                            backgroundColor: '#FAF5EE',
                            border: '1px solid #EFE8DF',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            flexShrink: 0
                          }}
                        >
                          <img
                            src={item.gambar}
                            alt={item.nama}
                            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                            onError={(e) => {
                              e.target.src = '/images/products/rengginang-udang.jpg';
                            }}
                          />
                        </div>

                        {/* Rincian Nama, Berat, Harga, & Sisa Stok */}
                        <div>
                          <h3
                            style={{
                              fontSize: '1.05rem',
                              fontWeight: 800,
                              color: '#43281C',
                              margin: '0 0 0.2rem 0'
                            }}
                          >
                            {item.nama}
                          </h3>
                          <div style={{ fontSize: '0.825rem', color: '#8E7D73', marginBottom: '0.35rem' }}>
                            Berat bersih: {item.beratBersih}
                          </div>
                          <div
                            style={{
                              fontSize: '0.925rem',
                              fontWeight: 700,
                              color: '#43281C',
                              marginBottom: '0.35rem'
                            }}
                          >
                            {formatRupiah(item.harga)} / bungkus
                          </div>
                          <div
                            style={{
                              fontSize: '0.825rem',
                              fontWeight: 600,
                              color: '#15803D' // Hijau sesuai Figma PPL-11
                            }}
                          >
                            Stok tersedia: {item.stokTersedia} bungkus
                          </div>
                        </div>
                      </div>

                      {/* Bagian Kanan: Stepper Kuantitas, Tombol Hapus, & Subtotal Produk */}
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '1.75rem',
                          flexWrap: 'wrap'
                        }}
                      >
                        {/* Kontrol Stepper & Tombol Hapus */}
                        <div>
                          <div
                            style={{
                              display: 'flex',
                              justifyContent: 'space-between',
                              alignItems: 'center',
                              marginBottom: '0.4rem'
                            }}
                          >
                            <span style={{ fontSize: '0.775rem', color: '#8E7D73' }}>Jumlah</span>

                            {/* Tombol Hapus Ikon Sampah Sesuai Figma */}
                            <button
                              type="button"
                              onClick={() => handleBukaModalHapus(item)}
                              title="Hapus produk dari keranjang"
                              style={{
                                background: '#FAF0EC',
                                border: '1px solid #F3DDD4',
                                borderRadius: '6px',
                                width: '28px',
                                height: '28px',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                cursor: 'pointer',
                                color: '#BA483A',
                                transition: 'all 0.2s',
                                padding: 0
                              }}
                              onMouseEnter={(e) => {
                                e.currentTarget.style.backgroundColor = '#FDE2D9';
                              }}
                              onMouseLeave={(e) => {
                                e.currentTarget.style.backgroundColor = '#FAF0EC';
                              }}
                            >
                              <Trash2 size={15} />
                            </button>
                          </div>

                          {/* Kotak Stepper (-  angka  +) */}
                          <div
                            style={{
                              display: 'flex',
                              alignItems: 'center',
                              border: '1px solid #D8C7B8',
                              borderRadius: '8px',
                              backgroundColor: '#FFFFFF',
                              height: '38px',
                              overflow: 'hidden'
                            }}
                          >
                            {/* Tombol Minus (-) */}
                            <button
                              type="button"
                              onClick={() => handleKurangJumlah(item.id)}
                              disabled={isMinQty}
                              style={{
                                width: '36px',
                                height: '100%',
                                border: 'none',
                                backgroundColor: isMinQty ? '#F9F7F5' : 'transparent',
                                color: isMinQty ? '#C5B7AB' : '#43281C',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                cursor: isMinQty ? 'not-allowed' : 'pointer',
                                transition: 'background-color 0.15s'
                              }}
                              onMouseEnter={(e) => {
                                if (!isMinQty) e.target.style.backgroundColor = '#F5ECE3';
                              }}
                              onMouseLeave={(e) => {
                                if (!isMinQty) e.target.style.backgroundColor = 'transparent';
                              }}
                            >
                              <Minus size={15} strokeWidth={2.5} />
                            </button>

                            {/* Angka Kuantitas Terpilih */}
                            <div
                              style={{
                                width: '42px',
                                textAlign: 'center',
                                fontSize: '0.95rem',
                                fontWeight: 800,
                                color: '#43281C',
                                userSelect: 'none'
                              }}
                            >
                              {item.jumlah}
                            </div>

                            {/* Tombol Plus (+) dengan pembatasan stok */}
                            <button
                              type="button"
                              onClick={() => handleTambahJumlah(item.id)}
                              disabled={isMaxStock}
                              style={{
                                width: '36px',
                                height: '100%',
                                border: 'none',
                                backgroundColor: isMaxStock ? '#F9F7F5' : 'transparent',
                                color: isMaxStock ? '#C5B7AB' : '#43281C',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                cursor: isMaxStock ? 'not-allowed' : 'pointer',
                                transition: 'background-color 0.15s'
                              }}
                              onMouseEnter={(e) => {
                                if (!isMaxStock) e.target.style.backgroundColor = '#F5ECE3';
                              }}
                              onMouseLeave={(e) => {
                                if (!isMaxStock) e.target.style.backgroundColor = 'transparent';
                              }}
                            >
                              <Plus size={15} strokeWidth={2.5} />
                            </button>
                          </div>
                        </div>

                        {/* Subtotal Item Produk */}
                        <div style={{ textAlign: 'right', minWidth: '100px' }}>
                          <div
                            style={{
                              fontSize: '1.15rem',
                              fontWeight: 800,
                              color: '#43281C',
                              marginBottom: '0.15rem'
                            }}
                          >
                            {formatRupiah(subtotalItem)}
                          </div>
                          <div style={{ fontSize: '0.75rem', color: '#8E7D73' }}>
                            Subtotal produk
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* Card Banner Info Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FAF5EE',
                  border: '1px solid #F0E6D8',
                  borderRadius: '12px',
                  padding: '1.15rem 1.35rem',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '0.35rem'
                }}
              >
                <div style={{ fontSize: '0.875rem', fontWeight: 800, color: '#43281C' }}>
                  Keranjang dapat dipakai tanpa login.
                </div>
                <div style={{ fontSize: '0.825rem', color: '#6B5E55', lineHeight: 1.5 }}>
                  Jumlah dan subtotal akan diperbarui otomatis; stok tetap divalidasi saat checkout.
                </div>
              </div>

              {/* Link Lanjut Belanja */}
              <div>
                <Link
                  to="/katalog"
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: '0.4rem',
                    fontSize: '0.9rem',
                    fontWeight: 700,
                    color: '#43281C',
                    textDecoration: 'none'
                  }}
                >
                  <ArrowLeft size={16} />
                  <span>Lanjut belanja</span>
                </Link>
                <div style={{ fontSize: '0.75rem', color: '#8E7D73', marginTop: '0.35rem' }}>
                  Contoh jumlah dan harga untuk kebutuhan desain UI/UX.
                </div>
              </div>
            </div>

            {/* ==================================================================
                KOLOM KANAN: RINGKASAN BELANJA SESUAI FIGMA PPL-11
                ================================================================== */}
            <div>
              <div
                style={{
                  backgroundColor: '#FFFFFF',
                  border: '1px solid #ECE4D9',
                  borderRadius: '16px',
                  padding: '2rem',
                  boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
                }}
              >
                {/* Header Ringkasan */}
                <div style={{ marginBottom: '1.5rem' }}>
                  <h2
                    style={{
                      fontSize: '1.35rem',
                      fontWeight: 800,
                      color: '#43281C',
                      marginBottom: '0.2rem'
                    }}
                  >
                    Ringkasan belanja
                  </h2>
                  <div style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
                    {totalJenisProduk} jenis produk
                  </div>
                </div>

                {/* Rincian Produk Sesuai Mockup */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.85rem', marginBottom: '1.5rem' }}>
                  {items.map((it) => (
                    <div
                      key={it.id}
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        fontSize: '0.9rem'
                      }}
                    >
                      <span style={{ color: '#43281C', fontWeight: 600 }}>{it.nama}</span>
                      <span style={{ color: '#6B5E55' }}>
                        {it.jumlah} × {formatRupiah(it.harga)}
                      </span>
                    </div>
                  ))}
                </div>

                {/* Garis Pemisah */}
                <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

                {/* Baris Subtotal */}
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Subtotal</span>
                  <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                    {formatRupiah(totalSubtotal)}
                  </span>
                </div>

                {/* Baris Ongkos Kirim */}
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Ongkos kirim</span>
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Dihitung saat checkout</span>
                </div>

                {/* Garis Pemisah */}
                <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

                {/* Baris Total Sementara */}
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: '1.75rem'
                  }}
                >
                  <span style={{ fontSize: '1rem', fontWeight: 700, color: '#43281C' }}>
                    Total sementara
                  </span>
                  <span style={{ fontSize: '1.45rem', fontWeight: 800, color: '#43281C' }}>
                    {formatRupiah(totalSubtotal)}
                  </span>
                </div>

                {/* Tombol Lanjutkan Checkout */}
                <button
                  type="button"
                  onClick={handleLanjutCheckout}
                  style={{
                    width: '100%',
                    backgroundColor: '#43281C',
                    color: '#FFFFFF',
                    border: 'none',
                    borderRadius: '10px',
                    padding: '1rem',
                    fontSize: '0.975rem',
                    fontWeight: 700,
                    cursor: 'pointer',
                    transition: 'background-color 0.2s',
                    boxShadow: '0 2px 8px rgba(67, 40, 28, 0.15)'
                  }}
                  onMouseEnter={(e) => (e.target.style.backgroundColor = '#2F1B12')}
                  onMouseLeave={(e) => (e.target.style.backgroundColor = '#43281C')}
                >
                  Lanjutkan Checkout
                </button>

                {/* Keterangan Bawah Tombol Sesuai Figma */}
                <p
                  style={{
                    fontSize: '0.775rem',
                    color: '#8E7D73',
                    textAlign: 'center',
                    marginTop: '0.75rem',
                    marginBottom: 0
                  }}
                >
                  Login tidak wajib untuk melanjutkan pesanan.
                </p>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* ======================================================================
          MODAL KONFIRMASI HAPUS PRODUK DARI KERANJANG SESUAI FIGMA (Image 4)
          ====================================================================== */}
      {itemAkanDihapus && (
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
          onClick={() => setItemAkanDihapus(null)}
        >
          <div
            style={{
              backgroundColor: '#FFFFFF',
              borderRadius: '20px',
              maxWidth: '540px',
              width: '100%',
              padding: '2.25rem',
              boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
              border: '1px solid #EFE8DF'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header Modal dengan Ikon Tempat Sampah Merah-Salmon */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '1.25rem', marginBottom: '1.5rem' }}>
              <div
                style={{
                  width: '56px',
                  height: '56px',
                  borderRadius: '14px',
                  backgroundColor: '#FDECE7',
                  color: '#BA483A',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  flexShrink: 0
                }}
              >
                <Trash2 size={28} strokeWidth={2.2} />
              </div>
              <div>
                <h3
                  style={{
                    fontSize: '1.35rem',
                    fontWeight: 800,
                    color: '#43281C',
                    margin: '0 0 0.2rem 0'
                  }}
                >
                  Hapus produk?
                </h3>
                <div style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
                  Konfirmasi tindakan
                </div>
              </div>
            </div>

            {/* Garis Pemisah Tipis Sesuai Mockup */}
            <div style={{ borderBottom: '1px solid #F0E8DF', marginBottom: '1.5rem' }} />

            {/* Konten Penjelasan Tindakan */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem', marginBottom: '2rem' }}>
              <p
                style={{
                  fontSize: '1.05rem',
                  color: '#43281C',
                  fontWeight: 600,
                  margin: 0
                }}
              >
                {itemAkanDihapus.nama} akan dihapus dari keranjang.
              </p>
              <p style={{ fontSize: '0.925rem', color: '#6B5E55', margin: 0, lineHeight: 1.5 }}>
                Subtotal otomatis dihitung ulang setelah produk dihapus.
              </p>
              <p style={{ fontSize: '0.85rem', color: '#8E7D73', margin: 0, lineHeight: 1.5 }}>
                Tindakan ini tidak menghapus data produk di katalog.
              </p>
            </div>

            {/* Tombol Aksi (Batal & Hapus) Sesuai Figma */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
              <button
                type="button"
                onClick={() => setItemAkanDihapus(null)}
                style={{
                  padding: '0.85rem 2rem',
                  border: '1px solid #D8C7B8',
                  borderRadius: '12px',
                  backgroundColor: '#FFFFFF',
                  color: '#43281C',
                  fontSize: '0.95rem',
                  fontWeight: 700,
                  cursor: 'pointer',
                  transition: 'background-color 0.15s'
                }}
                onMouseEnter={(e) => (e.target.style.backgroundColor = '#FAF5EE')}
                onMouseLeave={(e) => (e.target.style.backgroundColor = '#FFFFFF')}
              >
                Batal
              </button>
              <button
                type="button"
                onClick={handleKonfirmasiHapus}
                style={{
                  padding: '0.85rem 2.25rem',
                  border: 'none',
                  borderRadius: '12px',
                  backgroundColor: '#BA483A',
                  color: '#FFFFFF',
                  fontSize: '0.95rem',
                  fontWeight: 700,
                  cursor: 'pointer',
                  transition: 'background-color 0.15s',
                  boxShadow: '0 2px 6px rgba(186, 72, 58, 0.25)'
                }}
                onMouseEnter={(e) => (e.target.style.backgroundColor = '#A33B2E')}
                onMouseLeave={(e) => (e.target.style.backgroundColor = '#BA483A')}
              >
                Hapus
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
