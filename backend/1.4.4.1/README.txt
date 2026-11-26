WBS 1.4.4.1 (BE): Pengelompokan biaya pengeluaran operasional berdasarkan tanggal, kategori beban, dan periode.
Fitur: PPL-28 - Rekap Pengeluaran Usaha Bulanan

Kategori Beban:
1. OPERASIONAL: Beban operasional harian/rutin (listrik, air, gas, sewa, upah kerja).
2. BELANJA_PERSEDIAAN: Belanja bahan baku mentah (ketan, udang, bumbu, kemasan) - dipisahkan agar anti-duplikasi modal.
3. BIAYA_LAIN: Beban di luar operasional dan persediaan (pemeliharaan mesin, transportasi, tak terduga).

Endpoint:
GET /api/admin/laporan/pengeluaran?bulan={bulan}&tahun={tahun}
