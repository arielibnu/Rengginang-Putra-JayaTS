package katalogapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ==========================================
// WBS 1.1.8.2 - API PUBLIC CATALOG PRODUK
// Handler untuk menyajikan endpoint data produk
// dan varian yang aktif kepada calon pembeli.
// ==========================================

// HandlerKatalog mengelola seluruh request publik katalog produk.
type HandlerKatalog struct {
	DB *sql.DB
}

// Menangani routing request API katalog produk.
func (h *HandlerKatalog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Hanya mengizinkan metode GET untuk katalog publik.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsKatalog{
			Pesan: "Metode tidak diizinkan",
		})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/katalog")
	path = strings.Trim(path, "/")

	// Kasus 1: GET /api/katalog atau /api/katalog/produk
	if path == "" || path == "produk" {
		h.ambilSemuaProdukAktif(w, r)
		return
	}

	// Kasus 2: Detail varian atau produk /api/katalog/produk/{id} atau /api/katalog/{id}
	bagian := strings.Split(path, "/")
	var idTarget string

	if len(bagian) == 1 {
		idTarget = bagian[0]
	} else if len(bagian) == 2 && (bagian[0] == "produk" || bagian[0] == "varian") {
		idTarget = bagian[1]
	} else {
		kirimJSON(w, http.StatusNotFound, ResponsKatalog{
			Pesan: "Endpoint katalog tidak ditemukan",
		})
		return
	}

	h.ambilDetailProduk(w, r, idTarget)
}

// ==========================================
// AMBIL SEMUA PRODUK AKTIF
// ==========================================

// Menyajikan seluruh varian produk yang statusnya aktif.
func (h *HandlerKatalog) ambilSemuaProdukAktif(w http.ResponseWriter, r *http.Request) {
	daftar := make([]ItemKatalog, 0)

	// Jika koneksi database tersedia, ambil data dari tabel produk dan varian_produk.
	if h.DB != nil {
		query := `
			SELECT 
				v.id_varian,
				v.id_produk,
				p.nama_produk,
				v.nama_varian,
				v.harga_jual,
				v.berat_gram,
				v.stok,
				COALESCE(p.deskripsi, ''),
				v.status_aktif,
				v.diperbarui_pada
			FROM varian_produk v
			INNER JOIN produk p ON v.id_produk = p.id_produk
			WHERE v.status_aktif = TRUE AND p.status_aktif = TRUE
			ORDER BY v.id_varian ASC
		`

		rows, err := h.DB.Query(query)
		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var item ItemKatalog
				var statusAktif bool
				var diperbaruiPada time.Time

				scanErr := rows.Scan(
					&item.IDVarian,
					&item.IDProduk,
					&item.NamaProduk,
					&item.NamaVarian,
					&item.HargaJual,
					&item.BeratGram,
					&item.Stok,
					&item.Deskripsi,
					&statusAktif,
					&diperbaruiPada,
				)

				if scanErr == nil {
					item.StatusAktif = statusAktif
					item.Diperbarui = diperbaruiPada
					item.KodeVarian = strings.ToLower(item.NamaVarian)
					item.Kategori = "Rengginang tradisional"
					item.HargaFormat = FormatRupiah(item.HargaJual)
					item.BeratFormat = fmt.Sprintf("Berat bersih %d gr", item.BeratGram)
					item.StatusStok = TentukanStatusKetersediaan(item.Stok)

					gambarNormal, gambarDetail := CariGambarVarian(item.NamaVarian)
					item.Gambar = gambarNormal
					item.GambarDetail = gambarDetail

					daftar = append(daftar, item)
				}
			}
		}
	}

	// Jika database belum terisi atau tabel kosong, gunakan 4 varian resmi default.
	if len(daftar) == 0 {
		daftar = AmbilVarianDefault()
	}

	kirimJSON(w, http.StatusOK, ResponsKatalog{
		Pesan: "Data katalog produk aktif berhasil disajikan",
		Total: len(daftar),
		Data:  daftar,
	})
}

// ==========================================
// AMBIL DETAIL SATU PRODUK
// ==========================================

// Menyajikan data detail satu varian berdasarkan ID angka atau nama varian.
func (h *HandlerKatalog) ambilDetailProduk(w http.ResponseWriter, r *http.Request, target string) {
	targetBersih := strings.ToLower(strings.TrimSpace(target))

	// Cari di database terlebih dahulu jika DB aktif.
	if h.DB != nil {
		idVarian, errParse := strconv.ParseUint(targetBersih, 10, 64)

		query := `
			SELECT 
				v.id_varian,
				v.id_produk,
				p.nama_produk,
				v.nama_varian,
				v.harga_jual,
				v.berat_gram,
				v.stok,
				COALESCE(p.deskripsi, ''),
				v.status_aktif,
				v.diperbarui_pada
			FROM varian_produk v
			INNER JOIN produk p ON v.id_produk = p.id_produk
			WHERE (v.id_varian = ? OR LOWER(v.nama_varian) = ?)
			  AND v.status_aktif = TRUE AND p.status_aktif = TRUE
			LIMIT 1
		`

		var item ItemKatalog
		var statusAktif bool
		var diperbaruiPada time.Time

		err := h.DB.QueryRow(query, idVarian, targetBersih).Scan(
			&item.IDVarian,
			&item.IDProduk,
			&item.NamaProduk,
			&item.NamaVarian,
			&item.HargaJual,
			&item.BeratGram,
			&item.Stok,
			&item.Deskripsi,
			&statusAktif,
			&diperbaruiPada,
		)

		if err == nil && (errParse == nil || idVarian > 0 || item.IDVarian > 0) {
			item.StatusAktif = statusAktif
			item.Diperbarui = diperbaruiPada
			item.KodeVarian = strings.ToLower(item.NamaVarian)
			item.Kategori = "Rengginang tradisional"
			item.HargaFormat = FormatRupiah(item.HargaJual)
			item.BeratFormat = fmt.Sprintf("Berat bersih %d gr", item.BeratGram)
			item.StatusStok = TentukanStatusKetersediaan(item.Stok)

			gambarNormal, gambarDetail := CariGambarVarian(item.NamaVarian)
			item.Gambar = gambarNormal
			item.GambarDetail = gambarDetail

			kirimJSON(w, http.StatusOK, ResponsKatalog{
				Pesan: "Detail produk berhasil ditemukan",
				Total: 1,
				Data:  item,
			})
			return
		}
	}

	// Fallback ke daftar varian default jika tidak ditemukan di DB.
	daftarDefault := AmbilVarianDefault()
	for _, v := range daftarDefault {
		if strings.ToLower(v.KodeVarian) == targetBersih ||
			strings.ToLower(v.NamaVarian) == targetBersih ||
			strconv.FormatUint(v.IDVarian, 10) == targetBersih {

			kirimJSON(w, http.StatusOK, ResponsKatalog{
				Pesan: "Detail produk berhasil ditemukan",
				Total: 1,
				Data:  v,
			})
			return
		}
	}

	kirimJSON(w, http.StatusNotFound, ResponsKatalog{
		Pesan: "Varian produk tidak ditemukan atau sedang tidak aktif",
	})
}

// ==========================================
// FUNGSI BANTUAN LOGIKA KATALOG
// ==========================================

// TentukanStatusKetersediaan menghitung label ketersediaan
// berdasarkan sisa stok riil varian.
func TentukanStatusKetersediaan(stok uint64) string {
	if stok == 0 {
		return "Habis"
	}
	if stok <= 5 {
		return "Stok Menipis"
	}
	return "Tersedia • contoh"
}

// CariGambarVarian mencocokkan path file gambar horizontal
// dan vertikal berdasarkan nama varian produk.
func CariGambarVarian(namaVarian string) (string, string) {
	v := strings.ToLower(namaVarian)

	if strings.Contains(v, "bawang") {
		return "/images/products/rengginang-bawang.jpg", "/images/products/rengginang-bawang-vertical.jpg"
	}
	if strings.Contains(v, "udang") {
		return "/images/products/rengginang-udang.jpg", "/images/products/rengginang-udang-vertical.jpg"
	}
	if strings.Contains(v, "cumi") {
		return "/images/products/rengginang-cumi.jpg", "/images/products/rengginang-cumi-vertical.jpg"
	}
	if strings.Contains(v, "terasi") {
		return "/images/products/rengginang-terasi.jpg", "/images/products/rengginang-terasi-vertical.jpg"
	}

	return "/images/products/rengginang-all-hero.jpg", "/images/products/rengginang-all-hero.jpg"
}

// FormatRupiah memformat angka ke dalam teks Rupiah standar.
func FormatRupiah(angka float64) string {
	nilai := int64(angka)
	str := strconv.FormatInt(nilai, 10)
	n := len(str)
	if n <= 3 {
		return "Rp" + str
	}

	var hasil []string
	sisa := n % 3
	if sisa > 0 {
		hasil = append(hasil, str[:sisa])
	}
	for i := sisa; i < n; i += 3 {
		hasil = append(hasil, str[i:i+3])
	}
	return "Rp" + strings.Join(hasil, ".")
}

// AmbilVarianDefault menyajikan 4 varian rengginang resmi UMKM
// sesuai desain mockup Figma PPL-08.
func AmbilVarianDefault() []ItemKatalog {
	now := time.Now()

	return []ItemKatalog{
		{
			IDVarian:     1,
			IDProduk:     1,
			KodeVarian:   "udang",
			NamaProduk:   "Rengginang Putra Jaya",
			NamaVarian:   "Rengginang Udang",
			HargaJual:    28000,
			HargaFormat:  "Rp28.000",
			BeratGram:    450,
			BeratFormat:  "Berat bersih 450 gr",
			Stok:         34,
			StatusStok:   "Tersedia • contoh",
			Deskripsi:    "Varian rasa udang dalam kemasan 450 gram.",
			Kategori:     "Rengginang tradisional",
			Gambar:       "/images/products/rengginang-udang.jpg",
			GambarDetail: "/images/products/rengginang-udang-vertical.jpg",
			StatusAktif:  true,
			Diperbarui:   now,
		},
		{
			IDVarian:     2,
			IDProduk:     1,
			KodeVarian:   "bawang",
			NamaProduk:   "Rengginang Putra Jaya",
			NamaVarian:   "Rengginang Bawang",
			HargaJual:    28000,
			HargaFormat:  "Rp28.000",
			BeratGram:    450,
			BeratFormat:  "Berat bersih 450 gr",
			Stok:         40,
			StatusStok:   "Tersedia • contoh",
			Deskripsi:    "Varian rasa bawang dalam kemasan 450 gram.",
			Kategori:     "Rengginang tradisional",
			Gambar:       "/images/products/rengginang-bawang.jpg",
			GambarDetail: "/images/products/rengginang-bawang-vertical.jpg",
			StatusAktif:  true,
			Diperbarui:   now,
		},
		{
			IDVarian:     3,
			IDProduk:     1,
			KodeVarian:   "cumi",
			NamaProduk:   "Rengginang Putra Jaya",
			NamaVarian:   "Rengginang Cumi",
			HargaJual:    28000,
			HargaFormat:  "Rp28.000",
			BeratGram:    450,
			BeratFormat:  "Berat bersih 450 gr",
			Stok:         14,
			StatusStok:   "Tersedia • contoh",
			Deskripsi:    "Varian rasa cumi dalam kemasan 450 gram.",
			Kategori:     "Rengginang tradisional",
			Gambar:       "/images/products/rengginang-cumi.jpg",
			GambarDetail: "/images/products/rengginang-cumi-vertical.jpg",
			StatusAktif:  true,
			Diperbarui:   now,
		},
		{
			IDVarian:     4,
			IDProduk:     1,
			KodeVarian:   "terasi",
			NamaProduk:   "Rengginang Putra Jaya",
			NamaVarian:   "Rengginang Terasi",
			HargaJual:    28000,
			HargaFormat:  "Rp28.000",
			BeratGram:    450,
			BeratFormat:  "Berat bersih 450 gr",
			Stok:         40,
			StatusStok:   "Tersedia • contoh",
			Deskripsi:    "Varian rasa terasi dalam kemasan 450 gram.",
			Kategori:     "Rengginang tradisional",
			Gambar:       "/images/products/rengginang-terasi.jpg",
			GambarDetail: "/images/products/rengginang-terasi-vertical.jpg",
			StatusAktif:  true,
			Diperbarui:   now,
		},
	}
}

// Mengirim respons JSON ke client.
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
