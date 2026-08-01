package router

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
)

func Setup(router *gin.Engine, roleController *controller.RoleController) {
	api := router.Group("/api")

	roles := api.Group("/roles")
	{
		roles.POST("", roleController.Create)
		roles.GET("", roleController.GetAll)
		roles.GET("/:id", roleController.GetByID)
		roles.PUT("/:id", roleController.Update)
		roles.DELETE("/:id", roleController.Delete)
	}
}
