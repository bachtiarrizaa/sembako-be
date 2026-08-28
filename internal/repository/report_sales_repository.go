package repository

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
)

func (r *reportRepositoryImpl) GetSalesReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.SalesReportResponse, int64, error) {
	baseQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transactions"), startDate, endDate, cashierID, "")

	// 1. Ambil Data Summary (seluruh periode)
	type summaryResult struct {
		TotalGrossSales   float64 `gorm:"column:total_gross_sales"`
		TotalDiscount     float64 `gorm:"column:total_discount"`
		TotalNetSales     float64 `gorm:"column:total_net_sales"`
		TotalTransactions int64   `gorm:"column:total_transactions"`
	}

	var sumRes summaryResult
	err := baseQuery.Select(
		"COALESCE(SUM(subtotal), 0) AS total_gross_sales",
		"COALESCE(SUM(total_discount + points_discount_value), 0) AS total_discount",
		"COALESCE(SUM(total), 0) AS total_net_sales",
		"COUNT(id) AS total_transactions",
	).Scan(&sumRes).Error

	if err != nil {
		return nil, 0, err
	}

	avgTx := float64(0)
	if sumRes.TotalTransactions > 0 {
		avgTx = sumRes.TotalNetSales / float64(sumRes.TotalTransactions)
	}

	summary := model.SalesReportSummary{
		TotalGrossSales:         sumRes.TotalGrossSales,
		TotalDiscount:           sumRes.TotalDiscount,
		TotalNetSales:           sumRes.TotalNetSales,
		TotalTransactions:       sumRes.TotalTransactions,
		AverageTransactionValue: avgTx,
	}

	// 2. Hitung Total Hari yang Memiliki Transaksi (totalData untuk pagination)
	var totalDays int64
	countQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transactions"), startDate, endDate, cashierID, "")
	err = countQuery.Select("COUNT(DISTINCT TO_CHAR(created_at, 'YYYY-MM-DD'))").Scan(&totalDays).Error
	if err != nil {
		return nil, 0, err
	}

	// 3. Ambil Rincian Penjualan Harian dengan Paginasi
	var dailySales []model.DailySalesItem
	dailyQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transactions"), startDate, endDate, cashierID, "")
	offset, safeLimit := utils.CalculateOffset(page, limit)

	err = dailyQuery.Select(
		"TO_CHAR(created_at, 'YYYY-MM-DD') AS date",
		"COALESCE(SUM(subtotal), 0) AS gross_sales",
		"COALESCE(SUM(total_discount + points_discount_value), 0) AS discount",
		"COALESCE(SUM(total), 0) AS net_sales",
		"COUNT(id) AS total_transactions",
	).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date DESC").
		Offset(offset).
		Limit(safeLimit).
		Scan(&dailySales).Error

	if err != nil {
		return nil, 0, err
	}

	if dailySales == nil {
		dailySales = []model.DailySalesItem{}
	}

	return &model.SalesReportResponse{
		Summary:    summary,
		DailySales: dailySales,
	}, totalDays, nil
}
