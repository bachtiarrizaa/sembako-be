package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
)

type UnitRepository interface {
	Create(ctx context.Context, unit *entity.Unit) error
	FindUnits(ctx context.Context, req model.PaginationRequest) ([]entity.Unit, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Unit, error)
	FindByName(ctx context.Context, name string) (*entity.Unit, error)
	Update(ctx context.Context, unit *entity.Unit) error
	Delete(ctx context.Context, id string) error
}

type unitRepositoryImpl struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepositoryImpl{db: db}
}

func (r *unitRepositoryImpl) Create(ctx context.Context, unit *entity.Unit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *unitRepositoryImpl) FindUnits(ctx context.Context, req model.PaginationRequest) ([]entity.Unit, int64, error) {
	var units []entity.Unit
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Unit{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&units).Error; err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r *unitRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.WithContext(ctx).First(&unit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepositoryImpl) Update(ctx context.Context, unit *entity.Unit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *unitRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Unit{}, "id = ?", id).Error
}
