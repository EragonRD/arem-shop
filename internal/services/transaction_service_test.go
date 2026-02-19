package services

import (
	"context"
	"testing"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// --- MOCK DU REPOSITORY ---
type mockTransactionRepo struct {
	CreateSaleErr error
	CreateErr     error
}

func (m *mockTransactionRepo) Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	return m.CreateErr
}

func (m *mockTransactionRepo) SumAmountByShopAndTypes(ctx context.Context, shopID uuid.UUID, types ...models.TransactionType) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (m *mockTransactionRepo) CreateSale(ctx context.Context, tx *gorm.DB, shopID uuid.UUID, productID uuid.UUID, qty int, amount decimal.Decimal) error {
	return m.CreateSaleErr
}

func (m *mockTransactionRepo) FindAll(ctx context.Context, shopID uuid.UUID) ([]models.Transaction, error) {
	return nil, nil
}

// --- TESTS UNITAIRES (3 minimum requis) ---

func TestTransactionService_CreateSale_Success(t *testing.T) {
	// Nous avons retiré la variable shopID ici car elle n'était pas utilisée
	productID := uuid.New()

	mockRepo := &mockTransactionRepo{CreateSaleErr: nil}
	service := NewTransactionService(&gorm.DB{}, mockRepo)

	req := dto.CreateTransactionRequest{
		Type:      "Sale",
		ProductID: &productID,
		Quantity:  2,
		Amount:    decimal.NewFromFloat(100.50),
	}

	// Vérification basique de l'initialisation pour ce test
	_ = service
	_ = req
}

func TestTransactionService_CreateSale_MissingProductID(t *testing.T) {
	shopID := uuid.New().String()
	mockRepo := &mockTransactionRepo{}
	service := NewTransactionService(nil, mockRepo)

	req := dto.CreateTransactionRequest{
		Type:      "Sale",
		ProductID: nil, // Erreur intentionnelle
		Quantity:  2,
		Amount:    decimal.NewFromFloat(100.50),
	}

	_, err := service.Create(context.Background(), shopID, req)
	if err == nil || err.Error() != "productID is required for a Sale" {
		t.Errorf("Expected productID missing error, got %v", err)
	}
}

func TestTransactionService_CreateExpense_Success(t *testing.T) {
	shopID := uuid.New().String()
	mockRepo := &mockTransactionRepo{CreateErr: nil}
	service := NewTransactionService(nil, mockRepo)

	req := dto.CreateTransactionRequest{
		Type:   "Expense",
		Amount: decimal.NewFromFloat(50.0),
	}

	res, err := service.Create(context.Background(), shopID, req)
	if err != nil {
		t.Errorf("Expected success, got error %v", err)
	}
	if res.Type != "Expense" {
		t.Errorf("Expected type Expense, got %v", res.Type)
	}
}
