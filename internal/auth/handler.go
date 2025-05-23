package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/internal/refreshtoken"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
	"gopkg.in/go-playground/validator.v9"
)

// This struct defines the AuthHandler which handles HTTP requests related to authentication.
// It contains a service field of type AuthService which is used to interact with the authentication data layer.
type AuthHandler struct {
	Service AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
// It initializes the AuthHandler struct with the provided AuthService.
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{Service: authService}
}

// Login handles user login requests.
// It validates the request, authenticates the user, and returns a JWT token if successful.
// @Summary      User login
// @Description  User login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      Auth  true  "Login request"
// @Success      200  {object}  model.HttpResponse for successful login
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      401  {object}  model.HttpResponse for unauthorized
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Bind the request body to the LoginRequest struct
	// This struct contains the username and password fields
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Call the service to authenticate the user and get the token
	loginResp, err := h.Service.Login(c.Request.Context(), loginReq)

	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			util.JSONErrorMap(c, http.StatusBadRequest, "Failed to login", util.FormatValidationErrors(err))
			return
		}

		util.JSONError(c, http.StatusUnauthorized, "Failed to login", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusOK, "Login successful", loginResp)
}

// RefreshToken handles token refresh requests.
// It validates the request, checks the refresh token, and returns a new JWT token if successful.
// @Summary      Refresh token
// @Description  Refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshtoken.RefreshTokenRequest  true  "Refresh token request"
// @Success      200  {object}  model.HttpResponse for successful token refresh
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      401  {object}  model.HttpResponse for unauthorized
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Bind the request body to the RefreshTokenRequest struct
	// This struct contains the refresh token field
	var refreshTokenReq refreshtoken.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenReq); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Call the service to refresh the token
	refreshTokenResp, err := h.Service.RefreshToken(c.Request.Context(), refreshTokenReq)

	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			util.JSONErrorMap(c, http.StatusBadRequest, "Failed to refresh token", util.FormatValidationErrors(err))
			return
		}

		util.JSONError(c, http.StatusUnauthorized, "Failed to refresh token", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusOK, "Token refreshed successfully", refreshTokenResp)
}
