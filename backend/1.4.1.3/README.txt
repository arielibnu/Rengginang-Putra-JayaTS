WBS 1.4.1.3 (BE): Pemisahan transaksi penjualan riil dari simulasi sandbox dan pengecualian pesanan yang dibatalkan
Fitur: PPL-25 - Rekap Penjualan Bulanan

Tujuan:
1. Menjamin kas masuk nyata toko fisik (offline) tidak dicampur aduk dengan simulasi Midtrans Sandbox.
2. Memastikan pesanan yang berstatus 'batal', 'dibatalkan', atau 'kedaluwarsa' dikecualikan dari total barang terjual dan omzet riil.

Endpoint:
GET /api/admin/laporan/penjualan/audit-pemisahan?bulan={bulan}&tahun={tahun}
