package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerSupplierRoutes(
	router *gin.RouterGroup,
	controller *controller.SupplierController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	supplier := router.Group("/suppliers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		supplier.POST("", controller.CreateSupplier)
		supplier.GET("", controller.GetSuppliersWithPagination)
		supplier.GET("/:id", controller.GetSupplierById)
		supplier.PUT("/:id", controller.UpdateSupplier)
		supplier.PATCH("/:id/status", controller.UpdateStatus)
		supplier.DELETE("/:id", controller.DeleteSupplier)
	}
}
