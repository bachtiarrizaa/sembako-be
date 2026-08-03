package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type Unitusecase struct {
	repo repository.UnitRepository
}

func NewUnitUsecase(repo repository.UnitRepository) *Unitusecase {
	return &Unitusecase{repo: repo}
}

func (u *Unitusecase) CreateUnit(ctx context.Context, req model.CreateUnitRequest) (*model.UnitResponse, error) {
	existing, err := u.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing unit")
	}
	if existing != nil {
		return nil, errs.NewConflict("unit name already exists")
	}

	unit := &entity.Unit{Name: req.Name}
	if err := u.repo.Create(ctx, unit); err != nil {
		return nil, errs.NewInternal("failed to create unit")
	}
	return toUnitResponse(unit), nil
}

func (u *Unitusecase) GetUnitsWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.UnitResponse, utils.Pagination, error) {
	units, total, err := u.repo.FindUnits(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch untis")
	}

	var res []model.UnitResponse
	for _, r := range units {
		res = append(res, *toUnitResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *Unitusecase) GetUnitByID(ctx context.Context, id string) (*model.UnitResponse, error) {
	unit, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("unit not found")
	}
	return toUnitResponse(unit), nil
}

func (u *Unitusecase) UpdateUnit(ctx context.Context, id string, req model.UpdateUnitRequest) (*model.UnitResponse, error) {
	unit, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("unit not found")
	}

	existing, err := u.repo.FindByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != id {
		return nil, errs.NewConflict("unit name already exists")
	}

	unit.Name = req.Name
	if err := u.repo.Update(ctx, unit); err != nil {
		return nil, errs.NewInternal("failed to update unit")
	}
	return toUnitResponse(unit), nil
}

func (u *Unitusecase) DeleteUnit(ctx context.Context, id string) error {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return errs.NewNotFound("unit not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete unit")
	}
	return nil
}

func toUnitResponse(unit *entity.Unit) *model.UnitResponse {
	return &model.UnitResponse{
		ID:        unit.ID,
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}
}
