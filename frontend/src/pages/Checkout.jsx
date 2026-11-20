import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Check, Store, Package, Edit3 } from 'lucide-react';

const API_URL = 'http://localhost:8080';

export default function Checkout() {
  const location = useLocation();
  const navigate = useNavigate();

  // Status apakah data kontak sudah diisi dan disimpan
  // false = Tampilkan formulir input kontak (PPL-13)
  // true = Tampilkan ringkasan kontak yang sudah diisi (PPL-14)
  const [kontakTersimpan, setKontakTersimpan] = useState(false);

  // Field input data kontak (PPL-13)
  const [namaPenerima, setNamaPenerima] = useState('');
  const [nomorWhatsApp, setNomorWhatsApp] = useState('');

  // Status interaksi validasi
  const [namaDisentuh, setNamaDisentuh] = useState(false);
  const [nomorDisentuh, setNomorDisentuh] = useState(false);

  // PPL-14: Pilihan Metode Penerimaan ('ambil_sendiri' | 'kurir')
  const [metodePenerimaan, setMetodePenerimaan] = useState('ambil_sendiri');

  // Input alamat pengiriman jika memilih kurir
  const [alamatPengiriman, setAlamatPengiriman] = useState('');

  // Status proses checkout
  const [loadingCheckout, setLoadingCheckout] = useState(false);
  const [hasilPesanan, setHasilPesanan] = useState(null);

  // Data ringkasan produk checkout (default sesuai mockup Figma: Rengginang Udang)
  const [itemPesanan, setItemPesanan] = useState({
    nama: 'Rengginang Udang',
    varian: 'Udang',
    berat: '450 gr',
    harga: 28000,
    jumlah: 1,
    gambar: '/images/products/rengginang-udang-vertical.jpg'
  });

  // Validasi nomor WhatsApp format seluler Indonesia (08 / +62 / 62)
  const cekNomorHpValid = (nomor) => {
    if (!nomor) return false;
    const pola = /^(\+62|62|08)[0-9]{8,13}$/;
    return pola.test(nomor.trim().replace(/\s|-/g, ''));
  };

  const namaTerisi = namaPenerima.trim().length > 0;
  const nomorTerisi = nomorWhatsApp.trim().length > 0;
  const nomorValid = cekNomorHpValid(nomorWhatsApp);
  const formKontakValid = namaTerisi && nomorValid;

  const adaErrorNama = namaDisentuh && !namaTerisi;
  const adaErrorNomor = nomorDisentuh && (!nomorTerisi || !nomorValid);
  const adaErrorForm = adaErrorNama || adaErrorNomor;

  // Deteksi query parameter atau state jika dikirim dari katalog / keranjang
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const metodeQuery = searchParams.get('metode') || searchParams.get('metode_penerimaan');

    if (metodeQuery === 'kurir' || metodeQuery === 'ekspedisi') {
      setMetodePenerimaan('kurir');
    } else if (metodeQuery === 'ambil_sendiri' || metodeQuery === 'toko') {
      setMetodePenerimaan('ambil_sendiri');
    }

    // Jika membawa data kontak dari state navigasi
    if (location.state?.kontak) {
      if (location.state.kontak.nama_penerima) {
        setNamaPenerima(location.state.kontak.nama_penerima);
      }
      if (location.state.kontak.nomor_whatsapp) {
        setNomorWhatsApp(location.state.kontak.nomor_whatsapp);
      }
      if (location.state.kontak.nama_penerima && location.state.kontak.nomor_whatsapp) {
        setKontakTersimpan(true);
      }
    }

    // Ambil data item dari location.state jika dikirim dari katalog atau keranjang
    if (location.state?.produk) {
      const p = location.state.produk;
      setItemPesanan({
        nama: p.nama || 'Rengginang Udang',
        varian: p.varian || 'Udang',
        berat: p.beratBersih || '450 gr',
        harga: p.hargaAngka || 28000,
        jumlah: location.state.jumlahBeli || 1,
        gambar: p.gambarDetail || p.gambar || '/images/products/rengginang-udang-vertical.jpg'
      });
    }
  }, [location.search, location.state]);

  // Format angka ke Rupiah
  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  const subtotal = itemPesanan.harga * itemPesanan.jumlah;
  const isAmbilSendiri = metodePenerimaan === 'ambil_sendiri';

  // Handler tombol Simpan Data Kontak (PPL-13 -> PPL-14)
  const handleSimpanKontak = (e) => {
    if (e) e.preventDefault();
    if (!formKontakValid) {
      setNamaDisentuh(true);
      setNomorDisentuh(true);
      return;
    }
    setKontakTersimpan(true);
  };

  // Handler tombol Lanjutkan Checkout (WBS 1.2.5.2)
  const handleLanjutCheckout = async (e) => {
    if (e) e.preventDefault();
    if (!kontakTersimpan) {
      handleSimpanKontak(e);
      return;
    }

    setLoadingCheckout(true);

    const payloadCheckout = {
      kontak: {
        nama_penerima: namaPenerima.trim(),
        nomor_whatsapp: nomorWhatsApp.trim()
      },
      metode_penerimaan: metodePenerimaan,
      ongkos_kirim: isAmbilSendiri ? 0 : null,
      alamat_pengiriman: isAmbilSendiri ? null : alamatPengiriman,
      item_pesanan: [
        {
          nama_produk: itemPesanan.nama,
          varian: itemPesanan.varian,
          harga_satuan: itemPesanan.harga,
          jumlah: itemPesanan.jumlah,
          subtotal: subtotal
        }
      ],
      subtotal: subtotal,
      total_pembayaran: subtotal
    };

    try {
      const resp = await fetch(`${API_URL}/api/checkout/metode-penerimaan`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payloadCheckout)
      });

      const resJson = await resp.json();
      if (resp.ok && resJson.status === 'sukses') {
        setHasilPesanan(resJson.data);
      } else {
        simulasikanHasil(payloadCheckout);
      }
    } catch (err) {
      simulasikanHasil(payloadCheckout);
    } finally {
      setLoadingCheckout(false);
    }
  };

  const simulasikanHasil = (payload) => {
    const nomorPesanan = 'PJ-ORD-' + new Date().toISOString().slice(0, 10).replace(/-/g, '') + '-' + Math.floor(1000 + Math.random() * 9000);
    setHasilPesanan({
      nomor_pesanan: nomorPesanan,
      metode_penerimaan: payload.metode_penerimaan,
      ongkos_kirim: isAmbilSendiri ? 0 : 'Menyusul',
      total_pembayaran: subtotal,
      pesan: isAmbilSendiri
        ? 'Metode Ambil Sendiri di Toko tersimpan. Ongkir Rp0.'
        : 'Metode Kirim via Ekspedisi/Kurir tersimpan. Lanjut ke pengisian wilayah pengiriman (PPL-15).'
    });
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', paddingBottom: '5rem' }}>
      <div style={{ maxWidth: '1180px', margin: '0 auto', padding: '2rem 1.5rem' }}>

        {/* ======================================================================
            BREADCRUMB SESUAI FIGMA (Beranda / Keranjang / Checkout)
            ====================================================================== */}
        <nav style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.25rem' }}>
          <Link to="/" style={{ color: '#8E7D73', textDecoration: 'none' }}>
            Beranda
          </Link>
          <span style={{ margin: '0 0.5rem' }}>/</span>
          <Link to="/keranjang" style={{ color: '#8E7D73', textDecoration: 'none' }}>
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
            {kontakTersimpan
              ? 'Pilih cara menerima pesanan rengginang Anda.'
              : 'Selesaikan pesanan dengan cepat dan mudah.'}
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
            LAYOUT UTAMA 2 KOLOM
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
              KOLOM KIRI: FORM / RINGKASAN KONTAK & METODE PENERIMAAN
              ================================================================== */}
          <div>
            {/* KARTU 1: INFORMASI KONTAK */}
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
              {!kontakTersimpan ? (
                /* -------------------------------------------------------------
                   TAHAP SEBELUM CHECKOUT: ISI FORM INFORMASI KONTAK (PPL-13)
                   ------------------------------------------------------------- */
                <div>
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

                  <form onSubmit={handleSimpanKontak}>
                    {/* Input Nama Penerima */}
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
                      <div style={{ marginTop: '0.4rem', fontSize: '0.8rem' }}>
                        {adaErrorNama ? (
                          <span style={{ color: '#C62828', fontWeight: 600 }}>
                            Nama penerima wajib diisi.
                          </span>
                        ) : (
                          <span style={{ color: '#8E7D73' }}>
                            Masukkan nama sesuai penerima pesanan.
                          </span>
                        )}
                      </div>
                    </div>

                    {/* Input Nomor WhatsApp */}
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
                      <div style={{ marginTop: '0.4rem', fontSize: '0.8rem' }}>
                        {adaErrorNomor ? (
                          <span style={{ color: '#C62828', fontWeight: 600 }}>
                            Nomor WhatsApp tidak valid. Gunakan format 08 atau +62.
                          </span>
                        ) : (
                          <span style={{ color: '#8E7D73' }}>
                            Gunakan nomor HP Indonesia, diawali 08 atau +62.
                          </span>
                        )}
                      </div>
                    </div>

                    {/* Catatan Wajib Diisi */}
                    <div style={{ fontSize: '0.8rem', color: '#8E7D73', marginBottom: '0.4rem' }}>
                      * Wajib diisi untuk melanjutkan proses checkout.
                    </div>
                    <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '1.5rem' }}>
                      Data kontak digunakan untuk menghubungi pelanggan terkait pesanan.
                    </div>

                    {/* Banner Peringatan jika Error */}
                    {adaErrorForm && (
                      <div
                        style={{
                          backgroundColor: '#FDF2F2',
                          border: '1px solid #FFCDD2',
                          borderRadius: '8px',
                          padding: '0.75rem 1rem',
                          fontSize: '0.825rem',
                          fontWeight: 700,
                          color: '#C62828',
                          marginBottom: '1.25rem'
                        }}
                      >
                        Periksa data kontak sebelum melanjutkan.
                      </div>
                    )}

                    {/* Tombol Simpan Kontak */}
                    <button
                      type="submit"
                      disabled={!formKontakValid}
                      style={{
                        backgroundColor: formKontakValid ? '#43281C' : '#BFB5AC',
                        color: '#FFFFFF',
                        border: 'none',
                        borderRadius: '10px',
                        padding: '0.85rem 1.5rem',
                        fontSize: '0.925rem',
                        fontWeight: 700,
                        cursor: formKontakValid ? 'pointer' : 'not-allowed',
                        transition: 'background-color 0.2s',
                        boxShadow: formKontakValid ? '0 2px 8px rgba(67, 40, 28, 0.15)' : 'none'
                      }}
                    >
                      Simpan Data Kontak
                    </button>
                  </form>
                </div>
              ) : (
                /* -------------------------------------------------------------
                   SETELAH ISI: MUNCUL INFORMASI KONTAK YANG SUDAH DIISI (PPL-14)
                   ------------------------------------------------------------- */
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.25rem' }}>
                    <div>
                      <h2
                        style={{
                          fontSize: '1.15rem',
                          fontWeight: 800,
                          color: '#43281C',
                          marginBottom: '0.2rem'
                        }}
                      >
                        Informasi kontak
                      </h2>
                      <p style={{ fontSize: '0.875rem', color: '#8E7D73', margin: 0 }}>
                        Data yang sudah diisi pada tahap PPL-13.
                      </p>
                    </div>

                    <button
                      type="button"
                      onClick={() => setKontakTersimpan(false)}
                      style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.35rem',
                        backgroundColor: '#FAF5EE',
                        border: '1px solid #D8C7B8',
                        color: '#43281C',
                        fontSize: '0.8rem',
                        fontWeight: 700,
                        padding: '0.4rem 0.85rem',
                        borderRadius: '8px',
                        cursor: 'pointer'
                      }}
                    >
                      <Edit3 size={14} />
                      Ubah data kontak
                    </button>
                  </div>

                  {/* Box Rincian Kontak yang Sudah Diisi Sesuai Mockup */}
                  <div
                    style={{
                      display: 'grid',
                      gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
                      gap: '1rem',
                      backgroundColor: '#FFFFFF',
                      border: '1px solid #EFE8DF',
                      borderRadius: '12px',
                      padding: '1.15rem 1.35rem'
                    }}
                  >
                    <div>
                      <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.25rem' }}>
                        Nama penerima
                      </div>
                      <div style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                        {namaPenerima}
                      </div>
                    </div>

                    <div>
                      <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginBottom: '0.25rem' }}>
                        Nomor WhatsApp
                      </div>
                      <div style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                        {nomorWhatsApp}
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* KARTU 2: METODE PENERIMAAN (HANYA AKTIF / MUNCUL SETELAH KONTAK DIISI SESUAI FIGMA PPL-14) */}
            {kontakTersimpan && (
              <>
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
                    Metode penerimaan
                  </h2>
                  <p style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.5rem' }}>
                    Pilih satu metode sesuai kebutuhan Anda.
                  </p>

                  {/* OPSI 1: Ambil Sendiri di Toko */}
                  <div
                    onClick={() => setMetodePenerimaan('ambil_sendiri')}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      padding: '1.15rem 1.25rem',
                      borderRadius: '14px',
                      border: isAmbilSendiri ? '2px solid #43281C' : '1px solid #ECE4D9',
                      backgroundColor: '#FFFFFF',
                      marginBottom: '1rem',
                      cursor: 'pointer',
                      transition: 'all 0.2s',
                      boxShadow: isAmbilSendiri ? '0 2px 8px rgba(67, 40, 28, 0.06)' : 'none'
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                      {/* Radio Button Bulat */}
                      <div
                        style={{
                          width: '20px',
                          height: '20px',
                          borderRadius: '50%',
                          border: isAmbilSendiri ? '2px solid #43281C' : '2px solid #D8C7B8',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          flexShrink: 0
                        }}
                      >
                        {isAmbilSendiri && (
                          <div
                            style={{
                              width: '10px',
                              height: '10px',
                              borderRadius: '50%',
                              backgroundColor: '#43281C'
                            }}
                          />
                        )}
                      </div>

                      {/* Icon Toko */}
                      <div
                        style={{
                          width: '42px',
                          height: '42px',
                          borderRadius: '10px',
                          backgroundColor: '#FAF5EE',
                          border: '1px solid #EFE8DF',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          color: '#43281C',
                          flexShrink: 0
                        }}
                      >
                        <Store size={22} strokeWidth={1.8} />
                      </div>

                      {/* Teks Opsi */}
                      <div>
                        <div
                          style={{
                            fontSize: '0.975rem',
                            fontWeight: 800,
                            color: '#43281C',
                            marginBottom: '0.2rem'
                          }}
                        >
                          Ambil Sendiri di Toko
                        </div>
                        <div style={{ fontSize: '0.825rem', color: '#8E7D73' }}>
                          Pelanggan mengambil pesanan secara langsung.
                        </div>
                      </div>
                    </div>

                    {/* Badge Ongkir Rp0 */}
                    <div
                      style={{
                        backgroundColor: '#E8F5E9',
                        color: '#2E7D32',
                        fontSize: '0.775rem',
                        fontWeight: 700,
                        padding: '0.35rem 0.85rem',
                        borderRadius: '9999px',
                        whiteSpace: 'nowrap'
                      }}
                    >
                      Ongkir Rp0
                    </div>
                  </div>

                  {/* OPSI 2: Kirim via Ekspedisi / Kurir */}
                  <div
                    onClick={() => setMetodePenerimaan('kurir')}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      padding: '1.15rem 1.25rem',
                      borderRadius: '14px',
                      border: !isAmbilSendiri ? '2px solid #43281C' : '1px solid #ECE4D9',
                      backgroundColor: '#FFFFFF',
                      cursor: 'pointer',
                      transition: 'all 0.2s',
                      boxShadow: !isAmbilSendiri ? '0 2px 8px rgba(67, 40, 28, 0.06)' : 'none'
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                      {/* Radio Button Bulat */}
                      <div
                        style={{
                          width: '20px',
                          height: '20px',
                          borderRadius: '50%',
                          border: !isAmbilSendiri ? '2px solid #43281C' : '2px solid #D8C7B8',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          flexShrink: 0
                        }}
                      >
                        {!isAmbilSendiri && (
                          <div
                            style={{
                              width: '10px',
                              height: '10px',
                              borderRadius: '50%',
                              backgroundColor: '#43281C'
                            }}
                          />
                        )}
                      </div>

                      {/* Icon Kurir / Box */}
                      <div
                        style={{
                          width: '42px',
                          height: '42px',
                          borderRadius: '10px',
                          backgroundColor: '#FAF5EE',
                          border: '1px solid #EFE8DF',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          color: '#43281C',
                          flexShrink: 0
                        }}
                      >
                        <Package size={22} strokeWidth={1.8} />
                      </div>

                      {/* Teks Opsi */}
                      <div>
                        <div
                          style={{
                            fontSize: '0.975rem',
                            fontWeight: 800,
                            color: '#43281C',
                            marginBottom: '0.2rem'
                          }}
                        >
                          Kirim via Ekspedisi / Kurir
                        </div>
                        <div style={{ fontSize: '0.825rem', color: '#8E7D73' }}>
                          Pesanan dikirim ke alamat tujuan.
                        </div>
                      </div>
                    </div>

                    {/* Badge Ongkir menyusul */}
                    <div
                      style={{
                        backgroundColor: '#FAF5EE',
                        color: '#8E7D73',
                        border: '1px solid #EFE8DF',
                        fontSize: '0.775rem',
                        fontWeight: 700,
                        padding: '0.35rem 0.85rem',
                        borderRadius: '9999px',
                        whiteSpace: 'nowrap'
                      }}
                    >
                      Ongkir menyusul
                    </div>
                  </div>
                </div>

                {/* KARTU 3: KONDISIONAL BERDASARKAN METODE PENERIMAAN (WBS 1.2.5.3) */}
                {isAmbilSendiri ? (
                  /* FIGMA 1: Callout Ambil Sendiri Dipilih (Form alamat kurir disembunyikan) */
                  <div
                    style={{
                      backgroundColor: '#EDF7ED',
                      border: '1px solid #C8E6C9',
                      borderRadius: '14px',
                      padding: '1.25rem 1.5rem',
                      boxShadow: '0 2px 6px rgba(46, 125, 50, 0.04)'
                    }}
                  >
                    <div
                      style={{
                        fontSize: '0.925rem',
                        fontWeight: 800,
                        color: '#1B5E20',
                        marginBottom: '0.35rem'
                      }}
                    >
                      ✓ Ambil Sendiri dipilih
                    </div>
                    <div style={{ fontSize: '0.825rem', color: '#2E7D32', lineHeight: 1.55 }}>
                      Alamat pengiriman tidak diperlukan dan form alamat disembunyikan.
                    </div>
                    <div style={{ fontSize: '0.825rem', color: '#2E7D32', lineHeight: 1.55, marginTop: '0.2rem' }}>
                      Biaya pengiriman otomatis Rp0.
                    </div>
                  </div>
                ) : (
                  /* FIGMA 2: Alamat Pengiriman (Tampil hanya saat memilih Kirim via Ekspedisi/Kurir) */
                  <div
                    style={{
                      backgroundColor: '#FFFFFF',
                      border: '1px solid #ECE4D9',
                      borderRadius: '16px',
                      padding: '1.75rem 2rem',
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
                      Alamat pengiriman
                    </h2>
                    <p style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.25rem' }}>
                      Tampil hanya saat memilih Kirim via Ekspedisi/Kurir.
                    </p>

                    <div>
                      <label
                        style={{
                          display: 'block',
                          fontSize: '0.875rem',
                          fontWeight: 700,
                          color: '#43281C',
                          marginBottom: '0.5rem'
                        }}
                      >
                        Alamat lengkap *
                      </label>
                      <input
                        type="text"
                        value={alamatPengiriman}
                        onChange={(e) => setAlamatPengiriman(e.target.value)}
                        placeholder="Masukkan alamat penerima untuk pengiriman"
                        style={{
                          width: '100%',
                          padding: '0.85rem 1.15rem',
                          border: '1px solid #D8C7B8',
                          borderRadius: '10px',
                          fontSize: '0.95rem',
                          color: '#43281C',
                          outline: 'none',
                          backgroundColor: '#FFFFFF',
                          boxSizing: 'border-box'
                        }}
                      />
                      <div style={{ fontSize: '0.775rem', color: '#8E7D73', marginTop: '0.5rem' }}>
                        Detail alamat dan validasinya merupakan cakupan PPL-15.
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
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
                    style={{ width: '100%', height: '100%', objectFit: 'contain' }}
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

              {/* Rincian Subtotal, Metode, dan Ongkir Sesuai Figma */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Subtotal produk</span>
                <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                  {formatRupiah(subtotal)}
                </span>
              </div>

              {kontakTersimpan ? (
                <>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                    <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Metode penerimaan</span>
                    <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                      {isAmbilSendiri ? 'Ambil sendiri' : 'Kurir'}
                    </span>
                  </div>

                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
                    <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Ongkos kirim</span>
                    {isAmbilSendiri ? (
                      <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#2E7D32' }}>
                        Rp0
                      </span>
                    ) : (
                      <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>
                        Belum dihitung
                      </span>
                    )}
                  </div>
                </>
              ) : (
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Ongkos kirim</span>
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Dihitung nanti</span>
                </div>
              )}

              {/* Garis Pemisah */}
              <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

              {/* Total Pesanan / Total Sementara Sesuai Figma */}
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  marginBottom: '0.4rem'
                }}
              >
                <span style={{ fontSize: '1rem', fontWeight: 700, color: '#43281C' }}>
                  {kontakTersimpan && isAmbilSendiri ? 'Total pesanan' : 'Total sementara'}
                </span>
                <span style={{ fontSize: '1.45rem', fontWeight: 800, color: '#43281C' }}>
                  {formatRupiah(subtotal)}
                </span>
              </div>

              {/* Keterangan Bawah Total */}
              <div style={{ fontSize: '0.8rem', color: '#8E7D73', marginBottom: '1.75rem' }}>
                {kontakTersimpan
                  ? isAmbilSendiri
                    ? 'Tidak ada biaya pengiriman.'
                    : 'Belum termasuk ongkir dari PPL-16.'
                  : 'Selesaikan data kontak untuk melanjutkan.'}
              </div>

              {/* Tombol Aksi */}
              <button
                type="button"
                disabled={!kontakTersimpan || loadingCheckout}
                onClick={handleLanjutCheckout}
                style={{
                  width: '100%',
                  backgroundColor: !kontakTersimpan ? '#BFB5AC' : '#43281C',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '1rem',
                  fontSize: '0.975rem',
                  fontWeight: 700,
                  cursor: !kontakTersimpan || loadingCheckout ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: !kontakTersimpan ? 'none' : '0 2px 8px rgba(67, 40, 28, 0.15)'
                }}
                onMouseEnter={(e) => {
                  if (kontakTersimpan && !loadingCheckout) e.target.style.backgroundColor = '#2F1B12';
                }}
                onMouseLeave={(e) => {
                  if (kontakTersimpan && !loadingCheckout) e.target.style.backgroundColor = '#43281C';
                }}
              >
                {loadingCheckout
                  ? 'Memproses...'
                  : kontakTersimpan
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
                {kontakTersimpan
                  ? 'Pilihan disimpan pada alur pesanan yang sama.'
                  : 'Berikutnya: pilih metode penerimaan pesanan (PPL-14).'}
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
                  backgroundColor: '#E8F5E9',
                  color: '#2E7D32',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  margin: '0 auto 1rem',
                  fontWeight: 800,
                  fontSize: '1.5rem'
                }}
              >
                ✓
              </div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: 800, color: '#43281C', marginBottom: '0.5rem' }}>
                Metode Penerimaan Tersimpan
              </h3>
              <p style={{ fontSize: '0.875rem', color: '#6B5E55', lineHeight: 1.5 }}>
                {hasilPesanan.pesan}
              </p>
            </div>

            <div
              style={{
                backgroundColor: '#FAF5EE',
                borderRadius: '12px',
                padding: '1rem 1.25rem',
                marginBottom: '1.5rem',
                fontSize: '0.85rem',
                lineHeight: 1.6,
                color: '#43281C'
              }}
            >
              <div>
                <strong>Metode:</strong> {hasilPesanan.metode_penerimaan === 'ambil_sendiri' ? 'Ambil Sendiri di Toko' : 'Kirim via Ekspedisi/Kurir'}
              </div>
              <div>
                <strong>Ongkos Kirim:</strong> {hasilPesanan.ongkos_kirim === 0 ? 'Rp0 (Gratis)' : hasilPesanan.ongkos_kirim}
              </div>
              <div>
                <strong>Total Pembayaran:</strong> {formatRupiah(hasilPesanan.total_pembayaran)}
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
                fontSize: '0.925rem',
                fontWeight: 700,
                cursor: 'pointer'
              }}
            >
              Tutup
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
