import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  Circle,
  Square,
  LayoutList,
  Diamond,
  Plus,
  X,
  AlertCircle,
  CheckCircle,
  Trash2,
  Image as ImageIcon
} from 'lucide-react';
import '../styles/admin.css';

const API_URL = 'http://localhost:8080';
const STORAGE_KEY = 'rengginang_owner_products';

// 4 varian default sesuai mockup resmi (Image 2)
const DEFAULT_PRODUCTS = [
  {
    id: 'udang',
    name: 'Rengginang Udang',
    variant: 'Udang',
    subtitle: 'Varian Udang • Foto asli',
    price: 18000,
    priceDisplay: '18000 • demo',
    priceFormatted: 'Rp28.000',
    weight: '450 gr',
    status: 'Aktif',
    hasTransactions: false,
    image: '/images/products/rengginang-udang.jpg',
    description: 'Rengginang tradisional rasa udang'
  },
  {
    id: 'bawang',
    name: 'Rengginang Bawang',
    variant: 'Bawang',
    subtitle: 'Varian Bawang • Foto asli',
    price: 28000,
    priceFormatted: 'Rp28.000',
    weight: '450 gr',
    status: 'Aktif',
    hasTransactions: false,
    image: '/images/products/rengginang-bawang.jpg',
    description: 'Rengginang renyah beraroma bawang putih gurih alami resep tradisional.'
  },
  {
    id: 'cumi',
    name: 'Rengginang Cumi',
    variant: 'Cumi',
    subtitle: 'Varian Cumi • Foto asli',
    price: 28000,
    priceFormatted: 'Rp28.000',
    weight: '450 gr',
    status: 'Aktif',
    hasTransactions: false,
    image: '/images/products/rengginang-cumi.jpg',
    description: 'Rengginang spesial tinta cumi dengan potongan cumi asli beraroma sedap.'
  },
  {
    id: 'terasi',
    name: 'Rengginang Terasi',
    variant: 'Terasi',
    subtitle: 'Varian Terasi • Foto asli',
    price: 28000,
    priceFormatted: 'Rp28.000',
    weight: '450 gr',
    status: 'Aktif',
    hasTransactions: false,
    image: '/images/products/rengginang-terasi.jpg',
    description: 'Rengginang legendaris dengan aroma terasi khas pesisir yang harum dan otentik.'
  }
];

export default function KelolaProduk() {
  const navigate = useNavigate();

  // State autentikasi Owner
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [authErrorMessage, setAuthErrorMessage] = useState('');

  // State Produk (PPL-03)
  const [products, setProducts] = useState(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) {
        return JSON.parse(saved);
      }
    } catch (e) {
      console.warn('Gagal membaca data produk:', e);
    }
    return DEFAULT_PRODUCTS;
  });

  // State Modal
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    variant: '',
    subtitle: '',
    price: '28000',
    weight: '450 gr',
    status: 'Aktif',
    image: '/images/products/rengginang-udang.jpg',
    description: ''
  });

  const [formError, setFormError] = useState('');
  const [notification, setNotification] = useState(null);

  // Simpan produk ke localStorage setiap ada perubahan
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(products));
    } catch (e) {
      console.warn('Gagal menyimpan produk:', e);
    }
  }, [products]);

  // Validasi sesi Owner langsung ke backend
  useEffect(() => {
    const controller = new AbortController();

    const checkOwnerSession = async () => {
      try {
        const options = {
          method: 'GET',
          credentials: 'include',
          signal: controller.signal
        };

        const profileResponse = await fetch(`${API_URL}/api/pemilik/profil`, options);

        if (profileResponse.status === 401) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!profileResponse.ok) {
          throw new Error('Gagal memeriksa sesi Owner.');
        }

        const adminResponse = await fetch(`${API_URL}/api/admin/beranda`, options);

        if (adminResponse.status === 401 || adminResponse.status === 403) {
          localStorage.removeItem('ownerAuth');
          navigate('/login', { replace: true });
          return;
        }

        if (!adminResponse.ok) {
          throw new Error('Gagal memeriksa akses Admin.');
        }

        if (!controller.signal.aborted) {
          localStorage.removeItem('ownerAuth');
          setIsAuthorized(true);
          setAuthErrorMessage('');
        }
      } catch (error) {
        if (error.name !== 'AbortError' && !controller.signal.aborted) {
          console.warn('Backend offline, masuk mode frontend mandiri:', error.message);
          setIsAuthorized(true);
          setAuthErrorMessage('');
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsChecking(false);
        }
      }
    };

    checkOwnerSession();
    return () => controller.abort();
  }, [navigate]);

  // Handler Logout
  const handleLogout = async () => {
    if (isLoggingOut) return;
    setIsLoggingOut(true);

    try {
      await fetch(`${API_URL}/api/pemilik/keluar`, {
        method: 'POST',
        credentials: 'include'
      });
    } catch (error) {
      console.warn('Logout fallback:', error);
    } finally {
      localStorage.removeItem('ownerAuth');
      setIsAuthorized(false);
      navigate('/login', { replace: true });
      setIsLoggingOut(false);
    }
  };

  // Helper Notifikasi Toast
  const showToast = (message, type = 'success') => {
    setNotification({ message, type });
    setTimeout(() => {
      setNotification(null);
    }, 3500);
  };

  // Helper format rupiah
  const formatRupiah = (num) => {
    const val = parseInt(num, 10);
    if (isNaN(val)) return 'Rp0';
    return 'Rp' + val.toLocaleString('id-ID');
  };

  // Buka Halaman Tambah Produk Sesuai Mockup
  const handleOpenAddModal = () => {
    navigate('/admin/produk/tambah');
  };

  // Buka Halaman Edit Produk Sesuai Mockup
  const handleOpenEditModal = (product) => {
    navigate(`/admin/produk/edit/${product.id}`);
  };

  // Buka Modal Hapus / Nonaktifkan
  const handleOpenDeleteModal = (product) => {
    setSelectedProduct(product);
    setIsDeleteModalOpen(true);
  };

  // Submit Tambah Produk (PPL-03 AC 1)
  const handleSubmitAdd = (e) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      setFormError('Nama produk wajib diisi.');
      return;
    }
    const numPrice = parseInt(formData.price, 10);
    if (isNaN(numPrice) || numPrice <= 0) {
      setFormError('Harga jual harus berupa angka lebih dari 0.');
      return;
    }
    if (!formData.weight.trim()) {
      setFormError('Berat produk wajib diisi.');
      return;
    }

    const newId = 'prod-' + Date.now();
    const newProduct = {
      id: newId,
      name: formData.name.trim(),
      variant: formData.variant.trim() || formData.name.trim(),
      subtitle: formData.subtitle.trim() || `Varian ${formData.name.trim()} • Foto asli`,
      price: numPrice,
      priceFormatted: formatRupiah(numPrice),
      weight: formData.weight.trim(),
      status: formData.status,
      hasTransactions: false,
      image: formData.image,
      description: formData.description.trim()
    };

    setProducts((prev) => [...prev, newProduct]);
    setIsAddModalOpen(false);
    showToast(`Produk "${newProduct.name}" berhasil ditambahkan!`);
  };

  // Submit Edit Produk (PPL-03 AC 1 & 3)
  const handleSubmitEdit = (e) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      setFormError('Nama produk wajib diisi.');
      return;
    }
    const numPrice = parseInt(formData.price, 10);
    if (isNaN(numPrice) || numPrice <= 0) {
      setFormError('Harga jual harus berupa angka lebih dari 0.');
      return;
    }
    if (!formData.weight.trim()) {
      setFormError('Berat produk wajib diisi.');
      return;
    }

    setProducts((prev) =>
      prev.map((item) => {
        if (item.id === selectedProduct.id) {
          return {
            ...item,
            name: formData.name.trim(),
            variant: formData.variant.trim() || formData.name.trim(),
            subtitle: formData.subtitle.trim() || `Varian ${formData.name.trim()} • Foto asli`,
            price: numPrice,
            priceFormatted: formatRupiah(numPrice),
            weight: formData.weight.trim(),
            status: formData.status,
            image: formData.image,
            description: formData.description.trim()
          };
        }
        return item;
      })
    );

    setIsEditModalOpen(false);
    showToast(`Perubahan data "${formData.name.trim()}" berhasil disimpan!`);
  };

  // Nonaktifkan Produk (PPL-03 AC 2)
  const handleDeactivate = () => {
    if (!selectedProduct) return;
    setProducts((prev) =>
      prev.map((item) =>
        item.id === selectedProduct.id ? { ...item, status: 'Nonaktif' } : item
      )
    );
    setIsDeleteModalOpen(false);
    showToast(`Produk "${selectedProduct.name}" telah dinonaktifkan.`);
  };

  // Hapus Permanen Produk (PPL-03 AC 2)
  const handlePermanentDelete = () => {
    if (!selectedProduct) return;
    setProducts((prev) => prev.filter((item) => item.id !== selectedProduct.id));
    setIsDeleteModalOpen(false);
    showToast(`Produk "${selectedProduct.name}" telah dihapus secara permanen.`);
  };

  // Unggah foto lokal
  const handleImageFileChange = (e) => {
    const file = e.target.files && e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        setFormData((prev) => ({ ...prev, image: reader.result }));
      };
      reader.readAsDataURL(file);
    }
  };

  // Loading Screen
  if (isChecking) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Memeriksa akses Owner</h1>
          <p className="login-subtitle">Memvalidasi sesi login melalui backend...</p>
        </div>
      </div>
    );
  }

  // Not Authorized Screen
  if (!isAuthorized) {
    return (
      <div className="login-page">
        <div className="login-card">
          <div className="login-logo">PJ</div>
          <h1 className="login-title">Akses belum tersedia</h1>
          <p className="login-subtitle">{authErrorMessage || 'Sesi Owner tidak valid.'}</p>
          <button
            type="button"
            className="btn-submit-login"
            onClick={() => window.location.reload()}
          >
            Coba lagi
          </button>
          <Link to="/login" className="btn-back-home">
            Kembali ke login
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-layout">
      {/* Toast Notifikasi */}
      {notification && (
        <div className={`toast-notification ${notification.type}`}>
          {notification.type === 'success' ? (
            <CheckCircle size={18} />
          ) : (
            <AlertCircle size={18} />
          )}
          <span>{notification.message}</span>
        </div>
      )}

      {/* Sidebar Kiri Sesuai Figma Mockup (Image 2) */}
      <aside className="admin-sidebar">
        <div>
          {/* Brand Logo */}
          <div className="sidebar-brand">
            <div className="sidebar-logo">PJ</div>
            <div>
              <div className="sidebar-brand-name">PUTRA JAYA</div>
              <div className="sidebar-brand-sub">Dashboard Owner</div>
            </div>
          </div>

          {/* Menu Navigasi Admin */}
          <ul className="sidebar-menu">
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/dashboard')}
              >
                <Circle size={16} />
                <span>Ringkasan</span>
              </div>
            </li>
            <li>
              <div className="sidebar-item active">
                <Square size={16} />
                <span>Kelola Produk</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/stok')}
              >
                <LayoutList size={16} />
                <span>Data Stok</span>
              </div>
            </li>
            <li>
              <div
                className="sidebar-item"
                onClick={() => navigate('/admin/biaya')}
              >
                <Diamond size={16} />
                <span>Modal & Biaya</span>
              </div>
            </li>
          </ul>
        </div>

        {/* Footer Sidebar */}
        <div className="sidebar-footer">
          <Link to="/" className="sidebar-footer-link">
            Buka website publik
          </Link>

          <button
            type="button"
            onClick={handleLogout}
            disabled={isLoggingOut}
            className="sidebar-footer-link"
            style={{ color: '#A92631' }}
          >
            {isLoggingOut ? 'Memproses logout...' : 'Keluar'}
          </button>
        </div>
      </aside>

      {/* Konten Utama Kanan Sesuai Mockup Gambar 2 */}
      <main className="admin-main">
        {/* Breadcrumb: Dashboard Owner / Kelola Produk */}
        <div className="admin-breadcrumb">
          Dashboard Owner / Kelola Produk
        </div>

        {/* Header Halaman */}
        <div className="product-page-header">
          <div>
            <h1 className="product-page-title">Kelola produk</h1>
            <p className="product-page-subtitle">
              Tambah atau perbarui informasi varian rengginang.
            </p>
          </div>
          <button
            type="button"
            className="btn-add-product"
            onClick={handleOpenAddModal}
          >
            <Plus size={18} />
            <span>Tambah produk</span>
          </button>
        </div>

        {/* Tabel Produk (Card Putih) */}
        <div className="products-table-card">
          <div className="table-responsive">
            <table className="products-table">
              <thead>
                <tr>
                  <th style={{ width: '40%' }}>Produk</th>
                  <th style={{ width: '18%' }}>Harga</th>
                  <th style={{ width: '15%' }}>Berat</th>
                  <th style={{ width: '12%' }}>Status</th>
                  <th style={{ width: '15%', textAlign: 'right' }}></th>
                </tr>
              </thead>
              <tbody>
                {products.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="empty-table-cell">
                      Belum ada data produk. Klik <strong>+ Tambah produk</strong> untuk menambahkan varian baru.
                    </td>
                  </tr>
                ) : (
                  products.map((p) => (
                    <tr key={p.id}>
                      {/* Kolom Produk (Foto + Nama + Subtitle) */}
                      <td>
                        <div className="product-cell-flex">
                          <img
                            src={p.image}
                            alt={p.name}
                            className="product-table-img"
                            onError={(e) => {
                              e.target.onerror = null;
                              e.target.src = '/images/products/rengginang-udang.jpg';
                            }}
                          />
                          <div>
                            <div className="product-table-name">{p.name}</div>
                            <div className="product-table-sub">
                              {p.subtitle || `Varian ${p.name} • Foto asli`}
                            </div>
                          </div>
                        </div>
                      </td>

                      {/* Kolom Harga */}
                      <td>
                        <span className="product-table-price">
                          {p.priceFormatted || formatRupiah(p.price)}
                        </span>
                      </td>

                      {/* Kolom Berat */}
                      <td>
                        <span className="product-table-weight">{p.weight}</span>
                      </td>

                      {/* Kolom Status (Badge Aktif/Nonaktif) */}
                      <td>
                        <span
                          className={
                            p.status === 'Aktif'
                              ? 'badge-status-active'
                              : 'badge-status-inactive'
                          }
                        >
                          {p.status}
                        </span>
                      </td>

                      {/* Kolom Aksi (Hapus & Edit ›) */}
                      <td style={{ textAlign: 'right' }}>
                        <div className="table-actions-flex">
                          <button
                            type="button"
                            className="btn-action-delete"
                            onClick={() => handleOpenDeleteModal(p)}
                          >
                            Hapus
                          </button>
                          <button
                            type="button"
                            className="btn-action-edit"
                            onClick={() => handleOpenEditModal(p)}
                          >
                            <span>Edit</span>
                            <span className="chevron-symbol">›</span>
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Catatan Kaki Mockup */}
        <div className="table-footnote">
          Harga pada mockup hanya contoh dan perlu disesuaikan dengan data Owner.
        </div>
      </main>

      {/* ================================================================= */}
      {/* MODAL 1: TAMBAH PRODUK (PPL-03 AC 1)                               */}
      {/* ================================================================= */}
      {isAddModalOpen && (
        <div className="modal-backdrop" onClick={() => setIsAddModalOpen(false)}>
          <div
            className="modal-box"
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
          >
            <div className="modal-header">
              <h2 className="modal-title">Tambah Produk Baru</h2>
              <button
                type="button"
                className="modal-close-btn"
                onClick={() => setIsAddModalOpen(false)}
              >
                <X size={20} />
              </button>
            </div>

            {formError && (
              <div className="modal-alert-error">
                <AlertCircle size={16} />
                <span>{formError}</span>
              </div>
            )}

            <form onSubmit={handleSubmitAdd} className="modal-form">
              <div className="form-group">
                <label className="form-label">Nama Produk / Varian *</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="Contoh: Rengginang Terasi"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  required
                />
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label className="form-label">Harga Jual (Rp) *</label>
                  <input
                    type="number"
                    className="form-input"
                    placeholder="Contoh: 28000"
                    min="1"
                    value={formData.price}
                    onChange={(e) =>
                      setFormData({ ...formData, price: e.target.value })
                    }
                    required
                  />
                </div>

                <div className="form-group">
                  <label className="form-label">Berat Bersih *</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="Contoh: 450 gr"
                    value={formData.weight}
                    onChange={(e) =>
                      setFormData({ ...formData, weight: e.target.value })
                    }
                    required
                  />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Keterangan Subtitle</label>
                <input
                  type="text"
                  className="form-input"
                  placeholder="Contoh: Varian Terasi • Foto asli"
                  value={formData.subtitle}
                  onChange={(e) =>
                    setFormData({ ...formData, subtitle: e.target.value })
                  }
                />
              </div>

              <div className="form-group">
                <label className="form-label">Pilihan Foto Produk</label>
                <div className="photo-selection-grid">
                  {[
                    { label: 'Udang', src: '/images/products/rengginang-udang.jpg' },
                    { label: 'Bawang', src: '/images/products/rengginang-bawang.jpg' },
                    { label: 'Cumi', src: '/images/products/rengginang-cumi.jpg' },
                    { label: 'Terasi', src: '/images/products/rengginang-terasi.jpg' },
                    { label: 'Hero Mix', src: '/images/products/rengginang-all-hero.jpg' }
                  ].map((preset) => (
                    <button
                      key={preset.src}
                      type="button"
                      className={`photo-preset-btn ${formData.image === preset.src ? 'selected' : ''}`}
                      onClick={() => setFormData({ ...formData, image: preset.src })}
                    >
                      <img src={preset.src} alt={preset.label} />
                      <span>{preset.label}</span>
                    </button>
                  ))}
                </div>
                <div className="custom-photo-upload">
                  <label className="upload-file-label">
                    <ImageIcon size={16} />
                    <span>Atau unggah foto sendiri dari perangkat</span>
                    <input
                      type="file"
                      accept="image/*"
                      style={{ display: 'none' }}
                      onChange={handleImageFileChange}
                    />
                  </label>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Status Produk</label>
                <div className="status-radio-group">
                  <label className="status-radio-label">
                    <input
                      type="radio"
                      name="status"
                      value="Aktif"
                      checked={formData.status === 'Aktif'}
                      onChange={() => setFormData({ ...formData, status: 'Aktif' })}
                    />
                    <span className="radio-text">
                      <strong>Aktif</strong> (Dapat dilihat dan dibeli oleh pelanggan)
                    </span>
                  </label>
                  <label className="status-radio-label">
                    <input
                      type="radio"
                      name="status"
                      value="Nonaktif"
                      checked={formData.status === 'Nonaktif'}
                      onChange={() => setFormData({ ...formData, status: 'Nonaktif' })}
                    />
                    <span className="radio-text">
                      <strong>Nonaktif</strong> (Disembunyikan dari katalog pelanggan)
                    </span>
                  </label>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Deskripsi Rinci (Opsional)</label>
                <textarea
                  className="form-textarea"
                  rows="3"
                  placeholder="Tambahkan penjelasan rasa, bahan, atau keunggulan produk..."
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                />
              </div>

              <div className="modal-actions">
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => setIsAddModalOpen(false)}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-save">
                  Simpan Produk
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ================================================================= */}
      {/* MODAL 2: EDIT PRODUK (PPL-03 AC 1 & 3)                             */}
      {/* ================================================================= */}
      {isEditModalOpen && selectedProduct && (
        <div className="modal-backdrop" onClick={() => setIsEditModalOpen(false)}>
          <div
            className="modal-box"
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
          >
            <div className="modal-header">
              <h2 className="modal-title">Edit Produk</h2>
              <button
                type="button"
                className="modal-close-btn"
                onClick={() => setIsEditModalOpen(false)}
              >
                <X size={20} />
              </button>
            </div>

            {formError && (
              <div className="modal-alert-error">
                <AlertCircle size={16} />
                <span>{formError}</span>
              </div>
            )}

            <form onSubmit={handleSubmitEdit} className="modal-form">
              <div className="form-group">
                <label className="form-label">Nama Produk / Varian *</label>
                <input
                  type="text"
                  className="form-input"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  required
                />
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label className="form-label">Harga Jual (Rp) *</label>
                  <input
                    type="number"
                    className="form-input"
                    min="1"
                    value={formData.price}
                    onChange={(e) =>
                      setFormData({ ...formData, price: e.target.value })
                    }
                    required
                  />
                </div>

                <div className="form-group">
                  <label className="form-label">Berat Bersih *</label>
                  <input
                    type="text"
                    className="form-input"
                    value={formData.weight}
                    onChange={(e) =>
                      setFormData({ ...formData, weight: e.target.value })
                    }
                    required
                  />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Keterangan Subtitle</label>
                <input
                  type="text"
                  className="form-input"
                  value={formData.subtitle}
                  onChange={(e) =>
                    setFormData({ ...formData, subtitle: e.target.value })
                  }
                />
              </div>

              <div className="form-group">
                <label className="form-label">Foto Produk</label>
                <div className="photo-selection-grid">
                  {[
                    { label: 'Udang', src: '/images/products/rengginang-udang.jpg' },
                    { label: 'Bawang', src: '/images/products/rengginang-bawang.jpg' },
                    { label: 'Cumi', src: '/images/products/rengginang-cumi.jpg' },
                    { label: 'Terasi', src: '/images/products/rengginang-terasi.jpg' },
                    { label: 'Hero Mix', src: '/images/products/rengginang-all-hero.jpg' }
                  ].map((preset) => (
                    <button
                      key={preset.src}
                      type="button"
                      className={`photo-preset-btn ${formData.image === preset.src ? 'selected' : ''}`}
                      onClick={() => setFormData({ ...formData, image: preset.src })}
                    >
                      <img src={preset.src} alt={preset.label} />
                      <span>{preset.label}</span>
                    </button>
                  ))}
                </div>
                <div className="custom-photo-upload">
                  <label className="upload-file-label">
                    <ImageIcon size={16} />
                    <span>Atau unggah foto baru dari perangkat</span>
                    <input
                      type="file"
                      accept="image/*"
                      style={{ display: 'none' }}
                      onChange={handleImageFileChange}
                    />
                  </label>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Status Produk</label>
                <div className="status-radio-group">
                  <label className="status-radio-label">
                    <input
                      type="radio"
                      name="status-edit"
                      value="Aktif"
                      checked={formData.status === 'Aktif'}
                      onChange={() => setFormData({ ...formData, status: 'Aktif' })}
                    />
                    <span className="radio-text">
                      <strong>Aktif</strong> (Dapat dilihat dan dibeli oleh pelanggan)
                    </span>
                  </label>
                  <label className="status-radio-label">
                    <input
                      type="radio"
                      name="status-edit"
                      value="Nonaktif"
                      checked={formData.status === 'Nonaktif'}
                      onChange={() => setFormData({ ...formData, status: 'Nonaktif' })}
                    />
                    <span className="radio-text">
                      <strong>Nonaktif</strong> (Disembunyikan dari katalog pelanggan)
                    </span>
                  </label>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Deskripsi Rinci</label>
                <textarea
                  className="form-textarea"
                  rows="3"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                />
              </div>

              <div className="modal-actions">
                <button
                  type="button"
                  className="btn-modal-cancel"
                  onClick={() => setIsEditModalOpen(false)}
                >
                  Batal
                </button>
                <button type="submit" className="btn-modal-save">
                  Simpan Perubahan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ================================================================= */}
      {/* MODAL 3: HAPUS / NONAKTIFKAN (PPL-03 AC 2)                         */}
      {/* ================================================================= */}
      {isDeleteModalOpen && selectedProduct && (
        <div className="modal-backdrop" onClick={() => setIsDeleteModalOpen(false)}>
          <div
            className="modal-box modal-box-delete"
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
          >
            <div className="modal-delete-icon">
              <Trash2 size={28} />
            </div>

            <h2 className="modal-title" style={{ textAlign: 'center', marginBottom: '0.5rem' }}>
              Hapus atau Nonaktifkan Produk?
            </h2>

            <p className="modal-delete-target">
              <strong>{selectedProduct.name}</strong>
            </p>

            <div className="modal-delete-actions">
              <button
                type="button"
                className="btn-action-deactivate"
                onClick={handleDeactivate}
              >
                Nonaktifkan Produk (Aman)
              </button>

              <button
                type="button"
                className="btn-action-hard-delete"
                onClick={handlePermanentDelete}
              >
                Hapus Permanen
              </button>

              <button
                type="button"
                className="btn-modal-cancel"
                style={{ width: '100%', marginTop: '0.5rem' }}
                onClick={() => setIsDeleteModalOpen(false)}
              >
                Batal
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
