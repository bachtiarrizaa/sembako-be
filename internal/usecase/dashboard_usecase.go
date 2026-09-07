package usecase

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type DashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) *DashboardUsecase {
	return &DashboardUsecase{repo: repo}
}

func (u *DashboardUsecase) GetAdminDashboard(ctx context.Context, req model.DashboardQueryRequest) (*model.AdminDashboardResponse, error) {
	now := time.Now()
	startDate, endDate, prevStartDate, prevEndDate := calculateDateRange(now, req.Period)

	stats, err := u.repo.GetAdminStats(ctx, startDate, endDate, prevStartDate, prevEndDate)
	if err != nil {
		return nil, errs.NewInternal("failed to fetch dashboard statistics: " + err.Error())
	}

	marginAlerts, err := u.repo.GetMarginAlerts(ctx)
	if err != nil {
		marginAlerts = []model.MarginAlertSummaryDTO{}
	}

	pendingOpnameCount, err := u.repo.GetPendingOpnameCount(ctx)
	if err != nil {
		pendingOpnameCount = 0
	}

	rawTransactions, err := u.repo.GetRecentTransactions(ctx, startDate, endDate, 5)
	if err != nil {
		rawTransactions = nil
	}

	recentTransactions := make([]model.TransactionResponse, 0, len(rawTransactions))
	for _, trx := range rawTransactions {
		recentTransactions = append(recentTransactions, model.ToTransactionResponse(&trx))
	}

	lowStockItems, err := u.repo.GetLowStockItems(ctx, 10)
	if err != nil {
		lowStockItems = []model.LowStockItemDTO{}
	}

	res := &model.AdminDashboardResponse{
		MarginAlerts:       marginAlerts,
		PendingOpnameCount: pendingOpnameCount,
		Stats:              stats,
		RecentTransactions: recentTransactions,
		LowStockItems:      lowStockItems,
	}

	return res, nil
}

func (u *DashboardUsecase) GetCashierDashboard(ctx context.Context, cashierID string) (*model.CashierDashboardResponse, error) {
	activeShift, err := u.repo.GetActiveShiftByCashierID(ctx, cashierID)

	lowStockItems, _ := u.repo.GetLowStockItems(ctx, 5)
	if lowStockItems == nil {
		lowStockItems = []model.LowStockItemDTO{}
	}

	activePromos, _ := u.repo.GetActivePromos(ctx)
	if activePromos == nil {
		activePromos = []model.ActivePromoDTO{}
	}

	if err != nil || activeShift == nil {
		return &model.CashierDashboardResponse{
			ShiftOpen:   false,
			ActiveShift: nil,
			ShiftMetrics: model.CashierShiftMetricsDTO{
				TotalRevenue:      0,
				TotalTransactions: 0,
				CashInDrawer:      0,
				NonCashTotal:      0,
				RevenueChange:     0,
			},
			RecentTransactions: []model.TransactionResponse{},
			LowStockItems:      lowStockItems,
			ActivePromos:       activePromos,
		}, nil
	}

	shiftMetrics, _ := u.repo.GetCashierShiftMetrics(ctx, activeShift.ID, activeShift.OpeningBalance)
	rawTransactions, _ := u.repo.GetRecentTransactionsByShiftID(ctx, activeShift.ID, 5)

	recentTransactions := make([]model.TransactionResponse, 0, len(rawTransactions))
	for _, trx := range rawTransactions {
		recentTransactions = append(recentTransactions, model.ToTransactionResponse(&trx))
	}

	res := &model.CashierDashboardResponse{
		ShiftOpen: true,
		ActiveShift: &model.ActiveShiftInfoDTO{
			ShiftID:        activeShift.ID,
			OpenedAt:       activeShift.OpenedAt,
			OpeningBalance: activeShift.OpeningBalance,
			CashierName:    activeShift.Cashier.Name,
		},
		ShiftMetrics:       shiftMetrics,
		RecentTransactions: recentTransactions,
		LowStockItems:      lowStockItems,
		ActivePromos:       activePromos,
	}

	return res, nil
}

func calculateDateRange(now time.Time, period string) (startDate, endDate, prevStartDate, prevEndDate time.Time) {
	year, month, day := now.Date()
	location := now.Location()

	switch period {
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		daysToMonday := weekday - 1
		monday := time.Date(year, month, day-daysToMonday, 0, 0, 0, 0, location)

		startDate = monday
		endDate = now

		prevStartDate = monday.AddDate(0, 0, -7)
		prevEndDate = monday.Add(-time.Nanosecond)

	case "month":
		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)

		startDate = firstDayOfMonth
		endDate = now

		prevStartDate = firstDayOfMonth.AddDate(0, -1, 0)
		prevEndDate = firstDayOfMonth.Add(-time.Nanosecond)

	default: // "today"
		todayStart := time.Date(year, month, day, 0, 0, 0, 0, location)
		todayEnd := time.Date(year, month, day, 23, 59, 59, 999999999, location)

		startDate = todayStart
		endDate = todayEnd

		prevStartDate = todayStart.AddDate(0, 0, -1)
		prevEndDate = todayStart.Add(-time.Nanosecond)
	}

	return
}
