package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type AuthUsecase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	blacklistRepo    repository.BlacklistRepository
	jwtAccessSecret  string
	jwtAccessTTL     time.Duration
	jwtRefreshTTL    time.Duration
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	blacklistRepo repository.BlacklistRepository,
	jwtAccessSecret string,
	jwtAccessTTL time.Duration,
	jwtRefreshTTL time.Duration,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		blacklistRepo:    blacklistRepo,
		jwtAccessSecret:  jwtAccessSecret,
		jwtAccessTTL:     jwtAccessTTL,
		jwtRefreshTTL:    jwtRefreshTTL,
	}
}

type LoginResult struct {
	Response     *model.LoginResponse
	RefreshToken string
}

func (uc *AuthUsecase) Login(ctx context.Context, req model.LoginRequest) (*LoginResult, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewUnauthorized("invalid email or password")
		}
		return nil, errs.NewInternal("failed to fetch user data")
	}

	if !user.IsActive {
		return nil, errs.NewForbidden("account is inactive, please contact admin")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errs.NewUnauthorized("invalid email or password")
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, errs.NewInternal("invalid user id format")
	}

	accessToken, err := utils.GenerateAccessToken(userID, user.Role.Name, uc.jwtAccessSecret, uc.jwtAccessTTL)
	if err != nil {
		return nil, errs.NewInternal("failed to generate access token")
	}

	rawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errs.NewInternal("failed to generate refresh token")
	}

	refreshTokenRecord := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashRefreshToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(uc.jwtRefreshTTL),
	}
	if err := uc.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, errs.NewInternal("failed to store refresh token")
	}

	return &LoginResult{
		Response: &model.LoginResponse{
			AccessToken: accessToken,
			User:        model.ToUserResponse(user),
		},
		RefreshToken: rawRefreshToken,
	}, nil
}

func (uc *AuthUsecase) Refresh(ctx context.Context, rawToken string) (*model.RefreshResponse, string, error) {
	tokenHash := utils.HashRefreshToken(rawToken)

	record, err := uc.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errs.NewUnauthorized("invalid or expired refresh token")
		}
		return nil, "", errs.NewInternal("failed to validate refresh token")
	}

	if time.Now().After(record.ExpiresAt) {
		_ = uc.refreshTokenRepo.DeleteByUserID(ctx, record.UserID)
		return nil, "", errs.NewUnauthorized("refresh token expired, please login again")
	}

	user, err := uc.userRepo.FindByID(ctx, uuid.MustParse(record.UserID))
	if err != nil {
		return nil, "", errs.NewInternal("failed to fetch user data")
	}

	if !user.IsActive {
		return nil, "", errs.NewForbidden("account is inactive, please contact admin")
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, "", errs.NewInternal("invalid user id format")
	}

	if err := uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, "", errs.NewInternal("failed to rotate refresh token")
	}

	accessToken, err := utils.GenerateAccessToken(userID, user.Role.Name, uc.jwtAccessSecret, uc.jwtAccessTTL)
	if err != nil {
		return nil, "", errs.NewInternal("failed to generate access token")
	}

	newRawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, "", errs.NewInternal("failed to generate refresh token")
	}

	newRecord := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashRefreshToken(newRawRefreshToken),
		ExpiresAt: time.Now().Add(uc.jwtRefreshTTL),
	}
	if err := uc.refreshTokenRepo.Create(ctx, newRecord); err != nil {
		return nil, "", errs.NewInternal("failed to store refresh token")
	}

	return &model.RefreshResponse{AccessToken: accessToken}, newRawRefreshToken, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, userID uuid.UUID, accessToken string) error {
	if err := uc.refreshTokenRepo.DeleteByUserID(ctx, userID.String()); err != nil {
		return errs.NewInternal("failed to revoke session")
	}
	claims, err := utils.ParseAccessToken(accessToken, uc.jwtAccessSecret)
	if err != nil {
		return nil
	}

	tokenHash := utils.HashRefreshToken(accessToken)
	if err := uc.blacklistRepo.Blacklist(ctx, tokenHash, claims.ExpiresAt.Time); err != nil {
		return errs.NewInternal("failed to blacklist token")
	}

	return nil
}
