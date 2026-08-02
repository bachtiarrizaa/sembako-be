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
	categoryController *controller.CategoryController,
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
		users.POST("", userController.Create)
		users.GET("", userController.GetAll)
		users.GET("/:id", userController.GetByID)
		users.PUT("/:id", userController.Update)
		users.DELETE("/:id", userController.Delete)
	}

	categories := api.Group("/categories", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		categories.POST("", categoryController.Create)
		categories.GET("", categoryController.GetCategories)
		categories.GET("/:id", categoryController.GetCategoryById)
		categories.PUT("/:id", categoryController.UpdateCategory)
		categories.DELETE("/:id", categoryController.DeleteCategory)
	}
}
