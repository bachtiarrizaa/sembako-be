package model

type CreateProductUnitRequest struct {
	UnitID           string  `json:"unitId" validate:"required,uuid"`
	ConversionToBase float64 `json:"conversionToBase" validate:"required,gt=0"`
	SellingPrice     float64 `json:"sellingPrice" validate:"required,gt=0"`
	IsBaseUnit       bool    `json:"isBaseUnit"`
}

type UpdateProductUnitRequest struct {
	ID               *string `json:"id" validate:"omitempty,uuid"`
	UnitID           string  `json:"unitId" validate:"required,uuid"`
	ConversionToBase float64 `json:"conversionToBase" validate:"required,gt=0"`
	SellingPrice     float64 `json:"sellingPrice" validate:"required,gt=0"`
	IsBaseUnit       bool    `json:"isBaseUnit"`
	IsActive         bool    `json:"isActive"`
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
	DiscountAmount   float64               `json:"discountAmount"`
	DiscountedPrice  float64               `json:"discountedPrice"`
	IsBaseUnit       bool                  `json:"isBaseUnit"`
	IsActive         bool                  `json:"isActive"`
}
