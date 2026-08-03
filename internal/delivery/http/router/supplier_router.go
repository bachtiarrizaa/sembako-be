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
	roles := router.Group("/suppliers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		roles.POST("", controller.CreateSupplier)
		roles.GET("", controller.GetSuppliersWithPagination)
		roles.GET("/:id", controller.GetSupplierById)
		roles.PUT("/:id", controller.UpdateSupplier)
		roles.DELETE("/:id", controller.DeleteSupplier)
	}
}
