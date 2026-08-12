package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerRoleRoutes(
	router *gin.RouterGroup,
	controller *controller.RoleController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	role := router.Group("/roles", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		role.POST("", middleware.RequirePermission(permissionUsecase, "roles:create"), controller.CreateRole)
		role.GET("", middleware.RequirePermission(permissionUsecase, "roles:read"), controller.GetRolesWithPagination)
		role.GET("/:id", middleware.RequirePermission(permissionUsecase, "roles:read"), controller.GetRoleByID)
		role.PUT("/:id", middleware.RequirePermission(permissionUsecase, "roles:update"), controller.UpdateRole)
		role.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "roles:delete"), controller.DeleteRole)
	}
}
