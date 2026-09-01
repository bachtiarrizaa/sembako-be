package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

type StoreConfigurationController struct {
	storeConfigurationUsecase *usecase.StoreConfigurationUsecase
}

func NewStoreConfigurationController(
	storeConfigurationUsecase *usecase.StoreConfigurationUsecase,
) *StoreConfigurationController {
	return &StoreConfigurationController{
		storeConfigurationUsecase: storeConfigurationUsecase,
	}
}

func (c *StoreConfigurationController) Get(ctx *gin.Context) {
	res, err := c.storeConfigurationUsecase.Get(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "store configuration fetched successfully", res)
}
