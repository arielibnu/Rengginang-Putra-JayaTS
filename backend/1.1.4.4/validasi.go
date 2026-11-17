
package validasistok

import (
	"errors"
	"strings"
)

const MaksimumStok uint64 = 4294967295

// Menormalisasi dan memvalidasi mode stok.
func NormalisasiMode(mode string) (string, error) {
	mode = strings.ToUpper(strings.TrimSpace(mode))

	if mode == "" {
		mode = "RIIL"
	}

	if mode != "RIIL" && mode != "SIMULASI" {
		return "", errors.New("Mode stok tidak valid")
	}

	return mode, nil
}

// Menormalisasi dan memvalidasi jenis mutasi.
func NormalisasiJenis(jenis string) (string, error) {
	jenis = strings.ToUpper(strings.TrimSpace(jenis))

	switch jenis {
	case "MASUK", "KELUAR", "PENYESUAIAN":
		return jenis, nil
	default:
		return "", errors.New("Jenis mutasi tidak valid")
	}
}

// Menghitung perubahan stok dengan perlindungan
// agar stok tidak negatif dan tidak melebihi batas.
func HitungPerubahan(
	jenis string,
	stokSaatIni uint64,
	jumlah uint64,
	stokBaru *uint64,
) (uint64, uint64, error) {

	if stokSaatIni > MaksimumStok {
		return 0, 0, errors.New("Stok saat ini tidak valid")
	}

	switch jenis {

	case "MASUK":

		if jumlah == 0 || jumlah > MaksimumStok-stokSaatIni {
			return 0, 0, errors.New("Jumlah stok masuk tidak valid")
		}

		return jumlah, stokSaatIni + jumlah, nil

	case "KELUAR":

		if jumlah == 0 || jumlah > stokSaatIni {
			return 0, 0, errors.New("Stok tidak mencukupi")
		}

		return jumlah, stokSaatIni - jumlah, nil

	case "PENYESUAIAN":

		if stokBaru == nil || *stokBaru > MaksimumStok {
			return 0, 0, errors.New("Stok baru tidak valid")
		}

		if *stokBaru == stokSaatIni {
			return 0, 0, errors.New("Tidak ada perubahan stok")
		}

		var perubahan uint64

		if *stokBaru > stokSaatIni {
			perubahan = *stokBaru - stokSaatIni
		} else {
			perubahan = stokSaatIni - *stokBaru
		}

		return perubahan, *stokBaru, nil

	default:
		return 0, 0, errors.New("Jenis mutasi tidak valid")
	}
}
