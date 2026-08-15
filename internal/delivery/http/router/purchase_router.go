package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerPurchaseRoutes(
	router *gin.RouterGroup,
	controller *controller.PurchaseController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	purchase := router.Group("/purchases", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		purchase.POST("", middleware.RequirePermission(permissionUsecase, "purchases:create"), controller.CreatePurchase)
		purchase.GET("", middleware.RequirePermission(permissionUsecase, "purchases:read"), controller.GetPurchaseBatches)
		purchase.GET("/:id", middleware.RequirePermission(permissionUsecase, "purchases:read"), controller.GetPurchaseBatchByID)
		purchase.PUT("/:id", middleware.RequirePermission(permissionUsecase, "purchases:update"), controller.UpdatePurchase)
		purchase.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "purchases:delete"), controller.DeletePurchase)
	}
}
