package repository

import (
	"context"
	"sort"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
)

func (r *reportRepositoryImpl) GetProfitMarginReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.ProfitMarginReportResponse, int64, error) {
	type rawProductUnitMargin struct {
		ProductID        string  `gorm:"column:product_id"`
		ProductName      string  `gorm:"column:product_name"`
		ProductImage     *string `gorm:"column:product_image"`
		CategoryID       string  `gorm:"column:category_id"`
		CategoryName     string  `gorm:"column:category_name"`
		BaseUnitID       string  `gorm:"column:base_unit_id"`
		BaseUnitName     string  `gorm:"column:base_unit_name"`
		ProductUnitID    string  `gorm:"column:product_unit_id"`
		UnitID           string  `gorm:"column:unit_id"`
		UnitName         string  `gorm:"column:unit_name"`
		ConversionToBase float64 `gorm:"column:conversion_to_base"`
		UnitQtySold      float64 `gorm:"column:unit_qty_sold"`
		BaseQtySold      float64 `gorm:"column:base_qty_sold"`
		TotalSales       float64 `gorm:"column:total_sales"`
		TotalCost        float64 `gorm:"column:total_cost"`
		GrossProfit      float64 `gorm:"column:gross_profit"`
	}

	// 1. Ambil Summary Keseluruhan Periode
	type overallSummaryResult struct {
		TotalNetSales    float64 `gorm:"column:total_net_sales"`
		TotalCost        float64 `gorm:"column:total_cost"`
		TotalGrossProfit float64 `gorm:"column:total_gross_profit"`
	}

	summaryQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id"), startDate, endDate, cashierID, "t")

	var overallSum overallSummaryResult
	err := summaryQuery.Select(
		"COALESCE(SUM(ti.subtotal), 0) AS total_net_sales",
		"COALESCE(SUM(COALESCE(ti.total_cost, 0)), 0) AS total_cost",
		"COALESCE(SUM(COALESCE(ti.margin, ti.subtotal - COALESCE(ti.total_cost, 0))), 0) AS total_gross_profit",
	).Scan(&overallSum).Error

	if err != nil {
		return nil, 0, err
	}

	// 2. Hitung Total Produk Terjual (totalData untuk pagination)
	var totalProducts int64
	countProdQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id"), startDate, endDate, cashierID, "t")

	err = countProdQuery.Select("COUNT(DISTINCT pu.product_id)").Scan(&totalProducts).Error
	if err != nil {
		return nil, 0, err
	}

	// 3. Ambil Daftar Product ID yang Masuk ke Halaman Ini (Paginasi Produk)
	offset, safeLimit := utils.CalculateOffset(page, limit)

	type pagedProductID struct {
		ProductID   string  `gorm:"column:product_id"`
		GrossProfit float64 `gorm:"column:gross_profit"`
	}

	var pagedPIDs []pagedProductID
	pageQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id"), startDate, endDate, cashierID, "t")

	err = pageQuery.Select(
		"pu.product_id AS product_id",
		"COALESCE(SUM(COALESCE(ti.margin, ti.subtotal - COALESCE(ti.total_cost, 0))), 0) AS gross_profit",
	).
		Group("pu.product_id").
		Order("gross_profit DESC").
		Offset(offset).
		Limit(safeLimit).
		Scan(&pagedPIDs).Error

	if err != nil {
		return nil, 0, err
	}

	if len(pagedPIDs) == 0 {
		avgMarginPct := float64(0)
		if overallSum.TotalNetSales > 0 {
			avgMarginPct = (overallSum.TotalGrossProfit / overallSum.TotalNetSales) * 100
		}
		return &model.ProfitMarginReportResponse{
			Summary: model.ProfitMarginSummary{
				TotalNetSales:           overallSum.TotalNetSales,
				TotalCost:               overallSum.TotalCost,
				TotalGrossProfit:        overallSum.TotalGrossProfit,
				AverageMarginPercentage: avgMarginPct,
			},
			ProductMargins: []model.ProductMarginItem{},
		}, totalProducts, nil
	}

	targetPIDs := make([]string, 0, len(pagedPIDs))
	for _, p := range pagedPIDs {
		targetPIDs = append(targetPIDs, p.ProductID)
	}

	// 4. Ambil Detail Unit Breakdown untuk Produk Terpilih
	detailQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transaction_items ti").
		Joins("JOIN transactions t ON t.id = ti.transaction_id").
		Joins("JOIN product_units pu ON pu.id = ti.product_unit_id").
		Joins("JOIN products p ON p.id = pu.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Joins("LEFT JOIN units bu ON bu.id = p.base_unit_id").
		Joins("LEFT JOIN units u ON u.id = pu.unit_id").
		Where("p.id IN ?", targetPIDs), startDate, endDate, cashierID, "t")

	var rawRows []rawProductUnitMargin
	err = detailQuery.Select(
		"p.id AS product_id",
		"p.name AS product_name",
		"p.image AS product_image",
		"COALESCE(c.id::text, '') AS category_id",
		"COALESCE(c.name, '-') AS category_name",
		"COALESCE(bu.id::text, '') AS base_unit_id",
		"COALESCE(bu.name, '-') AS base_unit_name",
		"pu.id AS product_unit_id",
		"COALESCE(u.id::text, '') AS unit_id",
		"COALESCE(u.name, '-') AS unit_name",
		"pu.conversion_to_base AS conversion_to_base",
		"COALESCE(SUM(ti.qty), 0) AS unit_qty_sold",
		"COALESCE(SUM(ti.qty * pu.conversion_to_base), 0) AS base_qty_sold",
		"COALESCE(SUM(ti.subtotal), 0) AS total_sales",
		"COALESCE(SUM(COALESCE(ti.total_cost, 0)), 0) AS total_cost",
		"COALESCE(SUM(COALESCE(ti.margin, ti.subtotal - COALESCE(ti.total_cost, 0))), 0) AS gross_profit",
	).
		Group("p.id, p.name, p.image, c.id, c.name, bu.id, bu.name, pu.id, u.id, u.name, pu.conversion_to_base").
		Order("gross_profit DESC").
		Scan(&rawRows).Error

	if err != nil {
		return nil, 0, err
	}

	type productAgg struct {
		product     model.ReportProductDetailResponse
		qtySold     float64
		totalSales  float64
		totalCost   float64
		grossProfit float64
		units       []model.ReportProductUnitMarginResponse
	}

	aggMap := make(map[string]*productAgg)
	for _, row := range rawRows {
		unitMarginPct := float64(0)
		if row.TotalSales > 0 {
			unitMarginPct = (row.GrossProfit / row.TotalSales) * 100
		}

		unitItem := model.ReportProductUnitMarginResponse{
			ID: row.ProductUnitID,
			Unit: model.ReportItemInfoResponse{
				ID:   row.UnitID,
				Name: row.UnitName,
			},
			ConversionToBase: row.ConversionToBase,
			QtySold:          row.UnitQtySold,
			TotalSales:       row.TotalSales,
			TotalCost:        row.TotalCost,
			GrossProfit:      row.GrossProfit,
			MarginPercentage: unitMarginPct,
		}

		pAgg, exists := aggMap[row.ProductID]
		if !exists {
			pAgg = &productAgg{
				product: model.ReportProductDetailResponse{
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
				units: []model.ReportProductUnitMarginResponse{},
			}
			aggMap[row.ProductID] = pAgg
		}

		pAgg.qtySold += row.BaseQtySold
		pAgg.totalSales += row.TotalSales
		pAgg.totalCost += row.TotalCost
		pAgg.grossProfit += row.GrossProfit
		pAgg.units = append(pAgg.units, unitItem)
	}

	productMargins := make([]model.ProductMarginItem, 0, len(targetPIDs))
	for _, pid := range targetPIDs {
		pAgg, exists := aggMap[pid]
		if !exists {
			continue
		}
		pMarginPct := float64(0)
		if pAgg.totalSales > 0 {
			pMarginPct = (pAgg.grossProfit / pAgg.totalSales) * 100
		}

		productMargins = append(productMargins, model.ProductMarginItem{
			Product:          pAgg.product,
			QtySold:          pAgg.qtySold,
			TotalSales:       pAgg.totalSales,
			TotalCost:        pAgg.totalCost,
			GrossProfit:      pAgg.grossProfit,
			MarginPercentage: pMarginPct,
			Units:            pAgg.units,
		})
	}

	sort.Slice(productMargins, func(i, j int) bool {
		return productMargins[i].GrossProfit > productMargins[j].GrossProfit
	})

	avgMarginPct := float64(0)
	if overallSum.TotalNetSales > 0 {
		avgMarginPct = (overallSum.TotalGrossProfit / overallSum.TotalNetSales) * 100
	}

	return &model.ProfitMarginReportResponse{
		Summary: model.ProfitMarginSummary{
			TotalNetSales:           overallSum.TotalNetSales,
			TotalCost:               overallSum.TotalCost,
			TotalGrossProfit:        overallSum.TotalGrossProfit,
			AverageMarginPercentage: avgMarginPct,
		},
		ProductMargins: productMargins,
	}, totalProducts, nil
}
