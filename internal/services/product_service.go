package services

import (
	"context"
	"errors"
	"strings"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrInvalidShopID          = errors.New("invalid shopID")
	ErrInvalidProductID       = errors.New("invalid product ID")
	ErrProductNotFound        = errors.New("product not found")
	ErrPurchasePriceRequired  = errors.New("purchasePrice is required for SuperAdmin")
	ErrPurchasePriceNegative  = errors.New("purchasePrice must be greater than 0")
	ErrPurchasePriceForbidden = errors.New("purchasePrice is not allowed for Admin")
	ErrSellingPriceNegative   = errors.New("sellingPrice must be greater than 0")
	ErrStockNegative          = errors.New("stock cannot be negative")
	ErrInvalidProductName     = errors.New("name is required")
	ErrInvalidProductCategory = errors.New("category is required")
)

type ProductService struct {
	productRepo productServiceProductRepository
}

type productServiceProductRepository interface {
	ListByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Product, error)
	FindByIDAndShopID(ctx context.Context, productID, shopID uuid.UUID) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Save(ctx context.Context, product *models.Product) error
	DeleteByIDAndShopID(ctx context.Context, productID, shopID uuid.UUID) (bool, error)
}

func NewProductService(productRepo productServiceProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) List(ctx context.Context, shopID string, role models.UserRole) (interface{}, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	products, err := s.productRepo.ListByShopID(ctx, shopUUID)
	if err != nil {
		return nil, err
	}

	if role == models.RoleSuperAdmin {
		responses := make([]dto.ProductAdminResponse, 0, len(products))
		for _, product := range products {
			responses = append(responses, toProductAdminResponse(product))
		}
		return responses, nil
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}

	return responses, nil
}

func (s *ProductService) GetByID(ctx context.Context, shopID, productID string, role models.UserRole) (interface{}, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	productUUID, err := uuid.Parse(strings.TrimSpace(productID))
	if err != nil {
		return nil, ErrInvalidProductID
	}

	product, err := s.productRepo.FindByIDAndShopID(ctx, productUUID, shopUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if role == models.RoleSuperAdmin {
		resp := toProductAdminResponse(*product)
		return resp, nil
	}

	resp := toProductResponse(*product)
	return resp, nil
}

func (s *ProductService) Create(ctx context.Context, shopID string, role models.UserRole, req dto.CreateProductRequest) (interface{}, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	if err := validateProductCreateInput(role, req); err != nil {
		return nil, err
	}

	purchasePrice := decimal.Zero
	if role == models.RoleSuperAdmin && req.PurchasePrice != nil {
		purchasePrice = req.PurchasePrice.Round(2)
	}

	product := models.Product{
		Name:          strings.TrimSpace(req.Name),
		Description:   strings.TrimSpace(req.Description),
		Category:      strings.TrimSpace(req.Category),
		PurchasePrice: purchasePrice,
		SellingPrice:  req.SellingPrice.Round(2),
		Stock:         req.Stock,
		ImageURL:      strings.TrimSpace(req.ImageURL),
		ShopID:        shopUUID,
	}

	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	if role == models.RoleSuperAdmin {
		resp := toProductAdminResponse(product)
		return resp, nil
	}

	resp := toProductResponse(product)
	return resp, nil
}

func (s *ProductService) Update(ctx context.Context, shopID, productID string, role models.UserRole, req dto.UpdateProductRequest) (interface{}, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	productUUID, err := uuid.Parse(strings.TrimSpace(productID))
	if err != nil {
		return nil, ErrInvalidProductID
	}

	existing, err := s.productRepo.FindByIDAndShopID(ctx, productUUID, shopUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if err := validateProductUpdateInput(role, req); err != nil {
		return nil, err
	}

	existing.Name = strings.TrimSpace(req.Name)
	existing.Description = strings.TrimSpace(req.Description)
	existing.Category = strings.TrimSpace(req.Category)
	existing.SellingPrice = req.SellingPrice.Round(2)
	existing.Stock = req.Stock
	existing.ImageURL = strings.TrimSpace(req.ImageURL)

	if role == models.RoleSuperAdmin {
		if req.PurchasePrice != nil {
			existing.PurchasePrice = req.PurchasePrice.Round(2)
		}
	} else {
		existing.PurchasePrice = decimal.Zero
	}

	if err := s.productRepo.Save(ctx, existing); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if role == models.RoleSuperAdmin {
		resp := toProductAdminResponse(*existing)
		return resp, nil
	}

	resp := toProductResponse(*existing)
	return resp, nil
}

func (s *ProductService) Delete(ctx context.Context, shopID, productID string) error {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return ErrInvalidShopID
	}

	productUUID, err := uuid.Parse(strings.TrimSpace(productID))
	if err != nil {
		return ErrInvalidProductID
	}

	deleted, err := s.productRepo.DeleteByIDAndShopID(ctx, productUUID, shopUUID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrProductNotFound
	}

	return nil
}

func validateProductCreateInput(role models.UserRole, req dto.CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidProductName
	}
	if strings.TrimSpace(req.Category) == "" {
		return ErrInvalidProductCategory
	}
	if req.SellingPrice.LessThanOrEqual(decimal.Zero) {
		return ErrSellingPriceNegative
	}
	if req.Stock < 0 {
		return ErrStockNegative
	}
	if role == models.RoleSuperAdmin {
		if req.PurchasePrice == nil {
			return ErrPurchasePriceRequired
		}
		if req.PurchasePrice.LessThanOrEqual(decimal.Zero) {
			return ErrPurchasePriceNegative
		}
	} else if req.PurchasePrice != nil {
		return ErrPurchasePriceForbidden
	}
	return nil
}

func validateProductUpdateInput(role models.UserRole, req dto.UpdateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidProductName
	}
	if strings.TrimSpace(req.Category) == "" {
		return ErrInvalidProductCategory
	}
	if req.SellingPrice.LessThanOrEqual(decimal.Zero) {
		return ErrSellingPriceNegative
	}
	if req.Stock < 0 {
		return ErrStockNegative
	}
	if role == models.RoleSuperAdmin {
		if req.PurchasePrice != nil && req.PurchasePrice.LessThanOrEqual(decimal.Zero) {
			return ErrPurchasePriceNegative
		}
	} else if req.PurchasePrice != nil {
		return ErrPurchasePriceForbidden
	}
	return nil
}

func toProductResponse(product models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:           product.ID.String(),
		Name:         product.Name,
		Description:  product.Description,
		Category:     product.Category,
		SellingPrice: product.SellingPrice,
		Stock:        product.Stock,
		ImageURL:     product.ImageURL,
		ShopID:       product.ShopID.String(),
		CreatedAt:    product.CreatedAt,
	}
}

func toProductAdminResponse(product models.Product) dto.ProductAdminResponse {
	return dto.ProductAdminResponse{
		ProductResponse: toProductResponse(product),
		PurchasePrice:   product.PurchasePrice,
	}
}
