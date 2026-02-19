package services

import (
	"context"
	"strings"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"

	"github.com/google/uuid"
)

//ajout de la personne 3
/*
type ReportService struct {
	transactionRepo   repository.TransactionRepository // <-- PLUS D'ASTÉRISQUE ICI
	productRepo       repository.ProductRepository     // <-- PLUS D'ASTÉRISQUE ICI
	lowStockThreshold int
}
*/

// pas encore le travail de la personne 2
type ReportService struct {
	transactionRepo   repository.TransactionRepository // 🟢 Ton code : Interface (PAS d'astérisque)
	productRepo       *repository.ProductRepository    // 🔴 Code Personne 2 : Pointeur de struct (AVEC astérisque)
	lowStockThreshold int
}

//--

//----old----
// func NewReportService(transactionRepo *repository.TransactionRepository, productRepo *repository.ProductRepository, lowStockThreshold int) *ReportService {

// ----new---- Personn 3
// Dans les paramètres de la fonction aussi, enlève les astérisques devant les interfaces
/* func NewReportService(transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, lowStockThreshold int) *ReportService {
if lowStockThreshold < 0 {
	lowStockThreshold = 0
}
*/
// Dans les paramètres, on remet l'astérisque uniquement pour productRepo car pas encore le travail de la personne 2
func NewReportService(transactionRepo repository.TransactionRepository, productRepo *repository.ProductRepository, lowStockThreshold int) *ReportService {
	if lowStockThreshold < 0 {
		lowStockThreshold = 0
	}
	return &ReportService{
		transactionRepo:   transactionRepo,
		productRepo:       productRepo,
		lowStockThreshold: lowStockThreshold,
	}
}

func (s *ReportService) Dashboard(ctx context.Context, shopID string) (*dto.DashboardResponse, error) {
	shopUUID, err := uuid.Parse(strings.TrimSpace(shopID))
	if err != nil {
		return nil, ErrInvalidShopID
	}

	totalSales, err := s.transactionRepo.SumAmountByShopAndTypes(ctx, shopUUID, models.TransactionSale)
	if err != nil {
		return nil, err
	}

	// Les sorties financieres regroupent Expense + Withdrawal.
	totalExpenses, err := s.transactionRepo.SumAmountByShopAndTypes(ctx, shopUUID, models.TransactionExpense, models.TransactionWithdrawal)
	if err != nil {
		return nil, err
	}

	netProfit := totalSales.Sub(totalExpenses)

	lowStockCount, err := s.productRepo.CountLowStockByShopID(ctx, shopUUID, s.lowStockThreshold)
	if err != nil {
		return nil, err
	}

	resp := &dto.DashboardResponse{
		TotalSales:       totalSales.Round(2).InexactFloat64(),
		TotalExpenses:    totalExpenses.Round(2).InexactFloat64(),
		NetProfit:        netProfit.Round(2).InexactFloat64(),
		LowStockProducts: int(lowStockCount),
		ShopID:           shopUUID.String(),
	}

	return resp, nil
}

//---------- suite à la modification de transaction_repository.go-------
/*Gemini a dit
Ne t'inquiète pas, c'est une erreur ultra-classique en Go quand on met en place une architecture propre avec des interfaces ! Tu es tombé dans le "piège du pointeur d'interface".

🕵️ L'explication du problème
En Go, une interface est déjà un type référence (sous le capot, c'est comme un pointeur).
Si tu mets un astérisque * devant une interface (ex: *TransactionRepository), tu crées un "pointeur de pointeur", ce qui casse tout et empêche Go de trouver les méthodes. Les lignes rouges viennent toutes de là !
-*

//---- ci dessous ça a été autogénéré mais pas sur que ce soit correct (au dessus oui)----

// suite à la modification de transaction_repository.go, il faut aussi adapter les méthodes du service pour utiliser le struct privé transactionRepository et non l'interface TransactionRepository. Voici les modifications à apporter dans report_service.go :

/*
Dans report_service.go, tu dois remplacer les types de transactionRepo et productRepo pour qu'ils utilisent les structs privés transactionRepository et productRepository au lieu des interfaces TransactionRepository et ProductRepository.

Voici les modifications à apporter :

1. Importer les packages nécessaires pour les structs privés :
   - "arem-shop/internal/repository"
2. Modifier la struct ReportService pour utiliser les structs privés :
type ReportService struct {
	transactionRepo   *repository.transactionRepository
	productRepo       *repository.productRepository
	lowStockThreshold int
}
*/
