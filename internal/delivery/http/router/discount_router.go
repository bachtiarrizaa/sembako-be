package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerDiscountRoutes(
	router *gin.RouterGroup,
	controller *controller.DiscountController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	discount := router.Group("/discounts", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		discount.POST("", controller.CreateDiscount)
		discount.GET("", controller.GetDiscountWithPagination)
		discount.GET("/:id", controller.GetDiscountById)
		discount.PUT("/:id", controller.UpdateDiscount)
		discount.PATCH("/:id/status", controller.UpdateStatus)
		discount.DELETE("/:id", controller.DeleteDiscount)
	}
}
