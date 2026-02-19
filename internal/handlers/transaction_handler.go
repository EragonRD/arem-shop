package handlers

import (
	"errors"
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/middleware"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	shopID, ok := middleware.GetShopID(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_SHOP_CONTEXT", "shop context not found")
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	created, err := h.transactionService.Create(c.Request.Context(), shopID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		case errors.Is(err, services.ErrInvalidProductID):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_PRODUCT_ID", err.Error())
		case errors.Is(err, services.ErrInvalidTransactionType):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_TRANSACTION_TYPE", err.Error())
		case errors.Is(err, services.ErrTransactionProductRequired):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "PRODUCT_REQUIRED", err.Error())
		case errors.Is(err, services.ErrTransactionProductForbidden):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "PRODUCT_FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrTransactionQuantityRequired):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "QUANTITY_REQUIRED", err.Error())
		case errors.Is(err, services.ErrTransactionQuantityInvalid):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "QUANTITY_INVALID", err.Error())
		case errors.Is(err, services.ErrTransactionAmountInvalid):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "AMOUNT_INVALID", err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			utils.JSONErrorWithCode(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrInsufficientStock):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INSUFFICIENT_STOCK", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create transaction")
		}
		return
	}

	c.JSON(http.StatusCreated, created)
}
