package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/permissions"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(
			c.GetHeader("Authorization"),
		)

		parts := strings.Fields(header)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header",
				},
			)

			return
		}

		token := parts[1]

		claims, err := ValidateJWT(
			token,
			secret,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)

			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if !permissions.HasPermission(
			role,
			permission,
		) {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "permission denied",
				},
			)

			return
		}

		c.Next()
	}
}
