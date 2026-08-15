package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Load env from project root
	_ = godotenv.Load("../../.env")

	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	// Begin transaction to be rolled back
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to start test transaction: %v", tx.Error)
	}

	return tx
}

func TestPurchaseUsecase_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Rollback() // Rollback everything at the end of test

	purchaseBatchRepo := repository.NewPurchaseBatchRepository(db)
	productRepo := repository.NewProductRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)

	purchaseUsecase := usecase.NewPurchaseUsecase(db, purchaseBatchRepo, productRepo, supplierRepo)

	// Generate random UUIDs for clean isolation
	roleID := uuid.New().String()
	userID := uuid.New().String()
	supplierID := uuid.New().String()
	categoryID := uuid.New().String()
	unitBaseID := uuid.New().String()
	unitDusID := uuid.New().String()
	productID := uuid.New().String()

	// Seed required isolated data
	role := entity.Role{ID: roleID, Name: "admin_test_" + uuid.New().String()}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	user := entity.User{ID: userID, Name: "Test User", Email: "test_" + uuid.New().String() + "@user.com", PasswordHash: "hash", RoleID: role.ID, IsActive: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

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

	unitDus := entity.Unit{ID: unitDusID, Name: "Dus test " + uuid.New().String()}
	if err := db.Create(&unitDus).Error; err != nil {
		t.Fatalf("failed to seed dus unit: %v", err)
	}

	product := entity.Product{
		ID:         productID,
		Name:       "Test Product " + uuid.New().String(),
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
	productUnitDus := entity.ProductUnit{
		ID:               uuid.New().String(),
		ProductID:        product.ID,
		UnitID:           unitDus.ID,
		ConversionToBase: 24.0,
		SellingPrice:     110000,
		IsBaseUnit:       false,
		IsActive:         true,
	}
	if err := db.Create(&productUnitBase).Error; err != nil {
		t.Fatalf("failed to seed product unit base: %v", err)
	}
	if err := db.Create(&productUnitDus).Error; err != nil {
		t.Fatalf("failed to seed product unit dus: %v", err)
	}

	ctx := context.Background()
	var createdBatchID string

	t.Run("Create Purchase - Success", func(t *testing.T) {
		req := model.CreatePurchaseRequest{
			SupplierID:   supplier.ID,
			PurchaseDate: "2026-08-15",
			Items: []model.CreatePurchaseItem{
				{
					ProductID:     product.ID,
					UnitID:        unitDus.ID,
					Quantity:      2.0,
					PurchasePrice: 120000.0,
				},
			},
		}

		res, err := purchaseUsecase.CreatePurchase(ctx, user.ID, req)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(res) != 1 {
			t.Fatalf("expected 1 created batch, got: %d", len(res))
		}

		batch := res[0]
		createdBatchID = batch.ID
		if batch.InitialQuantity != 48.0 { // 2 Dus * 24 Pcs
			t.Errorf("expected initial quantity 48, got: %f", batch.InitialQuantity)
		}
		if batch.RemainingQuantity != 48.0 {
			t.Errorf("expected remaining quantity 48, got: %f", batch.RemainingQuantity)
		}
		if batch.PurchasePrice != 5000.0 { // 120000 / 24
			t.Errorf("expected purchase price 5000, got: %f", batch.PurchasePrice)
		}
		if batch.Unit == nil || batch.Unit.Name != unitDus.Name {
			t.Errorf("expected unit to be populated with %q, got: %+v", unitDus.Name, batch.Unit)
		}
		if batch.UnitPrice == nil || *batch.UnitPrice != 120000.0 {
			t.Errorf("expected unit price 120000, got: %v", batch.UnitPrice)
		}
		if batch.BaseUnit == nil || batch.BaseUnit.Name != unitBase.Name {
			t.Errorf("expected base unit to be populated with %q, got: %+v", unitBase.Name, batch.BaseUnit)
		}
	})

	t.Run("Create Purchase - Supplier Inactive - Failure", func(t *testing.T) {
		inactiveSupplier := entity.Supplier{ID: uuid.New().String(), Name: "Inactive Supplier " + uuid.New().String()}
		if err := db.Create(&inactiveSupplier).Error; err != nil {
			t.Fatalf("failed to seed inactive supplier: %v", err)
		}
		if err := db.Model(&inactiveSupplier).Update("is_active", false).Error; err != nil {
			t.Fatalf("failed to set supplier to inactive: %v", err)
		}

		req := model.CreatePurchaseRequest{
			SupplierID:   inactiveSupplier.ID,
			PurchaseDate: "2026-08-15",
			Items: []model.CreatePurchaseItem{
				{
					ProductID:     product.ID,
					UnitID:        unitDus.ID,
					Quantity:      2.0,
					PurchasePrice: 120000.0,
				},
			},
		}

		_, err := purchaseUsecase.CreatePurchase(ctx, user.ID, req)
		if err == nil {
			t.Errorf("expected error due to inactive supplier, got nil")
		}
	})

	t.Run("Update Purchase - Unsold - Success", func(t *testing.T) {
		// Update quantity to 3 Dus (72 Pcs)
		req := model.UpdatePurchaseRequest{
			SupplierID:    supplier.ID,
			PurchaseDate:  "2026-08-15",
			Quantity:      3.0,
			UnitID:        unitDus.ID,
			PurchasePrice: 120000.0,
		}

		updated, err := purchaseUsecase.UpdatePurchase(ctx, createdBatchID, req)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if updated.InitialQuantity != 72.0 {
			t.Errorf("expected updated initial qty 72, got: %f", updated.InitialQuantity)
		}
		if updated.RemainingQuantity != 72.0 {
			t.Errorf("expected updated remaining qty 72, got: %f", updated.RemainingQuantity)
		}
		if updated.Unit == nil || updated.Unit.Name != unitDus.Name {
			t.Errorf("expected unit to remain populated after update, got: %+v", updated.Unit)
		}
		if updated.UnitPrice == nil || *updated.UnitPrice != 120000.0 {
			t.Errorf("expected unit price 120000 after update, got: %v", updated.UnitPrice)
		}
	})

	t.Run("Update Purchase - Partially Sold - Fail modifying Qty", func(t *testing.T) {
		// Fetch the batch entity in the database
		var batch entity.PurchaseBatch
		if err := db.First(&batch, "id = ?", createdBatchID).Error; err != nil {
			t.Fatalf("failed to find batch: %v", err)
		}

		// Simulate partially sold: set remaining_qty to 50
		batch.RemainingQty = 50.0
		if err := db.Save(&batch).Error; err != nil {
			t.Fatalf("failed to update remaining qty: %v", err)
		}

		req := model.UpdatePurchaseRequest{
			SupplierID:    supplier.ID,
			PurchaseDate:  "2026-08-15",
			Quantity:      4.0,
			UnitID:        unitDus.ID,
			PurchasePrice: 120000.0,
		}

		_, err := purchaseUsecase.UpdatePurchase(ctx, createdBatchID, req)
		if err == nil {
			t.Errorf("expected error updating quantity of a partially sold batch, got nil")
		}
	})

	t.Run("Update Purchase - Partially Sold - Success modifying non-financials", func(t *testing.T) {
		invNum := "INV-NEW-123"
		req := model.UpdatePurchaseRequest{
			SupplierID:    supplier.ID,
			InvoiceNumber: &invNum,
			PurchaseDate:  "2026-08-16",
			Quantity:      3.0,
			UnitID:        unitDus.ID,
			PurchasePrice: 120000.0,
		}

		updated, err := purchaseUsecase.UpdatePurchase(ctx, createdBatchID, req)
		if err != nil {
			t.Fatalf("expected no error updating invoice of a partially sold batch, got: %v", err)
		}

		if updated.InvoiceNumber == nil || *updated.InvoiceNumber != invNum {
			t.Errorf("expected updated invoice number, got: %v", updated.InvoiceNumber)
		}
	})

	t.Run("Delete Purchase - Partially Sold - Failure", func(t *testing.T) {
		err := purchaseUsecase.DeletePurchase(ctx, createdBatchID)
		if err == nil {
			t.Errorf("expected error deleting partially sold batch, got nil")
		}
	})

	t.Run("Delete Purchase - Unsold - Success", func(t *testing.T) {
		var batch entity.PurchaseBatch
		if err := db.First(&batch, "id = ?", createdBatchID).Error; err != nil {
			t.Fatalf("failed to find batch: %v", err)
		}

		// Reset remaining qty to initial qty to simulate unsold
		batch.RemainingQty = batch.InitialQty
		if err := db.Save(&batch).Error; err != nil {
			t.Fatalf("failed to reset remaining qty: %v", err)
		}

		err := purchaseUsecase.DeletePurchase(ctx, createdBatchID)
		if err != nil {
			t.Fatalf("expected no error deleting unsold batch, got: %v", err)
		}

		var check entity.PurchaseBatch
		err = db.First(&check, "id = ?", createdBatchID).Error
		if err == nil {
			t.Errorf("expected batch to be deleted, but it was found")
		}
	})
}
