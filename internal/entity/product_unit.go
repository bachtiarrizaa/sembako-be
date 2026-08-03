package entity

import "time"

type ProductUnit struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID        string    `gorm:"column:product_id;type:uuid;not null" json:"product_id"`
	UnitID           string    `gorm:"column:unit_id;type:uuid;not null" json:"unit_id"`
	Unit             Unit      `gorm:"foreignKey:UnitID" json:"unit"`
	ConversionToBase float64   `gorm:"column:conversion_to_base;not null" json:"conversion_to_base"`
	SellingPrice     float64   `gorm:"column:selling_price;not null" json:"selling_price"`
	IsBaseUnit       bool      `gorm:"column:is_base_unit;not null;default:false" json:"is_base_unit"`
	IsActive         bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (ProductUnit) TableName() string {
	return "product_units"
}
