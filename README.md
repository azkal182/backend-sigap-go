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
  - Admin: `admin@example.com` / `admin123`
  - Super Admin: `superadmin@example.com` / `superadmin123`
- **Sample dormitories**

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
- `POST /api/users/:id/roles` - Assign role to user (requires `user:update` permission)
- `DELETE /api/users/:id/roles/:role_id` - Remove role from user (requires `user:update` permission)

### Roles (Protected)
- `GET /api/roles` - List roles (with pagination, requires `role:read` permission)
- `GET /api/roles/:id` - Get role by ID (requires `role:read` permission)
- `POST /api/roles` - Create role (requires `role:create` permission)
- `PUT /api/roles/:id` - Update role (requires `role:update` permission)
- `DELETE /api/roles/:id` - Delete role (requires `role:delete` permission, protected roles cannot be deleted)
- `POST /api/roles/:id/permissions` - Assign permission to role (requires `role:update` permission, protected roles cannot be modified)
- `DELETE /api/roles/:id/permissions` - Remove permission from role (requires `role:update` permission, protected roles cannot be modified)

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

## 🚀 Panduan Menambahkan Fitur Baru

Panduan lengkap untuk menambahkan fitur baru mengikuti Clean Architecture / Hexagonal Architecture pattern.

### 📋 Overview

Struktur Clean Architecture terdiri dari 4 layer:
1. **Domain Layer** - Business logic murni, tidak bergantung pada framework
2. **Application Layer** - Use cases dan business rules
3. **Infrastructure Layer** - Implementasi konkret (database, external services)
4. **Interface Layer** - HTTP handlers, API endpoints

### 🎯 Contoh: Menambahkan Fitur "Product Management"

Kita akan membuat CRUD untuk Product sebagai contoh lengkap.

---

### **Step 1: Define Domain Entity**

Buat entity di `internal/domain/entity/product.go`:

```go
package entity

import (
	"time"
	"github.com/google/uuid"
)

// Product represents a product entity in the domain
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "products"
}
```

**Best Practices:**
- Entity harus pure Go struct, tidak ada dependency pada framework
- Gunakan UUID untuk ID
- Tambahkan timestamps (CreatedAt, UpdatedAt)
- Untuk soft delete, gunakan DeletedAt (pointer)

---

### **Step 2: Create Repository Interface (Port)**

Buat interface di `internal/domain/repository/product_repository.go`:

```go
package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Product, int64, error)
}
```

**Best Practices:**
- Interface didefinisikan di domain layer (port)
- Gunakan context untuk cancellation dan timeout
- Return error untuk error handling yang konsisten

---

### **Step 3: Implement Repository (Adapter)**

Buat implementasi di `internal/infrastructure/repository/product_repository.go`:

```go
package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository() repository.ProductRepository {
	return &productRepository{
		db: database.DB,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	err := r.db.WithContext(ctx).Model(&entity.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}
```

**Best Practices:**
- Implementasi repository menggunakan GORM
- Selalu gunakan `WithContext(ctx)` untuk context propagation
- Handle error dengan baik

---

### **Step 4: Create DTOs**

Buat DTOs di `internal/application/dto/product_dto.go`:

```go
package dto

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ProductResponse represents product data in responses
type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ListProductsResponse represents paginated product list response
type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
```

**Best Practices:**
- DTOs untuk request/response terpisah dari entity
- Gunakan validation tags (binding:"required", dll)
- Untuk optional fields, gunakan pointer atau omitempty

---

### **Step 5: Create Use Case**

Buat use case di `internal/application/usecase/product_usecase.go`:

```go
package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// ProductUseCase handles product management use cases
type ProductUseCase struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &entity.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.toProductResponse(product), nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrProductNotFound
	}

	return uc.toProductResponse(product), nil
}

// UpdateProduct updates a product
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrProductNotFound
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	product.UpdatedAt = time.Now()

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.toProductResponse(product), nil
}

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return domainErrors.ErrProductNotFound
	}

	return uc.productRepo.Delete(ctx, id)
}

// ListProducts retrieves a paginated list of products
func (uc *ProductUseCase) ListProducts(ctx context.Context, page, pageSize int) (*dto.ListProductsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	products, total, err := uc.productRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	productResponses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, *uc.toProductResponse(product))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// toProductResponse converts entity.Product to dto.ProductResponse
func (uc *ProductUseCase) toProductResponse(product *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}
```

**Best Practices:**
- Use case berisi business logic
- Validasi input (page, pageSize, dll)
- Convert entity ke DTO untuk response
- Gunakan domain errors untuk error handling

**Jangan lupa tambahkan error di `internal/domain/errors/errors.go`:**
```go
ErrProductNotFound = errors.New("product not found")
```

---

### **Step 6: Create HTTP Handler**

Buat handler di `internal/interfaces/http/handler/product_handler.go`:

```go
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// ProductHandler handles product management requests
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.productUseCase.CreateProduct(c.Request.Context(), req)
	if err != nil {
		switch err {
		default:
			response.ErrorInternalServer(c, "Failed to create product", err.Error())
		}
		return
	}

	response.SuccessCreated(c, resp, "Product created successfully")
}

// GetProduct handles getting a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid product ID", err.Error())
		return
	}

	resp, err := h.productUseCase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrProductNotFound:
			response.ErrorNotFound(c, "Product not found")
		default:
			response.ErrorInternalServer(c, "Failed to get product", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Product retrieved successfully")
}

// UpdateProduct handles product update
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid product ID", err.Error())
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.productUseCase.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrProductNotFound:
			response.ErrorNotFound(c, "Product not found")
		default:
			response.ErrorInternalServer(c, "Failed to update product", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Product updated successfully")
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid product ID", err.Error())
		return
	}

	err = h.productUseCase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrProductNotFound:
			response.ErrorNotFound(c, "Product not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete product", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// ListProducts handles listing products with pagination
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.productUseCase.ListProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list products", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Products retrieved successfully")
}
```

**Best Practices:**
- Handler hanya untuk HTTP concerns (parsing request, calling use case, formatting response)
- Selalu gunakan standardized response helpers
- Handle error dengan switch case untuk error yang berbeda
- Parse dan validate input di handler

---

### **Step 7: Add Routes**

Update `internal/interfaces/http/router/router.go`:

```go
// Update SetupRouter function signature
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	dormitoryHandler *handler.DormitoryHandler,
	roleHandler *handler.RoleHandler,
	productHandler *handler.ProductHandler, // Add this
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	// ... existing code ...

	// Add product routes
	protected := api.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		// ... existing routes ...

		// Product routes
		products := protected.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("", authMiddleware.RequirePermission("product:create"), productHandler.CreateProduct)
			products.PUT("/:id", authMiddleware.RequirePermission("product:update"), productHandler.UpdateProduct)
			products.DELETE("/:id", authMiddleware.RequirePermission("product:delete"), productHandler.DeleteProduct)
		}
	}

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
