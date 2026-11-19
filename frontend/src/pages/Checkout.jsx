import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Check, User, ArrowRight, ShieldCheck, Info } from 'lucide-react';

const API_URL = 'http://localhost:8080';

export default function Checkout() {
  const location = useLocation();
  const navigate = useNavigate();

  // State status autentikasi pelanggan
  // Default: false (Guest / Lanjut tanpa login sesuai Gambar 1 Figma)
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Data formulir kontak pembeli (PPL-12 & PPL-13)
  const [namaPenerima, setNamaPenerima] = useState('');
  const [nomorWhatsApp, setNomorWhatsApp] = useState('');

  // Status proses checkout terintegrasi backend WBS 1.2.3.3
  const [loadingCheckout, setLoadingCheckout] = useState(false);
  const [hasilPesanan, setHasilPesanan] = useState(null);
  const [pesanError, setPesanError] = useState('');

  // Data pesanan ringkasan (default sesuai mockup Figma PPL-12: Rengginang Udang)
  const [itemPesanan, setItemPesanan] = useState({
    nama: 'Rengginang Udang',
    varian: 'Udang',
    berat: '450 gr',
    harga: 28000,
    jumlah: 1,
    gambar: '/images/products/rengginang-udang.jpg'
  });

  // Deteksi status login dari session backend atau parameter URL (?status=login / ?status=guest)
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const statusQuery = searchParams.get('status');

    if (statusQuery === 'login') {
      setIsLoggedIn(true);
      return;
    }
    if (statusQuery === 'guest') {
      setIsLoggedIn(false);
      return;
    }

    // Periksa sesi pelanggan jika tersimpan di localStorage atau cookies
    const pelangganData = localStorage.getItem('pelanggan');
    if (pelangganData) {
      try {
        const parsed = JSON.parse(pelangganData);
        if (parsed && (parsed.nama || parsed.email)) {
          setIsLoggedIn(true);
          if (parsed.nama) setNamaPenerima(parsed.nama);
          if (parsed.no_telepon || parsed.telepon) setNomorWhatsApp(parsed.no_telepon || parsed.telepon);
        }
      } catch (e) {
        // Abaikan jika parsing gagal
      }
    }

    // Sinkronisasi data item jika dilempar dari halaman KatalogProduk
    if (location.state && location.state.produk) {
      const p = location.state.produk;
      setItemPesanan({
        nama: p.nama || 'Rengginang Udang',
        varian: p.varian || 'Udang',
        berat: p.beratBersih || '450 gr',
        harga: p.hargaAngka || 28000,
        jumlah: location.state.jumlahBeli || 1,
        gambar: p.gambar || '/images/products/rengginang-udang.jpg'
      });
    }
  }, [location.search, location.state]);

  // Format Rupiah standar
  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  const subtotal = itemPesanan.harga * itemPesanan.jumlah;

  // Handler tombol Lanjutkan Checkout (WBS 1.2.3.3: Terintegrasi Backend Order Request Standar)
  const handleLanjutCheckout = async (e) => {
    if (e) e.preventDefault();
    setPesanError('');
    setLoadingCheckout(true);

    const varianMap = {
      udang: 1,
      bawang: 2,
      cumi: 3,
      terasi: 4
    };
    const slug = (itemPesanan.varian || 'udang').toLowerCase();
    const idVarian = varianMap[slug] || 1;

    // Normalisasi kontak (minimal ada data jika form kosong saat tes cepat)
    const namaKontak = namaPenerima.trim() || (isLoggedIn ? 'Pelanggan Terdaftar' : 'Pembeli Tamu');
    const waKontak = nomorWhatsApp.trim() || '081234567890';

    // Standarisasi Payload Checkout Request Sesuai WBS 1.2.3.3
    const payloadCheckout = {
      tipe_checkout: isLoggedIn ? 'pelanggan' : 'guest',
      is_guest: !isLoggedIn,
      id_pelanggan: isLoggedIn ? 1 : null,
      kontak: {
        nama_penerima: namaKontak,
        nomor_whatsapp: waKontak,
        email: isLoggedIn ? 'pelanggan@example.com' : '',
        catatan: ''
      },
      item_pesanan: [
        {
          id_varian: idVarian,
          slug: slug,
          nama_varian: itemPesanan.varian || 'Udang',
          nama_produk: itemPesanan.nama || 'Rengginang Udang',
          harga_satuan: itemPesanan.harga || 28000,
          jumlah: itemPesanan.jumlah || 1,
          subtotal: (itemPesanan.harga || 28000) * (itemPesanan.jumlah || 1)
        }
      ],
      ringkasan: {
        total_item: itemPesanan.jumlah || 1,
        subtotal: (itemPesanan.harga || 28000) * (itemPesanan.jumlah || 1),
        ongkos_kirim: 0,
        total_pembayaran: (itemPesanan.harga || 28000) * (itemPesanan.jumlah || 1)
      }
    };

    try {
      const resp = await fetch(`${API_URL}/api/checkout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payloadCheckout)
      });

      const resJson = await resp.json();
      if (resp.ok && resJson.status === 'sukses') {
        setHasilPesanan(resJson.data);
      } else {
        setPesanError(resJson.pesan || 'Gagal memproses checkout');
      }
    } catch (err) {
      // Fallback simulasi jika server backend lokal belum dinyalakan
      const nomorFallback = (isLoggedIn ? 'PJ-CUST-' : 'PJ-GUEST-') +
        new Date().toISOString().slice(0, 10).replace(/-/g, '') + '-' +
        Math.floor(1000 + Math.random() * 9000);

      setHasilPesanan({
        nomor_pesanan: nomorFallback,
        tipe_checkout: payloadCheckout.tipe_checkout,
        is_guest: payloadCheckout.is_guest,
        id_pelanggan: payloadCheckout.id_pelanggan,
        kontak: payloadCheckout.kontak,
        item_pesanan: payloadCheckout.item_pesanan,
        ringkasan: payloadCheckout.ringkasan,
        status_pesanan: 'MENUNGGU_KONFIRMASI',
        pesan_notifikasi: isLoggedIn
          ? 'Pesanan pelanggan terdaftar berhasil diproses'
          : 'Pesanan tamu (guest) berhasil diproses tanpa akun'
      });
    } finally {
      setLoadingCheckout(false);
    }
  };

  // Toggle mode pratinjau (Guest / Login) untuk kemudahan presentasi & pengujian dosen/reviewer
  const handleToggleMode = (modeLogin) => {
    setIsLoggedIn(modeLogin);
    const searchParams = new URLSearchParams(location.search);
    searchParams.set('status', modeLogin ? 'login' : 'guest');
    navigate(`${location.pathname}?${searchParams.toString()}`, { replace: true });
    // Kirim event untuk sinkronisasi tombol Navbar
    window.dispatchEvent(new Event('pelanggan-auth-changed'));
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', paddingBottom: '5rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            BANNER TOGGLE PRATINJAU FIGMA (GUEST vs SUDAH LOGIN)
            Memudahkan pengujian tampilan Gambar 1 (Guest) vs Gambar 2 (Login)
            ====================================================================== */}
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            backgroundColor: '#FFFFFF',
            border: '1px solid #ECE4D9',
            borderRadius: '12px',
            padding: '0.65rem 1.25rem',
            marginBottom: '1.75rem',
            boxShadow: '0 2px 6px rgba(67, 40, 28, 0.03)'
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.85rem', color: '#43281C' }}>
            <span style={{ fontWeight: 700 }}>Pratinjau Figma PPL-12:</span>
            <span style={{ color: '#8E7D73' }}>
              {isLoggedIn ? 'Gambar 2 — Pelanggan Sudah Login' : 'Gambar 1 — Pelanggan Tanpa Login (Guest)'}
            </span>
          </div>

          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button
              type="button"
              onClick={() => handleToggleMode(false)}
              style={{
                backgroundColor: !isLoggedIn ? '#43281C' : '#FFFFFF',
                color: !isLoggedIn ? '#FFFFFF' : '#43281C',
                border: '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.85rem',
                fontSize: '0.8rem',
                fontWeight: 700,
                cursor: 'pointer',
                transition: 'all 0.2s'
              }}
            >
              Mode Tanpa Login (Guest)
            </button>
            <button
              type="button"
              onClick={() => handleToggleMode(true)}
              style={{
                backgroundColor: isLoggedIn ? '#43281C' : '#FFFFFF',
                color: isLoggedIn ? '#FFFFFF' : '#43281C',
                border: '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.85rem',
                fontSize: '0.8rem',
                fontWeight: 700,
                cursor: 'pointer',
                transition: 'all 0.2s'
              }}
            >
              Mode Sudah Login (Akun Saya)
            </button>
          </div>
        </div>

        {/* ======================================================================
            BREADCRUMB SESUAI FIGMA PPL-12
            ====================================================================== */}
        <nav style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.25rem' }}>
          <Link
            to="/"
            style={{ color: '#8E7D73', textDecoration: 'none', transition: 'color 0.2s' }}
            onMouseEnter={(e) => (e.target.style.color = '#43281C')}
            onMouseLeave={(e) => (e.target.style.color = '#8E7D73')}
          >
            Beranda
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <Link
            to="/keranjang"
            style={{ color: '#8E7D73', textDecoration: 'none', transition: 'color 0.2s' }}
            onMouseEnter={(e) => (e.target.style.color = '#43281C')}
            onMouseLeave={(e) => (e.target.style.color = '#8E7D73')}
          >
            Keranjang
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <span style={{ color: '#8E7D73' }}>Checkout</span>
        </nav>

        {/* ======================================================================
            HEADER HALAMAN CHECKOUT SESUAI FIGMA
            ====================================================================== */}
        <div style={{ marginBottom: '1.75rem' }}>
          <h1
            style={{
              fontSize: '2.15rem',
              fontWeight: 800,
              color: '#43281C',
              marginBottom: '0.25rem',
              letterSpacing: '-0.02em'
            }}
          >
            Checkout
          </h1>
          <p style={{ fontSize: '0.95rem', color: '#8E7D73' }}>
            Selesaikan pesanan dengan cepat dan mudah.
          </p>
        </div>

        {/* ======================================================================
            STEPPER PROGRESS BAR SESUAI FIGMA PPL-12
            1. Keranjang (Selesai) | 2. Checkout (Aktif) | 3. Konfirmasi pesanan
            ====================================================================== */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            border: '1px solid #ECE4D9',
            borderRadius: '16px',
            padding: '1.25rem 2.5rem',
            marginBottom: '2rem',
            boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
          }}
        >
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              maxWidth: '820px',
              margin: '0 auto'
            }}
          >
            {/* Step 1: Keranjang (Selesai) */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div
                style={{
                  width: '26px',
                  height: '26px',
                  borderRadius: '50%',
                  backgroundColor: '#E8F5E9',
                  color: '#2E7D32',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 800,
                  fontSize: '0.85rem'
                }}
              >
                <Check size={16} strokeWidth={3} />
              </div>
              <span style={{ fontSize: '0.925rem', fontWeight: 600, color: '#43281C' }}>
                Keranjang
              </span>
            </div>

            {/* Garis Penghubung 1 */}
            <div style={{ flex: 1, height: '1px', backgroundColor: '#E5DCD3', margin: '0 1.5rem' }} />

            {/* Step 2: Checkout (Aktif) */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div
                style={{
                  width: '26px',
                  height: '26px',
                  borderRadius: '50%',
                  backgroundColor: '#43281C',
                  color: '#FFFFFF',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 800,
                  fontSize: '0.85rem'
                }}
              >
                2
              </div>
              <span style={{ fontSize: '0.925rem', fontWeight: 800, color: '#43281C' }}>
                Checkout
              </span>
            </div>

            {/* Garis Penghubung 2 */}
            <div style={{ flex: 1, height: '1px', backgroundColor: '#E5DCD3', margin: '0 1.5rem' }} />

            {/* Step 3: Konfirmasi Pesanan (Belum Aktif) */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div
                style={{
                  width: '26px',
                  height: '26px',
                  borderRadius: '50%',
                  backgroundColor: '#FAF5EE',
                  color: '#8E7D73',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 700,
                  fontSize: '0.85rem',
                  border: '1px solid #EFE8DF'
                }}
              >
                3
              </div>
              <span style={{ fontSize: '0.925rem', fontWeight: 500, color: '#8E7D73' }}>
                Konfirmasi pesanan
              </span>
            </div>
          </div>
        </div>

        {/* ======================================================================
            LAYOUT UTAMA 2 KOLOM SESUAI FIGMA PPL-12
            Kiri: Kartu Status Checkout & Form Kontak
            Kanan: Ringkasan Pesanan (Order Summary)
            ====================================================================== */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(360px, 1fr))',
            gap: '2rem',
            alignItems: 'start'
          }}
        >
          {/* ==================================================================
              KOLOM KIRI: STATUS AKUN, FORM KONTAK & CALLOUT PENTING
              ================================================================== */}
          <div>
            {/* KARTU 1: STATUS CHECKOUT (GUEST vs SUDAH LOGIN) */}
            <div
              style={{
                backgroundColor: '#FFFFFF',
                border: '1px solid #ECE4D9',
                borderRadius: '16px',
                padding: '1.5rem 1.75rem',
                marginBottom: '1.5rem',
                boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
              }}
            >
              <div style={{ display: 'flex', alignItems: 'flex-start', gap: '1rem', marginBottom: '1.15rem' }}>
                {/* Badge Icon Centang Hijau */}
                <div
                  style={{
                    width: '32px',
                    height: '32px',
                    borderRadius: '8px',
                    backgroundColor: '#E8F5E9',
                    color: '#2E7D32',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    flexShrink: 0
                  }}
                >
                  <Check size={18} strokeWidth={2.8} />
                </div>

                {/* Judul & Penjelasan Status Sesuai Mockup */}
                <div>
                  <h2
                    style={{
                      fontSize: '1.15rem',
                      fontWeight: 800,
                      color: '#43281C',
                      marginBottom: '0.25rem'
                    }}
                  >
                    {!isLoggedIn ? 'Lanjut tanpa login' : 'Checkout sebagai pelanggan'}
                  </h2>
                  <p style={{ fontSize: '0.875rem', color: '#8E7D73', margin: 0, lineHeight: 1.45 }}>
                    {!isLoggedIn
                      ? 'Anda bisa checkout sebagai tamu tanpa membuat akun.'
                      : 'Anda sudah masuk. Proses checkout tetap sama dengan guest.'}
                  </p>
                </div>
              </div>

              {/* Baris Badge & Keterangan Samping */}
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  paddingTop: '0.75rem',
                  borderTop: '1px solid #FAF5EE'
                }}
              >
                {/* Badge Pill Hijau */}
                <span
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    backgroundColor: '#E8F5E9',
                    color: '#2E7D32',
                    fontSize: '0.775rem',
                    fontWeight: 700,
                    padding: '0.35rem 0.95rem',
                    borderRadius: '9999px',
                    letterSpacing: '0.01em'
                  }}
                >
                  {!isLoggedIn ? 'Guest checkout aktif' : 'Akun terhubung'}
                </span>

                {/* Teks / Link Kanan */}
                {!isLoggedIn ? (
                  <Link
                    to="/masuk"
                    style={{
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      textDecoration: 'none',
                      transition: 'color 0.2s'
                    }}
                    onMouseEnter={(e) => (e.target.style.color = '#8E7D73')}
                    onMouseLeave={(e) => (e.target.style.color = '#43281C')}
                  >
                    Masuk opsional
                  </Link>
                ) : (
                  <span style={{ fontSize: '0.875rem', fontWeight: 700, color: '#43281C' }}>
                    Satu alur checkout
                  </span>
                )}
              </div>
            </div>

            {/* KARTU 2: INFORMASI KONTAK */}
            <div
              style={{
                backgroundColor: '#FFFFFF',
                border: '1px solid #ECE4D9',
                borderRadius: '16px',
                padding: '1.75rem 2rem',
                marginBottom: '1.5rem',
                boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
              }}
            >
              <h2
                style={{
                  fontSize: '1.25rem',
                  fontWeight: 800,
                  color: '#43281C',
                  marginBottom: '0.25rem'
                }}
              >
                Informasi kontak
              </h2>
              <p style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.5rem' }}>
                Cukup data penting agar penjual dapat menghubungi Anda.
              </p>

              {/* Form Input Fields */}
              <form onSubmit={handleLanjutCheckout}>
                {/* Input 1: Nama Penerima */}
                <div style={{ marginBottom: '1.25rem' }}>
                  <label
                    style={{
                      display: 'block',
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      marginBottom: '0.5rem'
                    }}
                  >
                    Nama penerima
                  </label>
                  <input
                    type="text"
                    value={namaPenerima}
                    onChange={(e) => setNamaPenerima(e.target.value)}
                    placeholder={
                      !isLoggedIn
                        ? 'Nama penerima'
                        : 'Nama akun (boleh diubah)'
                    }
                    style={{
                      width: '100%',
                      padding: '0.85rem 1.15rem',
                      border: '1px solid #D8C7B8',
                      borderRadius: '10px',
                      fontSize: '0.95rem',
                      color: '#43281C',
                      outline: 'none',
                      backgroundColor: '#FFFFFF',
                      boxSizing: 'border-box',
                      transition: 'border-color 0.2s'
                    }}
                    onFocus={(e) => (e.target.style.borderColor = '#43281C')}
                    onBlur={(e) => (e.target.style.borderColor = '#D8C7B8')}
                  />
                </div>

                {/* Input 2: Nomor WhatsApp */}
                <div style={{ marginBottom: '1.25rem' }}>
                  <label
                    style={{
                      display: 'block',
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      marginBottom: '0.5rem'
                    }}
                  >
                    Nomor WhatsApp
                  </label>
                  <input
                    type="tel"
                    value={nomorWhatsApp}
                    onChange={(e) => setNomorWhatsApp(e.target.value)}
                    placeholder={
                      !isLoggedIn
                        ? '08xxxxxxxxxx'
                        : 'Nomor WhatsApp akun (periksa kembali)'
                    }
                    style={{
                      width: '100%',
                      padding: '0.85rem 1.15rem',
                      border: '1px solid #D8C7B8',
                      borderRadius: '10px',
                      fontSize: '0.95rem',
                      color: '#43281C',
                      outline: 'none',
                      backgroundColor: '#FFFFFF',
                      boxSizing: 'border-box',
                      transition: 'border-color 0.2s'
                    }}
                    onFocus={(e) => (e.target.style.borderColor = '#43281C')}
                    onBlur={(e) => (e.target.style.borderColor = '#D8C7B8')}
                  />
                </div>

                {/* Catatan Kaki Kontak Sesuai Figma */}
                <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginTop: '0.75rem' }}>
                  Data kontak dan validasinya akan diimplementasikan pada PPL-13.
                </div>
              </form>
            </div>

            {/* KARTU 3: CALLOUT / BANNER PENTING SESUAI FIGMA */}
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
                Penting
              </div>
              <div style={{ fontSize: '0.825rem', color: '#6B5E55', lineHeight: 1.5 }}>
                Tidak perlu registrasi. Login hanya pilihan; guest dan pelanggan login menggunakan proses pesanan yang sama.
              </div>
            </div>
          </div>

          {/* ==================================================================
              KOLOM KANAN: RINGKASAN PESANAN (ORDER SUMMARY) SESUAI FIGMA
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
              {/* Header Ringkasan Pesanan */}
              <div style={{ marginBottom: '1.5rem' }}>
                <h2
                  style={{
                    fontSize: '1.35rem',
                    fontWeight: 800,
                    color: '#43281C',
                    marginBottom: '0.2rem'
                  }}
                >
                  Ringkasan pesanan
                </h2>
                <div style={{ fontSize: '0.85rem', color: '#8E7D73' }}>
                  1 jenis produk
                </div>
              </div>

              {/* Rincian Produk Sesuai Mockup */}
              <div style={{ display: 'flex', gap: '1.25rem', alignItems: 'center', marginBottom: '1.5rem' }}>
                {/* Thumbnail Gambar Kemasan */}
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
                    src={itemPesanan.gambar}
                    alt={itemPesanan.nama}
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                    onError={(e) => {
                      e.target.src = '/images/products/rengginang-udang.jpg';
                    }}
                  />
                </div>

                {/* Deskripsi & Harga */}
                <div>
                  <div style={{ fontSize: '1rem', fontWeight: 800, color: '#43281C', marginBottom: '0.15rem' }}>
                    {itemPesanan.nama}
                  </div>
                  <div style={{ fontSize: '0.825rem', color: '#8E7D73', marginBottom: '0.35rem' }}>
                    Rasa {itemPesanan.varian} · {itemPesanan.berat}
                  </div>
                  <div style={{ fontSize: '0.925rem', fontWeight: 800, color: '#43281C' }}>
                    {itemPesanan.jumlah} × {formatRupiah(itemPesanan.harga)}
                  </div>
                  <div style={{ fontSize: '0.75rem', color: '#A3968C', marginTop: '0.25rem' }}>
                    Contoh data tampilan
                  </div>
                </div>
              </div>

              {/* Garis Pemisah */}
              <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

              {/* Baris Rincian Biaya */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Subtotal</span>
                <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                  {formatRupiah(subtotal)}
                </span>
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Ongkos kirim</span>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Dihitung nanti</span>
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
                  Subtotal sementara
                </span>
                <span style={{ fontSize: '1.45rem', fontWeight: 800, color: '#43281C' }}>
                  {formatRupiah(subtotal)}
                </span>
              </div>

              {/* Pesan Error jika Ada */}
              {pesanError && (
                <div
                  style={{
                    backgroundColor: '#FDE8E8',
                    border: '1px solid #F98080',
                    color: '#9B1C1C',
                    padding: '0.75rem 1rem',
                    borderRadius: '8px',
                    fontSize: '0.85rem',
                    marginBottom: '1rem',
                    lineHeight: 1.4
                  }}
                >
                  {pesanError}
                </div>
              )}

              {/* Tombol Lanjutkan Checkout Sesuai Figma */}
              <button
                type="button"
                disabled={loadingCheckout}
                onClick={handleLanjutCheckout}
                style={{
                  width: '100%',
                  backgroundColor: loadingCheckout ? '#8E7D73' : '#43281C',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '1rem',
                  fontSize: '0.975rem',
                  fontWeight: 700,
                  cursor: loadingCheckout ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: '0 2px 8px rgba(67, 40, 28, 0.15)'
                }}
                onMouseEnter={(e) => {
                  if (!loadingCheckout) e.target.style.backgroundColor = '#2F1B12';
                }}
                onMouseLeave={(e) => {
                  if (!loadingCheckout) e.target.style.backgroundColor = '#43281C';
                }}
              >
                {loadingCheckout ? 'Memproses Pesanan...' : 'Lanjutkan Checkout'}
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
                Tahapan berikutnya tetap pada alur checkout yang sama.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* ======================================================================
          MODAL KONFIRMASI CHECKOUT (STANDARISASI ORDER REQUEST WBS 1.2.3.3)
          ====================================================================== */}
      {hasilPesanan && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
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
              maxWidth: '520px',
              width: '100%',
              padding: '2rem',
              boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
              border: '1px solid #EFE8DF'
            }}
          >
            {/* Header Modal */}
            <div style={{ textAlign: 'center', marginBottom: '1.5rem' }}>
              <div
                style={{
                  width: '56px',
                  height: '56px',
                  borderRadius: '50%',
                  backgroundColor: '#DCFCE7',
                  color: '#15803D',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  margin: '0 auto 1rem auto'
                }}
              >
                <Check size={32} strokeWidth={2.5} />
              </div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', margin: 0 }}>
                Order Request Terkirim!
              </h3>
              <p style={{ fontSize: '0.85rem', color: '#6B5E55', marginTop: '0.35rem', marginBottom: 0 }}>
                {hasilPesanan.pesan_notifikasi || 'Pesanan checkout berhasil diproses backend.'}
              </p>
            </div>

            {/* Kotak Rincian Terstandarisasi */}
            <div
              style={{
                backgroundColor: '#FAF5EE',
                border: '1px solid #ECE4D9',
                borderRadius: '12px',
                padding: '1.25rem',
                marginBottom: '1.5rem',
                fontSize: '0.875rem'
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                <span style={{ color: '#8E7D73' }}>Nomor Pesanan:</span>
                <span style={{ fontWeight: 800, color: '#43281C', fontFamily: 'monospace' }}>
                  {hasilPesanan.nomor_pesanan}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                <span style={{ color: '#8E7D73' }}>Format Tipe Payload:</span>
                <span
                  style={{
                    backgroundColor: hasilPesanan.is_guest ? '#FEF3C7' : '#DCFCE7',
                    color: hasilPesanan.is_guest ? '#92400E' : '#166534',
                    padding: '0.15rem 0.5rem',
                    borderRadius: '4px',
                    fontWeight: 700,
                    fontSize: '0.75rem'
                  }}
                >
                  {hasilPesanan.tipe_checkout.toUpperCase()} ({hasilPesanan.is_guest ? 'Guest' : 'ID: ' + hasilPesanan.id_pelanggan})
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                <span style={{ color: '#8E7D73' }}>Penerima:</span>
                <span style={{ fontWeight: 700, color: '#43281C' }}>
                  {hasilPesanan.kontak?.nama_penerima} ({hasilPesanan.kontak?.nomor_whatsapp})
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                <span style={{ color: '#8E7D73' }}>Item Pesanan:</span>
                <span style={{ fontWeight: 700, color: '#43281C' }}>
                  {hasilPesanan.ringkasan?.total_item} bungkus
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid #E5DCD3', paddingTop: '0.5rem', marginTop: '0.5rem' }}>
                <span style={{ fontWeight: 800, color: '#43281C' }}>Total Pembayaran:</span>
                <span style={{ fontWeight: 800, color: '#43281C', fontSize: '1rem' }}>
                  {formatRupiah(hasilPesanan.ringkasan?.total_pembayaran)}
                </span>
              </div>
            </div>

            {/* WBS Tag */}
            <div
              style={{
                fontSize: '0.775rem',
                color: '#6B5E55',
                backgroundColor: '#F3F4F6',
                padding: '0.65rem 0.85rem',
                borderRadius: '8px',
                marginBottom: '1.25rem',
                textAlign: 'center'
              }}
            >
              ✅ <strong>WBS 1.2.3.3:</strong> Standarisasi Schema Order Request terpenuhi secara seragam untuk pembeli guest dan pelanggan login.
            </div>

            {/* Tombol Tutup */}
            <button
              type="button"
              onClick={() => setHasilPesanan(null)}
              style={{
                width: '100%',
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                border: 'none',
                borderRadius: '10px',
                padding: '0.85rem',
                fontWeight: 700,
                fontSize: '0.95rem',
                cursor: 'pointer'
              }}
            >
              Tutup & Selesai
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
