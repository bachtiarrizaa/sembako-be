package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

const (
	ContextKeyUserID = "auth_user_id"
	ContextKeyRole   = "auth_role"
)

func AuthMiddleware(jwtAccessSecret string, blacklistRepo repository.BlacklistRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header not found")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header format must be 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseAccessToken(tokenString, jwtAccessSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Check if the token is blacklisted
		tokenHash := utils.HashRefreshToken(tokenString)
		isBlacklisted, err := blacklistRepo.IsBlacklisted(c.Request.Context(), tokenHash)
		if err != nil || isBlacklisted {
			utils.ErrorResponse(c, http.StatusUnauthorized, "token has been revoked")
			c.Abort()
			return
		}

		c.Set(utils.ContextKeyUserID, claims.UserID)
		c.Set(utils.ContextKeyRole, claims.Role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	return utils.GetUserID(c)
}

func GetRole(c *gin.Context) (string, bool) {
	return utils.GetRole(c)
}
