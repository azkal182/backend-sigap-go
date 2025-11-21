package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_ domainRepo.AttendanceSessionRepository = (*attendanceSessionRepository)(nil)
	_ domainRepo.StudentAttendanceRepository = (*studentAttendanceRepository)(nil)
	_ domainRepo.TeacherAttendanceRepository = (*teacherAttendanceRepository)(nil)
)

type attendanceSessionRepository struct {
	db *gorm.DB
}

type studentAttendanceRepository struct {
	db *gorm.DB
}

type teacherAttendanceRepository struct {
	db *gorm.DB
}

// Constructor helpers
func NewAttendanceSessionRepository() domainRepo.AttendanceSessionRepository {
	return &attendanceSessionRepository{db: database.DB}
}

func NewStudentAttendanceRepository() domainRepo.StudentAttendanceRepository {
	return &studentAttendanceRepository{db: database.DB}
}

func NewTeacherAttendanceRepository() domainRepo.TeacherAttendanceRepository {
	return &teacherAttendanceRepository{db: database.DB}
}

// Attendance session implementation
func (r *attendanceSessionRepository) Create(ctx context.Context, session *entity.AttendanceSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *attendanceSessionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.AttendanceSessionStatus, lockedAt *time.Time) error {
	updates := map[string]interface{}{"status": status}
	if lockedAt != nil {
		updates["locked_at"] = lockedAt
	}
	return r.db.WithContext(ctx).
		Model(&entity.AttendanceSession{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *attendanceSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AttendanceSession, error) {
	var session entity.AttendanceSession
	if err := r.db.WithContext(ctx).
		Preload("StudentAttendances").
		Preload("TeacherAttendances").
		Where("id = ?", id).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *attendanceSessionRepository) GetOpenByScheduleAndDate(ctx context.Context, scheduleID uuid.UUID, date time.Time) (*entity.AttendanceSession, error) {
	var session entity.AttendanceSession
	if err := r.db.WithContext(ctx).
		Where("class_schedule_id = ? AND date = ? AND status = ?",
			scheduleID, date, entity.AttendanceSessionStatusOpen).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *attendanceSessionRepository) List(ctx context.Context, filter domainRepo.AttendanceSessionFilter) ([]*entity.AttendanceSession, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.AttendanceSession{}).
		Preload("StudentAttendances").
		Preload("TeacherAttendances")

	if filter.ClassScheduleID != nil {
		query = query.Where("class_schedule_id = ?", filter.ClassScheduleID)
	}
	if filter.TeacherID != nil {
		query = query.Where("teacher_id = ?", filter.TeacherID)
	}
	if filter.Date != nil {
		query = query.Where("date = ?", filter.Date)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}

	limit := filter.Limit
	offset := filter.Offset
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sessions []*entity.AttendanceSession
	if err := query.Order("date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *attendanceSessionRepository) LockSessionsByDate(ctx context.Context, date time.Time) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.AttendanceSession{}).
		Where("date = ? AND status <> ?", date, entity.AttendanceSessionStatusLocked).
		Updates(map[string]interface{}{
			"status":    entity.AttendanceSessionStatusLocked,
			"locked_at": now,
		}).Error
}

// Student attendance implementation
func (r *studentAttendanceRepository) BulkUpsert(ctx context.Context, attendances []*entity.StudentAttendance) error {
	if len(attendances) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "attendance_session_id"}, {Name: "student_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "note", "updated_at"}),
		}).
		Create(&attendances).Error
}

func (r *studentAttendanceRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]*entity.StudentAttendance, error) {
	var records []*entity.StudentAttendance
	if err := r.db.WithContext(ctx).
		Where("attendance_session_id = ?", sessionID).
		Order("student_id ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// Teacher attendance implementation
func (r *teacherAttendanceRepository) Upsert(ctx context.Context, attendance *entity.TeacherAttendance) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "attendance_session_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"teacher_id", "status", "updated_at"}),
		}).
		Create(attendance).Error
}

func (r *teacherAttendanceRepository) GetBySession(ctx context.Context, sessionID uuid.UUID) (*entity.TeacherAttendance, error) {
	var record entity.TeacherAttendance
	if err := r.db.WithContext(ctx).
		Where("attendance_session_id = ?", sessionID).
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}
