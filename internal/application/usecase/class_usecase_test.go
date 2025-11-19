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

func newClassUseCaseForTest() (*ClassUseCase, *mocks.ClassRepositoryMock, *mocks.FanRepositoryMock, *mocks.MockStudentRepository, *mocks.StudentClassEnrollmentRepositoryMock, *mocks.ClassStaffRepositoryMock) {
	classRepo := new(mocks.ClassRepositoryMock)
	fanRepo := new(mocks.FanRepositoryMock)
	studentRepo := new(mocks.MockStudentRepository)
	enrollmentRepo := new(mocks.StudentClassEnrollmentRepositoryMock)
	staffRepo := new(mocks.ClassStaffRepositoryMock)
	uc := NewClassUseCase(classRepo, fanRepo, studentRepo, enrollmentRepo, staffRepo, &noopAuditLogger{})
	return uc, classRepo, fanRepo, studentRepo, enrollmentRepo, staffRepo
}

func TestClassUseCase_CreateClass(t *testing.T) {
	uc, classRepo, fanRepo, _, _, _ := newClassUseCaseForTest()
	fanID := uuid.New()

	fanRepo.On("GetByID", mock.Anything, fanID).Return(&entity.Fan{ID: fanID}, nil)
	classRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := uc.CreateClass(context.Background(), dto.CreateClassRequest{FanID: fanID.String(), Name: "Class A", Capacity: 10})
	assert.NoError(t, err)
	assert.Equal(t, "Class A", resp.Name)
	classRepo.AssertExpectations(t)
	fanRepo.AssertExpectations(t)

	t.Run("fan not found", func(t *testing.T) {
		uc, _, fanRepo, _, _, _ := newClassUseCaseForTest()
		fanRepo.On("GetByID", mock.Anything, fanID).Return(nil, assert.AnError)
		resp, err := uc.CreateClass(context.Background(), dto.CreateClassRequest{FanID: fanID.String(), Name: "Class"})
		assert.ErrorIs(t, err, domainErrors.ErrFanNotFound)
		assert.Nil(t, resp)
	})
}

func TestClassUseCase_GetClass(t *testing.T) {
	uc, classRepo, _, _, _, _ := newClassUseCaseForTest()
	classID := uuid.New()
	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID, Name: "Class"}, nil)

	resp, err := uc.GetClass(context.Background(), classID)
	assert.NoError(t, err)
	assert.Equal(t, "Class", resp.Name)
	classRepo.AssertExpectations(t)
}

func TestClassUseCase_ListClasses(t *testing.T) {
	uc, classRepo, fanRepo, _, _, _ := newClassUseCaseForTest()
	fanID := uuid.New()
	fanRepo.On("GetByID", mock.Anything, fanID).Return(&entity.Fan{ID: fanID}, nil)
	classRepo.On("ListByFan", mock.Anything, fanID, 10, 0).Return([]*entity.Class{{ID: uuid.New(), Name: "C1"}}, int64(1), nil)

	resp, err := uc.ListClassesByFan(context.Background(), fanID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, resp.Classes, 1)
	fanRepo.AssertExpectations(t)
	classRepo.AssertExpectations(t)
}

func TestClassUseCase_UpdateClass(t *testing.T) {
	uc, classRepo, _, _, _, _ := newClassUseCaseForTest()
	classID := uuid.New()
	name := "Updated"
	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID, Name: "Old"}, nil)
	classRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *entity.Class) bool { return c.Name == name })).Return(nil)

	resp, err := uc.UpdateClass(context.Background(), classID, dto.UpdateClassRequest{Name: &name})
	assert.NoError(t, err)
	assert.Equal(t, name, resp.Name)
}

func TestClassUseCase_DeleteClass(t *testing.T) {
	uc, classRepo, _, _, _, _ := newClassUseCaseForTest()
	classID := uuid.New()
	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	classRepo.On("Delete", mock.Anything, classID).Return(nil)

	assert.NoError(t, uc.DeleteClass(context.Background(), classID))
	classRepo.AssertExpectations(t)
}

func TestClassUseCase_EnrollStudent(t *testing.T) {
	uc, classRepo, _, studentRepo, enrollmentRepo, _ := newClassUseCaseForTest()
	classID := uuid.New()
	studentID := uuid.New()

	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
	enrollmentRepo.On("GetActiveByStudentAndClass", mock.Anything, studentID, classID).Return(nil, nil)
	enrollmentRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	req := dto.EnrollStudentRequest{StudentID: studentID.String(), StartDate: time.Now().Format(time.RFC3339)}
	assert.NoError(t, uc.EnrollStudent(context.Background(), classID, req))
	classRepo.AssertExpectations(t)
	studentRepo.AssertExpectations(t)
	enrollmentRepo.AssertExpectations(t)

	t.Run("already enrolled", func(t *testing.T) {
		uc, classRepo, _, studentRepo, enrollmentRepo, _ := newClassUseCaseForTest()
		classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
		studentRepo.On("GetByID", mock.Anything, studentID).Return(&entity.Student{ID: studentID}, nil)
		enrollmentRepo.On("GetActiveByStudentAndClass", mock.Anything, studentID, classID).Return(&entity.StudentClassEnrollment{ID: uuid.New()}, nil)

		req := dto.EnrollStudentRequest{StudentID: studentID.String(), StartDate: time.Now().Format(time.RFC3339)}
		err := uc.EnrollStudent(context.Background(), classID, req)
		assert.ErrorIs(t, err, domainErrors.ErrStudentAlreadyEnrolled)
	})
}

func TestClassUseCase_AssignStaff(t *testing.T) {
	uc, classRepo, _, _, _, staffRepo := newClassUseCaseForTest()
	classID := uuid.New()
	userID := uuid.New()

	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	staffRepo.On("Assign", mock.Anything, mock.Anything).Return(nil)

	req := dto.AssignClassStaffRequest{UserID: userID.String(), Role: "class_manager"}
	assert.NoError(t, uc.AssignStaff(context.Background(), classID, req))
}
