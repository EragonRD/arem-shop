package dto

import "github.com/shopspring/decimal"

// PublicProductResponse est exposee aux guests sans prix d'achat.
type PublicProductResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Category     string          `json:"category"`
	SellingPrice decimal.Decimal `json:"sellingPrice"`
	Stock        int             `json:"stock"`
	ImageURL     string          `json:"imageURL"`
	WhatsAppLink string          `json:"whatsappLink"`
}
