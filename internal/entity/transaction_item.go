package entity

import "time"

type TransactionItem struct {
	ID              string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TransactionID   string      `gorm:"column:transaction_id;type:uuid;not null" json:"transactionId"`
	ProductUnitID   string      `gorm:"column:product_unit_id;type:uuid;not null" json:"productUnitId"`
	ProductUnit     ProductUnit `gorm:"foreignKey:ProductUnitID" json:"productUnit"`
	Qty             float64     `gorm:"column:qty;type:numeric(14,4);not null" json:"qty"`
	UnitPrice       float64     `gorm:"column:unit_price;type:numeric(14,2);not null" json:"unitPrice"`
	DiscountApplied float64     `gorm:"column:discount_applied;type:numeric(14,2);not null;default:0" json:"discountApplied"`
	Subtotal        float64     `gorm:"column:subtotal;type:numeric(14,2);not null" json:"subtotal"`
	TotalCost       *float64    `gorm:"column:total_cost;type:numeric(14,2)" json:"totalCost"`
	Margin          *float64    `gorm:"column:margin;type:numeric(14,2)" json:"margin"`
	CreatedAt       time.Time   `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}
