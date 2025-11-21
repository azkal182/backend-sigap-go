# Phase 7 – Security (Leave Permit) & Health (Sick Status)

Tracking tasks and progress for Phase 7 implementation.

## Scope Recap
- Integrate security leave permits and health/UKS sick-status workflows into the attendance system so status auto-propagates to daily attendance (per `PHASE_PLANNING.md`).
- Follow the layering + permission rules described in [README.md](../../README.md), [ARCHITECTURE.md](../../ARCHITECTURE.md), and the feature checklist in [docs/adding_features.md](../adding_features.md).
- Ensure audit logging, permissions, and cron/CLI flows stay consistent with the existing attendance engine from Phase 6.

## Data Model Targets
- `leave_permits`
  - `id`, `student_id`, `type` (`home_leave`, `official_duty`), `start_date`, `end_date`, `reason`, `status` (`pending`, `approved`, `rejected`, `completed`), timestamps, `approved_by`.
  - Rule: overlapping permits should be prevented; attendance auto-marks `permit` within the active window.
- `health_statuses`
  - `id`, `student_id`, `diagnosis`, `notes`, `start_date`, `end_date`, `status` (`active`, `revoked`), timestamps, `created_by`, `revoked_by`.
  - Rule: students flagged as `active` health cases auto-mark `sick`; only health officers can revoke.

## API & Use Case Targets
- Leave Permit (Security)
  - `POST /api/leave-permits` – create permit (pending by default).
  - `PUT /api/leave-permits/:id/approve` & `/reject` & `/complete` – workflow transitions (guarded by permissions/roles).
  - `GET /api/leave-permits?student_id=...&status=...` – list/filter.
- Health Status (UKS)
  - `POST /api/health-statuses` – log a new sick status (auto-flags attendance as `sick`).
  - `PUT /api/health-statuses/:id/revoke` – mark recovered (health officer only).
  - `GET /api/health-statuses?student_id=...&status=...` – list/filter.
- Attendance integration hooks (use case level): whenever Phase 6 attendance use cases fetch sessions, they should consult leave/health windows to enforce status overrides.

## Permissions & Roles
- New permissions to seed via `cmd/seed`:
  - `leave_permits:read|create|approve|complete`.
  - `health_statuses:read|create|revoke`.
- Role guidance:
  - `security_officer` → leave permit create/approve/complete.
  - `health_officer` → health status create/revoke.
  - `admin` inherits all; `academic_sks` keeps read-only visibility for daily ops.

## Progress Log
- ✅ Domain layer artifacts landed: `leave_permits` & `health_statuses` entities plus related domain errors and repository contracts.
- ✅ Infrastructure completed: versioned migration `014_create_leave_health_tables` and GORM repositories with overlap/active queries.
- ✅ Application DTOs + LeavePermit/HealthStatus use cases implemented with audit logging, validation, and attendance helper hooks.
- ✅ Interface wired: dedicated HTTP handlers, routes, and permission guards registered in the router.
- ✅ Seeder refreshed so new leave/health permissions exist and admin/super_admin roles automatically inherit them.
- ⏳ Attendance use case integration & regression tests still pending.

## Implementation Guidelines (per docs/adding_features.md)
1. **Domain Layer**
   - Add entities + validation helpers under `internal/domain/entity`.
   - Extend repository interfaces in `internal/domain/repository` for leave/health CRUD + overlap checks.
   - Add domain errors (`ErrLeavePermitConflict`, `ErrHealthStatusActive`, etc.).
2. **Application Layer**
   - DTOs for create/update/list in `internal/application/dto` with validation tags.
   - Use cases orchestrating workflows, status transitions, attendance hooks, and audit logging.
3. **Infrastructure Layer**
   - Migrations (`internal/infrastructure/database/migrations.go`) for new tables/indexes.
   - Repository adapters under `internal/infrastructure/repository` (include overlap queries and status transitions).
   - Seeder updates for permissions + new roles.
4. **Interface Layer**
   - HTTP handlers (`internal/interfaces/http/handler`) + router wiring with permission guards.
   - Update integration tests in `internal/interfaces/http/integration_test.go` to cover CRUD + attendance impact.

## Integration Points with Attendance
- Extend `AttendanceUseCase` to query leave/health repositories during `OpenSessions` and `SubmitStudentAttendance` to auto-set status.
- CLI lock job should respect derived statuses (no change needed but add regression tests).
- Ensure audit logger captures leave/health approvals/revocations.

## Testing Targets
- Unit tests:
  - LeavePermitUseCase & HealthStatusUseCase (create, approve/revoke, overlap validation).
  - AttendanceUseCase updates verifying permit/sick overrides.
- Integration tests:
  - HTTP flows for leave/health endpoints (auth + permission guards).
  - Attendance endpoint verifying auto `permit`/`sick` when overlapping windows exist.

## TODO Checklist
- [x] Confirm detailed workflow (approval roles, notifications) with stakeholders.
- [x] Design domain entities + errors for leave permits & health statuses.
- [x] Define repository interfaces + filters; implement GORM adapters with overlap checks.
- [x] Create DTOs for create/list/update/revoke flows.
- [x] Implement LeavePermitUseCase & HealthStatusUseCase (audit logging, status transitions).
- [ ] Update AttendanceUseCase to consult leave/health repositories before saving attendance records.
- [x] Add HTTP handlers, routes, permission guards, and request/response documentation in README.
- [x] Write migrations + update `cmd/migrate` docs.
- [x] Extend seeders for new permissions/roles; update CLI docs if needed.
- [x] Add unit tests (use case level) and HTTP integration tests for leave/health + attendance overrides.

## Next Steps Snapshot
- Planning document created (this file) capturing scope, data model, endpoints, permissions, and TODO list for Phase 7.
- Upcoming work: finalize entity/repository design, migrations, and permission seeding before coding use cases.
