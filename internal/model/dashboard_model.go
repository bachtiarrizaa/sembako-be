package model

import "time"

type DashboardQueryRequest struct {
	Period string `form:"period"` // "today", "week", "month"
}

type RevenueStatDTO struct {
	Value            float64 `json:"value"`
	PercentageChange float64 `json:"percentageChange"`
}

type TransactionStatDTO struct {
	Value       float64 `json:"value"`
	CountChange int64   `json:"countChange"`
}

type ItemsSoldStatDTO struct {
	Value          float64 `json:"value"`
	TopProductName string  `json:"topProductName"`
}

type ActiveCustomerStatDTO struct {
	Value float64 `json:"value"`
}

type StockValueStatDTO struct {
	Value     float64 `json:"value"`
	TotalSKUs int64   `json:"totalSKUs"`
}

type MarginStatDTO struct {
	GrossProfit      float64 `json:"grossProfit"`
	MarginPercentage float64 `json:"marginPercentage"`
}

type DashboardStatsDTO struct {
	Revenue         RevenueStatDTO        `json:"revenue"`
	Transactions    TransactionStatDTO   `json:"transactions"`
	ItemsSold       ItemsSoldStatDTO      `json:"itemsSold"`
	ActiveCustomers ActiveCustomerStatDTO `json:"activeCustomers"`
	StockValue      StockValueStatDTO     `json:"stockValue"`
	Margin          MarginStatDTO         `json:"margin"`
}

type MarginAlertSummaryDTO struct {
	ProductName     string  `json:"productName"`
	CurrentMargin   float64 `json:"currentMargin"`
	ThresholdMargin float64 `json:"thresholdMargin"`
}

type LowStockItemDTO struct {
	ID             string  `json:"id"`
	ProductName    string  `json:"productName"`
	RemainingStock float64 `json:"remainingStock"`
	UnitName       string  `json:"unitName"`
}

type AdminDashboardResponse struct {
	MarginAlerts       []MarginAlertSummaryDTO `json:"marginAlerts"`
	PendingOpnameCount int64                   `json:"pendingOpnameCount"`
	Stats              DashboardStatsDTO       `json:"stats"`
	RecentTransactions []TransactionResponse   `json:"recentTransactions"`
	LowStockItems      []LowStockItemDTO       `json:"lowStockItems"`
}

type ActiveShiftInfoDTO struct {
	ShiftID        string    `json:"shiftId"`
	OpenedAt       time.Time `json:"openedAt"`
	OpeningBalance float64   `json:"openingBalance"`
	CashierName    string    `json:"cashierName"`
}

type CashierShiftMetricsDTO struct {
	TotalRevenue      float64 `json:"totalRevenue"`
	TotalTransactions int64   `json:"totalTransactions"`
	CashInDrawer      float64 `json:"cashInDrawer"`
	NonCashTotal      float64 `json:"nonCashTotal"`
	RevenueChange     float64 `json:"revenueChange"`
}

type ActivePromoDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	DiscountType  string   `json:"discountType"`
	DiscountValue float64  `json:"discountValue"`
	MinPurchase   *float64 `json:"minPurchase"`
}

type CashierDashboardResponse struct {
	ShiftOpen          bool                    `json:"shiftOpen"`
	ActiveShift        *ActiveShiftInfoDTO     `json:"activeShift"`
	ShiftMetrics       CashierShiftMetricsDTO  `json:"shiftMetrics"`
	RecentTransactions []TransactionResponse   `json:"recentTransactions"`
	LowStockItems      []LowStockItemDTO       `json:"lowStockItems"`
	ActivePromos       []ActivePromoDTO        `json:"activePromos"`
}
