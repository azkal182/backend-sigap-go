# Phase 3 – FAN & Class Structure

Tracking tasks and progress for Phase 3 implementation.

## Scope Recap
- Build layered academic structure: dormitory → FAN → classes.
- Tables:
  - `fans` (id, name, level, description)
  - `classes` (id, fan_id, name, capacity, is_active)
  - `student_class_enrollments` (id, class_id, student_id, enrolled_at, left_at)
  - `class_staff` (id, class_id, user_id, role (class_manager, homeroom_teacher))
- Endpoints:
  - FAN:
    - `POST /api/fans`
    - `GET /api/fans`
    - `PUT /api/fans/:id`
    - `DELETE /api/fans/:id`
  - Classes:
    - `POST /api/classes`
    - `GET /api/classes?fan_id=...`
    - `PUT /api/classes/:id`
    - `DELETE /api/classes/:id`
  - Enroll student:
    - `POST /api/classes/:id/students`
  - Assign staff:
    - `POST /api/classes/:id/staff`

## Completed
- [x] Define Phase 3 scope and documentation structure in this file.
- [x] Add domain entities for FAN & Class layer:
  - `Fan` (fans)
  - `Class` (classes)
  - `StudentClassEnrollment` (student_class_enrollments)
  - `ClassStaff` (class_staff)

## Todo
- [x] Design domain repository interfaces in `internal/domain/repository/`:
  - `FanRepository` (CRUD, list & filter, pagination).
  - `ClassRepository` (CRUD, list by fan, pagination).
  - Repositories for `StudentClassEnrollment` & `ClassStaff` (enroll/leave, assign staff, listing).
- [x] Add versioned migration for Phase 3 in `internal/infrastructure/database/migrations.go`:
  - Create tables `fans`, `classes`, `student_class_enrollments`, `class_staff` in `Up`.
  - Drop tables in correct reverse order in `Down`.
- [x] Implement GORM repositories in `internal/infrastructure/repository/` for all new entities.
- [x] Add DTOs in `internal/application/dto/`:
  - FAN: create/update/list/response DTOs.
  - Class: create/update/list/response DTOs.
  - Special requests for enroll student and assign staff.
- [x] Implement `FanUseCase` and `ClassUseCase` in `internal/application/usecase/` including:
  - CRUD FAN & Class.
  - List classes by FAN.
  - Enroll student to class and track enrollment history.
  - Assign staff (class_manager/homeroom_teacher) to class.
- [x] Add HTTP handlers in `internal/interfaces/http/handler/` and wire routes in `internal/interfaces/http/router/router.go`:
  - `/api/fans` and `/api/classes` endpoints per Phase 3 scope.
  - `/api/classes/:id/students` for enrollment.
  - `/api/classes/:id/staff` for staff assignment.
- [x] Extend permissions and seeding in `cmd/seed/main.go`:
  - Add `fans:*` and `classes:*` permissions (read/create/update/delete).
  - Assign to appropriate roles (e.g. admin, super_admin, academic roles from Phase 0).
- [x] Add unit tests for new use cases and repositories.
- [x] Extend HTTP integration tests to cover core Phase 3 flows:
  - Create & list FAN.
  - Create & list classes under FAN.
  - Enroll student into class.
  - Assign staff to class.
- [ ] Update `README.md` to document new FAN & Class endpoints and examples.
