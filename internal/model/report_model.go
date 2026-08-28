package model

// Request Query Params Filter
type GetReportRequest struct {
	PaginationRequest
	Type      string  `form:"type" binding:"required,oneof=sales profit_margin payment_method cashier top_selling inventory_valuation"`
	StartDate *string `form:"startDate"`
	EndDate   *string `form:"endDate"`
	CashierID *string `form:"cashierId"`
}

// --------------------------------------------------------------------------
// 1. Tab 1: Laporan Penjualan (sales)
// --------------------------------------------------------------------------
type SalesReportSummary struct {
	TotalGrossSales         float64 `json:"totalGrossSales"`
	TotalDiscount           float64 `json:"totalDiscount"`
	TotalNetSales           float64 `json:"totalNetSales"`
	TotalTransactions       int64   `json:"totalTransactions"`
	AverageTransactionValue float64 `json:"averageTransactionValue"`
}

type DailySalesItem struct {
	Date              string  `json:"date"`
	GrossSales        float64 `json:"grossSales"`
	Discount          float64 `json:"discount"`
	NetSales          float64 `json:"netSales"`
	TotalTransactions int64   `json:"totalTransactions"`
}

type SalesReportResponse struct {
	Summary    SalesReportSummary `json:"summary"`
	DailySales []DailySalesItem   `json:"dailySales"`
}

// --------------------------------------------------------------------------
// 2. Tab 2: Laporan Laba Kotor FIFO (profit_margin)
// --------------------------------------------------------------------------
type ReportItemInfoResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReportProductDetailResponse struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Image    *string                `json:"image"`
	Category ReportItemInfoResponse `json:"category"`
	BaseUnit ReportItemInfoResponse `json:"baseUnit"`
}

type ReportProductUnitMarginResponse struct {
	ID               string                 `json:"id"`
	Unit             ReportItemInfoResponse `json:"unit"`
	ConversionToBase float64                `json:"conversionToBase"`
	QtySold          float64                `json:"qtySold"`
	TotalSales       float64                `json:"totalSales"`
	TotalCost        float64                `json:"totalCost"`
	GrossProfit      float64                `json:"grossProfit"`
	MarginPercentage float64                `json:"marginPercentage"`
}

type ProfitMarginSummary struct {
	TotalNetSales           float64 `json:"totalNetSales"`
	TotalCost               float64 `json:"totalCost"`
	TotalGrossProfit        float64 `json:"totalGrossProfit"`
	AverageMarginPercentage float64 `json:"averageMarginPercentage"`
}

type ProductMarginItem struct {
	Product          ReportProductDetailResponse       `json:"product"`
	QtySold          float64                           `json:"qtySold"`
	TotalSales       float64                           `json:"totalSales"`
	TotalCost        float64                           `json:"totalCost"`
	GrossProfit      float64                           `json:"grossProfit"`
	MarginPercentage float64                           `json:"marginPercentage"`
	Units            []ReportProductUnitMarginResponse `json:"units"`
}

type ProfitMarginReportResponse struct {
	Summary        ProfitMarginSummary `json:"summary"`
	ProductMargins []ProductMarginItem `json:"productMargins"`
}

// --------------------------------------------------------------------------
// 3. Tab 3: Laporan Metode Pembayaran (payment_method)
// --------------------------------------------------------------------------
type PaymentMethodSummary struct {
	TotalTransactions int64   `json:"totalTransactions"`
	TotalAmount       float64 `json:"totalAmount"`
}

type PaymentMethodItem struct {
	PaymentMethod     string  `json:"paymentMethod"`
	TotalTransactions int64   `json:"totalTransactions"`
	TotalAmount       float64 `json:"totalAmount"`
	Percentage        float64 `json:"percentage"`
}

type PaymentMethodReportResponse struct {
	Summary PaymentMethodSummary `json:"summary"`
	Methods []PaymentMethodItem  `json:"methods"`
}

// --------------------------------------------------------------------------
// 4. Tab 4: Laporan Performa & Audit Kasir (cashier)
// --------------------------------------------------------------------------
type CashierReportSummary struct {
	TotalCashiers     int64   `json:"totalCashiers"`
	TotalShifts       int64   `json:"totalShifts"`
	TotalTransactions int64   `json:"totalTransactions"`
	TotalSales        float64 `json:"totalSales"`
	TotalDiscrepancy  float64 `json:"totalDiscrepancy"`
}

type CashierPerformanceItem struct {
	CashierID               string  `json:"cashierId"`
	CashierName             string  `json:"cashierName"`
	TotalShifts             int64   `json:"totalShifts"`
	TotalTransactions       int64   `json:"totalTransactions"`
	TotalSales              float64 `json:"totalSales"`
	AverageTransactionValue float64 `json:"averageTransactionValue"`
	TotalDiscrepancy        float64 `json:"totalDiscrepancy"`
}

type CashierReportResponse struct {
	Summary  CashierReportSummary     `json:"summary"`
	Cashiers []CashierPerformanceItem `json:"cashiers"`
}

// --------------------------------------------------------------------------
// 5. Tab 5: Laporan Produk Terlaris (top_selling)
// --------------------------------------------------------------------------
type TopSellingSummary struct {
	TotalItemsSold   float64 `json:"totalItemsSold"`
	TotalSalesAmount float64 `json:"totalSalesAmount"`
}

type TopSellingProductItem struct {
	Product                     ReportProductDetailResponse `json:"product"`
	QtySold                     float64                     `json:"qtySold"`
	TotalSales                  float64                     `json:"totalSales"`
	SalesContributionPercentage float64                     `json:"salesContributionPercentage"`
}

type TopSellingReportResponse struct {
	Summary  TopSellingSummary       `json:"summary"`
	Products []TopSellingProductItem `json:"products"`
}

// --------------------------------------------------------------------------
// 6. Tab 6: Laporan Valuasi Stok Toko (inventory_valuation)
// --------------------------------------------------------------------------
type InventoryValuationSummary struct {
	TotalProducts  int64   `json:"totalProducts"`
	TotalQuantity  float64 `json:"totalQuantity"`
	TotalValuation float64 `json:"totalValuation"`
}

type InventoryValuationItem struct {
	Product        ReportProductDetailResponse `json:"product"`
	CurrentStock   float64                     `json:"currentStock"`
	AverageCost    float64                     `json:"averageCost"`
	TotalValuation float64                     `json:"totalValuation"`
}

type InventoryValuationReportResponse struct {
	Summary InventoryValuationSummary `json:"summary"`
	Items   []InventoryValuationItem  `json:"items"`
}
