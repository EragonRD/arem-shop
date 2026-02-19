package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product est strictement isolé par shop et ne doit jamais fuiter cross-tenant.
type Product struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string          `gorm:"size:160;not null" json:"name"`
	Description   string          `gorm:"type:text;not null;default:''" json:"description"`
	Category      string          `gorm:"size:80;not null" json:"category"`
	PurchasePrice decimal.Decimal `gorm:"type:numeric(12,2);not null;check:purchase_price >= 0" json:"purchasePrice,omitempty"`
	SellingPrice  decimal.Decimal `gorm:"type:numeric(12,2);not null;check:selling_price >= 0" json:"sellingPrice"`
	Stock         int             `gorm:"not null;check:stock >= 0" json:"stock"`
	ImageURL      string          `gorm:"type:text;not null;default:''" json:"imageURL"`
	ShopID        uuid.UUID       `gorm:"type:uuid;not null;index:idx_products_shop_id;index:idx_products_shop_category,priority:1" json:"shopID"`
	Shop          Shop            `gorm:"foreignKey:ShopID" json:"-"`
	CreatedAt     time.Time       `gorm:"not null;default:now()" json:"createdAt"`
}

func (Product) TableName() string {
	return "products"
}
