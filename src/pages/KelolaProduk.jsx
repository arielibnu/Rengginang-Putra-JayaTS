
import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const API_URL = 'http://localhost:8080';

const produkKosong = {
  nama_produk: '',
  deskripsi: '',
  status_aktif: true
};

const varianKosong = {
  nama_varian: '',
  harga_jual: '',
  berat_gram: '',
  stok: '',
  status_aktif: true
};

const formatRupiah = (angka) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(Number(angka) || 0);

// Koneksi frontend dengan Backend Golang
async function apiRequest(endpoint, options = {}) {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  });

  const hasil = await response.json().catch(() => ({}));

  if (!response.ok) {
    const error = new Error(
      hasil.pesan || `Permintaan gagal (${response.status})`
    );
    error.status = response.status;
    throw error;
  }

  return hasil.data;
}

const tombol = {
  background: '#456b55',
  color: 'white',
  padding: '10px 15px',
  border: 'none',
  borderRadius: 7,
  cursor: 'pointer'
};

const tombolSekunder = {
  ...tombol,
  background: '#e9eee9',
  color: '#345443'
};

const tombolHapus = {
  ...tombol,
  background: '#fff0ee',
  color: '#ae433a'
};

const kartu = {
  background: '#ffffff',
  border: '1px solid #e2e9e2',
  borderRadius: 12,
  padding: 20,
  marginBottom: 20
};

const inputStyle = {
  display: 'block',
  width: '100%',
  padding: 10,
  marginTop: 6,
  border: '1px solid #d2ddd3',
  borderRadius: 7,
  boxSizing: 'border-box',
  font: 'inherit'
};

export default function KelolaProduk() {
  const navigate = useNavigate();

  const [produk, setProduk] = useState([]);
  const [varian, setVarian] = useState([]);
  const [produkDipilih, setProdukDipilih] = useState(null);

  const [formProduk, setFormProduk] = useState(produkKosong);
  const [formVarian, setFormVarian] = useState(varianKosong);

  const [editProduk, setEditProduk] = useState(null);
  const [editVarian, setEditVarian] = useState(null);

  const [tampilFormProduk, setTampilFormProduk] = useState(false);
  const [tampilFormVarian, setTampilFormVarian] = useState(false);

  const [memuat, setMemuat] = useState(true);
  const [menyimpan, setMenyimpan] = useState(false);
  const [memuatVarian, setMemuatVarian] = useState(false);

  const [pesan, setPesan] = useState('');
  const [kesalahan, setKesalahan] = useState('');

  const tampilkanError = (error) => {
    if (error.status === 401) {
      navigate('/login', { replace: true });
      return;
    }

    setKesalahan(error.message || 'Terjadi kesalahan.');
  };

  const bersihkanPesan = () => {
    setPesan('');
    setKesalahan('');
  };

  const ambilProduk = async () => {
    const data = await apiRequest('/api/admin/produk');

    const daftar = Array.isArray(data) ? data : [];
    setProduk(daftar);

    setProdukDipilih((sebelumnya) => {
      if (!sebelumnya) return null;

      return daftar.find(
        (item) => item.id_produk === sebelumnya.id_produk
      ) || null;
    });
  };

  const ambilVarian = async (idProduk) => {
    setMemuatVarian(true);

    try {
      const data = await apiRequest(
        `/api/admin/produk/${idProduk}/varian`
      );

      setVarian(Array.isArray(data) ? data : []);
    } finally {
      setMemuatVarian(false);
    }
  };

  // Periksa sesi Owner sebelum menampilkan data
  useEffect(() => {
    let aktif = true;

    const mulai = async () => {
      try {
        await apiRequest('/api/pemilik/profil');
        await apiRequest('/api/admin/beranda');

        const data = await apiRequest('/api/admin/produk');

        if (aktif) {
          setProduk(Array.isArray(data) ? data : []);
        }
      } catch (error) {
        if (!aktif) return;

        if (error.status === 401 || error.status === 403) {
          navigate('/login', { replace: true });
        } else {
          setKesalahan(error.message);
        }
      } finally {
        if (aktif) setMemuat(false);
      }
    };

    mulai();

    return () => {
      aktif = false;
    };
  }, [navigate]);

  // Pilih produk untuk menampilkan varian
  const pilihProduk = async (item) => {
    bersihkanPesan();

    setProdukDipilih(item);
    setVarian([]);
    setTampilFormVarian(false);

    try {
      await ambilVarian(item.id_produk);
    } catch (error) {
      tampilkanError(error);
    }
  };

  // FORM PRODUK
  const tambahProduk = () => {
    bersihkanPesan();
    setEditProduk(null);
    setFormProduk({ ...produkKosong });
    setTampilFormProduk(true);
  };

  const ubahProduk = (item) => {
    bersihkanPesan();

    setEditProduk(item);
    setFormProduk({
      nama_produk: item.nama_produk,
      deskripsi: item.deskripsi || '',
      status_aktif: Boolean(item.status_aktif)
    });

    setTampilFormProduk(true);
  };

  const simpanProduk = async (event) => {
    event.preventDefault();

    if (menyimpan) return;

    const nama = formProduk.nama_produk.trim();

    if (!nama || nama.length > 150) {
      setKesalahan('Nama produk wajib diisi, maksimal 150 karakter.');
      return;
    }

    setMenyimpan(true);
    bersihkanPesan();

    try {
      const endpoint = editProduk
        ? `/api/admin/produk/${editProduk.id_produk}`
        : '/api/admin/produk';

      await apiRequest(endpoint, {
        method: editProduk ? 'PUT' : 'POST',
        body: JSON.stringify({
          nama_produk: nama,
          deskripsi: formProduk.deskripsi.trim(),
          status_aktif: formProduk.status_aktif
        })
      });

      await ambilProduk();

      setTampilFormProduk(false);
      setEditProduk(null);
      setPesan('Produk berhasil disimpan.');
    } catch (error) {
      tampilkanError(error);
    } finally {
      setMenyimpan(false);
    }
  };

  const hapusProduk = async (item) => {
    if (menyimpan || !item.status_aktif) return;

    const yakin = window.confirm(
      `Nonaktifkan produk "${item.nama_produk}"?`
    );

    if (!yakin) return;

    setMenyimpan(true);
    bersihkanPesan();

    try {
      await apiRequest(
        `/api/admin/produk/${item.id_produk}`,
        { method: 'DELETE' }
      );

      await ambilProduk();

      if (produkDipilih?.id_produk === item.id_produk) {
        setTampilFormVarian(false);
      }

      setPesan('Produk berhasil dinonaktifkan.');
    } catch (error) {
      tampilkanError(error);
    } finally {
      setMenyimpan(false);
    }
  };

  // FORM VARIAN
  const tambahVarian = () => {
    if (!produkDipilih?.status_aktif) return;

    bersihkanPesan();
    setEditVarian(null);
    setFormVarian({ ...varianKosong });
    setTampilFormVarian(true);
  };

  const ubahVarian = (item) => {
    bersihkanPesan();

    setEditVarian(item);
    setFormVarian({
      nama_varian: item.nama_varian,
      harga_jual: String(item.harga_jual),
      berat_gram: String(item.berat_gram),
      stok: String(item.stok),
      status_aktif: Boolean(item.status_aktif)
    });

    setTampilFormVarian(true);
  };

  const simpanVarian = async (event) => {
    event.preventDefault();

    if (menyimpan || !produkDipilih) return;

    const nama = formVarian.nama_varian.trim();

    const harga = Number(formVarian.harga_jual);
    const berat = Number(formVarian.berat_gram);
    const stok = Number(formVarian.stok);

    if (!nama || nama.length > 100) {
      setKesalahan('Nama varian wajib diisi, maksimal 100 karakter.');
      return;
    }

    if (
      formVarian.harga_jual === '' ||
      !Number.isSafeInteger(harga) ||
      harga < 0 ||
      harga > 9999999999
    ) {
      setKesalahan('Harga jual harus berupa angka yang valid.');
      return;
    }

    if (
      formVarian.berat_gram === '' ||
      !Number.isSafeInteger(berat) ||
      berat <= 0 ||
      berat > 100000000
    ) {
      setKesalahan('Berat harus lebih dari 0 gram.');
      return;
    }

    if (
      formVarian.stok === '' ||
      !Number.isSafeInteger(stok) ||
      stok < 0 ||
      stok > 1000000000
    ) {
      setKesalahan('Stok harus berupa angka yang valid.');
      return;
    }

    setMenyimpan(true);
    bersihkanPesan();

    const idProduk = produkDipilih.id_produk;

    try {
      const endpoint = editVarian
        ? `/api/admin/varian/${editVarian.id_varian}`
        : `/api/admin/produk/${idProduk}/varian`;

      await apiRequest(endpoint, {
        method: editVarian ? 'PUT' : 'POST',
        body: JSON.stringify({
          nama_varian: nama,
          harga_jual: harga,
          berat_gram: berat,
          stok,
          status_aktif: formVarian.status_aktif
        })
      });

      await ambilVarian(idProduk);

      setTampilFormVarian(false);
      setEditVarian(null);
      setPesan('Varian berhasil disimpan.');
    } catch (error) {
      tampilkanError(error);
    } finally {
      setMenyimpan(false);
    }
  };

  const hapusVarian = async (item) => {
    if (menyimpan || !item.status_aktif || !produkDipilih) return;

    const yakin = window.confirm(
      `Nonaktifkan varian "${item.nama_varian}"?`
    );

    if (!yakin) return;

    setMenyimpan(true);
    bersihkanPesan();

    try {
      await apiRequest(
        `/api/admin/varian/${item.id_varian}`,
        { method: 'DELETE' }
      );

      await ambilVarian(produkDipilih.id_produk);

      setTampilFormVarian(false);
      setPesan('Varian berhasil dinonaktifkan.');
    } catch (error) {
      tampilkanError(error);
    } finally {
      setMenyimpan(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      background: '#f5f7f4',
      color: '#283b2e',
      padding: '30px 18px'
    }}>
      <main style={{
        maxWidth: 1050,
        margin: '0 auto'
      }}>

        <Link
          to="/admin/dashboard"
          style={{ color: '#456b55', textDecoration: 'none' }}
        >
          ← Kembali ke Dashboard
        </Link>

        <header style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 15,
          margin: '24px 0'
        }}>
          <div>
            <h1 style={{ marginBottom: 6 }}>Kelola Produk</h1>
            <p style={{ color: '#708073', margin: 0 }}>
              Pengelolaan produk dan varian Rengginang Putra Jaya
            </p>
          </div>

          <button
            type="button"
            style={tombol}
            onClick={tambahProduk}
            disabled={menyimpan || memuat}
          >
            + Tambah Produk
          </button>
        </header>

        {kesalahan && (
          <div role="alert" style={{
            ...kartu,
            color: '#a63232',
            borderColor: '#f0cece'
          }}>
            {kesalahan}
          </div>
        )}

        {pesan && (
          <div role="status" style={{
            ...kartu,
            color: '#246844',
            borderColor: '#c5e0ce'
          }}>
            {pesan}
          </div>
        )}

        {/* FORM TAMBAH / EDIT PRODUK */}
        {tampilFormProduk && (
          <form onSubmit={simpanProduk} style={kartu}>
            <h2>
              {editProduk ? 'Edit Produk' : 'Tambah Produk'}
            </h2>

            <div style={{ display: 'grid', gap: 15 }}>
              <label>
                Nama Produk
                <input
                  style={inputStyle}
                  required
                  maxLength={150}
                  value={formProduk.nama_produk}
                  onChange={(e) => setFormProduk({
                    ...formProduk,
                    nama_produk: e.target.value
                  })}
                />
              </label>

              <label>
                Deskripsi
                <textarea
                  style={{ ...inputStyle, minHeight: 80 }}
                  value={formProduk.deskripsi}
                  onChange={(e) => setFormProduk({
                    ...formProduk,
                    deskripsi: e.target.value
                  })}
                />
              </label>

              <label>
                <input
                  type="checkbox"
                  checked={formProduk.status_aktif}
                  onChange={(e) => setFormProduk({
                    ...formProduk,
                    status_aktif: e.target.checked
                  })}
                />
                {' '}Produk aktif
              </label>

              <div style={{ display: 'flex', gap: 10 }}>
                <button
                  type="submit"
                  style={tombol}
                  disabled={menyimpan}
                >
                  {menyimpan ? 'Menyimpan...' : 'Simpan Produk'}
                </button>

                <button
                  type="button"
                  style={tombolSekunder}
                  onClick={() => setTampilFormProduk(false)}
                  disabled={menyimpan}
                >
                  Batal
                </button>
              </div>
            </div>
          </form>
        )}

        {/* DAFTAR PRODUK */}
        <section style={kartu}>
          <h2>Daftar Produk</h2>

          {memuat ? (
            <p>Memuat data produk...</p>
          ) : produk.length === 0 ? (
            <p>Belum ada produk. Klik Tambah Produk.</p>
          ) : (
            produk.map((item) => (
              <div
                key={item.id_produk}
                style={{
                  padding: '16px 0',
                  borderBottom: '1px solid #e4e9e4',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  flexWrap: 'wrap',
                  gap: 15
                }}
              >
                <div>
                  <strong>{item.nama_produk}</strong>

                  <p style={{
                    margin: '5px 0',
                    color: '#68766b'
                  }}>
                    {item.deskripsi || 'Tanpa deskripsi'}
                  </p>

                  <small>
                    {item.status_aktif ? '● Aktif' : '○ Nonaktif'}
                  </small>
                </div>

                <div style={{
                  display: 'flex',
                  gap: 8,
                  flexWrap: 'wrap'
                }}>
                  <button
                    type="button"
                    style={tombolSekunder}
                    onClick={() => pilihProduk(item)}
                    disabled={menyimpan || memuatVarian}
                  >
                    Lihat Varian
                  </button>

                  <button
                    type="button"
                    style={tombolSekunder}
                    onClick={() => ubahProduk(item)}
                    disabled={menyimpan}
                  >
                    Edit
                  </button>

                  <button
                    type="button"
                    style={tombolHapus}
                    onClick={() => hapusProduk(item)}
                    disabled={menyimpan || !item.status_aktif}
                  >
                    Hapus
                  </button>
                </div>
              </div>
            ))
          )}
        </section>

        {/* DAFTAR VARIAN PRODUK */}
        {produkDipilih && (
          <section style={kartu}>
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              flexWrap: 'wrap',
              gap: 12
            }}>
              <h2>Varian: {produkDipilih.nama_produk}</h2>

              <button
                type="button"
                style={tombol}
                onClick={tambahVarian}
                disabled={
                  menyimpan ||
                  memuatVarian ||
                  !produkDipilih.status_aktif
                }
              >
                + Tambah Varian
              </button>
            </div>

            {memuatVarian ? (
              <p>Memuat data varian...</p>
            ) : varian.length === 0 ? (
              <p>Belum ada varian untuk produk ini.</p>
            ) : (
              varian.map((item) => (
                <div
                  key={item.id_varian}
                  style={{
                    padding: '15px 0',
                    borderBottom: '1px solid #e4e9e4',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    flexWrap: 'wrap',
                    gap: 12
                  }}
                >
                  <div>
                    <strong>{item.nama_varian}</strong>

                    <p style={{ margin: '6px 0' }}>
                      {formatRupiah(item.harga_jual)}
                      {' · '}{item.berat_gram} gram
                      {' · '}Stok: {item.stok}
                    </p>

                    <small>
                      {item.status_aktif ? 'Aktif' : 'Nonaktif'}
                    </small>
                  </div>

                  <div style={{ display: 'flex', gap: 8 }}>
                    <button
                      type="button"
                      style={tombolSekunder}
                      onClick={() => ubahVarian(item)}
                      disabled={menyimpan}
                    >
                      Edit
                    </button>

                    <button
                      type="button"
                      style={tombolHapus}
                      onClick={() => hapusVarian(item)}
                      disabled={menyimpan || !item.status_aktif}
                    >
                      Hapus
                    </button>
                  </div>
                </div>
              ))
            )}

            {/* FORM TAMBAH / EDIT VARIAN */}
            {tampilFormVarian && (
              <form
                onSubmit={simpanVarian}
                style={{
                  marginTop: 25,
                  padding: 20,
                  background: '#f7f9f6',
                  borderRadius: 10
                }}
              >
                <h3>
                  {editVarian ? 'Edit Varian' : 'Tambah Varian'}
                </h3>

                <div style={{ display: 'grid', gap: 15 }}>
                  <label>
                    Nama Varian
                    <input
                      style={inputStyle}
                      required
                      maxLength={100}
                      value={formVarian.nama_varian}
                      onChange={(e) => setFormVarian({
                        ...formVarian,
                        nama_varian: e.target.value
                      })}
                    />
                  </label>

                  {[
                    ['harga_jual', 'Harga Jual (Rp)'],
                    ['berat_gram', 'Berat (gram)'],
                    ['stok', 'Stok']
                  ].map(([field, label]) => (
                    <label key={field}>
                      {label}
                      <input
                        style={inputStyle}
                        type="number"
                        min={field === 'berat_gram' ? 1 : 0}
                        step="1"
                        required
                        value={formVarian[field]}
                        onChange={(e) => setFormVarian({
                          ...formVarian,
                          [field]: e.target.value
                        })}
                      />
                    </label>
                  ))}

                  <label>
                    <input
                      type="checkbox"
                      checked={formVarian.status_aktif}
                      onChange={(e) => setFormVarian({
                        ...formVarian,
                        status_aktif: e.target.checked
                      })}
                    />
                    {' '}Varian aktif
                  </label>

                  <div style={{ display: 'flex', gap: 10 }}>
                    <button
                      type="submit"
                      style={tombol}
                      disabled={menyimpan}
                    >
                      {menyimpan ? 'Menyimpan...' : 'Simpan Varian'}
                    </button>

                    <button
                      type="button"
                      style={tombolSekunder}
                      onClick={() => setTampilFormVarian(false)}
                      disabled={menyimpan}
                    >
                      Batal
                    </button>
                  </div>
                </div>
              </form>
            )}
          </section>
        )}
      </main>
    </div>
  );
}
