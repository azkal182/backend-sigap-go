# Phase 8 – Reports & Monitoring

Tracking tasks and design notes for the cross-module reporting initiative.

## Scope Recap
- Provide consolidated dashboards for attendance, leave permits, health statuses, SKS, and FAN/student mutations.
- Ensure every report follows clean architecture boundaries per [ARCHITECTURE.md](../../ARCHITECTURE.md) and reuses DTO/handler guidance from [docs/adding_features.md](../adding_features.md).
- Support filtering (date range, dormitory, FAN, role) while enforcing RBAC permissions.
- Produce export-friendly JSON payloads (future CSV export TBD) without duplicating core business logic.

## Data & Aggregation Targets
- **Attendance Sessions & Student Attendances**
  - Join `attendance_sessions` + `student_attendances` to summarize present/absent/permit/sick counts per dorm/class/day.
- **Teacher Attendance**
  - Reuse `teacher_attendances` to show punctuality/lock compliance metrics.
- **Leave Permits**
  - Aggregate active vs completed permits, grouped by type and dormitory.
- **Health Statuses**
  - Track active sick cases, revoke trends, and overlap with attendance submissions.
- **SKS Progress**
  - Pull from `student_sks_results` + optional `fan_completion_status` to highlight pass rates per FAN.
- **Mutation History**
  - Use `student_dormitory_history` + class enrollments to generate movement reports.

## API & Use Case Targets
- `GET /api/reports/attendance/students?date=YYYY-MM-DD&dormitory_id=...`
- `GET /api/reports/attendance/teachers?date=YYYY-MM-DD&slot_id=...`
- `GET /api/reports/leave-permits?status=active&type=...&date=...`
- `GET /api/reports/health-statuses?status=active&dormitory_id=...`
- `GET /api/reports/sks?fan_id=...&is_passed=...`
- `GET /api/reports/mutations?student_id=...&fan_id=...`

Each endpoint should:
1. Validate filters (use shared DTOs/pagination helpers).
2. Call dedicated report use cases that orchestrate repository queries (no ad-hoc SQL in handlers).
3. Return normalized response objects with metadata (e.g., totals, period label, grouping keys).

## Permissions & Security
- Introduce granular permissions, e.g., `reports:attendance:read`, `reports:security:read`, `reports:health:read`, `reports:academic:read`.
- Guard every report route using middleware per [docs/adding_features.md](../adding_features.md#step-7-routes) and keep actor context for audit logging.

## Stakeholder Requirements Summary
- **Security (Leave Permit):** Need daily/weekly summaries per dormitory with `type`, `status`, and `student` counts, plus ability to export active permits for gate officers.
- **UKS (Health):** Track active vs revoked cases, highlight overlaps with attendance per date, and flag students exceeding 3 consecutive sick days.
- **Academic (Attendance & SKS):** Require aggregate present/absent/permit/sick stats per class/FAN, teacher punctuality, and SKS pass-rate dashboards filtered by FAN and exam date.
- **Central Office:** Mutation history report must show dorm/class transitions with date ranges for audits; all reports should expose metadata (generated_at, filters) to facilitate CSV export later.

## Implementation Checklist
- [x] Finalize reporting requirements with Security, UKS, and Academic stakeholders (filters, grouping, export needs).
- [x] Define DTOs + validators for each report request/response in `internal/application/dto/report_dto.go`.
- [x] Add report-specific repository methods (or read-only query services) per clean architecture guidelines in [ARCHITECTURE.md](../../ARCHITECTURE.md).
- [x] Implement `ReportUseCase` modules (attendance, leave, health, SKS, mutation) reusing existing repositories.
- [x] Create HTTP handlers and routes under `internal/interfaces/http/handler/report_handler.go` and register them in router with new permissions.
- [x] Extend seeders to insert the new `reports:*` permissions and assign to roles (`admin`, `central_secretary`, etc.).
- [x] Document endpoints, filters, and sample payloads in `README.md` (Reports section) plus update `PHASE_PLANNING.md`.
- [x] Add unit tests for report use cases (filter validation, aggregation scenarios) and HTTP integration smoke tests covering each endpoint.
- [ ] Evaluate future CSV/Excel export hooks (optional) and note decision in this file.

## Next Steps Snapshot
- Phase 8 planning doc (this file) created to capture scope, data sources, endpoints, and TODO list before implementation begins.
- Pending stakeholder workshop to lock report filters and permission mapping.
