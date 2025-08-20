package middleware

import (
	"net/http"
	"strings"

	"vm-chan/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(authService domain.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token - accept both "Bearer <token>" format or just "<token>"
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			// If it has the Bearer prefix, extract the token part
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 {
				logger.Warn("Invalid authorization header format")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				c.Abort()
				return
			}
			token = tokenParts[1]
		}

		user, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Warn("Invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Next()
	}
}
