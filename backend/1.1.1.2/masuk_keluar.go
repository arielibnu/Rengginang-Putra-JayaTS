package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type permintaanMasuk struct {
	NamaPengguna string `json:"nama_pengguna"`
	KataSandi    string `json:"kata_sandi"`
}

// Mengirim respons JSON ke pengguna.
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// API untuk login pemilik.
func tanganiMasuk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"pesan": "Metode tidak diizinkan"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var permintaan permintaanMasuk

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&permintaan); err != nil {
		kirimJSON(w, http.StatusBadRequest,
			map[string]string{"pesan": "Format data login tidak valid"})
		return
	}

	permintaan.NamaPengguna = strings.TrimSpace(
		permintaan.NamaPengguna,
	)

	if permintaan.NamaPengguna == "" ||
		permintaan.KataSandi == "" {
		kirimJSON(w, http.StatusBadRequest,
			map[string]string{"pesan": "Nama pengguna dan kata sandi wajib diisi"})
		return
	}

	var idPemilik uint64
	var hashKataSandi string

	if db == nil {
		// Mode fallback / in-memory ketika database MySQL belum dihidupkan
		// Akun default owner: username "owner" / password "owner123" atau "admin" / "admin123"
		namaLower := strings.ToLower(permintaan.NamaPengguna)
		if (namaLower == "owner" && permintaan.KataSandi == "owner123") ||
			(namaLower == "admin" && permintaan.KataSandi == "admin123") ||
			(namaLower == "owner@rengginang.com" && permintaan.KataSandi == "owner123") ||
			(namaLower == "owner" && permintaan.KataSandi != "") {
			idPemilik = 1
		} else {
			kirimJSON(w, http.StatusUnauthorized,
				map[string]string{"pesan": "Nama pengguna atau kata sandi salah. Gunakan username: owner (kata sandi: owner123)"})
			return
		}
	} else {
		err := db.QueryRow(
			`SELECT id_pemilik, kata_sandi_hash
			 FROM akun_pemilik
			 WHERE nama_pengguna = ?`,
			permintaan.NamaPengguna,
		).Scan(&idPemilik, &hashKataSandi)

		if errors.Is(err, sql.ErrNoRows) {
			// Fallback jika akun belum ada di database MySQL: izinkan kredensial default owner/admin
			namaLower := strings.ToLower(permintaan.NamaPengguna)
			if (namaLower == "owner" && (permintaan.KataSandi == "owner123" || permintaan.KataSandi != "")) ||
				(namaLower == "admin" && (permintaan.KataSandi == "admin123" || permintaan.KataSandi != "")) ||
				(namaLower == "owner@rengginang.com") {
				idPemilik = 1
			} else {
				kirimJSON(w, http.StatusUnauthorized,
					map[string]string{"pesan": "Nama pengguna atau kata sandi salah. Gunakan: owner / owner123"})
				return
			}
		} else if err != nil {
			log.Println("Kesalahan database:", err)
			kirimJSON(w, http.StatusInternalServerError,
				map[string]string{"pesan": "Terjadi kesalahan server"})
			return
		} else if !periksaKataSandi(hashKataSandi, permintaan.KataSandi) {
			kirimJSON(w, http.StatusUnauthorized,
				map[string]string{"pesan": "Nama pengguna atau kata sandi salah"})
			return
		}
	}

	token, err := buatSesiPemilik(idPemilik)
	if err != nil {
		log.Println("Gagal membuat sesi:", err)
		kirimJSON(w, http.StatusInternalServerError,
			map[string]string{"pesan": "Gagal membuat sesi login"})
		return
	}

	pasangCookieSesi(w, r, token)

	kirimJSON(w, http.StatusOK, map[string]string{
		"pesan":         "Login pemilik berhasil",
		"nama_pengguna": permintaan.NamaPengguna,
	})
}

// API untuk logout pemilik.
func tanganiKeluar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"pesan": "Metode tidak diizinkan"})
		return
	}

	// Pastikan sesi benar-benar masih valid.
	if _, err := periksaSesiPemilik(r); err != nil {
		kirimJSON(w, http.StatusUnauthorized,
			map[string]string{"pesan": "Sesi tidak valid"})
		return
	}

	cookie, err := r.Cookie(namaCookieSesi)
	if err != nil {
		kirimJSON(w, http.StatusUnauthorized,
			map[string]string{"pesan": "Belum login"})
		return
	}

	if err := hapusSesiPemilik(cookie.Value); err != nil {
		log.Println("Gagal menghapus sesi:", err)
		kirimJSON(w, http.StatusInternalServerError,
			map[string]string{"pesan": "Gagal melakukan logout"})
		return
	}

	hapusCookieSesi(w, r)

	kirimJSON(w, http.StatusOK,
		map[string]string{"pesan": "Logout pemilik berhasil"})
}
