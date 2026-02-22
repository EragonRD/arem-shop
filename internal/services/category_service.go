package services

import (
	"context"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"

	"github.com/google/uuid"
)

type categoryRepository interface {
	ListByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Category, error)
}

type CategoryService struct {
	repo categoryRepository
}

func NewCategoryService(repo categoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context, shopID string) ([]dto.CategoryResponse, error) {
	shopUUID, err := uuid.Parse(shopID)
	if err != nil {
		return nil, ErrInvalidShopID
	}

	categories, err := s.repo.ListByShopID(ctx, shopUUID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = dto.CategoryResponse{
			ID:   c.ID.String(),
			Name: c.Name,
		}
	}

	return res, nil
}
