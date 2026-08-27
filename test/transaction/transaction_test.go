package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase/transaction"
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

func float64Ptr(v float64) *float64 {
	return &v
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
	transactionUsecase := transaction.NewTransactionUsecase(
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

	// Seed product 1 (tanpa diskon)
	product := entity.Product{
		CategoryID: category.ID,
		Name:       "Beras Pandan Wangi",
		BaseUnitID: unitKg.ID,
		IsActive:   true,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

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

	stock := entity.Stock{
		ProductID:   product.ID,
		QtyBaseUnit: 100.0,
	}
	if err := db.Create(&stock).Error; err != nil {
		t.Fatalf("failed to seed stock: %v", err)
	}

	// Seed FIFO Purchase Batches for Product 1
	batch1 := entity.PurchaseBatch{
		ProductID:     product.ID,
		SupplierID:    supplier.ID,
		UnitID:        &productUnit.ID,
		InitialQty:    1.0,
		RemainingQty:  1.0,
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
		RemainingQty:  50.0,
		PurchasePrice: 12000.0,
		PurchaseDate:  time.Now().AddDate(0, 0, -1),
		CreatedBy:     cashierID.String(),
	}
	if err := db.Create(&batch2).Error; err != nil {
		t.Fatalf("failed to seed batch 2: %v", err)
	}

	// Seed Product 2 (dengan Promo Diskon 10%)
	productWithDiscount := entity.Product{
		CategoryID: category.ID,
		Name:       "Beras Rojolele Diskon",
		BaseUnitID: unitKg.ID,
		IsActive:   true,
	}
	if err := db.Create(&productWithDiscount).Error; err != nil {
		t.Fatalf("failed to seed product with discount: %v", err)
	}

	productUnitWithDiscount := entity.ProductUnit{
		ProductID:        productWithDiscount.ID,
		UnitID:           unitKg.ID,
		ConversionToBase: 1.0,
		SellingPrice:     20000.0,
		IsBaseUnit:       true,
		IsActive:         true,
	}
	if err := db.Create(&productUnitWithDiscount).Error; err != nil {
		t.Fatalf("failed to seed product unit with discount: %v", err)
	}

	stockDiscount := entity.Stock{
		ProductID:   productWithDiscount.ID,
		QtyBaseUnit: 50.0,
	}
	if err := db.Create(&stockDiscount).Error; err != nil {
		t.Fatalf("failed to seed stock for discounted product: %v", err)
	}

	batchDiscount := entity.PurchaseBatch{
		ProductID:     productWithDiscount.ID,
		SupplierID:    supplier.ID,
		UnitID:        &productUnitWithDiscount.ID,
		InitialQty:    50.0,
		RemainingQty:  50.0,
		PurchasePrice: 15000.0,
		PurchaseDate:  time.Now().AddDate(0, 0, -1),
		CreatedBy:     cashierID.String(),
	}
	if err := db.Create(&batchDiscount).Error; err != nil {
		t.Fatalf("failed to seed batch for discounted product: %v", err)
	}

	// Seed Promo Diskon 10%
	promoDiscount := entity.Discount{
		Name:     "Promo Weekend 10%",
		Type:     entity.DiscountTypePercent,
		Value:    decimal.NewFromInt(10),
		IsActive: true,
	}
	if err := db.Create(&promoDiscount).Error; err != nil {
		t.Fatalf("failed to seed discount: %v", err)
	}

	productDiscountPivot := entity.ProductDiscount{
		DiscountID: promoDiscount.ID,
		ProductID:  productWithDiscount.ID,
		IsActive:   true,
	}
	if err := db.Create(&productDiscountPivot).Error; err != nil {
		t.Fatalf("failed to seed product discount pivot: %v", err)
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
	t.Run("Reject when no active shift", func(t *testing.T) {
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
	})

	// Buka shift kasir
	shift := entity.Shift{
		CashierID:      cashierID.String(),
		OpeningBalance: 100000,
		Status:         "open",
	}
	if err := db.Create(&shift).Error; err != nil {
		t.Fatalf("failed to seed shift: %v", err)
	}

	// Case 2: Checkout berhasil dengan payment method cash (tanpa diskon)
	t.Run("Successful cash transaction without discount", func(t *testing.T) {
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
		if res2.Subtotal != 30000.0 {
			t.Errorf("expected subtotal 30000, got %.2f", res2.Subtotal)
		}
		if res2.TotalDiscount != 0.0 {
			t.Errorf("expected total discount 0, got %.2f", res2.TotalDiscount)
		}
		if res2.Total != 30000.0 {
			t.Errorf("expected total 30000, got %.2f", res2.Total)
		}
		if *res2.ChangeGiven != 20000.0 {
			t.Errorf("expected change 20000, got %.2f", *res2.ChangeGiven)
		}

		// Verify detail item
		detailTrx, err := transactionUsecase.GetTransactionByID(ctx, res2.ID, cashierID.String(), "admin")
		if err != nil {
			t.Fatalf("failed to load detail: %v", err)
		}
		if len(detailTrx.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(detailTrx.Items))
		}
		item := detailTrx.Items[0]
		if item.DiscountApplied != 0.0 {
			t.Errorf("expected discountApplied 0, got %.2f", item.DiscountApplied)
		}
		if *item.TotalCost != 22000.0 {
			t.Errorf("expected TotalCost 22000, got %.2f", *item.TotalCost)
		}
		if *item.Margin != 8000.0 {
			t.Errorf("expected Margin 8000, got %.2f", *item.Margin)
		}
	})

	// Case 3: Checkout berhasil dengan Produk Diskon 10%
	t.Run("Successful transaction with active 10% discount", func(t *testing.T) {
		// Beli 2 unit @Rp 20.000. Diskon 10% = @Rp 2.000.
		// Gross Subtotal = 40.000
		// Total Discount = 4.000
		// Net Total = 36.000
		// HPP = 2 * 15.000 = 30.000
		// Margin = 36.000 - 30.000 = 6.000
		reqDiscount := model.CreateTransactionRequest{
			CustomerID:    &customer.ID,
			PaymentMethod: "cash",
			CashReceived:  float64Ptr(50000),
			Items: []model.CreateTransactionItem{
				{
					ProductUnitID: productUnitWithDiscount.ID,
					Qty:           2.0,
				},
			},
		}
		resDiscount, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqDiscount)
		if err != nil {
			t.Fatalf("expected transaction with discount to succeed, got error: %v", err)
		}
		if resDiscount.Subtotal != 40000.0 {
			t.Errorf("expected subtotal 40000, got %.2f", resDiscount.Subtotal)
		}
		if resDiscount.TotalDiscount != 4000.0 {
			t.Errorf("expected totalDiscount 4000, got %.2f", resDiscount.TotalDiscount)
		}
		if resDiscount.Total != 36000.0 {
			t.Errorf("expected net total 36000, got %.2f", resDiscount.Total)
		}
		if *resDiscount.ChangeGiven != 14000.0 {
			t.Errorf("expected change 14000, got %.2f", *resDiscount.ChangeGiven)
		}

		// Verify detail item record
		detailDiscount, err := transactionUsecase.GetTransactionByID(ctx, resDiscount.ID, cashierID.String(), "admin")
		if err != nil {
			t.Fatalf("failed to load detail: %v", err)
		}
		if len(detailDiscount.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(detailDiscount.Items))
		}
		itemDisc := detailDiscount.Items[0]
		if itemDisc.DiscountApplied != 4000.0 {
			t.Errorf("expected discountApplied 4000, got %.2f", itemDisc.DiscountApplied)
		}
		if itemDisc.Subtotal != 36000.0 {
			t.Errorf("expected item subtotal 36000, got %.2f", itemDisc.Subtotal)
		}
		if *itemDisc.TotalCost != 30000.0 {
			t.Errorf("expected item totalCost 30000, got %.2f", *itemDisc.TotalCost)
		}
		if *itemDisc.Margin != 6000.0 {
			t.Errorf("expected item margin 6000, got %.2f", *itemDisc.Margin)
		}
	})

	// Case 4: Ditolak karena stok kurang
	t.Run("Reject when stock is insufficient", func(t *testing.T) {
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
		_, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), req3)
		if err == nil {
			t.Errorf("expected error due to insufficient stock, got nil")
		}
	})

	// Case 5: 2 Produk dalam 1 transaksi, tapi Produk ke-2 kekurangan stok (Harus Rollback Total)
	t.Run("Reject and rollback when one of multiple products has insufficient stock", func(t *testing.T) {
		// Catat stok awal Product 1 dan Product 2
		stockBeforeP1, _ := stockRepo.GetByProductID(ctx, product.ID)
		stockBeforeP2, _ := stockRepo.GetByProductID(ctx, productWithDiscount.ID)

		reqMulti := model.CreateTransactionRequest{
			CustomerID:    &customer.ID,
			PaymentMethod: "cash",
			CashReceived:  float64Ptr(500000),
			Items: []model.CreateTransactionItem{
				{
					ProductUnitID: productUnit.ID, // Stok cukup (tersedia 98)
					Qty:           2.0,
				},
				{
					ProductUnitID: productUnitWithDiscount.ID, // Stok TIDAK cukup (minta 1000, tersedia 48)
					Qty:           1000.0,
				},
			},
		}

		_, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqMulti)
		if err == nil {
			t.Fatalf("expected transaction to fail due to product 2 insufficient stock, got nil")
		}

		// Pastikan stok Product 1 TETAP UTUH (tidak berkurang karena transaksi di-rollback)
		stockAfterP1, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stockAfterP1.QtyBaseUnit != stockBeforeP1.QtyBaseUnit {
			t.Errorf("expected Product 1 stock to remain %.2f after rollback, but got %.2f", stockBeforeP1.QtyBaseUnit, stockAfterP1.QtyBaseUnit)
		}

		// Pastikan stok Product 2 TETAP UTUH
		stockAfterP2, _ := stockRepo.GetByProductID(ctx, productWithDiscount.ID)
		if stockAfterP2.QtyBaseUnit != stockBeforeP2.QtyBaseUnit {
			t.Errorf("expected Product 2 stock to remain %.2f after rollback, but got %.2f", stockBeforeP2.QtyBaseUnit, stockAfterP2.QtyBaseUnit)
		}
	})

	// Case 6: Ditolak karena product unit tidak aktif
	t.Run("Reject when product unit is inactive", func(t *testing.T) {
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
		_, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), req4)
		if err == nil {
			t.Errorf("expected error due to inactive product unit, got nil")
		}
		db.Model(&productUnit).Update("is_active", true)
	})

	// Case 7: Pembayaran QRIS / Non-Tunai (CashReceived & ChangeGiven nil, ManualPaidConfirmation true)
	t.Run("Successful QRIS transaction with exact amount and confirmation", func(t *testing.T) {
		reqQRIS := model.CreateTransactionRequest{
			CustomerID:    &customer.ID,
			PaymentMethod: "qris",
			Items: []model.CreateTransactionItem{
				{
					ProductUnitID: productUnit.ID,
					Qty:           1.0,
				},
			},
		}
		resQRIS, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqQRIS)
		if err != nil {
			t.Fatalf("expected QRIS transaction to succeed, got: %v", err)
		}
		if resQRIS.PaymentMethod != "qris" {
			t.Errorf("expected paymentMethod qris, got %s", resQRIS.PaymentMethod)
		}
		if resQRIS.CashReceived != nil {
			t.Errorf("expected CashReceived to be nil for QRIS, got %v", resQRIS.CashReceived)
		}
		if resQRIS.ChangeGiven != nil {
			t.Errorf("expected ChangeGiven to be nil for QRIS, got %v", resQRIS.ChangeGiven)
		}
		if resQRIS.ManualPaidConfirmation == nil || !*resQRIS.ManualPaidConfirmation {
			t.Errorf("expected ManualPaidConfirmation to be true for QRIS")
		}
	})

	// Case 8: Penjualan Desimal / Eceran (misal 0.5 kg)
	t.Run("Successful decimal quantity transaction (0.5 kg eceran)", func(t *testing.T) {
		reqDecimal := model.CreateTransactionRequest{
			CustomerID:    &customer.ID,
			PaymentMethod: "cash",
			CashReceived:  float64Ptr(10000),
			Items: []model.CreateTransactionItem{
				{
					ProductUnitID: productUnit.ID,
					Qty:           0.5, // 0.5 kg @ Rp 15.000 = Rp 7.500
				},
			},
		}
		resDecimal, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqDecimal)
		if err != nil {
			t.Fatalf("expected decimal transaction to succeed, got: %v", err)
		}
		if resDecimal.Total != 7500.0 {
			t.Errorf("expected total 7500 for 0.5 kg, got %.2f", resDecimal.Total)
		}
		if *resDecimal.ChangeGiven != 2500.0 {
			t.Errorf("expected change 2500, got %.2f", *resDecimal.ChangeGiven)
		}
	})

	// Case 9: Multi-Satuan Produk yang Sama Dibeli Sekaligus (Cross-Unit FIFO)
	t.Run("Successful transaction with multiple units of the same product", func(t *testing.T) {
		// Buat unit baru: Karung 5kg (Conversion 5)
		karung5kg := entity.Unit{Name: "karung5kg_" + uuid.New().String()[:8]}
		db.Create(&karung5kg)

		productUnitKarung := entity.ProductUnit{
			ProductID:        product.ID,
			UnitID:           karung5kg.ID,
			ConversionToBase: 5.0,
			SellingPrice:     70000.0,
			IsBaseUnit:       false,
			IsActive:         true,
		}
		db.Create(&productUnitKarung)

		// Beli 1 Karung 5kg (= 5 kg) dan 2 Eceran 1kg (= 2 kg) -> Total 7 kg
		reqCross := model.CreateTransactionRequest{
			CustomerID:    &customer.ID,
			PaymentMethod: "cash",
			CashReceived:  float64Ptr(150000),
			Items: []model.CreateTransactionItem{
				{
					ProductUnitID: productUnitKarung.ID, // 1 * 70.000 = 70.000 (5 kg)
					Qty:           1.0,
				},
				{
					ProductUnitID: productUnit.ID, // 2 * 15.000 = 30.000 (2 kg)
					Qty:           2.0,
				},
			},
		}
		resCross, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqCross)
		if err != nil {
			t.Fatalf("expected cross-unit transaction to succeed, got: %v", err)
		}
		if resCross.Total != 100000.0 {
			t.Errorf("expected total 100000, got %.2f", resCross.Total)
		}
		if *resCross.ChangeGiven != 50000.0 {
			t.Errorf("expected change 50000, got %.2f", *resCross.ChangeGiven)
		}
	})

	// Case 10: Void Transaksi Berhasil (Stok & Batch FIFO Dikembalikan Utuh)
	t.Run("Successful Void Transaction with stock and FIFO batch restoration", func(t *testing.T) {
		// Buat transaksi baru untuk di-void
		reqToVoid := model.CreateTransactionRequest{
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

		stockBefore, _ := stockRepo.GetByProductID(ctx, product.ID)
		batchBefore, _ := purchaseBatchRepo.FindByID(ctx, batch2.ID)

		resTrx, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqToVoid)
		if err != nil {
			t.Fatalf("failed to create transaction to void: %v", err)
		}

		// Pastikan stok berkurang 2.0 setelah checkout
		stockMid, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stockMid.QtyBaseUnit != stockBefore.QtyBaseUnit-2.0 {
			t.Errorf("expected stock to decrease by 2.0, before: %.2f, mid: %.2f", stockBefore.QtyBaseUnit, stockMid.QtyBaseUnit)
		}

		// Lakukan Void Transaksi
		voidReq := model.VoidTransactionRequest{
			Reason: "Pembeli membatalkan pesanan karena salah varian",
		}
		voidRes, err := transactionUsecase.VoidTransaction(ctx, resTrx.ID, cashierID.String(), voidReq)
		if err != nil {
			t.Fatalf("expected void to succeed, got: %v", err)
		}

		// Validasi Status & Audit Trail Transaksi
		if voidRes.Status != "void" {
			t.Errorf("expected status 'void', got '%s'", voidRes.Status)
		}
		if voidRes.VoidReason == nil || *voidRes.VoidReason != voidReq.Reason {
			t.Errorf("expected voidReason '%s', got '%v'", voidReq.Reason, voidRes.VoidReason)
		}
		if voidRes.VoidedByUser == nil || voidRes.VoidedByUser.ID != cashierID.String() {
			t.Errorf("expected voidedByUser '%s', got '%v'", cashierID.String(), voidRes.VoidedByUser)
		}
		if voidRes.VoidedAt == nil {
			t.Errorf("expected voidedAt not to be nil")
		}

		// Validasi Stok Kembali Utuh
		stockAfter, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stockAfter.QtyBaseUnit != stockBefore.QtyBaseUnit {
			t.Errorf("expected stock to be restored to %.2f, got %.2f", stockBefore.QtyBaseUnit, stockAfter.QtyBaseUnit)
		}

		// Validasi Batch FIFO Kembali Utuh
		batchAfter, _ := purchaseBatchRepo.FindByID(ctx, batch2.ID)
		if batchAfter.RemainingQty != batchBefore.RemainingQty {
			t.Errorf("expected batch remaining qty to be restored to %.2f, got %.2f", batchBefore.RemainingQty, batchAfter.RemainingQty)
		}

		// Case 11: Ditolak jika transaksi sudah pernah di-void
		_, err = transactionUsecase.VoidTransaction(ctx, resTrx.ID, cashierID.String(), voidReq)
		if err == nil {
			t.Errorf("expected error when voiding already voided transaction, got nil")
		}
	})

	// Case 12: Ditolak jika Shift sudah ditutup (Closed Shift)
	t.Run("Reject Void when Shift is closed", func(t *testing.T) {
		// Buat transaksi pada shift saat ini
		reqTrx := model.CreateTransactionRequest{
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
		resTrx, err := transactionUsecase.CreateTransaction(ctx, cashierID.String(), reqTrx)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		// Tutup shift kasir
		db.Model(&shift).Update("status", "closed")

		// Coba void setelah shift tutup -> Harus Ditolak!
		voidReq := model.VoidTransactionRequest{
			Reason: "Mencoba void setelah shift ditutup",
		}
		_, err = transactionUsecase.VoidTransaction(ctx, resTrx.ID, cashierID.String(), voidReq)
		if err == nil {
			t.Errorf("expected error voiding transaction from a closed shift, got nil")
		}
	})
}
