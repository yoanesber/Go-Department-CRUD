package headers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader is a middleware function that generates a unique request ID for each incoming request.
// It sets the request ID in the response header "X-Request-Id".
func RequestIDHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		c.Writer.Header().Set("X-Request-Id", id.String())

		c.Next()
	}
}
