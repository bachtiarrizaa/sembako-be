package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(ctx context.Context, shift *entity.Shift) error
	FindActiveByCashierID(ctx context.Context, cashierID uuid.UUID) (*entity.Shift, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Shift, error)
	Update(ctx context.Context, shift *entity.Shift) error
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
		First(&shift, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepositoryImpl) Update(ctx context.Context, shift *entity.Shift) error {
	return r.db.WithContext(ctx).Save(shift).Error
}
