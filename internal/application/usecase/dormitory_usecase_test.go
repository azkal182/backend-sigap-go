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

func TestDormitoryUseCase_CreateDormitory(t *testing.T) {
	tests := []struct {
		name          string
		req           dto.CreateDormitoryRequest
		setupMocks    func(*mocks.MockDormitoryRepository)
		expectedError error
	}{
		{
			name: "success - create dormitory",
			req: dto.CreateDormitoryRequest{
				Name:        "Test Dormitory",
				Gender:      "male",
				Level:       "senior",
				Code:        "DRM01",
				Description: "A test dormitory",
			},
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("Create", mock.Anything, mock.MatchedBy(func(d *entity.Dormitory) bool {
					return d.Name == "Test Dormitory" && d.Description == "A test dormitory" && d.Gender == "male" && d.Level == "senior" && d.Code == "DRM01"
				})).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dormRepo := new(mocks.MockDormitoryRepository)
			userRepo := new(mocks.MockUserRepository)
			tt.setupMocks(dormRepo)

			auditLogger := &noopAuditLogger{}
			dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)
			resp, err := dormUseCase.CreateDormitory(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Name, resp.Name)
			}

			dormRepo.AssertExpectations(t)
		})
	}
}

func TestDormitoryUseCase_AssignAndRemoveUser(t *testing.T) {
	dormitoryID := uuid.New()
	userID := uuid.New()

	t.Run("assign user success", func(t *testing.T) {
		dormRepo := new(mocks.MockDormitoryRepository)
		userRepo := new(mocks.MockUserRepository)

		dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(&entity.Dormitory{ID: dormitoryID}, nil)
		userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID}, nil)
		dormRepo.On("AssignToUser", mock.Anything, userID, dormitoryID).Return(nil)

		auditLogger := &noopAuditLogger{}
		dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)

		err := dormUseCase.AssignUser(context.Background(), dormitoryID, userID)
		assert.NoError(t, err)
		dormRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("remove user success", func(t *testing.T) {
		dormRepo := new(mocks.MockDormitoryRepository)
		userRepo := new(mocks.MockUserRepository)

		dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(&entity.Dormitory{ID: dormitoryID}, nil)
		userRepo.On("GetByID", mock.Anything, userID).Return(&entity.User{ID: userID}, nil)
		dormRepo.On("RemoveFromUser", mock.Anything, userID, dormitoryID).Return(nil)

		auditLogger := &noopAuditLogger{}
		dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)

		err := dormUseCase.RemoveUser(context.Background(), dormitoryID, userID)
		assert.NoError(t, err)
		dormRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}

func strPtr(val string) *string {
	return &val
}

func TestDormitoryUseCase_GetDormitoryByID(t *testing.T) {
	dormitoryID := uuid.New()

	tests := []struct {
		name          string
		dormitoryID   uuid.UUID
		setupMocks    func(*mocks.MockDormitoryRepository)
		expectedError error
	}{
		{
			name:        "success - get dormitory by id",
			dormitoryID: dormitoryID,
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(&entity.Dormitory{
					ID:          dormitoryID,
					Name:        "Test Dormitory",
					Description: "A test dormitory",
					IsActive:    true,
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:        "failure - dormitory not found",
			dormitoryID: dormitoryID,
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(nil, domainErrors.ErrDormitoryNotFound)
			},
			expectedError: domainErrors.ErrDormitoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dormRepo := new(mocks.MockDormitoryRepository)
			userRepo := new(mocks.MockUserRepository)
			tt.setupMocks(dormRepo)

			auditLogger := &noopAuditLogger{}
			dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)
			resp, err := dormUseCase.GetDormitoryByID(context.Background(), tt.dormitoryID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.dormitoryID.String(), resp.ID)
			}

			dormRepo.AssertExpectations(t)
		})
	}
}

func TestDormitoryUseCase_UpdateDormitory(t *testing.T) {
	dormitoryID := uuid.New()

	tests := []struct {
		name          string
		dormitoryID   uuid.UUID
		req           dto.UpdateDormitoryRequest
		setupMocks    func(*mocks.MockDormitoryRepository)
		expectedError error
	}{
		{
			name:        "success - update dormitory fields",
			dormitoryID: dormitoryID,
			req: dto.UpdateDormitoryRequest{
				Name:   "Updated Name",
				Code:   strPtr("NEWCODE"),
				Level:  strPtr("junior"),
				Gender: strPtr("female"),
			},
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(&entity.Dormitory{
					ID:     dormitoryID,
					Name:   "Old Name",
					Code:   "OLDCODE",
					Level:  "senior",
					Gender: "male",
				}, nil)
				dormRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "failure - dormitory not found",
			dormitoryID: dormitoryID,
			req: dto.UpdateDormitoryRequest{
				Name: "Updated Name",
			},
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(nil, domainErrors.ErrDormitoryNotFound)
			},
			expectedError: domainErrors.ErrDormitoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dormRepo := new(mocks.MockDormitoryRepository)
			userRepo := new(mocks.MockUserRepository)
			tt.setupMocks(dormRepo)

			auditLogger := &noopAuditLogger{}
			dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)
			resp, err := dormUseCase.UpdateDormitory(context.Background(), tt.dormitoryID, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			dormRepo.AssertExpectations(t)
		})
	}
}

func TestDormitoryUseCase_DeleteDormitory(t *testing.T) {
	dormitoryID := uuid.New()

	tests := []struct {
		name          string
		dormitoryID   uuid.UUID
		setupMocks    func(*mocks.MockDormitoryRepository)
		expectedError error
	}{
		{
			name:        "success - delete dormitory",
			dormitoryID: dormitoryID,
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(&entity.Dormitory{
					ID: dormitoryID,
				}, nil)
				dormRepo.On("Delete", mock.Anything, dormitoryID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "failure - dormitory not found",
			dormitoryID: dormitoryID,
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormRepo.On("GetByID", mock.Anything, dormitoryID).Return(nil, domainErrors.ErrDormitoryNotFound)
			},
			expectedError: domainErrors.ErrDormitoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dormRepo := new(mocks.MockDormitoryRepository)
			userRepo := new(mocks.MockUserRepository)
			tt.setupMocks(dormRepo)

			auditLogger := &noopAuditLogger{}
			dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)
			err := dormUseCase.DeleteDormitory(context.Background(), tt.dormitoryID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			dormRepo.AssertExpectations(t)
		})
	}
}

func TestDormitoryUseCase_ListDormitories(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		pageSize      int
		setupMocks    func(*mocks.MockDormitoryRepository)
		expectedError error
	}{
		{
			name:     "success - list dormitories with pagination",
			page:     1,
			pageSize: 10,
			setupMocks: func(dormRepo *mocks.MockDormitoryRepository) {
				dormitories := []*entity.Dormitory{
					{ID: uuid.New(), Name: "Dormitory 1"},
					{ID: uuid.New(), Name: "Dormitory 2"},
				}
				dormRepo.On("List", mock.Anything, 10, 0).Return(dormitories, int64(2), nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dormRepo := new(mocks.MockDormitoryRepository)
			userRepo := new(mocks.MockUserRepository)
			tt.setupMocks(dormRepo)

			auditLogger := &noopAuditLogger{}
			dormUseCase := NewDormitoryUseCase(dormRepo, userRepo, auditLogger)
			resp, err := dormUseCase.ListDormitories(context.Background(), tt.page, tt.pageSize)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.page, resp.Page)
			}

			dormRepo.AssertExpectations(t)
		})
	}
}
