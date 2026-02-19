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

type fakeProductServiceRepo struct {
	listProducts []models.Product
	listErr      error

	findProduct *models.Product
	findErr     error

	createdProduct *models.Product
	createErr      error

	savedProduct *models.Product
	saveErr      error

	deleteResult bool
	deleteErr    error
}

func (r *fakeProductServiceRepo) ListByShopID(_ context.Context, _ uuid.UUID) ([]models.Product, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.listProducts, nil
}

func (r *fakeProductServiceRepo) FindByIDAndShopID(_ context.Context, productID, shopID uuid.UUID) (*models.Product, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.findProduct == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if r.findProduct.ID != productID || r.findProduct.ShopID != shopID {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *r.findProduct
	return &cp, nil
}

func (r *fakeProductServiceRepo) Create(_ context.Context, product *models.Product) error {
	if r.createErr != nil {
		return r.createErr
	}
	cp := *product
	r.createdProduct = &cp
	return nil
}

func (r *fakeProductServiceRepo) Save(_ context.Context, product *models.Product) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	cp := *product
	r.savedProduct = &cp
	return nil
}

func (r *fakeProductServiceRepo) DeleteByIDAndShopID(_ context.Context, _, _ uuid.UUID) (bool, error) {
	if r.deleteErr != nil {
		return false, r.deleteErr
	}
	return r.deleteResult, nil
}

func TestProductService_List_RoleBasedResponseShape(t *testing.T) {
	shopID := uuid.New()
	repo := &fakeProductServiceRepo{
		listProducts: []models.Product{
			{
				ID:            uuid.New(),
				Name:          "iPhone",
				Description:   "128GB",
				Category:      "Smartphones",
				PurchasePrice: decimal.RequireFromString("100.50"),
				SellingPrice:  decimal.RequireFromString("150.00"),
				Stock:         10,
				ImageURL:      "https://example.com/image.jpg",
				ShopID:        shopID,
			},
		},
	}
	service := NewProductService(repo)

	adminData, err := service.List(context.Background(), shopID.String(), models.RoleAdmin)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	adminProducts, ok := adminData.([]dto.ProductResponse)
	if !ok {
		t.Fatalf("expected []dto.ProductResponse, got %T", adminData)
	}
	if len(adminProducts) != 1 {
		t.Fatalf("expected 1 product, got %d", len(adminProducts))
	}

	superData, err := service.List(context.Background(), shopID.String(), models.RoleSuperAdmin)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	superProducts, ok := superData.([]dto.ProductAdminResponse)
	if !ok {
		t.Fatalf("expected []dto.ProductAdminResponse, got %T", superData)
	}
	if len(superProducts) != 1 {
		t.Fatalf("expected 1 product, got %d", len(superProducts))
	}
	if !superProducts[0].PurchasePrice.Equal(decimal.RequireFromString("100.50")) {
		t.Fatalf("unexpected purchasePrice value")
	}
}

func TestProductService_Create_SuperAdminRequiresPurchasePrice(t *testing.T) {
	shopID := uuid.New()
	service := NewProductService(&fakeProductServiceRepo{})

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

func TestProductService_Create_SuperAdminPurchasePriceMustBePositive(t *testing.T) {
	shopID := uuid.New()
	service := NewProductService(&fakeProductServiceRepo{})

	purchase := decimal.Zero
	req := dto.CreateProductRequest{
		Name:          "Galaxy",
		Category:      "Smartphones",
		PurchasePrice: &purchase,
		SellingPrice:  decimal.RequireFromString("200.00"),
		Stock:         5,
	}

	_, err := service.Create(context.Background(), shopID.String(), models.RoleSuperAdmin, req)
	if !errors.Is(err, ErrPurchasePriceNegative) {
		t.Fatalf("expected ErrPurchasePriceNegative, got %v", err)
	}
}

func TestProductService_Create_SellingPriceMustBePositive(t *testing.T) {
	shopID := uuid.New()
	service := NewProductService(&fakeProductServiceRepo{})

	req := dto.CreateProductRequest{
		Name:         "Mouse",
		Category:     "Accessories",
		SellingPrice: decimal.Zero,
		Stock:        5,
	}

	_, err := service.Create(context.Background(), shopID.String(), models.RoleAdmin, req)
	if !errors.Is(err, ErrSellingPriceNegative) {
		t.Fatalf("expected ErrSellingPriceNegative, got %v", err)
	}
}

func TestProductService_Create_AdminRejectsPurchasePriceInput(t *testing.T) {
	shopID := uuid.New()
	service := NewProductService(&fakeProductServiceRepo{})

	purchase := decimal.RequireFromString("40.00")
	req := dto.CreateProductRequest{
		Name:          "Mouse",
		Category:      "Accessories",
		PurchasePrice: &purchase,
		SellingPrice:  decimal.RequireFromString("60.00"),
		Stock:         5,
	}

	_, err := service.Create(context.Background(), shopID.String(), models.RoleAdmin, req)
	if !errors.Is(err, ErrPurchasePriceForbidden) {
		t.Fatalf("expected ErrPurchasePriceForbidden, got %v", err)
	}
}

func TestProductService_Create_AdminWithoutPurchasePricePersistsZero(t *testing.T) {
	shopID := uuid.New()
	repo := &fakeProductServiceRepo{}
	service := NewProductService(repo)

	req := dto.CreateProductRequest{
		Name:         "Mouse",
		Category:     "Accessories",
		SellingPrice: decimal.RequireFromString("60.00"),
		Stock:        5,
	}

	createdData, err := service.Create(context.Background(), shopID.String(), models.RoleAdmin, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if _, ok := createdData.(dto.ProductResponse); !ok {
		t.Fatalf("expected dto.ProductResponse, got %T", createdData)
	}
	if repo.createdProduct == nil {
		t.Fatalf("expected created product to be captured")
	}
	if !repo.createdProduct.PurchasePrice.Equal(decimal.Zero) {
		t.Fatalf("expected purchasePrice to be zero, got %s", repo.createdProduct.PurchasePrice.String())
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	repo := &fakeProductServiceRepo{
		findErr: gorm.ErrRecordNotFound,
	}
	service := NewProductService(repo)

	_, err := service.GetByID(context.Background(), shopID.String(), productID.String(), models.RoleAdmin)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_Update_StockCannotBeNegative(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	repo := &fakeProductServiceRepo{
		findProduct: &models.Product{
			ID:            productID,
			Name:          "Mouse",
			Category:      "Accessories",
			PurchasePrice: decimal.RequireFromString("20.00"),
			SellingPrice:  decimal.RequireFromString("40.00"),
			Stock:         8,
			ShopID:        shopID,
		},
	}
	service := NewProductService(repo)

	req := dto.UpdateProductRequest{
		Name:         "Mouse",
		Category:     "Accessories",
		SellingPrice: decimal.RequireFromString("45.00"),
		Stock:        -1,
	}

	_, err := service.Update(context.Background(), shopID.String(), productID.String(), models.RoleAdmin, req)
	if !errors.Is(err, ErrStockNegative) {
		t.Fatalf("expected ErrStockNegative, got %v", err)
	}
}

func TestProductService_Update_AdminRejectsPurchasePriceInput(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	repo := &fakeProductServiceRepo{
		findProduct: &models.Product{
			ID:            productID,
			Name:          "Keyboard",
			Category:      "Accessories",
			PurchasePrice: decimal.RequireFromString("20.00"),
			SellingPrice:  decimal.RequireFromString("40.00"),
			Stock:         8,
			ShopID:        shopID,
		},
	}
	service := NewProductService(repo)

	purchase := decimal.RequireFromString("25.00")
	req := dto.UpdateProductRequest{
		Name:          "Keyboard",
		Category:      "Accessories",
		PurchasePrice: &purchase,
		SellingPrice:  decimal.RequireFromString("45.00"),
		Stock:         7,
	}

	_, err := service.Update(context.Background(), shopID.String(), productID.String(), models.RoleAdmin, req)
	if !errors.Is(err, ErrPurchasePriceForbidden) {
		t.Fatalf("expected ErrPurchasePriceForbidden, got %v", err)
	}
}

func TestProductService_Update_AdminForcesPurchasePriceToZero(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	repo := &fakeProductServiceRepo{
		findProduct: &models.Product{
			ID:            productID,
			Name:          "Headphones",
			Category:      "Audio",
			PurchasePrice: decimal.RequireFromString("150.00"),
			SellingPrice:  decimal.RequireFromString("250.00"),
			Stock:         4,
			ShopID:        shopID,
		},
	}
	service := NewProductService(repo)

	req := dto.UpdateProductRequest{
		Name:         "Headphones",
		Category:     "Audio",
		SellingPrice: decimal.RequireFromString("260.00"),
		Stock:        3,
	}

	updatedData, err := service.Update(context.Background(), shopID.String(), productID.String(), models.RoleAdmin, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if _, ok := updatedData.(dto.ProductResponse); !ok {
		t.Fatalf("expected dto.ProductResponse, got %T", updatedData)
	}
	if repo.savedProduct == nil {
		t.Fatalf("expected saved product to be captured")
	}
	if !repo.savedProduct.PurchasePrice.Equal(decimal.Zero) {
		t.Fatalf("expected purchasePrice to be zero, got %s", repo.savedProduct.PurchasePrice.String())
	}
	if repo.savedProduct.ID != productID || repo.savedProduct.ShopID != shopID {
		t.Fatalf("expected tenant-scoped save with original id/shop")
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	shopID := uuid.New()
	productID := uuid.New()
	service := NewProductService(&fakeProductServiceRepo{
		deleteResult: false,
	})

	err := service.Delete(context.Background(), shopID.String(), productID.String())
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
