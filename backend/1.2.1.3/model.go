package kuotastok

import (
	"encoding/json"
	"strconv"
)

// ============================================================================
// WBS 1.2.1.3 - VALIDASI BATAS JUMLAH PEMBELIAN TERHADAP SISA STOK REALTIME
// Model data untuk validasi kuota stok realtime pada katalog dan detail produk.
// ============================================================================

// PermintaanValidasiKuota menampung input validasi kuota stok realtime.
// Mendukung fleksibilitas tipe input (string ID "udang", numeric ID 1, slug, dan alias jumlah/jumlah_beli).
type PermintaanValidasiKuota struct {
	IDVarian   uint64 `json:"id_varian"`
	Slug       string `json:"slug"`
	JumlahBeli int64  `json:"jumlah_beli"`
	Jumlah     int64  `json:"jumlah"` // Alias dari frontend
}

// UnmarshalJSON kustom untuk mendukung input fleksibel dari frontend React.
func (p *PermintaanValidasiKuota) UnmarshalJSON(data []byte) error {
	type Alias PermintaanValidasiKuota
	aux := struct {
		RawIDVarian interface{} `json:"id_varian"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Tangani format RawIDVarian (string "udang" atau numeric 1)
	switch v := aux.RawIDVarian.(type) {
	case string:
		if num, err := strconv.ParseUint(v, 10, 64); err == nil {
			p.IDVarian = num
		} else {
			p.Slug = v
		}
	case float64:
		p.IDVarian = uint64(v)
	}

	// Sinkronisasi alias jumlah dan jumlah_beli
	if p.JumlahBeli == 0 && p.Jumlah != 0 {
		p.JumlahBeli = p.Jumlah
	} else if p.Jumlah == 0 && p.JumlahBeli != 0 {
		p.Jumlah = p.JumlahBeli
	}

	return nil
}

// DataVarianStok entitas internal untuk informasi varian dan kuota stok dari database.
type DataVarianStok struct {
	IDVarian    uint64  `json:"id_varian"`
	IDProduk    uint64  `json:"id_produk"`
	Slug        string  `json:"slug"`
	NamaVarian  string  `json:"nama_varian"`
	NamaProduk  string  `json:"nama_produk"`
	HargaJual   uint64  `json:"harga_jual"`
	BeratGram   uint32  `json:"berat_gram"`
	Stok        int64   `json:"stok"`
	StatusAktif bool    `json:"status_aktif"`
}

// HasilValidasiKuota hasil evaluasi ketersediaan kuota stok secara realtime.
type HasilValidasiKuota struct {
	IDVarian           uint64 `json:"id_varian"`
	Slug               string `json:"slug"`
	NamaVarian         string `json:"nama_varian"`
	HargaSatuan        uint64 `json:"harga_satuan"`
	FormatHarga        string `json:"format_harga"`
	StokTersedia       int64  `json:"stok_tersedia"`
	FormatStok         string `json:"format_stok"`
	BatasMaksimalBeli  int64  `json:"batas_maksimal_beli"`
	JumlahDiminta      int64  `json:"jumlah_diminta"`
	TotalHarga         uint64 `json:"total_harga"`
	FormatTotalHarga   string `json:"format_total_harga"`
	SisaSetelahBeli    int64  `json:"sisa_setelah_beli"`
	StatusKetersediaan string `json:"status_ketersediaan"` // "TERSEDIA", "STOK HABIS", "MELEBIHI STOK", "JUMLAH TIDAK VALID"
	KuotaTersedia      bool   `json:"kuota_tersedia"`
	BisaDibeli         bool   `json:"bisa_dibeli"`
	Valid              bool   `json:"valid"`
	Pesan              string `json:"pesan"`
	Keterangan         string `json:"keterangan"`
}

// ResponsAPI format respon JSON backend seragam.
type ResponsAPI struct {
	Status string              `json:"status"` // "sukses", "ditolak", "gagal"
	Pesan  string              `json:"pesan"`
	Data   *HasilValidasiKuota `json:"data,omitempty"`
}
