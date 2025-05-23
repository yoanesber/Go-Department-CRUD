package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/metacontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
)

// RequestLogger is a middleware function that logs incoming HTTP requests.
// It initializes the logger, records the request details, and logs them after the request is processed.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(c.Request.Context())
		if !ok {
			// If metadata extraction fails, log an error and return
			logger.RequestLogger.Error("Failed to extract metadata from context")
			return
		}

		// Get the username from the context
		// This assumes that the username is set in the context by JWT validation middleware
		username := meta.UserName
		if username == "" {
			username = "unknown"
		}

		// Get the user roles from the metadata
		userRoles := meta.Roles

		// Then log the request details
		// This is done after the request is processed to capture the response status and duration
		duration := time.Since(start)
		logger.RequestLogger.WithFields(logrus.Fields{
			"content_length": c.Request.ContentLength,
			"content_type":   c.ContentType(),
			"duration":       duration.String(),
			"ip":             c.ClientIP(),
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"query":          c.Request.URL.Query(),
			"referer":        c.Request.Referer(),
			"request_id":     c.Writer.Header().Get("X-Request-Id"),
			"status":         c.Writer.Status(),
			"user_agent":     c.Request.UserAgent(),
			"username":       username,
			"roles":          userRoles,
		}).Info("Incoming request")
	}
}
