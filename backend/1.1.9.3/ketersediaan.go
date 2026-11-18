package ketersediaanproduk

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// HandlerDetail menangani endpoint ketersediaan dan rincian produk (WBS 1.1.9.3).
type HandlerDetail struct {
	DB *sql.DB
}

// TentukanStatusKetersediaan menguji apakah stok masih tersedia atau habis.
// Sesuai aturan bisnis Figma PPL-09:
// - Stok > 0  => "TERSEDIA", BisaDibeli: true
// - Stok <= 0 => "STOK HABIS", BisaDibeli: false
func TentukanStatusKetersediaan(stok int64) (string, bool, string, string) {
	if stok <= 0 {
		return "STOK HABIS", false, "Produk sedang habis", "Pembelian dinonaktifkan. Silakan pilih varian lain."
	}
	return "TERSEDIA", true, "Stok tersedia", "Jumlah stok mengikuti data persediaan dari Owner."
}

// ValidasiInteraksiBeli memvalidasi apakah tombol beli boleh ditekan atau harus dinonaktifkan (disabled).
func ValidasiInteraksiBeli(stokTersedia int64, jumlahBeli uint32) (bool, string, error) {
	if stokTersedia <= 0 {
		return false, "Pembelian dinonaktifkan karena produk sedang habis (out-of-stock)", errors.New("stok habis")
	}

	if jumlahBeli == 0 {
		return false, "Jumlah pembelian minimal 1 bungkus", errors.New("jumlah beli tidak valid")
	}

	if int64(jumlahBeli) > stokTersedia {
		pesan := fmt.Sprintf("Jumlah pesanan (%d bungkus) melebihi stok yang tersedia (%d bungkus)", jumlahBeli, stokTersedia)
		return false, pesan, errors.New("melebihi stok")
	}

	return true, "Produk siap dibeli", nil
}

// FormatRupiah memformat nominal angka menjadi format mata uang Rupiah standar.
func FormatRupiah(angka uint64) string {
	str := strconv.FormatUint(angka, 10)
	var hasil []byte
	pj := len(str)

	for i := 0; i < pj; i++ {
		if i > 0 && (pj-i)%3 == 0 {
			hasil = append(hasil, '.')
		}
		hasil = append(hasil, str[i])
	}
	return "Rp" + string(hasil)
}

// HandlerDetailProduk melayani permintaan GET rincian produk spesifik.
// Endpoint: GET /api/katalog/detail?id=cumi atau GET /api/katalog/detail?id=udang
func (h *HandlerDetail) HandlerDetailProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimResponJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan",
		})
		return
	}

	idParam := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("id")))
	if idParam == "" {
		idParam = "cumi" // Default sesuai Figma mockup
	}

	detail, err := h.ambilDetailProduk(idParam)
	if err != nil {
		kirimResponJSON(w, http.StatusNotFound, ResponsAPI{
			Status: "gagal",
			Pesan:  "Varian produk tidak ditemukan",
		})
		return
	}

	kirimResponJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Berhasil mengambil detail spesifikasi produk",
		Data:   detail,
	})
}

// HandlerValidasiBeli memproses validasi pembelian untuk memastikan tombol disable tidak dibypass.
// Endpoint: POST /api/katalog/validasi-beli
func (h *HandlerDetail) HandlerValidasiBeli(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimResponJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan",
		})
		return
	}

	var req PermintaanValidasiBeli
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimResponJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  "Format payload JSON tidak valid",
		})
		return
	}

	detail, err := h.ambilDetailProduk(req.Slug)
	if err != nil && req.IDVarian > 0 {
		detail, err = h.ambilDetailProdukByID(req.IDVarian)
	}

	if err != nil {
		kirimResponJSON(w, http.StatusNotFound, ResponsAPI{
			Status: "gagal",
			Pesan:  "Varian produk tidak ditemukan untuk divalidasi",
		})
		return
	}

	bisaDibeli, pesan, errVal := ValidasiInteraksiBeli(detail.StokTersedia, req.JumlahBeli)

	hasil := HasilValidasiBeli{
		IDVarian:      detail.IDVarian,
		Slug:          detail.Slug,
		BisaDibeli:    bisaDibeli,
		StatusStok:    detail.StatusKetersediaan,
		SisaStok:      detail.StokTersedia,
		PesanValidasi: pesan,
	}

	if errVal != nil {
		kirimResponJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "ditolak",
			Pesan:  pesan,
			Data:   hasil,
		})
		return
	}

	kirimResponJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  pesan,
		Data:   hasil,
	})
}

// ambilDetailProduk mencari rincian varian berdasarkan slug rasa atau basis data.
func (h *HandlerDetail) ambilDetailProduk(slug string) (*DetailProdukResponse, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))

	// Data default sesuai mockup resmi Figma PPL-09 (Cumi = 0 bungkus, Udang = 12 bungkus)
	dataDefault := map[string]DetailProdukResponse{
		"cumi": {
			IDVarian:     3,
			Slug:         "cumi",
			NamaProduk:   "Rengginang Cumi",
			NamaVarian:   "Cumi",
			HargaJual:    28000,
			BeratGram:    450,
			Deskripsi:    "Rengginang tradisional • kemasan 450 gram",
			Kategori:     "Rengginang tradisional",
			StokTersedia: 0, // STOK HABIS sesuai Gambar 1 Figma
			GambarUtama:  "/images/products/rengginang-cumi-vertical.jpg",
			GaleriFoto: []FotoGaleri{
				{Label: "Kemasan", URL: "/images/products/rengginang-cumi-vertical.jpg"},
				{Label: "Label", URL: "/images/products/rengginang-cumi.jpg"},
				{Label: "Detail isi", URL: "/images/products/rengginang-all-hero.jpg"},
			},
		},
		"udang": {
			IDVarian:     1,
			Slug:         "udang",
			NamaProduk:   "Rengginang Udang",
			NamaVarian:   "Udang",
			HargaJual:    28000,
			BeratGram:    450,
			Deskripsi:    "Rengginang tradisional • kemasan 450 gram",
			Kategori:     "Rengginang tradisional",
			StokTersedia: 12, // TERSEDIA sesuai Gambar 2 Figma
			GambarUtama:  "/images/products/rengginang-udang-vertical.jpg",
			GaleriFoto: []FotoGaleri{
				{Label: "Kemasan", URL: "/images/products/rengginang-udang-vertical.jpg"},
				{Label: "Label", URL: "/images/products/rengginang-udang.jpg"},
				{Label: "Detail isi", URL: "/images/products/rengginang-all-hero.jpg"},
			},
		},
		"bawang": {
			IDVarian:     2,
			Slug:         "bawang",
			NamaProduk:   "Rengginang Bawang",
			NamaVarian:   "Bawang",
			HargaJual:    28000,
			BeratGram:    450,
			Deskripsi:    "Rengginang tradisional • kemasan 450 gram",
			Kategori:     "Rengginang tradisional",
			StokTersedia: 38,
			GambarUtama:  "/images/products/rengginang-bawang-vertical.jpg",
			GaleriFoto: []FotoGaleri{
				{Label: "Kemasan", URL: "/images/products/rengginang-bawang-vertical.jpg"},
				{Label: "Label", URL: "/images/products/rengginang-bawang.jpg"},
				{Label: "Detail isi", URL: "/images/products/rengginang-all-hero.jpg"},
			},
		},
		"terasi": {
			IDVarian:     4,
			Slug:         "terasi",
			NamaProduk:   "Rengginang Terasi",
			NamaVarian:   "Terasi",
			HargaJual:    28000,
			BeratGram:    450,
			Deskripsi:    "Rengginang tradisional • kemasan 450 gram",
			Kategori:     "Rengginang tradisional",
			StokTersedia: 38,
			GambarUtama:  "/images/products/rengginang-terasi-vertical.jpg",
			GaleriFoto: []FotoGaleri{
				{Label: "Kemasan", URL: "/images/products/rengginang-terasi-vertical.jpg"},
				{Label: "Label", URL: "/images/products/rengginang-terasi.jpg"},
				{Label: "Detail isi", URL: "/images/products/rengginang-all-hero.jpg"},
			},
		},
	}

	item, ada := dataDefault[slug]
	if !ada {
		return nil, errors.New("varian tidak ditemukan")
	}

	// Jika koneksi MySQL aktif, sinkronkan data stok riil dari basis data
	if h.DB != nil {
		var stokDB uint64
		var hargaDB float64
		err := h.DB.QueryRow(`
			SELECT stok, harga_jual 
			FROM varian_produk 
			WHERE LOWER(nama_varian) LIKE ? AND status_aktif = TRUE
			LIMIT 1
		`, "%"+slug+"%").Scan(&stokDB, &hargaDB)
		if err == nil {
			item.StokTersedia = int64(stokDB)
			if hargaDB > 0 {
				item.HargaJual = uint64(hargaDB)
			}
		}
	}

	item.FormatHarga = FormatRupiah(item.HargaJual)
	item.FormatBerat = fmt.Sprintf("%d gram", item.BeratGram)
	item.FormatStok = fmt.Sprintf("%d bungkus", item.StokTersedia)
	item.StatusKetersediaan, item.BisaDibeli, item.PesanStatus, item.PesanKeterangan = TentukanStatusKetersediaan(item.StokTersedia)

	return &item, nil
}

func (h *HandlerDetail) ambilDetailProdukByID(idVarian uint64) (*DetailProdukResponse, error) {
	petaSlug := map[uint64]string{
		1: "udang",
		2: "bawang",
		3: "cumi",
		4: "terasi",
	}
	slug, ada := petaSlug[idVarian]
	if !ada {
		return nil, errors.New("id varian tidak terdaftar")
	}
	return h.ambilDetailProduk(slug)
}

func kirimResponJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
