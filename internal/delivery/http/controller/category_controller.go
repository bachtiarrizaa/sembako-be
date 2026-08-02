package controller

import (
	"errors"
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryController struct {
	usecase   *usecase.CategoryUsecase
	validator *validator.Validate
}

func NewCategoryController(usecase *usecase.CategoryUsecase) *CategoryController {
	return &CategoryController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *CategoryController) Create(ctx *gin.Context) {
	var req model.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreateCategory(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "category created successfully", res)
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	pageReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetCategories(ctx.Copy().Request.Context(), pageReq)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "categories fetched successfully", res, pagination)
}

func (c *CategoryController) GetCategoryById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := c.usecase.GetCategoryById(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "category fetched successfully", res)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateCategoryRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateCategory(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "category updated successfully", res)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.usecase.DeleteCategory(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "category deleted successfully", nil)
}

func handleError(ctx *gin.Context, err error) {
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		utils.ErrorResponse(ctx, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
}
