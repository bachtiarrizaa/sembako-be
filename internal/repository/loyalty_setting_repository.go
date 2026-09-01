package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type LoyaltySettingRepository interface {
	Get(ctx context.Context) (*entity.LoyaltySetting, error)
	Update(ctx context.Context, setting *entity.LoyaltySetting) error
}

type loyaltySettingRepositoryImpl struct {
	db *gorm.DB
}

func NewLoyaltySettingRepository(db *gorm.DB) LoyaltySettingRepository {
	return &loyaltySettingRepositoryImpl{db: db}
}

func (r *loyaltySettingRepositoryImpl) Get(ctx context.Context) (*entity.LoyaltySetting, error) {
	var loyaltySetting entity.LoyaltySetting
	if err := r.db.WithContext(ctx).First(&loyaltySetting).Error; err != nil {
		return nil, err
	}
	return &loyaltySetting, nil
}

func (r *loyaltySettingRepositoryImpl) Update(ctx context.Context, setting *entity.LoyaltySetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}
