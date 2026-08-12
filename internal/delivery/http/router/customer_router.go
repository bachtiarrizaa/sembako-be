package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerCustomerRoutes(
	router *gin.RouterGroup,
	controller *controller.CustomerController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	customer := router.Group("/customers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		customer.POST("", middleware.RequirePermission(permissionUsecase, "customers:create"), controller.CreateCustomer)
		customer.GET("", middleware.RequirePermission(permissionUsecase, "customers:read"), controller.GetCustomersWithPagination)
		customer.GET("/:id", middleware.RequirePermission(permissionUsecase, "customers:read"), controller.GetCustomerById)
		customer.PUT("/:id", middleware.RequirePermission(permissionUsecase, "customers:update"), controller.UpdateCustomer)
		customer.PATCH("/:id/status", middleware.RequirePermission(permissionUsecase, "customers:update"), controller.UpdateStatus)
		customer.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "customers:delete"), controller.DeleteCustomer)
	}
}
