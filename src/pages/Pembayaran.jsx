import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate, Link } from 'react-router-dom';
import { QRCodeSVG } from 'qrcode.react';
import { 
  QrCode, 
  CreditCard, 
  RefreshCw, 
  ExternalLink, 
  Copy, 
  Check, 
  AlertCircle, 
  Info, 
  ArrowLeft,
  Clock,
  CheckCircle2
} from 'lucide-react';

const API_URL = 'http://localhost:8080';

export default function Pembayaran() {
  const location = useLocation();
  const navigate = useNavigate();

  // Mode pembayaran: 'qris' | 'va' (sesuai Figma PPL19_01 & PPL19_02)
  const [metode, setMetode] = useState('qris');

  // Pilihan bank jika mode Virtual Account
  const [bankVA, setBankVA] = useState('bca');

  // State data pesanan & tagihan dari backend
  const [pesanan, setPesanan] = useState(null);
  const [nomorPesanan, setNomorPesanan] = useState('PJ-DEMO-0001');
  const [totalTagihan, setTotalTagihan] = useState(null);
  const [totalTagihanDisplay, setTotalTagihanDisplay] = useState('mengikuti hasil backend');
  
  // State status pembayaran dari backend
  const [statusPembayaran, setStatusPembayaran] = useState('menunggu_pembayaran');
  const [statusLabel, setStatusLabel] = useState('Menunggu bayar');

  // Status integrasi gateway
  const [isGatewayTersedia, setIsGatewayTersedia] = useState(false);
  const [dataGateway, setDataGateway] = useState(null);

  // Status interaksi UI
  const [loading, setLoading] = useState(true);
  const [isCheckingStatus, setIsCheckingStatus] = useState(false);
  const [notifikasiStatus, setNotifikasiStatus] = useState('');
  const [errorPesanan, setErrorPesanan] = useState('');
  const [copiedVA, setCopiedVA] = useState(false);
  const [tampilkanSpekBE, setTampilkanSpekBE] = useState(false);

  // Ambil parameter pesanan dan sinkronisasi ke backend
  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    const targetNomor = queryParams.get('nomor_pesanan') || 
                        queryParams.get('order_id') || 
                        queryParams.get('id') || 
                        location.state?.pesananData?.nomor_pesanan || 
                        'PJ-DEMO-0001';

    const targetMetode = queryParams.get('metode');
    if (targetMetode === 'va' || targetMetode === 'virtual_account') {
      setMetode('va');
    } else if (targetMetode === 'qris') {
      setMetode('qris');
    }

    setNomorPesanan(targetNomor);

    // Ambil detail pesanan dari backend API
    const loadDetailPesanan = async () => {
      setLoading(true);
      setErrorPesanan('');

      try {
        // Panggil endpoint detail pesanan backend (PPL-18 / WBS 1.2.9.2)
        const res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetNomor)}`);
        
        if (res.ok) {
          const json = await res.json();
          if (json.data) {
            setPesanan(json.data);
            setNomorPesanan(json.data.nomor_pesanan || targetNomor);
            if (json.data.total_pembayaran_format) {
              setTotalTagihan(json.data.total_pembayaran);
              setTotalTagihanDisplay(json.data.total_pembayaran_format);
            } else if (json.data.total_pembayaran) {
              setTotalTagihan(json.data.total_pembayaran);
              setTotalTagihanDisplay(`Rp${Number(json.data.total_pembayaran).toLocaleString('id-ID')}`);
            }

            // Status pembayaran dari backend
            if (json.data.status_pesanan === 'dibayar' || json.data.status_pesanan === 'settlement') {
              setStatusPembayaran('dibayar');
              setStatusLabel('Sudah dibayar');
            } else {
              setStatusPembayaran('menunggu_pembayaran');
              setStatusLabel('Menunggu bayar');
            }
          }
        } else {
          // Jika nomor PJ-DEMO-0001 atau data demo
          if (targetNomor.startsWith('PJ-DEMO')) {
            setTotalTagihanDisplay('mengikuti hasil backend');
            setStatusLabel('Menunggu bayar');
          } else {
            setErrorPesanan('Pesanan tidak ditemukan di database backend.');
          }
        }

        // Minta / sinkronisasi tagihan payment gateway Midtrans ke Backend (WBS 1.3.1.2)
        await mintaTagihanGateway(targetNomor, targetMetode || 'qris', bankVA);
      } catch (err) {
        console.warn('Backend pesanan offline atau belum dapat diakses:', err);
        // Fallback demo sesuai mockup Figma
        if (targetNomor.startsWith('PJ-DEMO')) {
          setTotalTagihanDisplay('mengikuti hasil backend');
          setStatusLabel('Menunggu bayar');
        } else {
          setErrorPesanan('Gagal terhubung ke backend Golang. Pastikan server aktif di port 8080.');
        }
      } finally {
        setLoading(false);
      }
    };

    loadDetailPesanan();
    mintaTagihanGateway(targetNomor, targetMetode || 'qris', bankVA);
  }, [location.search, location.state]);

  // Fungsi membuat / sinkronisasi tagihan payment gateway Midtrans ke Backend (WBS 1.3.1.2 & 1.3.1.4)
  const mintaTagihanGateway = async (orderId, targetMetode = metode, targetBank = bankVA) => {
    try {
      // 1. Coba minta backend buat/ambil tagihan gateway Midtrans Sandbox
      const res = await fetch(`${API_URL}/api/pembayaran/buat-tagihan`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nomor_pesanan: orderId,
          metode: targetMetode,
          bank: targetBank
        })
      });

      if (res.ok) {
        const json = await res.json();
        if (json && json.status === 'sukses' && json.data) {
          setIsGatewayTersedia(true);
          setDataGateway(json.data);
          if (json.data.gross_amount) {
            setTotalTagihan(json.data.gross_amount);
            setTotalTagihanDisplay(`Rp${Number(json.data.gross_amount).toLocaleString('id-ID')}`);
          }
          if (json.data.status_transaksi === 'settlement' || json.data.status_transaksi === 'capture') {
            setStatusPembayaran('dibayar');
            setStatusLabel('Sudah dibayar');
          }
          return;
        }
      }

      // 2. Jika belum, cek status transaksi yang ada
      const resStatus = await fetch(`${API_URL}/api/pembayaran/status?nomor_pesanan=${encodeURIComponent(orderId)}`);
      if (resStatus.ok) {
        const jsonStatus = await resStatus.json();
        if (jsonStatus && jsonStatus.status === 'sukses' && jsonStatus.data) {
          setIsGatewayTersedia(true);
          setDataGateway(jsonStatus.data);
          if (jsonStatus.data.status_transaksi === 'settlement' || jsonStatus.data.status_transaksi === 'capture') {
            setStatusPembayaran('dibayar');
            setStatusLabel('Sudah dibayar');
          }
          return;
        }
      }

      setIsGatewayTersedia(false);
      setDataGateway(null);
    } catch (err) {
      console.warn('Backend payment belum tersedia atau offline:', err);
      setIsGatewayTersedia(false);
      setDataGateway(null);
    }
  };

  // Handler tombol "Ganti metode"
  const handleGantiMetode = () => {
    const metodeBaru = metode === 'qris' ? 'va' : 'qris';
    setMetode(metodeBaru);
    setNotifikasiStatus('');
    mintaTagihanGateway(nomorPesanan, metodeBaru, bankVA);
  };

  // Handler ganti bank Virtual Account
  const handleGantiBankVA = (bankBaru) => {
    setBankVA(bankBaru);
    setNotifikasiStatus('');
    if (metode === 'va') {
      mintaTagihanGateway(nomorPesanan, 'va', bankBaru);
    }
  };

  // Handler simulasi bayar sukses (Midtrans Sandbox - WBS 1.3.1.4)
  const handleSimulasiBayar = async () => {
    setIsCheckingStatus(true);
    setNotifikasiStatus('Mengirim konfirmasi pembayaran ke backend Midtrans...');

    try {
      const res = await fetch(`${API_URL}/api/pembayaran/simulasi-bayar`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nomor_pesanan: nomorPesanan,
          status_transaksi: 'settlement'
        })
      });

      if (res.ok) {
        const json = await res.json();
        setStatusPembayaran('dibayar');
        setStatusLabel('Sudah dibayar');
        setNotifikasiStatus('✅ Pembayaran Berhasil! Database pesanan telah disinkronkan menjadi "dibayar".');
        if (json.data) {
          setDataGateway(json.data);
        }
      } else {
        setNotifikasiStatus('Gagal mengirim simulasi pembayaran.');
      }
    } catch {
      setNotifikasiStatus('Gagal terhubung ke backend simulasi.');
    } finally {
      setIsCheckingStatus(false);
    }
  };

  // Handler tombol "Cek Status Pembayaran" (Mengikuti hasil verifikasi backend/gateway - Ketentuan 9)
  const handleCekStatus = async () => {
    setIsCheckingStatus(true);
    setNotifikasiStatus('');

    try {
      // 1. Cek ke backend pesanan detail (WBS 1.2.9.2 & 1.2.10.2)
      const res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(nomorPesanan)}`);
      if (res.ok) {
        const json = await res.json();
        if (json.data) {
          const statusOrder = json.data.status_pesanan || '';
          if (statusOrder === 'dibayar' || statusOrder === 'settlement') {
            setStatusPembayaran('dibayar');
            setStatusLabel('Sudah dibayar');
            setNotifikasiStatus('Pembayaran berhasil diverifikasi oleh backend!');
            setIsCheckingStatus(false);
            return;
          }
        }
      }

      // 2. Cek ke endpoint status pembayaran Midtrans (WBS 1.3.1.4)
      const resBayar = await fetch(`${API_URL}/api/pembayaran/status?nomor_pesanan=${encodeURIComponent(nomorPesanan)}`);
      if (resBayar.ok) {
        const jsonBayar = await resBayar.json();
        if (jsonBayar.data && (jsonBayar.data.status_transaksi === 'settlement' || jsonBayar.data.status_transaksi === 'capture')) {
          setStatusPembayaran('dibayar');
          setStatusLabel('Sudah dibayar');
          setNotifikasiStatus('Pembayaran terkonfirmasi oleh payment gateway Midtrans!');
          setIsCheckingStatus(false);
          return;
        }
      }

      // Jika belum dibayar menurut backend
      setStatusPembayaran('menunggu_pembayaran');
      setStatusLabel('Menunggu bayar');
      setNotifikasiStatus('Status saat ini: Menunggu konfirmasi pembayaran dari payment gateway.');
    } catch {
      setNotifikasiStatus('Belum menerima verifikasi pembayaran dari backend.');
    } finally {
      setIsCheckingStatus(false);
    }
  };

  // Salin nomor Virtual Account jika ada
  const handleSalinVA = (nomor) => {
    if (!nomor) return;
    navigator.clipboard.writeText(nomor);
    setCopiedVA(true);
    setTimeout(() => setCopiedVA(false), 2000);
  };

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', paddingBottom: '6rem' }}>
      <div
        style={{
          maxWidth: '1120px',
          margin: '0 auto',
          padding: '2.5rem 1.5rem 1rem 1.5rem'
        }}
      >
        {/* ======================================================================
            HEADER HALAMAN SESUAI FIGMA
            Judul dinamis sesuai mode aktif: QRIS atau Virtual Account
            ====================================================================== */}
        <div style={{ marginBottom: '2.25rem' }}>
          <h1
            style={{
              fontSize: '1.85rem',
              fontWeight: 700,
              color: '#2D231E',
              margin: '0 0 0.5rem 0',
              letterSpacing: '-0.2px'
            }}
          >
            {metode === 'qris'
              ? 'Pembayaran QRIS (Sandbox)'
              : 'Pembayaran Virtual Account (Sandbox)'}
          </h1>

          {/* SUBTITLE TEKS MERAH PERSIS FIGMA */}
          <p
            style={{
              fontSize: '0.875rem',
              color: '#B91C1C',
              margin: 0,
              fontWeight: 500
            }}
          >
            Mode sandbox. Nomor/QR pembayaran hanya boleh berasal dari gateway.
          </p>
        </div>

        {/* ======================================================================
            KONDISI LOADING / ERROR
            ====================================================================== */}
        {loading ? (
          <div
            style={{
              backgroundColor: '#FFFFFF',
              borderRadius: '16px',
              border: '1px solid #EFEAE4',
              padding: '4rem 2rem',
              textAlign: 'center',
              boxShadow: '0 2px 8px rgba(67, 40, 28, 0.03)'
            }}
          >
            <div
              style={{
                display: 'inline-block',
                width: '36px',
                height: '36px',
                border: '3px solid #EFEAE4',
                borderTopColor: '#43281C',
                borderRadius: '50%',
                animation: 'spin 0.8s linear infinite',
                marginBottom: '1rem'
              }}
            />
            <style>{`@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }`}</style>
            <p style={{ color: '#6B5E55', fontSize: '0.95rem', margin: 0 }}>
              Memuat data tagihan dan parameter pembayaran sandbox...
            </p>
          </div>
        ) : errorPesanan && !nomorPesanan.startsWith('PJ-DEMO') ? (
          <div
            style={{
              backgroundColor: '#FFFFFF',
              borderRadius: '16px',
              border: '1px solid #FEE2E2',
              padding: '3rem 2rem',
              textAlign: 'center'
            }}
          >
            <AlertCircle size={44} color="#DC2626" style={{ margin: '0 auto 1rem auto' }} />
            <h3 style={{ fontSize: '1.25rem', color: '#991B1B', margin: '0 0 0.5rem 0' }}>
              Gagal Memuat Pesanan
            </h3>
            <p style={{ color: '#6B5E55', fontSize: '0.9rem', marginBottom: '1.5rem' }}>
              {errorPesanan}
            </p>
            <button
              type="button"
              onClick={() => window.location.reload()}
              style={{
                backgroundColor: '#43281C',
                color: '#FFFFFF',
                border: 'none',
                borderRadius: '8px',
                padding: '0.75rem 2rem',
                fontWeight: 600,
                cursor: 'pointer'
              }}
            >
              Coba Muat Ulang
            </button>
          </div>
        ) : (
          /* ======================================================================
              KONTEN UTAMA: LAYOUT 2 KOLOM SESUAI FIGMA
              Kolom Kiri: Metode Pembayaran (QRIS / VA)
              Kolom Kanan: Ringkasan Tagihan
              ====================================================================== */
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'minmax(0, 1.55fr) minmax(0, 1fr)',
              gap: '2rem',
              alignItems: 'start'
            }}
          >
            {/* ------------------------------------------------------------------
                KOLOM KIRI: METODE PEMBAYARAN SESUAI FIGMA
                ------------------------------------------------------------------ */}
            <div
              style={{
                backgroundColor: '#FFFFFF',
                borderRadius: '16px',
                border: '1px solid #ECE7E1',
                padding: '2rem 2.25rem',
                boxShadow: '0 2px 10px rgba(67, 40, 28, 0.03)'
              }}
            >
              <h2
                style={{
                  fontSize: '1.2rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  margin: '0 0 1.25rem 0'
                }}
              >
                Metode pembayaran
              </h2>

              {/* BARIS TOMBOL METODE PEMBAYARAN & GANTI METODE SESUAI FIGMA */}
              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: '1fr 1fr',
                  gap: '1rem',
                  marginBottom: '1.75rem'
                }}
              >
                {/* Tombol Metode Aktif (Cokelat Gelap) */}
                <button
                  type="button"
                  style={{
                    backgroundColor: '#43281C',
                    color: '#FFFFFF',
                    border: 'none',
                    borderRadius: '10px',
                    padding: '0.85rem 1.25rem',
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    cursor: 'default',
                    textAlign: 'center',
                    boxShadow: '0 2px 6px rgba(67, 40, 28, 0.12)'
                  }}
                >
                  {metode === 'qris' ? 'QRIS Sandbox' : 'Virtual Account'}
                </button>

                {/* Tombol Ganti Metode (Outline Sesuai Figma) */}
                <button
                  type="button"
                  onClick={handleGantiMetode}
                  style={{
                    backgroundColor: '#FFFFFF',
                    color: '#43281C',
                    border: '1px solid #DCD5CB',
                    borderRadius: '10px',
                    padding: '0.85rem 1.25rem',
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    cursor: 'pointer',
                    textAlign: 'center',
                    transition: 'all 0.2s ease'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = '#43281C';
                    e.currentTarget.style.backgroundColor = '#FAF5EE';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = '#DCD5CB';
                    e.currentTarget.style.backgroundColor = '#FFFFFF';
                  }}
                >
                  Ganti metode
                </button>
              </div>

              {/* ================================================================
                  AREA BOX UTAMA: KOTAK KREM PERSIS FIGMA
                  - Mode QRIS: "AREA QR PEMBAYARAN" / "PLACEHOLDER • DATA DARI GATEWAY"
                  - Mode VA:   "AREA KODE VIRTUAL ACCOUNT" / "PLACEHOLDER • DATA DARI GATEWAY"
                  ================================================================ */}
              {!isGatewayTersedia || !dataGateway ? (
                /* KONDISI 1: TAMPILAN PERSIS FIGMA PLACEHOLDER KREM */
                <div
                  style={{
                    backgroundColor: '#F5EFEB',
                    borderRadius: '14px',
                    minHeight: '270px',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '2.5rem 1.5rem',
                    textAlign: 'center',
                    marginBottom: '1.5rem',
                    border: '1px dashed #DFD4C8'
                  }}
                >
                  <div
                    style={{
                      fontSize: '1.35rem',
                      fontWeight: 800,
                      color: '#43281C',
                      letterSpacing: '0.5px',
                      marginBottom: '0.55rem'
                    }}
                  >
                    {metode === 'qris'
                      ? 'AREA QR PEMBAYARAN'
                      : 'AREA KODE VIRTUAL ACCOUNT'}
                  </div>

                  <div
                    style={{
                      fontSize: '0.8rem',
                      fontWeight: 700,
                      color: '#B91C1C',
                      letterSpacing: '0.75px'
                    }}
                  >
                    PLACEHOLDER • DATA DARI GATEWAY
                  </div>

                  {/* Keterangan Tambahan Mode VA (Pilihan Bank Simulator) */}
                  {metode === 'va' && (
                    <div
                      style={{
                        marginTop: '1.25rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        fontSize: '0.825rem',
                        color: '#6B5E55'
                      }}
                    >
                      <span>Bank simulasi:</span>
                      <select
                        value={bankVA}
                        onChange={(e) => handleGantiBankVA(e.target.value)}
                        style={{
                          backgroundColor: '#FFFFFF',
                          border: '1px solid #DCD5CB',
                          borderRadius: '6px',
                          padding: '0.25rem 0.5rem',
                          fontSize: '0.825rem',
                          fontWeight: 600,
                          color: '#43281C',
                          cursor: 'pointer'
                        }}
                      >
                        <option value="bca">BCA Virtual Account</option>
                        <option value="bni">BNI Virtual Account</option>
                        <option value="bri">BRI Virtual Account</option>
                        <option value="mandiri">Mandiri Bill Payment</option>
                      </select>
                    </div>
                  )}
                </div>
              ) : (
                /* KONDISI 2: DATA SUDAH BERASAL DARI GATEWAY SANDBOX */
                <div
                  style={{
                    backgroundColor: '#FAF7F2',
                    borderRadius: '14px',
                    padding: '2rem 1.5rem',
                    textAlign: 'center',
                    marginBottom: '1.5rem',
                    border: '1px solid #ECE7E1'
                  }}
                >
                  {metode === 'qris' ? (
                    <div>
                      <div
                        style={{
                          display: 'inline-block',
                          backgroundColor: '#FFFFFF',
                          padding: '1.25rem',
                          borderRadius: '12px',
                          border: '1px solid #EFEAE4',
                          boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
                          marginBottom: '1rem'
                        }}
                      >
                        <QRCodeSVG
                          value={dataGateway.qr_string || dataGateway.snap_redirect_url || `https://app.sandbox.midtrans.com/snap/v4/redirection/${dataGateway.transaction_id}`}
                          size={190}
                          level="M"
                          includeMargin={false}
                        />
                      </div>
                      <div style={{ fontSize: '0.85rem', color: '#6B5E55', fontWeight: 600, marginBottom: '0.75rem' }}>
                        Scan kode QR menggunakan aplikasi e-wallet simulator
                      </div>
                      {dataGateway.snap_token && (
                        <div>
                          <button
                            type="button"
                            onClick={() => {
                              if (window.snap) {
                                window.snap.pay(dataGateway.snap_token, {
                                  onSuccess: (result) => {
                                    console.log('Payment success:', result);
                                    setStatusPembayaran('dibayar');
                                    setStatusLabel('Sudah dibayar');
                                    setNotifikasiStatus('✅ Pembayaran Berhasil! Terverifikasi oleh Midtrans Sandbox.');
                                  },
                                  onPending: (result) => {
                                    console.log('Payment pending:', result);
                                    setNotifikasiStatus('Pembayaran sedang diproses oleh Midtrans Sandbox.');
                                  },
                                  onError: (result) => {
                                    console.log('Payment error:', result);
                                    setNotifikasiStatus('Pembayaran gagal atau dibatalkan di Midtrans.');
                                  }
                                });
                              } else if (dataGateway.snap_redirect_url) {
                                window.open(dataGateway.snap_redirect_url, '_blank');
                              }
                            }}
                            style={{
                              backgroundColor: '#43281C',
                              color: '#FFFFFF',
                              border: 'none',
                              borderRadius: '8px',
                              padding: '0.6rem 1.2rem',
                              fontSize: '0.85rem',
                              fontWeight: 700,
                              cursor: 'pointer',
                              display: 'inline-flex',
                              alignItems: 'center',
                              gap: '0.4rem',
                              boxShadow: '0 2px 6px rgba(67, 40, 28, 0.15)',
                              transition: 'all 0.2s ease'
                            }}
                          >
                            <ExternalLink size={16} /> Buka Layar Pembayaran Midtrans (Popup)
                          </button>
                        </div>
                      )}
                    </div>
                  ) : (
                    <div>
                      <div style={{ fontSize: '0.9rem', color: '#8E7D73', marginBottom: '0.35rem' }}>
                        Nomor Virtual Account ({bankVA.toUpperCase()}):
                      </div>
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          gap: '0.75rem',
                          marginBottom: '1rem'
                        }}
                      >
                        <span
                          style={{
                            fontSize: '1.45rem',
                            fontWeight: 800,
                            letterSpacing: '1px',
                            color: '#43281C'
                          }}
                        >
                          {dataGateway.va_number || `88012${nomorPesanan.replace(/[^0-9]/g, '') || '20261122'}`}
                        </span>
                        <button
                          type="button"
                          onClick={() => handleSalinVA(dataGateway.va_number || `88012${nomorPesanan.replace(/[^0-9]/g, '') || '20261122'}`)}
                          style={{
                            backgroundColor: '#FFFFFF',
                            border: '1px solid #DCD5CB',
                            borderRadius: '8px',
                            padding: '0.4rem 0.75rem',
                            cursor: 'pointer',
                            display: 'inline-flex',
                            alignItems: 'center',
                            gap: '0.35rem',
                            fontSize: '0.8rem',
                            fontWeight: 600,
                            color: '#43281C'
                          }}
                        >
                          {copiedVA ? <Check size={14} color="#15803D" /> : <Copy size={14} />}
                          <span>{copiedVA ? 'Tersalin' : 'Salin'}</span>
                        </button>
                      </div>
                      {dataGateway.snap_token && (
                        <div style={{ marginTop: '0.75rem' }}>
                          <button
                            type="button"
                            onClick={() => {
                              if (window.snap) {
                                window.snap.pay(dataGateway.snap_token, {
                                  onSuccess: (result) => {
                                    console.log('Payment success:', result);
                                    setStatusPembayaran('dibayar');
                                    setStatusLabel('Sudah dibayar');
                                    setNotifikasiStatus('✅ Pembayaran Berhasil! Terverifikasi oleh Midtrans Sandbox.');
                                  },
                                  onPending: (result) => {
                                    console.log('Payment pending:', result);
                                    setNotifikasiStatus('Pembayaran sedang diproses oleh Midtrans Sandbox.');
                                  },
                                  onError: (result) => {
                                    console.log('Payment error:', result);
                                    setNotifikasiStatus('Pembayaran gagal atau dibatalkan di Midtrans.');
                                  }
                                });
                              } else if (dataGateway.snap_redirect_url) {
                                window.open(dataGateway.snap_redirect_url, '_blank');
                              }
                            }}
                            style={{
                              backgroundColor: '#43281C',
                              color: '#FFFFFF',
                              border: 'none',
                              borderRadius: '8px',
                              padding: '0.6rem 1.2rem',
                              fontSize: '0.85rem',
                              fontWeight: 700,
                              cursor: 'pointer',
                              display: 'inline-flex',
                              alignItems: 'center',
                              gap: '0.4rem',
                              boxShadow: '0 2px 6px rgba(67, 40, 28, 0.15)',
                              transition: 'all 0.2s ease'
                            }}
                          >
                            <ExternalLink size={16} /> Buka Layar Pembayaran Midtrans (Popup)
                          </button>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              )}

              {/* CATATAN BAWAH KARTU SESUAI FIGMA */}
              <p
                style={{
                  fontSize: '0.875rem',
                  color: '#8E7D73',
                  margin: 0
                }}
              >
                Ikuti instruksi simulasi dari payment gateway.
              </p>
            </div>

            {/* ------------------------------------------------------------------
                KOLOM KANAN: RINGKASAN TAGIHAN SESUAI FIGMA
                ------------------------------------------------------------------ */}
            <div
              style={{
                backgroundColor: '#FFFFFF',
                borderRadius: '16px',
                border: '1px solid #ECE7E1',
                padding: '2rem 2.25rem',
                boxShadow: '0 2px 10px rgba(67, 40, 28, 0.03)'
              }}
            >
              <h2
                style={{
                  fontSize: '1.2rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  margin: '0 0 1.5rem 0'
                }}
              >
                Ringkasan tagihan
              </h2>

              {/* BARIS PESANAN SESUAI FIGMA */}
              <div
                style={{
                  fontSize: '0.95rem',
                  color: '#6B5E55',
                  marginBottom: '0.85rem'
                }}
              >
                Pesanan:{' '}
                <span style={{ fontWeight: 600, color: '#2D231E' }}>
                  {nomorPesanan}
                </span>
              </div>

              {/* BARIS TOTAL SESUAI FIGMA */}
              <div
                style={{
                  fontSize: '1rem',
                  fontWeight: 700,
                  color: '#2D231E',
                  marginBottom: '1.75rem'
                }}
              >
                Total:{' '}
                <span style={{ color: totalTagihanDisplay.startsWith('Rp') ? '#43281C' : '#2D231E' }}>
                  {totalTagihanDisplay}
                </span>
              </div>

              {/* STATUS BADGE / PILL SESUAI FIGMA */}
              <div style={{ marginBottom: '2rem' }}>
                <span
                  style={{
                    backgroundColor: statusPembayaran === 'dibayar' ? '#DCFCE7' : '#FEF3C7',
                    color: statusPembayaran === 'dibayar' ? '#15803D' : '#B45309',
                    padding: '0.45rem 1.25rem',
                    borderRadius: '9999px',
                    fontSize: '0.85rem',
                    fontWeight: 600,
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: '0.35rem'
                  }}
                >
                  {statusPembayaran === 'dibayar' ? (
                    <CheckCircle2 size={16} />
                  ) : (
                    <Clock size={16} />
                  )}
                  {statusLabel}
                </span>
              </div>

              {/* NOTIFIKASI HASIL CEK STATUS JIKA ADA */}
              {notifikasiStatus && (
                <div
                  style={{
                    backgroundColor: statusPembayaran === 'dibayar' ? '#F0FDF4' : '#FFFBEB',
                    border: `1px solid ${statusPembayaran === 'dibayar' ? '#BBF7D0' : '#FDE68A'}`,
                    color: statusPembayaran === 'dibayar' ? '#166534' : '#92400E',
                    padding: '0.75rem 1rem',
                    borderRadius: '10px',
                    fontSize: '0.825rem',
                    marginBottom: '1.25rem',
                    lineHeight: 1.4
                  }}
                >
                  {notifikasiStatus}
                </div>
              )}

              {/* ================================================================
                  TOMBOL-TOMBOL INTERAKTIF PPL-19 (Ketentuan 3, 8, 9)
                  ================================================================ */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                {/* Tombol Cek Status Pembayaran (Memanggil Verifikasi Backend) */}
                <button
                  type="button"
                  onClick={handleCekStatus}
                  disabled={isCheckingStatus}
                  style={{
                    width: '100%',
                    backgroundColor: '#43281C',
                    color: '#FFFFFF',
                    border: 'none',
                    borderRadius: '10px',
                    padding: '0.9rem 1.25rem',
                    fontSize: '0.925rem',
                    fontWeight: 700,
                    cursor: isCheckingStatus ? 'not-allowed' : 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.5rem',
                    boxShadow: '0 2px 6px rgba(67, 40, 28, 0.15)',
                    transition: 'background-color 0.2s'
                  }}
                  onMouseEnter={(e) => {
                    if (!isCheckingStatus) e.currentTarget.style.backgroundColor = '#2F1B12';
                  }}
                  onMouseLeave={(e) => {
                    if (!isCheckingStatus) e.currentTarget.style.backgroundColor = '#43281C';
                  }}
                >
                  <RefreshCw
                    size={16}
                    style={{
                      animation: isCheckingStatus ? 'spin 0.8s linear infinite' : 'none'
                    }}
                  />
                  <span>
                    {isCheckingStatus ? 'Memeriksa ke Backend...' : 'Cek Status Pembayaran'}
                  </span>
                </button>

                {/* Tombol Buka Simulator Midtrans Sandbox (WBS 1.3.1.3 FE) */}
                <a
                  href={
                    metode === 'qris'
                      ? 'https://simulator.sandbox.midtrans.com/qris/index'
                      : 'https://simulator.sandbox.midtrans.com/openapi/va/index'
                  }
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{
                    width: '100%',
                    backgroundColor: '#FFFFFF',
                    color: '#43281C',
                    border: '1px solid #DCD5CB',
                    borderRadius: '10px',
                    padding: '0.85rem 1.25rem',
                    fontSize: '0.9rem',
                    fontWeight: 600,
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.5rem',
                    textDecoration: 'none',
                    boxSizing: 'border-box',
                    transition: 'all 0.2s'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#FAF5EE';
                    e.currentTarget.style.borderColor = '#43281C';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = '#FFFFFF';
                    e.currentTarget.style.borderColor = '#DCD5CB';
                  }}
                >
                  <ExternalLink size={16} />
                  <span>Buka Simulator Sandbox Midtrans</span>
                </a>

                {/* Tombol Simulasi Pembayaran Sukses (WBS 1.3.1.4 Simulasi) */}
                <button
                  type="button"
                  onClick={handleSimulasiBayar}
                  disabled={isCheckingStatus || statusPembayaran === 'dibayar'}
                  style={{
                    width: '100%',
                    backgroundColor: statusPembayaran === 'dibayar' ? '#15803D' : '#22C55E',
                    color: '#FFFFFF',
                    border: 'none',
                    borderRadius: '10px',
                    padding: '0.85rem 1.25rem',
                    fontSize: '0.875rem',
                    fontWeight: 700,
                    cursor: (isCheckingStatus || statusPembayaran === 'dibayar') ? 'default' : 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.5rem',
                    boxShadow: '0 2px 6px rgba(34, 197, 94, 0.25)',
                    opacity: statusPembayaran === 'dibayar' ? 0.8 : 1,
                    transition: 'all 0.2s'
                  }}
                  onMouseEnter={(e) => {
                    if (statusPembayaran !== 'dibayar') e.currentTarget.style.backgroundColor = '#16A34A';
                  }}
                  onMouseLeave={(e) => {
                    if (statusPembayaran !== 'dibayar') e.currentTarget.style.backgroundColor = '#22C55E';
                  }}
                >
                  <CheckCircle2 size={16} />
                  <span>
                    {statusPembayaran === 'dibayar'
                      ? '✓ Terbayar Lunas di Gateway'
                      : 'Simulasikan Bayar Sukses (Midtrans)'}
                  </span>
                </button>

                {/* Tombol Kembali ke Konfirmasi Pesanan */}
                <Link
                  to={`/konfirmasi-pesanan?nomor_pesanan=${encodeURIComponent(nomorPesanan)}`}
                  style={{
                    textAlign: 'center',
                    fontSize: '0.85rem',
                    color: '#8E7D73',
                    textDecoration: 'none',
                    marginTop: '0.5rem',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.35rem'
                  }}
                >
                  <ArrowLeft size={14} />
                  <span>Kembali ke Detail Pesanan</span>
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
