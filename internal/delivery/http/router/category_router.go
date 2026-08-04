package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerCategoryRoutes(
	router *gin.RouterGroup,
	controller *controller.CategoryController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	category := router.Group("/categories", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		category.POST("", controller.CreateCategory)
		category.GET("", controller.GetCategoriesWithPagination)
		category.GET("/:id", controller.GetCategoryById)
		category.PUT("/:id", controller.UpdateCategory)
		category.DELETE("/:id", controller.DeleteCategory)
	}
}
