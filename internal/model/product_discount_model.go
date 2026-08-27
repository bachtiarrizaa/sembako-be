package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/shopspring/decimal"
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

type ProductUnitInDiscountResponse struct {
	ID               string                `json:"id"`
	Unit             UnitInProductResponse `json:"unit"`
	ConversionToBase float64               `json:"conversionToBase"`
	SellingPrice     float64               `json:"sellingPrice"`
	DiscountAmount   float64               `json:"discountAmount"`
	DiscountedPrice  float64               `json:"discountedPrice"`
	IsBaseUnit       bool                  `json:"isBaseUnit"`
	IsActive         bool                  `json:"isActive"`
}

func CalculateDiscountPrice(sellingPrice float64, discountType string, discountValue decimal.Decimal) (float64, float64) {
	var discountAmount float64
	val, _ := discountValue.Float64()

	if discountType == "percent" {
		discountAmount = sellingPrice * (val / 100.0)
	} else if discountType == "fixed" {
		discountAmount = val
	}

	if discountAmount > sellingPrice {
		discountAmount = sellingPrice
	}
	if discountAmount < 0 {
		discountAmount = 0
	}

	discountedPrice := sellingPrice - discountAmount
	if discountedPrice < 0 {
		discountedPrice = 0
	}

	return discountAmount, discountedPrice
}

type ProductInDiscountResponse struct {
	ID                string                          `json:"id"`
	ProductDiscountID string                          `json:"productDiscountId,omitempty"`
	Name              string                          `json:"name"`
	Image             *string                         `json:"image"`
	Category          CategoryInProductResponse       `json:"category"`
	IsActive          bool                            `json:"isActive"`
	Units             []ProductUnitInDiscountResponse `json:"units,omitempty"`
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
	var units []ProductUnitInDiscountResponse
	if len(d.Product.Units) > 0 {
		units = make([]ProductUnitInDiscountResponse, 0, len(d.Product.Units))
		for _, u := range d.Product.Units {
			discountAmount, discountedPrice := CalculateDiscountPrice(u.SellingPrice, d.Discount.Type, d.Discount.Value)
			units = append(units, ProductUnitInDiscountResponse{
				ID: u.ID,
				Unit: UnitInProductResponse{
					ID:   u.Unit.ID,
					Name: u.Unit.Name,
				},
				ConversionToBase: u.ConversionToBase,
				SellingPrice:     u.SellingPrice,
				DiscountAmount:   discountAmount,
				DiscountedPrice:  discountedPrice,
				IsBaseUnit:       u.IsBaseUnit,
				IsActive:         u.IsActive,
			})
		}
	}

	return ProductDiscountResponse{
		ID: d.ID,
		Product: ProductInDiscountResponse{
			ID:                d.ProductID,
			ProductDiscountID: d.ID,
			Name:              d.Product.Name,
			Image:             d.Product.Image,
			Category: CategoryInProductResponse{
				ID:   d.Product.CategoryID,
				Name: d.Product.Category.Name,
			},
			IsActive: d.IsActive,
			Units:    units,
		},
		Discount:  ToDiscountResponse(&d.Discount),
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
