package report_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func setupReportTestDB(t *testing.T) *gorm.DB {
	_ = godotenv.Load("../../.env")
	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to start test transaction: %v", tx.Error)
	}
	return tx
}

func floatPtr(v float64) *float64 {
	return &v
}

func stringPtr(s string) *string {
	return &s
}

func TestReportUsecase_GetSalesReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Seed roles & user
	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	if err := db.Create(&roleCashier).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	cashierID := uuid.New().String()
	cashierUser := entity.User{
		ID:           cashierID,
		Name:         "Cashier Report Test",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	if err := db.Create(&cashierUser).Error; err != nil {
		t.Fatalf("failed to seed cashier: %v", err)
	}

	shiftID := uuid.New().String()
	shift := entity.Shift{
		ID:             shiftID,
		CashierID:      cashierID,
		OpeningBalance: 100000,
		Status:         entity.ShiftStatusOpen,
		OpenedAt:       time.Now(),
	}
	if err := db.Create(&shift).Error; err != nil {
		t.Fatalf("failed to seed shift: %v", err)
	}

	// Seed Transactions
	t1Date := time.Date(2026, 8, 27, 10, 0, 0, 0, time.Local)
	t2Date := time.Date(2026, 8, 27, 14, 0, 0, 0, time.Local)
	t3Date := time.Date(2026, 8, 28, 9, 0, 0, 0, time.Local)

	trxList := []entity.Transaction{
		{
			ID:            uuid.New().String(),
			ReceiptNumber: "TRX-TEST-001",
			CashierID:     cashierID,
			ShiftID:       shiftID,
			PaymentMethod: "cash",
			Subtotal:      100000,
			TotalDiscount: 5000,
			Total:         95000,
			Status:        "completed",
			CreatedAt:     t1Date,
			UpdatedAt:     t1Date,
		},
		{
			ID:            uuid.New().String(),
			ReceiptNumber: "TRX-TEST-002",
			CashierID:     cashierID,
			ShiftID:       shiftID,
			PaymentMethod: "qris",
			Subtotal:      200000,
			TotalDiscount: 10000,
			Total:         190000,
			Status:        "completed",
			CreatedAt:     t2Date,
			UpdatedAt:     t2Date,
		},
		{
			ID:            uuid.New().String(),
			ReceiptNumber: "TRX-TEST-003",
			CashierID:     cashierID,
			ShiftID:       shiftID,
			PaymentMethod: "cash",
			Subtotal:      50000,
			TotalDiscount: 0,
			Total:         50000,
			Status:        "completed",
			CreatedAt:     t3Date,
			UpdatedAt:     t3Date,
		},
	}

	for _, tx := range trxList {
		if err := db.Create(&tx).Error; err != nil {
			t.Fatalf("failed to seed transaction: %v", err)
		}
	}

	// Execute Test
	startDateStr := "2026-08-27"
	endDateStr := "2026-08-28"
	req := model.GetReportRequest{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
		CashierID: &cashierID,
		Type:      "sales",
	}

	res, pag, err := reportUsecase.GetSalesReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetSalesReport returned error: %v", err)
	}

	if res == nil {
		t.Fatalf("expected report result, got nil")
	}



	if res.Summary.TotalGrossSales != 350000 {
		t.Errorf("expected TotalGrossSales 350000, got %f", res.Summary.TotalGrossSales)
	}
	if res.Summary.TotalDiscount != 15000 {
		t.Errorf("expected TotalDiscount 15000, got %f", res.Summary.TotalDiscount)
	}
	if res.Summary.TotalNetSales != 335000 {
		t.Errorf("expected TotalNetSales 335000, got %f", res.Summary.TotalNetSales)
	}
	if res.Summary.TotalTransactions != 3 {
		t.Errorf("expected TotalTransactions 3, got %d", res.Summary.TotalTransactions)
	}
	if len(res.DailySales) != 2 {
		t.Errorf("expected 2 daily sales rows, got %d", len(res.DailySales))
	}
	if pag.TotalData != 2 {
		t.Errorf("expected totalData 2, got %d", pag.TotalData)
	}
}

func TestReportUsecase_GetProfitMarginReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Seed units, category, product, product units
	unitKg := entity.Unit{Name: "Unit Kg " + uuid.New().String()[:8]}
	db.Create(&unitKg)

	unitSak := entity.Unit{Name: "Unit Sak " + uuid.New().String()[:8]}
	db.Create(&unitSak)

	cat := entity.Category{Name: "Cat " + uuid.New().String()[:8]}
	db.Create(&cat)

	prodID := uuid.New().String()
	prod := entity.Product{
		ID:         prodID,
		CategoryID: cat.ID,
		BaseUnitID: unitKg.ID,
		Name:       "Beras Pandan Wangi " + uuid.New().String()[:8],
		IsActive:   true,
	}
	db.Create(&prod)

	puKg := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        prodID,
		UnitID:           unitKg.ID,
		ConversionToBase: 1,
		SellingPrice:     15000,
	}
	db.Create(&puKg)

	puSak := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        prodID,
		UnitID:           unitSak.ID,
		ConversionToBase: 5,
		SellingPrice:     70000,
	}
	db.Create(&puSak)

	// Seed User, Shift, Transaction
	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	db.Create(&roleCashier)

	cashierID := uuid.New().String()
	cashierUser := entity.User{
		ID:           cashierID,
		Name:         "Cashier Margin Test",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	db.Create(&cashierUser)

	shiftID := uuid.New().String()
	shift := entity.Shift{
		ID:             shiftID,
		CashierID:      cashierID,
		OpeningBalance: 100000,
		Status:         entity.ShiftStatusOpen,
		OpenedAt:       time.Now(),
	}
	db.Create(&shift)

	tDate := time.Date(2026, 8, 27, 10, 0, 0, 0, time.Local)
	trx := entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-MARGIN-001",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "cash",
		Subtotal:      85000,
		TotalDiscount: 0,
		Total:         85000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	}
	db.Create(&trx)

	// Transaction Items
	// Item 1: 1 Sak (5 kg) @ 70000, Cost: 50000, Margin: 20000
	cost1 := float64(50000)
	margin1 := float64(20000)
	db.Create(&entity.TransactionItem{
		TransactionID: trx.ID,
		ProductUnitID: puSak.ID,
		Qty:           1,
		UnitPrice:     70000,
		Subtotal:      70000,
		TotalCost:     &cost1,
		Margin:        &margin1,
	})

	// Item 2: 1 Kg @ 15000, Cost: 10000, Margin: 5000
	cost2 := float64(10000)
	margin2 := float64(5000)
	db.Create(&entity.TransactionItem{
		TransactionID: trx.ID,
		ProductUnitID: puKg.ID,
		Qty:           1,
		UnitPrice:     15000,
		Subtotal:      15000,
		TotalCost:     &cost2,
		Margin:        &margin2,
	})

	// Execute Test
	startDateStr := "2026-08-27"
	endDateStr := "2026-08-27"
	req := model.GetReportRequest{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
		CashierID: &cashierID,
		Type:      "profit_margin",
	}

	res, pag, err := reportUsecase.GetProfitMarginReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetProfitMarginReport error: %v", err)
	}

	if res.Summary.TotalNetSales != 85000 {
		t.Errorf("expected TotalNetSales 85000, got %f", res.Summary.TotalNetSales)
	}
	if res.Summary.TotalCost != 60000 {
		t.Errorf("expected TotalCost 60000, got %f", res.Summary.TotalCost)
	}
	if res.Summary.TotalGrossProfit != 25000 {
		t.Errorf("expected TotalGrossProfit 25000, got %f", res.Summary.TotalGrossProfit)
	}

	if len(res.ProductMargins) != 1 {
		t.Fatalf("expected 1 product margin row, got %d", len(res.ProductMargins))
	}

	prodMargin := res.ProductMargins[0]
	if prodMargin.QtySold != 6 { // 1 sak (5kg) + 1 kg = 6kg
		t.Errorf("expected Base QtySold 6, got %f", prodMargin.QtySold)
	}
	if len(prodMargin.Units) != 2 {
		t.Errorf("expected 2 units breakdown, got %d", len(prodMargin.Units))
	}

	if pag.TotalData != 1 {
		t.Errorf("expected totalData 1, got %d", pag.TotalData)
	}
}

func TestReportUsecase_GetPaymentMethodReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	db.Create(&roleCashier)

	cashierID := uuid.New().String()
	cashierUser := entity.User{
		ID:           cashierID,
		Name:         "Cashier Pay Test",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	db.Create(&cashierUser)

	shiftID := uuid.New().String()
	shift := entity.Shift{
		ID:             shiftID,
		CashierID:      cashierID,
		OpeningBalance: 100000,
		Status:         entity.ShiftStatusOpen,
		OpenedAt:       time.Now(),
	}
	db.Create(&shift)

	tDate := time.Date(2026, 8, 27, 10, 0, 0, 0, time.Local)
	db.Create(&entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-PAY-001",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "cash",
		Subtotal:      100000,
		TotalDiscount: 0,
		Total:         100000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	})
	db.Create(&entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-PAY-002",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "cash",
		Subtotal:      50000,
		TotalDiscount: 0,
		Total:         50000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	})
	db.Create(&entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-PAY-003",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "qris",
		Subtotal:      50000,
		TotalDiscount: 0,
		Total:         50000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	})

	startDateStr := "2026-08-27"
	endDateStr := "2026-08-27"
	req := model.GetReportRequest{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
		CashierID: &cashierID,
		Type:      "payment_method",
	}

	res, err := reportUsecase.GetPaymentMethodReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetPaymentMethodReport error: %v", err)
	}

	if res.Summary.TotalTransactions != 3 {
		t.Errorf("expected 3 total transactions, got %d", res.Summary.TotalTransactions)
	}
	if res.Summary.TotalAmount != 200000 {
		t.Errorf("expected TotalAmount 200000, got %f", res.Summary.TotalAmount)
	}
	if len(res.Methods) != 2 {
		t.Fatalf("expected 2 payment methods, got %d", len(res.Methods))
	}

	// Verify cash method
	if res.Methods[0].PaymentMethod != "cash" || res.Methods[0].TotalAmount != 150000 || res.Methods[0].Percentage != 75 {
		t.Errorf("unexpected cash method data: %+v", res.Methods[0])
	}
}

func TestReportUsecase_GetCashierReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	db.Create(&roleCashier)

	cashierID := uuid.New().String()
	cashierUser := entity.User{
		ID:           cashierID,
		Name:         "Cashier Audit",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	db.Create(&cashierUser)

	openTime := time.Date(2026, 8, 27, 8, 0, 0, 0, time.Local)
	closeTime := time.Date(2026, 8, 27, 17, 0, 0, 0, time.Local)
	shiftID := uuid.New().String()
	note := "Uang kurang seribu"
	shift := entity.Shift{
		ID:              shiftID,
		CashierID:       cashierID,
		OpeningBalance:  100000,
		ClosingBalance:  floatPtr(299000),
		SystemBalance:   floatPtr(300000),
		Discrepancy:     floatPtr(-1000),
		DiscrepancyNote: &note,
		Status:          entity.ShiftStatusClosed,
		OpenedAt:        openTime,
		ClosedAt:        &closeTime,
	}
	db.Create(&shift)

	tDate := time.Date(2026, 8, 27, 10, 0, 0, 0, time.Local)
	db.Create(&entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-CASHIER-001",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "cash",
		Subtotal:      200000,
		TotalDiscount: 0,
		Total:         200000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	})

	startDateStr := "2026-08-27"
	endDateStr := "2026-08-27"
	req := model.GetReportRequest{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
		CashierID: &cashierID,
		Type:      "cashier",
	}

	res, pag, err := reportUsecase.GetCashierReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetCashierReport error: %v", err)
	}

	if res.Summary.TotalCashiers != 1 {
		t.Errorf("expected 1 cashier, got %d", res.Summary.TotalCashiers)
	}
	if res.Summary.TotalShifts != 1 {
		t.Errorf("expected 1 shift, got %d", res.Summary.TotalShifts)
	}
	if res.Summary.TotalSales != 200000 {
		t.Errorf("expected TotalSales 200000, got %f", res.Summary.TotalSales)
	}
	if res.Summary.TotalDiscrepancy != -1000 {
		t.Errorf("expected TotalDiscrepancy -1000, got %f", res.Summary.TotalDiscrepancy)
	}
	if len(res.Cashiers) != 1 {
		t.Fatalf("expected 1 cashier item, got %d", len(res.Cashiers))
	}
	if pag.TotalData != 1 {
		t.Errorf("expected totalData 1, got %d", pag.TotalData)
	}
}

func TestReportUsecase_GetTopSellingReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Seed units, category, products
	unitKg := entity.Unit{Name: "Unit Kg " + uuid.New().String()[:8]}
	db.Create(&unitKg)

	cat := entity.Category{Name: "Cat " + uuid.New().String()[:8]}
	db.Create(&cat)

	prod1ID := uuid.New().String()
	prod1 := entity.Product{
		ID:         prod1ID,
		CategoryID: cat.ID,
		BaseUnitID: unitKg.ID,
		Name:       "Beras Rojolele " + uuid.New().String()[:8],
		IsActive:   true,
	}
	db.Create(&prod1)

	prod2ID := uuid.New().String()
	prod2 := entity.Product{
		ID:         prod2ID,
		CategoryID: cat.ID,
		BaseUnitID: unitKg.ID,
		Name:       "Beras Pandan Wangi " + uuid.New().String()[:8],
		IsActive:   true,
	}
	db.Create(&prod2)

	pu1 := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        prod1ID,
		UnitID:           unitKg.ID,
		ConversionToBase: 1,
		SellingPrice:     15000,
	}
	db.Create(&pu1)

	pu2 := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        prod2ID,
		UnitID:           unitKg.ID,
		ConversionToBase: 1,
		SellingPrice:     16000,
	}
	db.Create(&pu2)

	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	db.Create(&roleCashier)

	cashierID := uuid.New().String()
	cashierUser := entity.User{
		ID:           cashierID,
		Name:         "Cashier Top",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	db.Create(&cashierUser)

	shiftID := uuid.New().String()
	shift := entity.Shift{
		ID:             shiftID,
		CashierID:      cashierID,
		OpeningBalance: 100000,
		Status:         entity.ShiftStatusOpen,
		OpenedAt:       time.Now(),
	}
	db.Create(&shift)

	tDate := time.Date(2026, 8, 27, 10, 0, 0, 0, time.Local)
	trx := entity.Transaction{
		ID:            uuid.New().String(),
		ReceiptNumber: "TRX-TOP-001",
		CashierID:     cashierID,
		ShiftID:       shiftID,
		PaymentMethod: "cash",
		Subtotal:      310000,
		TotalDiscount: 0,
		Total:         310000,
		Status:        "completed",
		CreatedAt:     tDate,
		UpdatedAt:     tDate,
	}
	db.Create(&trx)

	// Prod 1 sold 10 units
	cost1 := float64(100000)
	margin1 := float64(50000)
	db.Create(&entity.TransactionItem{
		TransactionID: trx.ID,
		ProductUnitID: pu1.ID,
		Qty:           10,
		UnitPrice:     15000,
		Subtotal:      150000,
		TotalCost:     &cost1,
		Margin:        &margin1,
	})

	// Prod 2 sold 10 units
	cost2 := float64(100000)
	margin2 := float64(60000)
	db.Create(&entity.TransactionItem{
		TransactionID: trx.ID,
		ProductUnitID: pu2.ID,
		Qty:           10,
		UnitPrice:     16000,
		Subtotal:      160000,
		TotalCost:     &cost2,
		Margin:        &margin2,
	})

	startDateStr := "2026-08-27"
	endDateStr := "2026-08-27"
	req := model.GetReportRequest{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
		CashierID: &cashierID,
		Type:      "top_selling",
	}

	res, pag, err := reportUsecase.GetTopSellingReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetTopSellingReport error: %v", err)
	}

	if res.Summary.TotalItemsSold != 20 {
		t.Errorf("expected TotalItemsSold 20, got %f", res.Summary.TotalItemsSold)
	}
	if res.Summary.TotalSalesAmount != 310000 {
		t.Errorf("expected TotalSalesAmount 310000, got %f", res.Summary.TotalSalesAmount)
	}
	if len(res.Products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(res.Products))
	}
	if pag.TotalData != 2 {
		t.Errorf("expected totalData 2, got %d", pag.TotalData)
	}
}

func TestReportUsecase_GetInventoryValuationReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer db.Rollback()

	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Seed role, user, supplier, units, category, product, stock, purchase batches
	roleAdmin := entity.Role{Name: "admin_" + uuid.New().String()[:8]}
	db.Create(&roleAdmin)

	adminUser := entity.User{
		ID:           uuid.New().String(),
		Name:         "Admin Val Test",
		Email:        "admin_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleAdmin.ID,
		IsActive:     true,
	}
	db.Create(&adminUser)

	supplier := entity.Supplier{
		ID:       uuid.New().String(),
		Name:     "Supplier Val " + uuid.New().String()[:8],
		IsActive: true,
	}
	db.Create(&supplier)

	unitKg := entity.Unit{Name: "Unit Val " + uuid.New().String()[:8]}
	db.Create(&unitKg)

	cat := entity.Category{Name: "Cat Val " + uuid.New().String()[:8]}
	db.Create(&cat)

	prodID := uuid.New().String()
	uniqueName := "Beras Pandan Wangi " + uuid.New().String()[:8]
	prod := entity.Product{
		ID:         prodID,
		CategoryID: cat.ID,
		BaseUnitID: unitKg.ID,
		Name:       uniqueName,
		IsActive:   true,
	}
	db.Create(&prod)

	// Stock
	stock := entity.Stock{
		ProductID:   prodID,
		QtyBaseUnit: 100,
	}
	db.Create(&stock)

	// Batch 1: 50 kg @ 10.000 = 500.000
	db.Create(&entity.PurchaseBatch{
		ID:            uuid.New().String(),
		ProductID:     prodID,
		SupplierID:    supplier.ID,
		CreatedBy:     adminUser.ID,
		PurchaseDate:  time.Now(),
		PurchasePrice: 10000,
		InitialQty:    50,
		RemainingQty:  50,
	})

	// Batch 2: 50 kg @ 14.000 = 700.000
	db.Create(&entity.PurchaseBatch{
		ID:            uuid.New().String(),
		ProductID:     prodID,
		SupplierID:    supplier.ID,
		CreatedBy:     adminUser.ID,
		PurchaseDate:  time.Now(),
		PurchasePrice: 14000,
		InitialQty:    50,
		RemainingQty:  50,
	})

	// Execute Test with search to isolate this product
	req := model.GetReportRequest{
		PaginationRequest: model.PaginationRequest{
			Search: uniqueName,
		},
		Type: "inventory_valuation",
	}

	res, pag, err := reportUsecase.GetInventoryValuationReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetInventoryValuationReport error: %v", err)
	}

	if res == nil {
		t.Fatalf("expected response, got nil")
	}

	if res.Summary.TotalValuation < 1200000 {
		t.Errorf("expected total valuation at least 1200000, got %f", res.Summary.TotalValuation)
	}

	if len(res.Items) != 1 {
		t.Fatalf("expected 1 item for filtered product, got %d", len(res.Items))
	}

	if res.Items[0].Product.Name != uniqueName {
		t.Errorf("expected product name '%s', got %s", uniqueName, res.Items[0].Product.Name)
	}
	if res.Items[0].CurrentStock != 100 {
		t.Errorf("expected current stock 100, got %f", res.Items[0].CurrentStock)
	}
	if res.Items[0].AverageCost != 12000 {
		t.Errorf("expected average cost 12000, got %f", res.Items[0].AverageCost)
	}
	if res.Items[0].TotalValuation != 1200000 {
		t.Errorf("expected total valuation 1200000, got %f", res.Items[0].TotalValuation)
	}
	if pag.TotalData != 1 {
		t.Errorf("expected totalData 1, got %d", pag.TotalData)
	}
}
