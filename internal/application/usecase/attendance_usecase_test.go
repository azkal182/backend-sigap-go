package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

type auditLoggerStub struct{}

func (auditLoggerStub) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestAttendanceUseCase_OpenSessions(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	teacherRepo := new(mocks.TeacherAttendanceRepositoryMock)
	classScheduleRepo := new(mocks.ClassScheduleRepositoryMock)
	uc := usecase.NewAttendanceUseCase(sessionRepo, studentRepo, teacherRepo, classScheduleRepo, auditLoggerStub{})

	scheduleID := uuidFromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	teacherID := uuidFromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	start := time.Now()

	classScheduleRepo.
		On("GetByID", mock.Anything, scheduleID).
		Return(&entity.ClassSchedule{ID: scheduleID, TeacherID: teacherID, StartTime: &start}, nil)

	sessionRepo.
		On("GetOpenByScheduleAndDate", mock.Anything, scheduleID, mock.AnythingOfType("time.Time")).
		Return(nil, gorm.ErrRecordNotFound)
	sessionRepo.
		On("Create", mock.Anything, mock.AnythingOfType("*entity.AttendanceSession")).
		Return(nil)

	err := uc.OpenSessions(context.Background(), dto.OpenAttendanceSessionRequest{
		ClassScheduleIDs: []string{scheduleID.String()},
		Date:             start.Format("2006-01-02"),
	})

	assert.NoError(t, err)
	classScheduleRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestAttendanceUseCase_SubmitStudentAttendance(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	teacherRepo := new(mocks.TeacherAttendanceRepositoryMock)
	classScheduleRepo := new(mocks.ClassScheduleRepositoryMock)
	uc := usecase.NewAttendanceUseCase(sessionRepo, studentRepo, teacherRepo, classScheduleRepo, auditLoggerStub{})

	sessionID := uuidFromString("cccccccc-cccc-cccc-cccc-cccccccccccc")
	studentID := uuidFromString("dddddddd-dddd-dddd-dddd-dddddddddddd")

	sessionRepo.On("GetByID", mock.Anything, sessionID).
		Return(&entity.AttendanceSession{ID: sessionID, Status: entity.AttendanceSessionStatusOpen}, nil)
	studentRepo.On("BulkUpsert", mock.Anything, mock.Anything).
		Return(nil)

	err := uc.SubmitStudentAttendance(context.Background(), sessionID, dto.SubmitStudentAttendanceRequest{
		Records: []dto.StudentAttendanceRecord{{StudentID: studentID.String(), Status: "present"}},
	})

	assert.NoError(t, err)
	studentRepo.AssertExpectations(t)
}

func TestAttendanceUseCase_SubmitStudentAttendance_Locked(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	uc := usecase.NewAttendanceUseCase(sessionRepo, new(mocks.StudentAttendanceRepositoryMock), new(mocks.TeacherAttendanceRepositoryMock), new(mocks.ClassScheduleRepositoryMock), auditLoggerStub{})

	sessionID := uuidFromString("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	sessionRepo.On("GetByID", mock.Anything, sessionID).
		Return(&entity.AttendanceSession{ID: sessionID, Status: entity.AttendanceSessionStatusLocked}, nil)

	err := uc.SubmitStudentAttendance(context.Background(), sessionID, dto.SubmitStudentAttendanceRequest{
		Records: []dto.StudentAttendanceRecord{{StudentID: uuidFromString("ffffffff-ffff-ffff-ffff-ffffffffffff").String(), Status: "present"}},
	})

	assert.ErrorIs(t, err, domainErrors.ErrAttendanceAlreadyLocked)
}

func TestAttendanceUseCase_LockSessions(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	uc := usecase.NewAttendanceUseCase(sessionRepo, new(mocks.StudentAttendanceRepositoryMock), new(mocks.TeacherAttendanceRepositoryMock), new(mocks.ClassScheduleRepositoryMock), auditLoggerStub{})

	sessionRepo.On("LockSessionsByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(nil)

	err := uc.LockSessions(context.Background(), dto.LockAttendanceRequest{Date: "2025-11-20"})

	assert.NoError(t, err)
	sessionRepo.AssertExpectations(t)
}

func uuidFromString(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		panic(err)
	}
	return parsed
}
