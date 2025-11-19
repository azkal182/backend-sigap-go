package router

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/handler"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/middleware"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// SetupRouter configures all routes
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	dormitoryHandler *handler.DormitoryHandler,
	studentHandler *handler.StudentHandler,
	roleHandler *handler.RoleHandler,
	locationHandler *handler.LocationHandler,
	permissionHandler *handler.PermissionHandler,
	auditLogHandler *handler.AuditLogHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// Global CORS middleware so all routes are covered
	router.Use(middleware.NewCORSMiddlewareFromEnv())
	// Audit context middleware to enrich context for audit logging
	router.Use(middleware.AuditContextMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.SuccessOK(c, gin.H{"status": "ok"}, "Service is healthy")
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Public location routes (no auth)
		api.GET("/provinces", locationHandler.ListProvinces)
		api.GET("/provinces/:id", locationHandler.GetProvince)
		api.GET("/regencies", locationHandler.ListRegencies)
		api.GET("/regencies/:id", locationHandler.GetRegency)
		api.GET("/districts", locationHandler.ListDistricts)
		api.GET("/districts/:id", locationHandler.GetDistrict)
		api.GET("/villages", locationHandler.ListVillages)
		api.GET("/villages/:id", locationHandler.GetVillage)

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// Current user
			protected.GET("/me", userHandler.Me)

			// Audit log routes (read-only)
			auditLogs := protected.Group("/audit-logs")
			{
				auditLogs.GET("", authMiddleware.RequirePermission("audit:read"), auditLogHandler.ListAuditLogs)
			}

			// Student routes
			students := protected.Group("/students")
			{
				students.GET("", authMiddleware.RequirePermission("student:read"), studentHandler.ListStudents)
				students.GET("/:id", authMiddleware.RequirePermission("student:read"), studentHandler.GetStudent)
				students.POST("", authMiddleware.RequirePermission("student:create"), studentHandler.CreateStudent)
				students.PUT("/:id", authMiddleware.RequirePermission("student:update"), studentHandler.UpdateStudent)
				students.PATCH("/:id/status", authMiddleware.RequirePermission("student:update"), studentHandler.UpdateStudentStatus)
				students.POST("/:id/mutate-dormitory", authMiddleware.RequirePermission("student:update"), studentHandler.MutateStudentDormitory)
			}

			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("", authMiddleware.RequirePermission("user:create"), userHandler.CreateUser)
				users.PUT("/:id", authMiddleware.RequirePermission("user:update"), userHandler.UpdateUser)
				users.DELETE("/:id", authMiddleware.RequirePermission("user:delete"), userHandler.DeleteUser)
				users.POST("/:id/roles", authMiddleware.RequirePermission("user:update"), userHandler.AssignRoleToUser)
				users.DELETE("/:id/roles/:role_id", authMiddleware.RequirePermission("user:update"), userHandler.RemoveRoleFromUser)
			}

			// Dormitory routes
			dormitories := protected.Group("/dormitories")
			{
				dormitories.GET("", dormitoryHandler.ListDormitories)
				dormitories.GET("/:id", authMiddleware.RequireDormitoryAccess(), dormitoryHandler.GetDormitory)
				dormitories.POST("", authMiddleware.RequirePermission("dorm:create"), dormitoryHandler.CreateDormitory)
				dormitories.PUT("/:id", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:update"), dormitoryHandler.UpdateDormitory)
				dormitories.DELETE("/:id", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:delete"), dormitoryHandler.DeleteDormitory)
				dormitories.POST("/:id/users", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:update"), dormitoryHandler.AssignDormitoryUser)
				dormitories.DELETE("/:id/users/:user_id", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:update"), dormitoryHandler.RemoveDormitoryUser)
			}

			// Role routes
			roles := protected.Group("/roles")
			{
				roles.GET("", authMiddleware.RequirePermission("role:read"), roleHandler.ListRoles)
				roles.GET("/:id", authMiddleware.RequirePermission("role:read"), roleHandler.GetRole)
				roles.POST("", authMiddleware.RequirePermission("role:create"), roleHandler.CreateRole)
				roles.PUT("/:id", authMiddleware.RequirePermission("role:update"), roleHandler.UpdateRole)
				roles.DELETE("/:id", authMiddleware.RequirePermission("role:delete"), roleHandler.DeleteRole)
				roles.POST("/:id/permissions", authMiddleware.RequirePermission("role:update"), roleHandler.AssignPermission)
				roles.DELETE("/:id/permissions", authMiddleware.RequirePermission("role:update"), roleHandler.RemovePermission)
			}

			// Permission routes (read-only)
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", authMiddleware.RequirePermission("role:read"), permissionHandler.ListPermissions)
			}
		}
	}

	return router
}
