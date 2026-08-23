package entity

import "time"

type Transaction struct {
	ID                     string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReceiptNumber          string            `gorm:"column:receipt_number;type:varchar(30);unique;not null" json:"receiptNumber"`
	CashierID              string            `gorm:"column:cashier_id;type:uuid;not null" json:"cashierId"`
	Cashier                User              `gorm:"foreignKey:CashierID" json:"cashier"`
	ShiftID                string            `gorm:"column:shift_id;type:uuid;not null" json:"shiftId"`
	Shift                  Shift             `gorm:"foreignKey:ShiftID" json:"shift"`
	CustomerID             *string           `gorm:"column:customer_id;type:uuid" json:"customerId"`
	Customer               *Customer         `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	PaymentMethod          string            `gorm:"column:payment_method;type:varchar(10);not null" json:"paymentMethod"`
	Subtotal               float64           `gorm:"column:subtotal;type:numeric(14,2);not null" json:"subtotal"`
	TotalDiscount          float64           `gorm:"column:total_discount;type:numeric(14,2);not null;default:0" json:"totalDiscount"`
	PointsUsed             int               `gorm:"column:points_used;type:integer;not null;default:0" json:"pointsUsed"`
	PointsDiscountValue    float64           `gorm:"column:points_discount_value;type:numeric(14,2);not null;default:0" json:"pointsDiscountValue"`
	PointsEarned           int               `gorm:"column:points_earned;type:integer;not null;default:0" json:"pointsEarned"`
	Total                  float64           `gorm:"column:total;type:numeric(14,2);not null" json:"total"`
	CashReceived           *float64          `gorm:"column:cash_received;type:numeric(14,2)" json:"cashReceived"`
	ChangeGiven            *float64          `gorm:"column:change_given;type:numeric(14,2)" json:"changeGiven"`
	ManualPaidConfirmation *bool             `gorm:"column:manual_paid_confirmation;type:boolean" json:"manualPaidConfirmation"`
	Status                 string            `gorm:"column:status;type:varchar(10);not null;default:completed" json:"status"`
	VoidReason             *string           `gorm:"column:void_reason;type:text" json:"voidReason"`
	VoidedBy               *string           `gorm:"column:voided_by;type:uuid" json:"voidedBy"`
	VoidedByUser           *User             `gorm:"foreignKey:VoidedBy" json:"voidedByUser,omitempty"`
	VoidedAt               *time.Time        `gorm:"column:voided_at" json:"voidedAt"`
	CreatedAt              time.Time         `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Items                  []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}
