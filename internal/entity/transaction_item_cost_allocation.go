package entity

import "time"

type TransactionItemCostAllocation struct {
	ID                  string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TransactionItemID   string          `gorm:"column:transaction_item_id;type:uuid;not null" json:"transactionItemId"`
	TransactionItem     TransactionItem `gorm:"foreignKey:TransactionItemID" json:"transactionItem"`
	PurchaseBatchID     string          `gorm:"column:purchase_batch_id;type:uuid;not null" json:"purchaseBatchId"`
	PurchaseBatch       PurchaseBatch   `gorm:"foreignKey:PurchaseBatchID" json:"purchaseBatch"`
	QtyAllocated        float64         `gorm:"column:qty_allocated;type:numeric(14,4);not null" json:"qtyAllocated"`
	PurchasePriceAtSale float64         `gorm:"column:purchase_price_at_sale;type:numeric(14,2);not null" json:"purchasePriceAtSale"`
	CostSubtotal        float64         `gorm:"column:cost_subtotal;type:numeric(14,2);not null" json:"costSubtotal"`
	CreatedAt           time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (TransactionItemCostAllocation) TableName() string {
	return "transaction_item_cost_allocations"
}
