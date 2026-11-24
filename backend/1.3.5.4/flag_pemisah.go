package flagpemisah

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

// LayananPemisahFlag mengelola validasi dan pemisahan tegas antara transaksi nyata dan simulasi sandbox
type LayananPemisahFlag struct {
	DB *sql.DB
	Mu sync.RWMutex
}

// InisialisasiLayanan membuat instance layanan pemisah flag
func InisialisasiLayanan(db *sql.DB) *LayananPemisahFlag {
	return &LayananPemisahFlag{
		DB: db,
	}
}

// EvaluasiFlag menentukan apakah transaksi berkategori nyata (offline kasir/tunai) atau simulasi sandbox
func (l *LayananPemisahFlag) EvaluasiFlag(sumber, flag string) (string, bool) {
	sumberClean := strings.ToLower(strings.TrimSpace(sumber))
	flagClean := strings.ToUpper(strings.TrimSpace(flag))

	// Transaksi kasir offline dan pembayaran tunai selalu PRODUKSI_NYATA (uang fisik toko)
	if sumberClean == SumberOfflineKasir || flagClean == FlagProduksiNyata || sumberClean == "tunai" {
		return FlagProduksiNyata, false
	}

	// Pembayaran Midtrans sandbox dan simulasi testing otomatis SANDBOX_SIMULASI
	if strings.Contains(sumberClean, "sandbox") || strings.Contains(sumberClean, "simulasi") || flagClean == FlagSandboxSimulasi {
		return FlagSandboxSimulasi, true
	}

	// Default fallback aman
	return FlagProduksiNyata, false
}

// AuditPemisahTransaksi mengumpulkan rekap perbandingan transaksi uang kas nyata vs simulasi dummy
func (l *LayananPemisahFlag) AuditPemisahTransaksi() (*RekapAuditFlag, error) {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	rekap := &RekapAuditFlag{
		DaftarPemisah: make([]InfoPemisahTransaksi, 0),
	}

	// 1. Ambil data dari tabel penjualan_offline (Semua berstatus PRODUKSI_NYATA / uang kas kasir)
	if l.DB != nil {
		rowsOffline, err := l.DB.Query(`
			SELECT nomor_transaksi, total_harga, metode_pembayaran, flag_lingkungan, is_sandbox, dibuat_pada
			FROM penjualan_offline ORDER BY id DESC LIMIT 20`)
		if err == nil {
			defer rowsOffline.Close()
			for rowsOffline.Next() {
				var noTrx, metode, flag string
				var total float64
				var isSbInt int
				var t time.Time
				if errScan := rowsOffline.Scan(&noTrx, &total, &metode, &flag, &isSbInt, &t); errScan == nil {
					isSb := isSbInt == 1
					rekap.TotalTransaksiNyata++
					totalInt := int64(total)
					rekap.TotalNominalNyata += totalInt

					rekap.DaftarPemisah = append(rekap.DaftarPemisah, InfoPemisahTransaksi{
						IDTransaksi:      noTrx,
						NomorReferensi:   noTrx,
						Sumber:           SumberOfflineKasir,
						FlagLingkungan:   FlagProduksiNyata,
						IsSandbox:        isSb,
						TotalNominal:     totalInt,
						TotalFormat:      fmt.Sprintf("Rp%s", formatRibuan(totalInt)),
						MetodePembayaran: metode,
						Keterangan:       "Penjualan Kasir Langsung di Toko (Uang Nyata)",
						WaktuTransaksi:   t,
					})
				}
			}
		}

		// 2. Ambil data dari tabel transaksi_pembayaran_gateway (Midtrans Sandbox)
		rowsSandbox, err := l.DB.Query(`
			SELECT nomor_pesanan, gross_amount, payment_type, transaction_status, dibuat_pada
			FROM transaksi_pembayaran_gateway ORDER BY id DESC LIMIT 20`)
		if err == nil {
			defer rowsSandbox.Close()
			for rowsSandbox.Next() {
				var noOrder, payType, trxStatus string
				var gross int64
				var t time.Time
				if errScan := rowsSandbox.Scan(&noOrder, &gross, &payType, &trxStatus, &t); errScan == nil {
					rekap.TotalTransaksiSandbox++
					rekap.TotalNominalSandbox += gross

					rekap.DaftarPemisah = append(rekap.DaftarPemisah, InfoPemisahTransaksi{
						IDTransaksi:      "SANDBOX-" + noOrder,
						NomorReferensi:   noOrder,
						Sumber:           SumberSandbox,
						FlagLingkungan:   FlagSandboxSimulasi,
						IsSandbox:        true,
						TotalNominal:     gross,
						TotalFormat:      fmt.Sprintf("Rp%s", formatRibuan(gross)),
						MetodePembayaran: payType,
						Keterangan:       fmt.Sprintf("Simulasi Payment Gateway Midtrans Sandbox (%s)", trxStatus),
						WaktuTransaksi:   t,
					})
				}
			}
		}
	}

	rekap.TotalNominalNyataFmt = fmt.Sprintf("Rp%s", formatRibuan(rekap.TotalNominalNyata))
	rekap.TotalNominalSandboxFmt = fmt.Sprintf("Rp%s", formatRibuan(rekap.TotalNominalSandbox))

	return rekap, nil
}

// HandlerAudit (GET /api/admin/transaksi/audit-flag)
func (l *LayananPemisahFlag) HandlerAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rekap, err := l.AuditPemisahTransaksi()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"pesan":  "Data pemisahan transaksi nyata dan sandbox berhasil diproses",
		"data":   rekap,
	})
}

// HandlerValidasiFlag (POST /api/admin/transaksi/validasi-flag)
func (l *LayananPemisahFlag) HandlerValidasiFlag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanValidasiFlag
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"pesan": "Format JSON tidak valid"})
		return
	}

	flag, isSandbox := l.EvaluasiFlag(req.SumberTransaksi, req.FlagLingkungan)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ResponsValidasiFlag{
		Valid:          true,
		Pesan:          fmt.Sprintf("Transaksi berhasil diklasifikasikan ke %s", flag),
		FlagLingkungan: flag,
		IsSandbox:      isSandbox,
	})
}

func formatRibuan(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var res []string
	for len(s) > 3 {
		res = append([]string{s[len(s)-3:]}, res...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		res = append([]string{s}, res...)
	}
	return strings.Join(res, ".")
}
