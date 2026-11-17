package stokapi

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	validasistok "rengginang-putra-jaya/auth/1.1.4.4"
)

type Handler struct {
	DB *sql.DB
}

// Menangani endpoint mutasi dan riwayat stok.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/api/admin/stok/mutasi":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			kirimJSON(w, http.StatusMethodNotAllowed,
				ResponsAPI{Pesan: "Metode tidak diizinkan"})
			return
		}

		h.tambahMutasi(w, r)

	case "/api/admin/stok/riwayat":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			kirimJSON(w, http.StatusMethodNotAllowed,
				ResponsAPI{Pesan: "Metode tidak diizinkan"})
			return
		}

		h.ambilRiwayat(w, r)

	default:
		kirimJSON(w, http.StatusNotFound,
			ResponsAPI{Pesan: "Endpoint tidak ditemukan"})
	}
}

// Mencatat stok masuk, keluar, dan penyesuaian.
func (h *Handler) tambahMutasi(w http.ResponseWriter, r *http.Request) {

	// Membatasi ukuran request maksimal 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req PermintaanMutasi

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Format JSON tidak valid"})
		return
	}

	// Memastikan hanya ada satu objek JSON.
	var tambahan interface{}

	if err := decoder.Decode(&tambahan); err != io.EOF {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Permintaan JSON tidak valid"})
		return
	}

	req.Catatan = strings.TrimSpace(req.Catatan)

	// ==========================================
	// VALIDASI INPUT
	// ==========================================

	// Validasi ID varian.
	if req.IDVarian == 0 {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "ID varian wajib diisi"})
		return
	}

	// Validasi dan normalisasi mode stok.
	// Mode kosong otomatis menjadi RIIL.
	modeStok, errMode := validasistok.NormalisasiMode(
		req.ModeStok,
	)

	if errMode != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: errMode.Error()})
		return
	}

	req.ModeStok = modeStok

	// Validasi dan normalisasi jenis mutasi.
	jenisMutasi, errJenis := validasistok.NormalisasiJenis(
		req.JenisMutasi,
	)

	if errJenis != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: errJenis.Error()})
		return
	}

	req.JenisMutasi = jenisMutasi

	// ==========================================
	// MEMULAI TRANSAKSI DATABASE
	// ==========================================

	tx, err := h.DB.Begin()

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal memulai transaksi database"})
		return
	}

	defer tx.Rollback()

	var stokSebelum uint64

	// ==========================================
	// MEMBACA STOK RIIL
	// ==========================================

	if req.ModeStok == "RIIL" {

		err = tx.QueryRow(`
			SELECT stok
			FROM varian_produk
			WHERE id_varian = ?
			FOR UPDATE
		`, req.IDVarian).Scan(&stokSebelum)

	} else {

		// ==========================================
		// MEMBACA STOK SIMULASI
		// ==========================================

		// Memastikan varian tersedia dan mengambil
		// stok riil sebagai nilai awal simulasi.
		var stokAwalRiil uint64

		err = tx.QueryRow(`
			SELECT stok
			FROM varian_produk
			WHERE id_varian = ?
			FOR UPDATE
		`, req.IDVarian).Scan(&stokAwalRiil)

		if err == nil {

			// Membuat stok simulasi jika belum ada.
			// Jika sudah ada, nilainya tidak diubah.
			_, err = tx.Exec(`
				INSERT INTO stok_simulasi (
					id_varian,
					stok_simulasi
				)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE
					stok_simulasi = stok_simulasi
			`, req.IDVarian, stokAwalRiil)
		}

		if err == nil {

			err = tx.QueryRow(`
				SELECT stok_simulasi
				FROM stok_simulasi
				WHERE id_varian = ?
				FOR UPDATE
			`, req.IDVarian).Scan(&stokSebelum)
		}
	}

	// Memeriksa apakah varian ditemukan.
	if err == sql.ErrNoRows {
		kirimJSON(w, http.StatusNotFound,
			ResponsAPI{Pesan: "Varian produk tidak ditemukan"})
		return
	}

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal membaca data stok"})
		return
	}

	// ==========================================
	// VALIDASI DAN PERHITUNGAN STOK
	// WBS 1.1.4.4
	// ==========================================

	jumlah, stokSesudah, errHitung := validasistok.HitungPerubahan(
		req.JenisMutasi,
		stokSebelum,
		req.Jumlah,
		req.StokBaru,
	)

	if errHitung != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: errHitung.Error()})
		return
	}

	// ==========================================
	// MEMPERBARUI STOK
	// ==========================================

	if req.ModeStok == "RIIL" {

		// Memperbarui stok asli produk.
		_, err = tx.Exec(`
			UPDATE varian_produk
			SET stok = ?,
				diperbarui_pada = CURRENT_TIMESTAMP
			WHERE id_varian = ?
		`, stokSesudah, req.IDVarian)

	} else {

		// Memperbarui stok simulasi saja.
		// Stok riil tidak ikut berubah.
		_, err = tx.Exec(`
			UPDATE stok_simulasi
			SET stok_simulasi = ?,
				diperbarui_pada = CURRENT_TIMESTAMP
			WHERE id_varian = ?
		`, stokSesudah, req.IDVarian)
	}

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal memperbarui stok"})
		return
	}

	// ==========================================
	// MENYIMPAN RIWAYAT MUTASI
	// ==========================================

	_, err = tx.Exec(`
		INSERT INTO riwayat_mutasi_stok (
			id_varian,
			jenis_mutasi,
			mode_stok,
			jumlah,
			stok_sebelum,
			stok_sesudah,
			catatan
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		req.IDVarian,
		req.JenisMutasi,
		req.ModeStok,
		jumlah,
		stokSebelum,
		stokSesudah,
		req.Catatan,
	)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menyimpan riwayat mutasi"})
		return
	}

	// ==========================================
	// MENYELESAIKAN TRANSAKSI
	// ==========================================

	if err := tx.Commit(); err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menyimpan transaksi"})
		return
	}

	// ==========================================
	// RESPONS BERHASIL
	// ==========================================

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Pesan: "Mutasi stok berhasil disimpan",
		Data: map[string]interface{}{
			"id_varian":    req.IDVarian,
			"jenis_mutasi": req.JenisMutasi,
			"mode_stok":    req.ModeStok,
			"stok_sebelum": stokSebelum,
			"stok_sesudah": stokSesudah,
			"jumlah":       jumlah,
		},
	})
}

// ==========================================
// MENGAMBIL RIWAYAT MUTASI STOK
// ==========================================

func (h *Handler) ambilRiwayat(w http.ResponseWriter, r *http.Request) {

	query := `
		SELECT
			id_mutasi,
			id_varian,
			jenis_mutasi,
			mode_stok,
			jumlah,
			stok_sebelum,
			stok_sesudah,
			COALESCE(catatan, ''),
			dibuat_pada
		FROM riwayat_mutasi_stok
		WHERE 1=1
	`

	var args []interface{}

	// ==========================================
	// FILTER BERDASARKAN ID VARIAN
	// ==========================================

	if id := r.URL.Query().Get("id_varian"); id != "" {

		idVarian, err := strconv.ParseUint(id, 10, 64)

		if err != nil || idVarian == 0 {
			kirimJSON(w, http.StatusBadRequest,
				ResponsAPI{Pesan: "ID varian tidak valid"})
			return
		}

		query += " AND id_varian = ?"
		args = append(args, idVarian)
	}

	// ==========================================
	// FILTER BERDASARKAN MODE STOK
	// ==========================================

	if mode := strings.ToUpper(
		strings.TrimSpace(r.URL.Query().Get("mode_stok")),
	); mode != "" {

		modeValid, err := validasistok.NormalisasiMode(mode)

		if err != nil {
			kirimJSON(w, http.StatusBadRequest,
				ResponsAPI{Pesan: err.Error()})
			return
		}

		query += " AND mode_stok = ?"
		args = append(args, modeValid)
	}

	// Mengambil maksimal 200 riwayat terbaru.
	query += " ORDER BY dibuat_pada DESC LIMIT 200"

	rows, err := h.DB.Query(query, args...)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal mengambil riwayat stok"})
		return
	}

	defer rows.Close()

	hasil := make([]RiwayatMutasi, 0)

	// ==========================================
	// MEMBACA DATA RIWAYAT
	// ==========================================

	for rows.Next() {

		var item RiwayatMutasi

		err := rows.Scan(
			&item.IDMutasi,
			&item.IDVarian,
			&item.JenisMutasi,
			&item.ModeStok,
			&item.Jumlah,
			&item.StokSebelum,
			&item.StokSesudah,
			&item.Catatan,
			&item.DibuatPada,
		)

		if err != nil {
			kirimJSON(w, http.StatusInternalServerError,
				ResponsAPI{Pesan: "Gagal membaca riwayat stok"})
			return
		}

		hasil = append(hasil, item)
	}

	if err := rows.Err(); err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Terjadi kesalahan saat membaca data"})
		return
	}

	// ==========================================
	// RESPONS RIWAYAT BERHASIL
	// ==========================================

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Pesan: "Riwayat mutasi stok berhasil diambil",
		Data:  hasil,
	})
}

// ==========================================
// MENGIRIM RESPONS JSON
// ==========================================

func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
