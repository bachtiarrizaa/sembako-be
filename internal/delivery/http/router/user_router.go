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
	user := router.Group("/users", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		user.POST("", controller.CreateUser)
		user.GET("", controller.GetUsersWithPagination)
		user.GET("/me", controller.GetProfileMe)
		user.PATCH("/me", controller.UpdateProfile)
		user.GET("/:id", controller.GetUserByID)
		user.PUT("/:id", controller.UpdateUser)
		user.PATCH("/:id/status", controller.UpdateStatus)
		user.DELETE("/:id", controller.DeleteUser)
	}
}
