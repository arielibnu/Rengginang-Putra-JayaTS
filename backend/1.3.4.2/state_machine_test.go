package statemachine

import (
	"testing"

	pengembalianstok "rengginang-putra-jaya/auth/1.3.4.4"
)

func TestValidasiTransisiStatus(t *testing.T) {
	layananStok := pengembalianstok.InisialisasiLayanan(nil)
	sm := InisialisasiLayanan(nil, layananStok)

	// Kasus Valid
	kasusValid := []struct {
		lama string
		baru string
	}{
		{"menunggu_diproses", "diproses"},
		{"baru", "diproses"},
		{"diproses", "dikirim"},
		{"dikirim", "selesai"},
		{"diproses", "batal"},
		{"menunggu_diproses", "batal"},
	}

	for _, k := range kasusValid {
		valid, alasan := sm.ValidasiTransisiStatus(k.lama, k.baru)
		if !valid {
			t.Errorf("Transisi valid %s -> %s ditolak: %s", k.lama, k.baru, alasan)
		}
	}

	// Kasus Tidak Valid (Harus Ditolak)
	kasusTidakValid := []struct {
		lama string
		baru string
	}{
		{"selesai", "batal"},   // status selesai tidak boleh batal
		{"selesai", "diproses"}, // status selesai tidak boleh mundur
		{"batal", "diproses"},   // status batal tidak boleh aktif lagi
		{"dikirim", "batal"},    // barang sudah dikirim kurir tidak boleh langsung batal tanpa retur
		{"baru", "selesai"},     // lompat status dilarang
	}

	for _, k := range kasusTidakValid {
		valid, _ := sm.ValidasiTransisiStatus(k.lama, k.baru)
		if valid {
			t.Errorf("Transisi tidak valid %s -> %s seharusnya ditolak!", k.lama, k.baru)
		}
	}
}

func TestPemisahanStatusPembayaran(t *testing.T) {
	layananStok := pengembalianstok.InisialisasiLayanan(nil)
	sm := InisialisasiLayanan(nil, layananStok)

	req := PermintaanUbahStatus{
		NomorPesanan:     "PJ-DEMO-0002",
		StatusPemrosesan: "dikirim",
		NomorResi:        "JNT-123456",
	}

	resp, err := sm.UbahStatusPemrosesan(req)
	if err != nil {
		t.Fatalf("Gagal ubah status: %v", err)
	}

	if resp.StatusPemrosesanBaru != "dikirim" {
		t.Errorf("Status pemrosesan baru harus 'dikirim', didapat: %s", resp.StatusPemrosesanBaru)
	}

	// Status pembayaran tetap utuh (terpisah)
	if resp.StatusPembayaranTetap != "dibayar" {
		t.Errorf("Status pembayaran harus tetap 'dibayar', didapat: %s", resp.StatusPembayaranTetap)
	}
}
