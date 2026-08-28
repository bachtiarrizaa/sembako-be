package repository

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
)

func (r *reportRepositoryImpl) GetTopSellingReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.TopSellingReportResponse, int64, error) {
	type rawTopProductRow struct {
		ProductID    string  `gorm:"column:product_id"`
		ProductName  string  `gorm:"column:product_name"`
		ProductImage *string `gorm:"column:product_image"`
		CategoryID   string  `gorm:"column:category_id"`
		CategoryName string  `gorm:"column:category_name"`
		BaseUnitID   string  `gorm:"column:base_unit_id"`
		BaseUnitName string  `gorm:"column:base_unit_name"`
		QtySold      float64 `gorm:"column:qty_sold"`
		TotalSales   float64 `gorm:"column:total_sales"`
	}

	// 1. Ambil Summary Total Items Sold & Total Sales
	type summaryResult struct {
		TotalItemsSold   float64 `gorm:"column:total_items_sold"`
		TotalSalesAmount float64 `gorm:"column:total_sales_amount"`
	}

	summaryQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id"), startDate, endDate, cashierID, "t")

	var sumRes summaryResult
	err := summaryQuery.Select(
		"COALESCE(SUM(ti.qty * pu.conversion_to_base), 0) AS total_items_sold",
		"COALESCE(SUM(ti.subtotal), 0) AS total_sales_amount",
	).Scan(&sumRes).Error

	if err != nil {
		return nil, 0, err
	}

	// 2. Hitung Total Produk Berbeda yang Terjual
	var totalProducts int64
	countQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id"), startDate, endDate, cashierID, "t")

	err = countQuery.Select("COUNT(DISTINCT pu.product_id)").Scan(&totalProducts).Error
	if err != nil {
		return nil, 0, err
	}

	// 3. Ambil Data Produk Terlaris Terurut dengan Paginasi
	offset, safeLimit := utils.CalculateOffset(page, limit)

	detailQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id").
		Joins("JOIN products p ON p.id = pu.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Joins("LEFT JOIN units bu ON bu.id = p.base_unit_id"), startDate, endDate, cashierID, "t")

	var rawRows []rawTopProductRow
	err = detailQuery.Select(
		"p.id AS product_id",
		"p.name AS product_name",
		"p.image AS product_image",
		"COALESCE(c.id::text, '') AS category_id",
		"COALESCE(c.name, '-') AS category_name",
		"COALESCE(bu.id::text, '') AS base_unit_id",
		"COALESCE(bu.name, '-') AS base_unit_name",
		"COALESCE(SUM(ti.qty * pu.conversion_to_base), 0) AS qty_sold",
		"COALESCE(SUM(ti.subtotal), 0) AS total_sales",
	).
		Group("p.id, p.name, p.image, c.id, c.name, bu.id, bu.name").
		Order("qty_sold DESC, total_sales DESC").
		Offset(offset).
		Limit(safeLimit).
		Scan(&rawRows).Error

	if err != nil {
		return nil, 0, err
	}

	products := make([]model.TopSellingProductItem, 0, len(rawRows))
	for _, row := range rawRows {
		contributionPct := float64(0)
		if sumRes.TotalSalesAmount > 0 {
			contributionPct = (row.TotalSales / sumRes.TotalSalesAmount) * 100
		}

		products = append(products, model.TopSellingProductItem{
			Product: model.ReportProductDetailResponse{
				ID:    row.ProductID,
				Name:  row.ProductName,
				Image: row.ProductImage,
				Category: model.ReportItemInfoResponse{
					ID:   row.CategoryID,
					Name: row.CategoryName,
				},
				BaseUnit: model.ReportItemInfoResponse{
					ID:   row.BaseUnitID,
					Name: row.BaseUnitName,
				},
			},
			QtySold:                     row.QtySold,
			TotalSales:                  row.TotalSales,
			SalesContributionPercentage: contributionPct,
		})
	}

	return &model.TopSellingReportResponse{
		Summary: model.TopSellingSummary{
			TotalItemsSold:   sumRes.TotalItemsSold,
			TotalSalesAmount: sumRes.TotalSalesAmount,
		},
		Products: products,
	}, totalProducts, nil
}
