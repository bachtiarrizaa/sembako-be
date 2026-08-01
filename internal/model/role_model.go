package model

import "time"

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type RoleResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
