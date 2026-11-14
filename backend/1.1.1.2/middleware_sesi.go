package main

import (
	"context"
	"log"
	"net/http"
)

type kunciKonteks string

const kunciIDPemilik kunciKonteks = "id_pemilik"

// Memeriksa sesi sebelum mengizinkan akses.
func middlewareSesiPemilik(
	berikutnya http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		idPemilik, err := periksaSesiPemilik(r)

		if err != nil {
			log.Println("Sesi ditolak:", err)

			kirimJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{
					"pesan": "Akses ditolak. Silakan login.",
				},
			)
			return
		}

		konteks := context.WithValue(
			r.Context(),
			kunciIDPemilik,
			idPemilik,
		)

		berikutnya.ServeHTTP(
			w,
			r.WithContext(konteks),
		)
	})
}

// Mengambil ID pemilik dari konteks permintaan.
func ambilIDPemilik(r *http.Request) (uint64, bool) {
	id, ada := r.Context().Value(kunciIDPemilik).(uint64)
	return id, ada
}
