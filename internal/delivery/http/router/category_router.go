package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

func registerCategoryRoutes(
	router *gin.RouterGroup,
	controller *controller.CategoryController,
	permissionUsecase *usecase.PermissionUsecase,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	category := router.Group("/categories", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		category.POST("", middleware.RequirePermission(permissionUsecase, "categories:create"), controller.CreateCategory)
		category.GET("", middleware.RequirePermission(permissionUsecase, "categories:read"), controller.GetCategoriesWithPagination)
		category.GET("/:id", middleware.RequirePermission(permissionUsecase, "categories:read"), controller.GetCategoryById)
		category.PUT("/:id", middleware.RequirePermission(permissionUsecase, "categories:update"), controller.UpdateCategory)
		category.DELETE("/:id", middleware.RequirePermission(permissionUsecase, "categories:delete"), controller.DeleteCategory)
	}
}
