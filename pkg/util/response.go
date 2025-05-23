package util

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the structure of an error response.
type HttpResponse struct {
	Message   string    `json:"message"`   // A user-friendly error message
	Error     any       `json:"error"`     // The actual error message (optional)
	Path      string    `json:"path"`      // The request path that caused the error (optional)
	Status    int       `json:"status"`    // HTTP status code (optional)
	Data      any       `json:"data"`      // Additional data related to the error (optional)
	Timestamp time.Time `json:"timestamp"` // The timestamp when the error occurred (optional)
}

func JSONSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, HttpResponse{
		Message:   message,
		Error:     nil,
		Path:      c.Request.URL.Path,
		Status:    status,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func JSONError(c *gin.Context, status int, message string, err string) {
	c.JSON(status, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    status,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func JSONErrorMap(c *gin.Context, status int, message string, err []map[string]string) {
	c.JSON(status, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    status,
		Data:      nil,
		Timestamp: time.Now(),
	})
}
