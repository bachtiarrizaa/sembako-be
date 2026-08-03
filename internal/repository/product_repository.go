package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindProducts(ctx context.Context, req model.PaginationRequest) ([]entity.Product, int64, error)
	FindById(ctx context.Context, id string) (*entity.Product, error)
	FindByName(ctx context.Context, name string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) ProductRepository
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: tx}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryImpl) FindProducts(ctx context.Context, req model.PaginationRequest) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Category").
		Preload("BaseUnit").
		Preload("Units").
		Preload("Units.Unit").
		Offset(offset).Limit(req.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepositoryImpl) FindById(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("BaseUnit").
		Preload("Units").
		Preload("Units.Unit").
		First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id).Error
}
