package model

import "time"

type GetProductsRequest struct {
	PaginationRequest
	IsActive   *bool  `form:"is_active"`
	CategoryID string `form:"category_id"`
}

type CreateProductUnitRequest struct {
	UnitID           string  `json:"unitId" validate:"required,uuid"`
	ConversionToBase float64 `json:"conversionToBase" validate:"required,gt=0"`
	SellingPrice     float64 `json:"sellingPrice" validate:"required,gt=0"`
	IsBaseUnit       bool    `json:"isBaseUnit"`
}

type CreateProductRequest struct {
	CategoryID             string                     `json:"categoryId" form:"categoryId" validate:"required,uuid"`
	Name                   string                     `json:"name" form:"name" validate:"required,min=2,max=150"`
	MinimumStock           *float64                   `json:"minimumStock" form:"minimumStock" validate:"omitempty,gte=0"`
	MarginThresholdPercent *float64                   `json:"marginThresholdPercent" form:"marginThresholdPercent" validate:"omitempty,gte=0,lte=100"`
	Units                  []CreateProductUnitRequest `json:"units" validate:"required,min=1,dive"`
}

type UpdateProductUnitRequest struct {
	ConversionToBase float64 `json:"conversionToBase" validate:"required,gt=0"`
	SellingPrice     float64 `json:"sellingPrice" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	CategoryID             string   `json:"categoryId" validate:"required,uuid"`
	Name                   string   `json:"name" validate:"required,min=2,max=150"`
	MinimumStock           *float64 `json:"minimumStock" validate:"omitempty,gte=0"`
	MarginThresholdPercent *float64 `json:"marginThresholdPercent" validate:"omitempty,gte=0,lte=100"`
}

type AddProductUnitRequest struct {
	UnitID           string  `json:"unitId" validate:"required,uuid"`
	ConversionToBase float64 `json:"conversionToBase" validate:"required,gt=0"`
	SellingPrice     float64 `json:"sellingPrice" validate:"required,gt=0"`
}

type UnitInProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductUnitResponse struct {
	ID               string                `json:"id"`
	Unit             UnitInProductResponse `json:"unit"`
	ConversionToBase float64               `json:"conversionToBase"`
	SellingPrice     float64               `json:"sellingPrice"`
	IsBaseUnit       bool                  `json:"isBaseUnit"`
	IsActive         bool                  `json:"isActive"`
}

type CategoryInProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID                     string                    `json:"id"`
	Category               CategoryInProductResponse `json:"category"`
	Name                   string                    `json:"name"`
	Image                  *string                   `json:"image"`
	BaseUnit               UnitInProductResponse     `json:"baseUnit"`
	MinimumStock           *float64                  `json:"minimumStock"`
	MarginThresholdPercent *float64                  `json:"marginThresholdPercent"`
	IsActive               bool                      `json:"isActive"`
	Units                  []ProductUnitResponse     `json:"units"`
	CreatedAt              time.Time                 `json:"createdAt"`
	UpdatedAt              time.Time                 `json:"updatedAt"`
}
