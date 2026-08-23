package usecase

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftUsecase struct {
	shiftRepo       repository.ShiftRepository
	transactionRepo repository.TransactionRepository
}

func NewShiftUsecase(shiftRepo repository.ShiftRepository, transactionRepo repository.TransactionRepository) *ShiftUsecase {
	return &ShiftUsecase{
		shiftRepo:       shiftRepo,
		transactionRepo: transactionRepo,
	}
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

const discrepancyTolerance = 1000

func (u *ShiftUsecase) CloseShift(ctx context.Context, shiftID uuid.UUID, cashierID uuid.UUID, req model.CloseShiftRequest) (*model.ShiftResponse, error) {
	shift, err := u.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("shift not found")
		}
		return nil, errs.NewInternal("failed to fetch shift")
	}

	if shift.CashierID != cashierID.String() {
		return nil, errs.NewForbidden("you can only close your own shift")
	}

	if shift.Status != entity.ShiftStatusOpen {
		return nil, errs.NewConflict("shift is already closed")
	}

	totalCashSales, err := u.getTotalCashSales(ctx, shiftID)
	if err != nil {
		return nil, errs.NewInternal("failed to calculate cash sales")
	}

	systemBalance := shift.OpeningBalance + totalCashSales
	discrepancy := req.ClosingBalance - systemBalance

	// SHIFT-07: wajib catatan alasan kalau selisih > Rp 1.000
	if math.Abs(discrepancy) > discrepancyTolerance && (req.DiscrepancyNote == nil || *req.DiscrepancyNote == "") {
		return nil, errs.NewBadRequest("discrepancy note is required when discrepancy exceeds Rp 1,000")
	}

	now := time.Now()
	shift.ClosingBalance = &req.ClosingBalance
	shift.SystemBalance = &systemBalance
	shift.Discrepancy = &discrepancy
	shift.DiscrepancyNote = req.DiscrepancyNote
	shift.Status = entity.ShiftStatusClosed
	shift.ClosedAt = &now

	rowsAffected, err := u.shiftRepo.Update(ctx, shift)
	if err != nil {
		return nil, errs.NewInternal("failed to close shift")
	}
	if rowsAffected == 0 {
		return nil, errs.NewConflict("shift was already closed by another request")
	}

	resp := model.ToShiftResponse(shift)
	return &resp, nil
}

func (u *ShiftUsecase) getTotalCashSales(ctx context.Context, shiftID uuid.UUID) (float64, error) {
	return u.transactionRepo.GetTotalCashSalesByShift(ctx, shiftID.String())
}

func (u *ShiftUsecase) ForceCloseShift(ctx context.Context, shiftID uuid.UUID, adminID uuid.UUID, req model.ForceCloseShiftRequest) (*model.ShiftResponse, error) {
	shift, err := u.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("shift not found")
		}
		return nil, errs.NewInternal("failed to fetch shift")
	}

	if shift.Status != entity.ShiftStatusOpen {
		return nil, errs.NewConflict("shift is already closed")
	}

	totalCashSales, err := u.getTotalCashSales(ctx, shiftID)
	if err != nil {
		return nil, errs.NewInternal("failed to calculate cash sales")
	}

	systemBalance := shift.OpeningBalance + totalCashSales
	discrepancy := req.ClosingBalance - systemBalance

	// SHIFT-07: wajib catatan alasan kalau selisih > Rp 1.000
	if math.Abs(discrepancy) > discrepancyTolerance && (req.DiscrepancyNote == nil || *req.DiscrepancyNote == "") {
		return nil, errs.NewBadRequest("discrepancy note is required when discrepancy exceeds Rp 1,000")
	}

	adminIDStr := adminID.String()
	now := time.Now()

	shift.ClosingBalance = &req.ClosingBalance
	shift.SystemBalance = &systemBalance
	shift.Discrepancy = &discrepancy
	shift.DiscrepancyNote = req.DiscrepancyNote
	shift.ForceCloseReason = &req.Reason
	shift.ForceClosedBy = &adminIDStr
	shift.Status = entity.ShiftStatusClosed
	shift.ClosedAt = &now

	rowsAffected, err := u.shiftRepo.Update(ctx, shift)
	if err != nil {
		return nil, errs.NewInternal("failed to close shift")
	}
	if rowsAffected == 0 {
		return nil, errs.NewConflict("shift was already closed by another request")
	}

	updated, err := u.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, errs.NewInternal("failed to load closed shift")
	}

	resp := model.ToShiftResponse(updated)
	return &resp, nil
}

func (u *ShiftUsecase) ListShifts(ctx context.Context, req model.ListShiftsRequest, userID uuid.UUID, role string) ([]model.ShiftListItemResponse, utils.Pagination, error) {
	var restrictToCashierID *uuid.UUID
	if strings.ToLower(role) == "cashier" {
		restrictToCashierID = &userID
	}

	shifts, total, err := u.shiftRepo.FindShifts(ctx, req, restrictToCashierID)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch shifts")
	}

	res := make([]model.ShiftListItemResponse, 0, len(shifts))
	for _, s := range shifts {
		res = append(res, model.ToShiftListItemResponse(&s))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)

	return res, pagination, nil
}

func (u *ShiftUsecase) GetShiftDetail(ctx context.Context, shiftID uuid.UUID, userID uuid.UUID, role string) (*model.ShiftResponse, error) {
	shift, err := u.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("shift not found")
		}
		return nil, errs.NewInternal("failed to fetch shift detail")
	}

	if strings.ToLower(role) == "cashier" && shift.CashierID != userID.String() {
		return nil, errs.NewForbidden("you can only view your own shift detail")
	}

	resp := model.ToShiftResponse(shift)
	return &resp, nil
}
