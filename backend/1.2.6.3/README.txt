================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.6.3
PPL-15: Pengisian Wilayah dan Alamat Pengiriman
Tugas: Validasi Kelengkapan Wilayah dan Alamat Tujuan Sesuai Parameter Kebutuhan API Kurir
Deliverable: Validasi Data Geografis Pengiriman
Role: BE, FE
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan logika backend untuk WBS 1.2.6.3 (BE):
   - Validasi kelengkapan wilayah dan alamat tujuan pengiriman agar memenuhi
     parameter standar yang dipersyaratkan oleh penyedia API Kurir ekspedisi
     (seperti RajaOngkir, Biteship, Komship, JNE, POS Indonesia, SiCepat, dll).
   - Validasi hierarki geografis bertingkat (Provinsi -> Kota/Kabupaten -> Kecamatan).
   - Validasi kelayakan kode pos (harus 5 digit angka Indonesia).
   - Validasi alamat lengkap detail (nama jalan, nomor rumah, RT/RW, patokan)
     dengan batas minimal karakter untuk mencegah kegagalan pengiriman kurir.
   - Validasi kontak penerima (nama penerima dan nomor WhatsApp berformat Indonesia).
   - Penanganan metode penerimaan pesanan:
     * Jika "kurir": seluruh parameter geografis dan kontak wajib lengkap.
     * Jika "ambil_sendiri": alamat kurir tidak diwajibkan (status "NOT_REQUIRED").
   - Standarisasi data geografis asal toko (Rengginang Putra Jaya Situbondo)
     dan data geografis tujuan untuk persiapan modul ongkos kirim (PPL-16).
   - Penyimpanan data alamat pengiriman ke database MySQL (tabel `alamat_pengiriman_pesanan`)
     dengan mode fallback in-memory yang tangguh bila database offline.
   - Penyediaan endpoint terintegrasi dengan Frontend React Checkout tanpa mengubah
     tampilan atau mockup antarmuka pengguna.

2. STRUKTUR FILE
   - model.go                 : Definisi struct data InformasiKontak, WilayahPengiriman,
                                PermintaanWilayahPengiriman, LokasiGeografis,
                                ParameterKurirStandar, HasilValidasiWilayah, dan ResponsAPI.
   - validasi_wilayah.go      : Logika validasi geografis, master data wilayah bertingkat,
                                normalisasi nomor telepon E.164, generator parameter kurir,
                                fungsi persistensi database, serta HTTP Handlers:
                                * HandlerSimpanDanValidasiWilayah (POST /api/checkout/wilayah-pengiriman)
                                * HandlerValidasiAlamatKurir       (POST /api/checkout/validasi-alamat)
                                * HandlerMasterWilayah            (GET  /api/checkout/master-wilayah)
                                * HandlerParameterKurir           (GET  /api/checkout/parameter-kurir)
   - validasi_wilayah_test.go : Unit test komprehensif menguji:
                                * Skenario validasi lengkap dan valid
                                * Skenario metode ambil sendiri
                                * Skenario error data tidak lengkap (provinsi kosong,
                                  kota kosong, kecamatan kosong, alamat pendek, kode pos salah,
                                  inkonsistensi hierarki wilayah, kontak tidak valid)
                                * Pengujian endpoint API dan pembatasan metode HTTP.
   - schema_wilayah.sql       : DDL tabel MySQL `alamat_pengiriman_pesanan` dan `master_wilayah_kurir`.
   - README.txt               : Dokumentasi teknis lengkap modul WBS 1.2.6.3.

3. SPESIFIKASI PARAMETER KEBUTUHAN API KURIR
   Untuk melakukan kalkulasi ongkos kirim atau penjemputan barang (pickup booking),
   API kurir mewajibkan data berikut:
   a. Asal Pengiriman (Origin Store):
      - Provinsi      : Jawa Timur
      - Kota/Kabupaten: Kabupaten Situbondo
      - Kecamatan     : Kecamatan Situbondo
      - Kode Pos      : 68312
   b. Tujuan Pengiriman (Destination Recipient):
      - Provinsi      : Nama provinsi tujuan (wajib)
      - Kota/Kabupaten: Nama kabupaten/kota tujuan (wajib, hierarkis sesuai provinsi)
      - Kecamatan     : Nama kecamatan tujuan (wajib, hierarkis sesuai kota)
      - Kode Pos      : 5 digit angka numerik (wajib untuk akurasi tarif zona pengiriman)
      - Alamat Lengkap: Minimal 10 karakter detail jalan, no. rumah, RT/RW, kelurahan
   c. Karakteristik Paket (Package Properties):
      - Berat         : Dalam gram (default 1350 gram untuk 2 Udang + 1 Bawang)
      - Kategori      : Makanan Ringan
      - Is Fragile    : True (rengginang renyah mudah hancur, membutuhkan penanganan hati-hati)

4. DOKUMENTASI ENDPOINT API
   
   a. POST /api/checkout/wilayah-pengiriman
      Deskripsi: Endpoint utama integrasi Frontend Checkout PPL-15.
      Request Header: Content-Type: application/json
      Request Body:
      {
        "kontak": {
          "nama_penerima": "Contoh Pelanggan",
          "nomor_whatsapp": "081234567890"
        },
        "metode_penerimaan": "kurir",
        "wilayah_pengiriman": {
          "provinsi": "Jawa Timur",
          "kota_kabupaten": "Kabupaten Sidoarjo",
          "kecamatan": "Kecamatan Sidoarjo",
          "kode_pos": "61211",
          "alamat_lengkap": "Jl. Contoh No. 10, RT 02 / RW 03, Kelurahan Contoh"
        },
        "subtotal_produk": 84000,
        "total_sementara": 84000
      }

      Response Sukses (200 OK):
      {
        "status": "sukses",
        "pesan": "Data wilayah dan alamat pengiriman valid dan siap dikirim ke API ongkos kirim (PPL-16).",
        "data": {
          "status": "sukses",
          "pesan": "Data wilayah dan alamat pengiriman valid dan siap dikirim ke API ongkos kirim (PPL-16).",
          "valid": true,
          "metode_penerimaan": "kurir",
          "wilayah": {
            "provinsi": "Jawa Timur",
            "kota_kabupaten": "Kabupaten Sidoarjo",
            "kecamatan": "Kecamatan Sidoarjo",
            "kode_pos": "61211",
            "alamat_lengkap": "Jl. Contoh No. 10, RT 02 / RW 03, Kelurahan Contoh"
          },
          "parameter_kurir": {
            "asal_pengiriman": {
              "provinsi": "Jawa Timur",
              "kota_kabupaten": "Kabupaten Situbondo",
              "kecamatan": "Kecamatan Situbondo",
              "kode_pos": "68312",
              "alamat_lengkap": "Jl. Raya Situbondo No. 45, Sentra UMKM Rengginang, Situbondo"
            },
            "tujuan_pengiriman": {
              "provinsi": "Jawa Timur",
              "kota_kabupaten": "Kabupaten Sidoarjo",
              "kecamatan": "Kecamatan Sidoarjo",
              "kode_pos": "61211",
              "alamat_lengkap": "Jl. Contoh No. 10, RT 02 / RW 03, Kelurahan Contoh"
            },
            "penerima": {
              "nama": "Contoh Pelanggan",
              "telepon": "081234567890",
              "telepon_standar": "+6281234567890"
            },
            "paket": {
              "berat_gram": 1350,
              "kategori": "Makanan Ringan",
              "is_fragile": true
            },
            "is_kurir_ready": true,
            "status_geografis": "VALID"
          }
        }
      }

      Response Error (400 Bad Request):
      {
        "status": "error",
        "pesan": "Data wilayah dan alamat tujuan belum lengkap atau tidak valid sesuai parameter API kurir",
        "errors": {
          "kota_kabupaten": "Kota/Kabupaten tujuan wajib dipilih",
          "alamat_lengkap": "Alamat lengkap tujuan pengiriman wajib diisi"
        }
      }

   b. GET /api/checkout/master-wilayah
      Mengembalikan data hierarki Provinsi -> Kota -> Kecamatan Indonesia.

   c. GET /api/checkout/parameter-kurir
      Mengembalikan dokumentasi spesifikasi parameter kebutuhan API kurir.

5. HASIL PENGUJIAN UNIT TEST
   Pengujian otomatis dapat dijalankan dengan perintah:
   $ cd backend/1.2.6.3
   $ go test -v .

   Hasil seluruh skenario (15 skenario validasi, boundary cases, dan HTTP compliance):
   PASS (100% Berhasil)
================================================================================
