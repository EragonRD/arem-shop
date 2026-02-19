package dto

// DashboardResponse est la vue financiere d'un shop (SuperAdmin only).
type DashboardResponse struct {
	TotalSales       float64 `json:"totalSales"`
	TotalExpenses    float64 `json:"totalExpenses"`
	NetProfit        float64 `json:"netProfit"`
	LowStockProducts int     `json:"lowStockProducts"`
	ShopID           string  `json:"shopID"`
}
