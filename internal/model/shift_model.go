package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type OpenShiftRequest struct {
	OpeningBalance float64 `json:"openingBalance" validate:"required,gte=0"`
}

type CloseShiftRequest struct {
	ClosingBalance  float64 `json:"closingBalance" validate:"required,gte=0"`
	DiscrepancyNote *string `json:"discrepancyNote" validate:"omitempty,max=500"`
}

type ShiftCashierResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ShiftResponse struct {
	ID              string               `json:"id"`
	CashierID       string               `json:"cashierId"`
	Cashier         ShiftCashierResponse `json:"cashier"`
	OpeningBalance  float64              `json:"openingBalance"`
	ClosingBalance  *float64             `json:"closingBalance"`
	SystemBalance   *float64             `json:"systemBalance"`
	Discrepancy     *float64             `json:"discrepancy"`
	DiscrepancyNote *string              `json:"discrepancyNote"`
	Status          string               `json:"status"`
	OpenedAt        time.Time            `json:"openedAt"`
	ClosedAt        *time.Time           `json:"closedAt"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
}

func ToShiftResponse(s *entity.Shift) ShiftResponse {
	return ShiftResponse{
		ID:        s.ID,
		CashierID: s.CashierID,
		Cashier: ShiftCashierResponse{
			ID:    s.Cashier.ID,
			Name:  s.Cashier.Name,
			Email: s.Cashier.Email,
		},
		OpeningBalance:  s.OpeningBalance,
		ClosingBalance:  s.ClosingBalance,
		SystemBalance:   s.SystemBalance,
		Discrepancy:     s.Discrepancy,
		DiscrepancyNote: s.DiscrepancyNote,
		Status:          string(s.Status),
		OpenedAt:        s.OpenedAt,
		ClosedAt:        s.ClosedAt,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}
