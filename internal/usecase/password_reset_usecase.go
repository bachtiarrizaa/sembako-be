package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/brevo"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type PasswordResetUsecase struct {
	userRepo                repository.UserRepository
	refreshTokenRepo        repository.RefreshTokenRepository
	passwordResetRepo       repository.PasswordResetRepository
	brevoService            brevo.BrevoService
	frontendResetUrl        string
	resetTokenExpireMinutes int
}

func NewPasswordResetUsecase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordResetRepo repository.PasswordResetRepository,
	brevoService brevo.BrevoService,
	frontendResetUrl string,
	resetTokenExpireMinutes int,
) *PasswordResetUsecase {
	return &PasswordResetUsecase{
		userRepo:                userRepo,
		refreshTokenRepo:        refreshTokenRepo,
		passwordResetRepo:       passwordResetRepo,
		brevoService:            brevoService,
		frontendResetUrl:        frontendResetUrl,
		resetTokenExpireMinutes: resetTokenExpireMinutes,
	}
}

func (uc *PasswordResetUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return errs.NewInternal("failed to process password reset request")
	}

	if !user.IsActive {
		return nil
	}

	_ = uc.passwordResetRepo.InvalidateAllByUserID(ctx, user.ID)

	rawToken, err := generateSecureToken()
	if err != nil {
		return errs.NewInternal("failed to generate reset token")
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(time.Duration(uc.resetTokenExpireMinutes) * time.Minute)

	resetTokenRecord := &entity.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := uc.passwordResetRepo.Create(ctx, resetTokenRecord); err != nil {
		return errs.NewInternal("failed to store reset token")
	}

	resetLink := uc.frontendResetUrl + "?token=" + rawToken
	if err := uc.brevoService.SendPasswordResetEmail(user.Email, user.Name, resetLink); err != nil {
		return errs.NewInternal("failed to send password reset email")
	}

	return nil
}

func (uc *PasswordResetUsecase) ResetPassword(ctx context.Context, token string, newPassword string) error {
	tokenHash := hashToken(token)

	resetRecord, err := uc.passwordResetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewBadRequest("invalid or expired reset token")
		}
		return errs.NewInternal("failed to validate reset token")
	}

	if resetRecord.UsedAt != nil || time.Now().After(resetRecord.ExpiresAt) {
		return errs.NewBadRequest("invalid or expired reset token")
	}

	userUUID, err := uuid.Parse(resetRecord.UserID)
	if err != nil {
		return errs.NewInternal("invalid user id format")
	}

	user, err := uc.userRepo.FindByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewBadRequest("user not found")
		}
		return errs.NewInternal("failed to fetch user data")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errs.NewInternal("failed to hash password")
	}

	user.PasswordHash = hashedPassword
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return errs.NewInternal("failed to update password")
	}

	if err := uc.passwordResetRepo.MarkAsUsed(ctx, resetRecord.ID); err != nil {
		return errs.NewInternal("failed to update reset token status")
	}

	_ = uc.refreshTokenRepo.DeleteByUserID(ctx, user.ID)

	return nil
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
