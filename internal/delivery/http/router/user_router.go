package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	router *gin.RouterGroup,
	controller *controller.UserController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	roles := router.Group("/users", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		roles.POST("", controller.CreateUser)
		roles.GET("", controller.GetUsersWithPagination)
		roles.GET("/:id", controller.GetUserByID)
		roles.PUT("/:id", controller.UpdateUser)
		roles.DELETE("/:id", controller.DeleteUser)
	}
}
