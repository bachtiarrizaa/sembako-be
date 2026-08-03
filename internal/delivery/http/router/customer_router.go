package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerCustomerRoutes(
	router *gin.RouterGroup,
	controller *controller.CustomerController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	unit := router.Group("/customers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		unit.POST("", controller.CreateCustomer)
		unit.GET("", controller.GetCustomersWithPagination)
		unit.GET("/:id", controller.GetCustomerById)
		unit.PUT("/:id", controller.UpdateCustomer)
		unit.DELETE("/:id", controller.DeleteCustomer)
	}
}
