import urllib.request
import json
import time

BASE_URL = 'http://localhost:8080'

print('======================================================================')
print('MEMULAI PENGUJIAN BACKEND PPL-23: STATE MACHINE & PENGEMBALIAN STOK')
print('======================================================================')

time.sleep(1)

def kirim_request(metode, endpoint, data=None):
    url = f"{BASE_URL}{endpoint}"
    req = urllib.request.Request(url, method=metode)
    req.add_header('Content-Type', 'application/json')
    body = json.dumps(data).encode('utf-8') if data else None
    try:
        with urllib.request.urlopen(req, data=body) as resp:
            return resp.status, json.loads(resp.read().decode('utf-8'))
    except urllib.error.HTTPError as e:
        try:
            res_body = json.loads(e.read().decode('utf-8'))
        except Exception:
            res_body = {'error': e.reason}
        return e.code, res_body

# 1. Test Endpoint Pilihan Opsi Transisi
code, res = kirim_request('GET', '/api/admin/pesanan/transisi-status?status=diproses')
print(f'1. Opsi Transisi Status dari diproses -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 200, f'Expected 200, got {code}'

# 2. Test Transisi Valid: diproses -> dikirim (PPL-23)
order_test = f"PJ-SM-{int(time.time())}"
payload_kirim = {
    'nomor_pesanan': order_test,
    'status_pemrosesan': 'dikirim',
    'nomor_resi': 'JNT-9988771122',
    'catatan': 'Pesanan diserahkan ke kurir J&T'
}
code, res = kirim_request('PUT', '/api/admin/pesanan/status', payload_kirim)
print(f'2. Transisi Valid diproses -> dikirim ({order_test}) -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 200, f'Expected 200, got {code}'
assert res.get('status_pemrosesan_baru') == 'dikirim'
assert res.get('status_pembayaran_tetap') == 'dibayar'

# 3. Test Transisi Valid: dikirim -> selesai
payload_selesai = {
    'nomor_pesanan': order_test,
    'status_pemrosesan': 'selesai',
    'catatan': 'Paket diterima oleh pelanggan'
}
code, res = kirim_request('PUT', '/api/admin/pesanan/status', payload_selesai)
print(f'3. Transisi Valid dikirim -> selesai ({order_test}) -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 200, f'Expected 200, got {code}'
assert res.get('status_pemrosesan_baru') == 'selesai'

# 4. Test Transisi Ilegal (Harus Ditolak State Machine 400 Bad Request): selesai -> batal
payload_batal_ilegal = {
    'nomor_pesanan': order_test,
    'status_pemrosesan': 'batal',
    'catatan': 'Mencoba membatalkan pesanan yang sudah selesai'
}
code, res = kirim_request('PUT', '/api/admin/pesanan/status', payload_batal_ilegal)
print(f'4. Transisi Ditolak State Machine (selesai -> batal) ({order_test}) -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 400, f'Expected 400, got {code}'

# 5. Test Pembatalan & Pengembalian Stok Otomatis (WBS 1.3.4.4)
order_test_cancel = f"PJ-CANCEL-{int(time.time())}"
payload_batal_stok = {
    'nomor_pesanan': order_test_cancel,
    'status_pemrosesan': 'batal',
    'catatan': 'Pelanggan minta batalkan order sebelum dikirim'
}
code, res = kirim_request('PUT', '/api/admin/pesanan/status', payload_batal_stok)
print(f'5. Pembatalan & Pengembalian Stok Otomatis ({order_test_cancel}) -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 200, f'Expected 200, got {code}'
assert res.get('stok_dikembalikan') == True
assert res.get('total_bungkus_kembali', 0) > 0

# 6. Test Idempotency Pengembalian Stok
code, res = kirim_request('POST', '/api/admin/pesanan/kembalikan-stok', {'nomor_pesanan': order_test_cancel, 'alasan': 'Pengujian idempotency'})
print(f'6. Idempotency Pengembalian Stok ({order_test_cancel}) -> HTTP {code}:', json.dumps(res, indent=2))
assert code == 200, f'Expected 200, got {code}'
assert res['data']['nomor_pesanan'] == order_test_cancel

print('======================================================================')
print('SEMUA PENGUJIAN BACKEND PPL-23 (1.3.4.2 & 1.3.4.4) BERHASIL 100%!')
print('======================================================================')
