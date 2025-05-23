package headers

import "github.com/gin-gonic/gin"

// RequestCorsHeader is a middleware function that sets CORS headers for incoming requests.
// It allows cross-origin requests from http://localhost and sets various CORS-related headers.
func RequestCorsHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", "http://localhost")
		header.Set("Access-Control-Max-Age", "86400")
		header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		header.Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		header.Set("Access-Control-Expose-Headers", "Content-Length")
		header.Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			// Handle preflight request
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}
