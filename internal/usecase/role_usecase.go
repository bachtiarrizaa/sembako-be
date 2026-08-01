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

type RoleUsecase struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) *RoleUsecase {
	return &RoleUsecase{repo: repo}
}

func (u *RoleUsecase) CreateRole(ctx context.Context, req model.CreateRoleRequest) (*model.RoleResponse, error) {
	existing, err := u.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing role")
	}
	if existing != nil {
		return nil, errs.NewConflict("role name already exists")
	}

	role := &entity.Role{Name: req.Name}
	if err := u.repo.Create(ctx, role); err != nil {
		return nil, errs.NewInternal("failed to create role")
	}
	return toRoleResponse(role), nil
}

func (u *RoleUsecase) GetAllRoles(ctx context.Context, req model.PaginationRequest) ([]model.RoleResponse, utils.Pagination, error) {
	roles, total, err := u.repo.FindAllPaginated(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch roles")
	}

	var res []model.RoleResponse
	for _, r := range roles {
		res = append(res, *toRoleResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *RoleUsecase) GetRoleByID(ctx context.Context, id string) (*model.RoleResponse, error) {
	role, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("role not found")
	}
	return toRoleResponse(role), nil
}

func (u *RoleUsecase) UpdateRole(ctx context.Context, id string, req model.UpdateRoleRequest) (*model.RoleResponse, error) {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("role not found")
	}

	existing, err := u.repo.FindByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != id {
		return nil, errs.NewConflict("role name already exists")
	}

	role := &entity.Role{
		ID:   id,
		Name: req.Name,
	}
	if err := u.repo.Update(ctx, role); err != nil {
		return nil, errs.NewInternal("failed to update role")
	}
	return toRoleResponse(role), nil
}

func (u *RoleUsecase) DeleteRole(ctx context.Context, id string) error {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return errs.NewNotFound("role not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete role")
	}
	return nil
}

func toRoleResponse(role *entity.Role) *model.RoleResponse {
	return &model.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
