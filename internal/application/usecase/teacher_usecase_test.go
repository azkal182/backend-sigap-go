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
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

type teacherNoopAuditLogger struct{}

func TestTeacherUseCase_CreateTeacher_ExistingUserNotFound(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, &teacherNoopAuditLogger{})

	teacherRepo.On("GetByCode", mock.Anything, "TCH-ERR").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "missing").Return(nil, assert.AnError)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode:      "TCH-ERR",
		FullName:         "Err Teacher",
		ExistingUsername: "Missing",
	})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrUserNotFound)
}

func TestTeacherUseCase_CreateTeacher_UserAlreadyAssigned(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, &teacherNoopAuditLogger{})

	userID := uuid.New()
	teacherRepo.On("GetByCode", mock.Anything, "TCH-DUP").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "duplicate").Return(&entity.User{ID: userID}, nil)
	teacherRepo.On("GetByUserID", mock.Anything, userID).Return(&entity.Teacher{ID: uuid.New()}, nil)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode:      "TCH-DUP",
		FullName:         "Dup Teacher",
		ExistingUsername: "duplicate",
	})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrTeacherUserAssigned)
}

func TestTeacherUseCase_CreateTeacher_MissingTeacherRole(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, &teacherNoopAuditLogger{})

	teacherRepo.On("GetByCode", mock.Anything, "TCH-NOROLE").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "norole").Return(nil, domainErrors.ErrUserNotFound)
	roleRepo.On("GetBySlug", mock.Anything, "teacher").Return(nil, assert.AnError)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode: "TCH-NOROLE",
		FullName:    "No Role",
	})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrRoleNotFound)
}

func TestTeacherUseCase_CreateTeacher_UserCreateFails(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, &teacherNoopAuditLogger{})

	teacherRepo.On("GetByCode", mock.Anything, "TCH-USRERR").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "userfails").Return(nil, domainErrors.ErrUserNotFound)
	roleRepo.On("GetBySlug", mock.Anything, "teacher").Return(&entity.Role{ID: uuid.New()}, nil)
	userRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode: "TCH-USRERR",
		FullName:    "User Fails",
	})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrInternalServer)
}

func TestTeacherUseCase_EnsureTeacherRole_Error(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, &teacherNoopAuditLogger{})

	roleRepo.On("GetBySlug", mock.Anything, "teacher").Return(nil, assert.AnError)
	err := uc.ensureTeacherRole(context.Background(), uuid.New())
	assert.ErrorIs(t, err, domainErrors.ErrRoleNotFound)
}

func TestTeacherUseCase_GetTeacher_WithUser(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherID := uuid.New()
	userID := uuid.New()
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, TeacherCode: "TCH-Get", FullName: "Get Teacher", UserID: &userID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)
	userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID, Username: "get-teacher"}, nil)

	resp, err := uc.GetTeacher(context.Background(), teacherID)
	assert.NoError(t, err)
	assert.Equal(t, "TCH-Get", resp.TeacherCode)
	assert.Equal(t, "get-teacher", resp.Username)
}

func TestTeacherUseCase_UpdateTeacher_SyncsUserName(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherID := uuid.New()
	userID := uuid.New()
	newName := "Updated Teacher"
	isActive := false
	joined := time.Now()

	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, FullName: "Old", UserID: &userID, IsActive: true}, nil)
	teacherRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID, Name: "Old", Username: "updatedteacher"}, nil)
	userRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	resp, err := uc.UpdateTeacher(context.Background(), teacherID, dto.UpdateTeacherRequest{
		FullName: &newName,
		IsActive: &isActive,
		JoinedAt: &joined,
	})

	assert.NoError(t, err)
	assert.Equal(t, newName, resp.FullName)
	assert.Equal(t, isActive, resp.IsActive)
	teacherRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestTeacherUseCase_DeactivateTeacher(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherID := uuid.New()
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID}, nil)
	teacherRepo.On("SoftDelete", mock.Anything, teacherID).Return(nil)

	assert.NoError(t, uc.DeactivateTeacher(context.Background(), teacherID))
	teacherRepo.AssertExpectations(t)
}

func TestTeacherUseCase_DeactivateTeacher_NotFound(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherID := uuid.New()
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(nil, assert.AnError)

	err := uc.DeactivateTeacher(context.Background(), teacherID)
	assert.ErrorIs(t, err, domainErrors.ErrTeacherNotFound)
}

func (n *teacherNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestTeacherUseCase_CreateTeacher_AutoUser(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherRepo.On("GetByCode", mock.Anything, "TCH-01").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "johndoe").Return(nil, domainErrors.ErrUserNotFound)
	teacherRoleID := uuid.New()
	roleRepo.On("GetBySlug", mock.Anything, "teacher").Return(&entity.Role{ID: teacherRoleID}, nil)
	userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
		return u.Username == "johndoe" && u.Name == "John Doe"
	})).Return(nil)
	teacherRepo.On("Create", mock.Anything, mock.MatchedBy(func(tch *entity.Teacher) bool {
		return tch.TeacherCode == "TCH-01" && tch.FullName == "John Doe"
	})).Return(nil)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode: "TCH-01",
		FullName:    "John Doe",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "TCH-01", resp.TeacherCode)
	assert.Equal(t, "John Doe", resp.FullName)
	teacherRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestTeacherUseCase_CreateTeacher_WithExistingUser(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	userID := uuid.New()
	teacherRoleID := uuid.New()

	teacherRepo.On("GetByCode", mock.Anything, "TCH-02").Return(nil, domainErrors.ErrTeacherNotFound)
	userRepo.On("GetByUsername", mock.Anything, "existinguser").Return(&entity.User{ID: userID, Username: "existinguser"}, nil)
	teacherRepo.On("GetByUserID", mock.Anything, userID).Return(nil, domainErrors.ErrTeacherNotFound)
	roleRepo.On("GetBySlug", mock.Anything, "teacher").Return(&entity.Role{ID: teacherRoleID}, nil)
	userRepo.On("AssignRole", mock.Anything, userID, teacherRoleID).Return(nil)
	teacherRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := uc.CreateTeacher(context.Background(), dto.CreateTeacherRequest{
		TeacherCode:      "TCH-02",
		FullName:         "Existing User",
		ExistingUsername: "ExistingUser",
	})

	assert.NoError(t, err)
	assert.Equal(t, "TCH-02", resp.TeacherCode)
	assert.Equal(t, "Existing User", resp.FullName)
	teacherRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestTeacherUseCase_ListTeachers_NormalizesPagination(t *testing.T) {
	teacherRepo := new(mocks.MockTeacherRepository)
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	logger := &teacherNoopAuditLogger{}
	uc := NewTeacherUseCase(teacherRepo, userRepo, roleRepo, logger)

	teacherID := uuid.New()
	userID := uuid.New()
	teacherEntities := []*entity.Teacher{{
		ID:          teacherID,
		TeacherCode: "TCH-03",
		FullName:    "Math Teacher",
		UserID:      &userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}}

	teacherRepo.On("List", mock.Anything, mock.MatchedBy(func(f repository.TeacherFilter) bool {
		return f.Page == 1 && f.PageSize == 10 && f.Keyword == "math"
	})).Return(teacherEntities, int64(1), nil)
	userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID, Username: "mathteacher"}, nil)

	resp, err := uc.ListTeachers(context.Background(), 0, 0, "math", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Len(t, resp.Teachers, 1)
	assert.Equal(t, "mathteacher", resp.Teachers[0].Username)
	teacherRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
