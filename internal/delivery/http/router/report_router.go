package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerReportRoutes(
	router *gin.RouterGroup,
	controller *controller.ReportController,
	permission *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	if controller == nil {
		return
	}
	reports := router.Group("/reports", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		reports.GET("", middleware.RequirePermission(permission, "reports:read"), controller.GetReport)
	}
}
