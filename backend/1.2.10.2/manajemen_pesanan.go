package pesananowner

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LayananPesananOwner mengelola antrean pesanan masuk dan detailnya untuk dashboard Owner.
type LayananPesananOwner struct {
	DB            *sql.DB
	mu            sync.RWMutex
	dataSimulasi  []PesananRingkasOwner
	detailSimulasi map[string]*DetailPesananOwner
}

// InisialisasiLayanan membuat instance baru LayananPesananOwner dengan data default Figma.
func InisialisasiLayanan(db *sql.DB) *LayananPesananOwner {
	layanan := &LayananPesananOwner{
		DB:             db,
		detailSimulasi: make(map[string]*DetailPesananOwner),
	}

	layanan.inisialisasiDataSimulasi()
	layanan.InisialisasiTabel()
	return layanan
}

// InisialisasiTabel memastikan kolom status_pemrosesan tersedia pada tabel pesanan.
func (l *LayananPesananOwner) InisialisasiTabel() {
	if l.DB == nil {
		return
	}

	// Tambahkan kolom status_pemrosesan jika belum ada
	query := `
		SELECT COUNT(*) FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'pesanan' 
		AND COLUMN_NAME = 'status_pemrosesan';
	`
	var count int
	err := l.DB.QueryRow(query).Scan(&count)
	if err == nil && count == 0 {
		_, _ = l.DB.Exec("ALTER TABLE pesanan ADD COLUMN status_pemrosesan VARCHAR(30) NOT NULL DEFAULT 'menunggu_diproses'")
	}
}

// inisialisasiDataSimulasi memuat data pesanan acuan Figma PPL21_01 & PPL21_02.
func (l *LayananPesananOwner) inisialisasiDataSimulasi() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.dataSimulasi = []PesananRingkasOwner{
		{
			ID:           101,
			NomorPesanan: "PJ-DEMO-0001",
			Kontak:       "Pelanggan A",
			WhatsApp:     "0812-3456-7801",
			Metode:       "Ambil Sendiri",
			Pembayaran:   "Menunggu",
			Status:       "Baru",
			Total:        56000,
			TotalFormat:  "Rp56.000",
			ItemsRingkas: "Udang x 2",
			Tanggal:      time.Now().Format("2006-01-02"),
			WaktuDibuat:  time.Now(),
		},
		{
			ID:           102,
			NomorPesanan: "PJ-DEMO-0002",
			Kontak:       "Pelanggan B",
			WhatsApp:     "08xx-xxxx-xxxx",
			Metode:       "Kirim Kurir",
			Pembayaran:   "Dibayar",
			Status:       "Diproses",
			Total:        99000,
			TotalFormat:  "Rp99.000",
			ItemsRingkas: "Udang x 2 • Bawang x 1",
			Tanggal:      time.Now().Format("2006-01-02"),
			WaktuDibuat:  time.Now().Add(-1 * time.Hour),
		},
		{
			ID:           103,
			NomorPesanan: "PJ-DEMO-0003",
			Kontak:       "Pelanggan C",
			WhatsApp:     "0812-3456-7803",
			Metode:       "Kirim Kurir",
			Pembayaran:   "Dibayar",
			Status:       "Selesai",
			Total:        132000,
			TotalFormat:  "Rp132.000",
			ItemsRingkas: "Terasi x 2 • Cumi x 2",
			Tanggal:      time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			WaktuDibuat:  time.Now().AddDate(0, 0, -1),
		},
		{
			ID:           104,
			NomorPesanan: "PJ-DEMO-0004",
			Kontak:       "Pelanggan D",
			WhatsApp:     "0812-3456-7804",
			Metode:       "Ambil Sendiri",
			Pembayaran:   "Gagal",
			Status:       "Batal",
			Total:        28000,
			TotalFormat:  "Rp28.000",
			ItemsRingkas: "Bawang x 1",
			Tanggal:      time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			WaktuDibuat:  time.Now().AddDate(0, 0, -2),
		},
	}

	for _, item := range l.dataSimulasi {
		l.detailSimulasi[item.NomorPesanan] = &DetailPesananOwner{
			NomorPesanan:   item.NomorPesanan,
			NomorInvoice:   "INV-" + strings.TrimPrefix(item.NomorPesanan, "PJ-"),
			Nama:           item.Kontak,
			WhatsApp:       item.WhatsApp,
			Alamat:         item.Metode,
			Pengiriman:     item.Metode,
			Metode:         item.Metode,
			Subtotal:       item.Total,
			SubtotalFormat: item.TotalFormat,
			Ongkir:         0,
			OngkirFormat:   "Rp0",
			Total:          item.Total,
			TotalFormat:    item.TotalFormat,
			Pembayaran:     item.Pembayaran,
			Pemrosesan:     item.Status,
			ItemsRingkasan: item.ItemsRingkas,
			Tanggal:        item.Tanggal,
		}
	}
}

// AmbilDaftarPesanan mengambil daftar antrean pesanan dengan filter status, tanggal, dan kata kunci.
func (l *LayananPesananOwner) AmbilDaftarPesanan(filter FilterPesananOwner) ([]PesananRingkasOwner, int, error) {
	var hasil []PesananRingkasOwner

	if l.DB != nil {
		query := `
			SELECT 
				p.id, p.nomor_pesanan, p.nama_penerima, p.nomor_whatsapp, 
				p.tipe_checkout, p.status_pesanan,
				COALESCE(p.status_pemrosesan, 'menunggu_diproses') AS status_pemrosesan,
				p.snapshot_total, p.dibuat_pada
			FROM pesanan p
			WHERE 1=1
		`
		var args []interface{}

		// 1. Filter Pencarian Nama / Nomor Pesanan / WhatsApp
		if strings.TrimSpace(filter.KataKunci) != "" {
			query += " AND (p.nomor_pesanan LIKE ? OR p.nama_penerima LIKE ? OR p.nomor_whatsapp LIKE ?)"
			likePattern := "%" + strings.TrimSpace(filter.KataKunci) + "%"
			args = append(args, likePattern, likePattern, likePattern)
		}

		// 2. Filter Status Pesanan / Pembayaran / Pemrosesan
		if filter.Status != "" && filter.Status != "Semua status" {
			st := strings.ToLower(filter.Status)
			switch st {
			case "baru":
				query += " AND (p.status_pemrosesan = 'menunggu_diproses' OR p.status_pemrosesan = 'baru' OR p.status_pesanan = 'menunggu_pembayaran')"
			case "diproses":
				query += " AND (p.status_pemrosesan = 'diproses')"
			case "selesai":
				query += " AND (p.status_pemrosesan = 'selesai')"
			case "batal":
				query += " AND (p.status_pemrosesan = 'batal' OR p.status_pesanan = 'dibatalkan')"
			case "menunggu":
				query += " AND (p.status_pesanan = 'menunggu_pembayaran')"
			case "dibayar":
				query += " AND (p.status_pesanan = 'dibayar')"
			}
		}

		// 3. Filter Rentang Tanggal
		if filter.Rentang != "" && filter.Rentang != "Semua tanggal" {
			switch strings.ToLower(filter.Rentang) {
			case "hari ini":
				query += " AND DATE(p.dibuat_pada) = CURDATE()"
			case "7 hari terakhir":
				query += " AND p.dibuat_pada >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
			case "30 hari terakhir":
				query += " AND p.dibuat_pada >= DATE_SUB(NOW(), INTERVAL 30 DAY)"
			}
		}

		query += " ORDER BY p.id DESC"

		rows, err := l.DB.Query(query, args...)
		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var r PesananRingkasOwner
				var tipeCheckout, statusPesanan, statusPemrosesan string
				var waktuDibuat time.Time

				if errScan := rows.Scan(
					&r.ID, &r.NomorPesanan, &r.Kontak, &r.WhatsApp,
					&tipeCheckout, &statusPesanan, &statusPemrosesan,
					&r.Total, &waktuDibuat,
				); errScan == nil {
					r.WaktuDibuat = waktuDibuat
					r.Tanggal = waktuDibuat.Format("2006-01-02")
					r.TotalFormat = fmt.Sprintf("Rp%s", formatRibuan(r.Total))

					if tipeCheckout == "pickup" {
						r.Metode = "Ambil Sendiri"
					} else {
						r.Metode = "Kirim Kurir"
					}

					// Status Pembayaran
					if strings.Contains(strings.ToLower(statusPesanan), "bayar") && !strings.Contains(strings.ToLower(statusPesanan), "menunggu") {
						r.Pembayaran = "Dibayar"
					} else if strings.Contains(strings.ToLower(statusPesanan), "gagal") || strings.Contains(strings.ToLower(statusPesanan), "batal") {
						r.Pembayaran = "Gagal"
					} else {
						r.Pembayaran = "Menunggu"
					}

					// Status Pemrosesan
					switch strings.ToLower(statusPemrosesan) {
					case "diproses":
						r.Status = "Diproses"
					case "selesai":
						r.Status = "Selesai"
					case "batal":
						r.Status = "Batal"
					default:
						r.Status = "Baru"
					}

					r.ItemsRingkas = l.ambilRingkasanItems(r.NomorPesanan)
					hasil = append(hasil, r)
				}
			}
		}
	}

	// Gabungkan dengan data simulasi Figma jika filter cocok
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, sim := range l.dataSimulasi {
		// Pastikan tidak duplikat dengan database
		sudahAda := false
		for _, dbItem := range hasil {
			if dbItem.NomorPesanan == sim.NomorPesanan {
				sudahAda = true
				break
			}
		}
		if sudahAda {
			continue
		}

		if l.cocokFilter(sim, filter) {
			hasil = append(hasil, sim)
		}
	}

	return hasil, len(hasil), nil
}

// ambilRingkasanItems mengambil format ringkas varian dari snapshot_item_pesanan.
func (l *LayananPesananOwner) ambilRingkasanItems(nomorPesanan string) string {
	if l.DB == nil {
		return "Rengginang"
	}

	query := `
		SELECT nama_varian, jumlah 
		FROM snapshot_item_pesanan 
		WHERE nomor_pesanan = ?
	`
	rows, err := l.DB.Query(query, nomorPesanan)
	if err != nil {
		return "Rengginang"
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var varian string
		var jumlah int
		if errScan := rows.Scan(&varian, &jumlah); errScan == nil {
			items = append(items, fmt.Sprintf("%s x %d", varian, jumlah))
		}
	}

	if len(items) == 0 {
		return "Rengginang"
	}
	return strings.Join(items, " • ")
}

// cocokFilter mengecek apakah pesanan simulasi sesuai dengan parameter filter.
func (l *LayananPesananOwner) cocokFilter(p PesananRingkasOwner, filter FilterPesananOwner) bool {
	if filter.KataKunci != "" {
		k := strings.ToLower(filter.KataKunci)
		if !strings.Contains(strings.ToLower(p.NomorPesanan), k) &&
			!strings.Contains(strings.ToLower(p.Kontak), k) &&
			!strings.Contains(strings.ToLower(p.WhatsApp), k) {
			return false
		}
	}

	if filter.Status != "" && filter.Status != "Semua status" {
		st := strings.ToLower(filter.Status)
		cocokStatus := strings.EqualFold(p.Status, st)
		cocokBayar := strings.EqualFold(p.Pembayaran, st)
		if !cocokStatus && !cocokBayar {
			return false
		}
	}

	if filter.Rentang != "" && filter.Rentang != "Semua tanggal" {
		today := time.Now().Format("2006-01-02")
		if strings.EqualFold(filter.Rentang, "hari ini") && p.Tanggal != today {
			return false
		}
	}

	return true
}

// AmbilDetailPesanan mengambil rincian lengkap pesanan beserta detail barang untuk frame PPL21_02.
func (l *LayananPesananOwner) AmbilDetailPesanan(nomorPesanan string) (*DetailPesananOwner, error) {
	if strings.TrimSpace(nomorPesanan) == "" {
		return nil, fmt.Errorf("nomor pesanan tidak boleh kosong")
	}

	// Cek database MySQL
	if l.DB != nil {
		query := `
			SELECT 
				p.nomor_pesanan, p.nomor_invoice, p.nama_penerima, p.nomor_whatsapp,
				p.alamat_lengkap, p.kurir_nama, p.kurir_layanan, p.tipe_checkout,
				p.snapshot_subtotal, p.snapshot_ongkir, p.snapshot_total,
				p.status_pesanan, COALESCE(p.status_pemrosesan, 'menunggu_diproses'),
				p.dibuat_pada
			FROM pesanan p
			WHERE p.nomor_pesanan = ?
			LIMIT 1
		`
		var d DetailPesananOwner
		var kurirNama, kurirLayanan, tipeCheckout, statusPesanan, statusPemrosesan string
		var dibuatPada time.Time

		err := l.DB.QueryRow(query, nomorPesanan).Scan(
			&d.NomorPesanan, &d.NomorInvoice, &d.Nama, &d.WhatsApp,
			&d.Alamat, &kurirNama, &kurirLayanan, &tipeCheckout,
			&d.Subtotal, &d.Ongkir, &d.Total,
			&statusPesanan, &statusPemrosesan, &dibuatPada,
		)

		if err == nil {
			d.Tanggal = dibuatPada.Format("2006-01-02")
			d.SubtotalFormat = fmt.Sprintf("Rp%s", formatRibuan(d.Subtotal))
			d.OngkirFormat = fmt.Sprintf("Rp%s", formatRibuan(d.Ongkir))
			d.TotalFormat = fmt.Sprintf("Rp%s", formatRibuan(d.Total))

			if tipeCheckout == "pickup" {
				d.Metode = "Ambil Sendiri"
				d.Pengiriman = "Ambil di Toko"
			} else {
				d.Metode = "Kirim Kurir"
				d.Pengiriman = fmt.Sprintf("%s - %s", kurirNama, kurirLayanan)
			}

			if strings.Contains(strings.ToLower(statusPesanan), "bayar") && !strings.Contains(strings.ToLower(statusPesanan), "menunggu") {
				d.Pembayaran = "Dibayar"
			} else {
				d.Pembayaran = "Menunggu Pembayaran"
			}

			switch strings.ToLower(statusPemrosesan) {
			case "diproses":
				d.Pemrosesan = "Diproses"
			case "selesai":
				d.Pemrosesan = "Selesai"
			case "batal":
				d.Pemrosesan = "Batal"
			default:
				d.Pemrosesan = "Menunggu Diproses"
			}

			// Ambil snapshot item barang
			d.Items = l.ambilSnapshotBarang(nomorPesanan)
			var itemStrs []string
			for _, it := range d.Items {
				itemStrs = append(itemStrs, fmt.Sprintf("%s x %d", it.NamaVarian, it.JumlahBungkus))
			}
			if len(itemStrs) > 0 {
				d.ItemsRingkasan = strings.Join(itemStrs, " • ")
			} else {
				d.ItemsRingkasan = "Rengginang Putra Jaya"
			}

			return &d, nil
		}
	}

	// Cek data simulasi Figma
	l.mu.RLock()
	defer l.mu.RUnlock()

	if detail, ada := l.detailSimulasi[nomorPesanan]; ada {
		return detail, nil
	}

	// Fallback jika tidak ditemukan
	return nil, fmt.Errorf("pesanan %s tidak ditemukan", nomorPesanan)
}

// ambilSnapshotBarang mengambil detail snapshot produk dari tabel snapshot_item_pesanan.
func (l *LayananPesananOwner) ambilSnapshotBarang(nomorPesanan string) []ItemBarangSnapshot {
	var items []ItemBarangSnapshot
	if l.DB == nil {
		return items
	}

	query := `
		SELECT 
			id_varian, nama_produk, nama_varian, jumlah, 
			harga_satuan_snapshot, subtotal_snapshot
		FROM snapshot_item_pesanan
		WHERE nomor_pesanan = ?
		ORDER BY id ASC
	`
	rows, err := l.DB.Query(query, nomorPesanan)
	if err != nil {
		return items
	}
	defer rows.Close()

	for rows.Next() {
		var it ItemBarangSnapshot
		if errScan := rows.Scan(
			&it.IDVarian, &it.NamaProduk, &it.NamaVarian,
			&it.JumlahBungkus, &it.HargaSatuan, &it.Subtotal,
		); errScan == nil {
			it.HargaSatuanFormat = fmt.Sprintf("Rp%s", formatRibuan(it.HargaSatuan))
			it.SubtotalFormat = fmt.Sprintf("Rp%s", formatRibuan(it.Subtotal))
			items = append(items, it)
		}
	}

	return items
}

// PerbaruiStatusPemrosesan mengupdate status pemrosesan pesanan pada database dan memori.
func (l *LayananPesananOwner) PerbaruiStatusPemrosesan(nomorPesanan string, statusBaru string) error {
	if strings.TrimSpace(nomorPesanan) == "" {
		return fmt.Errorf("nomor pesanan wajib diisi")
	}

	statusBaruValid := strings.TrimSpace(statusBaru)
	if statusBaruValid == "" {
		return fmt.Errorf("status baru wajib diisi")
	}

	// Update di database jika tersedia
	if l.DB != nil {
		query := `UPDATE pesanan SET status_pemrosesan = ? WHERE nomor_pesanan = ?`
		_, err := l.DB.Exec(query, statusBaruValid, nomorPesanan)
		if err != nil {
			return fmt.Errorf("gagal memperbarui status di database: %v", err)
		}
	}

	// Update di memori simulasi
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := range l.dataSimulasi {
		if l.dataSimulasi[i].NomorPesanan == nomorPesanan {
			l.dataSimulasi[i].Status = statusBaruValid
			break
		}
	}

	if detail, ada := l.detailSimulasi[nomorPesanan]; ada {
		detail.Pemrosesan = statusBaruValid
	}

	return nil
}

// ============================================================================
// HANDLER HTTP UNTUK OWNER API
// ============================================================================

// HandlerDaftarPesanan (GET /api/admin/pesanan)
func (l *LayananPesananOwner) HandlerDaftarPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, `{"status":"error","pesan":"Metode harus GET"}`, http.StatusMethodNotAllowed)
		return
	}

	filter := FilterPesananOwner{
		KataKunci: r.URL.Query().Get("q"),
		Status:    r.URL.Query().Get("status"),
		Rentang:   r.URL.Query().Get("rentang"),
	}

	data, total, err := l.AmbilDaftarPesanan(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","pesan":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ResponsDaftarPesanan{
		Status: "sukses",
		Data:   data,
		Total:  total,
	})
}

// HandlerDetailPesanan (GET /api/admin/pesanan/detail)
func (l *LayananPesananOwner) HandlerDetailPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, `{"status":"error","pesan":"Metode harus GET"}`, http.StatusMethodNotAllowed)
		return
	}

	nomorPesanan := r.URL.Query().Get("nomor_pesanan")
	if nomorPesanan == "" {
		nomorPesanan = r.URL.Query().Get("id")
	}

	if nomorPesanan == "" {
		http.Error(w, `{"status":"error","pesan":"Parameter nomor_pesanan atau id wajib diisi"}`, http.StatusBadRequest)
		return
	}

	detail, err := l.AmbilDetailPesanan(nomorPesanan)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","pesan":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ResponsDetailPesanan{
		Status: "sukses",
		Data:   detail,
	})
}

// HandlerUpdateStatusPesanan (PUT / POST /api/admin/pesanan/status)
func (l *LayananPesananOwner) HandlerUpdateStatusPesanan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode harus PUT atau POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanUpdateStatus
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	if err := l.PerbaruiStatusPemrosesan(req.NomorPesanan, req.StatusPemrosesan); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","pesan":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"pesan":  fmt.Sprintf("Status pesanan %s berhasil diperbarui menjadi %s", req.NomorPesanan, req.StatusPemrosesan),
	})
}

// formatRibuan memformat integer menjadi format ribuan dengan titik (contoh: 28000 -> 28.000).
func formatRibuan(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var res []string
	remainder := len(s) % 3
	if remainder > 0 {
		res = append(res, s[:remainder])
	}
	for i := remainder; i < len(s); i += 3 {
		res = append(res, s[i:i+3])
	}
	return strings.Join(res, ".")
}
