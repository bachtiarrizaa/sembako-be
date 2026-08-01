package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type RoleController struct {
	usecase   *usecase.RoleUsecase
	validator *validator.Validate
}

func NewRoleController(usecase *usecase.RoleUsecase) *RoleController {
	return &RoleController{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (c *RoleController) Create(ctx *gin.Context) {
	var req model.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := c.usecase.CreateRole(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "role created successfully", res)
}

func (c *RoleController) GetAll(ctx *gin.Context) {
	pagReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetAllRoles(ctx.Request.Context(), pagReq)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "roles fetched successfully", res, pagination)
}

func (c *RoleController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := c.usecase.GetRoleByID(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "role fetched successfully", res)
}

func (c *RoleController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := c.usecase.UpdateRole(ctx.Request.Context(), id, req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "role updated successfully", res)
}

func (c *RoleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.usecase.DeleteRole(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "role deleted successfully", nil)
}

func handleError(ctx *gin.Context, err error) {
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		utils.ErrorResponse(ctx, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
}
