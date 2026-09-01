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

func (u *StoreConfigurationUsecase) Update(ctx context.Context, req model.UpdateStoreConfigurationRequest) (*model.StoreConfigurationResponse, error) {
	storeConfig, err := u.storeConfigurationRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("store configuration not found")
		}
		return nil, errs.NewInternal("failed to fetch store configuration")
	}

	storeConfig.StoreName = req.StoreName
	storeConfig.StoreAddress = &req.StoreAddress
	storeConfig.StorePhone = &req.StorePhone
	storeConfig.ReceiptHeaderText = &req.ReceiptHeaderText
	storeConfig.ReceiptFooterText = &req.ReceiptFooterText
	storeConfig.ReceiptShowCashierName = req.ReceiptShowCashierName
	storeConfig.ReceiptShowCustomerName = req.ReceiptShowCustomerName
	storeConfig.ShiftDiscrepancyTolerance = req.ShiftDiscrepancyTolerance

	if err := u.storeConfigurationRepo.Update(ctx, storeConfig); err != nil {
		return nil, errs.NewInternal("failed to update store configuration")
	}

	resp := model.ToStoreConfigurationResponse(storeConfig)
	return &resp, nil
}
