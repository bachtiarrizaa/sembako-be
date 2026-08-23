package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(ctx context.Context, shift *entity.Shift) error
	FindActiveByCashierID(ctx context.Context, cashierID uuid.UUID) (*entity.Shift, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Shift, error)
	FindShifts(ctx context.Context, req model.ListShiftsRequest, restrictToCashierID *uuid.UUID) ([]entity.Shift, int64, error)
	Update(ctx context.Context, shift *entity.Shift) (int64, error)
	WithTx(tx *gorm.DB) ShiftRepository
}
type shiftRepositoryImpl struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepositoryImpl{db: db}
}

func (r *shiftRepositoryImpl) WithTx(tx *gorm.DB) ShiftRepository {
	return &shiftRepositoryImpl{db: tx}
}

func (r *shiftRepositoryImpl) Create(ctx context.Context, shift *entity.Shift) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

func (r *shiftRepositoryImpl) FindActiveByCashierID(ctx context.Context, cashierID uuid.UUID) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Where("cashier_id = ? AND status = ?", cashierID, entity.ShiftStatusOpen).
		First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Shift, error) {
	var shift entity.Shift
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Preload("ForceClosedByUser").
		First(&shift, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepositoryImpl) FindShifts(ctx context.Context, req model.ListShiftsRequest, restrictToCashierID *uuid.UUID) ([]entity.Shift, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.Shift{}).
		Preload("Cashier").
		Preload("ForceClosedByUser")

	// row-level access: kasir cuma boleh liat shift miliknya sendiri
	if restrictToCashierID != nil {
		query = query.Where("cashier_id = ?", *restrictToCashierID)
	}

	// filter opsional dari query param (khusus Admin/Owner)
	if req.CashierID != nil && *req.CashierID != "" {
		query = query.Where("cashier_id = ?", *req.CashierID)
	}
	if req.StartDate != nil && *req.StartDate != "" {
		query = query.Where("opened_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil && *req.EndDate != "" {
		query = query.Where("opened_at <= ?", *req.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var shifts []entity.Shift
	offset := (req.Page - 1) * req.Limit
	err := query.
		Order("opened_at DESC").
		Offset(offset).
		Limit(req.Limit).
		Find(&shifts).Error
	if err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

func (r *shiftRepositoryImpl) Update(ctx context.Context, shift *entity.Shift) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.Shift{}).
		Where("id = ? AND status = ?", shift.ID, entity.ShiftStatusOpen).
		Updates(shift)
	return result.RowsAffected, result.Error
}
