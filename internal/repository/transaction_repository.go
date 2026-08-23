package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id string) (*entity.Transaction, error)
	FindTransactions(ctx context.Context, req model.ListTransactionsRequest, restrictToCashierID *string) ([]entity.Transaction, int64, error)
	GetTotalCashSalesByShift(ctx context.Context, shiftID string) (float64, error)
	WithTx(tx *gorm.DB) TransactionRepository
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) WithTx(tx *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: tx}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.WithContext(ctx).
		Preload("Cashier").
		Preload("Shift").
		Preload("Customer").
		Preload("VoidedByUser").
		Preload("Items").
		Preload("Items.ProductUnit").
		Preload("Items.ProductUnit.Product").
		Preload("Items.ProductUnit.Unit").
		First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) FindTransactions(ctx context.Context, req model.ListTransactionsRequest, restrictToCashierID *string) ([]entity.Transaction, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.Transaction{}).
		Preload("Cashier").
		Preload("Customer")

	if restrictToCashierID != nil {
		query = query.Where("cashier_id = ?", *restrictToCashierID)
	} else if req.CashierID != nil && *req.CashierID != "" {
		query = query.Where("cashier_id = ?", *req.CashierID)
	}

	if req.CustomerID != nil && *req.CustomerID != "" {
		query = query.Where("customer_id = ?", *req.CustomerID)
	}
	if req.PaymentMethod != nil && *req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", *req.PaymentMethod)
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartDate != nil && *req.StartDate != "" {
		query = query.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil && *req.EndDate != "" {
		query = query.Where("created_at <= ?", *req.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var transactions []entity.Transaction
	offset := (req.Page - 1) * req.Limit
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(req.Limit).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepositoryImpl) GetTotalCashSalesByShift(ctx context.Context, shiftID string) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("shift_id = ? AND payment_method = 'cash' AND status = 'completed'", shiftID).
		Select("COALESCE(SUM(total), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
