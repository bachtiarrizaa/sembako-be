package repository

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
)

func (r *reportRepositoryImpl) GetPaymentMethodReport(ctx context.Context, startDate, endDate *time.Time, cashierID *string) (*model.PaymentMethodReportResponse, error) {
	type rawMethodResult struct {
		PaymentMethod     string  `gorm:"column:payment_method"`
		TotalTransactions int64   `gorm:"column:total_transactions"`
		TotalAmount       float64 `gorm:"column:total_amount"`
	}

	query := applyTrxFilters(r.db.WithContext(ctx).Table("transactions"), startDate, endDate, cashierID, "")

	var rawRows []rawMethodResult
	err := query.Select(
		"payment_method",
		"COUNT(id) AS total_transactions",
		"COALESCE(SUM(total), 0) AS total_amount",
	).
		Group("payment_method").
		Order("total_amount DESC").
		Scan(&rawRows).Error

	if err != nil {
		return nil, err
	}

	var overallTransactions int64
	var overallAmount float64
	for _, row := range rawRows {
		overallTransactions += row.TotalTransactions
		overallAmount += row.TotalAmount
	}

	methods := make([]model.PaymentMethodItem, 0, len(rawRows))
	for _, row := range rawRows {
		pct := float64(0)
		if overallAmount > 0 {
			pct = (row.TotalAmount / overallAmount) * 100
		}
		methods = append(methods, model.PaymentMethodItem{
			PaymentMethod:     row.PaymentMethod,
			TotalTransactions: row.TotalTransactions,
			TotalAmount:       row.TotalAmount,
			Percentage:        pct,
		})
	}

	return &model.PaymentMethodReportResponse{
		Summary: model.PaymentMethodSummary{
			TotalTransactions: overallTransactions,
			TotalAmount:       overallAmount,
		},
		Methods: methods,
	}, nil
}
