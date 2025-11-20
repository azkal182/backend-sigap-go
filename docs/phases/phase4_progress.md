# Phase 4 – Class Schedule & SKS Schedule

Tracking tasks and progress for Phase 4 implementation.

## Scope Recap
- **Teacher prerequisite (before schedule work):** add Teacher CRUD as a first-class feature linked to users so teachers can authenticate for future attendance submissions.
- Build scheduling layer for regular classes and SKS exams, including shared time slots per dormitory (e.g., Slot 1 = 08:00–09:00 for all classes in a dorm).
- Tables:
  - `subjects` (optional helper lookup: id, name, description)
  - `class_schedules` (id, class_id, subject_id, teacher_id, day_of_week, start_time, end_time)
  - `sks_definitions` (id, fan_id, code, name, kkm, description)
  - `sks_exam_schedules` (id, sks_id, exam_date, exam_time, location, examiner_id)
- Endpoints:
  - `POST /api/class-schedules`
  - `GET /api/class-schedules?class_id=...`
  - `POST /api/sks`
  - `GET /api/sks?fan_id=...`
  - `POST /api/sks-exams`
  - `GET /api/sks-exams?sks_id=...`

### Teacher Specification
- **Entity fields**: `id (uuid)`, `user_id (uuid)`, `teacher_code`, `full_name`, `gender`, `phone`, `email`, `specialization`, `employment_status`, `joined_at`, `is_active`, timestamps.
- **User linkage**: each teacher owns exactly one `users` row. Creating a teacher auto-generates a user with role `teacher` (default password strategy: random password emailed/logged for admin) or links to existing user if provided.
- **Business rules**:
  - `teacher_code` unique per system.
  - Cannot delete teacher if still assigned to class schedules; use `is_active` toggle.
  - Updating teacher info should sync selected fields (name/email) to linked user.
- **Endpoints**:
  - `POST /api/teachers` (auto create user + role assignment).
  - `GET /api/teachers` with filters (`is_active`, `keyword`).
  - `GET /api/teachers/:id`.
  - `PUT /api/teachers/:id` (update + optional password reset for linked user).
  - `DELETE /api/teachers/:id` or `PATCH` to deactivate/reactivate.
- **Permissions**: `teachers:read|create|update|delete`; assign to roles (`admin`, `academic_sks`).
- **Audit**: log create/update/delete events including user linkage metadata.

### Schedule Slot Specification
- **Purpose**: standardize daily schedule so each dormitory follows shared time slots (Slot 1 = 08:00–09:00, Slot 2 = 09:15–10:15, etc.).
- **Entity fields (`schedule_slots`)**: `id (uuid)`, `dormitory_id`, `slot_number`, `name`, `start_time`, `end_time`, `description`, `is_active`, timestamps.
- **Relationships**: `class_schedules` references slots via `slot_id`; slots optionally nullable to allow ad-hoc times.
- **Business rules**:
  - Slot numbers unique per dormitory.
  - Start/end times cannot overlap within same dormitory.
  - Class schedule referencing a slot must belong to a class under the same dormitory (via class → fan → dorm mapping).
- **Endpoints**:
  - `POST /api/schedule-slots` (per dorm). 
  - `GET /api/schedule-slots?dormitory_id=...`.
  - `PUT /api/schedule-slots/:id`.
  - `DELETE /api/schedule-slots/:id` (soft delete / deactivate).
- **Usage in schedules**:
  - When creating class schedules, choose either slot reference or free-form time; if slot provided, derive start/end from slot.
  - Provide helper endpoint `GET /api/class-schedules/slots?dormitory_id=...` for UI.

## Todo
- [x] Capture detailed Teacher requirements (fields, unique identifier, activation workflow, default password strategy) and validation rules.
- [x] Define Teacher domain model + relationship to `User` (1-1) including auto-user creation/association on teacher onboarding.
- [x] Register migrations for `teachers` table + foreign key to `users`, and ensure cascading/unique constraints (teacher per user).
- [x] Add teacher repository interfaces, GORM implementations, DTOs, and use cases (CRUD, list, soft-delete/activate) with proper validations.
- [x] Build teacher HTTP handlers + routes (`/api/teachers`) with permissions (`teachers:read|create|update|delete`) and transactional user creation + role assignment.
- [x] Update seeders to introduce default teacher role/permissions, and document the module in README + architecture notes.
- [x] Add unit tests (domain/use case) and integration tests (HTTP) covering teacher creation, auto-user linkage, and updates/deletes.
- [x] Ensure `Class`/schedule components reference teacher IDs (FK) once teacher module exists.
- [x] Update seeders/migrations/docs to include the `teachers`table and default teacher roles before schedule work begins.
- [ ] Finalize detailed requirements & validation rules for schedule entities (time windows, day enums, overlaps) plus slot definitions per dormitory.
- [x] Design & migrate `schedule_slots` (slot number, name, start/end time, dormitory_id).
- [x] Link upcoming class schedules to slots (enforce dormitory alignment).
- [x] Update repositories/use cases/handlers so class schedules can reuse slots and enforce dormitory-slot constraints.
- [x] Design domain entities for `Subject`, `ClassSchedule`, `SKSDefinition`, `SKSExamSchedule`.
- [x] Add repository interfaces in `internal/domain/repository/` for schedules and SKS data.
- [x] Register migration `010_create_schedule_slots` in `internal/infrastructure/database/migrations.go`.
- [x] Register new migrations for `class_schedules`, `subjects`, `sks_definitions`, `sks_exam_schedules`.
- [x] Implement GORM repositories in `internal/infrastructure/repository/` (including conflict/overlap checks where needed).
- [x] Define DTOs/validators in `internal/application/dto/` for schedule & SKS operations.
- [x] Implement use cases (create/list schedules, create/list SKS definitions, create/list SKS exams) with business rules.
- [x] Build HTTP handlers + routes, protecting endpoints with appropriate permissions (`class_schedules:*`, `sks:*`).
- [x] Add endpoint + handler for updating class schedules (`PUT /api/class-schedules/:id`) including validation and overlap checks.
- [x] Extend seed data with `schedule_slots:*` permissions and map them to admin role.
- [x] Add new permissions/roles for class schedules & SKS if/when required.
- [x] Add unit tests for schedule slot use case (create/overlap/list/update/delete) with mocks.
- [x] Add unit tests for class schedule & SKS use cases/repositories.
- [x] Add integration test covering `/api/schedule-slots` CRUD flow.
- [x] Extend HTTP integration tests for class schedules & SKS endpoints.
- [x] Update README & relevant docs after implementation.

### Recent Progress Snapshot
- ✅ Schedule slot domain entity, repository interface + GORM implementation, use case, DTOs, and handler completed (CRUD + overlap validation).
- ✅ Versioned migration `010_create_schedule_slots` created and wired into migrator.
- ✅ `/api/schedule-slots` endpoints (create/list/get/update/delete) registered with permission guards and integration-tested end-to-end.
- ✅ Seeder updated dengan permission `schedule_slots:*`, `class_schedules:*`, `sks_definitions:*`, dan `sks_exams:*` untuk mengamankan endpoint baru.
- ✅ README + docs diperbarui untuk merangkum fitur jadwal kelas & SKS beserta endpoint dan permission-nya.
