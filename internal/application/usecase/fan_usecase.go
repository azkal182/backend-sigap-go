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

// FanUseCase orchestrates FAN operations.
type FanUseCase struct {
	fanRepo       repository.FanRepository
	dormitoryRepo repository.DormitoryRepository
	auditLogger   appService.AuditLogger
}

// NewFanUseCase builds a FanUseCase instance.
func NewFanUseCase(
	fanRepo repository.FanRepository,
	dormitoryRepo repository.DormitoryRepository,
	auditLogger appService.AuditLogger,
) *FanUseCase {
	return &FanUseCase{fanRepo: fanRepo, dormitoryRepo: dormitoryRepo, auditLogger: auditLogger}
}

// CreateFan creates a new FAN entry.
func (uc *FanUseCase) CreateFan(ctx context.Context, req dto.CreateFanRequest) (*dto.FanResponse, error) {
	dormitoryID, err := uuid.Parse(req.DormitoryID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if err := uc.ensureDormitoryExists(ctx, dormitoryID); err != nil {
		return nil, err
	}

	now := time.Now()
	fan := &entity.Fan{
		ID:          uuid.New(),
		DormitoryID: dormitoryID,
		Name:        req.Name,
		Level:       req.Level,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.fanRepo.Create(ctx, fan); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "fan", "fan:create", fan.ID.String(), map[string]string{
		"name":  fan.Name,
		"level": fan.Level,
	})

	return uc.toFanResponse(fan), nil
}

// GetFan retrieves a fan by ID.
func (uc *FanUseCase) GetFan(ctx context.Context, id uuid.UUID) (*dto.FanResponse, error) {
	fan, err := uc.fanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrFanNotFound
	}
	return uc.toFanResponse(fan), nil
}

// ListFans returns paginated fans result.
func (uc *FanUseCase) ListFans(ctx context.Context, page, pageSize int) (*dto.ListFansResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	fans, total, err := uc.fanRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.buildListResponse(fans, total, page, pageSize), nil
}

// ListFansByDormitory returns fans scoped to a dormitory.
func (uc *FanUseCase) ListFansByDormitory(ctx context.Context, dormitoryID uuid.UUID, page, pageSize int) (*dto.ListFansResponse, error) {
	if err := uc.ensureDormitoryExists(ctx, dormitoryID); err != nil {
		return nil, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	fans, total, err := uc.fanRepo.ListByDormitory(ctx, dormitoryID, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.buildListResponse(fans, total, page, pageSize), nil
}

// UpdateFan updates an existing fan entry.
func (uc *FanUseCase) UpdateFan(ctx context.Context, id uuid.UUID, req dto.UpdateFanRequest) (*dto.FanResponse, error) {
	fan, err := uc.fanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrFanNotFound
	}

	if req.DormitoryID != nil {
		dormitoryID, parseErr := uuid.Parse(*req.DormitoryID)
		if parseErr != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if err := uc.ensureDormitoryExists(ctx, dormitoryID); err != nil {
			return nil, err
		}
		fan.DormitoryID = dormitoryID
	}

	if req.Name != nil {
		fan.Name = *req.Name
	}
	if req.Level != nil {
		fan.Level = *req.Level
	}
	if req.Description != nil {
		fan.Description = *req.Description
	}
	fan.UpdatedAt = time.Now()

	if err := uc.fanRepo.Update(ctx, fan); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "fan", "fan:update", fan.ID.String(), map[string]string{
		"name": fan.Name,
	})

	return uc.toFanResponse(fan), nil
}

// DeleteFan removes a fan by ID.
func (uc *FanUseCase) DeleteFan(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.fanRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrFanNotFound
	}

	if err := uc.fanRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "fan", "fan:delete", id.String(), nil)
	return nil
}

func (uc *FanUseCase) toFanResponse(fan *entity.Fan) *dto.FanResponse {
	return &dto.FanResponse{
		ID:          fan.ID.String(),
		DormitoryID: fan.DormitoryID.String(),
		Name:        fan.Name,
		Level:       fan.Level,
		Description: fan.Description,
		CreatedAt:   fan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   fan.UpdatedAt.Format(time.RFC3339),
	}
}

func (uc *FanUseCase) ensureDormitoryExists(ctx context.Context, dormitoryID uuid.UUID) error {
	if _, err := uc.dormitoryRepo.GetByID(ctx, dormitoryID); err != nil {
		return domainErrors.ErrDormitoryNotFound
	}
	return nil
}

func (uc *FanUseCase) buildListResponse(fans []*entity.Fan, total int64, page, pageSize int) *dto.ListFansResponse {
	responses := make([]dto.FanResponse, 0, len(fans))
	for _, fan := range fans {
		responses = append(responses, *uc.toFanResponse(fan))
	}

	return &dto.ListFansResponse{
		Fans:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}
}
