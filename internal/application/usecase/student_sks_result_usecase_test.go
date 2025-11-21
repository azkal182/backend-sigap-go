package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

func TestStudentSKSResultUseCase_CreateStudentSKSResult(t *testing.T) {
	ctx := context.Background()
	studentID := uuid.New()
	sksID := uuid.New()
	fanID := uuid.New()

	studentRepo := new(mocks.MockStudentRepository)
	resultRepo := new(mocks.StudentSKSResultRepositoryMock)
	fanRepo := new(mocks.FanCompletionStatusRepositoryMock)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	teacherRepo := new(mocks.TeacherRepositoryMock)

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID, FanID: fanID, KKM: 75}, nil)
	sksRepo.On("CountByFan", mock.Anything, fanID).Return(int64(1), nil)
	resultRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	resultRepo.On("CountPassedByStudentFan", mock.Anything, studentID, fanID).Return(int64(1), nil)
	fanRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	uc := NewStudentSKSResultUseCase(resultRepo, fanRepo, studentRepo, sksRepo, teacherRepo, &noopAuditLogger{})
	resp, err := uc.CreateStudentSKSResult(ctx, dto.CreateStudentSKSResultRequest{
		StudentID: studentID.String(),
		SKSID:     sksID.String(),
		Score:     80,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, studentID.String(), resp.StudentID)
	assert.Equal(t, fanID.String(), resp.FanID)

	studentRepo.AssertExpectations(t)
	sksRepo.AssertExpectations(t)
	resultRepo.AssertExpectations(t)
	fanRepo.AssertExpectations(t)
}

func TestStudentSKSResultUseCase_UpdateStudentSKSResult(t *testing.T) {
	ctx := context.Background()
	resultID := uuid.New()
	studentID := uuid.New()
	sksID := uuid.New()
	fanID := uuid.New()

	studentRepo := new(mocks.MockStudentRepository)
	resultRepo := new(mocks.StudentSKSResultRepositoryMock)
	fanRepo := new(mocks.FanCompletionStatusRepositoryMock)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	teacherRepo := new(mocks.TeacherRepositoryMock)

	resultRepo.On("GetByID", mock.Anything, resultID).Return(&entity.StudentSKSResult{
		ID:        resultID,
		StudentID: studentID,
		SKSID:     sksID,
		Score:     70,
	}, nil)
	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID, FanID: fanID, KKM: 75}, nil)
	resultRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	sksRepo.On("CountByFan", mock.Anything, fanID).Return(int64(2), nil)
	resultRepo.On("CountPassedByStudentFan", mock.Anything, studentID, fanID).Return(int64(1), nil)
	fanRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	uc := NewStudentSKSResultUseCase(resultRepo, fanRepo, studentRepo, sksRepo, teacherRepo, &noopAuditLogger{})
	score := 60.0
	resp, err := uc.UpdateStudentSKSResult(ctx, resultID, dto.UpdateStudentSKSResultRequest{Score: &score})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, score, resp.Score)
	assert.False(t, resp.IsPassed)

	resultRepo.AssertExpectations(t)
	fanRepo.AssertExpectations(t)
	sksRepo.AssertExpectations(t)
}

func TestStudentSKSResultUseCase_ListStudentSKSResults(t *testing.T) {
	ctx := context.Background()
	studentID := uuid.New()
	sksID := uuid.New()
	fanID := uuid.New()

	studentRepo := new(mocks.MockStudentRepository)
	resultRepo := new(mocks.StudentSKSResultRepositoryMock)
	fanRepo := new(mocks.FanCompletionStatusRepositoryMock)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	teacherRepo := new(mocks.TeacherRepositoryMock)

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	resultRepo.On("ListByStudent", mock.Anything, studentID, uuid.Nil, 10, 0).Return([]*entity.StudentSKSResult{{
		ID:        uuid.New(),
		StudentID: studentID,
		SKSID:     sksID,
	}}, int64(1), nil)
	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID, FanID: fanID}, nil)

	uc := NewStudentSKSResultUseCase(resultRepo, fanRepo, studentRepo, sksRepo, teacherRepo, &noopAuditLogger{})
	resp, err := uc.ListStudentSKSResults(ctx, studentID, "", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, resp.Results, 1)
	assert.Equal(t, fanID.String(), resp.Results[0].FanID)
	resultRepo.AssertExpectations(t)
}

func TestStudentSKSResultUseCase_ListFanCompletionStatuses(t *testing.T) {
	ctx := context.Background()
	studentID := uuid.New()
	studentRepo := new(mocks.MockStudentRepository)
	resultRepo := new(mocks.StudentSKSResultRepositoryMock)
	fanRepo := new(mocks.FanCompletionStatusRepositoryMock)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	teacherRepo := new(mocks.TeacherRepositoryMock)

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	fanRepo.On("ListByStudent", mock.Anything, studentID).Return([]*entity.FanCompletionStatus{{
		FanID:       uuid.New(),
		IsCompleted: true,
		CompletedAt: nil,
	}}, nil)

	uc := NewStudentSKSResultUseCase(resultRepo, fanRepo, studentRepo, sksRepo, teacherRepo, &noopAuditLogger{})
	resp, err := uc.ListFanCompletionStatuses(ctx, studentID)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	studentRepo.AssertExpectations(t)
	fanRepo.AssertExpectations(t)
}

func TestStudentSKSResultUseCase_CreateStudentSKSResult_InvalidStudent(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	studentRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, domainErrors.ErrStudentNotFound)

	uc := NewStudentSKSResultUseCase(new(mocks.StudentSKSResultRepositoryMock), new(mocks.FanCompletionStatusRepositoryMock), studentRepo, new(mocks.SKSDefinitionRepositoryMock), new(mocks.TeacherRepositoryMock), &noopAuditLogger{})
	resp, err := uc.CreateStudentSKSResult(ctx, dto.CreateStudentSKSResultRequest{
		StudentID: uuid.New().String(),
		SKSID:     uuid.New().String(),
		Score:     80,
	})

	assert.ErrorIs(t, err, domainErrors.ErrStudentNotFound)
	assert.Nil(t, resp)
}

func TestStudentSKSResultUseCase_CreateStudentSKSResult_InvalidSKSID(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	studentRepo.On("GetByID", mock.Anything, mock.Anything).Return(&entity.Student{ID: uuid.New()}, nil)

	uc := NewStudentSKSResultUseCase(new(mocks.StudentSKSResultRepositoryMock), new(mocks.FanCompletionStatusRepositoryMock), studentRepo, new(mocks.SKSDefinitionRepositoryMock), new(mocks.TeacherRepositoryMock), &noopAuditLogger{})
	resp, err := uc.CreateStudentSKSResult(ctx, dto.CreateStudentSKSResultRequest{
		StudentID: uuid.New().String(),
		SKSID:     "invalid-uuid",
		Score:     75,
	})

	assert.ErrorIs(t, err, domainErrors.ErrBadRequest)
	assert.Nil(t, resp)
}

func TestStudentSKSResultUseCase_CreateStudentSKSResult_InvalidExaminer(t *testing.T) {
	ctx := context.Background()
	studentID := uuid.New()
	sksID := uuid.New()
	fanID := uuid.New()

	studentRepo := new(mocks.MockStudentRepository)
	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID, FanID: fanID, KKM: 70}, nil)
	teacherRepo := new(mocks.TeacherRepositoryMock)
	teacherRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, domainErrors.ErrTeacherNotFound)

	uc := NewStudentSKSResultUseCase(new(mocks.StudentSKSResultRepositoryMock), new(mocks.FanCompletionStatusRepositoryMock), studentRepo, sksRepo, teacherRepo, &noopAuditLogger{})
	resp, err := uc.CreateStudentSKSResult(ctx, dto.CreateStudentSKSResultRequest{
		StudentID:  studentID.String(),
		SKSID:      sksID.String(),
		Score:      80,
		ExaminerID: ptrString(uuid.New().String()),
	})

	assert.ErrorIs(t, err, domainErrors.ErrTeacherNotFound)
	assert.Nil(t, resp)
}

func TestStudentSKSResultUseCase_CreateStudentSKSResult_InvalidExamDate(t *testing.T) {
	ctx := context.Background()
	studentID := uuid.New()
	sksID := uuid.New()
	fanID := uuid.New()

	studentRepo := new(mocks.MockStudentRepository)
	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID, FanID: fanID, KKM: 70}, nil)

	uc := NewStudentSKSResultUseCase(new(mocks.StudentSKSResultRepositoryMock), new(mocks.FanCompletionStatusRepositoryMock), studentRepo, sksRepo, new(mocks.TeacherRepositoryMock), &noopAuditLogger{})
	resp, err := uc.CreateStudentSKSResult(ctx, dto.CreateStudentSKSResultRequest{
		StudentID: studentID.String(),
		SKSID:     sksID.String(),
		Score:     75,
		ExamDate:  ptrString("bad-date"),
	})

	assert.ErrorIs(t, err, domainErrors.ErrBadRequest)
	assert.Nil(t, resp)
}

func ptrString(s string) *string {
	return &s
}
