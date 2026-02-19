//-----OLD-----
/*
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
*/

// -----NEW-----
package handlers

import (
	"net/http"

	"arem-shop/internal/services"

	"github.com/gin-gonic/gin"
)

// ReportHandler définit les routes HTTP pour les rapports financiers
type ReportHandler interface {
	GetDashboard(c *gin.Context)
}

type reportHandler struct {
	// On utilise le pointeur ici pour s'adapter à l'implémentation de ton report_service.go
	reportService *services.ReportService
}

// NewReportHandler instancie le contrôleur
func NewReportHandler(reportService *services.ReportService) ReportHandler {
	return &reportHandler{reportService: reportService}
}

// GetDashboard gère la route GET /reports/dashboard (Réservé au SuperAdmin via le middleware)
func (h *reportHandler) GetDashboard(c *gin.Context) {
	// SÉCURITÉ : Extraction du shopID du tenant
	shopIDAny, exists := c.Get("shopID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized, missing shop context"})
		return
	}
	shopID := shopIDAny.(string)

	// Appel du service pour calculer les revenus, dépenses et bénéfices
	res, err := h.reportService.Dashboard(c.Request.Context(), shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to generate dashboard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}
