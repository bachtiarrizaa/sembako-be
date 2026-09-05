package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type UpdateStoreConfigurationRequest struct {
	StoreName                 string  `json:"storeName" validate:"required,min=2,max=100"`
	StoreAddress              string  `json:"storeAddress"`
	StorePhone                string  `json:"storePhone"`
	ReceiptHeaderText         string  `json:"receiptHeaderText"`
	ReceiptFooterText         string  `json:"receiptFooterText"`
	ReceiptShowCashierName    bool    `json:"receiptShowCashierName"`
	ReceiptShowCustomerName   bool    `json:"receiptShowCustomerName"`
	ShiftDiscrepancyTolerance float64 `json:"shiftDiscrepancyTolerance" validate:"gte=0"`
}

type StoreConfigurationResponse struct {
	ID                        string    `json:"id"`
	StoreName                 string    `json:"storeName"`
	StoreAddress              *string   `json:"storeAddress"`
	StorePhone                *string   `json:"storePhone"`
	ReceiptHeaderText         *string   `json:"receiptHeaderText"`
	ReceiptFooterText         *string   `json:"receiptFooterText"`
	ReceiptShowCashierName    bool      `json:"receiptShowCashierName"`
	ReceiptShowCustomerName   bool      `json:"receiptShowCustomerName"`
	ShiftDiscrepancyTolerance float64   `json:"shiftDiscrepancyTolerance"`
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}

type PublicStoreInfoResponse struct {
	StoreName string `json:"storeName"`
}

func ToStoreConfigurationResponse(s *entity.StoreConfiguration) StoreConfigurationResponse {
	resp := StoreConfigurationResponse{
		ID:                        s.ID,
		StoreName:                 s.StoreName,
		StoreAddress:              s.StoreAddress,
		StorePhone:                s.StorePhone,
		ReceiptHeaderText:         s.ReceiptHeaderText,
		ReceiptFooterText:         s.ReceiptFooterText,
		ReceiptShowCashierName:    s.ReceiptShowCashierName,
		ReceiptShowCustomerName:   s.ReceiptShowCustomerName,
		ShiftDiscrepancyTolerance: s.ShiftDiscrepancyTolerance,
		CreatedAt:                 s.CreatedAt,
		UpdatedAt:                 s.UpdatedAt,
	}
	return resp
}
