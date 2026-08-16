package entity

import "time"

type StockMutation struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID   string    `gorm:"column:product_id;type:uuid;not null" json:"productId"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product"`
	Type        string    `gorm:"column:type;type:varchar(10);not null" json:"type"` // 'in' or 'out'
	Qty         float64   `gorm:"column:qty;type:numeric(14,4);not null" json:"qty"`
	QtyBefore   float64   `gorm:"column:qty_before;type:numeric(14,4);not null" json:"qtyBefore"`
	QtyAfter    float64   `gorm:"column:qty_after;type:numeric(14,4);not null" json:"qtyAfter"`
	Source      string    `gorm:"column:source;type:varchar(30);not null" json:"source"` // 'purchase', 'stock_count'
	ReferenceID *string   `gorm:"column:reference_id;type:uuid" json:"referenceId"`
	Note        *string   `gorm:"column:note;type:text" json:"note"`
	CreatedBy   string    `gorm:"column:created_by;type:uuid;not null" json:"createdBy"`
	Creator     User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (StockMutation) TableName() string {
	return "stock_mutations"
}
