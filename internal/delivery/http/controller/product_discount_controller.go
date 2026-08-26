package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductDiscountController struct {
	productDiscountUsecase *usecase.ProductDiscountUsecase
	validator              *validator.Validate
}

func NewProductDiscountController(usecase *usecase.ProductDiscountUsecase) *ProductDiscountController {
	return &ProductDiscountController{
		productDiscountUsecase: usecase,
		validator:              validator.New(),
	}
}
func (c *ProductDiscountController) Create(ctx *gin.Context) {
	var req model.CreateProductDiscountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.productDiscountUsecase.Create(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "product discount created successfully", res)
}
