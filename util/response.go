package util

import (
	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/model"
	"time"
)

func JSONSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, model.HttpResponse{
		Message:   message,
		Error:     nil,
		Path:      c.Request.URL.Path,
		Status:    status,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func JSONError(c *gin.Context, status int, message string, err string) {
	c.JSON(status, model.HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    status,
		Data:      nil,
		Timestamp: time.Now(),
	})
}
