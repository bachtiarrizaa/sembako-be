package bootstrap

import (
	"fmt"
	"time"

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
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	brevoService := brevo.NewBrevoService(cfg.BrevoApiKey, cfg.BrevoSenderEmail, cfg.BrevoSenderName)

	authUsecase := usecase.NewAuthUsecase(
		userRepo,
		refreshTokenRepo,
		blacklistRepo,
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
	userController := controller.NewUserController(userUsecase)

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

	app := gin.Default()
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
	)

	return app, nil
}
