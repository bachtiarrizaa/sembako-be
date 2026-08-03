package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerProductRoutes(
	router *gin.RouterGroup,
	controller *controller.ProductController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	product := router.Group("/products", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		product.POST("", controller.CreateProduct)
		product.GET("", controller.GetProducts)
		product.GET("/:id", controller.GetProductByID)
		product.PUT("/:id", controller.UpdateProduct)
		product.PATCH("/:id/status", controller.UpdateProductStatus)
		
		// Unit management
		product.POST("/:id/units", controller.AddProductUnit)
		product.PUT("/:id/units/:unitId", controller.UpdateProductUnit)
		product.PATCH("/:id/units/:unitId/status", controller.UpdateProductUnitStatus)
		product.DELETE("/:id/units/:unitId", controller.DeleteProductUnit)

		product.DELETE("/:id", controller.DeleteProduct)
	}
}
