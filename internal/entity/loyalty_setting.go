package entity

import "time"

type LoyaltySetting struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EarningRate    float64   `gorm:"column:earning_rate;type:numeric(14,2);not null" json:"earningRate"`
	RedemptionRate float64   `gorm:"column:redemption_rate;type:numeric(14,2);not null" json:"redemptionRate"`
	MinimumRedeem  int       `gorm:"column:minimum_redeem;not null" json:"minimumRedeem"`
	IsExpiryActive bool      `gorm:"column:is_expiry_active;not null" json:"isExpiryActive"`
	ExpiryMonths   int       `gorm:"column:expiry_months;not null" json:"expiryMonths"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (LoyaltySetting) TableName() string {
	return "loyalty_settings"
}
