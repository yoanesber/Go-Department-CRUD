package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/metacontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
)

// RoleBasedAccessControl is a middleware function that checks if the user has the required roles to access a specific route.
// It retrieves the user roles from the context and compares them with the allowed roles.
func RoleBasedAccessControl(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no allowed roles are provided, allow access
		if len(allowedRoles) == 0 {
			c.Next()
			return
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(c.Request.Context())
		if !ok {
			util.JSONError(c, http.StatusInternalServerError, "Failed to extract metadata", "Unable to extract user metadata from context")
			c.Abort()
			return
		}

		// Get the user roles from the metadata
		userRoles := meta.Roles
		if len(userRoles) == 0 {
			util.JSONError(c, http.StatusForbidden, "No roles found", "User does not have any roles")
			c.Abort()
			return
		}

		// Check if the user has any of the allowed roles
		// If the user has at least one allowed role, proceed to the next handler
		for _, role := range userRoles {
			for _, allowed := range allowedRoles {
				if role == allowed {
					c.Next()
					return
				}
			}
		}

		// If the user does not have any of the allowed roles, return a forbidden response
		// and abort the request
		util.JSONError(c, http.StatusForbidden, "Access denied", "User does not have the required role")
		c.Abort()
	}
}
