package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerUnitRoutes(
	router *gin.RouterGroup,
	controller *controller.UnitController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	unit := router.Group("/units", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		unit.POST("", controller.CreateUnit)
		unit.GET("", controller.GetUnitsWithPagination)
		unit.GET("/:id", controller.GetUnitByID)
		unit.PUT("/:id", controller.UpdateUnit)
		unit.DELETE("/:id", controller.DeleteUnit)
	}
}
