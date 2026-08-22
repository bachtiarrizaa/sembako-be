package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerShiftRoutes(
	router *gin.RouterGroup,
	controller *controller.ShiftController,
	permission *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	shift := router.Group("/shifts", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		shift.POST("/open", middleware.RequirePermission(permission, "shifts:create"), controller.OpenShift)

		// endpoint berikutnya, controller-nya belum dibuat:
		// shift.GET("/active", middleware.RequirePermission(permission, "shifts:read"), controller.GetActiveShift)
		// shift.POST("/:id/close", middleware.RequirePermission(permission, "shifts:close"), controller.CloseShift)
		// shift.POST("/:id/force-close", middleware.RequirePermission(permission, "shifts:force-close"), controller.ForceCloseShift)
		// shift.GET("", middleware.RequirePermission(permission, "shifts:read"), controller.ListShifts)
		// shift.GET("/:id", middleware.RequirePermission(permission, "shifts:read"), controller.GetShiftDetail)
	}
}
