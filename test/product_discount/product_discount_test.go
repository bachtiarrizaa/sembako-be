package product_discount_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/router"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
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

// 1. UNIT TEST: Usecase Layer
func TestProductDiscount_Usecase(t *testing.T) {
	db := setupTestDB(t)
	defer db.Rollback()

	productDiscountRepo := repository.NewProductDiscountRepository(db)
	productDiscountUsecase := usecase.NewProductDiscountUsecase(productDiscountRepo)

	ctx := context.Background()

	// Seed Category, Unit, Product & Discount
	category := entity.Category{Name: "Cat_" + uuid.New().String()[:8]}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	unit := entity.Unit{Name: "Unit_" + uuid.New().String()[:8]}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("failed to seed unit: %v", err)
	}

	product := entity.Product{
		CategoryID: category.ID,
		Name:       "Product_" + uuid.New().String()[:8],
		BaseUnitID: unit.ID,
		IsActive:   true,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(24 * 7 * time.Hour)
	discount := entity.Discount{
		Name:      "Promo_" + uuid.New().String()[:8],
		Type:      "percent",
		Value:     decimal.NewFromFloat(10.0),
		StartDate: &startDate,
		EndDate:   &endDate,
		IsActive:  true,
	}
	if err := db.Create(&discount).Error; err != nil {
		t.Fatalf("failed to seed discount: %v", err)
	}

	// Skenario 1: Create Product Discount Berhasil
	req := model.CreateProductDiscountRequest{
		ProductID:  product.ID,
		DiscountID: discount.ID,
	}

	res, err := productDiscountUsecase.Create(ctx, req)
	if err != nil {
		t.Fatalf("expected create to succeed, got: %v", err)
	}
	if res.Product.ID != product.ID || res.Discount.ID != discount.ID || !res.IsActive {
		t.Errorf("unexpected response data: %+v", res)
	}

	// Skenario 2: Duplikasi Product Discount (409 Conflict)
	_, errConflict := productDiscountUsecase.Create(ctx, req)
	if errConflict == nil {
		t.Fatalf("expected conflict error on duplicate create, got nil")
	}

	// Skenario 3: Get Product Discounts With Pagination
	listReq := model.GetProductDiscountsRequest{
		DiscountID: discount.ID,
	}
	discountsList, pagination, err := productDiscountUsecase.GetProductDiscounts(ctx, listReq)
	if err != nil {
		t.Fatalf("expected get product discounts to succeed, got: %v", err)
	}
	if len(discountsList) != 1 || pagination.TotalData != 1 {
		t.Errorf("expected 1 product discount, got %d (total: %d)", len(discountsList), pagination.TotalData)
	}

	// Skenario 4: Get Product Discount By ID Berhasil
	detail, err := productDiscountUsecase.GetProductDiscountByID(ctx, res.ID)
	if err != nil {
		t.Fatalf("expected get by ID to succeed, got: %v", err)
	}
	if detail.ID != res.ID || detail.Product.ID != product.ID {
		t.Errorf("unexpected detail response: %+v", detail)
	}

	// Skenario 5: Get Product Discount By ID Not Found
	_, errNotFound := productDiscountUsecase.GetProductDiscountByID(ctx, uuid.New().String())
	if errNotFound == nil {
		t.Fatalf("expected not found error for non-existent ID, got nil")
	}
}

// 2. INTEGRATION TEST: HTTP Controller & Routing Layer
func TestProductDiscount_API(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Rollback()

	jwtSecret := "test-secret-key-1234567890123456"

	// Seed Role & Permission
	role := entity.Role{Name: "Owner_" + uuid.New().String()[:8]}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	var permission entity.Permission
	if err := db.Where("name = ?", "discounts:create").First(&permission).Error; err != nil {
		permission = entity.Permission{
			Name:        "discounts:create",
			Description: "Create product discount",
			Type:        "action",
		}
		if err := db.Create(&permission).Error; err != nil {
			t.Fatalf("failed to seed permission: %v", err)
		}
	}

	var readPermission entity.Permission
	if err := db.Where("name = ?", "discounts:read").First(&readPermission).Error; err != nil {
		readPermission = entity.Permission{
			Name:        "discounts:read",
			Description: "Read product discount",
			Type:        "action",
		}
		if err := db.Create(&readPermission).Error; err != nil {
			t.Fatalf("failed to seed read permission: %v", err)
		}
	}

	if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?), (?, ?) ON CONFLICT DO NOTHING", role.ID, permission.ID, role.ID, readPermission.ID).Error; err != nil {
		t.Fatalf("failed to link role permission: %v", err)
	}

	// Seed User
	userUUID := uuid.New()
	user := entity.User{
		ID:           userUUID.String(),
		RoleID:       role.ID,
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		Name:         "Test User",
		PasswordHash: "hashedpassword",
		IsActive:     true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	// Generate Auth Token
	token, err := utils.GenerateAccessToken(userUUID, role.Name, jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Seed Product & Discount
	category := entity.Category{Name: "Cat_" + uuid.New().String()[:8]}
	_ = db.Create(&category)
	unit := entity.Unit{Name: "Unit_" + uuid.New().String()[:8]}
	_ = db.Create(&unit)
	product := entity.Product{CategoryID: category.ID, Name: "Product_" + uuid.New().String()[:8], BaseUnitID: unit.ID, IsActive: true}
	_ = db.Create(&product)
	discount := entity.Discount{Name: "Promo_" + uuid.New().String()[:8], Type: "fixed", Value: decimal.NewFromFloat(5000), IsActive: true}
	_ = db.Create(&discount)

	// Wiring Setup
	blacklistRepo := repository.NewBlacklistRepository(db)
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	permissionUsecase := usecase.NewPermissionUsecase(db, permissionRepo, userRepo)

	productDiscountRepo := repository.NewProductDiscountRepository(db)
	productDiscountUsecase := usecase.NewProductDiscountUsecase(productDiscountRepo)
	productDiscountController := controller.NewProductDiscountController(productDiscountUsecase)

	app := gin.New()
	router.Setup(
		app,
		jwtSecret,
		blacklistRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		productDiscountController,
		nil,
		permissionUsecase,
		nil, nil, nil, nil,
	)

	// Test Case 1: POST /api/product-discounts (201 Created)
	bodyReq, _ := json.Marshal(model.CreateProductDiscountRequest{
		ProductID:  product.ID,
		DiscountID: discount.ID,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/product-discounts", bytes.NewBuffer(bodyReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected HTTP 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var createdRes struct {
		Data model.ProductDiscountResponse `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &createdRes)

	// Test Case 2: Duplicate Product Discount (409 Conflict)
	wDuplicate := httptest.NewRecorder()
	reqDuplicate, _ := http.NewRequest(http.MethodPost, "/api/product-discounts", bytes.NewBuffer(bodyReq))
	reqDuplicate.Header.Set("Content-Type", "application/json")
	reqDuplicate.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wDuplicate, reqDuplicate)

	if wDuplicate.Code != http.StatusConflict {
		t.Fatalf("expected HTTP 409, got %d, body: %s", wDuplicate.Code, wDuplicate.Body.String())
	}

	// Test Case 3: Invalid Payload (422 Unprocessable Entity)
	invalidBody := []byte(`{"productId":"invalid-uuid","discountId":""}`)
	wInvalid := httptest.NewRecorder()
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/api/product-discounts", bytes.NewBuffer(invalidBody))
	reqInvalid.Header.Set("Content-Type", "application/json")
	reqInvalid.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected HTTP 422, got %d, body: %s", wInvalid.Code, wInvalid.Body.String())
	}

	// Test Case 4: Missing Auth Token (401 Unauthorized)
	wUnauthorized := httptest.NewRecorder()
	reqUnauthorized, _ := http.NewRequest(http.MethodPost, "/api/product-discounts", bytes.NewBuffer(bodyReq))
	reqUnauthorized.Header.Set("Content-Type", "application/json")
	app.ServeHTTP(wUnauthorized, reqUnauthorized)

	if wUnauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected HTTP 401, got %d, body: %s", wUnauthorized.Code, wUnauthorized.Body.String())
	}

	// Test Case 5: GET /api/product-discounts (200 OK)
	wGetList := httptest.NewRecorder()
	reqGetList, _ := http.NewRequest(http.MethodGet, "/api/product-discounts?discountId="+discount.ID, nil)
	reqGetList.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wGetList, reqGetList)

	if wGetList.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d, body: %s", wGetList.Code, wGetList.Body.String())
	}

	// Test Case 6: GET /api/product-discounts/:id (200 OK)
	wGetDetail := httptest.NewRecorder()
	reqGetDetail, _ := http.NewRequest(http.MethodGet, "/api/product-discounts/"+createdRes.Data.ID, nil)
	reqGetDetail.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wGetDetail, reqGetDetail)

	if wGetDetail.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d, body: %s", wGetDetail.Code, wGetDetail.Body.String())
	}

	// Test Case 7: GET /api/product-discounts/:id Not Found (404 Not Found)
	wGetNotFound := httptest.NewRecorder()
	reqGetNotFound, _ := http.NewRequest(http.MethodGet, "/api/product-discounts/"+uuid.New().String(), nil)
	reqGetNotFound.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(wGetNotFound, reqGetNotFound)

	if wGetNotFound.Code != http.StatusNotFound {
		t.Fatalf("expected HTTP 404, got %d, body: %s", wGetNotFound.Code, wGetNotFound.Body.String())
	}
}
