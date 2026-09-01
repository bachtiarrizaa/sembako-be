package entity

import "time"

type PointLedgerType string

const (
	PointLedgerTypeEarn       PointLedgerType = "earn"
	PointLedgerTypeRedeem     PointLedgerType = "redeem"
	PointLedgerTypeExpire     PointLedgerType = "expire"
	PointLedgerTypeAdjustment PointLedgerType = "adjustment"
)

type PointLedger struct {
	ID            string          `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CustomerID    string          `gorm:"column:customer_id;type:uuid;not null" json:"customerId"`
	Customer      Customer        `gorm:"foreignKey:CustomerID" json:"customer"`
	TransactionID *string         `gorm:"column:transaction_id;type:uuid" json:"transactionId"`
	Transaction   *Transaction    `gorm:"foreignKey:TransactionID" json:"transaction"`
	Type          PointLedgerType `gorm:"column:type;type:varchar(20);not null" json:"type"`
	Points        int             `gorm:"column:points;not null" json:"points"`
	Description   string          `gorm:"column:description;not null;default:''" json:"description"`
	ExpiredAt     *time.Time      `gorm:"column:expired_at" json:"expiredAt"`
	CreatedAt     time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PointLedger) TableName() string {
	return "point_ledgers"
}
