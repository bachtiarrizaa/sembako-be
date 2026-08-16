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

type PurchaseController struct {
	usecase   usecase.PurchaseUsecase
	validator *validator.Validate
}

func NewPurchaseController(usecase usecase.PurchaseUsecase) *PurchaseController {
	return &PurchaseController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *PurchaseController) CreatePurchase(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.CreatePurchase(ctx.Request.Context(), userID.String(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "purchase recorded successfully", res)
}

func (c *PurchaseController) GetPurchaseBatches(ctx *gin.Context) {
	var req model.GetPurchaseBatchesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	batches, pagination, err := c.usecase.GetPurchaseBatches(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "purchase batches retrieved successfully", batches, pagination)
}

func (c *PurchaseController) GetPurchaseBatchByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing purchase batch ID")
		return
	}

	batch, err := c.usecase.GetPurchaseBatchByID(ctx.Request.Context(), id)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "purchase batch retrieved successfully", batch)
}

func (c *PurchaseController) UpdatePurchase(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing purchase batch ID")
		return
	}

	var req model.UpdatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdatePurchase(ctx.Request.Context(), userID.String(), id, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "purchase batch updated successfully", res)
}

func (c *PurchaseController) DeletePurchase(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "missing purchase batch ID")
		return
	}

	err := c.usecase.DeletePurchase(ctx.Request.Context(), userID.String(), id)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "purchase batch deleted successfully", nil)
}
