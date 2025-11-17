package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// UserHandler handles user management requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.userUseCase.CreateUser(c.Request.Context(), req)
	if err != nil {
		switch err {
		case domainErrors.ErrUserAlreadyExists:
			response.ErrorConflict(c, "User already exists")
		default:
			response.ErrorInternalServer(c, "Failed to create user", err.Error())
		}
		return
	}

	response.SuccessCreated(c, resp, "User created successfully")
}

// GetUser handles getting a user by ID
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	resp, err := h.userUseCase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "User not found")
		default:
			response.ErrorInternalServer(c, "Failed to get user", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "User retrieved successfully")
}

// UpdateUser handles user update
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.userUseCase.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "User not found")
		case domainErrors.ErrUserAlreadyExists:
			response.ErrorConflict(c, "Email already taken")
		default:
			response.ErrorInternalServer(c, "Failed to update user", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "User updated successfully")
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete a user (soft delete)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	err = h.userUseCase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "User not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete user", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} map[string]string
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.userUseCase.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list users", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Users retrieved successfully")
}

// AssignRoleToUser handles assigning a role to a user
func (h *UserHandler) AssignRoleToUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	var req dto.AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	err = h.userUseCase.AssignRoleToUser(c.Request.Context(), userID, roleID)
	if err != nil {
		switch err {
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "User not found")
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		default:
			response.ErrorInternalServer(c, "Failed to assign role", err.Error())
		}
		return
	}

	response.SuccessOK(c, nil, "Role assigned successfully")
}

// RemoveRoleFromUser handles removing a role from a user
func (h *UserHandler) RemoveRoleFromUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	err = h.userUseCase.RemoveRoleFromUser(c.Request.Context(), userID, roleID)
	if err != nil {
		switch err {
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "User not found")
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		default:
			response.ErrorInternalServer(c, "Failed to remove role", err.Error())
		}
		return
	}

	response.SuccessOK(c, nil, "Role removed successfully")
}
