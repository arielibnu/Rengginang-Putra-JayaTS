================================================================================
DOKUMENTASI MODUL BACKEND WBS 1.2.3.3
PPL-12: Penyamaan Format Payload Data Checkout Bagi Pembeli Guest dan Pelanggan yang Login
Deliverable: Standarisasi Schema Order Request
================================================================================

1. DESKRIPSI PEKERJAAN
   Modul ini mengimplementasikan logika backend untuk WBS 1.2.3.3 (BE):
   - Standarisasi schema data order request (payload checkout) yang seragam baik bagi
     pembeli tamu tanpa login (Guest) maupun pelanggan yang telah memiliki akun (Logged-in Customer).
   - Penanganan field identifikasi pemesan (`is_guest`, `tipe_checkout`, `id_pelanggan`):
     * Guest     : is_guest = true,  tipe_checkout = "guest",     id_pelanggan = null
     * Pelanggan : is_guest = false, tipe_checkout = "pelanggan", id_pelanggan = uint64
   - Validasi struktur kontak penerima (nama penerima, nomor WhatsApp yang valid).
   - Validasi item pesanan dan ketersediaan stok aktual ke database MySQL (tabel `varian_produk`).
   - Kalkulasi ulang subtotal dan ringkasan pembayaran secara aman di sisi server.
   - Pembuatan nomor pesanan unik terstandarisasi (`PJ-GUEST-YYYYMMDD-XXXX` atau `PJ-CUST-YYYYMMDD-XXXX`).
   - Penyimpanan data pesanan terstandarisasi ke database MySQL (tabel `pesanan_checkout`).
   - Sinkronisasi endpoint order request dengan antarmuka Frontend React pada halaman Checkout.

2. STRUKTUR FILE
   - model.go                  : Definisi struct data PayloadCheckoutRequest, InformasiKontak,
                                 ItemPesananCheckout, RingkasanBiayaCheckout, HasilCheckoutResponse,
                                 dan ResponsAPI standar beserta Unmarshaler kustom.
   - payload_checkout.go       : Logika bisnis standarisasi payload, validasi kontak dan item,
                                 verifikasi stok database MySQL, penyimpanan pesanan, generator kode order,
                                 serta HTTP Handlers:
                                 * HandlerProsesCheckout (POST /api/checkout & POST /api/checkout/proses)
                                 * HandlerSkemaPayload   (GET  /api/checkout/skema)
   - payload_checkout_test.go  : Unit test komprehensif menguji skenario order guest, order pelanggan,
                                 validasi input (kontak kosong, WA invalid, kuantitas 0, over-stock, out-of-stock),
                                 skema dokumentasi, dan pembatasan metode HTTP.
   - schema_checkout.sql       : DDL tabel `pesanan_checkout` dan `detail_pesanan_checkout` MySQL.
   - README.txt                : Dokumentasi teknis lengkap modul WBS 1.2.3.3.

3. SKEMA STANDAR ORDER REQUEST (PAYLOAD CHECKOUT)

   Endpoint: POST /api/checkout
   Content-Type: application/json

   a. Format Payload Pembeli Tamu (Guest):
      {
        "tipe_checkout": "guest",
        "is_guest": true,
        "id_pelanggan": null,
        "kontak": {
          "nama_penerima": "Budi Santoso",
          "nomor_whatsapp": "081234567890",
          "email": "budi@example.com",
          "catatan": "Tolong packing rapi"
        },
        "item_pesanan": [
          {
            "id_varian": 1,
            "slug": "udang",
            "nama_varian": "Udang",
            "nama_produk": "Rengginang Udang",
            "harga_satuan": 28000,
            "jumlah": 2,
            "subtotal": 56000
          }
        ],
        "ringkasan": {
          "total_item": 2,
          "subtotal": 56000,
          "ongkos_kirim": 0,
          "total_pembayaran": 56000
        }
      }

   b. Format Payload Pelanggan Terdaftar (Logged-in):
      {
        "tipe_checkout": "pelanggan",
        "is_guest": false,
        "id_pelanggan": 42,
        "kontak": {
          "nama_penerima": "Siti Nurhaliza",
          "nomor_whatsapp": "085678912345",
          "email": "siti@example.com",
          "catatan": ""
        },
        "item_pesanan": [
          {
            "id_varian": 2,
            "slug": "bawang",
            "nama_varian": "Bawang",
            "nama_produk": "Rengginang Bawang",
            "harga_satuan": 28000,
            "jumlah": 1,
            "subtotal": 28000
          }
        ],
        "ringkasan": {
          "total_item": 1,
          "subtotal": 28000,
          "ongkos_kirim": 0,
          "total_pembayaran": 28000
        }
      }

4. ENDPOINT API

   a. POST /api/checkout atau POST /api/checkout/proses
      - Deskripsi : Menerima order request standar baik dari guest maupun pelanggan terdaftar.
      - Respon Sukses (HTTP 200 OK):
        {
          "status": "sukses",
          "pesan": "Standarisasi order request valid dan pesanan checkout berhasil diproses",
          "data": {
            "nomor_pesanan": "PJ-GUEST-20261119-4821",
            "tipe_checkout": "guest",
            "is_guest": true,
            "id_pelanggan": null,
            "kontak": { ... },
            "item_pesanan": [ ... ],
            "ringkasan": {
              "total_item": 2,
              "subtotal": 56000,
              "ongkos_kirim": 0,
              "total_pembayaran": 56000
            },
            "status_pesanan": "MENUNGGU_KONFIRMASI",
            "pesan_notifikasi": "Pesanan tamu (guest) berhasil diproses tanpa akun",
            "dibuat_pada": "2026-11-19T14:30:00Z"
          }
        }

      - Respon Ditolak (HTTP 400 Bad Request):
        {
          "status": "ditolak",
          "pesan": "Nama penerima wajib diisi (minimal 2 karakter)"
        }

   b. GET /api/checkout/skema
      - Deskripsi : Menyajikan dokumentasi resmi definisi skema standar order request bagi FE dan BE.
      - Respon HTTP 200 OK.

5. CARA MENJALANKAN & MENGUJI
   a. Menjalankan Unit Test WBS 1.2.3.3:
      cd backend
      go test -v ./1.2.3.3

   b. Menjalankan Server Backend:
      cd backend
      go run ./1.1.1.2
