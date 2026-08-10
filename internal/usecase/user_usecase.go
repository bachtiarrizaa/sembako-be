package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type UserUsecase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserUsecase(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.UserResponse, error) {
	_, err := u.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("role not found")
		}
		return nil, errs.NewInternal("failed to fetch role data")
	}

	existingEmail, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, errs.NewConflict("email already registered")
	}

	if req.Username != nil && *req.Username != "" {
		existingUsername, err := u.userRepo.FindByUsername(ctx, *req.Username)
		if err == nil && existingUsername != nil {
			return nil, errs.NewConflict("username already taken")
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errs.NewInternal("failed to hash password")
	}

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		IsActive:     *req.IsActive,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, errs.NewInternal("failed to create user")
	}

	return toUserResponse(user), nil
}

func (u *UserUsecase) GetUsersWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.UserResponse, utils.Pagination, error) {
	users, total, err := u.userRepo.FindUsers(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch users")
	}

	res := []model.UserResponse{}
	for _, usr := range users {
		res = append(res, *toUserResponse(&usr))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}
	return toUserResponse(user), nil
}

func (u *UserUsecase) GetMe(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}
	return toUserResponse(user), nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}

	_, err = u.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("role not found")
		}
		return nil, errs.NewInternal("failed to fetch role data")
	}

	if req.Email != user.Email {
		existingEmail, err := u.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingEmail != nil {
			return nil, errs.NewConflict("email already registered")
		}
	}

	if req.Username != nil && *req.Username != "" {
		if user.Username == nil || *req.Username != *user.Username {
			existingUsername, err := u.userRepo.FindByUsername(ctx, *req.Username)
			if err == nil && existingUsername != nil {
				return nil, errs.NewConflict("username already taken")
			}
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Username = req.Username
	user.RoleID = req.RoleID

	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errs.NewInternal("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, errs.NewInternal("failed to update user")
	}

	return toUserResponse(user), nil
}

func (u *UserUsecase) UpdateStatus(ctx context.Context, id uuid.UUID, req model.UpdateStatusUserRequest) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}

	user.IsActive = req.IsActive

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, errs.NewInternal("failed to update user status")
	}

	return toUserResponse(user), nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("user not found")
		}
		return errs.NewInternal("failed to fetch user data")
	}

	if err := u.userRepo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete user")
	}
	return nil
}

func toUserResponse(u *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Username: u.Username,
		Role: model.UserWithRole{
			ID:   u.Role.ID,
			Name: u.Role.Name,
		},
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
