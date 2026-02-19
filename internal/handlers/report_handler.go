package handlers

import (
	"errors"
	"net/http"

	"arem-shop/internal/middleware"
	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) Dashboard(c *gin.Context) {
	shopID, ok := middleware.GetShopID(c)
	if !ok {
		utils.JSONErrorWithCode(c, http.StatusUnauthorized, "MISSING_SHOP_CONTEXT", "shop context not found")
		return
	}

	resp, err := h.reportService.Dashboard(c.Request.Context(), shopID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidShopID):
			utils.JSONErrorWithCode(c, http.StatusUnauthorized, "INVALID_SHOP_ID", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to build dashboard")
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
