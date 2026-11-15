package snapshot

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	DB *sql.DB
}

type ItemPermintaan struct {
	IDVarian uint64 `json:"id_varian"`
	Jumlah   uint64 `json:"jumlah"`
}

type PermintaanTransaksi struct {
	Items []ItemPermintaan `json:"items"`
}

type DetailTransaksi struct {
	IDDetail           uint64 `json:"id_detail"`
	IDVarian           uint64 `json:"id_varian"`
	NamaVarianSnapshot string `json:"nama_varian_snapshot"`
	HargaJualSnapshot  string `json:"harga_jual_snapshot"`
	Jumlah             uint64 `json:"jumlah"`
	Subtotal           string `json:"subtotal"`
}

func balas(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// ServeHTTP menangani pembuatan dan pembacaan transaksi.
// Proteksi sesi Owner diterapkan di router utama.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		balas(w, 500, map[string]string{"pesan": "Database belum tersedia"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/admin/transaksi")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodPost:
			h.buatTransaksi(w, r)
		default:
			w.Header().Set("Allow", "POST")
			balas(w, 405, map[string]string{"pesan": "Metode tidak diizinkan"})
		}
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		balas(w, 405, map[string]string{"pesan": "Metode tidak diizinkan"})
		return
	}

	if strings.Contains(path, "/") {
		balas(w, 404, map[string]string{"pesan": "Endpoint tidak ditemukan"})
		return
	}

	id, err := strconv.ParseUint(path, 10, 64)
	if err != nil || id == 0 {
		balas(w, 400, map[string]string{"pesan": "ID transaksi tidak valid"})
		return
	}

	h.detailTransaksi(w, r, id)
}

func (h *Handler) buatTransaksi(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input PermintaanTransaksi
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		balas(w, 400, map[string]string{"pesan": "JSON tidak valid"})
		return
	}

	var tambahan interface{}
	if err := decoder.Decode(&tambahan); err != io.EOF {
		balas(w, 400, map[string]string{"pesan": "JSON harus satu objek"})
		return
	}

	if len(input.Items) == 0 || len(input.Items) > 50 {
		balas(w, 400, map[string]string{"pesan": "Jumlah item harus 1 sampai 50"})
		return
	}

	for _, item := range input.Items {
		if item.IDVarian == 0 || item.Jumlah == 0 || item.Jumlah > 100000 {
			balas(w, 400, map[string]string{"pesan": "ID varian atau jumlah tidak valid"})
			return
		}
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal memulai transaksi"})
		return
	}
	defer tx.Rollback()

	// Membuat transaksi baru.
	result, err := tx.ExecContext(
		r.Context(),
		"INSERT INTO transaksi () VALUES ()",
	)
	if err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal menyimpan transaksi"})
		return
	}

	idTransaksi, err := result.LastInsertId()
	if err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal membaca ID transaksi"})
		return
	}

	// Harga diambil dari database, bukan dari browser.
	for _, item := range input.Items {
		var nama string
		var harga string

		err := tx.QueryRowContext(
			r.Context(),
			`SELECT nama_varian, CAST(harga_jual AS CHAR)
			 FROM varian_produk
			 WHERE id_varian = ? AND status_aktif = TRUE
			 FOR UPDATE`,
			item.IDVarian,
		).Scan(&nama, &harga)

		if errors.Is(err, sql.ErrNoRows) {
			balas(w, 400, map[string]string{
				"pesan": "Varian tidak ditemukan atau tidak aktif",
			})
			return
		}
		if err != nil {
			balas(w, 500, map[string]string{"pesan": "Gagal mengambil harga varian"})
			return
		}

		// Hitung subtotal menggunakan operasi DECIMAL MySQL
		// agar nilai uang tidak bergantung pada float64.
		_, err = tx.ExecContext(
			r.Context(),
			`INSERT INTO detail_transaksi
			 (id_transaksi, id_varian, nama_varian_snapshot,
			  harga_jual_snapshot, jumlah, subtotal)
			 VALUES (?, ?, ?, ?, ?, CAST(? AS DECIMAL(12,2)) * ?)`,
			idTransaksi,
			item.IDVarian,
			nama,
			harga,
			item.Jumlah,
			harga,
			item.Jumlah,
		)
		if err != nil {
			balas(w, 500, map[string]string{"pesan": "Gagal menyimpan snapshot"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal menyelesaikan transaksi"})
		return
	}

	balas(w, 201, map[string]interface{}{
		"pesan":        "Transaksi dan snapshot harga berhasil disimpan",
		"id_transaksi": idTransaksi,
	})
}

func (h *Handler) detailTransaksi(
	w http.ResponseWriter,
	r *http.Request,
	id uint64,
) {
	var ada int

	err := h.DB.QueryRowContext(
		r.Context(),
		"SELECT 1 FROM transaksi WHERE id_transaksi = ?",
		id,
	).Scan(&ada)

	if errors.Is(err, sql.ErrNoRows) {
		balas(w, 404, map[string]string{"pesan": "Transaksi tidak ditemukan"})
		return
	}
	if err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal memeriksa transaksi"})
		return
	}

	rows, err := h.DB.QueryContext(
		r.Context(),
		`SELECT id_detail, id_varian, nama_varian_snapshot,
		        CAST(harga_jual_snapshot AS CHAR),
		        jumlah, CAST(subtotal AS CHAR)
		 FROM detail_transaksi
		 WHERE id_transaksi = ?
		 ORDER BY id_detail`,
		id,
	)
	if err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal mengambil detail transaksi"})
		return
	}
	defer rows.Close()

	daftar := make([]DetailTransaksi, 0)

	for rows.Next() {
		var detail DetailTransaksi

		err := rows.Scan(
			&detail.IDDetail,
			&detail.IDVarian,
			&detail.NamaVarianSnapshot,
			&detail.HargaJualSnapshot,
			&detail.Jumlah,
			&detail.Subtotal,
		)
		if err != nil {
			balas(w, 500, map[string]string{"pesan": "Gagal membaca detail transaksi"})
			return
		}

		daftar = append(daftar, detail)
	}

	if err := rows.Err(); err != nil {
		balas(w, 500, map[string]string{"pesan": "Gagal membaca seluruh detail transaksi"})
		return
	}

	balas(w, 200, map[string]interface{}{
		"pesan":        fmt.Sprintf("Detail transaksi %d berhasil diambil", id),
		"id_transaksi": id,
		"items":        daftar,
	})
}
