package router

import (
	"github.com/gin-gonic/gin"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

func registerPermissionRoutes(
	router *gin.RouterGroup,
	controller *controller.PermissionController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	menu := router.Group("/users/me/menu", middleware.AuthMiddleware(jwtSecret, blacklistRepo))
	{
		menu.GET("", controller.GetUserMenu)
	}
}
