package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(ctx context.Context, discount *entity.Discount) error
	FindDiscounts(ctx context.Context, req model.PaginationRequest) ([]entity.Discount, int64, error)
	FindById(ctx context.Context, id string) (*entity.Discount, error)
	FindByName(ctx context.Context, name string) (*entity.Discount, error)
	Update(ctx context.Context, discount *entity.Discount) error
	Delete(ctx context.Context, id string) error
}

type discountRepositoryImpl struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepositoryImpl{db: db}
}

func (r *discountRepositoryImpl) Create(ctx context.Context, discount *entity.Discount) error {
	return r.db.WithContext(ctx).Create(discount).Error
}

func (r *discountRepositoryImpl) FindDiscounts(ctx context.Context, req model.PaginationRequest) ([]entity.Discount, int64, error) {
	var discounts []entity.Discount
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Discount{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&discounts).Error; err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}

func (r *discountRepositoryImpl) FindById(ctx context.Context, id string) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.WithContext(ctx).First(&discount, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *discountRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&discount).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *discountRepositoryImpl) Update(ctx context.Context, discount *entity.Discount) error {
	return r.db.WithContext(ctx).Select("Name", "Type", "Value", "StartDate", "EndDate", "IsActive", "UpdatedAt").Save(discount).Error
}

func (r *discountRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Discount{}, "id = ?", id).Error
}
