package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreateProductRequest payload de creation produit (route privee).
type CreateProductRequest struct {
	Name          string           `json:"name" binding:"required,min=1,max=160"`
	Description   string           `json:"description" binding:"max=5000"`
	CategoryID    string           `json:"categoryID" binding:"required,uuid"`
	PurchasePrice *decimal.Decimal `json:"purchasePrice,omitempty"`
	SellingPrice  decimal.Decimal  `json:"sellingPrice" binding:"required"`
	Stock         int              `json:"stock" binding:"required,gte=0"`
	ImageURL      string           `json:"imageURL" binding:"max=2048"`
}

// UpdateProductRequest payload de mise a jour produit (route privee).
type UpdateProductRequest struct {
	Name          string           `json:"name" binding:"required,min=1,max=160"`
	Description   string           `json:"description" binding:"max=5000"`
	CategoryID    string           `json:"categoryID" binding:"required,uuid"`
	PurchasePrice *decimal.Decimal `json:"purchasePrice,omitempty"`
	SellingPrice  decimal.Decimal  `json:"sellingPrice" binding:"required"`
	Stock         int              `json:"stock" binding:"required,gte=0"`
	ImageURL      string           `json:"imageURL" binding:"max=2048"`
}

// ProductResponse representation API pour Admin (purchasePrice masque).
type ProductResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	CategoryID   string          `json:"categoryID"`
	Category     string          `json:"category"`
	SellingPrice decimal.Decimal `json:"sellingPrice"`
	Stock        int             `json:"stock"`
	ImageURL     string          `json:"imageURL"`
	ShopID       string          `json:"shopID"`
	CreatedAt    time.Time       `json:"createdAt"`
}

// ProductAdminResponse representation API pour SuperAdmin (purchasePrice visible).
type ProductAdminResponse struct {
	ProductResponse
	PurchasePrice decimal.Decimal `json:"purchasePrice"`
}
