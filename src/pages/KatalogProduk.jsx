import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { X, MessageSquare, MapPin, CheckCircle2, AlertCircle } from 'lucide-react';

const API_URL = 'http://localhost:8080';

// Data 4 varian rengginang sesuai mockup resmi Figma PPL-08 & PPL-09
const DAFTAR_PRODUK = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    varian: 'Udang',
    berat: 'Berat bersih 450 gr',
    beratBersih: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    stok: 12,
    stokText: '12 bungkus',
    statusKetersediaan: 'TERSEDIA',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Rengginang tradisional • kemasan 450 gram',
    gambar: '/images/products/rengginang-udang.jpg',
    gambarDetail: '/images/products/rengginang-udang-vertical.jpg'
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    varian: 'Bawang',
    berat: 'Berat bersih 450 gr',
    beratBersih: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    stok: 38,
    stokText: '38 bungkus',
    statusKetersediaan: 'TERSEDIA',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Rengginang tradisional • kemasan 450 gram',
    gambar: '/images/products/rengginang-bawang.jpg',
    gambarDetail: '/images/products/rengginang-bawang-vertical.jpg'
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    varian: 'Cumi',
    berat: 'Berat bersih 450 gr',
    beratBersih: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    stok: 0,
    stokText: '0 bungkus',
    statusKetersediaan: 'STOK HABIS',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Rengginang tradisional • kemasan 450 gram',
    gambar: '/images/products/rengginang-cumi.jpg',
    gambarDetail: '/images/products/rengginang-cumi-vertical.jpg'
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    varian: 'Terasi',
    berat: 'Berat bersih 450 gr',
    beratBersih: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    stok: 38,
    stokText: '38 bungkus',
    statusKetersediaan: 'TERSEDIA',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Rengginang tradisional • kemasan 450 gram',
    gambar: '/images/products/rengginang-terasi.jpg',
    gambarDetail: '/images/products/rengginang-terasi-vertical.jpg'
  }
];

export default function KatalogProduk() {
  const location = useLocation();
  const navigate = useNavigate();
  const params = useParams();

  // State untuk produk yang sedang dilihat detailnya (null = tampilkan katalog grid)
  const [produkDipilih, setProdukDipilih] = useState(null);

  // State galeri foto produk yang sedang aktif (0 = kemasan utama, 1 = label, 2 = detail isi)
  const [indexFotoAktif, setIndexFotoAktif] = useState(0);

  // Modal kontak / informasi pemesanan
  const [bukaModalKontak, setBukaModalKontak] = useState(false);
  const [judulModal, setJudulModal] = useState('Hubungi Penjual');

  // Deteksi jika URL membawa id produk (misal: /katalog?id=cumi atau /produk/udang)
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const idDariQuery = searchParams.get('id') || params.id;

    if (idDariQuery) {
      const produkCocok = DAFTAR_PRODUK.find(
        (p) => p.id.toLowerCase() === idDariQuery.toLowerCase()
      );
      if (produkCocok) {
        setProdukDipilih(produkCocok);
        setIndexFotoAktif(0);
        window.scrollTo({ top: 0, behavior: 'smooth' });
        return;
      }
    }

    if (!idDariQuery && !location.hash) {
      setProdukDipilih(null);
    }
  }, [location.search, params.id]);

  // Sinkronisasi data detail dan status ketersediaan dengan backend (WBS 1.1.9.3)
  useEffect(() => {
    if (!produkDipilih?.id) return;
    let isMounted = true;

    fetch(`${API_URL}/api/katalog/detail?id=${produkDipilih.id}`)
      .then((res) => {
        if (!res.ok) throw new Error('Gagal memuat status dari backend');
        return res.json();
      })
      .then((res) => {
        if (isMounted && res.status === 'sukses' && res.data) {
          const be = res.data;
          setProdukDipilih((prev) => {
            if (!prev || prev.id !== be.id_varian) return prev;
            return {
              ...prev,
              nama: be.nama,
              varian: be.nama_varian,
              beratBersih: be.berat_bersih_teks || prev.beratBersih,
              harga: be.harga_teks || prev.harga,
              stok: be.stok,
              stokText: be.stok_tersedia_teks || `${be.stok} bungkus`,
              statusKetersediaan: be.status_ketersediaan,
              bisaDibeli: be.bisa_dibeli,
              pesanStatus: be.pesan_status,
              subPesan: be.sub_pesan,
              deskripsi: be.deskripsi || prev.deskripsi
            };
          });
        }
      })
      .catch((err) => {
        // Fallback aman jika backend belum terhubung/offline
        console.warn('Backend detail produk offline, menggunakan fallback lokal:', err);
      });

    return () => {
      isMounted = false;
    };
  }, [produkDipilih?.id]);

  // Handler saat memilih varian dari grid katalog
  const handleBukaDetail = (produk) => {
    setProdukDipilih(produk);
    setIndexFotoAktif(0);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler kembali ke tampilan grid katalog
  const handleKembaliKeKatalog = () => {
    setProdukDipilih(null);
    navigate('/katalog');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler ganti varian rasa langsung pada halaman detail (Pilih Varian Rasa)
  const handleGantiVarian = (varianId) => {
    const target = DAFTAR_PRODUK.find((p) => p.id === varianId);
    if (target) {
      setProdukDipilih(target);
      setIndexFotoAktif(0);
    }
  };

  // Handler tombol Beli Produk dengan validasi pembatasan interaksi pembelian (WBS 1.1.9.3)
  const handleBeli = async (produk) => {
    if (produk.stok <= 0 || produk.bisaDibeli === false || produk.statusKetersediaan === 'STOK HABIS') {
      alert('Produk sedang habis. Pembelian dinonaktifkan.');
      return;
    }

    try {
      const res = await fetch(`${API_URL}/api/katalog/validasi-beli`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id_varian: produk.id,
          jumlah: 1
        })
      });

      const data = await res.json();
      if (!res.ok || (data && data.bisa_dibeli === false)) {
        alert(data.pesan || 'Produk sedang habis. Pembelian tidak dapat dilanjutkan.');
        return;
      }
    } catch (err) {
      console.warn('Backend validasi offline, fallback ke verifikasi lokal:', err);
    }

    setJudulModal(`Pesan ${produk.nama}`);
    setBukaModalKontak(true);
  };

  // Tentukan apakah produk yang sedang dilihat stoknya habis (WBS 1.1.9.3)
  const apakahStokHabis =
    produkDipilih &&
    (produkDipilih.stok <= 0 ||
      produkDipilih.statusKetersediaan === 'STOK HABIS' ||
      produkDipilih.bisaDibeli === false);

  // Daftar galeri foto untuk varian yang dipilih
  const getDaftarGaleri = (p) => [
    { label: 'Kemasan', src: p.gambarDetail || p.gambar },
    { label: 'Label', src: p.gambar },
    { label: 'Detail isi', src: '/images/products/rengginang-all-hero.jpg' }
  ];

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '85vh', paddingBottom: '4rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            TAMPILAN DETAIL PRODUK (PPL-09: Sesuai Gambar Mockup Figma)
            ====================================================================== */}
        {produkDipilih ? (
          <div>
            {/* Breadcrumb Sesuai Figma */}
            <nav style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.75rem' }}>
              <Link
                to="/"
                style={{ color: '#8E7D73', textDecoration: 'none', transition: 'color 0.2s' }}
                onMouseEnter={(e) => (e.target.style.color = '#43281C')}
                onMouseLeave={(e) => (e.target.style.color = '#8E7D73')}
              >
                Beranda
              </Link>
              <span style={{ margin: '0 0.5rem' }}>/</span>
              <span
                onClick={handleKembaliKeKatalog}
                style={{
                  color: '#8E7D73',
                  cursor: 'pointer',
                  textDecoration: 'none',
                  transition: 'color 0.2s'
                }}
                onMouseEnter={(e) => (e.target.style.color = '#43281C')}
                onMouseLeave={(e) => (e.target.style.color = '#8E7D73')}
              >
                Produk
              </span>
              <span style={{ margin: '0 0.5rem' }}>/</span>
              <span style={{ color: '#8E7D73' }}>
                {produkDipilih.nama}
              </span>
            </nav>

            {/* Header Detail Produk Sesuai Figma */}
            <div style={{ marginBottom: '2rem' }}>
              <h1
                style={{
                  fontSize: '2rem',
                  fontWeight: 800,
                  color: '#43281C',
                  marginBottom: '0.35rem',
                  letterSpacing: '-0.02em'
                }}
              >
                Detail Produk
              </h1>
              <p style={{ fontSize: '0.95rem', color: '#8E7D73' }}>
                Informasi varian, harga, berat, dan ketersediaan produk.
              </p>
            </div>

            {/* Layout 2 Kolom: Wadah Foto Kiri & Kartu Informasi Detail Kanan */}
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(360px, 1fr))',
                gap: '2.5rem',
                alignItems: 'start'
              }}
            >
              {/* KOLOM KIRI: Foto Produk & Galeri Thumbnail Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FAF5EE',
                  border: '1px solid #ECE4D9',
                  borderRadius: '24px',
                  padding: '1.75rem',
                  boxShadow: '0 2px 12px rgba(67, 40, 28, 0.03)'
                }}
              >
                {/* Foto Utama Produk */}
                <div
                  style={{
                    backgroundColor: '#FFFFFF',
                    borderRadius: '16px',
                    overflow: 'hidden',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: '420px',
                    marginBottom: '1.5rem',
                    border: '1px solid #EFE8DF'
                  }}
                >
                  <img
                    src={getDaftarGaleri(produkDipilih)[indexFotoAktif]?.src || produkDipilih.gambarDetail}
                    alt={produkDipilih.nama}
                    style={{
                      width: '100%',
                      height: '100%',
                      objectFit: 'contain',
                      display: 'block'
                    }}
                    onError={(e) => {
                      e.target.src = produkDipilih.gambar;
                    }}
                  />
                </div>

                {/* Baris Galeri Foto Produk Sesuai Figma */}
                <div style={{ display: 'flex', alignItems: 'center', gap: '1.25rem' }}>
                  <div style={{ display: 'flex', gap: '0.5rem' }}>
                    {getDaftarGaleri(produkDipilih).map((foto, idx) => (
                      <div
                        key={idx}
                        onClick={() => setIndexFotoAktif(idx)}
                        style={{
                          width: '62px',
                          height: '52px',
                          borderRadius: '8px',
                          overflow: 'hidden',
                          cursor: 'pointer',
                          border: indexFotoAktif === idx ? '2px solid #43281C' : '1px solid #E2D7CC',
                          transition: 'all 0.2s',
                          boxShadow: indexFotoAktif === idx ? '0 2px 6px rgba(67, 40, 28, 0.15)' : 'none',
                          backgroundColor: '#FFFFFF'
                        }}
                      >
                        <img
                          src={foto.src}
                          alt={foto.label}
                          style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                          onError={(e) => {
                            e.target.src = produkDipilih.gambar;
                          }}
                        />
                      </div>
                    ))}
                  </div>

                  <div>
                    <div style={{ fontSize: '0.875rem', fontWeight: 800, color: '#43281C', marginBottom: '0.15rem' }}>
                      Galeri foto produk
                    </div>
                    <div style={{ fontSize: '0.775rem', color: '#8E7D73' }}>
                      Kemasan &bull; label &bull; detail isi
                    </div>
                  </div>
                </div>
              </div>

              {/* KOLOM KANAN: Kartu Rincian & Pilihan Varian Sesuai Figma */}
              <div
                style={{
                  backgroundColor: '#FFFFFF',
                  border: '1px solid #ECE4D9',
                  borderRadius: '24px',
                  padding: '2rem 2.25rem',
                  boxShadow: '0 4px 16px rgba(67, 40, 28, 0.04)'
                }}
              >
                {/* Badge Status Ketersediaan Sesuai Figma */}
                <div style={{ marginBottom: '0.75rem' }}>
                  {apakahStokHabis ? (
                    <span
                      style={{
                        display: 'inline-block',
                        backgroundColor: '#FFEBEE',
                        color: '#C62828',
                        fontSize: '0.75rem',
                        fontWeight: 800,
                        letterSpacing: '0.04em',
                        padding: '0.35rem 0.95rem',
                        borderRadius: '9999px'
                      }}
                    >
                      STOK HABIS
                    </span>
                  ) : (
                    <span
                      style={{
                        display: 'inline-block',
                        backgroundColor: '#E8F5E9',
                        color: '#2E7D32',
                        fontSize: '0.75rem',
                        fontWeight: 800,
                        letterSpacing: '0.04em',
                        padding: '0.35rem 0.95rem',
                        borderRadius: '9999px'
                      }}
                    >
                      TERSEDIA
                    </span>
                  )}
                </div>

                {/* Nama Produk & Subjudul */}
                <h2
                  style={{
                    fontSize: '1.85rem',
                    fontWeight: 800,
                    color: '#43281C',
                    marginBottom: '0.2rem',
                    lineHeight: 1.2
                  }}
                >
                  {produkDipilih.nama}
                </h2>
                <p style={{ fontSize: '0.925rem', color: '#8E7D73', marginBottom: '1.35rem' }}>
                  {produkDipilih.deskripsi}
                </p>

                {/* Bagian Harga Jual Sesuai Figma */}
                <div style={{ marginBottom: '1.25rem' }}>
                  <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.2rem' }}>
                    Harga jual
                  </div>
                  <div style={{ display: 'flex', alignItems: 'baseline', gap: '1.25rem' }}>
                    <span style={{ fontSize: '1.85rem', fontWeight: 800, color: '#43281C' }}>
                      {produkDipilih.harga}
                    </span>
                    <span style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                      Contoh harga mockup
                    </span>
                  </div>
                  <div style={{ borderBottom: '1px solid #EFE8DF', marginTop: '1.15rem' }} />
                </div>

                {/* Pemilih Varian Rasa Sesuai Figma */}
                <div style={{ marginBottom: '1.5rem' }}>
                  <div
                    style={{
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      marginBottom: '0.65rem'
                    }}
                  >
                    Pilih varian rasa
                  </div>
                  <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
                    {DAFTAR_PRODUK.map((v) => {
                      const isActive = v.id === produkDipilih.id;
                      return (
                        <button
                          key={v.id}
                          type="button"
                          onClick={() => handleGantiVarian(v.id)}
                          style={{
                            backgroundColor: isActive ? '#43281C' : '#FFFFFF',
                            color: isActive ? '#FFFFFF' : '#43281C',
                            border: isActive ? 'none' : '1px solid #D8C7B8',
                            borderRadius: '10px',
                            padding: '0.65rem 1.35rem',
                            fontSize: '0.925rem',
                            fontWeight: 700,
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            boxShadow: isActive ? '0 2px 8px rgba(67, 40, 28, 0.15)' : 'none'
                          }}
                        >
                          {v.varian}
                        </button>
                      );
                    })}
                  </div>
                </div>

                {/* 2 Kotak Metrik: Berat Bersih & Stok Tersedia Sesuai Figma */}
                <div style={{ display: 'flex', gap: '1rem', marginBottom: '1.35rem' }}>
                  <div
                    style={{
                      flex: 1,
                      backgroundColor: '#FAF5EE',
                      borderRadius: '14px',
                      padding: '0.9rem 1.15rem'
                    }}
                  >
                    <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.25rem' }}>
                      Berat bersih
                    </div>
                    <div style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C' }}>
                      {produkDipilih.beratBersih}
                    </div>
                  </div>

                  <div
                    style={{
                      flex: 1,
                      backgroundColor: '#FAF5EE',
                      borderRadius: '14px',
                      padding: '0.9rem 1.15rem'
                    }}
                  >
                    <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.25rem' }}>
                      Stok tersedia
                    </div>
                    <div
                      style={{
                        fontSize: '1.25rem',
                        fontWeight: 800,
                        color: apakahStokHabis ? '#C62828' : '#2E7D32'
                      }}
                    >
                      {produkDipilih.stokText}
                    </div>
                  </div>
                </div>

                {/* Banner Status Ketersediaan Sesuai Figma */}
                <div style={{ marginBottom: '1.75rem' }}>
                  {apakahStokHabis ? (
                    <div
                      style={{
                        backgroundColor: '#FDF2F2',
                        border: '1px solid #FFCDD2',
                        borderRadius: '12px',
                        padding: '1rem 1.25rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.85rem'
                      }}
                    >
                      <div
                        style={{
                          width: '26px',
                          height: '26px',
                          borderRadius: '50%',
                          backgroundColor: '#C62828',
                          color: '#FFFFFF',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          fontWeight: 800,
                          fontSize: '0.9rem',
                          flexShrink: 0
                        }}
                      >
                        !
                      </div>
                      <div>
                        <div style={{ fontSize: '0.875rem', fontWeight: 800, color: '#C62828', marginBottom: '0.15rem' }}>
                          {produkDipilih.pesanStatus || 'Produk sedang habis'}
                        </div>
                        <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                          {produkDipilih.subPesan || 'Pembelian dinonaktifkan. Silakan pilih varian lain.'}
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div
                      style={{
                        backgroundColor: '#FAF5EE',
                        borderRadius: '12px',
                        padding: '1rem 1.25rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.85rem'
                      }}
                    >
                      <CheckCircle2 size={24} style={{ color: '#2E7D32', flexShrink: 0 }} />
                      <div>
                        <div style={{ fontSize: '0.875rem', fontWeight: 800, color: '#43281C', marginBottom: '0.15rem' }}>
                          {produkDipilih.pesanStatus || 'Stok tersedia'}
                        </div>
                        <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                          {produkDipilih.subPesan || 'Jumlah stok mengikuti data persediaan dari Owner.'}
                        </div>
                      </div>
                    </div>
                  )}
                </div>

                {/* Tombol Aksi Sesuai Figma */}
                <div style={{ display: 'flex', gap: '0.85rem' }}>
                  {apakahStokHabis ? (
                    <button
                      type="button"
                      disabled
                      style={{
                        flex: 1,
                        backgroundColor: '#D8CFC8',
                        color: '#8E7D73',
                        border: 'none',
                        borderRadius: '10px',
                        padding: '0.85rem 1.5rem',
                        fontSize: '0.95rem',
                        fontWeight: 700,
                        cursor: 'not-allowed'
                      }}
                    >
                      Stok Habis
                    </button>
                  ) : (
                    <button
                      type="button"
                      onClick={() => handleBeli(produkDipilih)}
                      style={{
                        flex: 1,
                        backgroundColor: '#43281C',
                        color: '#FFFFFF',
                        border: 'none',
                        borderRadius: '10px',
                        padding: '0.85rem 1.5rem',
                        fontSize: '0.95rem',
                        fontWeight: 700,
                        cursor: 'pointer',
                        transition: 'background-color 0.2s',
                        boxShadow: '0 2px 8px rgba(67, 40, 28, 0.15)'
                      }}
                      onMouseEnter={(e) => (e.target.style.backgroundColor = '#2F1B12')}
                      onMouseLeave={(e) => (e.target.style.backgroundColor = '#43281C')}
                    >
                      Beli Produk
                    </button>
                  )}

                  <button
                    type="button"
                    onClick={handleKembaliKeKatalog}
                    style={{
                      backgroundColor: '#FFFFFF',
                      color: '#43281C',
                      border: '1px solid #D8C7B8',
                      borderRadius: '10px',
                      padding: '0.85rem 1.75rem',
                      fontSize: '0.95rem',
                      fontWeight: 700,
                      cursor: 'pointer',
                      whiteSpace: 'nowrap',
                      transition: 'all 0.2s'
                    }}
                    onMouseEnter={(e) => {
                      e.target.style.borderColor = '#43281C';
                      e.target.style.backgroundColor = '#FAF5EE';
                    }}
                    onMouseLeave={(e) => {
                      e.target.style.borderColor = '#D8C7B8';
                      e.target.style.backgroundColor = '#FFFFFF';
                    }}
                  >
                    Kembali ke katalog
                  </button>
                </div>
              </div>
            </div>

            {/* Catatan Kaki Bawah Halaman Sesuai Figma */}
            <p
              style={{
                fontSize: '0.8rem',
                color: '#8E7D73',
                marginTop: '2.5rem',
                textAlign: 'left'
              }}
            >
              Harga dan stok pada mockup adalah data contoh; nilai aktual mengikuti data produk dan persediaan yang dikelola Owner.
            </p>
          </div>
        ) : (
          /* ======================================================================
             TAMPILAN GRID KATALOG PRODUK (Ketika belum ada produk dipilih)
             ====================================================================== */
          <div>
            {/* Header Halaman Katalog */}
            <div style={{ marginBottom: '2.5rem' }}>
              <h1
                style={{
                  fontSize: '2.25rem',
                  fontWeight: 800,
                  color: '#43281C',
                  marginBottom: '0.4rem',
                  letterSpacing: '-0.5px'
                }}
              >
                Semua varian rengginang
              </h1>
              <p style={{ fontSize: '1rem', color: '#8E7D73' }}>
                Temukan rasa yang paling kamu suka.
              </p>
            </div>

            {/* Sub-judul Grid */}
            <div style={{ marginBottom: '1.25rem' }}>
              <h2 style={{ fontSize: '1.15rem', fontWeight: 800, color: '#43281C' }}>
                Pilihan empat rasa
              </h2>
            </div>

            {/* Grid 4 Kartu Varian Rengginang */}
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
                gap: '1.5rem',
                marginBottom: '3rem'
              }}
            >
              {DAFTAR_PRODUK.map((item) => (
                <div
                  key={item.id}
                  style={{
                    backgroundColor: '#FFFFFF',
                    border: '1px solid #ECE4D9',
                    borderRadius: '16px',
                    padding: '1rem',
                    display: 'flex',
                    flexDirection: 'column',
                    transition: 'transform 0.2s, box-shadow 0.2s, border-color 0.2s',
                    boxShadow: '0 2px 8px rgba(67, 40, 28, 0.04)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.transform = 'translateY(-3px)';
                    e.currentTarget.style.boxShadow = '0 6px 18px rgba(67, 40, 28, 0.08)';
                    e.currentTarget.style.borderColor = '#D8C7B8';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.transform = 'translateY(0)';
                    e.currentTarget.style.boxShadow = '0 2px 8px rgba(67, 40, 28, 0.04)';
                    e.currentTarget.style.borderColor = '#ECE4D9';
                  }}
                >
                  {/* Foto Produk */}
                  <div
                    style={{
                      height: '240px',
                      borderRadius: '10px',
                      overflow: 'hidden',
                      backgroundColor: '#FAF5EE',
                      marginBottom: '1rem',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      cursor: 'pointer'
                    }}
                    onClick={() => handleBukaDetail(item)}
                  >
                    <img
                      src={item.gambar}
                      alt={item.nama}
                      style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover'
                      }}
                      onError={(e) => {
                        e.target.src = item.gambarDetail;
                      }}
                    />
                  </div>

                  {/* Keterangan Produk */}
                  <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                    <h3
                      style={{
                        fontSize: '1.15rem',
                        fontWeight: 800,
                        color: '#43281C',
                        marginBottom: '0.2rem',
                        cursor: 'pointer'
                      }}
                      onClick={() => handleBukaDetail(item)}
                    >
                      {item.nama}
                    </h3>

                    <p style={{ fontSize: '0.85rem', color: '#8E7D73', marginBottom: '0.5rem' }}>
                      {item.berat}
                    </p>

                    <div
                      style={{
                        fontSize: '1.15rem',
                        fontWeight: 800,
                        color: '#43281C',
                        marginBottom: '1rem'
                      }}
                    >
                      {item.harga}
                    </div>

                    {/* Tombol Aksi Kartu Grid */}
                    <div style={{ marginTop: 'auto', display: 'flex', gap: '0.5rem' }}>
                      <button
                        type="button"
                        onClick={() => handleBukaDetail(item)}
                        style={{
                          flex: 1,
                          backgroundColor: '#43281C',
                          color: '#FFFFFF',
                          border: 'none',
                          borderRadius: '8px',
                          padding: '0.65rem 0.75rem',
                          fontSize: '0.875rem',
                          fontWeight: 700,
                          cursor: 'pointer',
                          transition: 'background-color 0.2s'
                        }}
                        onMouseEnter={(e) => (e.target.style.backgroundColor = '#2F1B12')}
                        onMouseLeave={(e) => (e.target.style.backgroundColor = '#43281C')}
                      >
                        Lihat detail
                      </button>

                      {item.stok > 0 ? (
                        <button
                          type="button"
                          onClick={() => handleBeli(item)}
                          style={{
                            backgroundColor: '#FFFFFF',
                            color: '#43281C',
                            border: '1px solid #D8C7B8',
                            borderRadius: '8px',
                            padding: '0.65rem 0.85rem',
                            fontSize: '0.875rem',
                            fontWeight: 700,
                            cursor: 'pointer',
                            transition: 'all 0.2s'
                          }}
                          onMouseEnter={(e) => {
                            e.target.style.borderColor = '#43281C';
                            e.target.style.backgroundColor = '#FAF5EE';
                          }}
                          onMouseLeave={(e) => {
                            e.target.style.borderColor = '#D8C7B8';
                            e.target.style.backgroundColor = '#FFFFFF';
                          }}
                        >
                          Beli
                        </button>
                      ) : (
                        <span
                          style={{
                            backgroundColor: '#F5EBEB',
                            color: '#C62828',
                            borderRadius: '8px',
                            padding: '0.65rem 0.65rem',
                            fontSize: '0.75rem',
                            fontWeight: 700,
                            display: 'flex',
                            alignItems: 'center'
                          }}
                        >
                          Habis
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Banner Informasi Penutup Sesuai Figma */}
            <div
              style={{
                backgroundColor: '#FAF5EE',
                border: '1px solid #ECE4D9',
                borderRadius: '16px',
                padding: '1.5rem',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                flexWrap: 'wrap',
                gap: '1rem'
              }}
            >
              <div>
                <div style={{ fontSize: '1.05rem', fontWeight: 800, color: '#43281C', marginBottom: '0.2rem' }}>
                  Produk UMKM Rengginang Putra Jaya
                </div>
                <div style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
                  Dibuat secara tradisional di Situbondo dengan bahan pilihan tanpa pengawet.
                </div>
              </div>

              <div style={{ display: 'flex', gap: '0.75rem' }}>
                <button
                  type="button"
                  onClick={handleLihatInfoUsaha}
                  style={{
                    backgroundColor: '#FFFFFF',
                    border: '1px solid #D8C7B8',
                    borderRadius: '8px',
                    padding: '0.65rem 1.25rem',
                    fontSize: '0.875rem',
                    fontWeight: 700,
                    color: '#43281C',
                    cursor: 'pointer'
                  }}
                >
                  Informasi usaha
                </button>
                <button
                  type="button"
                  onClick={handleHubungiPenjual}
                  style={{
                    backgroundColor: '#43281C',
                    border: 'none',
                    borderRadius: '8px',
                    padding: '0.65rem 1.25rem',
                    fontSize: '0.875rem',
                    fontWeight: 700,
                    color: '#FFFFFF',
                    cursor: 'pointer'
                  }}
                >
                  Hubungi penjual
                </button>
              </div>
            </div>
          </div>
        )}

      </div>

      {/* Modal Kontak & Pemesanan */}
      {bukaModalKontak && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            backgroundColor: 'rgba(43, 27, 20, 0.55)',
            backdropFilter: 'blur(2px)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
            padding: '1.25rem'
          }}
          onClick={() => setBukaModalKontak(false)}
        >
          <div
            style={{
              backgroundColor: '#FFFFFF',
              borderRadius: '20px',
              padding: '2rem',
              width: '100%',
              maxWidth: '440px',
              boxShadow: '0 20px 40px rgba(67, 40, 28, 0.2)'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem' }}>
              <div style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C' }}>
                {judulModal}
              </div>
              <button
                type="button"
                onClick={() => setBukaModalKontak(false)}
                style={{
                  background: 'none',
                  border: 'none',
                  color: '#8E7D73',
                  cursor: 'pointer',
                  padding: '4px'
                }}
              >
                <X size={20} />
              </button>
            </div>

            <div
              style={{
                backgroundColor: '#FAF5EE',
                borderRadius: '12px',
                padding: '1.15rem',
                marginBottom: '1.5rem',
                fontSize: '0.9rem',
                color: '#66564C',
                lineHeight: 1.5
              }}
            >
              Untuk pemesanan dan informasi ketersediaan varian rengginang, silakan hubungi kontak resmi Owner:
              <div style={{ marginTop: '0.75rem', fontWeight: 700, color: '#43281C' }}>
                WhatsApp: +62 812-3456-7890<br />
                Lokasi: Situbondo, Jawa Timur
              </div>
            </div>

            <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
              <button
                type="button"
                onClick={() => setBukaModalKontak(false)}
                style={{
                  backgroundColor: '#FFFFFF',
                  border: '1px solid #D8C7B8',
                  borderRadius: '10px',
                  padding: '0.65rem 1.25rem',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  color: '#43281C',
                  cursor: 'pointer'
                }}
              >
                Tutup
              </button>
              <a
                href="https://wa.me/6281234567890"
                target="_blank"
                rel="noreferrer"
                style={{
                  backgroundColor: '#2E7D32',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '0.65rem 1.35rem',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  color: '#FFFFFF',
                  textDecoration: 'none',
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: '0.4rem'
                }}
              >
                <MessageSquare size={16} />
                Hubungi via WA
              </a>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
