package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerLoyaltySettingRoutes(
	router *gin.RouterGroup,
	controller *controller.LoyaltySettingController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	loyalty := router.Group("/loyalty-settings", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		loyalty.GET("", middleware.RequirePermission(permissionUsecase, "loyalty:read"), controller.Get)
		loyalty.PUT("", middleware.RequirePermission(permissionUsecase, "loyalty:write"), controller.Update)
	}
}
