package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type SupplierUsecase struct {
	repo repository.SupplierRepository
}

func NewSupplierUsecase(repo repository.SupplierRepository) *SupplierUsecase {
	return &SupplierUsecase{repo: repo}
}

func (u *SupplierUsecase) Create(ctx context.Context, req model.CreateSupplierRequest) (*model.SupplierResponse, error) {
	existing, err := u.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing supplier")
	}
	if existing != nil {
		return nil, errs.NewConflict("supplier name already exists")
	}

	supplier := &entity.Supplier{
		Name:        req.Name,
		ContactName: req.ContactName,
		Phone:       req.Phone,
		Address:     req.Address,
		IsActive:    true,
	}

	if err := u.repo.Create(ctx, supplier); err != nil {
		return nil, errs.NewInternal("failed to create supplier")
	}

	return toSupplierResponse(supplier), nil
}

func (u *SupplierUsecase) GetSuppliers(ctx context.Context, req model.PaginationRequest) ([]model.SupplierResponse, utils.Pagination, error) {
	suppliers, total, err := u.repo.FindSuppliers(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch suppliers")
	}

	var res []model.SupplierResponse
	for _, r := range suppliers {
		res = append(res, *toSupplierResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *SupplierUsecase) GetSupplierById(ctx context.Context, id string) (*model.SupplierResponse, error) {
	category, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("category not found")
	}
	return toSupplierResponse(category), nil
}

func (u *SupplierUsecase) Update(ctx context.Context, id string, req model.UpdateSupplierRequest) (*model.SupplierResponse, error) {
	supplier, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("supplier not found")
	}

	existing, err := u.repo.FindByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != id {
		return nil, errs.NewConflict("supplier name already exists")
	}

	supplier.Name = req.Name
	supplier.ContactName = req.ContactName
	supplier.Phone = req.Phone
	supplier.Address = req.Address

	if err := u.repo.Update(ctx, supplier); err != nil {
		return nil, errs.NewInternal("failed to update supplier")
	}
	return toSupplierResponse(supplier), nil
}

func (u *SupplierUsecase) UpdateStatus(ctx context.Context, id string, req model.UpdateStatusSupplierRequest) (*model.SupplierResponse, error) {
	supplier, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("supplier not found")
	}

	supplier.IsActive = req.IsActive

	if err := u.repo.Update(ctx, supplier); err != nil {
		return nil, errs.NewInternal("failed to update supplier status")
	}

	return toSupplierResponse(supplier), nil
}

func (u *SupplierUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.FindById(ctx, id)
	if err != nil {
		return errs.NewNotFound("supplier not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete supplier")
	}
	return nil
}

func toSupplierResponse(supplier *entity.Supplier) *model.SupplierResponse {
	return &model.SupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		ContactName: supplier.ContactName,
		Phone:       supplier.Phone,
		Address:     supplier.Address,
		IsActive:    supplier.IsActive,
		CreatedAt:   supplier.CreatedAt,
		UpdatedAt:   supplier.UpdatedAt,
	}
}
