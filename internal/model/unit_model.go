package model

import "time"

type CreateUnitRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateUnitRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UnitResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
