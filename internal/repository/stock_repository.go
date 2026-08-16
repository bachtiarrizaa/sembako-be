package repository

import (
	"context"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type StockRepository interface {
	GetByProductID(ctx context.Context, productID string) (*entity.Stock, error)
	UpsertStock(ctx context.Context, stock *entity.Stock) error
	WithTx(tx *gorm.DB) StockRepository
}

type stockRepositoryImpl struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepositoryImpl{db: db}
}

func (r *stockRepositoryImpl) WithTx(tx *gorm.DB) StockRepository {
	return &stockRepositoryImpl{db: tx}
}

func (r *stockRepositoryImpl) GetByProductID(ctx context.Context, productID string) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.db.WithContext(ctx).First(&stock, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepositoryImpl) UpsertStock(ctx context.Context, stock *entity.Stock) error {
	return r.db.WithContext(ctx).Save(stock).Error
}
