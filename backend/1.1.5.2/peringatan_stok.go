package peringatanstok

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// BatasMaksimumKonfigurasi batas nilai wajar konfigurasi stok minimum.
const BatasMaksimumKonfigurasi uint64 = 1000000

// Handler mengelola operasi batas ambang stok dan aturan peringatan kritis.
type Handler struct {
	DB *sql.DB
}

// EvaluasiStatusStok menentukan status persediaan berdasarkan aturan bisnis PPL-05:
// "Peringatan muncul saat stok tersedia kurang dari atau sama dengan batas minimum per varian."
func EvaluasiStatusStok(stokTersedia int64, batasMinimum uint64) (string, bool) {
	if stokTersedia <= int64(batasMinimum) {
		return "Stok menipis", true
	}
	return "Normal", false
}

// ValidasiInputBatasMinimum memvalidasi integritas data permintaan pengaturan batas minimum.
func ValidasiInputBatasMinimum(req PermintaanAturBatasMinimum) error {
	if req.IDVarian == 0 {
		return errors.New("ID varian produk wajib diisi dan harus lebih dari 0")
	}

	if req.BatasMinimum > BatasMaksimumKonfigurasi {
		return fmt.Errorf("Batas minimum stok tidak boleh melebihi %d bungkus", BatasMaksimumKonfigurasi)
	}

	return nil
}

// SusunRingkasanPeringatan menghasilkan ringkasan komprehensif dari daftar status varian.
func SusunRingkasanPeringatan(daftar []StatusStokVarian) RingkasanPeringatanStok {
	var varianMenipis []StatusStokVarian

	for _, v := range daftar {
		if v.ApakahKritis {
			varianMenipis = append(varianMenipis, v)
		}
	}

	jumlahMenipis := len(varianMenipis)
	adaPeringatan := jumlahMenipis > 0

	var pesanUtama string
	var keteranganRinci string

	if adaPeringatan {
		pesanUtama = fmt.Sprintf("Perhatian: %d varian stok menipis", jumlahMenipis)
		var bagianKeterangan []string
		for _, item := range varianMenipis {
			bagianKeterangan = append(bagianKeterangan,
				fmt.Sprintf("%s tersedia %d bungkus, sedangkan batas minimum %d bungkus.",
					item.NamaVarian, item.StokTersedia, item.BatasMinimum))
		}
		keteranganRinci = strings.Join(bagianKeterangan, " ")
	} else {
		pesanUtama = "Semua persediaan aman"
		keteranganRinci = "Seluruh persediaan varian produk berada di atas batas minimum."
	}

	return RingkasanPeringatanStok{
		TotalVarianAktif:    len(daftar),
		JumlahVarianMenipis: jumlahMenipis,
		AdaPeringatan:       adaPeringatan,
		PesanUtama:          pesanUtama,
		KeteranganRinci:     keteranganRinci,
		DaftarVarian:        daftar,
	}
}

// HandlerBatasMinimum melayani GET, POST, dan PUT konfigurasi batas ambang minimum stok.
// Endpoint: /api/admin/stok/batas-minimum
func (h *Handler) HandlerBatasMinimum(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ambilBatasMinimum(w, r)
	case http.MethodPost, http.MethodPut:
		h.simpanBatasMinimum(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, PUT")
		kirimResponJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan",
		})
	}
}

// HandlerPeringatanStok melayani GET evaluasi status stok dan peringatan kritis varian.
// Endpoint: /api/admin/stok/peringatan
func (h *Handler) HandlerPeringatanStok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimResponJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan",
		})
		return
	}

	daftar, err := h.ambilDaftarStatusStokDariDB()
	if err != nil {
		kirimResponJSON(w, http.StatusInternalServerError, ResponsAPI{
			Status: "error",
			Pesan:  "Gagal mengambil evaluasi peringatan stok: " + err.Error(),
		})
		return
	}

	ringkasan := SusunRingkasanPeringatan(daftar)

	kirimResponJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Berhasil memuat peringatan stok",
		Data:   ringkasan,
	})
}

// ambilBatasMinimum mengambil konfigurasi batas minimum varian produk dari basis data.
func (h *Handler) ambilBatasMinimum(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id_varian")

	// Ambil 1 varian spesifik jika query parameter id_varian disertakan
	if idParam != "" {
		idVarian, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil || idVarian == 0 {
			kirimResponJSON(w, http.StatusBadRequest, ResponsAPI{
				Status: "gagal",
				Pesan:  "Parameter id_varian tidak valid",
			})
			return
		}

		item, err := h.ambilBatasMinimumSatu(idVarian)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				kirimResponJSON(w, http.StatusNotFound, ResponsAPI{
					Status: "gagal",
					Pesan:  "Varian produk tidak ditemukan",
				})
				return
			}
			kirimResponJSON(w, http.StatusInternalServerError, ResponsAPI{
				Status: "error",
				Pesan:  "Gagal membaca konfigurasi batas minimum: " + err.Error(),
			})
			return
		}

		kirimResponJSON(w, http.StatusOK, ResponsAPI{
			Status: "sukses",
			Pesan:  "Berhasil mengambil batas minimum varian",
			Data:   item,
		})
		return
	}

	// Ambil seluruh konfigurasi varian
	daftar, err := h.ambilSemuaBatasMinimum()
	if err != nil {
		kirimResponJSON(w, http.StatusInternalServerError, ResponsAPI{
			Status: "error",
			Pesan:  "Gagal mengambil daftar batas minimum: " + err.Error(),
		})
		return
	}

	kirimResponJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Berhasil mengambil seluruh batas minimum varian",
		Data:   daftar,
	})
}

// simpanBatasMinimum memproses penyimpanan atau pembaruan batas minimum stok varian.
func (h *Handler) simpanBatasMinimum(w http.ResponseWriter, r *http.Request) {
	var req PermintaanAturBatasMinimum

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimResponJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  "Format payload JSON tidak valid",
		})
		return
	}

	if err := ValidasiInputBatasMinimum(req); err != nil {
		kirimResponJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  err.Error(),
		})
		return
	}

	if err := h.perbaruiBatasMinimumDB(req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			kirimResponJSON(w, http.StatusNotFound, ResponsAPI{
				Status: "gagal",
				Pesan:  "Varian produk tidak ditemukan",
			})
			return
		}
		kirimResponJSON(w, http.StatusInternalServerError, ResponsAPI{
			Status: "error",
			Pesan:  "Gagal menyimpan batas minimum: " + err.Error(),
		})
		return
	}

	// Ambil data status terkini setelah diperbarui
	statusTerkini, _ := h.ambilStatusStokSatu(req.IDVarian)

	kirimResponJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Batas minimum stok berhasil diperbarui",
		Data:   statusTerkini,
	})
}

// ====================================================================
// OPERASI BASIS DATA MYSQL
// ====================================================================

func (h *Handler) ambilBatasMinimumSatu(idVarian uint64) (*KonfigurasiBatasMinimum, error) {
	query := `
		SELECT 
			v.id_varian,
			v.nama_varian,
			COALESCE(b.batas_minimum, 15) AS batas_minimum,
			COALESCE(b.catatan, '') AS catatan,
			COALESCE(b.diperbarui_pada, v.diperbarui_pada) AS diperbarui_pada
		FROM varian_produk v
		LEFT JOIN batas_minimum_stok b ON v.id_varian = b.id_varian
		WHERE v.id_varian = ?
		LIMIT 1
	`
	var k KonfigurasiBatasMinimum
	err := h.DB.QueryRow(query, idVarian).Scan(
		&k.IDVarian,
		&k.NamaVarian,
		&k.BatasMinimum,
		&k.Catatan,
		&k.DiperbaruiPada,
	)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (h *Handler) ambilSemuaBatasMinimum() ([]KonfigurasiBatasMinimum, error) {
	query := `
		SELECT 
			v.id_varian,
			v.nama_varian,
			COALESCE(b.batas_minimum, 15) AS batas_minimum,
			COALESCE(b.catatan, '') AS catatan,
			COALESCE(b.diperbarui_pada, v.diperbarui_pada) AS diperbarui_pada
		FROM varian_produk v
		LEFT JOIN batas_minimum_stok b ON v.id_varian = b.id_varian
		WHERE v.status_aktif = TRUE
		ORDER BY v.id_varian ASC
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hasil []KonfigurasiBatasMinimum
	for rows.Next() {
		var k KonfigurasiBatasMinimum
		if err := rows.Scan(
			&k.IDVarian,
			&k.NamaVarian,
			&k.BatasMinimum,
			&k.Catatan,
			&k.DiperbaruiPada,
		); err != nil {
			return nil, err
		}
		hasil = append(hasil, k)
	}
	return hasil, nil
}

func (h *Handler) perbaruiBatasMinimumDB(req PermintaanAturBatasMinimum) error {
	// Pastikan varian ada di basis data
	var ada bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM varian_produk WHERE id_varian = ?)", req.IDVarian).Scan(&ada)
	if err != nil {
		return err
	}
	if !ada {
		return sql.ErrNoRows
	}

	// Simpan ke tabel batas_minimum_stok (Upsert)
	upsertQuery := `
		INSERT INTO batas_minimum_stok (id_varian, batas_minimum, catatan)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			batas_minimum = VALUES(batas_minimum),
			catatan = VALUES(catatan),
			diperbarui_pada = CURRENT_TIMESTAMP
	`
	_, err = h.DB.Exec(upsertQuery, req.IDVarian, req.BatasMinimum, req.Catatan)
	if err != nil {
		// Fallback: perbarui kolom varian_produk jika tabel terpisah belum dimigrasi
		_, errUpdate := h.DB.Exec("UPDATE varian_produk SET batas_minimum = ? WHERE id_varian = ?", req.BatasMinimum, req.IDVarian)
		if errUpdate != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) ambilDaftarStatusStokDariDB() ([]StatusStokVarian, error) {
	query := `
		SELECT 
			v.id_varian,
			v.nama_varian,
			v.stok AS stok_fisik,
			0 AS reservasi,
			COALESCE(b.batas_minimum, 15) AS batas_minimum
		FROM varian_produk v
		LEFT JOIN batas_minimum_stok b ON v.id_varian = b.id_varian
		WHERE v.status_aktif = TRUE
		ORDER BY v.id_varian ASC
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var daftar []StatusStokVarian
	for rows.Next() {
		var s StatusStokVarian
		if err := rows.Scan(
			&s.IDVarian,
			&s.NamaVarian,
			&s.StokFisik,
			&s.Reservasi,
			&s.BatasMinimum,
		); err != nil {
			return nil, err
		}

		s.StokTersedia = int64(s.StokFisik) - int64(s.Reservasi)
		s.Status, s.ApakahKritis = EvaluasiStatusStok(s.StokTersedia, s.BatasMinimum)
		daftar = append(daftar, s)
	}

	return daftar, nil
}

func (h *Handler) ambilStatusStokSatu(idVarian uint64) (*StatusStokVarian, error) {
	query := `
		SELECT 
			v.id_varian,
			v.nama_varian,
			v.stok AS stok_fisik,
			0 AS reservasi,
			COALESCE(b.batas_minimum, 15) AS batas_minimum
		FROM varian_produk v
		LEFT JOIN batas_minimum_stok b ON v.id_varian = b.id_varian
		WHERE v.id_varian = ?
		LIMIT 1
	`
	var s StatusStokVarian
	err := h.DB.QueryRow(query, idVarian).Scan(
		&s.IDVarian,
		&s.NamaVarian,
		&s.StokFisik,
		&s.Reservasi,
		&s.BatasMinimum,
	)
	if err != nil {
		return nil, err
	}

	s.StokTersedia = int64(s.StokFisik) - int64(s.Reservasi)
	s.Status, s.ApakahKritis = EvaluasiStatusStok(s.StokTersedia, s.BatasMinimum)
	return &s, nil
}

// kirimResponJSON pembantu serialisasi respons HTTP JSON.
func kirimResponJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
