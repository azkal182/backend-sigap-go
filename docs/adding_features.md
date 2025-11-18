# Panduan Menambahkan Fitur Baru

Dokumen ini memindahkan dan merapikan panduan yang sebelumnya ada di bagian bawah `README.md`.

Filosofi dan pattern mengikuti **Clean Architecture / Hexagonal Architecture** dengan 4 lapisan utama:

1. Domain Layer
2. Application Layer
3. Infrastructure Layer
4. Interface/Delivery Layer

Contoh di bawah menggunakan **Product Management** sebagaimana di README asli.

---

## Step 1: Domain Entity

File: `internal/domain/entity/product.go`

```go
// Lihat README.md asli untuk contoh lengkap entity Product
```

Prinsip:

- Entity adalah struct Go murni (tidak tergantung framework).
- ID menggunakan UUID (kecuali kasus khusus seperti fitur lokasi).
- Tambahkan timestamps (`CreatedAt`, `UpdatedAt`) dan `DeletedAt` jika perlu soft delete.

## Step 2: Repository Interface (Port)

File: `internal/domain/repository/product_repository.go`

Prinsip:

- Didefinisikan di **domain layer**.
- Semua method menerima `context.Context`.
- Kembalikan entity domain + `error`.

## Step 3: Implementasi Repository (Adapter)

File: `internal/infrastructure/repository/product_repository.go`

Prinsip:

- Implementasi interface dengan GORM.
- Selalu gunakan `db.WithContext(ctx)`.
- Pisahkan dengan jelas antara entity domain dan cara data disimpan.

## Step 4: DTOs

File: `internal/application/dto/product_dto.go`

Prinsip:

- DTO request/response terpisah dari entity domain.
- Gunakan tag `binding:"required"` dan validasi lain sesuai kebutuhan.
- Untuk field opsional gunakan pointer atau `omitempty`.

## Step 5: Use Case

File: `internal/application/usecase/product_usecase.go`

Prinsip:

- Use case adalah orchestrator business logic.
- Validasi input (misalnya `page`, `pageSize`, dll.).
- Mapping entity → DTO dilakukan di layer usecase.
- Gunakan `domainErrors` untuk error yang terstandarisasi.

## Step 6: HTTP Handler

File: `internal/interfaces/http/handler/product_handler.go`

Prinsip:

- Handler hanya menangani concern HTTP:
  - Parsing request
  - Memanggil use case
  - Mengirim response dengan helper standar di `response` package
- Gunakan `response.SuccessOK`, `response.ErrorBadRequest`, dll.

## Step 7: Routes

File: `internal/interfaces/http/router/router.go`

Prinsip:

- Tambahkan handler baru ke signature `SetupRouter` bila perlu.
- Tambahkan routes di dalam group yang sesuai (`/api/products`, dll.).
- Gunakan middleware auth/permission jika fitur protected.

---

## Best Practices Umum

- **Domain Layer** tidak boleh bergantung ke framework atau DB.
- **Application Layer** hanya bergantung ke **Domain Layer**.
- **Infrastructure Layer** mengimplementasikan interface dari Domain.
- **Interface Layer** (HTTP) memanggil use case dan membentuk response JSON.

Untuk detail contoh kode lengkap (entity, repository, usecase, handler, routes) silakan merujuk ke versi asli di README atau fitur yang sudah ada seperti **User**, **Role**, **Dormitory**, dan **Location**.

---

## Step 8: Migrations

Untuk perubahan schema database (tabel baru / kolom baru), gunakan sistem migration versi yang sudah ada.

- File: `internal/infrastructure/database/migrations.go`
- Pattern:

```go
// Migration 004: Example feature
RegisterMigration(
    "004_create_example",
    "Create example table",
    func(db *gorm.DB) error {
        return db.AutoMigrate(&entity.Example{})
    },
    func(db *gorm.DB) error {
        return db.Migrator().DropTable(&entity.Example{})
    },
)
```

Tips:

- Gunakan `AutoMigrate` untuk kasus sederhana.
- Untuk perubahan spesifik (rename/drop kolom), gunakan `Migrator()` seperti yang dilakukan pada migration `003_remove_dormitory_address_and_capacity`.
- Migration dieksekusi saat aplikasi start via `database.MigrateUpVersioned()` atau manual via `cmd/migrate`.

---

## Step 9: Seeder / Import CLI

Ada dua pola utama:

1. **Seed data aplikasi** (permissions, roles, user admin, dll.)
   - File: `cmd/seed/main.go`
   - Menggunakan repository layer (`infraRepo.New...Repository`).

2. **Import data referensi besar** (contoh: lokasi Indonesia)
   - Contoh: `cmd/location_import/main.go`
   - Menggunakan `database.DB` langsung dan membaca file JSON dari folder `data/...`.

Rekomendasi:

- Untuk data yang **berkaitan dengan domain dan permission** → gunakan seeder (`cmd/seed`).
- Untuk data referensi besar dan statis → buat CLI import khusus di folder `cmd/`.
- Pastikan import/idempotent: cek dulu apakah data sudah ada (mis `WHERE id = ?`) sebelum insert.

---

## Step 10: Testing

### 10.1 Unit Test Use Case

- Lokasi: `internal/application/usecase/..._usecase_test.go`
- Gunakan `testify/mock` dan mocks dari `internal/application/usecase/mocks`.

Pola umum:

- Definisikan table-driven tests (`tests := []struct { ... }{ ... }`).
- Untuk tiap case, siapkan `setupMocks` yang mengatur expectation pada repository dan service.
- Di dalam test:
  - Buat instance usecase dengan mock repo.
  - Panggil method usecase.
  - Assert error / response sesuai ekspektasi.
  - Panggil `AssertExpectations(t)` pada mock.

### 10.2 HTTP Integration Test

- Lokasi: `internal/interfaces/http/integration_test.go`.
- Pola:
  - Gunakan `testutil.SetupTestDB` untuk membuat sementara DB khusus test.
  - Override `database.DB` sementara dengan DB test.
  - Inisialisasi repository, usecase, handler, middleware seperti di `cmd/main.go`.
  - Panggil `router.SetupRouter(...)` untuk mendapatkan `*gin.Engine`.
  - Gunakan `httptest.NewRecorder()` dan `http.NewRequest` untuk memukul endpoint.
  - Parse body JSON dan assert field `success`, `data`, `message`, dsb.

Gunakan integration test ini untuk memastikan flow end-to-end (DB + usecase + handler + middleware) bekerja.

---

## Step 11: Permissions & Authorization

Untuk fitur yang perlu proteksi (bukan public seperti lokasi), ikuti pola berikut:

1. **Tambah Permission**
   - Tambahkan entry baru di seeder `cmd/seed/main.go` seperti:

   ```go
   {ID: uuid.New(), Name: "product:read", Slug: "product-read", Resource: "product", Action: "read", ...}
   ```

   - Hubungkan ke role yang relevan (misalnya admin) menggunakan `AssignPermission`.

2. **Gunakan di Router**
   - Di `internal/interfaces/http/router/router.go`, bungkus route dengan middleware:

   ```go
   products := protected.Group("/products")
   {
       products.GET("",   authMiddleware.RequirePermission("product:read"),   productHandler.ListProducts)
       products.POST("",  authMiddleware.RequirePermission("product:create"), productHandler.CreateProduct)
       products.PUT("/:id", authMiddleware.RequirePermission("product:update"), productHandler.UpdateProduct)
       products.DELETE("/:id", authMiddleware.RequirePermission("product:delete"), productHandler.DeleteProduct)
   }
   ```

3. **JWT & Guard (opsional)**
   - Untuk fitur yang perlu akses berbasis resource tertentu (seperti dormitory), gunakan middleware tambahan seperti `RequireDormitoryAccess`.

Dengan pola ini, fitur baru akan konsisten dengan sistem authorization yang sudah ada di project.


3. **JWT & Guard (opsional)**
   - Untuk fitur yang perlu akses berbasis resource tertentu (seperti dormitory), gunakan middleware tambahan seperti `RequireDormitoryAccess`.

Dengan pola ini, fitur baru akan konsisten dengan sistem authorization yang sudah ada di project.

---

## 12. (Optional) Audit Logging

Untuk fitur yang mengubah data penting (misalnya CRUD **User**, **Role**, **Dormitory**), kamu bisa menambahkan audit log agar setiap aksi terekam dengan jelas.

### 12.1 Kapan perlu audit log?

- Operasi `create`, `update`, `delete` pada entity penting.
- Perubahan permission pada role.
- Aksi administratif lain yang sensitif (misalnya menghapus data penting).

### 12.2 Komponen audit log yang sudah tersedia

Project ini sudah menyediakan:

- Entity `AuditLog` dan tabel `audit_logs` (migration `006_create_audit_logs`).
- Repository `AuditLogRepository`.
- Service `AuditLogger` di `internal/application/service/audit_logger.go`.
- Middleware `AuditContextMiddleware` untuk mengisi informasi HTTP ke `context`.
- Endpoint baca audit log: `GET /api/audit-logs` (requires permission `audit:read`).

### 12.3 Menggunakan AuditLogger di use case

1. **Inject AuditLogger ke use case**

   Di `cmd/main.go`:

   ```go
   auditLogRepo := infraRepo.NewAuditLogRepository()
   auditLogger := service.NewAuditLogger(auditLogRepo)

   userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, auditLogger)
   roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, auditLogger)
   dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo, auditLogger)
   ```

2. **Panggil AuditLogger di akhir operasi sukses**

   Contoh di dalam `UserUseCase.CreateUser`:

   ```go
   _ = uc.auditLogger.Log(ctx, "user", "user:create", user.ID.String(), map[string]string{
       "email": user.Email,
       "name":  user.Name,
   })
   ```

   Parameter:

   - `resource`: nama resource, misalnya `"user"`, `"role"`, `"dormitory"`.
   - `action`: aksi spesifik, misalnya `"user:create"`, `"role:delete"`, `"dorm:update"`.
   - `targetID`: ID entity yang diubah (biasanya `uuid.String()`).
   - `metadata`: map string kecil yang akan disimpan sebagai JSON (hanya field penting, bukan seluruh payload).

   AuditLogger akan otomatis melengkapi informasi lain dari `context` (path, method, IP, user-agent, status code, actor).

### 12.4 Middleware AuditContext

Pastikan `AuditContextMiddleware` dipasang di router:

```go
router := gin.Default()
router.Use(middleware.NewCORSMiddlewareFromEnv())
router.Use(middleware.AuditContextMiddleware())
```

Middleware ini mengisi nilai berikut ke dalam `context.Context`:

- `request_path` (FullPath)
- `request_method`
- `ip_address`
- `user_agent`
- `status_code` (diisi setelah handler selesai)

Sehingga setiap pemanggilan `AuditLogger.Log` akan menyimpan log dengan informasi HTTP yang lengkap.

### 12.5 Permissions untuk baca audit log

Permission `audit:read` sudah ditambahkan di `cmd/seed/main.go` dan diberikan ke role `admin` dan `super_admin`.

Route HTTP:

```text
GET /api/audit-logs?page=1&page_size=10&resource=user&action=user:create&actor_email=admin@example.com
```

Route ini menggunakan middleware:

```go
auditLogs := protected.Group("/audit-logs")
{
    auditLogs.GET("", authMiddleware.RequirePermission("audit:read"), auditLogHandler.ListAuditLogs)
}
```

Gunakan endpoint ini untuk debug dan monitoring perubahan data penting di sistem.
