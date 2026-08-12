package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerDiscountRoutes(
	router *gin.RouterGroup,
	controller *controller.DiscountController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	discount := router.Group("/discounts", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		discount.POST("", middleware.RequirePermission(permissionUsecase, "discounts:create"), controller.CreateDiscount)
		discount.GET("", middleware.RequirePermission(permissionUsecase, "discounts:read"), controller.GetDiscountWithPagination)
		discount.GET("/:id", middleware.RequirePermission(permissionUsecase, "discounts:read"), controller.GetDiscountById)
		discount.PUT("/:id", middleware.RequirePermission(permissionUsecase, "discounts:update"), controller.UpdateDiscount)
		discount.PATCH("/:id/status", middleware.RequirePermission(permissionUsecase, "discounts:update"), controller.UpdateStatus)
		discount.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "discounts:delete"), controller.DeleteDiscount)
	}
}
