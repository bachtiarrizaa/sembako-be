package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type UnitController struct {
	usecase   *usecase.Unitusecase
	validator *validator.Validate
}

func NewUnitController(usecase *usecase.Unitusecase) *UnitController {
	return &UnitController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *UnitController) CreateUnit(ctx *gin.Context) {
	var req model.CreateUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreateUnit(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "unit created successfully", res)
}

func (c *UnitController) GetUnitsWithPagination(ctx *gin.Context) {
	pagReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetUnitsWithPagination(ctx.Request.Context(), pagReq)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "units fetched successfully", res, pagination)
}

func (c *UnitController) GetUnitByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := c.usecase.GetUnitByID(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "unit fetched successfully", res)
}

func (c *UnitController) UpdateUnit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateUnit(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "unit updated successfully", res)
}

func (c *UnitController) DeleteUnit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.usecase.DeleteUnit(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "unit deleted successfully", nil)
}
