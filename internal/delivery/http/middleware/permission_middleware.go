package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/usecase"
)

func RequirePermission(usecase *usecase.PermissionUsecase, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		hasAccess, err := usecase.CheckUserPermission(c.Request.Context(), userID.String(), requiredPermission)
		if err != nil || !hasAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "you do not have permission to perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}
