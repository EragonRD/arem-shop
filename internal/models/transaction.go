package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Définition du type personnalisé attendu par le repository et le report_service
type TransactionType string

// Constantes pour l'énumération des types de transaction
const (
	TransactionSale       TransactionType = "Sale"
	TransactionExpense    TransactionType = "Expense"
	TransactionWithdrawal TransactionType = "Withdrawal"
)

// Transaction représente l'entité transaction dans la base de données
// Transaction trace tous les flux financiers d'un shop.
//----OLD----
/*
type Transaction struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type      TransactionType `gorm:"type:transaction_type;not null" json:"type"`
	ProductID *uuid.UUID      `gorm:"type:uuid" json:"productID,omitempty"`
	Product   *Product        `gorm:"foreignKey:ProductID" json:"-"`
	Quantity  int             `gorm:"not null;default:0;check:quantity >= 0" json:"quantity"`
	Amount    decimal.Decimal `gorm:"type:numeric(12,2);not null;check:amount >= 0" json:"amount"`
	ShopID    uuid.UUID       `gorm:"type:uuid;not null;index:idx_transactions_shop_created,priority:1;index:idx_transactions_shop_type,priority:1" json:"shopID"`
	Shop      Shop            `gorm:"foreignKey:ShopID" json:"-"`
	CreatedAt time.Time       `gorm:"not null;default:now();index:idx_transactions_shop_created,priority:2,sort:desc" json:"createdAt"`
}*/

// -----NEW-----
type Transaction struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type      string          `gorm:"type:varchar(20);not null" json:"type"`
	ProductID *uuid.UUID      `gorm:"type:uuid;index" json:"productID,omitempty"` // <-- Pointeur critique
	Quantity  int             `gorm:"not null;default:0" json:"quantity"`
	Amount    decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"amount"`
	ShopID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"shopID"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
