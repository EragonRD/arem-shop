package repository

import (
	"context"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

func (r *ShopRepository) FindByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error) {
	var shop models.Shop
	if err := r.db.WithContext(ctx).Where("id = ?", shopID).First(&shop).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}
