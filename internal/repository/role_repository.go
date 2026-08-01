package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	FindAllPaginated(ctx context.Context, req model.PaginationRequest) ([]entity.Role, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
}

type roleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}

func (r *roleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepositoryImpl) FindAllPaginated(ctx context.Context, req model.PaginationRequest) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Role{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, "id = ?", id).Error
}
