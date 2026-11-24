package flagpemisah

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluasiFlagTransaksiNyata(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	flag, isSandbox := layanan.EvaluasiFlag("offline_kasir", "")
	if flag != FlagProduksiNyata {
		t.Errorf("Ekspektasi %s, didapat: %s", FlagProduksiNyata, flag)
	}
	if isSandbox != false {
		t.Error("Ekspektasi isSandbox false untuk offline_kasir")
	}

	flagTunai, isSandboxTunai := layanan.EvaluasiFlag("tunai", "")
	if flagTunai != FlagProduksiNyata || isSandboxTunai != false {
		t.Error("Ekspektasi tunai diklasifikasikan sebagai transaksi nyata")
	}
}

func TestEvaluasiFlagTransaksiSandbox(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	flag, isSandbox := layanan.EvaluasiFlag("sandbox_midtrans", "SANDBOX_SIMULASI")
	if flag != FlagSandboxSimulasi {
		t.Errorf("Ekspektasi %s, didapat: %s", FlagSandboxSimulasi, flag)
	}
	if isSandbox != true {
		t.Error("Ekspektasi isSandbox true untuk sandbox_midtrans")
	}
}

func TestHandlerValidasiFlagHTTP(t *testing.T) {
	layanan := InisialisasiLayanan(nil)

	payload := PermintaanValidasiFlag{
		SumberTransaksi: "offline_kasir",
		FlagLingkungan:  "PRODUKSI_NYATA",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/transaksi/validasi-flag", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	layanan.HandlerValidasiFlag(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ekspektasi status 200, didapat: %d", rr.Code)
	}

	var resp ResponsValidasiFlag
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Gagal parse respons JSON: %v", err)
	}

	if !resp.Valid || resp.IsSandbox {
		t.Errorf("Respons tidak valid: %+v", resp)
	}
}
