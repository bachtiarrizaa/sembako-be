package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DiscountUsecase struct {
	repo repository.DiscountRepository
}

func NewDiscountUsecase(repo repository.DiscountRepository) *DiscountUsecase {
	return &DiscountUsecase{repo: repo}
}

func (u *DiscountUsecase) CreateDiscount(ctx context.Context, req model.CreateDiscountRequest) (*model.DiscountResponse, error) {
	existing, err := u.repo.FindByName(ctx, req.Name)
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

	discount := &entity.Discount{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  true,
	}

	if err := u.repo.Create(ctx, discount); err != nil {
		return nil, errs.NewInternal("failed to create discount")
	}

	return toDiscountResponse(discount), nil
}

func (u *DiscountUsecase) GetDiscountWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.DiscountResponse, utils.Pagination, error) {
	discounts, total, err := u.repo.FindDiscounts(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch discounts")
	}

	res := []model.DiscountResponse{}
	for _, r := range discounts {
		res = append(res, *toDiscountResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *DiscountUsecase) GetDiscountById(ctx context.Context, id string) (*model.DiscountResponse, error) {
	discount, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("discount not found")
	}

	return toDiscountResponse(discount), nil
}

func (u *DiscountUsecase) UpdateDiscount(ctx context.Context, id string, req model.UpdateDiscountRequest) (*model.DiscountResponse, error) {
	discount, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("discount not found")
	}

	if req.Name != nil && *req.Name != discount.Name {
		existing, err := u.repo.FindByName(ctx, *req.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewInternal("failed to check existing discount")
		}
		if existing != nil && existing.ID != id {
			return nil, errs.NewConflict("discount name already exists")
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
		return nil, errs.NewBadRequest("discount value must be greater than 0")
	}
	if discount.Type == "percent" && discount.Value.GreaterThan(decimal.NewFromInt(100)) {
		return nil, errs.NewBadRequest("percent discount cannot exceed 100")
	}
	if discount.StartDate != nil && discount.EndDate != nil {
		if discount.EndDate.Before(*discount.StartDate) {
			return nil, errs.NewBadRequest("end date must not be before start date")
		}
	}

	if err := u.repo.Update(ctx, discount); err != nil {
		return nil, errs.NewInternal("failed to update discount")
	}

	return toDiscountResponse(discount), nil
}

func (u *DiscountUsecase) UpdateStatus(ctx context.Context, id string, req model.UpdateStatusDiscountRequest) (*model.DiscountResponse, error) {
	discount, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("discount not found")
	}

	discount.IsActive = *req.IsActive

	if err := u.repo.Update(ctx, discount); err != nil {
		return nil, errs.NewInternal("failed to update discount status")
	}

	return toDiscountResponse(discount), nil
}

func (u *DiscountUsecase) DeleteDiscount(ctx context.Context, id string) error {
	_, err := u.repo.FindById(ctx, id)
	if err != nil {
		return errs.NewNotFound("discount not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete discount")
	}
	return nil
}

func toDiscountResponse(discount *entity.Discount) *model.DiscountResponse {
	return &model.DiscountResponse{
		ID:        discount.ID,
		Name:      discount.Name,
		Type:      discount.Type,
		Value:     discount.Value,
		StartDate: discount.StartDate,
		EndDate:   discount.EndDate,
		IsActive:  discount.IsActive,
		CreatedAt: discount.CreatedAt,
		UpdatedAt: discount.UpdatedAt,
	}
}
