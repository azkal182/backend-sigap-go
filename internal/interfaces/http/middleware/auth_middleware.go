package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/domain/service"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	tokenService service.TokenService
	userRepo     repository.UserRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(
	tokenService service.TokenService,
	userRepo repository.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}

// RequireAuth is a middleware that requires valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorUnauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			if err == domainErrors.ErrTokenExpired {
				response.ErrorUnauthorized(c, "Token expired")
			} else {
				response.ErrorUnauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// Get user with roles and dormitories
		user, err := m.userRepo.GetWithRolesAndDormitories(c.Request.Context(), claims.UserID)
		if err != nil {
			response.ErrorUnauthorized(c, "User not found")
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			response.ErrorForbidden(c, "User is inactive")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("user", user)

		c.Next()
	}
}

// RequirePermission is a middleware that requires specific permission
func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check auth
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			response.ErrorUnauthorized(c, "User not found in context")
			c.Abort()
			return
		}

		userEntity, ok := user.(*entity.User)
		if !ok {
			response.ErrorInternalServer(c, "Invalid user type")
			c.Abort()
			return
		}

		// Check permission
		if !userEntity.HasPermission(permission) {
			response.ErrorForbidden(c, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireDormitoryAccess is a middleware that checks if user can access a dormitory
func (m *AuthMiddleware) RequireDormitoryAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check auth
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Get dormitory ID from URL parameter or request body
		dormitoryIDStr := c.Param("id")
		if dormitoryIDStr == "" {
			// Try to get from query
			dormitoryIDStr = c.Query("dormitory_id")
		}

		if dormitoryIDStr == "" {
			response.ErrorBadRequest(c, "Dormitory ID required")
			c.Abort()
			return
		}

		dormitoryID, err := uuid.Parse(dormitoryIDStr)
		if err != nil {
			response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
			c.Abort()
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			response.ErrorUnauthorized(c, "User not found in context")
			c.Abort()
			return
		}

		userEntity, ok := user.(*entity.User)
		if !ok {
			response.ErrorInternalServer(c, "Invalid user type")
			c.Abort()
			return
		}

		// Check if user can access this dormitory
		if !userEntity.CanAccessDormitory(dormitoryID) {
			response.ErrorForbidden(c, "Access denied to this dormitory")
			c.Abort()
			return
		}

		// Store dormitory ID in context
		c.Set("dormitory_id", dormitoryID)

		c.Next()
	}
}
