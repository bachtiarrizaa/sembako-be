package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type PointLedgerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, ledger *entity.PointLedger) error
	FindByCustomerID(ctx context.Context, customerID string, req model.PaginationRequest) ([]entity.PointLedger, int64, error)
	FindExpiredEarnLedgersByCustomerID(ctx context.Context, tx *gorm.DB, customerID string) ([]entity.PointLedger, error)
	FindExpiredEarnLedgersGlobal(ctx context.Context) ([]entity.PointLedger, error)
}

type pointLedgerRepositoryImpl struct {
	db *gorm.DB
}

func NewPointLedgerRepository(db *gorm.DB) PointLedgerRepository {
	return &pointLedgerRepositoryImpl{db: db}
}

func (r *pointLedgerRepositoryImpl) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *pointLedgerRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, ledger *entity.PointLedger) error {
	return r.getDB(tx).WithContext(ctx).Create(ledger).Error
}

func (r *pointLedgerRepositoryImpl) FindByCustomerID(ctx context.Context, customerID string, req model.PaginationRequest) ([]entity.PointLedger, int64, error) {
	var ledgers []entity.PointLedger
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PointLedger{}).Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&ledgers).Error
	if err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *pointLedgerRepositoryImpl) FindExpiredEarnLedgersByCustomerID(ctx context.Context, tx *gorm.DB, customerID string) ([]entity.PointLedger, error) {
	var ledgers []entity.PointLedger
	err := r.getDB(tx).WithContext(ctx).
		Where("customer_id = ? AND type = ? AND expired_at IS NOT NULL AND expired_at < NOW()", customerID, entity.PointLedgerTypeEarn).
		Find(&ledgers).Error
	if err != nil {
		return nil, err
	}
	return ledgers, nil
}

func (r *pointLedgerRepositoryImpl) FindExpiredEarnLedgersGlobal(ctx context.Context) ([]entity.PointLedger, error) {
	var ledgers []entity.PointLedger
	err := r.db.WithContext(ctx).
		Where("type = ? AND expired_at IS NOT NULL AND expired_at < NOW()", entity.PointLedgerTypeEarn).
		Find(&ledgers).Error
	if err != nil {
		return nil, err
	}
	return ledgers, nil
}
