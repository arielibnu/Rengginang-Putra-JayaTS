package modalapi

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	validasibiaya "rengginang-putra-jaya/auth/1.1.6.4"
)

// Handler utama untuk endpoint modal dan pengeluaran.
type Handler struct {
	DB *sql.DB
}

// Menangani routing endpoint modal dan pengeluaran.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	// ==========================================
	// ENDPOINT MODAL PER VARIAN
	// ==========================================

	case "/api/admin/modal":
		switch r.Method {
		case http.MethodGet:
			h.ambilDaftarModal(w, r)
		case http.MethodPost:
			h.simpanModal(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			kirimJSON(w, http.StatusMethodNotAllowed,
				ResponsAPI{Pesan: "Metode tidak diizinkan"})
		}

	// ==========================================
	// ENDPOINT PENGELUARAN OPERASIONAL
	// ==========================================

	case "/api/admin/pengeluaran":
		switch r.Method {
		case http.MethodGet:
			h.ambilDaftarPengeluaran(w, r)
		case http.MethodPost:
			h.tambahPengeluaran(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			kirimJSON(w, http.StatusMethodNotAllowed,
				ResponsAPI{Pesan: "Metode tidak diizinkan"})
		}

	// ==========================================
	// ENDPOINT RINGKASAN BULANAN
	// ==========================================

	case "/api/admin/modal/ringkasan":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			kirimJSON(w, http.StatusMethodNotAllowed,
				ResponsAPI{Pesan: "Metode tidak diizinkan"})
			return
		}
		h.ambilRingkasanBulanan(w, r)

	default:
		kirimJSON(w, http.StatusNotFound,
			ResponsAPI{Pesan: "Endpoint tidak ditemukan"})
	}
}

// ==========================================
// SIMPAN MODAL VARIAN BARU
// ==========================================

// Menyimpan harga modal baru untuk suatu varian.
// Baris lama tidak dihapus supaya riwayat perubahan tetap ada.
func (h *Handler) simpanModal(w http.ResponseWriter, r *http.Request) {

	// Membatasi ukuran permintaan menjadi 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req PermintaanModal

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
	req.BerlakuMulai = strings.TrimSpace(req.BerlakuMulai)

	// ==========================================
	// VALIDASI INPUT
	// ==========================================

	if req.IDVarian == 0 {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "ID varian wajib diisi"})
		return
	}

	if req.HargaModal <= 0 {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Harga modal harus lebih dari 0"})
		return
	}

	// Validasi format tanggal berlaku.
	if req.BerlakuMulai == "" {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Tanggal berlaku mulai wajib diisi"})
		return
	}

	_, errTanggal := time.Parse("2006-01-02", req.BerlakuMulai)

	if errTanggal != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Format tanggal tidak valid, gunakan YYYY-MM-DD"})
		return
	}

	// ==========================================
	// MEMASTIKAN VARIAN ADA DI DATABASE
	// ==========================================

	var namaVarian string

	err := h.DB.QueryRow(`
		SELECT nama_varian
		FROM varian_produk
		WHERE id_varian = ?
	`, req.IDVarian).Scan(&namaVarian)

	if err == sql.ErrNoRows {
		kirimJSON(w, http.StatusNotFound,
			ResponsAPI{Pesan: "Varian produk tidak ditemukan"})
		return
	}

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal memeriksa varian produk"})
		return
	}

	// ==========================================
	// MENYIMPAN KE DATABASE
	// ==========================================

	hasil, err := h.DB.Exec(`
		INSERT INTO modal_varian (
			id_varian,
			harga_modal,
			berlaku_mulai,
			catatan
		)
		VALUES (?, ?, ?, ?)
	`,
		req.IDVarian,
		req.HargaModal,
		req.BerlakuMulai,
		req.Catatan,
	)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menyimpan modal varian"})
		return
	}

	idBaru, _ := hasil.LastInsertId()

	// ==========================================
	// RESPONS BERHASIL
	// ==========================================

	kirimJSON(w, http.StatusCreated, ResponsAPI{
		Pesan: "Modal varian berhasil disimpan",
		Data: map[string]interface{}{
			"id_modal":      idBaru,
			"id_varian":     req.IDVarian,
			"nama_varian":   namaVarian,
			"harga_modal":   req.HargaModal,
			"berlaku_mulai": req.BerlakuMulai,
			"catatan":       req.Catatan,
		},
	})
}

// ==========================================
// AMBIL DAFTAR MODAL PER VARIAN
// ==========================================

// Mengambil modal terakhir untuk setiap varian.
func (h *Handler) ambilDaftarModal(w http.ResponseWriter, r *http.Request) {

	query := `
		SELECT
			m.id_modal,
			m.id_varian,
			v.nama_varian,
			m.harga_modal,
			m.berlaku_mulai,
			COALESCE(m.catatan, ''),
			m.dibuat_pada
		FROM modal_varian m
		INNER JOIN varian_produk v
			ON m.id_varian = v.id_varian
		INNER JOIN (
			SELECT id_varian, MAX(id_modal) AS id_modal_terakhir
			FROM modal_varian
			GROUP BY id_varian
		) terbaru
			ON m.id_modal = terbaru.id_modal_terakhir
		ORDER BY v.nama_varian ASC
	`

	var args []interface{}

	// Filter berdasarkan ID varian tertentu.
	if id := r.URL.Query().Get("id_varian"); id != "" {

		idVarian, err := strconv.ParseUint(id, 10, 64)

		if err != nil || idVarian == 0 {
			kirimJSON(w, http.StatusBadRequest,
				ResponsAPI{Pesan: "ID varian tidak valid"})
			return
		}

		query = `
			SELECT
				m.id_modal,
				m.id_varian,
				v.nama_varian,
				m.harga_modal,
				m.berlaku_mulai,
				COALESCE(m.catatan, ''),
				m.dibuat_pada
			FROM modal_varian m
			INNER JOIN varian_produk v
				ON m.id_varian = v.id_varian
			WHERE m.id_varian = ?
			ORDER BY m.dibuat_pada DESC
			LIMIT 50
		`

		args = append(args, idVarian)
	}

	rows, err := h.DB.Query(query, args...)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal mengambil data modal"})
		return
	}

	defer rows.Close()

	hasil := make([]ModalVarian, 0)

	for rows.Next() {

		var item ModalVarian

		err := rows.Scan(
			&item.IDModal,
			&item.IDVarian,
			&item.NamaVarian,
			&item.HargaModal,
			&item.BerlakuMulai,
			&item.Catatan,
			&item.DibuatPada,
		)

		if err != nil {
			kirimJSON(w, http.StatusInternalServerError,
				ResponsAPI{Pesan: "Gagal membaca data modal"})
			return
		}

		hasil = append(hasil, item)
	}

	if err := rows.Err(); err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Terjadi kesalahan saat membaca data"})
		return
	}

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Pesan: "Data modal varian berhasil diambil",
		Data:  hasil,
	})
}

// ==========================================
// TAMBAH PENGELUARAN OPERASIONAL
// ==========================================

// Mencatat pengeluaran operasional usaha baru.
func (h *Handler) tambahPengeluaran(w http.ResponseWriter, r *http.Request) {

	// Membatasi ukuran permintaan menjadi 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req PermintaanPengeluaran

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

	req.TanggalPengeluaran = strings.TrimSpace(req.TanggalPengeluaran)
	req.Kategori = strings.TrimSpace(req.Kategori)
	req.Keterangan = strings.TrimSpace(req.Keterangan)

	// ==========================================
	// VALIDASI INPUT
	// ==========================================

	if req.TanggalPengeluaran == "" {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Tanggal pengeluaran wajib diisi"})
		return
	}

	_, errTanggal := time.Parse("2006-01-02", req.TanggalPengeluaran)

	if errTanggal != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Format tanggal tidak valid, gunakan YYYY-MM-DD"})
		return
	}

	if len(req.Kategori) > 100 {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Kategori maksimal 100 karakter"})
		return
	}

	if len(req.Keterangan) > 255 {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Keterangan maksimal 255 karakter"})
		return
	}

	// ==========================================
	// WBS 1.1.6.4 - ANTI DOUBLE-COUNTING
	// Cek apakah kategori termasuk komponen HPP
	// dan cek duplikasi pengeluaran yang sudah ada.
	// ==========================================

	// Ambil daftar pengeluaran yang sudah tercatat
	// di bulan yang sama untuk cek duplikasi.
	bulanPengeluaran := req.TanggalPengeluaran[:7]

	rowsCek, errCek := h.DB.Query(`
		SELECT tanggal_pengeluaran, kategori, nominal
		FROM pengeluaran_usaha
		WHERE DATE_FORMAT(tanggal_pengeluaran, '%Y-%m') = ?
	`, bulanPengeluaran)

	if errCek != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal memeriksa data pengeluaran"})
		return
	}

	defer rowsCek.Close()

	var daftarSudahAda []validasibiaya.DataPengeluaranCek

	for rowsCek.Next() {
		var cek validasibiaya.DataPengeluaranCek
		if err := rowsCek.Scan(
			&cek.Tanggal, &cek.Kategori, &cek.Nominal,
		); err == nil {
			daftarSudahAda = append(daftarSudahAda, cek)
		}
	}

	// Jalankan validasi lengkap anti double-counting.
	hasilValidasi := validasibiaya.ValidasiPengeluaranLengkap(
		req.TanggalPengeluaran,
		req.Kategori,
		req.Nominal,
		daftarSudahAda,
	)

	if !hasilValidasi.Valid {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: hasilValidasi.Pesan})
		return
	}

	// ==========================================
	// MENYIMPAN KE DATABASE
	// ==========================================

	hasil, err := h.DB.Exec(`
		INSERT INTO pengeluaran_usaha (
			tanggal_pengeluaran,
			kategori,
			nominal,
			keterangan
		)
		VALUES (?, ?, ?, ?)
	`,
		req.TanggalPengeluaran,
		req.Kategori,
		req.Nominal,
		req.Keterangan,
	)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menyimpan pengeluaran"})
		return
	}

	idBaru, _ := hasil.LastInsertId()

	// ==========================================
	// RESPONS BERHASIL
	// ==========================================

	kirimJSON(w, http.StatusCreated, ResponsAPI{
		Pesan: "Pengeluaran berhasil dicatat",
		Data: map[string]interface{}{
			"id_pengeluaran":      idBaru,
			"tanggal_pengeluaran": req.TanggalPengeluaran,
			"kategori":            req.Kategori,
			"nominal":             req.Nominal,
			"keterangan":          req.Keterangan,
			"peringatan":          hasilValidasi.Peringatan,
		},
	})
}

// ==========================================
// AMBIL DAFTAR PENGELUARAN
// ==========================================

// Mengambil daftar pengeluaran dengan filter bulan.
func (h *Handler) ambilDaftarPengeluaran(w http.ResponseWriter, r *http.Request) {

	query := `
		SELECT
			id_pengeluaran,
			tanggal_pengeluaran,
			kategori,
			nominal,
			COALESCE(keterangan, ''),
			dibuat_pada,
			diperbarui_pada
		FROM pengeluaran_usaha
		WHERE 1=1
	`

	var args []interface{}

	// ==========================================
	// FILTER BERDASARKAN BULAN
	// ==========================================

	// Format bulan: YYYY-MM (contoh: 2026-11)
	if bulan := strings.TrimSpace(
		r.URL.Query().Get("bulan"),
	); bulan != "" {

		_, errBulan := time.Parse("2006-01", bulan)

		if errBulan != nil {
			kirimJSON(w, http.StatusBadRequest,
				ResponsAPI{Pesan: "Format bulan tidak valid, gunakan YYYY-MM"})
			return
		}

		query += " AND DATE_FORMAT(tanggal_pengeluaran, '%Y-%m') = ?"
		args = append(args, bulan)
	}

	// ==========================================
	// FILTER BERDASARKAN KATEGORI
	// ==========================================

	if kategori := strings.TrimSpace(
		r.URL.Query().Get("kategori"),
	); kategori != "" {

		query += " AND kategori = ?"
		args = append(args, kategori)
	}

	// Mengambil maksimal 200 pengeluaran terbaru.
	query += " ORDER BY tanggal_pengeluaran DESC, dibuat_pada DESC LIMIT 200"

	rows, err := h.DB.Query(query, args...)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal mengambil data pengeluaran"})
		return
	}

	defer rows.Close()

	hasil := make([]PengeluaranUsaha, 0)

	for rows.Next() {

		var item PengeluaranUsaha

		err := rows.Scan(
			&item.IDPengeluaran,
			&item.TanggalPengeluaran,
			&item.Kategori,
			&item.Nominal,
			&item.Keterangan,
			&item.DibuatPada,
			&item.DiperbaruiPada,
		)

		if err != nil {
			kirimJSON(w, http.StatusInternalServerError,
				ResponsAPI{Pesan: "Gagal membaca data pengeluaran"})
			return
		}

		hasil = append(hasil, item)
	}

	if err := rows.Err(); err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Terjadi kesalahan saat membaca data"})
		return
	}

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Pesan: "Data pengeluaran berhasil diambil",
		Data:  hasil,
	})
}

// ==========================================
// RINGKASAN BULANAN (LAPORAN)
// ==========================================

// Mengambil ringkasan modal dan pengeluaran per bulan.
// Digunakan sebagai sumber perhitungan laporan bulanan.
func (h *Handler) ambilRingkasanBulanan(w http.ResponseWriter, r *http.Request) {

	// Parameter bulan wajib ada (format: YYYY-MM).
	bulan := strings.TrimSpace(r.URL.Query().Get("bulan"))

	if bulan == "" {
		// Kalau kosong, pakai bulan sekarang.
		bulan = time.Now().Format("2006-01")
	}

	_, errBulan := time.Parse("2006-01", bulan)

	if errBulan != nil {
		kirimJSON(w, http.StatusBadRequest,
			ResponsAPI{Pesan: "Format bulan tidak valid, gunakan YYYY-MM"})
		return
	}

	// ==========================================
	// MENGHITUNG RATA-RATA MODAL
	// ==========================================

	// Mengambil modal terbaru untuk setiap varian,
	// lalu menghitung rata-ratanya.
	var rataRataModal sql.NullFloat64
	var jumlahVarian int

	err := h.DB.QueryRow(`
		SELECT
			AVG(m.harga_modal),
			COUNT(DISTINCT m.id_varian)
		FROM modal_varian m
		INNER JOIN (
			SELECT id_varian, MAX(id_modal) AS id_modal_terakhir
			FROM modal_varian
			WHERE berlaku_mulai <= LAST_DAY(CONCAT(?, '-01'))
			GROUP BY id_varian
		) terbaru
			ON m.id_modal = terbaru.id_modal_terakhir
	`, bulan).Scan(&rataRataModal, &jumlahVarian)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menghitung rata-rata modal"})
		return
	}

	// ==========================================
	// MENGHITUNG TOTAL PENGELUARAN
	// ==========================================

	var totalPengeluaran sql.NullFloat64
	var jumlahPosBiaya int

	err = h.DB.QueryRow(`
		SELECT
			COALESCE(SUM(nominal), 0),
			COUNT(*)
		FROM pengeluaran_usaha
		WHERE DATE_FORMAT(tanggal_pengeluaran, '%Y-%m') = ?
	`, bulan).Scan(&totalPengeluaran, &jumlahPosBiaya)

	if err != nil {
		kirimJSON(w, http.StatusInternalServerError,
			ResponsAPI{Pesan: "Gagal menghitung total pengeluaran"})
		return
	}

	// ==========================================
	// RESPONS RINGKASAN
	// ==========================================

	nilaiModal := 0.0
	if rataRataModal.Valid {
		nilaiModal = rataRataModal.Float64
	}

	nilaiPengeluaran := 0.0
	if totalPengeluaran.Valid {
		nilaiPengeluaran = totalPengeluaran.Float64
	}

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Pesan: "Ringkasan bulanan berhasil diambil",
		Data: RingkasanBulanan{
			Bulan:            bulan,
			RataRataModal:    nilaiModal,
			JumlahVarian:     jumlahVarian,
			TotalPengeluaran: nilaiPengeluaran,
			JumlahPosBiaya:   jumlahPosBiaya,
		},
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
