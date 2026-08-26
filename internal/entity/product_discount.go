package entity

import "time"

type ProductDiscount struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DiscountID string    `gorm:"column:discount_id;type:uuid;not null" json:"discountId"`
	Discount   Discount  `gorm:"foreignKey:DiscountID" json:"discount"`
	ProductID  string    `gorm:"column:product_id;type:uuid;not null" json:"productId"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product"`
	IsActive   bool      `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ProductDiscount) TableName() string {
	return "product_discounts"
}
