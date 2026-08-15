package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type PurchaseBatchRepository interface {
	Create(ctx context.Context, batch *entity.PurchaseBatch) error
	FindPurchaseBatches(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]entity.PurchaseBatch, int64, error)
	FindByID(ctx context.Context, id string) (*entity.PurchaseBatch, error)
	Update(ctx context.Context, batch *entity.PurchaseBatch) error
	Delete(ctx context.Context, id string) error
	HasPurchaseReferences(ctx context.Context, productID string) (bool, error)
	WithTx(tx *gorm.DB) PurchaseBatchRepository
}

type purchaseBatchRepositoryImpl struct {
	db *gorm.DB
}

func NewPurchaseBatchRepository(db *gorm.DB) PurchaseBatchRepository {
	return &purchaseBatchRepositoryImpl{db: db}
}

func (r *purchaseBatchRepositoryImpl) WithTx(tx *gorm.DB) PurchaseBatchRepository {
	return &purchaseBatchRepositoryImpl{db: tx}
}

func (r *purchaseBatchRepositoryImpl) Create(ctx context.Context, batch *entity.PurchaseBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

func (r *purchaseBatchRepositoryImpl) FindPurchaseBatches(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]entity.PurchaseBatch, int64, error) {
	var batches []entity.PurchaseBatch
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PurchaseBatch{}).
		Joins("JOIN products ON products.id = purchase_batches.product_id").
		Joins("JOIN suppliers ON suppliers.id = purchase_batches.supplier_id")

	if req.Search != "" {
		query = query.Where("products.name ILIKE ? OR purchase_batches.invoice_number ILIKE ? OR suppliers.name ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.SupplierID != "" {
		query = query.Where("purchase_batches.supplier_id = ?", req.SupplierID)
	}
	if req.ProductID != "" {
		query = query.Where("purchase_batches.product_id = ?", req.ProductID)
	}
	if req.StartDate != "" {
		query = query.Where("purchase_batches.purchase_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("purchase_batches.purchase_date <= ?", req.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Product").
		Preload("Supplier").
		Preload("Creator").
		Order("purchase_batches.purchase_date DESC, purchase_batches.created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&batches).Error; err != nil {
		return nil, 0, err
	}

	return batches, total, nil
}

func (r *purchaseBatchRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.PurchaseBatch, error) {
	var batch entity.PurchaseBatch
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Supplier").
		Preload("Creator").
		First(&batch, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *purchaseBatchRepositoryImpl) Update(ctx context.Context, batch *entity.PurchaseBatch) error {
	return r.db.WithContext(ctx).Save(batch).Error
}

func (r *purchaseBatchRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.PurchaseBatch{}, "id = ?", id).Error
}

func (r *purchaseBatchRepositoryImpl) HasPurchaseReferences(ctx context.Context, productID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.PurchaseBatch{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
