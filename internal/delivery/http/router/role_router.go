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
	role := router.Group("/roles", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		role.POST("", controller.CreateRole)
		role.GET("", controller.GetRolesWithPagination)
		role.GET("/:id", controller.GetRoleByID)
		role.PUT("/:id", controller.UpdateRole)
		role.DELETE("/:id", controller.DeleteRole)
	}
}
