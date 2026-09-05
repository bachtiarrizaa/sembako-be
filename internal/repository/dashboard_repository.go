package repository

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetAdminStats(ctx context.Context, startDate, endDate, prevStartDate, prevEndDate time.Time) (model.DashboardStatsDTO, error)
	GetMarginAlerts(ctx context.Context) ([]model.MarginAlertSummaryDTO, error)
	GetPendingOpnameCount(ctx context.Context) (int64, error)
	GetRecentTransactions(ctx context.Context, startDate, endDate time.Time, limit int) ([]entity.Transaction, error)
	GetLowStockItems(ctx context.Context, limit int) ([]model.LowStockItemDTO, error)
	GetActiveShiftByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error)
	GetCashierShiftMetrics(ctx context.Context, shiftID string, openingBalance float64) (model.CashierShiftMetricsDTO, error)
	GetRecentTransactionsByShiftID(ctx context.Context, shiftID string, limit int) ([]entity.Transaction, error)
	GetActivePromos(ctx context.Context) ([]model.ActivePromoDTO, error)
}

type dashboardRepositoryImpl struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepositoryImpl{db: db}
}

func (r *dashboardRepositoryImpl) GetAdminStats(ctx context.Context, startDate, endDate, prevStartDate, prevEndDate time.Time) (model.DashboardStatsDTO, error) {
	var stats model.DashboardStatsDTO

	// 1. Current Period Sales & Transactions
	type SalesSummary struct {
		TotalRevenue      float64
		TotalTransactions int64
	}
	var currentSales SalesSummary
	r.db.WithContext(ctx).Table("transactions").
		Select("COALESCE(SUM(total), 0) as total_revenue, COUNT(id) as total_transactions").
		Where("status = ? AND created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Scan(&currentSales)

	stats.Revenue.Value = currentSales.TotalRevenue
	stats.Transactions.Value = float64(currentSales.TotalTransactions)

	// Previous Period Sales for Revenue % Change & Transaction Count Change
	var prevSales SalesSummary
	r.db.WithContext(ctx).Table("transactions").
		Select("COALESCE(SUM(total), 0) as total_revenue, COUNT(id) as total_transactions").
		Where("status = ? AND created_at BETWEEN ? AND ?", "completed", prevStartDate, prevEndDate).
		Scan(&prevSales)

	if prevSales.TotalRevenue > 0 {
		stats.Revenue.PercentageChange = ((currentSales.TotalRevenue - prevSales.TotalRevenue) / prevSales.TotalRevenue) * 100
	} else if currentSales.TotalRevenue > 0 {
		stats.Revenue.PercentageChange = 100.0
	} else {
		stats.Revenue.PercentageChange = 0.0
	}

	stats.Transactions.CountChange = currentSales.TotalTransactions - prevSales.TotalTransactions

	// 2. Items Sold & Top Product
	var totalItems float64
	r.db.WithContext(ctx).Table("transaction_items").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.status = ? AND transactions.created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Select("COALESCE(SUM(qty), 0)").
		Scan(&totalItems)

	type TopProduct struct {
		Name string
	}
	var topProd TopProduct
	r.db.WithContext(ctx).Table("transaction_items").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("JOIN product_units ON product_units.id = transaction_items.product_unit_id").
		Joins("JOIN products ON products.id = product_units.product_id").
		Where("transactions.status = ? AND transactions.created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Group("products.name").
		Order("SUM(qty) DESC").
		Limit(1).
		Select("products.name as name").
		Scan(&topProd)

	stats.ItemsSold.Value = totalItems
	stats.ItemsSold.TopProductName = topProd.Name

	// 3. Active Customers Count
	var activeCustCount int64
	r.db.WithContext(ctx).Table("transactions").
		Where("status = ? AND customer_id IS NOT NULL AND created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Select("COUNT(DISTINCT customer_id)").
		Scan(&activeCustCount)

	stats.ActiveCustomers.Value = float64(activeCustCount)

	// 4. Stock Value (Modal HPP Stock Fisik Gudang) & Total SKU
	var stockVal float64
	r.db.WithContext(ctx).Table("purchase_batches").
		Where("remaining_qty > 0").
		Select("COALESCE(SUM(remaining_qty * purchase_price), 0)").
		Scan(&stockVal)

	var totalSKUs int64
	r.db.WithContext(ctx).Table("products").Where("is_active = ?", true).Count(&totalSKUs)

	stats.StockValue.Value = stockVal
	stats.StockValue.TotalSKUs = totalSKUs

	// 5. Margin / Gross Profit (FIFO Costing)
	var totalCost float64
	r.db.WithContext(ctx).Table("transaction_item_cost_allocations").
		Joins("JOIN transaction_items ON transaction_items.id = transaction_item_cost_allocations.transaction_item_id").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.status = ? AND transactions.created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Select("COALESCE(SUM(transaction_item_cost_allocations.cost_subtotal), 0)").
		Scan(&totalCost)

	grossProfit := currentSales.TotalRevenue - totalCost
	stats.Margin.GrossProfit = grossProfit

	var marginPct float64
	if currentSales.TotalRevenue > 0 {
		marginPct = (grossProfit / currentSales.TotalRevenue) * 100
	}
	stats.Margin.MarginPercentage = marginPct

	return stats, nil
}

func (r *dashboardRepositoryImpl) GetMarginAlerts(ctx context.Context) ([]model.MarginAlertSummaryDTO, error) {
	var alerts []model.MarginAlertSummaryDTO

	query := `
		SELECT 
			p.name as product_name,
			ROUND(((pu.selling_price - pb.purchase_price) / pu.selling_price * 100)::numeric, 1) as current_margin,
			p.margin_threshold_percent as threshold_margin
		FROM products p
		JOIN product_units pu ON pu.product_id = p.id AND pu.is_base_unit = true
		JOIN purchase_batches pb ON pb.product_id = p.id AND pb.remaining_qty > 0
		WHERE p.is_active = true 
		  AND p.margin_threshold_percent IS NOT NULL 
		  AND ((pu.selling_price - pb.purchase_price) / pu.selling_price * 100) < p.margin_threshold_percent
		GROUP BY p.id, p.name, pu.selling_price, pb.purchase_price, p.margin_threshold_percent
		ORDER BY current_margin ASC
		LIMIT 5
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&alerts).Error
	return alerts, err
}

func (r *dashboardRepositoryImpl) GetPendingOpnameCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("stock_counts").Where("status = ?", "submitted").Count(&count).Error
	return count, err
}

func (r *dashboardRepositoryImpl) GetRecentTransactions(ctx context.Context, startDate, endDate time.Time, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Preload("Customer").
		Preload("Items.ProductUnit.Unit").
		Where("status = ? AND created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *dashboardRepositoryImpl) GetLowStockItems(ctx context.Context, limit int) ([]model.LowStockItemDTO, error) {
	var items []model.LowStockItemDTO

	query := `
		SELECT 
			p.id as id,
			p.name as product_name,
			COALESCE(SUM(s.qty_base_unit), 0) as remaining_stock,
			u.name as unit_name
		FROM products p
		LEFT JOIN stocks s ON s.product_id = p.id
		JOIN product_units pu ON pu.product_id = p.id AND pu.is_base_unit = true
		JOIN units u ON u.id = pu.unit_id
		WHERE p.is_active = true AND p.minimum_stock IS NOT NULL
		GROUP BY p.id, p.name, u.name, p.minimum_stock
		HAVING COALESCE(SUM(s.qty_base_unit), 0) <= p.minimum_stock
		ORDER BY remaining_stock ASC
		LIMIT ?
	`
	err := r.db.WithContext(ctx).Raw(query, limit).Scan(&items).Error
	return items, err
}

func (r *dashboardRepositoryImpl) GetActiveShiftByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Where("cashier_id = ? AND status = ?", cashierID, "open").
		Order("opened_at DESC").
		First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *dashboardRepositoryImpl) GetCashierShiftMetrics(ctx context.Context, shiftID string, openingBalance float64) (model.CashierShiftMetricsDTO, error) {
	var metrics model.CashierShiftMetricsDTO

	type ShiftSales struct {
		TotalRevenue      float64
		TotalTransactions int64
		CashRevenue       float64
		NonCashRevenue    float64
	}
	var sales ShiftSales

	r.db.WithContext(ctx).Table("transactions").
		Select(`
			COALESCE(SUM(total), 0) as total_revenue,
			COUNT(id) as total_transactions,
			COALESCE(SUM(CASE WHEN payment_method = 'cash' THEN total ELSE 0 END), 0) as cash_revenue,
			COALESCE(SUM(CASE WHEN payment_method != 'cash' THEN total ELSE 0 END), 0) as non_cash_revenue
		`).
		Where("shift_id = ? AND status = ?", shiftID, "completed").
		Scan(&sales)

	metrics.TotalRevenue = sales.TotalRevenue
	metrics.TotalTransactions = sales.TotalTransactions
	metrics.CashInDrawer = openingBalance + sales.CashRevenue
	metrics.NonCashTotal = sales.NonCashRevenue
	metrics.RevenueChange = 0.0

	return metrics, nil
}

func (r *dashboardRepositoryImpl) GetRecentTransactionsByShiftID(ctx context.Context, shiftID string, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Preload("Customer").
		Preload("Items.ProductUnit.Unit").
		Where("shift_id = ? AND status = ?", shiftID, "completed").
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *dashboardRepositoryImpl) GetActivePromos(ctx context.Context) ([]model.ActivePromoDTO, error) {
	var promos []model.ActivePromoDTO
	now := time.Now()

	query := `
		SELECT 
			id,
			name,
			type as discount_type,
			value::numeric as discount_value
		FROM discounts
		WHERE is_active = true
		  AND (start_date IS NULL OR start_date <= ?)
		  AND (end_date IS NULL OR end_date >= ?)
		ORDER BY created_at DESC
		LIMIT 5
	`
	err := r.db.WithContext(ctx).Raw(query, now, now).Scan(&promos).Error
	return promos, err
}
