package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type TransactionItemCostAllocationRepository interface {
	Create(ctx context.Context, allocation *entity.TransactionItemCostAllocation) error
	FindByTransactionItemID(ctx context.Context, itemID string) ([]entity.TransactionItemCostAllocation, error)
	WithTx(tx *gorm.DB) TransactionItemCostAllocationRepository
}

type transactionItemCostAllocationRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionItemCostAllocationRepository(db *gorm.DB) TransactionItemCostAllocationRepository {
	return &transactionItemCostAllocationRepositoryImpl{db: db}
}

func (r *transactionItemCostAllocationRepositoryImpl) WithTx(tx *gorm.DB) TransactionItemCostAllocationRepository {
	return &transactionItemCostAllocationRepositoryImpl{db: tx}
}

func (r *transactionItemCostAllocationRepositoryImpl) Create(ctx context.Context, allocation *entity.TransactionItemCostAllocation) error {
	return r.db.WithContext(ctx).Create(allocation).Error
}

func (r *transactionItemCostAllocationRepositoryImpl) FindByTransactionItemID(ctx context.Context, itemID string) ([]entity.TransactionItemCostAllocation, error) {
	var allocations []entity.TransactionItemCostAllocation
	err := r.db.WithContext(ctx).
		Where("transaction_item_id = ?", itemID).
		Find(&allocations).Error

	if err != nil {
		return nil, err
	}
	return allocations, nil
}
