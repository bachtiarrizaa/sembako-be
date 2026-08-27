package model

import (
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/shopspring/decimal"
)

type GetProductsRequest struct {
	PaginationRequest
	IsActive   *bool  `form:"is_active"`
	CategoryID string `form:"category_id"`
	Include    string `form:"include"`
}

type CreateProductRequest struct {
	CategoryID             string                     `json:"categoryId" form:"categoryId" validate:"required,uuid"`
	Name                   string                     `json:"name" form:"name" validate:"required,min=2,max=150"`
	MinimumStock           *float64                   `json:"minimumStock" form:"minimumStock" validate:"omitempty,gte=0"`
	MarginThresholdPercent *float64                   `json:"marginThresholdPercent" form:"marginThresholdPercent" validate:"omitempty,gte=0,lte=100"`
	Units                  []CreateProductUnitRequest `json:"units" validate:"required,min=1,dive"`
}

type UpdateProductRequest struct {
	CategoryID             string                     `json:"categoryId" form:"categoryId" validate:"required,uuid"`
	Name                   string                     `json:"name" form:"name" validate:"required,min=2,max=150"`
	MinimumStock           *float64                   `json:"minimumStock" form:"minimumStock" validate:"omitempty,gte=0"`
	MarginThresholdPercent *float64                   `json:"marginThresholdPercent" form:"marginThresholdPercent" validate:"omitempty,gte=0,lte=100"`
	Units                  []UpdateProductUnitRequest `json:"units" validate:"required,min=1,dive"`
}

type UpdateProductStatusRequest struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

type CategoryInProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DiscountInProductResponse struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value decimal.Decimal `json:"value"`
}

type ProductResponse struct {
	ID                     string                     `json:"id"`
	Category               CategoryInProductResponse  `json:"category"`
	Name                   string                     `json:"name"`
	Image                  *string                    `json:"image"`
	BaseUnit               UnitInProductResponse      `json:"baseUnit"`
	MinimumStock           *float64                   `json:"minimumStock"`
	MarginThresholdPercent *float64                   `json:"marginThresholdPercent"`
	IsActive               bool                       `json:"isActive"`
	Discount               *DiscountInProductResponse `json:"discount"`
	Units                  []ProductUnitResponse      `json:"units,omitempty"`
	Stock                  float64                    `json:"stock"`
	CreatedAt              time.Time                  `json:"createdAt"`
	UpdatedAt              time.Time                  `json:"updatedAt"`
}

func ToProductResponse(product *entity.Product) ProductResponse {
	var activeDiscount *entity.Discount
	now := time.Now()
	for _, pd := range product.ProductDiscounts {
		if !pd.IsActive {
			continue
		}
		d := pd.Discount
		if !d.IsActive {
			continue
		}
		if d.StartDate != nil && d.StartDate.After(now) {
			continue
		}
		if d.EndDate != nil && d.EndDate.Before(now) {
			continue
		}
		activeDiscount = &d
		break
	}

	var discountResp *DiscountInProductResponse
	if activeDiscount != nil {
		discountResp = &DiscountInProductResponse{
			ID:    activeDiscount.ID,
			Name:  activeDiscount.Name,
			Type:  string(activeDiscount.Type),
			Value: activeDiscount.Value,
		}
	}

	var units []ProductUnitResponse
	if len(product.Units) > 0 {
		units = make([]ProductUnitResponse, 0, len(product.Units))
		for _, pu := range product.Units {
			var discountAmount, discountedPrice float64
			if activeDiscount != nil {
				discountAmount, discountedPrice = CalculateDiscountPrice(pu.SellingPrice, string(activeDiscount.Type), activeDiscount.Value)
			} else {
				discountAmount = 0
				discountedPrice = pu.SellingPrice
			}

			units = append(units, ProductUnitResponse{
				ID: pu.ID,
				Unit: UnitInProductResponse{
					ID:   pu.Unit.ID,
					Name: pu.Unit.Name,
				},
				ConversionToBase: pu.ConversionToBase,
				SellingPrice:     pu.SellingPrice,
				DiscountAmount:   discountAmount,
				DiscountedPrice:  discountedPrice,
				IsBaseUnit:       pu.IsBaseUnit,
				IsActive:         pu.IsActive,
			})
		}
	}

	var stockVal float64 = 0
	if product.Stock != nil {
		stockVal = product.Stock.QtyBaseUnit
	}

	return ProductResponse{
		ID: product.ID,
		Category: CategoryInProductResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
		Name:  product.Name,
		Image: product.Image,
		BaseUnit: UnitInProductResponse{
			ID:   product.BaseUnit.ID,
			Name: product.BaseUnit.Name,
		},
		MinimumStock:           product.MinimumStock,
		MarginThresholdPercent: product.MarginThresholdPercent,
		IsActive:               product.IsActive,
		Discount:               discountResp,
		Units:                  units,
		Stock:                  stockVal,
		CreatedAt:              product.CreatedAt,
		UpdatedAt:              product.UpdatedAt,
	}
}
