package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StoreConfigurationController struct {
	storeConfigurationUsecase *usecase.StoreConfigurationUsecase
	validator                 *validator.Validate
}

func NewStoreConfigurationController(
	storeConfigurationUsecase *usecase.StoreConfigurationUsecase,
) *StoreConfigurationController {
	return &StoreConfigurationController{
		storeConfigurationUsecase: storeConfigurationUsecase,
		validator:                 validator.New(),
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

func (c *StoreConfigurationController) GetPublicStoreInfo(ctx *gin.Context) {
	res, err := c.storeConfigurationUsecase.GetPublicStoreInfo(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "store info fetched successfully", res)
}

func (c *StoreConfigurationController) Update(ctx *gin.Context) {
	var req model.UpdateStoreConfigurationRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.storeConfigurationUsecase.Update(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "store configuration updated successfully", res)
}
