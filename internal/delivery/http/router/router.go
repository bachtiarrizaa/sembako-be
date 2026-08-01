package router

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

func Setup(
	router *gin.Engine,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
	roleController *controller.RoleController,
	authController *controller.AuthController,
	userController *controller.UserController,
) {
	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", middleware.AuthMiddleware(jwtSecret, blacklistRepo), authController.Logout)
	}

	roles := api.Group("/roles")
	{
		roles.POST("", roleController.Create)
		roles.GET("", middleware.AuthMiddleware(jwtSecret, blacklistRepo), roleController.GetAll)
		roles.GET("/:id", roleController.GetByID)
		roles.PUT("/:id", roleController.Update)
		roles.DELETE("/:id", roleController.Delete)
	}

	users := api.Group("/users", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		users.GET("/me", userController.GetMe)
	}
}
