package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase/transaction"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionController struct {
	usecase   transaction.TransactionUsecase
	validator *validator.Validate
}

func NewTransactionController(usecase transaction.TransactionUsecase) *TransactionController {
	return &TransactionController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return
	}

	var req model.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreateTransaction(ctx.Request.Context(), userID.String(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "transaction created successfully", res)
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	userID, role, ok := utils.GetUserAndRole(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid transaction id")
		return
	}

	res, err := c.usecase.GetTransactionByID(ctx.Request.Context(), id, userID.String(), role)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "transaction fetched successfully", res)
}

func (c *TransactionController) ListTransactions(ctx *gin.Context) {
	userID, role, ok := utils.GetUserAndRole(ctx)
	if !ok {
		return
	}

	var req model.ListTransactionsRequest
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

	res, pagination, err := c.usecase.ListTransactions(ctx.Request.Context(), req, userID.String(), role)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "transactions fetched successfully", res, pagination)
}

func (c *TransactionController) VoidTransaction(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req model.VoidTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.VoidTransaction(ctx.Request.Context(), id, userID.String(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "transaction voided successfully", res)
}
