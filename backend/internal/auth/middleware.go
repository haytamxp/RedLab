package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/permissions"
)

func JWTMiddleware(secret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "missing authorization header"},
			)

			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		claims, err := ValidateJWT(token, secret)

		if err != nil {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token"},
			)

			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
func RequirePermission(permission string) gin.HandlerFunc {

	return func(c *gin.Context) {

		role := c.GetString("role")

		if !permissions.HasPermission(role, permission) {

			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
			})

			c.Abort()

			return
		}

		c.Next()
	}
}