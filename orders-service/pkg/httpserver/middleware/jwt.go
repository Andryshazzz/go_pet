package httpmiddleware

import (
	"net/http"
	"strings"

	"github.com/Andryshazzz/go_pet/orders-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// JWTAuth returns a middleware that validates JWT tokens.
// Sets "user_id" in the Gin context on success.
func JWTAuth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing authorization header",
				"message": "unauthorized",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid authorization format",
				"message": "unauthorized",
			})
			return
		}

		claims := &jwt.Claims{}
		if err := jwtService.ValidateToken(parts[1], claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid token",
				"message": "unauthorized",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
