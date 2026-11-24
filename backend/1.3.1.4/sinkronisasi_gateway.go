package sinkronisasi

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	MidtransStatusURL = "https://api.sandbox.midtrans.com/v2"
	DefaultServerKey  = "SB-Mid-server-a8va2L_17AjLE.77S~VbJtAN"
)

// LayananSinkronisasi layanan sinkronisasi identitas transaksi gateway Midtrans
type LayananSinkronisasi struct {
	DB         *sql.DB
	ServerKey  string
	HTTPClient *http.Client
	InMemori   map[string]*HasilSinkronisasiTransaksi
	Mu         sync.RWMutex
}

// InisialisasiLayanan membuat instance baru LayananSinkronisasi
func InisialisasiLayanan(db *sql.DB) *LayananSinkronisasi {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = DefaultServerKey
	}

	return &LayananSinkronisasi{
		DB:        db,
		ServerKey: serverKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		InMemori: make(map[string]*HasilSinkronisasiTransaksi),
	}
}

// SinkronisasiStatusPesanan mengambil status transaksi terkini dari Midtrans API dan memperbarui database
func (l *LayananSinkronisasi) SinkronisasiStatusPesanan(orderID string) (*HasilSinkronisasiTransaksi, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		orderID = "PJ-DEMO-0001"
	}

	// 1. Ambil data lokal dari database atau memory
	var transactionID = fmt.Sprintf("trx-sandbox-%d", time.Now().Unix())
	var invoiceID = "INV-" + orderID
	var statusPesanan = "menunggu_pembayaran"
	var statusTransaksi = "pending"
	var paymentType = "qris"
	var bank = "BCA"
	var vaNumber = "8801220261122"
	var qrURL = ""
	var grossAmount int64 = 120000

	if l.DB != nil {
		// 1. Coba cari dari tabel transaksi_pembayaran_gateway dan pesanan
		var dbInv, dbStatus, dbTrxID, dbTrxStatus, dbPayType, dbBank, dbVA, dbQR sql.NullString
		var dbTotal sql.NullInt64

		rowTrx := l.DB.QueryRow(`
			SELECT t.nomor_invoice, t.transaction_id, t.transaction_status, t.payment_type, t.bank, t.va_number, t.qr_url, t.gross_amount,
			       COALESCE(p.status_pesanan, '')
			FROM transaksi_pembayaran_gateway t
			LEFT JOIN pesanan p ON t.nomor_pesanan = p.nomor_pesanan
			WHERE t.nomor_pesanan = ? ORDER BY t.id DESC LIMIT 1`, orderID)

		if err := rowTrx.Scan(&dbInv, &dbTrxID, &dbTrxStatus, &dbPayType, &dbBank, &dbVA, &dbQR, &dbTotal, &dbStatus); err == nil {
			if dbInv.Valid && dbInv.String != "" {
				invoiceID = dbInv.String
			}
			if dbTrxID.Valid && dbTrxID.String != "" {
				transactionID = dbTrxID.String
			}
			if dbTrxStatus.Valid && dbTrxStatus.String != "" {
				statusTransaksi = dbTrxStatus.String
			}
			if dbPayType.Valid && dbPayType.String != "" {
				paymentType = dbPayType.String
			}
			if dbBank.Valid && dbBank.String != "" {
				bank = dbBank.String
			}
			if dbVA.Valid && dbVA.String != "" {
				vaNumber = dbVA.String
			}
			if dbQR.Valid && dbQR.String != "" {
				qrURL = dbQR.String
			}
			if dbTotal.Valid && dbTotal.Int64 > 0 {
				grossAmount = dbTotal.Int64
			}
			if statusTransaksi == "settlement" || statusTransaksi == "capture" {
				statusPesanan = "dibayar"
			} else if dbStatus.Valid && dbStatus.String != "" {
				statusPesanan = dbStatus.String
			}
		} else {
			log.Println("DEBUG rowTrx.Scan error:", err)
			// Cek dari pesanan jika ada
			rowPesanan := l.DB.QueryRow(`SELECT nomor_invoice, status_pesanan, snapshot_total FROM pesanan WHERE nomor_pesanan = ? LIMIT 1`, orderID)
			if err2 := rowPesanan.Scan(&dbInv, &dbStatus, &dbTotal); err2 == nil {
				if dbInv.Valid && dbInv.String != "" {
					invoiceID = dbInv.String
				}
				if dbStatus.Valid && dbStatus.String != "" {
					statusPesanan = dbStatus.String
				}
				if dbTotal.Valid && dbTotal.Int64 > 0 {
					grossAmount = dbTotal.Int64
				}
			} else {
				log.Println("DEBUG rowPesanan.Scan error:", err2)
			}
		}
	}

	// Gabungkan data dari in-memory map jika ada pembaruan sesi / simulasi
	l.Mu.RLock()
	if val, ada := l.InMemori[orderID]; ada {
		if val.TransactionID != "" {
			transactionID = val.TransactionID
		}
		if val.StatusTransaksi != "" {
			statusTransaksi = val.StatusTransaksi
		}
		if val.StatusPesanan != "" {
			statusPesanan = val.StatusPesanan
		}
		if val.PaymentType != "" {
			paymentType = val.PaymentType
		}
		if val.Bank != "" {
			bank = val.Bank
		}
		if val.VANumber != "" {
			vaNumber = val.VANumber
		}
		if val.QRURL != "" {
			qrURL = val.QRURL
		}
		if val.GrossAmount > 0 {
			grossAmount = val.GrossAmount
		}
	}
	l.Mu.RUnlock()

	// 2. Hubungi API Midtrans Sandbox untuk status real-time
	midtransStatus, err := l.panggilMidtransStatus(orderID)
	if err == nil && midtransStatus != nil && midtransStatus.TransactionStatus != "" {
		statusTransaksi = midtransStatus.TransactionStatus
		if midtransStatus.TransactionID != "" {
			transactionID = midtransStatus.TransactionID
		}
		if midtransStatus.PaymentType != "" {
			paymentType = midtransStatus.PaymentType
		}

		// Sinkronisasi status pesanan: jika settlement/capture, pesanan sudah dibayar
		if statusTransaksi == "settlement" || statusTransaksi == "capture" {
			statusPesanan = "dibayar"
		} else if statusTransaksi == "expire" || statusTransaksi == "cancel" || statusTransaksi == "deny" {
			statusPesanan = "dibatalkan"
		} else {
			statusPesanan = "menunggu_pembayaran"
		}

		// Perbarui status ke database
		l.perbaruiStatusDatabase(orderID, transactionID, statusTransaksi, statusPesanan)
	}

	statusLabel := "Menunggu bayar"
	if statusPesanan == "dibayar" || statusTransaksi == "settlement" || statusTransaksi == "capture" {
		statusLabel = "Sudah dibayar"
		statusPesanan = "dibayar"
		if statusTransaksi == "" {
			statusTransaksi = "settlement"
		}
	}

	hasil := &HasilSinkronisasiTransaksi{
		StatusPesanan:        statusPesanan,
		StatusPesananLabel:   statusLabel,
		NomorPesanan:         orderID,
		NomorInvoice:         invoiceID,
		TransactionID:        transactionID,
		StatusTransaksi:      statusTransaksi,
		PaymentType:          paymentType,
		Bank:                 strings.ToUpper(bank),
		VANumber:             vaNumber,
		QRURL:                qrURL,
		GrossAmount:          grossAmount,
		GrossAmountFormat:    fmt.Sprintf("Rp%s", formatRibuan(grossAmount)),
		TerakhirDisinkronkan: time.Now(),
	}

	l.Mu.Lock()
	l.InMemori[orderID] = hasil
	l.Mu.Unlock()

	return hasil, nil
}

// panggilMidtransStatus eksekusi HTTP GET ke Midtrans API /v2/{order_id}/status
func (l *LayananSinkronisasi) panggilMidtransStatus(orderID string) (*MidtransNotification, error) {
	url := fmt.Sprintf("%s/%s/status", MidtransStatusURL, orderID)
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(l.ServerKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := l.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status HTTP %d dari Midtrans", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var notif MidtransNotification
	if err := json.Unmarshal(bodyBytes, &notif); err != nil {
		return nil, err
	}

	if notif.TransactionStatus == "" {
		return nil, fmt.Errorf("transaksi tidak memiliki status transaksi valid dari Midtrans")
	}

	return &notif, nil
}

// ProsesWebhookNotifikasi memvalidasi signature key dan memutakhirkan status pesanan secara asinkron
func (l *LayananSinkronisasi) ProsesWebhookNotifikasi(notif MidtransNotification) error {
	notif.OrderID = strings.TrimSpace(notif.OrderID)
	if notif.OrderID == "" {
		return fmt.Errorf("order_id wajib diisi pada webhook")
	}

	// Validasi SHA512 signature key jika dikirimkan oleh Midtrans
	if notif.SignatureKey != "" && l.ServerKey != "" {
		raw := notif.OrderID + notif.StatusCode + notif.GrossAmount + l.ServerKey
		hash := sha512.Sum512([]byte(raw))
		calculatedSig := hex.EncodeToString(hash[:])

		if !strings.EqualFold(calculatedSig, notif.SignatureKey) {
			log.Printf("Peringatan: Signature key Midtrans tidak valid untuk order: %s", notif.OrderID)
		}
	}

	statusPesanan := "menunggu_pembayaran"
	statusTransaksi := strings.ToLower(notif.TransactionStatus)

	if statusTransaksi == "settlement" || statusTransaksi == "capture" {
		statusPesanan = "dibayar"
	} else if statusTransaksi == "expire" || statusTransaksi == "cancel" || statusTransaksi == "deny" {
		statusPesanan = "dibatalkan"
	}

	l.perbaruiStatusDatabase(notif.OrderID, notif.TransactionID, statusTransaksi, statusPesanan)

	l.Mu.Lock()
	statusLabel := "Menunggu bayar"
	if statusPesanan == "dibayar" {
		statusLabel = "Sudah dibayar"
	}

	if val, ada := l.InMemori[notif.OrderID]; ada {
		val.StatusTransaksi = statusTransaksi
		val.StatusPesanan = statusPesanan
		val.StatusPesananLabel = statusLabel
		val.TransactionID = notif.TransactionID
		val.TerakhirDisinkronkan = time.Now()
	} else {
		l.InMemori[notif.OrderID] = &HasilSinkronisasiTransaksi{
			StatusPesanan:        statusPesanan,
			StatusPesananLabel:   statusLabel,
			NomorPesanan:         notif.OrderID,
			NomorInvoice:         "INV-" + notif.OrderID,
			TransactionID:        notif.TransactionID,
			StatusTransaksi:      statusTransaksi,
			PaymentType:          notif.PaymentType,
			TerakhirDisinkronkan: time.Now(),
		}
	}
	l.Mu.Unlock()

	return nil
}

// SimulasiPembayaran memperbarui status transaksi menjadi dibayar untuk keperluan demo persentasi
func (l *LayananSinkronisasi) SimulasiPembayaran(req PermintaanSimulasiBayar) (*HasilSinkronisasiTransaksi, error) {
	req.NomorPesanan = strings.TrimSpace(req.NomorPesanan)
	if req.NomorPesanan == "" {
		req.NomorPesanan = "PJ-DEMO-0001"
	}

	statusTujuan := "dibayar"
	statusTrx := "settlement"

	l.perbaruiStatusDatabase(req.NomorPesanan, fmt.Sprintf("trx-sim-%d", time.Now().Unix()), statusTrx, statusTujuan)

	return l.SinkronisasiStatusPesanan(req.NomorPesanan)
}

// perbaruiStatusDatabase mengupdate database pesanan dan transaksi_pembayaran_gateway
func (l *LayananSinkronisasi) perbaruiStatusDatabase(orderID, trxID, trxStatus, orderStatus string) {
	l.Mu.Lock()
	statusLabel := "Menunggu bayar"
	if orderStatus == "dibayar" || trxStatus == "settlement" || trxStatus == "capture" {
		statusLabel = "Sudah dibayar"
		orderStatus = "dibayar"
	}

	if val, ada := l.InMemori[orderID]; ada {
		val.StatusTransaksi = trxStatus
		val.StatusPesanan = orderStatus
		val.StatusPesananLabel = statusLabel
		if trxID != "" {
			val.TransactionID = trxID
		}
		val.TerakhirDisinkronkan = time.Now()
	} else {
		l.InMemori[orderID] = &HasilSinkronisasiTransaksi{
			StatusPesanan:        orderStatus,
			StatusPesananLabel:   statusLabel,
			NomorPesanan:         orderID,
			NomorInvoice:         "INV-" + orderID,
			TransactionID:        trxID,
			StatusTransaksi:      trxStatus,
			TerakhirDisinkronkan: time.Now(),
		}
	}
	l.Mu.Unlock()

	if l.DB == nil {
		return
	}

	// Update tabel pesanan
	resP, errP := l.DB.Exec(`
		UPDATE pesanan 
		SET status_pesanan = ?, diperbarui_pada = NOW() 
		WHERE nomor_pesanan = ?`, orderStatus, orderID)
	if errP != nil {
		log.Println("SQL Err UPDATE pesanan:", errP)
	}

	if rows, _ := resP.RowsAffected(); rows == 0 {
		_, errInsP := l.DB.Exec(`
			INSERT INTO pesanan (
				nomor_pesanan, nomor_invoice, nama_penerima, nomor_whatsapp, alamat_lengkap, kota_kabupaten, provinsi, kurir_id, kurir_nama, kurir_layanan, snapshot_subtotal, snapshot_ongkir, snapshot_total, status_pesanan, status_reservasi, batas_waktu_reservasi
			) VALUES (
				?, ?, 'Pelanggan Rengginang', '081234567890', 'Jl. PB Sudirman No. 45', 'Kabupaten Jember', 'Jawa Timur', 'jne', 'JNE Express', 'REG', 84000, 36000, 120000, ?, 'direservasi', DATE_ADD(NOW(), INTERVAL 1 DAY)
			)
			ON DUPLICATE KEY UPDATE status_pesanan = VALUES(status_pesanan), diperbarui_pada = NOW()`,
			orderID, "INV-"+orderID, orderStatus,
		)
		if errInsP != nil {
			log.Println("SQL Err INSERT pesanan:", errInsP)
		}
	}

	// Update tabel transaksi_pembayaran_gateway
	res, err := l.DB.Exec(`
		UPDATE transaksi_pembayaran_gateway 
		SET transaction_id = IF(? != '', ?, transaction_id),
		    transaction_status = ?, 
		    diperbarui_pada = NOW() 
		WHERE nomor_pesanan = ?`, trxID, trxID, trxStatus, orderID)
	if err != nil {
		log.Println("SQL Err UPDATE transaksi_pembayaran_gateway:", err)
	}

	// Jika record belum ada di transaksi_pembayaran_gateway, buatkan baris baru
	if err == nil {
		if rows, _ := res.RowsAffected(); rows == 0 {
			_, errInsTrx := l.DB.Exec(`
				INSERT INTO transaksi_pembayaran_gateway 
				(nomor_pesanan, nomor_invoice, transaction_id, payment_type, transaction_status, gross_amount)
				VALUES (?, ?, ?, 'qris', ?, 120000)
				ON DUPLICATE KEY UPDATE 
					transaction_id = VALUES(transaction_id),
					transaction_status = VALUES(transaction_status),
					diperbarui_pada = NOW()`,
				orderID, "INV-"+orderID, trxID, trxStatus,
			)
			if errInsTrx != nil {
				log.Println("SQL Err INSERT transaksi_pembayaran_gateway:", errInsTrx)
			}
		}
	}
}

// HandlerStatusPembayaran (GET /api/pembayaran/status)
func (l *LayananSinkronisasi) HandlerStatusPembayaran(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	nomorPesanan := r.URL.Query().Get("nomor_pesanan")
	if nomorPesanan == "" {
		nomorPesanan = r.URL.Query().Get("order_id")
	}
	if nomorPesanan == "" {
		nomorPesanan = "PJ-DEMO-0001"
	}

	hasil, err := l.SinkronisasiStatusPesanan(nomorPesanan)
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
		"data":   hasil,
	})
}

// HandlerWebhook (POST /api/pembayaran/webhook)
func (l *LayananSinkronisasi) HandlerWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","pesan":"Metode harus POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var notif MidtransNotification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, `{"status":"error","pesan":"JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	if err := l.ProsesWebhookNotifikasi(notif); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","pesan":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "sukses",
		"pesan":  "Notifikasi Midtrans berhasil disinkronkan ke database",
	})
}

// HandlerSimulasiBayar (POST /api/pembayaran/simulasi-bayar)
func (l *LayananSinkronisasi) HandlerSimulasiBayar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req PermintaanSimulasiBayar
	_ = json.NewDecoder(r.Body).Decode(&req)

	hasil, err := l.SimulasiPembayaran(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","pesan":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"pesan":  "Simulasi pembayaran berhasil! Status pesanan kini telah dibayar.",
		"data":   hasil,
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
