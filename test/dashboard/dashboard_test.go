package dashboard_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/router"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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
	return tx
}

func TestAdminDashboard_API(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Rollback()

	jwtSecret := "test-secret"
	blacklistRepo := repository.NewBlacklistRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	permUsecase := usecase.NewPermissionUsecase(db, permRepo, userRepo)

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo)
	dashboardController := controller.NewDashboardController(dashboardUsecase)

	var adminRole entity.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = entity.Role{ID: uuid.New().String(), Name: "admin"}
		db.Create(&adminRole)
	}

	var perm entity.Permission
	if err := db.Where("name = ?", "dashboard:read").First(&perm).Error; err != nil {
		perm = entity.Permission{ID: uuid.New().String(), Name: "dashboard:read", Type: "menu"}
		db.Create(&perm)
	}
	_ = db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", adminRole.ID, perm.ID)

	userUUID := uuid.New()
	username := "admin_dash_" + userUUID.String()[:8]
	user := entity.User{
		ID:           userUUID.String(),
		Username:     &username,
		Email:        username + "@example.com",
		Name:         "Admin User",
		PasswordHash: "hashedpassword",
		RoleID:       adminRole.ID,
		IsActive:     true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	token, err := utils.GenerateAccessToken(userUUID, adminRole.Name, jwtSecret, 15*time.Minute)
	assert.NoError(t, err)

	app := gin.New()
	router.Setup(
		app,
		jwtSecret,
		blacklistRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil,
		permUsecase,
		nil, nil, nil, nil, nil, nil, nil, nil,
		dashboardController,
	)

	// Test GET /api/dashboard (Admin View - 200 OK)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/dashboard?period=today", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
}

func TestCashierDashboard_API(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Rollback()

	jwtSecret := "test-secret"
	blacklistRepo := repository.NewBlacklistRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	permUsecase := usecase.NewPermissionUsecase(db, permRepo, userRepo)

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo)
	dashboardController := controller.NewDashboardController(dashboardUsecase)

	var cashierRole entity.Role
	if err := db.Where("name = ?", "cashier").First(&cashierRole).Error; err != nil {
		cashierRole = entity.Role{ID: uuid.New().String(), Name: "cashier"}
		db.Create(&cashierRole)
	}

	var perm entity.Permission
	if err := db.Where("name = ?", "dashboard:read").First(&perm).Error; err != nil {
		perm = entity.Permission{ID: uuid.New().String(), Name: "dashboard:read", Type: "menu"}
		db.Create(&perm)
	}
	_ = db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", cashierRole.ID, perm.ID)

	userUUID := uuid.New()
	username := "cashier_dash_" + userUUID.String()[:8]
	user := entity.User{
		ID:           userUUID.String(),
		Username:     &username,
		Email:        username + "@example.com",
		Name:         "Cashier Budi",
		PasswordHash: "hashedpassword",
		RoleID:       cashierRole.ID,
		IsActive:     true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	token, err := utils.GenerateAccessToken(userUUID, cashierRole.Name, jwtSecret, 15*time.Minute)
	assert.NoError(t, err)

	app := gin.New()
	router.Setup(
		app,
		jwtSecret,
		blacklistRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil,
		permUsecase,
		nil, nil, nil, nil, nil, nil, nil, nil,
		dashboardController,
	)

	// Test GET /api/dashboard (Cashier View without active shift - 200 OK)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.False(t, data["shiftOpen"].(bool))
}
