package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                     string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID             string         `gorm:"column:category_id;type:uuid;not null" json:"category_id"`
	Category               Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Name                   string         `gorm:"column:name;not null" json:"name"`
	BaseUnitID             string         `gorm:"column:base_unit_id;type:uuid;not null" json:"base_unit_id"`
	BaseUnit               Unit           `gorm:"foreignKey:BaseUnitID" json:"base_unit"`
	MinimumStock           *float64       `gorm:"column:minimum_stock" json:"minimum_stock"`
	MarginThresholdPercent *float64       `gorm:"column:margin_threshold_percent" json:"margin_threshold_percent"`
	IsActive               bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	Units                  []ProductUnit  `gorm:"foreignKey:ProductID" json:"units,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
