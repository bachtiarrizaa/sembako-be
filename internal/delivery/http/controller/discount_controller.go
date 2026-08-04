package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DiscountController struct {
	usecase   *usecase.DiscountUsecase
	validator *validator.Validate
}

func NewDiscountController(usecase *usecase.DiscountUsecase) *DiscountController {
	return &DiscountController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *DiscountController) CreateDiscount(ctx *gin.Context) {
	var req model.CreateDiscountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreateDiscount(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "discount created successfully", res)
}

func (c *DiscountController) GetDiscountWithPagination(ctx *gin.Context) {
	pageReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetDiscountWithPagination(ctx.Request.Context(), pageReq)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "discounts fetched successfully", res, pagination)
}

func (c *DiscountController) GetDiscountById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := c.usecase.GetDiscountById(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "discount fetched successfully", res)
}

func (c *DiscountController) UpdateDiscount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateDiscountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateDiscount(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "discount updated successfully", res)
}

func (c *DiscountController) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateStatusDiscountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateStatus(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "discount status updated successfully", res)
}

func (c *DiscountController) DeleteDiscount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.usecase.DeleteDiscount(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "discount deleted successfully", nil)
}
