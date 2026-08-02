package model

import "time"

type CreateSupplierRequest struct {
	Name        string `json:"name" validate:"required"`
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

type UpdateSupplierRequest struct {
	Name        string `json:"name" validate:"required"`
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

type UpdateStatusSupplierRequest struct {
	IsActive bool `json:"isActive"`
}

type SupplierResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContactName string    `json:"contactName"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
