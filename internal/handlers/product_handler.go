package handlers

import (
	"context"
	"errors"
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/middleware"
	"arem-shop/internal/models"
	"arem-shop/internal/services"

	"github.com/gin-gonic/gin"
)

type productHandlerService interface {
	List(ctx context.Context, shopID string, role models.UserRole) (interface{}, error)
	Create(ctx context.Context, shopID string, role models.UserRole, req dto.CreateProductRequest) (interface{}, error)
	Update(ctx context.Context, shopID, productID string, role models.UserRole, req dto.UpdateProductRequest) (interface{}, error)
	Delete(ctx context.Context, shopID, productID string) error
}

type ProductHandler struct {
	productService productHandlerService
}

func NewProductHandler(productService productHandlerService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) List(c *gin.Context) {
	shopID, role, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	products, err := h.productService.List(c.Request.Context(), shopID, role)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			respondProductError(c, http.StatusUnauthorized, err.Error())
		default:
			respondProductError(c, http.StatusInternalServerError, "failed to list products")
		}
		return
	}

	respondProductSuccess(c, http.StatusOK, products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	shopID, role, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondProductError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	created, err := h.productService.Create(c.Request.Context(), shopID, role, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			respondProductError(c, http.StatusUnauthorized, err.Error())
		case errors.Is(err, services.ErrInvalidProductName):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrInvalidProductCategory):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPurchasePriceRequired):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPurchasePriceNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPurchasePriceForbidden):
			respondProductError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, services.ErrSellingPriceNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrStockNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		default:
			respondProductError(c, http.StatusInternalServerError, "failed to create product")
		}
		return
	}

	respondProductSuccess(c, http.StatusCreated, created)
}

func (h *ProductHandler) Update(c *gin.Context) {
	shopID, role, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondProductError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	updated, err := h.productService.Update(c.Request.Context(), shopID, c.Param("id"), role, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			respondProductError(c, http.StatusUnauthorized, err.Error())
		case errors.Is(err, services.ErrInvalidProductID):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			respondProductError(c, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrInvalidProductName):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrInvalidProductCategory):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPurchasePriceNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPurchasePriceForbidden):
			respondProductError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, services.ErrSellingPriceNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrStockNegative):
			respondProductError(c, http.StatusBadRequest, err.Error())
		default:
			respondProductError(c, http.StatusInternalServerError, "failed to update product")
		}
		return
	}

	respondProductSuccess(c, http.StatusOK, updated)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	shopID, _, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	err := h.productService.Delete(c.Request.Context(), shopID, c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			respondProductError(c, http.StatusUnauthorized, err.Error())
		case errors.Is(err, services.ErrInvalidProductID):
			respondProductError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			respondProductError(c, http.StatusNotFound, err.Error())
		default:
			respondProductError(c, http.StatusInternalServerError, "failed to delete product")
		}
		return
	}

	respondProductSuccess(c, http.StatusOK, gin.H{
		"message": "product deleted",
	})
}

func getTenantAuthContext(c *gin.Context) (string, models.UserRole, bool) {
	shopID, ok := middleware.GetShopID(c)
	if !ok {
		respondProductError(c, http.StatusUnauthorized, "shop context not found")
		return "", "", false
	}

	claims, ok := middleware.GetClaims(c)
	if !ok {
		respondProductError(c, http.StatusUnauthorized, "authentication claims not found")
		return "", "", false
	}

	return shopID, models.UserRole(claims.Role), true
}

type productSuccessEnvelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type productErrorEnvelope struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func respondProductSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, productSuccessEnvelope{
		Success: true,
		Data:    data,
	})
}

func respondProductError(c *gin.Context, status int, message string) {
	c.JSON(status, productErrorEnvelope{
		Success: false,
		Error:   message,
	})
}
