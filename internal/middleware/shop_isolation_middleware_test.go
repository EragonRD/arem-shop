package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"arem-shop/internal/config"
	"arem-shop/internal/models"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestShopIsolationMiddleware_AllowsMatchingPathShopID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.AppConfig{JWTSecret: "test-secret", JWTTTLHours: 1}
	shopID := uuid.New()
	token := mustToken(t, cfg, shopID)

	r := gin.New()
	r.Use(AuthMiddleware(cfg), ShopIsolationMiddleware())
	r.GET("/private/:shopID", func(c *gin.Context) {
		resolvedShopID, ok := GetShopID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing shop context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shopID": resolvedShopID})
	})

	req := httptest.NewRequest(http.MethodGet, "/private/"+shopID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse body: %v", err)
	}

	if body["shopID"] != shopID.String() {
		t.Fatalf("expected shopID %s, got %s", shopID.String(), body["shopID"])
	}
}

func TestShopIsolationMiddleware_RejectsCrossShopPathAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.AppConfig{JWTSecret: "test-secret", JWTTTLHours: 1}
	tokenShopID := uuid.New()
	differentShopID := uuid.New()
	token := mustToken(t, cfg, tokenShopID)

	r := gin.New()
	r.Use(AuthMiddleware(cfg), ShopIsolationMiddleware())
	r.GET("/private/:shopID", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/private/"+differentShopID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func mustToken(t *testing.T, cfg config.AppConfig, shopID uuid.UUID) string {
	t.Helper()

	user := models.User{
		ID:     uuid.New(),
		Email:  "admin@shop.com",
		Role:   models.RoleSuperAdmin,
		ShopID: shopID,
	}

	token, err := utils.GenerateToken(cfg, user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return token
}
