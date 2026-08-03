package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"gorm.io/gorm"
)

func (u *ProductUsecase) UpdateProductUnitStatus(ctx context.Context, productID string, unitID string) (*model.ProductResponse, error) {
	productUnit, err := u.productUnitRepo.FindByID(ctx, unitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product unit not found")
		}
		return nil, errs.NewInternal("failed to fetch product unit")
	}

	if productUnit.ProductID != productID {
		return nil, errs.NewNotFound("product unit not found for this product")
	}

	if productUnit.IsBaseUnit && productUnit.IsActive {
		return nil, errs.NewBadRequest("cannot deactivate base unit")
	}

	productUnit.IsActive = !productUnit.IsActive

	if err := u.productUnitRepo.Update(ctx, productUnit); err != nil {
		return nil, errs.NewInternal("failed to update product unit status")
	}

	updatedProduct, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated product")
	}

	return toProductResponse(updatedProduct), nil
}

func (u *ProductUsecase) AddProductUnit(ctx context.Context, productID string, req model.AddProductUnitRequest) (*model.ProductResponse, error) {
	product, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product not found")
		}
		return nil, errs.NewInternal("failed to fetch product")
	}

	for _, pu := range product.Units {
		if pu.UnitID == req.UnitID {
			return nil, errs.NewConflict("unit already exists for this product")
		}
	}

	newProductUnit := entity.ProductUnit{
		ProductID:        product.ID,
		UnitID:           req.UnitID,
		ConversionToBase: req.ConversionToBase,
		SellingPrice:     req.SellingPrice,
		IsBaseUnit:       false,
		IsActive:         true,
	}

	if err := u.productUnitRepo.Create(ctx, &newProductUnit); err != nil {
		return nil, errs.NewInternal("failed to add product unit")
	}

	updated, err := u.productRepo.FindById(ctx, product.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated product")
	}

	return toProductResponse(updated), nil
}

func (u *ProductUsecase) UpdateProductUnit(ctx context.Context, productID string, unitID string, req model.UpdateProductUnitRequest) (*model.ProductResponse, error) {
	_, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product not found")
		}
		return nil, errs.NewInternal("failed to fetch product")
	}

	productUnit, err := u.productUnitRepo.FindByID(ctx, unitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product unit not found")
		}
		return nil, errs.NewInternal("failed to fetch product unit")
	}

	if productUnit.ProductID != productID {
		return nil, errs.NewNotFound("product unit not found for this product")
	}

	if productUnit.IsBaseUnit && req.ConversionToBase != 1 {
		return nil, errs.NewBadRequest("base unit must have conversionToBase = 1")
	}

	productUnit.ConversionToBase = req.ConversionToBase
	productUnit.SellingPrice = req.SellingPrice

	if err := u.productUnitRepo.Update(ctx, productUnit); err != nil {
		return nil, errs.NewInternal("failed to update product unit")
	}

	updated, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated product")
	}

	return toProductResponse(updated), nil
}

func (u *ProductUsecase) DeleteProductUnit(ctx context.Context, productID string, unitID string) error {
	_, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("product not found")
		}
		return errs.NewInternal("failed to fetch product")
	}

	productUnit, err := u.productUnitRepo.FindByID(ctx, unitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("product unit not found")
		}
		return errs.NewInternal("failed to fetch product unit")
	}

	if productUnit.ProductID != productID {
		return errs.NewNotFound("product unit not found for this product")
	}

	if productUnit.IsBaseUnit {
		return errs.NewBadRequest("cannot delete base unit")
	}

	if u.db.Migrator().HasTable("transaction_items") {
		var count int64
		if err := u.db.Table("transaction_items").Where("product_unit_id = ?", unitID).Count(&count).Error; err != nil {
			return errs.NewInternal("failed to check transaction history")
		}
		if count > 0 {
			return errs.NewConflict("cannot delete product unit that has been used in transactions")
		}
	}

	if err := u.productUnitRepo.Delete(ctx, unitID); err != nil {
		return errs.NewInternal("failed to delete product unit")
	}

	return nil
}
