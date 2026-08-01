package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/controller"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/router"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func InitializeApp(cfg *config.Config) (*gin.Engine, error) {
	db, err := config.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	roleRepo := repository.NewRoleRepository(db)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)

	app := gin.Default()
	router.Setup(app, roleController)

	return app, nil
}
