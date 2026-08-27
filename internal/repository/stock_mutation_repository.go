package repository

import (
	"context"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type StockMutationRepository interface {
	Create(ctx context.Context, mutation *entity.StockMutation) error
	FindProductMutations(ctx context.Context, productID string, page int, limit int) ([]entity.StockMutation, int64, error)
	HasStockMutationReferences(ctx context.Context, productID string) (bool, error)
	WithTx(tx *gorm.DB) StockMutationRepository
}

type stockMutationRepositoryImpl struct {
	db *gorm.DB
}

func NewStockMutationRepository(db *gorm.DB) StockMutationRepository {
	return &stockMutationRepositoryImpl{db: db}
}

func (r *stockMutationRepositoryImpl) WithTx(tx *gorm.DB) StockMutationRepository {
	return &stockMutationRepositoryImpl{db: tx}
}

func (r *stockMutationRepositoryImpl) Create(ctx context.Context, mutation *entity.StockMutation) error {
	return r.db.WithContext(ctx).Create(mutation).Error
}

func (r *stockMutationRepositoryImpl) FindProductMutations(ctx context.Context, productID string, page int, limit int) ([]entity.StockMutation, int64, error) {
	var mutations []entity.StockMutation
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.StockMutation{}).Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Creator").
		Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&mutations).Error; err != nil {
		return nil, 0, err
	}

	return mutations, total, nil
}

func (r *stockMutationRepositoryImpl) HasStockMutationReferences(ctx context.Context, productID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.StockMutation{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
