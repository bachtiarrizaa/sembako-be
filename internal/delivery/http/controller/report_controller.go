package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReportController struct {
	usecase   usecase.ReportUsecase
	validator *validator.Validate
}

func NewReportController(usecase usecase.ReportUsecase) *ReportController {
	return &ReportController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *ReportController) GetReport(ctx *gin.Context) {
	var req model.GetReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	switch req.Type {
	case "sales":
		res, pagination, err := c.usecase.GetSalesReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponseWithPagination(ctx, http.StatusOK, "sales report fetched successfully", res, pagination)
	case "profit_margin":
		res, pagination, err := c.usecase.GetProfitMarginReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponseWithPagination(ctx, http.StatusOK, "profit margin report fetched successfully", res, pagination)
	case "payment_method":
		res, err := c.usecase.GetPaymentMethodReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponse(ctx, http.StatusOK, "payment method report fetched successfully", res)
	case "cashier":
		res, pagination, err := c.usecase.GetCashierReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponseWithPagination(ctx, http.StatusOK, "cashier report fetched successfully", res, pagination)
	case "top_selling":
		res, pagination, err := c.usecase.GetTopSellingReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponseWithPagination(ctx, http.StatusOK, "top selling products report fetched successfully", res, pagination)
	case "inventory_valuation":
		res, pagination, err := c.usecase.GetInventoryValuationReport(ctx.Request.Context(), req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}
		utils.SuccessResponseWithPagination(ctx, http.StatusOK, "inventory valuation report fetched successfully", res, pagination)
	default:
		utils.ErrorResponse(ctx, http.StatusBadRequest, "unsupported report type")
	}
}
