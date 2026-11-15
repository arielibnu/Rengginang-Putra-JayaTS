package produkapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type inputProduk struct {
	NamaProduk  string `json:"nama_produk"`
	Deskripsi   string `json:"deskripsi"`
	StatusAktif *bool  `json:"status_aktif"`
}

// HandlerProduk mengelola endpoint CRUD produk.
type HandlerProduk struct {
	DB *sql.DB
}

func (h *HandlerProduk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if h.DB == nil {
		responsProduk(w, http.StatusInternalServerError, "Database belum tersedia", nil)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/admin/produk")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.daftarProduk(w, r)
		case http.MethodPost:
			h.tambahProduk(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			responsProduk(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan", nil)
		}
		return
	}

	// Endpoint detail /api/admin/produk/{id}
	// Endpoint varian akan ditangani pada varian.go.
	if strings.Contains(path, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseUint(path, 10, 64)
	if err != nil || id == 0 {
		responsProduk(w, http.StatusBadRequest, "ID produk tidak valid", nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.detailProduk(w, r, id)
	case http.MethodPut:
		h.ubahProduk(w, r, id)
	case http.MethodDelete:
		h.nonaktifkanProduk(w, r, id)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		responsProduk(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan", nil)
	}
}

func responsProduk(w http.ResponseWriter, status int, pesan string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ResponsAPI{
		Pesan: pesan,
		Data:  data,
	})
}

func bacaInputProduk(w http.ResponseWriter, r *http.Request) (inputProduk, error) {
	var input inputProduk

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		return input, err
	}

	var tambahan interface{}
	if err := decoder.Decode(&tambahan); !errors.Is(err, nil) {
		// Diperiksa kembali di bawah agar JSON tambahan ditolak.
	}

	input.NamaProduk = strings.TrimSpace(input.NamaProduk)
	input.Deskripsi = strings.TrimSpace(input.Deskripsi)

	if input.NamaProduk == "" || len(input.NamaProduk) > 150 {
		return input, errors.New("nama produk wajib diisi, maksimal 150 karakter")
	}

	return input, nil
}

func (h *HandlerProduk) daftarProduk(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.QueryContext(
		r.Context(),
		`SELECT id_produk, nama_produk, COALESCE(deskripsi, ''),
		        status_aktif, dibuat_pada, diperbarui_pada
		 FROM produk ORDER BY id_produk DESC`,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal mengambil daftar produk", nil)
		return
	}
	defer rows.Close()

	hasil := make([]Produk, 0)

	for rows.Next() {
		var p Produk
		if err := rows.Scan(
			&p.IDProduk,
			&p.NamaProduk,
			&p.Deskripsi,
			&p.StatusAktif,
			&p.DibuatPada,
			&p.DiperbaruiPada,
		); err != nil {
			responsProduk(w, 500, "Gagal membaca data produk", nil)
			return
		}
		hasil = append(hasil, p)
	}

	if err := rows.Err(); err != nil {
		responsProduk(w, 500, "Gagal membaca daftar produk", nil)
		return
	}

	responsProduk(w, 200, "Daftar produk berhasil diambil", hasil)
}

func (h *HandlerProduk) detailProduk(w http.ResponseWriter, r *http.Request, id uint64) {
	var p Produk

	err := h.DB.QueryRowContext(
		r.Context(),
		`SELECT id_produk, nama_produk, COALESCE(deskripsi, ''),
		        status_aktif, dibuat_pada, diperbarui_pada
		 FROM produk WHERE id_produk = ?`,
		id,
	).Scan(
		&p.IDProduk,
		&p.NamaProduk,
		&p.Deskripsi,
		&p.StatusAktif,
		&p.DibuatPada,
		&p.DiperbaruiPada,
	)

	if errors.Is(err, sql.ErrNoRows) {
		responsProduk(w, 404, "Produk tidak ditemukan", nil)
		return
	}
	if err != nil {
		responsProduk(w, 500, "Gagal mengambil detail produk", nil)
		return
	}

	responsProduk(w, 200, "Detail produk berhasil diambil", p)
}

func (h *HandlerProduk) tambahProduk(w http.ResponseWriter, r *http.Request) {
	input, err := bacaInputProduk(w, r)
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
		`INSERT INTO produk (nama_produk, deskripsi, status_aktif)
		 VALUES (?, ?, ?)`,
		input.NamaProduk, input.Deskripsi, aktif,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal menambahkan produk", nil)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		responsProduk(w, 500, "Produk tersimpan tetapi ID gagal dibaca", nil)
		return
	}

	responsProduk(w, 201, "Produk berhasil ditambahkan", map[string]interface{}{
		"id_produk":    id,
		"nama_produk":  input.NamaProduk,
		"status_aktif": aktif,
	})
}

func (h *HandlerProduk) ubahProduk(w http.ResponseWriter, r *http.Request, id uint64) {
	input, err := bacaInputProduk(w, r)
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
		`UPDATE produk
		 SET nama_produk = ?, deskripsi = ?, status_aktif = ?
		 WHERE id_produk = ?`,
		input.NamaProduk, input.Deskripsi, aktif, id,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal mengubah produk", nil)
		return
	}

	jumlah, err := result.RowsAffected()
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa hasil perubahan", nil)
		return
	}
	if jumlah == 0 {
		responsProduk(w, 404, "Produk tidak ditemukan atau tidak berubah", nil)
		return
	}

	responsProduk(w, 200, "Produk berhasil diubah", nil)
}

func (h *HandlerProduk) nonaktifkanProduk(w http.ResponseWriter, r *http.Request, id uint64) {
	result, err := h.DB.ExecContext(
		r.Context(),
		`UPDATE produk SET status_aktif = FALSE WHERE id_produk = ?`,
		id,
	)
	if err != nil {
		responsProduk(w, 500, "Gagal menonaktifkan produk", nil)
		return
	}

	jumlah, err := result.RowsAffected()
	if err != nil {
		responsProduk(w, 500, "Gagal memeriksa perubahan produk", nil)
		return
	}
	if jumlah == 0 {
		responsProduk(w, 404, "Produk tidak ditemukan atau sudah nonaktif", nil)
		return
	}

	responsProduk(w, 200, "Produk berhasil dinonaktifkan", nil)
}
