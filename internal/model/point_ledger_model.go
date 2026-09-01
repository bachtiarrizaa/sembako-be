package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type PointLedgerResponse struct {
	ID            string     `json:"id"`
	CustomerID    string     `json:"customerId"`
	TransactionID *string    `json:"transactionId,omitempty"`
	Type          string     `json:"type"`
	Points        int        `json:"points"`
	Description   string     `json:"description"`
	ExpiredAt     *time.Time `json:"expiredAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func ToPointLedgerResponse(p *entity.PointLedger) PointLedgerResponse {
	return PointLedgerResponse{
		ID:            p.ID,
		CustomerID:    p.CustomerID,
		TransactionID: p.TransactionID,
		Type:          string(p.Type),
		Points:        p.Points,
		Description:   p.Description,
		ExpiredAt:     p.ExpiredAt,
		CreatedAt:     p.CreatedAt,
	}
}
