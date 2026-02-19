package services

import (
	"context"
	"errors"
	"testing"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type fakeProductServiceRepo struct {
	listProducts   []models.Product
	createdProduct *models.Product
	createErr      error
}

func (r *fakeProductServiceRepo) ListByShopID(_ context.Context, _ uuid.UUID) ([]models.Product, error) {
	return r.listProducts, nil
}

func (r *fakeProductServiceRepo) FindByIDAndShopID(_ context.Context, _, _ uuid.UUID) (*models.Product, error) {
	return nil, errors.New("not implemented")
}

func (r *fakeProductServiceRepo) Create(_ context.Context, product *models.Product) error {
	if r.createErr != nil {
		return r.createErr
	}
	cp := *product
	r.createdProduct = &cp
	return nil
}

func (r *fakeProductServiceRepo) Save(_ context.Context, _ *models.Product) error {
	return nil
}

func (r *fakeProductServiceRepo) DeleteByIDAndShopID(_ context.Context, _, _ uuid.UUID) (bool, error) {
	return false, nil
}

func TestProductService_List_HidesPurchasePriceForAdmin(t *testing.T) {
	shopID := uuid.New()
	repo := &fakeProductServiceRepo{
		listProducts: []models.Product{
			{
				ID:            uuid.New(),
				Name:          "iPhone",
				Category:      "Smartphones",
				PurchasePrice: decimal.RequireFromString("100.50"),
				SellingPrice:  decimal.RequireFromString("150.00"),
				Stock:         10,
				ShopID:        shopID,
			},
		},
	}

	service := NewProductService(repo)

	adminProducts, err := service.List(context.Background(), shopID.String(), models.RoleAdmin)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(adminProducts) != 1 {
		t.Fatalf("expected 1 product, got %d", len(adminProducts))
	}
	if adminProducts[0].PurchasePrice != nil {
		t.Fatalf("expected purchasePrice to be hidden for admin")
	}

	superAdminProducts, err := service.List(context.Background(), shopID.String(), models.RoleSuperAdmin)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if superAdminProducts[0].PurchasePrice == nil {
		t.Fatalf("expected purchasePrice for superadmin")
	}
	if !superAdminProducts[0].PurchasePrice.Equal(decimal.RequireFromString("100.50")) {
		t.Fatalf("unexpected purchasePrice value")
	}
}

func TestProductService_Create_SuperAdminRequiresPurchasePrice(t *testing.T) {
	shopID := uuid.New()
	repo := &fakeProductServiceRepo{}
	service := NewProductService(repo)

	req := dto.CreateProductRequest{
		Name:         "Galaxy",
		Category:     "Smartphones",
		SellingPrice: decimal.RequireFromString("200.00"),
		Stock:        5,
	}

	_, err := service.Create(context.Background(), shopID.String(), models.RoleSuperAdmin, req)
	if !errors.Is(err, ErrPurchasePriceRequired) {
		t.Fatalf("expected ErrPurchasePriceRequired, got %v", err)
	}
}

func TestProductService_Create_AdminForcesPurchasePriceToZero(t *testing.T) {
	shopID := uuid.New()
	repo := &fakeProductServiceRepo{}
	service := NewProductService(repo)

	purchase := decimal.RequireFromString("40.00")
	req := dto.CreateProductRequest{
		Name:          "Mouse",
		Category:      "Accessories",
		PurchasePrice: &purchase,
		SellingPrice:  decimal.RequireFromString("60.00"),
		Stock:         5,
	}

	_, err := service.Create(context.Background(), shopID.String(), models.RoleAdmin, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if repo.createdProduct == nil {
		t.Fatalf("expected created product to be captured")
	}
	if !repo.createdProduct.PurchasePrice.Equal(decimal.Zero) {
		t.Fatalf("expected purchasePrice to be zero for admin-created product, got %s", repo.createdProduct.PurchasePrice.String())
	}
}
