package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func TestStockUsecase_Opname(t *testing.T) {
	_ = godotenv.Load("../../.env")
	cfg := config.LoadConfig()
	dbConn, err := config.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	db := dbConn.Begin()
	defer db.Rollback()

	userRepo := repository.NewUserRepository(db)

	productRepo := repository.NewProductRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	stockRepo := repository.NewStockRepository(db)
	stockMutationRepo := repository.NewStockMutationRepository(db)
	stockCountRepo := repository.NewStockCountRepository(db)
	purchaseBatchRepo := repository.NewPurchaseBatchRepository(db)

	permissionUsecase := usecase.NewPermissionUsecase(db, repository.NewPermissionRepository(db), userRepo)
	stockUsecase := usecase.NewStockUsecase(db, stockRepo, stockMutationRepo, stockCountRepo, productRepo, permissionUsecase)
	purchaseUsecase := usecase.NewPurchaseUsecase(db, purchaseBatchRepo, productRepo, supplierRepo, stockRepo, stockMutationRepo)

	// Seeding roles, users, and permissions
	roleAdminID := uuid.New().String()
	roleCashierID := uuid.New().String()

	roleAdmin := entity.Role{ID: roleAdminID, Name: "Admin Role " + uuid.New().String()}
	roleCashier := entity.Role{ID: roleCashierID, Name: "Cashier Role " + uuid.New().String()}
	if err := db.Create(&roleAdmin).Error; err != nil {
		t.Fatalf("failed to seed roleAdmin: %v", err)
	}
	if err := db.Create(&roleCashier).Error; err != nil {
		t.Fatalf("failed to seed roleCashier: %v", err)
	}

	var permApprove entity.Permission
	if err := db.First(&permApprove, "name = ?", "opname:approve").Error; err != nil {
		t.Fatalf("failed to find permission: %v", err)
	}

	// Link admin role to permission
	if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleAdminID, permApprove.ID).Error; err != nil {
		t.Fatalf("failed to link role_permission: %v", err)
	}

	userAdminID := uuid.New().String()
	userCashierID := uuid.New().String()

	userAdmin := entity.User{
		ID:           userAdminID,
		Name:         "Admin User",
		Email:        "admin_" + uuid.New().String() + "@user.com",
		PasswordHash: "hash",
		RoleID:       roleAdmin.ID,
		IsActive:     true,
	}
	userCashier := entity.User{
		ID:           userCashierID,
		Name:         "Cashier User",
		Email:        "cashier_" + uuid.New().String() + "@user.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	if err := db.Create(&userAdmin).Error; err != nil {
		t.Fatalf("failed to seed userAdmin: %v", err)
	}
	if err := db.Create(&userCashier).Error; err != nil {
		t.Fatalf("failed to seed userCashier: %v", err)
	}

	supplierID := uuid.New().String()
	categoryID := uuid.New().String()
	unitBaseID := uuid.New().String()
	productID := uuid.New().String()

	supplier := entity.Supplier{ID: supplierID, Name: "Test Supplier " + uuid.New().String(), IsActive: true}
	if err := db.Create(&supplier).Error; err != nil {
		t.Fatalf("failed to seed supplier: %v", err)
	}
	category := entity.Category{ID: categoryID, Name: "Test Category " + uuid.New().String()}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}
	unitBase := entity.Unit{ID: unitBaseID, Name: "Pcs test " + uuid.New().String()}
	if err := db.Create(&unitBase).Error; err != nil {
		t.Fatalf("failed to seed base unit: %v", err)
	}

	product := entity.Product{
		ID:         productID,
		Name:       "Test Product Bimoli " + uuid.New().String(),
		CategoryID: category.ID,
		BaseUnitID: unitBase.ID,
		IsActive:   true,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

	productUnitBase := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        product.ID,
		UnitID:           unitBase.ID,
		ConversionToBase: 1.0,
		SellingPrice:     5000,
		IsBaseUnit:       true,
		IsActive:         true,
	}
	if err := db.Create(&productUnitBase).Error; err != nil {
		t.Fatalf("failed to seed product unit: %v", err)
	}

	ctx := context.Background()

	// 1. Record purchase to initialize stock
	t.Run("Initialize Stock via Purchase", func(t *testing.T) {
		req := model.CreatePurchaseRequest{
			SupplierID:   supplier.ID,
			PurchaseDate: "2026-08-15",
			Items: []model.CreatePurchaseItem{
				{
					ProductID:     product.ID,
					UnitID:        unitBase.ID,
					Quantity:      10.0,
					PurchasePrice: 4000.0, // Total = 40.000, UnitPrice = 4000
				},
			},
		}

		res, err := purchaseUsecase.CreatePurchase(ctx, userAdminID, req)
		if err != nil {
			t.Fatalf("failed to record purchase: %v", err)
		}

		if len(res) != 1 {
			t.Fatalf("expected 1 batch, got: %d", len(res))
		}

		// Check stock cache
		stock, err := stockRepo.GetByProductID(ctx, product.ID)
		if err != nil {
			t.Fatalf("failed to load stock cache: %v", err)
		}
		if stock.QtyBaseUnit != 10.0 {
			t.Errorf("expected stock 10.0, got: %f", stock.QtyBaseUnit)
		}
	})

	t.Run("Submit Opname by Cashier - Pending", func(t *testing.T) {
		req := model.SubmitStockCountRequest{
			ProductID:   product.ID,
			PhysicalQty: 8.0, // Discrepancy -2.0
			Note:        "2 pcs damaged",
		}

		res, err := stockUsecase.SubmitStockCount(ctx, userCashierID, req)
		if err != nil {
			t.Fatalf("failed to submit opname: %v", err)
		}

		if res.Status != "pending" {
			t.Errorf("expected status 'pending', got: %s", res.Status)
		}
		if res.SystemQty != 10.0 {
			t.Errorf("expected system qty 10.0, got: %f", res.SystemQty)
		}
		if res.Discrepancy != -2.0 {
			t.Errorf("expected discrepancy -2.0, got: %f", res.Discrepancy)
		}

		// Verify stock cache is still 10.0 (pending counts do not change stock)
		stock, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stock.QtyBaseUnit != 10.0 {
			t.Errorf("expected stock to remain 10.0, got: %f", stock.QtyBaseUnit)
		}
	})

	t.Run("Approve Opname (Loss) by Admin - Deduct FIFO", func(t *testing.T) {
		// Find pending opname
		counts, _, err := stockUsecase.GetStockCounts(ctx, model.GetStockCountsRequest{
			ProductID: product.ID,
			Status:    "pending",
		})
		if err != nil || len(counts) == 0 {
			t.Fatalf("failed to find pending stock count: %v", err)
		}

		targetID := counts[0].ID

		approveReq := model.ApproveStockCountRequest{
			Approve: true,
			Note:    "approved by owner",
		}

		res, err := stockUsecase.ApproveStockCount(ctx, userAdminID, targetID, approveReq)
		if err != nil {
			t.Fatalf("failed to approve stock count: %v", err)
		}

		if res.Status != "approved" {
			t.Errorf("expected status 'approved', got: %s", res.Status)
		}

		// Verify stock cache is updated to physical qty (8.0)
		stock, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stock.QtyBaseUnit != 8.0 {
			t.Errorf("expected stock to be adjusted to 8.0, got: %f", stock.QtyBaseUnit)
		}

		// Verify FIFO purchase batch remaining qty is reduced from 10 to 8
		var batch entity.PurchaseBatch
		if err := db.Order("purchase_date ASC, created_at ASC").First(&batch, "product_id = ?", product.ID).Error; err != nil {
			t.Fatalf("failed to load purchase batch: %v", err)
		}
		if batch.RemainingQty != 8.0 {
			t.Errorf("expected purchase batch remaining qty 8.0, got: %f", batch.RemainingQty)
		}

		// Verify stock mutation is logged
		mutations, _, _ := stockMutationRepo.FindProductMutations(ctx, product.ID, 1, 10)
		if len(mutations) < 2 {
			t.Fatalf("expected at least 2 mutations (1 purchase, 1 opname), got: %d", len(mutations))
		}
		latestMutation := mutations[0]
		if latestMutation.Type != "out" || latestMutation.Source != "stock_count" || latestMutation.Qty != 2.0 {
			t.Errorf("invalid latest stock mutation: %+v", latestMutation)
		}
	})

	t.Run("Submit Opname by Admin - Auto Approve", func(t *testing.T) {
		req := model.SubmitStockCountRequest{
			ProductID:   product.ID,
			PhysicalQty: 12.0, // Discrepancy +4.0 (surplus)
			Note:        "found extra boxes",
		}

		res, err := stockUsecase.SubmitStockCount(ctx, userAdminID, req)
		if err != nil {
			t.Fatalf("failed to submit opname: %v", err)
		}

		if res.Status != "approved" {
			t.Errorf("expected status 'approved' due to auto-approve, got: %s", res.Status)
		}

		// Verify stock cache is updated to 12.0
		stock, _ := stockRepo.GetByProductID(ctx, product.ID)
		if stock.QtyBaseUnit != 12.0 {
			t.Errorf("expected stock to be adjusted to 12.0, got: %f", stock.QtyBaseUnit)
		}

		// Verify a new virtual purchase batch is created for surplus
		var surplusBatch entity.PurchaseBatch
		err = db.First(&surplusBatch, "invoice_number = ?", "OPNAME-SURPLUS-"+res.ID).Error
		if err != nil {
			t.Fatalf("failed to find surplus purchase batch: %v", err)
		}
		if surplusBatch.InitialQty != 4.0 || surplusBatch.RemainingQty != 4.0 {
			t.Errorf("invalid surplus batch qty: %f/%f", surplusBatch.InitialQty, surplusBatch.RemainingQty)
		}
		if surplusBatch.PurchasePrice != 4000.0 { // should copy latest batch purchase price (4000)
			t.Errorf("expected surplus batch purchase price 4000.0, got: %f", surplusBatch.PurchasePrice)
		}
	})
}
