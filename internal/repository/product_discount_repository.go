package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"gorm.io/gorm"
)

type ProductDiscountRepository interface {
	Create(ctx context.Context, productDiscount *entity.ProductDiscount) error
	FindByID(ctx context.Context, id string) (*entity.ProductDiscount, error)
	FindByDiscountAndProduct(ctx context.Context, discountID string, productID string) (*entity.ProductDiscount, error)
	FindProductDiscounts(ctx context.Context, req model.GetProductDiscountsRequest) ([]entity.ProductDiscount, int64, error)
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

func (r *productDiscountRepositoryImpl) FindProductDiscounts(ctx context.Context, req model.GetProductDiscountsRequest) ([]entity.ProductDiscount, int64, error) {
	var productDiscounts []entity.ProductDiscount
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.ProductDiscount{}).
		Preload("Product.Category").
		Preload("Product").
		Preload("Discount")

	if req.DiscountID != "" {
		query = query.Where("discount_id = ?", req.DiscountID)
	}

	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.CalculateOffset(req.Page, req.Limit)
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&productDiscounts).
		Error; err != nil {
		return nil, 0, err
	}

	return productDiscounts, total, nil
}
