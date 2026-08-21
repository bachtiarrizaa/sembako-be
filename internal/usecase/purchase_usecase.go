package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"gorm.io/gorm"
)

type PurchaseUsecase interface {
	CreatePurchase(ctx context.Context, creatorID string, req model.CreatePurchaseRequest) (*model.PurchaseDetailResponse, error)
	GetPurchases(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]model.PurchaseSummaryResponse, utils.Pagination, error)
	GetPurchaseDetail(ctx context.Context, id string) (*model.PurchaseDetailResponse, error)
	UpdatePurchaseItem(ctx context.Context, updaterID string, id string, req model.UpdatePurchaseRequest) (*model.PurchaseBatchResponse, error)
	DeletePurchase(ctx context.Context, deleterID string, id string) error
	DeletePurchaseItem(ctx context.Context, deleterID string, id string) error
}

type purchaseUsecaseImpl struct {
	db                *gorm.DB
	purchaseRepo      repository.PurchaseRepository
	purchaseBatchRepo repository.PurchaseBatchRepository
	productRepo       repository.ProductRepository
	supplierRepo      repository.SupplierRepository
	stockRepo         repository.StockRepository
	stockMutationRepo repository.StockMutationRepository
}

func NewPurchaseUsecase(
	db *gorm.DB,
	purchaseRepo repository.PurchaseRepository,
	purchaseBatchRepo repository.PurchaseBatchRepository,
	productRepo repository.ProductRepository,
	supplierRepo repository.SupplierRepository,
	stockRepo repository.StockRepository,
	stockMutationRepo repository.StockMutationRepository,
) PurchaseUsecase {
	return &purchaseUsecaseImpl{
		db:                db,
		purchaseRepo:      purchaseRepo,
		purchaseBatchRepo: purchaseBatchRepo,
		productRepo:       productRepo,
		supplierRepo:      supplierRepo,
		stockRepo:         stockRepo,
		stockMutationRepo: stockMutationRepo,
	}
}

type preparedPurchaseItem struct {
	product      *entity.Product
	matchedUnit  *entity.ProductUnit
	qtyInBase    float64
	pricePerBase float64
}

func (u *purchaseUsecaseImpl) prepareItem(ctx context.Context, item model.CreatePurchaseItem) (*preparedPurchaseItem, error) {
	product, err := u.productRepo.FindById(ctx, item.ProductID)
	if err != nil {
		return nil, errs.NewNotFound("product not found: " + item.ProductID)
	}
	if !product.IsActive {
		return nil, errs.NewConflict("cannot purchase an inactive product: " + product.Name)
	}

	var matchedUnit *entity.ProductUnit
	for i := range product.Units {
		if product.Units[i].UnitID == item.UnitID {
			matchedUnit = &product.Units[i]
			break
		}
	}
	if matchedUnit == nil {
		return nil, errs.NewNotFound("product unit not found for the given product")
	}
	if !matchedUnit.IsActive {
		return nil, errs.NewConflict("cannot purchase with an inactive product unit")
	}

	qtyInBase := item.Quantity * matchedUnit.ConversionToBase
	pricePerBase := item.PurchasePrice / matchedUnit.ConversionToBase

	return &preparedPurchaseItem{
		product:      product,
		matchedUnit:  matchedUnit,
		qtyInBase:    qtyInBase,
		pricePerBase: pricePerBase,
	}, nil
}

func (u *purchaseUsecaseImpl) CreatePurchase(ctx context.Context, creatorID string, req model.CreatePurchaseRequest) (*model.PurchaseDetailResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return nil, errs.NewBadRequest("invalid purchase date format, must be YYYY-MM-DD")
	}

	// Validate supplier
	supplier, err := u.supplierRepo.FindById(ctx, req.SupplierID)
	if err != nil {
		return nil, errs.NewNotFound("supplier not found")
	}
	if !supplier.IsActive {
		return nil, errs.NewConflict("cannot record purchase from an inactive supplier")
	}

	// Pre-validate items and compute the invoice total
	prepared := make([]preparedPurchaseItem, 0, len(req.Items))
	var totalAmount float64 = 0
	for _, item := range req.Items {
		pi, err := u.prepareItem(ctx, item)
		if err != nil {
			return nil, err
		}
		totalAmount += pi.qtyInBase * pi.pricePerBase
		prepared = append(prepared, *pi)
	}

	var purchaseID string

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseRepoTx := u.purchaseRepo.WithTx(tx)
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		productRepoTx := u.productRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)

		purchase := &entity.Purchase{
			InvoiceNumber: req.InvoiceNumber,
			SupplierID:    req.SupplierID,
			PurchaseDate:  parsedDate,
			TotalAmount:   totalAmount,
			CreatedBy:     creatorID,
		}
		if err := purchaseRepoTx.Create(ctx, purchase); err != nil {
			return err
		}
		purchaseID = purchase.ID

		for i, item := range req.Items {
			pi := prepared[i]

			// Re-check product within the transaction to keep FIFO integrity
			product, err := productRepoTx.FindById(ctx, item.ProductID)
			if err != nil {
				return errs.NewNotFound("product not found: " + item.ProductID)
			}
			if !product.IsActive {
				return errs.NewConflict("cannot purchase an inactive product: " + product.Name)
			}

			batch := &entity.PurchaseBatch{
				PurchaseID:    &purchase.ID,
				ProductID:     item.ProductID,
				SupplierID:    req.SupplierID,
				UnitID:        &pi.matchedUnit.ID,
				UnitPrice:     &item.PurchasePrice,
				InitialQty:    pi.qtyInBase,
				RemainingQty:  pi.qtyInBase,
				PurchasePrice: pi.pricePerBase,
				InvoiceNumber: req.InvoiceNumber,
				PurchaseDate:  parsedDate,
				CreatedBy:     creatorID,
			}

			if err := purchaseBatchRepoTx.Create(ctx, batch); err != nil {
				return err
			}

			// 1. Get current stock
			var qtyBefore float64 = 0
			currentStock, err := stockRepoTx.GetByProductID(ctx, item.ProductID)
			if err == nil {
				qtyBefore = currentStock.QtyBaseUnit
			}

			// 2. Update stock cache
			qtyAfter := qtyBefore + pi.qtyInBase
			stock := &entity.Stock{
				ProductID:   item.ProductID,
				QtyBaseUnit: qtyAfter,
			}
			if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
				return err
			}

			// 3. Log stock mutation
			noteStr := ""
			if req.InvoiceNumber != nil {
				noteStr = "Invoice: " + *req.InvoiceNumber
			}
			mutation := &entity.StockMutation{
				ProductID:   item.ProductID,
				Type:        "in",
				Qty:         pi.qtyInBase,
				QtyBefore:   qtyBefore,
				QtyAfter:    qtyAfter,
				Source:      "purchase",
				ReferenceID: &batch.ID,
				CreatedBy:   creatorID,
			}
			if noteStr != "" {
				mutation.Note = &noteStr
			}
			if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return u.GetPurchaseDetail(ctx, purchaseID)
}

func (u *purchaseUsecaseImpl) GetPurchases(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]model.PurchaseSummaryResponse, utils.Pagination, error) {
	purchases, productNames, total, err := u.purchaseRepo.FindPurchases(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch purchases: " + err.Error())
	}

	responses := make([]model.PurchaseSummaryResponse, 0, len(purchases))
	for _, p := range purchases {
		responses = append(responses, model.ToPurchaseSummaryResponse(&p, productNames[p.ID]))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return responses, pagination, nil
}

func (u *purchaseUsecaseImpl) GetPurchaseDetail(ctx context.Context, id string) (*model.PurchaseDetailResponse, error) {
	purchase, err := u.purchaseRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("purchase not found")
		}
		return nil, errs.NewInternal("failed to fetch purchase: " + err.Error())
	}

	batches, err := u.purchaseBatchRepo.FindByPurchaseID(ctx, id)
	if err != nil {
		return nil, errs.NewInternal("failed to fetch purchase items: " + err.Error())
	}

	items := make([]model.PurchaseBatchResponse, 0, len(batches))
	for _, b := range batches {
		items = append(items, model.ToPurchaseBatchResponse(&b))
	}

	return &model.PurchaseDetailResponse{
		ID:            purchase.ID,
		InvoiceNumber: purchase.InvoiceNumber,
		PurchaseDate:  purchase.PurchaseDate,
		Supplier: model.PurchaseSupplierResponse{
			ID:   purchase.SupplierID,
			Name: purchase.Supplier.Name,
		},
		TotalAmount: purchase.TotalAmount,
		Creator: model.PurchaseCreatorResponse{
			ID:   purchase.CreatedBy,
			Name: purchase.Creator.Name,
		},
		CreatedAt: purchase.CreatedAt,
		Items:     items,
	}, nil
}

func (u *purchaseUsecaseImpl) UpdatePurchaseItem(ctx context.Context, updaterID string, id string, req model.UpdatePurchaseRequest) (*model.PurchaseBatchResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return nil, errs.NewBadRequest("invalid purchase date format, must be YYYY-MM-DD")
	}

	// Validate supplier
	supplier, err := u.supplierRepo.FindById(ctx, req.SupplierID)
	if err != nil {
		return nil, errs.NewNotFound("supplier not found")
	}
	if !supplier.IsActive {
		return nil, errs.NewConflict("cannot update purchase with an inactive supplier")
	}

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		purchaseRepoTx := u.purchaseRepo.WithTx(tx)
		productRepoTx := u.productRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)

		// Fetch existing batch
		batch, err := purchaseBatchRepoTx.FindByID(ctx, id)
		if err != nil {
			return errs.NewNotFound("purchase batch not found")
		}

		// Find product and unit conversion
		product, err := productRepoTx.FindById(ctx, batch.ProductID)
		if err != nil {
			return errs.NewNotFound("product not found")
		}

		var matchedUnit *entity.ProductUnit
		for i := range product.Units {
			if product.Units[i].UnitID == req.UnitID {
				matchedUnit = &product.Units[i]
				break
			}
		}
		if matchedUnit == nil {
			return errs.NewNotFound("product unit not found for the product")
		}
		if !matchedUnit.IsActive {
			return errs.NewConflict("cannot update with an inactive product unit")
		}

		newQtyInBase := req.Quantity * matchedUnit.ConversionToBase
		newPricePerBase := req.PurchasePrice / matchedUnit.ConversionToBase

		hasBeenSold := batch.RemainingQty < batch.InitialQty

		delta := newQtyInBase - batch.InitialQty

		if hasBeenSold {
			// If already sold, we cannot change product, quantity, or price!
			if newQtyInBase != batch.InitialQty || newPricePerBase != batch.PurchasePrice {
				return errs.NewConflict("cannot update quantity or price of a purchase batch that has been partially sold")
			}
			// Only update invoice number, date, and supplier
			batch.SupplierID = req.SupplierID
			batch.InvoiceNumber = req.InvoiceNumber
			batch.PurchaseDate = parsedDate
		} else {
			// Not sold, can update everything
			batch.SupplierID = req.SupplierID
			batch.InvoiceNumber = req.InvoiceNumber
			batch.PurchaseDate = parsedDate
			batch.UnitID = &matchedUnit.ID
			batch.UnitPrice = &req.PurchasePrice
			batch.InitialQty = newQtyInBase
			batch.RemainingQty = newQtyInBase
			batch.PurchasePrice = newPricePerBase

			if delta != 0 {
				var qtyBefore float64 = 0
				currentStock, err := stockRepoTx.GetByProductID(ctx, batch.ProductID)
				if err == nil {
					qtyBefore = currentStock.QtyBaseUnit
				}

				qtyAfter := qtyBefore + delta
				stock := &entity.Stock{
					ProductID:   batch.ProductID,
					QtyBaseUnit: qtyAfter,
				}
				if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
					return err
				}

				mType := "in"
				mQty := delta
				if delta < 0 {
					mType = "out"
					mQty = -delta
				}
				noteStr := "purchase update adjustment"
				if req.InvoiceNumber != nil {
					noteStr += " (Invoice: " + *req.InvoiceNumber + ")"
				}
				mutation := &entity.StockMutation{
					ProductID:   batch.ProductID,
					Type:        mType,
					Qty:         mQty,
					QtyBefore:   qtyBefore,
					QtyAfter:    qtyAfter,
					Source:      "purchase",
					ReferenceID: &batch.ID,
					Note:        &noteStr,
					CreatedBy:   updaterID,
				}
				if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
					return err
				}
			}
		}

		if err := purchaseBatchRepoTx.Update(ctx, batch); err != nil {
			return err
		}

		// Keep the parent header in sync (supplier/invoice/date + recompute total)
		if batch.PurchaseID != nil {
			purchase, err := purchaseRepoTx.FindByID(ctx, *batch.PurchaseID)
			if err != nil {
				return errs.NewInternal("purchase header not found")
			}
			purchase.SupplierID = req.SupplierID
			purchase.InvoiceNumber = req.InvoiceNumber
			purchase.PurchaseDate = parsedDate

			allBatches, err := purchaseBatchRepoTx.FindByPurchaseID(ctx, purchase.ID)
			if err != nil {
				return err
			}
			var total float64
			for _, b := range allBatches {
				total += b.InitialQty * b.PurchasePrice
			}
			purchase.TotalAmount = total
			if err := purchaseRepoTx.Update(ctx, purchase); err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	updated, err := u.purchaseBatchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated purchase batch")
	}

	res := model.ToPurchaseBatchResponse(updated)
	return &res, nil
}

func (u *purchaseUsecaseImpl) DeletePurchase(ctx context.Context, deleterID string, id string) error {
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseRepoTx := u.purchaseRepo.WithTx(tx)
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)

		// Fetch header
		purchase, err := purchaseRepoTx.FindByID(ctx, id)
		if err != nil {
			return errs.NewNotFound("purchase not found")
		}

		// Fetch all its batches
		batches, err := purchaseBatchRepoTx.FindByPurchaseID(ctx, id)
		if err != nil {
			return errs.NewInternal("failed to fetch purchase items: " + err.Error())
		}

		for _, b := range batches {
			if b.RemainingQty < b.InitialQty {
				return errs.NewConflict("cannot delete a purchase that has partially sold items")
			}
		}

		for i := range batches {
			batch := &batches[i]

			// 1. Get current stock
			var qtyBefore float64 = 0
			currentStock, err := stockRepoTx.GetByProductID(ctx, batch.ProductID)
			if err == nil {
				qtyBefore = currentStock.QtyBaseUnit
			}

			// 2. Update stock cache
			qtyAfter := qtyBefore - batch.InitialQty
			stock := &entity.Stock{
				ProductID:   batch.ProductID,
				QtyBaseUnit: qtyAfter,
			}
			if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
				return err
			}

			// 3. Log stock mutation
			noteStr := "purchase deleted"
			if batch.InvoiceNumber != nil {
				noteStr += " (Invoice: " + *batch.InvoiceNumber + ")"
			}
			mutation := &entity.StockMutation{
				ProductID:   batch.ProductID,
				Type:        "out",
				Qty:         batch.InitialQty,
				QtyBefore:   qtyBefore,
				QtyAfter:    qtyAfter,
				Source:      "purchase",
				ReferenceID: &batch.ID,
				Note:        &noteStr,
				CreatedBy:   deleterID,
			}
			if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
				return err
			}

			if err := purchaseBatchRepoTx.Delete(ctx, batch.ID); err != nil {
				return err
			}
		}

		return purchaseRepoTx.Delete(ctx, purchase.ID)
	})

	return txErr
}

func (u *purchaseUsecaseImpl) DeletePurchaseItem(ctx context.Context, deleterID string, id string) error {
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		purchaseRepoTx := u.purchaseRepo.WithTx(tx)
		stockRepoTx := u.stockRepo.WithTx(tx)
		stockMutationRepoTx := u.stockMutationRepo.WithTx(tx)

		// Fetch batch
		batch, err := purchaseBatchRepoTx.FindByID(ctx, id)
		if err != nil {
			return errs.NewNotFound("purchase batch not found")
		}

		if batch.RemainingQty < batch.InitialQty {
			return errs.NewConflict("cannot delete a purchase item that has been partially sold")
		}

		// 1. Get current stock
		var qtyBefore float64 = 0
		currentStock, err := stockRepoTx.GetByProductID(ctx, batch.ProductID)
		if err == nil {
			qtyBefore = currentStock.QtyBaseUnit
		}

		// 2. Update stock cache
		qtyAfter := qtyBefore - batch.InitialQty
		stock := &entity.Stock{
			ProductID:   batch.ProductID,
			QtyBaseUnit: qtyAfter,
		}
		if err := stockRepoTx.UpsertStock(ctx, stock); err != nil {
			return err
		}

		// 3. Log stock mutation
		noteStr := "purchase item deleted"
		if batch.InvoiceNumber != nil {
			noteStr += " (Invoice: " + *batch.InvoiceNumber + ")"
		}
		mutation := &entity.StockMutation{
			ProductID:   batch.ProductID,
			Type:        "out",
			Qty:         batch.InitialQty,
			QtyBefore:   qtyBefore,
			QtyAfter:    qtyAfter,
			Source:      "purchase",
			ReferenceID: &batch.ID,
			Note:        &noteStr,
			CreatedBy:   deleterID,
		}
		if err := stockMutationRepoTx.Create(ctx, mutation); err != nil {
			return err
		}

		if err := purchaseBatchRepoTx.Delete(ctx, batch.ID); err != nil {
			return err
		}

		// Recompute header total; delete header when the last item is removed
		if batch.PurchaseID != nil {
			purchase, err := purchaseRepoTx.FindByID(ctx, *batch.PurchaseID)
			if err != nil {
				return errs.NewInternal("purchase header not found")
			}
			allBatches, err := purchaseBatchRepoTx.FindByPurchaseID(ctx, purchase.ID)
			if err != nil {
				return err
			}
			if len(allBatches) == 0 {
				return purchaseRepoTx.Delete(ctx, purchase.ID)
			}
			var total float64
			for _, b := range allBatches {
				total += b.InitialQty * b.PurchasePrice
			}
			purchase.TotalAmount = total
			return purchaseRepoTx.Update(ctx, purchase)
		}

		return nil
	})

	return txErr
}