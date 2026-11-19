package kuotastok

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ============================================================================
// WBS 1.2.1.3 - HANDLER VALIDASI KUOTA STOK REALTIME
// Mengelola pengecekan kuota pembelian terhadap sisa stok aktual dari database.
// ============================================================================

// HandlerValidasiKuota mengelola validasi stok realtime dari MySQL.
type HandlerValidasiKuota struct {
	DB *sql.DB
}

// DataDefaultVarian fallback data resmi mockup Figma jika database belum dikonfigurasi.
var DataDefaultVarian = map[string]DataVarianStok{
	"udang": {
		IDVarian:    1,
		IDProduk:    1,
		Slug:        "udang",
		NamaVarian:  "Udang",
		NamaProduk:  "Rengginang Udang",
		HargaJual:   28000,
		BeratGram:   450,
		Stok:        12, // Kasus TERSEDIA (12 bungkus)
		StatusAktif: true,
	},
	"bawang": {
		IDVarian:    2,
		IDProduk:    1,
		Slug:        "bawang",
		NamaVarian:  "Bawang",
		NamaProduk:  "Rengginang Bawang",
		HargaJual:   28000,
		BeratGram:   450,
		Stok:        38, // Kasus TERSEDIA (38 bungkus)
		StatusAktif: true,
	},
	"cumi": {
		IDVarian:    3,
		IDProduk:    1,
		Slug:        "cumi",
		NamaVarian:  "Cumi",
		NamaProduk:  "Rengginang Cumi",
		HargaJual:   28000,
		BeratGram:   450,
		Stok:        0, // Kasus STOK HABIS (0 bungkus)
		StatusAktif: true,
	},
	"terasi": {
		IDVarian:    4,
		IDProduk:    1,
		Slug:        "terasi",
		NamaVarian:  "Terasi",
		NamaProduk:  "Rengginang Terasi",
		HargaJual:   28000,
		BeratGram:   450,
		Stok:        38, // Kasus TERSEDIA (38 bungkus)
		StatusAktif: true,
	},
}

// ValidasiKuotaStok mengevaluasi batas jumlah pembelian terhadap stok tersedia.
// Aturan bisnis:
// 1. Jika stok <= 0              => STOK HABIS, tidak valid, tidak bisa dibeli
// 2. Jika jumlah <= 0            => JUMLAH TIDAK VALID, tidak valid, minimal 1 bungkus
// 3. Jika jumlah > stokTersedia  => MELEBIHI STOK, kuota tidak mencukupi
// 4. Jika jumlah <= stokTersedia => TERSEDIA, kuota valid dan dapat dibeli
func ValidasiKuotaStok(stokTersedia int64, jumlahDiminta int64) (valid bool, bisaDibeli bool, statusKetersediaan string, pesan string, keterangan string) {
	if stokTersedia <= 0 {
		return false, false, "STOK HABIS",
			"Produk sedang habis (out-of-stock). Pembelian dinonaktifkan.",
			"Silakan pilih varian rasa lain yang masih memiliki persediaan."
	}

	if jumlahDiminta <= 0 {
		return false, false, "JUMLAH TIDAK VALID",
			"Jumlah pembelian minimal 1 bungkus.",
			"Tambahkan jumlah pembelian menggunakan tombol stepper (+)."
	}

	if jumlahDiminta > stokTersedia {
		pesanError := fmt.Sprintf(
			"Jumlah pembelian (%d bungkus) melebihi batas kuota sisa stok yang tersedia (%d bungkus).",
			jumlahDiminta,
			stokTersedia,
		)
		keteranganError := fmt.Sprintf(
			"Maksimal pembelian untuk varian ini adalah %d bungkus. Kurangi jumlah pesanan Anda.",
			stokTersedia,
		)
		return false, false, "MELEBIHI STOK", pesanError, keteranganError
	}

	pesanSukses := fmt.Sprintf(
		"Jumlah pembelian valid (%d bungkus) dan kuota stok tersedia.",
		jumlahDiminta,
	)
	sisa := stokTersedia - jumlahDiminta
	keteranganSukses := fmt.Sprintf(
		"Stok mencukupi untuk pesanan Anda. Sisa stok setelah pembelian: %d bungkus.",
		sisa,
	)

	return true, true, "TERSEDIA", pesanSukses, keteranganSukses
}

// FormatRupiah memformat nominal angka menjadi format Rupiah standar (contoh: Rp28.000).
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

// AmbilDataVarian mencari data varian dari database MySQL atau fallback in-memory.
func (h *HandlerValidasiKuota) AmbilDataVarian(idNum uint64, slugOrName string) (*DataVarianStok, error) {
	slugClean := strings.ToLower(strings.TrimSpace(slugOrName))

	// 1. Prioritaskan pembacaan stok realtime langsung dari database MySQL jika koneksi aktif
	if h.DB != nil {
		var (
			idVarian    uint64
			idProduk    uint64
			namaVarian  string
			namaProduk  sql.NullString
			hargaJual   float64
			beratGram   uint32
			stok        int64
			statusAktif bool
		)

		query := `
			SELECT 
				v.id_varian,
				v.id_produk,
				v.nama_varian,
				p.nama_produk,
				v.harga_jual,
				v.berat_gram,
				v.stok,
				v.status_aktif
			FROM varian_produk v
			LEFT JOIN produk p ON v.id_produk = p.id_produk
			WHERE (
				(? > 0 AND v.id_varian = ?) OR
				(? != '' AND (LOWER(v.nama_varian) = ? OR LOWER(v.nama_varian) LIKE ?))
			) AND v.status_aktif = TRUE
			LIMIT 1
		`

		likePattern := "%" + slugClean + "%"
		err := h.DB.QueryRow(query, idNum, idNum, slugClean, slugClean, likePattern).Scan(
			&idVarian,
			&idProduk,
			&namaVarian,
			&namaProduk,
			&hargaJual,
			&beratGram,
			&stok,
			&statusAktif,
		)

		if err == nil {
			namaProd := "Rengginang " + namaVarian
			if namaProduk.Valid && namaProduk.String != "" {
				namaProd = namaProduk.String
			}

			computedSlug := strings.ToLower(namaVarian)
			return &DataVarianStok{
				IDVarian:    idVarian,
				IDProduk:    idProduk,
				Slug:        computedSlug,
				NamaVarian:  namaVarian,
				NamaProduk:  namaProd,
				HargaJual:   uint64(hargaJual),
				BeratGram:   beratGram,
				Stok:        stok,
				StatusAktif: statusAktif,
			}, nil
		}
	}

	// 2. Fallback aman ke katalog default (apabila DB belum dijalankan atau data awal belum di-seed)
	if slugClean != "" {
		if item, ada := DataDefaultVarian[slugClean]; ada {
			return &item, nil
		}
		// Pencocokan substring (misal "rengginang udang")
		for k, v := range DataDefaultVarian {
			if strings.Contains(slugClean, k) || strings.Contains(k, slugClean) {
				return &v, nil
			}
		}
	}

	if idNum > 0 {
		for _, v := range DataDefaultVarian {
			if v.IDVarian == idNum {
				return &v, nil
			}
		}
	}

	return nil, errors.New("varian produk tidak ditemukan")
}

// HandlerCekKuotaRealtime menangani validasi kuota stok realtime.
// Mendukung:
// - GET  /api/katalog/cek-kuota?id=udang&jumlah=2
// - POST /api/katalog/validasi-kuota  dengan JSON body
// - POST /api/katalog/validasi-beli   dengan JSON body
func (h *HandlerValidasiKuota) HandlerCekKuotaRealtime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var (
		idNum      uint64
		slug       string
		jumlahBeli int64
	)

	switch r.Method {
	case http.MethodGet:
		// Parsing parameter URL query
		q := r.URL.Query()
		idParam := strings.TrimSpace(q.Get("id"))
		if idParam == "" {
			idParam = strings.TrimSpace(q.Get("slug"))
		}
		if idParam == "" {
			idParam = strings.TrimSpace(q.Get("id_varian"))
		}

		// Deteksi apakah idParam berupa angka atau nama varian
		if parsedNum, err := strconv.ParseUint(idParam, 10, 64); err == nil {
			idNum = parsedNum
		} else {
			slug = idParam
		}

		jumlahParam := strings.TrimSpace(q.Get("jumlah"))
		if jumlahParam == "" {
			jumlahParam = strings.TrimSpace(q.Get("jumlah_beli"))
		}

		if parsedJumlah, err := strconv.ParseInt(jumlahParam, 10, 64); err == nil {
			jumlahBeli = parsedJumlah
		} else {
			jumlahBeli = 1 // Default jika tidak disertakan
		}

	case http.MethodPost:
		// Parsing JSON body
		var req PermintaanValidasiKuota
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			kirimJSON(w, http.StatusBadRequest, ResponsAPI{
				Status: "gagal",
				Pesan:  "Format payload JSON tidak valid",
			})
			return
		}
		idNum = req.IDVarian
		slug = req.Slug
		jumlahBeli = req.JumlahBeli

	default:
		w.Header().Set("Allow", "GET, POST")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET atau POST.",
		})
		return
	}

	// Ambil data varian dan kuota stok aktual
	varian, err := h.AmbilDataVarian(idNum, slug)
	if err != nil {
		kirimJSON(w, http.StatusNotFound, ResponsAPI{
			Status: "gagal",
			Pesan:  "Varian produk tidak ditemukan untuk validasi kuota stok",
		})
		return
	}

	// Evaluasi aturan bisnis validasi kuota stok realtime
	valid, bisaDibeli, statusKetersediaan, pesan, keterangan := ValidasiKuotaStok(varian.Stok, jumlahBeli)

	// Hitung total harga dan sisa stok setelah beli
	var totalHarga uint64
	if valid {
		totalHarga = uint64(jumlahBeli) * varian.HargaJual
	} else if varian.Stok > 0 && jumlahBeli > 0 {
		totalHarga = uint64(jumlahBeli) * varian.HargaJual
	} else {
		totalHarga = 0
	}

	sisaSetelahBeli := varian.Stok - jumlahBeli

	hasil := &HasilValidasiKuota{
		IDVarian:           varian.IDVarian,
		Slug:               varian.Slug,
		NamaVarian:         varian.NamaVarian,
		HargaSatuan:        varian.HargaJual,
		FormatHarga:        FormatRupiah(varian.HargaJual),
		StokTersedia:       varian.Stok,
		FormatStok:         fmt.Sprintf("%d bungkus", varian.Stok),
		BatasMaksimalBeli:  varian.Stok,
		JumlahDiminta:      jumlahBeli,
		TotalHarga:         totalHarga,
		FormatTotalHarga:   FormatRupiah(totalHarga),
		SisaSetelahBeli:    sisaSetelahBeli,
		StatusKetersediaan: statusKetersediaan,
		KuotaTersedia:      valid,
		BisaDibeli:         bisaDibeli,
		Valid:              valid,
		Pesan:              pesan,
		Keterangan:         keterangan,
	}

	// Jika validasi gagal (melebihi stok, stok habis, atau jumlah nol), kembalikan HTTP 400
	if !valid {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "ditolak",
			Pesan:  pesan,
			Data:   hasil,
		})
		return
	}

	// Jika validasi sukses, kembalikan HTTP 200
	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  pesan,
		Data:   hasil,
	})
}

// kirimJSON helper mengirimkan respon JSON terformat.
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
