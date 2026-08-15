package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type CreatePurchaseRequest struct {
	SupplierID    string               `json:"supplierId" binding:"required,uuid"`
	InvoiceNumber *string              `json:"invoiceNumber"`
	PurchaseDate  string               `json:"purchaseDate" binding:"required"`
	Items         []CreatePurchaseItem `json:"items" binding:"required,min=1,dive"`
}

type CreatePurchaseItem struct {
	ProductID     string  `json:"productId" binding:"required,uuid"`
	UnitID        string  `json:"unitId" binding:"required,uuid"`
	Quantity      float64 `json:"quantity" binding:"required,gt=0"`
	PurchasePrice float64 `json:"purchasePrice" binding:"required,gte=0"`
}

type UpdatePurchaseRequest struct {
	SupplierID    string  `json:"supplierId" binding:"required,uuid"`
	InvoiceNumber *string `json:"invoiceNumber"`
	PurchaseDate  string  `json:"purchaseDate" binding:"required"`
	Quantity      float64 `json:"quantity" binding:"required,gt=0"`
	UnitID        string  `json:"unitId" binding:"required,uuid"`
	PurchasePrice float64 `json:"purchasePrice" binding:"required,gte=0"`
}

type PurchaseProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PurchaseSupplierResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PurchaseCreatorResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PurchaseBatchResponse struct {
	ID                string                   `json:"id"`
	Product           PurchaseProductResponse  `json:"product"`
	Supplier          PurchaseSupplierResponse `json:"supplier"`
	InitialQuantity   float64                  `json:"initialQuantity"`
	RemainingQuantity float64                  `json:"remainingQuantity"`
	PurchasePrice     float64                  `json:"purchasePrice"`
	InvoiceNumber     *string                  `json:"invoiceNumber"`
	PurchaseDate      time.Time                `json:"purchaseDate"`
	Creator           PurchaseCreatorResponse  `json:"creator"`
	CreatedAt         time.Time                `json:"createdAt"`
}

type GetPurchaseBatchesRequest struct {
	Search     string `form:"search"`
	SupplierID string `form:"supplierId"`
	ProductID  string `form:"productId"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}

func ToPurchaseBatchResponse(b *entity.PurchaseBatch) PurchaseBatchResponse {
	return PurchaseBatchResponse{
		ID: b.ID,
		Product: PurchaseProductResponse{
			ID:   b.ProductID,
			Name: b.Product.Name,
		},
		Supplier: PurchaseSupplierResponse{
			ID:   b.SupplierID,
			Name: b.Supplier.Name,
		},
		InitialQuantity:   b.InitialQty,
		RemainingQuantity: b.RemainingQty,
		PurchasePrice:     b.PurchasePrice,
		InvoiceNumber:     b.InvoiceNumber,
		PurchaseDate:      b.PurchaseDate,
		Creator: PurchaseCreatorResponse{
			ID:   b.CreatedBy,
			Name: b.Creator.Name,
		},
		CreatedAt: b.CreatedAt,
	}
}
