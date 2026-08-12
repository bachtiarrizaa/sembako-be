package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerUnitRoutes(
	router *gin.RouterGroup,
	controller *controller.UnitController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	unit := router.Group("/units", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		unit.POST("", middleware.RequirePermission(permissionUsecase, "units:create"), controller.CreateUnit)
		unit.GET("", middleware.RequirePermission(permissionUsecase, "units:read"), controller.GetUnitsWithPagination)
		unit.GET("/:id", middleware.RequirePermission(permissionUsecase, "units:read"), controller.GetUnitByID)
		unit.PUT("/:id", middleware.RequirePermission(permissionUsecase, "units:update"), controller.UpdateUnit)
		unit.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "units:delete"), controller.DeleteUnit)
	}
}
