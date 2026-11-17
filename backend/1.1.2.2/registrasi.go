
package pelanggan

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Handler menggunakan koneksi database MySQL dari backend utama.
type Handler struct {
	DB *sql.DB
}

// Data yang dikirim frontend saat pelanggan mendaftar.
type inputRegistrasi struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Registrasi menangani POST /api/pelanggan/daftar.
func (h *Handler) Registrasi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Registrasi hanya menerima metode POST.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimPesan(
			w,
			http.StatusMethodNotAllowed,
			"Metode tidak diizinkan",
		)
		return
	}

	// Pastikan koneksi database tersedia.
	if h.DB == nil {
		kirimPesan(
			w,
			http.StatusInternalServerError,
			"Layanan registrasi belum tersedia",
		)
		return
	}

	// Batasi ukuran request maksimal 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input inputRegistrasi

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Membaca data JSON dari frontend.
	if err := decoder.Decode(&input); err != nil {
		kirimPesan(
			w,
			http.StatusBadRequest,
			"Format JSON tidak valid",
		)
		return
	}

	// Pastikan request hanya berisi satu objek JSON.
	var tambahan interface{}

	if err := decoder.Decode(&tambahan); err != io.EOF {
		kirimPesan(
			w,
			http.StatusBadRequest,
			"JSON harus berisi satu objek saja",
		)
		return
	}

	// Membersihkan input nama dan email.
	input.Nama = strings.TrimSpace(input.Nama)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// Validasi nama pelanggan.
	if input.Nama == "" ||
		utf8.RuneCountInString(input.Nama) > 150 {

		kirimPesan(
			w,
			http.StatusBadRequest,
			"Nama pelanggan harus diisi dan maksimal 150 karakter",
		)
		return
	}

	// Validasi format email.
	alamat, err := mail.ParseAddress(input.Email)

	if err != nil ||
		alamat.Address != input.Email ||
		len(input.Email) > 191 ||
		strings.ContainsAny(input.Email, " \t\r\n") {

		kirimPesan(
			w,
			http.StatusBadRequest,
			"Format email tidak valid",
		)
		return
	}

	// bcrypt menerima password maksimal 72 byte.
	if len(input.Password) < 8 || len(input.Password) > 72 {
		kirimPesan(
			w,
			http.StatusBadRequest,
			"Password harus berisi 8 sampai 72 byte",
		)
		return
	}

	// Mengubah password menjadi hash bcrypt.
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		kirimPesan(
			w,
			http.StatusInternalServerError,
			"Gagal memproses registrasi",
		)
		return
	}

	// Menyimpan pelanggan ke database MySQL.
	hasil, err := h.DB.ExecContext(
		r.Context(),
		`INSERT INTO pelanggan (nama, email, password_hash)
		 VALUES (?, ?, ?)`,
		input.Nama,
		input.Email,
		string(hash),
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError

		// Error MySQL 1062 = email sudah terdaftar.
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			kirimPesan(
				w,
				http.StatusConflict,
				"Email sudah digunakan",
			)
			return
		}

		// Error database lain tidak ditampilkan ke pengguna.
		kirimPesan(
			w,
			http.StatusInternalServerError,
			"Registrasi gagal, silakan coba lagi",
		)
		return
	}

	// Mengambil ID pelanggan baru.
	id, err := hasil.LastInsertId()

	if err != nil {
		kirimPesan(
			w,
			http.StatusInternalServerError,
			"Akun tersimpan, tetapi gagal membaca ID pelanggan",
		)
		return
	}

	// Registrasi berhasil.
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		map[string]interface{}{
			"pesan": "Registrasi pelanggan berhasil",
			"id":    id,
		},
	)
}

// Helper untuk mengirim respons JSON yang konsisten.
// Bisa dipakai kembali untuk Login dan Logout PPL-02.
func kirimPesan(
	w http.ResponseWriter,
	status int,
	pesan string,
) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"pesan": pesan,
		},
	)
}
