package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DashboardResponse est la vue financiere d'un shop (SuperAdmin only).
type DashboardResponse struct {
	TotalSales       float64 `json:"totalSales"`
	TotalExpenses    float64 `json:"totalExpenses"`
	NetProfit        float64 `json:"netProfit"`
	LowStockProducts int     `json:"lowStockProducts"`
	ShopID           string  `json:"shopID"`
}

// TopProductResponse est une sous-structure pour les produits les plus vendus
type TopProductResponse struct {
	ProductID uuid.UUID       `json:"productID"`
	Name      string          `json:"name"`
	TotalSold int             `json:"totalSold"`
	Revenue   decimal.Decimal `json:"revenue"`
}
