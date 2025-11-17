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
- ✅ Contoh permission: `user:read`, `user:update`, `dorm:read`, `dorm:update`
- ✅ Role dapat memiliki banyak permission
- ✅ User dapat memiliki satu atau lebih role

### 3. User Management (CRUD Users)
- ✅ Create, Read, Update, Delete user
- ✅ Pagination support
- ✅ Role assignment

### 4. Dormitory Management (CRUD Dormitory)
- ✅ CRUD untuk data dormitory
- ✅ Setiap dormitory dapat dibatasi akses berdasarkan guard

### 5. Guard / Access Control
- ✅ Guard menentukan batas akses user terhadap dormitory:
  - **Access to specific dormitories only** — staff hanya dapat mengelola dormitory tertentu
  - **Access to all dormitories** — admin dapat mengelola seluruh dormitory

### 6. Standardized API Response
- ✅ Response format yang konsisten untuk semua endpoint
- ✅ Success response dengan struktur standar
- ✅ Error response dengan struktur standar
- ✅ Helper functions untuk berbagai HTTP status codes

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
```

### 4. Setup Database
```bash
# Create PostgreSQL database
createdb go_backend_db
```

### 5. Run Migrations
Migrations akan berjalan otomatis saat aplikasi start, atau bisa dijalankan manual:
```bash
go run cmd/main.go
```

### 6. Seed Data (Optional)
```bash
go run cmd/seed/main.go
```

Ini akan membuat:
- Permissions: `user:read`, `user:create`, `user:update`, `user:delete`, `dorm:read`, `dorm:create`, `dorm:update`, `dorm:delete`
- Roles: `admin`, `staff`, `user`
- Admin user: `admin@example.com` / `admin123`
- Sample dormitories

### 7. Run Application
```bash
# Using Make
make run

# Or directly
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080`

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

### Dormitories (Protected)
- `GET /api/dormitories` - List dormitories (with pagination)
- `GET /api/dormitories/:id` - Get dormitory by ID (requires dormitory access)
- `POST /api/dormitories` - Create dormitory (requires `dorm:create` permission)
- `PUT /api/dormitories/:id` - Update dormitory (requires dormitory access + `dorm:update` permission)
- `DELETE /api/dormitories/:id` - Delete dormitory (requires dormitory access + `dorm:delete` permission)

### Health Check
- `GET /health` - Health check endpoint

## 📤 API Response Format

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
- `409 Conflict` - `ErrorConflict()` - Konflik data (misal: email sudah terdaftar)
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
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
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
      "email": "admin@example.com",
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

### Default Permissions
- `user:read` - Read users
- `user:create` - Create users
- `user:update` - Update users
- `user:delete` - Delete users
- `dorm:read` - Read dormitories
- `dorm:create` - Create dormitories
- `dorm:update` - Update dormitories
- `dorm:delete` - Delete dormitories

### Default Roles
- **admin** - Has all permissions
- **staff** - Has `user:read`, `dorm:read`, `dorm:update`
- **user** - Has `dorm:read` only

## 🛡️ Guard System

Guard system mengontrol akses user ke dormitory tertentu:

1. **Admin/Super Admin** - Dapat mengakses semua dormitory
2. **Staff/User dengan assignment** - Hanya dapat mengakses dormitory yang di-assign ke mereka

Untuk assign dormitory ke user, gunakan endpoint:
```bash
# Assign dormitory to user (via database atau buat endpoint khusus)
```

## 🧪 Testing

```bash
# Run tests
make test

# Or
go test ./...
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

### Adding New Features
1. Define entity in `internal/domain/entity/`
2. Create repository interface in `internal/domain/repository/`
3. Implement repository in `internal/infrastructure/repository/`
4. Create use case in `internal/application/usecase/`
5. Create handler in `internal/interfaces/http/handler/`
6. Add routes in `internal/interfaces/http/router/router.go`
7. **Gunakan standardized response helpers** dari `internal/interfaces/http/response` untuk semua response

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
