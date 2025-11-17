

## **Rangkuman Starter Project Backend Go (Hexagonal Architecture)**

Starter project ini menyediakan fondasi backend berbasis **Golang** dengan **Clean Architecture / Hexagonal Architecture** yang mendukung skalabilitas, maintainability, dan pemisahan concern yang jelas antara domain, aplikasi, dan infrastruktur. Proyek ini ditujukan untuk sistem yang memiliki otentikasi JWT, manajemen user, dan entitas dormitory dengan pengaturan izin akses yang fleksibel.

---

## **Fitur Utama**

### **1. Authentication & Authorization**

* **JWT Authentication**

  * Access Token
  * Refresh Token
* **Registrasi dan Login User**
* **Refresh Token Endpoint**
* Middleware verifikasi JWT untuk melindungi endpoint.

### **2. Role & Permission System**

* Mendukung role-based dan permission-based authorization.
* Contoh permission:

  * `user:read`
  * `user:update`
  * `dorm:read`
  * `dorm:update`
* Role dapat memiliki banyak permission.
* User dapat memiliki satu atau lebih role.

### **3. User Management (CRUD Users)**

* Create, Read, Update, Delete user untuk admin.
* Pemisahan antara profile user pribadi vs admin management.

### **4. Dormitory Management (CRUD Dormitory)**

* CRUD untuk data dormitory.
* Setiap dormitory dapat dibatasi akses berdasarkan guard.

### **5. Guard / Access Control**

* Guard menentukan batas akses user terhadap dormitory:

  * **Access to specific dormitories only** — misal staff hanya dapat mengelola dormitory tertentu.
  * **Access to all dormitories** — misal admin pusat dapat membaca/mengelola seluruh dormitory.
* Guard bekerja berdasarkan ID dormitory yang terasosiasi ke user atau role.

---

## **Arsitektur Clean (Hexagonal Architecture)**

Struktur mengikuti prinsip **ports & adapters** agar mudah diuji, scalable, dan bebas ketergantungan dari framework.

### **Lapisan Utama**

#### **1. Domain Layer**

* Berisi **entity**, **value object**, **domain service**, dan **business rules**.
* Contoh entity:

  * `User`
  * `Role`
  * `Permission`
  * `Dormitory`
* Tidak bergantung pada database atau framework.

#### **2. Application Layer (Use Cases)**

* Berisi **service/use case** seperti:

  * `RegisterUser`
  * `LoginUser`
  * `RefreshToken`
  * `CreateDormitory`, `UpdateDormitory`, dll.
  * `AssignRoleToUser`
* Menggunakan **interface repository** (port) yang diimplementasikan di infrastruktur.
* Mengatur flow bisnis tetapi tidak tahu detail teknologi.

#### **3. Infrastructure Layer (Adapters)**

* Implementasi repository, misalnya:

  * PostgreSQL / MySQL / MongoDB repository
  * Token provider JWT
* Router/HTTP adapter (Fiber, Gin, Echo, atau standard net/http)

#### **4. Interface/Delivery Layer**

* Controller/handler HTTP
* Middleware:

  * JWT Auth
  * Permission checker
  * Dormitory guard
* Mapping antara request/response ke DTO.

---

## **Flow Authorization**

1. Request masuk → Middleware cek JWT.
2. Middleware cek **role & permission** sesuai endpoint.
3. Jika endpoint terkait dormitory → Guard cek:

   * User memiliki akses ke dormitory id tertentu
   * atau user memiliki akses global (“ALL_DORM_ACCESS”)
4. Jika lolos → dilanjutkan ke handler.

---

## **Keunggulan Desain**

* **Modular & scalable** → logic aplikasi terpisah dari detail teknis.
* **Testable** → domain dan application layer mudah untuk unit testing.
* **Maintainable** → perubahan pada database, framework, atau library JWT tidak mempengaruhi domain/application.
* **Flexible authorization** melalui kombinasi **role, permission, dan guard**.

---

Jika Anda ingin, saya dapat bantu:

✅ buatkan struktur folder lengkap
✅ buatkan blueprint kode (handler, use case, repository interface, entity)
✅ buatkan contoh implementasi JWT + middleware authorization
✅ buatkan ERD atau schema database

Tinggal beri tahu!
