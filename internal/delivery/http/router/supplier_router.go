package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerSupplierRoutes(
	router *gin.RouterGroup,
	controller *controller.SupplierController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	supplier := router.Group("/suppliers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		supplier.POST("", middleware.RequirePermission(permissionUsecase, "suppliers:create"), controller.CreateSupplier)
		supplier.GET("", middleware.RequirePermission(permissionUsecase, "suppliers:read"), controller.GetSuppliersWithPagination)
		supplier.GET("/:id", middleware.RequirePermission(permissionUsecase, "suppliers:read"), controller.GetSupplierById)
		supplier.PUT("/:id", middleware.RequirePermission(permissionUsecase, "suppliers:update"), controller.UpdateSupplier)
		supplier.PATCH("/:id/status", middleware.RequirePermission(permissionUsecase, "suppliers:update"), controller.UpdateStatus)
		supplier.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "suppliers:delete"), controller.DeleteSupplier)
	}
}
