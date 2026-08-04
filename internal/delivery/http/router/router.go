package router

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
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
) {
	api := router.Group("/api")

	registerAuthRoutes(api, authController, jwtSecret, blacklistRepo)

	registerRoleRoutes(api, roleController, jwtSecret, blacklistRepo)

	registerUserRoutes(api, userController, jwtSecret, blacklistRepo)

	registerCategoryRoutes(api, categoryController, jwtSecret, blacklistRepo)

	registerSupplierRoutes(api, supplierController, jwtSecret, blacklistRepo)

	registerUnitRoutes(api, unitController, jwtSecret, blacklistRepo)

	registerCustomerRoutes(api, customerController, jwtSecret, blacklistRepo)

	registerProductRoutes(api, productController, jwtSecret, blacklistRepo)

	registerDiscountRoutes(api, discountController, jwtSecret, blacklistRepo)
}
