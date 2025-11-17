package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// UserUseCase handles user management use cases
type UserUseCase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domainErrors.ErrUserAlreadyExists
	}

	// Create user
	user := &entity.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Assign roles if provided, otherwise assign default role
	if len(req.RoleIDs) > 0 {
		roles := make([]entity.Role, 0)
		for _, roleIDStr := range req.RoleIDs {
			roleID, err := uuid.Parse(roleIDStr)
			if err != nil {
				continue
			}
			role, err := uc.roleRepo.GetByID(ctx, roleID)
			if err != nil {
				continue
			}
			roles = append(roles, *role)
		}
		user.Roles = roles
	} else {
		// Assign default role (user role)
		defaultRole, err := uc.roleRepo.GetBySlug(ctx, "user")
		if err == nil && defaultRole != nil {
			user.Roles = []entity.Role{*defaultRole}
		}
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Get user with roles
	userWithRoles, err := uc.userRepo.GetWithRoles(ctx, user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.toUserResponse(userWithRoles), nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetWithRoles(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrUserNotFound
	}

	return uc.toUserResponse(user), nil
}

// UpdateUser updates a user
func (uc *UserUseCase) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrUserNotFound
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, domainErrors.ErrUserAlreadyExists
		}
		user.Email = req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()

	// Update roles if provided
	if len(req.RoleIDs) > 0 {
		roles := make([]entity.Role, 0)
		for _, roleIDStr := range req.RoleIDs {
			roleID, err := uuid.Parse(roleIDStr)
			if err != nil {
				continue
			}
			role, err := uc.roleRepo.GetByID(ctx, roleID)
			if err != nil {
				continue
			}
			roles = append(roles, *role)
		}
		user.Roles = roles
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Get updated user with roles
	userWithRoles, err := uc.userRepo.GetWithRoles(ctx, user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return uc.toUserResponse(userWithRoles), nil
}

// DeleteUser deletes a user (soft delete)
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return domainErrors.ErrUserNotFound
	}

	// Delete user
	return uc.userRepo.Delete(ctx, id)
}

// ListUsers retrieves a paginated list of users
func (uc *UserUseCase) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error) {
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

	users, total, err := uc.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userWithRoles, _ := uc.userRepo.GetWithRoles(ctx, user.ID)
		userResponses = append(userResponses, *uc.toUserResponse(userWithRoles))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// AssignRoleToUser assigns a role to a user
func (uc *UserUseCase) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domainErrors.ErrUserNotFound
	}

	// Check if role exists
	_, err = uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return domainErrors.ErrRoleNotFound
	}

	// Assign role
	return uc.userRepo.AssignRole(ctx, userID, roleID)
}

// RemoveRoleFromUser removes a role from a user
func (uc *UserUseCase) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domainErrors.ErrUserNotFound
	}

	// Check if role exists
	_, err = uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return domainErrors.ErrRoleNotFound
	}

	// Remove role
	return uc.userRepo.RemoveRole(ctx, userID, roleID)
}

// toUserResponse converts entity.User to dto.UserResponse
func (uc *UserUseCase) toUserResponse(user *entity.User) *dto.UserResponse {
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		Roles:     roles,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
