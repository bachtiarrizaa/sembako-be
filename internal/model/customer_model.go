package model

import "time"

type CreateCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Address     string `json:"address"`
}

type UpdateCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Address     string `json:"address"`
}

type UpdateStatusCustomerRequest struct {
	IsActive bool `json:"isActive"`
}

type CustomerResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	Address     string    `json:"address"`
	TotalPoints int       `json:"totalPoints"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
