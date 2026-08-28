package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
)

func (r *reportRepositoryImpl) GetInventoryValuationReport(ctx context.Context, search string, page, limit int) (*model.InventoryValuationReportResponse, int64, error) {
	type rawValuationRow struct {
		ProductID         string  `gorm:"column:product_id"`
		ProductName       string  `gorm:"column:product_name"`
		ProductImage      *string `gorm:"column:product_image"`
		CategoryID        string  `gorm:"column:category_id"`
		CategoryName      string  `gorm:"column:category_name"`
		BaseUnitID        string  `gorm:"column:base_unit_id"`
		BaseUnitName      string  `gorm:"column:base_unit_name"`
		CurrentStock      float64 `gorm:"column:current_stock"`
		BatchValuation    float64 `gorm:"column:batch_valuation"`
		BatchRemainingQty float64 `gorm:"column:batch_remaining_qty"`
	}

	// 1. Hitung Keseluruhan Total Valuasi & Total Kuantitas (Summary)
	type summaryValuationResult struct {
		TotalProducts  int64   `gorm:"column:total_products"`
		TotalQuantity  float64 `gorm:"column:total_quantity"`
		TotalValuation float64 `gorm:"column:total_valuation"`
	}

	var sumRes summaryValuationResult
	err := r.db.WithContext(ctx).Table("purchase_batches pb").
		Joins("JOIN products p ON p.id = pb.product_id").
		Where("p.is_active = true AND pb.remaining_qty > 0").
		Select("COALESCE(SUM(pb.remaining_qty * pb.purchase_price), 0) AS total_valuation").
		Scan(&sumRes).Error
	if err != nil {
		return nil, 0, err
	}

	type countQtyResult struct {
		TotalProducts int64   `gorm:"column:total_products"`
		TotalQuantity float64 `gorm:"column:total_quantity"`
	}
	var countQty countQtyResult
	err = r.db.WithContext(ctx).Table("products p").
		Joins("LEFT JOIN stocks s ON s.product_id = p.id").
		Where("p.is_active = true").
		Select(
			"COUNT(p.id) AS total_products",
			"COALESCE(SUM(s.qty_base_unit), 0) AS total_quantity",
		).Scan(&countQty).Error
	if err != nil {
		return nil, 0, err
	}
	sumRes.TotalProducts = countQty.TotalProducts
	sumRes.TotalQuantity = countQty.TotalQuantity

	// 2. Count total products matching search filter
	var totalFilteredProducts int64
	countQuery := r.db.WithContext(ctx).Table("products p").Where("p.is_active = true")
	if search != "" {
		countQuery = countQuery.Where("p.name ILIKE ?", "%"+search+"%")
	}
	err = countQuery.Count(&totalFilteredProducts).Error
	if err != nil {
		return nil, 0, err
	}

	// 3. Paged query
	offset, safeLimit := utils.CalculateOffset(page, limit)

	detailQuery := r.db.WithContext(ctx).Table("products p").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Joins("LEFT JOIN units bu ON bu.id = p.base_unit_id").
		Joins("LEFT JOIN stocks s ON s.product_id = p.id").
		Joins("LEFT JOIN purchase_batches pb ON pb.product_id = p.id AND pb.remaining_qty > 0").
		Where("p.is_active = true")

	if search != "" {
		detailQuery = detailQuery.Where("p.name ILIKE ?", "%"+search+"%")
	}

	var rawRows []rawValuationRow
	err = detailQuery.Select(
		"p.id AS product_id",
		"p.name AS product_name",
		"p.image AS product_image",
		"COALESCE(c.id::text, '') AS category_id",
		"COALESCE(c.name, '-') AS category_name",
		"COALESCE(bu.id::text, '') AS base_unit_id",
		"COALESCE(bu.name, '-') AS base_unit_name",
		"COALESCE(s.qty_base_unit, 0) AS current_stock",
		"COALESCE(SUM(pb.remaining_qty * pb.purchase_price), 0) AS batch_valuation",
		"COALESCE(SUM(pb.remaining_qty), 0) AS batch_remaining_qty",
	).
		Group("p.id, p.name, p.image, c.id, c.name, bu.id, bu.name, s.qty_base_unit").
		Order("batch_valuation DESC, p.name ASC").
		Offset(offset).
		Limit(safeLimit).
		Scan(&rawRows).Error

	if err != nil {
		return nil, 0, err
	}

	items := make([]model.InventoryValuationItem, 0, len(rawRows))
	for _, row := range rawRows {
		avgCost := float64(0)
		if row.BatchRemainingQty > 0 {
			avgCost = row.BatchValuation / row.BatchRemainingQty
		}

		items = append(items, model.InventoryValuationItem{
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
			CurrentStock:   row.CurrentStock,
			AverageCost:    avgCost,
			TotalValuation: row.BatchValuation,
		})
	}

	return &model.InventoryValuationReportResponse{
		Summary: model.InventoryValuationSummary{
			TotalProducts:  sumRes.TotalProducts,
			TotalQuantity:  sumRes.TotalQuantity,
			TotalValuation: sumRes.TotalValuation,
		},
		Items: items,
	}, totalFilteredProducts, nil
}
