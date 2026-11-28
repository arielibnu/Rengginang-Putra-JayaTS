package auditkeuntungan

import (
	"testing"
)

func TestAuditKeuntunganKelengkapanDanAntiDuplikasi(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	hasil, err := layanan.AuditKelengkapanDanDuplikasi(11, 2026)
	if err != nil {
		t.Fatalf("AuditKelengkapanDanDuplikasi gagal: %v", err)
	}

	if hasil.Periode.BulanNama != "November" {
		t.Errorf("Ekspektasi bulan November, dapat %s", hasil.Periode.BulanNama)
	}

	// Verifikasi audit anti duplikasi
	if hasil.AuditAntiDuplikasi.JumlahTransaksiValid != hasil.AuditAntiDuplikasi.JumlahTransaksiTerdeduplikasi {
		t.Errorf("Jumlah transaksi unik tidak sesuai")
	}

	// Verifikasi status integritas
	if hasil.AuditAntiDuplikasi.StatusIntegritas == "" {
		t.Errorf("Status integritas audit harus terisi")
	}
}
