package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
)

// responseWriterWrapper wraps gin.ResponseWriter to capture status code
type responseWriterWrapper struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	// Ensure status code is set even if WriteHeader not called explicitly
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// AuditContextMiddleware injects HTTP request context info into the context for audit logging
func AuditContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap ResponseWriter to capture status code
		w := &responseWriterWrapper{ResponseWriter: c.Writer}
		c.Writer = w

		// Inject request info into context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, appService.CtxKeyRequestPath, c.FullPath())
		ctx = context.WithValue(ctx, appService.CtxKeyRequestMethod, c.Request.Method)
		ctx = context.WithValue(ctx, appService.CtxKeyIPAddress, c.ClientIP())
		ctx = context.WithValue(ctx, appService.CtxKeyUserAgent, c.Request.UserAgent())
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// After handler, set status code in context
		ctx = c.Request.Context()
		ctx = context.WithValue(ctx, appService.CtxKeyStatusCode, w.statusCode)
		c.Request = c.Request.WithContext(ctx)
	}
}
