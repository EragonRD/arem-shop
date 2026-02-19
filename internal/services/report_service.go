package services

import (
	"context"
	"strings"

	"arem-shop/internal/dto"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"

	"github.com/google/uuid"
)

type ReportService struct {
	transactionRepo   *repository.TransactionRepository
	productRepo       *repository.ProductRepository
	lowStockThreshold int
}

func NewReportService(transactionRepo *repository.TransactionRepository, productRepo *repository.ProductRepository, lowStockThreshold int) *ReportService {
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
