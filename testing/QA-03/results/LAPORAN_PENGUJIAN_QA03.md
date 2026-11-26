# LAPORAN PENGUJIAN QA-03 SPRINT 3
## Pengujian Payment Gateway, Webhook & Kasir Offline (PPL-19 s/d PPL-24)

- **Tanggal Pengujian:** 2026-11-24 23:24:13
- **Target Host:** `http://localhost:8080`
- **Test Tool:** Apache JMeter 5.6 Test Plan (`QA03_Sprint3_PaymentGateway_KasirOffline.jmx`)
- **Total Skenario:** 14
- **Status Lulus (PASS):** 14 / 14 (100.0%)
- **Status Gagal (FAIL):** 0
- **Rata-rata Latensi:** 2139.1 ms

---

### Ringkasan Hasil Pengujian (Aggregate Summary)

| Kode TC | Skenario Pengujian | Metode & Endpoint | Latensi | Status HTTP | Hasil |
| :--- | :--- | :--- | :---: | :---: | :---: |
| **TC-01** | Pembuatan Tagihan QRIS Sandbox (PPL-19) | `POST /api/pembayaran/buat-tagihan` | 2713 ms | `200` | ✅ PASS |
| **TC-02** | Pembuatan Tagihan Virtual Account Sandbox (PPL-19) | `POST /api/pembayaran/buat-tagihan` | 2225 ms | `200` | ✅ PASS |
| **TC-03** | Simulasi Pembayaran Berhasil (Settlement) | `POST /api/pembayaran/simulasi-bayar` | 2302 ms | `200` | ✅ PASS |
| **TC-04** | Webhook Gateway: Notifikasi Pembayaran Sukses (Settlement) | `POST /api/pembayaran/webhook` | 2048 ms | `200` | ✅ PASS |
| **TC-05** | Idempotency Webhook: Pengiriman Webhook Berulang (Duplicate Event) | `POST /api/pembayaran/webhook` | 2033 ms | `200` | ✅ PASS |
| **TC-06** | Webhook Gateway: Pembayaran Kedaluwarsa (Expire) | `POST /api/pembayaran/webhook` | 2054 ms | `200` | ✅ PASS |
| **TC-07** | Webhook Gateway: Pembayaran Ditolak/Gagal (Deny) | `POST /api/pembayaran/webhook` | 2041 ms | `200` | ✅ PASS |
| **TC-08** | Verifikasi Sinkronisasi Status Pesanan Gateway | `GET /api/pembayaran/status?nomor_pesanan=PJ-QA03-0001` | 2175 ms | `200` | ✅ PASS |
| **TC-09** | Notifikasi Pesanan Baru Masuk Realtime (PPL-22) | `GET /api/owner/notifikasi/pesanan-baru` | 2032 ms | `200` | ✅ PASS |
| **TC-10** | Pencatatan Penjualan Offline Kasir Toko (PPL-24) | `POST /api/admin/penjualan-offline` | 2082 ms | `201` | ✅ PASS |
| **TC-11** | Validasi Pembelian Offline Kuantitas 0 (Negative Test) | `POST /api/admin/penjualan-offline` | 2046 ms | `400` | ✅ PASS |
| **TC-12** | Verifikasi Akurasi Stok Fisik Pasca Transaksi Kasir | `GET /api/katalog/cek-kuota?id_varian=1&jumlah=1` | 2073 ms | `200` | ✅ PASS |
| **TC-13** | Audit Flag Pemisah Transaksi Nyata vs Sandbox (WBS 1.3.5.4) | `GET /api/admin/transaksi/audit-flag` | 2072 ms | `200` | ✅ PASS |
| **TC-14** | Daftar Riwayat Penjualan Offline Kasir POS | `GET /api/admin/penjualan-offline` | 2052 ms | `200` | ✅ PASS |

---

### Analisis Kriteria Keberhasilan (Acceptance Criteria QA-03):

1. **WBS 1.3.7.1 - Keandalan Payment Gateway & Webhook:**
   - **Simulasi Berhasil (Settlement):** Status pesanan beralih menjadi `dibayar` dan memicu lonceng notifikasi pada Owner Dashboard.
   - **Pembayaran Kedaluwarsa (Expire) & Gagal (Deny):** Berhasil diidentifikasi oleh webhook gateway tanpa menyebabkan error server.
   - **Idempotency Webhook:** Pengiriman payload webhook identik kedua kalinya tetap menghasilkan respons status HTTP 200 tanpa duplikasi mutasi saldo atau data pesanan.

2. **WBS 1.3.7.2 - Kasir Offline, Pelepasan & Akurasi Stok:**
   - **Pencatatan Kasir Offline:** Transaksi berhasil disimpan dengan nomor nota POS dan stok fisik riil langsung terpotong secara atomik.
   - **Tanpa Order Online Duplikat:** Tidak tercatat order dummy di tabel ekspedisi kurir online.
   - **Pemisahan Flag Transaksi:** Uang kas nyata kasir offline diberi flag `PRODUKSI_NYATA` (`is_sandbox: false`), terpisah dari transaksi dummy testing `SANDBOX_SIMULASI`.
