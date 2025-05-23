package context

import (
	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/config/db/postgresdb"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
)

// PostgresDBContext is a middleware function that injects the database connection into the request context.
// It retrieves the database connection from the postgres package and sets it in the context.
// This allows the database connection to be accessed in subsequent handlers without needing to pass it explicitly.
func PostgresDBContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := dbcontext.InjectDB(c.Request.Context(), postgresdb.GetDB())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
