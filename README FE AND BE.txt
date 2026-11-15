
====================================================================
                     RENGGINANG PUTRA JAYA
                  DOKUMENTASI INTEGRASI FE & BE
                        WBS 1.1.1.3
====================================================================

INFORMASI PROYEK

Nama Proyek    : Website Rengginang Putra Jaya
Sprint         : Sprint 1
Product Backlog: PPL-01 - Login Owner dan Hak Akses Admin
WBS            : 1.1.1.3
Pekerjaan      : Pembatasan Akses Rute Admin dan Integrasi
                 Antarmuka Login Owner ke API Backend
Teknologi      : React, Vite, Golang, MySQL
Status         : Berhasil diuji secara lokal


====================================================================
A. PENJELASAN SINGKAT
====================================================================

Pada WBS 1.1.1.3, halaman Login Owner yang dibuat menggunakan
React dihubungkan dengan sistem autentikasi Backend Golang.

Sebelumnya, login pada frontend menggunakan akun simulasi.
Setelah integrasi, username dan password diperiksa oleh backend
melalui database MySQL.

Dashboard Admin juga dilindungi menggunakan validasi sesi.
Pengguna yang belum login atau sudah logout tidak dapat
mengakses Dashboard Admin.

Fitur yang diimplementasikan:

1. Login Owner melalui API Backend.
2. Validasi sesi login Owner.
3. Pembatasan akses Dashboard Admin.
4. Logout melalui API Backend.
5. Penolakan akses setelah logout.
6. Komunikasi Frontend dan Backend melalui CORS.

Tampilan halaman Login Owner dan Dashboard tetap menggunakan
desain yang dibuat oleh tim Frontend.


====================================================================
B. STRUKTUR FOLDER PROYEK
====================================================================

PPL-01/
|
|-- backend/
|   |-- 1.1.1.2/
|   |   |-- main.go
|   |   |-- autentikasi.go
|   |   |-- database.go
|   |   |-- masuk_keluar.go
|   |   |-- middleware_sesi.go
|   |   |-- pembatas_login.go
|   |   |-- perlindungan_csrf.go
|   |   |-- sesi.go
|   |   `-- README.txt
|   |
|   |-- 1.1.1.3/
|   |   |-- admin_terproteksi.go
|   |   |-- middleware_cors.go
|   |   `-- README.txt
|   |
|   |-- .env
|   |-- .env.example
|   |-- .gitignore
|   |-- go.mod
|   `-- go.sum
|
|-- frontend/
|   |-- public/
|   |-- src/
|   |   |-- pages/
|   |   |   |-- LoginOwner.jsx
|   |   |   `-- OwnerDashboard.jsx
|   |   |-- styles/
|   |   |-- App.jsx
|   |   `-- main.jsx
|   |
|   |-- index.html
|   |-- package.json
|   |-- package-lock.json
|   `-- vite.config.js
|
`-- README FE AND BE.txt


====================================================================
C. PERSIAPAN SEBELUM PENGUJIAN
====================================================================

PERANGKAT LUNAK

1. Visual Studio Code
2. Go / Golang
3. Node.js dan npm
4. XAMPP (MySQL)
5. Google Chrome

PERSIAPAN

1. Buka XAMPP dan aktifkan MySQL.
2. Pastikan database rengginang_putra_jaya tersedia.
3. Pastikan akun Owner sudah dibuat pada database.
4. Buka folder PPL-01 menggunakan Visual Studio Code.
5. Siapkan dua terminal untuk Backend dan Frontend.

CATATAN PENTING

Gunakan dua terminal yang berbeda karena Backend dan
Frontend harus berjalan secara bersamaan.

Backend  : http://localhost:8080
Frontend : http://localhost:5173


====================================================================
D. MENJALANKAN WEBSITE
====================================================================

--------------------------------------------------------------------
LANGKAH 1 - JALANKAN BACKEND GOLANG
--------------------------------------------------------------------

TUJUAN:
Mengaktifkan API autentikasi dan proteksi Admin.

CARA MENJALANKAN:

1. Buka Terminal baru di VS Code.
2. Pastikan berada di folder PPL-01.
3. Jalankan perintah:

[TERMINAL 1 - BACKEND]

cd backend
go run ./1.1.1.2

[AKHIR PERINTAH]

HASIL YANG DIHARAPKAN:

Berhasil terhubung ke MySQL!
Server berjalan di http://localhost:8080

Jangan tutup terminal ini.


--------------------------------------------------------------------
LANGKAH 2 - PERIKSA KONEKSI BACKEND
--------------------------------------------------------------------

Buka browser dan akses:

http://localhost:8080/api/kesehatan

HASIL YANG DIHARAPKAN:

Muncul respons JSON yang menunjukkan status "ok".

ARTINYA:

Backend berhasil menerima permintaan dari browser.


--------------------------------------------------------------------
LANGKAH 3 - JALANKAN FRONTEND REACT
--------------------------------------------------------------------

TUJUAN:
Menampilkan Website Rengginang Putra Jaya.

CARA MENJALANKAN:

1. Buka Terminal baru (Terminal 2).
2. Pastikan berada di folder PPL-01.
3. Jalankan perintah:

[TERMINAL 2 - FRONTEND]

cd frontend
npm run dev -- --port 5173 --strictPort

[AKHIR PERINTAH]

HASIL YANG DIHARAPKAN:

VITE ready
Local: http://localhost:5173/

CATATAN:

Jika dependensi belum terpasang, jalankan:

npm install

Setelah instalasi selesai, jalankan ulang perintah
npm run dev -- --port 5173 --strictPort

Jangan tutup terminal Frontend.

Jika port 5173 sedang digunakan aplikasi lain,
hentikan server Vite lama terlebih dahulu.


--------------------------------------------------------------------
LANGKAH 4 - BUKA WEBSITE
--------------------------------------------------------------------

Buka Google Chrome, kemudian akses:

HALAMAN UTAMA:
http://localhost:5173/

HALAMAN LOGIN OWNER:
http://localhost:5173/login

DASHBOARD ADMIN:
http://localhost:5173/admin/dashboard


====================================================================
E. PENGUJIAN INTEGRASI FRONTEND DAN BACKEND
====================================================================

Untuk melihat hasil komunikasi antara Frontend dan Backend:

1. Buka halaman Login Owner di Chrome.
2. Tekan Ctrl + Shift + I.
3. Pilih tab Network.
4. Pilih filter Fetch/XHR.
5. Pastikan perekaman Network aktif.

Setelah itu, lakukan pengujian berikut secara berurutan.


--------------------------------------------------------------------
PENGUJIAN 1 - LOGIN OWNER
--------------------------------------------------------------------

TUJUAN:
Memastikan Login React terhubung ke Backend Golang.

LANGKAH:

1. Buka http://localhost:5173/login.
2. Isi username Owner yang terdaftar di database.
3. Isi password akun Owner.
4. Klik tombol "Masuk ke dashboard".
5. Periksa tab Network - Fetch/XHR.

Contoh username pengujian:

owner1

Gunakan password akun yang sebelumnya dibuat pada
database MySQL.

HASIL YANG DIHARAPKAN:

Nama Request : masuk
Metode       : POST
Status       : 200 OK

API:
POST /api/pemilik/masuk

SETELAH BERHASIL:

Browser berpindah ke:

http://localhost:5173/admin/dashboard

PENJELASAN:

Frontend mengirimkan username dan password ke Backend.
Backend memvalidasi akun dan membuat sesi login Owner.


--------------------------------------------------------------------
PENGUJIAN 2 - VALIDASI SESI OWNER
--------------------------------------------------------------------

TUJUAN:
Memastikan sesi login diperiksa oleh Backend.

LANGKAH:

1. Setelah login, tetap buka tab Network.
2. Cari request bernama "profil".
3. Periksa status responsnya.

HASIL YANG DIHARAPKAN:

Nama Request : profil
Metode       : GET
Status       : 200 OK

API:
GET /api/pemilik/profil

PENJELASAN:

Backend memastikan Owner mempunyai sesi login valid.
Jika valid, pemeriksaan akses Dashboard dilanjutkan.


--------------------------------------------------------------------
PENGUJIAN 3 - AKSES ADMIN
--------------------------------------------------------------------

TUJUAN:
Memastikan hanya Owner dengan sesi valid yang dapat
mengakses API Admin.

LANGKAH:

1. Setelah login, tetap buka tab Network.
2. Cari request bernama "beranda".
3. Periksa status responsnya.

HASIL YANG DIHARAPKAN:

Nama Request : beranda
Metode       : GET
Status       : 200 OK

API:
GET /api/admin/beranda

PENJELASAN:

Backend memeriksa sesi Owner sebelum mengizinkan
akses ke endpoint Admin.

Jika berhasil, Dashboard Owner ditampilkan.


--------------------------------------------------------------------
PENGUJIAN 4 - LOGOUT OWNER
--------------------------------------------------------------------

TUJUAN:
Memastikan tombol Keluar mengakhiri sesi di Backend.

LANGKAH:

1. Pastikan berada di Dashboard Owner.
2. Tetap buka tab Network - Fetch/XHR.
3. Klik tombol "Keluar" pada sidebar.
4. Periksa request bernama "keluar".

HASIL YANG DIHARAPKAN:

Nama Request : keluar
Metode       : POST
Status       : 200 OK

API:
POST /api/pemilik/keluar

SETELAH BERHASIL:

Browser kembali ke:

http://localhost:5173/login

PENJELASAN:

Frontend mengirim permintaan logout ke Backend.
Backend mengakhiri sesi login Owner.


--------------------------------------------------------------------
PENGUJIAN 5 - AKSES DASHBOARD SETELAH LOGOUT
--------------------------------------------------------------------

TUJUAN:
Memastikan Dashboard tidak dapat diakses setelah logout.

LANGKAH:

1. Login menggunakan akun Owner.
2. Masuk ke Dashboard Admin.
3. Klik tombol "Keluar".
4. Setelah logout, buka alamat:

http://localhost:5173/admin/dashboard

5. Periksa tab Network - Fetch/XHR.

HASIL YANG DIHARAPKAN:

Nama Request : profil
Metode       : GET
Status       : 401 Unauthorized

Browser kembali ke halaman:

http://localhost:5173/login

PENJELASAN:

HTTP 401 menunjukkan sesi login sudah tidak valid.

Backend menolak permintaan akses, kemudian Frontend
mengarahkan pengguna kembali ke halaman Login Owner.

PENTING:

HTTP 401 pada tahap ini adalah hasil yang benar,
bukan kesalahan sistem.


====================================================================
F. PENGUJIAN CORS BACKEND
====================================================================

TUJUAN:

Memastikan Backend mengizinkan permintaan dari
Frontend React pada alamat http://localhost:5173.

CARA PENGUJIAN:

1. Pastikan Backend Golang masih berjalan.
2. Buka Terminal baru.
3. Jalankan perintah berikut dalam CMD atau PowerShell.

[TERMINAL 3 - UJI CORS]

curl.exe -i -X OPTIONS "http://127.0.0.1:8080/api/pemilik/masuk" -H "Origin: http://localhost:5173" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: content-type"

[AKHIR PERINTAH]

HASIL YANG DIHARAPKAN:

HTTP/1.1 204 No Content

Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Credentials: true

PENJELASAN:

Backend menerima permintaan preflight dari Origin
Frontend React yang sudah diizinkan.


====================================================================
G. REKAP HASIL PENGUJIAN
====================================================================

NO. PENGUJIAN                          HASIL
--------------------------------------------------------------------
01. Backend Golang berjalan           BERHASIL
02. Database MySQL terhubung          BERHASIL
03. Frontend React berjalan           BERHASIL
04. Halaman Login Owner tampil        BERHASIL
05. Login melalui API                 HTTP 200
06. Validasi sesi Owner               HTTP 200
07. Akses API Admin                   HTTP 200
08. Logout melalui API                HTTP 200
09. Akses Dashboard setelah logout    HTTP 401
10. Pengalihan ke halaman Login       BERHASIL
11. Pengujian CORS Backend            BERHASIL

KETERANGAN STATUS:

HTTP 200 = Permintaan berhasil diproses.
HTTP 204 = Permintaan berhasil tanpa isi respons.
HTTP 401 = Akses ditolak karena tidak terautentikasi.
HTTP 403 = Akses dilarang oleh server.


====================================================================
H. HASIL IMPLEMENTASI
====================================================================

Berdasarkan pengujian yang telah dilakukan:

1. Halaman Login Owner sudah terhubung dengan API Golang.

2. Username dan password diperiksa melalui autentikasi
   Backend yang menggunakan database MySQL.

3. Dashboard Admin memeriksa sesi Owner melalui API.

4. Endpoint Admin hanya memberikan akses kepada pengguna
   yang memiliki sesi login valid.

5. Tombol Keluar terhubung ke API Logout.

6. Setelah logout, sesi tidak dapat digunakan kembali
   untuk mengakses Dashboard Admin.

7. Tampilan Login Owner, Dashboard, gambar produk,
   dan desain website tetap dipertahankan.

Implementasi dilakukan dengan mengintegrasikan
sistem autentikasi Backend dan antarmuka Frontend.


====================================================================
I. KESIMPULAN
====================================================================

WBS 1.1.1.3 pada Website Rengginang Putra Jaya
berhasil diimplementasikan dan diuji secara lokal.

Frontend React berhasil terhubung dengan Backend
Golang untuk proses login, validasi sesi,
pembatasan akses Admin, dan logout.

Hasil pengujian menunjukkan bahwa permintaan login,
validasi sesi, akses Admin, dan logout memperoleh
status HTTP 200.

Setelah logout, Backend menolak akses dengan HTTP 401,
dan Frontend mengarahkan pengguna ke halaman Login.

Dengan demikian, alur utama integrasi Frontend dan
Backend pada WBS 1.1.1.3 sudah berjalan sesuai
tujuan implementasi.
