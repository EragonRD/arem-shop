package repository

import (
	"context"
	"errors"

	"arem-shop/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. DÉFINITION DE L'INTERFACE (Obligatoire selon l'architecture)
type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error
	SumAmountByShopAndTypes(ctx context.Context, shopID uuid.UUID, types ...models.TransactionType) (decimal.Decimal, error)
	CreateSale(ctx context.Context, tx *gorm.DB, shopID uuid.UUID, productID uuid.UUID, qty int, amount decimal.Decimal) error
	FindAll(ctx context.Context, shopID uuid.UUID) ([]models.Transaction, error)
}

// 2. STRUCT PRIVÉ
type transactionRepository struct {
	db *gorm.DB
}

// 3. CONSTRUCTEUR QUI RETOURNE L'INTERFACE
// Note : il renvoie TransactionRepository (l'interface) mais instancie &transactionRepository (le struct)
func NewTransactionRepository(db *gorm.DB) TransactionRepository { // <-- PLUS D'ASTÉRISQUE ICI
	return &transactionRepository{db: db}
}

// --- Tes méthodes existantes (avec 't' minuscule pour le receiver) ---
func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	execDB := tx
	if execDB == nil {
		execDB = r.db
	}

	return execDB.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) SumAmountByShopAndTypes(ctx context.Context, shopID uuid.UUID, types ...models.TransactionType) (decimal.Decimal, error) {
	if len(types) == 0 {
		return decimal.Zero, nil
	}

	var totalRaw string
	err := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount)::text, '0')").
		Where("shop_id = ? AND type IN ?", shopID, types).
		Scan(&totalRaw).Error
	if err != nil {
		return decimal.Zero, err
	}

	total, err := decimal.NewFromString(totalRaw)
	if err != nil {
		return decimal.Zero, err
	}

	return total, nil
}

// --- Nouvelles méthodes (avec 't' minuscule pour le receiver) ---

// CreateSale gère la création d'une vente de manière atomique avec un verrou pessimiste
func (r *transactionRepository) CreateSale(ctx context.Context, tx *gorm.DB, shopID uuid.UUID, productID uuid.UUID, qty int, amount decimal.Decimal) error {
	var product models.Product

	// SELECT FOR UPDATE : Verrouille la ligne du produit pour éviter les race conditions
	if err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND shop_id = ?", productID, shopID).
		First(&product).Error; err != nil {
		return err
	}

	// Vérification du stock
	if product.Stock < qty {
		return errors.New("insufficient stock")
	}

	// Mise à jour du stock (décrémentation)
	if err := tx.WithContext(ctx).Model(&product).Update("stock", product.Stock-qty).Error; err != nil {
		return err
	}

	// Insertion de la transaction en utilisant ta méthode Create existante
	// Note: On passe la string "Sale" ou la constante de ton modèle
	transaction := &models.Transaction{
		Type:      "Sale", // Ou utilise ta constante models.TransactionTypeSale si elle existe
		ProductID: &productID,
		Quantity:  qty,
		Amount:    amount,
		ShopID:    shopID,
	}

	return r.Create(ctx, tx, transaction)
}

// FindAll (nécessaire pour récupérer l'historique)
func (r *transactionRepository) FindAll(ctx context.Context, shopID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Filtre par shopID critique pour le multitenant
	err := r.db.WithContext(ctx).Where("shop_id = ?", shopID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

/*
TransactionRepository était souligné en rouge car il y avait un conflit de noms entre l'interface et le struct. En Go, on ne peut pas avoir une interface et un struct qui portent exactement le même nom. De plus, les méthodes doivent être attachées au struct privé (transactionRepository) et non à l'interface.

Voici les corrections apportées :
1. L'interface est renommée en transactionRepository (avec un t minuscule) pour la rendre privée.
2. Le struct est renommé en TransactionRepository (avec un T majuscule) pour le rendre public.
3. Toutes les méthodes sont attachées au struct privé (r *transactionRepository) et non à l'interface.

4. Le constructeur NewTransactionRepository retourne un pointeur vers TransactionRepository, qui implémente l'interface transactionRepository.

*/
