package entity

import "time"

type StockCount struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID   string     `gorm:"column:product_id;type:uuid;not null" json:"productId"`
	Product     Product    `gorm:"foreignKey:ProductID" json:"product"`
	CountDate   time.Time  `gorm:"column:count_date;type:date;not null" json:"countDate"`
	SystemQty   float64    `gorm:"column:system_qty;type:numeric(14,4);not null" json:"systemQty"`
	PhysicalQty float64    `gorm:"column:physical_qty;type:numeric(14,4);not null" json:"physicalQty"`
	Discrepancy float64    `gorm:"column:discrepancy;type:numeric(14,4);not null" json:"discrepancy"`
	Note        *string    `gorm:"column:note;type:text" json:"note"`
	Status      string     `gorm:"column:status;type:varchar(10);not null;default:pending" json:"status"` // 'pending', 'approved', 'rejected'
	SubmittedBy string     `gorm:"column:submitted_by;type:uuid;not null" json:"submittedBy"`
	Submitter   User       `gorm:"foreignKey:SubmittedBy" json:"submitter"`
	ApprovedBy  *string    `gorm:"column:approved_by;type:uuid" json:"approvedBy"`
	Approver    *User      `gorm:"foreignKey:ApprovedBy" json:"approver"`
	SubmittedAt time.Time  `gorm:"column:submitted_at;autoCreateTime" json:"submittedAt"`
	ApprovedAt  *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (StockCount) TableName() string {
	return "stock_counts"
}
