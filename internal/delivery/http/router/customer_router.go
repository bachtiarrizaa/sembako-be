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
	customer := router.Group("/customers", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		customer.POST("", controller.CreateCustomer)
		customer.GET("", controller.GetCustomersWithPagination)
		customer.GET("/:id", controller.GetCustomerById)
		customer.PUT("/:id", controller.UpdateCustomer)
		customer.DELETE("/:id", controller.DeleteCustomer)
	}
}
