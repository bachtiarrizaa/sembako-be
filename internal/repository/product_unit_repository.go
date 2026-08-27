package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type ProductUnitRepository interface {
	Create(ctx context.Context, unit *entity.ProductUnit) error
	CreateMany(ctx context.Context, units []entity.ProductUnit) error
	FindByProductID(ctx context.Context, productID string) ([]entity.ProductUnit, error)
	FindByID(ctx context.Context, id string) (*entity.ProductUnit, error)
	Update(ctx context.Context, unit *entity.ProductUnit) error
	Delete(ctx context.Context, id string) error
	DeleteByProductID(ctx context.Context, productID string) error
	DeactivateAllByProductID(ctx context.Context, productID string) error
	ActivateAllByProductID(ctx context.Context, productID string) error
	WithTx(tx *gorm.DB) ProductUnitRepository
}

type productUnitRepositoryImpl struct {
	db *gorm.DB
}

func NewProductUnitRepository(db *gorm.DB) ProductUnitRepository {
	return &productUnitRepositoryImpl{db: db}
}

func (r *productUnitRepositoryImpl) WithTx(tx *gorm.DB) ProductUnitRepository {
	return &productUnitRepositoryImpl{db: tx}
}

func (r *productUnitRepositoryImpl) Create(ctx context.Context, unit *entity.ProductUnit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *productUnitRepositoryImpl) CreateMany(ctx context.Context, units []entity.ProductUnit) error {
	if len(units) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&units).Error
}

func (r *productUnitRepositoryImpl) FindByProductID(ctx context.Context, productID string) ([]entity.ProductUnit, error) {
	var units []entity.ProductUnit
	if err := r.db.WithContext(ctx).
		Preload("Unit").
		Where("product_id = ?", productID).
		Order("product_units.is_base_unit DESC, product_units.conversion_to_base ASC, product_units.id ASC").
		Find(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

func (r *productUnitRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.ProductUnit, error) {
	var unit entity.ProductUnit
	if err := r.db.WithContext(ctx).
		Preload("Unit").
		First(&unit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *productUnitRepositoryImpl) Update(ctx context.Context, unit *entity.ProductUnit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *productUnitRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductUnit{}, "id = ?", id).Error
}

func (r *productUnitRepositoryImpl) DeleteByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&entity.ProductUnit{}).Error
}

func (r *productUnitRepositoryImpl) DeactivateAllByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Model(&entity.ProductUnit{}).Where("product_id = ?", productID).Update("is_active", false).Error
}

func (r *productUnitRepositoryImpl) ActivateAllByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Model(&entity.ProductUnit{}).Where("product_id = ?", productID).Update("is_active", true).Error
}
