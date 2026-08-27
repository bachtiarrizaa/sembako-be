package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DiscountUsecase struct {
	db                  *gorm.DB
	discountRepo        repository.DiscountRepository
	productDiscountRepo repository.ProductDiscountRepository
}

func NewDiscountUsecase(
	db *gorm.DB,
	discountRepo repository.DiscountRepository,
	productDiscountRepo repository.ProductDiscountRepository,
) *DiscountUsecase {
	return &DiscountUsecase{
		db:                  db,
		discountRepo:        discountRepo,
		productDiscountRepo: productDiscountRepo,
	}
}

func (u *DiscountUsecase) CreateDiscount(ctx context.Context, req model.CreateDiscountRequest) (*model.DiscountResponse, error) {
	existing, err := u.discountRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing discount")
	}
	if existing != nil {
		return nil, errs.NewConflict("discount name already exists")
	}

	if req.Value.LessThanOrEqual(decimal.Zero) {
		return nil, errs.NewBadRequest("discount value must be greater than 0")
	}

	if req.Type == "percent" && req.Value.GreaterThan(decimal.NewFromInt(100)) {
		return nil, errs.NewBadRequest("percent discount cannot exceed 100")
	}

	if req.StartDate != nil && req.EndDate != nil {
		if req.EndDate.Before(*req.StartDate) {
			return nil, errs.NewBadRequest("end date must not be before start date")
		}
	}

	if len(req.Products) > 0 {
		seenProductIDs := make(map[string]bool)
		for _, p := range req.Products {
			if seenProductIDs[p.ProductID] {
				return nil, errs.NewBadRequest("duplicate product in products list")
			}
			seenProductIDs[p.ProductID] = true
		}
	}

	var discount *entity.Discount

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		discountRepoTx := u.discountRepo.WithTx(tx)
		productDiscountRepoTx := u.productDiscountRepo.WithTx(tx)

		discount = &entity.Discount{
			Name:      req.Name,
			Type:      req.Type,
			Value:     req.Value,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			IsActive:  true,
		}

		if err := discountRepoTx.Create(ctx, discount); err != nil {
			return err
		}

		if len(req.Products) > 0 {
			productDiscounts := make([]entity.ProductDiscount, 0, len(req.Products))
			for _, p := range req.Products {
				productDiscounts = append(productDiscounts, entity.ProductDiscount{
					DiscountID: discount.ID,
					ProductID:  p.ProductID,
					IsActive:   true,
				})
			}

			if err := productDiscountRepoTx.CreateMany(ctx, productDiscounts); err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(txErr, &pgErr) {
			if pgErr.Code == "23503" {
				return nil, errs.NewBadRequest("one or more products do not exist")
			}
		}
		if appErr, ok := txErr.(*errs.AppError); ok {
			return nil, appErr
		}
		return nil, errs.NewInternal("failed to create discount: " + txErr.Error())
	}

	created, err := u.discountRepo.FindById(ctx, discount.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load created discount")
	}

	resp := model.ToDiscountResponse(created)
	return &resp, nil
}

func (u *DiscountUsecase) GetDiscountWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.DiscountResponse, utils.Pagination, error) {
	discounts, total, err := u.discountRepo.FindDiscounts(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch discounts")
	}

	res := make([]model.DiscountResponse, 0, len(discounts))
	for _, r := range discounts {
		res = append(res, model.ToDiscountResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (u *DiscountUsecase) GetDiscountById(ctx context.Context, id string) (*model.DiscountResponse, error) {
	discount, err := u.discountRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("discount not found")
		}
		return nil, errs.NewInternal("failed to fetch discount: " + err.Error())
	}

	resp := model.ToDiscountResponse(discount)
	return &resp, nil
}

func (u *DiscountUsecase) UpdateDiscount(ctx context.Context, id string, req model.UpdateDiscountRequest) (*model.DiscountResponse, error) {
	if req.Value != nil && req.Value.LessThanOrEqual(decimal.Zero) {
		return nil, errs.NewBadRequest("discount value must be greater than 0")
	}
	if req.Type != nil && *req.Type == "percent" && req.Value != nil && req.Value.GreaterThan(decimal.NewFromInt(100)) {
		return nil, errs.NewBadRequest("percent discount cannot exceed 100")
	}
	if req.StartDate.IsSet && req.EndDate.IsSet && req.StartDate.Value != nil && req.EndDate.Value != nil {
		if req.EndDate.Value.Before(*req.StartDate.Value) {
			return nil, errs.NewBadRequest("end date must not be before start date")
		}
	}

	if req.Products != nil {
		seenProductIDs := make(map[string]bool)
		for _, p := range req.Products {
			if seenProductIDs[p.ProductID] {
				return nil, errs.NewBadRequest("duplicate product in products list")
			}
			seenProductIDs[p.ProductID] = true
		}
	}

	var discount *entity.Discount

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		discountRepoTx := u.discountRepo.WithTx(tx)
		productDiscountRepoTx := u.productDiscountRepo.WithTx(tx)

		var err error
		discount, err = discountRepoTx.FindById(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.NewNotFound("discount not found")
			}
			return err
		}

		if req.Name != nil && *req.Name != discount.Name {
			existing, err := discountRepoTx.FindByName(ctx, *req.Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if existing != nil && existing.ID != id {
				return errs.NewConflict("discount name already exists")
			}
			discount.Name = *req.Name
		}

		if req.Type != nil {
			discount.Type = *req.Type
		}
		if req.Value != nil {
			discount.Value = *req.Value
		}
		if req.IsActive != nil {
			discount.IsActive = *req.IsActive
		}

		if req.StartDate.IsSet {
			discount.StartDate = req.StartDate.Value
		}
		if req.EndDate.IsSet {
			discount.EndDate = req.EndDate.Value
		}

		if discount.Value.LessThanOrEqual(decimal.Zero) {
			return errs.NewBadRequest("discount value must be greater than 0")
		}
		if discount.Type == "percent" && discount.Value.GreaterThan(decimal.NewFromInt(100)) {
			return errs.NewBadRequest("percent discount cannot exceed 100")
		}
		if discount.StartDate != nil && discount.EndDate != nil {
			if discount.EndDate.Before(*discount.StartDate) {
				return errs.NewBadRequest("end date must not be before start date")
			}
		}

		if err := discountRepoTx.Update(ctx, discount); err != nil {
			return err
		}

		if req.Products != nil {
			if err := productDiscountRepoTx.DeleteByDiscountID(ctx, discount.ID); err != nil {
				return err
			}

			if len(req.Products) > 0 {
				productDiscounts := make([]entity.ProductDiscount, 0, len(req.Products))
				for _, p := range req.Products {
					productDiscounts = append(productDiscounts, entity.ProductDiscount{
						DiscountID: discount.ID,
						ProductID:  p.ProductID,
						IsActive:   true,
					})
				}

				if err := productDiscountRepoTx.CreateMany(ctx, productDiscounts); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if txErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(txErr, &pgErr) {
			if pgErr.Code == "23503" {
				return nil, errs.NewBadRequest("one or more products do not exist")
			}
		}
		if appErr, ok := txErr.(*errs.AppError); ok {
			return nil, appErr
		}
		return nil, errs.NewInternal("failed to update discount: " + txErr.Error())
	}

	updated, err := u.discountRepo.FindById(ctx, discount.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated discount")
	}

	resp := model.ToDiscountResponse(updated)
	return &resp, nil
}

func (u *DiscountUsecase) UpdateStatus(ctx context.Context, id string, req model.UpdateStatusDiscountRequest) (*model.DiscountResponse, error) {
	discount, err := u.discountRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("discount not found")
		}
		return nil, errs.NewInternal("failed to fetch discount: " + err.Error())
	}

	discount.IsActive = *req.IsActive

	if err := u.discountRepo.Update(ctx, discount); err != nil {
		return nil, errs.NewInternal("failed to update discount status")
	}

	resp := model.ToDiscountResponse(discount)
	return &resp, nil
}

func (u *DiscountUsecase) DeleteDiscount(ctx context.Context, id string) error {
	_, err := u.discountRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("discount not found")
		}
		return errs.NewInternal("failed to fetch discount: " + err.Error())
	}

	if err := u.discountRepo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete discount")
	}
	return nil
}

