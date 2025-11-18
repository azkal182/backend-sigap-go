# Fitur Lokasi: Province, Regency, District, Village

Fitur ini menyediakan data lokasi Indonesia secara **read-only**:

- Province (Provinsi)
- Regency (Kabupaten/Kota)
- District (Kecamatan)
- Village (Desa/Kelurahan)

Data diambil dari tabel statis di database dan dapat diimport dari file JSON.

## Endpoint HTTP (Public)

Semua endpoint **tidak membutuhkan autentikasi** dan mendukung:

- Pagination: `page`, `page_size` (default `1`, `10`, max `100`)
- Pencarian: `search` (case-insensitive, berdasarkan `name`)
- Filter berdasarkan parent ID (untuk regency, district, village)

### 1. Provinces

- `GET /api/provinces`
  - Query params:
    - `page` (opsional)
    - `page_size` (opsional)
    - `search` (opsional, cari berdasarkan `name`)

- `GET /api/provinces/:id`
  - Path param: `id` (int)

**Response contoh (list):**

```json
{
  "success": true,
  "message": "Provinces retrieved successfully",
  "data": {
    "items": [
      { "id": 11, "name": "Aceh (NAD)", "code": "11" }
    ],
    "total": 34,
    "page": 1,
    "page_size": 10,
    "total_pages": 4
  }
}
```

### 2. Regencies

- `GET /api/regencies`
  - Query params:
    - `province_id` (opsional, filter berdasarkan provinsi)
    - `page`, `page_size` (opsional)
    - `search` (opsional)

- `GET /api/regencies/:id`
  - Path param: `id` (int)

### 3. Districts

- `GET /api/districts`
  - Query params:
    - `regency_id` (opsional)
    - `page`, `page_size` (opsional)
    - `search` (opsional)

- `GET /api/districts/:id`
  - Path param: `id` (int)

### 4. Villages

- `GET /api/villages`
  - Query params:
    - `district_id` (opsional)
    - `page`, `page_size` (opsional)
    - `search` (opsional)

- `GET /api/villages/:id`
  - Path param: `id` (int)

## Struktur Data

### Province

```json
{
  "id": 1,
  "name": "Aceh (NAD)",
  "code": "11"
}
```

### Regency

```json
{
  "id": 1,
  "type": "Kabupaten",
  "name": "Aceh Barat",
  "code": "05",
  "full_code": "1105",
  "provinsi_id": 1
}
```

### District

```json
{
  "id": 1,
  "name": "Air Majunto",
  "code": "13",
  "full_code": "170613",
  "kabupaten_id": 420
}
```

### Village

```json
{
  "id": 10,
  "name": "Yawosi (Fanindi)",
  "code": "2006",
  "full_code": "9106132006",
  "pos_code": "98552",
  "kecamatan_id": 7164
}
```

> Catatan: Di level database, foreign key dinormalisasi menjadi `province_id`, `regency_id`, dan `district_id`, tetapi JSON mengikuti format field yang kamu miliki.

## Import Data Lokasi dari JSON

Import dilakukan melalui command khusus, **bukan** bagian dari migration otomatis.

### Lokasi File JSON

Secara default, importer membaca file dari:

- `data/locations/provinces.json`
- `data/locations/regencies.json`
- `data/locations/districts.json`
- `data/locations/villages.json`

Masing-masing berupa **array JSON**:

```json
[
  { "id": 1, "name": "Aceh (NAD)", "code": "11" },
  { "id": 2, "name": "Sumatera Utara", "code": "12" }
]
```

### Menjalankan Import

Pastikan database sudah dibuat dan migration sudah dijalankan.

```bash
# Jalankan import lokasi
go run cmd/location_import/main.go
```

Importer akan menjalankan langkah berikut secara berurutan:

1. Import provinces
2. Import regencies
3. Import districts
4. Import villages

### Perilaku Import (Idempotent)

- Importer akan **cek `id` dulu**.
- Jika data dengan `id` yang sama sudah ada di tabel → data tersebut **dilewati** (tidak di-insert ulang).
- Jika terjadi error selain "record not found" saat pengecekan, error akan di-log dan importer lanjut ke record berikutnya.

Dengan demikian, command ini aman dijalankan berkali-kali tanpa menggandakan data.

## Catatan Tambahan

- Endpoint lokasi bersifat **read-only**.
- ID menggunakan tipe **int** khusus untuk fitur ini (karena mengikuti dataset yang ada), berbeda dengan fitur lain yang memakai UUID.
- Untuk pencarian, field yang digunakan saat ini adalah `name`.
