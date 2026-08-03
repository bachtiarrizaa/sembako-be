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
		roles.POST("", controller.CreateRole)
		roles.GET("", controller.GetRolesWithPagination)
		roles.GET("/:id", controller.GetRoleByID)
		roles.PUT("/:id", controller.UpdateRole)
		roles.DELETE("/:id", controller.DeleteRole)
	}
}
