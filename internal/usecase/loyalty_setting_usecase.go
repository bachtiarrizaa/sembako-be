package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type LoyaltySettingUsecase struct {
	loyaltySettingRepo repository.LoyaltySettingRepository
}

func NewLoyaltySettingUsecase(
	loyaltySettingRepo repository.LoyaltySettingRepository,
) *LoyaltySettingUsecase {
	return &LoyaltySettingUsecase{loyaltySettingRepo: loyaltySettingRepo}
}

func (u *LoyaltySettingUsecase) Get(ctx context.Context) (*model.LoyaltySettingResponse, error) {
	loyaltySetting, err := u.loyaltySettingRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("loyalty setting not found")
		}
		return nil, errs.NewInternal("failed to fetch loyalty setting")
	}
	resp := model.ToLoyaltySettingResponse(loyaltySetting)
	return &resp, nil
}

func (u *LoyaltySettingUsecase) Update(ctx context.Context, req model.UpdateLoyaltySettingRequest) (*model.LoyaltySettingResponse, error) {
	setting, err := u.loyaltySettingRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("loyalty setting not found")
		}
		return nil, errs.NewInternal("failed to fetch loyalty setting")
	}

	setting.EarningRate = req.EarningRate
	setting.RedemptionRate = req.RedemptionRate
	setting.MinimumRedeem = req.MinimumRedeem
	setting.IsExpiryActive = req.IsExpiryActive
	setting.ExpiryMonths = req.ExpiryMonths

	if err := u.loyaltySettingRepo.Update(ctx, setting); err != nil {
		return nil, errs.NewInternal("failed to update loyalty setting")
	}

	resp := model.ToLoyaltySettingResponse(setting)
	return &resp, nil
}
