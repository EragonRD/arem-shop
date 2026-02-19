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

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	claims, ok := middleware.GetClaims(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_CLAIMS", "authentication claims not found")
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req, models.UserRole(claims.Role), claims.ShopID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedRole):
			utils.JSONErrorWithCode(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrCrossShopRegister):
			utils.JSONErrorWithCode(c, http.StatusForbidden, "CROSS_SHOP_FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrShopNotFound):
			utils.JSONErrorWithCode(c, http.StatusNotFound, "SHOP_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrShopInactive):
			utils.JSONErrorWithCode(c, http.StatusForbidden, "SHOP_INACTIVE", err.Error())
		case errors.Is(err, services.ErrEmailAlreadyUsed):
			utils.JSONErrorWithCode(c, http.StatusConflict, "EMAIL_ALREADY_USED", err.Error())
		case errors.Is(err, services.ErrInvalidRole):
			utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_ROLE", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user")
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONValidationError(c, err)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		case errors.Is(err, services.ErrShopInactive):
			utils.JSONErrorWithCode(c, http.StatusForbidden, "SHOP_INACTIVE", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to login")
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
