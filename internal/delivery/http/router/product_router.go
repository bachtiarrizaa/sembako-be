package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerProductRoutes(
	router *gin.RouterGroup,
	controller *controller.ProductController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	product := router.Group("/products", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		product.POST("", middleware.RequirePermission(permissionUsecase, "products:create"), controller.CreateProduct)
		product.GET("", middleware.RequirePermission(permissionUsecase, "products:read"), controller.GetProducts)
		product.GET("/:id", middleware.RequirePermission(permissionUsecase, "products:read"), controller.GetProductByID)
		product.PUT("/:id", middleware.RequirePermission(permissionUsecase, "products:update"), controller.UpdateProduct)
		product.PATCH("/:id/status", middleware.RequirePermission(permissionUsecase, "products:update"), controller.UpdateProductStatus)
		product.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "products:delete"), controller.DeleteProduct)
	}
}
