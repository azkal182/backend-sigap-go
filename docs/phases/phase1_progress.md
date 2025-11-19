# Phase 1 – Dormitory Management Base

Tracking tasks and status for Phase 1 implementation.

## Completed
- [x] Requirements review & gap analysis between existing dormitory module and Phase 1 scope
- [x] Update dormitory schema/entities & DTOs (gender, level, code, validation)
- [x] Extend use cases + handlers for new dormitory fields
- [x] Implement staff assignment & removal endpoints (`POST/DELETE /api/dormitories/:id/users`)
- [x] Update router permissions for new endpoints
- [x] Add unit tests covering new use case behavior

## In Progress / Next
- [x] Review database migrations/seed data to ensure new dormitory attributes are aligned
- [x] Validate dormitory-user assignment flow via integration tests (HTTP level)
