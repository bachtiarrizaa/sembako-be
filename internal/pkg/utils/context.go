package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextKeyUserID = "auth_user_id"
	ContextKeyRole   = "auth_role"
)

// GetUserID retrieves the authenticated user ID from Gin context. If not found, it responds with 401 Unauthorized.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}
	return id, true
}

// GetRole retrieves the authenticated user's role from Gin context. If not found, it responds with 401 Unauthorized.
func GetRole(c *gin.Context) (string, bool) {
	val, exists := c.Get(ContextKeyRole)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return "", false
	}
	role, ok := val.(string)
	if !ok {
		ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return "", false
	}
	return role, true
}

// GetUserAndRole retrieves both authenticated user ID and role from Gin context. If not found, it responds with 401 Unauthorized.
func GetUserAndRole(c *gin.Context) (uuid.UUID, string, bool) {
	userID, ok := GetUserID(c)
	if !ok {
		return uuid.Nil, "", false
	}
	role, ok := GetRole(c)
	if !ok {
		return uuid.Nil, "", false
	}
	return userID, role, true
}
