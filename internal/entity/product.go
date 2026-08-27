package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                     string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID             string         `gorm:"column:category_id;type:uuid;not null" json:"categoryId"`
	Category               Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Name                   string         `gorm:"column:name;not null" json:"name"`
	Image                  *string        `gorm:"column:image" json:"image"`
	BaseUnitID             string         `gorm:"column:base_unit_id;type:uuid;not null" json:"baseUnitId"`
	BaseUnit               Unit           `gorm:"foreignKey:BaseUnitID" json:"baseUnit"`
	MinimumStock           *float64       `gorm:"column:minimum_stock" json:"minimumStock"`
	MarginThresholdPercent *float64       `gorm:"column:margin_threshold_percent" json:"marginThresholdPercent"`
	IsActive               bool           `gorm:"column:is_active;not null;default:true" json:"isActive"`
	Units                  []ProductUnit     `gorm:"foreignKey:ProductID" json:"units,omitempty"`
	Stock                  *Stock            `gorm:"foreignKey:ProductID" json:"stock,omitempty"`
	ProductDiscounts       []ProductDiscount `gorm:"foreignKey:ProductID" json:"productDiscounts,omitempty"`
	CreatedAt              time.Time         `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt              gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
