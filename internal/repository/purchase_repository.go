package repository

import (
	"context"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entity.Purchase) error
	FindPurchases(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]entity.Purchase, map[string][]string, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Purchase, error)
	Update(ctx context.Context, purchase *entity.Purchase) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) PurchaseRepository
}

type purchaseRepositoryImpl struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepositoryImpl{db: db}
}

func (r *purchaseRepositoryImpl) WithTx(tx *gorm.DB) PurchaseRepository {
	return &purchaseRepositoryImpl{db: tx}
}

func (r *purchaseRepositoryImpl) Create(ctx context.Context, purchase *entity.Purchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

func (r *purchaseRepositoryImpl) FindPurchases(ctx context.Context, req model.GetPurchaseBatchesRequest) ([]entity.Purchase, map[string][]string, int64, error) {
	var purchases []entity.Purchase
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Purchase{}).
		Joins("JOIN suppliers ON suppliers.id = purchases.supplier_id").
		Joins("JOIN users ON users.id = purchases.created_by")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where(
			"purchases.invoice_number ILIKE ? OR suppliers.name ILIKE ? OR EXISTS ("+
				"SELECT 1 FROM purchase_batches pb JOIN products p ON p.id = pb.product_id "+
				"WHERE pb.purchase_id = purchases.id AND p.name ILIKE ?)",
			search, search, search,
		)
	}
	if req.SupplierID != "" {
		query = query.Where("purchases.supplier_id = ?", req.SupplierID)
	}
	if req.ProductID != "" {
		query = query.Where("EXISTS (SELECT 1 FROM purchase_batches pb WHERE pb.purchase_id = purchases.id AND pb.product_id = ?)", req.ProductID)
	}
	if req.StartDate != "" {
		query = query.Where("purchases.purchase_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("purchases.purchase_date <= ?", req.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Supplier").
		Preload("Creator").
		Order("purchases.purchase_date DESC, purchases.created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&purchases).Error; err != nil {
		return nil, nil, 0, err
	}

	ids := make([]string, 0, len(purchases))
	for _, p := range purchases {
		ids = append(ids, p.ID)
	}

	productNames := make(map[string][]string)
	if len(ids) > 0 {
		type nameRow struct {
			PurchaseID string
			Name       string
		}
		var rows []nameRow
		if err := r.db.WithContext(ctx).Raw(
			`SELECT pb.purchase_id AS purchase_id, p.name AS name
			 FROM purchase_batches pb
			 JOIN products p ON p.id = pb.product_id
			 WHERE pb.purchase_id IN ?
			 ORDER BY p.name`,
			ids,
		).Scan(&rows).Error; err != nil {
			return nil, nil, 0, err
		}
		for _, row := range rows {
			productNames[row.PurchaseID] = append(productNames[row.PurchaseID], row.Name)
		}
	}

	return purchases, productNames, total, nil
}

func (r *purchaseRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Purchase, error) {
	var purchase entity.Purchase
	if err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Creator").
		First(&purchase, "purchases.id = ?", id).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *purchaseRepositoryImpl) Update(ctx context.Context, purchase *entity.Purchase) error {
	return r.db.WithContext(ctx).Save(purchase).Error
}

func (r *purchaseRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Purchase{}, "id = ?", id).Error
}