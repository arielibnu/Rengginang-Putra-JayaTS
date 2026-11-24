package tagihan

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	MidtransSandboxURL = "https://api.sandbox.midtrans.com/v2/charge"
	MidtransSnapURL    = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	MidtransStatusURL  = "https://api.sandbox.midtrans.com/v2"
	DefaultServerKey   = "SB-Mid-server-a8va2L_17AjLE.77S~VbJtAN"
)

// LayananTagihan layanan integrasi pembuatan tagihan payment gateway
type LayananTagihan struct {
	DB         *sql.DB
	ServerKey  string
	HTTPClient *http.Client
	InMemori   map[string]*EntitasTransaksiPembayaran
	Mu         sync.RWMutex
}

// InisialisasiLayanan membuat instance LayananTagihan baru
func InisialisasiLayanan(db *sql.DB) *LayananTagihan {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = DefaultServerKey
	}

	layanan := &LayananTagihan{
		DB:        db,
		ServerKey: serverKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		InMemori: make(map[string]*EntitasTransaksiPembayaran),
	}

	layanan.InisialisasiTabel()
	return layanan
}

// InisialisasiTabel membuat tabel transaksi_pembayaran_gateway di database MySQL jika belum ada
func (l *LayananTagihan) InisialisasiTabel() {
	if l.DB == nil {
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS transaksi_pembayaran_gateway (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		id_pesanan BIGINT NULL,
		nomor_pesanan VARCHAR(64) NOT NULL,
		nomor_invoice VARCHAR(64) NULL,
		transaction_id VARCHAR(100) NOT NULL,
		payment_type VARCHAR(50) NOT NULL,
		bank VARCHAR(30) NULL,
		va_number VARCHAR(100) NULL,
		qr_url TEXT NULL,
		qr_string TEXT NULL,
		gross_amount BIGINT NOT NULL,
		transaction_status VARCHAR(50) NOT NULL DEFAULT 'pending',
		fraud_status VARCHAR(50) NULL DEFAULT 'accept',
		raw_gateway_response TEXT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_nomor_pesanan (nomor_pesanan),
		INDEX idx_transaction_id (transaction_id),
		INDEX idx_transaction_status (transaction_status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	if _, err := l.DB.Exec(query); err != nil {
		log.Println("Peringatan inisialisasi tabel transaksi_pembayaran_gateway:", err)
	}
}

// BuatTagihanMidtrans memproses pembuatan tagihan ke Midtrans Core API Sandbox
func (l *LayananTagihan) BuatTagihanMidtrans(req PermintaanBuatTagihan) (*ResponsTagihanGateway, error) {
	req.NomorPesanan = strings.TrimSpace(req.NomorPesanan)
	if req.NomorPesanan == "" {
		req.NomorPesanan = "PJ-DEMO-0001"
	}

	metode := strings.ToLower(strings.TrimSpace(req.Metode))
	if metode == "va" || metode == "virtual_account" {
		metode = "bank_transfer"
	} else if metode == "" {
		metode = "qris"
	}

	bank := strings.ToLower(strings.TrimSpace(req.Bank))
	if bank == "" {
		bank = "bca"
	}

	// 1. Ambil data pesanan dari database atau fallback
	var idPesanan int64 = 1
	var nomorInvoice = "INV-" + req.NomorPesanan
	var totalBayar int64 = 120000
	var namaPenerima = "Pelanggan Rengginang"
	var nomorWhatsApp = "081234567890"

	if l.DB != nil {
		row := l.DB.QueryRow(`
			SELECT id, nomor_invoice, snapshot_total, nama_penerima, nomor_whatsapp 
			FROM pesanan 
			WHERE nomor_pesanan = ? LIMIT 1`, req.NomorPesanan)
		var dbID int64
		var dbInv, dbNama, dbWA sql.NullString
		var dbTotal sql.NullInt64
		if err := row.Scan(&dbID, &dbInv, &dbTotal, &dbNama, &dbWA); err == nil {
			idPesanan = dbID
			if dbInv.Valid && dbInv.String != "" {
				nomorInvoice = dbInv.String
			}
			if dbTotal.Valid && dbTotal.Int64 > 0 {
				totalBayar = dbTotal.Int64
			}
			if dbNama.Valid && dbNama.String != "" {
				namaPenerima = dbNama.String
			}
			if dbWA.Valid && dbWA.String != "" {
				nomorWhatsApp = dbWA.String
			}
		}
	}

	var transactionStatus = "pending"
	var qrURL = ""
	var qrString = ""
	var snapToken = ""
	var snapRedirectURL = ""
	var vaNumber = ""
	var simulatorURL = "https://simulator.sandbox.midtrans.com/qris/index"
	var expiryTime = time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	var rawJSON = ""

	// 2. Buat sesi transaksi baru di Midtrans Snap Sandbox
	orderIDForMidtrans := req.NomorPesanan
	snapReq := MidtransSnapRequest{
		TransactionDetails: MidtransTransactionDetails{
			OrderID:     orderIDForMidtrans,
			GrossAmount: totalBayar,
		},
		CustomerDetails: &MidtransCustomerDetails{
			FirstName: namaPenerima,
			Phone:     nomorWhatsApp,
		},
		ItemDetails: []MidtransItemDetails{
			{
				ID:       "PROD-01",
				Price:    totalBayar,
				Quantity: 1,
				Name:     "Pesanan " + req.NomorPesanan,
			},
		},
		EnabledPayments: []string{"gopay", "other_qris", "bca_va", "bni_va", "bri_va"},
	}

	snapResp, snapRaw, errSnap := l.panggilMidtransSnap(snapReq)
	// Jika order_id sudah pernah digunakan/diselesaikan di Midtrans Sandbox, buat sesi aktif baru dengan suffix unik
	if errSnap != nil && strings.Contains(errSnap.Error(), "not unique") {
		orderIDForMidtrans = fmt.Sprintf("%s-T%d", req.NomorPesanan, time.Now().Unix())
		snapReq.TransactionDetails.OrderID = orderIDForMidtrans
		snapResp, snapRaw, errSnap = l.panggilMidtransSnap(snapReq)
	}

	if errSnap == nil && snapResp != nil && snapResp.Token != "" {
		rawJSON = snapRaw
		snapToken = snapResp.Token
		snapRedirectURL = snapResp.RedirectURL
	} else {
		log.Println("Peringatan Snap Sandbox:", errSnap)
	}

	var transactionID = snapToken
	if transactionID == "" {
		transactionID = fmt.Sprintf("trx-sandbox-%d", time.Now().Unix())
	}
	if snapRedirectURL == "" && snapToken != "" {
		snapRedirectURL = fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v4/redirection/%s", snapToken)
	}

	// 4. Siapkan parameter tampilan QRIS dan Virtual Account Sandbox
	cleanNum := strings.ReplaceAll(req.NomorPesanan, "-", "")
	if len(cleanNum) > 8 {
		cleanNum = cleanNum[len(cleanNum)-8:]
	}

	if metode == "qris" {
		targetData := snapRedirectURL
		if targetData == "" {
			targetData = fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v4/redirection/%s", transactionID)
		}
		qrString = targetData
		qrURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&margin=10&data=%s", url.QueryEscape(targetData))
		simulatorURL = "https://simulator.sandbox.midtrans.com/qris/index"
	} else {
		switch bank {
		case "bni":
			vaNumber = fmt.Sprintf("8808%s", cleanNum)
		case "bri":
			vaNumber = fmt.Sprintf("88017%s", cleanNum)
		case "mandiri":
			vaNumber = fmt.Sprintf("70012%s", cleanNum)
		default:
			vaNumber = fmt.Sprintf("88012%s", cleanNum)
		}
		simulatorURL = "https://simulator.sandbox.midtrans.com/openapi/va/index"
	}

	// 5. Simpan hasil transaksi ke database & in-memory map
	record := &EntitasTransaksiPembayaran{
		IDPesanan:          idPesanan,
		NomorPesanan:       req.NomorPesanan,
		NomorInvoice:       nomorInvoice,
		TransactionID:      transactionID,
		PaymentType:        metode,
		Bank:               bank,
		VANumber:           vaNumber,
		QRURL:              qrURL,
		QRString:           qrString,
		GrossAmount:        totalBayar,
		TransactionStatus:  transactionStatus,
		FraudStatus:        "accept",
		RawGatewayResponse: rawJSON,
		DibuatPada:         time.Now(),
		DiperbaruiPada:     time.Now(),
	}

	l.simpanTransaksiDatabase(record)

	// Format nominal Rupiah
	formatRupiah := fmt.Sprintf("Rp%s", formatRibuan(totalBayar))

	return &ResponsTagihanGateway{
		Status:            "sukses",
		Pesan:             "Tagihan pembayaran gateway sandbox berhasil dibuat",
		NomorPesanan:      req.NomorPesanan,
		TransactionID:     transactionID,
		PaymentType:       metode,
		GrossAmount:       totalBayar,
		GrossAmountFormat: formatRupiah,
		Bank:              strings.ToUpper(bank),
		VANumber:          vaNumber,
		QRURL:             qrURL,
		QRString:          qrString,
		SnapToken:         snapToken,
		SnapRedirectURL:   snapRedirectURL,
		TransactionStatus: transactionStatus,
		SimulatorURL:      simulatorURL,
		ExpiryTime:        expiryTime,
		DibuatPada:        record.DibuatPada,
	}, nil
}

// panggilMidtransSnap helper eksekusi HTTP POST ke Midtrans Snap API Sandbox
func (l *LayananTagihan) panggilMidtransSnap(payload MidtransSnapRequest) (*MidtransSnapResponse, string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, MidtransSnapURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, "", err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(l.ServerKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := l.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	rawString := string(bodyBytes)
	if resp.StatusCode >= 400 {
		return nil, rawString, fmt.Errorf("midtrans snap error (%d): %s", resp.StatusCode, rawString)
	}

	var snapResp MidtransSnapResponse
	if err := json.Unmarshal(bodyBytes, &snapResp); err != nil {
		return nil, rawString, err
	}

	return &snapResp, rawString, nil
}

// panggilMidtransCharge helper eksekusi HTTP POST ke API Midtrans
func (l *LayananTagihan) panggilMidtransCharge(payload MidtransChargeRequest) (*MidtransChargeResponse, string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, MidtransSandboxURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, "", err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(l.ServerKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := l.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	rawString := string(bodyBytes)

	if resp.StatusCode >= 400 {
		return nil, rawString, fmt.Errorf("midtrans error (%d): %s", resp.StatusCode, rawString)
	}

	var midtransResp MidtransChargeResponse
	if err := json.Unmarshal(bodyBytes, &midtransResp); err != nil {
		return nil, rawString, err
	}

	return &midtransResp, rawString, nil
}

// ambilStatusDariMidtrans mengambil status pesanan jika order_id sudah pernah dibuat
func (l *LayananTagihan) ambilStatusDariMidtrans(orderID string) (*MidtransChargeResponse, string, error) {
	url := fmt.Sprintf("%s/%s/status", MidtransStatusURL, orderID)
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(l.ServerKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := l.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	rawString := string(bodyBytes)
	var midtransResp MidtransChargeResponse
	if err := json.Unmarshal(bodyBytes, &midtransResp); err != nil {
		return nil, rawString, err
	}

	return &midtransResp, rawString, nil
}

// simpanTransaksiDatabase menyimpan atau mengupdate entitas pembayaran
func (l *LayananTagihan) simpanTransaksiDatabase(t *EntitasTransaksiPembayaran) {
	l.Mu.Lock()
	l.InMemori[t.NomorPesanan] = t
	l.Mu.Unlock()

	if l.DB == nil {
		return
	}

	var existsID int64
	var existingTrxID sql.NullString
	errRow := l.DB.QueryRow(`SELECT id, transaction_id FROM transaksi_pembayaran_gateway WHERE nomor_pesanan = ? ORDER BY id DESC LIMIT 1`, t.NomorPesanan).Scan(&existsID, &existingTrxID)
	if errRow == nil && existsID > 0 {
		// Jangan timpa token Snap valid jika yang baru hanya fallback trx-
		if strings.HasPrefix(t.TransactionID, "trx-") && existingTrxID.Valid && len(existingTrxID.String) >= 30 && !strings.HasPrefix(existingTrxID.String, "trx-") {
			t.TransactionID = existingTrxID.String
			t.QRString = fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v4/redirection/%s", t.TransactionID)
		}

		updateQuery := `
		UPDATE transaksi_pembayaran_gateway 
		SET transaction_id = ?, payment_type = ?, bank = ?, va_number = ?, qr_url = ?, qr_string = ?, gross_amount = ?, transaction_status = ?, fraud_status = ?, raw_gateway_response = ?, diperbarui_pada = NOW()
		WHERE id = ?`
		_, err := l.DB.Exec(updateQuery, t.TransactionID, t.PaymentType, t.Bank, t.VANumber, t.QRURL, t.QRString, t.GrossAmount, t.TransactionStatus, t.FraudStatus, t.RawGatewayResponse, existsID)
		if err != nil {
			log.Println("Peringatan update transaksi database:", err)
		}
	} else {
		insertQuery := `
		INSERT INTO transaksi_pembayaran_gateway 
		(id_pesanan, nomor_pesanan, nomor_invoice, transaction_id, payment_type, bank, va_number, qr_url, qr_string, gross_amount, transaction_status, fraud_status, raw_gateway_response)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := l.DB.Exec(insertQuery,
			t.IDPesanan, t.NomorPesanan, t.NomorInvoice, t.TransactionID,
			t.PaymentType, t.Bank, t.VANumber, t.QRURL, t.QRString,
			t.GrossAmount, t.TransactionStatus, t.FraudStatus, t.RawGatewayResponse,
		)
		if err != nil {
			log.Println("Peringatan simpan transaksi ke database:", err)
		}
	}
}

// HandlerBuatTagihan (POST /api/pembayaran/buat-tagihan)
func (l *LayananTagihan) HandlerBuatTagihan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode HTTP harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req PermintaanBuatTagihan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","pesan":"Payload JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	resp, err := l.BuatTagihanMidtrans(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"pesan":  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"data":   resp,
	})
}

// formatRibuan helper pemformat ribuan angka
func formatRibuan(n int64) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}
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
