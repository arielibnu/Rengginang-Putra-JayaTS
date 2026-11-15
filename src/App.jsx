
import React, { useState } from 'react';
import {
  BrowserRouter,
  Routes,
  Route,
  useLocation
} from 'react-router-dom';

import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import CatalogPlaceholder from './pages/CatalogPlaceholder';
import ContactModal from './components/ContactModal';
import LoginInfoModal from './components/LoginInfoModal';

// Halaman Login dan Dashboard Owner
import LoginOwner from './pages/LoginOwner';
import OwnerDashboard from './pages/OwnerDashboard';

// PPL-03: Pengelolaan Data Produk dan Varian
import KelolaProduk from './pages/KelolaProduk';

import './styles/home.css';

function AppContent() {
  const [isContactOpen, setIsContactOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);

  const location = useLocation();

  // Login, Dashboard, dan Kelola Produk memiliki layout sendiri
  const isOwnerPage =
    location.pathname === '/login' ||
    location.pathname === '/admin/dashboard' ||
    location.pathname === '/admin/produk';

  if (isOwnerPage) {
    return (
      <Routes>
        <Route path="/login" element={<LoginOwner />} />

        <Route
          path="/admin/dashboard"
          element={<OwnerDashboard />}
        />

        {/* PPL-03: Halaman Kelola Produk dan Varian */}
        <Route
          path="/admin/produk"
          element={<KelolaProduk />}
        />
      </Routes>
    );
  }

  // Tampilan website publik tetap dipertahankan
  return (
    <div className="app-layout">
      <Navbar
        onOpenContact={() => setIsContactOpen(true)}
        onOpenLogin={() => setIsLoginOpen(true)}
      />

      <Routes>
        <Route
          path="/"
          element={
            <Home
              onOpenContact={() => setIsContactOpen(true)}
            />
          }
        />

        <Route
          path="/katalog"
          element={<CatalogPlaceholder />}
        />

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

      {/* Modal asli tetap dipertahankan */}
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
