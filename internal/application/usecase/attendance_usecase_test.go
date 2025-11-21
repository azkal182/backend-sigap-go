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

type fakeLeavePermitProvider struct {
	permit *entity.LeavePermit
	err    error
}

func (f fakeLeavePermitProvider) GetActivePermitForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error) {
	return f.permit, f.err
}

type fakeHealthStatusProvider struct {
	status *entity.HealthStatus
	err    error
}

func (f fakeHealthStatusProvider) GetActiveHealthStatusForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error) {
	return f.status, f.err
}

func TestAttendanceUseCase_OpenSessions(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	teacherRepo := new(mocks.TeacherAttendanceRepositoryMock)
	classScheduleRepo := new(mocks.ClassScheduleRepositoryMock)
	uc := usecase.NewAttendanceUseCase(sessionRepo, studentRepo, teacherRepo, classScheduleRepo, nil, nil, auditLoggerStub{})

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
	uc := usecase.NewAttendanceUseCase(sessionRepo, studentRepo, teacherRepo, classScheduleRepo, nil, nil, auditLoggerStub{})

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
	uc := usecase.NewAttendanceUseCase(sessionRepo, new(mocks.StudentAttendanceRepositoryMock), new(mocks.TeacherAttendanceRepositoryMock), new(mocks.ClassScheduleRepositoryMock), nil, nil, auditLoggerStub{})

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
	uc := usecase.NewAttendanceUseCase(sessionRepo, new(mocks.StudentAttendanceRepositoryMock), new(mocks.TeacherAttendanceRepositoryMock), new(mocks.ClassScheduleRepositoryMock), nil, nil, auditLoggerStub{})

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

func TestAttendanceUseCase_SubmitStudentAttendance_HealthOverride(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	uc := usecase.NewAttendanceUseCase(
		sessionRepo,
		studentRepo,
		new(mocks.TeacherAttendanceRepositoryMock),
		new(mocks.ClassScheduleRepositoryMock),
		fakeLeavePermitProvider{},
		fakeHealthStatusProvider{status: &entity.HealthStatus{ID: uuid.New()}},
		auditLoggerStub{},
	)

	sessionID := uuidFromString("11111111-1111-1111-1111-111111111111")
	studentID := uuidFromString("22222222-2222-2222-2222-222222222222")
	sessionRepo.On("GetByID", mock.Anything, sessionID).
		Return(&entity.AttendanceSession{ID: sessionID, Status: entity.AttendanceSessionStatusOpen, Date: time.Now()}, nil)
	studentRepo.
		On("BulkUpsert", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			records := args.Get(1).([]*entity.StudentAttendance)
			assert.Equal(t, entity.StudentAttendanceSick, records[0].Status)
		}).
		Return(nil)

	err := uc.SubmitStudentAttendance(context.Background(), sessionID, dto.SubmitStudentAttendanceRequest{
		Records: []dto.StudentAttendanceRecord{{StudentID: studentID.String(), Status: string(entity.StudentAttendancePresent)}},
	})

	assert.NoError(t, err)
	studentRepo.AssertExpectations(t)
}

func TestAttendanceUseCase_SubmitStudentAttendance_LeaveOverride(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	uc := usecase.NewAttendanceUseCase(
		sessionRepo,
		studentRepo,
		new(mocks.TeacherAttendanceRepositoryMock),
		new(mocks.ClassScheduleRepositoryMock),
		fakeLeavePermitProvider{permit: &entity.LeavePermit{ID: uuid.New()}},
		fakeHealthStatusProvider{},
		auditLoggerStub{},
	)

	sessionID := uuidFromString("33333333-3333-3333-3333-333333333333")
	studentID := uuidFromString("44444444-4444-4444-4444-444444444444")
	sessionRepo.On("GetByID", mock.Anything, sessionID).
		Return(&entity.AttendanceSession{ID: sessionID, Status: entity.AttendanceSessionStatusOpen, Date: time.Now()}, nil)
	studentRepo.
		On("BulkUpsert", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			records := args.Get(1).([]*entity.StudentAttendance)
			assert.Equal(t, entity.StudentAttendancePermit, records[0].Status)
		}).
		Return(nil)

	err := uc.SubmitStudentAttendance(context.Background(), sessionID, dto.SubmitStudentAttendanceRequest{
		Records: []dto.StudentAttendanceRecord{{StudentID: studentID.String(), Status: string(entity.StudentAttendancePresent)}},
	})

	assert.NoError(t, err)
	studentRepo.AssertExpectations(t)
}

func TestAttendanceUseCase_SubmitStudentAttendance_HealthPriorityOverLeave(t *testing.T) {
	sessionRepo := new(mocks.AttendanceSessionRepositoryMock)
	studentRepo := new(mocks.StudentAttendanceRepositoryMock)
	uc := usecase.NewAttendanceUseCase(
		sessionRepo,
		studentRepo,
		new(mocks.TeacherAttendanceRepositoryMock),
		new(mocks.ClassScheduleRepositoryMock),
		fakeLeavePermitProvider{permit: &entity.LeavePermit{ID: uuid.New()}},
		fakeHealthStatusProvider{status: &entity.HealthStatus{ID: uuid.New()}},
		auditLoggerStub{},
	)

	sessionID := uuidFromString("55555555-5555-5555-5555-555555555555")
	studentID := uuidFromString("66666666-6666-6666-6666-666666666666")
	sessionRepo.On("GetByID", mock.Anything, sessionID).
		Return(&entity.AttendanceSession{ID: sessionID, Status: entity.AttendanceSessionStatusOpen, Date: time.Now()}, nil)
	studentRepo.
		On("BulkUpsert", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			records := args.Get(1).([]*entity.StudentAttendance)
			assert.Equal(t, entity.StudentAttendanceSick, records[0].Status)
		}).
		Return(nil)

	err := uc.SubmitStudentAttendance(context.Background(), sessionID, dto.SubmitStudentAttendanceRequest{
		Records: []dto.StudentAttendanceRecord{{StudentID: studentID.String(), Status: string(entity.StudentAttendancePresent)}},
	})

	assert.NoError(t, err)
	studentRepo.AssertExpectations(t)
}
