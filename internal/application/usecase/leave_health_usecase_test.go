package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

type leaveHealthAuditLogger struct{}

func (leaveHealthAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func ctxWithActor() context.Context {
	// Production code uses service.CtxKeyActorID (a string) as the context key, so we must
	// match it in tests even though staticcheck warns about string keys.
	//lint:ignore SA1029 we must share the exact key with production code to exercise requireActorID logic

	return context.WithValue(context.Background(), service.CtxKeyActorID, uuid.New())
}

func stringPtrLH(val string) *string {
	return &val
}

func newLeavePermitUseCase(t *testing.T) (*LeavePermitUseCase, *mocks.LeavePermitRepositoryMock, *mocks.MockStudentRepository) {
	t.Helper()
	leaveRepo := new(mocks.LeavePermitRepositoryMock)
	studentRepo := new(mocks.MockStudentRepository)
	return NewLeavePermitUseCase(leaveRepo, studentRepo, leaveHealthAuditLogger{}), leaveRepo, studentRepo
}

func newHealthStatusUseCase(t *testing.T) (*HealthStatusUseCase, *mocks.HealthStatusRepositoryMock, *mocks.MockStudentRepository) {
	t.Helper()
	healthRepo := new(mocks.HealthStatusRepositoryMock)
	studentRepo := new(mocks.MockStudentRepository)
	return NewHealthStatusUseCase(healthRepo, studentRepo, leaveHealthAuditLogger{}), healthRepo, studentRepo
}

// ----- Leave permit tests -----

func TestLeavePermitUseCase_CreateLeavePermit_Success(t *testing.T) {
	ctx := ctxWithActor()
	uc, leaveRepo, studentRepo := newLeavePermitUseCase(t)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	leaveRepo.On("HasOverlap", mock.Anything, studentID, mock.Anything, mock.Anything, (*uuid.UUID)(nil)).Return(false, nil)
	leaveRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := uc.CreateLeavePermit(ctx, dto.CreateLeavePermitRequest{
		StudentID: studentID.String(),
		Type:      string(entity.LeavePermitTypeHomeLeave),
		Reason:    "Family",
		StartDate: "2025-11-01",
		EndDate:   "2025-11-02",
	})

	assert.NoError(t, err)
	assert.Equal(t, string(entity.LeavePermitTypeHomeLeave), resp.Type)
	leaveRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestLeavePermitUseCase_CreateLeavePermit_Overlap(t *testing.T) {
	ctx := ctxWithActor()
	uc, leaveRepo, studentRepo := newLeavePermitUseCase(t)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	leaveRepo.On("HasOverlap", mock.Anything, studentID, mock.Anything, mock.Anything, (*uuid.UUID)(nil)).Return(true, nil)

	resp, err := uc.CreateLeavePermit(ctx, dto.CreateLeavePermitRequest{
		StudentID: studentID.String(),
		Type:      string(entity.LeavePermitTypeOfficialDuty),
		StartDate: "2025-11-01",
		EndDate:   "2025-11-03",
	})
	assert.ErrorIs(t, err, domainErrors.ErrLeavePermitConflict)
	assert.Nil(t, resp)
	leaveRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestLeavePermitUseCase_UpdateLeavePermitStatus(t *testing.T) {
	ctx := ctxWithActor()
	uc, leaveRepo, studentRepo := newLeavePermitUseCase(t)
	studentID := uuid.New()
	permitID := uuid.New()
	permit := &entity.LeavePermit{ID: permitID, Status: entity.LeavePermitStatusPending}

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil).Maybe()
	leaveRepo.On("GetByID", mock.Anything, permitID).Return(permit, nil)
	leaveRepo.On("Update", mock.Anything, permit).Return(nil)

	resp, err := uc.UpdateLeavePermitStatus(ctx, permitID, dto.UpdateLeavePermitStatusRequest{Status: string(entity.LeavePermitStatusApproved)})
	assert.NoError(t, err)
	assert.Equal(t, string(entity.LeavePermitStatusApproved), resp.Status)
	leaveRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestLeavePermitUseCase_UpdateLeavePermitStatus_InvalidTransition(t *testing.T) {
	ctx := ctxWithActor()
	uc, leaveRepo, studentRepo := newLeavePermitUseCase(t)
	studentID := uuid.New()
	permitID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil).Maybe()
	leaveRepo.On("GetByID", mock.Anything, permitID).Return(&entity.LeavePermit{ID: permitID, Status: entity.LeavePermitStatusPending}, nil)

	resp, err := uc.UpdateLeavePermitStatus(ctx, permitID, dto.UpdateLeavePermitStatusRequest{Status: string(entity.LeavePermitStatusCompleted)})
	assert.ErrorIs(t, err, domainErrors.ErrLeavePermitStatus)
	assert.Nil(t, resp)
	leaveRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestLeavePermitUseCase_ListLeavePermits(t *testing.T) {
	uc, leaveRepo, _ := newLeavePermitUseCase(t)
	studentID := uuid.New()
	status := string(entity.LeavePermitStatusApproved)
	permitType := string(entity.LeavePermitTypeHomeLeave)
	date := "2025-11-05"

	leaveRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.LeavePermitFilter) bool {
		return filter.StudentID != nil && *filter.StudentID == studentID &&
			filter.Status != nil && *filter.Status == entity.LeavePermitStatus(status) &&
			filter.Type != nil && *filter.Type == entity.LeavePermitType(permitType) &&
			filter.Date != nil
	})).Return([]*entity.LeavePermit{{ID: uuid.New()}}, int64(1), nil)

	resp, err := uc.ListLeavePermits(context.Background(), dto.ListLeavePermitsRequest{
		StudentID: stringPtrLH(studentID.String()),
		Status:    &status,
		Type:      &permitType,
		Date:      &date,
		Page:      1,
		PageSize:  10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	leaveRepo.AssertExpectations(t)
}

func TestLeavePermitUseCase_GetActivePermitForDate_NotFound(t *testing.T) {
	uc, leaveRepo, _ := newLeavePermitUseCase(t)
	leaveRepo.On("ActiveByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	permit, err := uc.GetActivePermitForDate(context.Background(), uuid.New(), time.Now())
	assert.NoError(t, err)
	assert.Nil(t, permit)
	leaveRepo.AssertExpectations(t)
}

// ----- Health status tests -----

func TestHealthStatusUseCase_CreateHealthStatus_Success(t *testing.T) {
	ctx := ctxWithActor()
	uc, healthRepo, studentRepo := newHealthStatusUseCase(t)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	healthRepo.On("ActiveByDate", mock.Anything, studentID, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	healthRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	endDate := "2025-11-25"
	resp, err := uc.CreateHealthStatus(ctx, dto.CreateHealthStatusRequest{
		StudentID: studentID.String(),
		Diagnosis: "Influenza",
		Notes:     "Needs rest",
		StartDate: "2025-11-21",
		EndDate:   &endDate,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Influenza", resp.Diagnosis)
	healthRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestHealthStatusUseCase_CreateHealthStatus_ActiveExists(t *testing.T) {
	ctx := ctxWithActor()
	uc, healthRepo, studentRepo := newHealthStatusUseCase(t)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	healthRepo.On("ActiveByDate", mock.Anything, studentID, mock.Anything).Return(&entity.HealthStatus{ID: uuid.New()}, nil)

	resp, err := uc.CreateHealthStatus(ctx, dto.CreateHealthStatusRequest{
		StudentID: studentID.String(),
		Diagnosis: "Flu",
		StartDate: "2025-11-21",
	})
	assert.ErrorIs(t, err, domainErrors.ErrHealthStatusActive)
	assert.Nil(t, resp)
	healthRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestHealthStatusUseCase_ListHealthStatuses(t *testing.T) {
	uc, healthRepo, _ := newHealthStatusUseCase(t)
	studentID := uuid.New()
	status := string(entity.HealthStatusStateActive)
	date := "2025-11-22"

	healthRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.HealthStatusFilter) bool {
		return filter.StudentID != nil && *filter.StudentID == studentID &&
			filter.Status != nil && *filter.Status == entity.HealthStatusState(status) &&
			filter.Date != nil
	})).Return([]*entity.HealthStatus{{ID: uuid.New()}}, int64(1), nil)

	resp, err := uc.ListHealthStatuses(context.Background(), dto.ListHealthStatusesRequest{
		StudentID: stringPtrLH(studentID.String()),
		Status:    &status,
		Date:      &date,
		Page:      1,
		PageSize:  5,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	healthRepo.AssertExpectations(t)
}

func TestHealthStatusUseCase_RevokeHealthStatus_Success(t *testing.T) {
	ctx := ctxWithActor()
	uc, healthRepo, studentRepo := newHealthStatusUseCase(t)
	studentID := uuid.New()
	statusID := uuid.New()
	status := &entity.HealthStatus{ID: statusID, Status: entity.HealthStatusStateActive}

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil).Maybe()
	healthRepo.On("GetByID", mock.Anything, statusID).Return(status, nil)
	healthRepo.On("Update", mock.Anything, status).Return(nil)

	resp, err := uc.RevokeHealthStatus(ctx, statusID, dto.RevokeHealthStatusRequest{Reason: "Recovered"})
	assert.NoError(t, err)
	assert.Equal(t, string(entity.HealthStatusStateRevoked), resp.Status)
	healthRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestHealthStatusUseCase_RevokeHealthStatus_NotActive(t *testing.T) {
	ctx := ctxWithActor()
	uc, healthRepo, studentRepo := newHealthStatusUseCase(t)
	studentID := uuid.New()
	statusID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil).Maybe()
	healthRepo.On("GetByID", mock.Anything, statusID).Return(&entity.HealthStatus{ID: statusID, Status: entity.HealthStatusStateRevoked}, nil)

	resp, err := uc.RevokeHealthStatus(ctx, statusID, dto.RevokeHealthStatusRequest{})
	assert.ErrorIs(t, err, domainErrors.ErrHealthStatusForbidden)
	assert.Nil(t, resp)
	healthRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
}

func TestHealthStatusUseCase_GetActiveHealthStatusForDate_NotFound(t *testing.T) {
	uc, healthRepo, _ := newHealthStatusUseCase(t)
	healthRepo.On("ActiveByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	status, err := uc.GetActiveHealthStatusForDate(context.Background(), uuid.New(), time.Now())
	assert.NoError(t, err)
	assert.Nil(t, status)
	healthRepo.AssertExpectations(t)
}
