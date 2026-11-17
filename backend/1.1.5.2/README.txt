====================================================================
                     RENGGINANG PUTRA JAYA
             DOKUMENTASI BACKEND WBS 1.1.5.2
====================================================================

INFORMASI MODUL

Nama Proyek    : Website Rengginang Putra Jaya
Product Backlog: PPL-05 — Peringatan Stok Menipis Otomatis
WBS ID         : 1.1.5.2
Pekerjaan      : Penyediaan pengaturan batas minimum stok per varian
                 dan penerapan aturan peringatan stok kritis
Deliverable    : API & Form Batas Ambang Stok
Teknologi      : Go (Golang), MySQL
Status         : Berhasil Diuji & Siap Digunakan


====================================================================
A. PENJELASAN SINGKAT
====================================================================

Modul WBS 1.1.5.2 menyediakan backend untuk mengelola batas minimum
persediaan per varian rengginang serta mengevaluasi status stok kritis
secara otomatis berdasarkan aturan bisnis:

  "Peringatan muncul saat stok tersedia kurang dari atau sama
   dengan batas minimum per varian (stok_tersedia <= batas_minimum)."

Jika kondisi kritis terpenuhi:
- Status varian disetel menjadi "Stok menipis" (apakah_kritis = true).
- Backend menghasilkan pesan peringatan utama dan rincian varian
  untuk ditampilkan pada Dashboard Owner dan Halaman Data Stok.

Jika stok di atas batas minimum:
- Status varian disetel menjadi "Normal" (apakah_kritis = false).


====================================================================
B. DAFTAR ENDPOINT API
====================================================================

1. Evaluasi Peringatan Stok Otomatis
   - Endpoint : GET /api/admin/stok/peringatan
   - Akses    : Terproteksi (Sesi Owner & Hak Akses Admin)
   - Deskripsi: Mengembalikan rekapitulasi status stok seluruh varian,
                jumlah varian menipis, dan pesan peringatan otomatis.
   - Contoh Respons:
     {
       "status": "sukses",
       "pesan": "Berhasil memuat peringatan stok",
       "data": {
         "total_varian_aktif": 4,
         "jumlah_varian_menipis": 1,
         "ada_peringatan": true,
         "pesan_utama": "Perhatian: 1 varian stok menipis",
         "keterangan_rinci": "Rengginang Cumi tersedia 12 bungkus, sedangkan batas minimum 15 bungkus.",
         "daftar_varian": [
           {
             "id_varian": 1,
             "nama_varian": "Rengginang Udang",
             "stok_fisik": 34,
             "reservasi": 2,
             "stok_tersedia": 32,
             "batas_minimum": 15,
             "status": "Normal",
             "apakah_kritis": false
           },
           {
             "id_varian": 3,
             "nama_varian": "Rengginang Cumi",
             "stok_fisik": 14,
             "reservasi": 2,
             "stok_tersedia": 12,
             "batas_minimum": 15,
             "status": "Stok menipis",
             "apakah_kritis": true
           }
         ]
       }
     }

2. Pengambilan Konfigurasi Batas Minimum
   - Endpoint : GET /api/admin/stok/batas-minimum
   - Query    : ?id_varian=3 (opsional, untuk 1 varian)
   - Akses    : Terproteksi (Sesi Owner)

3. Pengaturan / Pembaruan Batas Minimum Varian
   - Endpoint : POST / PUT /api/admin/stok/batas-minimum
   - Akses    : Terproteksi (Sesi Owner & Proteksi CSRF)
   - Body     : JSON
     {
       "id_varian": 3,
       "batas_minimum": 15,
       "catatan": "Penyesuaian batas stok minimum Owner"
     }
   - Respons  : HTTP 200 OK beserta data status stok terkini varian.


====================================================================
C. CARA MENJALANKAN UNIT TEST
====================================================================

1. Masuk ke folder backend:
   cd backend

2. Jalankan perintah test untuk paket 1.1.5.2:
   go test -v ./1.1.5.2

Hasil: Seluruh pengujian (Evaluasi status stok, validasi batas minimum,
dan perakitan ringkasan peringatan otomatis) PASS 100%.
