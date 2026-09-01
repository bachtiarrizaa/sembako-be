package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type StoreConfigurationRepository interface {
	Get(ctx context.Context) (*entity.StoreConfiguration, error)
	Update(ctx context.Context, storeConfiguration *entity.StoreConfiguration) error
}

type storeConfigurationImpl struct {
	db *gorm.DB
}

func NewStoreConfigurationRepository(db *gorm.DB) StoreConfigurationRepository {
	return &storeConfigurationImpl{db: db}
}

func (r *storeConfigurationImpl) Get(ctx context.Context) (*entity.StoreConfiguration, error) {
	var config entity.StoreConfiguration
	if err := r.db.WithContext(ctx).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *storeConfigurationImpl) Update(ctx context.Context, storeConfiguration *entity.StoreConfiguration) error {
	return r.db.WithContext(ctx).Save(storeConfiguration).Error
}
