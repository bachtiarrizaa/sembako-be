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
		product.POST("", controller.Create)
		product.GET("", controller.GetProducts)
		product.GET("/:id", controller.GetByID)
	}
}
