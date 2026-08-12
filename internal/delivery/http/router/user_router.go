package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	router *gin.RouterGroup,
	controller *controller.UserController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	user := router.Group("/users", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		user.POST("", middleware.RequirePermission(permissionUsecase, "users:create"), controller.CreateUser)
		user.GET("", middleware.RequirePermission(permissionUsecase, "users:read"), controller.GetUsersWithPagination)
		user.GET("/me", controller.GetProfileMe)
		user.PATCH("/me", controller.UpdateProfile)
		user.GET("/:id", middleware.RequirePermission(permissionUsecase, "users:read"), controller.GetUserByID)
		user.PUT("/:id", middleware.RequirePermission(permissionUsecase, "users:update"), controller.UpdateUser)
		user.PATCH("/:id/status", middleware.RequirePermission(permissionUsecase, "users:update"), controller.UpdateStatus)
		user.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "users:delete"), controller.DeleteUser)
	}
}
