package transaction

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type preparedTrxItem struct {
	productUnit     *entity.ProductUnit
	product         *entity.Product
	qtyInBase       float64
	qty             float64
	unitPrice       float64
	discountApplied float64
	subtotal        float64
}

type allocationItem struct {
	batchID            string
	qtyAllocated       float64
	purchasePrice      float64
	costSubtotal       float64
	updatedBatchRecord *entity.PurchaseBatch
}

func (u *transactionUsecaseImpl) CreateTransaction(ctx context.Context, cashierID string, req model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errs.NewBadRequest("invalid cashier id")
	}

	// 1. Validasi Shift Aktif Kasir
	shift, err := u.shiftRepo.FindActiveByCashierID(ctx, cashierUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewConflict("no active shift found for this cashier")
		}
		return nil, errs.NewInternal("failed to check active shift")
	}

	// 2. Validasi Customer (jika ada)
	if req.CustomerID != nil && *req.CustomerID != "" {
		customer, err := u.customerRepo.FindById(ctx, *req.CustomerID)
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

	// 3. Validasi & Kalkulasi Item Beserta Diskon Aktif
	if len(req.Items) == 0 {
		return nil, errs.NewBadRequest("transaction items cannot be empty")
	}

	preparedItems, subtotal, totalDiscount, err := u.prepareItemsAndDiscounts(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	// 4. Validasi Keuangan & Loyalty Points Transaksi
	total := subtotal - totalDiscount
	var pointsUsed int = 0
	var pointsDiscountValue float64 = 0
	var pointsEarned int = 0
	var targetCustomer *entity.Customer

	loyaltySetting, _ := u.loyaltySettingRepo.Get(ctx)

	if req.CustomerID != nil && *req.CustomerID != "" {
		c, err := u.customerRepo.FindById(ctx, *req.CustomerID)
		if err == nil && c.IsActive {
			targetCustomer = c

			if req.UsePoints != nil && *req.UsePoints && loyaltySetting != nil && loyaltySetting.RedemptionRate > 0 {
				if targetCustomer.TotalPoints >= loyaltySetting.MinimumRedeem {
					maxPointsNeeded := int(total / loyaltySetting.RedemptionRate)
					if targetCustomer.TotalPoints <= maxPointsNeeded {
						pointsUsed = targetCustomer.TotalPoints
					} else {
						pointsUsed = maxPointsNeeded
					}
					pointsDiscountValue = float64(pointsUsed) * loyaltySetting.RedemptionRate
					total = total - pointsDiscountValue
					if total < 0 {
						total = 0
					}
				}
			}

			if loyaltySetting != nil && loyaltySetting.EarningRate > 0 {
				pointsEarned = int(total / loyaltySetting.EarningRate)
			}
		}
	}

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
		paidConfirm := true
		manualPaidConfirmation = &paidConfirm
	}

	// 5. Generate Receipt Number
	receiptNumber := fmt.Sprintf("TRX-%s-%s", time.Now().Format("20060102"), strings.ToUpper(uuid.New().String()[:6]))

	// 6. Eksekusi Database Transaction (Atomic)
	var transactionID string
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := u.transactionRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		costAllocationRepoTx := u.costAllocationRepo.WithTx(tx)

		// Create Header
		trx := &entity.Transaction{
			ReceiptNumber:          receiptNumber,
			CashierID:              cashierID,
			ShiftID:                shift.ID,
			CustomerID:             req.CustomerID,
			PaymentMethod:          req.PaymentMethod,
			Subtotal:               subtotal,
			TotalDiscount:          totalDiscount,
			PointsUsed:             pointsUsed,
			PointsDiscountValue:    pointsDiscountValue,
			PointsEarned:           pointsEarned,
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

		// Log Loyalty Point Ledgers & Update Customer Total Points
		if targetCustomer != nil {
			if pointsUsed > 0 {
				redeemLedger := &entity.PointLedger{
					CustomerID:    targetCustomer.ID,
					TransactionID: &trx.ID,
					Type:          entity.PointLedgerTypeRedeem,
					Points:        -pointsUsed,
					Description:   fmt.Sprintf("Points redeemed for transaction %s", receiptNumber),
				}
				if err := u.pointLedgerRepo.Create(ctx, tx, redeemLedger); err != nil {
					return errs.NewInternal("failed to log redeem point ledger: " + err.Error())
				}
			}

			if pointsEarned > 0 {
				var expiredAt *time.Time
				if loyaltySetting != nil && loyaltySetting.IsExpiryActive && loyaltySetting.ExpiryMonths > 0 {
					exp := time.Now().AddDate(0, loyaltySetting.ExpiryMonths, 0)
					expiredAt = &exp
				}

				earnLedger := &entity.PointLedger{
					CustomerID:    targetCustomer.ID,
					TransactionID: &trx.ID,
					Type:          entity.PointLedgerTypeEarn,
					Points:        pointsEarned,
					Description:   fmt.Sprintf("Points earned from transaction %s", receiptNumber),
					ExpiredAt:     expiredAt,
				}
				if err := u.pointLedgerRepo.Create(ctx, tx, earnLedger); err != nil {
					return errs.NewInternal("failed to log earn point ledger: " + err.Error())
				}
			}

			newTotalPoints := targetCustomer.TotalPoints - pointsUsed + pointsEarned
			if newTotalPoints < 0 {
				newTotalPoints = 0
			}
			targetCustomer.TotalPoints = newTotalPoints
			if err := u.customerRepo.Update(ctx, targetCustomer); err != nil {
				return errs.NewInternal("failed to update customer total points: " + err.Error())
			}
		}

		// Process Items, FIFO Costing, & Stock Reduction
		for _, pi := range preparedItems {
			// Cek kecukupan stok
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

			// Alokasikan HPP dengan FIFO
			totalCostForItem, allocations, err := u.allocateFIFOBatches(ctx, purchaseBatchRepoTx, pi.product, pi.qtyInBase)
			if err != nil {
				return err
			}

			// Simpan Detail Transaksi
			margin := pi.subtotal - totalCostForItem
			detail := entity.TransactionItem{
				TransactionID:   trx.ID,
				ProductUnitID:   pi.productUnit.ID,
				Qty:             pi.qty,
				UnitPrice:       pi.unitPrice,
				DiscountApplied: pi.discountApplied,
				Subtotal:        pi.subtotal,
				TotalCost:       &totalCostForItem,
				Margin:          &margin,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return errs.NewInternal("failed to create transaction detail: " + err.Error())
			}

			// Simpan Cost Allocations & Update Sisa Batch
			for _, alloc := range allocations {
				if err := purchaseBatchRepoTx.Update(ctx, alloc.updatedBatchRecord); err != nil {
					return errs.NewInternal("failed to update purchase batch: " + alloc.updatedBatchRecord.ID + " - " + err.Error())
				}

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

			// Potong Stok Agregat
			qtyBefore := currentStock.QtyBaseUnit
			qtyAfter := qtyBefore - pi.qtyInBase
			currentStock.QtyBaseUnit = qtyAfter

			if err := stockRepoTx.UpsertStock(ctx, currentStock); err != nil {
				return errs.NewInternal("failed to update stock cache: " + err.Error())
			}

			// Catat Kartu Stok Mutasi OUT
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
		if appErr, ok := txErr.(*errs.AppError); ok {
			return nil, appErr
		}
		return nil, errs.NewInternal("failed to process transaction: " + txErr.Error())
	}

	// 7. Muat Ulang Transaksi Lengkap untuk Response
	completeTrx, err := u.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, errs.NewInternal("failed to fetch created transaction")
	}

	res := model.ToTransactionResponse(completeTrx)
	return &res, nil
}

func (u *transactionUsecaseImpl) prepareItemsAndDiscounts(ctx context.Context, items []model.CreateTransactionItem) ([]preparedTrxItem, float64, float64, error) {
	preparedItems := make([]preparedTrxItem, 0, len(items))
	var subtotal float64 = 0
	var totalDiscount float64 = 0
	now := time.Now()

	for _, item := range items {
		unit, err := u.productUnitRepo.FindByID(ctx, item.ProductUnitID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, 0, errs.NewNotFound("product unit not found: " + item.ProductUnitID)
			}
			return nil, 0, 0, errs.NewInternal("failed to fetch product unit")
		}

		if !unit.IsActive {
			return nil, 0, 0, errs.NewConflict("cannot sell with an inactive product unit: " + unit.Unit.Name)
		}

		product, err := u.productRepo.FindById(ctx, unit.ProductID)
		if err != nil {
			return nil, 0, 0, errs.NewNotFound("parent product not found")
		}

		if !product.IsActive {
			return nil, 0, 0, errs.NewConflict("cannot sell an inactive product: " + product.Name)
		}

		// Cari diskon aktif untuk produk ini hari ini
		var activeDiscount *entity.Discount
		for _, pd := range product.ProductDiscounts {
			if !pd.IsActive {
				continue
			}
			d := pd.Discount
			if !d.IsActive {
				continue
			}
			if d.StartDate != nil && d.StartDate.After(now) {
				continue
			}
			if d.EndDate != nil && d.EndDate.Before(now) {
				continue
			}
			activeDiscount = &d
			break
		}

		var discountPerUnit float64 = 0
		if activeDiscount != nil {
			discountPerUnit, _ = model.CalculateDiscountPrice(
				unit.SellingPrice,
				string(activeDiscount.Type),
				activeDiscount.Value,
			)
		}

		qtyInBase := item.Qty * unit.ConversionToBase
		itemGrossSubtotal := item.Qty * unit.SellingPrice
		itemDiscount := discountPerUnit * item.Qty
		itemNetSubtotal := itemGrossSubtotal - itemDiscount

		subtotal += itemGrossSubtotal
		totalDiscount += itemDiscount

		preparedItems = append(preparedItems, preparedTrxItem{
			productUnit:     unit,
			product:         product,
			qtyInBase:       qtyInBase,
			qty:             item.Qty,
			unitPrice:       unit.SellingPrice,
			discountApplied: itemDiscount,
			subtotal:        itemNetSubtotal,
		})
	}

	return preparedItems, subtotal, totalDiscount, nil
}

func (u *transactionUsecaseImpl) allocateFIFOBatches(ctx context.Context, purchaseBatchRepoTx repository.PurchaseBatchRepository, product *entity.Product, qtyInBase float64) (float64, []allocationItem, error) {
	batches, err := purchaseBatchRepoTx.FindActiveBatchesByProductID(ctx, product.ID)
	if err != nil {
		return 0, nil, errs.NewInternal("failed to fetch active purchase batches: " + err.Error())
	}

	remainingQtyToAllocate := qtyInBase
	var totalCostForItem float64 = 0
	var allocations []allocationItem

	for i := range batches {
		if remainingQtyToAllocate <= 0 {
			break
		}
		batch := &batches[i]

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

	if remainingQtyToAllocate > 0 {
		return 0, nil, errs.NewConflict("insufficient purchase batch stock for product: " + product.Name)
	}

	return totalCostForItem, allocations, nil
}
