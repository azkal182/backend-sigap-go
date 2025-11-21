Berikut penulisan ulang seluruh **phase** dengan penamaan **tabel, field, dan endpoint semuanya dalam bahasa Inggris**, tanpa mengubah struktur domain pesantren.

---

# **PHASE 0 – Role & Permission Alignment**

### **Goals**

* Menyelaraskan domain pesantren dengan sistem RBAC yang sudah ada.
* Menentukan role dan permission dalam bahasa Inggris.

### **New Roles**

* `admin`
* `central_secretary`
* `dormitory_secretary`
* `class_manager`
* `academic_sks`
* `security_officer`
* `health_officer`
* `teacher`

### **Permissions (examples)**

* `students:read`, `students:create`, `students:update`, `students:mutate`
* `fans:read`, `fans:create`, `fans:update`
* `classes:read`, `classes:create`, `classes:update`
* `sks:read`, `sks:create`, `sks:update`
* `attendance:read`, `attendance:update`
* `leave_permits:create`, `leave_permits:read`
* `health_status:create`, `health_status:read`

### **Implementation**

* Seed roles & permissions.
* Set route guard using `RequirePermission`.

---

# **PHASE 1 – Dormitory Management Base (Asrama)**

### **Goals**

* Memakai entity dormitory untuk manajemen asrama.

### **Tables**

* `dormitories`

  * fields: `id`, `name`, `gender`, `level`, `code`, timestamps.
* `user_dormitories` (assignment staf ke asrama)

### **Endpoints**

* `GET /api/dormitories`
* `POST /api/dormitories`
* `PUT /api/dormitories/:id`
* `DELETE /api/dormitories/:id`
* Assign staff:

  * `POST /api/dormitories/:id/users`

---

# **PHASE 2 – Student Management (Kependudukan)**

### **Goals**

* Membuat modul inti santri dengan mutasi asrama & status aktif.

### **Tables**

* `students`

  * `id`, `student_number`, `full_name`, `birth_date`, `gender`, `parent_name`,
    `status` (`active`, `inactive`, `leave`, `graduated`), timestamps.
* `student_dormitory_history`

  * `id`, `student_id`, `dormitory_id`, `start_date`, `end_date`.

### **Endpoints**

* `POST /api/students`
* `GET /api/students`
* `GET /api/students/:id`
* `PUT /api/students/:id`
* Update status:

  * `PATCH /api/students/:id/status`
* Mutasi asrama:

  * `POST /api/students/:id/mutate-dormitory`

---

# **PHASE 3 – FAN & Class Structure**

### **Goals**

* Membuat struktur akademik berlapis: dormitory → FAN → classes.

### **Tables**

* `fans`

  * `id`, `name`, `level`, `description`.
* `classes`

  * `id`, `fan_id`, `name`, `capacity`, `is_active`.
* `student_class_enrollments`

  * `id`, `class_id`, `student_id`, `enrolled_at`, `left_at`.
* `class_staff`

  * `id`, `class_id`, `user_id`, `role` (`class_manager`, `homeroom_teacher`).

### **Endpoints**

* FAN:

  * `POST /api/fans`
  * `GET /api/fans`
  * `PUT /api/fans/:id`
  * `DELETE /api/fans/:id`
* Classes:

  * `POST /api/classes`
  * `GET /api/classes?fan_id=...`
  * `PUT /api/classes/:id`
  * `DELETE /api/classes/:id`
* Enroll student:

  * `POST /api/classes/:id/students`
* Assign staff:

  * `POST /api/classes/:id/staff`

---

# **PHASE 4 – Class Schedule & SKS Schedule**

### **Goals**

* Menyediakan jadwal pelajaran & jadwal ujian SKS, berbasis slot waktu dormitory.
* Memastikan integrasi dengan modul teacher (aktif), FAN, dan schedule slot.
* Semua mutasi jadwal diaudit dan diproteksi permission khusus.

### **Tables**

* `subjects` (optional)

  * `id`, `name`, `description`.
* `class_schedules`

  * `id`, `class_id`, `subject_id`, `teacher_id`,
    `day_of_week`, `start_time`, `end_time`, `slot_id`, `dormitory_id`, `is_active`.
* `sks_definitions`

  * `id`, `fan_id`, `code`, `name`, `kkm`, `description`.
* `sks_exam_schedules`

  * `id`, `sks_id`, `exam_date`, `exam_time`, `location`, `examiner_id`.

*Depends on*: `schedule_slots` (slot_number per dormitory dengan validasi overlap) dari tahap sebelumnya.

### **Business Rules**

1. Class schedule dapat menggunakan slot (otomatis ambil start/end) atau jam manual.
2. Slot harus berasal dari dormitory yang sama dengan kelas; slot harus aktif.
3. Teacher wajib status aktif sebelum dijadwalkan.
4. SKS definition optional subject; kode unik per sistem.
5. SKS exam schedule memerlukan tanggal (YYYY-MM-DD) dan waktu (HH:MM) valid serta examiner teacher aktif (opsional).
6. Semua operasi create/update/delete memicu audit log dan dilindungi permission `class_schedules:*`, `sks_definitions:*`, `sks_exams:*`.

### **Endpoints**

* Class schedule:

  * `GET /api/class-schedules?class_id=...&teacher_id=...&dormitory_id=...`
  * `GET /api/class-schedules/:id`
  * `POST /api/class-schedules`
  * `PUT /api/class-schedules/:id`
  * `DELETE /api/class-schedules/:id`
* SKS:

  * `GET /api/sks?fan_id=...`
  * `GET /api/sks/:id`
  * `POST /api/sks`
  * `PUT /api/sks/:id`
  * `DELETE /api/sks/:id`
* SKS exam schedule:

  * `GET /api/sks-exams?sks_id=...`
  * `GET /api/sks-exams/:id`
  * `POST /api/sks-exams`
  * `PUT /api/sks-exams/:id`
  * `DELETE /api/sks-exams/:id`

### **Status**

✅ Implemented (Nov 2025): domain entities, repositories, use cases, HTTP handlers, router wiring, permissions + seed, integration tests, dan dokumentasi (README + phase progress). Phase 5 adalah fokus berikutnya.

---

# **PHASE 5 – Academic SKS (Scores & Completion)**

### **Goals**

* Menyimpan nilai SKS, menentukan kelulusan setiap SKS dan FAN.

### **Tables**

* `student_sks_results`

  * `id`, `student_id`, `sks_id`, `score`, `is_passed`, `exam_date`, `examiner_id`.
* `fan_completion_status` (optional)

  * `id`, `student_id`, `fan_id`, `is_completed`, `completed_at`.

### **Endpoints**

* `POST /api/sks-results`
* `PUT /api/sks-results/:id`
* Capaian SKS santri:

  * `GET /api/students/:id/sks-results`
* Status kelulusan FAN:

  * `GET /api/students/:id/fans`

---

# **PHASE 6 – Attendance System (Santri & Teachers)**

### **Goals**

* Membuat sistem absensi otomatis, terkunci 23:59, terintegrasi izin & sakit.

### **Tables**

* `attendance_sessions`

  * `id`, `class_schedule_id`,
    `date`, `start_time`, `end_time`,
    `teacher_id`, `status` (`open`, `submitted`, `locked`), `locked_at`.
* `student_attendances`

  * `id`, `attendance_session_id`, `student_id`,
    `status` (`present`, `absent`, `permit`, `sick`), `note`.
* `teacher_attendances`

  * `id`, `attendance_session_id`, `teacher_id`, `status`.

### **Core Rules**

1. Auto-detect schedule → create/open session.
2. Absensi santri bulk update.
3. Teacher attendance auto-created.
4. Lock at 23:59 (cron job).
5. Only same-day edit allowed (except admin override).
6. Audit every change.

### **Endpoints**

* Start/open session automatically:

  * `POST /api/attendance-sessions/open`
* Submit/update attendance:

  * `POST /api/attendance-sessions/:id/students`
* Lock session (cron):

  * `POST /api/attendance-sessions/lock-day`
* Teacher attendance log (auto):

  * `POST /api/attendance-sessions/:id/teacher`

### **Permissions & Roles**

* New permissions: `attendance_sessions:read`, `attendance_sessions:create`, `attendance_sessions:update`, `attendance_sessions:lock`.
* `admin` role now includes all attendance permissions; `academic_sks` role spreads read/create/update for daily ops.
* `attendance_cron` role (service account) has read + lock, meant for `cmd/attendance_lock` nightly automation.
* CLI helper: `go run cmd/attendance_lock/main.go -date YYYY-MM-DD` (defaults to today) to enforce lock job.

---

# **PHASE 7 – Security (Permit) & Health (Sick Status)**

### **Goals**

* Integrasi izin keluar (security) dan status sakit (UKS) dengan absensi.

### **Tables**

* `leave_permits`

  * `id`, `student_id`, `type` (`home_leave`, `official_duty`),
    `start_date`, `end_date`, `status` (`approved`, `rejected`, `completed`).
* `health_statuses`

  * `id`, `student_id`, `diagnosis`, `start_date`, `end_date`, `status` (`active`, `revoked`).

### **Rules**

* Jika berada pada rentang `leave_permits` → status absensi = `permit`.
* Jika rentang `health_statuses` → status absensi = `sick`.
* Hanya health_officer yang bisa revoke sakit.
* Pengajar dilarang mengubah status izin/sakit.

### **Endpoints**

* Security:

  * `POST /api/leave-permits`
  * `GET /api/leave-permits?student_id=...`
* Health:

  * `POST /api/health-statuses`
  * `PUT /api/health-statuses/:id/revoke`

---

# **PHASE 8 – Reports & Monitoring**

### **Goals**

* Membuat laporan lengkap lintas modul.

### **Report Endpoints**

* Student attendance report:

  * `GET /api/reports/attendance/students?date=...&dormitory_id=...`
* Teacher attendance report:

  * `GET /api/reports/attendance/teachers?...`
* Leave permit report:

  * `GET /api/reports/leave-permits?...`
* Health status report:

  * `GET /api/reports/health-statuses?...`
* SKS report:

  * `GET /api/reports/sks?fan_id=...`
* FAN mutation report:

  * `GET /api/reports/mutations?student_id=...`

### **Permissions & Status**

* Permissions introduced: `reports:attendance:read`, `reports:security:read`, `reports:health:read`, `reports:academic:read`.
* `central_secretary`, `admin`, dan `super_admin` memiliki seluruh izin laporan; role lain dapat diberi subset sesuai kebutuhan.
* Status (Nov 2025): backend wiring (DTOs, repository, use case, handler, router, seed) selesai; dokumentasi & testing sedang berjalan.

---

# **PHASE 9 – Stability: Testing, Audit, Performance**

### **Goals**

* Memastikan sistem siap produksi.

### **Tasks**

* Unit test untuk usecase core: student, class, attendance, SKS.
* Integration test end-to-end.
* Audit logging semua perubahan sensitif.

