import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { X, MessageSquare, MapPin } from 'lucide-react';

// Data 4 varian rengginang sesuai mockup resmi Figma
const DAFTAR_PRODUK = [
  {
    id: 'udang',
    nama: 'Rengginang Udang',
    varian: 'Udang',
    berat: 'Berat bersih 450 gr',
    beratGram: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    ketersediaan: 'Tersedia • contoh',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Varian rasa udang dalam kemasan 450 gram.',
    gambar: '/images/products/rengginang-udang.jpg',
    gambarDetail: '/images/products/rengginang-udang-vertical.jpg'
  },
  {
    id: 'bawang',
    nama: 'Rengginang Bawang',
    varian: 'Bawang',
    berat: 'Berat bersih 450 gr',
    beratGram: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    ketersediaan: 'Tersedia • contoh',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Varian rasa bawang dalam kemasan 450 gram.',
    gambar: '/images/products/rengginang-bawang.jpg',
    gambarDetail: '/images/products/rengginang-bawang-vertical.jpg'
  },
  {
    id: 'cumi',
    nama: 'Rengginang Cumi',
    varian: 'Cumi',
    berat: 'Berat bersih 450 gr',
    beratGram: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    ketersediaan: 'Tersedia • contoh',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Varian rasa cumi dalam kemasan 450 gram.',
    gambar: '/images/products/rengginang-cumi.jpg',
    gambarDetail: '/images/products/rengginang-cumi-vertical.jpg'
  },
  {
    id: 'terasi',
    nama: 'Rengginang Terasi',
    varian: 'Terasi',
    berat: 'Berat bersih 450 gr',
    beratGram: '450 gram',
    harga: 'Rp28.000',
    hargaAngka: 28000,
    ketersediaan: 'Tersedia • contoh',
    kategori: 'Rengginang tradisional',
    deskripsi: 'Varian rasa terasi dalam kemasan 450 gram.',
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

  // Modal kontak / informasi usaha saat tombol beli atau hubungi penjual diklik
  const [bukaModalKontak, setBukaModalKontak] = useState(false);
  const [judulModal, setJudulModal] = useState('Hubungi Penjual');

  // Deteksi jika URL membawa id produk (misal: /katalog?id=bawang atau /produk/:id)
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const idDariQuery = searchParams.get('id') || params.id;

    if (idDariQuery) {
      const produkCocok = DAFTAR_PRODUK.find(
        (p) => p.id.toLowerCase() === idDariQuery.toLowerCase()
      );
      if (produkCocok) {
        setProdukDipilih(produkCocok);
        window.scrollTo({ top: 0, behavior: 'smooth' });
        return;
      }
    }

    // Jika tidak ada parameter spesifik dan navigasi ke root katalog
    if (!idDariQuery && !location.hash) {
      setProdukDipilih(null);
    }
  }, [location.search, params.id]);

  // Handler untuk membuka halaman detail produk (Gambar 2 Figma)
  const handleBukaDetail = (produk) => {
    setProdukDipilih(produk);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler kembali ke tampilan semua katalog (Gambar 1 Figma)
  const handleKembaliKeKatalog = () => {
    setProdukDipilih(null);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler tombol Beli atau Hubungi Penjual
  const handleBeli = (produk) => {
    setJudulModal(`Pesan ${produk.nama}`);
    setBukaModalKontak(true);
  };

  const handleHubungiPenjual = () => {
    setJudulModal('Hubungi Penjual');
    setBukaModalKontak(true);
  };

  const handleLihatInfoUsaha = () => {
    setJudulModal('Informasi Usaha UMKM');
    setBukaModalKontak(true);
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '85vh', paddingBottom: '4rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            TAMPILAN 1: DETAIL PRODUK (Sesuai Gambar Figma Kedua)
            Aktif jika ada varian produk yang dipilih
            ====================================================================== */}
        {produkDipilih ? (
          <div>
            {/* Breadcrumb Navigasi */}
            <nav style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '2.5rem' }}>
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
              <span style={{ color: '#43281C', fontWeight: 600 }}>
                {produkDipilih.nama}
              </span>
            </nav>

            {/* Layout 2 Kolom: Gambar Produk Kiri & Informasi Detail Kanan */}
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(340px, 1fr))',
                gap: '3.5rem',
                alignItems: 'start'
              }}
            >
              {/* Kolom Kiri: Wadah Foto Produk Memanjang dari Atas ke Bawah */}
              <div
                style={{
                  backgroundColor: '#FAF5EE',
                  border: '1px solid #ECE4D9',
                  borderRadius: '28px',
                  padding: '2rem 1.75rem',
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                  boxShadow: '0 4px 16px rgba(67, 40, 28, 0.03)'
                }}
              >
                <img
                  src={produkDipilih.gambarDetail || produkDipilih.gambar}
                  alt={produkDipilih.nama}
                  style={{
                    width: '100%',
                    maxWidth: '420px',
                    height: '520px',
                    objectFit: 'cover',
                    borderRadius: '16px',
                    display: 'block'
                  }}
                  onError={(e) => {
                    e.target.src = produkDipilih.gambar;
                  }}
                />
              </div>

              {/* Kolom Kanan: Rincian Produk */}
              <div>
                <h1
                  style={{
                    fontSize: '2.4rem',
                    fontWeight: 800,
                    color: '#43281C',
                    marginBottom: '0.4rem',
                    lineHeight: 1.2
                  }}
                >
                  {produkDipilih.nama}
                </h1>

                <p
                  style={{
                    fontSize: '1rem',
                    color: '#8E7D73',
                    marginBottom: '1.75rem',
                    lineHeight: 1.5
                  }}
                >
                  {produkDipilih.deskripsi}
                </p>

                {/* Harga Produk */}
                <div style={{ marginBottom: '1.75rem' }}>
                  <div
                    style={{
                      fontSize: '1.9rem',
                      fontWeight: 800,
                      color: '#43281C',
                      marginBottom: '0.2rem'
                    }}
                  >
                    {produkDipilih.harga}
                  </div>
                  <div style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
                    Harga contoh untuk desain
                  </div>
                </div>

                {/* Tabel Informasi Produk */}
                <div
                  style={{
                    borderTop: '1px solid #EFE6DC',
                    paddingTop: '1.25rem',
                    marginBottom: '2.25rem'
                  }}
                >
                  <h2
                    style={{
                      fontSize: '1.05rem',
                      fontWeight: 800,
                      color: '#43281C',
                      marginBottom: '1.1rem'
                    }}
                  >
                    Informasi produk
                  </h2>

                  <div style={{ display: 'flex', flexDirection: 'column' }}>
                    {/* Baris 1: Varian */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        padding: '0.85rem 0',
                        borderBottom: '1px solid #F4EDE5',
                        fontSize: '0.925rem'
                      }}
                    >
                      <span style={{ color: '#8E7D73' }}>Varian</span>
                      <span style={{ fontWeight: 700, color: '#43281C' }}>
                        {produkDipilih.varian}
                      </span>
                    </div>

                    {/* Baris 2: Berat Bersih */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        padding: '0.85rem 0',
                        borderBottom: '1px solid #F4EDE5',
                        fontSize: '0.925rem'
                      }}
                    >
                      <span style={{ color: '#8E7D73' }}>Berat bersih</span>
                      <span style={{ fontWeight: 700, color: '#43281C' }}>
                        {produkDipilih.beratGram}
                      </span>
                    </div>

                    {/* Baris 3: Ketersediaan */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        padding: '0.85rem 0',
                        borderBottom: '1px solid #F4EDE5',
                        fontSize: '0.925rem'
                      }}
                    >
                      <span style={{ color: '#8E7D73' }}>Ketersediaan</span>
                      <span style={{ fontWeight: 700, color: '#43281C' }}>
                        {produkDipilih.ketersediaan}
                      </span>
                    </div>

                    {/* Baris 4: Kategori */}
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        padding: '0.85rem 0',
                        borderBottom: '1px solid #F4EDE5',
                        fontSize: '0.925rem'
                      }}
                    >
                      <span style={{ color: '#8E7D73' }}>Kategori</span>
                      <span style={{ fontWeight: 700, color: '#43281C' }}>
                        {produkDipilih.kategori}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Tombol Aksi: Lihat Varian Lainnya & Hubungi Penjual */}
                <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    onClick={handleKembaliKeKatalog}
                    style={{
                      backgroundColor: '#4A3525',
                      color: '#FFFFFF',
                      border: 'none',
                      borderRadius: '10px',
                      padding: '0.85rem 1.85rem',
                      fontSize: '0.95rem',
                      fontWeight: 700,
                      cursor: 'pointer',
                      transition: 'background-color 0.2s',
                      boxShadow: '0 2px 6px rgba(74, 53, 37, 0.2)'
                    }}
                    onMouseEnter={(e) => (e.target.style.backgroundColor = '#382619')}
                    onMouseLeave={(e) => (e.target.style.backgroundColor = '#4A3525')}
                  >
                    Lihat varian lainnya
                  </button>

                  <button
                    type="button"
                    onClick={handleHubungiPenjual}
                    style={{
                      backgroundColor: '#FFFFFF',
                      color: '#43281C',
                      border: '1px solid #D8C7B8',
                      borderRadius: '10px',
                      padding: '0.85rem 1.85rem',
                      fontSize: '0.95rem',
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
                    Hubungi penjual
                  </button>
                </div>
              </div>
            </div>
          </div>
        ) : (
          /* ======================================================================
             TAMPILAN 2: KATALOG SEMUA VARIAN RENGGINANG (Sesuai Gambar Figma Pertama)
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
                    onClick={() => handleBukaDetail(item)}
                    style={{
                      backgroundColor: '#FAF5EE',
                      borderRadius: '12px',
                      overflow: 'hidden',
                      height: '205px',
                      cursor: 'pointer',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center'
                    }}
                  >
                    <img
                      src={item.gambar}
                      alt={item.nama}
                      style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover',
                        transition: 'transform 0.3s'
                      }}
                      onMouseEnter={(e) => (e.target.style.transform = 'scale(1.04)')}
                      onMouseLeave={(e) => (e.target.style.transform = 'scale(1)')}
                      onError={(e) => {
                        e.target.src = '/images/products/rengginang-all-hero.jpg';
                      }}
                    />
                  </div>

                  {/* Konten Kartu */}
                  <div style={{ marginTop: '0.85rem', flex: 1, display: 'flex', flexDirection: 'column' }}>
                    <h3
                      onClick={() => handleBukaDetail(item)}
                      style={{
                        fontSize: '1.05rem',
                        fontWeight: 800,
                        color: '#43281C',
                        marginBottom: '0.2rem',
                        cursor: 'pointer'
                      }}
                    >
                      {item.nama}
                    </h3>

                    <div style={{ fontSize: '0.825rem', color: '#8E7D73', marginBottom: '0.65rem' }}>
                      {item.berat}
                    </div>

                    <div
                      style={{
                        fontSize: '1.15rem',
                        fontWeight: 800,
                        color: '#43281C',
                        marginBottom: '1rem',
                        marginTop: 'auto'
                      }}
                    >
                      {item.harga}
                    </div>

                    {/* Tombol Beli & Lihat Detail */}
                    <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                      <button
                        type="button"
                        onClick={() => handleBeli(item)}
                        style={{
                          backgroundColor: '#4A3525',
                          color: '#FFFFFF',
                          border: 'none',
                          borderRadius: '8px',
                          padding: '0.5rem 1.25rem',
                          fontSize: '0.875rem',
                          fontWeight: 700,
                          cursor: 'pointer',
                          transition: 'background-color 0.2s',
                          flex: 1
                        }}
                        onMouseEnter={(e) => (e.target.style.backgroundColor = '#382619')}
                        onMouseLeave={(e) => (e.target.style.backgroundColor = '#4A3525')}
                      >
                        Beli
                      </button>

                      <button
                        type="button"
                        onClick={() => handleBukaDetail(item)}
                        style={{
                          backgroundColor: '#FFFFFF',
                          color: '#4A3525',
                          border: '1px solid #E6D8C8',
                          borderRadius: '8px',
                          padding: '0.5rem 0.85rem',
                          fontSize: '0.875rem',
                          fontWeight: 700,
                          cursor: 'pointer',
                          whiteSpace: 'nowrap',
                          transition: 'all 0.2s'
                        }}
                        onMouseEnter={(e) => {
                          e.target.style.borderColor = '#4A3525';
                          e.target.style.backgroundColor = '#FAF5EE';
                        }}
                        onMouseLeave={(e) => {
                          e.target.style.borderColor = '#E6D8C8';
                          e.target.style.backgroundColor = '#FFFFFF';
                        }}
                      >
                        Lihat detail ›
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Kotak Informasi Kualitas (Empat rasa, satu kualitas) */}
            <div
              style={{
                backgroundColor: '#FAF5EE',
                border: '1px solid #ECE4D9',
                borderRadius: '16px',
                padding: '1.75rem 2rem',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                flexWrap: 'wrap',
                gap: '1.25rem',
                boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
              }}
            >
              <div>
                <h3
                  style={{
                    fontSize: '1.25rem',
                    fontWeight: 800,
                    color: '#43281C',
                    marginBottom: '0.35rem'
                  }}
                >
                  Empat rasa, satu kualitas.
                </h3>
                <p style={{ fontSize: '0.875rem', color: '#8E7D73', margin: 0 }}>
                  Foto kemasan asli • Berat bersih 450 gr • Harga dan stok mengikuti data Owner.
                </p>
              </div>

              <button
                type="button"
                onClick={handleLihatInfoUsaha}
                style={{
                  backgroundColor: '#FFFFFF',
                  color: '#43281C',
                  border: '1px solid #D8C7B8',
                  borderRadius: '10px',
                  padding: '0.75rem 1.6rem',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  cursor: 'pointer',
                  transition: 'all 0.2s',
                  whiteSpace: 'nowrap'
                }}
                onMouseEnter={(e) => {
                  e.target.style.borderColor = '#43281C';
                  e.target.style.backgroundColor = '#FAF8F5';
                }}
                onMouseLeave={(e) => {
                  e.target.style.borderColor = '#D8C7B8';
                  e.target.style.backgroundColor = '#FFFFFF';
                }}
              >
                Lihat informasi usaha
              </button>
            </div>
          </div>
        )}
      </div>

      {/* ======================================================================
          MODAL INFORMASI KONTAK & PEMESANAN (Popup Sederhana)
          ====================================================================== */}
      {bukaModalKontak && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            backgroundColor: 'rgba(67, 40, 28, 0.45)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '1.25rem',
            zIndex: 1000
          }}
          onClick={() => setBukaModalKontak(false)}
        >
          <div
            style={{
              backgroundColor: '#FFFFFF',
              borderRadius: '20px',
              padding: '2rem 2.25rem',
              maxWidth: '480px',
              width: '100%',
              boxShadow: '0 20px 40px rgba(67, 40, 28, 0.2)',
              position: 'relative'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            {/* Tombol Tutup Silang */}
            <button
              type="button"
              onClick={() => setBukaModalKontak(false)}
              style={{
                position: 'absolute',
                top: '1.25rem',
                right: '1.25rem',
                border: 'none',
                background: 'none',
                cursor: 'pointer',
                color: '#8E7D73',
                padding: '4px'
              }}
            >
              <X size={20} />
            </button>

            {/* Header Modal */}
            <div style={{ textAlign: 'center', marginBottom: '1.5rem' }}>
              <div
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '52px',
                  height: '52px',
                  borderRadius: '50%',
                  backgroundColor: '#FAF5EE',
                  color: '#43281C',
                  marginBottom: '0.85rem'
                }}
              >
                <MessageSquare size={24} />
              </div>
              <h3 style={{ fontSize: '1.35rem', fontWeight: 800, color: '#43281C', marginBottom: '0.35rem' }}>
                {judulModal}
              </h3>
              <p style={{ fontSize: '0.875rem', color: '#8E7D73' }}>
                Rengginang Putra Jaya — Produksi Hj. Arifah
              </p>
            </div>

            {/* Konten Keterangan */}
            <div
              style={{
                backgroundColor: '#FAF5EE',
                border: '1px solid #ECE4D9',
                borderRadius: '12px',
                padding: '1.25rem',
                marginBottom: '1.5rem',
                fontSize: '0.9rem',
                color: '#43281C',
                lineHeight: 1.6
              }}
            >
              <div style={{ display: 'flex', gap: '0.65rem', marginBottom: '0.5rem' }}>
                <MapPin size={18} color="#8C5835" style={{ flexShrink: 0, marginTop: '2px' }} />
                <span>
                  <strong>Lokasi Usaha:</strong> Situbondo, Jawa Timur
                </span>
              </div>
              <div>
                Pemesanan dan info kemitraan dapat langsung menghubungi pengrajin UMKM kami. Semua produk diproduksi higienis dengan resep tradisional otentik.
              </div>
            </div>

            {/* Tombol Selesai */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
              <button
                type="button"
                onClick={() => setBukaModalKontak(false)}
                style={{
                  backgroundColor: '#4A3525',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '8px',
                  padding: '0.65rem 1.5rem',
                  fontSize: '0.9rem',
                  fontWeight: 700,
                  cursor: 'pointer',
                  width: '100%'
                }}
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
