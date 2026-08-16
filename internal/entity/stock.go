package entity

import "time"

type Stock struct {
	ProductID   string    `gorm:"column:product_id;type:uuid;primaryKey" json:"productId"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product"`
	QtyBaseUnit float64   `gorm:"column:qty_base_unit;type:numeric(14,4);not null;default:0" json:"qtyBaseUnit"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Stock) TableName() string {
	return "stocks"
}
