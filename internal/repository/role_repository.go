package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	FindAllPaginated(ctx context.Context, offset, limit int) ([]entity.Role, int64, error)
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

func (r *roleRepositoryImpl) FindAllPaginated(ctx context.Context, offset, limit int) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
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
