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

func TestStudentUseCase_CreateStudent(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)

	req := dto.CreateStudentRequest{
		StudentNumber: "STD001",
		FullName:      "John Doe",
		BirthDate:     time.Now(),
		Gender:        "male",
		ParentName:    "Parent",
	}

	studentRepo.On("GetByStudentNumber", mock.Anything, "STD001").Return(nil, domainErrors.ErrStudentNotFound)
	studentRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.CreateStudent(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.StudentNumber, resp.StudentNumber)
	studentRepo.AssertExpectations(t)

	// duplicate scenario
	dupRepo := new(mocks.MockStudentRepository)
	dupRepo.On("GetByStudentNumber", mock.Anything, "STD001").Return(&entity.Student{ID: uuid.New()}, nil)
	ucDup := NewStudentUseCase(dupRepo, dormRepo, &noopAuditLogger{})
	resp, err = ucDup.CreateStudent(ctx, req)
	assert.ErrorIs(t, err, domainErrors.ErrStudentAlreadyExists)
	assert.Nil(t, resp)
}

func TestStudentUseCase_UpdateStudent(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	studentID := uuid.New()
	fullName := "Updated"
	req := dto.UpdateStudentRequest{FullName: &fullName}

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID, FullName: "Old"}, nil)
	studentRepo.On("Update", mock.Anything, mock.MatchedBy(func(s *entity.Student) bool {
		return s.ID == studentID && s.FullName == fullName
	})).Return(nil)
	studentRepo.On("ListHistory", mock.Anything, studentID).Return([]*entity.StudentDormitoryHistory{}, nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.UpdateStudent(ctx, studentID, req)
	assert.NoError(t, err)
	assert.Equal(t, fullName, resp.FullName)
	studentRepo.AssertExpectations(t)

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		repo.On("GetByID", mock.Anything, studentID).Return(nil, domainErrors.ErrStudentNotFound)
		uc := NewStudentUseCase(repo, dormRepo, &noopAuditLogger{})
		resp, err := uc.UpdateStudent(ctx, studentID, req)
		assert.ErrorIs(t, err, domainErrors.ErrStudentNotFound)
		assert.Nil(t, resp)
	})

	t.Run("mutate failure", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		repo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError)
		uc := NewStudentUseCase(repo, dormRepo, &noopAuditLogger{})
		resp, err := uc.UpdateStudent(ctx, studentID, req)
		assert.ErrorIs(t, err, domainErrors.ErrInternalServer)
		assert.Nil(t, resp)
	})
}

func TestStudentUseCase_UpdateStudentStatus(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID, Status: entity.StudentStatusInactive}, nil)
	studentRepo.On("UpdateStatus", mock.Anything, studentID, entity.StudentStatusActive, true).Return(nil)
	studentRepo.On("ListHistory", mock.Anything, studentID).Return([]*entity.StudentDormitoryHistory{}, nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.UpdateStudentStatus(ctx, studentID, entity.StudentStatusActive)
	assert.NoError(t, err)
	assert.Equal(t, entity.StudentStatusActive, resp.Status)
	assert.True(t, resp.IsActive)
	studentRepo.AssertExpectations(t)

	t.Run("update status fails", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		repo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
		repo.On("UpdateStatus", mock.Anything, studentID, entity.StudentStatusInactive, false).Return(assert.AnError)
		ucErr := NewStudentUseCase(repo, dormRepo, &noopAuditLogger{})
		resp, err := ucErr.UpdateStudentStatus(ctx, studentID, entity.StudentStatusInactive)
		assert.ErrorIs(t, err, domainErrors.ErrInternalServer)
		assert.Nil(t, resp)
	})
}

func TestStudentUseCase_GetStudentByID(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	studentID := uuid.New()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID, StudentNumber: "STD"}, nil)
	studentRepo.On("ListHistory", mock.Anything, studentID).Return([]*entity.StudentDormitoryHistory{
		{ID: uuid.New(), StudentID: studentID, DormitoryID: uuid.New(), StartDate: time.Now()},
	}, nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.GetStudentByID(ctx, studentID)
	assert.NoError(t, err)
	assert.Equal(t, studentID.String(), resp.ID)
	assert.Len(t, resp.DormitoryHistory, 1)
	studentRepo.AssertExpectations(t)
}

func TestStudentUseCase_ListStudents(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)

	studentRepo.On("List", mock.Anything, 10, 0).Return([]*entity.Student{{ID: uuid.New(), StudentNumber: "S1"}}, int64(1), nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.ListStudents(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	studentRepo.AssertExpectations(t)

	t.Run("repo error", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		repo.On("List", mock.Anything, 10, 0).Return(nil, int64(0), assert.AnError)
		ucErr := NewStudentUseCase(repo, dormRepo, &noopAuditLogger{})
		resp, err := ucErr.ListStudents(ctx, 1, 10)
		assert.ErrorIs(t, err, domainErrors.ErrInternalServer)
		assert.Nil(t, resp)
	})
}

func TestStudentUseCase_MutateStudentDormitory(t *testing.T) {
	ctx := context.Background()
	studentRepo := new(mocks.MockStudentRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	studentID := uuid.New()
	dormID := uuid.New()
	startDate := time.Now()

	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)
	studentRepo.On("GetActiveHistory", mock.Anything, studentID).Return(&entity.StudentDormitoryHistory{ID: uuid.New()}, nil)
	studentRepo.On("CloseHistory", mock.Anything, mock.Anything, startDate).Return(nil)
	studentRepo.On("CreateHistory", mock.Anything, mock.Anything).Return(nil)
	studentRepo.On("ListHistory", mock.Anything, studentID).Return([]*entity.StudentDormitoryHistory{}, nil)

	uc := NewStudentUseCase(studentRepo, dormRepo, &noopAuditLogger{})
	resp, err := uc.MutateStudentDormitory(ctx, studentID, dormID, startDate)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	studentRepo.AssertExpectations(t)
	dormRepo.AssertExpectations(t)

	t.Run("dorm not found", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		dRepo := new(mocks.MockDormitoryRepository)
		repo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
		dRepo.On("GetByID", mock.Anything, dormID).Return(nil, domainErrors.ErrDormitoryNotFound)
		uc := NewStudentUseCase(repo, dRepo, &noopAuditLogger{})
		resp, err := uc.MutateStudentDormitory(ctx, studentID, dormID, time.Now())
		assert.ErrorIs(t, err, domainErrors.ErrDormitoryNotFound)
		assert.Nil(t, resp)
	})
}
