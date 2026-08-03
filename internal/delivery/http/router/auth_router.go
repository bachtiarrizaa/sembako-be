package router

import (
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(
	router *gin.RouterGroup,
	controller *controller.AuthController,
	jwtSecret string,
	blacklistRepo repository.BlacklistRepository,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/refresh", controller.Refresh)
		auth.POST("/logout", middleware.AuthMiddleware(jwtSecret, blacklistRepo), controller.Logout)
		auth.POST("/forgot-password", controller.ForgotPassword)
		auth.POST("/reset-password", controller.ResetPassword)
	}
}
