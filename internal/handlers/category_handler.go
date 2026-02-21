package handlers

import (
	"context"
	"errors"
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/services"

	"github.com/gin-gonic/gin"
)

type categoryHandlerService interface {
	List(ctx context.Context, shopID string) ([]dto.CategoryResponse, error)
}

type CategoryHandler struct {
	categoryService categoryHandlerService
}

func NewCategoryHandler(categoryService categoryHandlerService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(c *gin.Context) {
	shopID, _, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	categories, err := h.categoryService.List(c.Request.Context(), shopID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			respondProductError(c, http.StatusUnauthorized, err.Error())
		default:
			respondProductError(c, http.StatusInternalServerError, "failed to list categories")
		}
		return
	}

	respondProductSuccess(c, http.StatusOK, categories)
}
