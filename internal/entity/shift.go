package entity

import "time"

type ShiftStatus string

const (
	ShiftStatusOpen   ShiftStatus = "open"
	ShiftStatusClosed ShiftStatus = "closed"
)

type Shift struct {
	ID                string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CashierID         string      `gorm:"column:cashier_id;type:uuid;not null" json:"cashierId"`
	Cashier           User        `gorm:"foreignKey:CashierID" json:"cashier"`
	OpeningBalance    float64     `gorm:"column:opening_balance;type:numeric(14,2);not null" json:"openingBalance"`
	ClosingBalance    *float64    `gorm:"column:closing_balance;type:numeric(14,2)" json:"closingBalance"`
	SystemBalance     *float64    `gorm:"column:system_balance;type:numeric(14,2)" json:"systemBalance"`
	Discrepancy       *float64    `gorm:"column:discrepancy;type:numeric(14,2)" json:"discrepancy"`
	DiscrepancyNote   *string     `gorm:"column:discrepancy_note;type:text" json:"discrepancyNote"`
	Status            ShiftStatus `gorm:"column:status;type:varchar(10);not null;default:open" json:"status"`
	ForceCloseReason  *string     `gorm:"column:force_close_reason;type:text" json:"forceCloseReason"`
	ForceClosedBy     *string     `gorm:"column:force_closed_by;type:uuid" json:"forceClosedBy"`
	ForceClosedByUser *User       `gorm:"foreignKey:ForceClosedBy" json:"forceClosedByUser,omitempty"`
	OpenedAt          time.Time   `gorm:"column:opened_at;autoCreateTime" json:"openedAt"`
	ClosedAt          *time.Time  `gorm:"column:closed_at" json:"closedAt"`
	CreatedAt         time.Time   `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Shift) TableName() string {
	return "shifts"
}
