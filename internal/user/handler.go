package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
	"gopkg.in/go-playground/validator.v9"
)

// This struct defines the UserHandler which handles HTTP requests related to users.
// It contains a service field of type UserService which is used to interact with the user data layer.
type UserHandler struct {
	Service UserService
}

// NewUserHandler creates a new instance of UserHandler.
// It initializes the UserHandler struct with the provided UserService.
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{Service: userService}
}

// GetAllUsers retrieves all users from the database and returns them as JSON.
// @Summary      Get all users
// @Description  Get all users from the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers(c.Request.Context())
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to retrieve users", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusOK, "All Users retrieved successfully", users)
}

// GetUserByID retrieves a user by their ID from the database and returns it as JSON.
// @Summary      Get user by ID
// @Description  Get a user by their ID from the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Parse the ID from the URL parameter
	// and convert it to an int64
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	// Retrieve the user by ID from the service
	user, err := h.Service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to retrieve user", err.Error())
		return
	}

	if (user.Equals(&User{})) {
		util.JSONError(c, http.StatusNotFound, "User not found", "No user found with the given ID")
		return
	}

	util.JSONSuccess(c, http.StatusOK, "User retrieved successfully", user)
}

// CreateUser creates a new user in the database and returns it as JSON.
// @Summary      Create user
// @Description  Create a new user in the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "User object"
// @Success      201  {object}  model.HttpResponse for successful creation
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Bind the JSON request body to the user struct
	// and validate the input using ShouldBindJSON
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Create a new user in the database
	createdUser, err := h.Service.CreateUser(c.Request.Context(), user)
	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			util.JSONErrorMap(c, http.StatusBadRequest, "Failed to create user", util.FormatValidationErrors(err))
			return
		}

		util.JSONError(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusCreated, "User created successfully", createdUser)
}
