import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { X, MessageSquare, MapPin, CheckCircle2, AlertCircle } from 'lucide-react';

const API_URL = 'http://localhost:8080';

// Data 4 varian rengginang sesuai mockup resmi Figma PPL-08, PPL-09 & PPL-10
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
    deskripsi: 'Rengginang tradisional dengan rasa udang gurih.',
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
    deskripsi: 'Rengginang tradisional dengan rasa bawang gurih dan renyah.',
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
    deskripsi: 'Rengginang tradisional dengan cita rasa cumi gurih khas pesisir.',
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
    deskripsi: 'Rengginang tradisional dengan aroma terasi khas yang mantap.',
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

  // State jumlah pembelian (PPL-10)
  const [jumlahBeli, setJumlahBeli] = useState(1);

  // State validasi kuota stok realtime dari backend (WBS 1.2.1.3)
  const [validasiRealtime, setValidasiRealtime] = useState({
    tervalidasi: false,
    bisaDibeli: true,
    pesan: '',
    keterangan: '',
    sisaSetelahBeli: null
  });

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
        setJumlahBeli(produkCocok.stok > 0 ? 1 : 0);
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

  // Validasi kuota sisa stok secara realtime ke backend saat jumlah beli berubah (WBS 1.2.1.3)
  useEffect(() => {
    if (!produkDipilih?.id) return;
    let isMounted = true;
    const targetJumlah = jumlahBeli > 0 ? jumlahBeli : 1;

    fetch(`${API_URL}/api/katalog/cek-kuota?id=${produkDipilih.id}&jumlah=${targetJumlah}`)
      .then((res) => res.json())
      .then((res) => {
        if (!isMounted) return;
        if (res.data) {
          const info = res.data;
          setValidasiRealtime({
            tervalidasi: true,
            bisaDibeli: info.bisa_dibeli,
            pesan: info.pesan,
            keterangan: info.keterangan,
            sisaSetelahBeli: info.sisa_setelah_beli
          });

          // Sinkronisasi data stok aktual terkini dari database MySQL
          if (info.stok_tersedia !== undefined && info.stok_tersedia !== produkDipilih.stok) {
            setProdukDipilih((prev) => {
              if (!prev) return prev;
              return {
                ...prev,
                stok: info.stok_tersedia,
                stokText: info.format_stok || `${info.stok_tersedia} bungkus`
              };
            });
          }
        }
      })
      .catch((err) => {
        console.warn('Backend cek-kuota offline, verifikasi fallback lokal:', err);
      });

    return () => {
      isMounted = false;
    };
  }, [produkDipilih?.id, jumlahBeli]);

  // Handler saat memilih varian dari grid katalog
  const handleBukaDetail = (produk) => {
    setProdukDipilih(produk);
    setJumlahBeli(produk.stok > 0 ? 1 : 0);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler kembali ke tampilan grid katalog
  const handleKembaliKeKatalog = () => {
    setProdukDipilih(null);
    navigate('/katalog');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Handler ganti varian rasa langsung pada halaman detail (Pilih Varian)
  const handleGantiVarian = (varianId) => {
    const target = DAFTAR_PRODUK.find((p) => p.id === varianId);
    if (target) {
      setProdukDipilih(target);
      setJumlahBeli(target.stok > 0 ? 1 : 0);
    }
  };

  // Handler tombol Beli Produk dengan validasi kuota stok realtime (WBS 1.2.1.3 & PPL-11 Alur Keranjang)
  const handleBeli = async (produk) => {
    const qtyBeli = (produkDipilih && produkDipilih.id === produk.id) ? (jumlahBeli || 1) : 1;

    if (produk.stok <= 0) {
      alert('Produk sedang habis. Pembelian dinonaktifkan.');
      return;
    }
    if (qtyBeli > produk.stok) {
      alert(`Jumlah pembelian (${qtyBeli} bungkus) melebihi stok yang tersedia (${produk.stok} bungkus).`);
      return;
    }
    if (qtyBeli < 1) {
      alert('Jumlah pembelian minimal 1 bungkus.');
      return;
    }

    try {
      const res = await fetch(`${API_URL}/api/katalog/validasi-kuota`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id_varian: produk.id,
          jumlah: qtyBeli
        })
      });

      const resp = await res.json();
      if (!res.ok || (resp && resp.data && resp.data.bisa_dibeli === false)) {
        const pesanGagal = resp.pesan || resp.data?.pesan || 'Jumlah pembelian tidak valid atau melebihi kuota stok.';
        alert(pesanGagal);
        return;
      }
    } catch (err) {
      console.warn('Backend validasi offline, fallback ke verifikasi lokal:', err);
    }

    // Masukkan produk ke Keranjang Belanja (localStorage: keranjang_belanja)
    try {
      const STORAGE_KEY = 'keranjang_belanja';
      const saved = localStorage.getItem(STORAGE_KEY);
      let currentItems = [];
      if (saved) {
        try {
          const parsed = JSON.parse(saved);
          if (Array.isArray(parsed)) currentItems = parsed;
        } catch {}
      }

      // Cari apakah varian produk sudah ada di keranjang
      const existingIndex = currentItems.findIndex(
        (it) => it.slug === produk.id || it.id === produk.id || it.varian?.toLowerCase() === produk.varian?.toLowerCase()
      );

      if (existingIndex >= 0) {
        const itemLama = currentItems[existingIndex];
        const stokMaks = produk.stok || itemLama.stokTersedia || 99;
        const jumlahBaru = Math.min(stokMaks, (itemLama.jumlah || 0) + qtyBeli);
        currentItems[existingIndex] = {
          ...itemLama,
          jumlah: jumlahBaru,
          stokTersedia: stokMaks,
          harga: produk.hargaAngka || itemLama.harga || 28000
        };
      } else {
        const idNumerikMap = { udang: 1, bawang: 2, cumi: 3, terasi: 4, lorjuk: 5, manis: 6 };
        const idItem = idNumerikMap[produk.id] || Date.now();
        currentItems.push({
          id: idItem,
          slug: produk.id,
          nama: produk.nama,
          varian: produk.varian,
          beratBersih: produk.beratBersih || '450 gr',
          harga: produk.hargaAngka || 28000,
          jumlah: qtyBeli,
          stokTersedia: produk.stok,
          gambar: produk.gambar || '/images/products/rengginang-udang.jpg'
        });
      }

      localStorage.setItem(STORAGE_KEY, JSON.stringify(currentItems));
      window.dispatchEvent(new Event('storage'));
    } catch (e) {
      console.warn('Gagal menyimpan ke localStorage keranjang:', e);
    }

    // Arahkan pelanggan ke halaman Keranjang Belanja (PPL-11)
    navigate('/keranjang');
  };

  // Validasi ketersediaan stok & jumlah pembelian (PPL-10)
  const apakahStokHabis =
    !produkDipilih ||
    produkDipilih.stok <= 0 ||
    produkDipilih.statusKetersediaan === 'STOK HABIS' ||
    produkDipilih.bisaDibeli === false;

  const apakahMelebihiStok = !apakahStokHabis && jumlahBeli > produkDipilih.stok;
  const apakahJumlahKurang = !apakahStokHabis && jumlahBeli < 1;
  const apakahJumlahTidakValid = apakahMelebihiStok || apakahJumlahKurang;
  const tombolDisabled = apakahStokHabis || apakahJumlahTidakValid;

  const hargaPerBungkus = produkDipilih?.hargaAngka || 28000;
  const totalHarga = hargaPerBungkus * (apakahStokHabis ? 0 : jumlahBeli || 0);

  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  // Handler untuk banner informasi penutup katalog
  const handleHubungiPenjual = () => {
    setJudulModal('Hubungi Penjual');
    setBukaModalKontak(true);
  };

  const handleLihatInfoUsaha = () => {
    navigate('/');
    setTimeout(() => {
      const el = document.getElementById('tentang');
      if (el) el.scrollIntoView({ behavior: 'smooth' });
    }, 100);
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '85vh', paddingBottom: '4rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            TAMPILAN DETAIL PRODUK (PPL-10: Sesuai Gambar Mockup Figma)
            ====================================================================== */}
        {produkDipilih ? (
          <div>
            {/* Breadcrumb Sesuai Figma PPL-10 */}
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

            {/* Header Detail Produk Sesuai Figma PPL-10 */}
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
                Pilih varian dan jumlah pembelian sesuai kebutuhan Anda.
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
              {/* KOLOM KIRI: Foto Produk Kemasan Sesuai Figma PPL-10 */}
              <div
                style={{
                  backgroundColor: '#FFFFFF',
                  border: '1px solid #ECE4D9',
                  borderRadius: '24px',
                  padding: '1.25rem',
                  boxShadow: '0 4px 16px rgba(67, 40, 28, 0.04)',
                  display: 'flex',
                  flexDirection: 'column'
                }}
              >
                {/* Foto Kemasan Produk */}
                <div
                  style={{
                    backgroundColor: '#FAF5EE',
                    borderRadius: '16px',
                    overflow: 'hidden',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '440px',
                    position: 'relative',
                    border: '1px solid #EFE8DF'
                  }}
                >
                  <img
                    src={produkDipilih.gambarDetail || produkDipilih.gambar}
                    alt={produkDipilih.nama}
                    style={{
                      width: '100%',
                      height: '100%',
                      maxHeight: '480px',
                      objectFit: 'contain',
                      display: 'block'
                    }}
                    onError={(e) => {
                      e.target.src = produkDipilih.gambar;
                    }}
                  />
                </div>

                {/* Keterangan Foto Kemasan Sesuai Mockup */}
                <div
                  style={{
                    textAlign: 'center',
                    fontSize: '0.8rem',
                    color: '#8E7D73',
                    marginTop: '0.85rem'
                  }}
                >
                  Foto asli kemasan {produkDipilih.nama}
                </div>
              </div>

              {/* KOLOM KANAN: Kartu Rincian, Pemilih Varian, & Counter Jumlah Sesuai Figma PPL-10 */}
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
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.4rem',
                        backgroundColor: '#FFEBEE',
                        color: '#C62828',
                        fontSize: '0.75rem',
                        fontWeight: 800,
                        letterSpacing: '0.04em',
                        padding: '0.35rem 0.95rem',
                        borderRadius: '9999px'
                      }}
                    >
                      <span style={{ fontSize: '0.65rem' }}>●</span> STOK HABIS
                    </span>
                  ) : (
                    <span
                      style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.4rem',
                        backgroundColor: '#E8F5E9',
                        color: '#2E7D32',
                        fontSize: '0.75rem',
                        fontWeight: 800,
                        letterSpacing: '0.04em',
                        padding: '0.35rem 0.95rem',
                        borderRadius: '9999px'
                      }}
                    >
                      <span style={{ fontSize: '0.65rem' }}>●</span> TERSEDIA
                    </span>
                  )}
                </div>

                {/* Nama Produk & Deskripsi */}
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

                {/* Bagian Harga Jual Sesuai Figma PPL-10 */}
                <div style={{ marginBottom: '1.25rem' }}>
                  <div style={{ fontSize: '1.85rem', fontWeight: 800, color: '#43281C', lineHeight: 1 }}>
                    {produkDipilih.harga}
                  </div>
                  <div style={{ fontSize: '0.8rem', color: '#8E7D73', marginTop: '0.35rem' }}>
                    Harga per bungkus
                  </div>
                  <div style={{ borderBottom: '1px solid #EFE8DF', marginTop: '1.15rem' }} />
                </div>

                {/* 2 Kolom Metrik: Berat Bersih & Stok Tersedia Sesuai Figma */}
                <div style={{ display: 'flex', gap: '2rem', marginBottom: '1.5rem' }}>
                  <div>
                    <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.25rem' }}>
                      Berat bersih
                    </div>
                    <div style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C' }}>
                      {produkDipilih.beratBersih}
                    </div>
                  </div>

                  <div>
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

                {/* Pemilih Varian Sesuai Figma PPL-10 */}
                <div style={{ marginBottom: '1.5rem' }}>
                  <div
                    style={{
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      marginBottom: '0.65rem'
                    }}
                  >
                    Pilih varian
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

                {/* Bagian Jumlah Pembelian & Total Harga Sesuai Figma PPL-10 */}
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'flex-start',
                    marginBottom: '1.5rem',
                    gap: '1rem',
                    flexWrap: 'wrap'
                  }}
                >
                  {/* Sisi Kiri: Counter Jumlah Pembelian */}
                  <div>
                    <div
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: 700,
                        color: '#43281C',
                        marginBottom: '0.5rem'
                      }}
                    >
                      Jumlah pembelian
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.85rem' }}>
                      <div
                        style={{
                          display: 'inline-flex',
                          alignItems: 'center',
                          backgroundColor: '#FFFFFF',
                          border: apakahMelebihiStok ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                          borderRadius: '8px',
                          overflow: 'hidden'
                        }}
                      >
                        <button
                          type="button"
                          onClick={() => setJumlahBeli((prev) => Math.max(1, prev - 1))}
                          disabled={apakahStokHabis || jumlahBeli <= 1}
                          style={{
                            width: '40px',
                            height: '38px',
                            backgroundColor: '#FAF5EE',
                            border: 'none',
                            fontSize: '1.1rem',
                            fontWeight: 700,
                            color: '#43281C',
                            cursor: apakahStokHabis || jumlahBeli <= 1 ? 'not-allowed' : 'pointer',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            userSelect: 'none'
                          }}
                        >
                          -
                        </button>

                        <input
                          type="text"
                          inputMode="numeric"
                          pattern="[0-9]*"
                          value={apakahStokHabis ? 0 : jumlahBeli}
                          onChange={(e) => {
                            if (apakahStokHabis) return;
                            const val = e.target.value.replace(/[^0-9]/g, '');
                            setJumlahBeli(val === '' ? 0 : parseInt(val, 10));
                          }}
                          disabled={apakahStokHabis}
                          style={{
                            width: '56px',
                            height: '38px',
                            border: 'none',
                            textAlign: 'center',
                            fontSize: '1rem',
                            fontWeight: 800,
                            color: '#43281C',
                            outline: 'none',
                            backgroundColor: '#FFFFFF'
                          }}
                        />

                        <button
                          type="button"
                          onClick={() => setJumlahBeli((prev) => prev + 1)}
                          disabled={apakahStokHabis}
                          style={{
                            width: '40px',
                            height: '38px',
                            backgroundColor: '#FAF5EE',
                            border: 'none',
                            fontSize: '1.1rem',
                            fontWeight: 700,
                            color: '#43281C',
                            cursor: apakahStokHabis ? 'not-allowed' : 'pointer',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            userSelect: 'none'
                          }}
                        >
                          +
                        </button>
                      </div>

                      <span style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                        Maksimal {produkDipilih.stok} bungkus
                      </span>
                    </div>
                  </div>

                  {/* Sisi Kanan: Total Harga */}
                  <div style={{ textAlign: 'right' }}>
                    <div
                      style={{
                        fontSize: '0.775rem',
                        color: '#8E7D73',
                        marginBottom: '0.35rem'
                      }}
                    >
                      Total harga
                    </div>
                    <div
                      style={{
                        fontSize: '1.65rem',
                        fontWeight: 800,
                        color: '#43281C',
                        lineHeight: 1
                      }}
                    >
                      {formatRupiah(totalHarga)}
                    </div>
                  </div>
                </div>

                {/* Banner Status / Peringatan Stok Sesuai Figma PPL-10 */}
                <div style={{ marginBottom: '1.5rem' }}>
                  {apakahStokHabis ? (
                    <div
                      style={{
                        backgroundColor: '#FDF2F2',
                        border: '1px solid #FFCDD2',
                        borderRadius: '12px',
                        padding: '0.9rem 1.15rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.85rem'
                      }}
                    >
                      <div
                        style={{
                          width: '24px',
                          height: '24px',
                          borderRadius: '50%',
                          backgroundColor: '#C62828',
                          color: '#FFFFFF',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          fontWeight: 800,
                          fontSize: '0.85rem',
                          flexShrink: 0
                        }}
                      >
                        !
                      </div>
                      <div>
                        <div
                          style={{
                            fontSize: '0.875rem',
                            fontWeight: 800,
                            color: '#C62828',
                            marginBottom: '0.15rem'
                          }}
                        >
                          Produk sedang habis
                        </div>
                        <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                          Pembelian dinonaktifkan. Silakan pilih varian lain.
                        </div>
                      </div>
                    </div>
                  ) : apakahMelebihiStok ? (
                    <div
                      style={{
                        backgroundColor: '#FDF2F2',
                        border: '1px solid #FFCDD2',
                        borderRadius: '12px',
                        padding: '0.9rem 1.15rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.85rem'
                      }}
                    >
                      <div
                        style={{
                          width: '24px',
                          height: '24px',
                          borderRadius: '50%',
                          backgroundColor: '#C62828',
                          color: '#FFFFFF',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          fontWeight: 800,
                          fontSize: '0.85rem',
                          flexShrink: 0
                        }}
                      >
                        !
                      </div>
                      <div>
                        <div
                          style={{
                            fontSize: '0.875rem',
                            fontWeight: 800,
                            color: '#C62828',
                            marginBottom: '0.15rem'
                          }}
                        >
                          Jumlah melebihi stok tersedia
                        </div>
                        <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                          {validasiRealtime.pesan || `Maksimal pembelian ${produkDipilih.stok} bungkus. Kurangi jumlah pembelian.`}
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div
                      style={{
                        backgroundColor: '#EDF7ED',
                        border: '1px solid #C8E6C9',
                        borderRadius: '12px',
                        padding: '0.9rem 1.15rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.85rem'
                      }}
                    >
                      <CheckCircle2 size={22} style={{ color: '#2E7D32', flexShrink: 0 }} />
                      <div>
                        <div
                          style={{
                            fontSize: '0.875rem',
                            fontWeight: 800,
                            color: '#1B5E20',
                            marginBottom: '0.15rem'
                          }}
                        >
                          Stok tersedia
                        </div>
                        <div style={{ fontSize: '0.8rem', color: '#4E5D52' }}>
                          {validasiRealtime.keterangan || 'Produk dapat dibeli. Kuota stok tervalidasi secara realtime.'}
                        </div>
                      </div>
                    </div>
                  )}
                </div>

                {/* Tombol Aksi Sesuai Figma PPL-10 */}
                <div>
                  <button
                    type="button"
                    disabled={tombolDisabled}
                    onClick={() => handleBeli(produkDipilih)}
                    style={{
                      width: '100%',
                      backgroundColor: tombolDisabled ? '#BFB5AC' : '#43281C',
                      color: '#FFFFFF',
                      border: 'none',
                      borderRadius: '10px',
                      padding: '0.95rem 1.5rem',
                      fontSize: '0.975rem',
                      fontWeight: 700,
                      cursor: tombolDisabled ? 'not-allowed' : 'pointer',
                      transition: 'background-color 0.2s',
                      boxShadow: tombolDisabled ? 'none' : '0 2px 8px rgba(67, 40, 28, 0.15)'
                    }}
                    onMouseEnter={(e) => {
                      if (!tombolDisabled) e.target.style.backgroundColor = '#2F1B12';
                    }}
                    onMouseLeave={(e) => {
                      if (!tombolDisabled) e.target.style.backgroundColor = '#43281C';
                    }}
                  >
                    {apakahStokHabis
                      ? 'Stok Habis'
                      : apakahJumlahTidakValid
                      ? 'Jumlah Tidak Valid'
                      : 'Beli Produk'}
                  </button>

                  <p
                    style={{
                      fontSize: '0.775rem',
                      color: '#8E7D73',
                      textAlign: 'center',
                      marginTop: '0.65rem',
                      marginBottom: 0
                    }}
                  >
                    Pembelian diproses setelah jumlah dan stok tervalidasi.
                  </p>
                </div>
              </div>
            </div>
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
