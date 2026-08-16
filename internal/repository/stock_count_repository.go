package repository

import (
	"context"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type StockCountRepository interface {
	Create(ctx context.Context, sc *entity.StockCount) error
	FindByID(ctx context.Context, id string) (*entity.StockCount, error)
	FindStockCounts(ctx context.Context, productID string, status string, page int, limit int) ([]entity.StockCount, int64, error)
	Update(ctx context.Context, sc *entity.StockCount) error
	WithTx(tx *gorm.DB) StockCountRepository
}

type stockCountRepositoryImpl struct {
	db *gorm.DB
}

func NewStockCountRepository(db *gorm.DB) StockCountRepository {
	return &stockCountRepositoryImpl{db: db}
}

func (r *stockCountRepositoryImpl) WithTx(tx *gorm.DB) StockCountRepository {
	return &stockCountRepositoryImpl{db: tx}
}

func (r *stockCountRepositoryImpl) Create(ctx context.Context, sc *entity.StockCount) error {
	return r.db.WithContext(ctx).Create(sc).Error
}

func (r *stockCountRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.StockCount, error) {
	var sc entity.StockCount
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Submitter").
		Preload("Approver").
		First(&sc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (r *stockCountRepositoryImpl) FindStockCounts(ctx context.Context, productID string, status string, page int, limit int) ([]entity.StockCount, int64, error) {
	var counts []entity.StockCount
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.StockCount{})

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Preload("Product").
		Preload("Submitter").
		Preload("Approver").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&counts).Error; err != nil {
		return nil, 0, err
	}

	return counts, total, nil
}

func (r *stockCountRepositoryImpl) Update(ctx context.Context, sc *entity.StockCount) error {
	return r.db.WithContext(ctx).Save(sc).Error
}
