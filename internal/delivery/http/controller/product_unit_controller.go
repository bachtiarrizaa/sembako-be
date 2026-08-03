package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (c *ProductController) UpdateProductUnitStatus(ctx *gin.Context) {
	productID := ctx.Param("id")
	unitID := ctx.Param("unitId")

	if productID == "" || unitID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id or unitId")
		return
	}

	res, err := c.usecase.UpdateProductUnitStatus(ctx.Request.Context(), productID, unitID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product unit status updated successfully", res)
}

func (c *ProductController) AddProductUnit(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.AddProductUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.AddProductUnit(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "product unit added successfully", res)
}

func (c *ProductController) UpdateProductUnit(ctx *gin.Context) {
	id := ctx.Param("id")
	unitID := ctx.Param("unitId")
	if id == "" || unitID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id or unitId")
		return
	}

	var req model.UpdateProductUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateProductUnit(ctx.Request.Context(), id, unitID, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product unit updated successfully", res)
}

func (c *ProductController) DeleteProductUnit(ctx *gin.Context) {
	id := ctx.Param("id")
	unitID := ctx.Param("unitId")
	if id == "" || unitID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id or unitId")
		return
	}

	if err := c.usecase.DeleteProductUnit(ctx.Request.Context(), id, unitID); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product unit deleted successfully", nil)
}
