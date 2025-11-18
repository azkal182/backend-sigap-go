package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// DormitoryUseCase handles dormitory management use cases
type DormitoryUseCase struct {
	dormitoryRepo repository.DormitoryRepository
	userRepo      repository.UserRepository
	auditLogger   appService.AuditLogger
}

// NewDormitoryUseCase creates a new dormitory use case
func NewDormitoryUseCase(
	dormitoryRepo repository.DormitoryRepository,
	userRepo repository.UserRepository,
	auditLogger appService.AuditLogger,
) *DormitoryUseCase {
	return &DormitoryUseCase{
		dormitoryRepo: dormitoryRepo,
		userRepo:      userRepo,
		auditLogger:   auditLogger,
	}
}

// CreateDormitory creates a new dormitory
func (uc *DormitoryUseCase) CreateDormitory(ctx context.Context, req dto.CreateDormitoryRequest) (*dto.DormitoryResponse, error) {
	dormitory := &entity.Dormitory{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.dormitoryRepo.Create(ctx, dormitory); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "dormitory", "dorm:create", dormitory.ID.String(), map[string]string{
		"name":        dormitory.Name,
		"description": dormitory.Description,
	})

	return uc.toDormitoryResponse(dormitory), nil
}

// GetDormitoryByID retrieves a dormitory by ID
func (uc *DormitoryUseCase) GetDormitoryByID(ctx context.Context, id uuid.UUID) (*dto.DormitoryResponse, error) {
	dormitory, err := uc.dormitoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrDormitoryNotFound
	}

	return uc.toDormitoryResponse(dormitory), nil
}

// UpdateDormitory updates a dormitory
func (uc *DormitoryUseCase) UpdateDormitory(ctx context.Context, id uuid.UUID, req dto.UpdateDormitoryRequest) (*dto.DormitoryResponse, error) {
	dormitory, err := uc.dormitoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrDormitoryNotFound
	}

	if req.Name != "" {
		dormitory.Name = req.Name
	}
	if req.Description != "" {
		dormitory.Description = req.Description
	}
	if req.IsActive != nil {
		dormitory.IsActive = *req.IsActive
	}

	dormitory.UpdatedAt = time.Now()

	if err := uc.dormitoryRepo.Update(ctx, dormitory); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "dormitory", "dorm:update", dormitory.ID.String(), map[string]string{
		"name":        dormitory.Name,
		"description": dormitory.Description,
	})

	return uc.toDormitoryResponse(dormitory), nil
}

// DeleteDormitory deletes a dormitory (soft delete)
func (uc *DormitoryUseCase) DeleteDormitory(ctx context.Context, id uuid.UUID) error {
	dormitory, err := uc.dormitoryRepo.GetByID(ctx, id)
	if err != nil {
		return domainErrors.ErrDormitoryNotFound
	}

	if err := uc.dormitoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "dormitory", "dorm:delete", id.String(), map[string]string{
		"name":        dormitory.Name,
		"description": dormitory.Description,
	})

	return nil
}

// ListDormitories retrieves a paginated list of dormitories
func (uc *DormitoryUseCase) ListDormitories(ctx context.Context, page, pageSize int) (*dto.ListDormitoriesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	dormitories, total, err := uc.dormitoryRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	dormitoryResponses := make([]dto.DormitoryResponse, 0, len(dormitories))
	for _, dormitory := range dormitories {
		dormitoryResponses = append(dormitoryResponses, *uc.toDormitoryResponse(dormitory))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListDormitoriesResponse{
		Dormitories: dormitoryResponses,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

// toDormitoryResponse converts entity.Dormitory to dto.DormitoryResponse
func (uc *DormitoryUseCase) toDormitoryResponse(dormitory *entity.Dormitory) *dto.DormitoryResponse {
	return &dto.DormitoryResponse{
		ID:          dormitory.ID.String(),
		Name:        dormitory.Name,
		Description: dormitory.Description,
		IsActive:    dormitory.IsActive,
		CreatedAt:   dormitory.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   dormitory.UpdatedAt.Format(time.RFC3339),
	}
}
