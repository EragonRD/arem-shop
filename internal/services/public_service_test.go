package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type mockPublicShopRepository struct {
	shop *models.Shop
	err  error
}

func (m *mockPublicShopRepository) FindByID(_ context.Context, _ uuid.UUID) (*models.Shop, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.shop == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return m.shop, nil
}

type mockPublicProductRepository struct {
	products []models.Product
	err      error
}

func (m *mockPublicProductRepository) ListByShopID(_ context.Context, _ uuid.UUID) ([]models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.products, nil
}

func TestPublicServiceListProductsByShopID_InvalidShopID(t *testing.T) {
	service := NewPublicService(&mockPublicShopRepository{}, &mockPublicProductRepository{})

	_, err := service.ListProductsByShopID(context.Background(), "not-a-uuid")
	if !errors.Is(err, ErrInvalidShopID) {
		t.Fatalf("expected ErrInvalidShopID, got %v", err)
	}
}

func TestPublicServiceListProductsByShopID_InactiveShop(t *testing.T) {
	shopID := uuid.New()
	service := NewPublicService(
		&mockPublicShopRepository{
			shop: &models.Shop{
				ID:             shopID,
				Name:           "Demo Shop",
				Active:         false,
				WhatsAppNumber: "+212600000000",
			},
		},
		&mockPublicProductRepository{},
	)

	_, err := service.ListProductsByShopID(context.Background(), shopID.String())
	if !errors.Is(err, ErrShopInactive) {
		t.Fatalf("expected ErrShopInactive, got %v", err)
	}
}

func TestPublicServiceListProductsByShopID_HidesPurchasePrice(t *testing.T) {
	shopID := uuid.New()
	service := NewPublicService(
		&mockPublicShopRepository{
			shop: &models.Shop{
				ID:             shopID,
				Name:           "Demo Shop",
				Active:         true,
				WhatsAppNumber: "+212600000000",
			},
		},
		&mockPublicProductRepository{
			products: []models.Product{
				{
					ID:            uuid.New(),
					Name:          "Laptop",
					Description:   "Demo",
					Category:      "Laptops",
					PurchasePrice: decimal.RequireFromString("799.99"),
					SellingPrice:  decimal.RequireFromString("999.99"),
					Stock:         4,
					ImageURL:      "https://example.com/laptop.jpg",
					ShopID:        shopID,
				},
			},
		},
	)

	responses, err := service.ListProductsByShopID(context.Background(), shopID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}

	if responses[0].WhatsAppLink == "" {
		t.Fatalf("expected whatsapp link to be generated")
	}

	payload, err := json.Marshal(responses[0])
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	if containsPurchasePriceKey(payload) {
		t.Fatalf("public response must not expose purchasePrice: %s", string(payload))
	}
}

func containsPurchasePriceKey(payload []byte) bool {
	return string(payload) != "" && string(payload) != "null" &&
		jsonContainsKey(payload, "purchasePrice")
}

func jsonContainsKey(payload []byte, key string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal(payload, &obj); err != nil {
		return false
	}
	_, exists := obj[key]
	return exists
}
