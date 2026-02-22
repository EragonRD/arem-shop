package handlers

import (
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/middleware"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShopHandler struct {
	shopService *services.ShopService
}

func NewShopHandler(shopService *services.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

func (h *ShopHandler) UpdateShopInfo(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_CLAIMS", "authentication claims not found")
		return
	}

	var req dto.ShopUpdatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	shopUUID, err := uuid.Parse(claims.ShopID)
	if err != nil {
		utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_SHOP_ID", "invalid shop ID in token")
		return
	}

	shop, err := h.shopService.UpdateShopInfo(c.Request.Context(), shopUUID, req)
	if err != nil {
		utils.JSONErrorWithCode(c, http.StatusInternalServerError, "UPDATE_FAILED", "failed to update shop")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shop,
	})
}
