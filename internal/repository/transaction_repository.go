package repository

import (
	"context"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	execDB := tx
	if execDB == nil {
		execDB = r.db
	}

	return execDB.WithContext(ctx).Create(transaction).Error
}

func (r *TransactionRepository) SumAmountByShopAndTypes(ctx context.Context, shopID uuid.UUID, types ...models.TransactionType) (decimal.Decimal, error) {
	if len(types) == 0 {
		return decimal.Zero, nil
	}

	var totalRaw string
	err := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount)::text, '0')").
		Where("shop_id = ? AND type IN ?", shopID, types).
		Scan(&totalRaw).Error
	if err != nil {
		return decimal.Zero, err
	}

	total, err := decimal.NewFromString(totalRaw)
	if err != nil {
		return decimal.Zero, err
	}

	return total, nil
}
