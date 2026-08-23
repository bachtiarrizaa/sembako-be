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

type ForceCloseShiftRequest struct {
	ClosingBalance  float64 `json:"closingBalance" validate:"required,gte=0"`
	Reason          string  `json:"reason" validate:"required,max=500"`
	DiscrepancyNote *string `json:"discrepancyNote" validate:"omitempty,max=500"`
}

type ListShiftsRequest struct {
	PaginationRequest
	StartDate *string `form:"startDate"`
	EndDate   *string `form:"endDate"`
	CashierID *string `form:"cashierId"`
}

type ShiftCashierResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ShiftListItemResponse struct {
	ID             string               `json:"id"`
	CashierID      string               `json:"cashierId"`
	Cashier        ShiftCashierResponse `json:"cashier"`
	OpeningBalance float64              `json:"openingBalance"`
	ClosingBalance *float64             `json:"closingBalance"`
	Status         string               `json:"status"`
	OpenedAt       time.Time            `json:"openedAt"`
	ClosedAt       *time.Time           `json:"closedAt"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}

func ToShiftListItemResponse(s *entity.Shift) ShiftListItemResponse {
	return ShiftListItemResponse{
		ID:             s.ID,
		CashierID:      s.CashierID,
		Cashier: ShiftCashierResponse{
			ID:   s.Cashier.ID,
			Name: s.Cashier.Name,
		},
		OpeningBalance: s.OpeningBalance,
		ClosingBalance: s.ClosingBalance,
		Status:         string(s.Status),
		OpenedAt:       s.OpenedAt,
		ClosedAt:       s.ClosedAt,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

type ShiftResponse struct {
	ID                string                `json:"id"`
	CashierID         string                `json:"cashierId"`
	Cashier           ShiftCashierResponse  `json:"cashier"`
	OpeningBalance    float64               `json:"openingBalance"`
	ClosingBalance    *float64              `json:"closingBalance"`
	SystemBalance     *float64              `json:"systemBalance"`
	Discrepancy       *float64              `json:"discrepancy"`
	DiscrepancyNote   *string               `json:"discrepancyNote"`
	Status            string                `json:"status"`
	ForceCloseReason  *string               `json:"forceCloseReason"`
	ForceClosedByUser *ShiftCashierResponse `json:"forceClosedByUser"`
	OpenedAt          time.Time             `json:"openedAt"`
	ClosedAt          *time.Time            `json:"closedAt"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

func ToShiftResponse(s *entity.Shift) ShiftResponse {
	resp := ShiftResponse{
		ID:        s.ID,
		CashierID: s.CashierID,
		Cashier: ShiftCashierResponse{
			ID:   s.Cashier.ID,
			Name: s.Cashier.Name,
		},
		OpeningBalance:   s.OpeningBalance,
		ClosingBalance:   s.ClosingBalance,
		SystemBalance:    s.SystemBalance,
		Discrepancy:      s.Discrepancy,
		DiscrepancyNote:  s.DiscrepancyNote,
		Status:           string(s.Status),
		ForceCloseReason: s.ForceCloseReason,
		OpenedAt:         s.OpenedAt,
		ClosedAt:         s.ClosedAt,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}

	if s.ForceClosedByUser != nil {
		resp.ForceClosedByUser = &ShiftCashierResponse{
			ID:   s.ForceClosedByUser.ID,
			Name: s.ForceClosedByUser.Name,
		}
	}

	return resp
}
