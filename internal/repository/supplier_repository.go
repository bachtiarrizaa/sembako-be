package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *entity.Supplier) error
	FindSuppliers(ctx context.Context, req model.PaginationRequest) ([]entity.Supplier, int64, error)
	FindById(ctx context.Context, id string) (*entity.Supplier, error)
	FindByName(ctx context.Context, name string) (*entity.Supplier, error)
	Update(ctx context.Context, category *entity.Supplier) error
	Delete(ctx context.Context, id string) error
}

type supplierRepositoryImpl struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepositoryImpl{db: db}
}

func (r *supplierRepositoryImpl) Create(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

func (r *supplierRepositoryImpl) FindSuppliers(ctx context.Context, req model.PaginationRequest) ([]entity.Supplier, int64, error) {
	var suppliers []entity.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Supplier{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}
	return suppliers, total, nil
}

func (r *supplierRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepositoryImpl) FindById(ctx context.Context, id string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).First(&supplier, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepositoryImpl) Update(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

func (r *supplierRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Supplier{}, "id = ?", id).Error
}
