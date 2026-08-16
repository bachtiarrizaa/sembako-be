package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type StockController struct {
	usecase   usecase.StockUsecase
	validator *validator.Validate
}

func NewStockController(usecase usecase.StockUsecase) *StockController {
	return &StockController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *StockController) SubmitStockCount(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.SubmitStockCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.SubmitStockCount(ctx.Request.Context(), userID.String(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "stock count submitted successfully", res)
}

func (c *StockController) GetStockCounts(ctx *gin.Context) {
	var req model.GetStockCountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	// Apply default pagination values if empty
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	counts, pagination, err := c.usecase.GetStockCounts(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "stock counts retrieved successfully", counts, pagination)
}

func (c *StockController) ApproveStockCount(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing stock count ID")
		return
	}

	var req model.ApproveStockCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.ApproveStockCount(ctx.Request.Context(), userID.String(), id, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	msg := "stock count approved successfully"
	if !req.Approve {
		msg = "stock count rejected successfully"
	}

	utils.SuccessResponse(ctx, http.StatusOK, msg, res)
}

func (c *StockController) GetStockByProductID(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if productID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing product ID")
		return
	}

	res, err := c.usecase.GetStockByProductID(ctx.Request.Context(), productID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "stock summary retrieved successfully", res)
}

func (c *StockController) GetStockMutations(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if productID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing product ID")
		return
	}

	var req model.GetStockMutationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	mutations, pagination, err := c.usecase.GetStockMutations(ctx.Request.Context(), productID, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "stock mutations retrieved successfully", mutations, pagination)
}
