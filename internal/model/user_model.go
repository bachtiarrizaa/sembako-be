package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=150"`
	Email    string  `json:"email" validate:"required,email"`
	Username *string `json:"username" validate:"omitempty,min=2,max=100"`
	Password string  `json:"password" validate:"required,min=8"`
	RoleID   string  `json:"roleId" validate:"required,uuid"`
	IsActive *bool   `json:"isActive" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=150"`
	Email    string  `json:"email" validate:"required,email"`
	Username *string `json:"username" validate:"omitempty,min=2,max=100"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	RoleID   string  `json:"roleId" validate:"required,uuid"`
	IsActive *bool   `json:"isActive" validate:"required"`
}

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
