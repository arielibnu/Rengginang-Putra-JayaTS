import React, { useState } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import CatalogPlaceholder from './pages/CatalogPlaceholder';
import ContactModal from './components/ContactModal';
import LoginInfoModal from './components/LoginInfoModal';
import './styles/home.css';

export default function App() {
  const [isContactOpen, setIsContactOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);

  return (
    <BrowserRouter>
      <div className="app-layout">
        <Navbar 
          onOpenContact={() => setIsContactOpen(true)} 
          onOpenLogin={() => setIsLoginOpen(true)} 
        />

        <Routes>
          <Route path="/" element={<Home onOpenContact={() => setIsContactOpen(true)} />} />
          <Route path="/katalog" element={<CatalogPlaceholder />} />
          <Route path="*" element={<Home onOpenContact={() => setIsContactOpen(true)} />} />
        </Routes>

        <Footer onOpenContact={() => setIsContactOpen(true)} />

        {/* Modal Interaktif */}
        <ContactModal 
          isOpen={isContactOpen} 
          onClose={() => setIsContactOpen(false)} 
        />
        <LoginInfoModal 
          isOpen={isLoginOpen} 
          onClose={() => setIsLoginOpen(false)} 
        />
      </div>
    </BrowserRouter>
  );
}
