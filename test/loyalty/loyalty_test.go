package loyalty_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/router"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
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
	_ = tx.AutoMigrate(&entity.LoyaltySetting{}, &entity.PointLedger{})
	return tx
}

func TestLoyaltySetting_Usecase(t *testing.T) {
	db := setupTestDB(t)
	defer db.Rollback()

	loyaltySettingRepo := repository.NewLoyaltySettingRepository(db)
	loyaltySettingUsecase := usecase.NewLoyaltySettingUsecase(loyaltySettingRepo)

	ctx := context.Background()

	// Seed default loyalty setting
	setting := entity.LoyaltySetting{
		ID:             uuid.New().String(),
		EarningRate:    10000,
		RedemptionRate: 100,
		MinimumRedeem:  50,
		IsExpiryActive: false,
		ExpiryMonths:   12,
	}
	if err := db.Create(&setting).Error; err != nil {
		t.Fatalf("failed to seed loyalty setting: %v", err)
	}

	// 1. Test Get
	res, err := loyaltySettingUsecase.Get(ctx)
	if err != nil {
		t.Fatalf("failed to get loyalty setting: %v", err)
	}
	if res.EarningRate != 10000 || res.RedemptionRate != 100 {
		t.Errorf("unexpected rates: earning=%v, redemption=%v", res.EarningRate, res.RedemptionRate)
	}

	// 2. Test Update
	updateReq := model.UpdateLoyaltySettingRequest{
		EarningRate:    15000,
		RedemptionRate: 150,
		MinimumRedeem:  100,
		IsExpiryActive: true,
		ExpiryMonths:   6,
	}
	updatedRes, err := loyaltySettingUsecase.Update(ctx, updateReq)
	if err != nil {
		t.Fatalf("failed to update loyalty setting: %v", err)
	}
	if updatedRes.EarningRate != 15000 || updatedRes.RedemptionRate != 150 || updatedRes.MinimumRedeem != 100 || !updatedRes.IsExpiryActive || updatedRes.ExpiryMonths != 6 {
		t.Errorf("unexpected updated response: %+v", updatedRes)
	}
}

func TestLoyaltySetting_API(t *testing.T) {
	db := setupTestDB(t)
	defer db.Rollback()

	jwtSecret := "testsecret123456789012345678901234567890"

	// Seed Role
	role := entity.Role{
		ID:   uuid.New().String(),
		Name: "admin_" + uuid.New().String()[:8],
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	// Fetch or Create Permissions
	var permRead, permWrite entity.Permission
	if err := db.Where("name = ?", "loyalty:read").First(&permRead).Error; err != nil {
		permRead = entity.Permission{ID: uuid.New().String(), Name: "loyalty:read", Type: "menu"}
		db.Create(&permRead)
	}
	if err := db.Where("name = ?", "loyalty:write").First(&permWrite).Error; err != nil {
		permWrite = entity.Permission{ID: uuid.New().String(), Name: "loyalty:write", Type: "menu"}
		db.Create(&permWrite)
	}

	if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?), (?, ?) ON CONFLICT DO NOTHING", role.ID, permRead.ID, role.ID, permWrite.ID).Error; err != nil {
		t.Fatalf("failed to link role permission: %v", err)
	}

	userUUID := uuid.New()
	username := "admin_" + userUUID.String()[:8]
	user := entity.User{
		ID:           userUUID.String(),
		Username:     &username,
		Email:        "admin_" + userUUID.String()[:8] + "@example.com",
		Name:         "Admin User",
		PasswordHash: "hashedpassword",
		RoleID:       role.ID,
		IsActive:     true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	token, err := utils.GenerateAccessToken(userUUID, role.Name, jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Seed Loyalty Setting
	setting := entity.LoyaltySetting{
		ID:             uuid.New().String(),
		EarningRate:    10000,
		RedemptionRate: 100,
		MinimumRedeem:  50,
		IsExpiryActive: false,
		ExpiryMonths:   12,
	}
	db.Create(&setting)

	blacklistRepo := repository.NewBlacklistRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	permUsecase := usecase.NewPermissionUsecase(db, permRepo, userRepo)

	loyaltyRepo := repository.NewLoyaltySettingRepository(db)
	loyaltyUsecase := usecase.NewLoyaltySettingUsecase(loyaltyRepo)
	loyaltyController := controller.NewLoyaltySettingController(loyaltyUsecase)

	app := gin.New()
	router.Setup(
		app,
		jwtSecret,
		blacklistRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil,
		permUsecase,
		nil, nil, nil, nil, nil, nil,
		loyaltyController,
		nil,
		nil,
	)

	// Test 1: GET /api/loyalty-settings (200 OK)
	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/loyalty-settings", nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d, body: %s", wGet.Code, wGet.Body.String())
	}

	// Test 2: PUT /api/loyalty-settings (200 OK)
	updateBody := model.UpdateLoyaltySettingRequest{
		EarningRate:    20000,
		RedemptionRate: 200,
		MinimumRedeem:  100,
		IsExpiryActive: true,
		ExpiryMonths:   24,
	}
	bodyBytes, _ := json.Marshal(updateBody)

	wPut := httptest.NewRecorder()
	reqPut, _ := http.NewRequest(http.MethodPut, "/api/loyalty-settings", bytes.NewBuffer(bodyBytes))
	reqPut.Header.Set("Content-Type", "application/json")
	reqPut.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wPut, reqPut)

	if wPut.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d, body: %s", wPut.Code, wPut.Body.String())
	}
}
