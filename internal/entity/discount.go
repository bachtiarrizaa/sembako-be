package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type DiscountType string

const (
	DiscountTypePercent DiscountType = "percent"
	DiscountTypeFixed   DiscountType = "fixed"
)

type Discount struct {
	ID               string            `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name             string            `gorm:"column:name;type:varchar(150);not null" json:"name"`
	Type             DiscountType      `gorm:"column:type;type:varchar(10);not null" json:"type"`
	Value            decimal.Decimal   `gorm:"column:value;type:numeric(14,2);not null" json:"value"`
	StartDate        *time.Time        `gorm:"column:start_date" json:"startDate"`
	EndDate          *time.Time        `gorm:"column:end_date" json:"endDate"`
	IsActive         bool              `gorm:"column:is_active;not null;default:true" json:"isActive"`
	ProductDiscounts []ProductDiscount `gorm:"foreignKey:DiscountID"`
	CreatedAt        time.Time         `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Discount) TableName() string {
	return "discounts"
}
