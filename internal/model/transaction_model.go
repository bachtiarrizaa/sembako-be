package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type VoidTransactionRequest struct {
	Reason string `json:"reason" validate:"required,min=5,max=255"`
}

type CreateTransactionRequest struct {
	CustomerID    *string                 `json:"customerId" validate:"omitempty,uuid"`
	PaymentMethod string                  `json:"paymentMethod" validate:"required,oneof=cash qris transfer"`
	CashReceived  *float64                `json:"cashReceived" validate:"omitempty,gte=0"`
	UsePoints     *bool                   `json:"usePoints"`
	Items         []CreateTransactionItem `json:"items" validate:"required,dive,required"`
}

type ListTransactionsRequest struct {
	PaginationRequest
	StartDate     *string `form:"startDate"`
	EndDate       *string `form:"endDate"`
	CashierID     *string `form:"cashierId"`
	CustomerID    *string `form:"customerId"`
	PaymentMethod *string `form:"paymentMethod"`
	Status        *string `form:"status"`
}

type TransactionCashierResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransactionCustomerResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransactionResponse struct {
	ID                     string                       `json:"id"`
	ReceiptNumber          string                       `json:"receiptNumber"`
	Cashier                TransactionCashierResponse   `json:"cashier"`
	ShiftID                string                       `json:"shiftId"`
	Customer               *TransactionCustomerResponse `json:"customer"`
	PaymentMethod          string                       `json:"paymentMethod"`
	Subtotal               float64                      `json:"subtotal"`
	TotalDiscount          float64                      `json:"totalDiscount"`
	PointsUsed             int                          `json:"pointsUsed"`
	PointsDiscountValue    float64                      `json:"pointsDiscountValue"`
	PointsEarned           int                          `json:"pointsEarned"`
	Total                  float64                      `json:"total"`
	CashReceived           *float64                     `json:"cashReceived"`
	ChangeGiven            *float64                     `json:"changeGiven"`
	ManualPaidConfirmation *bool                        `json:"manualPaidConfirmation"`
	Status                 string                       `json:"status"`
	VoidReason             *string                      `json:"voidReason"`
	VoidedByUser           *TransactionCashierResponse  `json:"voidedByUser"`
	VoidedAt               *time.Time                   `json:"voidedAt"`
	CreatedAt              time.Time                    `json:"createdAt"`
	UpdatedAt              time.Time                    `json:"updatedAt"`
	Items                  []TransactionItemResponse    `json:"items"`
}

func ToTransactionResponse(t *entity.Transaction) TransactionResponse {
	resp := TransactionResponse{
		ID:            t.ID,
		ReceiptNumber: t.ReceiptNumber,
		Cashier: TransactionCashierResponse{
			ID:   t.CashierID,
			Name: t.Cashier.Name,
		},
		ShiftID:                t.ShiftID,
		PaymentMethod:          t.PaymentMethod,
		Subtotal:               t.Subtotal,
		TotalDiscount:          t.TotalDiscount,
		PointsUsed:             t.PointsUsed,
		PointsDiscountValue:    t.PointsDiscountValue,
		PointsEarned:           t.PointsEarned,
		Total:                  t.Total,
		CashReceived:           t.CashReceived,
		ChangeGiven:            t.ChangeGiven,
		ManualPaidConfirmation: t.ManualPaidConfirmation,
		Status:                 t.Status,
		VoidReason:             t.VoidReason,
		VoidedAt:               t.VoidedAt,
		CreatedAt:              t.CreatedAt,
		UpdatedAt:              t.UpdatedAt,
	}

	if t.Customer != nil {
		resp.Customer = &TransactionCustomerResponse{
			ID:   t.Customer.ID,
			Name: t.Customer.Name,
		}
	}

	if t.VoidedByUser != nil {
		resp.VoidedByUser = &TransactionCashierResponse{
			ID:   t.VoidedByUser.ID,
			Name: t.VoidedByUser.Name,
		}
	}

	if len(t.Items) > 0 {
		resp.Items = make([]TransactionItemResponse, 0, len(t.Items))
		for i := range t.Items {
			resp.Items = append(resp.Items, ToTransactionItemResponse(&t.Items[i]))
		}
	}

	return resp
}
