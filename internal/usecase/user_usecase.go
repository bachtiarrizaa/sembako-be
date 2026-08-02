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

func (uc *UserUsecase) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.UserResponse, error) {
	_, err := uc.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("role not found")
		}
		return nil, errs.NewInternal("failed to fetch role data")
	}

	existingEmail, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, errs.NewConflict("email already registered")
	}

	if req.Username != nil && *req.Username != "" {
		existingUsername, err := uc.userRepo.FindByUsername(ctx, *req.Username)
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

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, errs.NewInternal("failed to create user")
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (uc *UserUsecase) GetAllUsers(ctx context.Context, req model.PaginationRequest) ([]model.UserResponse, utils.Pagination, error) {
	users, total, err := uc.userRepo.FindAllPaginated(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch users")
	}

	var res []model.UserResponse
	for _, u := range users {
		res = append(res, model.ToUserResponse(&u))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}
	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (uc *UserUsecase) GetMe(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}
	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("user not found")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}

	// 1. Check if role exists
	_, err = uc.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("role not found")
		}
		return nil, errs.NewInternal("failed to fetch role data")
	}

	// 2. Check if email is unique (if changed)
	if req.Email != user.Email {
		existingEmail, err := uc.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingEmail != nil {
			return nil, errs.NewConflict("email already registered")
		}
	}

	// 3. Check if username is unique (if changed)
	if req.Username != nil && *req.Username != "" {
		if user.Username == nil || *req.Username != *user.Username {
			existingUsername, err := uc.userRepo.FindByUsername(ctx, *req.Username)
			if err == nil && existingUsername != nil {
				return nil, errs.NewConflict("username already taken")
			}
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Username = req.Username
	user.RoleID = req.RoleID
	user.IsActive = *req.IsActive

	// 4. Update password if provided
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errs.NewInternal("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, errs.NewInternal("failed to update user")
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("user not found")
		}
		return errs.NewInternal("failed to fetch user data")
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete user")
	}
	return nil
}
