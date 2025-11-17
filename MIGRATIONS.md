# Database Migrations

Sistem migration ini menggunakan versioned migrations untuk mengelola perubahan schema database secara terkontrol.

## Struktur Migration

- `internal/infrastructure/database/migration.go` - Core migration system
- `internal/infrastructure/database/migrations.go` - Registered migrations
- `cmd/migrate/main.go` - CLI tool untuk menjalankan migrations

## Cara Menggunakan

### 1. Menjalankan Migrations

**Apply semua pending migrations:**
```bash
make migrate-up
# atau
go run cmd/migrate/main.go -command up
```

**Rollback migration terakhir:**
```bash
make migrate-down
# atau
go run cmd/migrate/main.go -command down
```

**Cek status migrations:**
```bash
make migrate-status
# atau
go run cmd/migrate/main.go -command status
```

**Migrate ke versi tertentu:**
```bash
make migrate-to VERSION=001_initial_schema
# atau
go run cmd/migrate/main.go -command to -version 001_initial_schema
```

### 2. Menambahkan Migration Baru

Untuk menambahkan migration baru, edit file `internal/infrastructure/database/migrations.go`:

```go
RegisterMigration(
    "002_add_user_phone",
    "Add phone field to users table",
    func(db *gorm.DB) error {
        // Migration UP - apply changes
        return db.Migrator().AddColumn(&entity.User{}, "phone")
    },
    func(db *gorm.DB) error {
        // Migration DOWN - rollback changes
        return db.Migrator().DropColumn(&entity.User{}, "phone")
    },
)
```

**Aturan Penamaan:**
- Format: `XXX_description` (contoh: `001_initial_schema`, `002_add_user_phone`)
- Versi harus unik dan sequential
- Nama harus deskriptif

### 3. Migration Table

Sistem migration secara otomatis membuat tabel `schema_migrations` untuk tracking:
- `id` - Primary key
- `version` - Migration version (unique)
- `name` - Migration name
- `applied_at` - Timestamp ketika migration di-apply

## Migration yang Tersedia

### 001_initial_schema
**Deskripsi:** Create initial database schema

**Tables:**
- `users` - User accounts
- `roles` - User roles
- `permissions` - Permissions
- `dormitories` - Dormitories
- `user_roles` - Many-to-many: users ↔ roles
- `role_permissions` - Many-to-many: roles ↔ permissions
- `user_dormitories` - Many-to-many: users ↔ dormitories

**Rollback:** Drops all tables

## Best Practices

1. **Selalu test migration di development** sebelum apply ke production
2. **Selalu buat rollback function** untuk setiap migration
3. **Jangan edit migration yang sudah di-apply** di production
4. **Gunakan transaction** untuk migration yang kompleks (opsional)
5. **Backup database** sebelum menjalankan migration di production

## Troubleshooting

### Migration gagal di tengah jalan
Jika migration gagal, sistem akan stop di migration tersebut. Perbaiki masalahnya dan jalankan `migrate-up` lagi.

### Rollback migration
Gunakan `migrate-down` untuk rollback migration terakhir yang sudah di-apply.

### Reset semua migrations
```sql
-- Hapus semua data di schema_migrations
TRUNCATE TABLE schema_migrations;

-- Hapus semua tables (HATI-HATI!)
DROP TABLE IF EXISTS user_dormitories CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS dormitories CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS schema_migrations CASCADE;
```

Kemudian jalankan `migrate-up` untuk apply semua migrations dari awal.

## Integration dengan CI/CD

Contoh untuk GitHub Actions:

```yaml
- name: Run migrations
  run: |
    go run cmd/migrate/main.go -command up
  env:
    DB_HOST: ${{ secrets.DB_HOST }}
    DB_USER: ${{ secrets.DB_USER }}
    DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
    DB_NAME: ${{ secrets.DB_NAME }}
    DB_PORT: ${{ secrets.DB_PORT }}
    DB_SSLMODE: ${{ secrets.DB_SSLMODE }}
```
