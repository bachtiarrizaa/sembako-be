package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]entity.Permission, error)
	FindByID(ctx context.Context, id string) (*entity.Permission, error)
}

type permissionRepositoryImpl struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepositoryImpl{db: db}
}

func (r *permissionRepositoryImpl) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
