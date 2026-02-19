package repository

import (
	"context"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByIDAndShopID(ctx context.Context, productID, shopID uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND shop_id = ?", productID, shopID).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) Save(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) FindByIDAndShopIDForUpdate(ctx context.Context, tx *gorm.DB, productID, shopID uuid.UUID) (*models.Product, error) {
	execDB := tx
	if execDB == nil {
		execDB = r.db
	}

	var product models.Product
	err := execDB.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND shop_id = ?", productID, shopID).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) SaveWithTx(ctx context.Context, tx *gorm.DB, product *models.Product) error {
	execDB := tx
	if execDB == nil {
		execDB = r.db
	}

	return execDB.WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) CountLowStockByShopID(ctx context.Context, shopID uuid.UUID, threshold int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("shop_id = ? AND stock <= ?", shopID, threshold).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ProductRepository) DeleteByIDAndShopID(ctx context.Context, productID, shopID uuid.UUID) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND shop_id = ?", productID, shopID).
		Delete(&models.Product{})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
