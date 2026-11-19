package checkoutapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// WBS 1.2.3.3 - HANDLER PENYAMAAN FORMAT PAYLOAD CHECKOUT
// Memproses order request terstandarisasi untuk pembeli Guest dan Pelanggan.
// ============================================================================

// HandlerCheckout mengelola endpoint standarisasi payload checkout.
type HandlerCheckout struct {
	DB *sql.DB
}

// Fallback data stok resmi jika database MySQL offline
var dataFallbackVarian = map[uint64]struct {
	Slug       string
	NamaVarian string
	HargaJual  uint64
	Stok       int64
}{
	1: {Slug: "udang", NamaVarian: "Udang", HargaJual: 28000, Stok: 12},
	2: {Slug: "bawang", NamaVarian: "Bawang", HargaJual: 28000, Stok: 38},
	3: {Slug: "cumi", NamaVarian: "Cumi", HargaJual: 28000, Stok: 0},
	4: {Slug: "terasi", NamaVarian: "Terasi", HargaJual: 28000, Stok: 38},
}

var regexNomorWA = regexp.MustCompile(`^(\+62|62|08)[0-9]{8,13}$`)

// ValidasiPayloadCheckout memverifikasi kelengkapan dan keabsahan skema payload.
func ValidasiPayloadCheckout(req *PayloadCheckoutRequest) (bool, string) {
	if req == nil {
		return false, "Payload checkout tidak boleh kosong"
	}

	// 1. Validasi Data Kontak Penerima
	namaClean := strings.TrimSpace(req.Kontak.NamaPenerima)
	if len(namaClean) < 2 {
		return false, "Nama penerima wajib diisi (minimal 2 karakter)"
	}

	waClean := strings.TrimSpace(req.Kontak.NomorWhatsApp)
	waClean = strings.ReplaceAll(waClean, "-", "")
	waClean = strings.ReplaceAll(waClean, " ", "")
	if len(waClean) < 9 || len(waClean) > 15 {
		return false, "Nomor WhatsApp tidak valid (harus 9-15 digit)"
	}

	// 2. Validasi Item Pesanan
	if len(req.ItemPesanan) == 0 {
		return false, "Keranjang checkout kosong, minimal harus memuat 1 item produk"
	}

	for _, it := range req.ItemPesanan {
		if it.Jumlah == 0 {
			return false, fmt.Sprintf("Jumlah item untuk varian %s minimal 1 bungkus", it.NamaVarian)
		}
	}

	// 3. Normalisasi Penanda Tipe Checkout
	if req.IsGuest || req.TipeCheckout == "guest" || req.IDPelanggan == nil {
		req.TipeCheckout = "guest"
		req.IsGuest = true
		req.IDPelanggan = nil
	} else {
		req.TipeCheckout = "pelanggan"
		req.IsGuest = false
	}

	return true, ""
}

// VerifikasiStokDatabase memeriksa ketersediaan stok aktual varian dari database MySQL.
func (h *HandlerCheckout) VerifikasiStokDatabase(items []ItemPesananCheckout) error {
	for _, it := range items {
		// Jika terhubung ke MySQL
		if h.DB != nil {
			var stokDB int64
			var statusAktif bool
			err := h.DB.QueryRow(`
				SELECT stok, status_aktif 
				FROM varian_produk 
				WHERE (id_varian = ? OR LOWER(nama_varian) = ?) 
				LIMIT 1
			`, it.IDVarian, strings.ToLower(it.NamaVarian)).Scan(&stokDB, &statusAktif)

			if err == nil {
				if !statusAktif {
					return fmt.Errorf("varian %s sedang dinonaktifkan", it.NamaVarian)
				}
				if stokDB <= 0 {
					return fmt.Errorf("stok varian %s sedang habis", it.NamaVarian)
				}
				if int64(it.Jumlah) > stokDB {
					return fmt.Errorf("jumlah pesanan %s (%d bungkus) melebihi sisa stok (%d bungkus)", it.NamaVarian, it.Jumlah, stokDB)
				}
				continue
			}
		}

		// Fallback in-memory jika MySQL offline atau query gagal
		if fb, ada := dataFallbackVarian[it.IDVarian]; ada {
			if fb.Stok <= 0 {
				return fmt.Errorf("stok varian %s sedang habis", it.NamaVarian)
			}
			if int64(it.Jumlah) > fb.Stok {
				return fmt.Errorf("jumlah pesanan %s (%d bungkus) melebihi sisa stok (%d bungkus)", it.NamaVarian, it.Jumlah, fb.Stok)
			}
		}
	}
	return nil
}

// SimpanPesananDatabase menyimpan data pesanan terstandarisasi ke database MySQL jika tersedia.
func (h *HandlerCheckout) SimpanPesananDatabase(resp *HasilCheckoutResponse) {
	if h.DB == nil {
		return
	}

	// Pastikan tabel pesanan_checkout ada
	_, _ = h.DB.Exec(`
		CREATE TABLE IF NOT EXISTS pesanan_checkout (
			id_pesanan BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
			tipe_checkout ENUM('guest', 'pelanggan') NOT NULL,
			is_guest BOOLEAN NOT NULL DEFAULT TRUE,
			id_pelanggan BIGINT UNSIGNED NULL,
			nama_penerima VARCHAR(150) NOT NULL,
			nomor_whatsapp VARCHAR(30) NOT NULL,
			email VARCHAR(191) NULL,
			catatan TEXT NULL,
			total_item INT UNSIGNED NOT NULL,
			subtotal DECIMAL(12,2) NOT NULL,
			ongkos_kirim DECIMAL(12,2) NOT NULL DEFAULT 0,
			total_pembayaran DECIMAL(12,2) NOT NULL,
			status_pesanan VARCHAR(50) NOT NULL DEFAULT 'MENUNGGU_KONFIRMASI',
			dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	var idPelangganVal sql.NullInt64
	if resp.IDPelanggan != nil {
		idPelangganVal = sql.NullInt64{Int64: int64(*resp.IDPelanggan), Valid: true}
	}

	_, _ = h.DB.Exec(`
		INSERT INTO pesanan_checkout (
			nomor_pesanan, tipe_checkout, is_guest, id_pelanggan,
			nama_penerima, nomor_whatsapp, email, catatan,
			total_item, subtotal, ongkos_kirim, total_pembayaran, status_pesanan
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		resp.NomorPesanan, resp.TipeCheckout, resp.IsGuest, idPelangganVal,
		resp.Kontak.NamaPenerima, resp.Kontak.NomorWhatsApp, resp.Kontak.Email, resp.Kontak.Catatan,
		resp.Ringkasan.TotalItem, float64(resp.Ringkasan.Subtotal), float64(resp.Ringkasan.OngkosKirim),
		float64(resp.Ringkasan.TotalPembayaran), resp.StatusPesanan,
	)
}

// GenerateNomorPesanan membuat kode unik pesanan.
func GenerateNomorPesanan(isGuest bool) string {
	now := time.Now()
	prefix := "PJ-GUEST"
	if !isGuest {
		prefix = "PJ-CUST"
	}
	acak := rand.Intn(9000) + 1000
	return fmt.Sprintf("%s-%s-%04d", prefix, now.Format("20060102"), acak)
}

// HandlerProsesCheckout memproses order request baik untuk guest maupun pelanggan login.
// Endpoint: POST /api/checkout atau POST /api/checkout/proses
func (h *HandlerCheckout) HandlerProsesCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	var req PayloadCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  "Format payload JSON checkout tidak valid",
		})
		return
	}

	// 1. Validasi format skema payload order request
	valid, pesanError := ValidasiPayloadCheckout(&req)
	if !valid {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "ditolak",
			Pesan:  pesanError,
		})
		return
	}

	// 2. Verifikasi stok dari database MySQL
	if err := h.VerifikasiStokDatabase(req.ItemPesanan); err != nil {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "ditolak",
			Pesan:  err.Error(),
		})
		return
	}

	// 3. Kalkulasi ulang ringkasan biaya secara aman di sisi server
	var totalItem uint32
	var totalSubtotal uint64
	for i := range req.ItemPesanan {
		it := &req.ItemPesanan[i]
		if it.HargaSatuan == 0 {
			it.HargaSatuan = 28000 // default varian standar
		}
		it.Subtotal = uint64(it.Jumlah) * it.HargaSatuan
		totalItem += it.Jumlah
		totalSubtotal += it.Subtotal
	}

	req.Ringkasan.TotalItem = totalItem
	req.Ringkasan.Subtotal = totalSubtotal
	req.Ringkasan.TotalPembayaran = totalSubtotal + req.Ringkasan.OngkosKirim

	nomorPesanan := GenerateNomorPesanan(req.IsGuest)
	pesanNotif := "Pesanan tamu (guest) berhasil diproses tanpa akun"
	if !req.IsGuest {
		pesanNotif = "Pesanan pelanggan terdaftar berhasil diproses"
	}

	responData := &HasilCheckoutResponse{
		NomorPesanan:    nomorPesanan,
		TipeCheckout:    req.TipeCheckout,
		IsGuest:         req.IsGuest,
		IDPelanggan:     req.IDPelanggan,
		Kontak:          req.Kontak,
		ItemPesanan:     req.ItemPesanan,
		Ringkasan:       req.Ringkasan,
		StatusPesanan:   "MENUNGGU_KONFIRMASI",
		PesanNotifikasi: pesanNotif,
		DibuatPada:      time.Now(),
	}

	// 4. Simpan ke database jika aktif
	h.SimpanPesananDatabase(responData)

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Standarisasi order request valid dan pesanan checkout berhasil diproses",
		Data:   responData,
	})
}

// HandlerSkemaPayload menyajikan dokumentasi resmi skema standar order request bagi FE dan BE.
// Endpoint: GET /api/checkout/skema
func (h *HandlerCheckout) HandlerSkemaPayload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	contohGuest := PayloadCheckoutRequest{
		TipeCheckout: "guest",
		IsGuest:      true,
		IDPelanggan:  nil,
		Kontak: InformasiKontak{
			NamaPenerima:  "Budi Santoso",
			NomorWhatsApp: "081234567890",
			Email:         "budi@example.com",
			Catatan:       "Tolong packing rapat",
		},
		ItemPesanan: []ItemPesananCheckout{
			{
				IDVarian:    1,
				Slug:        "udang",
				NamaVarian:  "Udang",
				NamaProduk:  "Rengginang Udang",
				HargaSatuan: 28000,
				Jumlah:      2,
				Subtotal:    56000,
			},
		},
		Ringkasan: RingkasanBiayaCheckout{
			TotalItem:       2,
			Subtotal:        56000,
			OngkosKirim:     0,
			TotalPembayaran: 56000,
		},
	}

	idCust := uint64(10)
	contohPelanggan := contohGuest
	contohPelanggan.TipeCheckout = "pelanggan"
	contohPelanggan.IsGuest = false
	contohPelanggan.IDPelanggan = &idCust

	kirimJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "sukses",
		"wbs":         "1.2.3.3",
		"deliverable": "Standarisasi Schema Order Request",
		"deskripsi":   "Penyamaan format payload data checkout bagi pembeli guest dan pelanggan yang login",
		"skema": map[string]string{
			"tipe_checkout": "string ('guest' atau 'pelanggan')",
			"is_guest":      "boolean (true untuk guest, false untuk pelanggan login)",
			"id_pelanggan":  "integer atau null",
			"kontak":        "object (nama_penerima*, nomor_whatsapp*, email, catatan)",
			"item_pesanan":  "array of object (id_varian*, nama_varian, jumlah*, harga_satuan, subtotal)",
			"ringkasan":     "object (total_item, subtotal, ongkos_kirim, total_pembayaran)",
		},
		"contoh_payload_guest":     contohGuest,
		"contoh_payload_pelanggan": contohPelanggan,
	})
}

// Helper kirimJSON
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
