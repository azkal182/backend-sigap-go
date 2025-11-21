# Go Backend Starter - Hexagonal Architecture

Aplikasi backend starter berbasis **Golang** dengan **Clean Architecture / Hexagonal Architecture** yang mendukung skalabilitas, maintainability, dan pemisahan concern yang jelas antara domain, aplikasi, dan infrastruktur.

## 🚀 Fitur Utama

### 1. Authentication & Authorization
- ✅ JWT Authentication (Access Token & Refresh Token)
- ✅ User Registration & Login
- ✅ Refresh Token Endpoint
- ✅ Middleware JWT untuk proteksi endpoint

### 2. Role & Permission System
- ✅ Role-based dan Permission-based authorization
- ✅ **Role Management** - CRUD operations untuk roles
- ✅ **Protected Roles** - Role tertentu (admin, super_admin) tidak bisa diubah permission-nya
- ✅ **Assign/Remove Permissions** - Kelola permission pada role
- ✅ **Assign/Remove Roles to Users** - Kelola role assignment pada user
- ✅ **Default Role Assignment** - User baru otomatis mendapat role "user" jika tidak ditentukan
- ✅ Contoh permission: `user:read`, `user:update`, `dorm:read`, `dorm:update`, `role:read`, `role:create`, dll
- ✅ Role dapat memiliki banyak permission
- ✅ User dapat memiliki satu atau lebih role

### 3. User Management (CRUD Users)
- ✅ Create, Read, Update, Delete user
- ✅ Pagination support
- ✅ Role assignment saat create user
- ✅ Assign/Remove role ke user yang sudah ada
- ✅ Default role "user" untuk user baru

### 4. Dormitory Management (CRUD Dormitory)
- ✅ CRUD untuk data dormitory (dengan atribut `name`, `gender`, `level`, `code`, status aktif)
- ✅ Assign/Remove staff ke dormitory via endpoint khusus
- ✅ Setiap dormitory dapat dibatasi akses berdasarkan guard

### 5. Student Management (CRUD Students + SKS Results)
- ✅ CRUD santri dengan atribut `student_number`, `full_name`, `birth_date`, `gender`, `parent_name`, status aktif
- ✅ Patch status (`active`, `inactive`, `leave`, `graduated`) via endpoint khusus
- ✅ Mutasi asrama (riwayat `student_dormitory_history`) dengan audit logging
- ✅ Rekam hasil ujian SKS per santri dan pantau status kelulusan FAN secara otomatis
- ✅ Endpoint terproteksi permission `student:*` dan `student_sks_results:*`

### 6. Guard / Access Control
- ✅ Guard menentukan batas akses user terhadap dormitory:
  - **Access to specific dormitories only** — staff hanya dapat mengelola dormitory tertentu
  - **Access to all dormitories** — admin dapat mengelola seluruh dormitory

### 6. Standardized API Response
- ✅ Response format yang konsisten untuk semua endpoint
- ✅ Success response dengan struktur standar
- ✅ Error response dengan struktur standar
- ✅ Helper functions untuk berbagai HTTP status codes

### 7. Schedule Slot Management
- ✅ CRUD time-slot terstandarisasi per asrama/dormitory
- ✅ Validasi overlap slot per dormitory dan status aktif
- ✅ Reusable oleh fitur jadwal kelas sehingga jam dapat diwariskan otomatis

### 8. Class Schedule & SKS Scheduling
- ✅ CRUD jadwal kelas dengan opsi gunakan slot atau jam manual + validasi guru aktif
- ✅ CRUD definisi SKS per FAN dan map ke subject opsional
- ✅ CRUD jadwal ujian SKS dengan pemeriksa (teacher) aktif + validasi tanggal/waktu
- ✅ Permission khusus (`class_schedules:*`, `sks_definitions:*`, `sks_exams:*`) dan audit logging untuk seluruh mutasi

## 📁 Struktur Project (Hexagonal Architecture)

```
.
├── cmd/
│   ├── main.go              # Entry point aplikasi
│   └── seed/
│       └── main.go          # Seed data untuk development
├── internal/
│   ├── domain/              # Domain Layer (Core Business Logic)
│   │   ├── entity/          # Domain entities
│   │   ├── repository/      # Repository interfaces (ports)
│   │   ├── service/         # Domain service interfaces
│   │   └── errors/          # Domain errors
│   ├── application/         # Application Layer (Use Cases)
│   │   ├── usecase/         # Business use cases
│   │   └── dto/             # Data Transfer Objects
│   ├── infrastructure/      # Infrastructure Layer (Adapters)
│   │   ├── database/        # Database connection & migration
│   │   ├── repository/      # Repository implementations
│   │   └── service/         # Service implementations (JWT, etc)
│   └── interfaces/          # Interface/Delivery Layer
│       └── http/
│           ├── handler/     # HTTP handlers
│           ├── middleware/  # HTTP middleware
│           ├── response/    # Standardized response helpers
│           └── router/       # Route configuration
├── go.mod
├── go.sum
├── .env.example
├── Makefile
└── README.md
```

## 🏗️ Arsitektur Clean (Hexagonal Architecture)

### **1. Domain Layer**
Berisi **entity**, **value object**, **domain service**, dan **business rules**.
- `User`, `Role`, `Permission`, `Dormitory`
- Tidak bergantung pada database atau framework

### **2. Application Layer (Use Cases)**
Berisi **service/use case** seperti:
- `RegisterUser`, `LoginUser`, `RefreshToken`
- `CreateDormitory`, `UpdateDormitory`, dll.
- Menggunakan **interface repository** (port) yang diimplementasikan di infrastruktur

### **3. Infrastructure Layer (Adapters)**
Implementasi repository dan service:
- PostgreSQL repository (GORM)
- JWT token service
- Database connection

### **4. Interface/Delivery Layer**
Controller/handler HTTP:
- JWT Auth middleware
- Permission checker middleware
- Dormitory guard middleware
- Mapping request/response ke DTO
- Standardized response format untuk semua endpoint

## 🔐 Flow Authorization

1. Request masuk → Middleware cek JWT
2. Middleware cek **role & permission** sesuai endpoint
3. Jika endpoint terkait dormitory → Guard cek:
   - User memiliki akses ke dormitory id tertentu
   - atau user memiliki akses global (admin/super_admin)
4. Jika lolos → dilanjutkan ke handler

## 📋 Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Make (optional, untuk menggunakan Makefile)

## 🛠️ Installation

### 1. Clone Repository
```bash
git clone <repository-url>
cd go-backend-starter
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment Variables
```bash
cp .env.example .env
```

Edit `.env` file:
```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h

# Application
APP_ENV=development
LOG_LEVEL=debug

# CORS
# Comma-separated list of allowed origins, e.g.:
# CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com
CORS_ALLOWED_ORIGINS=
```

### 4. Setup Database
```bash
# Create PostgreSQL database
createdb go_backend_db
```

### 5. Run Migrations
Migrations menggunakan versioned migration system yang lebih aman dan dapat diulang:

```bash
# Run migrations (akan berjalan otomatis saat aplikasi start)
make run

# Atau jalankan migration manual
make migrate-up

# Check migration status
make migrate-status

# Rollback last migration
make migrate-down
```

**Catatan:** 
- Migrations akan otomatis berjalan saat aplikasi start
- Migration 001: Membuat schema database awal
- Migration 002: Menambahkan field `is_protected` pada roles table dan seed default roles
- Migration 003: Menghapus kolom `address` dan `capacity` dari tabel `dormitories` (schema dormitory sekarang hanya memuat `name`, `description`, `is_active`, timestamps, dan relasi)
- Default roles yang dibuat: `user`, `admin`, `super_admin`
- Role `admin` dan `super_admin` adalah protected roles

### 6. Seed Data (Optional)
```bash
go run cmd/seed/main.go
```

Ini akan membuat:
- **Permissions**: 
  - User: `user:read`, `user:create`, `user:update`, `user:delete`
  - Dormitory: `dorm:read`, `dorm:create`, `dorm:update`, `dorm:delete`
  - Role: `role:read`, `role:create`, `role:update`, `role:delete`
- **Roles**: 
  - `user` (default role, not protected) - memiliki `dorm:read`
  - `admin` (protected) - memiliki semua permissions
  - `super_admin` (protected) - memiliki semua permissions
- **Users**:
  - Admin: `admin` / `admin123`
  - Super Admin: `superadmin` / `superadmin123`
- **Sample dormitories**

### 7. Run Application
```bash
# Using Make
make run

# Or directly
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080`

### 8. OpenAPI Workflow (Kontributor)
- Jalankan `make openapi-sync` sebelum push untuk memastikan `docs/openapi.yaml` valid menurut Spectral dan sudah di-commit.
- CI (`.github/workflows/main.yml`) juga menjalankan lint yang sama, jadi pastikan lulus lokal agar pipeline tidak gagal.

## 📡 API Endpoints

### Authentication (Public)
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token

### Users (Protected)
- `GET /api/users` - List users (with pagination)
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create user (requires `user:create` permission)
- `PUT /api/users/:id` - Update user (requires `user:update` permission)
- `DELETE /api/users/:id` - Delete user (requires `user:delete` permission)
- `POST /api/users/:id/roles` - Assign role to user (requires `user:update` permission)
- `DELETE /api/users/:id/roles/:role_id` - Remove role from user (requires `user:update` permission)

### Current User (Protected)
- `GET /api/me` - Get current authenticated user (requires valid access token)

### Roles (Protected)
- `GET /api/roles` - List roles (with pagination, requires `role:read` permission)
- `GET /api/roles/:id` - Get role by ID (requires `role:read` permission)
- `POST /api/roles` - Create role (requires `role:create` permission)
- `PUT /api/roles/:id` - Update role (requires `role:update` permission)
- `DELETE /api/roles/:id` - Delete role (requires `role:delete` permission, protected roles cannot be deleted)
- `POST /api/roles/:id/permissions` - Assign permission to role (requires `role:update` permission, protected roles cannot be modified)
- `DELETE /api/roles/:id/permissions` - Remove permission from role (requires `role:update` permission, protected roles cannot be modified)

### Permissions (Protected)
- `GET /api/permissions` - List permissions (with pagination, requires `role:read` permission)

### Leave Permits & Health Statuses (Security / UKS)
- `GET /api/leave-permits` - List permits (filters: student, status, type, date; requires `leave_permits:read`)
- `POST /api/leave-permits` - Create pending permit (requires `leave_permits:create`); see README section for approve/reject/complete workflows
- `GET /api/health-statuses` - List active/revoked sick statuses (requires `health_statuses:read`)
- `POST /api/health-statuses` - Create sick status (requires `health_statuses:create`), overrides attendance to `sick`
- `PUT /api/health-statuses/:id/revoke` - Revoke status (requires `health_statuses:revoke`)

> Detail permissions & sample requests tersedia di bagian "Leave Permit Endpoints" dan "Health Status Endpoints" pada README ini.

### Reports (Protected)
- `GET /api/reports/attendance/students?date=YYYY-MM-DD&dormitory_id=...` – Aggregated student attendance (requires `reports:attendance:read`).
- `GET /api/reports/attendance/teachers?date=YYYY-MM-DD&slot_id=...` – Teacher punctuality summary (requires `reports:attendance:read`).
- `GET /api/reports/leave-permits?status=active&type=home_leave` – Security leave metrics (requires `reports:security:read`).
- `GET /api/reports/health-statuses?status=active&dormitory_id=...` – Active/revoked sickness counts (requires `reports:health:read`).
- `GET /api/reports/sks?fan_id=...&is_passed=true` – SKS pass/fail dashboard (requires `reports:academic:read`).
- `GET /api/reports/mutations?student_id=...` – Mutation trail across dorm/class history (requires `reports:academic:read`).

Example – Student Attendance Report

```bash
curl -X GET "http://localhost:8080/api/reports/attendance/students?date=2025-11-21&dormitory_id=<DORM_ID>" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200**

```json
{
  "success": true,
  "message": "Student attendance report retrieved successfully",
  "data": {
    "generated_at": "2025-11-21T23:59:00Z",
    "filters": {
      "date": "2025-11-21",
      "dormitory_id": "<DORM_ID>"
    },
    "rows": [
      {
        "dormitory_id": "<DORM_ID>",
        "class_id": "<CLASS_ID>",
        "fan_id": "<FAN_ID>",
        "total": 30,
        "present": 24,
        "absent": 2,
        "permit": 3,
        "sick": 1
      }
    ]
  }
}
```

### FAN (Protected)
- `GET /api/fans` - List FAN structures (requires `fans:read` permission)
- `GET /api/fans/:id` - Get FAN detail (requires `fans:read` permission)
- `POST /api/fans` - Create FAN (requires `fans:create` permission)
- `PUT /api/fans/:id` - Update FAN (requires `fans:update` permission)
- `DELETE /api/fans/:id` - Delete FAN (requires `fans:delete` permission)

### Classes (Protected)
- `GET /api/classes?fan_id=...` - List classes filtered by FAN (requires `classes:read` permission)
- `GET /api/classes/:id` - Get class detail (requires `classes:read` permission)
- `POST /api/classes` - Create class within a FAN (requires `classes:create` permission)
- `PUT /api/classes/:id` - Update class (requires `classes:update` permission)
- `DELETE /api/classes/:id` - Delete class (requires `classes:delete` permission)
- `POST /api/classes/:id/students` - Enroll student in a class (requires `classes:update` permission)
- `POST /api/classes/:id/staff` - Assign staff to class (requires `classes:update` permission)

### Schedule Slots (Protected)
- `GET /api/schedule-slots?dormitory_id=...` - List slots per dormitory (requires `schedule_slots:read`)
- `GET /api/schedule-slots/:id` - Get specific slot detail (requires `schedule_slots:read`)
- `POST /api/schedule-slots` - Create slot with overlap validation (requires `schedule_slots:create`)
- `PUT /api/schedule-slots/:id` - Update slot meta/time window (requires `schedule_slots:update`)
- `DELETE /api/schedule-slots/:id` - Delete/deactivate slot (requires `schedule_slots:delete`)

### Class Schedules (Protected)
- `GET /api/class-schedules?class_id=...` - List class schedules with filters (requires `class_schedules:read`)
- `GET /api/class-schedules/:id` - Get schedule detail (requires `class_schedules:read`)
- `POST /api/class-schedules` - Create schedule referencing slot or manual time (requires `class_schedules:create`)
- `PUT /api/class-schedules/:id` - Update schedule metadata/time/teacher (requires `class_schedules:update`)
- `DELETE /api/class-schedules/:id` - Delete schedule (requires `class_schedules:delete`)

### SKS Definitions (Protected)
- `GET /api/sks?fan_id=...` - List SKS definitions per FAN (requires `sks_definitions:read`)
- `GET /api/sks/:id` - Get SKS definition detail (requires `sks_definitions:read`)
- `POST /api/sks` - Create SKS definition (requires `sks_definitions:create`)
- `PUT /api/sks/:id` - Update definition attributes/status (requires `sks_definitions:update`)
- `DELETE /api/sks/:id` - Delete definition (requires `sks_definitions:delete`)

### SKS Exam Schedules (Protected)
- `GET /api/sks-exams?sks_id=...` - List exam schedules for a definition (requires `sks_exams:read`)
- `GET /api/sks-exams/:id` - Get exam schedule detail (requires `sks_exams:read`)
- `POST /api/sks-exams` - Create exam schedule with date/time validation (requires `sks_exams:create`)
- `PUT /api/sks-exams/:id` - Update examiner/date/time/location (requires `sks_exams:update`)
- `DELETE /api/sks-exams/:id` - Delete exam schedule (requires `sks_exams:delete`)

### Audit Logs (Protected)
- `GET /api/audit-logs` - List audit logs (with pagination and filters, requires `audit:read` permission)
- `DELETE /api/dormitories/:id` - Delete dormitory (requires dormitory access + `dorm:delete` permission)
- `POST /api/dormitories/:id/users` - Assign staff/user to dormitory (requires dormitory access + `dorm:update` permission)
- `DELETE /api/dormitories/:id/users/:user_id` - Remove staff/user assignment (requires dormitory access + `dorm:update` permission)

### Students (Protected)
- `GET /api/students` - List students (requires `student:read`)
- `GET /api/students/:id` - Get student detail (requires `student:read`)
- `POST /api/students` - Create student (requires `student:create`)
- `PUT /api/students/:id` - Update student (requires `student:update`)
- `PATCH /api/students/:id/status` - Update lifecycle status (requires `student:update`)
- `POST /api/students/:id/mutate-dormitory` - Mutate dormitory assignment and log history (requires `student:update`)
- `POST /api/students/:id/sks-results` - Record SKS result for a student (requires `student_sks_results:create`)
- `PUT /api/students/:id/sks-results/:result_id` - Update recorded SKS result (requires `student_sks_results:update`)
- `GET /api/students/:id/sks-results` - List SKS results for a student with pagination/filter per FAN (requires `student_sks_results:read`)
- `GET /api/students/:id/fans` - View FAN completion status derived from SKS results (requires `student_sks_results:read`)

### Health Check
- `GET /health` - Health check endpoint

## � Contoh Request & Response

Bagian ini memberikan contoh request dan response sukses (1 row data) untuk endpoint utama.

### 1. Authentication

#### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "password": "password123",
    "name": "John Doe"
  }'
```

**Response 201:**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "access_token": "<jwt-token>",
    "refresh_token": "<jwt-refresh-token>",
    "expires_at": "2024-01-01T12:15:00Z",
    "user": {
      "id": "uuid",
      "username": "user",
      "name": "John Doe",
      "roles": ["user"]
    }
  }
}
```

---

### 5. Students

#### Create Student

```bash
curl -X POST http://localhost:8080/api/students \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "student_number": "STD001",
    "full_name": "Integration Student",
    "birth_date": "2010-01-01T00:00:00Z",
    "gender": "male",
    "parent_name": "Parent Doe"
  }'
```

**Response 201:**

```json
{
  "success": true,
  "message": "Student created successfully",
  "data": {
    "id": "uuid",
    "student_number": "STD001",
    "full_name": "Integration Student",
    "birth_date": "2010-01-01T00:00:00Z",
    "gender": "male",
    "parent_name": "Parent Doe",
    "status": "active",
    "is_active": true,
    "created_at": "2025-11-19T02:56:53Z",
    "updated_at": "2025-11-19T02:56:53Z",
    "dormitory_history": []
  }
}
```

#### Patch Student Status

```bash
curl -X PATCH http://localhost:8080/api/students/<STUDENT_ID>/status \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "leave"
  }'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Student status updated successfully",
  "data": {
    "id": "uuid",
    "student_number": "STD001",
    "status": "leave",
    "is_active": false,
    "updated_at": "2025-11-19T03:05:00Z",
    "dormitory_history": []
  }
}
```

#### Mutate Dormitory

```bash
curl -X POST http://localhost:8080/api/students/<STUDENT_ID>/mutate-dormitory \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "dormitory_id": "<DORMITORY_ID>",
    "start_date": "2025-11-19T03:10:00Z"
  }'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Student dormitory mutated successfully",
  "data": {
    "id": "uuid",
    "student_number": "STD001",
    "status": "leave",
    "dormitory_history": [
      {
        "dormitory_id": "<DORMITORY_ID>",
        "start_date": "2025-11-19T03:10:00Z",
        "end_date": ""
      }
    ]
  }
}
```

#### Get Current User

```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Current user retrieved successfully",
  "data": {
    "id": "uuid",
    "username": "user",
    "name": "John Doe",
    "is_active": true,
    "roles": ["user"],
    "permissions": [
      "dorm:read"
    ],
    "dormitories": [],
    "created_at": "2025-11-18T06:04:28+07:00",
    "updated_at": "2025-11-18T06:04:28+07:00"
  }
}
```

### 6. Audit Logs

#### List Audit Logs

```bash
curl -X GET 'http://localhost:8080/api/audit-logs?page=1&page_size=10' \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Audit logs retrieved successfully",
  "data": {
    "logs": [
      {
        "id": "uuid",
        "actor_id": "uuid",
        "actor_username": "admin",
        "actor_roles": ["admin"],
        "action": "user:create",
        "resource": "user",
        "target_id": "uuid",
        "request_path": "/api/users",
        "request_method": "POST",
        "status_code": 201,
        "ip_address": "127.0.0.1",
        "user_agent": "curl/7.79.1",
        "metadata": "{\"username\":\"user\",\"name\":\"John Doe\"}",
        "created_at": "2025-11-18T06:10:00+07:00"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

### 3a. Permissions

#### List Permissions

```bash
curl -X GET 'http://localhost:8080/api/permissions?page=1&page_size=10' \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Permissions retrieved successfully",
  "data": {
    "permissions": [
      {
        "id": "uuid",
        "name": "user:read",
        "slug": "user-read",
        "resource": "user",
        "action": "read"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Response 200:** mirip dengan response register (berisi access token, refresh token, dan user).

---

### 2. Users

#### List Users

```bash
curl -X GET 'http://localhost:8080/api/users?page=1&page_size=10' \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "admin",
        "name": "Admin User",
        "is_active": true,
        "roles": ["admin"]
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### Get User by ID

```bash
curl -X GET http://localhost:8080/api/users/<USER_ID> \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:** satu objek user seperti di atas.

---

### 3. Roles

#### List Roles

```bash
curl -X GET 'http://localhost:8080/api/roles?page=1&page_size=10' \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Roles retrieved successfully",
  "data": {
    "roles": [
      {
        "id": "uuid",
        "name": "Admin",
        "slug": "admin",
        "is_active": true,
        "is_protected": true,
        "permissions": ["user:read", "user:create"]
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

---

### 4. Dormitories

#### List Dormitories

```bash
curl -X GET 'http://localhost:8080/api/dormitories?page=1&page_size=10' \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response 200:**

```json
{
  "success": true,
  "message": "Dormitories retrieved successfully",
  "data": {
    "dormitories": [
      {
        "id": "uuid",
        "name": "Dormitory A",
        "description": "Main dormitory building",
        "is_active": true
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

---

### 5. Location (Public)

#### List Provinces

```bash
curl -X GET 'http://localhost:8080/api/provinces?page=1&page_size=10&search=Aceh'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Provinces retrieved successfully",
  "data": {
    "items": [
      { "id": 1, "name": "Aceh (NAD)", "code": "11" }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### List Regencies by Province

```bash
curl -X GET 'http://localhost:8080/api/regencies?province_id=1&page=1&page_size=10'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Regencies retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "type": "Kabupaten",
        "name": "Aceh Barat",
        "code": "05",
        "full_code": "1105",
        "province_id": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### List Districts by Regency

```bash
curl -X GET 'http://localhost:8080/api/districts?regency_id=420&page=1&page_size=10'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Districts retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Air Majunto",
        "code": "13",
        "full_code": "170613",
        "regency_id": 420
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### List Villages by District

```bash
curl -X GET 'http://localhost:8080/api/villages?district_id=7164&page=1&page_size=10&search=Yawosi'
```

**Response 200:**

```json
{
  "success": true,
  "message": "Villages retrieved successfully",
  "data": {
    "items": [
      {
        "id": 10,
        "name": "Yawosi (Fanindi)",
        "code": "2006",
        "full_code": "9106132006",
        "pos_code": "98552",
        "district_id": 7164
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

## �📤 API Response Format

Semua endpoint menggunakan format response yang standar untuk memastikan konsistensi dan kemudahan integrasi.

### Success Response

Format response sukses mengikuti struktur berikut:

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    // Response data sesuai endpoint
  }
}
```

**Contoh:**
- `200 OK` - `SuccessOK()` - untuk GET, PUT yang berhasil
- `201 Created` - `SuccessCreated()` - untuk POST yang berhasil
- `204 No Content` - `SuccessNoContent()` - untuk DELETE yang berhasil

### Error Response

Format response error mengikuti struktur berikut:

```json
{
  "success": false,
  "message": "User not found",
  "error": "optional error detail"
}
```

**HTTP Status Codes:**
- `400 Bad Request` - `ErrorBadRequest()` - Request tidak valid
- `401 Unauthorized` - `ErrorUnauthorized()` - Tidak terautentikasi
- `403 Forbidden` - `ErrorForbidden()` - Tidak memiliki izin
- `404 Not Found` - `ErrorNotFound()` - Resource tidak ditemukan
- `409 Conflict` - `ErrorConflict()` - Konflik data (misal: username sudah terdaftar)
- `500 Internal Server Error` - `ErrorInternalServer()` - Error server

### Response Helper Functions

Semua helper functions tersedia di package `internal/interfaces/http/response`:

```go
// Success responses
response.SuccessOK(c, data, "message")
response.SuccessCreated(c, data, "message")
response.SuccessNoContent(c)

// Error responses
response.ErrorBadRequest(c, "message", "errorDetail")
response.ErrorUnauthorized(c, "message", "errorDetail")
response.ErrorForbidden(c, "message", "errorDetail")
response.ErrorNotFound(c, "message", "errorDetail")
response.ErrorConflict(c, "message", "errorDetail")
response.ErrorInternalServer(c, "message", "errorDetail")
```

## 🔑 Authentication

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-01T12:15:00Z",
    "user": {
      "id": "uuid",
      "username": "admin",
      "name": "Admin User",
      "roles": ["admin"]
    }
  }
}
```

**Error Response Example:**
```json
{
  "success": false,
  "message": "Invalid credentials",
  "error": ""
}
```

### Using Access Token
```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🎯 Permission System

### Attendance Endpoints

| Method & Path | Description | Required Permissions |
| --- | --- | --- |
| `POST /api/attendance-sessions/open` | Generate/open attendance sessions for the given class schedules on a specific date. | `attendance_sessions:create` |
| `POST /api/attendance-sessions/:id/students` | Bulk submit/update student attendance records for a session. | `attendance_sessions:update` |
| `POST /api/attendance-sessions/:id/teacher` | Submit teacher attendance (manual override). | `attendance_sessions:update` |
| `POST /api/attendance-sessions/lock-day` | Lock all sessions for a given date (cron-friendly). | `attendance_sessions:lock` |
| `GET /api/attendance-sessions` | List attendance sessions with filters (schedule, teacher, status, date). | `attendance_sessions:read` |

**Open Sessions Example**

```bash
curl -X POST http://localhost:8080/api/attendance-sessions/open \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "class_schedule_ids": ["9f0ad9cf-7db1-44fa-b522-4ddcbde9d6d7"],
    "date": "2025-11-20"
  }'
```

**Submit Student Attendance Example**

```bash
curl -X POST http://localhost:8080/api/attendance-sessions/{session_id}/students \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "records": [
      {"student_id": "63b2d272-4e61-4cc2-9fd1-6e327fa8f424", "status": "present"},
      {"student_id": "8a27a91a-7db7-4c64-b89f-5a8f82a8f457", "status": "permit", "note": "Competition"}
    ]
  }'
```

**Lock Sessions Example (cron CLI)**

```bash
go run cmd/attendance_lock/main.go -date 2025-11-20
```

### Leave Permit Endpoints

| Method & Path | Description | Required Permissions |
| --- | --- | --- |
| `GET /api/leave-permits` | List leave permits with optional filters (student, status, type, date). | `leave_permits:read` |
| `POST /api/leave-permits` | Create a leave permit (pending status). | `leave_permits:create` |
| `PUT /api/leave-permits/:id/approve` | Approve a pending permit. | `leave_permits:approve` |
| `PUT /api/leave-permits/:id/reject` | Reject a pending permit. | `leave_permits:approve` |
| `PUT /api/leave-permits/:id/complete` | Mark an approved permit as completed (student returned). | `leave_permits:complete` |

**Create Leave Permit Example**

```bash
curl -X POST http://localhost:8080/api/leave-permits \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "63b2d272-4e61-4cc2-9fd1-6e327fa8f424",
    "type": "official_duty",
    "reason": "National competition",
    "start_date": "2025-11-25",
    "end_date": "2025-11-28"
  }'
```

**Approve Leave Permit Example**

```bash
curl -X PUT http://localhost:8080/api/leave-permits/{permit_id}/approve \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json"
```

### Health Status Endpoints

| Method & Path | Description | Required Permissions |
| --- | --- | --- |
| `GET /api/health-statuses` | List health (sick) statuses with filters (student, status, date). | `health_statuses:read` |
| `POST /api/health-statuses` | Create an active health status (auto-marks student as sick in attendance). | `health_statuses:create` |
| `PUT /api/health-statuses/:id/revoke` | Revoke an active health status when the student recovers. | `health_statuses:revoke` |

**Create Health Status Example**

```bash
curl -X POST http://localhost:8080/api/health-statuses \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "63b2d272-4e61-4cc2-9fd1-6e327fa8f424",
    "diagnosis": "Influenza",
    "notes": "Needs bed rest",
    "start_date": "2025-11-21"
  }'
```

**Revoke Health Status Example**

```bash
curl -X PUT http://localhost:8080/api/health-statuses/{status_id}/revoke \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json"
```

### Default Permissions

**User Permissions:**
- `user:read` - Read users
- `user:create` - Create users
- `user:update` - Update users
- `user:delete` - Delete users

**Dormitory Permissions:**
- `dorm:read` - Read dormitories
- `dorm:create` - Create dormitories
- `dorm:update` - Update dormitories
- `dorm:delete` - Delete dormitories

**Role Permissions:**
- `role:read` - Read roles
- `role:create` - Create roles
- `role:update` - Update roles
- `role:delete` - Delete roles

**Attendance Permissions:**
- `attendance_sessions:read` – List attendance sessions and view details
- `attendance_sessions:create` – Open sessions automatically from schedules
- `attendance_sessions:update` – Submit/update student & teacher attendance
- `attendance_sessions:lock` – Lock sessions for a day (cron / admin action)

### Default Roles

- **user** (default role, not protected)
  - Has `dorm:read` only
  - Automatically assigned to new users if no role specified
  - Can be modified (permissions can be changed)

- **admin** (protected role)
  - Has all permissions (user:*, dorm:*, role:*)
  - Protected: Cannot modify permissions or delete
  - Use for administrative access

- **super_admin** (protected role)
  - Has all permissions (user:*, dorm:*, role:*)
  - Protected: Cannot modify permissions or delete
  - Use for super administrative access

- **academic_sks**
  - Focused on academic operations (SKS + attendance)
  - Granted `attendance_sessions:read`, `attendance_sessions:create`, `attendance_sessions:update`
  - Intended for academic officers managing day-to-day attendance submissions

- **attendance_cron**
  - Service account role for the cron/CLI job that locks sessions nightly
  - Granted `attendance_sessions:read` and `attendance_sessions:lock`
  - Assign to automation users or CI jobs invoking `cmd/attendance_lock`

### Protected Roles

Protected roles (`admin` dan `super_admin`) memiliki batasan:
- ❌ Tidak bisa diubah permission-nya (assign/remove permission)
- ❌ Tidak bisa dihapus
- ✅ Bisa diubah nama, slug, dan status aktif
- ✅ Bisa di-assign ke user

### Role Management Features

1. **Create Role**: Buat role baru dengan permission tertentu
2. **Update Role**: Ubah nama, slug, atau status role
3. **Delete Role**: Hapus role (kecuali protected roles)
4. **Assign Permission**: Tambahkan permission ke role
5. **Remove Permission**: Hapus permission dari role
6. **Assign Role to User**: Berikan role ke user
7. **Remove Role from User**: Hapus role dari user

### Default Role Assignment

Saat membuat user baru:
- Jika `role_ids` tidak diberikan → otomatis mendapat role "user"
- Jika `role_ids` diberikan → mendapat role sesuai yang ditentukan

## 🛡️ Guard System

Guard system mengontrol akses user ke dormitory tertentu:

1. **Admin/Super Admin** - Dapat mengakses semua dormitory
2. **Staff/User dengan assignment** - Hanya dapat mengakses dormitory yang di-assign ke mereka

## 📝 Role Management Examples

### Create Role
```bash
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Manager",
    "slug": "manager",
    "is_active": true,
    "is_protected": false,
    "permission_ids": ["permission-uuid-1", "permission-uuid-2"]
  }'
```

### Assign Permission to Role
```bash
curl -X POST http://localhost:8080/api/roles/{role_id}/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_id": "permission-uuid"
  }'
```

### Assign Role to User
```bash
curl -X POST http://localhost:8080/api/users/{user_id}/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": "role-uuid"
  }'
```

### Remove Role from User
```bash
curl -X DELETE http://localhost:8080/api/users/{user_id}/roles/{role_id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🧪 Testing

```bash
# Run tests
make test

# Or
go test ./...

# Run HTTP integration tests (includes dormitory & student flows)
go test ./internal/interfaces/http
```

## 📦 Build

```bash
# Build binary
make build

# Output akan di bin/server
```

## 🔧 Development

### Project Structure Best Practices
- **Domain Layer**: Pure business logic, no dependencies
- **Application Layer**: Use cases, depends only on domain
- **Infrastructure Layer**: External concerns (DB, HTTP, etc.)
- **Interface Layer**: HTTP handlers, depends on application layer

## 🌍 Fitur Lokasi (Province/Regency/District/Village)

Project ini menyediakan fitur lokasi Indonesia (provinsi, kabupaten/kota, kecamatan, desa/kelurahan) sebagai **data referensi read-only**.

Ringkasan:

- Endpoint public (tanpa auth) di bawah `/api`:
  - `GET /api/provinces`, `GET /api/provinces/:id`
  - `GET /api/regencies`, `GET /api/regencies/:id`
  - `GET /api/districts`, `GET /api/districts/:id`
  - `GET /api/villages`, `GET /api/villages/:id`
- Mendukung pagination (`page`, `page_size`) dan pencarian dengan `search` (berdasarkan `name`).
- Data diimport dari file JSON melalui command CLI khusus.

Detail lengkap schema, contoh JSON, dan cara import:

- Lihat **`docs/location_feature.md`**.

## 🧩 Panduan Menambahkan Fitur Baru

Panduan lengkap dan cukup panjang tentang cara menambah fitur baru (entity, repository, usecase, handler, routes) telah dipindahkan ke dokumen terpisah agar README tetap ringkas.

- Lihat **`docs/adding_features.md`** untuk panduan step-by-step menambahkan fitur baru mengikuti pola Clean/Hexagonal Architecture.

	return router
}
```

---

### **Step 8: Register in main.go**

Update `cmd/main.go`:

```go
// Initialize repositories
productRepo := infraRepo.NewProductRepository()
// ... existing repos ...

// Initialize use cases
productUseCase := usecase.NewProductUseCase(productRepo)
// ... existing use cases ...

// Initialize handlers
productHandler := handler.NewProductHandler(productUseCase)
// ... existing handlers ...

// Setup router
r := router.SetupRouter(
	authHandler, 
	userHandler, 
	dormitoryHandler, 
	roleHandler,
	productHandler, // Add this
	authMiddleware,
)
```

---

### **Step 9: Add Migration**

Update `internal/infrastructure/database/migrations.go`:

```go
// Add new migration in init() function
RegisterMigration(
	"003_create_products_table",
	"Create products table",
	func(db *gorm.DB) error {
		return db.AutoMigrate(&entity.Product{})
	},
	func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entity.Product{})
	},
)
```

---

### **Step 10: Add Permissions (Optional)**

Jika fitur memerlukan permission, tambahkan di seed (`cmd/seed/main.go`):

```go
// Add product permissions
{ID: uuid.New(), Name: "product:read", Slug: "product-read", Resource: "product", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
{ID: uuid.New(), Name: "product:create", Slug: "product-create", Resource: "product", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
{ID: uuid.New(), Name: "product:update", Slug: "product-update", Resource: "product", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
{ID: uuid.New(), Name: "product:delete", Slug: "product-delete", Resource: "product", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
```

---

### **Step 11: Testing**

Buat test file di `internal/application/usecase/product_usecase_test.go`:

```go
package usecase

import (
	"context"
	"testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	req := dto.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Stock:       10,
	}

	// Mock expectations
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil)

	// Execute
	result, err := uc.CreateProduct(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	mockRepo.AssertExpectations(t)
}
```

---

### **📝 Checklist Menambahkan Fitur Baru**

- [ ] 1. Buat entity di `internal/domain/entity/`
- [ ] 2. Buat repository interface di `internal/domain/repository/`
- [ ] 3. Implement repository di `internal/infrastructure/repository/`
- [ ] 4. Buat DTOs di `internal/application/dto/`
- [ ] 5. Buat use case di `internal/application/usecase/`
- [ ] 6. Tambahkan domain errors jika diperlukan
- [ ] 7. Buat handler di `internal/interfaces/http/handler/`
- [ ] 8. Tambahkan routes di `internal/interfaces/http/router/router.go`
- [ ] 9. Register di `cmd/main.go`
- [ ] 10. Tambahkan migration di `internal/infrastructure/database/migrations.go`
- [ ] 11. Tambahkan permissions di seed (jika diperlukan)
- [ ] 12. Buat unit tests
- [ ] 13. Update README dengan dokumentasi endpoint baru

---

### **🎯 Best Practices**

1. **Separation of Concerns**: Setiap layer hanya fokus pada concern-nya sendiri
2. **Dependency Rule**: Inner layer tidak boleh bergantung pada outer layer
3. **Use Interfaces**: Repository interface di domain, implementasi di infrastructure
4. **Error Handling**: Gunakan domain errors untuk error yang konsisten
5. **Standardized Response**: Selalu gunakan response helpers untuk konsistensi
6. **Validation**: Validasi input di handler dan use case
7. **Context**: Selalu gunakan context untuk cancellation dan timeout
8. **Testing**: Buat test untuk use case dan handler

---

### **🔍 Contoh Alur Request**

```
1. HTTP Request → Handler (parse request, validate)
2. Handler → Use Case (business logic)
3. Use Case → Repository Interface (port)
4. Repository Implementation → Database (adapter)
5. Database → Repository Implementation
6. Repository Implementation → Use Case
7. Use Case → Handler (convert to DTO)
8. Handler → HTTP Response (standardized format)
```

---

### **📚 Referensi File Structure**

```
internal/
├── domain/
│   ├── entity/
│   │   └── product.go          # Step 1
│   ├── repository/
│   │   └── product_repository.go # Step 2
│   └── errors/
│       └── errors.go            # Add errors
├── application/
│   ├── dto/
│   │   └── product_dto.go       # Step 4
│   └── usecase/
│       └── product_usecase.go   # Step 5
├── infrastructure/
│   ├── repository/
│   │   └── product_repository.go # Step 3
│   └── database/
│       └── migrations.go        # Step 9
└── interfaces/
    └── http/
        ├── handler/
        │   └── product_handler.go # Step 6
        └── router/
            └── router.go          # Step 7
```

---

### **💡 Tips**

- **Mulai dari Domain**: Selalu mulai dari domain layer (entity, repository interface)
- **Test Incrementally**: Test setiap layer setelah dibuat
- **Follow Naming Convention**: Gunakan naming yang konsisten
- **Document Complex Logic**: Tambahkan comment untuk logic yang kompleks
- **Keep It Simple**: Jangan over-engineer, mulai dengan yang sederhana

### Using Standardized Responses

Saat membuat handler baru, selalu gunakan helper functions dari package `response`:

```go
import "github.com/your-org/go-backend-starter/internal/interfaces/http/response"

// Success response
response.SuccessOK(c, data, "Operation successful")
response.SuccessCreated(c, data, "Resource created")

// Error response
response.ErrorBadRequest(c, "Invalid input", err.Error())
response.ErrorNotFound(c, "Resource not found")
```

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Support

Untuk pertanyaan atau support, silakan buat issue di repository ini.

---

**Happy Coding! 🚀**
