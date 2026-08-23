package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (c *ShiftController) CloseShift(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	shiftID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid shift id")
		return
	}

	var req model.CloseShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.shiftUsecase.CloseShift(ctx.Request.Context(), shiftID, userID, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "shift closed successfully", res)
}

func (c *ShiftController) ForceCloseShift(ctx *gin.Context) {
	adminID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	shiftID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid shift id")
		return
	}

	var req model.ForceCloseShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.shiftUsecase.ForceCloseShift(ctx.Request.Context(), shiftID, adminID, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "shift force closed successfully", res)
}

func (c *ShiftController) ListShifts(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	role, exists := middleware.GetRole(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.ListShiftsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	res, pagination, err := c.shiftUsecase.ListShifts(ctx.Request.Context(), req, userID, role)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "shifts fetched successfully", res, pagination)
}

func (c *ShiftController) GetShiftDetail(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	role, exists := middleware.GetRole(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	shiftID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid shift id")
		return
	}

	res, err := c.shiftUsecase.GetShiftDetail(ctx.Request.Context(), shiftID, userID, role)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "shift detail fetched successfully", res)
}
