package services

import (
	"context"
	"errors"
	"testing"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type fakeTransactionProductRepo struct {
	product *models.Product
}

func (r *fakeTransactionProductRepo) FindByIDAndShopIDForUpdate(_ context.Context, _ *gorm.DB, productID, shopID uuid.UUID) (*models.Product, error) {
	if r.product == nil || r.product.ID != productID || r.product.ShopID != shopID {
		return nil, gorm.ErrRecordNotFound
	}
	return r.product, nil
}

func (r *fakeTransactionProductRepo) SaveWithTx(_ context.Context, _ *gorm.DB, product *models.Product) error {
	r.product = product
	return nil
}

type fakeTransactionRepo struct {
	created []models.Transaction
}

func (r *fakeTransactionRepo) Create(_ context.Context, _ *gorm.DB, transaction *models.Transaction) error {
	copyTx := *transaction
	r.created = append(r.created, copyTx)
	return nil
}

func TestTransactionService_CreateSale_DecrementsStockAndCreatesTransaction(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	quantity := 2
	productIDStr := productID.String()

	productRepo := &fakeTransactionProductRepo{
		product: &models.Product{
			ID:     productID,
			ShopID: shopID,
			Stock:  5,
		},
	}
	txRepo := &fakeTransactionRepo{}
	service := NewTransactionService(nil, productRepo, txRepo)

	req := dto.CreateTransactionRequest{
		Type:      "Sale",
		ProductID: &productIDStr,
		Quantity:  &quantity,
		Amount:    decimal.RequireFromString("199.99"),
	}

	resp, err := service.Create(context.Background(), shopID.String(), req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if resp.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", resp.Quantity)
	}
	if productRepo.product.Stock != 3 {
		t.Fatalf("expected stock 3, got %d", productRepo.product.Stock)
	}
	if len(txRepo.created) != 1 {
		t.Fatalf("expected 1 transaction created, got %d", len(txRepo.created))
	}
}

func TestTransactionService_CreateSale_RejectsInsufficientStock(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	quantity := 4
	productIDStr := productID.String()

	productRepo := &fakeTransactionProductRepo{
		product: &models.Product{
			ID:     productID,
			ShopID: shopID,
			Stock:  3,
		},
	}
	txRepo := &fakeTransactionRepo{}
	service := NewTransactionService(nil, productRepo, txRepo)

	req := dto.CreateTransactionRequest{
		Type:      "Sale",
		ProductID: &productIDStr,
		Quantity:  &quantity,
		Amount:    decimal.RequireFromString("10"),
	}

	_, err := service.Create(context.Background(), shopID.String(), req)
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	if len(txRepo.created) != 0 {
		t.Fatalf("expected no transaction creation when stock is insufficient")
	}
	if productRepo.product.Stock != 3 {
		t.Fatalf("expected stock unchanged at 3, got %d", productRepo.product.Stock)
	}
}

func TestTransactionService_ExpenseRejectsProductID(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New().String()
	quantity := 0

	service := NewTransactionService(nil, &fakeTransactionProductRepo{}, &fakeTransactionRepo{})

	req := dto.CreateTransactionRequest{
		Type:      "Expense",
		ProductID: &productID,
		Quantity:  &quantity,
		Amount:    decimal.RequireFromString("50.00"),
	}

	_, err := service.Create(context.Background(), shopID.String(), req)
	if !errors.Is(err, ErrTransactionProductForbidden) {
		t.Fatalf("expected ErrTransactionProductForbidden, got %v", err)
	}
}
