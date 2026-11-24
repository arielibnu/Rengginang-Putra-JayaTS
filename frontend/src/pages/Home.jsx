import React from 'react';
import Hero from '../components/Hero';
import FlavorPreview from '../components/FlavorPreview';
import WhyChooseUs from '../components/WhyChooseUs';

export default function Home({ onOpenContact }) {
  return (
    <main>
      {/* 1. Hero Section */}
      <Hero onOpenContact={onOpenContact} />

      {/* 2. Varian Rengginang */}
      <FlavorPreview />

      {/* 3. Mengapa Pilih Putra Jaya? */}
      <WhyChooseUs />
    </main>
  );
}
