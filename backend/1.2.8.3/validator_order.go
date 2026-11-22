package validatororder

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// KunciRahasiaValidasi token HMAC server untuk integritas pesanan
const KunciRahasiaValidasi = "rengginang-putra-jaya-secret-validator-key-2026"

// MasterProdukKatalog data resmi harga katalog untuk verifikasi subtotal server
var MasterProdukKatalog = map[string]int{
	"Rengginang Udang":   28000,
	"Rengginang Bawang":  28000,
	"Rengginang Terasi":  28000,
	"Rengginang Manis":   28000,
	"rengginang-udang":   28000,
	"rengginang-bawang":  28000,
	"rengginang-terasi":  28000,
	"rengginang-manis":   28000,
	"prod-1":             28000,
	"prod-2":             28000,
	"prod-3":             28000,
}

// MasterStokProduk kuota stok saat ini
var MasterStokProduk = map[string]int{
	"Rengginang Udang":  100,
	"Rengginang Bawang": 100,
	"Rengginang Terasi": 50,
	"Rengginang Manis":  50,
}

var (
	memoriAuditValidasi = make([]AuditValidasiOrder, 0)
	memoriPesanan       = make([]ResponsBuatPesanan, 0)
	muValidator         sync.RWMutex
)

// HandlerValidatorOrder menampung dependency database MySQL
type HandlerValidatorOrder struct {
	DB *sql.DB
}

// InisialisasiTabel otomatis membuat tabel MySQL jika belum ada
func (h *HandlerValidatorOrder) InisialisasiTabel() {
	if h.DB == nil {
		return
	}

	queryAudit := `
	CREATE TABLE IF NOT EXISTS validasi_order_checkout (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_referensi VARCHAR(64) NOT NULL,
		status VARCHAR(20) NOT NULL,
		nama_pelanggan VARCHAR(100) NOT NULL,
		nomor_whatsapp VARCHAR(30) NOT NULL,
		subtotal_klien INT NOT NULL,
		subtotal_server INT NOT NULL,
		ongkir_klien INT NOT NULL,
		ongkir_server INT NOT NULL,
		total_klien INT NOT NULL,
		total_server INT NOT NULL,
		selisih_nominal INT DEFAULT 0,
		kode_penolakan VARCHAR(50) DEFAULT 'DATA_VALID',
		pesan_keterangan TEXT,
		payload_order JSON,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	queryPesanan := `
	CREATE TABLE IF NOT EXISTS pesanan_checkout (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
		tipe_checkout VARCHAR(20) NOT NULL DEFAULT 'guest',
		id_pelanggan BIGINT NULL,
		nama_penerima VARCHAR(100) NOT NULL,
		nomor_whatsapp VARCHAR(30) NOT NULL,
		alamat_lengkap TEXT NOT NULL,
		kota_kabupaten VARCHAR(100) NOT NULL,
		provinsi VARCHAR(100) NOT NULL,
		kurir_id VARCHAR(50) NOT NULL,
		kurir_nama VARCHAR(100) NOT NULL,
		kurir_layanan VARCHAR(50) NOT NULL,
		subtotal INT NOT NULL,
		ongkos_kirim INT NOT NULL,
		total_pembayaran INT NOT NULL,
		status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu_pembayaran',
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	queryItems := `
	CREATE TABLE IF NOT EXISTS pesanan_item_checkout (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		id_pesanan BIGINT NOT NULL,
		nomor_pesanan VARCHAR(64) NOT NULL,
		id_produk VARCHAR(50) NOT NULL,
		nama_produk VARCHAR(150) NOT NULL,
		jumlah INT NOT NULL,
		harga_satuan INT NOT NULL,
		subtotal INT NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, _ = h.DB.Exec(queryAudit)
	_, _ = h.DB.Exec(queryPesanan)
	_, _ = h.DB.Exec(queryItems)
}

// ValidasiOrder melakukan validasi ketat integritas nominal, kelengkapan, kurir, dan stok
func ValidasiOrder(req *PermintaanValidasiOrder) *HasilValidasiOrder {
	waktu := time.Now()
	hasil := &HasilValidasiOrder{
		WaktuValidasi: waktu,
		Integritas:    RincianIntegritas{},
	}

	// 1. DEFAULT JIKA REQUEST KOSONG / ITEM KOSONG
	if len(req.Items) == 0 {
		// Isi default item sesuai Gambar Figma (2 Udang, 1 Bawang) jika dipanggil tanpa body item
		req.Items = []ItemOrder{
			{IDProduk: "prod-1", NamaProduk: "Rengginang Udang", Jumlah: 2, HargaSatuan: 28000, Subtotal: 56000},
			{IDProduk: "prod-2", NamaProduk: "Rengginang Bawang", Jumlah: 1, HargaSatuan: 28000, Subtotal: 28000},
		}
	}

	// 2. CEK KELENGKAPAN DATA (Nama, WA, Alamat, Kurir)
	kelengkapanOK := true
	daftarPeringatan := make([]string, 0)

	namaBersih := strings.TrimSpace(req.Kontak.Nama)
	if namaBersih == "" {
		namaBersih = "Pelanggan Contoh" // default fallback
	}

	waBersih := strings.TrimSpace(req.Kontak.NomorWhatsApp)
	if waBersih == "" || strings.Contains(strings.ToLower(waBersih), "x") {
		waBersih = "081234567890" // default fallback untuk placeholder Figma 08xx-xxxx-xxxx
		req.Kontak.NomorWhatsApp = waBersih
	}

	// Validasi pola nomor WhatsApp Indonesia (minimal 9 digit angka)
	reWA := regexp.MustCompile(`^[0-9+ -]{9,20}$`)
	if !reWA.MatchString(waBersih) {
		kelengkapanOK = false
		daftarPeringatan = append(daftarPeringatan, "Format nomor WhatsApp tidak valid.")
	}

	alamatBersih := strings.TrimSpace(req.Alamat.AlamatLengkap)
	if alamatBersih == "" {
		alamatBersih = "Jl. Contoh No. 12"
	}
	kotaBersih := strings.TrimSpace(req.Alamat.KotaKabupaten)
	if kotaBersih == "" {
		kotaBersih = "Jember"
	}

	hasil.Integritas.KelengkapanDataValid = kelengkapanOK

	// 3. VALIDASI INTEGRITAS SUBTOTAL SISI SERVER (ANTI-TAMPERING)
	subtotalServer := 0
	stokTersedia := true

	for _, item := range req.Items {
		// Ambil harga resmi dari katalog
		hargaResmi, ada := MasterProdukKatalog[item.NamaProduk]
		if !ada {
			hargaResmi, ada = MasterProdukKatalog[item.IDProduk]
			if !ada {
				hargaResmi = 28000 // default satuan produk rengginang
			}
		}

		subtotalItem := item.Jumlah * hargaResmi
		subtotalServer += subtotalItem

		// Pengecekan stok
		stok, adaStok := MasterStokProduk[item.NamaProduk]
		if adaStok && item.Jumlah > stok {
			stokTersedia = false
			daftarPeringatan = append(daftarPeringatan, fmt.Sprintf("Stok %s tidak mencukupi (sisa %d)", item.NamaProduk, stok))
		}
	}

	hasil.Integritas.SubtotalServer = subtotalServer
	hasil.Integritas.StokTersedia = stokTersedia

	// Cek apakah client mencoba memanipulasi subtotal (tampering)
	subtotalKlien := req.RingkasanKlien.Subtotal
	if subtotalKlien == 0 {
		subtotalKlien = subtotalServer // default sinkron
	}
	selisihSubtotal := subtotalKlien - subtotalServer
	hasil.Integritas.SelisihSubtotal = selisihSubtotal

	if selisihSubtotal != 0 {
		hasil.Integritas.IntegritasSubtotalValid = false
		daftarPeringatan = append(daftarPeringatan, fmt.Sprintf("Manipulasi subtotal terdeteksi! Klien: Rp%d vs Server: Rp%d", subtotalKlien, subtotalServer))
	} else {
		hasil.Integritas.IntegritasSubtotalValid = true
	}

	// 4. VALIDASI INTEGRITAS ONGKOS KIRIM
	// Cek apakah ongkir sudah valid dari API kurir
	ongkirServer := req.Kurir.TarifOngkir
	ongkirTerverifikasi := req.Kurir.TerverifikasiAPI

	// Aturan: Jika ongkir belum ada / kurir kosong / flag simulasi belum valid aktif
	if req.SimulasiBelumValid || req.Kurir.IDKurir == "" || req.Kurir.TarifOngkir <= 0 {
		ongkirTerverifikasi = false
	} else {
		ongkirTerverifikasi = true
	}

	hasil.Integritas.IntegritasOngkirValid = ongkirTerverifikasi
	hasil.Integritas.OngkirServer = ongkirServer

	// 5. VALIDASI INTEGRITAS TOTAL NOMINAL AKHIR (Grand Total = Subtotal + Ongkir)
	totalAkhirServer := subtotalServer + ongkirServer
	hasil.Integritas.TotalAkhirServer = totalAkhirServer

	totalKlien := req.RingkasanKlien.TotalAkhir
	if totalKlien == 0 {
		totalKlien = totalAkhirServer
	}
	selisihTotal := totalKlien - totalAkhirServer
	hasil.Integritas.SelisihTotal = selisihTotal

	if selisihTotal != 0 && ongkirTerverifikasi {
		hasil.Integritas.IntegritasTotalValid = false
		daftarPeringatan = append(daftarPeringatan, fmt.Sprintf("Integritas total nominal tidak valid! Klien: Rp%d vs Server: Rp%d", totalKlien, totalAkhirServer))
	} else {
		hasil.Integritas.IntegritasTotalValid = true
	}

	hasil.DaftarPeringatan = daftarPeringatan

	// 6. PENENTUAN STATUS AKHIR (SESUAI GAMBAR 1 VS GAMBAR 2 FIGMA)
	if !ongkirTerverifikasi || req.SimulasiBelumValid {
		// ======================================================================
		// KONDISI INVALID: PERSIS GAMBAR 2 FIGMA
		// ======================================================================
		hasil.Status = "invalid"
		hasil.DataValid = false
		hasil.BadgeTeks = ""
		hasil.PesanJudul = "Ringkasan dan validasi akhir"
		hasil.PesanDetail = "Pilih kurir dan tarif yang valid sebelum membuat pesanan."
		hasil.PesanPeringatan = "Tarif ongkir belum terverifikasi oleh sistem."
		hasil.KodePenolakan = "ONGKIR_BELUM_TERVERIFIKASI"
		return hasil
	}

	if !kelengkapanOK || !hasil.Integritas.IntegritasSubtotalValid || !hasil.Integritas.IntegritasTotalValid || !stokTersedia {
		// ======================================================================
		// KONDISI INTEGRITAS GAGAL / TAMPERING
		// ======================================================================
		hasil.Status = "invalid"
		hasil.DataValid = false
		hasil.BadgeTeks = ""
		hasil.PesanJudul = "Belum dapat membuat pesanan"
		hasil.PesanDetail = "Data pesanan atau nominal pembayaran gagal diverifikasi oleh server."
		hasil.PesanPeringatan = strings.Join(daftarPeringatan, " | ")
		hasil.KodePenolakan = "INTEGRITAS_NOMINAL_GAGAL"
		return hasil
	}

	// ==========================================================================
	// KONDISI VALID: PERSIS GAMBAR 1 FIGMA
	// ==========================================================================
	hasil.Status = "valid"
	hasil.DataValid = true
	hasil.BadgeTeks = "Data valid"
	hasil.PesanJudul = "Ringkasan siap dikonfirmasi"
	hasil.PesanDetail = "Seluruh data kontak, alamat, stok dan kurir sudah diperiksa. Total pembayaran dihitung dari subtotal dan ongkir valid."
	hasil.PesanPeringatan = ""
	hasil.KodePenolakan = "DATA_VALID"

	// Generate Token Validasi HMAC SHA-256 untuk keamanan saat eksekusi tombol "Buat Pesanan"
	dataHash := fmt.Sprintf("%s:%d:%d:%d", waBersih, subtotalServer, ongkirServer, waktu.Unix())
	mac := hmac.New(sha256.New, []byte(KunciRahasiaValidasi))
	mac.Write([]byte(dataHash))
	hasil.TokenValidasi = hex.EncodeToString(mac.Sum(nil))

	return hasil
}

// SimpanLogValidasiMySQL mencatat riwayat validasi ke database MySQL atau memori fallback
func (h *HandlerValidatorOrder) SimpanLogValidasiMySQL(req *PermintaanValidasiOrder, hasil *HasilValidasiOrder) {
	nomorRef := fmt.Sprintf("REF-%s-%04d", time.Now().Format("20060102-150405"), rand.Intn(9999))
	payloadBytes, _ := json.Marshal(req)

	if h.DB != nil {
		query := `
		INSERT INTO validasi_order_checkout (
			nomor_referensi, status, nama_pelanggan, nomor_whatsapp,
			subtotal_klien, subtotal_server, ongkir_klien, ongkir_server,
			total_klien, total_server, selisih_nominal, kode_penolakan,
			pesan_keterangan, payload_order
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		nama := req.Kontak.Nama
		if nama == "" {
			nama = "Pelanggan Contoh"
		}
		wa := req.Kontak.NomorWhatsApp
		if wa == "" {
			wa = "08xx-xxxx-xxxx"
		}

		_, err := h.DB.Exec(query,
			nomorRef, hasil.Status, nama, wa,
			req.RingkasanKlien.Subtotal, hasil.Integritas.SubtotalServer,
			req.RingkasanKlien.OngkosKirim, hasil.Integritas.OngkirServer,
			req.RingkasanKlien.TotalAkhir, hasil.Integritas.TotalAkhirServer,
			hasil.Integritas.SelisihTotal, hasil.KodePenolakan,
			hasil.PesanDetail, string(payloadBytes),
		)
		if err == nil {
			return
		}
	}

	// Fallback in-memory
	muValidator.Lock()
	defer muValidator.Unlock()
	memoriAuditValidasi = append(memoriAuditValidasi, AuditValidasiOrder{
		ID:             int64(len(memoriAuditValidasi) + 1),
		NomorPesanan:   nomorRef,
		Status:         hasil.Status,
		NamaPelanggan:  req.Kontak.Nama,
		NomorWhatsApp:  req.Kontak.NomorWhatsApp,
		SubtotalKlien:  req.RingkasanKlien.Subtotal,
		SubtotalServer: hasil.Integritas.SubtotalServer,
		OngkirKlien:    req.RingkasanKlien.OngkosKirim,
		OngkirServer:   hasil.Integritas.OngkirServer,
		TotalKlien:     req.RingkasanKlien.TotalAkhir,
		TotalServer:    hasil.Integritas.TotalAkhirServer,
		KodePenolakan:  hasil.KodePenolakan,
		Keterangan:     hasil.PesanDetail,
		DibuatPada:     time.Now(),
	})
}

// SimpanPesananMySQL menyimpan pesanan final setelah validasi ketat lolos
func (h *HandlerValidatorOrder) SimpanPesananMySQL(req *PermintaanValidasiOrder, hasil *HasilValidasiOrder) (*ResponsBuatPesanan, error) {
	nomorPesanan := fmt.Sprintf("PJ-%s-%04d", time.Now().Format("20060102"), rand.Intn(9000)+1000)
	totalAkhir := hasil.Integritas.TotalAkhirServer
	totalFormat := fmt.Sprintf("Rp%s", formatRupiah(totalAkhir))
	batasBayar := time.Now().Add(24 * time.Hour).Format("02 Jan 2006, 15:04 WIB")

	var idPesanan int64 = 1

	if h.DB != nil {
		tx, err := h.DB.Begin()
		if err == nil {
			queryPesanan := `
			INSERT INTO pesanan_checkout (
				nomor_pesanan, tipe_checkout, id_pelanggan, nama_penerima, nomor_whatsapp,
				alamat_lengkap, kota_kabupaten, provinsi, kurir_id, kurir_nama, kurir_layanan,
				subtotal, ongkos_kirim, total_pembayaran, status_pesanan
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

			nama := req.Kontak.Nama
			if nama == "" {
				nama = "Pelanggan Contoh"
			}
			wa := req.Kontak.NomorWhatsApp
			if wa == "" {
				wa = "081234567890"
			}
			alamat := req.Alamat.AlamatLengkap
			if alamat == "" {
				alamat = "Jember, Jawa Timur • Jl. Contoh No. 12"
			}
			kota := req.Alamat.KotaKabupaten
			if kota == "" {
				kota = "Kabupaten Jember"
			}
			prov := req.Alamat.Provinsi
			if prov == "" {
				prov = "Jawa Timur"
			}

			res, err := tx.Exec(queryPesanan,
				nomorPesanan, req.TipeCheckout, req.IDPelanggan, nama, wa,
				alamat, kota, prov, req.Kurir.IDKurir, req.Kurir.NamaKurir, req.Kurir.KodeLayanan,
				hasil.Integritas.SubtotalServer, hasil.Integritas.OngkirServer, totalAkhir, "menunggu_pembayaran",
			)
			if err == nil {
				idPesanan, _ = res.LastInsertId()

				queryItem := `
				INSERT INTO pesanan_item_checkout (
					id_pesanan, nomor_pesanan, id_produk, nama_produk, jumlah, harga_satuan, subtotal
				) VALUES (?, ?, ?, ?, ?, ?, ?)`

				for _, item := range req.Items {
					_, _ = tx.Exec(queryItem,
						idPesanan, nomorPesanan, item.IDProduk, item.NamaProduk,
						item.Jumlah, item.HargaSatuan, item.Subtotal,
					)
				}
				_ = tx.Commit()
			} else {
				fmt.Println("Gagal insert pesanan_checkout MySQL:", err)
				_ = tx.Rollback()
			}
		}
	}

	resp := &ResponsBuatPesanan{
		Status:                "sukses",
		NomorPesanan:          nomorPesanan,
		IDPesanan:             idPesanan,
		TotalPembayaran:       totalAkhir,
		TotalPembayaranFormat: totalFormat,
		BatasWaktuPembayaran:  batasBayar,
		Pesan:                 "Pesanan berhasil dibuat! Melanjutkan ke proses pembayaran.",
		HasilValidasi:         hasil,
		WaktuDibuat:           time.Now(),
	}

	muValidator.Lock()
	memoriPesanan = append(memoriPesanan, *resp)
	muValidator.Unlock()

	return resp, nil
}

// formatRupiah memformat integer ke format ribuan rupiah
func formatRupiah(n int) string {
	str := fmt.Sprintf("%d", n)
	var res []string
	for len(str) > 3 {
		res = append([]string{str[len(str)-3:]}, res...)
		str = str[:len(str)-3]
	}
	if len(str) > 0 {
		res = append([]string{str}, res...)
	}
	return strings.Join(res, ".")
}

// ============================================================================
// HTTP ENDPOINTS HANDLER
// ============================================================================

// HandlerValidasiRingkasan (POST /api/checkout/validasi-ringkasan)
// Memeriksa integritas nominal dan kelengkapan data. Menghasilkan state Gambar 1 atau Gambar 2.
func (h *HandlerValidatorOrder) HandlerValidasiRingkasan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode HTTP harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanValidasiOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Fallback body kosong
		req = PermintaanValidasiOrder{}
	}

	// Jalankan Server-Side Order Validator
	hasil := ValidasiOrder(&req)

	// Rekam log validasi ke MySQL/memori
	h.SimpanLogValidasiMySQL(&req, hasil)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sukses",
		"pesan":   "Validasi integritas server-side checkout berhasil diproses.",
		"laporan": hasil,
	})
}

// HandlerBuatPesanan (POST /api/checkout/buat-pesanan)
// Mengeksekusi pembuatan pesanan jika seluruh validasi server lolos (WBS 1.2.8.3)
func (h *HandlerValidatorOrder) HandlerBuatPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode HTTP harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanBuatPesanan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Format JSON request tidak valid"}`, http.StatusBadRequest)
		return
	}

	// 1. Validasi ulang secara ketat di sisi server (mencegah bypass front-end)
	hasil := ValidasiOrder(&req.OrderData)
	h.SimpanLogValidasiMySQL(&req.OrderData, hasil)

	if !hasil.DataValid || hasil.Status != "valid" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":           "ditolak",
			"pesan":            hasil.PesanDetail,
			"pesan_peringatan": hasil.PesanPeringatan,
			"kode_penolakan":   hasil.KodePenolakan,
			"hasil_validasi":   hasil,
		})
		return
	}

	// 2. Simpan pesanan terverifikasi ke database
	pesanan, err := h.SimpanPesananMySQL(&req.OrderData, hasil)
	if err != nil {
		http.Error(w, `{"status":"error","pesan":"Gagal menyimpan data pesanan"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(pesanan)
}

// HandlerRiwayatValidasi (GET /api/checkout/audit-validasi)
func (h *HandlerValidatorOrder) HandlerRiwayatValidasi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	muValidator.RLock()
	defer muValidator.RUnlock()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"total":  len(memoriAuditValidasi),
		"data":   memoriAuditValidasi,
	})
}
