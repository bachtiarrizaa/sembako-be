package model

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type CreateProductInDiscountRequest struct {
	ProductID string `json:"productId" validate:"required,uuid"`
}

type CreateDiscountRequest struct {
	Name      string                           `json:"name" validate:"required,min=2,max=150"`
	Type      string                           `json:"type" validate:"required,oneof=percent fixed"`
	Value     decimal.Decimal                  `json:"value" validate:"required"`
	StartDate *time.Time                       `json:"startDate" validate:"omitempty"`
	EndDate   *time.Time                       `json:"endDate" validate:"omitempty"`
	Products  []CreateProductInDiscountRequest `json:"products" validate:"omitempty,dive"`
}

type UpdateDiscountRequest struct {
	Name      *string                          `json:"name" validate:"omitempty,min=2,max=150"`
	Type      *string                          `json:"type" validate:"omitempty,oneof=percent fixed"`
	Value     *decimal.Decimal                 `json:"value" validate:"omitempty"`
	StartDate Nullable[time.Time]              `json:"startDate"`
	EndDate   Nullable[time.Time]              `json:"endDate"`
	IsActive  *bool                            `json:"isActive"`
	Products  []CreateProductInDiscountRequest `json:"products" validate:"omitempty,dive"`
}

type UpdateStatusDiscountRequest struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

type DiscountResponse struct {
	ID        string                      `json:"id"`
	Name      string                      `json:"name"`
	Type      string                      `json:"type"`
	Value     decimal.Decimal             `json:"value"`
	StartDate *time.Time                  `json:"startDate"`
	EndDate   *time.Time                  `json:"endDate"`
	IsActive  bool                        `json:"isActive"`
	Products  []ProductInDiscountResponse `json:"products"`
	CreatedAt time.Time                   `json:"createdAt"`
	UpdatedAt time.Time                   `json:"updatedAt"`
}

func ToDiscountResponse(d *entity.Discount) DiscountResponse {
	products := []ProductInDiscountResponse{}
	if len(d.ProductDiscounts) > 0 {
		products = make([]ProductInDiscountResponse, 0, len(d.ProductDiscounts))
		for _, pd := range d.ProductDiscounts {
			var units []ProductUnitInDiscountResponse
			if len(pd.Product.Units) > 0 {
				units = make([]ProductUnitInDiscountResponse, 0, len(pd.Product.Units))
				for _, u := range pd.Product.Units {
					discountAmount, discountedPrice := CalculateDiscountPrice(u.SellingPrice, d.Type, d.Value)
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

			products = append(products, ProductInDiscountResponse{
				ID:                pd.ProductID,
				ProductDiscountID: pd.ID,
				Name:              pd.Product.Name,
				Image:             pd.Product.Image,
				Category: CategoryInProductResponse{
					ID:   pd.Product.CategoryID,
					Name: pd.Product.Category.Name,
				},
				IsActive: pd.IsActive,
				Units:    units,
			})
		}
	}

	return DiscountResponse{
		ID:        d.ID,
		Name:      d.Name,
		Type:      d.Type,
		Value:     d.Value,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
		IsActive:  d.IsActive,
		Products:  products,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
