package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, cashierID string, req model.CreateTransactionRequest) (*model.TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id string, cashierID string, role string) (*model.TransactionResponse, error)
	ListTransactions(ctx context.Context, req model.ListTransactionsRequest, cashierID string, role string) ([]model.TransactionResponse, utils.Pagination, error)
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
	}
}

type preparedTrxItem struct {
	productUnit *entity.ProductUnit
	product     *entity.Product
	qtyInBase   float64
	qty         float64
	unitPrice   float64
	subtotal    float64
}

func (u *transactionUsecaseImpl) CreateTransaction(ctx context.Context, cashierID string, req model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errs.NewBadRequest("invalid cashier id")
	}

	// 1. Validasi Shift Aktif untuk Kasir
	shift, err := u.shiftRepo.FindActiveByCashierID(ctx, cashierUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewConflict("no active shift found for this cashier")
		}
		return nil, errs.NewInternal("failed to check active shift")
	}

	// 2. Validasi Customer (jika ada)
	var customer *entity.Customer
	if req.CustomerID != nil && *req.CustomerID != "" {
		customer, err = u.customerRepo.FindById(ctx, *req.CustomerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errs.NewNotFound("customer not found")
			}
			return nil, errs.NewInternal("failed to fetch customer")
		}
		if !customer.IsActive {
			return nil, errs.NewConflict("cannot record transaction for an inactive customer")
		}
	}

	// 3. Validasi & Kalkulasi Item
	if len(req.Items) == 0 {
		return nil, errs.NewBadRequest("transaction items cannot be empty")
	}

	preparedItems := make([]preparedTrxItem, 0, len(req.Items))
	var subtotal float64 = 0

	for _, item := range req.Items {
		// Fetch product unit
		unit, err := u.productUnitRepo.FindByID(ctx, item.ProductUnitID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errs.NewNotFound("product unit not found: " + item.ProductUnitID)
			}
			return nil, errs.NewInternal("failed to fetch product unit")
		}

		if !unit.IsActive {
			return nil, errs.NewConflict("cannot sell with an inactive product unit: " + unit.Unit.Name)
		}

		// Fetch parent product
		product, err := u.productRepo.FindById(ctx, unit.ProductID)
		if err != nil {
			return nil, errs.NewNotFound("parent product not found")
		}

		if !product.IsActive {
			return nil, errs.NewConflict("cannot sell an inactive product: " + product.Name)
		}

		qtyInBase := item.Qty * unit.ConversionToBase
		itemSubtotal := item.Qty * unit.SellingPrice
		subtotal += itemSubtotal

		preparedItems = append(preparedItems, preparedTrxItem{
			productUnit: unit,
			product:     product,
			qtyInBase:   qtyInBase,
			qty:         item.Qty,
			unitPrice:   unit.SellingPrice,
			subtotal:    itemSubtotal,
		})
	}

	// 4. Validasi Keuangan Transaksi
	var total float64 = subtotal // Di Fase 1, belum ada diskon per total transaksi atau poin
	var cashReceived *float64
	var changeGiven *float64
	var manualPaidConfirmation *bool

	if req.PaymentMethod == "cash" {
		if req.CashReceived == nil {
			return nil, errs.NewBadRequest("cash received must be provided for cash payment method")
		}
		if *req.CashReceived < total {
			return nil, errs.NewBadRequest(fmt.Sprintf("cash received (%.2f) must be equal or greater than total (%.2f)", *req.CashReceived, total))
		}
		cashReceived = req.CashReceived
		change := *req.CashReceived - total
		changeGiven = &change
	} else {
		// QRIS/Transfer
		paidConfirm := true
		manualPaidConfirmation = &paidConfirm
	}

	// 5. Generate Receipt Number: TRX-YYYYMMDD-XXXXXX
	receiptNumber := fmt.Sprintf("TRX-%s-%s", time.Now().Format("20060102"), strings.ToUpper(uuid.New().String()[:6]))

	var transactionID string

	// DB Transaction block
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := u.transactionRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		costAllocationRepoTx := u.costAllocationRepo.WithTx(tx)

		// Create Transaction Header
		trx := &entity.Transaction{
			ReceiptNumber:          receiptNumber,
			CashierID:              cashierID,
			ShiftID:                shift.ID,
			CustomerID:             req.CustomerID,
			PaymentMethod:          req.PaymentMethod,
			Subtotal:               subtotal,
			TotalDiscount:          0, // Fase 1 default 0
			PointsUsed:             0, // Fase 1 default 0
			PointsDiscountValue:    0, // Fase 1 default 0
			PointsEarned:           0, // Fase 1 default 0
			Total:                  total,
			CashReceived:           cashReceived,
			ChangeGiven:            changeGiven,
			ManualPaidConfirmation: manualPaidConfirmation,
			Status:                 "completed",
		}

		if err := transactionRepoTx.Create(ctx, trx); err != nil {
			return errs.NewInternal("failed to create transaction: " + err.Error())
		}
		transactionID = trx.ID

		// Process items
		for _, pi := range preparedItems {
			// Check stock sufficiency
			currentStock, err := stockRepoTx.GetByProductID(ctx, pi.product.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errs.NewConflict("insufficient stock for product: " + pi.product.Name)
				}
				return errs.NewInternal("failed to check stock: " + err.Error())
			}

			if currentStock.QtyBaseUnit < pi.qtyInBase {
				return errs.NewConflict("insufficient stock for product: " + pi.product.Name)
			}

			// FIFO Costing Engine:
			// Fetch active batches for this product sorted by purchase_date ASC, created_at ASC (FIFO)
			batches, err := purchaseBatchRepoTx.FindActiveBatchesByProductID(ctx, pi.product.ID)
			if err != nil {
				return errs.NewInternal("failed to fetch active purchase batches: " + err.Error())
			}

			var remainingQtyToAllocate = pi.qtyInBase
			var totalCostForItem float64 = 0
			type allocationItem struct {
				batchID            string
				qtyAllocated       float64
				purchasePrice      float64
				costSubtotal       float64
				updatedBatchRecord *entity.PurchaseBatch
			}
			var allocations []allocationItem

			for i := range batches {
				if remainingQtyToAllocate <= 0 {
					break
				}
				batch := &batches[i]
				
				// Determine how much we can allocate from this batch
				var qtyAllocated float64
				if batch.RemainingQty >= remainingQtyToAllocate {
					qtyAllocated = remainingQtyToAllocate
					batch.RemainingQty = batch.RemainingQty - remainingQtyToAllocate
					remainingQtyToAllocate = 0
				} else {
					qtyAllocated = batch.RemainingQty
					remainingQtyToAllocate = remainingQtyToAllocate - batch.RemainingQty
					batch.RemainingQty = 0
				}

				costSubtotal := qtyAllocated * batch.PurchasePrice
				totalCostForItem += costSubtotal

				allocations = append(allocations, allocationItem{
					batchID:            batch.ID,
					qtyAllocated:       qtyAllocated,
					purchasePrice:      batch.PurchasePrice,
					costSubtotal:       costSubtotal,
					updatedBatchRecord: batch,
				})
			}

			// If we still have remaining quantity to allocate, it means the granular batch stocks are out of sync/insufficient compared to the global stock.
			if remainingQtyToAllocate > 0 {
				return errs.NewConflict("insufficient purchase batch stock for product: " + pi.product.Name)
			}

			// Create Detail Item (TransactionItem) with computed HPP/Cost & Margin
			margin := pi.subtotal - totalCostForItem
			detail := entity.TransactionItem{
				TransactionID:   trx.ID,
				ProductUnitID:   pi.productUnit.ID,
				Qty:             pi.qty,
				UnitPrice:       pi.unitPrice,
				DiscountApplied: 0, // Fase 1 default 0
				Subtotal:        pi.subtotal,
				TotalCost:       &totalCostForItem,
				Margin:          &margin,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return errs.NewInternal("failed to create transaction detail: " + err.Error())
			}

			// Save allocations and update batches in DB
			for _, alloc := range allocations {
				// Update batch remaining qty
				if err := purchaseBatchRepoTx.Update(ctx, alloc.updatedBatchRecord); err != nil {
					return errs.NewInternal("failed to update purchase batch: " + alloc.updatedBatchRecord.ID + " - " + err.Error())
				}

				// Create TransactionItemCostAllocation
				costAlloc := &entity.TransactionItemCostAllocation{
					TransactionItemID:   detail.ID,
					PurchaseBatchID:     alloc.batchID,
					QtyAllocated:        alloc.qtyAllocated,
					PurchasePriceAtSale: alloc.purchasePrice,
					CostSubtotal:        alloc.costSubtotal,
				}

				if err := costAllocationRepoTx.Create(ctx, costAlloc); err != nil {
					return errs.NewInternal("failed to create cost allocation record: " + err.Error())
				}
			}

			// Update Stock Cache
			qtyBefore := currentStock.QtyBaseUnit
			qtyAfter := qtyBefore - pi.qtyInBase
			currentStock.QtyBaseUnit = qtyAfter

			if err := stockRepoTx.UpsertStock(ctx, currentStock); err != nil {
				return errs.NewInternal("failed to update stock cache: " + err.Error())
			}

			// Create Stock Mutation log
			noteStr := "Sale: " + receiptNumber
			mutation := &entity.StockMutation{
				ProductID:   pi.product.ID,
				Type:        "out",
				Qty:         pi.qtyInBase,
				QtyBefore:   qtyBefore,
				QtyAfter:    qtyAfter,
				Source:      "sale",
				ReferenceID: &trx.ID,
				Note:        &noteStr,
				CreatedBy:   cashierID,
			}

			if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
				return errs.NewInternal("failed to log stock mutation: " + err.Error())
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return u.GetTransactionByID(ctx, transactionID, cashierID, "cashier")
}

func (u *transactionUsecaseImpl) GetTransactionByID(ctx context.Context, id string, cashierID string, role string) (*model.TransactionResponse, error) {
	trx, err := u.transactionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("transaction not found")
		}
		return nil, errs.NewInternal("failed to fetch transaction: " + err.Error())
	}

	// Restrict cashier to view only their own transaction
	if strings.ToLower(role) == "cashier" && trx.CashierID != cashierID {
		return nil, errs.NewForbidden("you can only view your own transactions")
	}

	resp := model.ToTransactionResponse(trx)
	return &resp, nil
}

func (u *transactionUsecaseImpl) ListTransactions(ctx context.Context, req model.ListTransactionsRequest, cashierID string, role string) ([]model.TransactionResponse, utils.Pagination, error) {
	var restrictToCashierID *string
	if strings.ToLower(role) == "cashier" {
		restrictToCashierID = &cashierID
	}

	transactions, total, err := u.transactionRepo.FindTransactions(ctx, req, restrictToCashierID)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to list transactions: " + err.Error())
	}

	resp := make([]model.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		resp = append(resp, model.ToTransactionResponse(&transactions[i]))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return resp, pagination, nil
}
