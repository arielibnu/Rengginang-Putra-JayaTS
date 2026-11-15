package produkapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type inputVarian struct {
	NamaVarian  string `json:"nama_varian"`
	HargaJual   int64  `json:"harga_jual"`
	BeratGram   int64  `json:"berat_gram"`
	Stok        int64  `json:"stok"`
	StatusAktif *bool  `json:"status_aktif"`
}

// HandlerVarian menangani permintaan API varian.
type HandlerVarian struct {
	DB *sql.DB
}

// Endpoint:
// GET, POST /api/admin/produk/{id}/varian
// PUT, DELETE /api/admin/varian/{id}

func (h *HandlerVarian) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		responsProduk(w, 500, "Database belum tersedia", nil)
		return
	}

	path := strings.Trim(
		strings.TrimPrefix(r.URL.Path, "/api/admin/varian/"),
		"/",
	)

	id, err := strconv.ParseUint(path, 10, 64)
	if err != nil || id == 0 || strings.Contains(path, "/") {
		responsProduk(w, 400, "ID varian tidak valid", nil)
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.ubahVarian(w, r, id)
	case http.MethodDelete:
		h.nonaktifkanVarian(w, r, id)
	default:
		w.Header().Set("Allow", "PUT, DELETE")
		responsProduk(w, 405, "Metode tidak diizinkan", nil)
	}
}

// HandlerVarianProduk menangani varian milik produk tertentu.
func (h *HandlerVarian) HandlerVarianProduk(
	w http.ResponseWriter,
	r *http.Request,
	idProduk uint64,
) {
	switch r.Method {
	case http.MethodGet:
		h.daftarVarian(w, r, idProduk)
	case http.MethodPost:
		h.tambahVarian(w, r, idProduk)
	default:
		w.Header().Set("Allow", "GET, POST")
		responsProduk(w, 405, "Metode tidak diizinkan", nil)
	}
}

func bacaInputVarian(w http.ResponseWriter, r *http.Request) (inputVarian, error) {
	var input inputVarian

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		return input, errors.New("Format JSON varian tidak valid")
	}

	var tambahan interface{}
	if err := decoder.Decode(&tambahan); err != io.EOF {
		return input, errors.New("JSON harus berisi tepat satu objek")
	}

	input.NamaVarian = strings.TrimSpace(input.NamaVarian)

	if input.NamaVarian == "" ||
		len(input.NamaVarian) > 100 {
		return input, errors.New(
			"Nama varian wajib diisi, maksimal 100 karakter",
		)
	}

	if input.HargaJual < 0 || input.HargaJual > 9999999999 {
		return input, errors.New("Harga jual tidak valid")
	}

	if input.BeratGram <= 0 || input.BeratGram > 100000000 {
		return input, errors.New("Berat harus lebih dari 0 gram")
	}

	if input.Stok < 0 || input.Stok > 1000000000 {
		return input, errors.New("Stok tidak valid")
	}

	return input, nil
}

func (h *HandlerVarian) daftarVarian(
	w http.ResponseWriter,
	r *http.Request,
	idProduk uint64,
) {
	var ada int

	err := h.DB.QueryRowContext(
		r.Context(),
		"SELECT 1 FROM produk WHERE id_produk = ?",
		idProduk,
	).Scan(&ada)

	if errors.Is(err, sql.ErrNoRows) {
		responsProduk(w, 404, "Produk tidak ditemukan", nil)
		return
	}
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa produk", nil)
		return
	}

	rows, err := h.DB.QueryContext(
		r.Context(),
		`SELECT id_varian, id_produk, nama_varian,
		        harga_jual, berat_gram, stok, status_aktif,
		        dibuat_pada, diperbarui_pada
		 FROM varian_produk
		 WHERE id_produk = ?
		 ORDER BY id_varian DESC`,
		idProduk,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal mengambil varian", nil)
		return
	}
	defer rows.Close()

	hasil := make([]VarianProduk, 0)

	for rows.Next() {
		var v VarianProduk

		if err := rows.Scan(
			&v.IDVarian,
			&v.IDProduk,
			&v.NamaVarian,
			&v.HargaJual,
			&v.BeratGram,
			&v.Stok,
			&v.StatusAktif,
			&v.DibuatPada,
			&v.DiperbaruiPada,
		); err != nil {
			responsProduk(w, 500, "Gagal membaca data varian", nil)
			return
		}

		hasil = append(hasil, v)
	}

	if err := rows.Err(); err != nil {
		responsProduk(w, 500, "Gagal membaca daftar varian", nil)
		return
	}

	responsProduk(w, 200, "Daftar varian berhasil diambil", hasil)
}

func (h *HandlerVarian) tambahVarian(
	w http.ResponseWriter,
	r *http.Request,
	idProduk uint64,
) {
	input, err := bacaInputVarian(w, r)
	if err != nil {
		responsProduk(w, 400, err.Error(), nil)
		return
	}

	var ada int
	err = h.DB.QueryRowContext(
		r.Context(),
		"SELECT 1 FROM produk WHERE id_produk = ?",
		idProduk,
	).Scan(&ada)

	if errors.Is(err, sql.ErrNoRows) {
		responsProduk(w, 404, "Produk tidak ditemukan", nil)
		return
	}
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa produk", nil)
		return
	}

	aktif := true
	if input.StatusAktif != nil {
		aktif = *input.StatusAktif
	}

	result, err := h.DB.ExecContext(
		r.Context(),
		`INSERT INTO varian_produk
		 (id_produk, nama_varian, harga_jual, berat_gram,
		  stok, status_aktif)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		idProduk,
		input.NamaVarian,
		input.HargaJual,
		input.BeratGram,
		input.Stok,
		aktif,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal menambahkan varian. Periksa nama varian.", nil)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		responsProduk(w, 500, "Varian tersimpan tetapi ID gagal dibaca", nil)
		return
	}

	responsProduk(w, 201, "Varian berhasil ditambahkan",
		map[string]interface{}{
			"id_varian":    id,
			"id_produk":    idProduk,
			"nama_varian":  input.NamaVarian,
			"harga_jual":   input.HargaJual,
			"berat_gram":   input.BeratGram,
			"stok":         input.Stok,
			"status_aktif": aktif,
		},
	)
}

func (h *HandlerVarian) ubahVarian(
	w http.ResponseWriter,
	r *http.Request,
	id uint64,
) {
	input, err := bacaInputVarian(w, r)
	if err != nil {
		responsProduk(w, 400, err.Error(), nil)
		return
	}

	aktif := true
	if input.StatusAktif != nil {
		aktif = *input.StatusAktif
	}

	result, err := h.DB.ExecContext(
		r.Context(),
		`UPDATE varian_produk
		 SET nama_varian = ?, harga_jual = ?,
		     berat_gram = ?, stok = ?, status_aktif = ?
		 WHERE id_varian = ?`,
		input.NamaVarian,
		input.HargaJual,
		input.BeratGram,
		input.Stok,
		aktif,
		id,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal mengubah varian", nil)
		return
	}

	jumlah, err := result.RowsAffected()
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa perubahan varian", nil)
		return
	}
	if jumlah == 0 {
		responsProduk(w, 404, "Varian tidak ditemukan atau tidak berubah", nil)
		return
	}

	responsProduk(w, 200, "Varian berhasil diubah", nil)
}

func (h *HandlerVarian) nonaktifkanVarian(
	w http.ResponseWriter,
	r *http.Request,
	id uint64,
) {
	result, err := h.DB.ExecContext(
		r.Context(),
		`UPDATE varian_produk
		 SET status_aktif = FALSE
		 WHERE id_varian = ?`,
		id,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal menonaktifkan varian", nil)
		return
	}

	jumlah, err := result.RowsAffected()
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa status varian", nil)
		return
	}
	if jumlah == 0 {
		responsProduk(w, 404, "Varian tidak ditemukan atau sudah nonaktif", nil)
		return
	}

	responsProduk(w, 200, "Varian berhasil dinonaktifkan", nil)
}
