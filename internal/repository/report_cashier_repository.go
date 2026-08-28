package repository

import (
	"context"
	"sort"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
)

func (r *reportRepositoryImpl) GetCashierReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string, page, limit int) (*model.CashierReportResponse, int64, error) {
	type rawTxAgg struct {
		CashierID         string  `gorm:"column:cashier_id"`
		CashierName       string  `gorm:"column:cashier_name"`
		TotalTransactions int64   `gorm:"column:total_transactions"`
		TotalSales        float64 `gorm:"column:total_sales"`
	}

	txQuery := applyTrxFilters(r.db.WithContext(ctx).Table("transactions t").
		Joins("JOIN users u ON u.id = t.cashier_id"), startDate, endDate, cashierID, "t")

	var rawTxList []rawTxAgg
	err := txQuery.Select(
		"t.cashier_id",
		"u.name AS cashier_name",
		"COUNT(t.id) AS total_transactions",
		"COALESCE(SUM(t.total), 0) AS total_sales",
	).
		Group("t.cashier_id, u.name").
		Scan(&rawTxList).Error

	if err != nil {
		return nil, 0, err
	}

	type rawShiftAgg struct {
		CashierID        string  `gorm:"column:cashier_id"`
		CashierName      string  `gorm:"column:cashier_name"`
		TotalShifts      int64   `gorm:"column:total_shifts"`
		TotalDiscrepancy float64 `gorm:"column:total_discrepancy"`
	}

	shiftQuery := r.db.WithContext(ctx).Table("shifts s").
		Joins("JOIN users u ON u.id = s.cashier_id")

	if startDate != nil {
		shiftQuery = shiftQuery.Where("s.opened_at >= ?", *startDate)
	}
	if endDate != nil {
		shiftQuery = shiftQuery.Where("s.opened_at <= ?", *endDate)
	}
	if cashierID != nil && *cashierID != "" {
		shiftQuery = shiftQuery.Where("s.cashier_id = ?", *cashierID)
	}

	var rawShiftList []rawShiftAgg
	err = shiftQuery.Select(
		"s.cashier_id",
		"u.name AS cashier_name",
		"COUNT(s.id) AS total_shifts",
		"COALESCE(SUM(s.discrepancy), 0) AS total_discrepancy",
	).
		Group("s.cashier_id, u.name").
		Scan(&rawShiftList).Error

	if err != nil {
		return nil, 0, err
	}

	type cashierData struct {
		id           string
		name         string
		shifts       int64
		transactions int64
		sales        float64
		discrepancy  float64
	}

	cashierMap := make(map[string]*cashierData)
	var cashierOrder []string

	for _, tx := range rawTxList {
		cData, exists := cashierMap[tx.CashierID]
		if !exists {
			cData = &cashierData{
				id:   tx.CashierID,
				name: tx.CashierName,
			}
			cashierMap[tx.CashierID] = cData
			cashierOrder = append(cashierOrder, tx.CashierID)
		}
		cData.transactions += tx.TotalTransactions
		cData.sales += tx.TotalSales
	}

	for _, sh := range rawShiftList {
		cData, exists := cashierMap[sh.CashierID]
		if !exists {
			cData = &cashierData{
				id:   sh.CashierID,
				name: sh.CashierName,
			}
			cashierMap[sh.CashierID] = cData
			cashierOrder = append(cashierOrder, sh.CashierID)
		}
		cData.shifts += sh.TotalShifts
		cData.discrepancy += sh.TotalDiscrepancy
	}

	var totalShifts int64
	var totalTx int64
	var totalSales float64
	var totalDiscrepancy float64

	allCashiers := make([]model.CashierPerformanceItem, 0, len(cashierMap))
	for _, cid := range cashierOrder {
		cData := cashierMap[cid]
		totalShifts += cData.shifts
		totalTx += cData.transactions
		totalSales += cData.sales
		totalDiscrepancy += cData.discrepancy

		avgTx := float64(0)
		if cData.transactions > 0 {
			avgTx = cData.sales / float64(cData.transactions)
		}

		allCashiers = append(allCashiers, model.CashierPerformanceItem{
			CashierID:               cData.id,
			CashierName:             cData.name,
			TotalShifts:             cData.shifts,
			TotalTransactions:       cData.transactions,
			TotalSales:              cData.sales,
			AverageTransactionValue: avgTx,
			TotalDiscrepancy:        cData.discrepancy,
		})
	}

	sort.Slice(allCashiers, func(i, j int) bool {
		return allCashiers[i].TotalSales > allCashiers[j].TotalSales
	})

	totalCashiersCount := int64(len(allCashiers))

	offset, safeLimit := utils.CalculateOffset(page, limit)
	pagedCashiers := []model.CashierPerformanceItem{}
	if offset < len(allCashiers) {
		end := offset + safeLimit
		if end > len(allCashiers) {
			end = len(allCashiers)
		}
		pagedCashiers = allCashiers[offset:end]
	}

	return &model.CashierReportResponse{
		Summary: model.CashierReportSummary{
			TotalCashiers:     totalCashiersCount,
			TotalShifts:       totalShifts,
			TotalTransactions: totalTx,
			TotalSales:        totalSales,
			TotalDiscrepancy:  totalDiscrepancy,
		},
		Cashiers: pagedCashiers,
	}, totalCashiersCount, nil
}
