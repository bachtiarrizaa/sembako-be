package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, token *entity.PasswordResetToken) error
	FindByTokenHash(ctx context.Context, hash string) (*entity.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, id string) error
	InvalidateAllByUserID(ctx context.Context, userID string) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, hash string) (*entity.PasswordResetToken, error) {
	var token entity.PasswordResetToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.PasswordResetToken{}).Where("id = ?", id).Updates(map[string]interface{}{
		"used_at":    &now,
		"updated_at": now,
	}).Error
}

func (r *passwordResetRepository) InvalidateAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Updates(map[string]interface{}{
			"used_at":    &now,
			"updated_at": now,
		}).Error
}
