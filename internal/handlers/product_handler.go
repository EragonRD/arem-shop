package handlers

import (
	"errors"
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/middleware"
	"arem-shop/internal/models"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
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
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list products")
		}
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	shopID, role, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	created, err := h.productService.Create(c.Request.Context(), shopID, role, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		case errors.Is(err, services.ErrInvalidProductName):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_NAME", err.Error())
		case errors.Is(err, services.ErrInvalidProductCategory):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_CATEGORY", err.Error())
		case errors.Is(err, services.ErrPurchasePriceRequired):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "PURCHASE_PRICE_REQUIRED", err.Error())
		case errors.Is(err, services.ErrPurchasePriceNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "PURCHASE_PRICE_NEGATIVE", err.Error())
		case errors.Is(err, services.ErrSellingPriceNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "SELLING_PRICE_NEGATIVE", err.Error())
		case errors.Is(err, services.ErrStockNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "STOCK_NEGATIVE", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create product")
		}
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ProductHandler) Update(c *gin.Context) {
	shopID, role, ok := getTenantAuthContext(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	updated, err := h.productService.Update(c.Request.Context(), shopID, c.Param("id"), role, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		case errors.Is(err, services.ErrInvalidProductID):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_PRODUCT_ID", err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			utils.JSONErrorWithCode(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrInvalidProductName):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_NAME", err.Error())
		case errors.Is(err, services.ErrInvalidProductCategory):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_CATEGORY", err.Error())
		case errors.Is(err, services.ErrPurchasePriceNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "PURCHASE_PRICE_NEGATIVE", err.Error())
		case errors.Is(err, services.ErrSellingPriceNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "SELLING_PRICE_NEGATIVE", err.Error())
		case errors.Is(err, services.ErrStockNegative):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "STOCK_NEGATIVE", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product")
		}
		return
	}

	c.JSON(http.StatusOK, updated)
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
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		case errors.Is(err, services.ErrInvalidProductID):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_PRODUCT_ID", err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			utils.JSONErrorWithCode(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete product")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func getTenantAuthContext(c *gin.Context) (string, models.UserRole, bool) {
	shopID, ok := middleware.GetShopID(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_SHOP_CONTEXT", "shop context not found")
		return "", "", false
	}

	claims, ok := middleware.GetClaims(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_CLAIMS", "authentication claims not found")
		return "", "", false
	}

	return shopID, models.UserRole(claims.Role), true
}
