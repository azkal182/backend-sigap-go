# Phase 5 – Academic SKS (Scores & Completion)

Tracking tasks and progress for Phase 5 implementation.

## Scope Recap
- Extend the academic module so each student's SKS exams and FAN completion status are persisted and queryable.
- Allow admins/academic staff to record SKS exam results (score, pass/fail, examiner, date) and view historical outcomes per student.
- Derive FAN completion by checking whether all SKS inside a FAN are passed for a student; expose this via read-only endpoints.
- Protect every create/update endpoint with the appropriate SKS permissions and ensure audit logs capture data changes.

## Data Model Targets
- `student_sks_results`
  - Fields: `id`, `student_id`, `sks_id`, `score`, `is_passed`, `exam_date`, `examiner_id`, timestamps.
  - Business rules:
    - `sks_id` must exist in `sks_definitions` and belong to an active FAN.
    - `examiner_id` must refer to an active teacher (optional but validated if provided).
    - `score` vs `kkm` determines `is_passed` when not explicitly supplied.
- `fan_completion_status` (optional but recommended cache for quick lookups)
  - Fields: `id`, `student_id`, `fan_id`, `is_completed`, `completed_at`, timestamps.
  - Automatically recomputed once all SKS in the FAN are passed.

## API & Use Case Targets
- `POST /api/sks-results` – create SKS exam result for a student.
- `PUT /api/sks-results/:id` – update score, pass/fail, or examiner.
- `GET /api/students/:id/sks-results` – list SKS outcomes for a student (with pagination/filter by FAN/SKS optional).
- `GET /api/students/:id/fans` – show FAN completion statuses (including progress metadata).
- Optional helper: `GET /api/sks-results?fan_id=...` for admin overview.
- Permissions: introduce `sks_results:read|create|update`; gate handlers using existing middleware per docs/adding_features.md.

## TODO
- [x] Confirm detailed requirements with stakeholders (score validation rules, retake policy, whether `fan_completion_status` is required or derived on the fly).
- [x] Design/extend domain entities (`StudentSKSResult`, optional `FANCompletionStatus`) under `internal/domain/entity`.
- [x] Add repository interfaces in `internal/domain/repository` plus GORM implementations in `internal/infrastructure/repository` for the new tables.
- [x] Create/extend DTOs in `internal/application/dto` for request/response payloads (create/update SKS results, list responses, FAN completion view).
- [x] Implement `StudentSKSResultUseCase` (or extend existing SKS use case) handling create/update/list logic, score-to-pass calculation, and FAN completion updates.
- [x] Add audit logging inside use cases for create/update operations (resource `sks_result`, action `sks_result:create|update`).
- [x] Build HTTP handlers + routes in `internal/interfaces/http/handler` and register them in the router with permission guards and pagination/query validation.
- [x] Register migrations for `student_sks_results` and optional `fan_completion_status` tables; update `cmd/migrate` + ensure down migrations available.
- [x] Update seeders (`cmd/seed`) to add new permissions (`sks_results:*`) and assign to relevant roles.
- [x] Extend documentation (README, PHASE_PLANNING.md) summarizing the new endpoints, tables, and permissions.
- [x] Write unit tests for the new use cases (success/error flows, fan completion logic) and HTTP integration tests covering CRUD/list endpoints.

### Recent Progress Snapshot
- Phase 5 planning document created (this file) capturing scope, data model, endpoints, and TODO checklist.
