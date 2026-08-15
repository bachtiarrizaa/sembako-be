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
	CreatePurchase(ctx context.Context, creatorID string, req model.CreatePurchaseRequest) ([]model.PurchaseBatchResponse, error)
	GetPurchaseBatches(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]model.PurchaseBatchResponse, utils.Pagination, error)
	GetPurchaseBatchByID(ctx context.Context, id string) (*model.PurchaseBatchResponse, error)
	UpdatePurchase(ctx context.Context, id string, req model.UpdatePurchaseRequest) (*model.PurchaseBatchResponse, error)
	DeletePurchase(ctx context.Context, id string) error
}

type purchaseUsecaseImpl struct {
	db                *gorm.DB
	purchaseBatchRepo repository.PurchaseBatchRepository
	productRepo       repository.ProductRepository
	supplierRepo      repository.SupplierRepository
}

func NewPurchaseUsecase(
	db *gorm.DB,
	purchaseBatchRepo repository.PurchaseBatchRepository,
	productRepo repository.ProductRepository,
	supplierRepo repository.SupplierRepository,
) PurchaseUsecase {
	return &purchaseUsecaseImpl{
		db:                db,
		purchaseBatchRepo: purchaseBatchRepo,
		productRepo:       productRepo,
		supplierRepo:      supplierRepo,
	}
}

func (u *purchaseUsecaseImpl) CreatePurchase(ctx context.Context, creatorID string, req model.CreatePurchaseRequest) ([]model.PurchaseBatchResponse, error) {
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

	var createdIDs []string

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)
		productRepoTx := u.productRepo.WithTx(tx)

		for _, item := range req.Items {
			// Find product
			product, err := productRepoTx.FindById(ctx, item.ProductID)
			if err != nil {
				return errs.NewNotFound("product not found: " + item.ProductID)
			}
			if !product.IsActive {
				return errs.NewConflict("cannot purchase an inactive product: " + product.Name)
			}

			// Find matching unit
			var matchedUnit *entity.ProductUnit
			for _, pu := range product.Units {
				if pu.UnitID == item.UnitID {
					matchedUnit = &pu
					break
				}
			}
			if matchedUnit == nil {
				return errs.NewNotFound("product unit not found for the given product")
			}
			if !matchedUnit.IsActive {
				return errs.NewConflict("cannot purchase with an inactive product unit")
			}

			// Conversion logic
			qtyInBase := item.Quantity * matchedUnit.ConversionToBase
			pricePerBase := item.PurchasePrice / matchedUnit.ConversionToBase

			batch := &entity.PurchaseBatch{
				ProductID:     item.ProductID,
				SupplierID:    req.SupplierID,
				UnitID:        &matchedUnit.ID,
				UnitPrice:     &item.PurchasePrice,
				InitialQty:    qtyInBase,
				RemainingQty:  qtyInBase,
				PurchasePrice: pricePerBase,
				InvoiceNumber: req.InvoiceNumber,
				PurchaseDate:  parsedDate,
				CreatedBy:     creatorID,
			}

			if err := purchaseBatchRepoTx.Create(ctx, batch); err != nil {
				return err
			}
			createdIDs = append(createdIDs, batch.ID)
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	responses := make([]model.PurchaseBatchResponse, 0, len(createdIDs))
	for _, id := range createdIDs {
		b, err := u.purchaseBatchRepo.FindByID(ctx, id)
		if err != nil {
			return nil, errs.NewInternal("failed to load created purchase batch")
		}
		responses = append(responses, model.ToPurchaseBatchResponse(b))
	}

	return responses, nil
}

func (u *purchaseUsecaseImpl) GetPurchaseBatches(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]model.PurchaseBatchResponse, utils.Pagination, error) {
	batches, total, err := u.purchaseBatchRepo.FindPurchaseBatches(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch purchase batches: " + err.Error())
	}

	responses := make([]model.PurchaseBatchResponse, 0, len(batches))
	for _, b := range batches {
		responses = append(responses, model.PurchaseBatchResponse{
			ID: b.ID,
			Product: model.PurchaseProductResponse{
				ID:   b.ProductID,
				Name: b.Product.Name,
			},
			Supplier: model.PurchaseSupplierResponse{
				ID:   b.SupplierID,
				Name: b.Supplier.Name,
			},
			Unit:              model.ToPurchaseUnitResponse(&b),
			UnitPrice:         b.UnitPrice,
			BaseUnit:          model.ToBaseUnitResponse(&b.Product.BaseUnit),
			InitialQuantity:   b.InitialQty,
			RemainingQuantity: b.RemainingQty,
			PurchasePrice:     b.PurchasePrice,
			InvoiceNumber:     b.InvoiceNumber,
			PurchaseDate:      b.PurchaseDate,
			Creator: model.PurchaseCreatorResponse{
				ID:   b.CreatedBy,
				Name: b.Creator.Name,
			},
			CreatedAt: b.CreatedAt,
		})
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return responses, pagination, nil
}

func (u *purchaseUsecaseImpl) GetPurchaseBatchByID(ctx context.Context, id string) (*model.PurchaseBatchResponse, error) {
	batch, err := u.purchaseBatchRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("purchase batch not found")
		}
		return nil, errs.NewInternal("failed to fetch purchase batch: " + err.Error())
	}
	res := model.ToPurchaseBatchResponse(batch)
	return &res, nil
}

func (u *purchaseUsecaseImpl) UpdatePurchase(ctx context.Context, id string, req model.UpdatePurchaseRequest) (*model.PurchaseBatchResponse, error) {
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
		productRepoTx := u.productRepo.WithTx(tx)

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
		for _, pu := range product.Units {
			if pu.UnitID == req.UnitID {
				matchedUnit = &pu
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
		}

		return purchaseBatchRepoTx.Update(ctx, batch)
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

func (u *purchaseUsecaseImpl) DeletePurchase(ctx context.Context, id string) error {
	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		purchaseBatchRepoTx := u.purchaseBatchRepo.WithTx(tx)

		// Fetch existing batch
		batch, err := purchaseBatchRepoTx.FindByID(ctx, id)
		if err != nil {
			return errs.NewNotFound("purchase batch not found")
		}

		hasBeenSold := batch.RemainingQty < batch.InitialQty
		if hasBeenSold {
			return errs.NewConflict("cannot delete a purchase batch that has been partially sold")
		}

		return purchaseBatchRepoTx.Delete(ctx, id)
	})

	return txErr
}
