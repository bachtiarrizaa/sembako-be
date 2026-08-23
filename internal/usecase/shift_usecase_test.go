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

func setupShiftTestDB(t *testing.T) *gorm.DB {
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

func TestShiftUsecase_Workflow(t *testing.T) {
	db := setupShiftTestDB(t)
	defer db.Rollback()

	shiftRepo := repository.NewShiftRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	shiftUsecase := usecase.NewShiftUsecase(shiftRepo, transactionRepo)

	// Seed roles & users
	roleAdmin := entity.Role{Name: "admin_" + uuid.New().String()[:8]}
	roleCashier := entity.Role{Name: "cashier_" + uuid.New().String()[:8]}
	if err := db.Create(&roleAdmin).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}
	if err := db.Create(&roleCashier).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	cashierID := uuid.New()
	adminID := uuid.New()

	cashierUser := entity.User{
		ID:           cashierID.String(),
		Name:         "Test Cashier",
		Email:        "cashier_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleCashier.ID,
		IsActive:     true,
	}
	adminUser := entity.User{
		ID:           adminID.String(),
		Name:         "Test Admin",
		Email:        "admin_" + uuid.New().String() + "@test.com",
		PasswordHash: "hash",
		RoleID:       roleAdmin.ID,
		IsActive:     true,
	}

	if err := db.Create(&cashierUser).Error; err != nil {
		t.Fatalf("failed to seed cashier: %v", err)
	}
	if err := db.Create(&adminUser).Error; err != nil {
		t.Fatalf("failed to seed admin: %v", err)
	}

	ctx := context.Background()

	// 1. Open Shift
	openReq := model.OpenShiftRequest{OpeningBalance: 500000}
	shiftRes, err := shiftUsecase.OpenShift(ctx, cashierID, openReq)
	if err != nil {
		t.Fatalf("failed to open shift: %v", err)
	}
	if shiftRes.Status != "open" {
		t.Errorf("expected status 'open', got '%s'", shiftRes.Status)
	}

	shiftID, err := uuid.Parse(shiftRes.ID)
	if err != nil {
		t.Fatalf("failed to parse shift ID: %v", err)
	}

	// 2. Open Shift again should conflict
	_, err = shiftUsecase.OpenShift(ctx, cashierID, openReq)
	if err == nil {
		t.Errorf("expected error when opening duplicate active shift, got nil")
	}

	// 3. Get Active Shift
	activeRes, err := shiftUsecase.GetActiveShift(ctx, cashierID)
	if err != nil {
		t.Errorf("failed to get active shift: %v", err)
	}
	if activeRes.ID != shiftRes.ID {
		t.Errorf("expected active shift ID %s, got %s", shiftRes.ID, activeRes.ID)
	}

	// 4. Close Shift with discrepancy > 1000 without note should fail
	closeReqFail := model.CloseShiftRequest{ClosingBalance: 515000}
	_, err = shiftUsecase.CloseShift(ctx, shiftID, cashierID, closeReqFail)
	if err == nil {
		t.Errorf("expected error when discrepancy exceeds 1000 without note, got nil")
	}

	// 5. Close Shift with note should succeed
	note := "Found extra cash"
	closeReqSuccess := model.CloseShiftRequest{ClosingBalance: 515000, DiscrepancyNote: &note}
	closedRes, err := shiftUsecase.CloseShift(ctx, shiftID, cashierID, closeReqSuccess)
	if err != nil {
		t.Fatalf("failed to close shift: %v", err)
	}
	if closedRes.Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", closedRes.Status)
	}
}
