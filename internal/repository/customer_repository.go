package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	FindCustomers(ctx context.Context, req model.PaginationRequest) ([]entity.Customer, int64, error)
	FindById(ctx context.Context, id string) (*entity.Customer, error)
	FindByName(ctx context.Context, name string) (*entity.Customer, error)
	FindByPhoneNumber(ctx context.Context, phone string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id string) error
}

type customerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepositoryImpl{db: db}
}

func (r *customerRepositoryImpl) Create(ctx context.Context, customer *entity.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepositoryImpl) FindCustomers(ctx context.Context, req model.PaginationRequest) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Customer{})

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR phone_number ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *customerRepositoryImpl) FindById(ctx context.Context, id string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepositoryImpl) FindByPhoneNumber(ctx context.Context, phone string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Where("phone_number = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepositoryImpl) Update(ctx context.Context, customer *entity.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Customer{}, "id = ?", id).Error
}
