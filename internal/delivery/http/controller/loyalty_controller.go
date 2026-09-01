package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoyaltySettingController struct {
	loyaltySettingUsecase *usecase.LoyaltySettingUsecase
	validator             *validator.Validate
}

func NewLoyaltySettingController(
	loyaltySettingUsecase *usecase.LoyaltySettingUsecase,
) *LoyaltySettingController {
	return &LoyaltySettingController{
		loyaltySettingUsecase: loyaltySettingUsecase,
		validator:             validator.New(),
	}
}

func (c *LoyaltySettingController) Get(ctx *gin.Context) {
	res, err := c.loyaltySettingUsecase.Get(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "loyalty setting fetched successfully", res)
}

func (c *LoyaltySettingController) Update(ctx *gin.Context) {
	var req model.UpdateLoyaltySettingRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.loyaltySettingUsecase.Update(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "loyalty setting updated successfully", res)
}
