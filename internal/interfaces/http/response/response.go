package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Success sends a standardized success response
func Success(c *gin.Context, statusCode int, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, statusCode int, message string, errorDetail ...string) {
	errorDetailStr := ""
	if len(errorDetail) > 0 {
		errorDetailStr = errorDetail[0]
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorDetailStr,
	})
}

// SuccessOK sends a 200 OK success response
func SuccessOK(c *gin.Context, data interface{}, message ...string) {
	Success(c, http.StatusOK, data, message...)
}

// SuccessCreated sends a 201 Created success response
func SuccessCreated(c *gin.Context, data interface{}, message ...string) {
	Success(c, http.StatusCreated, data, message...)
}

// SuccessNoContent sends a 204 No Content success response
func SuccessNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ErrorBadRequest sends a 400 Bad Request error response
func ErrorBadRequest(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Bad request"
	}
	Error(c, http.StatusBadRequest, message, errorDetail...)
}

// ErrorUnauthorized sends a 401 Unauthorized error response
func ErrorUnauthorized(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Unauthorized"
	}
	Error(c, http.StatusUnauthorized, message, errorDetail...)
}

// ErrorForbidden sends a 403 Forbidden error response
func ErrorForbidden(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Forbidden"
	}
	Error(c, http.StatusForbidden, message, errorDetail...)
}

// ErrorNotFound sends a 404 Not Found error response
func ErrorNotFound(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Not found"
	}
	Error(c, http.StatusNotFound, message, errorDetail...)
}

// ErrorConflict sends a 409 Conflict error response
func ErrorConflict(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Conflict"
	}
	Error(c, http.StatusConflict, message, errorDetail...)
}

// ErrorInternalServer sends a 500 Internal Server Error response
func ErrorInternalServer(c *gin.Context, message string, errorDetail ...string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message, errorDetail...)
}
