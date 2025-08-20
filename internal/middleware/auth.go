package middleware

import (
	"net/http"
	"strings"

	"vm-chan/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(authService domain.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
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

		c.Set("user", user)
		c.Next()
	}
}
