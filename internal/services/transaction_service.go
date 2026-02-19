// -----NEW-----
package services

import (
	"context"
	"errors"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionService définit le contrat du service
type TransactionService interface {
	Create(ctx context.Context, shopID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
}

type transactionService struct {
	db              *gorm.DB
	transactionRepo repository.TransactionRepository
}

// NewTransactionService instancie le service
func NewTransactionService(db *gorm.DB, transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		db:              db,
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) Create(ctx context.Context, shopID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// 1. Sécurité Multi-tenant : validation du ShopID
	shopUUID, err := uuid.Parse(shopID)
	if err != nil {
		return nil, errors.New("invalid shop ID")
	}

	// 2. Routage selon le type de transaction
	switch req.Type {
	case "Sale":
		if req.ProductID == nil {
			return nil, errors.New("productID is required for a Sale")
		}
		if req.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0 for a Sale")
		}

		// Ouverture de la transaction SQL
		tx := s.db.WithContext(ctx).Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}

		// Appel de la méthode atomique du repo
		err := s.transactionRepo.CreateSale(ctx, tx, shopUUID, *req.ProductID, req.Quantity, req.Amount)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Si tout s'est bien passé, on valide la transaction
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}

		return &dto.TransactionResponse{
			Type:      req.Type,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Amount:    req.Amount,
			ShopID:    shopUUID,
		}, nil

	case "Expense", "Withdrawal":
		transaction := &models.Transaction{
			Type:   req.Type,
			Amount: req.Amount,
			ShopID: shopUUID,
		}

		err := s.transactionRepo.Create(ctx, nil, transaction)
		if err != nil {
			return nil, err
		}

		return &dto.TransactionResponse{
			ID:     transaction.ID,
			Type:   transaction.Type,
			Amount: transaction.Amount,
			ShopID: transaction.ShopID,
		}, nil

	default:
		return nil, errors.New("invalid transaction type")
	}
}

//-----OLD-----
/*
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
	ErrInvalidTransactionType      = errors.New("invalid transaction type")
	ErrTransactionProductRequired  = errors.New("productID is required for Sale")
	ErrTransactionProductForbidden = errors.New("productID must be empty for Expense/Withdrawal")
	ErrTransactionQuantityRequired = errors.New("quantity must be greater than 0 for Sale")
	ErrTransactionQuantityInvalid  = errors.New("quantity must be 0 for Expense/Withdrawal")
	ErrTransactionAmountInvalid    = errors.New("amount must be greater than 0")
	ErrInsufficientStock           = errors.New("insufficient stock")
)

type TransactionService struct {
	productRepo     transactionServiceProductRepository
	transactionRepo transactionServiceTransactionRepository
	runInTx         func(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type transactionServiceProductRepository interface {
	FindByIDAndShopIDForUpdate(ctx context.Context, tx *gorm.DB, productID, shopID uuid.UUID) (*models.Product, error)
	SaveWithTx(ctx context.Context, tx *gorm.DB, product *models.Product) error
}

type transactionServiceTransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error
}

func NewTransactionService(db *gorm.DB, productRepo transactionServiceProductRepository, transactionRepo transactionServiceTransactionRepository) *TransactionService {
	service := &TransactionService{
		productRepo:     productRepo,
		transactionRepo: transactionRepo,
	}

	if db != nil {
		service.runInTx = func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return db.WithContext(ctx).Transaction(fn)
		}
	} else {
		service.runInTx = func(_ context.Context, fn func(tx *gorm.DB) error) error {
			return fn(nil)
		}
	}

	return service
}

func (s *TransactionService) Create(ctx context.Context, shopID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	transactionType := models.TransactionType(strings.TrimSpace(req.Type))
	if transactionType != models.TransactionSale && transactionType != models.TransactionExpense && transactionType != models.TransactionWithdrawal {
		return nil, ErrInvalidTransactionType
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrTransactionAmountInvalid
	}

	quantity := 0
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	var productUUID *uuid.UUID
	if req.ProductID != nil && strings.TrimSpace(*req.ProductID) != "" {
		parsedProductID, parseErr := uuid.Parse(strings.TrimSpace(*req.ProductID))
		if parseErr != nil {
			return nil, ErrInvalidProductID
		}
		productUUID = &parsedProductID
	}

	if transactionType == models.TransactionSale {
		if productUUID == nil {
			return nil, ErrTransactionProductRequired
		}
		if quantity <= 0 {
			return nil, ErrTransactionQuantityRequired
		}
	} else {
		if productUUID != nil {
			return nil, ErrTransactionProductForbidden
		}
		if quantity != 0 {
			return nil, ErrTransactionQuantityInvalid
		}
	}

	transaction := models.Transaction{
		Type:      transactionType,
		ProductID: productUUID,
		Quantity:  quantity,
		Amount:    req.Amount.Round(2),
		ShopID:    shopUUID,
	}

	err = s.runInTx(ctx, func(tx *gorm.DB) error {
		if transactionType == models.TransactionSale {
			product, lockErr := s.productRepo.FindByIDAndShopIDForUpdate(ctx, tx, *productUUID, shopUUID)
			if lockErr != nil {
				if errors.Is(lockErr, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return lockErr
			}

			if product.Stock < quantity {
				return ErrInsufficientStock
			}

			product.Stock -= quantity
			if saveErr := s.productRepo.SaveWithTx(ctx, tx, product); saveErr != nil {
				return saveErr
			}
		}

		if createErr := s.transactionRepo.Create(ctx, tx, &transaction); createErr != nil {
			return createErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := toTransactionResponse(transaction)
	return &resp, nil
}

func toTransactionResponse(transaction models.Transaction) dto.TransactionResponse {
	var productID *string
	if transaction.ProductID != nil {
		productIDValue := transaction.ProductID.String()
		productID = &productIDValue
	}

	return dto.TransactionResponse{
		ID:        transaction.ID.String(),
		Type:      string(transaction.Type),
		ProductID: productID,
		Quantity:  transaction.Quantity,
		Amount:    transaction.Amount,
		ShopID:    transaction.ShopID.String(),
		CreatedAt: transaction.CreatedAt,
	}
}
*/

//-----Commenaire Personne 3-----
//L'objectif du Service est d'orchestrer, pas d'exécuter les requêtes DB complexes.
// ICI
//ouvrir la transaction SQL (db.Begin()), appeler notre fameux CreateSale, et valider (Commit) ou annuler (Rollback) en fonction du résultat.
//
