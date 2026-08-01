package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (ctrl *UserController) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	res, err := ctrl.userUsecase.GetMe(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil mengambil profil", res)
}
