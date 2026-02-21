package services

import (
	"context"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"

	"github.com/google/uuid"
)

type ShopService struct {
	shopRepo *repository.ShopRepository
}

func NewShopService(shopRepo *repository.ShopRepository) *ShopService {
	return &ShopService{shopRepo: shopRepo}
}

func (s *ShopService) UpdateShopInfo(ctx context.Context, shopID uuid.UUID, payload dto.ShopUpdatePayload) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByID(ctx, shopID)
	if err != nil {
		return nil, err
	}

	shop.Name = payload.Name
	if payload.WhatsAppNumber != "" {
		shop.WhatsAppNumber = payload.WhatsAppNumber
	}

	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return nil, err
	}

	return shop, nil
}
