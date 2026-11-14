package main

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type dataPercobaan struct {
	jumlah         int
	mulai          time.Time
	diblokirSampai time.Time
}

var (
	muPercobaan     sync.Mutex
	daftarPercobaan = make(map[string]*dataPercobaan)
)

const (
	maksimalPercobaan = 5
	durasiPengamatan  = 10 * time.Minute
	durasiBlokir      = 10 * time.Minute
)

// Membatasi percobaan login berdasarkan alamat IP.
func middlewarePembatasLogin(berikutnya http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			kirimJSON(w, http.StatusBadRequest,
				map[string]string{"pesan": "Alamat pengguna tidak valid"})
			return
		}

		sekarang := time.Now()

		muPercobaan.Lock()

		data, ada := daftarPercobaan[ip]
		if !ada {
			data = &dataPercobaan{mulai: sekarang}
			daftarPercobaan[ip] = data
		}

		if sekarang.Before(data.diblokirSampai) {
			sisa := time.Until(data.diblokirSampai)
			muPercobaan.Unlock()

			w.Header().Set("Retry-After", strconv.Itoa(int(sisa.Seconds())+1))
			kirimJSON(w, http.StatusTooManyRequests,
				map[string]string{"pesan": "Terlalu banyak percobaan login. Coba lagi nanti."})
			return
		}

		if sekarang.Sub(data.mulai) >= durasiPengamatan {
			data.jumlah = 0
			data.mulai = sekarang
			data.diblokirSampai = time.Time{}
		}

		data.jumlah++

		if data.jumlah > maksimalPercobaan {
			data.diblokirSampai = sekarang.Add(durasiBlokir)
			muPercobaan.Unlock()

			w.Header().Set("Retry-After", strconv.Itoa(int(durasiBlokir.Seconds())))
			kirimJSON(w, http.StatusTooManyRequests,
				map[string]string{"pesan": "Batas percobaan login tercapai."})
			return
		}

		muPercobaan.Unlock()

		berikutnya.ServeHTTP(w, r)
	})
}
