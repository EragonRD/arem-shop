package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"arem-shop/internal/config"
	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"
	"arem-shop/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorizedRole   = errors.New("unauthorized role")
	ErrCrossShopRegister  = errors.New("cross-shop registration forbidden")
	ErrShopNotFound       = errors.New("shop not found")
	ErrShopInactive       = errors.New("shop is inactive")
	ErrEmailAlreadyUsed   = errors.New("email already exists in this shop")
	ErrInvalidRole        = errors.New("invalid role")
)

type AuthService struct {
	cfg      config.AppConfig
	userRepo *repository.UserRepository
	shopRepo *repository.ShopRepository
}

func NewAuthService(cfg config.AppConfig, userRepo *repository.UserRepository, shopRepo *repository.ShopRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
		shopRepo: shopRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest, actorRole models.UserRole, actorShopID string) (*dto.AuthUserResponse, error) {
	if actorRole != models.RoleSuperAdmin {
		return nil, ErrUnauthorizedRole
	}

	targetShopID, err := uuid.Parse(strings.TrimSpace(req.ShopID))
	if err != nil {
		return nil, fmt.Errorf("parse req shopID: %w", err)
	}

	callerShopID, err := uuid.Parse(strings.TrimSpace(actorShopID))
	if err != nil {
		return nil, fmt.Errorf("parse actor shopID: %w", err)
	}

	// Protection multi-tenant: un SuperAdmin ne gere que les users de son shop.
	if callerShopID != targetShopID {
		return nil, ErrCrossShopRegister
	}

	shop, err := s.shopRepo.FindByID(ctx, targetShopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}
	if !shop.Active {
		return nil, ErrShopInactive
	}

	exists, err := s.userRepo.ExistsByEmailAndShopID(ctx, req.Email, targetShopID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyUsed
	}

	role := models.UserRole(strings.TrimSpace(req.Role))
	if role != models.RoleSuperAdmin && role != models.RoleAdmin {
		return nil, ErrInvalidRole
	}

	hashedPassword, err := utils.HashPassword(req.Password, s.cfg.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
		Role:     role,
		ShopID:   targetShopID,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "uq_users_shop_email") {
			return nil, ErrEmailAlreadyUsed
		}
		return nil, err
	}

	resp := toAuthUserResponse(user)
	return &resp, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	shopID, err := uuid.Parse(strings.TrimSpace(req.ShopID))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	shop, err := s.shopRepo.FindByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !shop.Active {
		return nil, ErrShopInactive
	}

	user, err := s.userRepo.FindByEmailAndShopID(ctx, req.Email, shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(s.cfg, *user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	resp := &dto.LoginResponse{
		Token: token,
		User:  toAuthUserResponse(*user),
	}
	return resp, nil
}

func toAuthUserResponse(user models.User) dto.AuthUserResponse {
	return dto.AuthUserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		ShopID:    user.ShopID.String(),
		CreatedAt: user.CreatedAt,
	}
}
