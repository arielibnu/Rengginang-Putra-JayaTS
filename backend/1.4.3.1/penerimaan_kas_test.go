package penerimaankas

import (
	"testing"
)

func TestParsingPeriodePenerimaanKas(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	p, err := layanan.ParsingPeriode("November", "2026")
	if err != nil {
		t.Fatalf("Gagal parsing periode: %v", err)
	}

	if p.BulanNama != "November" || p.BulanAngka != 11 || p.Tahun != 2026 {
		t.Errorf("Periode tidak cocok: %+v", p)
	}
}

func TestHitungRekapPenerimaanKas(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	laporan, err := layanan.HitungRekapPenerimaanKas("November", "2026")
	if err != nil {
		t.Fatalf("Gagal menghitung rekap penerimaan kas: %v", err)
	}

	if laporan.KasMasukNyata <= 0 {
		t.Errorf("Kas masuk nyata harus > 0, didapat: %f", laporan.KasMasukNyata)
	}

	if laporan.SimulasiSandbox <= 0 {
		t.Errorf("Simulasi sandbox harus > 0, didapat: %f", laporan.SimulasiSandbox)
	}

	// Anti-Pencampuran: Kas tercatat nyata HANYA berasal dari transaksi riil lunas
	if laporan.TotalKasTercatat != laporan.KasMasukNyata {
		t.Errorf("Total kas tercatat tidak boleh mencampurkan sandbox! Nyata: %f, Total: %f", laporan.KasMasukNyata, laporan.TotalKasTercatat)
	}
}
