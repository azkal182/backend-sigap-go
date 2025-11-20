package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

type sksNoopAuditLogger struct{}

func (n *sksNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func newSKSDefinitionUC() (*SKSDefinitionUseCase, *mocks.SKSDefinitionRepositoryMock, *mocks.FanRepositoryMock, *mocks.SubjectRepositoryMock) {
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	fanRepo := new(mocks.FanRepositoryMock)
	subjectRepo := new(mocks.SubjectRepositoryMock)
	uc := NewSKSDefinitionUseCase(sksRepo, fanRepo, subjectRepo, &sksNoopAuditLogger{})
	return uc, sksRepo, fanRepo, subjectRepo
}

func newSKSExamUC() (*SKSExamScheduleUseCase, *mocks.SKSExamScheduleRepositoryMock, *mocks.SKSDefinitionRepositoryMock, *mocks.MockTeacherRepository) {
	examRepo := new(mocks.SKSExamScheduleRepositoryMock)
	sksRepo := new(mocks.SKSDefinitionRepositoryMock)
	teacherRepo := new(mocks.MockTeacherRepository)
	uc := NewSKSExamScheduleUseCase(examRepo, sksRepo, teacherRepo, &sksNoopAuditLogger{})
	return uc, examRepo, sksRepo, teacherRepo
}

func TestSKSDefinitionUseCase_Create(t *testing.T) {
	uc, sksRepo, fanRepo, subjectRepo := newSKSDefinitionUC()
	fanID := uuid.New()
	subjectID := uuid.New()

	fanRepo.On("GetByID", mock.Anything, fanID).Return(&entity.Fan{ID: fanID}, nil)
	subjectRepo.On("GetByID", mock.Anything, subjectID).Return(&entity.Subject{ID: subjectID}, nil)
	sksRepo.On("GetByCode", mock.Anything, "SKS-01").Return(nil, domainErrors.ErrSKSDefinitionNotFound)
	sksRepo.On("Create", mock.Anything, mock.MatchedBy(func(def *entity.SKSDefinition) bool {
		return def.FanID == fanID && def.SubjectID != nil && *def.SubjectID == subjectID && def.IsActive
	})).Return(nil)

	isActive := true
	req := dto.CreateSKSDefinitionRequest{
		FanID:       fanID.String(),
		SubjectID:   stringPtrSKS(subjectID.String()),
		Code:        "SKS-01",
		Name:        "Tahfidz",
		KKM:         75,
		Description: "hafalan juz 30",
		IsActive:    &isActive,
	}

	resp, err := uc.CreateSKSDefinition(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "SKS-01", resp.Code)
	assert.NotNil(t, resp.SubjectID)
	sksRepo.AssertExpectations(t)
}

func TestSKSDefinitionUseCase_CreateDuplicate(t *testing.T) {
	uc, sksRepo, fanRepo, _ := newSKSDefinitionUC()
	fanID := uuid.New()

	fanRepo.On("GetByID", mock.Anything, fanID).Return(&entity.Fan{ID: fanID}, nil)
	sksRepo.On("GetByCode", mock.Anything, "SKS-02").Return(&entity.SKSDefinition{ID: uuid.New(), Code: "SKS-02"}, nil)

	_, err := uc.CreateSKSDefinition(context.Background(), dto.CreateSKSDefinitionRequest{
		FanID: fanID.String(),
		Code:  "SKS-02",
		Name:  "Fiqh",
	})
	assert.ErrorIs(t, err, domainErrors.ErrSKSDefinitionAlreadyExist)
}

func TestSKSDefinitionUseCase_UpdateAndClearSubject(t *testing.T) {
	uc, sksRepo, _, subjectRepo := newSKSDefinitionUC()
	id := uuid.New()
	currentSubject := uuid.New()
	definition := &entity.SKSDefinition{ID: id, FanID: uuid.New(), SubjectID: &currentSubject, Code: "SKS-03", Name: "Akidah", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	sksRepo.On("GetByID", mock.Anything, id).Return(definition, nil)
	sksRepo.On("Update", mock.Anything, definition).Return(nil)

	resp, err := uc.UpdateSKSDefinition(context.Background(), id, dto.UpdateSKSDefinitionRequest{SubjectID: stringPtrSKS("")})
	assert.NoError(t, err)
	assert.Nil(t, resp.SubjectID)
	subjectRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestSKSDefinitionUseCase_ListAndDelete(t *testing.T) {
	uc, sksRepo, _, _ := newSKSDefinitionUC()
	fanID := uuid.New()
	definition := &entity.SKSDefinition{ID: uuid.New(), FanID: fanID, Code: "SKS-04", Name: "Hadits", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	sksRepo.On("List", mock.Anything, fanID, 10, 0).Return([]*entity.SKSDefinition{definition}, int64(1), nil)
	resp, err := uc.ListSKSDefinitions(context.Background(), fanID.String(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, resp.Definitions, 1)

	sksRepo.On("GetByID", mock.Anything, definition.ID).Return(definition, nil)
	sksRepo.On("Delete", mock.Anything, definition.ID).Return(nil)
	assert.NoError(t, uc.DeleteSKSDefinition(context.Background(), definition.ID))
}

func TestSKSExamScheduleUseCase_Create(t *testing.T) {
	uc, examRepo, sksRepo, teacherRepo := newSKSExamUC()
	sksID := uuid.New()
	teacherID := uuid.New()

	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, IsActive: true}, nil)
	examRepo.On("Create", mock.Anything, mock.MatchedBy(func(ex *entity.SKSExamSchedule) bool {
		return ex.SKSID == sksID && ex.ExaminerID != nil && *ex.ExaminerID == teacherID
	})).Return(nil)

	req := dto.CreateSKSExamScheduleRequest{
		SKSID:      sksID.String(),
		ExaminerID: stringPtrSKS(teacherID.String()),
		ExamDate:   "2025-01-10",
		ExamTime:   "08:30",
		Location:   "Aula",
		Notes:      "Gelombang 1",
	}

	resp, err := uc.CreateSKSExamSchedule(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, sksID.String(), resp.SKSID)
	assert.Equal(t, "2025-01-10", resp.ExamDate)
	examRepo.AssertExpectations(t)
}

func TestSKSExamScheduleUseCase_Create_InvalidTeacher(t *testing.T) {
	uc, _, sksRepo, teacherRepo := newSKSExamUC()
	sksID := uuid.New()
	teacherID := uuid.New()

	sksRepo.On("GetByID", mock.Anything, sksID).Return(&entity.SKSDefinition{ID: sksID}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(nil, assert.AnError)

	_, err := uc.CreateSKSExamSchedule(context.Background(), dto.CreateSKSExamScheduleRequest{
		SKSID:      sksID.String(),
		ExaminerID: stringPtrSKS(teacherID.String()),
		ExamDate:   "2025-02-01",
		ExamTime:   "09:00",
	})
	assert.ErrorIs(t, err, domainErrors.ErrTeacherNotFound)
}

func TestSKSExamScheduleUseCase_UpdateAndList(t *testing.T) {
	uc, examRepo, _, teacherRepo := newSKSExamUC()
	examID := uuid.New()
	sksID := uuid.New()
	teacherID := uuid.New()
	exam := &entity.SKSExamSchedule{ID: examID, SKSID: sksID, ExamDate: time.Now(), ExamTime: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}

	examRepo.On("GetByID", mock.Anything, examID).Return(exam, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, IsActive: true}, nil)
	examRepo.On("Update", mock.Anything, exam).Return(nil)

	date := "2025-03-01"
	timeStr := "14:00"
	req := dto.UpdateSKSExamScheduleRequest{
		ExamDate:   stringPtrSKS(date),
		ExamTime:   stringPtrSKS(timeStr),
		ExaminerID: stringPtrSKS(teacherID.String()),
	}

	resp, err := uc.UpdateSKSExamSchedule(context.Background(), examID, req)
	assert.NoError(t, err)
	assert.Equal(t, date, resp.ExamDate)
	assert.Equal(t, timeStr, resp.ExamTime)

	examRepo.On("ListBySKS", mock.Anything, sksID, 10, 0).Return([]*entity.SKSExamSchedule{exam}, int64(1), nil)
	listResp, err := uc.ListSKSExamSchedules(context.Background(), sksID.String(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, listResp.Exams, 1)
}

func TestSKSExamScheduleUseCase_Delete(t *testing.T) {
	uc, examRepo, _, _ := newSKSExamUC()
	examID := uuid.New()
	exam := &entity.SKSExamSchedule{ID: examID}

	examRepo.On("GetByID", mock.Anything, examID).Return(exam, nil)
	examRepo.On("Delete", mock.Anything, examID).Return(nil)
	assert.NoError(t, uc.DeleteSKSExamSchedule(context.Background(), examID))
}

func stringPtrSKS(val string) *string {
	return &val
}
