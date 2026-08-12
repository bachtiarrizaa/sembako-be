package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type PermissionController struct {
	usecase *usecase.PermissionUsecase
}

func NewPermissionController(usecase *usecase.PermissionUsecase) *PermissionController {
	return &PermissionController{usecase: usecase}
}

func (c *PermissionController) GetUserMenu(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := c.usecase.GetUserMenu(ctx.Request.Context(), userID.String())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to fetch user menu")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "user menus fetched successfully", res)
}
