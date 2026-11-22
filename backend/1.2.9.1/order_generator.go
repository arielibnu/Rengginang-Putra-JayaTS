package ordergenerator

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	memoriPesananPPL18 = make([]EntitasPesanan, 0)
	muPesanan          sync.RWMutex
)

// GeneratorOrder layanan pembuatan ID pesanan unik dan penyimpanan snapshot transaksi
type GeneratorOrder struct {
	DB *sql.DB
}

// InisialisasiTabel membuat tabel pesanan dan snapshot jika belum ada di MySQL
func (g *GeneratorOrder) InisialisasiTabel() {
	if g.DB == nil {
		return
	}

	queryPesanan := `
	CREATE TABLE IF NOT EXISTS pesanan (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		nomor_pesanan VARCHAR(64) NOT NULL UNIQUE,
		nomor_invoice VARCHAR(64) NOT NULL UNIQUE,
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
		snapshot_subtotal INT NOT NULL,
		snapshot_ongkir INT NOT NULL,
		snapshot_total INT NOT NULL,
		status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu_pembayaran',
		status_reservasi VARCHAR(30) NOT NULL DEFAULT 'direservasi',
		batas_waktu_reservasi DATETIME NOT NULL,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		diperbarui_pada DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_nomor_pesanan (nomor_pesanan),
		INDEX idx_nomor_invoice (nomor_invoice),
		INDEX idx_status_pesanan (status_pesanan)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	querySnapshot := `
	CREATE TABLE IF NOT EXISTS snapshot_item_pesanan (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		id_pesanan BIGINT NOT NULL,
		nomor_pesanan VARCHAR(64) NOT NULL,
		id_varian BIGINT NOT NULL DEFAULT 1,
		id_produk VARCHAR(50) NOT NULL,
		nama_produk VARCHAR(150) NOT NULL,
		nama_varian VARCHAR(100) NOT NULL,
		jumlah INT NOT NULL,
		harga_satuan_snapshot INT NOT NULL,
		subtotal_snapshot INT NOT NULL,
		berat_satuan_gram INT NOT NULL DEFAULT 500,
		dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_pesanan_snapshot (nomor_pesanan)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, _ = g.DB.Exec(queryPesanan)
	_, _ = g.DB.Exec(querySnapshot)
}

// GenerateNomorPesananUnik membuat format ID pesanan resmi: PJ-YYYYMMDD-XXXX
func GenerateNomorPesananUnik() string {
	tanggal := time.Now().Format("20060102")
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	kodeAcak := 1000 + randomNum.Int64()
	return fmt.Sprintf("PJ-%s-%04d", tanggal, kodeAcak)
}

// GenerateNomorInvoiceUnik membuat nomor invoice resmi: INV-YYYYMMDD-XXXX
func GenerateNomorInvoiceUnik() string {
	tanggal := time.Now().Format("20060102")
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	kodeAcak := 1000 + randomNum.Int64()
	return fmt.Sprintf("INV-%s-%04d", tanggal, kodeAcak)
}

// BuatSnapshotDanSimpanOrder mengeksekusi penyimpanan snapshot harga dan pembuatan order
func (g *GeneratorOrder) BuatSnapshotDanSimpanOrder(
	tipeCheckout string,
	idPelanggan *uint64,
	kontak DataKontakPengiriman,
	kurir DataKurirSnapshot,
	items []ItemSnapshotPesanan,
) (*EntitasPesanan, error) {
	now := time.Now()
	nomorPesanan := GenerateNomorPesananUnik()
	nomorInvoice := GenerateNomorInvoiceUnik()
	batasWaktu := now.Add(24 * time.Hour) // Holding time 24 jam

	// Hitung snapshot subtotal
	subtotalSnapshot := 0
	for i := range items {
		if items[i].HargaSatuanSnapshot <= 0 {
			items[i].HargaSatuanSnapshot = 28000
		}
		items[i].SubtotalSnapshot = items[i].Jumlah * items[i].HargaSatuanSnapshot
		subtotalSnapshot += items[i].SubtotalSnapshot
	}

	totalSnapshot := subtotalSnapshot + kurir.TarifOngkir

	pesanan := &EntitasPesanan{
		NomorPesanan:        nomorPesanan,
		NomorInvoice:        nomorInvoice,
		TipeCheckout:        tipeCheckout,
		IDPelanggan:         idPelanggan,
		Kontak:              kontak,
		Kurir:               kurir,
		Items:               items,
		SnapshotSubtotal:    subtotalSnapshot,
		SnapshotOngkir:      kurir.TarifOngkir,
		SnapshotTotal:       totalSnapshot,
		StatusPesanan:       "menunggu_pembayaran",
		StatusReservasi:     "direservasi",
		BatasWaktuReservasi: batasWaktu,
		BatasWaktuFormat:    batasWaktu.Format("02 Jan 2006, 15:04 WIB"),
		DibuatPada:          now,
		DiperbaruiPada:      now,
	}

	// Simpan ke MySQL jika tersedia
	if g.DB != nil {
		tx, err := g.DB.Begin()
		if err == nil {
			queryOrder := `
			INSERT INTO pesanan (
				nomor_pesanan, nomor_invoice, tipe_checkout, id_pelanggan,
				nama_penerima, nomor_whatsapp, alamat_lengkap, kota_kabupaten, provinsi,
				kurir_id, kurir_nama, kurir_layanan, snapshot_subtotal, snapshot_ongkir,
				snapshot_total, status_pesanan, status_reservasi, batas_waktu_reservasi, dibuat_pada
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

			nama := kontak.NamaPenerima
			if strings.TrimSpace(nama) == "" {
				nama = "Pelanggan Contoh"
			}
			wa := kontak.NomorWhatsApp
			if strings.TrimSpace(wa) == "" || strings.Contains(strings.ToLower(wa), "x") {
				wa = "081234567890"
			}

			res, err := tx.Exec(queryOrder,
				nomorPesanan, nomorInvoice, tipeCheckout, idPelanggan,
				nama, wa, kontak.AlamatLengkap, kontak.KotaKabupaten, kontak.Provinsi,
				kurir.IDKurir, kurir.NamaKurir, kurir.KodeLayanan,
				subtotalSnapshot, kurir.TarifOngkir, totalSnapshot,
				"menunggu_pembayaran", "direservasi", batasWaktu, now,
			)

			if err == nil {
				idBaru, _ := res.LastInsertId()
				pesanan.ID = idBaru

				queryItem := `
				INSERT INTO snapshot_item_pesanan (
					id_pesanan, nomor_pesanan, id_varian, id_produk,
					nama_produk, nama_varian, jumlah, harga_satuan_snapshot,
					subtotal_snapshot, berat_satuan_gram, dibuat_pada
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

				for _, item := range items {
					_, _ = tx.Exec(queryItem,
						idBaru, nomorPesanan, item.IDVarian, item.IDProduk,
						item.NamaProduk, item.NamaVarian, item.Jumlah,
						item.HargaSatuanSnapshot, item.SubtotalSnapshot,
						item.BeratSatuanGram, now,
					)
				}
				_ = tx.Commit()
			} else {
				_ = tx.Rollback()
			}
		}
	}

	muPesanan.Lock()
	if pesanan.ID == 0 {
		pesanan.ID = int64(len(memoriPesananPPL18) + 1)
	}
	memoriPesananPPL18 = append(memoriPesananPPL18, *pesanan)
	muPesanan.Unlock()

	return pesanan, nil
}

// AmbilPesananBerdasarkanNomor mencari pesanan berdasarkan Nomor Pesanan resmi
func (g *GeneratorOrder) AmbilPesananBerdasarkanNomor(nomorPesanan string) (*EntitasPesanan, bool) {
	if g.DB != nil {
		query := `
		SELECT id, nomor_pesanan, nomor_invoice, tipe_checkout, nama_penerima, nomor_whatsapp,
		       alamat_lengkap, kota_kabupaten, provinsi, kurir_id, kurir_nama, kurir_layanan,
		       snapshot_subtotal, snapshot_ongkir, snapshot_total, status_pesanan, status_reservasi,
		       batas_waktu_reservasi, dibuat_pada
		FROM pesanan
		WHERE nomor_pesanan = ?`

		var p EntitasPesanan
		err := g.DB.QueryRow(query, nomorPesanan).Scan(
			&p.ID, &p.NomorPesanan, &p.NomorInvoice, &p.TipeCheckout,
			&p.Kontak.NamaPenerima, &p.Kontak.NomorWhatsApp, &p.Kontak.AlamatLengkap,
			&p.Kontak.KotaKabupaten, &p.Kontak.Provinsi,
			&p.Kurir.IDKurir, &p.Kurir.NamaKurir, &p.Kurir.KodeLayanan,
			&p.SnapshotSubtotal, &p.SnapshotOngkir, &p.SnapshotTotal,
			&p.StatusPesanan, &p.StatusReservasi, &p.BatasWaktuReservasi, &p.DibuatPada,
		)
		if err == nil {
			p.BatasWaktuFormat = p.BatasWaktuReservasi.Format("02 Jan 2006, 15:04 WIB")
			return &p, true
		}
	}

	muPesanan.RLock()
	defer muPesanan.RUnlock()
	for _, p := range memoriPesananPPL18 {
		if p.NomorPesanan == nomorPesanan {
			return &p, true
		}
	}

	return nil, false
}
