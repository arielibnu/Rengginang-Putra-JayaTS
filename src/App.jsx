
import React, { useState } from 'react';
import {
  BrowserRouter,
  Routes,
  Route,
  useLocation
} from 'react-router-dom';

// Komponen website publik
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import CatalogPlaceholder from './pages/CatalogPlaceholder';
import ContactModal from './components/ContactModal';
import LoginInfoModal from './components/LoginInfoModal';

// PPL-01: Login Owner dan Dashboard
import LoginOwner from './pages/LoginOwner';
import OwnerDashboard from './pages/OwnerDashboard';

// PPL-03: Pengelolaan Produk dan Varian
import KelolaProduk from './pages/KelolaProduk';

// PPL-04: Pengelolaan Stok dan Riwayat Persediaan
import DataStok from './pages/DataStok';
import CatatMutasi from './pages/CatatMutasi';
import CatatStokKeluar from './pages/CatatStokKeluar';
import MutasiTersimpan from './pages/MutasiTersimpan';
import StokKeluarError from './pages/StokKeluarError';
import PenyesuaianStok from './pages/PenyesuaianStok';
import RiwayatStok from './pages/RiwayatStok';

// PPL-02: Registrasi dan Login Pelanggan
import LoginPelanggan from './pages/LoginPelanggan';
import DaftarPelanggan from './pages/DaftarPelanggan';
import DaftarEmailError from './pages/DaftarEmailError';
import AkunPelanggan from './pages/AkunPelanggan';

// CSS website dan halaman pelanggan
import './styles/home.css';
import './styles/auth-pelanggan.css';

function AppContent() {
  const [isContactOpen, setIsContactOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);

  const location = useLocation();

  // Area Owner tetap menggunakan layout sendiri.
  const isOwnerPage =
    location.pathname === '/login' ||
    location.pathname === '/login-owner' ||
    location.pathname === '/admin' ||
    location.pathname.startsWith('/admin/');

  if (isOwnerPage) {
    return (
      <Routes>
        {/* Login Owner PPL-01 */}
        <Route
          path="/login"
          element={<LoginOwner />}
        />

        <Route
          path="/login-owner"
          element={<LoginOwner />}
        />

        {/* Dashboard Owner */}
        <Route
          path="/admin"
          element={<OwnerDashboard />}
        />

        <Route
          path="/admin/dashboard"
          element={<OwnerDashboard />}
        />

        {/* Kelola Produk PPL-03 */}
        <Route
          path="/admin/produk"
          element={<KelolaProduk />}
        />

        {/* Data Stok PPL-04 */}
        <Route
          path="/admin/stok"
          element={<DataStok />}
        />
        <Route
          path="/admin/persediaan"
          element={<DataStok />}
        />
        <Route
          path="/admin/data-stok"
          element={<DataStok />}
        />
        <Route
          path="/admin/stok/mutasi"
          element={<CatatMutasi />}
        />
        <Route
          path="/admin/stok/catat"
          element={<CatatMutasi />}
        />
        <Route
          path="/admin/stok/keluar"
          element={<CatatStokKeluar />}
        />
        <Route
          path="/admin/stok/catat-keluar"
          element={<CatatStokKeluar />}
        />
        <Route
          path="/admin/stok/tersimpan"
          element={<MutasiTersimpan />}
        />
        <Route
          path="/admin/stok/sukses"
          element={<MutasiTersimpan />}
        />
        <Route
          path="/admin/stok/error"
          element={<StokKeluarError />}
        />
        <Route
          path="/admin/stok/penyesuaian"
          element={<PenyesuaianStok />}
        />
        <Route
          path="/admin/stok/riwayat"
          element={<RiwayatStok />}
        />
      </Routes>
    );
  }

  // Website publik dan pelanggan
  return (
    <div className="app-layout">
      <Navbar
        onOpenContact={() => setIsContactOpen(true)}
        onOpenLogin={() => setIsLoginOpen(true)}
      />

      <Routes>
        {/* Beranda */}
        <Route
          path="/"
          element={
            <Home
              onOpenContact={() => setIsContactOpen(true)}
            />
          }
        />

        {/* Katalog dapat diakses tanpa login */}
        <Route
          path="/katalog"
          element={<CatalogPlaceholder />}
        />

        {/* Login pelanggan PPL-02 */}
        <Route
          path="/masuk"
          element={<LoginPelanggan />}
        />

        {/* Registrasi pelanggan PPL-02 */}
        <Route
          path="/daftar"
          element={<DaftarPelanggan />}
        />

        {/* Email sudah terdaftar */}
        <Route
          path="/daftar-error"
          element={<DaftarEmailError />}
        />

        {/* Akun pelanggan */}
        <Route
          path="/akun"
          element={<AkunPelanggan />}
        />

        <Route
          path="/akun-saya"
          element={<AkunPelanggan />}
        />

        {/* Alamat tidak dikenal kembali ke beranda */}
        <Route
          path="*"
          element={
            <Home
              onOpenContact={() => setIsContactOpen(true)}
            />
          }
        />
      </Routes>

      <Footer
        onOpenContact={() => setIsContactOpen(true)}
      />

      <ContactModal
        isOpen={isContactOpen}
        onClose={() => setIsContactOpen(false)}
      />

      <LoginInfoModal
        isOpen={isLoginOpen}
        onClose={() => setIsLoginOpen(false)}
      />
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
}