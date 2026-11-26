import React, { useEffect, useState } from 'react';
import { useParams, useLocation, useNavigate, Link } from 'react-router-dom';
import { Printer, ArrowLeft, CheckCircle, Clock, XCircle, FileText } from 'lucide-react';

const API_URL = 'http://localhost:8080';

// Data fallback demo sesuai mockup Figma PPL23_01_Invoice_Digital__HALAMAN_INVOICE
const DEMO_INVOICE_LOOKUP = {
  'PJ-DEMO-0002': {
    nomor_pesanan: 'PJ-DEMO-0002',
    is_demo: true,
    pembeli: 'Pelanggan Contoh',
    status_pembayaran: 'Dibayar',
    status_pemrosesan: 'Diproses',
    items: [
      { nama_produk: 'Rengginang Udang', jumlah_bungkus: 2, subtotal_format: 'Rp56.000' },
      { nama_produk: 'Rengginang Bawang', jumlah_bungkus: 1, subtotal_format: 'Rp28.000' }
    ],
    subtotal_ongkir: 'Dari data order',
    total_pembayaran: 'Dari backend'
  },
  'PJ-DEMO-0001': {
    nomor_pesanan: 'PJ-DEMO-0001',
    is_demo: true,
    pembeli: 'Pelanggan A',
    status_pembayaran: 'Menunggu',
    status_pemrosesan: 'Baru',
    items: [
      { nama_produk: 'Rengginang Udang', jumlah_bungkus: 2, subtotal_format: 'Rp56.000' }
    ],
    subtotal_ongkir: 'Rp56.000 (Ongkir Rp0)',
    total_pembayaran: 'Rp56.000'
  }
};

export default function InvoiceDigital() {
  const { id } = useParams();
  const location = useLocation();
  const navigate = useNavigate();

  // Query parameter support
  const queryParams = new URLSearchParams(location.search);
  const targetId = id || queryParams.get('id') || queryParams.get('nomor_pesanan') || 'PJ-DEMO-0002';

  const [invoice, setInvoice] = useState(null);
  const [loading, setLoading] = useState(true);

  // Ambil data invoice dari backend API
  useEffect(() => {
    let isMounted = true;
    setLoading(true);

    const fetchInvoiceData = async () => {
      try {
        // Coba endpoint publik detail pesanan
        let res = await fetch(`${API_URL}/api/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetId)}`);
        
        // Jika gagal, coba endpoint admin
        if (!res.ok) {
          res = await fetch(`${API_URL}/api/admin/pesanan/detail?nomor_pesanan=${encodeURIComponent(targetId)}`, {
            credentials: 'include'
          });
        }

        if (res.ok) {
          const json = await res.json();
          const d = json.data || {};
          const detail = d.pesanan_detail || {};

          const formatRibuan = (num) => 'Rp' + Number(num || 0).toLocaleString('id-ID');

          // Parse items
          let itemRows = [];
          if (Array.isArray(detail.items) && detail.items.length > 0) {
            itemRows = detail.items.map((it) => ({
              nama_produk: it.nama_varian || it.nama_produk || 'Rengginang Putra Jaya',
              jumlah_bungkus: it.jumlah_bungkus || it.jumlah || 1,
              subtotal_format: it.subtotal ? formatRibuan(it.subtotal) : '-'
            }));
          } else if (d.items) {
            // String fallback
            itemRows = [
              { nama_produk: d.items, jumlah_bungkus: d.total_bungkus || 1, subtotal_format: d.total_pembayaran_format || '-' }
            ];
          }

          // Hitung subtotal & ongkir format
          const subtotalOngkirText = detail.snapshot_subtotal
            ? `${formatRibuan(detail.snapshot_subtotal)} (Ongkir: ${formatRibuan(detail.snapshot_ongkir || 0)})`
            : d.ongkir
            ? `Subtotal: ${formatRibuan(d.subtotal || 0)} + Ongkir: ${formatRibuan(d.ongkir)}`
            : 'Dari data order';

          const totalText = d.total_pembayaran_format
            || (detail.snapshot_total ? formatRibuan(detail.snapshot_total) : (d.total ? formatRibuan(d.total) : 'Dari backend'));

          if (isMounted) {
            setInvoice({
              nomor_pesanan: d.nomor_pesanan || targetId,
              is_demo: (d.nomor_pesanan || targetId).startsWith('PJ-DEMO'),
              pembeli: d.penerima || detail.kontak?.nama_penerima || 'Pelanggan',
              status_pembayaran: d.status_pembayaran || (d.status_pesanan === 'dibayar' ? 'Dibayar' : 'Menunggu Pembayaran'),
              status_pemrosesan: d.status_pemrosesan || d.status || 'Diproses',
              items: itemRows.length > 0 ? itemRows : DEMO_INVOICE_LOOKUP['PJ-DEMO-0002'].items,
              subtotal_ongkir: subtotalOngkirText,
              total_pembayaran: totalText
            });
            setLoading(false);
            return;
          }
        }
      } catch (err) {
        console.warn('Gagal memuat invoice dari backend, menggunakan data demo:', err);
      }

      // Fallback ke data demo Figma jika pesanan demo atau backend offline
      if (isMounted) {
        const fallback = DEMO_INVOICE_LOOKUP[targetId] || {
          ...DEMO_INVOICE_LOOKUP['PJ-DEMO-0002'],
          nomor_pesanan: targetId
        };
        setInvoice(fallback);
        setLoading(false);
      }
    };

    fetchInvoiceData();

    return () => {
      isMounted = false;
    };
  }, [targetId]);

  const handlePrint = () => {
    window.print();
  };

  const inv = invoice || DEMO_INVOICE_LOOKUP['PJ-DEMO-0002'];

  return (
    <div style={{ backgroundColor: '#FAF8F5', minHeight: '90vh', padding: '2.5rem 1rem 5rem 1rem' }}>
      <div style={{ maxWidth: '840px', margin: '0 auto' }}>
        
        {/* Tombol Aksi Atas (Kembali & Cetak) */}
        <div
          className="no-print"
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: '1.5rem'
          }}
        >
          <button
            type="button"
            onClick={() => navigate(-1)}
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: '0.4rem',
              backgroundColor: 'transparent',
              border: 'none',
              color: '#6B5E55',
              fontSize: '0.9rem',
              fontWeight: 600,
              cursor: 'pointer'
            }}
          >
            <ArrowLeft size={16} />
            <span>Kembali</span>
          </button>

          <button
            type="button"
            onClick={handlePrint}
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: '0.5rem',
              backgroundColor: '#FFFFFF',
              border: '1px solid #D8C7B8',
              color: '#43281C',
              padding: '0.5rem 1.15rem',
              borderRadius: '8px',
              fontSize: '0.875rem',
              fontWeight: 600,
              cursor: 'pointer',
              boxShadow: '0 2px 4px rgba(0,0,0,0.03)'
            }}
          >
            <Printer size={16} />
            <span>Cetak / Unduh PDF</span>
          </button>
        </div>

        {/* JUDUL HALAMAN SESUAI FIGMA FRAME PPL23_01_Invoice_Digital__HALAMAN_INVOICE */}
        <h1
          style={{
            fontSize: '1.75rem',
            fontWeight: 800,
            color: '#2D231E',
            margin: '0 0 2rem 0'
          }}
        >
          Invoice Digital
        </h1>

        {/* KARTU INVOICE SESUAI DESAIN FIGMA */}
        <div
          style={{
            backgroundColor: '#FFFFFF',
            borderRadius: '16px',
            border: '1px solid #ECE7E1',
            boxShadow: '0 4px 24px rgba(67, 40, 28, 0.04)',
            padding: '3rem 2.5rem',
            position: 'relative'
          }}
        >
          {/* HEADER KARTU INVOICE */}
          <div style={{ marginBottom: '2rem' }}>
            <h2
              style={{
                fontSize: '1.35rem',
                fontWeight: 800,
                color: '#2D231E',
                margin: '0 0 0.5rem 0',
                letterSpacing: '0.2px'
              }}
            >
              RENGGINANG PUTRA JAYA — INVOICE
            </h2>
            <p
              style={{
                fontSize: '0.9rem',
                color: '#7A6E65',
                margin: 0
              }}
            >
              Pesanan: <strong>{inv.nomor_pesanan}</strong> |{' '}
              {inv.is_demo ? 'Contoh data layar' : 'Data transaksi backend'}
            </p>
          </div>

          {/* TABEL RINCIAN PRODUK & JUMLAH */}
          <div style={{ marginTop: '1.5rem', marginBottom: '2rem' }}>
            {/* Header Kolom */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                paddingBottom: '0.75rem',
                borderBottom: '1px solid #EFE8DF',
                fontSize: '0.95rem',
                fontWeight: 700,
                color: '#2D231E'
              }}
            >
              <span>Produk</span>
              <span>Jumlah</span>
            </div>

            {/* Baris Daftar Produk */}
            {loading ? (
              <div style={{ padding: '2rem 0', textAlign: 'center', color: '#8E7D73' }}>
                Memuat rincian invoice...
              </div>
            ) : (
              inv.items.map((item, idx) => (
                <div
                  key={idx}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '1rem 0',
                    borderBottom: '1px solid #F6F1EA',
                    fontSize: '0.95rem',
                    color: '#3B302A'
                  }}
                >
                  <span style={{ fontWeight: 500 }}>{item.nama_produk}</span>
                  <span style={{ fontWeight: 700, color: '#2D231E' }}>
                    {item.jumlah_bungkus} bungkus
                  </span>
                </div>
              ))
            )}

            {/* Baris Subtotal & Ongkir */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '1.25rem 0 0.85rem 0',
                fontSize: '0.95rem',
                color: '#554A43'
              }}
            >
              <span>Subtotal & ongkir</span>
              <span style={{ fontWeight: 700, color: '#2D231E' }}>
                {inv.subtotal_ongkir}
              </span>
            </div>

            {/* Divider Ringan */}
            <div style={{ height: '1px', backgroundColor: '#EFE8DF', margin: '0.5rem 0' }} />

            {/* Baris Total */}
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '1rem 0',
                fontSize: '1.05rem',
                color: '#2D231E'
              }}
            >
              <span style={{ fontWeight: 700 }}>Total</span>
              <span style={{ fontWeight: 800, fontSize: '1.15rem', color: '#2D231E' }}>
                {inv.total_pembayaran}
              </span>
            </div>
          </div>

          {/* CATATAN FOOTNOTE SESUAI FIGMA */}
          <div
            style={{
              paddingTop: '1.5rem',
              borderTop: '1px solid #F0ECE6',
              fontSize: '0.85rem',
              color: '#8E7D73',
              lineHeight: 1.5
            }}
          >
            Status bayar tidak disatukan dengan status proses.
          </div>
        </div>

        {/* Tautan Pintas Tambahan untuk Navigasi */}
        <div
          className="no-print"
          style={{
            marginTop: '2rem',
            display: 'flex',
            justifyContent: 'center',
            gap: '1.5rem'
          }}
        >
          <Link
            to="/katalog"
            style={{
              color: '#43281C',
              fontSize: '0.875rem',
              fontWeight: 600,
              textDecoration: 'none'
            }}
          >
            &larr; Belanja Lagi
          </Link>

          <Link
            to="/owner/pesanan"
            style={{
              color: '#8E7D73',
              fontSize: '0.875rem',
              fontWeight: 500,
              textDecoration: 'none'
            }}
          >
            Halaman Pesanan Owner &rarr;
          </Link>
        </div>

      </div>

      {/* Style Print Mode */}
      <style>{`
        @media print {
          .no-print {
            display: none !important;
          }
          body {
            background-color: #FFFFFF !important;
          }
        }
      `}</style>
    </div>
  );
}
