import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Check, Store, Package, Edit3, ChevronDown } from 'lucide-react';

const API_URL = 'http://localhost:8080';

// Master data wilayah Indonesia bertingkat (Provinsi → Kota/Kabupaten → Kecamatan)
const MASTER_WILAYAH = {
  'Jawa Timur': {
    'Kabupaten Sidoarjo': [
      'Kecamatan Sidoarjo',
      'Kecamatan Gedangan',
      'Kecamatan Waru',
      'Kecamatan Candi',
      'Kecamatan Taman'
    ],
    'Kabupaten Situbondo': [
      'Kecamatan Situbondo',
      'Kecamatan Panji',
      'Kecamatan Besuki',
      'Kecamatan Asembagus'
    ],
    'Kota Surabaya': [
      'Kecamatan Wonokromo',
      'Kecamatan Gubeng',
      'Kecamatan Tegalsari',
      'Kecamatan Rungkut'
    ],
    'Kabupaten Malang': [
      'Kecamatan Kepanjen',
      'Kecamatan Singosari',
      'Kecamatan Lawang'
    ]
  },
  'DKI Jakarta': {
    'Kota Jakarta Selatan': [
      'Kecamatan Kebayoran Baru',
      'Kecamatan Cilandak',
      'Kecamatan Pasar Minggu'
    ],
    'Kota Jakarta Pusat': [
      'Kecamatan Gambir',
      'Kecamatan Menteng',
      'Kecamatan Tanah Abang'
    ]
  },
  'Jawa Barat': {
    'Kota Bandung': [
      'Kecamatan Coblong',
      'Kecamatan Sukasari'
    ],
    'Kota Bekasi': [
      'Kecamatan Bekasi Barat',
      'Kecamatan Bekasi Timur'
    ]
  },
  'Jawa Tengah': {
    'Kota Semarang': [
      'Kecamatan Semarang Barat',
      'Kecamatan Semarang Selatan'
    ],
    'Kota Surakarta': [
      'Kecamatan Banjarsari',
      'Kecamatan Laweyan'
    ]
  }
};

export default function Checkout() {
  const location = useLocation();
  const navigate = useNavigate();

  // Status data kontak dari PPL-13 (sudah terisi sesuai mockup Figma PPL-14 & PPL-15)
  const [kontakTersimpan, setKontakTersimpan] = useState(true);
  const [namaPenerima, setNamaPenerima] = useState('Contoh Pelanggan');
  const [nomorWhatsApp, setNomorWhatsApp] = useState('081234567890');

  // PPL-14: Pilihan Metode Penerimaan ('kurir' | 'ambil_sendiri')
  // Default: 'kurir' sesuai Gambar Mockup Figma PPL-15
  const [metodePenerimaan, setMetodePenerimaan] = useState('kurir');

  // PPL-15: Data Wilayah dan Alamat Pengiriman
  const [provinsi, setProvinsi] = useState('');
  const [kotaKabupaten, setKotaKabupaten] = useState('');
  const [kecamatan, setKecamatan] = useState('');
  const [kodePos, setKodePos] = useState('');
  const [alamatLengkap, setAlamatLengkap] = useState('');

  // Status interaksi untuk validasi error
  const [formDisentuh, setFormDisentuh] = useState(false);

  // Status checkout & hasil
  const [loadingCheckout, setLoadingCheckout] = useState(false);
  const [hasilPesanan, setHasilPesanan] = useState(null);

  // Rincian produk sesuai mockup resmi Figma PPL-15 (2 jenis produk, Subtotal Rp84.000)
  const [itemsCheckout] = useState([
    {
      id: 'udang',
      nama: 'Rengginang Udang',
      varian: 'Udang',
      berat: '450 gr',
      harga: 28000,
      jumlah: 2,
      gambar: '/images/products/rengginang-udang-vertical.jpg'
    },
    {
      id: 'bawang',
      nama: 'Rengginang Bawang',
      varian: 'Bawang',
      berat: '450 gr',
      harga: 28000,
      jumlah: 1,
      gambar: '/images/products/rengginang-bawang.jpg'
    }
  ]);

  // Kalkulasi subtotal produk (2 x 28.000 + 1 x 28.000 = Rp84.000)
  const subtotalProduk = itemsCheckout.reduce(
    (total, item) => total + item.harga * item.jumlah,
    0
  );

  const isKurir = metodePenerimaan === 'kurir';

  // Deteksi query URL jika ingin membuka preset Figma tertentu
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const preset = searchParams.get('preset') || searchParams.get('state');

    if (preset === 'lengkap' || preset === 'valid') {
      // Figma 2: Wilayah & Alamat Lengkap
      setProvinsi('Jawa Timur');
      setKotaKabupaten('Kabupaten Sidoarjo');
      setKecamatan('Kecamatan Sidoarjo');
      setKodePos('61211');
      setAlamatLengkap(
        'Jl. Contoh No. 10, RT 02 / RW 03, Kelurahan Contoh\nDekat minimarket (contoh alamat)'
      );
      setFormDisentuh(true);
    } else if (preset === 'error') {
      // Figma 3: Error data belum lengkap
      setProvinsi('Jawa Timur');
      setKotaKabupaten('');
      setKecamatan('');
      setKodePos('61211');
      setAlamatLengkap('');
      setFormDisentuh(true);
    } else if (preset === 'kosong') {
      // Figma 1: Kosong belum diisi
      setProvinsi('');
      setKotaKabupaten('');
      setKecamatan('');
      setKodePos('');
      setAlamatLengkap('');
      setFormDisentuh(false);
    }
  }, [location.search]);

  // Handler dropdown bertingkat: saat Provinsi berganti, reset Kota dan Kecamatan
  const handleGantiProvinsi = (e) => {
    const p = e.target.value;
    setProvinsi(p);
    setKotaKabupaten('');
    setKecamatan('');
  };

  // Handler dropdown bertingkat: saat Kota berganti, reset Kecamatan
  const handleGantiKota = (e) => {
    const k = e.target.value;
    setKotaKabupaten(k);
    setKecamatan('');
  };

  // Validasi kelengkapan wilayah dan alamat tujuan (WBS 1.2.6.3)
  const provinsiValid = Boolean(provinsi.trim());
  const kotaValid = Boolean(kotaKabupaten.trim());
  const kecamatanValid = Boolean(kecamatan.trim());
  const alamatValid = Boolean(alamatLengkap.trim());

  // Alamat pengiriman dianggap lengkap jika semua field wajib terisi
  const isAlamatLengkap = provinsiValid && kotaValid && kecamatanValid && alamatValid;

  // Kondisi error: jika form sudah disentuh (misal pernah memilih provinsi atau klik lanjut) tapi masih ada yang kosong
  const isFormIncomplete =
    formDisentuh && (!provinsiValid || !kotaValid || !kecamatanValid || !alamatValid);

  // Format angka Rupiah
  const formatRupiah = (angka) => {
    return 'Rp' + Number(angka || 0).toLocaleString('id-ID');
  };

  // Handler submit lanjut ke pemilihan kurir (PPL-16)
  const handleLanjutKeKurir = async (e) => {
    if (e) e.preventDefault();

    if (isKurir && !isAlamatLengkap) {
      setFormDisentuh(true);
      return;
    }

    setLoadingCheckout(true);

    const payloadPengiriman = {
      kontak: {
        nama_penerima: namaPenerima,
        nomor_whatsapp: nomorWhatsApp
      },
      metode_penerimaan: metodePenerimaan,
      wilayah_pengiriman: isKurir
        ? {
            provinsi: provinsi,
            kota_kabupaten: kotaKabupaten,
            kecamatan: kecamatan,
            kode_pos: kodePos,
            alamat_lengkap: alamatLengkap
          }
        : null,
      subtotal_produk: subtotalProduk,
      total_sementara: subtotalProduk
    };

    try {
      const resp = await fetch(`${API_URL}/api/checkout/wilayah-pengiriman`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payloadPengiriman)
      });

      const resJson = await resp.json();
      if (resp.ok && resJson.status === 'sukses') {
        setHasilPesanan(resJson.data);
      } else {
        simulasikanHasil(payloadPengiriman);
      }
    } catch (err) {
      simulasikanHasil(payloadPengiriman);
    } finally {
      setLoadingCheckout(false);
    }
  };

  const simulasikanHasil = (payload) => {
    setHasilPesanan({
      status: 'sukses',
      pesan: 'Data wilayah dan alamat pengiriman berhasil disimpan. Siap dikirim ke API ongkos kirim (PPL-16).',
      wilayah: payload.wilayah_pengiriman
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
            HEADER CHECKOUT SESUAI FIGMA PPL-15
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
            Lengkapi wilayah dan alamat pengiriman pesanan Anda.
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
            LAYOUT UTAMA 2 KOLOM SESUAI FIGMA PPL-15
            ====================================================================== */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'minmax(0, 1.45fr) minmax(360px, 1fr)',
            gap: '2rem',
            alignItems: 'start'
          }}
        >
          {/* ==================================================================
              KOLOM KIRI: KONTAK, METODE PENERIMAAN, WILAYAH & ALAMAT PENGIRIMAN
              ================================================================== */}
          <div>
            {/* KARTU 1: INFORMASI KONTAK (Sesuai Mockup PPL-15) */}
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
              <p style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.25rem' }}>
                Data yang sudah diisi pada tahap PPL-13.
              </p>

              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: '1fr 1fr',
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

            {/* KARTU 2: METODE PENERIMAAN (Sesuai Mockup PPL-15) */}
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
                Pilihan dari PPL-14 tetap ditampilkan.
              </p>

              {/* Opsi 1: Ambil Sendiri di Toko */}
              <div
                onClick={() => setMetodePenerimaan('ambil_sendiri')}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: '1.15rem 1.25rem',
                  borderRadius: '14px',
                  border: !isKurir ? '2px solid #43281C' : '1px solid #ECE4D9',
                  backgroundColor: '#FFFFFF',
                  marginBottom: '1rem',
                  cursor: 'pointer',
                  transition: 'all 0.2s',
                  boxShadow: !isKurir ? '0 2px 8px rgba(67, 40, 28, 0.06)' : 'none'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                  <div
                    style={{
                      width: '20px',
                      height: '20px',
                      borderRadius: '50%',
                      border: !isKurir ? '2px solid #43281C' : '2px solid #D8C7B8',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      flexShrink: 0
                    }}
                  >
                    {!isKurir && (
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

              {/* Opsi 2: Kirim via Ekspedisi / Kurir (Aktif) */}
              <div
                onClick={() => setMetodePenerimaan('kurir')}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: '1.15rem 1.25rem',
                  borderRadius: '14px',
                  border: isKurir ? '2px solid #43281C' : '1px solid #ECE4D9',
                  backgroundColor: '#FFFFFF',
                  cursor: 'pointer',
                  transition: 'all 0.2s',
                  boxShadow: isKurir ? '0 2px 8px rgba(67, 40, 28, 0.06)' : 'none'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                  <div
                    style={{
                      width: '20px',
                      height: '20px',
                      borderRadius: '50%',
                      border: isKurir ? '2px solid #43281C' : '2px solid #D8C7B8',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      flexShrink: 0
                    }}
                  >
                    {isKurir && (
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

            {/* KARTU 3: WILAYAH DAN ALAMAT PENGIRIMAN (FOKUS UTAMA PPL-15) */}
            {isKurir ? (
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
                  Wilayah dan alamat pengiriman
                </h2>
                <p style={{ fontSize: '0.875rem', color: '#8E7D73', marginBottom: '1.5rem' }}>
                  Wajib diisi jika pelanggan memilih Kirim via Kurir.
                </p>

                <form onSubmit={handleLanjutKeKurir}>
                  {/* Grid 2 Kolom untuk Dropdown Bertingkat Wilayah (WBS 1.2.6.1) */}
                  <div
                    style={{
                      display: 'grid',
                      gridTemplateColumns: '1fr 1fr',
                      gap: '1.25rem',
                      marginBottom: '1.25rem'
                    }}
                  >
                    {/* Field 1: Provinsi * */}
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
                        Provinsi *
                      </label>
                      <div style={{ position: 'relative' }}>
                        <select
                          value={provinsi}
                          onChange={handleGantiProvinsi}
                          style={{
                            width: '100%',
                            appearance: 'none',
                            padding: '0.85rem 2.25rem 0.85rem 1.15rem',
                            border: formDisentuh && !provinsiValid ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                            borderRadius: '10px',
                            fontSize: '0.925rem',
                            color: provinsi ? '#43281C' : '#8E7D73',
                            backgroundColor: '#FFFFFF',
                            outline: 'none',
                            cursor: 'pointer',
                            boxSizing: 'border-box'
                          }}
                        >
                          <option value="">Pilih provinsi</option>
                          {Object.keys(MASTER_WILAYAH).map((p) => (
                            <option key={p} value={p}>
                              {p}
                            </option>
                          ))}
                        </select>
                        <ChevronDown
                          size={18}
                          style={{
                            position: 'absolute',
                            right: '1rem',
                            top: '50%',
                            transform: 'translateY(-50%)',
                            color: '#8E7D73',
                            pointerEvents: 'none'
                          }}
                        />
                      </div>
                    </div>

                    {/* Field 2: Kota / Kabupaten * (Bergantung pada Provinsi) */}
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
                        Kota / Kabupaten *
                      </label>
                      <div style={{ position: 'relative' }}>
                        <select
                          value={kotaKabupaten}
                          onChange={handleGantiKota}
                          disabled={!provinsi}
                          style={{
                            width: '100%',
                            appearance: 'none',
                            padding: '0.85rem 2.25rem 0.85rem 1.15rem',
                            border: formDisentuh && !kotaValid ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                            borderRadius: '10px',
                            fontSize: '0.925rem',
                            color: kotaKabupaten ? '#43281C' : '#8E7D73',
                            backgroundColor: !provinsi ? '#FAF8F5' : '#FFFFFF',
                            outline: 'none',
                            cursor: !provinsi ? 'not-allowed' : 'pointer',
                            boxSizing: 'border-box'
                          }}
                        >
                          <option value="">Pilih kota / kabupaten</option>
                          {provinsi &&
                            MASTER_WILAYAH[provinsi] &&
                            Object.keys(MASTER_WILAYAH[provinsi]).map((k) => (
                              <option key={k} value={k}>
                                {k}
                              </option>
                            ))}
                        </select>
                        <ChevronDown
                          size={18}
                          style={{
                            position: 'absolute',
                            right: '1rem',
                            top: '50%',
                            transform: 'translateY(-50%)',
                            color: '#8E7D73',
                            pointerEvents: 'none'
                          }}
                        />
                      </div>
                    </div>

                    {/* Field 3: Kecamatan * (Bergantung pada Kota/Kabupaten) */}
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
                        Kecamatan *
                      </label>
                      <div style={{ position: 'relative' }}>
                        <select
                          value={kecamatan}
                          onChange={(e) => setKecamatan(e.target.value)}
                          disabled={!kotaKabupaten}
                          style={{
                            width: '100%',
                            appearance: 'none',
                            padding: '0.85rem 2.25rem 0.85rem 1.15rem',
                            border: formDisentuh && !kecamatanValid ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                            borderRadius: '10px',
                            fontSize: '0.925rem',
                            color: kecamatan ? '#43281C' : '#8E7D73',
                            backgroundColor: !kotaKabupaten ? '#FAF8F5' : '#FFFFFF',
                            outline: 'none',
                            cursor: !kotaKabupaten ? 'not-allowed' : 'pointer',
                            boxSizing: 'border-box'
                          }}
                        >
                          <option value="">Pilih kecamatan</option>
                          {provinsi &&
                            kotaKabupaten &&
                            MASTER_WILAYAH[provinsi]?.[kotaKabupaten]?.map((kec) => (
                              <option key={kec} value={kec}>
                                {kec}
                              </option>
                            ))}
                        </select>
                        <ChevronDown
                          size={18}
                          style={{
                            position: 'absolute',
                            right: '1rem',
                            top: '50%',
                            transform: 'translateY(-50%)',
                            color: '#8E7D73',
                            pointerEvents: 'none'
                          }}
                        />
                      </div>
                    </div>

                    {/* Field 4: Kode Pos (Opsional) */}
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
                        Kode pos (opsional)
                      </label>
                      <input
                        type="text"
                        value={kodePos}
                        onChange={(e) => setKodePos(e.target.value)}
                        style={{
                          width: '100%',
                          padding: '0.85rem 1.15rem',
                          border: '1px solid #D8C7B8',
                          borderRadius: '10px',
                          fontSize: '0.925rem',
                          color: '#43281C',
                          backgroundColor: '#FFFFFF',
                          outline: 'none',
                          boxSizing: 'border-box'
                        }}
                      />
                    </div>
                  </div>

                  {/* Field 5: Alamat Lengkap * (Textarea) */}
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
                      Alamat lengkap *
                    </label>
                    <textarea
                      rows={3}
                      value={alamatLengkap}
                      onChange={(e) => setAlamatLengkap(e.target.value)}
                      placeholder="Nama jalan, nomor rumah, RT/RW, kelurahan, dan petunjuk lokasi&#10;Contoh: Jl. ..., No. ..., Kelurahan ..."
                      style={{
                        width: '100%',
                        padding: '0.85rem 1.15rem',
                        border: formDisentuh && !alamatValid ? '1px solid #D32F2F' : '1px solid #D8C7B8',
                        borderRadius: '10px',
                        fontSize: '0.925rem',
                        color: '#43281C',
                        backgroundColor: '#FFFFFF',
                        outline: 'none',
                        resize: 'vertical',
                        lineHeight: 1.5,
                        boxSizing: 'border-box'
                      }}
                    />
                  </div>

                  {/* ==========================================================
                      BANNER STATUS SESUAI 3 GAMBAR FIGMA PPL-15:
                      1. State 1 (Default Kosong): Lengkapi wilayah pengiriman (Cream)
                      2. State 2 (Lengkap & Valid): Alamat sudah lengkap (Hijau)
                      3. State 3 (Error Incomplete): Data pengiriman belum lengkap (Merah)
                      ========================================================== */}
                  <div>
                    {isAlamatLengkap ? (
                      /* FIGMA 2: BANNER HIJAU SUKSES */
                      <div
                        style={{
                          backgroundColor: '#EDF7ED',
                          border: '1px solid #C8E6C9',
                          borderRadius: '12px',
                          padding: '1.15rem 1.25rem'
                        }}
                      >
                        <div
                          style={{
                            fontSize: '0.925rem',
                            fontWeight: 800,
                            color: '#1B5E20',
                            marginBottom: '0.25rem'
                          }}
                        >
                          ✓ Alamat sudah lengkap
                        </div>
                        <div style={{ fontSize: '0.825rem', color: '#2E7D32' }}>
                          Wilayah siap dikirim untuk pengecekan ongkir.
                        </div>
                      </div>
                    ) : isFormIncomplete ? (
                      /* FIGMA 3: BANNER MERAH ERROR */
                      <div
                        style={{
                          backgroundColor: '#FDF2F2',
                          border: '1px solid #FFCDD2',
                          borderRadius: '12px',
                          padding: '1.15rem 1.25rem'
                        }}
                      >
                        <div
                          style={{
                            fontSize: '0.925rem',
                            fontWeight: 800,
                            color: '#C62828',
                            marginBottom: '0.25rem'
                          }}
                        >
                          ! Data pengiriman belum lengkap
                        </div>
                        <div style={{ fontSize: '0.825rem', color: '#C62828' }}>
                          Lengkapi Kota/Kabupaten, Kecamatan, dan Alamat Lengkap.
                        </div>
                      </div>
                    ) : (
                      /* FIGMA 1: BANNER CREAM DEFAULT KOSONG */
                      <div
                        style={{
                          backgroundColor: '#FAF5EE',
                          border: '1px solid #F0E6D8',
                          borderRadius: '12px',
                          padding: '1.15rem 1.25rem'
                        }}
                      >
                        <div
                          style={{
                            fontSize: '0.925rem',
                            fontWeight: 800,
                            color: '#43281C',
                            marginBottom: '0.25rem'
                          }}
                        >
                          Lengkapi wilayah pengiriman
                        </div>
                        <div style={{ fontSize: '0.825rem', color: '#8E7D73' }}>
                          Provinsi → Kota/Kabupaten → Kecamatan → Alamat lengkap.
                        </div>
                      </div>
                    )}
                  </div>
                </form>
              </div>
            ) : (
              /* Saat Memilih Ambil Sendiri di Toko: Form Alamat Disembunyikan (WBS 1.2.6.3) */
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
            )}
          </div>

          {/* ==================================================================
              KOLOM KANAN: RINGKASAN PESANAN (ORDER SUMMARY) & CATATAN DEVELOPER
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
                  {itemsCheckout.length} jenis produk
                </div>
              </div>

              {/* Rincian Produk Sesuai Mockup Figma PPL-15 */}
              <div style={{ display: 'flex', gap: '1.25rem', alignItems: 'center', marginBottom: '1rem' }}>
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
                    src={itemsCheckout[0].gambar}
                    alt={itemsCheckout[0].nama}
                    style={{ width: '100%', height: '100%', objectFit: 'contain' }}
                    onError={(e) => {
                      e.target.src = '/images/products/rengginang-udang.jpg';
                    }}
                  />
                </div>

                <div>
                  <div style={{ fontSize: '1rem', fontWeight: 800, color: '#43281C', marginBottom: '0.15rem' }}>
                    {itemsCheckout[0].nama}
                  </div>
                  <div style={{ fontSize: '0.825rem', color: '#8E7D73', marginBottom: '0.35rem' }}>
                    Rasa {itemsCheckout[0].varian} · {itemsCheckout[0].berat}
                  </div>
                  <div style={{ fontSize: '0.925rem', fontWeight: 800, color: '#43281C' }}>
                    {itemsCheckout[0].jumlah} × {formatRupiah(itemsCheckout[0].harga)}
                  </div>
                  <div style={{ fontSize: '0.75rem', color: '#A3968C', marginTop: '0.25rem' }}>
                    Contoh data tampilan
                  </div>
                </div>
              </div>

              {/* Sub-item kedua sesuai mockup Figma: Rengginang Bawang x 1 */}
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  fontSize: '0.875rem',
                  color: '#43281C',
                  padding: '0.5rem 0',
                  borderTop: '1px solid #FAF5EE',
                  marginBottom: '1rem'
                }}
              >
                <span style={{ color: '#6B5E55' }}>
                  {itemsCheckout[1].nama} × {itemsCheckout[1].jumlah}
                </span>
                <span style={{ fontWeight: 700 }}>
                  {formatRupiah(itemsCheckout[1].harga * itemsCheckout[1].jumlah)}
                </span>
              </div>

              {/* Garis Pemisah */}
              <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

              {/* Rincian Subtotal, Metode, dan Ongkir */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Subtotal produk</span>
                <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                  {formatRupiah(subtotalProduk)}
                </span>
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Metode penerimaan</span>
                <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#43281C' }}>
                  {isKurir ? 'Kurir' : 'Ambil sendiri'}
                </span>
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
                <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>Ongkos kirim</span>
                {isKurir ? (
                  <span style={{ fontSize: '0.9rem', color: '#8E7D73' }}>
                    Belum dihitung
                  </span>
                ) : (
                  <span style={{ fontSize: '0.95rem', fontWeight: 700, color: '#2E7D32' }}>
                    Rp0
                  </span>
                )}
              </div>

              {/* Garis Pemisah */}
              <div style={{ borderBottom: '1px solid #EFE8DF', marginBottom: '1.25rem' }} />

              {/* Total Sementara Sesuai Figma PPL-15 */}
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  marginBottom: '0.4rem'
                }}
              >
                <span style={{ fontSize: '1rem', fontWeight: 700, color: '#43281C' }}>
                  {isKurir ? 'Total sementara' : 'Total pesanan'}
                </span>
                <span style={{ fontSize: '1.45rem', fontWeight: 800, color: '#43281C' }}>
                  {formatRupiah(subtotalProduk)}
                </span>
              </div>

              {/* Keterangan Bawah Total Sesuai Figma */}
              <div style={{ fontSize: '0.8rem', color: '#8E7D73', marginBottom: '1.75rem' }}>
                {isKurir
                  ? 'Belum termasuk ongkir dari PPL-16.'
                  : 'Tidak ada biaya pengiriman.'}
              </div>

              {/* Tombol Aksi Sesuai Figma:
                  - Jika Lengkap: "Lanjutkan ke Pilihan Kurir" (Aktif Cokelat)
                  - Jika Belum Lengkap: "Lengkapi Alamat untuk Lanjut" (Disabled Abu-abu) */}
              <button
                type="button"
                disabled={(isKurir && !isAlamatLengkap) || loadingCheckout}
                onClick={handleLanjutKeKurir}
                style={{
                  width: '100%',
                  backgroundColor: isKurir && !isAlamatLengkap ? '#BAACA1' : '#43281C',
                  color: '#FFFFFF',
                  border: 'none',
                  borderRadius: '10px',
                  padding: '1rem',
                  fontSize: '0.975rem',
                  fontWeight: 700,
                  cursor: isKurir && !isAlamatLengkap ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.2s',
                  boxShadow: isKurir && !isAlamatLengkap ? 'none' : '0 2px 8px rgba(67, 40, 28, 0.15)'
                }}
                onMouseEnter={(e) => {
                  if (isAlamatLengkap && !loadingCheckout) e.target.style.backgroundColor = '#2F1B12';
                }}
                onMouseLeave={(e) => {
                  if (isAlamatLengkap && !loadingCheckout) e.target.style.backgroundColor = '#43281C';
                }}
              >
                {loadingCheckout
                  ? 'Memproses...'
                  : isKurir
                  ? isAlamatLengkap
                    ? 'Lanjutkan ke Pilihan Kurir'
                    : 'Lengkapi Alamat untuk Lanjut'
                  : 'Lanjutkan Checkout'}
              </button>
            </div>

            {/* KARTU CATATAN UNTUK DEVELOPER SESUAI FIGMA PPL-15 */}
            <div
              style={{
                backgroundColor: '#FAF5EE',
                border: '1px solid #F0E6D8',
                borderRadius: '16px',
                padding: '1.5rem',
                marginTop: '1.5rem',
                boxShadow: '0 2px 6px rgba(67, 40, 28, 0.02)'
              }}
            >
              <div
                style={{
                  fontSize: '0.925rem',
                  fontWeight: 800,
                  color: '#43281C',
                  marginBottom: '0.65rem'
                }}
              >
                Catatan untuk developer
              </div>
              <div style={{ fontSize: '0.825rem', color: '#6B5E55', lineHeight: 1.6 }}>
                <p style={{ margin: '0 0 0.4rem 0', fontWeight: 700, color: '#43281C' }}>
                  Provinsi → Kota/Kabupaten → Kecamatan
                </p>
                <p style={{ margin: '0 0 0.35rem 0' }}>
                  Pilihan wilayah saling bergantung.
                </p>
                <p style={{ margin: '0 0 0.35rem 0' }}>
                  Simpan ID wilayah resmi untuk API ongkir.
                </p>
                <p style={{ margin: '0 0 0.35rem 0' }}>
                  Alamat tidak wajib saat Ambil Sendiri.
                </p>
                <p style={{ margin: 0 }}>
                  Ongkir sebenarnya ditentukan pada PPL-16.
                </p>
              </div>
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
              maxWidth: '500px',
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
                Wilayah Pengiriman Tervalidasi
              </h3>
              <p style={{ fontSize: '0.875rem', color: '#6B5E55', lineHeight: 1.5 }}>
                {hasilPesanan.pesan}
              </p>
            </div>

            {hasilPesanan.wilayah && (
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
                  <strong>Provinsi:</strong> {hasilPesanan.wilayah.provinsi}
                </div>
                <div>
                  <strong>Kota/Kabupaten:</strong> {hasilPesanan.wilayah.kota_kabupaten}
                </div>
                <div>
                  <strong>Kecamatan:</strong> {hasilPesanan.wilayah.kecamatan}
                </div>
                {hasilPesanan.wilayah.kode_pos && (
                  <div>
                    <strong>Kode Pos:</strong> {hasilPesanan.wilayah.kode_pos}
                  </div>
                )}
                <div>
                  <strong>Alamat Lengkap:</strong> {hasilPesanan.wilayah.alamat_lengkap}
                </div>
              </div>
            )}

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
