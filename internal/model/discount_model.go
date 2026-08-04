package model

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type CreateDiscountRequest struct {
	Name      string          `json:"name" validate:"required,min=2,max=150"`
	Type      string          `json:"type" validate:"required,oneof=percent fixed"`
	Value     decimal.Decimal `json:"value" validate:"required"`
	StartDate *time.Time      `json:"startDate" validate:"omitempty"`
	EndDate   *time.Time      `json:"endDate" validate:"omitempty"`
}

type UpdateDiscountRequest struct {
	Name      *string             `json:"name" validate:"omitempty,min=2,max=150"`
	Type      *string             `json:"type" validate:"omitempty,oneof=percent fixed"`
	Value     *decimal.Decimal    `json:"value" validate:"omitempty"`
	StartDate Nullable[time.Time] `json:"startDate"`
	EndDate   Nullable[time.Time] `json:"endDate"`
	IsActive  *bool               `json:"isActive"`
}

type UpdateStatusDiscountRequest struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

type DiscountResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Value     decimal.Decimal `json:"value"`
	StartDate *time.Time      `json:"startDate"`
	EndDate   *time.Time      `json:"endDate"`
	IsActive  bool            `json:"isActive"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

func ToDiscountResponse(d *entity.Discount) DiscountResponse {
	return DiscountResponse{
		ID:        d.ID,
		Name:      d.Name,
		Type:      d.Type,
		Value:     d.Value,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
