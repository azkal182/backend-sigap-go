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
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

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

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("", authMiddleware.RequirePermission("user:create"), userHandler.CreateUser)
				users.PUT("/:id", authMiddleware.RequirePermission("user:update"), userHandler.UpdateUser)
				users.DELETE("/:id", authMiddleware.RequirePermission("user:delete"), userHandler.DeleteUser)
			}

			// Dormitory routes
			dormitories := protected.Group("/dormitories")
			{
				dormitories.GET("", dormitoryHandler.ListDormitories)
				dormitories.GET("/:id", authMiddleware.RequireDormitoryAccess(), dormitoryHandler.GetDormitory)
				dormitories.POST("", authMiddleware.RequirePermission("dorm:create"), dormitoryHandler.CreateDormitory)
				dormitories.PUT("/:id", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:update"), dormitoryHandler.UpdateDormitory)
				dormitories.DELETE("/:id", authMiddleware.RequireDormitoryAccess(), authMiddleware.RequirePermission("dorm:delete"), dormitoryHandler.DeleteDormitory)
			}
		}
	}

	return router
}
