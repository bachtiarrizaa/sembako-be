package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type StoreConfigurationUsecase struct {
	storeConfigurationRepo repository.StoreConfigurationRepository
}

func NewStoreConfigurationUsecase(
	storeConfigurationRepo repository.StoreConfigurationRepository,
) *StoreConfigurationUsecase {
	return &StoreConfigurationUsecase{storeConfigurationRepo: storeConfigurationRepo}
}

func (u *StoreConfigurationUsecase) Get(ctx context.Context) (*model.StoreConfigurationResponse, error) {
	storeConfig, err := u.storeConfigurationRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("store configuration not found")
		}
		return nil, errs.NewInternal("failed to fetch store configuration")
	}
	resp := model.ToStoreConfigurationResponse(storeConfig)
	return &resp, nil
}
