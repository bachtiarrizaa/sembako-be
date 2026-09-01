package controller

import (
	"net/http"

	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PointLedgerController struct {
	usecase *usecase.PointLedgerUsecase
}

func NewPointLedgerController(usecase *usecase.PointLedgerUsecase) *PointLedgerController {
	return &PointLedgerController{usecase: usecase}
}

func (c *PointLedgerController) GetCustomerLedgers(ctx *gin.Context) {
	customerID := ctx.Param("id")
	pageReq, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid query params")
		return
	}

	res, pagination, err := c.usecase.GetCustomerLedgers(ctx.Request.Context(), customerID, pageReq)
	if err != nil {
		handleError(ctx, err)
		return
	}

	utils.SuccessResponseWithPagination(ctx, http.StatusOK, "customer point ledgers fetched successfully", res, pagination)
}
