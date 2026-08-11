package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
	validator   *validator.Validate
	uploadDir   string
}

func NewUserController(userUsecase *usecase.UserUsecase, uploadDir string) *UserController {
	return &UserController{
		userUsecase: userUsecase,
		validator:   validator.New(),
		uploadDir:   uploadDir,
	}
}

func (c *UserController) GetProfileMe(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := c.userUsecase.GetMe(ctx.Request.Context(), userID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "profile fetched successfully", res)
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.UpdateProfileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	uploadCfg := utils.DefaultImageConfig(filepath.Join(c.uploadDir, "profiles"))
	result, err := utils.HandleFileUpload(ctx, uploadCfg)
	if err != nil {
		if uploadErr, ok := err.(*utils.UploadError); ok {
			utils.ErrorResponse(ctx, http.StatusBadRequest, uploadErr.Message)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to process uploaded file")
		return
	}

	var imagePath *string
	if result != nil {
		imagePath = &result.FilePath
	}

	res, err := c.userUsecase.UpdateProfile(ctx.Request.Context(), userID, req, imagePath)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "profile updated successfully", res)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.userUsecase.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "user created successfully", res)
}

func (c *UserController) GetUsersWithPagination(ctx *gin.Context) {
	pagReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.userUsecase.GetUsersWithPagination(ctx.Request.Context(), pagReq)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "users fetched successfully", res, pagination)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}

	res, err := c.userUsecase.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "user fetched successfully", res)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}

	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.userUsecase.UpdateUser(ctx.Request.Context(), id, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "user updated successfully", res)
}

func (c *UserController) UpdateStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateStatusUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := c.validator.Struct(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	res, err := c.userUsecase.UpdateStatus(ctx.Request.Context(), id, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "user status updated successfully", res)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}

	if err := c.userUsecase.DeleteUser(ctx.Request.Context(), id); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "user deleted successfully", nil)
}
