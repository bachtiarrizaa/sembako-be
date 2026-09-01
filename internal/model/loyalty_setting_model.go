package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type UpdateLoyaltySettingRequest struct {
	EarningRate    float64 `json:"earningRate" validate:"required,gt=0"`
	RedemptionRate float64 `json:"redemptionRate" validate:"required,gt=0"`
	MinimumRedeem  int     `json:"minimumRedeem" validate:"gte=0"`
	IsExpiryActive bool    `json:"isExpiryActive"`
	ExpiryMonths   int     `json:"expiryMonths" validate:"gte=0"`
}

type LoyaltySettingResponse struct {
	ID             string    `json:"id"`
	EarningRate    float64   `json:"earningRate"`
	RedemptionRate float64   `json:"redemptionRate"`
	MinimumRedeem  int       `json:"minimumRedeem"`
	IsExpiryActive bool      `json:"isExpiryActive"`
	ExpiryMonths   int       `json:"expiryMonths"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func ToLoyaltySettingResponse(m *entity.LoyaltySetting) LoyaltySettingResponse {
	resp := LoyaltySettingResponse{
		ID:             m.ID,
		EarningRate:    m.EarningRate,
		RedemptionRate: m.RedemptionRate,
		MinimumRedeem:  m.MinimumRedeem,
		IsExpiryActive: m.IsExpiryActive,
		ExpiryMonths:   m.ExpiryMonths,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	return resp
}
