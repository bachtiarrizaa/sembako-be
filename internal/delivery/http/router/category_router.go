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
	roles := router.Group("/categories", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		roles.POST("", controller.CreateCategory)
		roles.GET("", controller.GetCategoriesWithPagination)
		roles.GET("/:id", controller.GetCategoryById)
		roles.PUT("/:id", controller.UpdateCategory)
		roles.DELETE("/:id", controller.DeleteCategory)
	}
}
