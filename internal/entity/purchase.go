package entity

import (
	"time"
)

type Purchase struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InvoiceNumber *string   `gorm:"column:invoice_number" json:"invoiceNumber"`
	SupplierID    string    `gorm:"column:supplier_id;type:uuid;not null" json:"supplierId"`
	Supplier      Supplier  `gorm:"foreignKey:SupplierID" json:"supplier"`
	PurchaseDate  time.Time `gorm:"column:purchase_date;type:date;not null" json:"purchaseDate"`
	TotalAmount   float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	CreatedBy     string    `gorm:"column:created_by;type:uuid;not null" json:"createdBy"`
	Creator       User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Purchase) TableName() string {
	return "purchases"
}