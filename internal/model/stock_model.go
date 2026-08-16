package model

import "time"

type SubmitStockCountRequest struct {
	ProductID   string  `json:"productId" validate:"required,uuid"`
	PhysicalQty float64 `json:"physicalQty" validate:"required,gte=0"`
	Note        string  `json:"note" validate:"omitempty,max=500"`
}

type ApproveStockCountRequest struct {
	Approve bool   `json:"approve"`
	Note    string `json:"note" validate:"omitempty,max=500"`
}

type GetStockCountsRequest struct {
	PaginationRequest
	ProductID string `form:"productId"`
	Status    string `form:"status"` // 'pending', 'approved', 'rejected'
}

type GetStockMutationsRequest struct {
	PaginationRequest
}

type StockCountProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StockCountUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StockCountResponse struct {
	ID          string                     `json:"id"`
	ProductID   string                     `json:"productId"`
	Product     StockCountProductResponse  `json:"product"`
	CountDate   time.Time                  `json:"countDate"`
	SystemQty   float64                    `json:"systemQty"`
	PhysicalQty float64                    `json:"physicalQty"`
	Discrepancy float64                    `json:"discrepancy"`
	Note        *string                    `json:"note"`
	Status      string                     `json:"status"`
	SubmittedBy string                     `json:"submittedBy"`
	Submitter   StockCountUserResponse     `json:"submitter"`
	ApprovedBy  *string                    `json:"approvedBy"`
	Approver    *StockCountUserResponse    `json:"approver,omitempty"`
	SubmittedAt time.Time                  `json:"submittedAt"`
	ApprovedAt  *time.Time                 `json:"approvedAt,omitempty"`
}

type StockMutationUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StockMutationResponse struct {
	ID          string                     `json:"id"`
	ProductID   string                     `json:"productId"`
	Type        string                     `json:"type"` // 'in' or 'out'
	Qty         float64                    `json:"qty"`
	QtyBefore   float64                    `json:"qtyBefore"`
	QtyAfter    float64                    `json:"qtyAfter"`
	Source      string                     `json:"source"`
	ReferenceID *string                    `json:"referenceId,omitempty"`
	Note        *string                    `json:"note,omitempty"`
	Creator     StockMutationUserResponse  `json:"creator"`
	CreatedAt   time.Time                  `json:"createdAt"`
}

type StockBaseUnitResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StockSummaryResponse struct {
	ProductID   string                `json:"productId"`
	QtyBaseUnit float64               `json:"qtyBaseUnit"`
	BaseUnit    StockBaseUnitResponse `json:"baseUnit"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}
