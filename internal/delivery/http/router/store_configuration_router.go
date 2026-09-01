package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerStoreConfigurationRoutes(
	router *gin.RouterGroup,
	controller *controller.StoreConfigurationController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	storeConfiguration := router.Group("/settings", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		storeConfiguration.GET("", middleware.RequirePermission(permissionUsecase, "settings:read"), controller.Get)
	}
}
