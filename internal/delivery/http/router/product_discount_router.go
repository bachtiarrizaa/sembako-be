package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerProductDiscountRoutes(
	router *gin.RouterGroup,
	controller *controller.ProductDiscountController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	productDiscount := router.Group("/product-discounts", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		productDiscount.POST("", middleware.RequirePermission(permissionUsecase, "discounts:create"), controller.Create)
		productDiscount.GET("", middleware.RequirePermission(permissionUsecase, "discounts:read"), controller.GetProductDiscounts)
		productDiscount.GET("/:id", middleware.RequirePermission(permissionUsecase, "discounts:read"), controller.GetProductDiscountByID)
	}
}
