WBS 1.4.2.1 - Perhitungan Sisa Stok Akhir Berdasarkan Riwayat Mutasi Barang (PPL-26)
=====================================================================================

Modul ini bertanggung jawab menghitung sisa stok fisik, reservasi aktif, stok tersedia, dan nilai persediaan barang pada akhir periode penutupan (bulanan/tahunan) dari riwayat mutasi resmi.

Fitur Utama:
1. Agregasi mutasi barang (MASUK & KELUAR) dari tabel riwayat_mutasi_stok.
2. Pemisahan ketat transaksi riil vs sandbox ('mode_stok = RIIL', mengabaikan SIMULASI).
3. Perhitungan stok tersedia akurat: Stok Fisik - Reservasi Aktif.
4. Perhitungan nilai aset persediaan barang (Stok Fisik * Harga Jual).
5. Endpoint REST API: GET /api/admin/laporan/stok-akhir?bulan=November&tahun=2026.
