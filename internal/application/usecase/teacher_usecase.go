package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// TeacherUseCase orchestrates teacher operations.
type TeacherUseCase struct {
	teacherRepo repository.TeacherRepository
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	auditLogger appService.AuditLogger
}

// NewTeacherUseCase constructs TeacherUseCase.
func NewTeacherUseCase(
	teacherRepo repository.TeacherRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	auditLogger appService.AuditLogger,
) *TeacherUseCase {
	return &TeacherUseCase{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		auditLogger: auditLogger,
	}
}

// CreateTeacher creates teacher and links/creates corresponding user.
func (uc *TeacherUseCase) CreateTeacher(ctx context.Context, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error) {
	if existing, _ := uc.teacherRepo.GetByCode(ctx, req.TeacherCode); existing != nil {
		return nil, domainErrors.ErrTeacherAlreadyExists
	}

	var user *entity.User
	var err error

	if req.ExistingUsername != "" {
		user, err = uc.userRepo.GetByUsername(ctx, strings.ToLower(req.ExistingUsername))
		if err != nil || user == nil {
			return nil, domainErrors.ErrUserNotFound
		}
		if linked, _ := uc.teacherRepo.GetByUserID(ctx, user.ID); linked != nil {
			return nil, domainErrors.ErrTeacherUserAssigned
		}
		if err := uc.ensureTeacherRole(ctx, user.ID); err != nil {
			return nil, err
		}
	} else {
		username := uc.deriveUsername(ctx, req.FullName)
		user = &entity.User{
			ID:        uuid.New(),
			Username:  username,
			Password:  "ppdf2025",
			Name:      req.FullName,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := user.HashPassword(); err != nil {
			return nil, domainErrors.ErrInternalServer
		}
		teacherRole, err := uc.roleRepo.GetBySlug(ctx, "teacher")
		if err != nil || teacherRole == nil {
			return nil, domainErrors.ErrRoleNotFound
		}
		user.Roles = []entity.Role{*teacherRole}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, domainErrors.ErrInternalServer
		}
	}

	teacher := &entity.Teacher{
		ID:               uuid.New(),
		UserID:           &user.ID,
		TeacherCode:      req.TeacherCode,
		FullName:         req.FullName,
		Gender:           req.Gender,
		Phone:            req.Phone,
		Email:            req.Email,
		Specialization:   req.Specialization,
		EmploymentStatus: req.EmploymentStatus,
		JoinedAt:         req.JoinedAt,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := uc.teacherRepo.Create(ctx, teacher); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "teacher", "teachers:create", teacher.ID.String(), map[string]string{
		"teacher_code": teacher.TeacherCode,
	})

	return uc.toTeacherResponse(teacher, user), nil
}

// GetTeacher retrieves teacher by id.
func (uc *TeacherUseCase) GetTeacher(ctx context.Context, id uuid.UUID) (*dto.TeacherResponse, error) {
	teacher, err := uc.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrTeacherNotFound
	}

	var user *entity.User
	if teacher.UserID != nil {
		user, _ = uc.userRepo.GetByID(ctx, *teacher.UserID)
	}

	return uc.toTeacherResponse(teacher, user), nil
}

// ListTeachers lists teachers with pagination/filter.
func (uc *TeacherUseCase) ListTeachers(ctx context.Context, page, pageSize int, keyword string, isActive *bool) (*dto.ListTeachersResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	filter := repository.TeacherFilter{Keyword: keyword, IsActive: isActive, Page: page, PageSize: pageSize}
	teachers, total, err := uc.teacherRepo.List(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.TeacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		var user *entity.User
		if teacher.UserID != nil {
			user, _ = uc.userRepo.GetByID(ctx, *teacher.UserID)
		}
		responses = append(responses, *uc.toTeacherResponse(teacher, user))
	}

	return &dto.ListTeachersResponse{
		Teachers:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateTeacher updates teacher profile.
func (uc *TeacherUseCase) UpdateTeacher(ctx context.Context, id uuid.UUID, req dto.UpdateTeacherRequest) (*dto.TeacherResponse, error) {
	teacher, err := uc.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrTeacherNotFound
	}

	if req.FullName != nil {
		teacher.FullName = *req.FullName
	}
	if req.Gender != nil {
		teacher.Gender = *req.Gender
	}
	if req.Phone != nil {
		teacher.Phone = *req.Phone
	}
	if req.Email != nil {
		teacher.Email = *req.Email
	}
	if req.Specialization != nil {
		teacher.Specialization = *req.Specialization
	}
	if req.EmploymentStatus != nil {
		teacher.EmploymentStatus = *req.EmploymentStatus
	}
	if req.JoinedAt != nil {
		teacher.JoinedAt = req.JoinedAt
	}
	if req.IsActive != nil {
		teacher.IsActive = *req.IsActive
	}
	teacher.UpdatedAt = time.Now()

	if err := uc.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	var user *entity.User
	if teacher.UserID != nil {
		user, _ = uc.userRepo.GetByID(ctx, *teacher.UserID)
		if user != nil {
			updated := false
			if req.FullName != nil && user.Name != *req.FullName {
				user.Name = *req.FullName
				updated = true
			}
			if req.Email != nil && user.Username == strings.ToLower(user.Username) {
				// keep username stable; only sync email not stored on user entity currently
			}
			if updated {
				user.UpdatedAt = time.Now()
				_ = uc.userRepo.Update(ctx, user)
			}
		}
	}

	_ = uc.auditLogger.Log(ctx, "teacher", "teachers:update", teacher.ID.String(), map[string]string{
		"teacher_code": teacher.TeacherCode,
	})

	return uc.toTeacherResponse(teacher, user), nil
}

// DeactivateTeacher soft-deletes teacher.
func (uc *TeacherUseCase) DeactivateTeacher(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.teacherRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrTeacherNotFound
	}
	if err := uc.teacherRepo.SoftDelete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}
	_ = uc.auditLogger.Log(ctx, "teacher", "teachers:deactivate", id.String(), nil)
	return nil
}

func (uc *TeacherUseCase) ensureTeacherRole(ctx context.Context, userID uuid.UUID) error {
	teacherRole, err := uc.roleRepo.GetBySlug(ctx, "teacher")
	if err != nil || teacherRole == nil {
		return domainErrors.ErrRoleNotFound
	}
	return uc.userRepo.AssignRole(ctx, userID, teacherRole.ID)
}

func (uc *TeacherUseCase) deriveUsername(ctx context.Context, fullName string) string {
	base := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(fullName), " ", ""))
	if base == "" {
		base = fmt.Sprintf("teacher%s", uuid.New().String()[:8])
	}

	username := base
	counter := 1
	for {
		if existing, _ := uc.userRepo.GetByUsername(ctx, username); existing == nil {
			return username
		}
		username = fmt.Sprintf("%s%d", base, counter)
		counter++
	}
}

func (uc *TeacherUseCase) toTeacherResponse(teacher *entity.Teacher, user *entity.User) *dto.TeacherResponse {
	resp := &dto.TeacherResponse{
		ID:               teacher.ID.String(),
		TeacherCode:      teacher.TeacherCode,
		FullName:         teacher.FullName,
		Gender:           teacher.Gender,
		Phone:            teacher.Phone,
		Email:            teacher.Email,
		Specialization:   teacher.Specialization,
		EmploymentStatus: teacher.EmploymentStatus,
		IsActive:         teacher.IsActive,
		CreatedAt:        teacher.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        teacher.UpdatedAt.Format(time.RFC3339),
	}
	if teacher.JoinedAt != nil {
		resp.JoinedAt = teacher.JoinedAt.Format(time.RFC3339)
	}
	if user != nil {
		resp.UserID = user.ID.String()
		resp.Username = user.Username
	}
	return resp
}
