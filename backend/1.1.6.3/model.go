package modalapi

import "time"

// ==========================================
// STRUKTUR DATA MODAL PER VARIAN
// ==========================================

// Data modal produk yang disimpan di database.
type ModalVarian struct {
	IDModal      uint64    `json:"id_modal"`
	IDVarian     uint64    `json:"id_varian"`
	NamaVarian   string    `json:"nama_varian"`
	HargaModal   float64   `json:"harga_modal"`
	BerlakuMulai string    `json:"berlaku_mulai"`
	Catatan      string    `json:"catatan"`
	DibuatPada   time.Time `json:"dibuat_pada"`
}

// Permintaan untuk menyimpan atau mengubah modal varian.
type PermintaanModal struct {
	IDVarian     uint64  `json:"id_varian"`
	HargaModal   float64 `json:"harga_modal"`
	BerlakuMulai string  `json:"berlaku_mulai"`
	Catatan      string  `json:"catatan"`
}

// ==========================================
// STRUKTUR DATA PENGELUARAN OPERASIONAL
// ==========================================

// Data pengeluaran usaha yang disimpan di database.
type PengeluaranUsaha struct {
	IDPengeluaran       uint64    `json:"id_pengeluaran"`
	TanggalPengeluaran  string    `json:"tanggal_pengeluaran"`
	Kategori            string    `json:"kategori"`
	Nominal             float64   `json:"nominal"`
	Keterangan          string    `json:"keterangan"`
	DibuatPada          time.Time `json:"dibuat_pada"`
	DiperbaruiPada      time.Time `json:"diperbarui_pada"`
}

// Permintaan untuk mencatat pengeluaran baru.
type PermintaanPengeluaran struct {
	TanggalPengeluaran string  `json:"tanggal_pengeluaran"`
	Kategori           string  `json:"kategori"`
	Nominal            float64 `json:"nominal"`
	Keterangan         string  `json:"keterangan"`
}

// ==========================================
// STRUKTUR RINGKASAN LAPORAN BULANAN
// ==========================================

// Ringkasan modal dan pengeluaran untuk satu bulan.
type RingkasanBulanan struct {
	Bulan              string  `json:"bulan"`
	RataRataModal      float64 `json:"rata_rata_modal"`
	JumlahVarian       int     `json:"jumlah_varian"`
	TotalPengeluaran   float64 `json:"total_pengeluaran"`
	JumlahPosBiaya     int     `json:"jumlah_pos_biaya"`
}

// ==========================================
// RESPONS API
// ==========================================

// Format standar respons JSON untuk semua endpoint.
type ResponsAPI struct {
	Pesan string      `json:"pesan"`
	Data  interface{} `json:"data,omitempty"`
}
