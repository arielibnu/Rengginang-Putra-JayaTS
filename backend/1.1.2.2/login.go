
package pelanggan

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const namaCookiePelanggan = "sesi_pelanggan"

type inputLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login menangani POST /api/pelanggan/masuk.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimPesan(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	if h.DB == nil {
		kirimPesan(w, http.StatusInternalServerError, "Layanan login belum tersedia")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input inputLogin
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		kirimPesan(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	var tambahan interface{}
	if err := decoder.Decode(&tambahan); err != io.EOF {
		kirimPesan(w, http.StatusBadRequest, "JSON harus berisi satu objek")
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Email == "" || input.Password == "" {
		kirimPesan(w, http.StatusBadRequest, "Email dan password wajib diisi")
		return
	}

	var (
		id           uint64
		nama         string
		passwordHash string
		aktif        bool
	)

	err := h.DB.QueryRowContext(
		r.Context(),
		`SELECT id, nama, password_hash, aktif
		 FROM pelanggan WHERE email = ?`,
		input.Email,
	).Scan(&id, &nama, &passwordHash, &aktif)

	if err != nil {
		if err == sql.ErrNoRows {
			kirimPesan(w, http.StatusUnauthorized, "Email atau password salah")
			return
		}
		kirimPesan(w, http.StatusInternalServerError, "Gagal memproses login")
		return
	}

	if !aktif || bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(input.Password),
	) != nil {
		kirimPesan(w, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Token sesi acak 32 byte.
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		kirimPesan(w, http.StatusInternalServerError, "Gagal membuat sesi")
		return
	}

	token := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))
	expiredAt := time.Now().Add(24 * time.Hour)

	_, err = h.DB.ExecContext(
		r.Context(),
		`INSERT INTO sesi_pelanggan
		 (pelanggan_id, token_hash, kedaluwarsa_pada)
		 VALUES (?, ?, ?)`,
		id,
		hex.EncodeToString(tokenHash[:]),
		expiredAt,
	)
	if err != nil {
		kirimPesan(w, http.StatusInternalServerError, "Gagal menyimpan sesi")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     namaCookiePelanggan,
		Value:    token,
		Path:     "/api/pelanggan",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
		Expires:  expiredAt,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"pesan": "Login pelanggan berhasil",
		"pelanggan": map[string]interface{}{
			"id":    id,
			"nama":  nama,
			"email": input.Email,
		},
	})
}