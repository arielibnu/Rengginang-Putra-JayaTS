package pengembalianstok

import (
	"testing"
)

func TestKembalikanStokPesananBatal(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	// Test 1: Pengembalian stok pertama kali (Sukses)
	hasil1, err1 := layanan.KembalikanStokPesananBatal("PJ-DEMO-0002", "Uji coba pembatalan")
	if err1 != nil {
		t.Fatalf("Gagal mengembalikan stok: %v", err1)
	}

	if hasil1.Status != "sukses" {
		t.Errorf("Status harus sukses, didapat: %s", hasil1.Status)
	}

	if hasil1.TotalBungkusKembali <= 0 {
		t.Errorf("Total bungkus kembali harus > 0, didapat: %d", hasil1.TotalBungkusKembali)
	}

	// Test 2: Idempotency (Pengembalian kedua kalinya untuk nomor pesanan yang sama harus dicegah)
	hasil2, err2 := layanan.KembalikanStokPesananBatal("PJ-DEMO-0002", "Uji coba pengulangan")
	if err2 != nil {
		t.Fatalf("Idempotency error: %v", err2)
	}

	if hasil2 != hasil1 {
		t.Errorf("Hasil idempotency harus mengembalikan instance hasil awal yang sama")
	}
}
