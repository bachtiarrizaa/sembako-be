package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type ProductDiscountUsecase struct {
	productDiscountRepo repository.ProductDiscountRepository
}

func NewProductDiscountUsecase(repo repository.ProductDiscountRepository) *ProductDiscountUsecase {
	return &ProductDiscountUsecase{productDiscountRepo: repo}
}

func (u *ProductDiscountUsecase) Create(ctx context.Context, req model.CreateProductDiscountRequest) (*model.ProductDiscountResponse, error) {
	existing, err := u.productDiscountRepo.FindByDiscountAndProduct(ctx, req.DiscountID, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing product discount")
	}
	if existing != nil {
		return nil, errs.NewConflict("this discount is already applied to the product")
	}

	productDiscount := &entity.ProductDiscount{
		ProductID:  req.ProductID,
		DiscountID: req.DiscountID,
		IsActive:   true,
	}

	if err := u.productDiscountRepo.Create(ctx, productDiscount); err != nil {
		return nil, errs.NewInternal("failed to create product discount")
	}

	created, err := u.productDiscountRepo.FindByID(ctx, productDiscount.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load created product discount: " + err.Error())
	}

	resp := model.ToProductDiscountResponse(created)
	return &resp, nil
}
