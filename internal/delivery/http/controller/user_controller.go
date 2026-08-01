package controller

import (
	"net/http"

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
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
		validator:   validator.New(),
	}
}

func (ctrl *UserController) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := ctrl.userUsecase.GetMe(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "profile fetched successfully", res)
}

func (ctrl *UserController) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := ctrl.validator.Struct(req); err != nil {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := ctrl.userUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "user created successfully", res)
}

func (ctrl *UserController) GetAll(c *gin.Context) {
	pagReq, err := utils.ParsePaginationQuery(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := ctrl.userUsecase.GetAllUsers(c.Request.Context(), pagReq)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponseWithPagination(c, http.StatusOK, "users fetched successfully", res, pagination)
}

func (ctrl *UserController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	res, err := ctrl.userUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user fetched successfully", res)
}

func (ctrl *UserController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := ctrl.validator.Struct(req); err != nil {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := ctrl.userUsecase.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user updated successfully", res)
}

func (ctrl *UserController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	if err := ctrl.userUsecase.DeleteUser(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user deleted successfully", nil)
}
