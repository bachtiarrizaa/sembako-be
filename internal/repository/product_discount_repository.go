package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type ProductDiscountRepository interface {
	Create(ctx context.Context, productDiscount *entity.ProductDiscount) error
	FindByID(ctx context.Context, id string) (*entity.ProductDiscount, error)
	FindByDiscountAndProduct(ctx context.Context, discountID string, productID string) (*entity.ProductDiscount, error)
	WithTx(tx *gorm.DB) ProductDiscountRepository
}

type productDiscountRepositoryImpl struct {
	db *gorm.DB
}

func NewProductDiscountRepository(db *gorm.DB) ProductDiscountRepository {
	return &productDiscountRepositoryImpl{db: db}
}

func (r *productDiscountRepositoryImpl) WithTx(tx *gorm.DB) ProductDiscountRepository {
	return &productDiscountRepositoryImpl{db: tx}
}

func (r *productDiscountRepositoryImpl) Create(ctx context.Context, productDiscount *entity.ProductDiscount) error {
	return r.db.WithContext(ctx).Create(productDiscount).Error
}

func (r *productDiscountRepositoryImpl) FindByDiscountAndProduct(ctx context.Context, discountID string, productID string) (*entity.ProductDiscount, error) {
	var productDiscount entity.ProductDiscount
	if err := r.db.WithContext(ctx).
		Where("discount_id = ? AND product_id = ?", discountID, productID).
		First(&productDiscount).Error; err != nil {
		return nil, err
	}
	return &productDiscount, nil
}

func (r *productDiscountRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.ProductDiscount, error) {
	var productDiscount entity.ProductDiscount
	if err := r.db.WithContext(ctx).
		Preload("Product.Category").
		Preload("Product").
		Preload("Discount").
		First(&productDiscount, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &productDiscount, nil
}
