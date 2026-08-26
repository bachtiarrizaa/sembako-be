package usecase_test

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

func setupTransactionTestDB(t *testing.T) *gorm.DB {
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

func TestTransactionUsecase_CreateTransaction(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Rollback()

	// Initialize repos & usecases
	shiftRepo := repository.NewShiftRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	productUnitRepo := repository.NewProductUnitRepository(db)
	productRepo := repository.NewProductRepository(db)
	stockRepo := repository.NewStockRepository(db)
	stockMutationRepo := repository.NewStockMutationRepository(db)
	purchaseBatchRepo := repository.NewPurchaseBatchRepository(db)
	costAllocationRepo := repository.NewTransactionItemCostAllocationRepository(db)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionUsecase(
		db,
		transactionRepo,
		shiftRepo,
		customerRepo,
		productUnitRepo,
		productRepo,
		stockRepo,
		stockMutationRepo,
		purchaseBatchRepo,
		costAllocationRepo,
	)

	// Seed roles & users
	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	if err := db.Create(&roleCashier).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	cashierID := uuid.New()
	cashierUser := entity.User{
		ID:           cashierID.String(),
		Name:         "Test Cashier",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	if err := db.Create(&cashierUser).Error; err != nil {
		t.Fatalf("failed to seed cashier: %v", err)
	}

	// Seed supplier
	supplier := entity.Supplier{
		Name: "Test Supplier",
	}
	if err := db.Create(&supplier).Error; err != nil {
		t.Fatalf("failed to seed supplier: %v", err)
	}

	// Seed category & unit
	category := entity.Category{Name: "Beras_" + uuid.New().String()[:8]}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	unitKg := entity.Unit{Name: "kg_" + uuid.New().String()[:8]}
	if err := db.Create(&unitKg).Error; err != nil {
		t.Fatalf("failed to seed unit: %v", err)
	}

	// Seed product
	product := entity.Product{
		CategoryID: category.ID,
		Name:       "Beras Pandan Wangi",
		BaseUnitID: unitKg.ID,
		IsActive:   true,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

	// Seed product unit
	productUnit := entity.ProductUnit{
		ProductID:        product.ID,
		UnitID:           unitKg.ID,
		ConversionToBase: 1.0,
		SellingPrice:     15000.0,
		IsBaseUnit:       true,
		IsActive:         true,
	}
	if err := db.Create(&productUnit).Error; err != nil {
		t.Fatalf("failed to seed product unit: %v", err)
	}

	// Seed stock
	stock := entity.Stock{
		ProductID:   product.ID,
		QtyBaseUnit: 100.0,
	}
	if err := db.Create(&stock).Error; err != nil {
		t.Fatalf("failed to seed stock: %v", err)
	}

	// Seed FIFO Purchase Batches
	batch1 := entity.PurchaseBatch{
		ProductID:     product.ID,
		SupplierID:    supplier.ID,
		UnitID:        &productUnit.ID,
		InitialQty:    1.0,
		RemainingQty:  1.0, // Only 1 kg left on older batch
		PurchasePrice: 10000.0,
		PurchaseDate:  time.Now().AddDate(0, 0, -2),
		CreatedBy:     cashierID.String(),
	}
	if err := db.Create(&batch1).Error; err != nil {
		t.Fatalf("failed to seed batch 1: %v", err)
	}

	batch2 := entity.PurchaseBatch{
		ProductID:     product.ID,
		SupplierID:    supplier.ID,
		UnitID:        &productUnit.ID,
		InitialQty:    50.0,
		RemainingQty:  50.0, // 50 kg on newer batch
		PurchasePrice: 12000.0,
		PurchaseDate:  time.Now().AddDate(0, 0, -1),
		CreatedBy:     cashierID.String(),
	}
	if err := db.Create(&batch2).Error; err != nil {
		t.Fatalf("failed to seed batch 2: %v", err)
	}

	// Seed customer
	customer := entity.Customer{
		Name:        "John Doe",
		PhoneNumber: "08123456789",
		Address:     "Jakarta",
		IsActive:    true,
	}
	if err := db.Create(&customer).Error; err != nil {
		t.Fatalf("failed to seed customer: %v", err)
	}

	ctx := context.Background()

	// Case 1: Ditolak karena tidak ada shift aktif
	req1 := model.CreateTransactionRequest{
		CustomerID:    &customer.ID,
		PaymentMethod: "cash",
		CashReceived:  float64Ptr(50000),
		Items: []model.CreateTransactionItem{
			{
				ProductUnitID: productUnit.ID,
				Qty:           2.0,
			},
		},
	}
	_, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), req1)
	if err == nil {
		t.Errorf("expected error due to inactive shift, got nil")
	}

	// Buka shift kasir
	shift := entity.Shift{
		CashierID:      cashierID.String(),
		OpeningBalance: 100000,
		Status:         "open",
	}
	if err := db.Create(&shift).Error; err != nil {
		t.Fatalf("failed to seed shift: %v", err)
	}

	// Case 2: Checkout berhasil dengan payment method cash
	// Qty = 2.0. FIFO logic should allocate:
	// - 1.0 from batch1 (price 10,000) -> remaining 0
	// - 1.0 from batch2 (price 12,000) -> remaining 49
	// Total Cost = (1 * 10,000) + (1 * 12,000) = 22,000
	// Margin = 30,000 - 22,000 = 8,000
	req2 := model.CreateTransactionRequest{
		CustomerID:    &customer.ID,
		PaymentMethod: "cash",
		CashReceived:  float64Ptr(50000),
		Items: []model.CreateTransactionItem{
			{
				ProductUnitID: productUnit.ID,
				Qty:           2.0,
			},
		},
	}
	res2, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), req2)
	if err != nil {
		t.Fatalf("expected transaction success, got error: %v", err)
	}
	if res2.Total != 30000.0 {
		t.Errorf("expected total 30000, got %.2f", res2.Total)
	}
	if *res2.ChangeGiven != 20000.0 {
		t.Errorf("expected change 20000, got %.2f", *res2.ChangeGiven)
	}

	// Load detail again to check populated items with cost and margin
	detailTrx, err := transactionUsecase.GetTransactionByID(ctx, res2.ID, cashierID.String(), "cashier")
	if err != nil {
		t.Fatalf("failed to load detail: %v", err)
	}

	if len(detailTrx.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(detailTrx.Items))
	}
	item := detailTrx.Items[0]
	if *item.TotalCost != 22000.0 {
		t.Errorf("expected HPP/TotalCost 22000, got %.2f", *item.TotalCost)
	}
	if *item.Margin != 8000.0 {
		t.Errorf("expected Margin 8000, got %.2f", *item.Margin)
	}

	// Pastikan stok berkurang
	updatedStock, err := stockRepo.GetByProductID(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to get updated stock: %v", err)
	}
	if updatedStock.QtyBaseUnit != 98.0 {
		t.Errorf("expected stock to be 98, got %.2f", updatedStock.QtyBaseUnit)
	}

	// Pastikan sisa batch berkurang
	dbBatch1, _ := purchaseBatchRepo.FindByID(ctx, batch1.ID)
	if dbBatch1.RemainingQty != 0.0 {
		t.Errorf("expected batch1 remaining qty to be 0, got %.2f", dbBatch1.RemainingQty)
	}

	dbBatch2, _ := purchaseBatchRepo.FindByID(ctx, batch2.ID)
	if dbBatch2.RemainingQty != 49.0 {
		t.Errorf("expected batch2 remaining qty to be 49, got %.2f", dbBatch2.RemainingQty)
	}

	// Case 3: Ditolak karena stok kurang
	req3 := model.CreateTransactionRequest{
		CustomerID:    &customer.ID,
		PaymentMethod: "cash",
		CashReceived:  float64Ptr(2000000),
		Items: []model.CreateTransactionItem{
			{
				ProductUnitID: productUnit.ID,
				Qty:           200.0,
			},
		},
	}
	_, err = transactionUsecase.CreateTransaction(ctx, cashierID.String(), req3)
	if err == nil {
		t.Errorf("expected error due to insufficient stock, got nil")
	}

	// Case 4: Ditolak karena product unit tidak aktif
	db.Model(&productUnit).Update("is_active", false)
	req4 := model.CreateTransactionRequest{
		CustomerID:    &customer.ID,
		PaymentMethod: "cash",
		CashReceived:  float64Ptr(50000),
		Items: []model.CreateTransactionItem{
			{
				ProductUnitID: productUnit.ID,
				Qty:           1.0,
			},
		},
	}
	_, err = transactionUsecase.CreateTransaction(ctx, cashierID.String(), req4)
	if err == nil {
		t.Errorf("expected error due to inactive product unit, got nil")
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}
