package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bachtiarrizaa/sembako-be/internal/delivery/http/middleware"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

const refreshTokenCookieName = "refreshToken"

type AuthController struct {
	authUsecase          *usecase.AuthUsecase
	passwordResetUsecase *usecase.PasswordResetUsecase
	validator            *validator.Validate
	isProduction         bool
	refreshTokenTTL      time.Duration
}

func NewAuthController(
	authUsecase *usecase.AuthUsecase,
	passwordResetUsecase *usecase.PasswordResetUsecase,
	isProduction bool,
	refreshTokenTTL time.Duration,
) *AuthController {
	return &AuthController{
		authUsecase:          authUsecase,
		passwordResetUsecase: passwordResetUsecase,
		validator:            validator.New(),
		isProduction:         isProduction,
		refreshTokenTTL:      refreshTokenTTL,
	}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := ctrl.validator.Struct(req); err != nil {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	result, err := ctrl.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	ctrl.setRefreshCookie(c, result.RefreshToken)
	utils.SuccessResponse(c, http.StatusOK, "login successful", result.Response)
}

func (ctrl *AuthController) Refresh(c *gin.Context) {
	rawToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || rawToken == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "refresh token not found")
		return
	}

	resp, newRawToken, err := ctrl.authUsecase.Refresh(c.Request.Context(), rawToken)
	if err != nil {
		handleError(c, err)
		return
	}

	ctrl.setRefreshCookie(c, newRawToken)
	utils.SuccessResponse(c, http.StatusOK, "token refreshed", resp)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	authHeader := c.GetHeader("Authorization")
	var accessToken string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	}

	if err := ctrl.authUsecase.Logout(c.Request.Context(), userID, accessToken); err != nil {
		handleError(c, err)
		return
	}

	ctrl.clearRefreshCookie(c)
	utils.SuccessResponse(c, http.StatusOK, "logout successful", nil)
}

func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := ctrl.validator.Struct(req); err != nil {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	if err := ctrl.passwordResetUsecase.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "if the email exists, a password reset link has been sent", nil)
}

func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := ctrl.validator.Struct(req); err != nil {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, utils.FormatValidationError(err))
		return
	}

	if err := ctrl.passwordResetUsecase.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		handleError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "password reset successful", nil)
}

func (ctrl *AuthController) setRefreshCookie(c *gin.Context, rawToken string) {
	maxAge := int(ctrl.refreshTokenTTL.Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		rawToken,
		maxAge,
		"/api/auth",
		"",
		ctrl.isProduction,
		true,
	)
}

func (ctrl *AuthController) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		"/api/auth",
		"",
		ctrl.isProduction,
		true,
	)
}
