# OpenAPI Hybrid TODO Plan (Option 3)

Dokumen ini memetakan seluruh rute HTTP (berdasarkan `internal/interfaces/http/router/router.go`) dan menyusun backlog bertahap untuk menyiapkan spesifikasi OpenAPI hybrid: kombinasi dokumentasi manual + referensi langsung ke DTO (`internal/application/dto`). Gunakan checklist ini agar pekerjaan terdokumentasi, terukur, dan mudah dibagi lintas tim.

## 1. Route Inventory (Referensi Awal)
| Area | Prefix / Rute | Catatan Referensi |
| --- | --- | --- |
| Health | `GET /health` | Respons standar @router/router.go#41-45 |
| Auth (public) | `POST /api/auth/register`, `login`, `refresh` | Handler dengan anotasi Swagger awal @handler/auth_handler.go#31-125 |
| Locations (public) | `GET /api/provinces`, `/regencies`, `/districts`, `/villages` (termasuk detail by `:id`) | Tidak butuh auth; cocok sebagai contoh endpoint publik |
| Current user | `GET /api/me` | Profil user terautentikasi |
| Audit logs | `GET /api/audit-logs` | Wajib permission `audit:read` |
| Students | CRUD + `PATCH /status`, `POST /mutate-dormitory`, SKS result sub-routes | Banyak hubungan ke DTO mahasiswa & SKS |
| Users | CRUD + assign/remove roles | |
| Roles | CRUD + assign/remove permissions | |
| Permissions | `GET /api/permissions` | Reuse data tabel permissions |
| Dormitories | CRUD + assign/remove users, guard middleware | Perlu dokumentasi guard rules |
| FAN | CRUD FAN | |
| Classes | CRUD, enroll student, assign staff | |
| Teachers | CRUD/deactivate teacher | |
| Leave Permits | CRUD-like workflow (approve/reject/complete) | |
| Health Statuses | Create/revoke + filter listing | |
| Schedule Slots | CRUD | |
| Class Schedules | CRUD | |
| Attendance Sessions | List, open, submit student/teacher, lock-day | Perlu deskripsi filter/query |
| SKS Definitions | CRUD | |
| SKS Exams | CRUD | |
| Reports | `/reports/attendance/students|teachers`, `/reports/leave-permits`, `/reports/health-statuses`, `/reports/sks`, `/reports/mutations` | Semua `GET` dengan filter ketat |

> Gunakan README untuk contoh payload dan respon dasar (@README.md#249-570, #770-886, dsb.), lalu kaitkan ke DTO terkait.

## 2. Backlog TODO (Hybrid Spec)
### Phase 0 – Fondasi Format
- [x] Pilih format utama (`openapi.yaml` di `docs/api_spec.md` atau file baru) dan tentukan struktur folder.
- [x] Tambahkan metadata dasar (title, version, servers, securitySchemes bearer). 
- [x] Buat script `make openapi-sync` untuk regenerasi/validasi (placeholder bolehan jalankan templating Go + lint). 

### Phase 1 – Identitas & Publik
- [ ] Dokumentasikan auth endpoints (register/login/refresh) lengkap dengan schema `AuthResponse` merujuk `dto.AuthResponse`.
- [ ] Tambahkan section `CurrentUser` (`GET /api/me`) + `Permissions` (list) dengan referensi ke DTO user summary.
- [ ] Dokumentasikan lokasi publik (`/provinces`, dll.) termasuk query `page`, `page_size`, `search` & struktur data JSON di `data/locations/*.json`.

### Phase 2 – User, Role, Permission Domain
- [x] Ekstrak DTO terkait (`UserDTO`, `RoleDTO`, `PermissionDTO`) sebagai komponen schemas.
- [x] Mapping seluruh routes `/users`, `/roles`, `/permissions` dalam spec (method, params, body, responses, error cases `409`, `403`).
- [x] Catat constraint khusus (role protected, default role `user`).

### Phase 3 – Dormitories & Guarded Resources
- [x] Tambahkan schema `Dormitory`, `DormitoryAssignmentRequest` dari DTO.
- [x] Dokumentasikan middleware guard behavior di deskripsi endpoint `GET/PUT/DELETE /dormitories/:id`.
- [x] Sertakan parameter path/query (mis. `dormitory_id` untuk listing) serta header requirement.

### Phase 4 – Academic Structure (Students, FAN, Classes, Teachers)
- [x] Definisikan schema `Student`, `StudentStatusUpdate`, `StudentSKSResult`, `Fan`, `Class`, `Teacher`.
- [x] Dokumentasikan sub-routes nested (mis. `/students/:id/sks-results`, `/classes/:id/students`).
- [x] Tambahkan contoh request/respons baru (boleh merujuk contoh README) dan error per validasi.

### Phase 5 – Schedule & Attendance
- [x] Schema `ScheduleSlot`, `ClassSchedule`, `AttendanceSession`, `AttendanceRecord`.
- [x] Dokumentasikan endpoints `/schedule-slots`, `/class-schedules`, `/attendance-sessions` termasuk query filter (`date`, `teacher_id`, `status`).
- [x] Jelaskan alur `open -> submit student -> submit teacher -> lock-day` sebagai sub-bab.

### Phase 6 – SKS Definisi & Ujian
- [x] Schema `SKSDefinition`, `SKSExamSchedule` lengkap field relasional (`fan_id`, `subject_id`).
- [x] Dokumentasikan query param (`fan_id`, `sks_id`, `is_active`).

### Phase 7 – Leave, Health, Reports
- [x] Schema `LeavePermit`, `LeavePermitApproval`, `HealthStatus`, `HealthStatusRevoke`.
- [x] Pastikan workflow status (pending -> approved/rejected -> completed) dijabarkan.
- [x] Untuk reports, definisikan schema agregasi (mis. `AttendanceReportRow`, `LeaveReportRow`). Jelaskan filter wajib/opsional dan contoh response ringkas.

### Phase 8 – QA & Automasi
- [x] Tambahkan lint step (mis. `spectral lint docs/api_spec.yaml`) di GitHub Actions.
- [x] Buat guideline review: perubahan `internal/application/dto` atau handler wajib menjalankan `make openapi-sync` dan commit diff spec.
- [x] Dokumentasikan prosedur di README/CONTRIBUTING agar setiap PR menjaga sinkronisasi.

## 3. Tracking & Pembagian Kerja
- Gunakan checkbox di atas untuk menandai progres. 
- Bisa pecah per fase menjadi tiket terpisah (Phase 1-8). 
- Setiap fase minimal menghasilkan PR: (1) update spec, (2) update docs/api_spec.md dengan contoh tambahan, (3) update automation bila perlu.

Dengan roadmap ini, tim bisa mulai dari endpoint kritikal (Auth/User) lalu merambah modul lainnya secara bertahap tanpa melewatkan referensi DTO atau aturan permission yang penting.
