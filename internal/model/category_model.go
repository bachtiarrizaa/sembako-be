package model

import "time"

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
