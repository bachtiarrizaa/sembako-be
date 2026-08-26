package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type CreateProductDiscountRequest struct {
	ProductID  string `json:"productId" validate:"required,uuid"`
	DiscountID string `json:"discountId" validate:"required,uuid"`
}

type UpdateProductDiscountRequest struct {
	ProductID  string `json:"productId" validate:"required,uuid"`
	DiscountID string `json:"discountId" validate:"required,uuid"`
}

type UpdateProductDiscountStatusRequest struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

type GetProductDiscountsRequest struct {
	PaginationRequest
	DiscountID string `form:"discountId"`
	ProductID  string `form:"productId"`
	IsActive   *bool  `form:"isActive"`
}

type ProductInDiscountResponse struct {
	ID       string                    `json:"id"`
	Name     string                    `json:"name"`
	Category CategoryInProductResponse `json:"category"`
}

type ProductDiscountResponse struct {
	ID        string                    `json:"id"`
	Product   ProductInDiscountResponse `json:"product"`
	Discount  DiscountResponse          `json:"discount"`
	IsActive  bool                      `json:"isActive"`
	CreatedAt time.Time                 `json:"createdAt"`
	UpdatedAt time.Time                 `json:"updatedAt"`
}

func ToProductDiscountResponse(d *entity.ProductDiscount) ProductDiscountResponse {
	return ProductDiscountResponse{
		ID: d.ID,
		Product: ProductInDiscountResponse{
			ID:   d.ProductID,
			Name: d.Product.Name,
			Category: CategoryInProductResponse{
				ID:   d.Product.CategoryID,
				Name: d.Product.Category.Name,
			},
		},
		Discount:  ToDiscountResponse(&d.Discount),
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
