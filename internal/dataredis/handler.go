package dataredis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
)

// This struct defines the DataRedisHandler which handles HTTP requests related to Redis data.
// It contains a service field of type DataRedisService which is used to interact with the Redis data layer.
type DataRedisHandler struct {
	Service DataRedisService
}

// NewDataRedisHandler creates a new instance of DataRedisHandler.
// It initializes the DataRedisHandler struct with the provided DataRedisService.
func NewDataRedisHandler(dataRedisService DataRedisService) *DataRedisHandler {
	return &DataRedisHandler{Service: dataRedisService}
}

// GetStringValue retrieves a string value from Redis by its key and returns it as JSON.
// @Summary      Get string value from Redis
// @Description  Get a string value from Redis by its key
// @Tags         dataredis
// @Accept       json
// @Produce      json
// @Param        key   path      string  true  "Redis key"
// @Success      200  {object}  HttpResponse for successful retrieval
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /dataredis/string/{key} [get]
func (h *DataRedisHandler) GetStringValue(c *gin.Context) {
	// Parse the key from the URL parameter
	key := c.Param("key")
	if key == "" {
		util.JSONError(c, http.StatusBadRequest, "Invalid key", "Key cannot be empty")
		return
	}

	// Call the service to get the string value from Redis
	value, err := h.Service.GetStringValue(c.Request.Context(), key)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to get string value", err.Error())
		return
	}

	// Check if the value is empty
	if value == "" {
		util.JSONError(c, http.StatusNotFound, "Value not found", "Value is empty")
		return
	}

	// Return the string value as JSON
	util.JSONSuccess(c, http.StatusOK, "String value retrieved successfully", value)
}

// GetJSONValue retrieves a JSON value from Redis by its key and returns it as JSON.
// @Summary      Get JSON value from Redis
// @Description  Get a JSON value from Redis by its key
// @Tags         dataredis
// @Accept       json
// @Produce      json
// @Param        key   path      string  true  "Redis key"
// @Success      200  {object}  HttpResponse for successful retrieval
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /dataredis/json/{key} [get]
func (h *DataRedisHandler) GetJSONValue(c *gin.Context) {
	// Parse the key from the URL parameter
	key := c.Param("key")
	if key == "" {
		util.JSONError(c, http.StatusBadRequest, "Invalid key", "Key cannot be empty")
		return
	}

	// Call the service to get the JSON value from Redis
	value, err := h.Service.GetJSONValue(c.Request.Context(), key)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to get JSON value", err.Error())
		return
	}

	// Check if the value is empty
	if value == nil {
		util.JSONError(c, http.StatusNotFound, "Value not found", "Value is empty")
		return
	}

	// Return the JSON value as JSON
	util.JSONSuccess(c, http.StatusOK, "JSON value retrieved successfully", value)
}
