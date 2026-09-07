package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	usecase *usecase.DashboardUsecase
}

func NewDashboardController(usecase *usecase.DashboardUsecase) *DashboardController {
	return &DashboardController{usecase: usecase}
}

func (c *DashboardController) GetDashboard(ctx *gin.Context) {
	role, _ := utils.GetRole(ctx)
	userID, _ := utils.GetUserID(ctx)

	if role == "cashier" {
		res, err := c.usecase.GetCashierDashboard(ctx.Request.Context(), userID.String())
		if err != nil {
			handleError(ctx, err)
			return
		}
		utils.SuccessResponse(ctx, http.StatusOK, "cashier dashboard data fetched successfully", res)
		return
	}

	var req model.DashboardQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, err := c.usecase.GetAdminDashboard(ctx.Request.Context(), req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "admin dashboard data fetched successfully", res)
}
