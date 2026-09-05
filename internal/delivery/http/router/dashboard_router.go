package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerDashboardRoutes(
	router *gin.RouterGroup,
	controller *controller.DashboardController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	dashboard := router.Group("/dashboard", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		dashboard.GET("", middleware.RequirePermission(permissionUsecase, "dashboard:read"), controller.GetDashboard)
	}
}
