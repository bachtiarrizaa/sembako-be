package usecase

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type PermissionUsecase struct {
	db             *gorm.DB
	permissionRepo repository.PermissionRepository
	userRepo       repository.UserRepository
}

func NewPermissionUsecase(
	db *gorm.DB,
	permissionRepo repository.PermissionRepository,
	userRepo repository.UserRepository,
) *PermissionUsecase {
	return &PermissionUsecase{
		db:             db,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (u *PermissionUsecase) CheckUserPermission(ctx context.Context, userIDStr string, permissionName string) (bool, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return false, err
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	rolePermissions, err := u.permissionRepo.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return false, err
	}

	permMap := make(map[string]bool)
	for _, p := range rolePermissions {
		permMap[p.Name] = true
	}

	if permMap[permissionName] {
		return true, nil
	}

	return false, nil
}

func (u *PermissionUsecase) GetUserMenu(ctx context.Context, userIDStr string) ([]model.PermissionTreeResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rolePermissions, err := u.permissionRepo.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	return buildMenuTree(rolePermissions), nil
}

func buildMenuTree(menus []entity.Permission) []model.PermissionTreeResponse {
	menuMap := make(map[string]entity.Permission)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	var roots []entity.Permission
	for _, m := range menus {
		if m.ParentID == nil {
			roots = append(roots, m)
		} else {
			_, parentExists := menuMap[*m.ParentID]
			if !parentExists {
				roots = append(roots, m)
			}
		}
	}

	var buildNode func(node entity.Permission) model.PermissionTreeResponse
	buildNode = func(node entity.Permission) model.PermissionTreeResponse {
		var children []model.PermissionTreeResponse
		for _, m := range menus {
			if m.ParentID != nil && *m.ParentID == node.ID {
				children = append(children, buildNode(m))
			}
		}
		return model.PermissionTreeResponse{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
			ParentID:    node.ParentID,
			Type:        node.Type,
			Path:        node.Path,
			Children:    children,
		}
	}

	var tree []model.PermissionTreeResponse
	for _, r := range roots {
		tree = append(tree, buildNode(r))
	}
	return tree
}
