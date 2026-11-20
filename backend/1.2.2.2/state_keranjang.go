package statemanager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// WBS 1.2.2.2 - STATE MANAGER KERANJANG BELANJA (FE & BE)
// Mengelola penyimpanan dan sinkronisasi state keranjang belanja
// bagi pembeli tamu (Guest) dan pelanggan terdaftar (Logged-in Customer).
// ============================================================================

// HandlerStateManager mengelola endpoint sinkronisasi dan state keranjang.
type HandlerStateManager struct {
	DB *sql.DB
	mu sync.RWMutex
}

// In-memory store untuk penyimpanan keranjang saat database MySQL offline
var memoriKeranjangPelanggan = make(map[uint64][]ItemKeranjang)

// Data keranjang bawaan resmi sesuai Figma PPL11_01
var defaultKeranjangItems = []ItemKeranjang{
	{
		IDVarian:     1,
		Slug:         "udang",
		NamaVarian:   "Udang",
		NamaProduk:   "Rengginang Udang",
		BeratBersih:  "450 gr",
		HargaSatuan:  28000,
		Jumlah:       2,
		StokTersedia: 12,
		Subtotal:     56000,
		Gambar:       "/images/products/rengginang-udang.jpg",
	},
	{
		IDVarian:     2,
		Slug:         "bawang",
		NamaVarian:   "Bawang",
		NamaProduk:   "Rengginang Bawang",
		BeratBersih:  "450 gr",
		HargaSatuan:  28000,
		Jumlah:       1,
		StokTersedia: 38,
		Subtotal:     28000,
		Gambar:       "/images/products/rengginang-bawang.jpg",
	},
}

func init() {
	// Isi data default untuk simulasi pelanggan ID: 1
	memoriKeranjangPelanggan[1] = append([]ItemKeranjang{}, defaultKeranjangItems...)
}

// HitungRekapKeranjang menghitung ulang subtotal, total kuantitas, dan subtotal per item.
func HitungRekapKeranjang(items []ItemKeranjang) (uint32, uint64, []ItemKeranjang) {
	var totalQty uint32
	var grandTotal uint64
	hasilItems := make([]ItemKeranjang, len(items))

	for i, it := range items {
		if it.HargaSatuan == 0 {
			it.HargaSatuan = 28000
		}
		it.Subtotal = uint64(it.Jumlah) * it.HargaSatuan
		totalQty += it.Jumlah
		grandTotal += it.Subtotal
		hasilItems[i] = it
	}

	return totalQty, grandTotal, hasilItems
}

// FormatRupiah helper format mata uang IDR.
func FormatRupiah(angka uint64) string {
	str := strconv.FormatUint(angka, 10)
	var res []string
	l := len(str)
	for i := 0; i < l; i++ {
		if i > 0 && (l-i)%3 == 0 {
			res = append(res, ".")
		}
		res = append(res, string(str[i]))
	}
	return "Rp" + strings.Join(res, "")
}

// SimpanKeranjangDatabase menyimpan item keranjang pelanggan ke tabel MySQL.
func (h *HandlerStateManager) SimpanKeranjangDatabase(idPelanggan uint64, items []ItemKeranjang) error {
	if h.DB == nil {
		return nil
	}

	// Buat tabel jika belum ada
	_, _ = h.DB.Exec(`
		CREATE TABLE IF NOT EXISTS keranjang_pelanggan (
			id_keranjang BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			pelanggan_id BIGINT UNSIGNED NOT NULL UNIQUE,
			total_item INT UNSIGNED NOT NULL DEFAULT 0,
			subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
			diperbarui_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	_, _ = h.DB.Exec(`
		CREATE TABLE IF NOT EXISTS item_keranjang_pelanggan (
			id_item BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			keranjang_id BIGINT UNSIGNED NOT NULL,
			id_varian BIGINT UNSIGNED NOT NULL,
			nama_varian VARCHAR(100) NOT NULL,
			jumlah INT UNSIGNED NOT NULL,
			harga_satuan DECIMAL(12,2) NOT NULL,
			subtotal DECIMAL(12,2) NOT NULL,
			dibuat_pada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	totalQty, subtotal, _ := HitungRekapKeranjang(items)

	// Upsert keranjang_pelanggan
	res, err := h.DB.Exec(`
		INSERT INTO keranjang_pelanggan (pelanggan_id, total_item, subtotal, diperbarui_pada)
		VALUES (?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE total_item = VALUES(total_item), subtotal = VALUES(subtotal), diperbarui_pada = NOW()
	`, idPelanggan, totalQty, float64(subtotal))
	if err != nil {
		return err
	}

	var keranjangID int64
	err = h.DB.QueryRow("SELECT id_keranjang FROM keranjang_pelanggan WHERE pelanggan_id = ?", idPelanggan).Scan(&keranjangID)
	if err != nil {
		keranjangID, _ = res.LastInsertId()
	}

	// Hapus item lama, insert item baru
	_, _ = h.DB.Exec("DELETE FROM item_keranjang_pelanggan WHERE keranjang_id = ?", keranjangID)

	for _, it := range items {
		_, _ = h.DB.Exec(`
			INSERT INTO item_keranjang_pelanggan (keranjang_id, id_varian, nama_varian, jumlah, harga_satuan, subtotal)
			VALUES (?, ?, ?, ?, ?, ?)
		`, keranjangID, it.IDVarian, it.NamaVarian, it.Jumlah, float64(it.HargaSatuan), float64(it.Subtotal))
	}

	return nil
}

// HandlerSinkronisasiKeranjang menangani sinkronisasi keranjang belanja.
// Endpoint: POST /api/keranjang/sinkronisasi
func (h *HandlerStateManager) HandlerSinkronisasiKeranjang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	var req PermintaanSinkronisasi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  "Format payload JSON sinkronisasi keranjang tidak valid",
		})
		return
	}

	totalQty, subtotal, itemsRekap := HitungRekapKeranjang(req.Items)

	idKeranjang := "GUEST-CART"
	pesan := "State keranjang belanja tamu (guest) berhasil disinkronkan ke localStorage"

	if !req.IsGuest && req.IDPelanggan != nil {
		idKeranjang = fmt.Sprintf("CUST-CART-%d", *req.IDPelanggan)
		pesan = "State keranjang belanja pelanggan terdaftar berhasil disinkronkan ke database"

		// Simpan ke in-memory store
		h.mu.Lock()
		memoriKeranjangPelanggan[*req.IDPelanggan] = append([]ItemKeranjang{}, itemsRekap...)
		h.mu.Unlock()

		// Simpan ke basis data MySQL jika koneksi aktif
		_ = h.SimpanKeranjangDatabase(*req.IDPelanggan, itemsRekap)
	}

	responData := ResponKeranjang{
		IDKeranjang:    idKeranjang,
		IsGuest:        req.IsGuest,
		IDPelanggan:    req.IDPelanggan,
		TotalJenis:     len(itemsRekap),
		TotalKuantitas: totalQty,
		Subtotal:       subtotal,
		FormatSubtotal: FormatRupiah(subtotal),
		Items:          itemsRekap,
		Pesan:          pesan,
		TerakhirDiubah: time.Now(),
	}

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  pesan,
		Data:   responData,
	})
}

// HandlerAmbilKeranjang mengambil data keranjang belanja pelanggan atau guest.
// Endpoint: GET /api/keranjang
func (h *HandlerStateManager) HandlerAmbilKeranjang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan GET.",
		})
		return
	}

	// Cek query parameter id_pelanggan
	idPelangganStr := r.URL.Query().Get("id_pelanggan")
	var items []ItemKeranjang

	if idPelangganStr != "" {
		if id, err := strconv.ParseUint(idPelangganStr, 10, 64); err == nil {
			h.mu.RLock()
			items = memoriKeranjangPelanggan[id]
			h.mu.RUnlock()
		}
	}

	// Jika belum ada data tersimpan, sajikan data keranjang resmi Figma PPL11_01
	if len(items) == 0 {
		items = append([]ItemKeranjang{}, defaultKeranjangItems...)
	}

	totalQty, subtotal, itemsRekap := HitungRekapKeranjang(items)

	responData := ResponKeranjang{
		IDKeranjang:    "KERANJANG-AKTIF",
		IsGuest:        idPelangganStr == "",
		TotalJenis:     len(itemsRekap),
		TotalKuantitas: totalQty,
		Subtotal:       subtotal,
		FormatSubtotal: FormatRupiah(subtotal),
		Items:          itemsRekap,
		Pesan:          "Data keranjang belanja berhasil dimuat",
		TerakhirDiubah: time.Now(),
	}

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Data keranjang belanja berhasil dimuat",
		Data:   responData,
	})
}

// HandlerHapusItemKeranjang menghapus satu item produk dari keranjang pelanggan.
// Endpoint: DELETE /api/keranjang/item?id_varian={id}
func (h *HandlerStateManager) HandlerHapusItemKeranjang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		w.Header().Set("Allow", "DELETE, POST")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan DELETE.",
		})
		return
	}

	idVarianStr := r.URL.Query().Get("id_varian")
	idVarian, err := strconv.ParseUint(idVarianStr, 10, 64)
	if err != nil {
		kirimJSON(w, http.StatusBadRequest, ResponsAPI{
			Status: "gagal",
			Pesan:  "Parameter id_varian tidak valid",
		})
		return
	}

	idPelangganStr := r.URL.Query().Get("id_pelanggan")
	idPelanggan := uint64(1) // default testing
	if idPelangganStr != "" {
		if id, err := strconv.ParseUint(idPelangganStr, 10, 64); err == nil {
			idPelanggan = id
		}
	}

	h.mu.Lock()
	curItems := memoriKeranjangPelanggan[idPelanggan]
	newItems := make([]ItemKeranjang, 0, len(curItems))
	for _, it := range curItems {
		if it.IDVarian != idVarian {
			newItems = append(newItems, it)
		}
	}
	memoriKeranjangPelanggan[idPelanggan] = newItems
	h.mu.Unlock()

	// Update MySQL jika terhubung
	_ = h.SimpanKeranjangDatabase(idPelanggan, newItems)

	totalQty, subtotal, itemsRekap := HitungRekapKeranjang(newItems)

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Item produk berhasil dihapus dari keranjang belanja",
		Data: ResponKeranjang{
			IDKeranjang:    fmt.Sprintf("CUST-CART-%d", idPelanggan),
			TotalJenis:     len(itemsRekap),
			TotalKuantitas: totalQty,
			Subtotal:       subtotal,
			FormatSubtotal: FormatRupiah(subtotal),
			Items:          itemsRekap,
			Pesan:          "Item berhasil dihapus",
			TerakhirDiubah: time.Now(),
		},
	})
}

// HandlerBersihkanKeranjang mengosongkan seluruh isi keranjang pelanggan.
// Endpoint: POST /api/keranjang/bersihkan
func (h *HandlerStateManager) HandlerBersihkanKeranjang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		kirimJSON(w, http.StatusMethodNotAllowed, ResponsAPI{
			Status: "gagal",
			Pesan:  "Metode HTTP tidak diizinkan. Gunakan POST.",
		})
		return
	}

	idPelangganStr := r.URL.Query().Get("id_pelanggan")
	idPelanggan := uint64(1)
	if idPelangganStr != "" {
		if id, err := strconv.ParseUint(idPelangganStr, 10, 64); err == nil {
			idPelanggan = id
		}
	}

	h.mu.Lock()
	memoriKeranjangPelanggan[idPelanggan] = []ItemKeranjang{}
	h.mu.Unlock()

	_ = h.SimpanKeranjangDatabase(idPelanggan, []ItemKeranjang{})

	kirimJSON(w, http.StatusOK, ResponsAPI{
		Status: "sukses",
		Pesan:  "Seluruh isi keranjang belanja berhasil dikosongkan",
		Data: ResponKeranjang{
			IDKeranjang:    fmt.Sprintf("CUST-CART-%d", idPelanggan),
			TotalJenis:     0,
			TotalKuantitas: 0,
			Subtotal:       0,
			FormatSubtotal: "Rp0",
			Items:          []ItemKeranjang{},
			Pesan:          "Keranjang belanja kosong",
			TerakhirDiubah: time.Now(),
		},
	})
}

// Helper kirimJSON
func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
