package bootstrap

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/router"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/brevo"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func InitializeApp(cfg *config.Config) (*gin.Engine, error) {
	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET must not be empty")
	}
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET must not be empty")
	}

	db, err := config.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	refreshTTL := time.Duration(cfg.JWTRefreshExpireDays) * 24 * time.Hour
	isProduction := cfg.AppEnv == "production"

	roleRepo := repository.NewRoleRepository(db)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	blacklistRepo := repository.NewBlacklistRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	brevoService := brevo.NewBrevoService(cfg.BrevoApiKey, cfg.BrevoSenderEmail, cfg.BrevoSenderName)

	authUsecase := usecase.NewAuthUsecase(
		userRepo,
		refreshTokenRepo,
		blacklistRepo,
		permissionRepo,
		cfg.JWTAccessSecret,
		time.Duration(cfg.JWTAccessExpireMinutes)*time.Minute,
		refreshTTL,
	)

	passwordResetUsecase := usecase.NewPasswordResetUsecase(
		userRepo,
		refreshTokenRepo,
		passwordResetRepo,
		brevoService,
		cfg.FrontendResetUrl,
		cfg.ResetTokenExpireMinutes,
	)

	authController := controller.NewAuthController(authUsecase, passwordResetUsecase, isProduction, refreshTTL)

	userUsecase := usecase.NewUserUsecase(userRepo, roleRepo)
	userController := controller.NewUserController(userUsecase, cfg.UploadDir)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUsecase)

	supplierRepo := repository.NewSupplierRepository(db)
	supplierUsecase := usecase.NewSupplierUsecase(supplierRepo)
	supplierController := controller.NewSupplierController(supplierUsecase)

	unitRepo := repository.NewUnitRepository(db)
	unitUsecase := usecase.NewUnitUsecase(unitRepo)
	unitController := controller.NewUnitController(unitUsecase)

	customerRepo := repository.NewCustomerRepository(db)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerController := controller.NewCustomerController(customerUsecase)

	productRepo := repository.NewProductRepository(db)
	productUnitRepo := repository.NewProductUnitRepository(db)
	productUsecase := usecase.NewProductUsecase(db, productRepo, productUnitRepo)
	productController := controller.NewProductController(productUsecase, cfg.UploadDir)

	discountRepo := repository.NewDiscountRepository(db)
	discountUsecase := usecase.NewDiscountUsecase(discountRepo)
	discountController := controller.NewDiscountController(discountUsecase)

	permissionUsecase := usecase.NewPermissionUsecase(db, permissionRepo, userRepo)
	permissionController := controller.NewPermissionController(permissionUsecase)

	stockRepo := repository.NewStockRepository(db)
	stockMutationRepo := repository.NewStockMutationRepository(db)
	stockCountRepo := repository.NewStockCountRepository(db)

	purchaseBatchRepo := repository.NewPurchaseBatchRepository(db)
	purchaseUsecase := usecase.NewPurchaseUsecase(db, purchaseBatchRepo, productRepo, supplierRepo, stockRepo, stockMutationRepo)
	purchaseController := controller.NewPurchaseController(purchaseUsecase)

	stockUsecase := usecase.NewStockUsecase(db, stockRepo, stockMutationRepo, stockCountRepo, productRepo, permissionUsecase)
	stockController := controller.NewStockController(stockUsecase)

	app := gin.Default()
	app.Static("/uploads", cfg.UploadDir)
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Setup(
		app,
		cfg.JWTAccessSecret,
		blacklistRepo,
		roleController,
		authController,
		userController,
		categoryController,
		supplierController,
		unitController,
		customerController,
		productController,
		discountController,
		permissionController,
		permissionUsecase,
		purchaseController,
		stockController,
	)

	return app, nil
}
