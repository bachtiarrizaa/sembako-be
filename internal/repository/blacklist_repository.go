package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

type BlacklistRepository interface {
	Blacklist(ctx context.Context, tokenHash string, expiresAt time.Time) error
	IsBlacklisted(ctx context.Context, tokenHash string) (bool, error)
}

type blacklistRepository struct {
	db *gorm.DB
}

func NewBlacklistRepository(db *gorm.DB) BlacklistRepository {
	return &blacklistRepository{db: db}
}

func (r *blacklistRepository) Blacklist(ctx context.Context, tokenHash string, expiresAt time.Time) error {
	item := &entity.BlacklistedToken{
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *blacklistRepository) IsBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BlacklistedToken{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
