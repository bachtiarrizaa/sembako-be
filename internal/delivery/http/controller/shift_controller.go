package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShiftController struct {
	shiftUsecase *usecase.ShiftUsecase
	validator    *validator.Validate
}

func NewShiftController(shiftUsecase *usecase.ShiftUsecase) *ShiftController {
	return &ShiftController{
		shiftUsecase: shiftUsecase,
		validator:    validator.New(),
	}
}

func (c *ShiftController) OpenShift(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.OpenShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.shiftUsecase.OpenShift(ctx.Request.Context(), userID, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "shift opened successfully", res)
}

func (c *ShiftController) GetActiveShift(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := c.shiftUsecase.GetActiveShift(ctx.Request.Context(), userID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "active shift fetched successfully", res)
}
