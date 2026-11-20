import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Check } from 'lucide-react';

const API_URL = 'http://localhost:8080';

export default function Checkout() {
  const location = useLocation();
  const navigate = useNavigate();

  // Status login pelanggan (true = Pelanggan login, false = Guest)
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Field input data kontak (PPL-13)
  const [namaPenerima, setNamaPenerima] = useState('');
  const [nomorWhatsApp, setNomorWhatsApp] = useState('');

  // Status interaksi user (apakah field sudah pernah disentuh / disubmit)
  const [namaDisentuh, setNamaDisentuh] = useState(false);
  const [nomorDisentuh, setNomorDisentuh] = useState(false);
  const [isAutoFilled, setIsAutoFilled] = useState(false);

  // State untuk pengiriman data / pesanan
  const [loadingCheckout, setLoadingCheckout] = useState(false);
  const [hasilPesanan, setHasilPesanan] = useState(null);

  // Data ringkasan produk checkout (default sesuai mockup Figma: Rengginang Udang)
  const [itemPesanan, setItemPesanan] = useState({
    nama: 'Rengginang Udang',
    varian: 'Udang',
    berat: '450 gr',
    harga: 28000,
    jumlah: 1,
    gambar: '/images/products/rengginang-udang.jpg'
  });

  // Deteksi status login awal dan data keranjang jika ada
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const statusQuery = searchParams.get('status');

    if (statusQuery === 'login') {
      terapkanPreset('login_autofill');
      return;
    } else if (statusQuery === 'guest_valid') {
      terapkanPreset('guest_valid');
      return;
    } else if (statusQuery === 'error') {
      terapkanPreset('guest_error');
      return;
    } else if (statusQuery === 'guest') {
      terapkanPreset('guest_kosong');
      return;
    }

    // Ambil data item dari location.state jika dikirim dari katalog atau keranjang
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

  // Fungsi validasi nomor telepon seluler Indonesia
  const cekNomorHpValid = (nomor) => {
    if (!nomor) return false;
    // Format diawali 08 atau +62 atau 62, total 10 sampai 15 digit angka
    const pola = /^(\+62|62|08)[0-9]{8,13}$/;
    return pola.test(nomor.trim().replace(/\s|-/g, ''));
  };

  // Status validasi masing-masing field
  const namaTerisi = namaPenerima.trim().length > 0;
  const nomorTerisi = nomorWhatsApp.trim().length > 0;
  const nomorValid = cekNomorHpValid(nomorWhatsApp);

  // Deteksi apakah sedang dalam kondisi error (State 2)
  const adaErrorNama = namaDisentuh && !namaTerisi;
  const adaErrorNomor = nomorDisentuh && (!nomorTerisi || !nomorValid);
  const adaErrorForm = adaErrorNama || adaErrorNomor;

  // Form dianggap valid jika nama terisi dan format nomor HP valid (State 3 & 4)
  const formValid = namaTerisi && nomorValid;

  // Format angka ke Rupiah
  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  const subtotal = itemPesanan.harga * itemPesanan.jumlah;

  // Preset pengujian 4 state sesuai mockup Figma PPL-13
  const terapkanPreset = (mode) => {
    if (mode === 'guest_kosong') {
      // Gambar 2: Form Kosong (Guest)
      setIsLoggedIn(false);
      setNamaPenerima('');
      setNomorWhatsApp('');
      setNamaDisentuh(false);
      setNomorDisentuh(false);
      setIsAutoFilled(false);
    } else if (mode === 'guest_error') {
      // Gambar 3: Validasi Error
      setIsLoggedIn(false);
      setNamaPenerima('');
      setNomorWhatsApp('0812abc');
      setNamaDisentuh(true);
      setNomorDisentuh(true);
      setIsAutoFilled(false);
    } else if (mode === 'guest_valid') {
      // Gambar 4: Data Valid Sebagai Guest
      setIsLoggedIn(false);
      setNamaPenerima('Contoh Pelanggan');
      setNomorWhatsApp('081234567890');
      setNamaDisentuh(true);
      setNomorDisentuh(true);
      setIsAutoFilled(false);
    } else if (mode === 'login_autofill') {
      // Gambar 5: Checkout Pelanggan Login (Auto-fill)
      setIsLoggedIn(true);
      setNamaPenerima('Contoh Pelanggan');
      setNomorWhatsApp('081234567890');
      setNamaDisentuh(false);
      setNomorDisentuh(false);
      setIsAutoFilled(true);
    }
  };

  // Handler tombol Lanjutkan Checkout
  const handleLanjutCheckout = async (e) => {
    if (e) e.preventDefault();

    // Jika form belum valid, aktifkan tanda error
    if (!formValid) {
      setNamaDisentuh(true);
      setNomorDisentuh(true);
      return;
    }

    setLoadingCheckout(true);

    const payloadCheckout = {
      tipe_checkout: isLoggedIn ? 'pelanggan' : 'guest',
      is_guest: !isLoggedIn,
      kontak: {
        nama_penerima: namaPenerima.trim(),
        nomor_whatsapp: nomorWhatsApp.trim()
      },
      item_pesanan: [
        {
          nama_produk: itemPesanan.nama,
          varian: itemPesanan.varian,
          harga_satuan: itemPesanan.harga,
          jumlah: itemPesanan.jumlah,
          subtotal: subtotal
        }
      ],
      subtotal: subtotal
    };

    try {
      const resp = await fetch(`${API_URL}/api/checkout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payloadCheckout)
      });

      const resJson = await resp.json();
      if (resp.ok && resJson.status === 'sukses') {
        setHasilPesanan(resJson.data);
      } else {
        // Fallback simulasi jika backend belum menangani route checkout
        simulasikanHasil(payloadCheckout);
      }
    } catch (err) {
      // Fallback offline
      simulasikanHasil(payloadCheckout);
    } finally {
      setLoadingCheckout(false);
    }
  };

  const simulasikanHasil = (payload) => {
    const nomorPesanan = 'PJ-ORD-' + new Date().toISOString().slice(0, 10).replace(/-/g, '') + '-' + Math.floor(1000 + Math.random() * 9000);
    setHasilPesanan({
      nomor_pesanan: nomorPesanan,
      tipe_checkout: payload.tipe_checkout,
      kontak: payload.kontak,
      total_pembayaran: subtotal,
      pesan: 'Data kontak berhasil divalidasi. Lanjut ke pemilihan metode pengiriman (PPL-14).'
    });
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', paddingBottom: '5rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            TOOLBAR NAVIGASI 4 GAMBAR FIGMA PPL-13 (Untuk Kemudahan Review & Cek)
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
            boxShadow: '0 2px 6px rgba(67, 40, 28, 0.03)',
            flexWrap: 'wrap',
            gap: '0.65rem'
          }}
        >
          <div style={{ fontSize: '0.85rem', fontWeight: 700, color: '#43281C' }}>
            Pratinjau 4 Kondisi Mockup Figma PPL-13:
          </div>

          <div style={{ display: 'flex', gap: '0.45rem', flexWrap: 'wrap' }}>
            <button
              type="button"
              onClick={() => terapkanPreset('guest_kosong')}
              style={{
                backgroundColor: !isLoggedIn && !namaTerisi && !nomorTerisi && !nomorDisentuh ? '#FAF5EE' : '#FFFFFF',
                color: '#43281C',
                border: !isLoggedIn && !namaTerisi && !nomorTerisi && !nomorDisentuh ? '1.5px solid #43281C' : '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.75rem',
                fontSize: '0.8rem',
                fontWeight: 600,
                cursor: 'pointer'
              }}
            >
              1. Form Kosong (Guest)
            </button>
            <button
              type="button"
              onClick={() => terapkanPreset('guest_error')}
              style={{
                backgroundColor: adaErrorForm ? '#FDF2F2' : '#FFFFFF',
                color: adaErrorForm ? '#C62828' : '#43281C',
                border: adaErrorForm ? '1.5px solid #C62828' : '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.75rem',
                fontSize: '0.8rem',
                fontWeight: 600,
                cursor: 'pointer'
              }}
            >
              2. Validasi Error
            </button>
            <button
              type="button"
              onClick={() => terapkanPreset('guest_valid')}
              style={{
                backgroundColor: !isLoggedIn && formValid ? '#EDF7ED' : '#FFFFFF',
                color: !isLoggedIn && formValid ? '#1B5E20' : '#43281C',
                border: !isLoggedIn && formValid ? '1.5px solid #1B5E20' : '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.75rem',
                fontSize: '0.8rem',
                fontWeight: 600,
                cursor: 'pointer'
              }}
            >
              3. Data Valid (Guest)
            </button>
            <button
              type="button"
              onClick={() => terapkanPreset('login_autofill')}
              style={{
                backgroundColor: isLoggedIn ? '#FAF5EE' : '#FFFFFF',
                color: '#43281C',
                border: isLoggedIn ? '1.5px solid #43281C' : '1px solid #D8C7B8',
                borderRadius: '8px',
                padding: '0.35rem 0.75rem',
                fontSize: '0.8rem',
                fontWeight: 600,
                cursor: 'pointer'
              }}
            >
              4. Pelanggan Login (Auto-fill)
            </button>
          </div>
        </div>

        {/* ======================================================================
            BREADCRUMB SESUAI FIGMA (Beranda / Keranjang / Checkout)
            ====================================================================== */}
        <nav style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.25rem' }}>
          <Link
            to="/"
            style={{ color: '#8E7D73', textDecoration: 'none' }}
          >
            Beranda
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <Link
            to="/keranjang"
            style={{ color: '#8E7D73', textDecoration: 'none' }}
          >
            Keranjang
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <span style={{ color: '#8E7D73' }}>Checkout</span>
        </nav>

        {/* ======================================================================
            HEADER CHECKOUT SESUAI FIGMA
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
          <p style={{ fontSize: '0.95rem', color: '#8E7D73', margin: 0 }}>
            Selesaikan pesanan dengan cepat dan mudah.
          </p>
        </div>

        {/* ======================================================================
            STEPPER PROGRESS SESUAI FIGMA
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
                  width: '24px',
                  height: '24px',
                  borderRadius: '50%',
                  backgroundColor: '#E8F5E9',
                  color: '#2E7D32',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 800,
                  fontSize: '0.8rem'
                }}
              >
                <Check size={15} strokeWidth={3} />
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
                  width: '24px',
                  height: '24px',
                  borderRadius: '50%',
                  backgroundColor: '#43281C',
                  color: '#FFFFFF',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 800,
                  fontSize: '0.8rem'
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

            {/* Step 3: Konfirmasi Pesanan */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div
                style={{
                  width: '24px',
                  height: '24px',
                  borderRadius: '50%',
                  backgroundColor: '#FAF5EE',
                  color: '#8E7D73',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 700,
                  fontSize: '0.8rem',
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
            LAYOUT UTAMA 2 KOLOM SESUAI 4 GAMBAR FIGMA PPL-13
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
            {/* KARTU 1: STATUS CHECKOUT (Lanjut tanpa login vs Checkout sebagai pelanggan) */}
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
              <div style={{ display: 'flex', alignItems: 'flex-start', gap: '1rem', marginBottom: '1rem' }}>
                {/* Ikon Centang Hijau Kotak Bulat */}
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
                <span
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    backgroundColor: '#E8F5E9',
                    color: '#2E7D32',
                    fontSize: '0.775rem',
                    fontWeight: 700,
                    padding: '0.35rem 0.95rem',
                    borderRadius: '9999px'
                  }}
                >
                  {!isLoggedIn ? 'Guest checkout aktif' : 'Akun terhubung'}
                </span>

                {!isLoggedIn ? (
                  <span style={{ fontSize: '0.875rem', color: '#8E7D73' }}>
                    Masuk opsional
                  </span>
                ) : (
                  <span style={{ fontSize: '0.875rem', color: '#8E7D73' }}>
                    Satu alur checkout
                  </span>
                )}
              </div>
            </div>

            {/* KARTU 2: INFORMASI KONTAK (FOKUS UTAMA PPL-13) */}
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
                Isi nama penerima dan nomor WhatsApp aktif. Keduanya wajib diisi.
              </p>

              {/* Form Input Fields */}
              <form onSubmit={handleLanjutCheckout}>
                {/* Field 1: Nama Penerima */}
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
                    Nama penerima *
                  </label>
                  <input
                    type="text"
                    value={namaPenerima}
                    onChange={(e) => {
                      setNamaPenerima(e.target.value);
                      setNamaDisentuh(true);
                      setIsAutoFilled(false);
                    }}
                    onBlur={() => setNamaDisentuh(true)}
                    placeholder="Nama penerima"
                    style={{
                      width: '100%',
                      padding: '0.85rem 1.15rem',
                      border: adaErrorNama ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                      borderRadius: '10px',
                      fontSize: '0.95rem',
                      color: '#43281C',
                      outline: 'none',
                      backgroundColor: '#FFFFFF',
                      boxSizing: 'border-box'
                    }}
                  />

                  {/* Helper / Error Text Nama Penerima */}
                  <div style={{ marginTop: '0.4rem', fontSize: '0.8rem' }}>
                    {adaErrorNama ? (
                      <span style={{ color: '#C62828', fontWeight: 600 }}>
                        Nama penerima wajib diisi.
                      </span>
                    ) : isAutoFilled ? (
                      <span style={{ color: '#2E7D32', fontWeight: 600 }}>
                        Diisi otomatis dari akun. Tetap dapat diedit.
                      </span>
                    ) : namaTerisi ? (
                      <span style={{ color: '#2E7D32' }}>
                        Nama penerima sudah diisi.
                      </span>
                    ) : (
                      <span style={{ color: '#8E7D73' }}>
                        Masukkan nama sesuai penerima pesanan.
                      </span>
                    )}
                  </div>
                </div>

                {/* Field 2: Nomor WhatsApp */}
                <div style={{ marginBottom: '1.5rem' }}>
                  <label
                    style={{
                      display: 'block',
                      fontSize: '0.875rem',
                      fontWeight: 700,
                      color: '#43281C',
                      marginBottom: '0.5rem'
                    }}
                  >
                    Nomor WhatsApp *
                  </label>
                  <input
                    type="tel"
                    value={nomorWhatsApp}
                    onChange={(e) => {
                      setNomorWhatsApp(e.target.value);
                      setNomorDisentuh(true);
                    }}
                    onBlur={() => setNomorDisentuh(true)}
                    placeholder="08xxxxxxxxxx"
                    style={{
                      width: '100%',
                      padding: '0.85rem 1.15rem',
                      border: adaErrorNomor ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                      borderRadius: '10px',
                      fontSize: '0.95rem',
                      color: '#43281C',
                      outline: 'none',
                      backgroundColor: '#FFFFFF',
                      boxSizing: 'border-box'
                    }}
                  />

                  {/* Helper / Error Text Nomor WhatsApp */}
                  <div style={{ marginTop: '0.4rem', fontSize: '0.8rem' }}>
                    {adaErrorNomor ? (
                      <span style={{ color: '#C62828', fontWeight: 600 }}>
                        Nomor WhatsApp tidak valid. Gunakan format 08 atau +62.
                      </span>
                    ) : isAutoFilled ? (
                      <span style={{ color: '#8E7D73' }}>
                        Periksa kembali nomor WhatsApp sebelum melanjutkan.
                      </span>
                    ) : nomorValid ? (
                      <span style={{ color: '#2E7D32' }}>
                        Nomor WhatsApp sesuai format nomor seluler Indonesia.
                      </span>
                    ) : (
                      <span style={{ color: '#8E7D73' }}>
                        Gunakan nomor HP Indonesia, diawali 08 atau +62.
                      </span>
                    )}
                  </div>
                </div>

                {/* Banner Status Validasi Kontak Sesuai 4 Gambar Figma */}
                <div style={{ marginBottom: '1rem' }}>
                  {adaErrorForm ? (
                    /* Gambar 3: Banner Error Merah */
                    <div
                      style={{
                        backgroundColor: '#FDF2F2',
                        border: '1px solid #FFCDD2',
                        borderRadius: '8px',
                        padding: '0.75rem 1rem',
                        fontSize: '0.825rem',
                        fontWeight: 700,
                        color: '#C62828'
                      }}
                    >
                      Periksa data kontak sebelum melanjutkan.
                    </div>
                  ) : isAutoFilled ? (
                    /* Gambar 5: Banner Profil Terisi Hijau */
                    <div
                      style={{
                        backgroundColor: '#EDF7ED',
                        border: '1px solid #C8E6C9',
                        borderRadius: '8px',
                        padding: '0.75rem 1rem',
                        fontSize: '0.825rem',
                        fontWeight: 700,
                        color: '#1B5E20'
                      }}
                    >
                      Data profil terisi — pelanggan dapat mengubahnya.
                    </div>
                  ) : formValid ? (
                    /* Gambar 4: Banner Data Valid Hijau */
                    <div
                      style={{
                        backgroundColor: '#EDF7ED',
                        border: '1px solid #C8E6C9',
                        borderRadius: '8px',
                        padding: '0.75rem 1rem',
                        fontSize: '0.825rem',
                        fontWeight: 700,
                        color: '#1B5E20'
                      }}
                    >
                      Data valid — lanjut ke pemilihan metode penerimaan.
                    </div>
                  ) : (
                    /* Gambar 2: Catatan Wajib Diisi */
                    <div style={{ fontSize: '0.8rem', color: '#8E7D73' }}>
                      * Wajib diisi untuk melanjutkan proses checkout.
                    </div>
                  )}
                </div>

                {/* Keterangan Umum Bawah Form */}
                <div style={{ fontSize: '0.775rem', color: '#8E7D73' }}>
                  Data kontak digunakan untuk menghubungi pelanggan terkait pesanan.
                </div>
              </form>
            </div>

            {/* KARTU 3: CALLOUT PENTING */}
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
              KOLOM KANAN: RINGKASAN PESANAN (ORDER SUMMARY)
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

              {/* Rincian Subtotal & Ongkir */}
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

              {/* Subtotal Sementara */}
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

              {/* Tombol Aksi Sesuai Validasi Figma PPL-13 */}
              <button
                type="button"
                disabled={!formValid || loadingCheckout}
                onClick={handleLanjutCheckout}
                style={{
                  width: '100%',
                  backgroundColor: !formValid ? '#BFB5AC' : '#43281C',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '1rem',
                  fontSize: '0.975rem',
                  fontWeight: 700,
                  cursor: !formValid || loadingCheckout ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: !formValid ? 'none' : '0 2px 8px rgba(67, 40, 28, 0.15)'
                }}
                onMouseEnter={(e) => {
                  if (formValid && !loadingCheckout) e.target.style.backgroundColor = '#2F1B12';
                }}
                onMouseLeave={(e) => {
                  if (formValid && !loadingCheckout) e.target.style.backgroundColor = '#43281C';
                }}
              >
                {loadingCheckout
                  ? 'Memproses...'
                  : formValid
                  ? 'Lanjutkan Checkout'
                  : 'Lengkapi Data Kontak'}
              </button>

              {/* Catatan Kaki Bawah Tombol Sesuai Figma */}
              <p
                style={{
                  fontSize: '0.775rem',
                  color: '#8E7D73',
                  textAlign: 'center',
                  marginTop: '0.75rem',
                  marginBottom: 0
                }}
              >
                Berikutnya: pilih metode penerimaan pesanan (PPL-14).
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* ======================================================================
          MODAL KONFIRMASI CHECKOUT
          ====================================================================== */}
      {hasilPesanan && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
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
              maxWidth: '480px',
              width: '100%',
              padding: '2rem',
              boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
              border: '1px solid #EFE8DF'
            }}
          >
            <div style={{ textAlign: 'center', marginBottom: '1.25rem' }}>
              <div
                style={{
                  width: '52px',
                  height: '52px',
                  borderRadius: '50%',
                  backgroundColor: '#DCFCE7',
                  color: '#15803D',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  margin: '0 auto 1rem auto'
                }}
              >
                <Check size={28} strokeWidth={2.5} />
              </div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', margin: 0 }}>
                Data Kontak Tervalidasi!
              </h3>
              <p style={{ fontSize: '0.85rem', color: '#6B5E55', marginTop: '0.35rem', marginBottom: 0 }}>
                {hasilPesanan.pesan}
              </p>
            </div>

            <div
              style={{
                backgroundColor: '#FAF5EE',
                borderRadius: '12px',
                padding: '1rem',
                marginBottom: '1.25rem',
                fontSize: '0.875rem'
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.4rem' }}>
                <span style={{ color: '#8E7D73' }}>Penerima:</span>
                <span style={{ fontWeight: 700, color: '#43281C' }}>
                  {hasilPesanan.kontak?.nama_penerima}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.4rem' }}>
                <span style={{ color: '#8E7D73' }}>WhatsApp:</span>
                <span style={{ fontWeight: 700, color: '#43281C' }}>
                  {hasilPesanan.kontak?.nomor_whatsapp}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.4rem' }}>
                <span style={{ color: '#8E7D73' }}>Status Checkout:</span>
                <span style={{ fontWeight: 700, color: '#2E7D32' }}>
                  {hasilPesanan.tipe_checkout === 'pelanggan' ? 'Pelanggan Login' : 'Guest (Tamu)'}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid #EFE8DF', paddingTop: '0.5rem', marginTop: '0.5rem' }}>
                <span style={{ fontWeight: 800, color: '#43281C' }}>Subtotal:</span>
                <span style={{ fontWeight: 800, color: '#43281C' }}>
                  {formatRupiah(hasilPesanan.total_pembayaran)}
                </span>
              </div>
            </div>

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
              Tutup & Lanjutkan
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
