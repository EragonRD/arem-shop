package handlers

import (
	"context"
	"errors"
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

type publicService interface {
	ListProductsByShopID(ctx context.Context, shopID string) ([]dto.PublicProductResponse, error)
}

type PublicHandler struct {
	publicService publicService
}

func NewPublicHandler(publicService publicService) *PublicHandler {
	return &PublicHandler{publicService: publicService}
}

func (h *PublicHandler) ListPublicProducts(c *gin.Context) {
	shopID := c.Param("shopID")

	products, err := h.publicService.ListProductsByShopID(c.Request.Context(), shopID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_SHOP_ID", err.Error())
		case errors.Is(err, services.ErrShopNotFound):
			utils.JSONErrorWithCode(c, http.StatusNotFound, "SHOP_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrShopInactive):
			utils.JSONErrorWithCode(c, http.StatusForbidden, "SHOP_INACTIVE", err.Error())
		case errors.Is(err, services.ErrInvalidShopWhatsApp):
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INVALID_SHOP_WHATSAPP", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list public products")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": products})
}
