package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerRoleRoutes(
	router *gin.RouterGroup,
	controller *controller.RoleController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	roles := router.Group("/roles", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		roles.POST("", controller.Create)
		roles.GET("", controller.GetRolesWithPagination)
		roles.GET("/:id", controller.GetByID)
		roles.PUT("/:id", controller.Update)
		roles.DELETE("/:id", controller.Delete)
	}
}
