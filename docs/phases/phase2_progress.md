# Phase 2 – Student Management (Kependudukan)

Tracking tasks and progress for Phase 2 implementation.

## Scope Recap
- Build core `students` table/entity with fields: `id`, `student_number`, `full_name`, `birth_date`, `gender`, `parent_name`, `status` (`active`, `inactive`, `leave`, `graduated`), timestamps.
- Track historical dormitory assignments via `student_dormitory_history` (student_id, dormitory_id, start/end date).
- Deliver endpoints:
  - `POST /api/students`
  - `GET /api/students`
  - `GET /api/students/:id`
  - `PUT /api/students/:id`
  - `PATCH /api/students/:id/status`
  - `POST /api/students/:id/mutate-dormitory`

## Todo
- [x] Design domain entities + repository interfaces for Student & StudentDormitoryHistory.
- [x] Add migrations / database schema changes for new tables.
- [x] Implement infrastructure repositories (GORM) for students + history queries.
- [x] Create DTOs & validations (create/update/status/mutate requests and responses).
- [x] Implement Student use case (CRUD, status update, dormitory mutation with history logging).
- [x] Build HTTP handler & routes, including status patch + mutate endpoints, with proper permissions/guard integration.
- [x] Add unit tests for use case logic and handler/service layers.
- [x] Extend integration tests (HTTP) covering happy-path student creation, status change, and dorm mutation.
- [x] Update README (features + endpoints) after implementation.
