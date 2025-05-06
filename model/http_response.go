package model

import "time"

// ErrorResponse represents the structure of an error response.
type HttpResponse struct {
	Message   string    `json:"message"`   // A user-friendly error message
	Error     any       `json:"error"`     // The actual error message (optional)
	Path      string    `json:"path"`      // The request path that caused the error (optional)
	Status    int       `json:"status"`    // HTTP status code (optional)
	Data      any       `json:"data"`      // Additional data related to the error (optional)
	Timestamp time.Time `json:"timestamp"` // The timestamp when the error occurred (optional)
}
