package repository

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetSalesReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.SalesReportResponse, int64, error)
	GetProfitMarginReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.ProfitMarginReportResponse, int64, error)
	GetPaymentMethodReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string) (*model.PaymentMethodReportResponse, error)
	GetCashierReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.CashierReportResponse, int64, error)
	GetTopSellingReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.TopSellingReportResponse, int64, error)
	GetInventoryValuationReport(ctx context.Context, search string, page, limit int) (*model.InventoryValuationReportResponse, int64, error)
}

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepositoryImpl{db: db}
}

// Helper untuk filter transaksi standar
func applyTrxFilters(db *gorm.DB, startDate, endDate *time.Time, cashierID *string, prefix string) *gorm.DB {
	p := ""
	if prefix != "" {
		p = prefix + "."
	}
	db = db.Where(p+"status = ?", "completed")
	if startDate != nil {
		db = db.Where(p+"created_at >= ?", *startDate)
	}
	if endDate != nil {
		db = db.Where(p+"created_at <= ?", *endDate)
	}
	if cashierID != nil && *cashierID != "" {
		db = db.Where(p+"cashier_id = ?", *cashierID)
	}
	return db
}
