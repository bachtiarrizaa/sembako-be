package router

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func Setup(
	router *gin.Engine,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
	roleController *controller.RoleController,
	authController *controller.AuthController,
	userController *controller.UserController,
	categoryController *controller.CategoryController,
	supplierController *controller.SupplierController,
	unitController *controller.UnitController,
	customerController *controller.CustomerController,
	productController *controller.ProductController,
	discountController *controller.DiscountController,
	permissionController *controller.PermissionController,
	permissionUsecase *usecase.PermissionUsecase,
	purchaseController *controller.PurchaseController,
	stockController *controller.StockController,
) {
	api := router.Group("/api")

	registerAuthRoutes(api, authController, jwtSecret, blacklistRepo)

	registerRoleRoutes(api, roleController, permissionUsecase, jwtSecret, blacklistRepo)

	registerUserRoutes(api, userController, permissionUsecase, jwtSecret, blacklistRepo)

	registerCategoryRoutes(api, categoryController, permissionUsecase, jwtSecret, blacklistRepo)

	registerSupplierRoutes(api, supplierController, permissionUsecase, jwtSecret, blacklistRepo)

	registerUnitRoutes(api, unitController, permissionUsecase, jwtSecret, blacklistRepo)

	registerCustomerRoutes(api, customerController, permissionUsecase, jwtSecret, blacklistRepo)

	registerProductRoutes(api, productController, permissionUsecase, jwtSecret, blacklistRepo)

	registerDiscountRoutes(api, discountController, permissionUsecase, jwtSecret, blacklistRepo)

	registerPermissionRoutes(api, permissionController, jwtSecret, blacklistRepo)

	registerPurchaseRoutes(api, purchaseController, permissionUsecase, jwtSecret, blacklistRepo)

	registerStockRoutes(api, stockController, permissionUsecase, jwtSecret, blacklistRepo)
}
