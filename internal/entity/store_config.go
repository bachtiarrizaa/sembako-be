package entity

import "time"

type StoreConfiguration struct {
	ID                        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreName                 string    `gorm:"column:store_name;type:varchar(100);not null" json:"storeName"`
	StoreAddress              *string   `gorm:"column:store_address;type:text" json:"storeAddress"`
	StorePhone                *string   `gorm:"column:store_phone;type:varchar(20)" json:"storePhone"`
	ReceiptHeaderText         *string   `gorm:"column:receipt_header_text;type:text" json:"receiptHeaderText"`
	ReceiptFooterText         *string   `gorm:"column:receipt_footer_text;type:text" json:"receiptFooterText"`
	ReceiptShowCashierName    bool      `gorm:"column:receipt_show_cashier_name" json:"receiptShowCashierName"`
	ReceiptShowCustomerName   bool      `gorm:"column:receipt_show_customer_name" json:"receiptShowCustomerName"`
	ShiftDiscrepancyTolerance float64   `gorm:"column:shift_discrepancy_tolerance;type:numeric(14,2);not null;default:1000" json:"shiftDiscrepancyTolerance"`
	CreatedAt                 time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (StoreConfiguration) TableName() string {
	return "store_configurations"
}
