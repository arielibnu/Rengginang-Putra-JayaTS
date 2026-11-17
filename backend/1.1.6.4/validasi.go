package validasibiaya

import (
	"errors"
	"strings"
)

// ==========================================
// WBS 1.1.6.4
// Aturan Anti Double-Counting Biaya
//
// Memastikan modal barang terjual (HPP) dan
// pengeluaran operasional tidak terhitung ganda.
// ==========================================

// Daftar kategori yang termasuk komponen modal/HPP.
// Kategori ini tidak boleh dimasukkan sebagai
// pengeluaran operasional karena sudah dihitung
// di dalam harga modal per varian.
var kategoriModalHPP = []string{
	"bahan baku",
	"bahan mentah",
	"ketan",
	"beras ketan",
	"udang",
	"cumi",
	"terasi",
	"bawang",
	"bawang putih",
	"bawang merah",
	"bumbu produksi",
	"rempah",
	"tinta cumi",
	"garam produksi",
}

// Daftar kategori pengeluaran operasional yang sah.
// Kategori ini bukan bagian dari HPP sehingga aman
// dicatat sebagai pengeluaran usaha terpisah.
var kategoriOperasionalSah = []string{
	"listrik",
	"gas",
	"listrik & gas",
	"air",
	"sewa tempat",
	"gaji karyawan",
	"transportasi",
	"kemasan",
	"plastik kemasan",
	"label kemasan",
	"minyak goreng",
	"perawatan alat",
	"kebersihan",
	"pemasaran",
	"iklan",
	"lain-lain",
}

// ==========================================
// VALIDASI KATEGORI PENGELUARAN
// ==========================================

// CekKategoriPengeluaran memvalidasi apakah kategori
// yang dimasukkan bukan bagian dari komponen modal/HPP.
// Kalau kategorinya termasuk HPP, pengeluaran ditolak
// supaya biaya tidak dihitung dua kali.
func CekKategoriPengeluaran(kategori string) error {

	kategori = strings.ToLower(strings.TrimSpace(kategori))

	if kategori == "" {
		return errors.New("Kategori pengeluaran wajib diisi")
	}

	// Periksa apakah kategori ini termasuk komponen HPP.
	for _, modal := range kategoriModalHPP {
		if strings.Contains(kategori, modal) {
			return errors.New(
				"Kategori \"" + kategori + "\" termasuk komponen modal/HPP. " +
					"Biaya bahan baku sudah dihitung dalam harga modal per varian. " +
					"Jangan masukkan ulang sebagai pengeluaran operasional",
			)
		}
	}

	return nil
}

// ==========================================
// VALIDASI DUPLIKASI PENGELUARAN
// ==========================================

// DataPengeluaranCek berisi informasi pengeluaran
// yang sudah tercatat untuk pengecekan duplikasi.
type DataPengeluaranCek struct {
	Tanggal  string
	Kategori string
	Nominal  float64
}

// CekDuplikasiPengeluaran memeriksa apakah pengeluaran
// baru kemungkinan duplikat dengan yang sudah tercatat.
// Dianggap duplikat jika tanggal, kategori, dan nominal
// semuanya sama persis.
func CekDuplikasiPengeluaran(
	tanggal string,
	kategori string,
	nominal float64,
	daftarSudahAda []DataPengeluaranCek,
) error {

	kategori = strings.ToLower(strings.TrimSpace(kategori))
	tanggal = strings.TrimSpace(tanggal)

	for _, ada := range daftarSudahAda {
		tanggalAda := strings.TrimSpace(ada.Tanggal)
		kategoriAda := strings.ToLower(strings.TrimSpace(ada.Kategori))

		if tanggalAda == tanggal &&
			kategoriAda == kategori &&
			ada.Nominal == nominal {

			return errors.New(
				"Pengeluaran dengan tanggal, kategori, dan nominal yang sama " +
					"sudah pernah dicatat. Pastikan ini bukan pencatatan ganda",
			)
		}
	}

	return nil
}

// ==========================================
// CEK APAKAH KATEGORI OPERASIONAL SAH
// ==========================================

// ApakahKategoriOperasionalSah mengecek apakah kategori
// termasuk dalam daftar kategori operasional yang dikenal.
// Mengembalikan true kalau dikenal, false kalau tidak.
// Kategori yang tidak dikenal tetap diizinkan tapi
// bisa diberi peringatan di frontend.
func ApakahKategoriOperasionalSah(kategori string) bool {

	kategori = strings.ToLower(strings.TrimSpace(kategori))

	for _, sah := range kategoriOperasionalSah {
		if strings.Contains(kategori, sah) {
			return true
		}
	}

	return false
}

// ==========================================
// VALIDASI NOMINAL
// ==========================================

// CekNominalPengeluaran memvalidasi bahwa nominal
// pengeluaran masuk akal dan tidak negatif.
func CekNominalPengeluaran(nominal float64) error {

	if nominal <= 0 {
		return errors.New("Nominal pengeluaran harus lebih dari 0")
	}

	// Batas maksimal nominal Rp 1 miliar per transaksi.
	// Angka di atas itu kemungkinan salah ketik.
	if nominal > 1000000000 {
		return errors.New(
			"Nominal melebihi batas maksimal Rp 1.000.000.000 per transaksi",
		)
	}

	return nil
}

// ==========================================
// VALIDASI LENGKAP PENGELUARAN
// ==========================================

// HasilValidasi menyimpan hasil pengecekan lengkap.
type HasilValidasi struct {
	Valid     bool   `json:"valid"`
	Pesan     string `json:"pesan"`
	Peringatan string `json:"peringatan,omitempty"`
}

// ValidasiPengeluaranLengkap menjalankan semua aturan
// validasi sekaligus untuk satu pengeluaran baru.
func ValidasiPengeluaranLengkap(
	tanggal string,
	kategori string,
	nominal float64,
	daftarSudahAda []DataPengeluaranCek,
) HasilValidasi {

	// 1. Cek nominal dulu.
	if err := CekNominalPengeluaran(nominal); err != nil {
		return HasilValidasi{
			Valid: false,
			Pesan: err.Error(),
		}
	}

	// 2. Cek apakah kategori termasuk komponen HPP.
	if err := CekKategoriPengeluaran(kategori); err != nil {
		return HasilValidasi{
			Valid: false,
			Pesan: err.Error(),
		}
	}

	// 3. Cek duplikasi pengeluaran.
	if err := CekDuplikasiPengeluaran(
		tanggal, kategori, nominal, daftarSudahAda,
	); err != nil {
		return HasilValidasi{
			Valid: false,
			Pesan: err.Error(),
		}
	}

	// 4. Kasih peringatan kalau kategori tidak dikenal.
	peringatan := ""
	if !ApakahKategoriOperasionalSah(kategori) {
		peringatan = "Kategori \"" + kategori +
			"\" belum ada di daftar standar. " +
			"Pastikan ini bukan komponen modal/HPP"
	}

	return HasilValidasi{
		Valid:      true,
		Pesan:      "Pengeluaran valid untuk dicatat",
		Peringatan: peringatan,
	}
}
