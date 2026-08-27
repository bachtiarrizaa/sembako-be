package model

import "github.com/bachtiarrizaa/sembako-be/internal/entity"

type CreateTransactionItem struct {
	ProductUnitID string  `json:"productUnitId" validate:"required,uuid"`
	Qty           float64 `json:"qty" validate:"required,gt=0"`
}

type TransactionItemResponse struct {
	ID              string   `json:"id"`
	ProductUnitID   string   `json:"productUnitId"`
	ProductName     string   `json:"productName"`
	UnitName        string   `json:"unitName"`
	Qty             float64  `json:"qty"`
	UnitPrice       float64  `json:"unitPrice"`
	DiscountApplied float64  `json:"discountApplied"`
	Subtotal        float64  `json:"subtotal"`
	TotalCost       *float64 `json:"totalCost"`
	Margin          *float64 `json:"margin"`
}

func ToTransactionItemResponse(ti *entity.TransactionItem) TransactionItemResponse {
	return TransactionItemResponse{
		ID:              ti.ID,
		ProductUnitID:   ti.ProductUnitID,
		ProductName:     ti.ProductUnit.Product.Name,
		UnitName:        ti.ProductUnit.Unit.Name,
		Qty:             ti.Qty,
		UnitPrice:       ti.UnitPrice,
		DiscountApplied: ti.DiscountApplied,
		Subtotal:        ti.Subtotal,
		TotalCost:       ti.TotalCost,
		Margin:          ti.Margin,
	}
}
