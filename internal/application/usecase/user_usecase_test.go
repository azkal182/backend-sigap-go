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

type noopAuditLogger struct{}

func (n *noopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestUserUseCase_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		req           dto.CreateUserRequest
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockRoleRepository)
		expectedError error
	}{
		{
			name: "success - create user without roles",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, domainErrors.ErrUserNotFound)
				// Default role lookup ("user")
				roleRepo.On("GetBySlug", mock.Anything, "user").Return(&entity.Role{
					ID:   uuid.New(),
					Name: "User",
				}, nil)
				userRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				userRepo.On("GetWithRoles", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       uuid.New(),
					Username: "newuser",
					Name:     "New User",
					IsActive: true,
					Roles:    []entity.Role{},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "success - create user with roles",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
				RoleIDs:  []string{uuid.New().String()},
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleID := uuid.New()
				// Inject roleID into request
				// (we rely on CreateUser reading req.RoleIDs and then calling GetByID with this ID)
				// Note: we can't modify tt.req here easily, so we just relax the ID matching below.
				userRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, domainErrors.ErrUserNotFound)
				roleRepo.On("GetByID", mock.Anything, mock.Anything).Return(&entity.Role{
					ID:   roleID,
					Name: "user",
				}, nil)
				userRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				userRepo.On("GetWithRoles", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       uuid.New(),
					Username: "newuser",
					Name:     "New User",
					IsActive: true,
					Roles: []entity.Role{
						{ID: roleID, Name: "user"},
					},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "failure - user already exists",
			req: dto.CreateUserRequest{
				Username: "existing",
				Password: "password123",
				Name:     "Existing User",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.On("GetByUsername", mock.Anything, "existing").Return(&entity.User{
					ID:       uuid.New(),
					Username: "existing",
				}, nil)
			},
			expectedError: domainErrors.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			roleRepo := new(mocks.MockRoleRepository)
			tt.setupMocks(userRepo, roleRepo)

			auditLogger := &noopAuditLogger{}
			userUseCase := NewUserUseCase(userRepo, roleRepo, auditLogger)
			resp, err := userUseCase.CreateUser(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Username, resp.Username)
				assert.Equal(t, tt.req.Name, resp.Name)
			}

			userRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMocks    func(*mocks.MockUserRepository)
		expectedError error
	}{
		{
			name:   "success - get user by id",
			userID: userID,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(&entity.User{
					ID:       userID,
					Username: "user",
					Name:     "Test User",
					IsActive: true,
					Roles:    []entity.Role{{ID: uuid.New(), Name: "user"}},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "failure - user not found",
			userID: userID,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(nil, domainErrors.ErrUserNotFound)
			},
			expectedError: domainErrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			roleRepo := new(mocks.MockRoleRepository)
			tt.setupMocks(userRepo)
			auditLogger := &noopAuditLogger{}
			userUseCase := NewUserUseCase(userRepo, roleRepo, auditLogger)
			resp, err := userUseCase.GetUserByID(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.userID.String(), resp.ID)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		req           dto.UpdateUserRequest
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockRoleRepository)
		expectedError error
	}{
		{
			name:   "success - update user name",
			userID: userID,
			req: dto.UpdateUserRequest{
				Name: "Updated Name",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:       userID,
					Username: "user",
					Name:     "Old Name",
					IsActive: true,
				}, nil)
				userRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(&entity.User{
					ID:       userID,
					Username: "user",
					Name:     "Updated Name",
					IsActive: true,
					Roles:    []entity.Role{},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "failure - user not found",
			userID: userID,
			req: dto.UpdateUserRequest{
				Name: "Updated Name",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.On("GetByID", mock.Anything, userID).Return(nil, domainErrors.ErrUserNotFound)
			},
			expectedError: domainErrors.ErrUserNotFound,
		},
		{
			name:   "failure - Username already taken",
			userID: userID,
			req: dto.UpdateUserRequest{
				Username: "taken",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:       userID,
					Username: "user",
				}, nil)
				otherUserID := uuid.New()
				userRepo.On("GetByUsername", mock.Anything, "taken").Return(&entity.User{
					ID:       otherUserID,
					Username: "taken",
				}, nil)
			},
			expectedError: domainErrors.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			roleRepo := new(mocks.MockRoleRepository)
			tt.setupMocks(userRepo, roleRepo)
			auditLogger := &noopAuditLogger{}
			userUseCase := NewUserUseCase(userRepo, roleRepo, auditLogger)
			resp, err := userUseCase.UpdateUser(context.Background(), tt.userID, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			userRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMocks    func(*mocks.MockUserRepository)
		expectedError error
	}{
		{
			name:   "success - delete user",
			userID: userID,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID: userID,
				}, nil)
				userRepo.On("Delete", mock.Anything, userID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "failure - user not found",
			userID: userID,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.On("GetByID", mock.Anything, userID).Return(nil, domainErrors.ErrUserNotFound)
			},
			expectedError: domainErrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			roleRepo := new(mocks.MockRoleRepository)
			tt.setupMocks(userRepo)

			auditLogger := &noopAuditLogger{}
			userUseCase := NewUserUseCase(userRepo, roleRepo, auditLogger)
			err := userUseCase.DeleteUser(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_ListUsers(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		pageSize      int
		setupMocks    func(*mocks.MockUserRepository)
		expectedError error
	}{
		{
			name:     "success - list users with pagination",
			page:     1,
			pageSize: 10,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				users := []*entity.User{
					{ID: uuid.New(), Username: "user1", Name: "User 1"},
					{ID: uuid.New(), Username: "user2", Name: "User 2"},
				}
				userRepo.On("List", mock.Anything, 10, 0).Return(users, int64(2), nil)
				userRepo.On("GetWithRoles", mock.Anything, mock.Anything).Return(&entity.User{
					ID:    uuid.New(),
					Roles: []entity.Role{},
				}, nil).Times(2)
			},
			expectedError: nil,
		},
		{
			name:     "success - default page and pageSize",
			page:     0,
			pageSize: 0,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				users := []*entity.User{}
				userRepo.On("List", mock.Anything, 10, 0).Return(users, int64(0), nil)
			},
			expectedError: nil,
		},
		{
			name:     "success - max pageSize capped at 100",
			page:     1,
			pageSize: 200,
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				users := []*entity.User{}
				userRepo.On("List", mock.Anything, 100, 0).Return(users, int64(0), nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			roleRepo := new(mocks.MockRoleRepository)
			tt.setupMocks(userRepo)

			auditLogger := &noopAuditLogger{}
			userUseCase := NewUserUseCase(userRepo, roleRepo, auditLogger)
			resp, err := userUseCase.ListUsers(context.Background(), tt.page, tt.pageSize)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				// ListUsers normalizes page < 1 to 1, so for default case (page=0) we expect page 1
				expectedPage := tt.page
				if expectedPage < 1 {
					expectedPage = 1
				}
				assert.Equal(t, expectedPage, resp.Page)
			}

			userRepo.AssertExpectations(t)
		})
	}
}
