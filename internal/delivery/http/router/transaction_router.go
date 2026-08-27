package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerTransactionRoutes(
	router *gin.RouterGroup,
	controller *controller.TransactionController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	transaction := router.Group("/transactions", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		transaction.POST("", middleware.RequirePermission(permissionUsecase, "pos:create"), controller.CreateTransaction)
		transaction.GET("", middleware.RequirePermission(permissionUsecase, "transactions:read"), controller.ListTransactions)
		transaction.GET("/:id", middleware.RequirePermission(permissionUsecase, "transactions:read"), controller.GetTransactionByID)
		transaction.POST("/:id/void", middleware.RequirePermission(permissionUsecase, "transactions:void"), controller.VoidTransaction)
	}
}
