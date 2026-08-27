package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"gorm.io/gorm"
)

func (u *transactionUsecaseImpl) VoidTransaction(ctx context.Context, id string, cashierID string, req model.VoidTransactionRequest) (*model.TransactionResponse, error) {
	// 1. Ambil data transaksi beserta item & satuan produknya
	trx, err := u.transactionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("transaction not found")
		}
		return nil, errs.NewInternal("failed to fetch transaction")
	}

	// 2. Validasi: Hanya transaksi 'completed' yang boleh di-void
	if trx.Status != "completed" {
		return nil, errs.NewConflict("only completed transactions can be voided")
	}

	// 3. Validasi: Shift transaksi harus masih 'open'
	if trx.Shift.Status != "open" {
		return nil, errs.NewConflict("cannot void transaction from a closed shift")
	}

	// 4. Eksekusi Pembatalan Secara Transaksional (Atomic DB Transaction)
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := u.transactionRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		costAllocationRepoTx := u.costAllocationRepo.WithTx(tx)

		for _, item := range trx.Items {
			qtyInBase := item.Qty * item.ProductUnit.ConversionToBase

			// A. Ambil stok saat ini dan kembalikan (tambah kembali)
			currentStock, err := stockRepoTx.GetByProductID(ctx, item.ProductUnit.ProductID)
			if err != nil {
				return errs.NewInternal("failed to fetch stock for product: " + item.ProductUnit.Product.Name)
			}

			qtyBefore := currentStock.QtyBaseUnit
			qtyAfter := qtyBefore + qtyInBase
			currentStock.QtyBaseUnit = qtyAfter

			if err := stockRepoTx.UpsertStock(ctx, currentStock); err != nil {
				return errs.NewInternal("failed to restore stock: " + err.Error())
			}

			// B. Catat mutasi stok masuk (IN) karena transaksi dibatalkan
			noteStr := "Void Transaction: " + trx.ReceiptNumber
			mutation := &entity.StockMutation{
				ProductID:   item.ProductUnit.ProductID,
				Type:        "in",
				Qty:         qtyInBase,
				QtyBefore:   qtyBefore,
				QtyAfter:    qtyAfter,
				Source:      "sale",
				ReferenceID: &trx.ID,
				Note:        &noteStr,
				CreatedBy:   cashierID,
			}
			if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
				return errs.NewInternal("failed to create void stock mutation: " + err.Error())
			}

			// C. Kembalikan sisa kuantitas pada batch pembelian FIFO
			allocations, err := costAllocationRepoTx.FindByTransactionItemID(ctx, item.ID)
			if err != nil {
				return errs.NewInternal("failed to fetch cost allocations: " + err.Error())
			}

			for _, alloc := range allocations {
				batch, err := purchaseBatchRepoTx.FindByID(ctx, alloc.PurchaseBatchID)
				if err != nil {
					return errs.NewInternal("failed to fetch purchase batch to restore: " + err.Error())
				}
				batch.RemainingQty = batch.RemainingQty + alloc.QtyAllocated
				if err := purchaseBatchRepoTx.Update(ctx, batch); err != nil {
					return errs.NewInternal("failed to restore purchase batch: " + err.Error())
				}
			}
		}

		// D. Update status transaksi menjadi 'void'
		now := time.Now()
		trx.Status = "void"
		trx.VoidReason = &req.Reason
		trx.VoidedBy = &cashierID
		trx.VoidedAt = &now

		if err := transactionRepoTx.Update(ctx, trx); err != nil {
			return errs.NewInternal("failed to update transaction status to voided: " + err.Error())
		}

		return nil
	})

	if txErr != nil {
		if appErr, ok := txErr.(*errs.AppError); ok {
			return nil, appErr
		}
		return nil, errs.NewInternal("failed to void transaction: " + txErr.Error())
	}

	// 5. Muat ulang transaksi yang sudah di-void untuk response
	updatedTrx, err := u.transactionRepo.FindByID(ctx, trx.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to fetch updated voided transaction")
	}

	res := model.ToTransactionResponse(updatedTrx)
	return &res, nil
}
