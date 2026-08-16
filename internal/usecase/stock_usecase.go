package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type StockUsecase interface {
	GetStockByProductID(ctx context.Context, productID string) (*model.StockSummaryResponse, error)
	GetStockMutations(ctx context.Context, productID string, req model.GetStockMutationsRequest) ([]model.StockMutationResponse, utils.Pagination, error)
	SubmitStockCount(ctx context.Context, userID string, req model.SubmitStockCountRequest) (*model.StockCountResponse, error)
	GetStockCounts(ctx context.Context, req model.GetStockCountsRequest) ([]model.StockCountResponse, utils.Pagination, error)
	ApproveStockCount(ctx context.Context, userID string, id string, req model.ApproveStockCountRequest) (*model.StockCountResponse, error)
}

type stockUsecaseImpl struct {
	db                *gorm.DB
	stockRepo         repository.StockRepository
	stockMutationRepo repository.StockMutationRepository
	stockCountRepo    repository.StockCountRepository
	productRepo       repository.ProductRepository
	permissionUsecase *PermissionUsecase
}

func NewStockUsecase(
	db *gorm.DB,
	stockRepo repository.StockRepository,
	stockMutationRepo repository.StockMutationRepository,
	stockCountRepo repository.StockCountRepository,
	productRepo repository.ProductRepository,
	permissionUsecase *PermissionUsecase,
) StockUsecase {
	return &stockUsecaseImpl{
		db:                db,
		stockRepo:         stockRepo,
		stockMutationRepo: stockMutationRepo,
		stockCountRepo:    stockCountRepo,
		productRepo:       productRepo,
		permissionUsecase: permissionUsecase,
	}
}

func (u *stockUsecaseImpl) GetStockByProductID(ctx context.Context, productID string) (*model.StockSummaryResponse, error) {
	product, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		return nil, errs.NewNotFound("product not found")
	}

	qty := 0.0
	updatedAt := time.Now()
	stock, err := u.stockRepo.GetByProductID(ctx, productID)
	if err == nil {
		qty = stock.QtyBaseUnit
		updatedAt = stock.UpdatedAt
	}

	return &model.StockSummaryResponse{
		ProductID:   productID,
		QtyBaseUnit: qty,
		BaseUnit: model.StockBaseUnitResponse{
			ID:   product.BaseUnit.ID,
			Name: product.BaseUnit.Name,
		},
		UpdatedAt: updatedAt,
	}, nil
}

func (u *stockUsecaseImpl) GetStockMutations(ctx context.Context, productID string, req model.GetStockMutationsRequest) ([]model.StockMutationResponse, utils.Pagination, error) {
	_, err := u.productRepo.FindById(ctx, productID)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewNotFound("product not found")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	mutations, total, err := u.stockMutationRepo.FindProductMutations(ctx, productID, req.Page, req.Limit)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch stock mutations: " + err.Error())
	}

	responses := make([]model.StockMutationResponse, 0, len(mutations))
	for _, m := range mutations {
		responses = append(responses, model.StockMutationResponse{
			ID:          m.ID,
			ProductID:   m.ProductID,
			Type:        m.Type,
			Qty:         m.Qty,
			QtyBefore:   m.QtyBefore,
			QtyAfter:    m.QtyAfter,
			Source:      m.Source,
			ReferenceID: m.ReferenceID,
			Note:        m.Note,
			Creator: model.StockMutationUserResponse{
				ID:   m.CreatedBy,
				Name: m.Creator.Name,
			},
			CreatedAt: m.CreatedAt,
		})
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return responses, pagination, nil
}

func (u *stockUsecaseImpl) SubmitStockCount(ctx context.Context, userID string, req model.SubmitStockCountRequest) (*model.StockCountResponse, error) {
	_, err := u.productRepo.FindById(ctx, req.ProductID)
	if err != nil {
		return nil, errs.NewNotFound("product not found")
	}

	hasApprovePermission, err := u.permissionUsecase.CheckUserPermission(ctx, userID, "opname:approve")
	if err != nil {
		hasApprovePermission = false
	}

	var sc *entity.StockCount
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockCountRepoTx := u.stockCountRepo.WithTx(tx)

		// 1. Get current system qty (locked for write)
		var systemQty float64 = 0
		stock, err := stockRepoTx.GetByProductID(ctx, req.ProductID)
		if err == nil {
			systemQty = stock.QtyBaseUnit
		}

		disc := req.PhysicalQty - systemQty
		note := req.Note

		sc = &entity.StockCount{
			ProductID:   req.ProductID,
			CountDate:   time.Now(),
			SystemQty:   systemQty,
			PhysicalQty: req.PhysicalQty,
			Discrepancy: disc,
			Status:      "pending",
			SubmittedBy: userID,
		}
		if note != "" {
			sc.Note = &note
		}

		if err := stockCountRepoTx.Create(ctx, sc); err != nil {
			return err
		}

		// Auto approve if submitted by user with approve permission
		if hasApprovePermission {
			err := u.approveStockCountTx(ctx, tx, userID, sc, true, "Auto-approved (Submitted by Owner/Admin)")
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// Reload with preloads for response
	reloaded, err := u.stockCountRepo.FindByID(ctx, sc.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load stock count response")
	}

	res := u.toStockCountResponse(reloaded)
	return &res, nil
}

func (u *stockUsecaseImpl) GetStockCounts(ctx context.Context, req model.GetStockCountsRequest) ([]model.StockCountResponse, utils.Pagination, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	counts, total, err := u.stockCountRepo.FindStockCounts(ctx, req.ProductID, req.Status, req.Page, req.Limit)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch stock counts: " + err.Error())
	}

	responses := make([]model.StockCountResponse, 0, len(counts))
	for _, c := range counts {
		responses = append(responses, u.toStockCountResponse(&c))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return responses, pagination, nil
}

func (u *stockUsecaseImpl) ApproveStockCount(ctx context.Context, userID string, id string, req model.ApproveStockCountRequest) (*model.StockCountResponse, error) {
	var sc *entity.StockCount
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		stockCountRepoTx := u.stockCountRepo.WithTx(tx)

		var err error
		sc, err = stockCountRepoTx.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.NewNotFound("stock count record not found")
			}
			return err
		}

		if sc.Status != "pending" {
			return errs.NewConflict("stock count has already been processed")
		}

		return u.approveStockCountTx(ctx, tx, userID, sc, req.Approve, req.Note)
	})

	if txErr != nil {
		return nil, txErr
	}

	reloaded, err := u.stockCountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewInternal("failed to reload stock count")
	}

	res := u.toStockCountResponse(reloaded)
	return &res, nil
}

// Internal helper for transaction processing
func (u *stockUsecaseImpl) approveStockCountTx(ctx context.Context, tx *gorm.DB, approverID string, sc *entity.StockCount, approve bool, note string) error {
	stockRepoTx := u.stockRepo.WithTx(tx)
	stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)
	stockCountRepoTx := u.stockCountRepo.WithTx(tx)

	now := time.Now()
	sc.ApprovedBy = &approverID
	sc.ApprovedAt = &now

	fullNote := ""
	if sc.Note != nil {
		fullNote = *sc.Note
	}
	if note != "" {
		if fullNote != "" {
			fullNote += " - Note Approval: " + note
		} else {
			fullNote = note
		}
	}
	if fullNote != "" {
		sc.Note = &fullNote
	}

	if !approve {
		sc.Status = "rejected"
		return stockCountRepoTx.Update(ctx, sc)
	}

	sc.Status = "approved"

	// Fetch current stock
	var qtyBefore float64 = 0
	currentStock, err := stockRepoTx.GetByProductID(ctx, sc.ProductID)
	if err == nil {
		qtyBefore = currentStock.QtyBaseUnit
	}

	disc := sc.Discrepancy

	if disc < 0 {
		// 1. Loss: Consume from oldest purchase batches (FIFO)
		qtyToDeduct := math.Abs(disc)
		remainingToDeduct := qtyToDeduct

		for remainingToDeduct > 0 {
			var batch entity.PurchaseBatch
			err := tx.WithContext(ctx).
				Order("purchase_date ASC, created_at ASC").
				First(&batch, "product_id = ? AND remaining_qty > 0", sc.ProductID).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// In case of physical stock mismatch but not enough purchase batches (e.g. initial setup issue),
					// we allow it but stop deducting since remaining_qty is exhausted.
					break
				}
				return err
			}

			var deduct float64
			if batch.RemainingQty >= remainingToDeduct {
				deduct = remainingToDeduct
			} else {
				deduct = batch.RemainingQty
			}

			batch.RemainingQty -= deduct
			if err := tx.Save(&batch).Error; err != nil {
				return err
			}

			remainingToDeduct -= deduct
		}

		// Update stock cache
		qtyAfter := qtyBefore + disc // disc is negative, so this subtracts
		if qtyAfter < 0 {
			qtyAfter = 0 // Safeguard
		}
		stock := &entity.Stock{
			ProductID:   sc.ProductID,
			QtyBaseUnit: qtyAfter,
		}
		if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
			return err
		}

		// Log mutation
		noteMutation := fmt.Sprintf("Penyesuaian opname fisik (selisih kurang: %f)", disc)
		mutation := &entity.StockMutation{
			ProductID:   sc.ProductID,
			Type:        "out",
			Qty:         qtyToDeduct,
			QtyBefore:   qtyBefore,
			QtyAfter:    qtyAfter,
			Source:      "stock_count",
			ReferenceID: &sc.ID,
			Note:        &noteMutation,
			CreatedBy:   approverID,
		}
		if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
			return err
		}

	} else if disc > 0 {
		// 2. Surplus: Get last purchase price and create surplus batch
		var lastPrice float64 = 0
		var latestBatch entity.PurchaseBatch
		err := tx.WithContext(ctx).
			Model(&entity.PurchaseBatch{}).
			Where("product_id = ?", sc.ProductID).
			Order("purchase_date DESC, created_at DESC").
			First(&latestBatch).Error
		
		var supplierID string
		if err == nil {
			lastPrice = latestBatch.PurchasePrice
			supplierID = latestBatch.SupplierID
		} else {
			// Find first active supplier
			var supplier entity.Supplier
			errSupplier := tx.WithContext(ctx).
				Model(&entity.Supplier{}).
				Where("is_active = ?", true).
				First(&supplier).Error
			if errSupplier != nil {
				return errs.NewConflict("no active supplier found to attribute the surplus batch to")
			}
			supplierID = supplier.ID
		}

		invNum := "OPNAME-SURPLUS-" + sc.ID
		surplusBatch := &entity.PurchaseBatch{
			ProductID:     sc.ProductID,
			SupplierID:    supplierID,
			InitialQty:    disc,
			RemainingQty:  disc,
			PurchasePrice: lastPrice,
			InvoiceNumber: &invNum,
			PurchaseDate:  now,
			CreatedBy:     approverID,
		}

		if err := tx.Create(surplusBatch).Error; err != nil {
			return err
		}

		// Update stock cache
		qtyAfter := qtyBefore + disc
		stock := &entity.Stock{
			ProductID:   sc.ProductID,
			QtyBaseUnit: qtyAfter,
		}
		if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
			return err
		}

		// Log mutation
		noteMutation := fmt.Sprintf("Penyesuaian opname fisik (selisih lebih: +%f)", disc)
		mutation := &entity.StockMutation{
			ProductID:   sc.ProductID,
			Type:        "in",
			Qty:         disc,
			QtyBefore:   qtyBefore,
			QtyAfter:    qtyAfter,
			Source:      "stock_count",
			ReferenceID: &sc.ID,
			Note:        &noteMutation,
			CreatedBy:   approverID,
		}
		if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
			return err
		}
	}

	return stockCountRepoTx.Update(ctx, sc)
}

func (u *stockUsecaseImpl) toStockCountResponse(sc *entity.StockCount) model.StockCountResponse {
	res := model.StockCountResponse{
		ID:          sc.ID,
		ProductID:   sc.ProductID,
		Product: model.StockCountProductResponse{
			ID:   sc.Product.ID,
			Name: sc.Product.Name,
		},
		CountDate:   sc.CountDate,
		SystemQty:   sc.SystemQty,
		PhysicalQty: sc.PhysicalQty,
		Discrepancy: sc.Discrepancy,
		Note:        sc.Note,
		Status:      sc.Status,
		SubmittedBy: sc.SubmittedBy,
		Submitter: model.StockCountUserResponse{
			ID:   sc.Submitter.ID,
			Name: sc.Submitter.Name,
		},
		SubmittedAt: sc.SubmittedAt,
	}

	if sc.ApprovedBy != nil && sc.Approver != nil {
		res.ApprovedBy = sc.ApprovedBy
		res.Approver = &model.StockCountUserResponse{
			ID:   sc.Approver.ID,
			Name: sc.Approver.Name,
		}
		res.ApprovedAt = sc.ApprovedAt
	}

	return res
}
