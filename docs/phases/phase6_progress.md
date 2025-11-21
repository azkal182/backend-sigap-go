# Phase 6 – Attendance System (Santri & Teachers)

Tracking tasks and progress for Phase 6 implementation.

## Scope Recap
- Implement an attendance engine for class schedules that automatically opens sessions, records teacher/student presence, and enforces a lock at 23:59.
- Integrate attendance flow with existing leave-permit and health-status modules (Phase 7) while ensuring permissions and audit logging remain consistent with the patterns in [README.md](../README.md) and [ARCHITECTURE.md](../ARCHITECTURE.md).
- Follow the layered approach defined in [docs/adding_features.md](../docs/adding_features.md) for new entities, repositories, use cases, handlers, and routes.

## Data Model Targets
- `attendance_sessions`
  - `id`, `class_schedule_id`, `date`, `start_time`, `end_time`, `teacher_id`, `status` (`open`, `submitted`, `locked`), `locked_at`, timestamps.
  - Business rules: each class schedule can have at most one open session per day; session auto-created from schedule.
- `student_attendances`
  - `id`, `attendance_session_id`, `student_id`, `status` (`present`, `absent`, `permit`, `sick`), `note`, timestamps.
- `teacher_attendances`
  - `id`, `attendance_session_id`, `teacher_id`, `status`, timestamps.

## API & Use Case Targets
- `POST /api/attendance-sessions/open` – open session(s) automatically based on class schedules.
- `POST /api/attendance-sessions/:id/students` – submit/update student attendance in bulk.
- `POST /api/attendance-sessions/:id/teacher` – mark teacher attendance automatically/manual override.
- `POST /api/attendance-sessions/lock-day` – lock all sessions for a specific date (cron-friendly endpoint).
- `GET /api/attendance-sessions` + `GET /api/attendance-sessions/:id` – list/detail sessions with filters (class, dormitory, teacher, date).
- All mutations require dedicated permissions (e.g., `attendance_sessions:create`, `attendance_sessions:update`, `attendance_sessions:lock`).

- `cmd/attendance_lock` CLI added for cron locking (assign `attendance_sessions:lock` via `attendance_cron` role).

## TODO
- [ ] Confirm detailed requirements with stakeholders (bulk import constraints, cron locking policy, integration with leave/health data from Phase 7).
- [x] Design domain entities (`AttendanceSession`, `StudentAttendance`, `TeacherAttendance`) under `internal/domain/entity` plus related errors.
- [x] Add repository interfaces in `internal/domain/repository` and implement GORM adapters in `internal/infrastructure/repository`.
- [x] Create DTOs in `internal/application/dto` for session creation, student attendance submission, teacher attendance, and listing responses.
- [x] Implement use cases for opening sessions, submitting attendance, locking sessions, and querying history (with audit logging in each mutation).
- [x] Ensure background-friendly helpers (for cron lock) follow patterns in `docs/adding_features.md` (e.g., CLI/cron entry in `cmd/`).
- [x] Build HTTP handlers + routes in `internal/interfaces/http/handler` and register via router with permission guards consistent with README.
- [x] Register new migrations for attendance tables; update migration docs and `cmd/migrate` references.
- [x] Update seeders to include attendance permissions (read/create/update/lock) and assign to appropriate roles (admin, academic_sks, attendance_cron).
- [ ] Extend documentation (README, PHASE_PLANNING.md) summarizing the new endpoints, tables, cron flow, and permissions.
- [ ] Add unit tests (use case level) and HTTP integration tests covering session open/submit/lock flows.

### Recent Progress Snapshot
- Planning document created (this file) outlining Phase 6 scope, data model, endpoints, and TODO checklist.
