WBS 1.4.5.1 (BE): Perhitungan laba bersih: (Penjualan Bersih - HPP Terjual - Biaya Operasional), mengecualikan ongkir titipan.
Fitur:
- Penjualan bersih dihitung dari transaksi valid non-sandbox (kasir offline / toko).
- Ongkos kirim titipan ekspedisi dieksklusikan dari pendapatan usaha.
- Modal barang terjual (HPP) dan Biaya operasional dikurangkan untuk menghitung laba bersih.
- Endpoint: GET /api/admin/laporan/estimasi-keuntungan?bulan=...&tahun=...
