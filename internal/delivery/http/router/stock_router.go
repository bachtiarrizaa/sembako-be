package router

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func registerStockRoutes(
	router *gin.RouterGroup,
	controller *controller.StockController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	stocks := router.Group("/stocks", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		stocks.POST("/opname", middleware.RequirePermission(permissionUsecase, "opname:create"), controller.SubmitStockCount)
		stocks.GET("/opname", middleware.RequirePermission(permissionUsecase, "opname:read"), controller.GetStockCounts)
		stocks.POST("/opname/:id/approve", middleware.RequirePermission(permissionUsecase, "opname:approve"), controller.ApproveStockCount)
		stocks.GET("/:productID", middleware.RequirePermission(permissionUsecase, "stocks:read"), controller.GetStockByProductID)
		stocks.GET("/:productID/mutations", middleware.RequirePermission(permissionUsecase, "stocks:read"), controller.GetStockMutations)
	}
}
