package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftUsecase struct {
	shiftRepo repository.ShiftRepository
}

func NewShiftUsecase(shiftRepo repository.ShiftRepository) *ShiftUsecase {
	return &ShiftUsecase{shiftRepo: shiftRepo}
}

func (u *ShiftUsecase) OpenShift(ctx context.Context, cashierID uuid.UUID, req model.OpenShiftRequest) (*model.ShiftResponse, error) {
	existing, err := u.shiftRepo.FindActiveByCashierID(ctx, cashierID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check active shift")
	}
	if existing != nil {
		return nil, errs.NewConflict("cashier already has an active shift")
	}

	shift := &entity.Shift{
		CashierID:      cashierID.String(),
		OpeningBalance: req.OpeningBalance,
		Status:         entity.ShiftStatusOpen,
	}

	if err := u.shiftRepo.Create(ctx, shift); err != nil {
		return nil, errs.NewInternal("failed to open shift")
	}

	shiftID, err := uuid.Parse(shift.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to parse created shift id")
	}

	created, err := u.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, errs.NewInternal("failed to load created shift")
	}

	resp := model.ToShiftResponse(created)
	return &resp, nil
}

func (u *ShiftUsecase) GetActiveShift(ctx context.Context, cashierID uuid.UUID) (*model.ShiftResponse, error) {
	shift, err := u.shiftRepo.FindActiveByCashierID(ctx, cashierID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("no active shift found")
		}
		return nil, errs.NewInternal("failed to fetch active shift")
	}

	resp := model.ToShiftResponse(shift)
	return &resp, nil
}
