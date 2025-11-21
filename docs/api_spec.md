# SIGAP API Specification (Work in Progress)

This document tracks endpoint specifications, sample requests/responses, and phased documentation tasks.

## 1. Base Information
- **Base URL (Prod):** `https://<your-domain>/api`
- **Authentication:** Bearer token via `Authorization: Bearer <token>` header.
- **Content Type:** `application/json` unless stated otherwise.
- **Pagination Pattern:** `?page=<n>&limit=<m>` or `?offset=<n>&limit=<m>` depending on handler (see individual sections).

## 2. Documentation Tasks (Phase Plan)
| Section | Status | Notes |
| --- | --- | --- |
| Auth | ✅ Done (login/register) |
| Current User & Permissions | ✅ `/me`, `/permissions` |
| Users & Roles | ✅ Core CRUD + role assignment |
| Dormitories | ✅ CRUD + user assignment |
| Students & SKS | ✅ CRUD + SKS results/fan status |
| FAN & Classes | ✅ FAN hierarchy + classes |
| Teachers | ✅ CRUD |
| Attendance | ✅ Sessions + student/teacher logs |
| SKS Definitions & Exams | ✅ Definition + exam schedules |
| Leave & Health | ✅ Security/UKS flows |
| Reports | ✅ Filtered GET endpoints |
| Locations (Public) | ✅ Provinces/Regencies/etc |
| Audit Logs | ✅ Read-only list |

### Task Breakdown
- [x] Phase 1 – Auth endpoints (register/login/refresh) with sample payloads.
- [x] Phase 2 – Core directory endpoints (Users, Roles, Permissions, Dormitories).
- [x] Phase 3 – Academic structure (Students, FAN, Classes, SKS definitions/exams).
- [x] Phase 4 – Attendance sessions + student/teacher attendance flows.
- [x] Phase 5 – Security/UKS (Leave permits, Health statuses) plus audit logs.
- [x] Phase 6 – Reports (attendance, leave, health, SKS, mutation) including filters.
- [x] Phase 7 – Misc endpoints (locations, schedule slots, class schedules).

> Update this checklist as sections are completed. Each phase should document request/response schema, query params, and error cases.

## 3. Endpoint Catalog (from `internal/interfaces/http/router/router.go`)
- **Auth:** `/auth/register`, `/auth/login`, `/auth/refresh`
- **Users:** `/users`, `/users/:id`, `/users/:id/roles`
- **Roles:** `/roles`, `/roles/:id`, `/roles/:id/permissions`
- **Permissions:** `/permissions`
- **Dormitories:** `/dormitories`, nested user management
- **Students:** `/students`, nested SKS result actions
- **Teachers:** `/teachers`
- **FAN / Classes:** `/fans`, `/classes`
- **Schedule Slots & Schedules:** `/schedule-slots`, `/class-schedules`
- **SKS Definitions/Exams:** `/sks`, `/sks-exams`
- **Attendance:** `/attendance` (sessions scoped under other routes)
- **Leave & Health:** `/leave-permits`, `/health-statuses`
- **Reports:** `/reports/...`

Each section below will be fleshed out per phase.

## 4. Auth Endpoints (Phase 1 ✅)
### 4.1 Register
- **Method/URL:** `POST /api/auth/register`
- **Request**
```json
{
  "username": "john.doe",
  "password": "password123",
  "name": "John Doe"
}
```

**Create User – Request**
```json
{
  "username": "central.sec",
  "password": "supersecret",
  "name": "Central Secretary",
  "role_ids": ["central-secretary-role-uuid"],
  "dormitory_ids": []
}
```
- **Response 201**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<refresh>",
    "expires_at": "2025-11-21T00:00:00Z",
    "user": {
      "id": "uuid",
      "username": "john.doe",
      "name": "John Doe",
      "roles": ["user"]
    }
  }
}
```

### 4.2 Login
- **Method/URL:** `POST /api/auth/login`
- **Request**
```json
{
  "username": "admin",
  "password": "admin123"
}
```
- **Response 200**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<refresh>",
    "expires_at": "2025-11-21T01:00:00Z"
  }
}
```

### 4.3 Refresh Token
- **Method/URL:** `POST /api/auth/refresh`
- **Request**
```json
{
  "refresh_token": "<refresh>"
}
```
- **Response 200** – same payload as login.

## 5. Dormitories (Phase 2b ✅)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/dormitories` | `dorm:read` | List dormitories accessible to caller (supports pagination). |
| GET | `/api/dormitories/:id` | `dorm:read` + dorm access | Retrieve dorm detail. |
| POST | `/api/dormitories` | `dorm:create` | Create dormitory. |
| PUT | `/api/dormitories/:id` | `dorm:update` + dorm access | Update metadata. |
| DELETE | `/api/dormitories/:id` | `dorm:delete` + dorm access | Soft-delete/disable. |
| POST | `/api/dormitories/:id/users` | `dorm:update` + dorm access | Assign staff/user to dorm. |
| DELETE | `/api/dormitories/:id/users/:user_id` | `dorm:update` + dorm access | Remove assignment. |

**Create Dormitory – Request**
```json
{
  "name": "Dorm A",
  "gender": "male",
  "level": "senior",
  "code": "DORMA",
  "description": "Main dorm",
  "is_active": true
}
```

**Assign Dormitory User – Request**
```bash
POST /api/dormitories/dorm-uuid/users
{
  "user_id": "user-uuid"
}
```

## 6. Current User & Permissions

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/me` | Authenticated | Returns current user profile (roles, permissions, dormitories). |
| GET | `/api/permissions` | `role:read` | Paginated list of permissions (used for role editors). |

**Sample `/me` Response**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "admin",
    "name": "Admin User",
    "roles": ["admin"],
    "permissions": ["user:read", "reports:attendance:read"]
  }
}
```

## 7. Users & Roles (Phase 2 ✅)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/users` | `user:read` | List users with pagination. |
| GET | `/api/users/:id` | `user:read` | Retrieve user detail (roles/dorms). |
| POST | `/api/users` | `user:create` | Create user. |
| PUT | `/api/users/:id` | `user:update` | Update user profile. |
| DELETE | `/api/users/:id` | `user:delete` | Soft delete user. |
| POST | `/api/users/:id/roles` | `user:update` | Assign role(s) to user. |
| DELETE | `/api/users/:id/roles/:role_id` | `user:update` | Remove a role. |

**List Users – Request**
```
GET /api/users?limit=20&offset=0
Authorization: Bearer <token>
```

**Response 200**
```json
{
  "success": true,
  "message": "Users retrieved",
  "data": {
    "total": 2,
    "items": [
      {
        "id": "uuid",
        "username": "admin",
        "name": "Admin User",
        "roles": ["admin"],
        "is_active": true,
        "created_at": "2025-11-20T01:00:00Z"
      }
    ]
  }
}
```

**Assign Role – Request**
```bash
POST /api/users/9b3.../roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "role_id": "2b1..."
}
```

### Roles & Permissions

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/roles` | `role:read` | Paginated role list. |
| POST | `/api/roles` | `role:create` | Create role (non-protected). |
| PUT | `/api/roles/:id` | `role:update` | Update role metadata. |
| DELETE | `/api/roles/:id` | `role:delete` | Delete role (unless protected). |
| POST | `/api/roles/:id/permissions` | `role:update` | Assign permissions. |

**Create Role – Request**
```json
{
  "name": "Central Secretary",
  "slug": "central_secretary",
  "permission_ids": ["perm-id-1", "perm-id-2"]
}
```

**Response 201**
```json
{
  "success": true,
  "message": "Role created successfully",
  "data": {
    "id": "uuid",
    "name": "Central Secretary",
    "permissions": ["reports:attendance:read", "reports:academic:read"]
  }
}
```

## 8. FAN & Classes (Phase 3a ✅)

### FAN Endpoints

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/fans` | `fans:read` | List FAN structures. |
| GET | `/api/fans/:id` | `fans:read` | Detail (levels, classes). |
| POST | `/api/fans` | `fans:create` | Create FAN hierarchy root. |
| PUT | `/api/fans/:id` | `fans:update` | Update metadata. |
| DELETE | `/api/fans/:id` | `fans:delete` | Remove FAN (if unused). |

**Create FAN – Request**
```json
{
  "name": "FAN 2026",
  "description": "Batch 2026",
  "is_active": true
}
```

**Create Class – Request**
```json
{
  "fan_id": "fan-uuid",
  "name": "Class Alpha",
  "level": "grade_10",
  "homeroom_teacher_id": "teacher-uuid"
}
```

### Class Endpoints

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/classes?fan_id=...` | `classes:read` | List classes optionally filtered by FAN. |
| GET | `/api/classes/:id` | `classes:read` | Class detail with enrollments/staff. |
| POST | `/api/classes` | `classes:create` | Create class under FAN. |
| PUT | `/api/classes/:id` | `classes:update` | Update class info. |
| DELETE | `/api/classes/:id` | `classes:delete` | Delete class. |
| POST | `/api/classes/:id/students` | `classes:update` | Enroll student. |
| POST | `/api/classes/:id/staff` | `classes:update` | Assign staff. |

**Enroll Student – Request**
```json
{
  "student_id": "student-uuid"
}
```

**Assign Staff – Request**
```json
{
  "teacher_id": "teacher-uuid",
  "role": "advisor"
}
```

## 9. Students & SKS (Phase 3 ✅)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/students` | `student:read` | List students (supports `status`, `dormitory_id`, pagination). |
| GET | `/api/students/:id` | `student:read` | Student detail including dorm history. |
| POST | `/api/students` | `student:create` | Create student profile. |
| PUT | `/api/students/:id` | `student:update` | Update student attributes. |
| PATCH | `/api/students/:id/status` | `student:update` | Change lifecycle status. |
| POST | `/api/students/:id/mutate-dormitory` | `student:update` | Move student between dorms. |
| POST | `/api/students/:id/sks-results` | `student_sks_results:create` | Record SKS exam result. |
| PUT | `/api/students/:id/sks-results/:result_id` | `student_sks_results:update` | Update SKS result. |
| GET | `/api/students/:id/sks-results` | `student_sks_results:read` | List SKS results (filter by FAN/SKS). |
| GET | `/api/students/:id/fans` | `student_sks_results:read` | FAN completion status summary. |

**Create Student – Request**
```json
{
  "student_number": "STD001",
  "full_name": "Integration Student",
  "birth_date": "2010-01-01",
  "gender": "male",
  "parent_name": "Parent Doe"
}
```

**Response 201**
```json
{
  "success": true,
  "message": "Student created successfully",
  "data": {
    "id": "uuid",
    "student_number": "STD001",
    "status": "active",
    "created_at": "2025-11-21T03:00:00Z"
  }
}
```

**Record SKS Result – Request**
```json
{
  "fan_id": "fan-uuid",
  "sks_id": "sks-uuid",
  "score": 84,
  "exam_date": "2025-11-15",
  "examiner_id": "teacher-uuid"
}
```

**Response 201** – Contains `is_passed`, metadata, and audit info.

## 10. Teachers (Phase 3b ✅)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/teachers` | `teachers:read` | List teachers. |
| GET | `/api/teachers/:id` | `teachers:read` | Teacher detail. |
| POST | `/api/teachers` | `teachers:create` | Create teacher record. |
| PUT | `/api/teachers/:id` | `teachers:update` | Update teacher info. |
| DELETE | `/api/teachers/:id` | `teachers:delete` | Deactivate teacher. |

**Create Teacher – Request**
```json
{
  "user_id": "user-uuid",
  "subject_ids": ["sub-1"],
  "status": "active"
}
```

## 11. Attendance (Phase 4 ✅)

### Schedule Slots

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/schedule-slots?dormitory_id=...` | `schedule_slots:read` | List slots by dorm. |
| POST | `/api/schedule-slots` | `schedule_slots:create` | Create slot with conflict checks. |
| PUT | `/api/schedule-slots/:id` | `schedule_slots:update` | Update slot window/meta. |
| DELETE | `/api/schedule-slots/:id` | `schedule_slots:delete` | Delete/deactivate slot. |

### Class Schedules

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/class-schedules?class_id=...` | `class_schedules:read` | List schedule entries. |
| GET | `/api/class-schedules/:id` | `class_schedules:read` | Detail. |
| POST | `/api/class-schedules` | `class_schedules:create` | Create schedule referencing slot/time. |
| PUT | `/api/class-schedules/:id` | `class_schedules:update` | Update schedule info. |
| DELETE | `/api/class-schedules/:id` | `class_schedules:delete` | Delete schedule. |

### Attendance Sessions & Logs

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| POST | `/api/attendance/sessions`* | `attendance_sessions:create` | (handled via Attendance use case) create session per schedule. |
| GET | `/api/attendance/sessions` | `attendance_sessions:read` | Filter by schedule/teacher/date/status. |
| PUT | `/api/attendance/sessions/:id/lock` | `attendance_sessions:lock` | Lock session for given date. |
| POST | `/api/attendance/sessions/:id/student-attendance` | `attendance_sessions:update` | Bulk upsert student attendance statuses. |
| POST | `/api/attendance/sessions/:id/teacher-attendance` | `attendance_sessions:update` | Upsert teacher attendance. |

**List Sessions – Request**
```
GET /api/attendance/sessions?teacher_id=<uuid>&date=2025-11-21&limit=10
```

**Bulk Student Attendance – Request**
```json
{
  "attendances": [
    { "student_id": "stu-1", "status": "present" },
    { "student_id": "stu-2", "status": "permit", "note": "family" }
  ]
}
```

**Create Schedule Slot – Request**
```json
{
  "dormitory_id": "dorm-uuid",
  "name": "Subuh Prep",
  "start_time": "05:00",
  "end_time": "06:00",
  "days": ["monday", "tuesday", "wednesday"],
  "is_active": true
}
```

**Create Class Schedule – Request**
```json
{
  "class_id": "class-uuid",
  "teacher_id": "teacher-uuid",
  "subject_id": "subject-uuid",
  "schedule_slot_id": "slot-uuid",
  "start_time": "2025-11-22T07:30:00Z",
  "end_time": "2025-11-22T09:00:00Z"
}
```

\*Actual routes are grouped under `/students/:id/sks-results` or `/attendance` handlers; reference router for exact nesting when integrating.

## 12. SKS Definitions & Exam Schedules (Phase 3 ✅)

### SKS Definitions

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/sks` | `sks_definitions:read` | List SKS definitions (supports `fan_id`). |
| GET | `/api/sks/:id` | `sks_definitions:read` | Detail. |
| POST | `/api/sks` | `sks_definitions:create` | Create definition referencing FAN & subject. |
| PUT | `/api/sks/:id` | `sks_definitions:update` | Update metadata/status. |
| DELETE | `/api/sks/:id` | `sks_definitions:delete` | Delete definition. |

### SKS Exam Schedules

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/sks-exams?sks_id=...` | `sks_exams:read` | List schedules per definition. |
| GET | `/api/sks-exams/:id` | `sks_exams:read` | Detail. |
| POST | `/api/sks-exams` | `sks_exams:create` | Create exam schedule. |
| PUT | `/api/sks-exams/:id` | `sks_exams:update` | Update date/time/examiner. |
| DELETE | `/api/sks-exams/:id` | `sks_exams:delete` | Delete schedule. |

## 13. Leave Permits & Health Statuses (Phase 5 ✅)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/leave-permits` | `leave_permits:read` | Filter by `student_id`, `status`, `type`, date range. |
| POST | `/api/leave-permits` | `leave_permits:create` | Create pending permit. |
| PUT | `/api/leave-permits/:id/approve` | `leave_permits:approve` | Approve permit. |
| PUT | `/api/leave-permits/:id/reject` | `leave_permits:approve` | Reject permit. |
| PUT | `/api/leave-permits/:id/complete` | `leave_permits:complete` | Mark as completed on return. |
| GET | `/api/health-statuses` | `health_statuses:read` | Filter by `student_id`, `status`, date range. |
| POST | `/api/health-statuses` | `health_statuses:create` | Create sick status. |
| PUT | `/api/health-statuses/:id/revoke` | `health_statuses:revoke` | Revoke status (sets `end_date`). |

**Create Leave Permit – Request**
```json
{
  "student_id": "student-uuid",
  "type": "home_leave",
  "reason": "Family matters",
  "start_date": "2025-11-22",
  "end_date": "2025-11-24",
  "created_by": "staff-uuid"
}
```

**List Leave Permits – Query Example**
```
GET /api/leave-permits?status=approved&dormitory_id=<uuid>&limit=20
```

**Health Status Create – Request**
```json
{
  "student_id": "student-uuid",
  "diagnosis": "Flu",
  "notes": "Needs rest",
  "start_date": "2025-11-21",
  "created_by": "uks-staff-uuid"
}
```

## 14. Reports (Phase 6 ✅)

| Endpoint | Permission | Required Params | Notes |
| --- | --- | --- | --- |
| `GET /api/reports/attendance/students` | `reports:attendance:read` | `date` (YYYY-MM-DD); optional `dormitory_id`, `class_id`, `fan_id`. | Returns aggregated attendance counts. |
| `GET /api/reports/attendance/teachers` | `reports:attendance:read` | `date`; optional `slot_id`, `teacher_id`. | Teacher punctuality summary. |
| `GET /api/reports/leave-permits` | `reports:security:read` | Optional `status`, `type`, `dormitory_id`, `date_range[start|end]`. | Aggregated leave totals. |
| `GET /api/reports/health-statuses` | `reports:health:read` | Optional `status`, `dormitory_id`, `date_range`. | Counts of active/revoked cases. |
| `GET /api/reports/sks` | `reports:academic:read` | Optional `fan_id`, `sks_id`, `is_passed`, `date_range`. | Pass/fail summary + averages. |
| `GET /api/reports/mutations` | `reports:academic:read` | Optional `student_id`, `fan_id`, `dormitory_id`, `date_range`. | Dorm/class mutation history. |

**Student Attendance Report – Request**
```
GET /api/reports/attendance/students?date=2025-11-21&dormitory_id=<uuid>
Authorization: Bearer <token>
```

**Response 200 (trimmed)**
```json
{
  "success": true,
  "data": {
    "generated_at": "2025-11-21T23:59:00Z",
    "filters": {
      "date": "2025-11-21",
      "dormitory_id": "<uuid>"
    },
    "rows": [
      {
        "dormitory_id": "<uuid>",
        "class_id": "<uuid>",
        "fan_id": "<uuid>",
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

## 15. Locations (Public) & Audit Logs (Phase 7 ✅)

### Location Endpoints (no auth)

| Method | URL | Description |
| --- | --- | --- |
| GET | `/api/provinces` | List provinces. |
| GET | `/api/provinces/:id` | Province detail. |
| GET | `/api/regencies` | Filter by `province_id`. |
| GET | `/api/regencies/:id` | Regency detail. |
| GET | `/api/districts` | Filter by `regency_id`. |
| GET | `/api/districts/:id` | District detail. |
| GET | `/api/villages` | Filter by `district_id`. |
| GET | `/api/villages/:id` | Village detail. |

### Audit Logs (protected)

| Method | URL | Permission | Description |
| --- | --- | --- | --- |
| GET | `/api/audit-logs` | `audit:read` | Paginated audit trail filtered by actor/resource/timestamp. |

**Audit Logs – Request**
```
GET /api/audit-logs?page=1&limit=20&action=update
Authorization: Bearer <token>
```

## 16. Contribution Checklist
When updating this doc:
1. Mark the relevant phase checkbox.
2. Add endpoint descriptions with tables (method, URL, perms, request fields, response fields).
3. Include at least one sample request/response per endpoint.
4. Reference validation rules from DTOs when helpful.
