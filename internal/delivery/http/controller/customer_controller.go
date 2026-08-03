package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type CustomerController struct {
	usecase   *usecase.CustomerUsecase
	validator *validator.Validate
}

func NewCustomerController(usecase *usecase.CustomerUsecase) *CustomerController {
	return &CustomerController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var req model.CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreateCustomer(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "customer created successfully", res)
}

func (c *CustomerController) GetCustomersWithPagination(ctx *gin.Context) {
	pageReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetCustomersWithPagination(ctx.Request.Context(), pageReq)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "customers fetched successfully", res, pagination)
}

func (c *CustomerController) GetCustomerById(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.usecase.GetCustomerById(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "customer fetched successfully", res)
}

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	id := ctx.Param("id")

	var req model.UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateCustomer(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "customer updated successfully", res)
}

func (c *CustomerController) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req model.UpdateStatusCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := c.usecase.UpdateStatus(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "customer status updated successfully", res)
}

func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.usecase.DeleteCustomer(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "customer deleted successfully", nil)
}
