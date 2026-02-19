package middleware

import (
	"strings"

	"arem-shop/internal/config"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

const claimsContextKey = "auth_claims"

// AuthMiddleware valide le JWT et injecte les claims dans le contexte request.
func AuthMiddleware(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authorizationHeader == "" {
			utils.JSONErrorWithCode(c, 401, "MISSING_TOKEN", "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.JSONErrorWithCode(c, 401, "INVALID_AUTH_HEADER", "invalid Authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg.JWTSecret)
		if err != nil {
			utils.JSONErrorWithCode(c, 401, "INVALID_TOKEN", "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) (*utils.Claims, bool) {
	rawClaims, exists := c.Get(claimsContextKey)
	if !exists {
		return nil, false
	}

	claims, ok := rawClaims.(*utils.Claims)
	if !ok {
		return nil, false
	}

	return claims, true
}
