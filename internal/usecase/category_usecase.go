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

type CategoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (u *CategoryUsecase) CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	existing, err := u.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing category")
	}
	if existing != nil {
		return nil, errs.NewConflict("category name already exists")
	}

	category := &entity.Category{Name: req.Name}
	if err := u.repo.Create(ctx, category); err != nil {
		return nil, errs.NewInternal("failed to create category")
	}

	return toCategoryResponse(category), nil
}

func (u *CategoryUsecase) GetCategoriesWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.CategoryResponse, utils.Pagination, error) {
	categories, total, err := u.repo.FindCategories(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch categories")
	}

	res := []model.CategoryResponse{}
	for _, r := range categories {
		res = append(res, *toCategoryResponse(&r))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *CategoryUsecase) GetCategoryById(ctx context.Context, id string) (*model.CategoryResponse, error) {
	category, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("category not found")
	}
	return toCategoryResponse(category), nil
}

func (u *CategoryUsecase) UpdateCategory(ctx context.Context, id string, req model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	category, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("category not found")
	}

	existing, err := u.repo.FindByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != id {
		return nil, errs.NewConflict("category name already exist")
	}

	category.Name = req.Name
	if err := u.repo.Update(ctx, category); err != nil {
		return nil, errs.NewInternal("failed to update category")
	}
	return toCategoryResponse(category), nil
}

func (u *CategoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	_, err := u.repo.FindById(ctx, id)
	if err != nil {
		return errs.NewNotFound("category not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete category")
	}
	return nil
}

func toCategoryResponse(category *entity.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
