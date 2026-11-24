/**
 * Konfigurasi Informasi Usaha UMKM RENGGINANG PUTRA JAYA
 * Sesuai Desain Mockup Resmi Beranda
 */

export const BUSINESS_CONFIG = {
  name: "Rengginang Putra Jaya",
  producer: "Hj. Arifah",
  location: "Situbondo, Jawa Timur",
  taglineText: "CAMILAN KHAS SITUBONDO",
  heroHeading: "Rengginang tradisional, rasa istimewa.",
  heroDescription: "Nikmati empat pilihan rasa rengginang khas Situbondo, dibuat dengan resep tradisional.",
  
  // Nomor WhatsApp (ubah saat nomor resmi tersedia)
  whatsappNumber: null,
  whatsappDefaultMessage: "Halo Rengginang Putra Jaya, saya ingin menanyakan informasi pemesanan produk rengginang.",

  netWeight: "Berat bersih 450 gr",
  priceFormatted: "Rp28.000",
  priceNumber: 28000,

  routes: {
    home: "/",
    catalog: "/katalog",
    about: "/#tentang",
    contact: "/#kontak",
    login: "/login"
  },

  features: [
    {
      number: "1",
      title: "4 varian rasa",
      desc: "Pilih rasa favoritmu."
    },
    {
      number: "2",
      title: "Informasi jelas",
      desc: "Harga, berat, dan stok."
    },
    {
      number: "3",
      title: "Belanja lebih mudah",
      desc: "Tampilan nyaman di HP."
    }
  ],

  flavors: [
    {
      id: "udang",
      name: "Rengginang Udang",
      image: "/images/products/rengginang-udang.jpg",
      weight: "Berat bersih 450 gr",
      price: "Rp28.000"
    },
    {
      id: "bawang",
      name: "Rengginang Bawang",
      image: "/images/products/rengginang-bawang.jpg",
      weight: "Berat bersih 450 gr",
      price: "Rp28.000"
    },
    {
      id: "cumi",
      name: "Rengginang Cumi",
      image: "/images/products/rengginang-cumi.jpg",
      weight: "Berat bersih 450 gr",
      price: "Rp28.000"
    },
    {
      id: "terasi",
      name: "Rengginang Terasi",
      image: "/images/products/rengginang-terasi.jpg",
      weight: "Berat bersih 450 gr",
      price: "Rp28.000"
    }
  ]
};
