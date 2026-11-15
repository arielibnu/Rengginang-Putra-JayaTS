
PANDUAN PENGUJIAN BACKEND PROTEKSI AKSES ADMIN
WEBSITE RENGGINANG PUTRA JAYA

------------------------------------------------------------

Proyek          : Website Rengginang Putra Jaya
Product Backlog : PPL-01 - Login Owner dan Hak Akses Admin
WBS             : 1.1.1.3
Sprint          : Sprint 1
Bagian          : Backend (BE)
Hasil Pekerjaan : API Admin Terproteksi dan Middleware CORS

------------------------------------------------------------

DESKRIPSI TUGAS

Mengimplementasikan pembatasan akses API Admin agar hanya
Owner yang memiliki sesi login valid yang dapat mengaksesnya.

Backend juga menyediakan middleware CORS sebagai persiapan
integrasi halaman login frontend dengan API.

Pengujian dilakukan melalui terminal Visual Studio Code.

CATATAN

Dokumentasi ini khusus pengujian backend WBS 1.1.1.3.
Integrasi halaman login dan proteksi navigasi di browser
masih memerlukan kode frontend.


============================================================
1. MENJALANKAN BACKEND
============================================================

TUJUAN

Memastikan aplikasi backend dapat berjalan setelah
penambahan proteksi Admin dan middleware CORS.

LANGKAH PENGUJIAN

1. Aktifkan MySQL melalui XAMPP.
2. Buka folder backend di Visual Studio Code.
3. Buka Terminal, kemudian New Terminal.
4. Pada Terminal 1, jalankan perintah berikut.

------------------ KODE TERMINAL 1 -------------------------

go build ./1.1.1.2
go run ./1.1.1.2

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

Berhasil terhubung ke MySQL!
Server berjalan di http://localhost:8080

PENJELASAN

Backend berhasil dikompilasi dan terhubung ke database.
Server menerima permintaan API melalui port 8080.

PENTING

Biarkan Terminal 1 tetap aktif selama pengujian.
Jika muncul error port 8080 sudah digunakan, hentikan
server lama terlebih dahulu.


============================================================
2. PERSIAPAN TERMINAL PENGUJIAN
============================================================

TUJUAN

Menyiapkan PowerShell untuk menjalankan pengujian API.

LANGKAH PENGUJIAN

1. Buka Terminal 2 melalui Terminal, kemudian New Terminal.
2. Jika Terminal 2 masih menggunakan CMD, ketik:

---------------- KODE MASUK POWERSHELL ---------------------

powershell.exe

-------------------- AKHIR KODE ----------------------------

3. Tekan Enter.
4. Pastikan terminal menampilkan awalan PS.
5. Jalankan seluruh pengujian berikutnya pada Terminal 2.

PENTING

Tidak perlu menjalankan powershell.exe kembali apabila
terminal sudah berada dalam mode PowerShell.


============================================================
3. PENGUJIAN API ADMIN TANPA LOGIN
============================================================

TUJUAN

Membuktikan bahwa API Admin tidak dapat diakses
oleh pengguna yang belum login.

LANGKAH PENGUJIAN

Jalankan kode berikut di Terminal 2.

------------------ KODE TERMINAL 2 -------------------------

try {
    Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/admin/beranda" -Method GET -ErrorAction Stop
} catch {
    Write-Host "HTTP Status:" ([int]$_.Exception.Response.StatusCode)
}

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

HTTP Status: 401

PENJELASAN

HTTP 401 menunjukkan middleware sesi berhasil menolak
pengguna yang tidak memiliki sesi login valid.


============================================================
4. PENGUJIAN LOGIN OWNER
============================================================

TUJUAN

Membuat sesi login Owner untuk menguji akses API Admin.

LANGKAH PENGUJIAN

Jalankan seluruh kode berikut pada Terminal 2.

------------------ KODE TERMINAL 2 -------------------------

$kredensial = Get-Credential -UserName "owner1" -Message "Pengujian Akses Admin"

$data = @{
    nama_pengguna = $kredensial.UserName
    kata_sandi = $kredensial.GetNetworkCredential().Password
} | ConvertTo-Json

$sesi = New-Object Microsoft.PowerShell.Commands.WebRequestSession

$header = @{
    Origin = "http://localhost:5173"
}

Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/pemilik/masuk" -Method POST -ContentType "application/json" -Headers $header -Body $data -WebSession $sesi

-------------------- AKHIR KODE ----------------------------

Masukkan kata sandi akun owner1 pada jendela yang muncul.

HASIL YANG DIHARAPKAN

Login pemilik berhasil

PENJELASAN

Backend memverifikasi akun Owner dan membuat sesi login.
Sesi tersebut digunakan untuk mengakses API Admin.


============================================================
5. PENGUJIAN AKSES ADMIN SETELAH LOGIN
============================================================

TUJUAN

Memastikan Owner yang memiliki sesi valid dapat
mengakses API Admin.

LANGKAH PENGUJIAN

Gunakan Terminal 2 yang sama dengan pengujian login.

------------------ KODE TERMINAL 2 -------------------------

Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/admin/beranda" -Method GET -Headers $header -WebSession $sesi | ConvertTo-Json -Depth 5

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

"status": "berhasil"
"pesan": "Akses admin diizinkan"
"halaman": "Beranda Admin"
"akses": "Owner"

PENJELASAN

Middleware memeriksa cookie sesi yang dikirim Owner.
Apabila sesi valid, sistem memberikan respons API Admin.

Pengujian ini membuktikan bahwa hanya Owner dengan
sesi valid yang diberikan akses ke endpoint Admin.


============================================================
6. PENGUJIAN CORS DENGAN ORIGIN DIIZINKAN
============================================================

TUJUAN

Memastikan backend menerima permintaan preflight CORS
dari alamat frontend yang telah dikonfigurasi.

LANGKAH PENGUJIAN

Jalankan kode berikut pada Terminal 2.

------------------ KODE TERMINAL 2 -------------------------

curl.exe -i -X OPTIONS "http://127.0.0.1:8080/api/pemilik/masuk" -H "Origin: http://localhost:5173" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: content-type"

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

HTTP/1.1 204 No Content
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type

PENJELASAN

HTTP 204 menunjukkan backend berhasil menangani
permintaan preflight CORS.

Header yang ditampilkan menunjukkan bahwa origin
frontend yang dikonfigurasi mendapat izin melakukan
permintaan dengan cookie sesi.


============================================================
7. PENGUJIAN CORS DENGAN ORIGIN TIDAK DIIZINKAN
============================================================

TUJUAN

Memastikan middleware CORS menolak Origin yang tidak
terdaftar dalam konfigurasi backend.

LANGKAH PENGUJIAN

Jalankan kode berikut pada Terminal 2.

------------------ KODE TERMINAL 2 -------------------------

curl.exe -i -X OPTIONS "http://127.0.0.1:8080/api/pemilik/masuk" -H "Origin: http://website-tidak-diizinkan.test" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: content-type"

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

HTTP/1.1 403 Forbidden

Asal permintaan tidak diizinkan

PENJELASAN

Backend menolak Origin yang berbeda dari alamat frontend
yang telah diizinkan melalui konfigurasi CORS.


============================================================
8. PENGUJIAN LOGOUT OWNER
============================================================

TUJUAN

Memastikan Owner dapat mengakhiri sesi login.

LANGKAH PENGUJIAN

Gunakan Terminal 2 yang masih menyimpan variabel sesi.

------------------ KODE TERMINAL 2 -------------------------

Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/pemilik/keluar" -Method POST -Headers $header -WebSession $sesi

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

Logout pemilik berhasil

PENJELASAN

Backend menghapus sesi login Owner dari database.


============================================================
9. PENGUJIAN AKSES ADMIN SETELAH LOGOUT
============================================================

TUJUAN

Membuktikan bahwa Owner yang telah logout tidak dapat
mengakses API Admin menggunakan sesi sebelumnya.

LANGKAH PENGUJIAN

Jalankan kode berikut pada Terminal 2.

------------------ KODE TERMINAL 2 -------------------------

try {
    Invoke-RestMethod -Uri "http://127.0.0.1:8080/api/admin/beranda" -Method GET -WebSession $sesi -ErrorAction Stop
} catch {
    Write-Host "HTTP Status Setelah Logout:" ([int]$_.Exception.Response.StatusCode)
}

-------------------- AKHIR KODE ----------------------------

HASIL YANG DIHARAPKAN

HTTP Status Setelah Logout: 401

PENJELASAN

Middleware menolak akses karena sesi Owner telah
dihapus melalui proses logout.


============================================================
10. RINGKASAN PENGUJIAN
============================================================

Pengujian 1
Menjalankan backend
Hasil: Server aktif dan MySQL terhubung

Pengujian 2
Persiapan PowerShell
Hasil: Terminal siap digunakan

Pengujian 3
Akses Admin tanpa login
Hasil: HTTP 401

Pengujian 4
Login Owner
Hasil: Berhasil

Pengujian 5
Akses Admin setelah login
Hasil: Akses diizinkan

Pengujian 6
CORS Origin diizinkan
Hasil: HTTP 204

Pengujian 7
CORS Origin tidak diizinkan
Hasil: HTTP 403

Pengujian 8
Logout Owner
Hasil: Berhasil

Pengujian 9
Akses Admin setelah logout
Hasil: HTTP 401


============================================================
11. KESIMPULAN
============================================================

Backend WBS 1.1.1.3 pada Website Rengginang Putra Jaya
berhasil menerapkan pembatasan akses API Admin
menggunakan middleware sesi Owner.

Pengujian membuktikan bahwa akses Admin ditolak ketika
pengguna belum login, diizinkan ketika sesi valid, dan
kembali ditolak setelah logout.

Middleware CORS juga telah diuji untuk menangani
permintaan dari Origin yang diizinkan serta menolak
Origin yang tidak diizinkan.

Dengan demikian, bagian Backend WBS 1.1.1.3 telah
diimplementasikan dan diuji secara lokal.

Integrasi antarmuka login Owner dengan API backend,
pengiriman cookie dari browser, dan proteksi navigasi
halaman Admin masih menunggu implementasi Frontend.
