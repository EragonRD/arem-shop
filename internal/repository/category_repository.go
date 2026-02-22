package repository

import (
	"context"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) ListByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("name ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
