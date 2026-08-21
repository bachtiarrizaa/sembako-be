package controller

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type ProductController struct {
	usecase   *usecase.ProductUsecase
	validator *validator.Validate
	uploadDir string
}

func NewProductController(usecase *usecase.ProductUsecase, uploadDir string) *ProductController {
	return &ProductController{
		usecase:   usecase,
		validator: validator.New(),
		uploadDir: uploadDir,
	}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req model.CreateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := json.Unmarshal([]byte(ctx.PostForm("units")), &req.Units); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid units format")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	uploadCfg := utils.DefaultImageConfig(filepath.Join(c.uploadDir, "products"))
	result, err := utils.HandleFileUpload(ctx, uploadCfg)
	if err != nil {
		if uploadErr, ok := err.(*utils.UploadError); ok {
			utils.ErrorResponse(ctx, http.StatusBadRequest, uploadErr.Message)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to process upload file")
		return
	}

	var imagePath *string
	if result != nil {
		imagePath = &result.FilePath
	}

	res, err := c.usecase.CreateProduct(ctx.Request.Context(), req, imagePath)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "product created successfully", res)
}

func (c *ProductController) GetProducts(ctx *gin.Context) {
	var req model.GetProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	res, pagination, err := c.usecase.GetProducts(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "products fetched successfully", res, pagination)
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.usecase.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product fetched successfully", res)
}

func (c *ProductController) UpdateProductStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateProductStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.usecase.UpdateProductStatus(ctx.Request.Context(), id, *req.IsActive)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product status updated successfully", res)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := json.Unmarshal([]byte(ctx.PostForm("units")), &req.Units); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid units format")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	uploadCfg := utils.DefaultImageConfig(filepath.Join(c.uploadDir, "products"))
	result, err := utils.HandleFileUpload(ctx, uploadCfg)
	if err != nil {
		if uploadErr, ok := err.(*utils.UploadError); ok {
			utils.ErrorResponse(ctx, http.StatusBadRequest, uploadErr.Message)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to process upload file")
		return
	}

	var imagePath *string
	if result != nil {
		imagePath = &result.FilePath
	}

	res, err := c.usecase.UpdateProduct(ctx.Request.Context(), id, req, imagePath)
	if err != nil {
		// If upload succeeded but product update failed, clean up the newly uploaded image
		if imagePath != nil {
			utils.DeleteFile(*imagePath)
		}
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product updated successfully", res)
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.usecase.DeleteProduct(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "product deleted successfully", nil)
}
