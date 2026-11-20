package validasiwilayah

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// ============================================================================
// WBS 1.2.6.3 - LOGIKA VALIDASI GEOGRAFIS PENGIRIMAN API KURIR
// ============================================================================

// HandlerValidasiWilayah menangani request verifikasi kelengkapan wilayah dan alamat.
type HandlerValidasiWilayah struct {
	DB *sql.DB
}

// Master data asal pengiriman resmi toko Rengginang Putra Jaya
var AsalTokoRengginang = LokasiGeografis{
	Provinsi:      "Jawa Timur",
	KotaKabupaten: "Kabupaten Situbondo",
	Kecamatan:     "Kecamatan Situbondo",
	KodePos:       "68312",
	AlamatLengkap: "Jl. Raya Situbondo No. 45, Sentra UMKM Rengginang, Situbondo",
	KodeWilayahID: "ID-JI-SIT-01",
}

// MasterWilayahIndonesia master data hierarki wilayah bertingkat (Provinsi -> Kota/Kabupaten -> Kecamatan)
var MasterWilayahIndonesia = map[string]map[string][]string{
	"Jawa Timur": {
		"Kabupaten Sidoarjo": {
			"Kecamatan Sidoarjo",
			"Kecamatan Gedangan",
			"Kecamatan Waru",
			"Kecamatan Candi",
			"Kecamatan Taman",
		},
		"Kabupaten Situbondo": {
			"Kecamatan Situbondo",
			"Kecamatan Panji",
			"Kecamatan Besuki",
			"Kecamatan Asembagus",
		},
		"Kota Surabaya": {
			"Kecamatan Wonokromo",
			"Kecamatan Gubeng",
			"Kecamatan Tegalsari",
			"Kecamatan Rungkut",
		},
		"Kabupaten Malang": {
			"Kecamatan Kepanjen",
			"Kecamatan Singosari",
			"Kecamatan Lawang",
		},
	},
	"DKI Jakarta": {
		"Kota Jakarta Selatan": {
			"Kecamatan Kebayoran Baru",
			"Kecamatan Cilandak",
			"Kecamatan Pasar Minggu",
		},
		"Kota Jakarta Pusat": {
			"Kecamatan Gambir",
			"Kecamatan Menteng",
			"Kecamatan Tanah Abang",
		},
	},
	"Jawa Barat": {
		"Kota Bandung": {
			"Kecamatan Coblong",
			"Kecamatan Sukasari",
		},
		"Kota Bekasi": {
			"Kecamatan Bekasi Barat",
			"Kecamatan Bekasi Timur",
		},
	},
	"Jawa Tengah": {
		"Kota Semarang": {
			"Kecamatan Semarang Barat",
			"Kecamatan Semarang Selatan",
		},
		"Kota Surakarta": {
			"Kecamatan Banjarsari",
			"Kecamatan Laweyan",
		},
	},
}

var regexNomorTeleponIndo = regexp.MustCompile(`^(\+62|62|08)[0-9]{8,13}$`)
var regexKodePos = regexp.MustCompile(`^[0-9]{5}$`)

// FormatNomorE164 menstandarisasi nomor telepon ke format internasional E.164 (+62...)
func FormatNomorE164(no string) string {
	clean := strings.TrimSpace(no)
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")

	if strings.HasPrefix(clean, "+62") {
		return clean
	}
	if strings.HasPrefix(clean, "62") {
		return "+" + clean
	}
	if strings.HasPrefix(clean, "08") {
		return "+62" + strings.TrimPrefix(clean, "0")
	}
	return clean
}

// ValidasiParameterKurir melakukan pemeriksaan menyeluruh terhadap seluruh parameter geografis dan kontak
// yang dibutuhkan oleh API Kurir ekspedisi (RajaOngkir, Biteship, Komship, JNE, SiCepat, dll).
func ValidasiParameterKurir(req *PermintaanWilayahPengiriman) (bool, map[string]string, *ParameterKurirStandar) {
	errors := make(map[string]string)

	if req == nil {
		errors["payload"] = "Payload permintaan tidak boleh kosong"
		return false, errors, nil
	}

	// 1. Validasi Kontak Penerima (Dibutuhkan kurir untuk konfirmasi serah terima)
	namaClean := strings.TrimSpace(req.Kontak.NamaPenerima)
	if len(namaClean) < 2 {
		errors["nama_penerima"] = "Nama penerima wajib diisi (minimal 2 karakter)"
	}

	waClean := strings.TrimSpace(req.Kontak.NomorWhatsApp)
	waClean = strings.ReplaceAll(waClean, "-", "")
	waClean = strings.ReplaceAll(waClean, " ", "")
	if waClean == "" {
		errors["nomor_whatsapp"] = "Nomor kontak penerima (WhatsApp/Telepon) wajib diisi"
	} else if !regexNomorTeleponIndo.MatchString(waClean) {
		errors["nomor_whatsapp"] = "Nomor WhatsApp tidak valid (harus nomor Indonesia dengan format 08xx/628xx, 10-15 digit)"
	}

	// 2. Pemeriksaan Metode Penerimaan
	metode := strings.ToLower(strings.TrimSpace(req.MetodePenerimaan))
	if metode == "" {
		metode = "kurir"
	}

	// Jika pembeli memilih ambil sendiri di toko (PPL-14)
	if metode == "ambil_sendiri" {
		if len(errors) > 0 {
			return false, errors, nil
		}
		paramKurir := &ParameterKurirStandar{
			AsalPengiriman:   AsalTokoRengginang,
			TujuanPengiriman: AsalTokoRengginang, // Ambil di toko
			Penerima: PenerimaKurir{
				Nama:           namaClean,
				Telepon:        waClean,
				TeleponStandar: FormatNomorE164(waClean),
			},
			Paket: PaketKurir{
				BeratGram: 1350,
				Kategori:  "Makanan Ringan",
				IsFragile: true,
			},
			IsKurirReady:    false,
			StatusGeografis: "NOT_REQUIRED",
		}
		return true, nil, paramKurir
	}

	// 3. Validasi Geografis Wilayah Pengiriman Kurir
	if req.WilayahPengiriman == nil {
		errors["wilayah_pengiriman"] = "Data wilayah pengiriman wajib diisi untuk pengiriman via kurir"
		return false, errors, nil
	}

	provinsi := strings.TrimSpace(req.WilayahPengiriman.Provinsi)
	if provinsi == "" {
		errors["provinsi"] = "Provinsi tujuan wajib dipilih"
	}

	kota := strings.TrimSpace(req.WilayahPengiriman.KotaKabupaten)
	if kota == "" {
		errors["kota_kabupaten"] = "Kota/Kabupaten tujuan wajib dipilih"
	}

	kecamatan := strings.TrimSpace(req.WilayahPengiriman.Kecamatan)
	if kecamatan == "" {
		errors["kecamatan"] = "Kecamatan tujuan wajib dipilih"
	}

	kodePos := strings.TrimSpace(req.WilayahPengiriman.KodePos)
	if kodePos != "" && !regexKodePos.MatchString(kodePos) {
		errors["kode_pos"] = "Kode pos harus terdiri dari tepat 5 digit angka"
	}

	alamat := strings.TrimSpace(req.WilayahPengiriman.AlamatLengkap)
	if alamat == "" {
		errors["alamat_lengkap"] = "Alamat lengkap tujuan pengiriman wajib diisi"
	} else if len(alamat) < 10 {
		errors["alamat_lengkap"] = "Alamat lengkap terlalu singkat (minimal 10 karakter agar dapat dijangkau kurir)"
	}

	// 4. Validasi Hierarki Wilayah Bertingkat (jika ada pada Master Data)
	if provinsi != "" && kota != "" {
		if kotaMap, provAda := MasterWilayahIndonesia[provinsi]; provAda {
			if kecList, kotaAda := kotaMap[kota]; kotaAda {
				if kecamatan != "" {
					kecCocok := false
					for _, k := range kecList {
						if strings.EqualFold(k, kecamatan) {
							kecCocok = true
							break
						}
					}
					if !kecCocok {
						errors["kecamatan"] = fmt.Sprintf("Kecamatan '%s' tidak ditemukan dalam wilayah '%s'", kecamatan, kota)
					}
				}
			} else {
				errors["kota_kabupaten"] = fmt.Sprintf("Kota/Kabupaten '%s' tidak berada di provinsi '%s'", kota, provinsi)
			}
		}
	}

	if len(errors) > 0 {
		return false, errors, nil
	}

	// 5. Kalkulasi Berat Paket (jika belum ditentukan, gunakan default 1350gr sesuai item checkout 2 udang + 1 bawang)
	beratPaket := req.TotalBeratGram
	if beratPaket == 0 {
		beratPaket = 1350
	}

	paramKurir := &ParameterKurirStandar{
		AsalPengiriman: AsalTokoRengginang,
		TujuanPengiriman: LokasiGeografis{
			Provinsi:      provinsi,
			KotaKabupaten: kota,
			Kecamatan:     kecamatan,
			KodePos:       kodePos,
			AlamatLengkap: alamat,
			KodeWilayahID: fmt.Sprintf("DEST-%s-%s", strings.ToUpper(provinsi[:2]), kodePos),
		},
		Penerima: PenerimaKurir{
			Nama:           namaClean,
			Telepon:        waClean,
			TeleponStandar: FormatNomorE164(waClean),
		},
		Paket: PaketKurir{
			BeratGram: beratPaket,
			Kategori:  "Makanan Ringan",
			IsFragile: true,
		},
		IsKurirReady:    true,
		StatusGeografis: "VALID",
	}

	return true, nil, paramKurir
}

// SimpanAlamatDatabase mencatat data alamat pengiriman ke tabel MySQL jika koneksi aktif.
func (h *HandlerValidasiWilayah) SimpanAlamatDatabase(req *PermintaanWilayahPengiriman, param *ParameterKurirStandar) error {
	if h.DB == nil || req == nil {
		return nil
	}

	// Buat tabel jika belum ada
	_, _ = h.DB.Exec(`
		CREATE TABLE IF NOT EXISTS alamat_pengiriman_pesanan (
			id_alamat BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			id_pesanan VARCHAR(64) NULL,
			nama_penerima VARCHAR(150) NOT NULL,
			nomor_whatsapp VARCHAR(25) NOT NULL,
			metode_penerimaan VARCHAR(20) NOT NULL DEFAULT 'kurir',
			provinsi VARCHAR(100) NOT NULL DEFAULT '',
			kota_kabupaten VARCHAR(100) NOT NULL DEFAULT '',
			kecamatan VARCHAR(100) NOT NULL DEFAULT '',
			kode_pos VARCHAR(10) NOT NULL DEFAULT '',
			alamat_lengkap TEXT NOT NULL,
			catatan_alamat TEXT NULL,
			kurir_ready BOOLEAN NOT NULL DEFAULT TRUE,
			waktu_dibuat DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	query := `
		INSERT INTO alamat_pengiriman_pesanan (
			nama_penerima, nomor_whatsapp, metode_penerimaan,
			provinsi, kota_kabupaten, kecamatan, kode_pos, alamat_lengkap, kurir_ready
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	prov := ""
	kota := ""
	kec := ""
	pos := ""
	alm := ""
	ready := false

	if req.WilayahPengiriman != nil {
		prov = req.WilayahPengiriman.Provinsi
		kota = req.WilayahPengiriman.KotaKabupaten
		kec = req.WilayahPengiriman.Kecamatan
		pos = req.WilayahPengiriman.KodePos
		alm = req.WilayahPengiriman.AlamatLengkap
	}
	if param != nil {
		ready = param.IsKurirReady
	}

	_, err := h.DB.Exec(query,
		req.Kontak.NamaPenerima,
		req.Kontak.NomorWhatsApp,
		req.MetodePenerimaan,
		prov, kota, kec, pos, alm, ready,
	)
	return err
}

// HandlerSimpanDanValidasiWilayah menangani request POST /api/checkout/wilayah-pengiriman dari Frontend
func (h *HandlerValidasiWilayah) HandlerSimpanDanValidasiWilayah(w http.ResponseWriter, r *http.Request) {
	// Dukungan CORS preflight
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSONWilayah(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "error",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	var req PermintaanWilayahPengiriman
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimJSONWilayah(w, http.StatusBadRequest, ResponsAPI{
			Status: "error",
			Pesan:  "Format payload JSON tidak valid",
			Errors: map[string]string{"json": err.Error()},
		})
		return
	}

	valid, errors, paramKurir := ValidasiParameterKurir(&req)
	if !valid {
		kirimJSONWilayah(w, http.StatusBadRequest, ResponsAPI{
			Status: "error",
			Pesan:  "Data wilayah dan alamat tujuan belum lengkap atau tidak valid sesuai parameter API kurir",
			Errors: errors,
		})
		return
	}

	// Simpan jika terhubung ke MySQL
	_ = h.SimpanAlamatDatabase(&req, paramKurir)

	pesanSukses := "Data wilayah dan alamat pengiriman valid dan siap dikirim ke API ongkos kirim (PPL-16)."
	if req.MetodePenerimaan == "ambil_sendiri" {
		pesanSukses = "Metode ambil sendiri dipilih. Pesanan dapat diambil langsung di outlet Situbondo."
	}

	hasil := HasilValidasiWilayah{
		Status:           "sukses",
		Pesan:            pesanSukses,
		Valid:            true,
		MetodePenerimaan: req.MetodePenerimaan,
		Wilayah:          req.WilayahPengiriman,
		ParameterKurir:   paramKurir,
	}

	kirimJSONWilayah(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  pesanSukses,
		Data:   hasil,
	})
}

// HandlerValidasiAlamatKurir menangani request POST /api/checkout/validasi-alamat
func (h *HandlerValidasiWilayah) HandlerValidasiAlamatKurir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSONWilayah(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "error",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	var req PermintaanWilayahPengiriman
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimJSONWilayah(w, http.StatusBadRequest, ResponsAPI{
			Status: "error",
			Pesan:  "Format payload JSON tidak valid",
		})
		return
	}

	valid, errors, paramKurir := ValidasiParameterKurir(&req)
	if !valid {
		kirimJSONWilayah(w, http.StatusBadRequest, ResponsAPI{
			Status: "error",
			Pesan:  "Validasi data geografis kurir gagal",
			Errors: errors,
		})
		return
	}

	kirimJSONWilayah(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Validasi data geografis kurir berhasil",
		Data: map[string]interface{}{
			"valid":           true,
			"parameter_kurir": paramKurir,
		},
	})
}

// HandlerMasterWilayah menangani request GET /api/checkout/master-wilayah
func (h *HandlerValidasiWilayah) HandlerMasterWilayah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSONWilayah(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "error",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	kirimJSONWilayah(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Master data wilayah bertingkat Indonesia",
		Data:   MasterWilayahIndonesia,
	})
}

// HandlerParameterKurir menangani request GET /api/checkout/parameter-kurir
func (h *HandlerValidasiWilayah) HandlerParameterKurir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSONWilayah(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "error",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	spek := map[string]interface{}{
		"modul":       "WBS 1.2.6.3 - Validasi Data Geografis Pengiriman API Kurir",
		"asal_toko":   AsalTokoRengginang,
		"kurir_mitra": []string{"JNE", "SiCepat", "J&T", "POS Indonesia", "Biteship", "RajaOngkir"},
		"parameter_wajib": map[string]string{
			"provinsi":        "Nama provinsi tujuan (string, terdaftar)",
			"kota_kabupaten":  "Nama kota/kabupaten tujuan (string, sesuai provinsi)",
			"kecamatan":       "Nama kecamatan tujuan (string, sesuai kota/kabupaten)",
			"kode_pos":        "5 digit angka kode pos Indonesia",
			"alamat_lengkap":  "Detail jalan, nomor rumah, RT/RW, kelurahan (min 10 karakter)",
			"nama_penerima":   "Nama orang yang menerima kiriman (min 2 karakter)",
			"nomor_telepon":   "Nomor aktif WhatsApp/telepon format Indonesia",
			"berat_gram":      "Total berat paket pengiriman dalam satuan gram",
		},
	}

	kirimJSONWilayah(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Spesifikasi parameter kebutuhan API kurir",
		Data:   spek,
	})
}

// Helper untuk kirim JSON response
func kirimJSONWilayah(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
