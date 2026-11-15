package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const namaCookieSesi = "sesi_pemilik"

const durasiSesi = 8 * time.Hour

// Membuat token sesi acak yang aman.
func buatTokenSesi() (string, error) {
	data := make([]byte, 32)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}

// Mengubah token menjadi hash SHA-256.
func hashTokenSesi(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Menyimpan sesi baru ke database.
func buatSesiPemilik(idPemilik uint64) (string, error) {
	token, err := buatTokenSesi()
	if err != nil {
		return "", err
	}

	kedaluwarsa := time.Now().Add(durasiSesi)

	_, err = db.Exec(
		`INSERT INTO sesi_pemilik
		(id_pemilik, token_hash, kedaluwarsa_pada)
		VALUES (?, ?, ?)`,
		idPemilik,
		hashTokenSesi(token),
		kedaluwarsa,
	)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan sesi: %w", err)
	}

	return token, nil
}

// Mengirim cookie sesi ke browser.
func pasangCookieSesi(
	w http.ResponseWriter,
	r *http.Request,
	token string,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     namaCookieSesi,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(durasiSesi.Seconds()),
	})
}

// Mengambil dan memvalidasi sesi Owner.
func periksaSesiPemilik(r *http.Request) (uint64, error) {
	cookie, err := r.Cookie(namaCookieSesi)
	if err != nil || cookie.Value == "" {
		return 0, errors.New("sesi tidak ditemukan")
	}

	var idPemilik uint64

	err = db.QueryRow(
		`SELECT id_pemilik
		FROM sesi_pemilik
		WHERE token_hash = ?
		AND kedaluwarsa_pada > NOW()`,
		hashTokenSesi(cookie.Value),
	).Scan(&idPemilik)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("sesi tidak valid atau kedaluwarsa")
	}

	if err != nil {
		return 0, err
	}

	return idPemilik, nil
}

// Menghapus sesi dari database saat logout.
func hapusSesiPemilik(token string) error {
	_, err := db.Exec(
		`DELETE FROM sesi_pemilik WHERE token_hash = ?`,
		hashTokenSesi(token),
	)
	return err
}

// Menghapus cookie sesi dari browser.
func hapusCookieSesi(
	w http.ResponseWriter,
	r *http.Request,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     namaCookieSesi,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
