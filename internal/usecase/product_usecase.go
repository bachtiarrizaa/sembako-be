package usecase

import (
	"context"
	"errors"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type ProductUsecase struct {
	db              *gorm.DB
	productRepo     repository.ProductRepository
	productUnitRepo repository.ProductUnitRepository
}

func NewProductUsecase(
	db *gorm.DB,
	productRepo repository.ProductRepository,
	productUnitRepo repository.ProductUnitRepository,
) *ProductUsecase {
	return &ProductUsecase{
		db:              db,
		productRepo:     productRepo,
		productUnitRepo: productUnitRepo,
	}
}

func (u *ProductUsecase) CreateProduct(
	ctx context.Context,
	req model.CreateProductRequest,
	imagePath *string,
) (*model.ProductResponse, error) {
	if err := validateProductUnits(req.Units); err != nil {
		return nil, err
	}

	var baseUnitID string
	for _, ur := range req.Units {
		if ur.IsBaseUnit {
			baseUnitID = ur.UnitID
			break
		}
	}

	existing, err := u.productRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing product")
	}
	if existing != nil {
		return nil, errs.NewConflict("product name already exists")
	}

	var product *entity.Product

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		productRepoTx := u.productRepo.WithTx(tx)
		productUnitRepoTx := u.productUnitRepo.WithTx(tx)

		product = &entity.Product{
			CategoryID:             req.CategoryID,
			Name:                   req.Name,
			Image:                  imagePath,
			BaseUnitID:             baseUnitID,
			MinimumStock:           req.MinimumStock,
			MarginThresholdPercent: req.MarginThresholdPercent,
			IsActive:               true,
		}
		if err := productRepoTx.Create(ctx, product); err != nil {
			return err
		}

		units := make([]entity.ProductUnit, 0, len(req.Units))
		for _, ur := range req.Units {
			units = append(units, entity.ProductUnit{
				ProductID:        product.ID,
				UnitID:           ur.UnitID,
				ConversionToBase: ur.ConversionToBase,
				SellingPrice:     ur.SellingPrice,
				IsBaseUnit:       ur.IsBaseUnit,
				IsActive:         true,
			})
		}
		if err := productUnitRepoTx.CreateMany(ctx, units); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(txErr, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.NewConflict("product unit has duplicate unit assignment")
		}
		return nil, errs.NewInternal("failed to create product")
	}

	created, err := u.productRepo.FindById(ctx, product.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load created product")
	}

	res := model.ToProductResponse(created)
	return &res, nil
}

func (u *ProductUsecase) GetProducts(ctx context.Context, req model.GetProductsRequest) ([]model.ProductResponse, utils.Pagination, error) {
	products, total, err := u.productRepo.FindProducts(ctx, req)
	if err != nil {
		return nil, utils.Pagination{}, errs.NewInternal("failed to fetch products")
	}

	res := []model.ProductResponse{}
	for _, p := range products {
		res = append(res, model.ToProductResponse(&p))
	}

	pagination := utils.BuildPagination(req.Page, req.Limit, total)
	return res, pagination, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id string) (*model.ProductResponse, error) {
	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product not found")
		}
		return nil, errs.NewInternal("failed to fetch product")
	}
	res := model.ToProductResponse(product)
	return &res, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id string, req model.UpdateProductRequest, imagePath *string) (*model.ProductResponse, error) {
	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product not found")
		}
		return nil, errs.NewInternal("failed to fetch product")
	}

	existing, err := u.productRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewInternal("failed to check existing product")
	}
	if existing != nil && existing.ID != id {
		return nil, errs.NewConflict("product name already exists")
	}

	var oldImage *string
	if imagePath != nil {
		if product.Image != nil {
			oldImage = product.Image
		}
		product.Image = imagePath
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.MinimumStock = req.MinimumStock
	product.MarginThresholdPercent = req.MarginThresholdPercent

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, errs.NewInternal("failed to update product")
	}

	if oldImage != nil {
		utils.DeleteFile(*oldImage)
	}

	updated, err := u.productRepo.FindById(ctx, product.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated product")
	}

	res := model.ToProductResponse(updated)
	return &res, nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFound("product not found")
		}
		return errs.NewInternal("failed to fetch product")
	}

	// Cek referensi: produk yang sudah dipakai di transaksi/pembelian/stok tidak boleh dihapus,
	// hanya bisa dinonaktifkan via PATCH /:id/status.
	hasTransactions, err := u.productRepo.HasTransactionReferences(ctx, id)
	if err != nil {
		return errs.NewInternal("failed to check transaction references")
	}
	if hasTransactions {
		return errs.NewConflict("produk sudah digunakan dalam transaksi, gunakan deactivate")
	}

	hasPurchases, err := u.productRepo.HasPurchaseReferences(ctx, id)
	if err != nil {
		return errs.NewInternal("failed to check purchase references")
	}
	if hasPurchases {
		return errs.NewConflict("produk sudah digunakan dalam pembelian, gunakan deactivate")
	}

	hasMutations, err := u.productRepo.HasStockMutationReferences(ctx, id)
	if err != nil {
		return errs.NewInternal("failed to check stock mutation references")
	}
	if hasMutations {
		return errs.NewConflict("produk sudah memiliki riwayat stok, gunakan deactivate")
	}

	if err := u.productRepo.Delete(ctx, product.ID); err != nil {
		return errs.NewInternal("failed to delete product")
	}

	return nil
}

func validateProductUnits(units []model.CreateProductUnitRequest) error {
	baseUnitCount := 0
	seenUnitID := make(map[string]bool)

	for _, u := range units {
		if seenUnitID[u.UnitID] {
			return errs.NewBadRequest("duplicate unit in product units")
		}
		seenUnitID[u.UnitID] = true

		if u.IsBaseUnit {
			baseUnitCount++
			if u.ConversionToBase != 1 {
				return errs.NewBadRequest("base unit must have conversionToBase = 1")
			}
		}
	}

	if baseUnitCount == 0 {
		return errs.NewBadRequest("exactly one unit must be marked as base unit")
	}
	if baseUnitCount > 1 {
		return errs.NewBadRequest("only one unit can be marked as base unit")
	}

	return nil
}

func (u *ProductUsecase) UpdateProductStatus(ctx context.Context, id string) (*model.ProductResponse, error) {
	product, err := u.productRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("product not found")
		}
		return nil, errs.NewInternal("failed to fetch product")
	}

	product.IsActive = !product.IsActive

	txErr := u.db.Transaction(func(tx *gorm.DB) error {
		productRepoTx := u.productRepo.WithTx(tx)
		productUnitRepoTx := u.productUnitRepo.WithTx(tx)

		if err := productRepoTx.Update(ctx, product); err != nil {
			return err
		}

		if !product.IsActive {
			if err := productUnitRepoTx.DeactivateAllByProductID(ctx, product.ID); err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, errs.NewInternal("failed to update product status")
	}

	updated, err := u.productRepo.FindById(ctx, product.ID)
	if err != nil {
		return nil, errs.NewInternal("failed to load updated product")
	}

	res := model.ToProductResponse(updated)
	return &res, nil
}
