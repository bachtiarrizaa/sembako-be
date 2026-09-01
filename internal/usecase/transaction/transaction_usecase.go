package transaction

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, cashierID string, req model.CreateTransactionRequest) (*model.TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id string, cashierID string, role string) (*model.TransactionResponse, error)
	ListTransactions(ctx context.Context, req model.ListTransactionsRequest, cashierID string, role string) ([]model.TransactionResponse, utils.Pagination, error)
	VoidTransaction(ctx context.Context, id string, cashierID string, req model.VoidTransactionRequest) (*model.TransactionResponse, error)
}

type transactionUsecaseImpl struct {
	db                 *gorm.DB
	transactionRepo    repository.TransactionRepository
	shiftRepo          repository.ShiftRepository
	customerRepo       repository.CustomerRepository
	productUnitRepo    repository.ProductUnitRepository
	productRepo        repository.ProductRepository
	stockRepo          repository.StockRepository
	stockMutationRepo  repository.StockMutationRepository
	purchaseBatchRepo  repository.PurchaseBatchRepository
	costAllocationRepo repository.TransactionItemCostAllocationRepository
	loyaltySettingRepo repository.LoyaltySettingRepository
	pointLedgerRepo    repository.PointLedgerRepository
}

func NewTransactionUsecase(
	db *gorm.DB,
	transactionRepo repository.TransactionRepository,
	shiftRepo repository.ShiftRepository,
	customerRepo repository.CustomerRepository,
	productUnitRepo repository.ProductUnitRepository,
	productRepo repository.ProductRepository,
	stockRepo repository.StockRepository,
	stockMutationRepo repository.StockMutationRepository,
	purchaseBatchRepo repository.PurchaseBatchRepository,
	costAllocationRepo repository.TransactionItemCostAllocationRepository,
	loyaltySettingRepo repository.LoyaltySettingRepository,
	pointLedgerRepo repository.PointLedgerRepository,
) TransactionUsecase {
	return &transactionUsecaseImpl{
		db:                 db,
		transactionRepo:    transactionRepo,
		shiftRepo:          shiftRepo,
		customerRepo:       customerRepo,
		productUnitRepo:    productUnitRepo,
		productRepo:        productRepo,
		stockRepo:          stockRepo,
		stockMutationRepo:  stockMutationRepo,
		purchaseBatchRepo:  purchaseBatchRepo,
		costAllocationRepo: costAllocationRepo,
		loyaltySettingRepo: loyaltySettingRepo,
		pointLedgerRepo:    pointLedgerRepo,
	}
}

func (u *transactionUsecaseImpl) GetTransactionByID(ctx context.Context, id string, cashierID string, role string) (*model.TransactionResponse, error) {
	trx, err := u.transactionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("transaction not found")
		}
		return nil, errs.NewInternal("failed to fetch transaction")
	}

	if role == "cashier" && trx.CashierID != cashierID {
		return nil, errs.NewForbidden("you can only view your own transactions")
	}

	res := model.ToTransactionResponse(trx)

	// Sembunyikan HPP dan Margin untuk kasir (hanya untuk admin / owner)
	if role == "cashier" {
		for i := range res.Items {
			res.Items[i].TotalCost = nil
			res.Items[i].Margin = nil
		}
	}

	return &res, nil
}

func (u *transactionUsecaseImpl) ListTransactions(ctx context.Context, req model.ListTransactionsRequest, cashierID string, role string) ([]model.TransactionResponse, utils.Pagination, error) {
	var restrictToCashierID *string
	if role == "cashier" {
		restrictToCashierID = &cashierID
	}

	transactions, total, err := u.transactionRepo.FindTransactions(ctx, req, restrictToCashierID)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to list transactions: " + err.Error())
	}

	responses := make([]model.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		responses = append(responses, model.ToTransactionResponse(&transactions[i]))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return responses, pagination, nil
}
