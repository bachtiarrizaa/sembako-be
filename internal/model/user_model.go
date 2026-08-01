package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type UserWithRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Username  *string      `json:"username"`
	Role      UserWithRole `json:"role"`
	IsActive  bool         `json:"isActive"`
	CreatedAt time.Time    `json:"createdAt"`
}

func ToUserResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Username: u.Username,
		Role: UserWithRole{
			ID:   u.Role.ID,
			Name: u.Role.Name,
		},
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
