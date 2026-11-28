WBS 1.4.5.3 (BE): Pemberian indikator disclaimer jika data biaya belum lengkap dan audit anti-duplikasi hitungan.
Fitur:
- Pemeriksaan kelengkapan data pengeluaran operasional dan HPP.
- Peringatan disclaimer transparan apabila salah satu komponen biaya belum tercatat.
- Audit anti-duplikasi: Menjamin setiap transaksi dihitung tepat 1 kali menggunakan nomor_transaksi unik.
- Pengecualian otomatis terhadap transaksi online Midtrans Sandbox dan ongkos kirim titipan ekspedisi.
- Endpoint: GET /api/admin/laporan/estimasi-keuntungan/audit?bulan=...&tahun=...
