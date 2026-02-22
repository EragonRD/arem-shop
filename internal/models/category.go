package models

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a product category isolated by shop.
type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:80;not null" json:"name"`
	ShopID    uuid.UUID `gorm:"type:uuid;not null;index:idx_categories_shop_id" json:"shopID"`
	Shop      Shop      `gorm:"foreignKey:ShopID" json:"-"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
}

func (Category) TableName() string {
	return "categories"
}
