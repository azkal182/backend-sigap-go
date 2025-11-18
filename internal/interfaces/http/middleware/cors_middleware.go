package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewCORSMiddlewareFromEnv creates a CORS middleware using CORS_ALLOWED_ORIGINS env var.
//
// CORS_ALLOWED_ORIGINS should be a comma-separated list, for example:
//
//	CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com
//
// If CORS_ALLOWED_ORIGINS is empty, the middleware will allow all origins ("*")
// which is convenient for local development but should be configured properly
// in production.
func NewCORSMiddlewareFromEnv() gin.HandlerFunc {
	value := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowed []string
	if value != "" {
		parts := strings.Split(value, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				allowed = append(allowed, trimmed)
			}
		}
	}

	// If no origins configured, allow all (development-friendly default)
	allowAll := len(allowed) == 0

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Determine if origin is allowed
		allowedOrigin := ""
		if allowAll {
			allowedOrigin = "*"
		} else if origin != "" {
			for _, o := range allowed {
				if o == origin {
					allowedOrigin = origin
					break
				}
			}
		}

		// If we have a match (or allowAll), set CORS headers
		if allowedOrigin != "" {
			log.Printf("CORS middleware: origin=%q allowedOrigin=%q", origin, allowedOrigin)
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
