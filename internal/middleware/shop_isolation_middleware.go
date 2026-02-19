package middleware

import (
	"strings"

	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const shopIDContextKey = "shop_id"

// ShopIsolationMiddleware impose l'isolation multi-tenant sur toutes les routes privees.
func ShopIsolationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			utils.JSONErrorWithCode(c, 401, "MISSING_CLAIMS", "authentication claims not found")
			c.Abort()
			return
		}

		jwtShopID := strings.TrimSpace(claims.ShopID)
		if _, err := uuid.Parse(jwtShopID); err != nil {
			utils.JSONErrorWithCode(c, 401, "INVALID_TOKEN_CLAIMS", "invalid shopID in token")
			c.Abort()
			return
		}

		// Si la route contient :shopID (ou ?shopID), il doit matcher le tenant du JWT.
		pathShopID := strings.TrimSpace(c.Param("shopID"))
		if pathShopID != "" && !strings.EqualFold(pathShopID, jwtShopID) {
			utils.JSONErrorWithCode(c, 403, "CROSS_SHOP_FORBIDDEN", "cross-shop access forbidden")
			c.Abort()
			return
		}

		queryShopID := strings.TrimSpace(c.Query("shopID"))
		if queryShopID != "" && !strings.EqualFold(queryShopID, jwtShopID) {
			utils.JSONErrorWithCode(c, 403, "CROSS_SHOP_FORBIDDEN", "cross-shop access forbidden")
			c.Abort()
			return
		}

		c.Set(shopIDContextKey, jwtShopID)
		c.Next()
	}
}

func GetShopID(c *gin.Context) (string, bool) {
	shopIDRaw, exists := c.Get(shopIDContextKey)
	if !exists {
		return "", false
	}

	shopID, ok := shopIDRaw.(string)
	if !ok || strings.TrimSpace(shopID) == "" {
		return "", false
	}

	return shopID, true
}
