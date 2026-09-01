package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type PointLedgerUsecase struct {
	pointLedgerRepo repository.PointLedgerRepository
	loyaltyRepo     repository.LoyaltySettingRepository
	customerRepo    repository.CustomerRepository
}

func NewPointLedgerUsecase(
	pointLedgerRepo repository.PointLedgerRepository,
	loyaltyRepo repository.LoyaltySettingRepository,
	customerRepo repository.CustomerRepository,
) *PointLedgerUsecase {
	return &PointLedgerUsecase{
		pointLedgerRepo: pointLedgerRepo,
		loyaltyRepo:     loyaltyRepo,
		customerRepo:    customerRepo,
	}
}

func (u *PointLedgerUsecase) GetCustomerLedgers(ctx context.Context, customerID string, req model.PaginationRequest) ([]model.PointLedgerResponse, utils.Pagination, error) {
	_, err := u.customerRepo.FindById(ctx, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.Pagination{}, errs.NewNotFound("customer not found")
		}
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch customer")
	}

	ledgers, total, err := u.pointLedgerRepo.FindByCustomerID(ctx, customerID, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch point ledgers")
	}

	res := []model.PointLedgerResponse{}
	for _, l := range ledgers {
		res = append(res, model.ToPointLedgerResponse(&l))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (u *PointLedgerUsecase) ProcessLazyExpiry(ctx context.Context, tx *gorm.DB, customerID string) error {
	setting, err := u.loyaltyRepo.Get(ctx)
	if err != nil || !setting.IsExpiryActive {
		return nil
	}

	expiredEarns, err := u.pointLedgerRepo.FindExpiredEarnLedgersByCustomerID(ctx, tx, customerID)
	if err != nil || len(expiredEarns) == 0 {
		return nil
	}

	customer, err := u.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return nil
	}

	totalExpired := 0
	for _, earn := range expiredEarns {
		if earn.Points > 0 {
			totalExpired += earn.Points
			earn.Points = 0
			_ = u.pointLedgerRepo.Create(ctx, tx, &entity.PointLedger{
				CustomerID:  customerID,
				Type:        entity.PointLedgerTypeExpire,
				Points:      -earn.Points,
				Description: "Point expired",
			})
		}
	}

	if totalExpired > 0 {
		newPoints := customer.TotalPoints - totalExpired
		if newPoints < 0 {
			newPoints = 0
		}
		customer.TotalPoints = newPoints
		_ = u.customerRepo.Update(ctx, customer)
	}

	return nil
}
