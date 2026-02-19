package services

import (
	"context"
	"errors"
	"strings"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInvalidShopWhatsApp = errors.New("invalid shop whatsapp number")

type PublicService struct {
	shopRepo    publicServiceShopRepository
	productRepo publicServiceProductRepository
}

type publicServiceShopRepository interface {
	FindByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error)
}

type publicServiceProductRepository interface {
	ListByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Product, error)
}

func NewPublicService(shopRepo publicServiceShopRepository, productRepo publicServiceProductRepository) *PublicService {
	return &PublicService{
		shopRepo:    shopRepo,
		productRepo: productRepo,
	}
}

func (s *PublicService) ListProductsByShopID(ctx context.Context, shopID string) ([]dto.PublicProductResponse, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	shop, err := s.shopRepo.FindByID(ctx, shopUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}
	if !shop.Active {
		return nil, ErrShopInactive
	}

	products, err := s.productRepo.ListByShopID(ctx, shopUUID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PublicProductResponse, 0, len(products))
	for _, product := range products {
		whatsappLink, linkErr := utils.GenerateWhatsAppLink(shop.WhatsAppNumber, product.Name)
		if linkErr != nil {
			return nil, ErrInvalidShopWhatsApp
		}

		responses = append(responses, dto.PublicProductResponse{
			ID:           product.ID.String(),
			Name:         product.Name,
			Description:  product.Description,
			Category:     product.Category,
			SellingPrice: product.SellingPrice,
			Stock:        product.Stock,
			ImageURL:     product.ImageURL,
			WhatsAppLink: whatsappLink,
		})
	}

	return responses, nil
}
