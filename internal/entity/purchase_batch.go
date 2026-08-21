package entity

import (
	"time"
)

type PurchaseBatch struct {
	ID             string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PurchaseID     *string      `gorm:"column:purchase_id;type:uuid" json:"purchaseId"`
	ProductID      string       `gorm:"column:product_id;type:uuid;not null" json:"productId"`
	Product        Product      `gorm:"foreignKey:ProductID" json:"product"`
	SupplierID     string       `gorm:"column:supplier_id;type:uuid;not null" json:"supplierId"`
	Supplier       Supplier     `gorm:"foreignKey:SupplierID" json:"supplier"`
	UnitID         *string      `gorm:"column:unit_id;type:uuid" json:"unitId"`
	Unit           *ProductUnit `gorm:"foreignKey:UnitID" json:"unit"`
	UnitPrice      *float64     `gorm:"column:unit_price" json:"unitPrice"`
	UnitSourceID   string       `gorm:"->;column:unit_source_id;-:migration" json:"-"`
	UnitSourceName string       `gorm:"->;column:unit_source_name;-:migration" json:"-"`
	InitialQty     float64      `gorm:"column:initial_qty;not null" json:"initialQuantity"`
	RemainingQty   float64      `gorm:"column:remaining_qty;not null" json:"remainingQuantity"`
	PurchasePrice  float64      `gorm:"column:purchase_price;not null" json:"purchasePrice"`
	InvoiceNumber  *string      `gorm:"column:invoice_number" json:"invoiceNumber"`
	PurchaseDate   time.Time    `gorm:"column:purchase_date;type:date;not null" json:"purchaseDate"`
	CreatedBy      string       `gorm:"column:created_by;type:uuid;not null" json:"createdBy"`
	Creator        User         `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PurchaseBatch) TableName() string {
	return "purchase_batches"
}
