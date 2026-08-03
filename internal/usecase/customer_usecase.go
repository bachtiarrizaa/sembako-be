package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type CustomerUsecase struct {
	repo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (u *CustomerUsecase) CreateCustomer(ctx context.Context, req model.CreateCustomerRequest) (*model.CustomerResponse, error) {
	existingName, err := u.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing customer name")
	}
	if existingName != nil {
		return nil, errs.NewConflict("customer name already exists")
	}

	existingPhone, err := u.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing phone number")
	}
	if existingPhone != nil {
		return nil, errs.NewConflict("phone number already registered")
	}

	customer := &entity.Customer{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		IsActive:    true,
		TotalPoints: 0,
	}

	if err := u.repo.Create(ctx, customer); err != nil {
		return nil, errs.NewInternal("failed to create customer")
	}

	return toCustomerResponse(customer), nil
}

func (u *CustomerUsecase) GetCustomersWithPagination(ctx context.Context, req model.PaginationRequest) ([]model.CustomerResponse, utils.Pagination, error) {
	customers, total, err := u.repo.FindCustomers(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch customers")
	}

	res := []model.CustomerResponse{}
	for _, c := range customers {
		res = append(res, *toCustomerResponse(&c))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (u *CustomerUsecase) GetCustomerById(ctx context.Context, id string) (*model.CustomerResponse, error) {
	customer, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("customer not found")
	}
	return toCustomerResponse(customer), nil
}

func (u *CustomerUsecase) UpdateCustomer(ctx context.Context, id string, req model.UpdateCustomerRequest) (*model.CustomerResponse, error) {
	customer, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("customer not found")
	}

	existingName, err := u.repo.FindByName(ctx, req.Name)
	if err == nil && existingName != nil && existingName.ID != id {
		return nil, errs.NewConflict("customer name already exists")
	}

	existingPhone, err := u.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingPhone != nil && existingPhone.ID != id {
		return nil, errs.NewConflict("phone number already registered")
	}

	customer.Name = req.Name
	customer.PhoneNumber = req.PhoneNumber
	customer.Address = req.Address

	if err := u.repo.Update(ctx, customer); err != nil {
		return nil, errs.NewInternal("failed to update customer")
	}

	return toCustomerResponse(customer), nil
}

func (u *CustomerUsecase) UpdateStatus(ctx context.Context, id string, req model.UpdateStatusCustomerRequest) (*model.CustomerResponse, error) {
	customer, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, errs.NewNotFound("customer not found")
	}

	customer.IsActive = req.IsActive

	if err := u.repo.Update(ctx, customer); err != nil {
		return nil, errs.NewInternal("failed to update customer status")
	}

	return toCustomerResponse(customer), nil
}

func (u *CustomerUsecase) DeleteCustomer(ctx context.Context, id string) error {
	_, err := u.repo.FindById(ctx, id)
	if err != nil {
		return errs.NewNotFound("customer not found")
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return errs.NewInternal("failed to delete customer")
	}
	return nil
}

func toCustomerResponse(customer *entity.Customer) *model.CustomerResponse {
	return &model.CustomerResponse{
		ID:          customer.ID,
		Name:        customer.Name,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		TotalPoints: customer.TotalPoints,
		IsActive:    customer.IsActive,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}
}
