package entity

import "time"

type ProductUnit struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID        string    `gorm:"column:product_id;type:uuid;not null" json:"productId"`
	UnitID           string    `gorm:"column:unit_id;type:uuid;not null" json:"unitId"`
	Unit             Unit      `gorm:"foreignKey:UnitID" json:"unit"`
	Product          Product   `gorm:"foreignKey:ProductID" json:"product"`
	ConversionToBase float64   `gorm:"column:conversion_to_base;not null" json:"conversionToBase"`
	SellingPrice     float64   `gorm:"column:selling_price;not null" json:"sellingPrice"`
	IsBaseUnit       bool      `gorm:"column:is_base_unit;not null;default:false" json:"isBaseUnit"`
	IsActive         bool      `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ProductUnit) TableName() string {
	return "product_units"
}
