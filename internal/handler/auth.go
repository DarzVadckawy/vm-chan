package handler

import (
	"net/http"

	"vm-chan/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service domain.AuthService
	logger  *zap.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(service domain.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse "Login successful"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Authentication failed"
// @Failure 500 {object} domain.ErrorResponse "Server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:       "Invalid request format",
			Code:        "validation_error",
			Description: "The request body does not match the expected format",
		})
		return
	}

	response, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Error("Login failed", zap.String("username", req.Username), zap.Error(err))
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error:       "Invalid credentials",
			Code:        "authentication_failed",
			Description: "The provided username or password is incorrect",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
