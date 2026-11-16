
package pelanggan

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

// Membaca sesi pelanggan dari cookie dan database.
func (h *Handler) ambilPelanggan(r *http.Request) (uint64, string, string, error) {
	cookie, err := r.Cookie(namaCookiePelanggan)
	if err != nil || len(cookie.Value) != 64 {
		return 0, "", "", sql.ErrNoRows
	}

	tokenHash := sha256.Sum256([]byte(cookie.Value))

	var id uint64
	var nama, email string

	err = h.DB.QueryRowContext(
		r.Context(),
		`SELECT p.id, p.nama, p.email
		 FROM sesi_pelanggan s
		 JOIN pelanggan p ON p.id = s.pelanggan_id
		 WHERE s.token_hash = ?
		   AND s.kedaluwarsa_pada > NOW()
		   AND p.aktif = TRUE`,
		hex.EncodeToString(tokenHash[:]),
	).Scan(&id, &nama, &email)

	return id, nama, email, err
}

// Profil menangani GET /api/pelanggan/profil.
func (h *Handler) Profil(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimPesan(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	if h.DB == nil {
		kirimPesan(w, http.StatusInternalServerError, "Layanan belum tersedia")
		return
	}

	id, nama, email, err := h.ambilPelanggan(r)
	if err != nil {
		kirimPesan(w, http.StatusUnauthorized, "Sesi pelanggan tidak valid")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"pesan": "Sesi pelanggan valid",
		"pelanggan": map[string]interface{}{
			"id":    id,
			"nama":  nama,
			"email": email,
		},
	})
}

// Logout menangani POST /api/pelanggan/keluar.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimPesan(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	if h.DB == nil {
		kirimPesan(w, http.StatusInternalServerError, "Layanan belum tersedia")
		return
	}

	cookie, err := r.Cookie(namaCookiePelanggan)
	if err == nil && len(cookie.Value) == 64 {
		tokenHash := sha256.Sum256([]byte(cookie.Value))

		_, err = h.DB.ExecContext(
			r.Context(),
			"DELETE FROM sesi_pelanggan WHERE token_hash = ?",
			hex.EncodeToString(tokenHash[:]),
		)
		if err != nil {
			kirimPesan(w, http.StatusInternalServerError, "Logout gagal")
			return
		}
	}

	// Menghapus cookie dari browser pelanggan.
	http.SetCookie(w, &http.Cookie{
		Name:     namaCookiePelanggan,
		Value:    "",
		Path:     "/api/pelanggan",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	kirimPesan(w, http.StatusOK, "Logout pelanggan berhasil")
}
