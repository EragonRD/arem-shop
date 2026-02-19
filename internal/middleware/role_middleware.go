package middleware

import (
	"arem-shop/internal/models"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware filtre l'acces selon les roles autorises.
func RoleMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
	allowed := make(map[models.UserRole]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			utils.JSONErrorWithCode(c, 401, "MISSING_CLAIMS", "authentication claims not found")
			c.Abort()
			return
		}

		role := models.UserRole(claims.Role)
		if _, exists := allowed[role]; !exists {
			utils.JSONErrorWithCode(c, 403, "FORBIDDEN", "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
