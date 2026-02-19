//-----OLD-----
/*
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
*/

// -----NEW-----
package handlers

import (
	"net/http"

	"arem-shop/internal/dto"
	"arem-shop/internal/services"

	"github.com/gin-gonic/gin"
)

// TransactionHandler définit les routes HTTP pour les transactions
type TransactionHandler interface {
	CreateTransaction(c *gin.Context)
}

type transactionHandler struct {
	transactionService services.TransactionService
}

// NewTransactionHandler instancie le contrôleur
func NewTransactionHandler(transactionService services.TransactionService) TransactionHandler {
	return &transactionHandler{transactionService: transactionService}
}

//-----OLD-----
/*
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
*/

// -----NEW-----
// CreateTransaction gère la route POST /transactions
func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	// 1. SÉCURITÉ CRITIQUE : Extraction du shopID depuis le contexte Gin (injecté par le JWT)
	// On ne fait JAMAIS confiance au body de la requête pour le shopID
	shopIDAny, exists := c.Get("shopID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized, missing shop context"})
		return
	}
	shopID := shopIDAny.(string)

	// 2. Validation du payload utilisateur avec notre DTO
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body: " + err.Error()})
		return
	}

	// 3. Appel de la logique métier (le Service)
	res, err := h.transactionService.Create(c.Request.Context(), shopID, req)
	if err != nil {
		// S'il y a une erreur (ex: stock insuffisant), on renvoie une 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// 4. Réponse de succès formatée selon les standards du projet
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}
